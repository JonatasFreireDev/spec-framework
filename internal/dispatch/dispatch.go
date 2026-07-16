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
func Return(root, work, id, agent, summary string, evidence []string) (Envelope, error) {
	var e Envelope
	path := filepath.Join(dir(root, work), id+".json")
	if err := read(path, &e); err != nil {
		return e, err
	}
	if e.Agent != agent || e.Status != "assigned" {
		return e, errors.New("dispatch is not assigned to this agent")
	}
	if strings.TrimSpace(summary) == "" || len(evidence) == 0 {
		return e, errors.New("return requires summary and evidence")
	}
	e.Status = "returned"
	e.Summary = summary
	e.Evidence = evidence
	e.ReturnedAt = time.Now().UTC().Format(time.RFC3339)
	if err := write(path, e); err != nil {
		return e, err
	}
	return e, workflow.ReleaseLease(root, e.TaskID, agent)
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
