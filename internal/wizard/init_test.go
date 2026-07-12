package wizard

import (
	"bytes"
	"os"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/install"
)

func TestProgramOptionsLeaveStandardInputToBubbleTea(t *testing.T) {
	output := &bytes.Buffer{}
	if got := len(programOptions(os.Stdin, output)); got != 1 {
		t.Fatalf("standard input options = %d, want output option only", got)
	}
	if got := len(programOptions(bytes.NewBufferString("input"), output)); got != 2 {
		t.Fatalf("injected input options = %d, want input and output options", got)
	}
}

func TestValidateSources(t *testing.T) {
	tests := []struct {
		name          string
		startingPoint string
		sources       string
		wantError     bool
	}{
		{name: "existing documents without sources", startingPoint: "existing-documents", wantError: true},
		{name: "existing documents with sources", startingPoint: "existing-documents", sources: "docs/prd.md"},
		{name: "new product without sources", startingPoint: "new-product"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSources(tt.startingPoint, tt.sources)
			if (err != nil) != tt.wantError {
				t.Fatalf("validateSources() error = %v, wantError = %t", err, tt.wantError)
			}
		})
	}
}

func TestShowSourcePathsOnlyForExistingDocuments(t *testing.T) {
	for _, startingPoint := range []string{"new-product", "existing-product", "existing-feature", "existing-implementation", "audit-only"} {
		if showSourcePaths(startingPoint) {
			t.Errorf("showSourcePaths(%q) = true, want false", startingPoint)
		}
	}
	if !showSourcePaths("existing-documents") {
		t.Error("showSourcePaths(existing-documents) = false, want true")
	}
}

func TestAgentOptionsMapChoices(t *testing.T) {
	opts := agentOptions()
	if len(opts) != len(choices) {
		t.Fatalf("got %d options, want %d", len(opts), len(choices))
	}
	for i, opt := range opts {
		want := choices[i]
		if opt.Value != want.Agent {
			t.Errorf("option %d value = %q, want %q", i, opt.Value, want.Agent)
		}
		if opt.Key != want.Name {
			t.Errorf("option %d key = %q, want %q", i, opt.Key, want.Name)
		}
	}
}

func TestResultAgentNames(t *testing.T) {
	r := Result{Agents: []install.Agent{install.Codex, install.Claude}}
	got := r.AgentNames()
	want := []string{"codex", "claude"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
