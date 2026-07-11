package workflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func setupProduct(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	r := Registry{Artifacts: []Artifact{{ID: "DOMAIN-1", Type: "domain", Status: "approved", Path: "domains/events/context.md"}, {ID: "GOAL-1", Type: "user-goal", Status: "approved", Path: "domains/events/goals/manage/context.md", ParentIDs: []string{"DOMAIN-1"}}, {ID: "FT-1", Type: "feature", Status: "approved", Path: "domains/events/goals/manage/features/invites/context.md", ParentIDs: []string{"GOAL-1"}}}}
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), r); err != nil {
		t.Fatal(err)
	}
	for _, a := range r.Artifacts {
		path := filepath.Join(root, filepath.FromSlash(a.Path))
		_ = os.MkdirAll(filepath.Dir(path), 0755)
		_ = os.WriteFile(path, []byte("id: "+a.ID+"\nstatus: "+a.Status+"\n"), 0644)
	}
	return root
}
func TestWorkspaceSelectionAndStatus(t *testing.T) {
	root := setupProduct(t)
	w, err := CreateWorkspace(root, "FT-1", "", "", "", "tester")
	if err != nil {
		t.Fatal(err)
	}
	if w.ID != "WORK-001" {
		t.Fatal(w.ID)
	}
	s, err := WorkspaceStatus(root, w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Next != "use-case" {
		t.Fatalf("next=%s", s.Next)
	}
}
func TestApproveUpdatesStatusRegistryAndRecord(t *testing.T) {
	root := setupProduct(t)
	path := filepath.Join(root, "domains", "events", "goals", "manage", "features", "invites", "context.md")
	if _, err := Approve(root, path, "in_progress", "owner", "started"); err != nil {
		t.Fatal(err)
	}
	rec, err := Approve(root, path, "implemented", "owner", "reviewed")
	if err != nil {
		t.Fatal(err)
	}
	if rec.StatusGranted != "implemented" || rec.ContentHash == "" {
		t.Fatal(rec)
	}
	var r Registry
	if err := readJSON(filepath.Join(root, ".product", "artifacts.json"), &r); err != nil {
		t.Fatal(err)
	}
	if r.Artifacts[2].Status != "implemented" {
		t.Fatal(r.Artifacts[2].Status)
	}
}
func TestApproveBlocksUnapprovedParent(t *testing.T) {
	root := setupProduct(t)
	var r Registry
	_ = readJSON(filepath.Join(root, ".product", "artifacts.json"), &r)
	r.Artifacts[1].Status = "draft"
	_ = writeJSON(filepath.Join(root, ".product", "artifacts.json"), r)
	_, err := Approve(root, filepath.Join(root, filepath.FromSlash(r.Artifacts[2].Path)), "approved", "owner", "")
	if err == nil {
		t.Fatal("expected parent blocker")
	}
}
func TestGateReadiness(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "knowledge", "conventions")
	_ = os.MkdirAll(path, 0755)
	_ = os.WriteFile(filepath.Join(path, "gates.md"), []byte("| `GATE-TEST` | `TBD by adopter` | x |\n| `GATE-LINT` | `go vet` | x |"), 0644)
	missing, err := GateReadiness(root)
	if err != nil || len(missing) != 1 || missing[0] != "GATE-TEST" {
		t.Fatalf("%v %v", missing, err)
	}
}
func TestGraphClaimsAndCompletion(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	path := filepath.Join(root, "graph.json")
	g := Graph{ID: "G", Nodes: []Node{{ID: "T1", Path: "tasks/T1.md", Status: "pending"}, {ID: "T2", Path: "tasks/T2.md", Status: "pending", DependsOn: []string{"T1"}}}}
	_ = writeJSON(path, g)
	ready, _ := Ready(path)
	if len(ready) != 1 || ready[0].ID != "T1" {
		t.Fatal(ready)
	}
	if _, err := ClaimTask(root, path, "T1", "codex"); err != nil {
		t.Fatal(err)
	}
	if _, err := ClaimTask(root, path, "T1", "claude"); err == nil {
		t.Fatal("duplicate claim allowed")
	}
	if err := Complete(root, path, "T1", "codex"); err != nil {
		t.Fatal(err)
	}
	ready, _ = Ready(path)
	if len(ready) != 1 || ready[0].ID != "T2" {
		data, _ := json.Marshal(ready)
		t.Fatal(string(data))
	}
}

func TestConcurrentClaimHasSingleWinner(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	path := filepath.Join(root, "graph.json")
	_ = writeJSON(path, Graph{ID: "G", Nodes: []Node{{ID: "T1", Status: "pending"}}})
	var wins atomic.Int32
	var wg sync.WaitGroup
	for range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := ClaimTask(root, path, "T1", "agent"); err == nil {
				wins.Add(1)
			}
		}()
	}
	wg.Wait()
	if wins.Load() != 1 {
		t.Fatalf("winners=%d", wins.Load())
	}
}

func TestClaimRejectsParallelScopeConflict(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	path := filepath.Join(root, "graph.json")
	_ = writeJSON(path, Graph{ID: "G", Nodes: []Node{{ID: "T1", Status: "pending", WriteScope: []string{"src/events"}}, {ID: "T2", Status: "pending", WriteScope: []string{"src/events/api"}}}})
	if _, err := ClaimTask(root, path, "T1", "codex"); err != nil {
		t.Fatal(err)
	}
	if _, err := ClaimTask(root, path, "T2", "claude"); err == nil {
		t.Fatal("overlapping scope claim allowed")
	}
}
