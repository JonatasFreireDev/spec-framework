// Package dispatch persists supervised subagent assignments. It never starts a
// process, grants approval, or performs delivery operations.
package dispatch

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/sourceimport"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

type Envelope struct {
	Version            int      `json:"version"`
	ID                 string   `json:"id"`
	WorkspaceID        string   `json:"workspace_id"`
	TaskID             string   `json:"task_id"`
	UnitKind           string   `json:"unit_kind"`
	UnitPath           string   `json:"unit_path,omitempty"`
	ImportRun          string   `json:"import_run,omitempty"`
	ImportChunk        string   `json:"import_chunk,omitempty"`
	Role               string   `json:"role"`
	Agent              string   `json:"agent"`
	Graph              string   `json:"graph"`
	TaskPath           string   `json:"task_path"`
	InputHash          string   `json:"input_hash"`
	DiffHash           string   `json:"diff_hash,omitempty"`
	ParentID           string   `json:"parent_id,omitempty"`
	Dependencies       []string `json:"dependencies,omitempty"`
	Phase              int      `json:"phase,omitempty"`
	ContextPolicy      string   `json:"context_policy,omitempty"`
	RequiredReading    []string `json:"required_reading"`
	WriteScope         []string `json:"write_scope"`
	ExpectedEvidence   []string `json:"expected_evidence"`
	Status             string   `json:"status"`
	CreatedAt          string   `json:"created_at"`
	ReturnedAt         string   `json:"returned_at,omitempty"`
	Summary            string   `json:"summary,omitempty"`
	Evidence           []string `json:"evidence,omitempty"`
	OutputHashes       []string `json:"output_hashes,omitempty"`
	Blockers           []string `json:"blockers,omitempty"`
	DecisionCandidates []string `json:"decision_candidates,omitempty"`
	Forbidden          []string `json:"forbidden"`
}
type Candidate struct {
	TaskID, Role, Path string
	WriteScope         []string
	Ready              bool
	Blockers           []string
}
type Finding struct {
	Kind       string `json:"kind"`
	DispatchID string `json:"dispatch_id"`
	Detail     string `json:"detail"`
	Owner      string `json:"owner"`
}
type Transcript struct {
	DispatchID string `json:"dispatch_id"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	ExitCode   int    `json:"exit_code"`
	OutputHash string `json:"output_hash"`
}
type WaveResult struct {
	ID         string      `json:"id"`
	Transcript *Transcript `json:"transcript,omitempty"`
	Error      string      `json:"error,omitempty"`
}
type Recommendation struct {
	Kind                 string `json:"kind"`
	Detail               string `json:"detail"`
	RequiresConfirmation bool   `json:"requires_confirmation"`
}
type Config struct {
	Version             int      `json:"version"`
	Enabled             bool     `json:"enabled"`
	Harnesses           []string `json:"harnesses"`
	MaxParallel         int      `json:"max_parallel"`
	TranscriptRetention int      `json:"transcript_retention"`
}

type engineeringExecution struct {
	Mode          string `json:"mode"`
	ContextPolicy string `json:"context_policy"`
	MaxParallel   int    `json:"max_parallel"`
	Fallback      string `json:"fallback"`
}

type engineeringRoute struct {
	Skill      string   `json:"skill"`
	Phase      int      `json:"phase"`
	DependsOn  []string `json:"depends_on"`
	WriteScope []string `json:"write_scope"`
	Status     string   `json:"status"`
}

type engineeringHandoff struct {
	SchemaVersion int                  `json:"schema_version"`
	Execution     engineeringExecution `json:"execution"`
	Routes        []engineeringRoute   `json:"routes"`
}

var engineeringRoles = map[string]bool{
	"technical-landscape":   true,
	"engineering-standards": true,
	"operations-baseline":   true,
	"engineering-evidence":  true,
	"engineering-system":    true,
}

var engineeringWriteScopes = map[string][]string{
	"technical-landscape":   {"engineering/architecture", "engineering/catalog"},
	"engineering-standards": {"engineering/standards"},
	"operations-baseline":   {"engineering/operations"},
	"engineering-evidence":  {"engineering/evidence"},
	"engineering-system":    {"engineering/engineering-system.md", "engineering/engineering-system.yaml", "engineering/quality"},
}

func dir(root, work string) string {
	return filepath.Join(root, ".product", "workspaces", work, "dispatches")
}
func configPath(root string) string { return filepath.Join(root, ".product", "dispatch.json") }
func LoadConfig(root string) (Config, error) {
	data, err := os.ReadFile(configPath(root))
	if os.IsNotExist(err) {
		return Config{Version: 1, Enabled: false, MaxParallel: 1, TranscriptRetention: 100}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}
	if c.MaxParallel < 1 {
		c.MaxParallel = 1
	}
	return c, nil
}
func SaveConfig(root string, c Config) error {
	if c.Version == 0 {
		c.Version = 1
	}
	if c.MaxParallel < 1 {
		c.MaxParallel = 1
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(configPath(root)), 0755); err != nil {
		return err
	}
	return os.WriteFile(configPath(root), append(data, '\n'), 0644)
}

// AssignEngineering persists one harness-native engineering specialist
// assignment. The CLI validates the bounded context and write scope but never
// starts a subagent itself.
func AssignEngineering(root, work, handoffPath, role, agent string, dependencies []string) (Envelope, error) {
	if !engineeringRoles[role] {
		return Envelope{}, errors.New("unsupported engineering specialist role")
	}
	if strings.TrimSpace(work) == "" || strings.TrimSpace(agent) == "" || strings.TrimSpace(handoffPath) == "" {
		return Envelope{}, errors.New("engineering assignment requires work, handoff path, role, and agent")
	}
	if err := validateEngineeringHandoffPath(root, work, handoffPath); err != nil {
		return Envelope{}, err
	}
	data, err := os.ReadFile(handoffPath)
	if err != nil {
		return Envelope{}, err
	}
	var handoff engineeringHandoff
	if err := json.Unmarshal(data, &handoff); err != nil {
		return Envelope{}, fmt.Errorf("read engineering handoff: %w", err)
	}
	if handoff.SchemaVersion != 1 || handoff.Execution.Mode != "delegated" {
		return Envelope{}, errors.New("engineering handoff must use schema_version 1 and delegated execution")
	}
	if handoff.Execution.ContextPolicy != "minimal" {
		return Envelope{}, errors.New("delegated engineering context_policy must be minimal")
	}
	if handoff.Execution.MaxParallel < 1 {
		return Envelope{}, errors.New("delegated engineering max_parallel must be at least 1")
	}
	var route *engineeringRoute
	for index := range handoff.Routes {
		if handoff.Routes[index].Skill == role {
			route = &handoff.Routes[index]
			break
		}
	}
	if route == nil {
		return Envelope{}, errors.New("engineering role is not present in the handoff routes")
	}
	if route.Status != "pending" && route.Status != "ready" {
		return Envelope{}, errors.New("engineering route is not assignable")
	}
	if route.Phase < 1 || len(route.WriteScope) == 0 {
		return Envelope{}, errors.New("engineering route has no write scope")
	}
	for _, scope := range route.WriteScope {
		if !insideAnyScope(scope, engineeringWriteScopes[role]) {
			return Envelope{}, fmt.Errorf("engineering route write scope %s is not owned by %s", scope, role)
		}
	}
	if err := os.MkdirAll(dir(root, work), 0755); err != nil {
		return Envelope{}, err
	}
	release, err := acquireEngineeringAssignmentLock(root, work)
	if err != nil {
		return Envelope{}, err
	}
	defer release()
	resolved, err := validateEngineeringDependencies(root, work, route.DependsOn, dependencies)
	if err != nil {
		return Envelope{}, err
	}
	active, duplicate, err := activeEngineeringAssignments(root, work, role)
	if err != nil {
		return Envelope{}, err
	}
	if duplicate {
		return Envelope{}, errors.New("engineering specialist role already has an active assignment")
	}
	if active >= handoff.Execution.MaxParallel {
		return Envelope{}, errors.New("delegated engineering assignments reached max_parallel")
	}
	sum := sha256.Sum256(data)
	id := fmt.Sprintf("DISPATCH-%d", time.Now().UTC().UnixNano())
	e := Envelope{
		Version: 1, ID: id, WorkspaceID: work, UnitKind: "engineering-specialist", UnitPath: filepath.ToSlash(handoffPath),
		Role: role, Agent: agent, InputHash: hex.EncodeToString(sum[:]), Dependencies: resolved, Phase: route.Phase,
		ContextPolicy: "minimal", RequiredReading: []string{filepath.ToSlash(handoffPath), "FRAMEWORK.md", "skills/" + role + "/SKILL.md"},
		WriteScope: route.WriteScope, ExpectedEvidence: []string{"summary", "evidence", "output hashes", "blockers"},
		Status: "assigned", CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Forbidden: []string{"approval", "product-decision", "application-code", "commit", "push", "merge", "release", "write-outside-scope"},
	}
	return e, write(filepath.Join(dir(root, work), id+".json"), e)
}

func validateEngineeringDependencies(root, work string, expectedRoles, ids []string) ([]string, error) {
	if len(expectedRoles) != len(ids) {
		return nil, errors.New("engineering assignment dependencies do not match the handoff route")
	}
	roles := map[string]bool{}
	resolved := append([]string(nil), ids...)
	for _, id := range resolved {
		if !validDispatchID(id) {
			return nil, errors.New("engineering dependency has invalid dispatch id")
		}
		var dependency Envelope
		if err := read(filepath.Join(dir(root, work), id+".json"), &dependency); err != nil {
			return nil, fmt.Errorf("read engineering dependency %s: %w", id, err)
		}
		if dependency.UnitKind != "engineering-specialist" || dependency.Status != "returned" {
			return nil, fmt.Errorf("engineering dependency %s is not returned", id)
		}
		if len(dependency.Blockers) > 0 {
			return nil, fmt.Errorf("engineering dependency %s returned blockers", id)
		}
		if err := verifyEngineeringOutputs(root, dependency); err != nil {
			return nil, fmt.Errorf("engineering dependency %s is stale: %w", id, err)
		}
		roles[dependency.Role] = true
	}
	for _, role := range expectedRoles {
		if !roles[role] {
			return nil, fmt.Errorf("engineering dependency role %s is missing", role)
		}
	}
	sort.Strings(resolved)
	return resolved, nil
}

func activeEngineeringAssignments(root, work, role string) (int, bool, error) {
	xs, err := Observe(root, work)
	if err != nil {
		return 0, false, err
	}
	active := 0
	duplicate := false
	for _, item := range xs {
		if item.UnitKind == "engineering-specialist" && item.Status == "assigned" {
			active++
			if item.Role == role {
				duplicate = true
			}
		}
	}
	return active, duplicate, nil
}

func acquireEngineeringAssignmentLock(root, work string) (func(), error) {
	lock := filepath.Join(dir(root, work), ".engineering-assign.lock")
	for attempt := 0; attempt < 50; attempt++ {
		if err := os.Mkdir(lock, 0755); err == nil {
			return func() { _ = os.Remove(lock) }, nil
		} else if !os.IsExist(err) {
			return nil, err
		}
		if info, err := os.Stat(lock); err == nil && time.Since(info.ModTime()) > 30*time.Second {
			_ = os.Remove(lock)
			continue
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil, errors.New("engineering assignment lock is busy")
}

// ReturnEngineering records a compact specialist result after verifying that
// every declared output is inside the assignment scope and still matches its
// SHA-256 hash. It does not grant approval or merge specialist contracts.
func ReturnEngineering(root, work, id, agent, summary string, evidence, outputHashes, blockers, decisionCandidates []string) (Envelope, error) {
	var e Envelope
	if !validDispatchID(id) {
		return e, errors.New("engineering return has invalid dispatch id")
	}
	release, err := acquireEngineeringAssignmentLock(root, work)
	if err != nil {
		return e, err
	}
	defer release()
	path := filepath.Join(dir(root, work), id+".json")
	if err := read(path, &e); err != nil {
		return e, err
	}
	if e.UnitKind != "engineering-specialist" || e.Agent != agent || e.Status != "assigned" {
		return e, errors.New("dispatch is not an assigned engineering specialist for this agent")
	}
	if strings.TrimSpace(summary) == "" || len(evidence) == 0 || len(outputHashes) == 0 {
		return e, errors.New("engineering return requires summary, evidence, and output hashes")
	}
	currentInput, err := os.ReadFile(filepath.FromSlash(e.UnitPath))
	if err != nil {
		return e, fmt.Errorf("read current engineering handoff: %w", err)
	}
	inputSum := sha256.Sum256(currentInput)
	if e.InputHash != hex.EncodeToString(inputSum[:]) {
		return e, errors.New("engineering assignment input handoff is stale")
	}
	for _, item := range outputHashes {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 || !insideAnyScope(parts[0], e.WriteScope) {
			return e, fmt.Errorf("engineering output %s escapes or omits its write scope", item)
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(parts[0])))
		if err != nil {
			return e, fmt.Errorf("read engineering output %s: %w", parts[0], err)
		}
		sum := sha256.Sum256(data)
		if !strings.EqualFold(parts[1], hex.EncodeToString(sum[:])) {
			return e, fmt.Errorf("engineering output hash mismatch for %s", parts[0])
		}
	}
	e.Status = "returned"
	e.Summary = summary
	e.Evidence = evidence
	e.OutputHashes = append([]string(nil), outputHashes...)
	e.Blockers = append([]string(nil), blockers...)
	e.DecisionCandidates = append([]string(nil), decisionCandidates...)
	e.ReturnedAt = time.Now().UTC().Format(time.RFC3339)
	return e, write(path, e)
}

func insideAnyScope(path string, scopes []string) bool {
	path = filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
	if path == "." || path == ".." || strings.HasPrefix(path, "../") || filepath.IsAbs(filepath.FromSlash(path)) {
		return false
	}
	for _, scope := range scopes {
		scope = filepath.ToSlash(strings.TrimSuffix(filepath.Clean(filepath.FromSlash(scope)), "/"))
		if path == scope || strings.HasPrefix(path, scope+"/") {
			return true
		}
	}
	return false
}

func validateEngineeringHandoffPath(root, work, handoffPath string) error {
	workspace := filepath.Join(root, ".product", "workspaces", work)
	relative, err := filepath.Rel(workspace, handoffPath)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return errors.New("engineering handoff must be inside its product workspace")
	}
	if strings.ToLower(filepath.Ext(relative)) != ".json" {
		return errors.New("engineering handoff must be a JSON file")
	}
	return nil
}

func validDispatchID(id string) bool {
	return strings.HasPrefix(id, "DISPATCH-") && filepath.Base(id) == id && !strings.ContainsAny(id, `/\\`)
}

func verifyEngineeringOutputs(root string, e Envelope) error {
	if len(e.OutputHashes) == 0 {
		return errors.New("returned dependency has no output hashes")
	}
	for _, output := range e.OutputHashes {
		parts := strings.SplitN(output, "=", 2)
		if len(parts) != 2 || !insideAnyScope(parts[0], e.WriteScope) {
			return errors.New("returned dependency output is malformed or outside scope")
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(parts[0])))
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		if !strings.EqualFold(parts[1], hex.EncodeToString(sum[:])) {
			return errors.New("returned dependency output hash mismatch")
		}
	}
	return nil
}
func Plan(root, graph string) ([]Candidate, error) {
	nodes, err := workflow.ReadyUnclaimed(root, graph)
	if err != nil {
		return nil, err
	}
	items := make([]Candidate, 0, len(nodes))
	for _, node := range nodes {
		readiness, err := workflow.CheckTaskReadiness(root, graph, node.ID)
		item := Candidate{TaskID: node.ID, Role: "code-runner", Path: node.Path, WriteScope: node.WriteScope, Ready: err == nil && readiness.Ready}
		if err != nil {
			item.Blockers = []string{err.Error()}
		} else {
			for _, check := range readiness.Checks {
				if check.Status == "block" {
					item.Blockers = append(item.Blockers, check.Detail)
				}
			}
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].TaskID < items[j].TaskID })
	return items, nil
}
func Assign(root, work, graph, task, role, agent string) (Envelope, error) {
	if role == "" {
		role = "code-runner"
	}
	if role != "code-runner" {
		return Envelope{}, errors.New("only code-runner assignments are enabled; independent review dispatch remains read-only")
	}
	readiness, err := workflow.CheckTaskReadiness(root, graph, task)
	if err != nil {
		return Envelope{}, err
	}
	if !readiness.Ready {
		return Envelope{}, errors.New("task is not ready for dispatch")
	}
	if _, err = workflow.ClaimLease(root, graph, task, agent, 30*time.Minute); err != nil {
		return Envelope{}, err
	}
	taskPath := ""
	for _, node := range mustNodes(graph) {
		if node.ID == task {
			taskPath = filepath.Join(filepath.Dir(graph), filepath.FromSlash(node.Path))
			break
		}
	}
	data, err := os.ReadFile(taskPath)
	if err != nil {
		_ = workflow.ReleaseLease(root, task, agent)
		return Envelope{}, err
	}
	sum := sha256.Sum256(data)
	if err = os.MkdirAll(dir(root, work), 0755); err != nil {
		return Envelope{}, err
	}
	id := fmt.Sprintf("DISPATCH-%d", time.Now().UTC().UnixNano())
	e := Envelope{Version: 1, ID: id, WorkspaceID: work, TaskID: task, UnitKind: "task", UnitPath: filepath.ToSlash(taskPath), Role: role, Agent: agent, Graph: filepath.ToSlash(graph), TaskPath: filepath.ToSlash(taskPath), InputHash: hex.EncodeToString(sum[:]), RequiredReading: []string{filepath.ToSlash(taskPath), filepath.ToSlash(graph)}, ExpectedEvidence: []string{"diff hash", "test log"}, Status: "assigned", CreatedAt: time.Now().UTC().Format(time.RFC3339), Forbidden: []string{"approval", "commit", "push", "merge", "release", "review-resolution"}}
	for _, node := range mustNodes(graph) {
		if node.ID == task {
			e.WriteScope = node.WriteScope
		}
	}
	return e, write(filepath.Join(dir(root, work), id+".json"), e)
}
func Return(root, work, id, agent, summary, diffHash string, evidence []string) (Envelope, error) {
	var e Envelope
	path := filepath.Join(dir(root, work), id+".json")
	if err := read(path, &e); err != nil {
		return e, err
	}
	if e.Agent != agent || e.Status != "assigned" {
		return e, errors.New("dispatch is not assigned to this agent")
	}
	if strings.TrimSpace(summary) == "" || strings.TrimSpace(diffHash) == "" || len(evidence) == 0 {
		return e, errors.New("return requires summary, diff hash, and evidence")
	}
	e.Status = "returned"
	e.Summary = summary
	e.Evidence = evidence
	e.DiffHash = diffHash
	e.ReturnedAt = time.Now().UTC().Format(time.RFC3339)
	if err := write(path, e); err != nil {
		return e, err
	}
	if e.TaskID == "" {
		return e, nil
	}
	return e, workflow.ReleaseLease(root, e.TaskID, agent)
}
func ReturnImport(root, work, id, agent, summary string, review sourceimport.ChunkReview) (Envelope, error) {
	var e Envelope
	path := filepath.Join(dir(root, work), id+".json")
	if err := read(path, &e); err != nil {
		return e, err
	}
	if e.UnitKind != "import-chunk" || e.Agent != agent || e.Status != "assigned" {
		return e, errors.New("dispatch is not an assigned import chunk for this agent")
	}
	if strings.TrimSpace(summary) == "" {
		return e, errors.New("import return requires summary")
	}
	if err := sourceimport.RecordChunkReview(root, e.ImportRun, e.ImportChunk, agent, review); err != nil {
		return e, err
	}
	e.Status = "returned"
	e.Summary = summary
	e.Evidence = []string{"structured source evidence"}
	e.ReturnedAt = time.Now().UTC().Format(time.RFC3339)
	return e, write(path, e)
}

// AssignReview creates an independent read-only review envelope for the exact
// diff returned by a code-runner. It never claims write ownership.
func AssignReview(root, work, parentID, role, agent string) (Envelope, error) {
	if role != "qa" && role != "code-review" && role != "security-review" {
		return Envelope{}, errors.New("review role must be qa, code-review, or security-review")
	}
	var parent Envelope
	if err := read(filepath.Join(dir(root, work), parentID+".json"), &parent); err != nil {
		return Envelope{}, err
	}
	if parent.Role != "code-runner" || parent.Status != "returned" || parent.DiffHash == "" {
		return Envelope{}, errors.New("review requires a returned code-runner dispatch with diff hash")
	}
	id := fmt.Sprintf("DISPATCH-%d", time.Now().UTC().UnixNano())
	e := parent
	e.ID = id
	e.Role = role
	e.Agent = agent
	e.ParentID = parent.ID
	e.Status = "assigned"
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	e.ReturnedAt = ""
	e.Summary = ""
	e.Evidence = nil
	e.WriteScope = nil
	e.Forbidden = []string{"approval", "code-change", "commit", "push", "merge", "release", "review-resolution"}
	e.ExpectedEvidence = []string{"review verdict", "findings", "diff hash: " + parent.DiffHash}
	return e, write(filepath.Join(dir(root, work), id+".json"), e)
}

// AssignResearch creates a read-only Technical Discovery assignment. Results
// are proposals with cited evidence, never decisions or Engineering Proposals.
func AssignResearch(root, work, unitPath, agent string) (Envelope, error) {
	if strings.TrimSpace(unitPath) == "" || strings.TrimSpace(agent) == "" {
		return Envelope{}, errors.New("research dispatch requires unit path and agent")
	}
	data, err := os.ReadFile(unitPath)
	if err != nil {
		return Envelope{}, err
	}
	sum := sha256.Sum256(data)
	if err := os.MkdirAll(dir(root, work), 0755); err != nil {
		return Envelope{}, err
	}
	id := fmt.Sprintf("DISPATCH-%d", time.Now().UTC().UnixNano())
	e := Envelope{Version: 1, ID: id, WorkspaceID: work, UnitKind: "technical-research", UnitPath: filepath.ToSlash(unitPath), Role: "technical-discovery", Agent: agent, InputHash: hex.EncodeToString(sum[:]), RequiredReading: []string{filepath.ToSlash(unitPath)}, ExpectedEvidence: []string{"sources", "options", "uncertainties"}, Status: "assigned", CreatedAt: time.Now().UTC().Format(time.RFC3339), Forbidden: []string{"decision", "approval", "engineering-proposal", "commit", "push", "merge", "release"}}
	return e, write(filepath.Join(dir(root, work), id+".json"), e)
}

// AssignThreatModel is a read-only assignment for one named trust boundary.
func AssignThreatModel(root, work, boundaryPath, agent string) (Envelope, error) {
	e, err := AssignResearch(root, work, boundaryPath, agent)
	if err != nil {
		return e, err
	}
	e.UnitKind = "security-boundary"
	e.Role = "threat-modeler"
	e.ExpectedEvidence = []string{"threats", "controls", "residual risks"}
	e.Forbidden = []string{"risk-acceptance", "approval", "release", "commit", "push", "merge"}
	return e, write(filepath.Join(dir(root, work), e.ID+".json"), e)
}

// AssignImportChunk claims one existing scalable-import chunk. It cannot record
// review evidence or materialize mappings; those remain explicit import steps.
func AssignImportChunk(root, work, run, chunkID, agent string) (Envelope, error) {
	chunk, err := sourceimport.Resume(root, run, chunkID, agent)
	if err != nil {
		return Envelope{}, err
	}
	id := fmt.Sprintf("DISPATCH-%d", time.Now().UTC().UnixNano())
	unit := filepath.Join(root, "knowledge", "imports", "runs", run, "chunks", chunk.ID+".json")
	data, err := os.ReadFile(unit)
	if err != nil {
		return Envelope{}, err
	}
	sum := sha256.Sum256(data)
	if err := os.MkdirAll(dir(root, work), 0755); err != nil {
		return Envelope{}, err
	}
	e := Envelope{Version: 1, ID: id, WorkspaceID: work, UnitKind: "import-chunk", UnitPath: filepath.ToSlash(unit), ImportRun: run, ImportChunk: chunk.ID, Role: "artifact-importer", Agent: agent, InputHash: hex.EncodeToString(sum[:]), RequiredReading: []string{filepath.ToSlash(unit)}, ExpectedEvidence: []string{"evidence per source", "gaps"}, Status: "assigned", CreatedAt: time.Now().UTC().Format(time.RFC3339), Forbidden: []string{"materialize", "approval", "mapping-selection", "commit", "push", "merge", "release"}}
	return e, write(filepath.Join(dir(root, work), id+".json"), e)
}
func Observe(root, work string) ([]Envelope, error) {
	entries, err := os.ReadDir(dir(root, work))
	if os.IsNotExist(err) {
		return []Envelope{}, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Envelope
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		var e Envelope
		if err := read(filepath.Join(dir(root, work), entry.Name()), &e); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Reconcile is read-only and never changes an envelope, task, or lease.
func Reconcile(root, work string) ([]Finding, error) {
	xs, err := Observe(root, work)
	if err != nil {
		return nil, err
	}
	byID := map[string]Envelope{}
	for _, x := range xs {
		byID[x.ID] = x
	}
	var out []Finding
	for _, x := range xs {
		if x.Role != "code-runner" && x.ParentID != "" {
			p, ok := byID[x.ParentID]
			if !ok {
				out = append(out, Finding{"orphaned-review", x.ID, "parent dispatch missing", "delivery-orchestrator"})
			} else if p.DiffHash == "" || p.DiffHash != x.DiffHash {
				out = append(out, Finding{"review-diff-mismatch", x.ID, "review does not match parent diff", "code-review"})
			}
		}
		if x.Status == "assigned" && (x.Role == "qa" || x.Role == "code-review" || x.Role == "security-review") && x.DiffHash == "" {
			out = append(out, Finding{"review-missing-diff", x.ID, "independent review lacks diff hash", "delivery-orchestrator"})
		}
		if x.UnitKind == "engineering-specialist" {
			for _, dependencyID := range x.Dependencies {
				dependency, ok := byID[dependencyID]
				if !ok || dependency.UnitKind != "engineering-specialist" || dependency.Status != "returned" {
					out = append(out, Finding{"engineering-dependency-not-returned", x.ID, "dependency " + dependencyID + " is missing or no longer returned", "engineering-orchestrator"})
				}
			}
			if x.Status == "returned" {
				for _, output := range x.OutputHashes {
					parts := strings.SplitN(output, "=", 2)
					if len(parts) != 2 || !insideAnyScope(parts[0], x.WriteScope) {
						out = append(out, Finding{"engineering-output-invalid", x.ID, "returned output hash is malformed", "engineering-orchestrator"})
						continue
					}
					data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(parts[0])))
					if err != nil {
						out = append(out, Finding{"engineering-output-missing", x.ID, "returned output " + parts[0] + " is missing", "engineering-orchestrator"})
						continue
					}
					sum := sha256.Sum256(data)
					if !strings.EqualFold(parts[1], hex.EncodeToString(sum[:])) {
						out = append(out, Finding{"engineering-output-stale", x.ID, "returned output " + parts[0] + " changed after return", "engineering-orchestrator"})
					}
				}
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DispatchID < out[j].DispatchID })
	return out, nil
}

// Run executes only an assigned code-runner envelope after explicit enablement.
// It deliberately refuses Git delivery commands and records a transcript.
func Run(root, work, id string, enabled bool, command string, args []string) (Transcript, error) {
	if !enabled {
		return Transcript{}, errors.New("supervised dispatch execution is disabled")
	}
	cfg, err := LoadConfig(root)
	if err != nil {
		return Transcript{}, err
	}
	if !cfg.Enabled {
		return Transcript{}, errors.New("dispatch capability is disabled for this product")
	}
	allowed := false
	for _, h := range cfg.Harnesses {
		if strings.EqualFold(filepath.Base(command), filepath.Base(h)) {
			allowed = true
		}
	}
	if !allowed {
		return Transcript{}, errors.New("command is not an enabled dispatch harness")
	}
	var e Envelope
	path := filepath.Join(dir(root, work), id+".json")
	if err := read(path, &e); err != nil {
		return Transcript{}, err
	}
	if e.Role != "code-runner" || e.Status != "assigned" {
		return Transcript{}, errors.New("only assigned code-runner dispatches can run")
	}
	if strings.EqualFold(filepath.Base(command), "git") {
		return Transcript{}, errors.New("dispatch runner cannot invoke git delivery commands")
	}
	started := time.Now().UTC()
	cmd := exec.Command(command, args...)
	cmd.Dir = filepath.Dir(root)
	output, err := cmd.CombinedOutput()
	if err == nil {
		if scopeErr := validateWriteScope(filepath.Dir(root), e.WriteScope); scopeErr != nil {
			err = scopeErr
		}
	}
	exit := 0
	if x, ok := err.(*exec.ExitError); ok {
		exit = x.ExitCode()
	}
	sum := sha256.Sum256(output)
	t := Transcript{DispatchID: id, StartedAt: started.Format(time.RFC3339), FinishedAt: time.Now().UTC().Format(time.RFC3339), ExitCode: exit, OutputHash: hex.EncodeToString(sum[:])}
	data, _ := json.MarshalIndent(t, "", "  ")
	tdir := filepath.Join(dir(root, work), "transcripts")
	if mk := os.MkdirAll(tdir, 0755); mk != nil {
		return t, mk
	}
	if writeErr := os.WriteFile(filepath.Join(tdir, id+".json"), append(data, '\n'), 0644); writeErr != nil {
		return t, writeErr
	}
	if cfg.TranscriptRetention > 0 {
		_ = retainTranscripts(tdir, cfg.TranscriptRetention)
	}
	return t, err
}
func validateWriteScope(repo string, scopes []string) error {
	if len(scopes) == 0 {
		return errors.New("dispatch task has no write scope")
	}
	out, err := exec.Command("git", "-C", repo, "diff", "--name-only").CombinedOutput()
	if err != nil {
		return err
	}
	for _, raw := range strings.Fields(string(out)) {
		ok := false
		for _, scope := range scopes {
			scope = filepath.ToSlash(strings.TrimSuffix(scope, "/"))
			if raw == scope || strings.HasPrefix(raw, scope+"/") {
				ok = true
			}
		}
		if !ok {
			return fmt.Errorf("working-tree path %s escapes dispatch write scope", raw)
		}
	}
	return nil
}
func retainTranscripts(path string, max int) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for len(names) > max {
		if err := os.Remove(filepath.Join(path, names[0])); err != nil {
			return err
		}
		names = names[1:]
	}
	return nil
}

// RunWave runs already-assigned envelopes with bounded local concurrency.
func RunWave(root, work string, ids []string, max int, enabled bool, command string, args []string) []WaveResult {
	if max < 1 {
		max = 1
	}
	sem := make(chan struct{}, max)
	out := make(chan WaveResult, len(ids))
	for _, id := range ids {
		id := id
		go func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			t, e := Run(root, work, id, enabled, command, args)
			r := WaveResult{ID: id, Transcript: &t}
			if e != nil {
				r.Error = e.Error()
			}
			out <- r
		}()
	}
	results := make([]WaveResult, 0, len(ids))
	for range ids {
		results = append(results, <-out)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
	return results
}

// WaveIDs derives assigned code-runner envelopes from the persisted scheduler wave.
func WaveIDs(root, work, wave string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, ".product", "scheduler", "waves", work+".json"))
	if err != nil {
		return nil, err
	}
	var schedule workflow.Schedule
	if err := json.Unmarshal(data, &schedule); err != nil {
		return nil, err
	}
	var tasks map[string]bool = map[string]bool{}
	for _, w := range schedule.Waves {
		if w.ID == wave {
			for _, t := range w.Tasks {
				tasks[t] = true
			}
		}
	}
	if len(tasks) == 0 {
		return nil, errors.New("scheduled wave not found")
	}
	xs, err := Observe(root, work)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, e := range xs {
		if e.Role == "code-runner" && e.Status == "assigned" && tasks[e.TaskID] {
			ids = append(ids, e.ID)
		}
	}
	sort.Strings(ids)
	return ids, nil
}

// Recommend is advisory only; it never assigns, reprioritizes, or executes.
func Recommend(root, work string, max int) ([]Recommendation, error) {
	xs, err := Observe(root, work)
	if err != nil {
		return nil, err
	}
	active := 0
	var out []Recommendation
	for _, x := range xs {
		if x.Status == "assigned" {
			active++
		}
		if x.Status == "returned" && x.Role == "code-runner" {
			out = append(out, Recommendation{"review-ready", "assign independent QA and Code Review for " + x.ID, true})
		}
	}
	if max > 0 && active >= max {
		out = append(out, Recommendation{"capacity", "active dispatches reach configured capacity", true})
	}
	if active == 0 {
		out = append(out, Recommendation{"idle", "no active dispatches; inspect dispatch plan", true})
	}
	return out, nil
}
func mustNodes(graph string) []workflow.Node {
	var raw struct {
		Nodes []workflow.Node `json:"nodes"`
	}
	data, _ := os.ReadFile(graph)
	_ = json.Unmarshal(data, &raw)
	return raw.Nodes
}
func write(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}
func read(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
