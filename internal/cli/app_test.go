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
	for _, command := range []string{"init", "validate", "move", "update", "uninstall", "upgrade", "version"} {
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

func TestCobraCommandTreeKeepsStableTopLevelCommands(t *testing.T) {
	root := cli.New("test").NewCommand(&bytes.Buffer{}, &bytes.Buffer{})
	for _, name := range []string{"init", "validate", "graph", "runtime", "update", "uninstall", "upgrade", "version"} {
		command, _, err := root.Find([]string{name})
		if err != nil || command == root || command.Name() != name {
			t.Errorf("Cobra command %q was not registered: command=%v err=%v", name, command, err)
		}
	}
}

func TestCobraLeafHelpDoesNotInvokeLegacyCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := cli.New("test").Run([]string{"graph", "--help"}, &stdout, &stderr)
	if exitCode != 0 || !strings.Contains(stdout.String(), "Inspect and operate execution graphs.") || strings.Contains(stdout.String(), "unknown graph command") || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", exitCode, stdout.String(), stderr.String())
	}
}

func TestLegacySubcommandHelpIncludesFlagsAndSucceeds(t *testing.T) {
	for _, args := range [][]string{
		{"init", "--help"},
		{"task", "readiness", "--help"},
		{"graph", "ready", "--help"},
	} {
		var stdout, stderr bytes.Buffer
		if code := cli.New("test").Run(args, &stdout, &stderr); code != 0 {
			t.Fatalf("args=%v exit=%d stdout=%q stderr=%q", args, code, stdout.String(), stderr.String())
		}
		if !strings.Contains(stdout.String(), "Usage") || stderr.Len() != 0 {
			t.Fatalf("args=%v stdout=%q stderr=%q", args, stdout.String(), stderr.String())
		}
		if args[0] == "init" && !strings.Contains(stdout.String(), "-starting-point string") {
			t.Fatalf("init help missing starting-point flag: %q", stdout.String())
		}
		if args[0] == "task" && !strings.Contains(stdout.String(), "-graph string") {
			t.Fatalf("task help missing graph flag: %q", stdout.String())
		}
	}
}
