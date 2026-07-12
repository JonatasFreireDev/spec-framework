package workflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
		text := "id: " + a.ID + "\nstatus: " + a.Status + "\n"
		_ = os.WriteFile(path, []byte(text), 0644)
		_ = os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755)
		_ = writeJSON(filepath.Join(root, ".product", "history", "approval-"+strings.ToLower(a.ID)+".json"), Approval{ArtifactID: a.ID, Path: a.Path, ContentHash: Hash(text), StatusGranted: a.Status, ApprovedBy: "Test Owner"})
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

func TestApprovedDocumentRequiresCurrentApprovalRecord(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "artifact.md")
	text := "# Artifact\n\n| Field | Value |\n| --- | --- |\n| Status | `approved` |\n"
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	if approvedDocument(root, path) {
		t.Fatal("status text bypassed approval history")
	}
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	record := Approval{ArtifactID: "ART-1", Path: "artifact.md", ContentHash: Hash(text), StatusGranted: "approved", ApprovedBy: "Human"}
	if err := writeJSON(filepath.Join(root, ".product", "history", "approval-art-1.json"), record); err != nil {
		t.Fatal(err)
	}
	if !approvedDocument(root, path) {
		t.Fatal("current approval record was not accepted")
	}
}

func TestPassedEngineeringReviewBecomesStaleAfterProposalChange(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "use-case")
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	proposalPath := filepath.Join(dir, "engineering-proposal.md")
	reviewPath := filepath.Join(dir, "engineering-review.md")
	proposal := "# Proposal\n"
	if err := os.WriteFile(proposalPath, []byte(proposal), 0o644); err != nil {
		t.Fatal(err)
	}
	review := "# Review\n\n| Field | Value |\n| --- | --- |\n| Status | `approved` |\n| Verdict | `passed` |\n| Proposal hash | `" + Hash(proposal) + "` |\n"
	if err := os.WriteFile(reviewPath, []byte(review), 0o644); err != nil {
		t.Fatal(err)
	}
	record := Approval{ArtifactID: "ENGREV-1", Path: "use-case/engineering-review.md", ContentHash: Hash(review), StatusGranted: "approved", ApprovedBy: "Human"}
	if err := writeJSON(filepath.Join(root, ".product", "history", "approval-engrev-1.json"), record); err != nil {
		t.Fatal(err)
	}
	if !passedEngineeringReview(root, reviewPath, proposalPath) {
		t.Fatal("current passed review was rejected")
	}
	if err := os.WriteFile(proposalPath, []byte("# Changed Proposal\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if passedEngineeringReview(root, reviewPath, proposalPath) {
		t.Fatal("stale review was accepted after proposal change")
	}
}

func TestEngineeringSystemApprovalIsAtomicAndComposite(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product"), 0o755); err != nil {
		t.Fatal(err)
	}
	engineering := filepath.Join(root, "engineering")
	if err := os.MkdirAll(filepath.Join(engineering, "architecture"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"context.md":              "---\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering-system.md":   "| Field | Value |\n| --- | --- |\n| ID | `ENGSYS-001` |\n| Status | `draft` |\n| Version | `1.0.0` |\n",
		"engineering-system.yaml": "schema_version: 1\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: product\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: baseline\n    evidence: []\n",
		"architecture/modules.md": "# Modules\n",
	}
	for name, body := range files {
		path := filepath.Join(engineering, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	registry := Registry{Artifacts: []Artifact{{ID: "ENGSYS-001", Type: "engineering-system", Status: "draft", Path: "engineering/engineering-system.md"}}}
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	canonical := filepath.Join(engineering, "engineering-system.md")
	if _, err := Approve(root, canonical, "approved", "Human", "Reviewed composite system"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"context.md", "engineering-system.md", "engineering-system.yaml"} {
		data, _ := os.ReadFile(filepath.Join(engineering, name))
		if !strings.Contains(string(data), "approved") {
			t.Fatalf("%s was not synchronized: %s", name, data)
		}
	}
	if !hasCurrentApproval(root, canonical, "approved") {
		t.Fatal("composite approval was not current")
	}
	if err := os.WriteFile(filepath.Join(engineering, "architecture", "modules.md"), []byte("# Changed Modules\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if hasCurrentApproval(root, canonical, "approved") {
		t.Fatal("engineering contract change did not stale composite approval")
	}
}

func TestStructuredNotApplicableRejectsIncidentalProse(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "design.md")
	if err := os.WriteFile(path, []byte("Status: draft\n\nNot applicable was rejected; rationale pending.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if notApplicableDocument(path) {
		t.Fatal("incidental prose bypassed Not applicable gate")
	}
	if err := os.WriteFile(path, []byte("Status: not_applicable\nRationale: This delivery has no user interface.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !notApplicableDocument(path) {
		t.Fatal("structured Not applicable contract was rejected")
	}
}

func TestArchitectureGateRequiresCurrentScopedDecision(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	discoveryRel := "domains/events/use-cases/check-in/technical-discovery.md"
	discoveryPath := filepath.Join(root, filepath.FromSlash(discoveryRel))
	if err := os.MkdirAll(filepath.Dir(discoveryPath), 0o755); err != nil {
		t.Fatal(err)
	}
	discovery := "## Architecture Gate\n\n| Field | Value |\n| --- | --- |\n| Verdict | Decision required |\n| Decision | DEC-001 |\n| Rationale | The delivery changes a service boundary. |\n"
	if err := os.WriteFile(discoveryPath, []byte(discovery), 0o644); err != nil {
		t.Fatal(err)
	}
	decisionRel := "knowledge/decisions/DEC-001.md"
	decisionPath := filepath.Join(root, filepath.FromSlash(decisionRel))
	if err := os.MkdirAll(filepath.Dir(decisionPath), 0o755); err != nil {
		t.Fatal(err)
	}
	decision := "# Decision\nstatus: approved\n"
	if err := os.WriteFile(decisionPath, []byte(decision), 0o644); err != nil {
		t.Fatal(err)
	}
	index := map[string]any{"decisions": []any{map[string]any{"id": "DEC-001", "status": "approved", "scope": "architecture", "path": decisionRel}}}
	if err := writeJSON(filepath.Join(root, ".product", "decisions.json"), index); err != nil {
		t.Fatal(err)
	}
	if architectureResolved(root, discoveryPath) {
		t.Fatal("decision without approval record resolved Architecture Gate")
	}
	record := Approval{ArtifactID: "DEC-001", Path: decisionRel, ContentHash: Hash(decision), StatusGranted: "approved", ApprovedBy: "Human"}
	if err := writeJSON(filepath.Join(root, ".product", "history", "approval-dec-001.json"), record); err != nil {
		t.Fatal(err)
	}
	if !architectureResolved(root, discoveryPath) {
		t.Fatal("current scope-compatible decision did not resolve Architecture Gate")
	}
}
