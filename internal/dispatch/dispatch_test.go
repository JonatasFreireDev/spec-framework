package dispatch

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestReviewEnvelopePinsParentDiffAndCannotRun(t *testing.T) {
	root := t.TempDir()
	work := "WORK-001"
	if err := os.MkdirAll(dir(root, work), 0755); err != nil {
		t.Fatal(err)
	}
	parent := Envelope{Version: 1, ID: "DISPATCH-1", WorkspaceID: work, TaskID: "TK-1", Role: "code-runner", Agent: "runner", Status: "returned", DiffHash: "abc", InputHash: "input"}
	if err := write(filepath.Join(dir(root, work), parent.ID+".json"), parent); err != nil {
		t.Fatal(err)
	}
	review, err := AssignReview(root, work, parent.ID, "qa", "qa-1")
	if err != nil {
		t.Fatal(err)
	}
	if review.DiffHash != "abc" || review.ParentID != parent.ID || len(review.WriteScope) != 0 {
		t.Fatalf("review=%+v", review)
	}
	if _, err := Run(root, work, review.ID, false, "echo", nil); err == nil {
		t.Fatal("review run accepted")
	}
}

func TestReconcileReportsOrphanReview(t *testing.T) {
	root := t.TempDir()
	work := "WORK-001"
	if err := os.MkdirAll(dir(root, work), 0755); err != nil {
		t.Fatal(err)
	}
	if err := write(filepath.Join(dir(root, work), "DISPATCH-2.json"), Envelope{ID: "DISPATCH-2", WorkspaceID: work, Role: "code-review", ParentID: "missing", DiffHash: "abc", Status: "assigned"}); err != nil {
		t.Fatal(err)
	}
	items, err := Reconcile(root, work)
	if err != nil || len(items) != 1 || items[0].Kind != "orphaned-review" {
		t.Fatalf("items=%+v err=%v", items, err)
	}
}

func TestEngineeringDelegationEnforcesPhasesContextScopeAndReturns(t *testing.T) {
	root := t.TempDir()
	work := "WORK-ENG-001"
	handoffPath := filepath.Join(root, ".product", "workspaces", work, "engineering-handoff.json")
	handoff := engineeringHandoff{
		SchemaVersion: 1,
		Execution:     engineeringExecution{Mode: "delegated", ContextPolicy: "minimal", MaxParallel: 2, Fallback: "sequential"},
		Routes: []engineeringRoute{
			{Skill: "technical-landscape", Phase: 1, WriteScope: []string{"engineering/architecture", "engineering/catalog"}, Status: "pending"},
			{Skill: "engineering-standards", Phase: 2, DependsOn: []string{"technical-landscape"}, WriteScope: []string{"engineering/standards"}, Status: "pending"},
			{Skill: "operations-baseline", Phase: 2, DependsOn: []string{"technical-landscape"}, WriteScope: []string{"engineering/operations"}, Status: "pending"},
			{Skill: "engineering-evidence", Phase: 3, DependsOn: []string{"engineering-standards", "operations-baseline"}, WriteScope: []string{"engineering/evidence"}, Status: "pending"},
			{Skill: "engineering-system", Phase: 4, DependsOn: []string{"technical-landscape", "engineering-standards", "operations-baseline", "engineering-evidence"}, WriteScope: []string{"engineering/engineering-system.md", "engineering/engineering-system.yaml", "engineering/quality"}, Status: "pending"},
		},
	}
	data, err := json.MarshalIndent(handoff, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(handoffPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(handoffPath, append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}

	technical, err := AssignEngineering(root, work, handoffPath, "technical-landscape", "landscape-agent", nil)
	if err != nil {
		t.Fatal(err)
	}
	if technical.ContextPolicy != "minimal" || technical.Phase != 1 || technical.UnitKind != "engineering-specialist" {
		t.Fatalf("technical envelope=%+v", technical)
	}
	if _, err := AssignEngineering(root, work, handoffPath, "engineering-standards", "standards-agent", nil); err == nil {
		t.Fatal("standards assignment accepted without returned technical dependency")
	}
	technicalOutput := writeEngineeringOutput(t, root, "engineering/catalog/catalog.yaml", "schema_version: 1\n")
	if _, err := ReturnEngineering(root, work, technical.ID, technical.Agent, "technical graph mapped", []string{"catalog evidence"}, []string{"../escape=" + technicalOutput}, nil, nil); err == nil {
		t.Fatal("engineering return accepted output outside write scope")
	}
	technical, err = ReturnEngineering(root, work, technical.ID, technical.Agent, "technical graph mapped", []string{"catalog evidence"}, []string{"engineering/catalog/catalog.yaml=" + technicalOutput}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	standards, err := AssignEngineering(root, work, handoffPath, "engineering-standards", "standards-agent", []string{technical.ID})
	if err != nil {
		t.Fatal(err)
	}
	operations, err := AssignEngineering(root, work, handoffPath, "operations-baseline", "operations-agent", []string{technical.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := AssignEngineering(root, work, handoffPath, "engineering-evidence", "evidence-agent", []string{standards.ID, operations.ID}); err == nil {
		t.Fatal("assignment exceeded max_parallel or accepted unreturned dependencies")
	}
	standardsHash := writeEngineeringOutput(t, root, "engineering/standards/standards.yaml", "schema_version: 1\n")
	standardsReturned, err := ReturnEngineering(root, work, standards.ID, standards.Agent, "standards mapped", []string{"standards evidence"}, []string{"engineering/standards/standards.yaml=" + standardsHash}, nil, []string{"DEC-STANDARDS-001"})
	if err != nil {
		t.Fatal(err)
	}
	if len(standardsReturned.DecisionCandidates) != 1 || standardsReturned.DecisionCandidates[0] != "DEC-STANDARDS-001" {
		t.Fatalf("decision candidates were not persisted: %+v", standardsReturned)
	}
	operationsHash := writeEngineeringOutput(t, root, "engineering/operations/operations.yaml", "schema_version: 1\n")
	if _, err := ReturnEngineering(root, work, operations.ID, operations.Agent, "operations mapped", []string{"operations evidence"}, []string{"engineering/operations/operations.yaml=" + operationsHash}, nil, nil); err != nil {
		t.Fatal(err)
	}
	evidence, err := AssignEngineering(root, work, handoffPath, "engineering-evidence", "evidence-agent", []string{standards.ID, operations.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(evidence.Dependencies) != 2 || evidence.Phase != 3 {
		t.Fatalf("evidence envelope=%+v", evidence)
	}
	if err := os.WriteFile(filepath.Join(root, "engineering", "standards", "standards.yaml"), []byte("changed\n"), 0644); err != nil {
		t.Fatal(err)
	}
	findings, err := Reconcile(root, work)
	if err != nil {
		t.Fatal(err)
	}
	foundStale := false
	for _, finding := range findings {
		if finding.Kind == "engineering-output-stale" && finding.DispatchID == standards.ID {
			foundStale = true
		}
		if finding.Kind == "review-missing-diff" && finding.DispatchID == evidence.ID {
			t.Fatal("engineering assignment was misclassified as a review")
		}
	}
	if !foundStale {
		t.Fatalf("missing stale engineering output finding: %+v", findings)
	}
}

func TestEngineeringDelegationRejectsSequentialHandoff(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".product", "workspaces", "WORK-001", "handoff.json")
	data := []byte(`{"schema_version":1,"execution":{"mode":"sequential","context_policy":"minimal","max_parallel":1,"fallback":"sequential"},"routes":[{"skill":"technical-landscape","phase":1,"write_scope":["engineering/catalog"],"status":"pending"}]}`)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := AssignEngineering(root, "WORK-001", path, "technical-landscape", "agent", nil); err == nil {
		t.Fatal("sequential engineering handoff accepted for subagent assignment")
	}
}

func TestEngineeringReturnRejectsStaleHandoff(t *testing.T) {
	root := t.TempDir()
	work := "WORK-001"
	path := filepath.Join(root, ".product", "workspaces", work, "handoff.json")
	data := []byte(`{"schema_version":1,"execution":{"mode":"delegated","context_policy":"minimal","max_parallel":1,"fallback":"sequential"},"routes":[{"skill":"technical-landscape","phase":1,"write_scope":["engineering/catalog"],"status":"pending"}]}`)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	envelope, err := AssignEngineering(root, work, path, "technical-landscape", "agent", nil)
	if err != nil {
		t.Fatal(err)
	}
	hash := writeEngineeringOutput(t, root, "engineering/catalog/catalog.yaml", "schema_version: 1\n")
	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReturnEngineering(root, work, envelope.ID, envelope.Agent, "done", []string{"evidence"}, []string{"engineering/catalog/catalog.yaml=" + hash}, nil, nil); err == nil {
		t.Fatal("engineering return accepted a stale handoff")
	}
}

func TestEngineeringAssignmentRejectsEscapedHandoffAndWriteScope(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "outside.json")
	data := []byte(`{"schema_version":1,"execution":{"mode":"delegated","context_policy":"minimal","max_parallel":1,"fallback":"sequential"},"routes":[{"skill":"technical-landscape","phase":1,"write_scope":["engineering"],"status":"pending"}]}`)
	if err := os.WriteFile(outside, data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := AssignEngineering(root, "WORK-001", outside, "technical-landscape", "agent", nil); err == nil {
		t.Fatal("engineering assignment accepted a handoff outside its workspace")
	}
	inside := filepath.Join(root, ".product", "workspaces", "WORK-001", "handoff.json")
	if err := os.MkdirAll(filepath.Dir(inside), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(inside, data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := AssignEngineering(root, "WORK-001", inside, "technical-landscape", "agent", nil); err == nil {
		t.Fatal("engineering assignment accepted a write scope outside specialist ownership")
	}
}

func TestEngineeringBlockedReturnDoesNotUnlockDependentPhase(t *testing.T) {
	root := t.TempDir()
	work := "WORK-001"
	path := filepath.Join(root, ".product", "workspaces", work, "handoff.json")
	data := []byte(`{"schema_version":1,"execution":{"mode":"delegated","context_policy":"minimal","max_parallel":1,"fallback":"sequential"},"routes":[{"skill":"technical-landscape","phase":1,"depends_on":[],"write_scope":["engineering/catalog"],"status":"pending"},{"skill":"engineering-standards","phase":2,"depends_on":["technical-landscape"],"write_scope":["engineering/standards"],"status":"pending"}]}`)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	technical, err := AssignEngineering(root, work, path, "technical-landscape", "agent", nil)
	if err != nil {
		t.Fatal(err)
	}
	hash := writeEngineeringOutput(t, root, "engineering/catalog/catalog.yaml", "schema_version: 1\n")
	technical, err = ReturnEngineering(root, work, technical.ID, technical.Agent, "partial", []string{"catalog"}, []string{"engineering/catalog/catalog.yaml=" + hash}, []string{"repository boundary unresolved"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := AssignEngineering(root, work, path, "engineering-standards", "standards", []string{technical.ID}); err == nil {
		t.Fatal("blocked engineering return unlocked a dependent phase")
	}
}

func TestEngineeringAssignmentCapacityIsAtomic(t *testing.T) {
	root := t.TempDir()
	work := "WORK-001"
	path := filepath.Join(root, ".product", "workspaces", work, "handoff.json")
	data := []byte(`{"schema_version":1,"execution":{"mode":"delegated","context_policy":"minimal","max_parallel":1,"fallback":"sequential"},"routes":[{"skill":"technical-landscape","phase":1,"depends_on":[],"write_scope":["engineering/catalog"],"status":"pending"},{"skill":"engineering-standards","phase":1,"depends_on":[],"write_scope":["engineering/standards"],"status":"pending"}]}`)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	var wait sync.WaitGroup
	for _, assignment := range []struct{ role, agent string }{{"technical-landscape", "agent-1"}, {"engineering-standards", "agent-2"}} {
		assignment := assignment
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			_, err := AssignEngineering(root, work, path, assignment.role, assignment.agent, nil)
			results <- err
		}()
	}
	close(start)
	wait.Wait()
	close(results)
	succeeded := 0
	failed := 0
	for err := range results {
		if err == nil {
			succeeded++
		} else {
			failed++
		}
	}
	if succeeded != 1 || failed != 1 {
		t.Fatalf("atomic capacity succeeded=%d failed=%d", succeeded, failed)
	}
}

func writeEngineeringOutput(t *testing.T, root, relative, content string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
