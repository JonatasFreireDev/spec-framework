package decisioncheck

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunChecksDomainsUnindexedReferencesAndFixLinks(t *testing.T) {
	root := t.TempDir()
	write := func(path, text string) {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(".product/decisions.json", mustJSON(map[string]any{"decisionDomains": map[string]any{"security": "knowledge/security-decisions"}, "decisions": []any{
		map[string]any{"id": "DEC-101", "domain": "design", "type": "architecture", "status": "proposed", "path": "design/decisions/DEC-101.md", "affectedArtifacts": []any{"design.md"}},
		map[string]any{"id": "DEC-102", "domain": "engineering", "type": "data", "status": "proposed", "path": "engineering/decisions/DEC-102.md", "affectedArtifacts": []any{"engineering.md"}},
	}}))
	write("design/decisions/DEC-101.md", "# DEC-101\n")
	write("engineering/decisions/DEC-102.md", "# DEC-102\n")
	write("design.md", "---\ndecisions:\n  - DEC-101\n---\n\nUses DEC-101.\n")
	report, err := Run(Options{Root: root, FrameworkRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	if !has(report, "decision-links") {
		t.Fatalf("missing link finding: %+v", report.Diagnostics)
	}
	if _, err := Run(Options{Root: root, FrameworkRoot: root, FixLinks: true}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(root, "design.md"))
	if strings.Contains(string(data), "[DEC-101](design/decisions/DEC-101.md)") {
		t.Fatalf("link preview changed files unexpectedly: %s", data)
	}
	if _, err := Run(Options{Root: root, FrameworkRoot: root, FixLinks: true, Yes: true}); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(filepath.Join(root, "design.md"))
	if strings.Contains(string(data), "- [DEC-101]") {
		t.Fatalf("frontmatter was modified: %s", data)
	}
	if !strings.Contains(string(data), "[DEC-101](design/decisions/DEC-101.md)") {
		t.Fatalf("link was not fixed: %s", data)
	}
}

func TestRunReportsInvalidDomainAndUnindexedDecision(t *testing.T) {
	root := t.TempDir()
	write := func(path, text string) {
		full := filepath.Join(root, filepath.FromSlash(path))
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		_ = os.WriteFile(full, []byte(text), 0o644)
	}
	write(".product/decisions.json", `{"decisions":[{"id":"DEC-201","domain":"design","type":"data","status":"proposed","path":"engineering/decisions/DEC-201.md"}]}`)
	write("engineering/decisions/DEC-201.md", "# DEC-201\n")
	write("design/decisions/DEC-202.md", "# DEC-202\n")
	report, err := Run(Options{Root: root, FrameworkRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	if !has(report, "decision-domain") || !has(report, "decisions-unindexed") {
		t.Fatalf("diagnostics=%+v", report.Diagnostics)
	}
}

func has(report Report, check string) bool {
	for _, d := range report.Diagnostics {
		if d.Check == check {
			return true
		}
	}
	return false
}
func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
