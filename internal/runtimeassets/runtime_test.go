package runtimeassets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverRequiresCanonicalManifest(t *testing.T) {
	root := t.TempDir()
	if _, _, err := Discover(root); err == nil {
		t.Fatal("mention-free directory activated without a manifest")
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("Spec Framework"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Discover(root); err == nil {
		t.Fatal("text mention activated the framework")
	}
	manifest := filepath.Join(root, "product", ".product", "framework.json")
	if err := os.MkdirAll(filepath.Dir(manifest), 0755); err != nil {
		t.Fatal(err)
	}
	data := []byte(`{"schema_version":3,"framework":"spec-framework","version":"1.2.3","activation":{"mode":"manifest-only"}}`)
	if err := os.WriteFile(manifest, data, 0644); err != nil {
		t.Fatal(err)
	}
	got, value, err := Discover(filepath.Join(root, "product"))
	if err != nil || got != root || value.Version != "1.2.3" {
		t.Fatalf("got root=%q version=%q err=%v", got, value.Version, err)
	}
}

func TestEnsureMaterializesVersionedAssets(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	root, err := Ensure("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"FRAMEWORK.md", "AGENTS.framework.md", "docs/execution-runtime.md", "docs/engineering-systems.md", "docs/lifecycle-and-approvals.md", "docs/artifact-registry-modules.md", "init/schema.json", "init/catalog.json", "init/bootstrap.json", "init/contracts/new-product.json", "skills/code-runner/SKILL.md", "skills/grill-me/SKILL.md", "skills/grill-me/agents/openai.yaml", "skills/framework-guide/agents/openai.yaml", "skills/framework-guide/scripts/inspect-workspace.ps1", "skills/framework-guide/scripts/inspect-workspace.sh", "skills/framework-guide/scripts/inspect-code-roots.ps1", "skills/product-orchestrator/scripts/inventory-product-landscape.sh", "skills/engineering-orchestrator/SKILL.md", "skills/technical-landscape/SKILL.md", "skills/technical-landscape/assets/technical-catalog-template.yaml", "skills/technical-landscape/assets/topology-template.yaml", "skills/technical-landscape/scripts/inventory-technical-landscape.ps1", "skills/engineering-standards/SKILL.md", "skills/engineering-standards/assets/standards-catalog-template.yaml", "skills/operations-baseline/SKILL.md", "skills/operations-baseline/assets/operations-catalog-template.yaml", "skills/engineering-system/assets/quality-system-template.yaml", "skills/engineering-evidence/SKILL.md", "skills/engineering-evidence/scripts/inventory-engineering-evidence.ps1", "skills/design-system/scripts/inventory-design-evidence.sh", "skills/artifact-importer/scripts/record-review-and-validate.ps1", "skills/documentation-writer/scripts/validate-artifacts.sh", "skills/discovery-and-challenge.md", "skills/problem-discovery/assets/interview-note-template.md", "skills/problem-discovery/assets/research-summary-template.md", "skills/specification/assets/specification-template.md", "skills/specification/assets/behavior-contract-template.md", "skills/specification/assets/security-contract-template.md", "skills/specification/references/contracts.yaml", "skills/specification/references/migration-v2.md", "examples/events/domains/events/domain.md", ".complete"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(name))); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"skills/engineering-orchestrator/assets/engineering-baseline-handoff-template.json", "skills/subagent-return-reviewer/assets/engineering-specialist-return-template.json"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(name))); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "engineering-catalog-and-standards.md")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "AGENTS.framework.md"))
	if err != nil || !strings.Contains(string(data), "Common Agent Rules") {
		t.Fatalf("runtime common agent rules missing: %v", err)
	}
	qualityTemplate, err := os.ReadFile(filepath.Join(root, "skills", "engineering-system", "assets", "quality-system-template.yaml"))
	if err != nil || !strings.Contains(string(qualityTemplate), "records: []") || strings.Contains(string(qualityTemplate), "QEX-001") {
		t.Fatalf("runtime quality template must not ship placeholder evidence: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "decisions")); !os.IsNotExist(err) {
		t.Fatalf("obsolete framework archive must not be materialized: %v", err)
	}
}
