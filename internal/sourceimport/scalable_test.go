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

func TestScalableRunRejectsBinaryBeforeCopyAndHonorsDefaultExcludes(t *testing.T) {
	root, source := t.TempDir(), t.TempDir()
	dependency := filepath.Join(source, "node_modules")
	if err := os.MkdirAll(dependency, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dependency, "ignored.md"), []byte("ignored"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "bad.pdf"), []byte("binary"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateScalableRun(root, []string{source}, CreateOptions{MaxFiles: 5, MaxTotalBytes: 100, MaxFileBytes: 100, ChunkSize: 1, BinaryPolicy: "reject"}); err == nil {
		t.Fatal("binary accepted")
	}
	if _, err := os.Stat(filepath.Join(root, "knowledge", "imports", "sources", "IMPORT-001")); !os.IsNotExist(err) {
		t.Fatal("partial sources preserved after rejected binary")
	}
}

func TestScalableReviewRequiresEvidenceAndGuardsMaterialization(t *testing.T) {
	root, source := t.TempDir(), t.TempDir()
	file := filepath.Join(source, "a.md")
	if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	run, err := CreateScalableRun(root, []string{file}, CreateOptions{MaxFiles: 1, MaxTotalBytes: 100, MaxFileBytes: 100, ChunkSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := Resume(root, run, "CHUNK-0001", "agent")
	if err != nil {
		t.Fatal(err)
	}
	if err := RecordChunkReview(root, run, chunk.ID, "agent", ChunkReview{}); err == nil {
		t.Fatal("review without evidence accepted")
	}
	if err := RecordChunkReview(root, run, chunk.ID, "agent", ChunkReview{SourceEvidence: map[string][]Evidence{"SRC-000001": {{Locator: "line 1", Claim: "content"}}}}); err != nil {
		t.Fatal(err)
	}
	if err := requireReviewedChunks(root, run); err != nil {
		t.Fatal(err)
	}
}

func TestScalableResumeSkipsLockedChunk(t *testing.T) {
	root, source := t.TempDir(), t.TempDir()
	file := filepath.Join(source, "a.md")
	if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	run, err := CreateScalableRun(root, []string{file}, CreateOptions{MaxFiles: 1, MaxTotalBytes: 100, MaxFileBytes: 100, ChunkSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "knowledge", "imports", "runs", run, "chunks")
	unlock, err := lockChunk(dir, "CHUNK-0001")
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	if _, err := Resume(root, run, "CHUNK-0001", "agent"); err == nil {
		t.Fatal("locked chunk claimed")
	}
}
