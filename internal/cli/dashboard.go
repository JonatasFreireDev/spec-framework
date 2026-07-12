package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func runDashboard(args []string, out, errout io.Writer) int {
	fs := flag.NewFlagSet("dashboard", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	work := fs.String("work", "", "workspace id")
	asJSON := fs.Bool("json", false, "JSON output")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	if *work == "" {
		fmt.Fprintln(errout, "dashboard requires --work WORK-NNN")
		return 2
	}
	return writeDashboard(productPath(*root), *work, *asJSON, out, errout)
}
func writeDashboard(root, work string, asJSON bool, out, errout io.Writer) int {
	d, e := workflow.BuildDashboard(root, work)
	if e != nil {
		fmt.Fprintln(errout, e)
		return 1
	}
	if asJSON {
		b, _ := json.MarshalIndent(d, "", "  ")
		fmt.Fprintln(out, string(b))
		return 0
	}
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	done := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	current := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
	blocked := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	fmt.Fprintln(out, title.Render("Spec Framework · "+d.WorkspaceID))
	fmt.Fprintf(out, "Feature: %s\nUse case: %s\n\n", d.Feature, valueOr(d.UseCase, "not selected"))
	var flow []string
	for _, s := range d.Stages {
		label := "○ " + s.Label
		switch s.Status {
		case "done":
			label = done.Render("✓ " + s.Label)
		case "current":
			label = current.Render("● " + s.Label)
		case "blocked":
			label = blocked.Render("× " + s.Label)
		default:
			label = muted.Render(label)
		}
		flow = append(flow, label)
	}
	fmt.Fprintln(out, strings.Join(flow, "  →  "))
	fmt.Fprintf(out, "\nCurrent: %s · Skill: %s · Expected: %s\n", d.CurrentStep, d.RecommendedSkill, d.ExpectedArtifact)
	fmt.Fprintf(out, "Graph: %s · Tasks: %d total / %d ready\n", valueOr(d.GraphStatus, "not created"), d.TaskTotal, d.TaskReady)
	if d.LatestCheckpoint != "" {
		fmt.Fprintln(out, "Checkpoint:", d.LatestCheckpoint)
	}
	if d.LatestHandoff != "" {
		fmt.Fprintln(out, "Handoff:", d.LatestHandoff)
	}
	for _, x := range d.Decisions {
		fmt.Fprintln(out, "DECISION", x)
	}
	for _, x := range d.ActiveLeases {
		fmt.Fprintln(out, "LEASE", x)
	}
	for _, x := range d.Blockers {
		fmt.Fprintln(out, blocked.Render("BLOCKED "+x))
	}
	for _, x := range d.RequiredReading {
		fmt.Fprintln(out, "READ", x)
	}
	for _, x := range d.NextCommands {
		fmt.Fprintln(out, current.Render("NEXT "+x))
	}
	if len(d.Blockers) > 0 {
		return 1
	}
	return 0
}
func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
