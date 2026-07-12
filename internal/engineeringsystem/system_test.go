package engineeringsystem

import (
	"os"
	"path/filepath"
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
	write("engineering-system.md", "# Engineering System\n")
	write("architecture/modules.md", "# Modules\n")
	write("engineering-system.yaml", "scope: product\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: verified\n    evidence: []\ndecisions: []\nstandards: []\nfitness_functions: []\n")
	inspection, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(inspection.Blockers) != 1 || inspection.Blockers[0] != "area modules maturity verified requires evidence" {
		t.Fatalf("inspection=%+v", inspection)
	}
	write("engineering-system.yaml", "scope: product\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: verified\n    evidence: [tests/modules]\ndecisions: []\nstandards: []\nfitness_functions: []\n")
	inspection, err = Inspect(root)
	if err != nil || len(inspection.Blockers) != 0 {
		t.Fatalf("inspection=%+v err=%v", inspection, err)
	}
}
