package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditOnlyMutationMatrix(t *testing.T) {
	root := auditOnlyFixture(t)
	tests := []struct {
		name    string
		args    []string
		blocked bool
	}{
		{"validate read only", []string{"validate"}, false},
		{"validate registry", []string{"validate", "--write-registry"}, true},
		{"validate report", []string{"validate", "--write-report"}, true},
		{"validate disabled report flag", []string{"validate", "--write-report=false"}, false},
		{"list features", []string{"work"}, false},
		{"create workspace", []string{"work", "--feature", "FT-1"}, true},
		{"approve", []string{"approve", "--artifact", "foundation/problem/problem.md"}, true},
		{"review stage", []string{"review", "--stage", "specification"}, false},
		{"preview stage approval", []string{"approve-stage", "--stage", "specification"}, false},
		{"apply stage approval", []string{"approve-stage", "--stage", "specification", "--yes"}, true},
		{"graph ready", []string{"graph", "ready"}, false},
		{"graph materialize", []string{"graph", "materialize"}, true},
		{"design inspect", []string{"design", "inspect"}, false},
		{"design init", []string{"design", "init"}, true},
		{"design dry migration", []string{"design", "migrate", "--dry-run"}, false},
		{"design explicitly non-dry migration", []string{"design", "migrate", "--dry-run=false"}, true},
		{"decision preview", []string{"decisions", "migrate"}, false},
		{"decision apply", []string{"decisions", "migrate", "--yes"}, true},
		{"runtime resume", []string{"resume"}, false},
		{"runtime checkpoint", []string{"checkpoint"}, true},
		{"schedule writes waves", []string{"schedule"}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := append(append([]string{}, test.args...), "--product-root", root)
			blocked, _ := auditOnlyMutation(args)
			if blocked != test.blocked {
				t.Fatalf("blocked=%t want=%t args=%v", blocked, test.blocked, args)
			}
		})
	}
}

func TestAuditOnlyGuardResolvesCommandSpecificRoots(t *testing.T) {
	repository := t.TempDir()
	product := filepath.Join(repository, "product")
	writeAuditManifest(t, product)
	for _, args := range [][]string{
		{"adapters", "install", "impeccable", "--root", repository, "--yes"},
		{"migrate", "external-runtime", "--target", repository, "--yes"},
	} {
		if blocked, _ := auditOnlyMutation(args); !blocked {
			t.Fatalf("command-specific root bypassed audit guard: %v", args)
		}
	}
}

func TestAppRejectsAuditOnlyWriteBeforeHandler(t *testing.T) {
	root := auditOnlyFixture(t)
	var stdout, stderr bytes.Buffer
	code := New("test").Run([]string{"validate", "--product-root", root, "--write-report"}, &stdout, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), "audit-only blocks product mutation") {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(filepath.Join(root, "audits")); !os.IsNotExist(err) {
		t.Fatalf("blocked command wrote audits directory: %v", err)
	}
}

func auditOnlyFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeAuditManifest(t, root)
	return root
}

func writeAuditManifest(t *testing.T, root string) {
	t.Helper()
	manifestPath := filepath.Join(root, ".product", "framework.json")
	if err := os.MkdirAll(filepath.Dir(manifestPath), 0755); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(map[string]any{"framework": "spec-framework", "starting_point": "audit-only"})
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		t.Fatal(err)
	}
}
