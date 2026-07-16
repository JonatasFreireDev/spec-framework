package validator

import (
	"encoding/json"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/reviewfinding"
)

// validateReviewFindings checks only imported, provider-neutral review evidence.
// Finding files deliberately live outside product artifacts and cannot express an
// approval, task completion, or provider-side resolution.
func validateReviewFindings(snap Snapshot) []Diagnostic {
	const prefix = ".product/reviews/findings/"
	seen := map[string]string{}
	var diagnostics []Diagnostic
	for path, raw := range snap.JSON {
		if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, ".json") {
			continue
		}
		data, err := json.Marshal(raw)
		if err != nil {
			diagnostics = append(diagnostics, Diagnostic{Severity: Error, Check: "review-findings", File: path, Message: "Review finding cannot be read", Fix: "Write one valid provider-neutral finding JSON object."})
			continue
		}
		var finding reviewfinding.Finding
		if err := json.Unmarshal(data, &finding); err != nil {
			diagnostics = append(diagnostics, Diagnostic{Severity: Error, Check: "review-findings", File: path, Message: "Review finding is not a JSON object", Fix: "Use the normalized finding schema."})
			continue
		}
		if err := finding.Validate(); err != nil {
			diagnostics = append(diagnostics, Diagnostic{Severity: Error, Check: "review-findings", File: path, Message: err.Error(), Fix: "Provide required provenance, evidence, owner, severity, and status."})
			continue
		}
		if prior, exists := seen[finding.ID]; exists {
			diagnostics = append(diagnostics, Diagnostic{Severity: Error, Check: "review-findings", File: path, Message: "Duplicate review finding id also appears in " + prior, Fix: "Keep one immutable imported record per finding id."})
			continue
		}
		seen[finding.ID] = path
	}
	return diagnostics
}
