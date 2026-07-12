package workflow

import (
	"os"
	"path/filepath"
	"strings"
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
	if len(d.Stages) != 12 || d.CurrentStep != "use-case" || d.Stages[6].ID != "engineering-proposal" || d.Stages[7].ID != "engineering-review" {
		t.Fatalf("%+v", d)
	}
}

func TestBuildDashboardSurfacesEngineeringSystemBlockers(t *testing.T) {
	root := setupProduct(t)
	engineering := filepath.Join(root, "engineering")
	if err := os.MkdirAll(filepath.Join(engineering, "architecture"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"context.md":                     "---\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering-system.md":          "# System\n",
		"architecture/system-context.md": "# Context\n",
		"engineering-system.yaml":        "schema_version: 1\nid: ENGSYS-OTHER\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: product\nareas:\n  context:\n    contract: architecture/system-context.md\n    maturity: baseline\n    evidence: []\n",
	} {
		if err := os.WriteFile(filepath.Join(engineering, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	w, err := CreateWorkspace(root, "FT-1", "", "", "", "tester")
	if err != nil {
		t.Fatal(err)
	}
	dashboard, err := BuildDashboard(root, w.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, blocker := range dashboard.Blockers {
		if strings.Contains(blocker, "context and catalog id do not match") {
			found = true
		}
	}
	if !found {
		t.Fatalf("blockers=%v", dashboard.Blockers)
	}
}
