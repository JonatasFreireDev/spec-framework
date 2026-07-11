package validator

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportValidationDetectsChangedSourceAndDuplicateTargets(t *testing.T) {
	root := t.TempDir()
	run := filepath.Join(root, "knowledge", "imports", "runs", "IMPORT-001")
	sourceRel := "knowledge/imports/sources/epic.md"
	source := filepath.Join(root, filepath.FromSlash(sourceRel))
	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(run, 0755); err != nil {
		t.Fatal(err)
	}
	original := []byte("original")
	sum := sha256.Sum256(original)
	if err := os.WriteFile(source, []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	write := func(name string, value any) {
		data, _ := json.Marshal(value)
		if name == "mapping.json" {
			data = append([]byte{0xef, 0xbb, 0xbf}, data...)
		}
		if err := os.WriteFile(filepath.Join(run, name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("inventory.json", map[string]any{"schema_version": 1, "import_id": "IMPORT-001", "sources": []any{map[string]any{"path": sourceRel, "sha256": fmt.Sprintf("%x", sum[:])}}})
	write("import-plan.json", map[string]any{"materialization_approved": false})
	write("mapping.json", map[string]any{"mappings": []any{map[string]any{"id": "MAP-1", "selected": true, "target": "domains/a/domain.md", "source_documents": []any{sourceRel}}, map[string]any{"id": "MAP-2", "selected": true, "target": "domains/a/domain.md", "source_documents": []any{sourceRel}}}})
	for _, name := range []string{"conflicts.md", "import-report.md"} {
		if err := os.WriteFile(filepath.Join(run, name), []byte("# Report"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	changed, duplicate := false, false
	for _, d := range result.Diagnostics {
		if d.Check == "imports" && strings.Contains(d.Message, "Source changed") {
			changed = true
		}
		if d.Check == "imports" && strings.Contains(d.Message, "Multiple selected mappings") {
			duplicate = true
		}
	}
	if !changed || !duplicate {
		t.Fatalf("changed=%v duplicate=%v diagnostics=%+v", changed, duplicate, result.Diagnostics)
	}
}

func TestDeliveryClosureRejectsLegacyAndUnknownHandoffs(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "framework", "skills", "known"), 0755)
	s := Snapshot{Root: root, FrameworkRoot: root, Text: map[string]string{"framework/skills/a/SKILL.md": "## Handoff\nNext: 05-old.md\n", "framework/skills/b/SKILL.md": "## Handoff\nNext: missing-skill.\n"}}
	d := validateDeliveryClosure(s)
	legacy, unknown := false, false
	for _, x := range d {
		if strings.Contains(x.Message, "Legacy numbered") {
			legacy = true
		}
		if strings.Contains(x.Message, "Unknown next skill") {
			unknown = true
		}
	}
	if !legacy || !unknown {
		t.Fatalf("legacy=%v unknown=%v diagnostics=%+v", legacy, unknown, d)
	}
}

func TestValidatedTaskRequiresMatchingDiffHashes(t *testing.T) {
	text := "# Task\n\n| Field | Value |\n| --- | --- |\n| Status | validated |\n| Branch | feature/x |\n| Base commit | abcdef1 |\n| Diff hash | hash-a |\n| Changed paths | src/x.go |\n| Test status | passed |\n| Commits | abcdef2 |\n| Code paths | src/x.go |\n| Code Review diff hash | hash-a |\n| QA diff hash | hash-b |\n"
	s := Snapshot{Text: map[string]string{"domains/x/use-cases/u/tasks/TK-1.md": text}}
	d := validateDeliveryClosure(s)
	found := false
	for _, x := range d {
		if x.Check == "diff-staleness" {
			found = true
		}
	}
	if !found {
		t.Fatalf("diagnostics=%+v", d)
	}
}

func TestDiagnosticsAreDeterministic(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "domains", "a"), 0755)
	_ = os.WriteFile(filepath.Join(root, "domains", "a", "context.md"), []byte("status: draft\n"), 0644)
	first, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	second, _ := Validate(context.Background(), root, root)
	if len(first.Diagnostics) == 0 || len(first.Diagnostics) != len(second.Diagnostics) {
		t.Fatalf("diagnostics=%v", first.Diagnostics)
	}
	for i := range first.Diagnostics {
		if first.Diagnostics[i] != second.Diagnostics[i] {
			t.Fatal("unstable diagnostics")
		}
	}
}

func TestBlocksBrokenMarkdownLink(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "index.md"), []byte("[Missing](missing.md)\n"), 0644)
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Check == "links" {
			found = true
		}
	}
	if !found {
		t.Fatalf("%+v", result)
	}
}

func TestRequiresMatchingApprovalRecord(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	artifact := "# Artifact\n"
	_ = os.WriteFile(filepath.Join(root, "artifact.md"), []byte(artifact), 0644)
	registry := map[string]any{"artifacts": []any{map[string]any{"id": "ART-1", "status": "approved", "path": "artifact.md"}}}
	data, _ := json.Marshal(registry)
	_ = os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0644)
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range result.Diagnostics {
		if d.Check == "approval-records" {
			found = true
		}
	}
	if !found {
		t.Fatalf("%+v", result)
	}
}
