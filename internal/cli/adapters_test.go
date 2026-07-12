package cli

import (
	"bytes"
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
