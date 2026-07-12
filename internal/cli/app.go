package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/install"
	"github.com/JonatasFreireDev/spec-framework/internal/moveartifact"
	"github.com/JonatasFreireDev/spec-framework/internal/sourceimport"
	"github.com/JonatasFreireDev/spec-framework/internal/validator"
	"github.com/JonatasFreireDev/spec-framework/internal/wizard"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

type App struct {
	version string
}

func New(version string) App {
	if version == "" {
		version = "dev"
	}
	return App{version: version}
}

func (app App) Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		writeHelp(stdout)
		return 0
	}

	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "spec-framework %s\n", app.version)
		return 0
	case "move":
		return runMove(args[1:], stdout, stderr)
	case "init":
		return app.runInit(args[1:], stdout, stderr)
	case "upgrade":
		return app.runUpgrade(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "import":
		return runImport(args[1:], stdout, stderr)
	case "work":
		return runWork(args[1:], stdout, stderr)
	case "status", "next":
		return runWorkStatus(args[0], args[1:], stdout, stderr)
	case "approve":
		return runApprove(args[1:], stdout, stderr)
	case "gates":
		return runGates(args[1:], stdout, stderr)
	case "graph":
		return runGraph(args[1:], stdout, stderr)
	case "task":
		return runTask(args[1:], stdout, stderr)
	case "guide":
		return runGuide(args[1:], stdout, stderr)
	case "review", "approve-stage":
		return runStage(args[0], args[1:], stdout, stderr)
	case "resume", "handoff", "checkpoint", "lease", "commands", "schedule", "integrate", "runtime":
		return runRuntime(args[0], args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		writeHelp(stderr)
		return 2
	}
}

func runWork(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("work", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	feature := flags.String("feature", "", "feature path or id")
	domain := flags.String("domain", "", "domain scope")
	goal := flags.String("goal", "", "goal scope")
	useCase := flags.String("use-case", "", "optional use-case slug or product-relative path")
	createdBy := flags.String("created-by", "human", "workspace creator")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	if strings.TrimSpace(*feature) == "" {
		items, err := workflow.Features(p)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Available features:")
		for _, a := range items {
			fmt.Fprintf(stdout, "- %s  %s  [%s]\n", a.ID, a.Path, a.Status)
		}
		fmt.Fprintln(stdout, "Select with: spec-framework work --feature <id-or-path>")
		return 0
	}
	w, err := workflow.CreateWorkspace(p, *feature, *domain, *goal, *useCase, *createdBy)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Created %s\n- Feature: %s\n- Next skill: %s\n", w.ID, w.Scope["feature"], w.RecommendedSkill)
	return 0
}
func runWorkStatus(command string, args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	id := flags.String("work", "", "workspace id")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *id == "" {
		fmt.Fprintln(stderr, command+" requires --work WORK-NNN")
		return 2
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	s, err := workflow.WorkspaceStatus(p, *id)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Workspace: %s\nArtifact: %s (%s)\nStatus: %s\nNext skill: %s\n", s.Workspace.ID, s.Artifact.ID, s.Artifact.Path, s.Artifact.Status, s.Next)
	for _, b := range s.Blockers {
		fmt.Fprintf(stdout, "BLOCKED: %s\n", b)
	}
	if len(s.Blockers) > 0 {
		return 1
	}
	return 0
}
func runApprove(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("approve", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	artifact := flags.String("artifact", "", "artifact path")
	grant := flags.String("grant", "approved", "status to grant")
	by := flags.String("approved-by", "", "approving human")
	notes := flags.String("notes", "", "approval notes")
	yes := flags.Bool("yes", false, "confirm approval")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *artifact == "" || *by == "" {
		fmt.Fprintln(stderr, "approve requires --artifact and --approved-by")
		return 2
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	a := *artifact
	if !filepath.IsAbs(a) {
		a = filepath.Join(p, filepath.FromSlash(a))
	}
	if !*yes {
		preview, err := workflow.PreviewApproval(p, a, *grant)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "Approval preview\n- Artifact: %s (%s)\n- Current status: %s\n- Grant: %s\n- Result hash: %s\n", preview.Artifact.ID, preview.Artifact.Path, preview.Artifact.Status, preview.Grant, preview.CurrentHash)
		for _, b := range preview.ParentBlockers {
			fmt.Fprintf(stdout, "BLOCKED: %s\n", b)
		}
		fmt.Fprintln(stdout, "Re-run with --yes to apply.")
		if len(preview.ParentBlockers) > 0 {
			return 1
		}
		return 0
	}
	rec, err := workflow.Approve(p, a, *grant, *by, *notes)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Approved %s as %s\n- Path: %s\n- Hash: %s\n", rec.ArtifactID, rec.StatusGranted, rec.Path, rec.ContentHash)
	return 0
}
func runGates(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("gates", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	missing, err := workflow.GateReadiness(p)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if len(missing) == 0 {
		fmt.Fprintln(stdout, "Gates: ready")
		return 0
	}
	for _, id := range missing {
		fmt.Fprintf(stdout, "MISSING %s: TBD blocks implementation\n", id)
	}
	return 1
}
func runGraph(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "graph requires ready, materialize, claim, release, or complete")
		return 2
	}
	command := args[0]
	flags := flag.NewFlagSet("graph "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	graph := flags.String("graph", "", "execution graph path")
	task := flags.String("task", "", "task id")
	agent := flags.String("agent", "", "agent id")
	yes := flags.Bool("yes", false, "confirm mutation")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	g := *graph
	if g != "" && !filepath.IsAbs(g) {
		g = filepath.Join(p, filepath.FromSlash(g))
	}
	if g != "" {
		rel, err := filepath.Rel(p, g)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			fmt.Fprintln(stderr, "graph path escapes product root")
			return 2
		}
	}
	switch command {
	case "materialize":
		if g == "" {
			fmt.Fprintln(stderr, "graph materialize requires --graph")
			return 2
		}
		if !*yes {
			fmt.Fprintln(stdout, "Preview: materialize missing task files and tasks.md; re-run with --yes")
			return 0
		}
		result, err := workflow.MaterializeTasks(g)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if _, err := validator.WriteRegistry(p); err != nil {
			fmt.Fprintln(stderr, "materialized but registry update failed:", err)
			return 1
		}
		fmt.Fprintf(stdout, "MATERIALIZED %d tasks\n", len(result.Tasks))
		for _, p := range result.Tasks {
			fmt.Fprintln(stdout, "-", p)
		}
		return 0
	case "ready":
		nodes, err := workflow.ReadyUnclaimed(p, g)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		for _, n := range nodes {
			fmt.Fprintf(stdout, "READY %s %s\n", n.ID, n.Path)
		}
		return 0
	case "claim":
		if *task == "" || *agent == "" || g == "" {
			fmt.Fprintln(stderr, "graph claim requires --graph, --task, and --agent")
			return 2
		}
		c, err := workflow.ClaimTask(p, g, *task, *agent)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "CLAIMED %s by %s\n", c.TaskID, c.Agent)
		return 0
	case "release":
		if *task == "" || *agent == "" {
			fmt.Fprintln(stderr, "graph release requires --task and --agent")
			return 2
		}
		if err := workflow.ReleaseClaim(p, *task, *agent); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "RELEASED %s\n", *task)
		return 0
	case "complete":
		if *task == "" || *agent == "" || g == "" {
			fmt.Fprintln(stderr, "graph complete requires --graph, --task, and --agent")
			return 2
		}
		if err := workflow.Complete(p, g, *task, *agent); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "COMPLETED %s\n", *task)
		return 0
	default:
		fmt.Fprintln(stderr, "unknown graph command", command)
		return 2
	}
}

func runImport(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "materialize" {
		fmt.Fprintln(stderr, "usage: spec-framework import materialize --run IMPORT-NNN --approved-by <name> --yes")
		return 2
	}
	flags := flag.NewFlagSet("import materialize", flag.ContinueOnError)
	flags.SetOutput(stderr)
	productRoot := flags.String("product-root", "product", "product root")
	runID := flags.String("run", "", "import run id")
	approvedBy := flags.String("approved-by", "", "human approving the selected mappings")
	yes := flags.Bool("yes", false, "confirm materialization")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if !*yes || strings.TrimSpace(*runID) == "" || strings.TrimSpace(*approvedBy) == "" {
		fmt.Fprintln(stderr, "materialization requires --run, --approved-by, and --yes")
		return 2
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	root := *productRoot
	if !filepath.IsAbs(root) {
		root = filepath.Join(cwd, root)
	}
	created, err := sourceimport.Materialize(root, *runID, *approvedBy)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Materialized %s as draft (%d files)\n", *runID, len(created))
	for _, path := range created {
		fmt.Fprintf(stdout, "- product/%s\n", path)
	}
	return 0
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	productRoot := flags.String("product-root", "product", "product root")
	frameworkRoot := flags.String("framework-root", ".spec-framework", "framework root")
	writeReport := flags.Bool("write-report", false, "write validation and readiness reports")
	writeRegistry := flags.Bool("write-registry", false, "rebuild the artifact registry")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	productPath := *productRoot
	if !filepath.IsAbs(productPath) {
		productPath = filepath.Join(cwd, productPath)
	}
	frameworkPath := *frameworkRoot
	if !filepath.IsAbs(frameworkPath) {
		frameworkPath = filepath.Join(cwd, frameworkPath)
	}
	if *writeRegistry {
		path, err := validator.WriteRegistry(productPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		rel, _ := filepath.Rel(productPath, path)
		fmt.Fprintf(stdout, "Wrote %s\n", filepath.ToSlash(rel))
	}
	result, err := validator.Validate(context.Background(), productPath, frameworkPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *writeReport {
		paths, err := validator.WriteReport(productPath, result)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		for _, path := range paths {
			rel, _ := filepath.Rel(productPath, path)
			fmt.Fprintf(stdout, "Wrote %s\n", filepath.ToSlash(rel))
		}
	}
	for _, d := range result.Diagnostics {
		fmt.Fprintf(stdout, "%s %s %s: %s\n", strings.ToUpper(string(d.Severity)), d.Check, d.File, d.Message)
	}
	icon := "✅"
	if result.Errors > 0 {
		icon = "🔴"
	} else if result.Warnings > 0 {
		icon = "🟡"
	}
	fmt.Fprintf(stdout, "Verdict: %s %s (%d errors, %d warnings, %d notes)\n", icon, result.Verdict(), result.Errors, result.Warnings, result.Notes)
	if result.Errors > 0 {
		return 1
	}
	return 0
}

func (app App) runInit(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("target", "", "target directory")
	agentsValue := flags.String("agents", "", "comma-separated agents")
	startingPoint := flags.String("starting-point", "new-product", "new-product, existing-product, existing-documents, existing-feature, existing-implementation, or audit-only")
	sourcesValue := flags.String("sources", "", "comma-separated source files or directories for existing-documents")
	sourceDir := flags.String("source-dir", "", "source directory for existing-documents")
	force := flags.Bool("force", false, "allow non-empty target")
	yes := flags.Bool("yes", false, "run headlessly")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if !*yes {
		result, err := wizard.RunInit(os.Stdin, stdout)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if result.Cancelled || !result.Confirmed {
			return 0
		}
		*target = result.Target
		selected := result.AgentNames()
		*agentsValue = strings.Join(selected, ",")
		*startingPoint = result.StartingPoint
		*sourcesValue = strings.Join(result.Sources, ",")
	}
	if *target == "" {
		fmt.Fprintln(stderr, "init requires --target")
		return 2
	}
	agents, err := install.ParseAgents(*agentsValue)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	point, err := install.ParseStartingPoint(*startingPoint)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	sources := splitCSV(*sourcesValue)
	if strings.TrimSpace(*sourceDir) != "" {
		sources = append(sources, *sourceDir)
	}
	result, err := install.Init(install.Options{Target: *target, Version: app.version, Agents: agents, StartingPoint: point, Sources: sources, Force: *force})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Initialized Spec Framework product at %s\n- Product root: product\n- Framework assets: .spec-framework\n- Agent integrations: %s\n- Starting point: %s\n", result.Target, *agentsValue, result.StartingPoint)
	if result.ImportID != "" {
		fmt.Fprintf(stdout, "- Import inventory: product/knowledge/imports/runs/%s\n", result.ImportID)
	}
	return 0
}

func splitCSV(value string) []string {
	var out []string
	for _, part := range strings.Split(value, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func (app App) runUpgrade(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("upgrade", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("target", ".", "target directory")
	agentsValue := flags.String("agents", "", "comma-separated agents; defaults to the installed manifest")
	yes := flags.Bool("yes", false, "confirm upgrade")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if !*yes {
		fmt.Fprintln(stderr, "upgrade requires --yes in headless mode")
		return 2
	}
	var agents []install.Agent
	var err error
	if strings.TrimSpace(*agentsValue) == "" {
		agents, err = install.InstalledAgents(*target)
	} else {
		agents, err = install.ParseAgents(*agentsValue)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	result, err := install.Upgrade(install.Options{Target: *target, Version: app.version, Agents: agents})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Upgraded Spec Framework assets at %s\n- Product root preserved: product\n- Framework assets updated: .spec-framework\n- Version: %s\n", result.Target, app.version)
	return 0
}

func runMove(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("move", flag.ContinueOnError)
	flags.SetOutput(stderr)
	from := flags.String("from", "", "source path")
	to := flags.String("to", "", "target path")
	dryRun := flags.Bool("dry-run", false, "plan without writing")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *from == "" || *to == "" {
		fmt.Fprintln(stderr, "move requires --from and --to")
		return 2
	}
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	plan, err := moveartifact.Build(root, *from, *to)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if !*dryRun {
		if err := moveartifact.Apply(plan); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	verb := "Moved"
	if *dryRun {
		verb = "Dry run"
	}
	fmt.Fprintf(stdout, "%s: %s -> %s\n", verb, plan.OldRel, plan.NewRel)
	fmt.Fprintf(stdout, "Rewritten files: %d\n", len(plan.Rewrites))
	for _, item := range plan.Rewrites {
		rel, _ := filepath.Rel(plan.Root, item.Path)
		fmt.Fprintf(stdout, "- %s %s\n", filepath.ToSlash(rel), item.Kind)
	}
	fmt.Fprintf(stdout, "Free-text mentions requiring review: %d\n", len(plan.Mentions))
	for _, item := range plan.Mentions {
		fmt.Fprintf(stdout, "- %s\n", item)
	}
	return 0
}

func writeHelp(output io.Writer) {
	fmt.Fprint(output, `Usage: spec-framework <command> [options]

Commands:
  init       Initialize a product repository.
  import     Materialize approved source mappings as drafts.
  work       Select a feature and create a concurrent workspace.
  status     Show workspace readiness and blockers.
  next       Show the next skill for a workspace.
  approve    Review and record an explicit artifact approval.
  gates      Check whether implementation gates are configured.
  graph      Inspect and operate execution graph claims.
  task       Inspect task readiness.
  guide      Explain the current workspace gate and next action.
  review     Preview a workspace stage approval.
  approve-stage Approve every eligible artifact in a stage atomically.
  resume     Resume a persisted runtime workspace.
  handoff    Persist an agent/orchestrator handoff.
  checkpoint Persist a resumable checkpoint.
  lease      Claim, heartbeat, or recover task leases.
  commands   Plan or execute shell-free R0/R1 commands.
  schedule   Build deterministic parallel execution waves.
  integrate  Plan or apply local commit integration.
  runtime    Migrate a v1 workspace to runtime v2.
  validate   Validate a product repository.
  move       Move an artifact and update references.
  upgrade    Refresh installed framework assets.
  version    Print the CLI version.
  help       Show this help.
`)
}
