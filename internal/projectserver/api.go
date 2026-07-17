package projectserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

// projectView is the UI contract. It deliberately translates framework records
// into product-facing information without adding presentation fields to workflow.Artifact.
type projectView struct {
	Revision  uint64               `json:"revision"`
	Generated string               `json:"generatedAt"`
	Metrics   viewMetrics          `json:"metrics"`
	Artifacts []artifactView       `json:"artifacts"`
	Files     []worktreeFile       `json:"files"`
	Types     []artifactTypeConfig `json:"types"`
	Git       gitState             `json:"git"`
}

type viewMetrics struct {
	Total      int `json:"total"`
	Approved   int `json:"approved"`
	Pending    int `json:"pending"`
	Rejected   int `json:"rejected"`
	InProgress int `json:"inProgress"`
	Blocked    int `json:"blocked"`
	Stale      int `json:"stale"`
	Untracked  int `json:"untracked"`
	Changed    int `json:"changed"`
}
type artifactView struct {
	ID              string             `json:"id"`
	Type            string             `json:"type"`
	Title           string             `json:"title"`
	Path            string             `json:"path"`
	Folder          string             `json:"folder"`
	Status          string             `json:"status"`
	Maturity        string             `json:"maturity,omitempty"`
	TargetFeature   string             `json:"targetFeature,omitempty"`
	ApprovalAdapter string             `json:"approvalAdapter,omitempty"`
	Parents         []string           `json:"parents"`
	Children        []string           `json:"children"`
	Updated         string             `json:"updated,omitempty"`
	Content         string             `json:"content,omitempty"`
	DerivedStatus   string             `json:"derivedStatus"`
	Blockers        []string           `json:"blockers"`
	LatestApproval  *approvalSummary   `json:"latestApproval,omitempty"`
	View            artifactTypeConfig `json:"view"`
}
type approvalSummary struct {
	By     string `json:"by"`
	At     string `json:"at"`
	Status string `json:"status"`
	Notes  string `json:"notes,omitempty"`
}
type worktreeFile struct {
	Path       string `json:"path"`
	Folder     string `json:"folder"`
	Extension  string `json:"extension,omitempty"`
	Registered bool   `json:"registered"`
	Changed    bool   `json:"changed"`
	State      string `json:"state"`
	Updated    string `json:"updated"`
	Error      string `json:"error,omitempty"`
}
type artifactTypeConfig struct {
	Type     string   `json:"type"`
	Label    string   `json:"label"`
	Renderer string   `json:"renderer"`
	Tabs     []string `json:"tabs"`
	Fields   []string `json:"fields"`
}
type gitState struct {
	Available    bool     `json:"available"`
	Branch       string   `json:"branch,omitempty"`
	ChangedPaths []string `json:"changedPaths"`
}

// Legacy status is retained for existing clients; the new UI must use /api/project-view.
type projectStatus struct {
	Documents []document `json:"documents"`
	Metrics   metrics    `json:"metrics"`
}
type metrics struct {
	Total    int `json:"total"`
	Approved int `json:"approved"`
	Pending  int `json:"pending"`
	Rejected int `json:"rejected"`
}
type document struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Path    string `json:"path"`
	Folder  string `json:"folder"`
	Status  string `json:"status"`
	Updated string `json:"updated"`
	Content string `json:"content"`
}
type transitionRequest struct {
	ArtifactID string `json:"artifactId"`
	Status     string `json:"status"`
	ApprovedBy string `json:"approvedBy"`
	Notes      string `json:"notes"`
	Confirmed  bool   `json:"confirmed"`
}
type batchApprovalRequest struct {
	ArtifactIDs []string `json:"artifactIds"`
	ApprovedBy  string   `json:"approvedBy"`
	Notes       string   `json:"notes"`
	Confirmed   bool     `json:"confirmed"`
}

type userPreferences struct {
	ReviewerName string `json:"reviewerName"`
}

func readPreferences(root string) (userPreferences, error) {
	data, err := os.ReadFile(filepath.Join(root, ".product", "ui-preferences.json"))
	if errors.Is(err, os.ErrNotExist) {
		return userPreferences{}, nil
	}
	if err != nil {
		return userPreferences{}, err
	}
	var preferences userPreferences
	if err := json.Unmarshal(data, &preferences); err != nil {
		return userPreferences{}, fmt.Errorf("read UI preferences: %w", err)
	}
	return preferences, nil
}

func writePreferences(root string, preferences userPreferences) error {
	dir := filepath.Join(root, ".product")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(preferences, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(dir, ".ui-preferences-*.tmp")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err := temporary.Write(append(data, '\n')); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(name, filepath.Join(dir, "ui-preferences.json"))
}

func readProjectView(root string, revision uint64) (projectView, error) {
	registry, err := workflow.LoadRegistry(root)
	if err != nil {
		return projectView{}, fmt.Errorf("read artifact registry: %w", err)
	}
	git := readGitState(root)
	changed := map[string]bool{}
	for _, path := range git.ChangedPaths {
		changed[path] = true
	}
	latest := readApprovalHistory(root)
	stale := readStaleness(root, registry)
	byID := map[string]workflow.Artifact{}
	children := map[string][]string{}
	registered := map[string]bool{}
	for _, a := range registry.Artifacts {
		byID[a.ID] = a
		registered[filepath.ToSlash(a.Path)] = true
		for _, p := range a.ParentIDs {
			children[p] = append(children[p], a.ID)
		}
	}
	view := projectView{Revision: revision, Generated: time.Now().UTC().Format(time.RFC3339), Artifacts: make([]artifactView, 0, len(registry.Artifacts)), Types: allTypeConfigs(), Git: git}
	for _, a := range registry.Artifacts {
		path := filepath.Join(root, filepath.FromSlash(a.Path))
		content, readErr := os.ReadFile(path)
		info, statErr := os.Stat(path)
		item := artifactView{ID: a.ID, Type: a.Type, Title: a.ID, Path: filepath.ToSlash(a.Path), Folder: folderFor(a.Path), Status: a.Status, Maturity: a.Maturity, TargetFeature: a.TargetFeature, ApprovalAdapter: a.ApprovalAdapter, Parents: a.ParentIDs, Children: children[a.ID], DerivedStatus: "current", Blockers: []string{}, View: typeConfig(a.Type)}
		if readErr != nil || statErr != nil {
			item.DerivedStatus = "missing"
			item.Blockers = []string{"artifact file is missing or unreadable"}
		} else {
			item.Content = string(content)
			item.Title = titleFor(a, content)
			item.Updated = info.ModTime().UTC().Format(time.RFC3339)
		}
		if stale[a.ID] {
			item.DerivedStatus = "stale"
		}
		if a.Status == "rejected" {
			item.DerivedStatus = "rejected"
		}
		if len(item.Blockers) > 0 {
			item.DerivedStatus = "blocked"
		}
		for _, parent := range a.ParentIDs {
			if p, ok := byID[parent]; !ok || p.Status != "approved" {
				item.Blockers = append(item.Blockers, "parent "+parent+" is not approved")
			}
		}
		if len(item.Blockers) > 0 && item.DerivedStatus == "current" {
			item.DerivedStatus = "blocked"
		}
		if x, ok := latest[a.ID]; ok {
			item.LatestApproval = &x
		}
		view.Artifacts = append(view.Artifacts, item)
		view.Metrics.Total++
		switch a.Status {
		case "approved":
			view.Metrics.Approved++
		case "rejected":
			view.Metrics.Rejected++
		case "in_progress", "implemented", "validated":
			view.Metrics.InProgress++
		default:
			view.Metrics.Pending++
		}
		if item.DerivedStatus == "stale" {
			view.Metrics.Stale++
		}
		if item.DerivedStatus == "blocked" {
			view.Metrics.Blocked++
		}
	}
	files, err := scanWorktree(root, registered, changed)
	if err != nil {
		return projectView{}, err
	}
	view.Files = files
	for _, f := range files {
		if !f.Registered {
			view.Metrics.Untracked++
		}
		if f.Changed {
			view.Metrics.Changed++
		}
	}
	sort.Slice(view.Artifacts, func(i, j int) bool { return view.Artifacts[i].Path < view.Artifacts[j].Path })
	return view, nil
}

func readStatus(root string) (projectStatus, error) {
	v, err := readProjectView(root, 0)
	if err != nil {
		return projectStatus{}, err
	}
	out := projectStatus{}
	for _, a := range v.Artifacts {
		out.Documents = append(out.Documents, document{a.ID, a.Title, a.Path, a.Folder, a.Status, a.Updated, a.Content})
		out.Metrics.Total++
		switch a.Status {
		case "approved":
			out.Metrics.Approved++
		case "rejected":
			out.Metrics.Rejected++
		default:
			out.Metrics.Pending++
		}
	}
	return out, nil
}
func scanWorktree(root string, registered, changed map[string]bool) ([]worktreeFile, error) {
	out := []worktreeFile{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(rel, ".git/") {
			return nil
		}
		info, e := d.Info()
		f := worktreeFile{Path: rel, Folder: folderFor(rel), Extension: strings.TrimPrefix(filepath.Ext(rel), "."), Registered: registered[rel], Changed: changed[rel], State: "registered", Updated: info.ModTime().UTC().Format(time.RFC3339)}
		if !f.Registered {
			f.State = "untracked"
		}
		if f.Changed {
			f.State = "changed"
		}
		if e != nil {
			f.Error = e.Error()
			f.State = "unreadable"
		}
		out = append(out, f)
		return nil
	})
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, err
}
func readApprovalHistory(root string) map[string]approvalSummary {
	out := map[string]approvalSummary{}
	entries, err := os.ReadDir(filepath.Join(root, ".product", "history"))
	if err != nil {
		return out
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, ".product", "history", e.Name()))
		if err != nil {
			continue
		}
		var a workflow.Approval
		if json.Unmarshal(data, &a) != nil || a.ArtifactID == "" {
			continue
		}
		if old, ok := out[a.ArtifactID]; !ok || a.ApprovedAt > old.At {
			out[a.ArtifactID] = approvalSummary{a.ApprovedBy, a.ApprovedAt, a.StatusGranted, a.Notes}
		}
	}
	return out
}
func readStaleness(root string, r workflow.Registry) map[string]bool {
	result := map[string]bool{}
	data, err := os.ReadFile(filepath.Join(root, ".product", "derivations.json"))
	if err != nil {
		return result
	}
	var raw struct {
		Derivations []struct {
			ArtifactID  string `json:"artifact_id"`
			DerivedFrom []struct {
				ArtifactID  string `json:"artifact_id"`
				Path        string `json:"path"`
				ContentHash string `json:"content_hash"`
			} `json:"derived_from"`
		} `json:"derivations"`
	}
	if json.Unmarshal(data, &raw) != nil {
		return result
	}
	for _, d := range raw.Derivations {
		for _, s := range d.DerivedFrom {
			data, e := os.ReadFile(filepath.Join(root, filepath.FromSlash(s.Path)))
			if e != nil || workflow.Hash(string(data)) != s.ContentHash {
				result[d.ArtifactID] = true
			}
		}
	}
	return result
}
func readGitState(root string) gitState {
	out := gitState{ChangedPaths: []string{}}
	cmd := exec.Command("git", "-C", root, "status", "--porcelain=v1", "-b")
	data, err := cmd.Output()
	if err != nil {
		return out
	}
	out.Available = true
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.HasPrefix(line, "## ") {
			out.Branch = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			continue
		}
		if len(line) >= 4 {
			out.ChangedPaths = append(out.ChangedPaths, filepath.ToSlash(strings.TrimSpace(line[3:])))
		}
	}
	return out
}
func allTypeConfigs() []artifactTypeConfig {
	keys := []string{"problem", "vision", "product-principles", "north-star", "strategy", "feature-brief", "product-baseline", "implementation-assessment", "domain", "user-goal", "feature", "use-case", "design", "design-system", "engineering-system", "engineering-proposal", "engineering-review", "tests", "qa-evidence", "security-review", "audit", "execution-graph", "taskset", "task", "decision", "release", "security-baseline"}
	out := make([]artifactTypeConfig, 0, len(keys))
	for _, key := range keys {
		out = append(out, typeConfig(key))
	}
	return out
}
func typeConfig(kind string) artifactTypeConfig {
	c := artifactTypeConfig{Type: kind, Label: strings.ReplaceAll(kind, "-", " "), Renderer: "markdown", Tabs: []string{"content", "details", "relations", "history"}, Fields: []string{"status", "maturity", "approval"}}
	switch kind {
	case "execution-graph":
		c.Renderer = "graph"
		c.Tabs = []string{"graph", "details", "relations", "history"}
		c.Fields = []string{"nodes", "dependencies", "writeScope"}
	case "design-system":
		c.Renderer = "design-system"
		c.Tabs = []string{"overview", "tokens", "themes", "components", "history"}
		c.Fields = []string{"tokens", "themes", "version"}
	case "design":
		c.Renderer = "design"
		c.Fields = []string{"maturity", "fidelity", "sources", "screens", "mappings"}
	case "task", "taskset":
		c.Renderer = "task"
		c.Fields = []string{"dependencies", "acceptanceCriteria", "writeScope", "evidence"}
	case "decision":
		c.Renderer = "decision"
		c.Fields = []string{"scope", "affectedArtifacts", "workflowEffects"}
	case "qa-evidence", "security-review", "audit":
		c.Renderer = "evidence"
		c.Fields = []string{"findings", "evidence", "verdict"}
	}
	return c
}
func findArtifact(root, id string) (workflow.Artifact, error) {
	r, err := workflow.LoadRegistry(root)
	if err != nil {
		return workflow.Artifact{}, err
	}
	for _, a := range r.Artifacts {
		if a.ID == id {
			return a, nil
		}
	}
	return workflow.Artifact{}, errors.New("artifact not found")
}
func titleFor(a workflow.Artifact, content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return a.ID
}
func folderFor(path string) string {
	f := filepath.ToSlash(filepath.Dir(path))
	if f == "." {
		return "Raiz"
	}
	return f
}
func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	d := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	d.DisallowUnknownFields()
	if d.Decode(target) != nil {
		return errors.New("invalid request body")
	}
	return nil
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeAPIError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
func methodNotAllowed(w http.ResponseWriter) {
	writeAPIError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
}
