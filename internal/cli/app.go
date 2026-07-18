package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/adapters"
	"github.com/JonatasFreireDev/spec-framework/internal/install"
	"github.com/JonatasFreireDev/spec-framework/internal/moveartifact"
	"github.com/JonatasFreireDev/spec-framework/internal/runtimeassets"
	"github.com/JonatasFreireDev/spec-framework/internal/sourceimport"
	"github.com/JonatasFreireDev/spec-framework/internal/validator"
	"github.com/JonatasFreireDev/spec-framework/internal/wizard"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

type App struct {
	version string
}
type multiFlag []string

func (m *multiFlag) String() string         { return strings.Join(*m, ",") }
func (m *multiFlag) Set(value string) error { *m = append(*m, value); return nil }

func absoluteProductRoot(value string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(value) {
		return value, nil
	}
	return filepath.Join(cwd, value), nil
}
func parseByteLimit(value string) (int64, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	multiplier := int64(1)
	for _, unit := range []struct {
		suffix string
		size   int64
	}{{"GB", 1 << 30}, {"MB", 1 << 20}, {"KB", 1 << 10}, {"B", 1}} {
		if strings.HasSuffix(value, unit.suffix) {
			value, multiplier = strings.TrimSpace(strings.TrimSuffix(value, unit.suffix)), unit.size
			break
		}
	}
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil || number < 1 {
		return 0, fmt.Errorf("invalid byte limit %q", value)
	}
	return number * multiplier, nil
}

func New(version string) App {
	if version == "" {
		version = "dev"
	}
	return App{version: version}
}

func (app App) Run(args []string, stdout, stderr io.Writer) int {
	root := app.NewCommand(stdout, stderr)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		if exit, ok := err.(commandExitError); ok {
			return exit.code
		}
		fmt.Fprintln(stderr, err)
		return 2
	}
	return 0
}

func runSkill(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 || args[0] != "path" {
		fmt.Fprintln(stderr, "usage: spec-framework skill path <skill-name>")
		return 2
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	_, _, frameworkRoot, _, err := runtimeassets.Resolve(cwd)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	path := filepath.Join(frameworkRoot, "skills", args[1], "SKILL.md")
	if _, err := os.Stat(path); err != nil {
		fmt.Fprintf(stderr, "unknown framework skill %q\n", args[1])
		return 1
	}
	fmt.Fprintln(stdout, path)
	return 0
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
	graphView := flags.Bool("graph", false, "show consolidated workflow")
	asJSON := flags.Bool("json", false, "JSON output")
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
	if *graphView {
		return writeDashboard(p, *id, *asJSON, stdout, stderr)
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
		for _, b := range preview.ValidationBlockers {
			fmt.Fprintf(stdout, "BLOCKED: %s\n", b)
		}
		fmt.Fprintln(stdout, "Re-run with --yes to apply.")
		if len(preview.ParentBlockers) > 0 || len(preview.ValidationBlockers) > 0 {
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
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: spec-framework import create|status|resume|materialize")
		return 2
	}
	if args[0] == "create" {
		flags := flag.NewFlagSet("import create", flag.ContinueOnError)
		flags.SetOutput(stderr)
		productRoot := flags.String("product-root", "product", "product root")
		include := multiFlag{}
		exclude := multiFlag{}
		flags.Var(&include, "include", "include glob (repeatable)")
		flags.Var(&exclude, "exclude", "exclude glob (repeatable)")
		maxFiles := flags.Int("max-files", 500, "maximum matched files")
		maxTotal := flags.String("max-total-bytes", "200MB", "maximum copied bytes")
		maxFile := flags.String("max-file-bytes", "10MB", "maximum bytes per file")
		chunkSize := flags.Int("chunk-size", 25, "sources per analysis chunk")
		binaryPolicy := flags.String("binary-policy", "inventory_only", "inventory_only or reject")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if len(flags.Args()) == 0 {
			fmt.Fprintln(stderr, "import create requires one or more source paths")
			return 2
		}
		total, err := parseByteLimit(*maxTotal)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		file, err := parseByteLimit(*maxFile)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		root, err := absoluteProductRoot(*productRoot)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		run, err := sourceimport.CreateScalableRun(root, flags.Args(), sourceimport.CreateOptions{Include: include, Exclude: exclude, MaxFiles: *maxFiles, MaxTotalBytes: total, MaxFileBytes: file, ChunkSize: *chunkSize, BinaryPolicy: *binaryPolicy})
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Created scalable import run", run)
		return 0
	}
	if args[0] == "status" || args[0] == "resume" || args[0] == "record-review" {
		flags := flag.NewFlagSet("import "+args[0], flag.ContinueOnError)
		flags.SetOutput(stderr)
		productRoot := flags.String("product-root", "product", "product root")
		runID := flags.String("run", "", "import run id")
		agent := flags.String("agent", "", "importer identity")
		chunk := flags.String("chunk", "", "chunk id")
		input := flags.String("input", "", "review JSON input")
		yes := flags.Bool("yes", false, "confirm review record")
		jsonOutput := flags.Bool("json", false, "structured output")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if strings.TrimSpace(*runID) == "" {
			fmt.Fprintln(stderr, "import "+args[0]+" requires --run")
			return 2
		}
		root, err := absoluteProductRoot(*productRoot)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if args[0] == "resume" {
			claimed, err := sourceimport.Resume(root, *runID, *chunk, *agent)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			if *jsonOutput {
				data, _ := json.Marshal(claimed)
				fmt.Fprintln(stdout, string(data))
			} else {
				fmt.Fprintln(stdout, "RESUMED", claimed.ID, "for", claimed.Agent)
			}
			return 0
		}
		if args[0] == "record-review" {
			if !*yes || *chunk == "" || *agent == "" || *input == "" {
				fmt.Fprintln(stderr, "import record-review requires --chunk, --agent, --input, and --yes")
				return 2
			}
			data, err := os.ReadFile(*input)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			var review sourceimport.ChunkReview
			if err := json.Unmarshal(data, &review); err != nil {
				fmt.Fprintln(stderr, err)
				return 2
			}
			if err := sourceimport.RecordChunkReview(root, *runID, *chunk, *agent, review); err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			fmt.Fprintln(stdout, "REVIEWED", *chunk)
			return 0
		}
		status, err := sourceimport.ImportStatus(root, *runID)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if *jsonOutput {
			data, _ := json.Marshal(status)
			fmt.Fprintln(stdout, string(data))
		} else {
			fmt.Fprintf(stdout, "%s: %d source(s), %d chunk(s): %d queued, %d reviewing, %d reviewed, %d blocked, %d excluded\n", status.ImportID, status.Sources, status.Chunks, status.Queued, status.Reviewing, status.Reviewed, status.Blocked, status.Excluded)
		}
		return 0
	}
	if args[0] != "materialize" {
		fmt.Fprintln(stderr, "usage: spec-framework import create|status|resume|materialize")
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
	productRoot := flags.String("product-root", "", "product root; discovered from product/.product/framework.json by default")
	frameworkRoot := flags.String("framework-root", "", "framework root; resolved from the versioned user cache by default")
	writeReport := flags.Bool("write-report", false, "write validation and readiness reports")
	writeRegistry := flags.Bool("write-registry", false, "rebuild the artifact registry")
	strict := flags.Bool("strict", false, "promote approved-artifact delivery warnings to errors")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	productPath, frameworkPath := *productRoot, *frameworkRoot
	if productPath == "" && frameworkPath == "" {
		_, productPath, frameworkPath, _, err = runtimeassets.Resolve(cwd)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	} else {
		if productPath == "" {
			productPath = "product"
		}
		if frameworkPath == "" {
			resolveFrom := cwd
			if productPath != "" {
				resolveFrom = productPath
			}
			_, _, frameworkPath, _, err = runtimeassets.Resolve(resolveFrom)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
		}
		if !filepath.IsAbs(productPath) {
			productPath = filepath.Join(cwd, productPath)
		}
		if !filepath.IsAbs(frameworkPath) {
			frameworkPath = filepath.Join(cwd, frameworkPath)
		}
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
	var result validator.Result
	if *strict {
		result, err = validator.ValidateStrict(context.Background(), productPath, frameworkPath)
	} else {
		result, err = validator.Validate(context.Background(), productPath, frameworkPath)
	}
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

func runTemplate(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || (args[0] != "audit" && args[0] != "normalize") {
		fmt.Fprintln(stderr, "usage: spec-framework template audit|normalize --artifact <path> [--skill <owner-skill>] [--product-root product] [--framework-root <path>]")
		return 2
	}
	flags := flag.NewFlagSet("template audit", flag.ContinueOnError)
	flags.SetOutput(stderr)
	productRoot := flags.String("product-root", "product", "product root")
	frameworkRoot := flags.String("framework-root", "", "framework root; resolved from the versioned user cache by default")
	artifact := flags.String("artifact", "", "registered artifact path")
	skill := flags.String("skill", "", "owning skill used for provenance normalization")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if strings.TrimSpace(*artifact) == "" {
		fmt.Fprintln(stderr, "template audit requires --artifact")
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
	if frameworkPath == "" {
		_, _, frameworkPath, _, err = runtimeassets.Resolve(productPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	} else if !filepath.IsAbs(frameworkPath) {
		frameworkPath = filepath.Join(cwd, frameworkPath)
	}
	artifactPath := *artifact
	if !filepath.IsAbs(artifactPath) {
		artifactPath = filepath.Join(productPath, filepath.FromSlash(artifactPath))
	}
	if args[0] == "normalize" {
		findings, err := validator.AuditTemplate(context.Background(), productPath, frameworkPath, artifactPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if len(findings) > 0 {
			fmt.Fprintln(stdout, "Template normalization blocked; resolve audit findings first:")
			for _, finding := range findings {
				fmt.Fprintf(stdout, "%s %s %s: %s\n", strings.ToUpper(string(finding.Severity)), finding.Check, finding.File, finding.Message)
			}
			return 1
		}
		data, err := os.ReadFile(artifactPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		updated, err := sourceimport.NormalizeProvenance(string(data), *skill)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		tmp := artifactPath + ".normalize.tmp"
		if err := os.WriteFile(tmp, []byte(updated), 0644); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := os.Rename(tmp, artifactPath); err != nil {
			_ = os.Remove(tmp)
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := sourceimport.RecordNormalization(productPath, artifactPath, updated); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "Normalized provenance: %s\n- Skill: %s\n", filepath.ToSlash(*artifact), *skill)
		return 0
	}
	findings, err := validator.AuditTemplate(context.Background(), productPath, frameworkPath, artifactPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if len(findings) == 0 {
		fmt.Fprintf(stdout, "Template audit passed: %s\n", filepath.ToSlash(*artifact))
		return 0
	}
	fmt.Fprintf(stdout, "Template audit blocked: %s\n", filepath.ToSlash(*artifact))
	for _, finding := range findings {
		fmt.Fprintf(stdout, "%s %s %s: %s\n", strings.ToUpper(string(finding.Severity)), finding.Check, finding.File, finding.Message)
	}
	return 1
}

func (app App) runInit(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("target", "", "target directory")
	agentsValue := flags.String("agents", "", "comma-separated agents")
	startingPoint := flags.String("starting-point", "new-product", "new-product, existing-product, existing-documents, existing-feature, existing-implementation, or audit-only")
	sourcesValue := flags.String("sources", "", "comma-separated source files or directories for existing-documents")
	codeRootsValue := flags.String("code-roots", "", "comma-separated implementation roots as path:role (for example web:web,api:api)")
	sourceDir := flags.String("source-dir", "", "source directory for existing-documents")
	importMaxFiles := flags.Int("import-max-files", 500, "maximum files for existing-documents import")
	importMaxTotal := flags.String("import-max-total-bytes", "200MB", "maximum copied bytes for existing-documents import")
	importMaxFile := flags.String("import-max-file-bytes", "10MB", "maximum bytes per imported file")
	importChunkSize := flags.Int("import-chunk-size", 25, "sources per import review chunk")
	importBinaryPolicy := flags.String("import-binary-policy", "inventory_only", "inventory_only or reject")
	force := flags.Bool("force", false, "compatibility flag; never overwrites an existing product directory")
	installImpeccable := flags.Bool("install-impeccable", false, "install the optional Impeccable adapter after init")
	impeccableVersion := flags.String("impeccable-version", "", "explicit Impeccable CLI version")
	yes := flags.Bool("yes", false, "run headlessly")
	positionalTarget := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		positionalTarget, args = args[0], args[1:]
	}
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *target == "" {
		*target = positionalTarget
	}
	if *target == "" && flags.NArg() == 1 {
		*target = flags.Arg(0)
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
		*installImpeccable = result.InstallImpeccable
		*impeccableVersion = result.ImpeccableVersion
	}
	if *target == "" {
		fmt.Fprintln(stderr, "init requires --target")
		return 2
	}
	if *installImpeccable && strings.TrimSpace(*impeccableVersion) == "" {
		fmt.Fprintln(stderr, "--install-impeccable requires --impeccable-version")
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
	codeRoots, err := parseCodeRoots(*codeRootsValue)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if strings.TrimSpace(*sourceDir) != "" {
		sources = append(sources, *sourceDir)
	}
	importTotal, err := parseByteLimit(*importMaxTotal)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	importFile, err := parseByteLimit(*importMaxFile)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	result, err := install.Init(install.Options{Target: *target, Version: app.version, Agents: agents, StartingPoint: point, Sources: sources, CodeRoots: codeRoots, ImportOptions: sourceimport.CreateOptions{MaxFiles: *importMaxFiles, MaxTotalBytes: importTotal, MaxFileBytes: importFile, ChunkSize: *importChunkSize, BinaryPolicy: *importBinaryPolicy}, Force: *force})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Initialized Spec Framework product at %s\n- Product root: product\n- Framework runtime: %s\n- Repository-local agent trees: none\n- Starting point: %s\n", result.Target, result.SpecRoot, result.StartingPoint)
	if result.ImportID != "" {
		fmt.Fprintf(stdout, "- Import inventory: product/knowledge/imports/runs/%s\n", result.ImportID)
	}
	if len(codeRoots) > 0 {
		fmt.Fprintln(stdout, "- Declared code roots: product/knowledge/assessments/product-landscape.md")
	}
	if *installImpeccable {
		fmt.Fprintln(stdout, "[1/3] Resolving Impeccable version...")
		resolved, resolveErr := adapters.ResolveVersion("impeccable", *impeccableVersion)
		if resolveErr != nil {
			fmt.Fprintln(stderr, "Product initialized, but Impeccable version resolution failed:", resolveErr)
			return 1
		}
		argv, _ := adapters.ProviderArgv("impeccable", "install", resolved)
		fmt.Fprintf(stdout, "[1/3] Resolved Impeccable %s\n", resolved)
		fmt.Fprintf(stdout, "[2/3] Installing optional adapter\n- Provider: pbakaus/impeccable\n- Working directory: %s\n- Command: npx %s\n", result.Target, strings.Join(argv, " "))
		if err := adapters.Execute(result.Target, "impeccable", "install", resolved, stdout, stderr); err != nil {
			fmt.Fprintln(stderr, "Product initialized, but optional Impeccable installation failed:", err)
			return 1
		}
		fmt.Fprintln(stdout, "[3/3] Impeccable installed")
		fmt.Fprintf(stdout, "Next: reload the selected agent harness, open the repository root (%s), then run /impeccable init; keep generated product assets under product/design/. Do not use product/ as the working directory for adapter installation.\n", result.Target)
	}
	return 0
}

func (app App) runMigrate(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "external-runtime" {
		fmt.Fprintln(stderr, "usage: spec-framework migrate external-runtime [--target <path>] [--dry-run|--yes]")
		return 2
	}
	flags := flag.NewFlagSet("migrate external-runtime", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("target", ".", "legacy repository root")
	dryRun := flags.Bool("dry-run", false, "preview without writing")
	yes := flags.Bool("yes", false, "apply migration")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	legacy := filepath.Join(*target, ".spec-framework", "manifest.json")
	data, err := os.ReadFile(legacy)
	if err != nil {
		fmt.Fprintln(stderr, "legacy manifest not found:", legacy)
		return 1
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	version, _ := value["version"].(string)
	if version == "" {
		version = app.version
	}
	fmt.Fprintf(stdout, "Migration preview\n- Runtime version: %s\n- Create: product/.product/framework.json\n- Preserve for manual review: .spec-framework and local agent trees\n", version)
	if *dryRun || !*yes {
		fmt.Fprintln(stdout, "No files changed. Re-run with --yes to apply.")
		return 0
	}
	agents := []install.Agent{install.Codex}
	if raw, ok := value["agents"].([]any); ok {
		var names []string
		for _, item := range raw {
			if name, ok := item.(string); ok {
				names = append(names, name)
			}
		}
		if parsed, err := install.ParseAgents(strings.Join(names, ",")); err == nil {
			agents = parsed
		}
	}
	result, err := install.Upgrade(install.Options{Target: *target, Version: version, Agents: agents})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Migrated manifest; runtime available at %s\n", result.SpecRoot)
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

func parseCodeRoots(value string) ([]runtimeassets.CodeRoot, error) {
	var roots []runtimeassets.CodeRoot
	for _, item := range splitCSV(value) {
		parts := strings.SplitN(item, ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("invalid --code-roots entry %q; use path:role", item)
		}
		roots = append(roots, runtimeassets.CodeRoot{Path: strings.TrimSpace(parts[0]), Role: strings.TrimSpace(parts[1])})
	}
	return roots, nil
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
	fmt.Fprintf(stdout, "Upgraded Spec Framework runtime at %s\n- Product root preserved: product\n- Framework runtime: %s\n- Version: %s\n", result.Target, result.SpecRoot, app.version)
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
