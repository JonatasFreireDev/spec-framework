package reviewfinding

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

// Import stores a provider-neutral snapshot of review findings. It has no
// provider client and never resolves a remote thread, edits a task, or changes
// approval state. Existing immutable snapshots may only be imported again when
// their complete content matches.
func Import(root, source string, findings []Finding) ([]Finding, error) {
	source = strings.TrimSpace(source)
	if source == "" || strings.ContainsAny(source, `/\\`) {
		return nil, errors.New("review import source is required and must be a simple provider name")
	}
	dir := filepath.Join(root, ".product", "reviews", "findings")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, finding := range findings {
		if err := finding.Validate(); err != nil {
			return nil, err
		}
		if finding.Source != source {
			return nil, fmt.Errorf("review finding %s source %q does not match import source %q", finding.ID, finding.Source, source)
		}
		if !safeID(finding.ID) {
			return nil, fmt.Errorf("review finding id %q is unsafe", finding.ID)
		}
		if seen[finding.ID] {
			return nil, fmt.Errorf("duplicate review finding id %s in import", finding.ID)
		}
		seen[finding.ID] = true
		path := filepath.Join(dir, finding.ID+".json")
		payload, err := json.MarshalIndent(finding, "", "  ")
		if err != nil {
			return nil, err
		}
		payload = append(payload, '\n')
		if prior, err := os.ReadFile(path); err == nil {
			if hash(prior) != hash(payload) {
				return nil, fmt.Errorf("review finding %s already exists with different immutable evidence", finding.ID)
			}
			continue
		} else if !os.IsNotExist(err) {
			return nil, err
		}
		if err := atomicWrite(path, payload); err != nil {
			return nil, err
		}
	}
	return findings, nil
}

func safeID(id string) bool {
	return id != "" && filepath.Base(id) == id && !strings.ContainsAny(id, `/\\`)
}
func hash(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func atomicWrite(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
