package workflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const RuntimeVersion = 2

type RuntimeState struct {
	Version     int      `json:"version"`
	WorkspaceID string   `json:"workspace_id"`
	Phase       string   `json:"phase"`
	Status      string   `json:"status"`
	UpdatedAt   string   `json:"updated_at"`
	Attempts    int      `json:"attempts"`
	Blockers    []string `json:"blockers,omitempty"`
}
type Handoff struct {
	Version         int      `json:"version"`
	ID              string   `json:"id"`
	WorkspaceID     string   `json:"workspace_id"`
	From            string   `json:"from"`
	To              string   `json:"to"`
	Summary         string   `json:"summary"`
	CreatedAt       string   `json:"created_at"`
	RequiredReading []string `json:"required_reading,omitempty"`
	Blockers        []string `json:"blockers,omitempty"`
}
type Checkpoint struct {
	Version                                                             int `json:"version"`
	ID, WorkspaceID, Step, BaseCommit, InputHash, OutputHash, CreatedAt string
	Stale                                                               bool `json:"stale"`
}
type Lease struct {
	Version                                                 int `json:"version"`
	TaskID, Graph, Agent, ClaimedAt, HeartbeatAt, ExpiresAt string
	Attempt                                                 int `json:"attempt"`
}
type CommandPlan struct {
	Version                                                                      int `json:"version"`
	ID, WorkspaceID, TaskID, Cwd, Source, Risk, BaseCommit, InputHash, CreatedAt string
	Argv                                                                         []string `json:"argv"`
	TimeoutSeconds                                                               int      `json:"timeout_seconds"`
	AllowedWrites                                                                []string `json:"allowed_writes"`
	EnvAllowlist                                                                 []string `json:"env_allowlist"`
}
type CommandEvidence struct {
	Version                                   int `json:"version"`
	PlanID, StartedAt, FinishedAt, OutputHash string
	ExitCode                                  int    `json:"exit_code"`
	Success                                   bool   `json:"success"`
	Output                                    string `json:"output"`
}
type Wave struct {
	ID    string   `json:"id"`
	Tasks []string `json:"tasks"`
}
type Schedule struct {
	Version                       int `json:"version"`
	WorkspaceID, Graph, CreatedAt string
	MaxParallel                   int    `json:"max_parallel"`
	Waves                         []Wave `json:"waves"`
}
type Integration struct {
	Version                                        int `json:"version"`
	ID, WorkspaceID, BaseCommit, Status, CreatedAt string
	Commits                                        []string `json:"commits"`
	IntegratedDiffHash                             string   `json:"integrated_diff_hash,omitempty"`
	RequiresIntegratedQA                           bool     `json:"requires_integrated_qa"`
}

func workspaceDir(root, id string) string { return filepath.Join(root, ".product", "workspaces", id) }
func LoadWorkspace(root, id string) (Workspace, error) {
	var w Workspace
	if err := readJSON(filepath.Join(workspaceDir(root, id), "workspace.json"), &w); err == nil {
		return w, nil
	}
	return w, readJSON(filepath.Join(root, ".product", "workspaces", id+".json"), &w)
}
func Resume(root, id string) (RuntimeState, error) {
	var s RuntimeState
	if err := readJSON(filepath.Join(workspaceDir(root, id), "state.json"), &s); err == nil {
		return s, nil
	}
	w, err := LoadWorkspace(root, id)
	if err != nil {
		return s, err
	}
	return RuntimeState{Version: 1, WorkspaceID: id, Phase: w.CurrentStep, Status: "legacy", UpdatedAt: w.CreatedAt, Blockers: w.BlockedBy}, nil
}
func MigrateWorkspace(root, id string, dry bool) (string, error) {
	w, err := LoadWorkspace(root, id)
	if err != nil {
		return "", err
	}
	dst := workspaceDir(root, id)
	if _, err = os.Stat(dst); err == nil {
		return "already v2", nil
	}
	if dry {
		return "migrate " + id + " to runtime v2", nil
	}
	if err = os.MkdirAll(dst, 0755); err != nil {
		return "", err
	}
	if err = writeJSON(filepath.Join(dst, "workspace.json"), w); err != nil {
		return "", err
	}
	for _, d := range []string{"handoffs", "checkpoints", "command-plans", "evidence", "tasks"} {
		if err = os.MkdirAll(filepath.Join(dst, d), 0755); err != nil {
			return "", err
		}
	}
	s := RuntimeState{Version: 2, WorkspaceID: id, Phase: w.CurrentStep, Status: "active", UpdatedAt: time.Now().UTC().Format(time.RFC3339), Blockers: w.BlockedBy}
	if err = writeJSON(filepath.Join(dst, "state.json"), s); err != nil {
		return "", err
	}
	old := filepath.Join(root, ".product", "workspaces", id+".json")
	_ = os.Rename(old, old+".v1.bak")
	return "migrated " + id, nil
}
func WriteHandoff(root, id, from, to, summary string) (Handoff, error) {
	dir := filepath.Join(workspaceDir(root, id), "handoffs")
	_ = os.MkdirAll(dir, 0755)
	n, _ := nextID(dir, "HANDOFF-")
	h := Handoff{Version: 2, ID: n, WorkspaceID: id, From: from, To: to, Summary: summary, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	return h, writeJSON(filepath.Join(dir, n+".json"), h)
}
func WriteCheckpoint(root, id, step, base, input, output string) (Checkpoint, error) {
	dir := filepath.Join(workspaceDir(root, id), "checkpoints")
	_ = os.MkdirAll(dir, 0755)
	n, _ := nextID(dir, "CHECKPOINT-")
	c := Checkpoint{Version: 2, ID: n, WorkspaceID: id, Step: step, BaseCommit: base, InputHash: input, OutputHash: output, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	return c, writeJSON(filepath.Join(dir, n+".json"), c)
}

func leasePath(root, task string) string {
	return filepath.Join(root, ".product", "claims", task+".json")
}
func ClaimLease(root, graph, task, agent string, ttl time.Duration) (Lease, error) {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	now := time.Now().UTC()
	var old Lease
	if readJSON(leasePath(root, task), &old) == nil {
		exp, _ := time.Parse(time.RFC3339, old.ExpiresAt)
		if exp.After(now) {
			return Lease{}, fmt.Errorf("task %s is leased by %s until %s", task, old.Agent, old.ExpiresAt)
		}
	}
	if _, err := ClaimTask(root, graph, task, agent); err != nil && !strings.Contains(err.Error(), "already claimed") {
		return Lease{}, err
	}
	_ = os.MkdirAll(filepath.Dir(leasePath(root, task)), 0755)
	l := Lease{Version: 2, TaskID: task, Agent: agent, Graph: filepath.ToSlash(graph), ClaimedAt: now.Format(time.RFC3339), HeartbeatAt: now.Format(time.RFC3339), ExpiresAt: now.Add(ttl).Format(time.RFC3339), Attempt: old.Attempt + 1}
	if l.Attempt > 3 {
		return Lease{}, fmt.Errorf("task %s exceeded three attempts", task)
	}
	return l, writeJSON(leasePath(root, task), l)
}
func Heartbeat(root, task, agent string, ttl time.Duration) (Lease, error) {
	var l Lease
	if err := readJSON(leasePath(root, task), &l); err != nil {
		return l, err
	}
	if l.Agent != agent {
		return l, fmt.Errorf("task %s is not leased by %s", task, agent)
	}
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	now := time.Now().UTC()
	l.HeartbeatAt = now.Format(time.RFC3339)
	l.ExpiresAt = now.Add(ttl).Format(time.RFC3339)
	return l, writeJSON(leasePath(root, task), l)
}
func RecoverLeases(root string) ([]string, error) {
	dir := filepath.Join(root, ".product", "claims")
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	var out []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		var l Lease
		if readJSON(filepath.Join(dir, e.Name()), &l) != nil {
			continue
		}
		exp, _ := time.Parse(time.RFC3339, l.ExpiresAt)
		if !exp.After(now) {
			_ = os.Remove(filepath.Join(dir, e.Name()))
			_ = ReleaseClaim(root, l.TaskID, l.Agent)
			out = append(out, l.TaskID)
		}
	}
	sort.Strings(out)
	return out, nil
}

func CreateTaskWorktree(repoRoot, work, task string) (string, error) {
	if work == "" || task == "" {
		return "", fmt.Errorf("work and task are required")
	}
	branch := "codex/" + strings.ToLower(work+"-"+task)
	path := filepath.Join(repoRoot, ".worktrees", work, task)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	if out, err := gitOutput(repoRoot, "worktree", "add", "-b", branch, path); err != nil {
		return "", fmt.Errorf("git worktree add: %s", strings.TrimSpace(out))
	}
	return path, nil
}

func CreateCommandPlan(root, work, task, cwd, source, risk string, argv []string, timeout int) (CommandPlan, error) {
	if regexp.MustCompile(`(?i)^DEC-\d+$`).MatchString(strings.TrimSpace(source)) {
		return CommandPlan{}, fmt.Errorf("decision text is not an executable command source; use a validated gate")
	}
	risk = strings.ToUpper(risk)
	if risk != "R0" && risk != "R1" {
		return CommandPlan{}, fmt.Errorf("risk %s is disabled", risk)
	}
	if len(argv) == 0 {
		return CommandPlan{}, fmt.Errorf("argv is required")
	}
	if timeout <= 0 {
		timeout = 300
	}
	dir := filepath.Join(workspaceDir(root, work), "command-plans")
	_ = os.MkdirAll(dir, 0755)
	id, _ := nextID(dir, "CMDPLAN-")
	base, _ := gitOutput(root, "rev-parse", "HEAD")
	raw, _ := json.Marshal(argv)
	sum := sha256.Sum256(raw)
	p := CommandPlan{Version: 2, ID: id, WorkspaceID: work, TaskID: task, Cwd: filepath.ToSlash(cwd), Source: source, Risk: risk, Argv: argv, TimeoutSeconds: timeout, BaseCommit: strings.TrimSpace(base), InputHash: hex.EncodeToString(sum[:]), CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	return p, writeJSON(filepath.Join(dir, id+".json"), p)
}
func ExecuteCommandPlan(root, work, id string) (CommandEvidence, error) {
	var p CommandPlan
	if err := readJSON(filepath.Join(workspaceDir(root, work), "command-plans", id+".json"), &p); err != nil {
		return CommandEvidence{}, err
	}
	if p.Risk != "R0" && p.Risk != "R1" {
		return CommandEvidence{}, fmt.Errorf("risk %s is disabled", p.Risk)
	}
	cwd := filepath.Join(root, filepath.FromSlash(p.Cwd))
	rel, err := filepath.Rel(root, cwd)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return CommandEvidence{}, fmt.Errorf("cwd escapes repository")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.TimeoutSeconds)*time.Second)
	defer cancel()
	started := time.Now().UTC()
	cmd := exec.CommandContext(ctx, p.Argv[0], p.Argv[1:]...)
	cmd.Dir = cwd
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "SYSTEMROOT=" + os.Getenv("SYSTEMROOT"), "TEMP=" + os.TempDir()}
	out, runErr := cmd.CombinedOutput()
	exit := 0
	if runErr != nil {
		exit = -1
		if x, ok := runErr.(*exec.ExitError); ok {
			exit = x.ExitCode()
		}
	}
	sum := sha256.Sum256(out)
	e := CommandEvidence{Version: 2, PlanID: id, StartedAt: started.Format(time.RFC3339), FinishedAt: time.Now().UTC().Format(time.RFC3339), ExitCode: exit, Success: runErr == nil, OutputHash: hex.EncodeToString(sum[:]), Output: string(out)}
	dir := filepath.Join(workspaceDir(root, work), "evidence")
	_ = os.MkdirAll(dir, 0755)
	if err := writeJSON(filepath.Join(dir, id+".json"), e); err != nil {
		return e, err
	}
	return e, runErr
}

func BuildSchedule(root, work, graph string, max int) (Schedule, error) {
	if max <= 0 {
		max = 4
	}
	var g Graph
	if err := readJSON(graph, &g); err != nil {
		return Schedule{}, err
	}
	done := map[string]bool{}
	remaining := map[string]Node{}
	for _, n := range g.Nodes {
		if n.Status == "complete" || n.Status == "validated" {
			done[n.ID] = true
		} else {
			remaining[n.ID] = n
		}
	}
	s := Schedule{Version: 2, WorkspaceID: work, Graph: filepath.ToSlash(graph), MaxParallel: max, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	for len(remaining) > 0 {
		var ids []string
		for id, n := range remaining {
			ok := true
			for _, d := range n.DependsOn {
				if !done[d] {
					ok = false
				}
			}
			if ok {
				ids = append(ids, id)
			}
		}
		sort.Strings(ids)
		if len(ids) == 0 {
			return s, fmt.Errorf("graph has a cycle or unresolved dependency")
		}
		var chosen []string
		for _, id := range ids {
			n := remaining[id]
			conflict := false
			for _, x := range chosen {
				o := remaining[x]
				if scopesOverlap(n.WriteScope, o.WriteScope) || resourcesOverlap(n.SharedResources, o.SharedResources) {
					conflict = true
				}
			}
			if !conflict && len(chosen) < max {
				chosen = append(chosen, id)
			}
		}
		s.Waves = append(s.Waves, Wave{ID: fmt.Sprintf("WAVE-%03d", len(s.Waves)+1), Tasks: chosen})
		for _, id := range chosen {
			done[id] = true
			delete(remaining, id)
		}
	}
	dir := filepath.Join(root, ".product", "scheduler", "waves")
	_ = os.MkdirAll(dir, 0755)
	return s, writeJSON(filepath.Join(dir, work+".json"), s)
}
func CreateIntegration(root, work, base string, commits []string) (Integration, error) {
	dir := filepath.Join(root, ".product", "integrations")
	_ = os.MkdirAll(dir, 0755)
	id, _ := nextID(dir, "INTEGRATION-")
	i := Integration{Version: 2, ID: id, WorkspaceID: work, BaseCommit: base, Commits: commits, Status: "planned", CreatedAt: time.Now().UTC().Format(time.RFC3339), RequiresIntegratedQA: true}
	return i, writeJSON(filepath.Join(dir, id+".json"), i)
}
func ApplyIntegration(root, id string) (Integration, error) {
	path := filepath.Join(root, ".product", "integrations", id+".json")
	var i Integration
	if err := readJSON(path, &i); err != nil {
		return i, err
	}
	for _, c := range i.Commits {
		if _, err := gitOutput(root, "cherry-pick", c); err != nil {
			_, _ = gitOutput(root, "cherry-pick", "--abort")
			i.Status = "conflict"
			_ = writeJSON(path, i)
			return i, fmt.Errorf("integration conflict at %s: %w", c, err)
		}
	}
	diff, _ := gitOutput(root, "diff", i.BaseCommit+"..HEAD")
	sum := sha256.Sum256([]byte(diff))
	i.IntegratedDiffHash = hex.EncodeToString(sum[:])
	i.Status = "awaiting_integrated_qa"
	return i, writeJSON(path, i)
}
func gitOutput(dir string, args ...string) (string, error) {
	b, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
	return string(b), err
}
