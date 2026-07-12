package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecisionImpactReportsValidityReferencesAndGaps(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "knowledge", "decisions"), 0755)
	_ = os.MkdirAll(filepath.Join(root, ".product", "history"), 0755)
	decisionRel := "knowledge/decisions/DEC-010-queue.md"
	body := "# Decision\n\n- ID: DEC-010\n- Status: approved\n"
	_ = os.WriteFile(filepath.Join(root, filepath.FromSlash(decisionRel)), []byte(body), 0644)
	_ = os.WriteFile(filepath.Join(root, "plan.md"), []byte("Uses DEC-010\n"), 0644)
	_ = writeJSON(filepath.Join(root, ".product", "decisions.json"), map[string]any{"decisions": []any{map[string]any{"id": "DEC-010", "type": "architecture", "status": "approved", "path": decisionRel, "affectedArtifacts": []any{"plan.md", "missing.md"}, "workflowEffects": map[string]any{"requiredGates": []any{"GATE-QUEUE"}}}}})
	_ = writeJSON(filepath.Join(root, ".product", "history", "approval-dec-010.json"), Approval{ArtifactID: "DEC-010", Path: decisionRel, ContentHash: Hash(body), StatusGranted: "approved"})
	r, err := DecisionImpactReport(root, "DEC-010")
	if err != nil {
		t.Fatal(err)
	}
	if !r.Valid || len(r.References) != 1 || len(r.PropagationGaps) != 1 {
		t.Fatalf("%+v", r)
	}
}

func TestCommandPlanRejectsDecisionTextAsSource(t *testing.T) {
	root := t.TempDir()
	_, err := CreateCommandPlan(root, "W", "T", ".", "DEC-010", "R0", []string{"go", "test"}, 1)
	if err == nil {
		t.Fatal("decision was accepted as command source")
	}
}
