package validator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteReportUsesAtomicCanonicalPaths(t *testing.T) {
	root := t.TempDir()
	paths, err := WriteReport(root, Result{})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Fatalf("paths=%v", paths)
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Fatal(err)
		}
	}
}
func TestWriteRegistrySortsArtifacts(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "docs"), 0755)
	_ = os.WriteFile(filepath.Join(root, "docs", "b.md"), []byte("| ID | B-1 |\n| Type | test |\n| Status | draft |\n"), 0644)
	_ = os.WriteFile(filepath.Join(root, "docs", "a.md"), []byte("| ID | A-1 |\n| Type | test |\n| Status | draft |\n"), 0644)
	path, err := WriteRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	var value struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	if json.Unmarshal(data, &value) != nil || len(value.Artifacts) != 2 {
		t.Fatal(string(data))
	}
	if value.Artifacts[0]["id"] != "A-1" {
		t.Fatalf("%v", value.Artifacts)
	}
}
