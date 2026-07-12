package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DecisionMigrationItem struct {
	ID           string   `json:"id"`
	InferredType string   `json:"inferred_type"`
	Scope        []string `json:"scope"`
	NeedsReview  bool     `json:"needs_review"`
	Reasons      []string `json:"reasons"`
}

type DecisionMigrationPlan struct {
	Root, IndexPath string
	Items           []DecisionMigrationItem
	Original        map[string]any
}

type DecisionMigrationResult struct {
	Changed           int
	Backup, IndexPath string
}

func PlanDecisionMigration(root string) (DecisionMigrationPlan, error) {
	path := filepath.Join(root, ".product", "decisions.json")
	var index map[string]any
	if err := readJSON(path, &index); err != nil {
		return DecisionMigrationPlan{}, err
	}
	plan := DecisionMigrationPlan{Root: root, IndexPath: path, Original: index}
	for _, d := range stringMapSlice(index["decisions"]) {
		id := fmt.Sprint(d["id"])
		_, hasType := d["type"].(string)
		_, hasEffects := d["workflowEffects"].(map[string]any)
		_, hasStructuredScope := d["scope"].([]any)
		scopeValues := stringAnySlice(d["scope"])
		if len(scopeValues) == 0 {
			scopeValues = decisionScopePaths(d)
		}
		if hasType && hasEffects && hasStructuredScope {
			continue
		}
		kind, ambiguous := inferDecisionType(fmt.Sprint(d["scope"]))
		item := DecisionMigrationItem{ID: id, InferredType: kind, Scope: scopeValues, NeedsReview: ambiguous}
		if !hasType {
			item.Reasons = append(item.Reasons, "missing type")
		}
		if len(scopeValues) == 0 {
			item.NeedsReview = true
			item.Reasons = append(item.Reasons, "missing path scope")
		}
		if !hasEffects {
			item.Reasons = append(item.Reasons, "missing workflowEffects")
		}
		plan.Items = append(plan.Items, item)
	}
	return plan, nil
}

func ApplyDecisionMigration(plan DecisionMigrationPlan, items []DecisionMigrationItem) (DecisionMigrationResult, error) {
	selected := map[string]DecisionMigrationItem{}
	for _, item := range items {
		selected[item.ID] = item
	}
	cloneData, err := json.Marshal(plan.Original)
	if err != nil {
		return DecisionMigrationResult{}, err
	}
	var index map[string]any
	if err = json.Unmarshal(cloneData, &index); err != nil {
		return DecisionMigrationResult{}, err
	}
	decisions, _ := index["decisions"].([]any)
	changed := 0
	for _, raw := range decisions {
		d, _ := raw.(map[string]any)
		item, ok := selected[fmt.Sprint(d["id"])]
		if !ok {
			continue
		}
		if _, ok = d["type"].(string); !ok {
			d["type"] = item.InferredType
		}
		if legacy, ok := d["scope"].(string); ok {
			d["legacyScope"] = legacy
		}
		d["scope"] = stringsToAny(item.Scope)
		if _, ok = d["workflowEffects"].(map[string]any); !ok {
			d["workflowEffects"] = emptyDecisionEffects()
		}
		changed++
	}
	if changed == 0 {
		return DecisionMigrationResult{IndexPath: plan.IndexPath}, nil
	}
	index["schemaVersion"] = 2
	backupDir := filepath.Join(plan.Root, ".product", "migrations")
	if err = os.MkdirAll(backupDir, 0755); err != nil {
		return DecisionMigrationResult{}, err
	}
	backup := filepath.Join(backupDir, "decisions-v1-"+time.Now().UTC().Format("20060102T150405.000000000Z")+".json")
	original, err := os.ReadFile(plan.IndexPath)
	if err != nil {
		return DecisionMigrationResult{}, err
	}
	if err = atomicWrite(backup, original); err != nil {
		return DecisionMigrationResult{}, err
	}
	if err = writeJSON(plan.IndexPath, index); err != nil {
		_ = atomicWrite(plan.IndexPath, original)
		_ = os.Remove(backup)
		return DecisionMigrationResult{}, err
	}
	return DecisionMigrationResult{Changed: changed, Backup: backup, IndexPath: plan.IndexPath}, nil
}

func decisionScopePaths(d map[string]any) []string {
	var out []string
	for _, item := range stringAnySlice(d["affectedArtifacts"]) {
		item = filepath.ToSlash(strings.TrimSpace(item))
		if item != "" && !strings.HasPrefix(item, "../") && !strings.HasPrefix(item, ".product/") {
			out = append(out, item)
		}
	}
	return uniqueSorted(out)
}

func inferDecisionType(scope string) (string, bool) {
	scope = strings.ToLower(scope)
	var matches []string
	for _, kind := range []string{"architecture", "security", "data", "delivery", "product"} {
		if strings.Contains(scope, kind) || kind == "delivery" && strings.Contains(scope, "release") {
			matches = append(matches, kind)
		}
	}
	if len(matches) == 1 {
		return matches[0], false
	}
	if len(matches) > 1 {
		return matches[0], true
	}
	return "product", true
}

func emptyDecisionEffects() map[string]any {
	return map[string]any{"requiredTaskTypes": []any{}, "requiredGates": []any{}, "requiredEvidence": []any{}, "requiredWriteScopes": []any{}, "sharedResources": []any{}}
}

func stringsToAny(items []string) []any {
	out := make([]any, len(items))
	for i, item := range items {
		out[i] = item
	}
	return out
}
