package decisions

import "testing"

func TestDomainPathsSupportDefaultsAndExtensions(t *testing.T) {
	paths := DomainPaths(map[string]any{"decisionDomains": map[string]any{"security": "knowledge/security-decisions"}})
	if paths["design"] != "design/decisions/" || paths["security"] != "knowledge/security-decisions/" {
		t.Fatalf("paths=%v", paths)
	}
	if DomainForPath("engineering/decisions/ADR-001.md", paths) != "engineering" {
		t.Fatal("engineering path was not classified")
	}
	if err := ValidatePath("design", "knowledge/decisions/DEC-001.md", paths); err == nil {
		t.Fatal("misplaced design decision was accepted")
	}
}
