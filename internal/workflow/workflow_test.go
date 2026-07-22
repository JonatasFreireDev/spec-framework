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
	if err := writeJSON(filepath.Join(root, ".product", "framework.json"), map[string]any{"starting_point": "new-product"}); err != nil {
		t.Fatal(err)
	}
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

func TestExistingFeatureWorkspaceRequiresCurrentFeatureBriefApproval(t *testing.T) {
	root := setupProduct(t)
	if err := writeJSON(filepath.Join(root, ".product", "framework.json"), map[string]any{"starting_point": "existing-feature"}); err != nil {
		t.Fatal(err)
	}
	brief := Artifact{ID: "FEATURE-BRIEF-1", Type: "feature-brief", Status: "draft", Path: "foundation/feature-brief.md", TargetFeature: "FT-1"}
	registry, err := LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	registry.Artifacts = append([]Artifact{brief}, registry.Artifacts...)
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, filepath.FromSlash(brief.Path))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("| ID | FEATURE-BRIEF-1 |\n| Type | feature-brief |\n| Status | draft |\n| Target Feature | FT-1 |\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "lacks current approval evidence") {
		t.Fatalf("expected Feature Brief approval blocker, got %v", err)
	}
	if _, err := Approve(root, path, "approved", "Product Owner", "bounded feature approved"); err != nil {
		t.Fatal(err)
	}
	registry, err = LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	registry.Artifacts[0].TargetFeature = "FT-OTHER"
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "registry target") {
		t.Fatalf("expected Feature Brief registry tampering blocker, got %v", err)
	}
	registry.Artifacts[0].TargetFeature = "FT-1"
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err != nil {
		t.Fatalf("workspace remained blocked after approval: %v", err)
	}
}

func TestExistingFeatureWorkspaceRequiresSharedBaselines(t *testing.T) {
	root := setupProduct(t)
	if err := writeJSON(filepath.Join(root, ".product", "framework.json"), map[string]any{"starting_point": "existing-feature"}); err != nil {
		t.Fatal(err)
	}
	registry, err := LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	baselines := []Artifact{
		{ID: "FEATURE-BRIEF-1", Type: "feature-brief", Status: "approved", Path: "foundation/feature-brief.md", TargetFeature: "FT-1"},
		{ID: "LANDSCAPE-1", Type: "product-landscape", Status: "draft", Path: "knowledge/assessments/product-landscape.md"},
		{ID: "ENGSYS-1", Type: "engineering-system", Status: "approved", Path: "engineering/engineering-system.md"},
		{ID: "DSYS-1", Type: "design-system", Status: "approved", Path: "design/system/design-system.md"},
	}
	registry.RequiredArtifacts = []ArtifactRequirement{
		{Type: "feature-brief", Path: "foundation/feature-brief.md"},
		{Type: "product-landscape", Path: "knowledge/assessments/product-landscape.md"},
		{Type: "engineering-system", Path: "engineering/engineering-system.md"},
		{Type: "design-system", Path: "design/system/design-system.md"},
	}
	registry.Artifacts = append(baselines, registry.Artifacts...)
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	for _, artifact := range baselines {
		path := filepath.Join(root, filepath.FromSlash(artifact.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		text := "id: " + artifact.ID + "\nstatus: " + artifact.Status + "\n"
		if artifact.Type == "feature-brief" {
			text = "| ID | FEATURE-BRIEF-1 |\n| Status | approved |\n| Target Feature | FT-1 |\n"
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
		if artifact.Status == "approved" {
			if err := writeJSON(filepath.Join(root, ".product", "history", "approval-"+strings.ToLower(artifact.ID)+".json"), Approval{ArtifactID: artifact.ID, Path: artifact.Path, ContentHash: Hash(text), StatusGranted: artifact.Status, ApprovedBy: "Test Owner"}); err != nil {
				t.Fatal(err)
			}
		}
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "product landscape lacks current approval evidence") {
		t.Fatalf("expected Product Landscape blocker, got %v", err)
	}
}

func TestFeatureTargetMatchesCanonicalAndContextPaths(t *testing.T) {
	feature := Artifact{ID: "FT-1", Path: "domains/d/goals/g/features/f/context.md"}
	for _, target := range []string{"FT-1", "domains/d/goals/g/features/f/context.md", "domains/d/goals/g/features/f/feature.md"} {
		if !featureTargetMatches(target, feature) {
			t.Errorf("target %q should match", target)
		}
	}
	if featureTargetMatches("FT-OTHER", feature) {
		t.Fatal("unrelated feature target matched")
	}
}

func TestExistingImplementationWorkspaceRequiresCurrentAssessmentApproval(t *testing.T) {
	root := setupProduct(t)
	if err := writeJSON(filepath.Join(root, ".product", "framework.json"), map[string]any{"starting_point": "existing-implementation"}); err != nil {
		t.Fatal(err)
	}
	artifacts := []Artifact{
		{ID: "IMPL-ASSESS-1", Type: "implementation-assessment", Status: "draft", Path: "knowledge/assessments/implementation-assessment.md"},
		{ID: "PROBLEM-1", Type: "problem", Status: "draft", Path: "foundation/problem/problem.md", ParentIDs: []string{"IMPL-ASSESS-1"}},
		{ID: "VISION-1", Type: "vision", Status: "draft", Path: "foundation/vision/vision.md", ParentIDs: []string{"PROBLEM-1"}},
		{ID: "PRINCIPLES-1", Type: "product-principles", Status: "draft", Path: "foundation/vision/principles.md", ParentIDs: []string{"VISION-1"}},
		{ID: "NORTH-STAR-1", Type: "north-star", Status: "draft", Path: "foundation/vision/north-star.md", ParentIDs: []string{"VISION-1"}},
		{ID: "STRATEGY-1", Type: "strategy", Status: "draft", Path: "foundation/strategy/strategy.md", ParentIDs: []string{"VISION-1", "PRINCIPLES-1", "NORTH-STAR-1"}},
	}
	registry, err := LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	registry.Artifacts = append(artifacts, registry.Artifacts...)
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	for _, artifact := range artifacts {
		path := filepath.Join(root, filepath.FromSlash(artifact.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("| ID | "+artifact.ID+" |\n| Type | "+artifact.Type+" |\n| Status | draft |\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if map[string]bool{"problem": true, "vision": true, "strategy": true}[artifact.Type] {
			if err := os.WriteFile(filepath.Join(filepath.Dir(path), "context.md"), []byte("status: draft\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "implementation assessment lacks current approval evidence") {
		t.Fatalf("expected assessment approval blocker, got %v", err)
	}
	if _, err := Approve(root, filepath.Join(root, filepath.FromSlash(artifacts[0].Path)), "approved", "Product Owner", "implementation evidence reviewed"); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "problem lacks current approval evidence") {
		t.Fatalf("expected full Foundation blocker after assessment, got %v", err)
	}
	for _, artifact := range artifacts[1:] {
		if _, err := Approve(root, filepath.Join(root, filepath.FromSlash(artifact.Path)), "approved", "Product Owner", "Foundation reviewed"); err != nil {
			t.Fatalf("approve %s: %v", artifact.Type, err)
		}
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err != nil {
		t.Fatalf("workspace remained blocked after full Foundation approval: %v", err)
	}
}

func TestExistingProductWorkspaceRequiresBaselineAndStrategyApprovals(t *testing.T) {
	root := setupProduct(t)
	if err := writeJSON(filepath.Join(root, ".product", "framework.json"), map[string]any{"starting_point": "existing-product"}); err != nil {
		t.Fatal(err)
	}
	baseline := Artifact{ID: "PRODUCT-BASELINE-1", Type: "product-baseline", Status: "draft", Path: "foundation/product-baseline.md"}
	strategy := Artifact{ID: "STRATEGY-1", Type: "strategy", Status: "draft", Path: "foundation/strategy/strategy.md", ParentIDs: []string{baseline.ID}}
	registry, err := LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	registry.Artifacts = append([]Artifact{baseline, strategy}, registry.Artifacts...)
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	for _, artifact := range []Artifact{baseline, strategy} {
		path := filepath.Join(root, filepath.FromSlash(artifact.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("| ID | "+artifact.ID+" |\n| Type | "+artifact.Type+" |\n| Status | draft |\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if artifact.Type == "strategy" {
			if err := os.WriteFile(filepath.Join(filepath.Dir(path), "context.md"), []byte("status: draft\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "product baseline lacks current approval evidence") {
		t.Fatalf("expected baseline blocker, got %v", err)
	}
	if _, err := Approve(root, filepath.Join(root, filepath.FromSlash(baseline.Path)), "approved", "Product Owner", "current product reviewed"); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "strategy lacks current approval evidence") {
		t.Fatalf("expected Strategy blocker, got %v", err)
	}
	if _, err := Approve(root, filepath.Join(root, filepath.FromSlash(strategy.Path)), "approved", "Product Owner", "future direction approved"); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err != nil {
		t.Fatalf("workspace remained blocked after baseline and Strategy approvals: %v", err)
	}
}

func TestExistingDocumentsWorkspaceRequiresLatestRunMaterialization(t *testing.T) {
	root := setupProduct(t)
	manifest := map[string]any{"starting_point": "existing-documents", "import": map[string]any{"latest_run": "IMPORT-002"}}
	if err := writeJSON(filepath.Join(root, ".product", "framework.json"), manifest); err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(root, "knowledge", "imports", "runs", "IMPORT-002", "import-plan.json")
	if err := writeJSON(planPath, map[string]any{"status": "draft", "materialization_approved": false}); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "latest import run IMPORT-002 is not materially complete") {
		t.Fatalf("expected latest import blocker, got %v", err)
	}
	materialized := map[string]any{
		"status": "materialized", "materialization_approved": true,
		"materialization_approved_by": "Product Owner", "materialization_approved_at": "2026-07-12T12:00:00Z",
		"materialized_paths": []string{"domains/payments/domain.md"},
	}
	if err := writeJSON(planPath, materialized); err != nil {
		t.Fatal(err)
	}
	mappingPath := filepath.Join(filepath.Dir(planPath), "mapping.json")
	if err := writeJSON(mappingPath, map[string]any{"mappings": []any{map[string]any{"target": "domains/payments/domain.md", "selected": true, "draft_content": "# Payments\n"}}}); err != nil {
		t.Fatal(err)
	}
	targetPath := filepath.Join(root, "domains", "payments", "domain.md")
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte("tampered\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err == nil || !strings.Contains(err.Error(), "differs from the approved draft") {
		t.Fatalf("expected tampered import blocker, got %v", err)
	}
	if err := os.WriteFile(targetPath, []byte("# Payments\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateWorkspace(root, "FT-1", "", "", "", "tester"); err != nil {
		t.Fatalf("workspace remained blocked after latest import materialization: %v", err)
	}
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

func TestApproveSupportsModularArtifactTypeWithGenericAdapter(t *testing.T) {
	root := setupProduct(t)
	path := filepath.Join(root, "audits", "custom-review.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	text := "| ID | REVIEW-1 |\n| Type | custom-review |\n| Status | draft |\n"
	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	registry, err := LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	registry.Artifacts = append(registry.Artifacts, Artifact{ID: "REVIEW-1", Type: "custom-review", Status: "draft", Path: "audits/custom-review.md"})
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), registry); err != nil {
		t.Fatal(err)
	}
	if _, err := Approve(root, path, "approved", "owner", ""); err != nil {
		t.Fatal(err)
	}
	entries, err := filepath.Glob(filepath.Join(root, ".product", "history", "approval-review-1-approved-*.json"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("generic approval record missing: %v %v", err, entries)
	}
	updated, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(updated), "approved") {
		t.Fatalf("custom artifact was not approved: %v %s", err, updated)
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

func TestApproveBlocksNonConformantApprovedCandidate(t *testing.T) {
	root := t.TempDir()
	tasks := filepath.Join(root, "tasks")
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(tasks, 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(tasks, "TK-001.md")
	if err := os.WriteFile(path, []byte("# Task\n\n| Field | Value |\n| --- | --- |\n| Status | `draft` |\n\n## Objective\nDo the thing.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), Registry{Artifacts: []Artifact{{ID: "TK-001", Type: "task", Status: "draft", Path: "tasks/TK-001.md"}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := Approve(root, path, "approved", "Human", ""); err == nil || !strings.Contains(err.Error(), "template-conformance") {
		t.Fatalf("expected template conformance blocker, got %v", err)
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
		"context.md":                     "---\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering-system.md":          "| Field | Value |\n| --- | --- |\n| ID | `ENGSYS-001` |\n| Status | `draft` |\n| Version | `1.0.0` |\n",
		"engineering-system.yaml":        "schema_version: 1\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: product\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: baseline\n    evidence: []\n  quality:\n    contract: quality/quality-system.md\n    maturity: baseline\n    evidence: []\n",
		"architecture/modules.md":        "# Modules\n",
		"quality/quality-system.md":      "| Field | Value |\n| --- | --- |\n| Engineering System | `ENGSYS-001 @ 1.0.0` |\n| Status | `draft` |\n\n| Area | Policy | Evidence | Maturity |\n| --- | --- | --- | --- |\n| Behavioral | strategy | none | baseline |\n| Accessibility | strategy | none | baseline |\n| Security and privacy | strategy | none | baseline |\n| Performance and reliability | model | none | baseline |\n| Observability | model | none | baseline |\n",
		"quality/quality-system.yaml":    "schema_version: 1\nengineering_system: ENGSYS-001\nversion: 1.0.0\nstatus: draft\nareas:\n  behavioral: {maturity: baseline, policy: test-strategy.md}\n  accessibility: {maturity: baseline, policy: test-strategy.md}\n  security_privacy: {maturity: baseline, policy: test-strategy.md, delegated_gate: security-review}\n  performance_reliability: {maturity: baseline, policy: quality-model.md}\n  observability: {maturity: baseline, policy: quality-model.md}\ngate_source: knowledge/conventions/gates.md\nexceptions:\n  require_owner: true\n  require_residual_risk: true\n  require_expiry_or_review: true\n",
		"quality/fitness-functions.yaml": "version: 1\nfunctions: []\n",
		"quality/test-strategy.md":       "# Strategy\n",
		"quality/quality-model.md":       "# Model\n",
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
	qualityCatalogPath := filepath.Join(engineering, "quality", "quality-system.yaml")
	qualityCatalogData, _ := os.ReadFile(qualityCatalogPath)
	if err := os.Remove(qualityCatalogPath); err != nil {
		t.Fatal(err)
	}
	if _, err := PreviewApproval(root, canonical, "approved"); err == nil {
		t.Fatal("approval preview accepted a missing configured Quality System catalog")
	}
	if err := os.WriteFile(qualityCatalogPath, qualityCatalogData, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Approve(root, canonical, "approved", "Human", "Reviewed composite system"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"context.md", "engineering-system.md", "engineering-system.yaml", "quality/quality-system.md", "quality/quality-system.yaml"} {
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
