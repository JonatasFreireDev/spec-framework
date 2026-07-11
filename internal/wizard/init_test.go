package wizard

import (
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/install"
)

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
