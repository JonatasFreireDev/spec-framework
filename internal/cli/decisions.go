package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/wizard"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func runDecisions(args []string, out, errout io.Writer) int {
	if len(args) == 0 || args[0] != "migrate" {
		fmt.Fprintln(errout, "decisions requires migrate")
		return 2
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
