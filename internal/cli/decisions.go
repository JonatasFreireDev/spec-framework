package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/decisioncheck"
	"github.com/JonatasFreireDev/spec-framework/internal/validator"
	"github.com/JonatasFreireDev/spec-framework/internal/wizard"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func runDecisions(args []string, out, errout io.Writer) int {
	if len(args) == 0 || (args[0] != "migrate" && args[0] != "check") {
		fmt.Fprintln(errout, "decisions requires check or migrate")
		return 2
	}
	if args[0] == "check" {
		return runDecisionCheck(args[1:], out, errout)
	}
	fs := flag.NewFlagSet("decisions migrate", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	interactive := fs.Bool("interactive", false, "review choices interactively")
	yes := fs.Bool("yes", false, "apply inferred migration")
	accept := fs.Bool("accept-inferred", false, "accept ambiguous inferred values")
	asJSON := fs.Bool("json", false, "JSON preview")
	if e := fs.Parse(args[1:]); e != nil {
		return 2
	}
	plan, e := workflow.PlanDecisionMigration(productPath(*root))
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	if *asJSON {
		b, _ := json.MarshalIndent(plan.Items, "", "  ")
		fmt.Fprintln(out, string(b))
	} else {
		fmt.Fprintf(out, "Decision migration preview: %d entries\n", len(plan.Items))
		for _, x := range plan.Items {
			review := ""
			if x.NeedsReview {
				review = " REVIEW"
			}
			fmt.Fprintf(out, "- %s type=%s scope=%v%s (%s)\n", x.ID, x.InferredType, x.Scope, review, strings.Join(x.Reasons, ", "))
		}
	}
	if len(plan.Items) == 0 {
		fmt.Fprintln(out, "No legacy decisions require migration.")
		return 0
	}
	items := plan.Items
	if *interactive {
		var confirmed bool
		items, confirmed, e = wizard.RunDecisionMigration(os.Stdin, out, plan)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		if !confirmed {
			return 0
		}
	} else if !*yes {
		fmt.Fprintln(out, "Preview only. Re-run with --interactive or --yes.")
		return 0
	} else if !*accept {
		for _, x := range items {
			if x.NeedsReview {
				fmt.Fprintln(errout, "ambiguous inference requires --interactive or --accept-inferred:", x.ID)
				return 1
			}
		}
	}
	result, e := workflow.ApplyDecisionMigration(plan, items)
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	fmt.Fprintf(out, "Migrated %d decisions\nBackup: %s\nIndex: %s\n", result.Changed, result.Backup, result.IndexPath)
	return 0
}

func runDecisionCheck(args []string, out, errout io.Writer) int {
	fs := flag.NewFlagSet("decisions check", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	frameworkRoot := fs.String("framework-root", ".", "framework root")
	domain := fs.String("domain", "", "only inspect one decision domain")
	strict := fs.Bool("strict", false, "fail when errors are found; warnings remain non-blocking")
	asJSON := fs.Bool("json", false, "JSON output")
	fixLinks := fs.Bool("fix-links", false, "preview or fix mechanically resolvable decision links")
	yes := fs.Bool("yes", false, "apply link fixes")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	rootPath := productPath(*root)
	frameworkPath := *frameworkRoot
	if !filepath.IsAbs(frameworkPath) {
		cwd, _ := os.Getwd()
		frameworkPath = filepath.Join(cwd, frameworkPath)
	}
	report, err := decisioncheck.Run(decisioncheck.Options{Root: rootPath, FrameworkRoot: frameworkPath, Domain: *domain, FixLinks: *fixLinks, Yes: *yes})
	if err != nil {
		fmt.Fprintln(errout, err)
		return 1
	}
	if *fixLinks && !*yes {
		fmt.Fprintln(out, "Preview only: re-run with --fix-links --yes to apply mechanically resolvable links.")
	}
	if *asJSON {
		b, _ := json.MarshalIndent(report, "", "  ")
		fmt.Fprintln(out, string(b))
	} else {
		for _, d := range report.Diagnostics {
			fmt.Fprintf(out, "%s %s %s: %s\n", strings.ToUpper(string(d.Severity)), d.Check, d.File, d.Message)
			if d.Fix != "" {
				fmt.Fprintf(out, "  Fix: %s\n", d.Fix)
			}
		}
		for _, path := range report.ChangedFiles {
			fmt.Fprintf(out, "CHANGED %s\n", path)
		}
		fmt.Fprintf(out, "Decision check: %d decisions, %d changed files\n", report.DecisionCount, len(report.ChangedFiles))
	}
	if *strict {
		for _, d := range report.Diagnostics {
			if d.Severity == validator.Error {
				return 1
			}
		}
	}
	return 0
}
