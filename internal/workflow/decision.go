package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type DecisionImpact struct {
	ID                                                             string `json:"id"`
	Type                                                           string `json:"type"`
	Status                                                         string `json:"status"`
	Path                                                           string `json:"path"`
	Valid                                                          bool   `json:"valid"`
	AffectedArtifacts []string `json:"affected_artifacts,omitempty"`
	References []string `json:"references,omitempty"`
	StaleArtifacts []string `json:"stale_artifacts,omitempty"`
	PropagationGaps []string `json:"propagation_gaps,omitempty"`
	WorkflowEffects                                                map[string]any `json:"workflow_effects,omitempty"`
	Blockers                                                       []string       `json:"blockers,omitempty"`
}

func DecisionImpactReport(root, id string) (DecisionImpact, error) {
	var index map[string]any
	if err := readJSON(filepath.Join(root, ".product", "decisions.json"), &index); err != nil {
		return DecisionImpact{}, err
	}
	items, _ := index["decisions"].([]any)
	var d map[string]any
	for _, raw := range items {
		x, _ := raw.(map[string]any)
		if fmt.Sprint(x["id"]) == id {
			d = x
		}
	}
	if d == nil {
		return DecisionImpact{}, fmt.Errorf("decision %s is not indexed", id)
	}
	decisionType := fmt.Sprint(d["type"]); if decisionType == "<nil>" || decisionType == "" { scope:=fmt.Sprint(d["scope"]); switch { case strings.Contains(scope,"architecture"): decisionType="architecture"; case strings.Contains(scope,"security"): decisionType="security"; case strings.Contains(scope,"data"): decisionType="data"; case strings.Contains(scope,"release")||strings.Contains(scope,"delivery"):decisionType="delivery"; default:decisionType="product" } }
	r := DecisionImpact{ID: id, Type: decisionType, Status: fmt.Sprint(d["status"]), Path: filepath.ToSlash(fmt.Sprint(d["path"])), WorkflowEffects: map[string]any{}}
	r.AffectedArtifacts = stringAnySlice(d["affectedArtifacts"])
	if x, ok := d["workflowEffects"].(map[string]any); ok {
		r.WorkflowEffects = x
	}
	decisionPath := filepath.Join(root, filepath.FromSlash(r.Path))
	r.Valid = r.Status == "approved" && hasCurrentApproval(root, decisionPath, "approved")
	if r.Status != "approved" {
		r.Blockers = append(r.Blockers, "decision is not approved")
	}
	if !hasCurrentApproval(root, decisionPath, "approved") {
		r.Blockers = append(r.Blockers, "current approval record is missing")
	}
	_ = filepath.WalkDir(root, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if e.IsDir() {
			if rel == ".git" || rel == ".product" {
				return filepath.SkipDir
			}
			return nil
		}
		if rel == r.Path || !(strings.HasSuffix(rel, ".md") || strings.HasSuffix(rel, ".json")) {
			return nil
		}
		b, _ := os.ReadFile(path)
		if regexp.MustCompile(`\b` + regexp.QuoteMeta(id) + `\b`).Match(b) {
			r.References = append(r.References, rel)
		}
		return nil
	})
	var deriv map[string]any
	_ = readJSON(filepath.Join(root, ".product", "derivations.json"), &deriv)
	entries, _ := deriv["derivations"].([]any)
	for _, raw := range entries {
		entry, _ := raw.(map[string]any)
		sources, _ := entry["derived_from"].([]any)
		for _, sr := range sources {
			s, _ := sr.(map[string]any)
			if fmt.Sprint(s["artifact_id"]) == id {
				b, _ := os.ReadFile(decisionPath)
				if Hash(string(b)) != fmt.Sprint(s["content_hash"]) {
					r.StaleArtifacts = append(r.StaleArtifacts, fmt.Sprint(entry["path"]))
				}
			}
		}
	}
	ref := map[string]bool{}
	for _, x := range r.References {
		ref[x] = true
	}
	for _, x := range r.AffectedArtifacts {
		x = filepath.ToSlash(x)
		if strings.HasPrefix(x, "../") || strings.HasSuffix(x, "/") {
			continue
		}
		if !ref[x] {
			r.PropagationGaps = append(r.PropagationGaps, x)
		}
	}
	sort.Strings(r.AffectedArtifacts)
	sort.Strings(r.References)
	sort.Strings(r.StaleArtifacts)
	sort.Strings(r.PropagationGaps)
	return r, nil
}
func stringAnySlice(v any) []string {
	raw, _ := v.([]any)
	var out []string
	for _, x := range raw {
		out = append(out, fmt.Sprint(x))
	}
	return out
}

func decisionEffectsFor(root string, ids []string) []map[string]any {
	var index map[string]any
	if readJSON(filepath.Join(root, ".product", "decisions.json"), &index) != nil {
		return nil
	}
	wanted := map[string]bool{}
	for _, id := range ids {
		wanted[id] = true
	}
	var out []map[string]any
	for _, raw := range stringMapSlice(index["decisions"]) {
		if wanted[fmt.Sprint(raw["id"])] && fmt.Sprint(raw["status"]) == "approved" {
			out = append(out, raw)
		}
	}
	return out
}
func stringMapSlice(v any) []map[string]any {
	raw, _ := v.([]any)
	var out []map[string]any
	for _, x := range raw {
		if m, ok := x.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}
