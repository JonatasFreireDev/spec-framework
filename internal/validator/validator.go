package validator

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Note    Severity = "note"
)

type Diagnostic struct {
	Severity                  Severity `json:"severity"`
	Check, File, Message, Fix string
}
type Snapshot struct {
	Root, FrameworkRoot string
	Files               []string
	Text                map[string]string
	JSON                map[string]any
}
type Result struct {
	Diagnostics             []Diagnostic
	Errors, Warnings, Notes int
}

func Scan(ctx context.Context, root, frameworkRoot string, workers int) (Snapshot, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Snapshot{}, err
	}
	frameworkRoot, err = filepath.Abs(frameworkRoot)
	if err != nil {
		return Snapshot{}, err
	}
	var files []string
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return Snapshot{}, err
	}
	sort.Strings(files)
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	if workers > 8 {
		workers = 8
	}
	if workers > len(files) {
		workers = len(files)
	}
	if workers < 1 {
		workers = 1
	}
	type item struct {
		index int
		text  string
		json  any
		err   error
	}
	jobs := make(chan int)
	results := make(chan item, len(files))
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				select {
				case <-ctx.Done():
					results <- item{index: index, err: ctx.Err()}
					continue
				default:
				}
				data, err := os.ReadFile(files[index])
				it := item{index: index, err: err, text: strings.TrimPrefix(string(data), "\ufeff")}
				if err == nil && strings.HasSuffix(files[index], ".json") {
					_ = json.Unmarshal(bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf}), &it.json)
				}
				results <- it
			}
		}()
	}
	go func() {
		for i := range files {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()
	text := map[string]string{}
	jsonFiles := map[string]any{}
	for result := range results {
		if result.err != nil {
			return Snapshot{}, result.err
		}
		rel, _ := filepath.Rel(root, files[result.index])
		rel = filepath.ToSlash(rel)
		text[rel] = result.text
		if result.json != nil {
			jsonFiles[rel] = result.json
		}
	}
	return Snapshot{Root: root, FrameworkRoot: frameworkRoot, Files: files, Text: text, JSON: jsonFiles}, nil
}

func Validate(ctx context.Context, root, frameworkRoot string) (Result, error) {
	snap, err := Scan(ctx, root, frameworkRoot, 0)
	if err != nil {
		return Result{}, err
	}
	var d []Diagnostic
	for rel, text := range snap.Text {
		if strings.HasPrefix(rel, "domains/") && strings.HasSuffix(rel, "context.md") {
			d = append(d, validateContextFull(rel, text)...)
		}
		if strings.HasSuffix(rel, "execution-graph.json") {
			d = append(d, validateGraph(rel, snap.JSON[rel], snap)...)
		}
	}
	d = append(d, validateApprovalRecords(snap)...)
	d = append(d, validateMarkdownLinks(snap)...)
	d = append(d, validateUseCaseBundles(snap)...)
	d = append(d, validateIdentity(snap)...)
	d = append(d, validateEvidence(snap)...)
	d = append(d, validateQualityGates(snap)...)
	d = append(d, validateStatusAndStaleness(snap)...)
	d = append(d, validateDecisions(snap)...)
	d = append(d, validateSkillReferences(snap)...)
	d = append(d, validateDeliveryAndRigor(snap)...)
	d = append(d, validateRegistryAndApprovalGates(snap)...)
	d = append(d, validateImportRuns(snap)...)
	sort.Slice(d, func(i, j int) bool {
		a, b := d[i], d[j]
		if a.Severity != b.Severity {
			return rank(a.Severity) < rank(b.Severity)
		}
		if a.Check != b.Check {
			return a.Check < b.Check
		}
		if a.File != b.File {
			return a.File < b.File
		}
		return a.Message < b.Message
	})
	r := Result{Diagnostics: d}
	for _, x := range d {
		switch x.Severity {
		case Error:
			r.Errors++
		case Warning:
			r.Warnings++
		case Note:
			r.Notes++
		}
	}
	return r, nil
}
func validateContext(file, text string) []Diagnostic {
	var out []Diagnostic
	for _, field := range []string{"id:", "type:", "name:", "status:", "owner_skill:", "slug:"} {
		if !containsLinePrefix(text, field) {
			out = append(out, Diagnostic{Error, "contexts", file, "Missing required context field: " + strings.TrimSuffix(field, ":"), "Add the field to context.md."})
		}
	}
	return out
}
func validateGraph(file string, value any, snap Snapshot) []Diagnostic {
	object, ok := value.(map[string]any)
	if !ok {
		return []Diagnostic{{Error, "execution-graph", file, "Invalid JSON execution graph", "Write a valid JSON object."}}
	}
	nodes, ok := object["nodes"].([]any)
	if !ok {
		return []Diagnostic{{Error, "execution-graph", file, "Execution graph must contain nodes[]", "Add a nodes array."}}
	}
	ids := map[string]bool{}
	objects := map[string]map[string]any{}
	var out []Diagnostic
	for _, raw := range nodes {
		node, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id, _ := node["id"].(string)
		if id == "" {
			out = append(out, Diagnostic{Error, "execution-graph", file, "Node is missing id", "Add a stable task id."})
		} else if ids[id] {
			out = append(out, Diagnostic{Error, "execution-graph", file, "Duplicate node id: " + id, "Use unique ids."})
		}
		ids[id] = true
		objects[id] = node
		if path, _ := node["path"].(string); path == "" {
			out = append(out, Diagnostic{Error, "execution-graph", file, fmt.Sprintf("Node %s is missing path", id), "Point to tasks/<task-id>.md."})
		}
	}
	base := filepath.ToSlash(filepath.Dir(file))
	for id, node := range objects {
		path, _ := node["path"].(string)
		if path != "" {
			full := filepath.ToSlash(filepath.Join(base, path))
			if _, ok := snap.Text[full]; !ok {
				out = append(out, Diagnostic{Error, "execution-graph", file, fmt.Sprintf("Node %s path does not exist: %s", id, path), "Create the canonical task file."})
			}
		}
		deps, ok := node["dependsOn"].([]any)
		if !ok {
			out = append(out, Diagnostic{Error, "execution-graph", file, fmt.Sprintf("Node %s dependsOn must be an array.", id), "Set dependsOn to an array of task ids."})
			continue
		}
		for _, raw := range deps {
			dep, _ := raw.(string)
			if !ids[dep] {
				out = append(out, Diagnostic{Error, "execution-graph", file, fmt.Sprintf("Node %s depends on missing node %s.", id, dep), "Add or remove the dependency."})
			}
		}
		if _, ok := node["writeScope"].([]any); !ok {
			out = append(out, Diagnostic{Error, "execution-graph", file, fmt.Sprintf("Node %s writeScope must be an array.", id), "Declare concrete write scopes."})
		}
	}
	ordered := make([]string, 0, len(objects))
	for id := range objects {
		ordered = append(ordered, id)
	}
	sort.Strings(ordered)
	for i, left := range ordered {
		for _, right := range ordered[i+1:] {
			if dependencyPath(left, right, objects, map[string]bool{}) || dependencyPath(right, left, objects, map[string]bool{}) {
				continue
			}
			for _, a := range scopes(objects[left]) {
				for _, b := range scopes(objects[right]) {
					if scopeOverlap(a, b) {
						out = append(out, Diagnostic{Warning, "write-scope", file, fmt.Sprintf("Parallel nodes %s and %s have overlapping writeScope: %s <> %s.", left, right, a, b), "Add a dependency or separate write scopes."})
					}
				}
			}
		}
	}
	return out
}

func scopes(node map[string]any) []string {
	raw, _ := node["writeScope"].([]any)
	var out []string
	for _, v := range raw {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, strings.Trim(filepath.ToSlash(filepath.Clean(s)), "/"))
		}
	}
	return out
}
func scopeOverlap(a, b string) bool {
	return a == b || strings.HasPrefix(a, b+"/") || strings.HasPrefix(b, a+"/")
}
func dependencyPath(from, to string, nodes map[string]map[string]any, seen map[string]bool) bool {
	if from == to {
		return true
	}
	if seen[from] {
		return false
	}
	seen[from] = true
	deps, _ := nodes[from]["dependsOn"].([]any)
	for _, raw := range deps {
		dep, _ := raw.(string)
		if dep == to || dependencyPath(dep, to, nodes, seen) {
			return true
		}
	}
	return false
}
func containsLinePrefix(text, prefix string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(strings.ToLower(line)), prefix) {
			return true
		}
	}
	return false
}
func rank(s Severity) int {
	if s == Error {
		return 0
	}
	if s == Warning {
		return 1
	}
	return 2
}

var mdLink = regexp.MustCompile(`(?m)(?P<image>!)?\[[^\]\n]+\]\(([^)\n]+)\)`)

func validateMarkdownLinks(s Snapshot) []Diagnostic {
	var out []Diagnostic
	for rel, text := range s.Text {
		if !strings.HasSuffix(rel, ".md") {
			continue
		}
		text = regexp.MustCompile("(?s)```.*?```").ReplaceAllString(text, "")
		matches := mdLink.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			target := strings.TrimSpace(match[2])
			if target == "" || strings.HasPrefix(target, "#") || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") || strings.Contains(target, "<") {
				continue
			}
			if i := strings.Index(target, "#"); i >= 0 {
				target = target[:i]
			}
			decoded, err := url.PathUnescape(strings.Trim(target, "<>"))
			if err == nil {
				target = decoded
			}
			candidate := filepath.Clean(filepath.Join(s.Root, filepath.Dir(filepath.FromSlash(rel)), filepath.FromSlash(target)))
			if _, err := os.Stat(candidate); err != nil {
				out = append(out, Diagnostic{Error, "links", rel, "Broken Markdown link: " + target, "Create the target or update the link."})
			}
		}
	}
	return out
}

func validateApprovalRecords(s Snapshot) []Diagnostic {
	registry, ok := s.JSON[".product/artifacts.json"].(map[string]any)
	if !ok {
		return nil
	}
	items, _ := registry["artifacts"].([]any)
	records := map[string][]map[string]any{}
	for rel, value := range s.JSON {
		if !strings.HasPrefix(rel, ".product/history/approval-") {
			continue
		}
		if record, ok := value.(map[string]any); ok {
			id, _ := record["artifact_id"].(string)
			records[id] = append(records[id], record)
		}
	}
	var out []Diagnostic
	for _, raw := range items {
		artifact, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id, _ := artifact["id"].(string)
		status, _ := artifact["status"].(string)
		path, _ := artifact["path"].(string)
		if !requiresApproval(status) || id == "" || path == "" {
			continue
		}
		text, exists := s.Text[filepath.ToSlash(path)]
		if !exists {
			continue
		}
		expected := Hash(text)
		matched := false
		for _, record := range records[id] {
			if record["path"] == path && record["status_granted"] == status && record["content_hash"] == expected {
				matched = true
				break
			}
		}
		if !matched {
			out = append(out, Diagnostic{Error, "approval-records", path, fmt.Sprintf("%s is %s, but no matching approval record exists in .product/history/.", id, status), "Do not auto-fix approval records. Ask the approving human to create a matching record."})
		}
	}
	return out
}
func requiresApproval(status string) bool {
	switch status {
	case "approved", "in_progress", "implemented", "validated", "released":
		return true
	}
	return false
}
func (r Result) Verdict() string {
	if r.Errors > 0 {
		return "blocked"
	}
	if r.Warnings > 0 {
		return "ready_with_warnings"
	}
	return "ready"
}
func Hash(text string) string {
	normalized := strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(normalized, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	normalized = strings.Join(lines, "\n")
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}
