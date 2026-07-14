package install

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitShipsLeanReadmeSurfaceAndUpgradePreservesAdopterReadmes(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	product := filepath.Join(target, "product")
	readmes := 0
	if err := filepath.WalkDir(product, func(path string, entry os.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && entry.Name() == "README.md" {
			readmes++
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if readmes != 16 {
		t.Fatalf("initialized product has %d READMEs, want 16", readmes)
	}
	obsolete := filepath.Join(product, "knowledge", "prompts", "README.md")
	if _, err := os.Stat(obsolete); !os.IsNotExist(err) {
		t.Fatalf("leaf placeholder README was materialized: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(obsolete), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(obsolete, []byte("adopter-owned guidance"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(obsolete)
	if err != nil || string(data) != "adopter-owned guidance" {
		t.Fatalf("upgrade changed adopter README: %q %v", data, err)
	}
}

func TestInitExistingFeatureActivatesOnlyFeatureBriefFoundation(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}, StartingPoint: "existing-feature"}); err != nil {
		t.Fatal(err)
	}
	brief := filepath.Join(target, "product", "foundation", "feature-brief.md")
	if data, err := os.ReadFile(brief); err != nil || !strings.Contains(string(data), "FEATURE-BRIEF-TBD") {
		t.Fatalf("feature brief missing or invalid: %v", err)
	}
	contextData, err := os.ReadFile(filepath.Join(target, "product", "context.md"))
	if err != nil || !strings.Contains(string(contextData), "complete and approve `foundation/feature-brief.md`") || strings.Contains(string(contextData), "Do not create domains or features until") {
		t.Fatalf("product context retained contradictory bootstrap guidance: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(target, "product", ".product", "artifacts.json"))
	if err != nil {
		t.Fatal(err)
	}
	var registry struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatal(err)
	}
	for _, artifact := range registry.Artifacts {
		kind, _ := artifact["type"].(string)
		if map[string]bool{"problem": true, "vision": true, "product-principles": true, "north-star": true, "strategy": true}[kind] {
			t.Fatalf("full Foundation artifact remained active: %v", artifact)
		}
	}
	serialized := string(data)
	if !strings.Contains(serialized, `"id": "FEATURE-BRIEF-TBD"`) || !strings.Contains(serialized, `"targetFeature": "FT-TEMPLATE"`) || !strings.Contains(serialized, `"parentIds": [`) || !strings.Contains(serialized, `"FEATURE-BRIEF-TBD"`) {
		t.Fatalf("feature brief registry linkage missing: %s", serialized)
	}
}

func TestInitExistingImplementationRegistersAssessmentBeforeProblem(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}, StartingPoint: "existing-implementation"}); err != nil {
		t.Fatal(err)
	}
	assessmentPath := filepath.Join(target, "product", "knowledge", "assessments", "implementation-assessment.md")
	if data, err := os.ReadFile(assessmentPath); err != nil || !strings.Contains(string(data), "IMPL-ASSESS-TBD") {
		t.Fatalf("implementation assessment missing or invalid: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(target, "product", ".product", "artifacts.json"))
	if err != nil {
		t.Fatal(err)
	}
	var registry struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatal(err)
	}
	serialized := string(data)
	if !strings.Contains(serialized, `"id": "IMPL-ASSESS-TBD"`) || !strings.Contains(serialized, `"type": "implementation-assessment"`) {
		t.Fatalf("assessment registry entry missing: %s", serialized)
	}
	foundProblemLink := false
	for _, artifact := range registry.Artifacts {
		kind, _ := artifact["type"].(string)
		if kind == "problem" && strings.Contains(string(mustJSON(t, artifact["parentIds"])), "IMPL-ASSESS-TBD") {
			foundProblemLink = true
		}
	}
	if !foundProblemLink {
		t.Fatal("Problem is not linked to Implementation Assessment")
	}
}

func TestInitExistingProductRegistersBaselineAndStrategy(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}, StartingPoint: "existing-product"}); err != nil {
		t.Fatal(err)
	}
	baselinePath := filepath.Join(target, "product", "foundation", "product-baseline.md")
	if data, err := os.ReadFile(baselinePath); err != nil || !strings.Contains(string(data), "PRODUCT-BASELINE-TBD") {
		t.Fatalf("product baseline missing or invalid: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(target, "product", ".product", "artifacts.json"))
	if err != nil {
		t.Fatal(err)
	}
	var registry struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatal(err)
	}
	foundStrategy := false
	for _, artifact := range registry.Artifacts {
		kind, _ := artifact["type"].(string)
		if map[string]bool{"problem": true, "vision": true, "product-principles": true, "north-star": true}[kind] {
			t.Fatalf("consolidated Foundation artifact remained active: %v", artifact)
		}
		if kind == "strategy" {
			foundStrategy = strings.Contains(string(mustJSON(t, artifact["parentIds"])), "PRODUCT-BASELINE-TBD")
		}
	}
	if !foundStrategy || !strings.Contains(string(data), `"type": "product-baseline"`) {
		t.Fatalf("baseline to Strategy registry contract missing: %s", data)
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

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
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		contract := string(data)
		if !strings.Contains(contract, "Resolve framework-guide first unless") || !strings.Contains(contract, "concrete artifact or workspace scope") || !strings.Contains(contract, "is not direct-route evidence by itself") {
			t.Fatalf("dispatcher is not Guide-first: %s", path)
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
	bootstrap, err := os.ReadFile(filepath.Join(target, "product", "BOOTSTRAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(bootstrap), "--run IMPORT-001") || strings.Contains(string(bootstrap), "<IMPORT-NNN>") || strings.Contains(string(bootstrap), "<latest-run>") {
		t.Fatalf("bootstrap did not pin active import run: %s", bootstrap)
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

func TestUpgradeRefreshesInstalledDispatcher(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	agentHome := filepath.Join(t.TempDir(), "agents")
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", agentHome)
	target := filepath.Join(t.TempDir(), "product")
	if _, err := Init(Options{Target: target, Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	dispatcherPath := filepath.Join(agentHome, ".codex", "skills", "spec-framework", "SKILL.md")
	if err := os.WriteFile(dispatcherPath, []byte("legacy dispatcher"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dispatcherPath)
	if err != nil {
		t.Fatal(err)
	}
	contract := string(data)
	if strings.Contains(contract, "legacy dispatcher") || !strings.Contains(contract, "Resolve framework-guide first unless") {
		t.Fatalf("upgrade did not refresh Guide-first dispatcher: %s", contract)
	}
}

func TestUpgradePreservesExistingFeatureBrief(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product")
	if _, err := Init(Options{Target: target, Agents: []Agent{Codex}, StartingPoint: "existing-feature"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(target, "product", "foundation", "feature-brief.md")
	if err := os.WriteFile(path, []byte("adopter-owned feature brief"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "adopter-owned feature brief" {
		t.Fatalf("upgrade changed feature brief: %q, %v", data, err)
	}
}

func TestUpgradePreservesImplementationAssessment(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product")
	if _, err := Init(Options{Target: target, Agents: []Agent{Codex}, StartingPoint: "existing-implementation"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(target, "product", "knowledge", "assessments", "implementation-assessment.md")
	if err := os.WriteFile(path, []byte("adopter-owned assessment"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "adopter-owned assessment" {
		t.Fatalf("upgrade changed implementation assessment: %q, %v", data, err)
	}
}

func TestUpgradePreservesProductBaseline(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product")
	if _, err := Init(Options{Target: target, Agents: []Agent{Codex}, StartingPoint: "existing-product"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(target, "product", "foundation", "product-baseline.md")
	if err := os.WriteFile(path, []byte("adopter-owned product baseline"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "adopter-owned product baseline" {
		t.Fatalf("upgrade changed product baseline: %q, %v", data, err)
	}
}
