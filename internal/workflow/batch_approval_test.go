package workflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildAndApplyFoundationBatchApproval(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := Registry{Artifacts: []Artifact{
		{ID: "PRB-001", Type: "problem", Status: "draft", Path: "foundation/problem/problem.md"},
		{ID: "VIS-001", Type: "vision", Status: "draft", Path: "foundation/vision/vision.md", ParentIDs: []string{"PRB-001"}},
	}}
	data, _ := json.Marshal(registry)
	if err := os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"foundation/problem/problem.md", "foundation/vision/vision.md"} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("status: draft\n# Artifact\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	plan, err := BuildBatchApprovalPlan(root, BatchScope{Foundation: true}, "approved")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.ToApprove) != 2 || len(plan.Blockers) != 0 {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	records, err := ApproveBatch(root, plan, "Product Owner", "Foundation reviewed")
	if err != nil || len(records) != 2 {
		t.Fatalf("records=%d err=%v", len(records), err)
	}
	for _, path := range []string{"foundation/problem/problem.md", "foundation/vision/vision.md"} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil || !strings.Contains(string(data), "status: approved") {
			t.Fatalf("artifact %s was not approved: %v\n%s", path, err, data)
		}
	}
}

func TestBuildBatchBlocksUnapprovedParentOutsideScope(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := Registry{Artifacts: []Artifact{
		{ID: "PRB-001", Type: "problem", Status: "draft", Path: "foundation/problem/problem.md"},
		{ID: "VIS-001", Type: "vision", Status: "draft", Path: "foundation/vision/vision.md", ParentIDs: []string{"PRB-001"}},
	}}
	data, _ := json.Marshal(registry)
	_ = os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0o644)
	for _, path := range []string{"foundation/problem/problem.md", "foundation/vision/vision.md"} {
		full := filepath.Join(root, filepath.FromSlash(path))
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		_ = os.WriteFile(full, []byte("status: draft\n"), 0o644)
	}
	plan, err := BuildBatchApprovalPlan(root, BatchScope{IDs: []string{"VIS-001"}}, "approved")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.ToApprove) != 0 || len(plan.Blockers) != 1 || !strings.Contains(plan.Blockers[0].Reason, "parent PRB-001") {
		t.Fatalf("expected parent blocker: %+v", plan)
	}
}

func TestBuildBatchBlocksStaleDerivedArtifact(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := Registry{Artifacts: []Artifact{{ID: "DES-001", Type: "design", Status: "draft", Path: "design.md"}}}
	data, _ := json.Marshal(registry)
	_ = os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0o644)
	_ = os.WriteFile(filepath.Join(root, "design.md"), []byte("status: draft\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "source.md"), []byte("changed\n"), 0o644)
	derivations := map[string]any{"derivations": []any{map[string]any{
		"artifact_id": "DES-001", "path": "design.md",
		"derived_from": []any{map[string]any{"path": "source.md", "content_hash": "old-hash"}},
	}}}
	data, _ = json.Marshal(derivations)
	_ = os.WriteFile(filepath.Join(root, ".product", "derivations.json"), data, 0o644)
	plan, err := BuildBatchApprovalPlan(root, BatchScope{IDs: []string{"DES-001"}}, "approved")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.ToApprove) != 0 || len(plan.Blockers) != 1 || !strings.Contains(plan.Blockers[0].Reason, "stale") {
		t.Fatalf("expected stale blocker: %+v", plan)
	}
}
