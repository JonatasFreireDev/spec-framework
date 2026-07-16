package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRejectRequiresRationaleAndReturnsToDraft(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	artifact := Artifact{ID: "ART-1", Type: "feature", Status: "draft", Path: "artifact.md"}
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), Registry{Artifacts: []Artifact{artifact}}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, artifact.Path)
	if err := os.WriteFile(path, []byte("status: draft\n# Artifact\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Approve(root, path, "rejected", "Product Owner", ""); err == nil || !strings.Contains(err.Error(), "requires notes") {
		t.Fatalf("empty rejection err=%v", err)
	}
	record, err := Approve(root, path, "rejected", "Product Owner", "Clarify the failure path before approval.")
	if err != nil {
		t.Fatal(err)
	}
	if record.StatusGranted != "rejected" || record.Notes == "" {
		t.Fatalf("record=%+v", record)
	}
	registry, err := LoadRegistry(root)
	if err != nil || registry.Artifacts[0].Status != "rejected" {
		t.Fatalf("registry=%+v err=%v", registry, err)
	}
	if _, err := Approve(root, path, "draft", "Product Owner", "Revisions started."); err != nil {
		t.Fatal(err)
	}
}

func TestRejectCanReopenApprovedArtifact(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	artifact := Artifact{ID: "ART-1", Type: "feature", Status: "approved", Path: "artifact.md"}
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), Registry{Artifacts: []Artifact{artifact}}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, artifact.Path)
	if err := os.WriteFile(path, []byte("status: approved\n# Artifact\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Approve(root, path, "rejected", "Product Owner", "The approved version needs a correction."); err != nil {
		t.Fatal(err)
	}
	registry, err := LoadRegistry(root)
	if err != nil || registry.Artifacts[0].Status != "rejected" {
		t.Fatalf("registry=%+v err=%v", registry, err)
	}
}

func TestRejectedArtifactCanBeApprovedAgain(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	artifact := Artifact{ID: "ART-1", Type: "feature", Status: "rejected", Path: "artifact.md"}
	if err := writeJSON(filepath.Join(root, ".product", "artifacts.json"), Registry{Artifacts: []Artifact{artifact}}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, artifact.Path)
	if err := os.WriteFile(path, []byte("status: rejected\n# Artifact\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Approve(root, path, "approved", "Product Owner", "Revised after the rejection."); err != nil {
		t.Fatal(err)
	}
	registry, err := LoadRegistry(root)
	if err != nil || registry.Artifacts[0].Status != "approved" {
		t.Fatalf("registry=%+v err=%v", registry, err)
	}
}
