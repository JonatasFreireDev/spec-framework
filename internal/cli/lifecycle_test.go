package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/cli"
)

func TestUninstallWithoutYesIsPreviewOnly(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.New("1.0.0").Run([]string{"uninstall"}, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	for _, expected := range []string{"CLI uninstall", "Product repositories: never touched", "Re-run with --yes"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("preview missing %q: %s", expected, stdout.String())
		}
	}
}

func TestLifecycleCommandsHaveNativeHelp(t *testing.T) {
	for _, name := range []string{"update", "uninstall"} {
		var stdout, stderr bytes.Buffer
		code := cli.New("1.0.0").Run([]string{name, "--help"}, &stdout, &stderr)
		if code != 0 || stderr.Len() != 0 || !strings.Contains(stdout.String(), "--yes") {
			t.Fatalf("%s help code=%d stdout=%q stderr=%q", name, code, stdout.String(), stderr.String())
		}
	}
}
