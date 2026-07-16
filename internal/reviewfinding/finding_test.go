package reviewfinding

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindingRequiresProvenanceAndValidSeverity(t *testing.T) {
	f := Finding{ID: "RF-1", Source: "github", Reference: "https://example.test/pr/1#discussion", Severity: "warning", Description: "Missing assertion", Status: "open", Scope: "src/x.go", Evidence: "review comment", Owner: "qa"}
	if err := f.Validate(); err != nil {
		t.Fatal(err)
	}
	f.Source = ""
	if err := f.Validate(); err == nil {
		t.Fatal("missing provenance accepted")
	}
	f.Source, f.Severity = "github", "urgent"
	if err := f.Validate(); err == nil {
		t.Fatal("invalid severity accepted")
	}
}

func TestImportPreservesImmutableProviderEvidence(t *testing.T) {
	root := t.TempDir()
	finding := Finding{ID: "RF-1", Source: "github", Reference: "https://example.test/pr/1", Severity: "warning", Description: "Missing assertion", Status: "open", Scope: "src/x.go", Evidence: "review comment", Owner: "qa"}
	if _, err := Import(root, "github", []Finding{finding}); err != nil {
		t.Fatal(err)
	}
	if _, err := Import(root, "github", []Finding{finding}); err != nil {
		t.Fatal(err)
	}
	finding.Description = "changed"
	if _, err := Import(root, "github", []Finding{finding}); err == nil {
		t.Fatal("changed evidence accepted")
	}
	if _, err := os.Stat(filepath.Join(root, ".product", "reviews", "findings", "RF-1.json")); err != nil {
		t.Fatal(err)
	}
}

func TestFindingRoutesWithoutChangingAuthority(t *testing.T) {
	if got := (Finding{Scope: "missing security permission check"}).Route(); got != "security-review" {
		t.Fatalf("route=%s", got)
	}
	if got := (Finding{Description: "coverage missing for negative case"}).Route(); got != "qa" {
		t.Fatalf("route=%s", got)
	}
}
