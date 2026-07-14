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
		{"codex", ".codex", "request_user_input"},
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
