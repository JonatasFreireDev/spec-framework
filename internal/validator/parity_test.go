package validator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type fixture struct {
	root      string
	artifacts []map[string]any
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	root := t.TempDir()
	writeFixture(t, root, ".product/ids.json", map[string]any{"policy": "slug-scoped", "deprecated_counters": true})
	writeFixture(t, root, ".product/decisions.json", map[string]any{"decisions": []any{}})
	return &fixture{root: root}
}
func writeFixture(t *testing.T, root, path string, value any) {
	t.Helper()
	name := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		t.Fatal(err)
	}
	var data []byte
	if text, ok := value.(string); ok {
		data = []byte(text)
	} else {
		data, _ = json.MarshalIndent(value, "", "  ")
	}
	if err := os.WriteFile(name, data, 0644); err != nil {
		t.Fatal(err)
	}
}
func addArtifact(t *testing.T, f *fixture, id, kind, status, path, content string, extra map[string]any) {
	t.Helper()
	writeFixture(t, f.root, path, content)
	a := map[string]any{"id": id, "type": kind, "status": status, "path": path, "parentIds": []any{}}
	for k, v := range extra {
		a[k] = v
	}
	f.artifacts = append(f.artifacts, a)
}
func (f *fixture) validate(t *testing.T) Result {
	t.Helper()
	writeFixture(t, f.root, ".product/artifacts.json", map[string]any{"artifacts": f.artifacts})
	frameworkRoot := filepath.Clean(filepath.Join("..", ".."))
	result, err := Validate(context.Background(), f.root, frameworkRoot)
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func hasCheck(result Result, check string) bool {
	for _, d := range result.Diagnostics {
		if d.Check == check {
			return true
		}
	}
	return false
}

func TestParityMissingApprovalRecord(t *testing.T) {
	f := newFixture(t)
	addArtifact(t, f, "ART-1", "specification", "approved", "spec.md", "# Spec\n", nil)
	if result := f.validate(t); !hasCheck(result, "approval-records") {
		t.Fatalf("%+v", result)
	}
}
func TestParityStaleDerivation(t *testing.T) {
	f := newFixture(t)
	addArtifact(t, f, "SRC", "specification", "draft", "source.md", "changed\n", nil)
	addArtifact(t, f, "DER", "design", "proposed", "derived.md", "derived\n", nil)
	writeFixture(t, f.root, ".product/derivations.json", map[string]any{"derivations": []any{map[string]any{"artifact_id": "DER", "path": "derived.md", "derived_from": []any{map[string]any{"artifact_id": "SRC", "path": "source.md", "content_hash": "stale"}}}}})
	if result := f.validate(t); !hasCheck(result, "staleness") {
		t.Fatalf("%+v", result)
	}
}
func TestParityMissingSkillReference(t *testing.T) {
	f := newFixture(t)
	content := "| Field | Value |\n| --- | --- |\n| ID | TK-1 |\n| Type | task |\n| Status | draft |\n| Next skill | missing-skill |\n"
	addArtifact(t, f, "TK-1", "task", "draft", "task.md", content, map[string]any{"ownerSkill": "code-runner"})
	if result := f.validate(t); !hasCheck(result, "skill-reference") {
		t.Fatalf("%+v", result)
	}
}
func TestParityQAPlaceholderAndFindingRouting(t *testing.T) {
	f := newFixture(t)
	qa := "# QA\n\n| Field | Value |\n| --- | --- |\n| Test command | pending |\n| Gate logs | pending |\n| Environment | CI |\n| Limitations | none |\n| Verdict | passed |\n\n## Findings\n\n| Severity | Finding | Evidence | Required Fix |\n| --- | --- | --- | --- |\n| blocker | Broken | log | Fix |\n"
	addArtifact(t, f, "QA-1", "qa-evidence", "approved", "qa.md", qa, nil)
	result := f.validate(t)
	if !hasCheck(result, "qa-evidence") || !hasCheck(result, "failure-routing") {
		t.Fatalf("%+v", result)
	}
}
func TestParityTaskCommitAndPRReferences(t *testing.T) {
	f := newFixture(t)
	task := "# Task\n\n| Field | Value |\n| --- | --- |\n| Branch | feature/x |\n| Commits | not-a-hash |\n| PR | pending review |\n| Code paths | src/app.go |\n| Test status | passed |\n| Gate logs | test.log |\n| QA evidence | qa.md |\n"
	addArtifact(t, f, "TK-1", "task", "validated", "task.md", task, nil)
	result := f.validate(t)
	if !hasCheck(result, "code-evidence") {
		t.Fatalf("%+v", result)
	}
}
