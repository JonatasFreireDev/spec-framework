package engineeringsystem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTriggersParsesAllowedAndRejectsUnknown(t *testing.T) {
	valid, invalid := Triggers("---\nengineering_triggers:\n  - migration\n  - new_dependency\n  - magic_change\n---\n")
	if len(valid) != 2 || valid[0] != "migration" || valid[1] != "new_dependency" {
		t.Fatalf("valid=%v", valid)
	}
	if len(invalid) != 1 || invalid[0] != "magic_change" {
		t.Fatalf("invalid=%v", invalid)
	}
}

func TestInspectValidatesCatalogContractsAndMaturityEvidence(t *testing.T) {
	root := t.TempDir()
	engineering := filepath.Join(root, "engineering")
	if err := os.MkdirAll(filepath.Join(engineering, "architecture"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(path, text string) {
		if err := os.WriteFile(filepath.Join(engineering, filepath.FromSlash(path)), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("context.md", "---\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 1.2.3\norigin_mode: generate\n---\n")
	write("engineering-system.md", "| Field | Value |\n| --- | --- |\n| ID | `ENGSYS-TEST-001` |\n| Status | `draft` |\n| Version | `1.2.3` |\n")
	write("architecture/modules.md", "# Modules\n")
	write("engineering-system.yaml", "schema_version: 1\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 1.2.3\norigin_mode: generate\nscope: product\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: verified\n    evidence: []\ndecisions: []\nstandards: []\nfitness_functions: []\n")
	inspection, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(inspection.Blockers) != 1 || inspection.Blockers[0] != "area modules maturity verified requires evidence" {
		t.Fatalf("inspection=%+v", inspection)
	}
	write("engineering-system.yaml", "schema_version: 1\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 1.2.3\norigin_mode: generate\nscope: product\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: verified\n    evidence:\n      - tests/modules\ndecisions: []\nstandards: []\nfitness_functions: []\n")
	inspection, err = Inspect(root)
	if err != nil || len(inspection.Blockers) != 0 {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
	if err := os.WriteFile(filepath.Join(root, "outside.md"), []byte("outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	write("engineering-system.yaml", "schema_version: 1\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 1.2.3\norigin_mode: generate\nscope: product\nareas:\n  modules:\n    contract: ../outside.md\n    maturity: baseline\n    evidence: []\n")
	inspection, err = Inspect(root)
	foundEscape := false
	for _, blocker := range inspection.Blockers {
		foundEscape = foundEscape || strings.Contains(blocker, "escapes engineering")
	}
	if err != nil || !foundEscape {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
}

func TestTriggersAcceptsInlineYAMLList(t *testing.T) {
	valid, invalid := Triggers("---\nengineering_triggers: [migration, external_integration]\n---\n")
	if len(invalid) != 0 || len(valid) != 2 || valid[0] != "external_integration" || valid[1] != "migration" {
		t.Fatalf("valid=%v invalid=%v", valid, invalid)
	}
}

func TestInspectRejectsCatalogIdentityMismatch(t *testing.T) {
	root := t.TempDir()
	engineering := filepath.Join(root, "engineering")
	files := map[string]string{
		"context.md":                     "---\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering-system.md":          "| Field | Value |\n| --- | --- |\n| ID | `ENGSYS-TEST-001` |\n| Status | `draft` |\n| Version | `1.0.0` |\n",
		"architecture/system-context.md": "# Context\n",
		"engineering-system.yaml":        "schema_version: 1\nid: ENGSYS-OTHER-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: product\nareas:\n  context:\n    contract: architecture/system-context.md\n    maturity: baseline\n    evidence: []\n",
	}
	for name, body := range files {
		path := filepath.Join(engineering, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	inspection, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(inspection.Blockers) != 1 || inspection.Blockers[0] != "context and catalog id do not match" {
		t.Fatalf("blockers=%v", inspection.Blockers)
	}
}

func TestMigrateAddsSchemaVersionWithoutChangingOtherFields(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "engineering")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "engineering-system.yaml")
	original := "id: ENGSYS-001\nstatus: draft\ncustom_field: preserve-me\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	items, err := Migrate(root, true)
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%v err=%v", items, err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != original {
		t.Fatal("dry-run changed the catalog")
	}
	if _, err := Migrate(root, false); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	text := string(data)
	if !strings.Contains(text, "schema_version: 1") || !strings.Contains(text, "custom_field: preserve-me") {
		t.Fatalf("catalog=%s", text)
	}
}

func TestMigrateMaterializesLegacyQualitySystemWithoutOverwriting(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "engineering")
	if err := os.MkdirAll(filepath.Join(dir, "quality"), 0o755); err != nil {
		t.Fatal(err)
	}
	catalog := "schema_version: 1\nid: ENGSYS-001\nstatus: draft\nversion: 1.2.3\norigin_mode: evolve\nscope: product\nareas:\n  quality:\n    contract: quality/quality-model.md\n    maturity: baseline\n    evidence: []\n"
	if err := os.WriteFile(filepath.Join(dir, "engineering-system.yaml"), []byte(catalog), 0o644); err != nil {
		t.Fatal(err)
	}
	model := filepath.Join(dir, "quality", "quality-model.md")
	if err := os.WriteFile(model, []byte("owned model\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	changes, err := Migrate(root, true)
	if err != nil || len(changes) != 4 {
		t.Fatalf("changes=%v err=%v", changes, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "quality", "quality-system.md")); !os.IsNotExist(err) {
		t.Fatal("dry-run materialized quality files")
	}
	if _, err := Migrate(root, false); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "engineering-system.yaml"))
	if !strings.Contains(string(data), "contract: quality/quality-system.md") {
		t.Fatalf("catalog=%s", data)
	}
	for _, name := range []string{"quality-system.md", "quality-system.yaml", "test-strategy.md"} {
		if _, err := os.Stat(filepath.Join(dir, "quality", name)); err != nil {
			t.Fatal(err)
		}
	}
	data, _ = os.ReadFile(model)
	if string(data) != "owned model\n" {
		t.Fatal("migration overwrote legacy quality model")
	}
}

func TestMigrateRollsBackGeneratedFilesWhenMaterializationFails(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "engineering", "quality")
	if err := os.MkdirAll(filepath.Join(dir, "quality-system.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	original := "schema_version: 1\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: evolve\nscope: product\nareas:\n  quality:\n    contract: quality/quality-model.md\n    maturity: baseline\n    evidence: []\n"
	catalogPath := filepath.Join(root, "engineering", "engineering-system.yaml")
	if err := os.WriteFile(catalogPath, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Migrate(root, false); err == nil {
		t.Fatal("expected migration failure")
	}
	data, _ := os.ReadFile(catalogPath)
	if string(data) != original {
		t.Fatalf("catalog was not rolled back: %s", data)
	}
	for _, name := range []string{"quality-system.yaml", "test-strategy.md"} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Fatalf("generated file survived rollback: %s", name)
		}
	}
}

func TestCompositeHashChangesWithAnyEngineeringContract(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "engineering")
	if err := os.MkdirAll(filepath.Join(dir, "architecture"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "engineering-system.md"), []byte("system\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	contract := filepath.Join(dir, "architecture", "modules.md")
	if err := os.WriteFile(contract, []byte("modules-v1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := CompositeHash(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(contract, []byte("modules-v2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	after, err := CompositeHash(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if before == after {
		t.Fatal("composite hash ignored engineering contract change")
	}
}

func TestInspectValidatesConfiguredQualitySystem(t *testing.T) {
	root := t.TempDir()
	engineering := filepath.Join(root, "engineering")
	files := map[string]string{
		"context.md":                  "---\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering-system.md":       "| Field | Value |\n| --- | --- |\n| ID | `ENGSYS-TEST-001` |\n| Status | `draft` |\n| Version | `1.0.0` |\n",
		"quality/quality-system.md":   "| Field | Value |\n| --- | --- |\n| Engineering System | `ENGSYS-TEST-001 @ 1.0.0` |\n| Status | `draft` |\n\n| Area | Policy | Evidence | Maturity |\n| --- | --- | --- | --- |\n| Behavioral | strategy | none | baseline |\n| Accessibility | strategy | none | baseline |\n| Security and privacy | strategy | none | baseline |\n| Performance and reliability | model | none | baseline |\n| Observability | model | none | baseline |\n",
		"quality/quality-model.md":    "# Model\n",
		"quality/test-strategy.md":    "# Strategy\n",
		"engineering-system.yaml":     "schema_version: 1\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: product\nareas:\n  quality:\n    contract: quality/quality-system.md\n    maturity: baseline\n    evidence: []\n",
		"quality/quality-system.yaml": "schema_version: 1\nengineering_system: ENGSYS-TEST-001\nversion: 1.0.0\nstatus: draft\nareas:\n  behavioral: {maturity: baseline, policy: test-strategy.md, required_evidence: []}\n  accessibility: {maturity: baseline, policy: test-strategy.md, required_evidence: []}\n  security_privacy: {maturity: baseline, policy: test-strategy.md, delegated_gate: security-review, required_evidence: []}\n  performance_reliability: {maturity: baseline, policy: quality-model.md, required_evidence: []}\n  observability: {maturity: baseline, policy: quality-model.md, required_evidence: []}\ngate_source: knowledge/conventions/gates.md\nexceptions:\n  require_owner: true\n  require_residual_risk: true\n  require_expiry_or_review: true\n",
	}
	for name, body := range files {
		path := filepath.Join(engineering, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	inspection, err := Inspect(root)
	if err != nil || len(inspection.Blockers) != 0 || !inspection.QualitySystem {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
	catalogPath := filepath.Join(engineering, "quality", "quality-system.yaml")
	catalogData, _ := os.ReadFile(catalogPath)
	escaping := strings.Replace(string(catalogData), "policy: test-strategy.md", "policy: ../architecture/modules.md", 1)
	if err := os.WriteFile(catalogPath, []byte(escaping), 0o644); err != nil {
		t.Fatal(err)
	}
	inspection, err = Inspect(root)
	foundEscape := false
	for _, blocker := range inspection.Blockers {
		foundEscape = foundEscape || strings.Contains(blocker, "escapes engineering/quality")
	}
	if err != nil || !foundEscape {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
	if err := os.WriteFile(catalogPath, catalogData, 0o644); err != nil {
		t.Fatal(err)
	}
	incompleteException := strings.Replace(string(catalogData), "  require_expiry_or_review: true", "  require_expiry_or_review: true\n  records:\n    - id: QEX-TEST\n      owner: qa\n      status: open", 1)
	if err := os.WriteFile(catalogPath, []byte(incompleteException), 0o644); err != nil {
		t.Fatal(err)
	}
	inspection, err = Inspect(root)
	foundIncomplete := false
	for _, blocker := range inspection.Blockers {
		foundIncomplete = foundIncomplete || strings.Contains(blocker, "lacks scope, owner, rationale")
	}
	if err != nil || !foundIncomplete {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
	if err := os.WriteFile(catalogPath, catalogData, 0o644); err != nil {
		t.Fatal(err)
	}
	lifecycleExceptions := strings.Replace(string(catalogData), "  require_expiry_or_review: true", "  require_expiry_or_review: true\n  records:\n    - id: QEX-ACTIVE\n      scope: product\n      owner: qa\n      rationale: temporary\n      residual_risk: low\n      mitigation: monitor\n      expiry_or_review: 2999-01-01\n      reentry_gate: qa\n      status: open\n    - id: QEX-CLOSED\n      scope: product\n      owner: qa\n      rationale: resolved\n      residual_risk: none\n      mitigation: fixed\n      expiry_or_review: 2999-01-01\n      reentry_gate: qa\n      status: closed\n    - id: QEX-EXPIRED\n      scope: product\n      owner: qa\n      rationale: overdue\n      residual_risk: high\n      mitigation: none\n      expiry_or_review: 2000-01-01\n      reentry_gate: qa\n      status: open", 1)
	if err := os.WriteFile(catalogPath, []byte(lifecycleExceptions), 0o644); err != nil {
		t.Fatal(err)
	}
	inspection, err = Inspect(root)
	foundExpired := false
	for _, blocker := range inspection.Blockers {
		foundExpired = foundExpired || strings.Contains(blocker, "QEX-EXPIRED is open but expired")
	}
	if err != nil || len(inspection.QualityExceptions) != 1 || inspection.QualityExceptions[0] != "QEX-ACTIVE" || !foundExpired {
		t.Fatalf("active exceptions=%v blockers=%v err=%v", inspection.QualityExceptions, inspection.Blockers, err)
	}
	if err := os.WriteFile(catalogPath, catalogData, 0o644); err != nil {
		t.Fatal(err)
	}
	humanPath := filepath.Join(engineering, "quality", "quality-system.md")
	humanData, _ := os.ReadFile(humanPath)
	governedCatalog := strings.Replace(string(catalogData), "maturity: baseline, policy: test-strategy.md, required_evidence: []", "maturity: governed, policy: test-strategy.md, required_evidence: [missing.log]", 1)
	governedHuman := strings.Replace(string(humanData), "| Behavioral | strategy | none | baseline |", "| Behavioral | strategy | missing.log | governed |", 1)
	if err := os.WriteFile(catalogPath, []byte(governedCatalog), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(humanPath, []byte(governedHuman), 0o644); err != nil {
		t.Fatal(err)
	}
	inspection, err = Inspect(root)
	foundMissingEvidence := false
	for _, blocker := range inspection.Blockers {
		foundMissingEvidence = foundMissingEvidence || strings.Contains(blocker, "invalid or missing evidence missing.log")
	}
	if err != nil || !foundMissingEvidence {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
	if err := os.WriteFile(catalogPath, catalogData, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(humanPath, humanData, 0o644); err != nil {
		t.Fatal(err)
	}
	mismatchedHuman := strings.Replace(string(humanData), "| Behavioral | strategy | none | baseline |", "| Behavioral | strategy | none | governed |", 1)
	if err := os.WriteFile(humanPath, []byte(mismatchedHuman), 0o644); err != nil {
		t.Fatal(err)
	}
	inspection, err = Inspect(root)
	foundMismatch := false
	for _, blocker := range inspection.Blockers {
		foundMismatch = foundMismatch || strings.Contains(blocker, "human and mechanical maturity differ for behavioral")
	}
	if err != nil || !foundMismatch {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
	if err := os.WriteFile(humanPath, humanData, 0o644); err != nil {
		t.Fatal(err)
	}
	os.Remove(filepath.Join(engineering, "quality", "test-strategy.md"))
	inspection, err = Inspect(root)
	if err != nil || len(inspection.Blockers) == 0 {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
}
