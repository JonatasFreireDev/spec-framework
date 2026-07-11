package fsx

import (
	"path/filepath"
	"testing"
)

func TestInsideConfinesPathsToRoot(t *testing.T) {
	root := t.TempDir()
	if !Inside(root, filepath.Join(root, "nested", "file")) {
		t.Fatal("nested path should be inside")
	}
	if Inside(root, filepath.Join(root, "..", "escaped")) {
		t.Fatal("escaped path should be outside")
	}
}
