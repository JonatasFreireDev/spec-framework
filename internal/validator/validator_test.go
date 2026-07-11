package validator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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
