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
