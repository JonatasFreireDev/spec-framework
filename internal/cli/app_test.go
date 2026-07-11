package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/cli"
)

func TestHelpListsStableCommands(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := cli.New("test-version").Run([]string{"help"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%q", exitCode, stderr.String())
	}
	for _, command := range []string{"init", "validate", "move", "upgrade", "version"} {
		if !strings.Contains(stdout.String(), "  "+command) {
			t.Errorf("help does not list %q:\n%s", command, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestVersionWritesConfiguredVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := cli.New("1.2.3-test").Run([]string{"version"}, &stdout, &stderr)

	if exitCode != 0 || stdout.String() != "spec-framework 1.2.3-test\n" || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", exitCode, stdout.String(), stderr.String())
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := cli.New("test").Run([]string{"unknown"}, &stdout, &stderr)

	if exitCode != 2 || !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("exit=%d stderr=%q", exitCode, stderr.String())
	}
}
