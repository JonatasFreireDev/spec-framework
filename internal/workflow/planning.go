package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type MaterializeResult struct {
	Graph string   `json:"graph"`
	Tasks []string `json:"tasks"`
	Index string   `json:"index"`
}
type ReadinessCheck struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}
type TaskReadiness struct {
	Ready  bool             `json:"ready"`
	TaskID string           `json:"task_id"`
	Checks []ReadinessCheck `json:"checks"`
}
type Guide struct {
	WorkspaceID, CurrentStep, RecommendedSkill, ExpectedArtifact string
	RequiredReading, Blockers, Commands                          []string
}
type StageReview struct {
	WorkspaceID, Stage string
	Artifacts          []Artifact
	Blockers           []string
}

func MaterializeTasks(graphPath string) (MaterializeResult, error) {
	original, err := os.ReadFile(graphPath)
	if err != nil {
		return MaterializeResult{}, err
	}
	var raw map[string]any
	if err = json.Unmarshal(original, &raw); err != nil {
		return MaterializeResult{}, err
	}
	status := strings.ToLower(fmt.Sprint(raw["status"]))
	if status != "draft" && status != "proposed" {
		return MaterializeResult{}, fmt.Errorf("graph status must be draft or proposed, got %s", status)
	}
	nodes, ok := raw["nodes"].([]any)
	if !ok || len(nodes) == 0 {
		return MaterializeResult{}, fmt.Errorf("graph has no nodes")
	}
	base := filepath.Dir(graphPath)
	var paths []string
	seen := map[string]bool{}
	for _, item := range nodes {
		n, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id := fmt.Sprint(n["id"])
		rel := filepath.ToSlash(fmt.Sprint(n["path"]))
		if id == "" || rel == "" {
			return MaterializeResult{}, fmt.Errorf("graph node requires id and path")
		}
		if !strings.HasPrefix(rel, "tasks/") || filepath.Ext(rel) != ".md" {
			return MaterializeResult{}, fmt.Errorf("node %s path must be tasks/<id>.md", id)
		}
		if seen[rel] {
			return MaterializeResult{}, fmt.Errorf("duplicate task path %s", rel)
		}
		seen[rel] = true
		path := filepath.Join(base, filepath.FromSlash(rel))
		if _, e := os.Stat(path); e == nil {
			return MaterializeResult{}, fmt.Errorf("refusing to overwrite existing task %s", rel)
		}
		paths = append(paths, path)
	}
	created := []string{}
	rollback := func() {
		for _, p := range created {
			_ = os.Remove(p)
		}
		_ = atomicWrite(graphPath, original)
	}
	for _, item := range nodes {
		n := item.(map[string]any)
		rel := filepath.ToSlash(fmt.Sprint(n["path"]))
		path := filepath.Join(base, filepath.FromSlash(rel))
		body := renderTask(raw, n)
		if err = atomicWrite(path, []byte(body)); err != nil {
			rollback()
			return MaterializeResult{}, err
		}
		created = append(created, path)
	}
	indexPath := filepath.Join(base, "tasks.md")
	if _, e := os.Stat(indexPath); e == nil {
		rollback()
		return MaterializeResult{}, fmt.Errorf("refusing to overwrite existing tasks.md")
	}
	index := renderTaskIndex(raw, nodes)
	if err = atomicWrite(indexPath, []byte(index)); err != nil {
		rollback()
		return MaterializeResult{}, err
	}
	created = append(created, indexPath)
	raw["status"] = "materialized"
	raw["updatedAt"] = time.Now().UTC().Format("2006-01-02")
	if err = writeJSON(graphPath, raw); err != nil {
		rollback()
		return MaterializeResult{}, err
	}
	result := MaterializeResult{Graph: graphPath, Index: indexPath}
	for _, p := range paths {
		result.Tasks = append(result.Tasks, p)
	}
	return result, nil
}
func renderTask(graph, node map[string]any) string {
	id := fmt.Sprint(node["id"])
	title := fmt.Sprint(node["title"])
	if title == "<nil>" || title == "" {
		title = id
	}
	list := func(v any) string {
		a, _ := v.([]any)
		var x []string
		for _, z := range a {
			x = append(x, fmt.Sprint(z))
		}
		if len(x) == 0 {
			return "None"
		}
		return strings.Join(x, ", ")
	}
	delivery, _ := node["delivery"].(map[string]any)
	return fmt.Sprintf("# Task: %s\n\n## Snapshot\n\n| Field | Value |\n| --- | --- |\n| ID | `%s` |\n| Status | `draft` |\n| Source graph | `%s` |\n| Source specification | `%s` |\n| Source node | `%s` |\n| Owner skill | `%s` |\n| Next skill | `code-runner` |\n\n## Delivery\n\n| Field | Value |\n| --- | --- |\n| Level | `%s` |\n| Priority | `%s` |\n| Depends on | `%s` |\n| Rationale | Inherited from the approved graph. |\n\n## Task Contract\n\n| Field | Value |\n| --- | --- |\n| Title | %s |\n| Type | `%s` |\n| Depends on | `%s` |\n| Source sections | `%s` |\n| Requirements | `%s` |\n| Acceptance criteria | `%s` |\n| Planned tests | `%s` |\n| Applicable decisions | `%s` |\n| Write scope | `%s` |\n| Shared resources | `%s` |\n| Graph node status | `%s` |\n\n## Objective\n\n%s\n\n## Acceptance Checks\n\n- %s\n\n## Blockers\n\n- None.\n\n## Handoff\n\n| Field | Value |\n| --- | --- |\n| Ready for implementation | `no; requires approval and readiness check` |\n| Required next skill | `code-runner` |\n| Notes | Materialized from the proposed execution graph. |\n", title, id, graph["id"], graph["sourceSpecification"], id, node["ownerSkill"], delivery["level"], delivery["priority"], list(node["dependsOn"]), title, node["type"], list(node["dependsOn"]), list(node["sourceSections"]), list(node["requirements"]), list(node["acceptanceCriteria"]), list(node["plannedTests"]), list(node["decisions"]), list(node["writeScope"]), list(node["sharedResources"]), node["status"], title, list(node["acceptanceChecks"]))
}
func renderTaskIndex(graph map[string]any, nodes []any) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Tasks: %s\n\n| Field | Value |\n| --- | --- |\n| ID | `TASKSET-%s` |\n| Status | `draft` |\n| Owner skill | `task-generator` |\n| Next skill | `code-runner` |\n\nGenerated from `execution-graph.json`. Do not edit this index manually.\n\n| Task | Title | Depends on |\n| --- | --- | --- |\n", graph["id"], graph["id"])
	for _, item := range nodes {
		n := item.(map[string]any)
		fmt.Fprintf(&b, "| [%s](%s) | %s | %v |\n", n["id"], n["path"], n["title"], n["dependsOn"])
	}
	return b.String()
}

func CheckTaskReadiness(root, graphPath, taskID string) (TaskReadiness, error) {
	var g Graph
	if err := readJSON(graphPath, &g); err != nil {
		return TaskReadiness{}, err
	}
	r := TaskReadiness{Ready: true, TaskID: taskID}
	add := func(id string, ok bool, detail string) {
		s := "pass"
		if !ok {
			s = "block"
			r.Ready = false
		}
		r.Checks = append(r.Checks, ReadinessCheck{ID: id, Status: s, Detail: detail})
	}
	var n *Node
	for i := range g.Nodes {
		if g.Nodes[i].ID == taskID {
			n = &g.Nodes[i]
		}
	}
	if n == nil {
		return r, fmt.Errorf("task %s not found", taskID)
	}
	var graphRaw map[string]any
	_ = readJSON(graphPath, &graphRaw)
	add("graph-status", isApproved(strings.ToLower(fmt.Sprint(graphRaw["status"]))), "execution graph must be approved")
	taskPath := filepath.Join(filepath.Dir(graphPath), filepath.FromSlash(n.Path))
	data, err := os.ReadFile(taskPath)
	add("task-file", err == nil, n.Path)
	text := string(data)
	taskStatus := extractStatus(text)
	add("task-status", isApproved(taskStatus), "task must be approved")
	add("approval-record", hasCurrentApproval(root, taskPath, taskStatus), "task status must have a matching current approval record")
	for _, d := range n.DependsOn {
		done := false
		for _, x := range g.Nodes {
			if x.ID == d && (x.Status == "complete" || x.Status == "validated") {
				done = true
			}
		}
		add("dependency-"+d, done, "dependency must be complete")
	}
	add("write-scope", len(n.WriteScope) > 0, "writeScope is required")
	trace := strings.Contains(text, "REQ-") && strings.Contains(text, "AC-") && (strings.Contains(text, "TEST-") || strings.Contains(strings.ToLower(text), "evidence"))
	add("traceability", trace, "REQ, AC, and TEST/evidence are required")
	graphBytes, _ := os.ReadFile(graphPath)
	decisionIDs := uniqueDecisionIDs(text + "\n" + string(graphBytes))
	for _, decision := range decisionEffectsFor(root, decisionIDs) {
		id := fmt.Sprint(decision["id"])
		effects, _ := decision["workflowEffects"].(map[string]any)
		for _, required := range stringAnySlice(effects["requiredTaskTypes"]) {
			found := false
			for _, x := range g.Nodes {
				if x.Type == required {
					found = true
				}
			}
			add("decision-"+id+"-task-type-"+required, found, "graph must cover required task type")
		}
		for _, required := range stringAnySlice(effects["requiredWriteScopes"]) {
			found := false
			for _, x := range g.Nodes {
				for _, scope := range x.WriteScope {
					if scope == required || strings.HasPrefix(scope, strings.TrimSuffix(required, "/")+"/") {
						found = true
					}
				}
			}
			add("decision-"+id+"-write-scope", found, "graph must cover "+required)
		}
		for _, required := range stringAnySlice(effects["sharedResources"]) {
			found := false
			for _, x := range g.Nodes {
				for _, resource := range x.SharedResources {
					if resource == required {
						found = true
					}
				}
			}
			add("decision-"+id+"-resource", found, "graph must declare shared resource "+required)
		}
		gatesText, _ := os.ReadFile(filepath.Join(root, "knowledge", "conventions", "gates.md"))
		for _, required := range stringAnySlice(effects["requiredGates"]) {
			configured := strings.Contains(string(gatesText), required) && !lineForIDContains(string(gatesText), required, "TBD")
			add("decision-"+id+"-gate-"+required, configured, "required gate must be configured")
		}
		for _, required := range stringAnySlice(effects["requiredEvidence"]) {
			add("decision-"+id+"-evidence", strings.Contains(strings.ToLower(text), strings.ToLower(required)), "task must declare evidence "+required)
		}
	}
	missing, e := GateReadiness(root)
	add("technical-gates", e == nil && len(missing) == 0, "all product gates must be configured")
	if _, e = os.Stat(leasePath(root, taskID)); e == nil {
		add("lease", false, "task already leased")
	} else {
		add("lease", true, "lease available")
	}
	return r, nil
}
func uniqueDecisionIDs(text string) []string {
	re := regexp.MustCompile(`\bDEC-\d+\b`)
	seen := map[string]bool{}
	var out []string
	for _, x := range re.FindAllString(text, -1) {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	sort.Strings(out)
	return out
}
func lineForIDContains(text, id, needle string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, id) && strings.Contains(strings.ToUpper(line), strings.ToUpper(needle)) {
			return true
		}
	}
	return false
}

func hasCurrentApproval(root, artifactPath, status string) bool {
	if !isApproved(status) {
		return false
	}
	rel, err := filepath.Rel(root, artifactPath)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return false
	}
	hash := Hash(string(data))
	entries, _ := os.ReadDir(filepath.Join(root, ".product", "history"))
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		var a Approval
		if readJSON(filepath.Join(root, ".product", "history", e.Name()), &a) == nil && filepath.ToSlash(a.Path) == rel && a.StatusGranted == status && a.ContentHash == hash {
			return true
		}
	}
	return false
}

func WorkspaceGuide(root, id string) (Guide, error) {
	s, err := WorkspaceStatus(root, id)
	if err != nil {
		return Guide{}, err
	}
	g := Guide{WorkspaceID: id, CurrentStep: s.Next, RecommendedSkill: s.Next, Blockers: s.Blockers}
	m := map[string][]string{"use-case": {"feature context", "feature.md"}, "specification": {"use-case context", "use-case.md"}, "ux-ui": {"specification.md", "contracts/"}, "technical-discovery": {"specification.md", "design.md", "engineering/"}, "product-historian": {"technical-discovery.md", "knowledge/decisions/"}, "engineering-proposal": {"technical-discovery.md", "engineering/", "knowledge/decisions/"}, "engineering-review": {"engineering-proposal.md", "technical-discovery.md", "engineering/", "knowledge/decisions/"}, "implementation-planner": {"engineering-proposal.md", "engineering-review.md", "technical-discovery.md", "design.md", "specification.md"}, "execution-graph": {"implementation-plan.md"}, "task-generator": {"execution-graph.json"}, "code-runner": {"task file", "knowledge/conventions/gates.md"}}
	g.RequiredReading = m[s.Next]
	g.ExpectedArtifact = map[string]string{"feature": "approved feature scope", "use-case": "use-cases/<slug>/", "specification": "specification.md and contracts/", "ux-ui": "design.md", "technical-discovery": "technical-discovery.md", "product-historian": "resolved Architecture Gate", "engineering-proposal": "engineering-proposal.md", "engineering-review": "engineering-review.md with a current verdict", "implementation-planner": "implementation-plan.md", "execution-graph": "execution-graph.json", "task-generator": "tasks/*.md and tasks.md", "code-runner": "working-tree evidence"}[s.Next]
	g.Commands = []string{"spec-framework status --work " + id, "Use skill: " + s.Next}
	return g, nil
}

func ReviewStage(root, work, stage string) (StageReview, error) {
	w, err := LoadWorkspace(root, work)
	if err != nil {
		return StageReview{}, err
	}
	reg, err := LoadRegistry(root)
	if err != nil {
		return StageReview{}, err
	}
	prefix := strings.TrimSuffix(filepath.ToSlash(w.Scope["feature"]), "context.md")
	types := map[string][]string{"feature": {"feature"}, "use-cases": {"use-case"}, "specification": {"specification", "specification-contract"}, "design": {"design"}, "technical-discovery": {"technical-discovery"}, "engineering": {"engineering-proposal", "engineering-review"}, "planning": {"implementation-plan"}, "tasks": {"execution-graph", "task", "taskset"}}[stage]
	if len(types) == 0 {
		return StageReview{}, fmt.Errorf("unknown stage %s", stage)
	}
	allowed := map[string]bool{}
	for _, t := range types {
		allowed[t] = true
	}
	r := StageReview{WorkspaceID: work, Stage: stage}
	for _, a := range reg.Artifacts {
		kind := strings.ReplaceAll(a.Type, "_", "-")
		if allowed[kind] && strings.HasPrefix(filepath.ToSlash(a.Path), prefix) {
			r.Artifacts = append(r.Artifacts, a)
		}
	}
	sort.Slice(r.Artifacts, func(i, j int) bool { return r.Artifacts[i].Path < r.Artifacts[j].Path })
	if len(r.Artifacts) == 0 {
		r.Blockers = append(r.Blockers, "stage has no registered artifacts")
	}
	return r, nil
}

func ApproveStage(root, work, stage, by, notes string) ([]Approval, error) {
	review, err := ReviewStage(root, work, stage)
	if err != nil || len(review.Blockers) > 0 {
		return nil, fmt.Errorf("stage blocked: %s", strings.Join(review.Blockers, "; "))
	}
	type backup struct {
		path string
		data []byte
	}
	var backups []backup
	history := filepath.Join(root, ".product", "history")
	before, _ := os.ReadDir(history)
	known := map[string]bool{}
	for _, e := range before {
		known[e.Name()] = true
	}
	regPath := filepath.Join(root, ".product", "artifacts.json")
	regData, _ := os.ReadFile(regPath)
	for _, a := range review.Artifacts {
		p := filepath.Join(root, filepath.FromSlash(a.Path))
		d, e := os.ReadFile(p)
		if e != nil {
			return nil, e
		}
		backups = append(backups, backup{p, d})
	}
	rollback := func() {
		for _, b := range backups {
			_ = atomicWrite(b.path, b.data)
		}
		_ = atomicWrite(regPath, regData)
		after, _ := os.ReadDir(history)
		for _, e := range after {
			if !known[e.Name()] {
				_ = os.Remove(filepath.Join(history, e.Name()))
			}
		}
	}
	var out []Approval
	for _, a := range review.Artifacts {
		if a.Status == "approved" {
			continue
		}
		rec, e := Approve(root, filepath.Join(root, filepath.FromSlash(a.Path)), "approved", by, notes)
		if e != nil {
			rollback()
			return nil, e
		}
		out = append(out, rec)
	}
	return out, nil
}
