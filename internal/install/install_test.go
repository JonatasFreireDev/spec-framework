package install

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitKeepsRuntimeAndDispatchersOutsideRepository(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	agentHome := filepath.Join(t.TempDir(), "agents")
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", agentHome)
	target := filepath.Join(t.TempDir(), "product")
	result, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex, Cursor, Claude}})
	if err != nil {
		t.Fatal(err)
	}
	if result.SkillCount != 0 {
		t.Fatalf("repository-local skill target count=%d", result.SkillCount)
	}
	for _, file := range []string{"product/.product/framework.json", "product/BOOTSTRAP.md"} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(file))); err != nil {
			t.Errorf("missing %s: %v", file, err)
		}
	}
	for _, file := range []string{".spec-framework", ".agents", ".cursor", ".claude", ".github"} {
		if _, err := os.Stat(filepath.Join(target, file)); !os.IsNotExist(err) {
			t.Fatalf("repository was polluted with %s", file)
		}
	}
	for _, path := range []string{
		filepath.Join(agentHome, ".codex", "skills", "spec-framework", "SKILL.md"),
		filepath.Join(agentHome, ".cursor", "skills", "spec-framework", "SKILL.md"),
		filepath.Join(agentHome, ".claude", "skills", "spec-framework", "SKILL.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInitFromExistingDocumentsCreatesImportRun(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	source := filepath.Join(t.TempDir(), "epic.md")
	if err := os.WriteFile(source, []byte("# Epic\n\nPayment and event features."), 0644); err != nil {
		t.Fatal(err)
	}
	result, err := Init(Options{Target: target, Version: "v0.3.0", Agents: []Agent{Codex}, StartingPoint: "existing-documents", Sources: []string{source}})
	if err != nil {
		t.Fatal(err)
	}
	if result.ImportID != "IMPORT-001" {
		t.Fatalf("import=%q", result.ImportID)
	}
	for _, name := range []string{"inventory.json", "import-plan.json", "mapping.json", "conflicts.md", "import-report.md"} {
		if _, err := os.Stat(filepath.Join(target, "product", "knowledge", "imports", "runs", result.ImportID, name)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseStartingPoint(t *testing.T) {
	if got, err := ParseStartingPoint(""); err != nil || got != "new-product" {
		t.Fatalf("got=%q err=%v", got, err)
	}
	if _, err := ParseStartingPoint("unknown"); err == nil {
		t.Fatal("expected unsupported starting point")
	}
}

func TestInstalledAgentsReadsManifestSelection(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product")
	if _, err := Init(Options{Target: target, Agents: []Agent{Codex, Cursor, Claude}}); err != nil {
		t.Fatal(err)
	}
	agents, err := InstalledAgents(target)
	if err != nil {
		t.Fatal(err)
	}
	if len(agents) != 3 {
		t.Fatalf("agents=%v", agents)
	}
}

func TestUpgradePreservesProductContent(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
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
	bootstrap := filepath.Join(target, "client-owned.txt")
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
