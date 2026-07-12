package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecisionMigrationPlansAppliesBacksUpAndIsIdempotent(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	legacy := map[string]any{"decisions": []any{map[string]any{"id": "DEC-001", "status": "approved", "scope": "architecture/security", "affectedArtifacts": []any{"domains/events/plan.md"}}}}
	_ = writeJSON(filepath.Join(root, ".product", "decisions.json"), legacy)
	plan, err := PlanDecisionMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Items) != 1 || !plan.Items[0].NeedsReview {
		t.Fatalf("%+v", plan.Items)
	}
	plan.Items[0].InferredType = "architecture"
	result, err := ApplyDecisionMigration(plan, plan.Items)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 1 {
		t.Fatal(result)
	}
	if _, err = os.Stat(result.Backup); err != nil {
		t.Fatal(err)
	}
	next, err := PlanDecisionMigration(root)
	if err != nil || len(next.Items) != 0 {
		t.Fatalf("%+v %v", next.Items, err)
	}
}

func TestBuildDashboardShowsCanonicalStages(t *testing.T) {
	root := setupProduct(t)
	w, err := CreateWorkspace(root, "FT-1", "", "", "", "tester")
	if err != nil {
		t.Fatal(err)
	}
	d, err := BuildDashboard(root, w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Stages) != 10 || d.CurrentStep != "use-case" {
		t.Fatalf("%+v", d)
	}
}
