package install

import (
	"strings"
	"testing"
)

var bootstrapStartingPoints = []string{"new-product", "existing-product", "existing-documents", "existing-feature", "existing-implementation", "audit-only"}

func TestDeclarativeBootstrapProfilesAreComplete(t *testing.T) {
	for _, startingPoint := range bootstrapStartingPoints {
		t.Run(startingPoint, func(t *testing.T) {
			bootstrap, err := declarativeBootstrapFor(startingPoint)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(bootstrap, "# Product Bootstrap") || !strings.Contains(bootstrap, "## Guided steps") || !strings.Contains(bootstrap, "## Rules for the agent") {
				t.Fatalf("bootstrap is missing required sections: %s", bootstrap)
			}
			if !strings.Contains(bootstrap, "Starting point: **"+startingPoint+"**") {
				t.Fatalf("bootstrap does not identify starting point")
			}
		})
	}
}

func TestBootstrapContainsProfileSpecificReadingAndPrompts(t *testing.T) {
	cases := map[string][]string{
		"new-product":             {"product/context.md", "problem.md", "Vision", "Stop before approval"},
		"existing-product":        {"product-baseline.md", "runtime configuration", "future bets"},
		"existing-documents":      {"traceability.json", "mapping.json", "Do not automatically convert Epic to Feature", "materialize"},
		"existing-feature":        {"feature-brief.md", "walking skeleton", "Feature Brief"},
		"existing-implementation": {"implementation-assessment.md", "database", "full Foundation"},
		"audit-only":              {"read-only", "Do not edit product files", "Human decision"},
	}
	for startingPoint, expected := range cases {
		t.Run(startingPoint, func(t *testing.T) {
			bootstrap := bootstrapFor(startingPoint)
			for _, text := range expected {
				if !strings.Contains(strings.ToLower(bootstrap), strings.ToLower(text)) {
					t.Errorf("bootstrap missing %q", text)
				}
			}
		})
	}
}

func TestExistingDocumentsBootstrapPinsImportRun(t *testing.T) {
	bootstrap := bootstrapFor("existing-documents")
	if !strings.Contains(bootstrap, "<latest-run>") || !strings.Contains(bootstrap, "traceability.json") || !strings.Contains(bootstrap, "import materialize") {
		t.Fatalf("existing-documents bootstrap is not import-oriented: %s", bootstrap)
	}
}

func TestAuditOnlyBootstrapHasNoMutationInstructions(t *testing.T) {
	bootstrap := bootstrapFor("audit-only")
	for _, forbidden := range []string{"approve --", "import materialize", "work --feature", "write-registry"} {
		if strings.Contains(bootstrap, forbidden) {
			t.Fatalf("audit-only bootstrap contains mutation guidance %q", forbidden)
		}
	}
}
