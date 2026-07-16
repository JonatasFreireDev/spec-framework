package reviewfinding

import (
	"errors"
	"strings"
)

// Finding is provider-neutral review evidence. It is deliberately unable to
// express approval or delivery state.
type Finding struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Reference   string `json:"reference"`
	DiffHash    string `json:"diff_hash,omitempty"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Scope       string `json:"scope"`
	Evidence    string `json:"evidence"`
	Owner       string `json:"owner"`
}

func (f Finding) Validate() error {
	for name, value := range map[string]string{"id": f.ID, "source": f.Source, "reference": f.Reference, "description": f.Description, "status": f.Status, "scope": f.Scope, "evidence": f.Evidence, "owner": f.Owner} {
		if strings.TrimSpace(value) == "" {
			return errors.New("review finding " + name + " is required")
		}
	}
	if !map[string]bool{"blocker": true, "required_fix": true, "warning": true, "note": true}[strings.ToLower(strings.TrimSpace(f.Severity))] {
		return errors.New("review finding severity is invalid")
	}
	if !map[string]bool{"open": true, "routed": true, "resolved": true, "wontfix": true}[strings.ToLower(strings.TrimSpace(f.Status))] {
		return errors.New("review finding status is invalid")
	}
	return nil
}

// Route proposes the owning framework skill. It is advisory and never changes
// a task, review thread, approval, or provider-side state.
func (f Finding) Route() string {
	scope := strings.ToLower(f.Scope + " " + f.Description)
	if strings.Contains(scope, "security") || strings.Contains(scope, "permission") || strings.Contains(scope, "privacy") {
		return "security-review"
	}
	if strings.Contains(scope, "test") || strings.Contains(scope, "coverage") {
		return "qa"
	}
	if strings.Contains(scope, "requirement") || strings.Contains(scope, "decision") || strings.Contains(scope, "scope") {
		return "product-historian"
	}
	return "bug-fixer"
}
