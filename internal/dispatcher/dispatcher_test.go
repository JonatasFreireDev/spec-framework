package dispatcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallMapsNativeQuestionToolPerHarness(t *testing.T) {
	home := t.TempDir()
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", home)
	for _, test := range []struct{ agent, root, tool string }{
		{"codex", ".agents", "request_user_input"},
		{"cursor", ".cursor", "Cursor's native user-question tool"},
		{"claude", ".claude", "AskUserQuestion"},
	} {
		path, err := Install(test.agent)
		if err != nil {
			t.Fatalf("%s: %v", test.agent, err)
		}
		wantPath := filepath.Join(home, test.root, "skills", "spec-framework", "SKILL.md")
		if path != wantPath {
			t.Fatalf("%s path=%s want=%s", test.agent, path, wantPath)
		}
		data, err := os.ReadFile(path)
		if err != nil || !strings.Contains(string(data), test.tool) || !strings.Contains(string(data), "native_user_question") || !strings.Contains(string(data), "AGENTS.framework.md") {
			t.Fatalf("%s dispatcher mapping missing: %v %s", test.agent, err, data)
		}
	}
}

func TestInstallMigratesLegacyCodexDispatcherWithoutTouchingOtherSkills(t *testing.T) {
	home := t.TempDir()
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", home)
	legacy := filepath.Join(home, ".codex", "skills", "spec-framework", "SKILL.md")
	other := filepath.Join(home, ".codex", "skills", "other-skill", "SKILL.md")
	for _, path := range []string{legacy, other} {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("existing"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	path, err := Install("codex")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(home, ".agents", "skills", "spec-framework", "SKILL.md"); path != want {
		t.Fatalf("path=%s want=%s", path, want)
	}
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatalf("legacy dispatcher remains: %v", err)
	}
	if data, err := os.ReadFile(other); err != nil || string(data) != "existing" {
		t.Fatalf("unrelated legacy skill changed: %q %v", data, err)
	}
}
