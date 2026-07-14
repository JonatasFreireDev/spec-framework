package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdaptersListAndInstallPreview(t *testing.T) {
	var out, errout bytes.Buffer
	code := New("test").Run([]string{"adapters", "list"}, &out, &errout)
	if code != 0 || !strings.Contains(out.String(), "impeccable") {
		t.Fatalf("list=%d out=%s err=%s", code, out.String(), errout.String())
	}
	out.Reset()
	errout.Reset()
	code = New("test").Run([]string{"adapters", "install", "impeccable", "--version", "2.3.2"}, &out, &errout)
	if code != 0 || !strings.Contains(out.String(), "impeccable@2.3.2 skills install") || !strings.Contains(out.String(), "Re-run with --yes") {
		t.Fatalf("preview=%d out=%s err=%s", code, out.String(), errout.String())
	}
}

func TestInitRequiresPinnedImpeccableVersion(t *testing.T) {
	var out, errout bytes.Buffer
	code := New("test").Run([]string{"init", "--target", t.TempDir(), "--agents", "codex", "--install-impeccable", "--yes"}, &out, &errout)
	if code != 2 || !strings.Contains(errout.String(), "--impeccable-version") {
		t.Fatalf("init=%d out=%s err=%s", code, out.String(), errout.String())
	}
}

func TestAdaptersDiscoverRepositoryRootFromProductDirectory(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "product", ".product"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "product", ".product", "framework.json"), []byte(`{"framework":"spec-framework"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(repo, "product", "domains", "events")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	var out, errout bytes.Buffer
	code := New("test").Run([]string{"adapters", "install", "impeccable", "--version", "2.3.2"}, &out, &errout)
	if code != 0 {
		t.Fatalf("code=%d out=%s err=%s", code, out.String(), errout.String())
	}
	if !strings.Contains(out.String(), "Working directory: "+repo) {
		t.Fatalf("adapter root not discovered: %s", out.String())
	}
}
