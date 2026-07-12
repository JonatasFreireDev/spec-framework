package adapters

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryAndPinnedArgv(t *testing.T) {
	argv, err := ProviderArgv("impeccable", "install", "2.3.2")
	if err != nil || len(argv) != 3 || argv[0] != "impeccable@2.3.2" || argv[1] != "skills" || argv[2] != "install" {
		t.Fatalf("unexpected argv: %v %v", argv, err)
	}
	if _, err := ProviderArgv("impeccable", "install", ""); err == nil {
		t.Fatal("expected explicit version requirement")
	}
}

func TestResolveExactVersionWithoutNetwork(t *testing.T) {
	resolved, err := ResolveVersion("impeccable", "2.3.2")
	if err != nil || resolved != "2.3.2" {
		t.Fatalf("resolved=%q err=%v", resolved, err)
	}
	if _, err := ResolveVersion("impeccable", "banana"); err == nil {
		t.Fatal("expected invalid version error")
	}
}

func TestInspectFindsProjectSkill(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".agents", "skills", "impeccable", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("---\nname: impeccable\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, err := Inspect(root, "impeccable")
	if err != nil || !status.Installed || len(status.Paths) != 1 {
		t.Fatalf("unexpected status: %+v %v", status, err)
	}
}
