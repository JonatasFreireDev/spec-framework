package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEngineeringSystemInspectJSON(t *testing.T) {
	root := t.TempDir()
	engineering := filepath.Join(root, "engineering")
	if err := os.MkdirAll(filepath.Join(engineering, "architecture"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"context.md":                     "---\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 0.1.0\norigin_mode: generate\n---\n",
		"engineering-system.md":          "# Engineering System\n",
		"architecture/system-context.md": "# Context\n",
		"engineering-system.yaml":        "schema_version: 1\nid: ENGSYS-TEST-001\nstatus: draft\nversion: 0.1.0\norigin_mode: generate\nscope: product\nareas:\n  system_context:\n    contract: architecture/system-context.md\n    maturity: baseline\n    evidence: []\ndecisions: []\nstandards: []\nfitness_functions: []\n",
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
	var stdout, stderr bytes.Buffer
	if code := runEngineeringSystem([]string{"inspect", "--product-root", root, "--json"}, &stdout, &stderr); code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"id": "ENGSYS-TEST-001"`) || !strings.Contains(stdout.String(), `"areas"`) {
		t.Fatalf("stdout=%s", stdout.String())
	}
}

func TestEngineeringSystemTriggersListsCanonicalValues(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runEngineeringSystem([]string{"triggers"}, &stdout, &stderr); code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "architecture_boundary_change") || !strings.Contains(stdout.String(), "migration") {
		t.Fatalf("stdout=%s", stdout.String())
	}
}
