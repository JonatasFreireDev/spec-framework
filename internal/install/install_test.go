package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitGeneratesSelectedAgentSkillTrees(t *testing.T) {
	target := filepath.Join(t.TempDir(), "product")
	result, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex, Cursor, Claude}})
	if err != nil {
		t.Fatal(err)
	}
	if result.SkillCount != 3 {
		t.Fatalf("skill target count=%d", result.SkillCount)
	}
	for _, file := range []string{".agents/skills/code-runner/SKILL.md", ".cursor/skills/code-runner/SKILL.md", ".claude/skills/code-runner/SKILL.md", ".spec-framework/manifest.json", "product/.product/framework.json"} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(file))); err != nil {
			t.Errorf("missing %s: %v", file, err)
		}
	}
	for _, file := range []string{"README.md", "BOOTSTRAP.md"} {
		data, err := os.ReadFile(filepath.Join(target, file))
		if err != nil {
			t.Fatalf("missing generated %s: %v", file, err)
		}
		if len(data) < 200 {
			t.Fatalf("generated %s is unexpectedly empty", file)
		}
	}
	workflow, err := os.ReadFile(filepath.Join(target, ".github", "workflows", "framework-validation.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(workflow), "releases/download/vdev") {
		t.Fatal("development workflow points to nonexistent vdev release")
	}
	if _, err := os.Stat(filepath.Join(target, ".claude", "skills", "threat-modeler", "agents", "openai.yaml")); !os.IsNotExist(err) {
		t.Fatal("Codex metadata leaked into Claude skills")
	}
}

func TestUpgradePreservesProductContent(t *testing.T) {
	target := filepath.Join(t.TempDir(), "product")
	_, err := Init(Options{Target: target, Agents: []Agent{Codex}})
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(target, "product", "foundation", "problem", "problem.md")
	if err = os.WriteFile(file, []byte("owned"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err = Upgrade(Options{Target: target, Agents: []Agent{Cursor}}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(file)
	if string(data) != "owned" {
		t.Fatal("product content changed")
	}
	readme := filepath.Join(target, "README.md")
	bootstrap := filepath.Join(target, "BOOTSTRAP.md")
	if err := os.WriteFile(readme, []byte("adopter readme"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bootstrap, []byte("adopter bootstrap"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err = Upgrade(Options{Target: target, Agents: []Agent{Cursor}}); err != nil {
		t.Fatal(err)
	}
	for file, want := range map[string]string{readme: "adopter readme", bootstrap: "adopter bootstrap"} {
		got, _ := os.ReadFile(file)
		if string(got) != want {
			t.Fatalf("upgrade overwrote %s", file)
		}
	}
}
