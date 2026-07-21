package install

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/runtimeassets"
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
	if readmes != 22 {
		t.Fatalf("initialized product has %d READMEs, want 22", readmes)
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

func TestUpgradeDoesNotOptExistingSpecificationIntoV2(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	useCase := filepath.Join(target, "product", "domains", "d", "goals", "g", "features", "f", "use-cases", "u")
	if err := os.MkdirAll(useCase, 0755); err != nil {
		t.Fatal(err)
	}
	contextBody := []byte("---\nid: UC-1\nrigor_tier: S\n---\n")
	specificationBody := []byte("# Adopter-owned legacy Specification\n")
	if err := os.WriteFile(filepath.Join(useCase, "context.md"), contextBody, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(useCase, "specification.md"), specificationBody, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	for path, expected := range map[string][]byte{
		filepath.Join(useCase, "context.md"):       contextBody,
		filepath.Join(useCase, "specification.md"): specificationBody,
	} {
		actual, err := os.ReadFile(path)
		if err != nil || string(actual) != string(expected) {
			t.Fatalf("upgrade changed adopter Specification file %s: %q %v", path, actual, err)
		}
	}
}

func TestInitDiscoversSiblingCodeRootsAndPreservesThemOnUpgrade(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if err := os.MkdirAll(filepath.Join(target, "web"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "web", "package.json"), []byte(`{"name":"web"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	manifest, err := os.ReadFile(filepath.Join(target, "product", ".product", "framework.json"))
	if err != nil || !strings.Contains(string(manifest), `"path": "web"`) || !strings.Contains(string(manifest), `"role": "web"`) || !strings.Contains(string(manifest), `"mode": "cli-fallback"`) || !strings.Contains(string(manifest), `"status": "needs-agent-review"`) {
		t.Fatalf("code root missing from manifest: %s %v", manifest, err)
	}
	landscape, err := os.ReadFile(filepath.Join(target, "product", "knowledge", "assessments", "product-landscape.md"))
	if err != nil || !strings.Contains(string(landscape), "`web/`") || !strings.Contains(string(landscape), "pending comprehensive inventory") {
		t.Fatalf("landscape missing discovered root: %s %v", landscape, err)
	}
	if _, err := Upgrade(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	manifest, err = os.ReadFile(filepath.Join(target, "product", ".product", "framework.json"))
	if err != nil || !strings.Contains(string(manifest), `"path": "web"`) {
		t.Fatalf("upgrade discarded code roots: %s %v", manifest, err)
	}
}

func TestInitUsesAgentDeclaredRootsAsAuthoritative(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	for _, root := range []string{"web", "services/api"} {
		if err := os.MkdirAll(filepath.Join(target, filepath.FromSlash(root)), 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(target, "web", "package.json"), []byte(`{"name":"web"}`), 0644); err != nil {
		t.Fatal(err)
	}
	result, err := Init(Options{
		Target: target, Version: "test", Agents: []Agent{Codex},
		CodeRoots:             []runtimeassets.CodeRoot{{Path: "services/api", Role: "api"}},
		CodeRootDiscoveryMode: CodeRootDiscoveryAgentDeclared,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.CodeRootDiscovery.Mode != CodeRootDiscoveryAgentDeclared || result.CodeRootDiscovery.Status != "confirmed" {
		t.Fatalf("unexpected discovery result: %+v", result.CodeRootDiscovery)
	}
	manifest, err := os.ReadFile(filepath.Join(target, "product", ".product", "framework.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifest), `"path": "services/api"`) || strings.Contains(string(manifest), `"path": "web"`) || !strings.Contains(string(manifest), `"mode": "agent-declared"`) {
		t.Fatalf("agent-declared roots were not authoritative: %s", manifest)
	}
}

func TestInitAgentConfirmedNoCodeOverridesCLICandidates(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if err := os.MkdirAll(filepath.Join(target, "web"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "web", "package.json"), []byte(`{"name":"web"}`), 0644); err != nil {
		t.Fatal(err)
	}
	result, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}, CodeRootDiscoveryMode: CodeRootDiscoveryAgentConfirmedNone})
	if err != nil {
		t.Fatal(err)
	}
	if result.CodeRootDiscovery.Mode != CodeRootDiscoveryAgentConfirmedNone || result.CodeRootDiscovery.Status != "confirmed" {
		t.Fatalf("agent no-code decision was not authoritative: %+v", result.CodeRootDiscovery)
	}
	manifest, err := os.ReadFile(filepath.Join(target, "product", ".product", "framework.json"))
	if err != nil || strings.Contains(string(manifest), `"path": "web"`) {
		t.Fatalf("CLI candidate leaked into agent-confirmed no-code map: %s %v", manifest, err)
	}
}

func TestUpgradeConfirmsFallbackRootsWithoutOverwritingLandscape(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if err := os.MkdirAll(filepath.Join(target, "web"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "web", "package.json"), []byte(`{"name":"web"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	landscape := filepath.Join(target, "product", "knowledge", "assessments", "product-landscape.md")
	if err := os.WriteFile(landscape, []byte("adopter-owned landscape\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Version: "test", Agents: []Agent{Codex}, CodeRoots: []runtimeassets.CodeRoot{{Path: "web", Role: "frontend"}}, CodeRootDiscoveryMode: CodeRootDiscoveryAgentDeclared}); err != nil {
		t.Fatal(err)
	}
	manifest, err := os.ReadFile(filepath.Join(target, "product", ".product", "framework.json"))
	if err != nil || !strings.Contains(string(manifest), `"role": "frontend"`) || !strings.Contains(string(manifest), `"mode": "agent-declared"`) {
		t.Fatalf("upgrade did not confirm roots: %s %v", manifest, err)
	}
	data, err := os.ReadFile(landscape)
	if err != nil || string(data) != "adopter-owned landscape\n" {
		t.Fatalf("upgrade overwrote Product Landscape: %q %v", data, err)
	}
}

func TestUpgradeDoesNotOptLegacyProductsIntoBaselinePolicy(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(target, "product", ".product", "framework.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	delete(manifest, "baseline_policy")
	updated, _ := json.Marshal(manifest)
	if err := os.WriteFile(manifestPath, updated, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(manifestPath)
	if err != nil || strings.Contains(string(data), "baseline_policy") {
		t.Fatalf("upgrade opted legacy product into policy: %s %v", data, err)
	}
}

func TestUpgradeMarksPolicyEnabledLegacyDiscoveryForAgentReview(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product-repo")
	if _, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(target, "product", ".product", "framework.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	delete(manifest, "code_root_discovery")
	legacy, _ := json.Marshal(manifest)
	if err := os.WriteFile(manifestPath, legacy, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Version: "test", Agents: []Agent{Codex}}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(manifestPath)
	if err != nil || !strings.Contains(string(data), `"mode": "legacy-unclassified"`) || !strings.Contains(string(data), `"status": "needs-agent-review"`) {
		t.Fatalf("legacy discovery was not marked for review: %s %v", data, err)
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
		filepath.Join(agentHome, ".agents", "skills", "spec-framework", "SKILL.md"),
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

func TestInitShipsEngineeringCatalogRootsForEveryAgentTarget(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "product")
	result, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex, Cursor, Claude}})
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{
		"product/engineering/architecture/topology.yaml",
		"product/engineering/catalog/catalog.yaml",
		"product/engineering/standards/standards.yaml",
		"product/engineering/standards/profiles/product-default.yaml",
		"product/engineering/operations/operations.yaml",
		"product/engineering/evidence/inventory.md",
	} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(file))); err != nil {
			t.Errorf("missing %s: %v", file, err)
		}
	}
	for _, skill := range []string{"engineering-orchestrator", "technical-landscape", "engineering-standards", "operations-baseline", "engineering-evidence", "engineering-system"} {
		if _, err := os.Stat(filepath.Join(result.SpecRoot, "skills", skill, "SKILL.md")); err != nil {
			t.Errorf("runtime missing engineering skill %s: %v", skill, err)
		}
	}
	handoff, err := os.ReadFile(filepath.Join(result.SpecRoot, "skills", "engineering-orchestrator", "assets", "engineering-baseline-handoff-template.json"))
	if err != nil || !strings.Contains(string(handoff), `"mode": "sequential"`) || !strings.Contains(string(handoff), `"context_policy": "minimal"`) || !strings.Contains(string(handoff), `"skill": "engineering-system"`) {
		t.Fatalf("runtime engineering delegation contract missing: %v", err)
	}
	aggregate, err := os.ReadFile(filepath.Join(target, "product", "engineering", "engineering-system.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, owner := range []string{"technical-landscape", "engineering-standards", "operations-baseline", "engineering-evidence", "engineering-system"} {
		if !strings.Contains(string(aggregate), "owner_skill: "+owner) {
			t.Errorf("engineering aggregate missing owner %s", owner)
		}
	}
	owned := filepath.Join(target, "product", "engineering", "catalog", "catalog.yaml")
	if err := os.WriteFile(owned, []byte("adopter-owned catalog\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Upgrade(Options{Target: target, Version: "test", Agents: []Agent{Codex, Cursor, Claude}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(owned)
	if err != nil || string(data) != "adopter-owned catalog\n" {
		t.Fatalf("upgrade changed adopter engineering catalog: %q %v", data, err)
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
	for _, name := range []string{"inventory.json", "import-plan.json", "mapping.json", "traceability.json", "conflicts.md", "import-report.md"} {
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
	dispatcherPath := filepath.Join(agentHome, ".agents", "skills", "spec-framework", "SKILL.md")
	legacyDispatcherPath := filepath.Join(agentHome, ".codex", "skills", "spec-framework", "SKILL.md")
	if err := os.WriteFile(dispatcherPath, []byte("legacy dispatcher"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(legacyDispatcherPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacyDispatcherPath, []byte("legacy location"), 0644); err != nil {
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
	if _, err := os.Stat(legacyDispatcherPath); !os.IsNotExist(err) {
		t.Fatalf("upgrade did not remove legacy dispatcher: %v", err)
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
