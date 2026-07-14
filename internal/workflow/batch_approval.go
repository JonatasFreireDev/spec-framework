package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BatchScope describes a user-visible approval selection. Exactly one of
// Artifact, IDs, Foundation, Stage, or AllEligible must be selected. Until
// narrows AllEligible to the ordered delivery stage sequence.
type BatchScope struct {
	Artifact    []string
	IDs         []string
	Foundation  bool
	Stage       string
	Until       string
	AllEligible bool
}

type BatchItem struct {
	Artifact Artifact `json:"artifact"`
	Action   string   `json:"action"`
	Reason   string   `json:"reason,omitempty"`
	Hash     string   `json:"hash,omitempty"`
}

type BatchPlan struct {
	Grant     string      `json:"grant"`
	Scope     BatchScope  `json:"scope"`
	Items     []BatchItem `json:"items"`
	ToApprove []Artifact  `json:"to_approve"`
	Ignored   []BatchItem `json:"ignored"`
	Blockers  []BatchItem `json:"blockers"`
	NextGate  string      `json:"next_gate,omitempty"`
}

var approvalStageOrder = []string{"foundation", "domains", "feature", "use-cases", "specification", "design", "engineering", "planning", "tasks"}

var approvalStageTypes = map[string][]string{
	"foundation":    {"problem", "vision", "product-principles", "north-star", "strategy", "product-baseline", "feature-brief", "implementation-assessment"},
	"domains":       {"domain", "user-goal"},
	"feature":       {"feature"},
	"use-cases":     {"use-case"},
	"specification": {"specification", "specification-contract"},
	"design":        {"design", "design-system"},
	"engineering":   {"technical-discovery", "engineering-proposal", "engineering-review", "engineering-system"},
	"planning":      {"implementation-plan"},
	"tasks":         {"execution-graph", "task", "taskset"},
}

func BuildBatchApprovalPlan(root string, scope BatchScope, grant string) (BatchPlan, error) {
	grant = strings.TrimSpace(grant)
	if grant != "approved" {
		return BatchPlan{}, fmt.Errorf("approve-batch only supports --grant approved")
	}
	selected, err := selectBatchArtifacts(root, scope)
	if err != nil {
		return BatchPlan{}, err
	}
	ordered := orderBatchArtifacts(selected)
	plan := BatchPlan{Grant: grant, Scope: scope}
	plannedApproved := map[string]bool{}
	for _, a := range ordered {
		item := BatchItem{Artifact: a}
		if stale, reason := batchArtifactStale(root, a); stale {
			item.Action, item.Reason = "blocked", reason
			plan.Blockers = append(plan.Blockers, item)
			continue
		}
		if a.Status == "approved" {
			path := filepath.Join(root, filepath.FromSlash(a.Path))
			if !hasCurrentApproval(root, path, "approved") {
				item.Action, item.Reason = "blocked", "artifact is marked approved but lacks current approval evidence"
				plan.Blockers = append(plan.Blockers, item)
			} else {
				item.Action, item.Reason = "ignored", "already approved; current approval evidence is not re-created"
				plan.Ignored = append(plan.Ignored, item)
			}
			continue
		}
		preview, previewErr := PreviewApproval(root, filepath.Join(root, filepath.FromSlash(a.Path)), grant)
		if previewErr != nil {
			item.Action, item.Reason = "blocked", previewErr.Error()
			plan.Blockers = append(plan.Blockers, item)
			continue
		}
		blockers := make([]string, 0, len(preview.ParentBlockers))
		for _, blocker := range preview.ParentBlockers {
			parentID := parentIDFromBlocker(blocker)
			if parentID == "" || !plannedApproved[parentID] {
				blockers = append(blockers, blocker)
			}
		}
		if len(blockers) > 0 {
			item.Action, item.Reason = "blocked", strings.Join(blockers, "; ")
			plan.Blockers = append(plan.Blockers, item)
			continue
		}
		item.Action, item.Hash = "approve", preview.CurrentHash
		plan.Items = append(plan.Items, item)
		plan.ToApprove = append(plan.ToApprove, a)
		plannedApproved[a.ID] = true
	}
	if len(plan.ToApprove) == 0 && len(plan.Blockers) > 0 {
		return plan, nil
	}
	plan.NextGate = nextApprovalGate(scope)
	return plan, nil
}

func ApproveBatch(root string, plan BatchPlan, approvedBy, notes string) ([]Approval, error) {
	if strings.TrimSpace(approvedBy) == "" {
		return nil, fmt.Errorf("approve-batch requires an approver identity")
	}
	if len(plan.Blockers) > 0 {
		return nil, fmt.Errorf("approval batch has %d blocker(s)", len(plan.Blockers))
	}
	snapshot, err := snapshotApprovalRoot(root)
	if err != nil {
		return nil, err
	}
	var records []Approval
	for _, a := range plan.ToApprove {
		record, err := Approve(root, filepath.Join(root, filepath.FromSlash(a.Path)), plan.Grant, approvedBy, notes)
		if err != nil {
			if rollbackErr := snapshot.restore(root); rollbackErr != nil {
				return records, fmt.Errorf("approval batch stopped after %d artifact(s): %w; rollback failed: %v", len(records), err, rollbackErr)
			}
			return nil, fmt.Errorf("approval batch rolled back after %d artifact(s): %w", len(records), err)
		}
		records = append(records, record)
	}
	return records, nil
}

type approvalSnapshot struct{ Files map[string][]byte }

func snapshotApprovalRoot(root string) (approvalSnapshot, error) {
	s := approvalSnapshot{Files: map[string][]byte{}}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		s.Files[filepath.ToSlash(rel)] = data
		return nil
	})
	return s, err
}

func (s approvalSnapshot) restore(root string) error {
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if _, ok := s.Files[filepath.ToSlash(rel)]; !ok {
			return os.Remove(path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for rel, data := range s.Files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, filepath.FromSlash(rel))), 0o755); err != nil {
			return err
		}
		if err := atomicWrite(filepath.Join(root, filepath.FromSlash(rel)), data); err != nil {
			return err
		}
	}
	return nil
}

func selectBatchArtifacts(root string, scope BatchScope) ([]Artifact, error) {
	selectors := 0
	if len(scope.Artifact) > 0 {
		selectors++
	}
	if len(scope.IDs) > 0 {
		selectors++
	}
	if scope.Foundation {
		selectors++
	}
	if scope.Stage != "" {
		selectors++
	}
	if scope.AllEligible {
		selectors++
	}
	if selectors != 1 {
		return nil, fmt.Errorf("approve-batch requires exactly one scope: --artifact, --ids, --foundation, --stage, or --all-eligible")
	}
	if scope.Until != "" && !scope.AllEligible {
		return nil, fmt.Errorf("--until requires --all-eligible")
	}
	r, err := LoadRegistry(root)
	if err != nil {
		return nil, err
	}
	byID := map[string]Artifact{}
	for _, a := range r.Artifacts {
		byID[a.ID] = a
	}
	selected := make([]Artifact, 0)
	seen := map[string]bool{}
	add := func(a Artifact) {
		if !seen[a.ID] {
			seen[a.ID] = true
			selected = append(selected, a)
		}
	}
	if len(scope.Artifact) > 0 {
		for _, value := range scope.Artifact {
			value = strings.TrimSpace(value)
			found := false
			for _, a := range r.Artifacts {
				if a.ID == value || filepath.ToSlash(a.Path) == filepath.ToSlash(value) {
					add(a)
					found = true
				}
			}
			if !found {
				return nil, fmt.Errorf("artifact not found in registry: %s", value)
			}
		}
	} else if len(scope.IDs) > 0 {
		for _, id := range scope.IDs {
			id = strings.TrimSpace(id)
			a, ok := byID[id]
			if !ok {
				return nil, fmt.Errorf("artifact id not found in registry: %s", id)
			}
			add(a)
		}
	} else if scope.Foundation {
		for _, a := range r.Artifacts {
			if batchContains(approvalStageTypes["foundation"], normalizedType(a.Type)) {
				add(a)
			}
		}
	} else if scope.Stage != "" {
		stage := normalizeApprovalStage(scope.Stage)
		if stage == "" {
			return nil, fmt.Errorf("unknown approval stage: %s", scope.Stage)
		}
		for _, a := range r.Artifacts {
			if batchContains(approvalStageTypes[stage], normalizedType(a.Type)) {
				add(a)
			}
		}
	} else {
		if scope.Until == "" {
			for _, a := range r.Artifacts {
				add(a)
			}
		} else {
			idx := approvalStageIndex(scope.Until)
			if idx < 0 {
				return nil, fmt.Errorf("unknown approval stage: %s", scope.Until)
			}
			for _, a := range r.Artifacts {
				stageIndex := artifactApprovalStage(a)
				if stageIndex >= 0 && stageIndex <= idx {
					add(a)
				}
			}
		}
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("approval scope selected no registered artifacts")
	}
	return selected, nil
}

func orderBatchArtifacts(items []Artifact) []Artifact {
	byID := map[string]Artifact{}
	for _, a := range items {
		byID[a.ID] = a
	}
	out := make([]Artifact, 0, len(items))
	remaining := map[string]bool{}
	for _, a := range items {
		remaining[a.ID] = true
	}
	for len(remaining) > 0 {
		var ready []Artifact
		for id := range remaining {
			a := byID[id]
			ok := true
			for _, parent := range a.ParentIDs {
				if remaining[parent] {
					ok = false
					break
				}
			}
			if ok {
				ready = append(ready, a)
			}
		}
		if len(ready) == 0 { // Cycles are reported by the validator; keep preview deterministic.
			for id := range remaining {
				ready = append(ready, byID[id])
			}
		}
		sort.Slice(ready, func(i, j int) bool { return ready[i].Path < ready[j].Path })
		for _, a := range ready {
			out = append(out, a)
			delete(remaining, a.ID)
		}
	}
	return out
}

func batchArtifactStale(root string, artifact Artifact) (bool, string) {
	var raw struct {
		Derivations []struct {
			ArtifactID string `json:"artifact_id"`
			Path       string `json:"path"`
			Sources    []struct {
				Path        string `json:"path"`
				ContentHash string `json:"content_hash"`
			} `json:"derived_from"`
		} `json:"derivations"`
	}
	data, err := os.ReadFile(filepath.Join(root, ".product", "derivations.json"))
	if err != nil || json.Unmarshal(data, &raw) != nil {
		return false, ""
	}
	for _, d := range raw.Derivations {
		if d.ArtifactID != artifact.ID && filepath.ToSlash(d.Path) != filepath.ToSlash(artifact.Path) {
			continue
		}
		for _, source := range d.Sources {
			content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(source.Path)))
			if err != nil || Hash(string(content)) != source.ContentHash {
				return true, "artifact is stale because a derived source changed or is missing"
			}
		}
	}
	return false, ""
}

func parentIDFromBlocker(blocker string) string {
	const prefix = "parent "
	if !strings.HasPrefix(blocker, prefix) {
		return ""
	}
	value := strings.TrimPrefix(blocker, prefix)
	if idx := strings.Index(value, " "); idx >= 0 {
		value = value[:idx]
	}
	return value
}
func normalizedType(value string) string { return strings.ReplaceAll(strings.ToLower(value), "_", "-") }
func batchContains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}
func normalizeApprovalStage(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "usecase" {
		value = "use-cases"
	}
	if value == "domain" {
		value = "domains"
	}
	return value
}
func approvalStageIndex(value string) int {
	value = normalizeApprovalStage(value)
	for i, item := range approvalStageOrder {
		if item == value {
			return i
		}
	}
	return -1
}
func artifactApprovalStage(a Artifact) int {
	kind := normalizedType(a.Type)
	for i, stage := range approvalStageOrder {
		if batchContains(approvalStageTypes[stage], kind) {
			return i
		}
	}
	return -1
}
func nextApprovalGate(scope BatchScope) string {
	if scope.Until != "" {
		idx := approvalStageIndex(scope.Until)
		if idx >= 0 && idx+1 < len(approvalStageOrder) {
			return approvalStageOrder[idx+1]
		}
	}
	if scope.Stage != "" {
		idx := approvalStageIndex(scope.Stage)
		if idx >= 0 && idx+1 < len(approvalStageOrder) {
			return approvalStageOrder[idx+1]
		}
	}
	return "next eligible gate"
}
