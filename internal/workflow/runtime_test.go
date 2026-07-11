package workflow

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRuntimeV2WorkspaceResumeAndArtifacts(t *testing.T) {
	root := setupProduct(t)
	w, err := CreateWorkspace(root, "FT-1", "", "", "", "tester")
	if err != nil {
		t.Fatal(err)
	}
	s, err := Resume(root, w.ID)
	if err != nil || s.Version != 2 {
		t.Fatalf("%+v %v", s, err)
	}
	if _, err = WriteHandoff(root, w.ID, "delivery-orchestrator", "use-case", "continue"); err != nil {
		t.Fatal(err)
	}
	if _, err = WriteCheckpoint(root, w.ID, "use-case", "abc", "in", "out"); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeMigratesLegacyWorkspace(t *testing.T) {
	root := setupProduct(t)
	dir := filepath.Join(root, ".product", "workspaces")
	_ = os.MkdirAll(dir, 0755)
	w := Workspace{ID: "WORK-009", CurrentStep: "feature", CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	if err := writeJSON(filepath.Join(dir, "WORK-009.json"), w); err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateWorkspace(root, "WORK-009", false); err != nil {
		t.Fatal(err)
	}
	s, err := Resume(root, "WORK-009")
	if err != nil || s.Version != 2 {
		t.Fatalf("%+v %v", s, err)
	}
}

func TestLeaseHeartbeatAndRecovery(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	graph := filepath.Join(root, "graph.json")
	_ = writeJSON(graph, Graph{ID: "G", Nodes: []Node{{ID: "T1", Status: "pending"}}})
	l, err := ClaimLease(root, graph, "T1", "agent", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = Heartbeat(root, "T1", "agent", time.Minute); err != nil {
		t.Fatal(err)
	}
	l.ExpiresAt = time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	_ = writeJSON(leasePath(root, "T1"), l)
	x, err := RecoverLeases(root)
	if err != nil || len(x) != 1 {
		t.Fatalf("%v %v", x, err)
	}
}

func TestScheduleSerializesConflicts(t *testing.T) {
	root := t.TempDir()
	graph := filepath.Join(root, "graph.json")
	_ = writeJSON(graph, Graph{ID: "G", Nodes: []Node{{ID: "A", Status: "pending", WriteScope: []string{"src"}}, {ID: "B", Status: "pending", WriteScope: []string{"src/api"}}, {ID: "C", Status: "pending", WriteScope: []string{"docs"}}}})
	s, err := BuildSchedule(root, "WORK-1", graph, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Waves) != 2 || len(s.Waves[0].Tasks) != 2 {
		t.Fatalf("%+v", s.Waves)
	}
}

func TestCommandPlanRejectsR2(t *testing.T) {
	root := t.TempDir()
	_, err := CreateCommandPlan(root, "W", "T", ".", "test", "R2", []string{"go", "test"}, 1)
	if err == nil {
		t.Fatal("R2 accepted")
	}
}
