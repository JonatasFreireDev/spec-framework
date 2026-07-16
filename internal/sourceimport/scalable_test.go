package sourceimport

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScalableRunFiltersBudgetsPagesAndResumes(t *testing.T) {
	root, source := t.TempDir(), t.TempDir()
	for _, name := range []string{"a.md", "b.md", "c.md", "skip.bin"} {
		if err := os.WriteFile(filepath.Join(source, name), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	run, err := CreateScalableRun(root, []string{source}, CreateOptions{Include: []string{"*.md", "**/*.md"}, MaxFiles: 3, MaxTotalBytes: 100, MaxFileBytes: 20, ChunkSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	status, err := ImportStatus(root, run)
	if err != nil {
		t.Fatal(err)
	}
	if status.Sources != 3 || status.Chunks != 2 || status.Queued != 2 {
		t.Fatalf("status=%+v", status)
	}
	chunk, err := Resume(root, run, "", "importer")
	if err != nil || chunk.Status != "reviewing" {
		t.Fatalf("chunk=%+v err=%v", chunk, err)
	}
	if _, err := Resume(root, run, chunk.ID, "other"); err == nil {
		t.Fatal("active chunk was claimed twice")
	}
}

func TestScalableRunRejectsBudgetBeforeCopy(t *testing.T) {
	root, source := t.TempDir(), t.TempDir()
	file := filepath.Join(source, "large.md")
	if err := os.WriteFile(file, []byte("too large"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateScalableRun(root, []string{file}, CreateOptions{MaxFiles: 1, MaxTotalBytes: 1, MaxFileBytes: 20, ChunkSize: 1}); err == nil {
		t.Fatal("budget overflow accepted")
	}
	entries, _ := os.ReadDir(filepath.Join(root, "knowledge", "imports", "sources"))
	if len(entries) != 0 {
		t.Fatal("source copied after failed budget")
	}
}
