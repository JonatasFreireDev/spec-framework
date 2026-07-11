package moveartifact

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanAndApplyMove(t *testing.T) {
	root := t.TempDir()
	write := func(name, content string) {
		file := filepath.Join(root, filepath.FromSlash(name))
		_ = os.MkdirAll(filepath.Dir(file), 0755)
		_ = os.WriteFile(file, []byte(content), 0644)
	}
	write("domains/old/file.md", "# Target\n")
	write("docs/link.md", "[Target](../domains/old/file.md#target)\n")
	write("docs/note.md", "Review domains/old manually.\n")
	write("registry.json", `{"path":"domains/old/file.md"}`)
	p, err := Build(root, "domains/old", "domains/new")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Rewrites) != 2 || len(p.Mentions) != 1 {
		t.Fatalf("rewrites=%v mentions=%v", p.Rewrites, p.Mentions)
	}
	if _, err = os.Stat(filepath.Join(root, "domains/old")); err != nil {
		t.Fatal("plan mutated source")
	}
	if err = Apply(p); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(root, "domains/new/file.md")); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(root, "docs/link.md"))
	if !strings.Contains(string(data), "../domains/new/file.md#target") {
		t.Fatal(string(data))
	}
}

func TestBuildRejectsEscape(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "source"), []byte("x"), 0644)
	if _, err := Build(root, "source", "../outside"); err == nil {
		t.Fatal("expected confinement error")
	}
}
