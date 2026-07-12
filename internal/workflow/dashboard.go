package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type WorkflowStage struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}
type WorkflowDashboard struct {
	WorkspaceID      string          `json:"workspace_id"`
	Feature          string          `json:"feature"`
	UseCase          string          `json:"use_case,omitempty"`
	CurrentStep      string          `json:"current_step"`
	RecommendedSkill string          `json:"recommended_skill"`
	ExpectedArtifact string          `json:"expected_artifact"`
	Stages           []WorkflowStage `json:"stages"`
	Blockers         []string        `json:"blockers,omitempty"`
	RequiredReading  []string        `json:"required_reading,omitempty"`
	NextCommands     []string        `json:"next_commands,omitempty"`
	Decisions        []string        `json:"decisions,omitempty"`
	ActiveLeases     []string        `json:"active_leases,omitempty"`
	GraphStatus      string          `json:"graph_status,omitempty"`
	TaskTotal        int             `json:"task_total"`
	TaskReady        int             `json:"task_ready"`
	LatestCheckpoint string          `json:"latest_checkpoint,omitempty"`
	LatestHandoff    string          `json:"latest_handoff,omitempty"`
}

func BuildDashboard(root, id string) (WorkflowDashboard, error) {
	guide, err := WorkspaceGuide(root, id)
	if err != nil {
		return WorkflowDashboard{}, err
	}
	w, err := LoadWorkspace(root, id)
	if err != nil {
		return WorkflowDashboard{}, err
	}
	d := WorkflowDashboard{WorkspaceID: id, Feature: w.Scope["feature"], UseCase: w.Scope["use_case"], CurrentStep: guide.CurrentStep, RecommendedSkill: guide.RecommendedSkill, ExpectedArtifact: guide.ExpectedArtifact, Blockers: guide.Blockers, RequiredReading: guide.RequiredReading, NextCommands: guide.Commands}
	order := []struct{ id, label, skill string }{{"feature", "Feature", "feature"}, {"use-case", "Use Case", "use-case"}, {"specification", "Specification", "specification"}, {"design", "Design", "ux-ui"}, {"technical-discovery", "Technical Discovery", "technical-discovery"}, {"architecture-gate", "Architecture Gate", "product-historian"}, {"implementation-plan", "Implementation Plan", "implementation-planner"}, {"execution-graph", "Execution Graph", "execution-graph"}, {"tasks", "Tasks", "task-generator"}, {"code-runner", "Ready for Code", "code-runner"}}
	current := len(order) - 1
	for i, x := range order {
		if x.skill == guide.CurrentStep {
			current = i
			break
		}
	}
	for i, x := range order {
		status := "pending"
		if i < current {
			status = "done"
		} else if i == current {
			status = "current"
			if len(guide.Blockers) > 0 {
				status = "blocked"
			}
		}
		d.Stages = append(d.Stages, WorkflowStage{ID: x.id, Label: x.label, Status: status})
	}
	uc := resolveDashboardUseCase(root, w)
	if uc != "" {
		rel, _ := filepath.Rel(root, uc)
		d.UseCase = filepath.ToSlash(rel)
		graph := filepath.Join(uc, "execution-graph.json")
		d.GraphStatus = jsonStatus(graph)
		var g Graph
		if readJSON(graph, &g) == nil {
			d.TaskTotal = len(g.Nodes)
			ready, _ := ReadyUnclaimed(root, graph)
			d.TaskReady = len(ready)
		}
	}
	d.Decisions = dashboardDecisions(root, d.Feature, d.UseCase)
	d.ActiveLeases = dashboardLeases(root)
	d.LatestCheckpoint = latestRuntimeFile(workspaceDir(root, id), "checkpoints")
	d.LatestHandoff = latestRuntimeFile(workspaceDir(root, id), "handoffs")
	return d, nil
}
func resolveDashboardUseCase(root string, w Workspace) string {
	if s := w.Scope["use_case"]; s != "" {
		p := filepath.Join(root, filepath.FromSlash(s))
		if info, e := os.Stat(p); e == nil {
			if info.IsDir() {
				return p
			}
			return filepath.Dir(p)
		}
	}
	base := filepath.Dir(filepath.Join(root, filepath.FromSlash(w.Scope["feature"])))
	entries, _ := os.ReadDir(filepath.Join(base, "use-cases"))
	var dirs []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), "_") {
			dirs = append(dirs, filepath.Join(base, "use-cases", e.Name()))
		}
	}
	if len(dirs) == 1 {
		return dirs[0]
	}
	return ""
}
func dashboardDecisions(root string, scopes ...string) []string {
	var index map[string]any
	if readJSON(filepath.Join(root, ".product", "decisions.json"), &index) != nil {
		return nil
	}
	var out []string
	for _, d := range stringMapSlice(index["decisions"]) {
		id := fmt.Sprint(d["id"])
		for _, a := range stringAnySlice(d["affectedArtifacts"]) {
			for _, scope := range scopes {
				scope = strings.TrimSuffix(filepath.ToSlash(scope), "context.md")
				if scope != "" && strings.HasPrefix(filepath.ToSlash(a), scope) {
					out = append(out, id+" ["+fmt.Sprint(d["status"])+"]")
				}
			}
		}
	}
	return uniqueSorted(out)
}
func dashboardLeases(root string) []string {
	entries, _ := os.ReadDir(filepath.Join(root, ".product", "claims"))
	var out []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		var l Lease
		if readJSON(filepath.Join(root, ".product", "claims", e.Name()), &l) == nil {
			out = append(out, l.TaskID+" → "+l.Agent+" until "+l.ExpiresAt)
		}
	}
	sort.Strings(out)
	return out
}
func latestRuntimeFile(base, child string) string {
	entries, _ := os.ReadDir(filepath.Join(base, child))
	if len(entries) == 0 {
		return ""
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries[len(entries)-1].Name()
}
func uniqueSorted(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range items {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	sort.Strings(out)
	return out
}
