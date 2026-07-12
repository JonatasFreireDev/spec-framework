package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMaterializeTasksClosesGraphCycle(t *testing.T) {
	root := t.TempDir()
	graph := filepath.Join(root, "execution-graph.json")
	raw := map[string]any{"id": "GRAPH-1", "status": "proposed", "sourceSpecification": "SPEC-1", "nodes": []any{map[string]any{"id": "TK-1", "path": "tasks/TK-1.md", "title": "First slice", "type": "backend", "ownerSkill": "code-runner", "status": "pending", "dependsOn": []any{}, "sourceSections": []any{"Behavior"}, "requirements": []any{"REQ-001"}, "acceptanceCriteria": []any{"AC-001"}, "plannedTests": []any{"TEST-001"}, "writeScope": []any{"internal/x"}, "sharedResources": []any{}, "acceptanceChecks": []any{"works"}, "delivery": map[string]any{"level": "L1", "priority": "P0"}}}}
	if err := writeJSON(graph, raw); err != nil {
		t.Fatal(err)
	}
	r, err := MaterializeTasks(graph)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Tasks) != 1 {
		t.Fatal(r)
	}
	if _, err = os.Stat(filepath.Join(root, "tasks", "TK-1.md")); err != nil {
		t.Fatal(err)
	}
	var updated map[string]any
	_ = readJSON(graph, &updated)
	if updated["status"] != "materialized" {
		t.Fatal(updated["status"])
	}
	if _, err = MaterializeTasks(graph); err == nil {
		t.Fatal("second materialization should fail")
	}
}

func TestReadinessReportsBlockers(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "knowledge", "conventions"), 0755)
	_ = os.WriteFile(filepath.Join(root, "knowledge", "conventions", "gates.md"), []byte("| `GATE-TEST` | `TBD` | |\n"), 0644)
	graph := filepath.Join(root, "execution-graph.json")
	_ = writeJSON(graph, Graph{ID: "G", Nodes: []Node{{ID: "TK-1", Path: "tasks/TK-1.md", Status: "pending", WriteScope: []string{"src"}}}})
	_ = os.MkdirAll(filepath.Join(root, "tasks"), 0755)
	_ = os.WriteFile(filepath.Join(root, "tasks", "TK-1.md"), []byte("| Status | `draft` |\n"), 0644)
	r, err := CheckTaskReadiness(root, graph, "TK-1")
	if err != nil {
		t.Fatal(err)
	}
	if r.Ready {
		t.Fatal("draft task with TBD gates was ready")
	}
}

func TestWorkspaceGuideUsesCanonicalStateMachine(t *testing.T) {
	root := setupProduct(t)
	w, err := CreateWorkspace(root, "FT-1", "", "", "", "tester")
	if err != nil {
		t.Fatal(err)
	}
	g, err := WorkspaceGuide(root, w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if g.RecommendedSkill != "use-case" || g.ExpectedArtifact == "" {
		t.Fatalf("%+v", g)
	}
}

func TestApproveStageUsesIndividualApprovalRecords(t *testing.T) {
	root := setupProduct(t)
	var reg Registry
	_ = readJSON(filepath.Join(root, ".product", "artifacts.json"), &reg)
	reg.Artifacts[2].Status = "draft"
	_ = writeJSON(filepath.Join(root, ".product", "artifacts.json"), reg)
	_ = os.WriteFile(filepath.Join(root, filepath.FromSlash(reg.Artifacts[2].Path)), []byte("id: FT-1\nstatus: draft\n"), 0644)
	w, err := CreateWorkspace(root, "FT-1", "", "", "", "tester")
	if err != nil {
		t.Fatal(err)
	}
	records, err := ApproveStage(root, w.ID, "feature", "owner", "stage review")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].ArtifactID != "FT-1" {
		t.Fatalf("%+v", records)
	}
}

func TestReadinessEnforcesDecisionEffects(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	_ = os.MkdirAll(filepath.Join(root, "knowledge", "conventions"), 0755)
	_ = os.WriteFile(filepath.Join(root, "knowledge", "conventions", "gates.md"), []byte("| `GATE-BASE` | `go test` | |\n"), 0644)
	_ = writeJSON(filepath.Join(root, ".product", "decisions.json"), map[string]any{"decisions": []any{map[string]any{"id": "DEC-010", "status": "approved", "workflowEffects": map[string]any{"requiredGates": []any{"GATE-DECISION"}}}}})
	graph := filepath.Join(root, "execution-graph.json")
	_ = writeJSON(graph, map[string]any{"id": "G", "status": "approved", "nodes": []any{map[string]any{"id": "TK-1", "path": "tasks/TK-1.md", "type": "backend", "status": "pending", "dependsOn": []any{}, "writeScope": []any{"src"}, "sharedResources": []any{}}}})
	_ = os.MkdirAll(filepath.Join(root, "tasks"), 0755)
	_ = os.WriteFile(filepath.Join(root, "tasks", "TK-1.md"), []byte("| Status | `draft` |\nREQ-001 AC-001 TEST-001 DEC-010\n"), 0644)
	r, err := CheckTaskReadiness(root, graph, "TK-1")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, c := range r.Checks {
		if c.ID == "decision-DEC-010-gate-GATE-DECISION" && c.Status == "block" {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing decision effect check: %+v", r.Checks)
	}
}
