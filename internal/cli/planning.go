package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/JonatasFreireDev/spec-framework/internal/designsystem"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func productPath(value string) string {
	if filepath.IsAbs(value) {
		return value
	}
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, value)
}
func runTask(args []string, out, errout io.Writer) int {
	if len(args) == 0 || args[0] != "readiness" {
		fmt.Fprintln(errout, "task requires readiness")
		return 2
	}
	fs := flag.NewFlagSet("task readiness", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	graph := fs.String("graph", "", "graph path")
	task := fs.String("task", "", "task id")
	asJSON := fs.Bool("json", false, "JSON output")
	if e := fs.Parse(args[1:]); e != nil {
		return 2
	}
	p := productPath(*root)
	g := *graph
	if !filepath.IsAbs(g) {
		g = filepath.Join(p, filepath.FromSlash(g))
	}
	r, e := workflow.CheckTaskReadiness(p, g, *task)
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	if *asJSON {
		b, _ := json.MarshalIndent(r, "", "  ")
		fmt.Fprintln(out, string(b))
	} else {
		for _, c := range r.Checks {
			fmt.Fprintf(out, "%s %s: %s\n", c.Status, c.ID, c.Detail)
		}
	}
	if !r.Ready {
		return 1
	}
	return 0
}
func runGuide(args []string, out, errout io.Writer) int {
	fs := flag.NewFlagSet("guide", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	work := fs.String("work", "", "workspace id")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	g, e := workflow.WorkspaceGuide(productPath(*root), *work)
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	fmt.Fprintf(out, "Workspace: %s\nFeature scope: %s\n", g.WorkspaceID, g.FeatureScope)
	if g.UseCaseScope != "" {
		fmt.Fprintf(out, "Use case scope: %s\n", g.UseCaseScope)
	}
	fmt.Fprintf(out, "Current step: %s\nSkill: %s\nExpected artifact: %s\n", g.CurrentStep, g.RecommendedSkill, g.ExpectedArtifact)
	for _, x := range g.RequiredReading {
		fmt.Fprintln(out, "READ:", x)
	}
	for _, x := range g.Blockers {
		fmt.Fprintln(out, "BLOCKED:", x)
	}
	for _, x := range g.Commands {
		fmt.Fprintln(out, "NEXT:", x)
	}
	if len(g.Blockers) > 0 {
		return 1
	}
	return 0
}
func runStage(command string, args []string, out, errout io.Writer) int {
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	work := fs.String("work", "", "workspace id")
	stage := fs.String("stage", "", "stage")
	by := fs.String("approved-by", "", "approver")
	notes := fs.String("notes", "", "notes")
	yes := fs.Bool("yes", false, "confirm")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	p := productPath(*root)
	r, e := workflow.ReviewStage(p, *work, *stage)
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	fmt.Fprintf(out, "Stage %s: %d artifacts\n", r.Stage, len(r.Artifacts))
	for _, a := range r.Artifacts {
		fmt.Fprintf(out, "- %s %s [%s]\n", a.ID, a.Path, a.Status)
	}
	for _, b := range r.Blockers {
		fmt.Fprintln(out, "BLOCKED:", b)
	}
	if command == "review" || !*yes {
		if command == "approve-stage" {
			fmt.Fprintln(out, "Re-run with --yes to approve atomically.")
		}
		if len(r.Blockers) > 0 {
			return 1
		}
		return 0
	}
	if *by == "" {
		fmt.Fprintln(errout, "approve-stage requires --approved-by")
		return 2
	}
	records, e := workflow.ApproveStage(p, *work, *stage, *by, *notes)
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	fmt.Fprintf(out, "APPROVED %d artifacts\n", len(records))
	return 0
}

func runImpact(args []string, out, errout io.Writer) int {
	fs := flag.NewFlagSet("impact", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	decision := fs.String("decision", "", "decision id")
	designSystem := fs.String("design-system", "", "Design System id")
	asJSON := fs.Bool("json", false, "JSON output")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	if *decision == "" && *designSystem == "" {
		fmt.Fprintln(errout, "impact requires --decision DEC-NNN or --design-system DSYS-NNN")
		return 2
	}
	if *designSystem != "" {
		r, e := designsystem.ImpactReport(productPath(*root), *designSystem)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		if *asJSON {
			b, _ := json.MarshalIndent(r, "", "  ")
			fmt.Fprintln(out, string(b))
		} else {
			fmt.Fprintf(out, "Design System: %s@%s\n", r.ID, r.Version)
			for _, item := range r.Consumers {
				fmt.Fprintln(out, "CONSUMER", item)
			}
			for _, blocker := range r.Blockers {
				fmt.Fprintln(out, "BLOCKED", blocker)
			}
		}
		if len(r.Blockers) > 0 {
			return 1
		}
		return 0
	}
	r, e := workflow.DecisionImpactReport(productPath(*root), *decision)
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	if *asJSON {
		b, _ := json.MarshalIndent(r, "", "  ")
		fmt.Fprintln(out, string(b))
	} else {
		fmt.Fprintf(out, "Decision: %s\nType: %s\nStatus: %s\nValid: %t\n", r.ID, r.Type, r.Status, r.Valid)
		for _, x := range r.AffectedArtifacts {
			fmt.Fprintln(out, "AFFECTS", x)
		}
		for _, x := range r.References {
			fmt.Fprintln(out, "REFERENCED", x)
		}
		for _, x := range r.StaleArtifacts {
			fmt.Fprintln(out, "STALE", x)
		}
		for _, x := range r.PropagationGaps {
			fmt.Fprintln(out, "MISSING", x)
		}
		for _, x := range r.Blockers {
			fmt.Fprintln(out, "BLOCKED", x)
		}
	}
	if !r.Valid {
		return 1
	}
	return 0
}
