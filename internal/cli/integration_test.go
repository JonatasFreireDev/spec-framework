package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/cli"
)

func TestGoCLIInitValidateUpgradeAndMove(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "product")
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex,cursor,claude", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stderr=%s", code, stderr.String())
	}
	for _, path := range []string{".agents/skills/code-runner/SKILL.md", ".cursor/skills/code-runner/SKILL.md", ".claude/skills/code-runner/SKILL.md", ".spec-framework/manifest.json"} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); err != nil {
			t.Fatal(err)
		}
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chdir(target); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate"}, &stdout, &stderr); code != 0 {
		t.Fatalf("validate=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	moveSource := filepath.Join(target, "product", "move-source")
	if err = os.MkdirAll(moveSource, 0755); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(moveSource, "context.md"), []byte("# Move\n"), 0644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"move", "--from", "product/move-source", "--to", "product/moved"}, &stdout, &stderr); code != 0 {
		t.Fatalf("move=%d stderr=%s", code, stderr.String())
	}
	if _, err = os.Stat(filepath.Join(target, "product", "moved", "context.md")); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"upgrade", "--target", target, "--agents", "codex,cursor,claude", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("upgrade=%d stderr=%s", code, stderr.String())
	}
}
