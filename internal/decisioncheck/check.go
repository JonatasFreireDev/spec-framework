package decisioncheck

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/decisions"
	"github.com/JonatasFreireDev/spec-framework/internal/validator"
)

type Options struct {
	Root, FrameworkRoot, Domain string
	FixLinks, Yes               bool
}

type Report struct {
	Root          string                 `json:"root"`
	Diagnostics   []validator.Diagnostic `json:"diagnostics"`
	ChangedFiles  []string               `json:"changed_files,omitempty"`
	DecisionCount int                    `json:"decision_count"`
}

var decisionID = regexp.MustCompile(`\bDEC-[0-9]+\b`)
var markdownLink = regexp.MustCompile(`\[([^\]\n]+)\]\(([^)\n]+)\)`)

func Run(options Options) (Report, error) {
	snap, err := validator.Scan(context.Background(), options.Root, options.FrameworkRoot, 0)
	if err != nil {
		return Report{}, err
	}
	report := inspect(snap, options.Domain)
	if options.FixLinks {
		if !options.Yes {
			return report, nil
		}
		changed, err := fixLinks(snap, options.Domain)
		if err != nil {
			return Report{}, err
		}
		report.ChangedFiles = changed
		if len(changed) > 0 {
			snap, err = validator.Scan(context.Background(), options.Root, options.FrameworkRoot, 0)
			if err != nil {
				return Report{}, err
			}
			report = inspect(snap, options.Domain)
			report.ChangedFiles = changed
		}
	}
	return report, nil
}

func inspect(s validator.Snapshot, domainFilter string) Report {
	r := Report{Root: s.Root}
	index, ok := s.JSON[".product/decisions.json"].(map[string]any)
	if !ok {
		r.Diagnostics = append(r.Diagnostics, validator.Diagnostic{Severity: validator.Error, Check: "decisions-index", File: ".product/decisions.json", Message: "Decision index is missing or invalid.", Fix: "Create .product/decisions.json with a decisions array."})
		return r
	}
	paths := decisions.DomainPaths(index)
	known := map[string]map[string]any{}
	knownPaths := map[string]string{}
	selectedIDs := map[string]bool{}
	for _, d := range stringMaps(index["decisions"]) {
		id := fmt.Sprint(d["id"])
		path := filepath.ToSlash(fmt.Sprint(d["path"]))
		if id == "<nil>" || id == "" {
			continue
		}
		domain := fmt.Sprint(d["domain"])
		if domain == "<nil>" || domain == "" {
			domain = decisions.DomainForPath(path, paths)
		}
		selected := domainFilter == "" || domain == domainFilter
		if selected {
			r.DecisionCount++
			selectedIDs[id] = true
		}
		if _, exists := known[id]; exists {
			if selected {
				r.Diagnostics = append(r.Diagnostics, diag(validator.Error, "decisions-duplicate", ".product/decisions.json", id+" is indexed more than once.", "Keep one index entry per decision id."))
			}
		}
		known[id] = d
		knownPaths[id] = path
		if selected && !fileExists(s, path) {
			r.Diagnostics = append(r.Diagnostics, diag(validator.Error, "decisions-path", path, id+" is indexed but its file does not exist.", "Create the file or correct the index path."))
		}
		if selected && (fmt.Sprint(d["domain"]) == "<nil>" || strings.TrimSpace(fmt.Sprint(d["domain"])) == "") {
			r.Diagnostics = append(r.Diagnostics, diag(validator.Warning, "decision-domain", path, id+" has no domain; it is using legacy path inference.", "Add domain to the index without moving the existing file."))
		}
		if selected && domain == "" {
			r.Diagnostics = append(r.Diagnostics, diag(validator.Error, "decision-domain", path, id+" has no registered domain.", "Add domain metadata or configure decisionDomains."))
		} else if selected {
			if err := decisions.ValidatePath(domain, path, paths); err != nil {
				r.Diagnostics = append(r.Diagnostics, diag(validator.Error, "decision-domain", path, id+" is stored in the wrong domain.", err.Error()))
			}
		}
		if selected {
			if issue := typeDomainIssue(fmt.Sprint(d["type"]), domain); issue != "" {
				r.Diagnostics = append(r.Diagnostics, diag(validator.Warning, "decision-type-domain", path, id+" has inconsistent type/domain: "+issue, "Align type with the decision owner or use cross-cutting."))
			}
		}
		if selected && fmt.Sprint(d["status"]) == "approved" {
			checkApproval(&r, s, id, path)
		}
	}
	for path := range s.Text {
		for domain, root := range paths {
			if !strings.HasPrefix(path, root) || filepath.Base(path) == "README.md" || !strings.HasSuffix(path, ".md") {
				continue
			}
			if domainFilter != "" && domain != domainFilter {
				continue
			}
			found := false
			for _, p := range knownPaths {
				if p == path {
					found = true
					break
				}
			}
			if !found {
				r.Diagnostics = append(r.Diagnostics, diag(validator.Warning, "decisions-unindexed", path, "Decision document is not indexed.", "Add it to .product/decisions.json before using it."))
			}
		}
	}
	for path, text := range s.Text {
		if !strings.HasSuffix(path, ".md") || isDecisionPath(path, paths) {
			continue
		}
		for _, id := range unique(decisionID.FindAllString(text, -1)) {
			if domainFilter != "" && !selectedIDs[id] {
				continue
			}
			d, exists := known[id]
			if !exists {
				r.Diagnostics = append(r.Diagnostics, diag(validator.Error, "decision-reference", path, "Unknown decision reference: "+id, "Index the decision or fix the reference."))
				continue
			}
			if !hasCanonicalLink(text, id, filepath.Dir(path), fmt.Sprint(d["path"]), s.Root) {
				r.Diagnostics = append(r.Diagnostics, diag(validator.Warning, "decision-links", path, id+" is referenced without a valid navigable link.", "Link the ID to its canonical path."))
			}
			if !affected(d, path) {
				r.Diagnostics = append(r.Diagnostics, diag(validator.Warning, "decision-scope", path, id+" is referenced but is not listed in affectedArtifacts.", "Add the artifact path to affectedArtifacts or remove the reference."))
			}
		}
	}
	sort.Slice(r.Diagnostics, func(i, j int) bool {
		a, b := r.Diagnostics[i], r.Diagnostics[j]
		if a.Severity != b.Severity {
			return a.Severity < b.Severity
		}
		if a.Check != b.Check {
			return a.Check < b.Check
		}
		return a.File < b.File
	})
	return r
}

func checkApproval(r *Report, s validator.Snapshot, id, path string) {
	text, ok := s.Text[path]
	if !ok {
		return
	}
	expected := validator.Hash(text)
	matched := false
	invalid := false
	for rel, value := range s.JSON {
		if !strings.HasPrefix(rel, ".product/history/approval-") {
			continue
		}
		rec, _ := value.(map[string]any)
		if fmt.Sprint(rec["artifact_id"]) != id {
			continue
		}
		invalid = true
		if fmt.Sprint(rec["path"]) == path && fmt.Sprint(rec["status_granted"]) == "approved" && fmt.Sprint(rec["content_hash"]) == expected {
			matched = true
		}
	}
	if !matched {
		msg := id + " has no current hash-matching approval record."
		check := "decision-approval-stale"
		if !invalid {
			check = "decision-approval-missing"
			msg = id + " is approved but has no approval record."
		}
		r.Diagnostics = append(r.Diagnostics, diag(validator.Error, check, path, msg, "Re-approve the current decision through the human approval flow."))
	}
}

func fileExists(s validator.Snapshot, path string) bool {
	_, ok := s.Text[filepath.ToSlash(path)]
	return ok
}

func fixLinks(s validator.Snapshot, domainFilter string) ([]string, error) {
	index, _ := s.JSON[".product/decisions.json"].(map[string]any)
	paths := decisions.DomainPaths(index)
	canonical := map[string]string{}
	for _, d := range stringMaps(index["decisions"]) {
		id := fmt.Sprint(d["id"])
		path := filepath.ToSlash(fmt.Sprint(d["path"]))
		domain := fmt.Sprint(d["domain"])
		if domain == "<nil>" || domain == "" {
			domain = decisions.DomainForPath(path, paths)
		}
		if domainFilter == "" || domain == domainFilter {
			canonical[id] = path
		}
	}
	var changed []string
	for path, text := range s.Text {
		if !strings.HasSuffix(path, ".md") || isDecisionPath(path, paths) {
			continue
		}
		original := text
		lines := strings.Split(text, "\n")
		inFence := false
		inFrontmatter := len(lines) > 0 && strings.TrimSpace(lines[0]) == "---"
		for i, line := range lines {
			if i > 0 && strings.TrimSpace(line) == "---" && inFrontmatter {
				inFrontmatter = false
				continue
			}
			if inFrontmatter {
				continue
			}
			if strings.HasPrefix(strings.TrimSpace(line), "```") {
				inFence = !inFence
				continue
			}
			if inFence {
				continue
			}
			for id, target := range canonical {
				if !strings.Contains(line, id) || hasCanonicalLink(line, id, filepath.Dir(path), target, s.Root) {
					continue
				}
				idx := strings.Index(line, id)
				if idx < 0 || strings.Contains(line[:idx], "[") && strings.Contains(line[idx:], "](") {
					continue
				}
				rel, _ := filepath.Rel(filepath.Dir(filepath.Join(s.Root, filepath.FromSlash(path))), filepath.Join(s.Root, filepath.FromSlash(target)))
				rel = filepath.ToSlash(rel)
				line = strings.Replace(line, id, "["+id+"]("+rel+")", 1)
				lines[i] = line
			}
		}
		text = strings.Join(lines, "\n")
		if text != original {
			if err := os.WriteFile(filepath.Join(s.Root, filepath.FromSlash(path)), []byte(text), 0644); err != nil {
				return nil, err
			}
			changed = append(changed, path)
		}
	}
	sort.Strings(changed)
	return changed, nil
}

func hasCanonicalLink(text, id, fromDir, target, root string) bool {
	for _, m := range markdownLink.FindAllStringSubmatch(text, -1) {
		if !strings.Contains(m[1], id) {
			continue
		}
		parts := strings.SplitN(m[2], "#", 2)
		candidate := filepath.Clean(filepath.Join(root, filepath.FromSlash(fromDir), filepath.FromSlash(parts[0])))
		if filepath.ToSlash(candidate) == filepath.ToSlash(filepath.Join(root, filepath.FromSlash(target))) {
			if len(parts) == 2 {
				data, err := os.ReadFile(candidate)
				if err != nil || !markdownAnchors(string(data))[parts[1]] {
					continue
				}
			}
			return true
		}
	}
	return false
}

var heading = regexp.MustCompile(`(?m)^#{1,6}\s+(.+?)\s*#*\s*$`)

func markdownAnchors(text string) map[string]bool {
	out := map[string]bool{}
	for _, match := range heading.FindAllStringSubmatch(text, -1) {
		value := strings.ToLower(strings.TrimSpace(match[1]))
		value = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(value, "")
		value = regexp.MustCompile(`[^\pL\pN\s-]`).ReplaceAllString(value, "")
		value = regexp.MustCompile(`[\s-]+`).ReplaceAllString(value, "-")
		out[strings.Trim(value, "-")] = true
	}
	return out
}

func isDecisionPath(path string, roots map[string]string) bool {
	for _, root := range roots {
		if strings.HasPrefix(path, root) {
			return true
		}
	}
	return false
}
func affected(d map[string]any, path string) bool {
	for _, v := range stringSlice(d["affectedArtifacts"]) {
		v = filepath.ToSlash(v)
		if strings.HasPrefix(v, "../") || v == "product" || strings.HasSuffix(v, "/") && strings.HasPrefix(path, v) || v == path {
			return true
		}
	}
	return false
}
func typeDomainIssue(typ, domain string) string {
	allowed := map[string]map[string]bool{"product": {"product": true, "cross-cutting": true}, "architecture": {"design": true, "engineering": true, "cross-cutting": true}, "security": {"product": true, "design": true, "engineering": true, "cross-cutting": true}, "data": {"engineering": true, "cross-cutting": true}, "delivery": {"product": true, "cross-cutting": true}}
	if typ == "<nil>" || typ == "" || domain == "" {
		return ""
	}
	if !allowed[typ][domain] {
		return typ + " does not normally belong to " + domain
	}
	return ""
}
func diag(sev validator.Severity, check, file, message, fix string) validator.Diagnostic {
	return validator.Diagnostic{Severity: sev, Check: check, File: file, Message: message, Fix: fix}
}
func stringMaps(v any) []map[string]any {
	var out []map[string]any
	if raw, ok := v.([]any); ok {
		for _, x := range raw {
			if m, ok := x.(map[string]any); ok {
				out = append(out, m)
			}
		}
	}
	return out
}
func stringSlice(v any) []string {
	var out []string
	if raw, ok := v.([]any); ok {
		for _, x := range raw {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
	}
	return out
}
func unique(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
