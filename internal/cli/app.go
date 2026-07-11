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
	"github.com/JonatasFreireDev/spec-framework/internal/validator"
	"github.com/JonatasFreireDev/spec-framework/internal/wizard"
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
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		writeHelp(stderr)
		return 2
	}
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
	productPath := filepath.Join(cwd, *productRoot)
	frameworkPath := filepath.Join(cwd, *frameworkRoot)
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
	result, err := install.Init(install.Options{Target: *target, Version: app.version, Agents: agents, Force: *force})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Initialized Spec Framework product at %s\n- Product root: product\n- Framework assets: .spec-framework\n- Agent integrations: %s\n", result.Target, *agentsValue)
	return 0
}

func (app App) runUpgrade(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("upgrade", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("target", ".", "target directory")
	agentsValue := flags.String("agents", "codex", "comma-separated agents")
	yes := flags.Bool("yes", false, "confirm upgrade")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if !*yes {
		fmt.Fprintln(stderr, "upgrade requires --yes in headless mode")
		return 2
	}
	agents, err := install.ParseAgents(*agentsValue)
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
  validate   Validate a product repository.
  move       Move an artifact and update references.
  upgrade    Refresh installed framework assets.
  version    Print the CLI version.
  help       Show this help.
`)
}
