package workflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/engineeringsystem"
	"github.com/JonatasFreireDev/spec-framework/internal/validator"
)

type Artifact struct {
	ID              string   `json:"id"`
	Type            string   `json:"type"`
	Status          string   `json:"status"`
	Path            string   `json:"path"`
	ParentIDs       []string `json:"parentIds"`
	Maturity        string   `json:"maturity,omitempty"`
	TargetFeature   string   `json:"targetFeature,omitempty"`
	ApprovalAdapter string   `json:"approval_adapter,omitempty"`
}
type ArtifactRequirement struct {
	Type string `json:"type"`
	Path string `json:"path"`
}
type Registry struct {
	Artifacts         []Artifact            `json:"artifacts"`
	Modules           []string              `json:"modules,omitempty"`
	RequiredArtifacts []ArtifactRequirement `json:"required_artifacts,omitempty"`
}
type Workspace struct {
	ID               string            `json:"id"`
	Scope            map[string]string `json:"scope"`
	CurrentStep      string            `json:"current_step"`
	RecommendedSkill string            `json:"recommended_skill"`
	BlockedBy        []string          `json:"blocked_by"`
	CreatedBy        string            `json:"created_by"`
	CreatedAt        string            `json:"created_at"`
}
type Status struct {
	Workspace Workspace
	Artifact  Artifact
	Next      string
	Blockers  []string
}
type Approval struct {
	ArtifactID    string `json:"artifact_id"`
	Path          string `json:"path"`
	ContentHash   string `json:"content_hash"`
	StatusGranted string `json:"status_granted"`
	ApprovedBy    string `json:"approved_by"`
	ApprovedAt    string `json:"approved_at"`
	Notes         string `json:"notes"`
}
type ApprovalPreview struct {
	Artifact           Artifact
	Grant              string
	CurrentHash        string
	ParentBlockers     []string
	ValidationBlockers []string
}
type Claim struct {
	TaskID    string `json:"task_id"`
	Graph     string `json:"graph"`
	Agent     string `json:"agent"`
	ClaimedAt string `json:"claimed_at"`
}
type Claims struct {
	Claims []Claim `json:"claims"`
}

func LoadRegistry(root string) (Registry, error) {
	var r Registry
	err := readJSON(filepath.Join(root, ".product", "artifacts.json"), &r)
	return r, err
}

func CreateWorkspace(root, selector, domain, goal, useCase, createdBy string) (Workspace, error) {
	r, err := LoadRegistry(root)
	if err != nil {
		return Workspace{}, err
	}
	a, err := resolveFeature(r, selector, domain, goal)
	if err != nil {
		return Workspace{}, err
	}
	if err := requireStartingPointApproval(root, r, a); err != nil {
		return Workspace{}, err
	}
	if hasImplementationReadyUseCase(r, a.ID) {
		missing, gateErr := GateReadiness(root)
		if gateErr != nil {
			return Workspace{}, fmt.Errorf("implementation-ready use case requires configured implementation gates: %w", gateErr)
		}
		if len(missing) > 0 {
			return Workspace{}, fmt.Errorf("implementation-ready use case requires configured implementation gates: %s", strings.Join(missing, ", "))
		}
	}
	dir := filepath.Join(root, ".product", "workspaces")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Workspace{}, err
	}
	id, err := nextID(dir, "WORK-")
	if err != nil {
		return Workspace{}, err
	}
	scope := map[string]string{"feature": a.Path}
	if strings.TrimSpace(useCase) != "" {
		scope["use_case"] = filepath.ToSlash(strings.TrimSpace(useCase))
	}
	for _, x := range r.Artifacts {
		if contains(a.ParentIDs, x.ID) {
			if x.Type == "user-goal" {
				scope["goal"] = x.Path
			}
			if x.Type == "domain" {
				scope["domain"] = x.Path
			}
		}
	}
	next, blockers := nextFor(root, a, scope["use_case"])
	w := Workspace{ID: id, Scope: scope, CurrentStep: next, RecommendedSkill: next, BlockedBy: blockers, CreatedBy: createdBy, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	wd := filepath.Join(dir, id)
	if err := os.MkdirAll(wd, 0755); err != nil {
		return Workspace{}, err
	}
	for _, child := range []string{"handoffs", "checkpoints", "command-plans", "evidence", "tasks"} {
		if err := os.MkdirAll(filepath.Join(wd, child), 0755); err != nil {
			return Workspace{}, err
		}
	}
	if err := writeJSON(filepath.Join(wd, "workspace.json"), w); err != nil {
		return Workspace{}, err
	}
	_ = writeJSON(filepath.Join(wd, "state.json"), RuntimeState{Version: RuntimeVersion, WorkspaceID: id, Phase: next, Status: "active", UpdatedAt: time.Now().UTC().Format(time.RFC3339), Blockers: blockers})
	return w, nil
}

func hasImplementationReadyUseCase(registry Registry, featureID string) bool {
	for _, artifact := range registry.Artifacts {
		if artifact.Type == "use-case" && artifact.Maturity == "implementation-ready" && contains(artifact.ParentIDs, featureID) {
			return true
		}
	}
	return false
}

func requireStartingPointApproval(root string, registry Registry, selectedFeature Artifact) error {
	var manifest map[string]any
	if err := readJSON(filepath.Join(root, ".product", "framework.json"), &manifest); err != nil {
		return err
	}
	point, _ := manifest["starting_point"].(string)
	if point == "existing-documents" {
		return requireLatestImportMaterialized(root, manifest)
	}
	requirements := registry.RequiredArtifacts
	if len(requirements) == 0 {
		requirements = legacyStartingPointRequirements(point)
	}
	if len(requirements) == 0 {
		return nil
	}
	for _, required := range requirements {
		found := false
		for _, artifact := range registry.Artifacts {
			if artifact.Type != required.Type || (required.Path != "" && filepath.ToSlash(artifact.Path) != filepath.ToSlash(required.Path)) {
				continue
			}
			found = true
			path := filepath.Join(root, filepath.FromSlash(artifact.Path))
			if !isApproved(artifact.Status) || !hasCurrentApproval(root, path, artifact.Status) {
				return fmt.Errorf("%s lacks current approval evidence: %s", strings.ReplaceAll(required.Type, "-", " "), artifact.Path)
			}
			if required.Type == "feature-brief" {
				target, err := featureBriefTarget(path)
				if err != nil {
					return err
				}
				if artifact.TargetFeature != target {
					return fmt.Errorf("feature brief registry target %q differs from approved artifact target %q; rebuild the registry", artifact.TargetFeature, target)
				}
				if !featureTargetMatches(target, selectedFeature) {
					return fmt.Errorf("feature brief target %q does not match selected feature %s (%s)", target, selectedFeature.ID, selectedFeature.Path)
				}
			}
			break
		}
		if !found {
			return fmt.Errorf("%s requires a registered %s", point, required.Path)
		}
	}
	return nil
}

func legacyStartingPointRequirements(point string) []ArtifactRequirement {
	switch point {
	case "existing-feature":
		return []ArtifactRequirement{{Type: "feature-brief", Path: "foundation/feature-brief.md"}}
	case "existing-product":
		return []ArtifactRequirement{{Type: "product-baseline", Path: "foundation/product-baseline.md"}, {Type: "strategy", Path: "foundation/strategy/strategy.md"}}
	case "existing-implementation":
		return []ArtifactRequirement{{Type: "implementation-assessment", Path: "knowledge/assessments/implementation-assessment.md"}, {Type: "problem", Path: "foundation/problem/problem.md"}, {Type: "vision", Path: "foundation/vision/vision.md"}, {Type: "product-principles", Path: "foundation/vision/principles.md"}, {Type: "north-star", Path: "foundation/vision/north-star.md"}, {Type: "strategy", Path: "foundation/strategy/strategy.md"}}
	default:
		return nil
	}
}

func featureBriefTarget(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Split(line, "|")
		if len(parts) < 3 || !strings.EqualFold(strings.TrimSpace(parts[1]), "Target Feature") {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[2]), "`"), nil
	}
	return "", errors.New("feature brief requires a Target Feature field")
}

func featureTargetMatches(target string, feature Artifact) bool {
	target = filepath.ToSlash(strings.TrimSpace(strings.Trim(target, "`")))
	if target == "" || strings.Contains(strings.ToUpper(target), "TBD") || strings.Contains(target, "[") {
		return false
	}
	path := filepath.ToSlash(feature.Path)
	if strings.EqualFold(target, feature.ID) || strings.EqualFold(target, path) {
		return true
	}
	if strings.HasSuffix(path, "/context.md") {
		return strings.EqualFold(target, strings.TrimSuffix(path, "/context.md")+"/feature.md")
	}
	if strings.HasSuffix(path, "/feature.md") {
		return strings.EqualFold(target, strings.TrimSuffix(path, "/feature.md")+"/context.md")
	}
	return false
}

func requireLatestImportMaterialized(root string, manifest map[string]any) error {
	importConfig, _ := manifest["import"].(map[string]any)
	runID, _ := importConfig["latest_run"].(string)
	if strings.TrimSpace(runID) == "" {
		return errors.New("existing-documents requires a latest import run in .product/framework.json")
	}
	planPath := filepath.Join(root, "knowledge", "imports", "runs", runID, "import-plan.json")
	var plan map[string]any
	if err := readJSON(planPath, &plan); err != nil {
		return fmt.Errorf("read latest import plan %s: %w", runID, err)
	}
	status, _ := plan["status"].(string)
	approved, _ := plan["materialization_approved"].(bool)
	approvedBy, _ := plan["materialization_approved_by"].(string)
	approvedAt, _ := plan["materialization_approved_at"].(string)
	paths, _ := plan["materialized_paths"].([]any)
	if status != "materialized" || !approved || strings.TrimSpace(approvedBy) == "" || strings.TrimSpace(approvedAt) == "" || len(paths) == 0 {
		return fmt.Errorf("latest import run %s is not materially complete; review mappings and run spec-framework import materialize", runID)
	}
	var mapping struct {
		Mappings []struct {
			Target       string `json:"target"`
			Selected     bool   `json:"selected"`
			DraftContent string `json:"draft_content"`
		} `json:"mappings"`
	}
	if err := readJSON(filepath.Join(root, "knowledge", "imports", "runs", runID, "mapping.json"), &mapping); err != nil {
		return fmt.Errorf("read latest import mapping %s: %w", runID, err)
	}
	reported := map[string]bool{}
	for _, raw := range paths {
		path, _ := raw.(string)
		if path == "" {
			return fmt.Errorf("latest import run %s has an invalid materialized path", runID)
		}
		reported[filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))] = true
	}
	selected := map[string]string{}
	for _, item := range mapping.Mappings {
		if !item.Selected {
			continue
		}
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(item.Target)))
		if clean == ".." || strings.HasPrefix(clean, "../") || filepath.IsAbs(filepath.FromSlash(item.Target)) {
			return fmt.Errorf("latest import run %s selected target escapes product root: %s", runID, item.Target)
		}
		selected[clean] = item.DraftContent
	}
	if len(selected) == 0 || len(selected) != len(reported) {
		return fmt.Errorf("latest import run %s materialized paths do not match selected mappings", runID)
	}
	recordedHashes := map[string]string{}
	if raw, ok := plan["materialized_hashes"].(map[string]any); ok {
		for path, value := range raw {
			recordedHashes[filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))] = fmt.Sprint(value)
		}
	}
	normalizedHashes := map[string]string{}
	if raw, ok := plan["normalized_hashes"].(map[string]any); ok {
		for path, value := range raw {
			normalizedHashes[filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))] = fmt.Sprint(value)
		}
	}
	for path, expected := range selected {
		if !reported[path] {
			return fmt.Errorf("latest import run %s did not report selected target %s", runID, path)
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			return fmt.Errorf("latest import run %s materialized target %s is missing: %w", runID, path, err)
		}
		if normalized, exists := normalizedHashes[path]; exists {
			_ = normalized // The artifact may legitimately change status after normalization.
			continue
		}
		expectedHash := Hash(expected)
		if recorded, exists := recordedHashes[path]; exists {
			if recorded != expectedHash {
				return fmt.Errorf("latest import run %s recorded hash for %s does not match the approved draft", runID, path)
			}
		} else if string(data) != expected {
			return fmt.Errorf("latest import run %s legacy target %s differs from the approved draft", runID, path)
		}
	}
	return nil
}
func Features(root string) ([]Artifact, error) {
	r, err := LoadRegistry(root)
	if err != nil {
		return nil, err
	}
	var out []Artifact
	for _, a := range r.Artifacts {
		if a.Type == "feature" {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

func PreviewApproval(root, artifactPath, grant string) (ApprovalPreview, error) {
	r, err := LoadRegistry(root)
	if err != nil {
		return ApprovalPreview{}, err
	}
	rel, err := filepath.Rel(root, artifactPath)
	if err != nil {
		return ApprovalPreview{}, err
	}
	rel = filepath.ToSlash(rel)
	var a Artifact
	for _, item := range r.Artifacts {
		if filepath.ToSlash(item.Path) == rel {
			a = item
			break
		}
	}
	if a.ID == "" {
		return ApprovalPreview{}, fmt.Errorf("artifact is not registered: %s", rel)
	}
	if !validTransition(a.Status, grant) {
		return ApprovalPreview{}, fmt.Errorf("invalid status transition %s -> %s", a.Status, grant)
	}
	var blockers []string
	if grant != "rejected" && grant != "draft" && grant != "proposed" {
		for _, pid := range a.ParentIDs {
			for _, p := range r.Artifacts {
				if p.ID == pid && (!isApproved(p.Status) || !hasCurrentApproval(root, filepath.Join(root, filepath.FromSlash(p.Path)), p.Status)) {
					blockers = append(blockers, fmt.Sprintf("parent %s lacks current approval evidence", pid))
				}
			}
		}
	}
	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return ApprovalPreview{}, err
	}
	updated, err := setStatus(string(data), grant)
	if err != nil {
		return ApprovalPreview{}, err
	}
	validationBlockers := make([]string, 0)
	if grant != "rejected" && grant != "draft" && grant != "proposed" {
		validation, err := validator.ValidateCandidate(context.Background(), root, root, artifactPath, updated, grant)
		if err != nil {
			return ApprovalPreview{}, err
		}
		for _, diagnostic := range validation.Diagnostics {
			if filepath.ToSlash(diagnostic.File) == rel && (diagnostic.Check == "template-conformance" || diagnostic.Check == "import-provenance" || diagnostic.Check == "status-coherence" || diagnostic.Check == "qa-evidence" || diagnostic.Check == "delivery") && diagnostic.Severity == validator.Error {
				validationBlockers = append(validationBlockers, fmt.Sprintf("%s: %s", diagnostic.Check, diagnostic.Message))
			}
		}
	}
	updates, err := approvalUpdatesForArtifact(root, a, artifactPath, updated, grant)
	if err != nil {
		return ApprovalPreview{}, err
	}
	currentHash, err := approvalHash(root, a, artifactPath, updates)
	if err != nil {
		return ApprovalPreview{}, err
	}
	return ApprovalPreview{Artifact: a, Grant: grant, CurrentHash: currentHash, ParentBlockers: blockers, ValidationBlockers: validationBlockers}, nil
}

func WorkspaceStatus(root, id string) (Status, error) {
	var w Workspace
	if loaded, err := LoadWorkspace(root, id); err != nil {
		return Status{}, err
	} else {
		w = loaded
	}
	r, err := LoadRegistry(root)
	if err != nil {
		return Status{}, err
	}
	path := w.Scope["feature"]
	var a Artifact
	for _, x := range r.Artifacts {
		if filepath.ToSlash(x.Path) == filepath.ToSlash(path) {
			a = x
			break
		}
	}
	if a.ID == "" {
		return Status{}, fmt.Errorf("workspace feature not found in registry: %s", path)
	}
	next, blockers := nextFor(root, a, w.Scope["use_case"])
	previous := w.CurrentStep
	w.CurrentStep = next
	w.RecommendedSkill = next
	w.BlockedBy = blockers
	if err := writeJSON(filepath.Join(workspaceDir(root, id), "workspace.json"), w); err != nil {
		return Status{}, err
	}
	state := RuntimeState{Version: RuntimeVersion, WorkspaceID: id, Phase: next, Status: "active", UpdatedAt: time.Now().UTC().Format(time.RFC3339), Blockers: blockers}
	if len(blockers) > 0 {
		state.Status = "blocked"
	}
	if err := writeJSON(filepath.Join(workspaceDir(root, id), "state.json"), state); err != nil {
		return Status{}, err
	}
	if previous != "" && previous != next {
		baseCommit, _ := gitOutput(filepath.Dir(root), "rev-parse", "HEAD")
		input := Hash(a.Path + "\n" + a.Status)
		_, _ = WriteCheckpoint(root, id, next, strings.TrimSpace(baseCommit), input, Hash(strings.Join(blockers, "\n")))
		_, _ = WriteHandoff(root, id, previous, next, "Workflow advanced after gate evaluation.")
	}
	return Status{Workspace: w, Artifact: a, Next: next, Blockers: blockers}, nil
}

func nextFor(root string, a Artifact, useCaseSelector string) (string, []string) {
	if a.Status != "approved" && a.Status != "in_progress" && a.Status != "implemented" && a.Status != "validated" && a.Status != "released" {
		return "feature", []string{"feature is not approved"}
	}
	base := filepath.Dir(filepath.Join(root, filepath.FromSlash(a.Path)))
	ucRoot := filepath.Join(base, "use-cases")
	entries, _ := os.ReadDir(ucRoot)
	var useCases []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), "_") {
			useCases = append(useCases, filepath.Join(ucRoot, e.Name()))
		}
	}
	sort.Strings(useCases)
	if len(useCases) == 0 {
		return "use-case", nil
	}
	if useCaseSelector == "" && len(useCases) > 1 {
		return "use-case", []string{"multiple use cases exist; select one with --use-case"}
	}
	uc := useCases[0]
	if useCaseSelector != "" {
		found := ""
		for _, candidate := range useCases {
			rel, _ := filepath.Rel(root, candidate)
			if filepath.Base(candidate) == useCaseSelector || filepath.ToSlash(rel) == filepath.ToSlash(useCaseSelector) {
				found = candidate
			}
		}
		if found == "" {
			return "use-case", []string{"selected use case was not found"}
		}
		uc = found
	}
	if !approvedDocument(root, filepath.Join(uc, "context.md")) {
		return "use-case", []string{"use case is not approved"}
	}
	if !approvedDocument(root, filepath.Join(uc, "specification.md")) {
		return "specification", []string{"specification is missing or not approved"}
	}
	for _, contract := range requiredContracts(filepath.Join(uc, "context.md")) {
		path := filepath.Join(uc, "contracts", contract+".md")
		if !approvedDocument(root, path) && !notApplicableDocument(path) {
			return "specification", []string{"required " + contract + " contract is missing or not approved"}
		}
	}
	if !approvedDocument(root, filepath.Join(uc, "design.md")) && !notApplicableDocument(filepath.Join(uc, "design.md")) {
		return "ux-ui", []string{"design is missing, not approved, or lacks structured not_applicable status and rationale"}
	}
	if !approvedDocument(root, filepath.Join(uc, "technical-discovery.md")) {
		return "technical-discovery", []string{"technical discovery is missing or not approved"}
	}
	if !architectureResolved(root, filepath.Join(uc, "technical-discovery.md")) {
		return "product-historian", []string{"Architecture Gate is unresolved"}
	}
	if engineeringReviewApplies(filepath.Join(uc, "context.md")) {
		if !approvedDocument(root, filepath.Join(uc, "engineering-proposal.md")) {
			return "engineering-proposal", []string{"applicable engineering proposal is missing or not approved"}
		}
		if !passedEngineeringReview(root, filepath.Join(uc, "engineering-review.md"), filepath.Join(uc, "engineering-proposal.md")) {
			return "engineering-review", []string{"applicable engineering review is missing, not approved, or not passed"}
		}
	}
	if !approvedDocument(root, filepath.Join(uc, "implementation-plan.md")) {
		return "implementation-planner", []string{"implementation plan is missing or not approved"}
	}
	graph := filepath.Join(uc, "execution-graph.json")
	graphState := jsonStatus(graph)
	switch graphState {
	case "":
		return "execution-graph", []string{"execution graph is missing"}
	case "draft":
		return "execution-graph", []string{"execution graph is still draft"}
	case "proposed":
		return "task-generator", nil
	case "materialized":
		return "task-generator", []string{"materialized graph and tasks require final approval"}
	default:
		if !isApproved(graphState) {
			return "execution-graph", []string{"execution graph has unsupported status " + graphState}
		}
	}
	tasks, _ := os.ReadDir(filepath.Join(uc, "tasks"))
	if len(tasks) == 0 {
		return "task-generator", []string{"task files are missing"}
	}
	for _, task := range tasks {
		if task.IsDir() || filepath.Ext(task.Name()) != ".md" {
			continue
		}
		if !approvedDocument(root, filepath.Join(uc, "tasks", task.Name())) {
			return "task-generator", []string{"task " + task.Name() + " is not approved"}
		}
	}
	if missing, err := GateReadiness(root); err != nil || len(missing) > 0 {
		return "code-runner", []string{"implementation gates are not configured"}
	}
	return "code-runner", nil
}

func approvedDocument(root, path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	status := extractStatus(string(data))
	return isApproved(status) && hasCurrentApproval(root, path, status)
}
func notApplicableDocument(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	text := string(data)
	if extractStatus(text) != "not_applicable" {
		return false
	}
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?mi)^\s*Rationale:\s*(.+)$`),
		regexp.MustCompile(`(?mi)^\|\s*Rationale\s*\|\s*` + "`?" + `([^|` + "`" + `]+)` + "`?" + `\s*\|$`),
	}
	for _, pattern := range patterns {
		if match := pattern.FindStringSubmatch(text); len(match) == 2 && meaningful(match[1]) {
			return true
		}
	}
	return false
}
func architectureResolved(root, path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	text := string(data)
	verdict, decision, rationale := architectureGateFields(text)
	if verdict == "Not required" {
		return meaningful(rationale)
	}
	if verdict != "Decision required" {
		return false
	}
	id := regexp.MustCompile(`\bDEC-\d+\b`).FindString(decision)
	if id == "" {
		return false
	}
	var index map[string]any
	if readJSON(filepath.Join(root, ".product", "decisions.json"), &index) != nil {
		return false
	}
	for _, item := range stringMapSlice(index["decisions"]) {
		if fmt.Sprint(item["id"]) != id || fmt.Sprint(item["status"]) != "approved" {
			continue
		}
		decisionPath := filepath.Join(root, filepath.FromSlash(fmt.Sprint(item["path"])))
		if !hasCurrentApproval(root, decisionPath, "approved") || !workflowDecisionApplies(item, path) {
			return false
		}
		return true
	}
	return false
}

func architectureGateFields(text string) (string, string, string) {
	field := func(name string) string {
		pattern := regexp.MustCompile(`(?mi)^\|\s*` + regexp.QuoteMeta(name) + `\s*\|\s*` + "`?" + `([^|` + "`" + `]+)` + "`?" + `\s*\|$`)
		match := pattern.FindStringSubmatch(text)
		if len(match) == 2 {
			return strings.TrimSpace(match[1])
		}
		return ""
	}
	verdict, decision, rationale := field("Verdict"), field("Decision"), field("Rationale")
	if verdict != "" {
		return verdict, decision, rationale
	}
	pattern := regexp.MustCompile(`(?mi)^\|\s*(Decision required|Not required)\s*\|\s*([^|]*)\|\s*([^|]*)\|$`)
	match := pattern.FindStringSubmatch(text)
	if len(match) == 4 {
		return strings.TrimSpace(match[1]), strings.TrimSpace(match[2]), strings.TrimSpace(match[3])
	}
	return "", "", ""
}

func workflowDecisionApplies(decision map[string]any, path string) bool {
	path = filepath.ToSlash(path)
	var scopes []string
	if scope, ok := decision["scope"].(string); ok {
		scopes = append(scopes, strings.Split(scope, "/")...)
	}
	scopes = append(scopes, stringAnySlice(decision["affectedArtifacts"])...)
	for _, scope := range scopes {
		scope = filepath.ToSlash(strings.TrimSpace(scope))
		if scope == "architecture" || scope == "security" || scope == "data" || scope == "product" || strings.HasPrefix(path, strings.TrimSuffix(scope, "/")) || strings.HasPrefix(scope, filepath.ToSlash(filepath.Dir(path))) {
			return true
		}
	}
	return false
}

func meaningful(value string) bool {
	value = strings.ToLower(strings.Trim(strings.TrimSpace(value), "`[]"))
	return value != "" && value != "n/a" && value != "none" && value != "tbd" && value != "pending" && !strings.Contains(value, "placeholder")
}
func tierLDocument(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	text := strings.ToUpper(string(data))
	return strings.Contains(text, "RIGOR_TIER: L") || strings.Contains(text, "| L |")
}
func engineeringReviewApplies(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	triggers, _ := engineeringsystem.Triggers(string(data))
	return tierLDocument(path) || len(triggers) > 0
}
func passedEngineeringReview(root, path, proposalPath string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	text := string(data)
	if !approvedDocument(root, path) || !regexp.MustCompile(`(?mi)^\|\s*Verdict\s*\|\s*`+"`?"+`passed`+"`?"+`\s*\|`).MatchString(text) {
		return false
	}
	proposal, err := os.ReadFile(proposalPath)
	if err != nil {
		return false
	}
	match := regexp.MustCompile(`(?mi)^\|\s*Proposal hash\s*\|\s*` + "`?" + `([a-f0-9]{64})` + "`?" + `\s*\|`).FindStringSubmatch(text)
	return len(match) == 2 && match[1] == Hash(string(proposal))
}
func approvedJSON(path string) bool {
	var raw map[string]any
	if readJSON(path, &raw) != nil {
		return false
	}
	return isApproved(fmt.Sprint(raw["status"]))
}
func jsonStatus(path string) string {
	var raw map[string]any
	if readJSON(path, &raw) != nil {
		return ""
	}
	return strings.ToLower(fmt.Sprint(raw["status"]))
}
func extractStatus(text string) string {
	re := regexp.MustCompile(`(?mi)^\s*status:\s*` + "`?" + `([a-z_]+)` + "`?")
	if m := re.FindStringSubmatch(text); len(m) > 1 {
		return strings.ToLower(m[1])
	}
	re = regexp.MustCompile(`(?mi)^\|\s*Status\s*\|\s*` + "`?" + `([a-z_]+)` + "`?")
	if m := re.FindStringSubmatch(text); len(m) > 1 {
		return strings.ToLower(m[1])
	}
	return ""
}
func requiredContracts(contextPath string) []string {
	data, _ := os.ReadFile(contextPath)
	text := strings.ToUpper(string(data))
	tier := "S"
	if strings.Contains(text, "RIGOR_TIER: L") || strings.Contains(text, "| L |") {
		tier = "L"
	} else if strings.Contains(text, "RIGOR_TIER: M") || strings.Contains(text, "| M |") {
		tier = "M"
	}
	out := []string{"behavior", "quality"}
	if tier == "M" || tier == "L" {
		out = append(out, "product", "ux", "api", "data", "rollout")
	}
	if tier == "L" {
		out = append(out, "security", "observability")
	}
	return out
}

func Approve(root, artifactPath, grant, approvedBy, notes string) (Approval, error) {
	grant = strings.TrimSpace(grant)
	if !map[string]bool{"draft": true, "proposed": true, "approved": true, "rejected": true, "in_progress": true, "implemented": true, "validated": true, "released": true}[grant] {
		return Approval{}, fmt.Errorf("unsupported grant %q", grant)
	}
	if grant == "rejected" && strings.TrimSpace(notes) == "" {
		return Approval{}, errors.New("rejection requires notes explaining what must be revised")
	}
	preview, err := PreviewApproval(root, artifactPath, grant)
	if err != nil {
		return Approval{}, err
	}
	if grant != "rejected" && grant != "draft" && grant != "proposed" && (len(preview.ParentBlockers) > 0 || len(preview.ValidationBlockers) > 0) {
		return Approval{}, errors.New(strings.Join(append(preview.ParentBlockers, preview.ValidationBlockers...), "; "))
	}
	r, err := LoadRegistry(root)
	if err != nil {
		return Approval{}, err
	}
	rel, err := filepath.Rel(root, artifactPath)
	if err != nil {
		return Approval{}, err
	}
	rel = filepath.ToSlash(rel)
	var idx = -1
	for i, a := range r.Artifacts {
		if filepath.ToSlash(a.Path) == rel {
			idx = i
			break
		}
	}
	if idx < 0 {
		return Approval{}, fmt.Errorf("artifact is not registered: %s", rel)
	}
	for _, pid := range r.Artifacts[idx].ParentIDs {
		if grant == "rejected" || grant == "draft" || grant == "proposed" {
			break
		}
		for _, p := range r.Artifacts {
			if p.ID == pid && (!isApproved(p.Status) || !hasCurrentApproval(root, filepath.Join(root, filepath.FromSlash(p.Path)), p.Status)) {
				return Approval{}, fmt.Errorf("parent %s lacks current approval evidence", pid)
			}
		}
	}
	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return Approval{}, err
	}
	updated, err := setStatus(string(data), grant)
	if err != nil {
		return Approval{}, err
	}
	registryPath := filepath.Join(root, ".product", "artifacts.json")
	registryData, err := os.ReadFile(registryPath)
	if err != nil {
		return Approval{}, err
	}
	updates, err := approvalUpdatesForArtifact(root, r.Artifacts[idx], artifactPath, updated, grant)
	if err != nil {
		return Approval{}, err
	}
	backups := map[string][]byte{}
	for path := range updates {
		backups[path], err = os.ReadFile(path)
		if err != nil {
			return Approval{}, err
		}
	}
	rollback := func() {
		for path, backup := range backups {
			_ = atomicWrite(path, backup)
		}
		_ = atomicWrite(registryPath, registryData)
	}
	for path, content := range updates {
		if err := atomicWrite(path, content); err != nil {
			rollback()
			return Approval{}, err
		}
	}
	r.Artifacts[idx].Status = grant
	if err := writeJSON(registryPath, r); err != nil {
		rollback()
		return Approval{}, err
	}
	contentHash, err := approvalHash(root, r.Artifacts[idx], artifactPath, updates)
	if err != nil {
		rollback()
		return Approval{}, err
	}
	rec := Approval{ArtifactID: r.Artifacts[idx].ID, Path: rel, ContentHash: contentHash, StatusGranted: grant, ApprovedBy: approvedBy, ApprovedAt: time.Now().UTC().Format(time.RFC3339), Notes: notes}
	dir := filepath.Join(root, ".product", "history")
	if err := os.MkdirAll(dir, 0755); err != nil {
		rollback()
		return Approval{}, err
	}
	name := fmt.Sprintf("approval-%s-%s-%d.json", strings.ToLower(rec.ArtifactID), grant, time.Now().UTC().UnixNano())
	if err := writeJSON(filepath.Join(dir, name), rec); err != nil {
		rollback()
		return Approval{}, err
	}
	return rec, nil
}

func approvalUpdatesForArtifact(root string, artifact Artifact, artifactPath, updated, status string) (map[string][]byte, error) {
	updates := map[string][]byte{artifactPath: []byte(updated)}
	adapter := artifact.ApprovalAdapter
	if adapter == "" {
		if map[string]bool{"problem": true, "vision": true, "product-principles": true, "north-star": true, "strategy": true}[artifact.Type] {
			adapter = "foundation-context"
		} else if artifact.Type == "engineering-system" {
			adapter = "engineering-system"
		}
	}
	if adapter == "foundation-context" {
		contextPath := filepath.Join(filepath.Dir(artifactPath), "context.md")
		contextData, err := os.ReadFile(contextPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil {
			contextUpdated, updateErr := setStatus(string(contextData), status)
			if updateErr != nil {
				return nil, updateErr
			}
			updates[contextPath] = []byte(contextUpdated)
		}
	}
	if adapter != "engineering-system" {
		return updates, nil
	}
	inspection, inspectErr := engineeringsystem.Inspect(root)
	if inspectErr != nil {
		return nil, inspectErr
	}
	if len(inspection.Blockers) > 0 {
		return nil, fmt.Errorf("Engineering System is not approvable: %s", strings.Join(inspection.Blockers, "; "))
	}
	contextPath := filepath.Join(root, "engineering", "context.md")
	catalogPath := filepath.Join(root, "engineering", "engineering-system.yaml")
	contextData, err := os.ReadFile(contextPath)
	if err != nil {
		return nil, err
	}
	catalogData, err := os.ReadFile(catalogPath)
	if err != nil {
		return nil, err
	}
	contextUpdated, catalogUpdated, err := engineeringsystem.SynchronizeStatus(string(contextData), catalogData, status)
	if err != nil {
		return nil, err
	}
	updates[contextPath] = []byte(contextUpdated)
	updates[catalogPath] = catalogUpdated
	qualityMarkdownPath := filepath.Join(root, "engineering", "quality", "quality-system.md")
	qualityCatalogPath := filepath.Join(root, "engineering", "quality", "quality-system.yaml")
	if qualityMarkdown, readErr := os.ReadFile(qualityMarkdownPath); readErr == nil {
		qualityMarkdownUpdated, updateErr := setStatus(string(qualityMarkdown), status)
		if updateErr != nil {
			return nil, updateErr
		}
		qualityCatalog, readCatalogErr := os.ReadFile(qualityCatalogPath)
		if readCatalogErr != nil {
			return nil, readCatalogErr
		}
		qualityCatalogUpdated, updateCatalogErr := engineeringsystem.SynchronizeQualityStatus(qualityCatalog, status)
		if updateCatalogErr != nil {
			return nil, updateCatalogErr
		}
		updates[qualityMarkdownPath] = []byte(qualityMarkdownUpdated)
		updates[qualityCatalogPath] = qualityCatalogUpdated
	} else if !os.IsNotExist(readErr) {
		return nil, readErr
	}
	return updates, nil
}

func approvalHash(root string, artifact Artifact, artifactPath string, updates map[string][]byte) (string, error) {
	adapter := artifact.ApprovalAdapter
	if adapter == "" && (artifact.Type == "engineering-system" || filepath.Base(artifactPath) == "engineering-system.md") {
		adapter = "engineering-system"
	}
	if adapter != "engineering-system" {
		if updates != nil {
			if content, exists := updates[artifactPath]; exists {
				return Hash(string(content)), nil
			}
		}
		data, err := os.ReadFile(artifactPath)
		if err != nil {
			return "", err
		}
		return Hash(string(data)), nil
	}
	relative := map[string][]byte{}
	for path, content := range updates {
		rel, _ := filepath.Rel(root, path)
		relative[filepath.ToSlash(rel)] = content
	}
	return engineeringsystem.CompositeHash(root, relative)
}

func GateReadiness(root string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, "knowledge", "conventions", "gates.md"))
	if err != nil {
		return nil, err
	}
	var missing []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "| `GATE-") && strings.Contains(strings.ToUpper(line), "TBD") {
			parts := strings.Split(line, "|")
			if len(parts) > 1 {
				missing = append(missing, strings.Trim(strings.TrimSpace(parts[1]), "`"))
			}
		}
	}
	sort.Strings(missing)
	return missing, nil
}

type Graph struct {
	ID    string `json:"id"`
	Nodes []Node `json:"nodes"`
}
type Node struct {
	ID              string   `json:"id"`
	Type            string   `json:"type"`
	Path            string   `json:"path"`
	DependsOn       []string `json:"dependsOn"`
	Status          string   `json:"status"`
	WriteScope      []string `json:"writeScope"`
	SharedResources []string `json:"sharedResources"`
}

func Ready(graphPath string) ([]Node, error) {
	var g Graph
	if err := readJSON(graphPath, &g); err != nil {
		return nil, err
	}
	done := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Status == "complete" || n.Status == "validated" {
			done[n.ID] = true
		}
	}
	var out []Node
	for _, n := range g.Nodes {
		if n.Status != "" && n.Status != "pending" && n.Status != "ready" {
			continue
		}
		ok := true
		for _, d := range n.DependsOn {
			if !done[d] {
				ok = false
			}
		}
		if ok {
			out = append(out, n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func ReadyUnclaimed(root, graphPath string) ([]Node, error) {
	nodes, err := Ready(graphPath)
	if err != nil {
		return nil, err
	}
	var out []Node
	for _, n := range nodes {
		if _, err := os.Stat(filepath.Join(root, ".product", "claims", n.ID+".lock")); os.IsNotExist(err) {
			out = append(out, n)
		}
	}
	return out, nil
}
func ClaimTask(root, graphPath, taskID, agent string) (Claim, error) {
	if err := validateRuntimeComponent(taskID, "task"); err != nil {
		return Claim{}, err
	}
	if err := validateRuntimeComponent(agent, "agent"); err != nil {
		return Claim{}, err
	}
	ready, err := Ready(graphPath)
	if err != nil {
		return Claim{}, err
	}
	allowed := false
	for _, n := range ready {
		if n.ID == taskID {
			allowed = true
		}
	}
	if !allowed {
		return Claim{}, fmt.Errorf("task %s is not ready", taskID)
	}
	lockDir := filepath.Join(root, ".product", "claims")
	if err := os.MkdirAll(lockDir, 0755); err != nil {
		return Claim{}, err
	}
	lockPath := filepath.Join(lockDir, taskID+".lock")
	lock, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return Claim{}, fmt.Errorf("task %s is already claimed", taskID)
		}
		return Claim{}, err
	}
	_ = lock.Close()
	locked := true
	defer func() {
		if locked {
			_ = os.Remove(lockPath)
		}
	}()
	unlock, err := acquireClaimLock(root)
	if err != nil {
		return Claim{}, err
	}
	defer unlock()
	path := filepath.Join(root, ".product", "claims.json")
	var all Claims
	if err := readJSON(path, &all); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Claim{}, err
	}
	for _, c := range all.Claims {
		if c.TaskID == taskID {
			return Claim{}, fmt.Errorf("task %s is already claimed by %s", taskID, c.Agent)
		}
	}
	var graph Graph
	if err := readJSON(graphPath, &graph); err != nil {
		return Claim{}, err
	}
	nodes := map[string]Node{}
	for _, n := range graph.Nodes {
		nodes[n.ID] = n
	}
	target := nodes[taskID]
	relGraph, _ := filepath.Rel(root, graphPath)
	for _, existing := range all.Claims {
		if filepath.ToSlash(existing.Graph) != filepath.ToSlash(relGraph) {
			continue
		}
		other := nodes[existing.TaskID]
		if scopesOverlap(target.WriteScope, other.WriteScope) || resourcesOverlap(target.SharedResources, other.SharedResources) {
			return Claim{}, fmt.Errorf("task %s conflicts with claimed task %s", taskID, existing.TaskID)
		}
	}
	rel, _ := filepath.Rel(root, graphPath)
	c := Claim{TaskID: taskID, Graph: filepath.ToSlash(rel), Agent: agent, ClaimedAt: time.Now().UTC().Format(time.RFC3339)}
	all.Claims = append(all.Claims, c)
	if err := writeJSON(path, all); err != nil {
		return Claim{}, err
	}
	locked = false
	return c, nil
}
func ReleaseClaim(root, taskID, agent string) error {
	if err := validateRuntimeComponent(taskID, "task"); err != nil {
		return err
	}
	if err := validateRuntimeComponent(agent, "agent"); err != nil {
		return err
	}
	unlock, err := acquireClaimLock(root)
	if err != nil {
		return err
	}
	defer unlock()
	path := filepath.Join(root, ".product", "claims.json")
	var all Claims
	if err := readJSON(path, &all); err != nil {
		return err
	}
	var kept []Claim
	found := false
	for _, c := range all.Claims {
		if c.TaskID == taskID && c.Agent == agent {
			found = true
		} else {
			kept = append(kept, c)
		}
	}
	if !found {
		return fmt.Errorf("task %s is not claimed by %s", taskID, agent)
	}
	all.Claims = kept
	if err := writeJSON(path, all); err != nil {
		return err
	}
	_ = os.Remove(filepath.Join(root, ".product", "claims", taskID+".lock"))
	return nil
}
func acquireClaimLock(root string) (func(), error) {
	path := filepath.Join(root, ".product", "claims.lock")
	for range 200 {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err == nil {
			_ = f.Close()
			return func() { _ = os.Remove(path) }, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		time.Sleep(5 * time.Millisecond)
	}
	return nil, fmt.Errorf("timed out waiting for claims lock")
}
func Complete(root, graphPath, taskID, agent string) error {
	var claims Claims
	if err := readJSON(filepath.Join(root, ".product", "claims.json"), &claims); err != nil {
		return err
	}
	owned := false
	for _, c := range claims.Claims {
		if c.TaskID == taskID && c.Agent == agent {
			owned = true
		}
	}
	if !owned {
		return fmt.Errorf("task %s is not claimed by %s", taskID, agent)
	}
	original, err := os.ReadFile(graphPath)
	if err != nil {
		return err
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(strings.TrimPrefix(string(original), "\ufeff")), &raw); err != nil {
		return err
	}
	nodes, _ := raw["nodes"].([]any)
	found := false
	for _, item := range nodes {
		if node, ok := item.(map[string]any); ok && fmt.Sprint(node["id"]) == taskID {
			node["status"] = "complete"
			found = true
		}
	}
	if !found {
		return fmt.Errorf("task %s not found", taskID)
	}
	raw["updatedAt"] = time.Now().UTC().Format("2006-01-02")
	if err := writeJSON(graphPath, raw); err != nil {
		return err
	}
	if err := ReleaseClaim(root, taskID, agent); err != nil {
		_ = atomicWrite(graphPath, original)
		return err
	}
	return nil
}

func resolveFeature(r Registry, selector, domain, goal string) (Artifact, error) {
	selector = filepath.ToSlash(strings.TrimSpace(selector))
	domain = filepath.ToSlash(strings.TrimSpace(domain))
	goal = filepath.ToSlash(strings.TrimSpace(goal))
	var matches []Artifact
	for _, a := range r.Artifacts {
		if a.Type != "feature" {
			continue
		}
		matched := a.ID == selector || a.Path == selector || strings.TrimSuffix(a.Path, "/context.md") == strings.TrimSuffix(selector, "/")
		if !matched {
			continue
		}
		path := filepath.ToSlash(a.Path)
		if domain != "" && !strings.Contains(path, "domains/"+strings.Trim(domain, "/")+"/") {
			continue
		}
		if goal != "" && !strings.Contains(path, "/goals/"+strings.Trim(goal, "/")+"/") {
			continue
		}
		matches = append(matches, a)
	}
	if len(matches) == 0 {
		return Artifact{}, fmt.Errorf("feature not found: %s", selector)
	}
	if len(matches) > 1 {
		return Artifact{}, fmt.Errorf("feature selector is ambiguous; add --domain and --goal")
	}
	return matches[0], nil
}
func setStatus(text, status string) (string, error) {
	re := regexp.MustCompile(`(?m)^(\s*status:\s*)[^\r\n]+`)
	if re.MatchString(text) {
		return re.ReplaceAllString(text, "${1}"+status), nil
	}
	re = regexp.MustCompile(`(?mi)^(\|\s*Status\s*\|\s*)[^|]+(\|)`)
	if re.MatchString(text) {
		return re.ReplaceAllString(text, "${1}`"+status+"` ${2}"), nil
	}
	return "", fmt.Errorf("artifact has no editable status field")
}
func Hash(text string) string {
	normalized := strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(normalized, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return hex.EncodeToString(sum[:])
}
func nextID(dir, prefix string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	max := 0
	for _, e := range entries {
		var n int
		if _, err := fmt.Sscanf(e.Name(), prefix+"%03d", &n); err == nil && n > max {
			max = n
		}
	}
	return fmt.Sprintf("%s%03d", prefix, max+1), nil
}
func contains(items []string, value string) bool {
	for _, x := range items {
		if x == value {
			return true
		}
	}
	return false
}
func scopesOverlap(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			x = filepath.ToSlash(strings.TrimSpace(x))
			y = filepath.ToSlash(strings.TrimSpace(y))
			if x != "" && y != "" && (x == y || strings.HasPrefix(x, y+"/") || strings.HasPrefix(y, x+"/")) {
				return true
			}
		}
	}
	return false
}
func resourcesOverlap(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y && x != "" && !strings.EqualFold(x, "N/A") && !strings.EqualFold(x, "none") {
				return true
			}
		}
	}
	return false
}
func isApproved(s string) bool {
	return map[string]bool{"approved": true, "in_progress": true, "implemented": true, "validated": true, "released": true}[s]
}
func validTransition(from, to string) bool {
	if from == to {
		return true
	}
	allowed := map[string][]string{"draft": {"proposed", "approved", "rejected"}, "proposed": {"approved", "rejected"}, "materialized": {"approved", "rejected"}, "approved": {"in_progress", "rejected"}, "in_progress": {"implemented"}, "implemented": {"validated"}, "validated": {"released"}, "rejected": {"draft", "proposed", "approved"}}
	return contains(allowed[from], to)
}
func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(strings.TrimPrefix(string(data), "\ufeff")), value)
}
func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(path, append(data, '\n'))
}
func atomicWrite(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
