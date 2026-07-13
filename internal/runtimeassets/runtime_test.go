package runtimeassets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverRequiresCanonicalManifest(t *testing.T) {
	root := t.TempDir()
	if _, _, err := Discover(root); err == nil {
		t.Fatal("mention-free directory activated without a manifest")
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("Spec Framework"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Discover(root); err == nil {
		t.Fatal("text mention activated the framework")
	}
	manifest := filepath.Join(root, "product", ".product", "framework.json")
	if err := os.MkdirAll(filepath.Dir(manifest), 0755); err != nil {
		t.Fatal(err)
	}
	data := []byte(`{"schema_version":3,"framework":"spec-framework","version":"1.2.3","activation":{"mode":"manifest-only"}}`)
	if err := os.WriteFile(manifest, data, 0644); err != nil {
		t.Fatal(err)
	}
	got, value, err := Discover(filepath.Join(root, "product"))
	if err != nil || got != root || value.Version != "1.2.3" {
		t.Fatalf("got root=%q version=%q err=%v", got, value.Version, err)
	}
}

func TestEnsureMaterializesVersionedAssets(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	root, err := Ensure("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"FRAMEWORK.md", "init/schema.json", "init/catalog.json", "init/contracts/new-product.json", "skills/code-runner/SKILL.md", "skills/discovery-and-challenge.md", "templates/specification-template.md", "examples/events/domains/events/domain.md", ".complete"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(name))); err != nil {
			t.Fatal(err)
		}
	}
}
