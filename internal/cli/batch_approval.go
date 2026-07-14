package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	huh "charm.land/huh/v2"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func runApproveBatch(args []string, out, errout io.Writer) int {
	fs := flag.NewFlagSet("approve-batch", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	artifact := fs.String("artifact", "", "comma-separated artifact paths or IDs")
	ids := fs.String("ids", "", "comma-separated artifact IDs")
	foundation := fs.Bool("foundation", false, "select Foundation artifacts")
	stage := fs.String("stage", "", "select one stage: foundation, domains, feature, use-cases, specification, design, engineering, planning, or tasks")
	until := fs.String("until", "", "with --all-eligible, select through this stage")
	allEligible := fs.Bool("all-eligible", false, "select all registered artifacts eligible through the optional --until stage")
	by := fs.String("approved-by", "", "approving human")
	notes := fs.String("notes", "", "approval notes")
	dryRun := fs.Bool("dry-run", false, "preview without mutating")
	yes := fs.Bool("yes", false, "confirm and apply the batch")
	interactive := fs.Bool("interactive", false, "use the Bubble Tea confirmation UI")
	asJSON := fs.Bool("json", false, "JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *yes && strings.TrimSpace(*by) == "" {
		fmt.Fprintln(errout, "approve-batch --yes requires --approved-by")
		return 2
	}
	if *yes && *interactive {
		fmt.Fprintln(errout, "approve-batch cannot combine --yes and --interactive")
		return 2
	}
	scope := workflow.BatchScope{Artifact: splitCSV(*artifact), IDs: splitCSV(*ids), Foundation: *foundation, Stage: *stage, Until: *until, AllEligible: *allEligible}
	plan, err := workflow.BuildBatchApprovalPlan(productPath(*root), scope, "approved")
	if err != nil {
		fmt.Fprintln(errout, err)
		return 2
	}
	if *asJSON && !*yes {
		data, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Fprintln(out, string(data))
	} else if !*yes {
		printBatchPlan(out, plan)
	}
	if len(plan.Blockers) > 0 {
		if *yes {
			fmt.Fprintln(errout, "approval batch blocked; resolve the listed blockers before applying")
		}
		return 1
	}
	if *interactive {
		if strings.TrimSpace(*by) == "" {
			fmt.Fprintln(errout, "approve-batch --interactive requires --approved-by")
			return 2
		}
		confirmed := false
		form := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(fmt.Sprintf("Approve %d artifact(s) as %s?", len(plan.ToApprove), *by)).Affirmative("Approve").Negative("Cancel").Value(&confirmed)))
		if err := form.Run(); err != nil {
			fmt.Fprintln(errout, err)
			return 1
		}
		if !confirmed {
			fmt.Fprintln(out, "Approval batch cancelled.")
			return 0
		}
		*yes = true
	}
	if !*yes || *dryRun {
		if !*asJSON {
			fmt.Fprintln(out, "Re-run with --yes --approved-by <name> to apply this batch.")
		}
		return 0
	}
	records, err := workflow.ApproveBatch(productPath(*root), plan, *by, *notes)
	if *asJSON {
		result := map[string]any{"plan": plan, "approved": records}
		if err != nil {
			result["error"] = err.Error()
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprintln(out, string(data))
		if err != nil {
			return 1
		}
		return 0
	}
	for _, record := range records {
		fmt.Fprintf(out, "APPROVED %s [%s] hash=%s\n", record.ArtifactID, record.Path, record.ContentHash)
	}
	if err != nil {
		fmt.Fprintln(errout, err)
		return 1
	}
	fmt.Fprintf(out, "APPROVED %d artifact(s); next gate: %s\n", len(records), plan.NextGate)
	return 0
}

func printBatchPlan(out io.Writer, plan workflow.BatchPlan) {
	fmt.Fprintf(out, "Approval batch preview\n- Grant: %s\n- Next gate: %s\n", plan.Grant, plan.NextGate)
	for _, item := range plan.Items {
		fmt.Fprintf(out, "APPROVE %s (%s) hash=%s\n", item.Artifact.ID, item.Artifact.Path, item.Hash)
	}
	for _, item := range plan.Ignored {
		fmt.Fprintf(out, "IGNORE %s (%s): %s\n", item.Artifact.ID, item.Artifact.Path, item.Reason)
	}
	for _, item := range plan.Blockers {
		fmt.Fprintf(out, "BLOCKED %s (%s): %s\n", item.Artifact.ID, item.Artifact.Path, item.Reason)
	}
	fmt.Fprintf(out, "Summary: approve=%d ignored=%d blocked=%d\n", len(plan.ToApprove), len(plan.Ignored), len(plan.Blockers))
}
