package reviewfinding

import "testing"

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
