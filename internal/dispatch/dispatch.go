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

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

type Envelope struct {
	Version          int      `json:"version"`
	ID               string   `json:"id"`
	WorkspaceID      string   `json:"workspace_id"`
	TaskID           string   `json:"task_id"`
	Role             string   `json:"role"`
	Agent            string   `json:"agent"`
	Graph            string   `json:"graph"`
	TaskPath         string   `json:"task_path"`
	InputHash        string   `json:"input_hash"`
	DiffHash         string   `json:"diff_hash,omitempty"`
	ParentID         string   `json:"parent_id,omitempty"`
	RequiredReading  []string `json:"required_reading"`
	WriteScope       []string `json:"write_scope"`
	ExpectedEvidence []string `json:"expected_evidence"`
	Status           string   `json:"status"`
	CreatedAt        string   `json:"created_at"`
	ReturnedAt       string   `json:"returned_at,omitempty"`
	Summary          string   `json:"summary,omitempty"`
	Evidence         []string `json:"evidence,omitempty"`
	Forbidden        []string `json:"forbidden"`
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

func dir(root, work string) string {
	return filepath.Join(root, ".product", "workspaces", work, "dispatches")
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
	e := Envelope{Version: 1, ID: id, WorkspaceID: work, TaskID: task, Role: role, Agent: agent, Graph: filepath.ToSlash(graph), TaskPath: filepath.ToSlash(taskPath), InputHash: hex.EncodeToString(sum[:]), RequiredReading: []string{filepath.ToSlash(taskPath), filepath.ToSlash(graph)}, ExpectedEvidence: []string{"diff hash", "test log"}, Status: "assigned", CreatedAt: time.Now().UTC().Format(time.RFC3339), Forbidden: []string{"approval", "commit", "push", "merge", "release", "review-resolution"}}
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
	return e, workflow.ReleaseLease(root, e.TaskID, agent)
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
		if x.Status == "assigned" && x.Role != "code-runner" && x.DiffHash == "" {
			out = append(out, Finding{"review-missing-diff", x.ID, "independent review lacks diff hash", "delivery-orchestrator"})
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
	return t, err
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
