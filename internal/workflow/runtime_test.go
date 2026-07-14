package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func TestClaimLeaseRejectsExistingClaim(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	graph := filepath.Join(root, "graph.json")
	_ = writeJSON(graph, Graph{ID: "G", Nodes: []Node{{ID: "T1", Status: "pending"}}})
	if _, err := ClaimLease(root, graph, "T1", "agent-1", time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := ClaimLease(root, graph, "T1", "agent-2", time.Minute); err == nil || (!strings.Contains(err.Error(), "leased") && !strings.Contains(err.Error(), "claimed")) {
		t.Fatalf("expected second agent to be rejected, got %v", err)
	}
}

func TestExecuteCommandPlanRejectsMalformedArgv(t *testing.T) {
	root := t.TempDir()
	planDir := filepath.Join(root, ".product", "workspaces", "WORK-1", "command-plans")
	if err := os.MkdirAll(planDir, 0755); err != nil {
		t.Fatal(err)
	}
	plan := CommandPlan{Version: RuntimeVersion, ID: "CMDPLAN-001", WorkspaceID: "WORK-1", TaskID: "T1", Cwd: ".", Risk: "R0"}
	data, _ := json.Marshal(plan)
	if err := os.WriteFile(filepath.Join(planDir, "CMDPLAN-001.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ExecuteCommandPlan(root, "WORK-1", "CMDPLAN-001", "agent"); err == nil || !strings.Contains(err.Error(), "argv") {
		t.Fatalf("expected malformed argv error, got %v", err)
	}
}

func TestExecuteCommandPlanRejectsTamperedArgv(t *testing.T) {
	root := t.TempDir()
	planDir := filepath.Join(root, ".product", "workspaces", "WORK-1", "command-plans")
	claimsDir := filepath.Join(root, ".product", "claims")
	if err := os.MkdirAll(planDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(claimsDir, 0755); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal([]string{"echo", "original"})
	sum := sha256.Sum256(raw)
	plan := CommandPlan{Version: RuntimeVersion, ID: "CMDPLAN-001", WorkspaceID: "WORK-1", TaskID: "T1", Cwd: ".", Risk: "R0", Argv: []string{"echo", "tampered"}, InputHash: hex.EncodeToString(sum[:])}
	if err := writeJSON(filepath.Join(planDir, "CMDPLAN-001.json"), plan); err != nil {
		t.Fatal(err)
	}
	lease := Lease{Version: RuntimeVersion, TaskID: "T1", Agent: "agent", ExpiresAt: time.Now().Add(time.Minute).UTC().Format(time.RFC3339)}
	if err := writeJSON(filepath.Join(claimsDir, "T1.json"), lease); err != nil {
		t.Fatal(err)
	}
	if err := writeJSON(filepath.Join(root, ".product", "claims.json"), Claims{Claims: []Claim{{TaskID: "T1", Agent: "agent"}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := ExecuteCommandPlan(root, "WORK-1", "CMDPLAN-001", "agent"); err == nil || !strings.Contains(err.Error(), "hash") {
		t.Fatalf("expected hash mismatch error, got %v", err)
	}
}

func TestRuntimeComponentsCannotEscapeWorktree(t *testing.T) {
	if _, err := CreateTaskWorktree(t.TempDir(), "../WORK-1", "T1"); err == nil {
		t.Fatal("expected worktree path traversal to be rejected")
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
