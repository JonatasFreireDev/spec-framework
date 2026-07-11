package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var allowedStatuses = map[string]bool{"draft": true, "proposed": true, "approved": true, "in_progress": true, "implemented": true, "validated": true, "released": true, "deprecated": true, "superseded": true, "rejected": true, "not_applicable": true}
var contextFieldPattern = regexp.MustCompile(`(?m)^\s*([a-zA-Z_]+):\s*(.*?)\s*$`)
var tableRowPattern = regexp.MustCompile(`(?m)^\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*$`)

func metadata(text string) map[string]string {
	out := map[string]string{}
	for _, m := range contextFieldPattern.FindAllStringSubmatch(text, -1) {
		key := strings.ToLower(strings.TrimSpace(m[1]))
		if _, exists := out[key]; !exists {
			out[key] = strings.Trim(strings.TrimSpace(m[2]), `"'`)
		}
	}
	return out
}
func tableFields(text string) map[string]string {
	out := map[string]string{}
	for _, m := range tableRowPattern.FindAllStringSubmatch(text, -1) {
		key := strings.ToLower(strings.TrimSpace(m[1]))
		if key == "field" || strings.Trim(m[1], " -") == "" {
			continue
		}
		out[key] = strings.TrimSpace(m[2])
	}
	return out
}

func validateContextFull(file, text string) []Diagnostic {
	meta := metadata(text)
	var out []Diagnostic
	for _, field := range []string{"id", "type", "name", "status", "owner_skill", "slug"} {
		if meta[field] == "" {
			check := "context"
			if field == "slug" {
				check = "identity"
			}
			out = append(out, Diagnostic{Error, check, file, "Missing context field: " + field + ".", "Add " + field + "."})
		}
	}
	if slug := meta["slug"]; slug != "" && slug != filepath.Base(filepath.Dir(filepath.FromSlash(file))) {
		out = append(out, Diagnostic{Error, "identity", file, fmt.Sprintf("Context slug %s does not match folder %s.", slug, filepath.Base(filepath.Dir(file))), "Keep slug equal to the immutable folder name."})
	}
	if status := meta["status"]; status != "" && !allowedStatuses[status] {
		out = append(out, Diagnostic{Error, "context", file, "Invalid status: " + status + ".", "Use a framework-approved status."})
	}
	if !strings.Contains(text, "## Handoff") {
		out = append(out, Diagnostic{Warning, "context", file, "Missing Handoff section.", "Add next skill and required reading."})
	}
	return out
}

func validateUseCaseBundles(s Snapshot) []Diagnostic {
	dirs := map[string]bool{}
	for rel := range s.Text {
		parts := strings.Split(filepath.ToSlash(rel), "/")
		for i := 0; i+2 < len(parts); i++ {
			if parts[i] == "use-cases" && parts[i+2] == "context.md" && i+3 == len(parts) {
				dirs[strings.Join(parts[:i+2], "/")] = true
			}
		}
	}
	base := []string{"context.md", "use-case.md", "specification.md", "tasks.md", "tests.md"}
	tierFiles := map[string][]string{"M": {"design.md", "implementation-plan.md", "execution-graph.json"}, "L": {"design.md", "implementation-plan.md", "execution-graph.json", "analytics.md", "audit.md", "qa-evidence.md", "security-review.md"}, "N/A": {"design.md", "implementation-plan.md", "execution-graph.json", "analytics.md", "audit.md"}}
	var out []Diagnostic
	var ordered []string
	for dir := range dirs {
		ordered = append(ordered, dir)
	}
	sort.Strings(ordered)
	for _, dir := range ordered {
		contextText := s.Text[dir+"/context.md"]
		tier := strings.ToUpper(metadata(contextText)["rigor_tier"])
		if tier == "" {
			out = append(out, Diagnostic{Error, "rigor-tier", dir + "/context.md", "Missing rigor_tier.", "Set rigor_tier to S, M, L, or N/A."})
		}
		required := append([]string{}, base...)
		required = append(required, tierFiles[tier]...)
		for _, name := range required {
			path := dir + "/" + name
			if _, ok := s.Text[path]; !ok {
				if _, jsonOK := s.JSON[path]; !jsonOK {
					out = append(out, Diagnostic{Error, "use-case-bundle", path, "Required use-case file is missing.", "Create " + name + " from the framework template."})
				}
			}
		}
	}
	return out
}

func validateIdentity(s Snapshot) []Diagnostic {
	hasProduct := false
	for rel := range s.Text {
		if strings.HasPrefix(rel, ".product/") || strings.HasPrefix(rel, "domains/") {
			hasProduct = true
			break
		}
	}
	if !hasProduct {
		return nil
	}
	value, ok := s.JSON[".product/ids.json"].(map[string]any)
	if !ok {
		return []Diagnostic{{Error, "identity", ".product/ids.json", ".product/ids.json is missing or invalid.", "Add identity policy metadata."}}
	}
	var out []Diagnostic
	if value["policy"] != "slug-scoped" {
		out = append(out, Diagnostic{Error, "identity", ".product/ids.json", "Expected slug-scoped identity policy.", "Set policy to slug-scoped."})
	}
	if value["deprecated_counters"] != true {
		out = append(out, Diagnostic{Error, "identity", ".product/ids.json", "Central counters must be deprecated.", "Set deprecated_counters to true."})
	}
	for key, item := range value {
		switch item.(type) {
		case float64:
			out = append(out, Diagnostic{Error, "identity", ".product/ids.json", "Central numeric counter remains: " + key, "Remove central counters."})
		}
	}
	return out
}

func validateEvidence(s Snapshot) []Diagnostic {
	registry, ok := s.JSON[".product/artifacts.json"].(map[string]any)
	if !ok {
		return nil
	}
	items, _ := registry["artifacts"].([]any)
	var out []Diagnostic
	for _, raw := range items {
		artifact, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		kind, _ := artifact["type"].(string)
		status, _ := artifact["status"].(string)
		path, _ := artifact["path"].(string)
		text := s.Text[filepath.ToSlash(path)]
		fields := tableFields(text)
		if kind == "task" && (status == "implemented" || status == "validated" || status == "released") {
			for _, field := range []string{"branch", "commits", "code paths"} {
				if placeholder(fields[field]) {
					out = append(out, Diagnostic{Error, "code-evidence", path, "Task is " + status + " but has no concrete " + field + ".", "Record structured implementation evidence."})
				}
			}
			if commit := fields["commits"]; commit != "" && !regexp.MustCompile(`(?i)([0-9a-f]{7,40}|https?://\S+/commit/\S+)`).MatchString(commit) {
				out = append(out, Diagnostic{Error, "code-evidence", path, "Commits is not a traceable commit hash or commit URL.", "Record a commit hash or URL."})
			}
		}
		if kind == "task" && (status == "validated" || status == "released") {
			for _, field := range []string{"pr", "test status", "gate logs", "qa evidence"} {
				if placeholder(fields[field]) {
					out = append(out, Diagnostic{Error, "code-evidence", path, "Validated task has no concrete " + field + ".", "Record passing validation evidence."})
				}
			}
			if pr := fields["pr"]; pr != "" && !regexp.MustCompile(`(?i)(https?://\S+/(pull|merge_requests)/\d+|#\d+)`).MatchString(pr) {
				out = append(out, Diagnostic{Error, "code-evidence", path, "PR is not a traceable PR reference.", "Record a PR URL or number."})
			}
		}
	}
	return out
}
func placeholder(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	return v == "" || v == "n/a" || v == "none" || v == "pending" || strings.Contains(v, "placeholder") || strings.Contains(v, "not-a-") || strings.Contains(v, "until validation")
}

func validateQualityGates(s Snapshot) []Diagnostic {
	registry, ok := s.JSON[".product/artifacts.json"].(map[string]any)
	if !ok {
		return nil
	}
	items, _ := registry["artifacts"].([]any)
	var out []Diagnostic
	artifacts := map[string]map[string]any{}
	for _, raw := range items {
		if a, ok := raw.(map[string]any); ok {
			if id, _ := a["id"].(string); id != "" {
				artifacts[id] = a
			}
		}
	}
	for _, artifact := range artifacts {
		kind, _ := artifact["type"].(string)
		status, _ := artifact["status"].(string)
		path, _ := artifact["path"].(string)
		text := s.Text[filepath.ToSlash(path)]
		fields := tableFields(text)
		if kind == "qa-evidence" && status == "approved" {
			for _, field := range []string{"test command", "gate logs", "environment", "limitations"} {
				if placeholder(fields[field]) {
					out = append(out, Diagnostic{Error, "qa-evidence", path, "Approved QA evidence has no real gate output or limitation is recorded for " + field + ".", "Record concrete QA gate evidence."})
				}
			}
			if strings.ToLower(fields["verdict"]) != "passed" {
				out = append(out, Diagnostic{Error, "qa-evidence", path, "Approved QA evidence verdict must be passed.", "Resolve QA blockers before approval."})
			}
		}
		if kind == "code-review" && status == "approved" {
			for _, field := range []string{"verdict", "completeness passed", "adherence passed", "quality passed"} {
				if placeholder(fields[field]) {
					out = append(out, Diagnostic{Error, "code-review", path, "Approved Code Review is missing " + field + ".", "Record the review verdict and dimensions."})
				}
			}
		}
		if requiresApproval(status) {
			out = append(out, validateFindingRoutes(path, text)...)
		}
		if kind == "use-case" && (status == "validated" || status == "released") {
			required := []string{"tests", "qa-evidence", "security-review", "audit", "code-review"}
			children := map[string]bool{}
			for _, candidate := range artifacts {
				parents, _ := candidate["parentIds"].([]any)
				for _, p := range parents {
					if p == artifact["id"] && candidate["status"] == "approved" {
						if typ, _ := candidate["type"].(string); typ != "" {
							children[typ] = true
						}
					}
				}
			}
			for _, typ := range required {
				if !children[typ] {
					name := typ + ".md"
					if typ == "tests" {
						name = "tests.md"
					}
					out = append(out, Diagnostic{Error, "validation-gates", path, "Validated use case is missing approved " + name + ".", "Create and approve the required validation artifact."})
				}
			}
		}
	}
	return out
}

func validateFindingRoutes(path, text string) []Diagnostic {
	lines := strings.Split(text, "\n")
	var headers []string
	inFindings := false
	var out []Diagnostic
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "## ") {
			inFindings = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(trim, "## ")), "Findings")
			headers = nil
			continue
		}
		if !inFindings || !strings.HasPrefix(trim, "|") {
			continue
		}
		cells := splitTable(trim)
		if len(cells) == 0 {
			continue
		}
		if headers == nil {
			headers = make([]string, len(cells))
			for i, c := range cells {
				headers[i] = strings.ToLower(c)
			}
			continue
		}
		separator := true
		for _, c := range cells {
			if strings.Trim(c, " -:") != "" {
				separator = false
			}
		}
		if separator {
			continue
		}
		row := map[string]string{}
		for i, c := range cells {
			if i < len(headers) {
				row[headers[i]] = c
			}
		}
		sev := strings.ToLower(row["severity"])
		if sev == "blocker" || sev == "critical" || sev == "required_fix" {
			finding := row["finding"]
			if placeholder(row["route"]) {
				out = append(out, Diagnostic{Error, "failure-routing", path, "Blocking finding is missing Route: " + finding, "Add the owning artifact or workflow route."})
			}
			if placeholder(row["owner"]) {
				out = append(out, Diagnostic{Error, "failure-routing", path, "Blocking finding is missing Owner: " + finding, "Assign the responsible skill or human."})
			}
		}
	}
	return out
}
func splitTable(line string) []string {
	line = strings.Trim(line, "|")
	parts := strings.Split(line, "|")
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = strings.TrimSpace(p)
	}
	return out
}

func artifactList(s Snapshot) []map[string]any {
	registry, _ := s.JSON[".product/artifacts.json"].(map[string]any)
	raw, _ := registry["artifacts"].([]any)
	byPath := map[string]map[string]any{}
	for _, item := range raw {
		if a, ok := item.(map[string]any); ok {
			if path, _ := a["path"].(string); path != "" {
				byPath[filepath.ToSlash(path)] = a
			}
		}
	}
	for _, a := range buildRegistry(s) {
		path, _ := a["path"].(string)
		id, _ := a["id"].(string)
		if path == "" {
			continue
		}
		if _, exists := byPath[filepath.ToSlash(path)]; exists {
			continue
		}
		for existingPath, existing := range byPath {
			existingID, _ := existing["id"].(string)
			if existingID == id && filepath.ToSlash(filepath.Dir(existingPath)) == filepath.ToSlash(filepath.Dir(path)) {
				delete(byPath, existingPath)
			}
		}
		byPath[filepath.ToSlash(path)] = a
	}
	var out []map[string]any
	for _, a := range byPath {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return firstString(out[i]["path"], "") < firstString(out[j]["path"], "") })
	return out
}
func stringSlice(value any) []string {
	raw, _ := value.([]any)
	var out []string
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
func feeds(status string) bool {
	return status == "approved" || status == "in_progress" || status == "implemented" || status == "validated" || status == "released" || status == "not_applicable"
}
func needsParent(status string) bool { return status == "proposed" || feeds(status) }

func validateStatusAndStaleness(s Snapshot) []Diagnostic {
	items := artifactList(s)
	byID := map[string]map[string]any{}
	for _, a := range items {
		if id, _ := a["id"].(string); id != "" {
			byID[id] = a
		}
	}
	var out []Diagnostic
	for _, a := range items {
		id, _ := a["id"].(string)
		status, _ := a["status"].(string)
		path, _ := a["path"].(string)
		for _, parentID := range stringSlice(a["parentIds"]) {
			parent := byID[parentID]
			if parent == nil {
				continue
			}
			parentStatus, _ := parent["status"].(string)
			if feeds(status) && !feeds(parentStatus) {
				out = append(out, Diagnostic{Error, "status-policy", path, fmt.Sprintf("%s is %s, but parent %s is %s.", id, status, parentID, parentStatus), "Approve the parent before advancing the child."})
			}
		}
	}
	derivations, _ := s.JSON[".product/derivations.json"].(map[string]any)
	entries, _ := derivations["derivations"].([]any)
	for index, raw := range entries {
		entry, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id, _ := entry["artifact_id"].(string)
		path, _ := entry["path"].(string)
		artifact := byID[id]
		if artifact == nil {
			out = append(out, Diagnostic{Warning, "staleness", ".product/derivations.json", fmt.Sprintf("Derivation entry %d references unknown artifact %s.", index+1, id), "Regenerate derivations."})
			continue
		}
		if _, ok := s.Text[filepath.ToSlash(path)]; !ok {
			if _, jsonOK := s.JSON[filepath.ToSlash(path)]; !jsonOK {
				out = append(out, Diagnostic{Error, "staleness", path, "Derived artifact path does not exist.", "Fix the path or remove the derivation entry."})
				continue
			}
		}
		sources, _ := entry["derived_from"].([]any)
		for _, sourceRaw := range sources {
			source, ok := sourceRaw.(map[string]any)
			if !ok {
				continue
			}
			sourceID, _ := source["artifact_id"].(string)
			sourcePath, _ := source["path"].(string)
			expected, _ := source["content_hash"].(string)
			text, exists := s.Text[filepath.ToSlash(sourcePath)]
			if !exists {
				if rawJSON, jsonOK := s.Text[filepath.ToSlash(sourcePath)]; jsonOK {
					text = rawJSON
					exists = true
				}
			}
			if !exists {
				out = append(out, Diagnostic{Error, "staleness", path, fmt.Sprintf("%s source path does not exist: %s", id, sourcePath), "Fix the source path."})
				continue
			}
			if Hash(text) != expected {
				severity := Warning
				if status, _ := artifact["status"].(string); needsParent(status) {
					severity = Error
				}
				out = append(out, Diagnostic{severity, "staleness", path, fmt.Sprintf("%s is stale because source %s changed since derivation.", id, sourceID), "Regenerate or re-approve the derived artifact."})
			}
		}
	}
	return out
}

func validateDecisions(s Snapshot) []Diagnostic {
	index, ok := s.JSON[".product/decisions.json"].(map[string]any)
	if !ok {
		return nil
	}
	raw, _ := index["decisions"].([]any)
	known := map[string]map[string]any{}
	var out []Diagnostic
	for _, item := range raw {
		d, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id, _ := d["id"].(string)
		if id == "" {
			out = append(out, Diagnostic{Error, "decisions-index", ".product/decisions.json", "Decision is missing id.", "Add the canonical decision id."})
			continue
		}
		if known[id] != nil {
			out = append(out, Diagnostic{Error, "decisions-index", ".product/decisions.json", "Duplicate decision id: " + id, "Keep one index entry per decision."})
		}
		known[id] = d
		path, _ := d["path"].(string)
		if path == "" {
			out = append(out, Diagnostic{Error, "decisions-index", ".product/decisions.json", id + " is missing path.", "Point to knowledge/decisions/."})
		} else if _, exists := s.Text[filepath.ToSlash(path)]; !exists {
			out = append(out, Diagnostic{Error, "decisions-index", path, id + " decision path does not exist.", "Create the decision or fix its path."})
		}
	}
	decisionRef := regexp.MustCompile(`\bDEC-\d+\b`)
	for path, text := range s.Text {
		if strings.HasPrefix(path, "knowledge/decisions/") {
			continue
		}
		for _, id := range decisionRef.FindAllString(text, -1) {
			if known[id] == nil {
				out = append(out, Diagnostic{Error, "decision-references", path, "Unknown decision reference: " + id, "Add it to .product/decisions.json or fix the reference."})
			}
		}
	}
	return dedupeDiagnostics(out)
}
func dedupeDiagnostics(items []Diagnostic) []Diagnostic {
	seen := map[string]bool{}
	var out []Diagnostic
	for _, d := range items {
		key := string(d.Severity) + d.Check + d.File + d.Message
		if !seen[key] {
			seen[key] = true
			out = append(out, d)
		}
	}
	return out
}

func validateSkillReferences(s Snapshot) []Diagnostic {
	known := map[string]bool{}
	for _, root := range []string{filepath.Join(s.Root, ".codex", "skills"), filepath.Join(s.Root, ".agents", "skills"), filepath.Join(s.FrameworkRoot, "skills"), filepath.Join(s.FrameworkRoot, "framework", "skills")} {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if _, err := os.Stat(filepath.Join(root, entry.Name(), "SKILL.md")); err == nil {
					known[entry.Name()] = true
				}
			}
		}
	}
	if len(known) == 0 {
		return []Diagnostic{{Error, "skill-reference", filepath.ToSlash(filepath.Join(s.FrameworkRoot, "skills")), "No framework skills were found.", "Install framework skills before validation."}}
	}
	var out []Diagnostic
	for _, artifact := range artifactList(s) {
		path, _ := artifact["path"].(string)
		id, _ := artifact["id"].(string)
		if strings.Contains(id, "TEMPLATE") || strings.Contains(path, "_template") {
			continue
		}
		text := s.Text[filepath.ToSlash(path)]
		fields := tableFields(text)
		values := map[string]string{"Owner skill": firstString(artifact["ownerSkill"], fields["owner skill"]), "Next skill": fields["next skill"]}
		for field, value := range values {
			value = strings.Trim(strings.TrimSpace(value), "`")
			if value == "" || placeholder(value) || strings.Contains(strings.ToLower(value), " or ") {
				continue
			}
			name := normalizeSkill(value)
			if name != "" && !known[name] {
				out = append(out, Diagnostic{Error, "skill-reference", path, fmt.Sprintf("%s references missing skill %s.", field, value), "Use a skill installed in the framework skills directory."})
			}
		}
	}
	return out
}
func normalizeSkill(value string) string {
	value = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(value), " AI"))
	return strings.NewReplacer("/", "-", "_", "-", " ", "-").Replace(strings.ToLower(value))
}

func validateDeliveryAndRigor(s Snapshot) []Diagnostic {
	required := map[string]bool{"domain": true, "goal": true, "feature": true, "use-case": true, "specification": true, "implementation-plan": true, "execution-graph": true, "taskset": true, "task": true}
	levels := map[string]bool{"L0": true, "L1": true, "L2": true, "L3": true, "L4": true, "L5": true, "N/A": true}
	priorities := map[string]bool{"P0": true, "P1": true, "P2": true, "P3": true, "N/A": true}
	var out []Diagnostic
	for _, a := range artifactList(s) {
		kind, _ := a["type"].(string)
		if !required[kind] {
			continue
		}
		id, _ := a["id"].(string)
		path, _ := a["path"].(string)
		if strings.Contains(id, "TEMPLATE") || strings.Contains(path, "_template") {
			continue
		}
		delivery, _ := a["delivery"].(map[string]any)
		level := normalizeDelivery(firstString(delivery["level"], ""), "L")
		priority := normalizeDelivery(firstString(delivery["priority"], ""), "P")
		if !levels[level] {
			out = append(out, Diagnostic{Warning, "delivery", path, id + " is missing or has invalid delivery.level.", "Use L0-L5 or N/A."})
		}
		if !priorities[priority] {
			out = append(out, Diagnostic{Warning, "delivery", path, id + " is missing or has invalid delivery.priority.", "Use P0-P3 or N/A."})
		}
		if strings.TrimSpace(firstString(delivery["rationale"], "")) == "" {
			out = append(out, Diagnostic{Warning, "delivery", path, id + " is missing delivery.rationale.", "Explain the level and priority."})
		}
		deps := delivery["depends_on"]
		if deps == nil {
			deps = delivery["dependsOn"]
		}
		if deps != nil {
			if _, ok := deps.([]any); !ok {
				if _, ok := deps.([]string); !ok {
					out = append(out, Diagnostic{Warning, "delivery", path, id + " delivery dependencies must be a list.", "Use delivery.depends_on as an array."})
				}
			}
		}
	}
	triggers := regexp.MustCompile(`(?i)\bauth(?:entication)?\b|\blogin\b|\bpermission(?:s)?\b|\bauthori[sz]ation\b|\bpayment(?:s)?\b|\bPII\b|\bprivacy\b|\bupload\b|\bUGC\b|\bpublic endpoint\b|\bRLS\b|\bmigration\b`)
	for path, text := range s.Text {
		if !strings.HasSuffix(path, "/context.md") || !strings.Contains(path, "/use-cases/") {
			continue
		}
		tier := strings.ToUpper(metadata(text)["rigor_tier"])
		dir := filepath.ToSlash(filepath.Dir(path))
		var combined strings.Builder
		for candidate, body := range s.Text {
			if strings.HasPrefix(candidate, dir+"/") {
				combined.WriteString(body)
				combined.WriteByte('\n')
			}
		}
		if tier != "L" && tier != "N/A" && triggers.MatchString(combined.String()) {
			out = append(out, Diagnostic{Error, "rigor-tier", path, "Sensitive use case requires rigor tier L.", "Set rigor_tier to L and add required validation artifacts."})
		}
	}
	return out
}
func normalizeDelivery(value, prefix string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "N/A" {
		return value
	}
	if strings.HasPrefix(value, prefix) && len(value) >= 2 {
		return value[:2]
	}
	return ""
}

func validateRegistryAndApprovalGates(s Snapshot) []Diagnostic {
	items := artifactList(s)
	if len(items) == 0 {
		return []Diagnostic{{Warning, "artifacts-registry", ".product/artifacts.json", "Artifacts registry is missing or empty.", "Run validate --write-registry."}}
	}
	var out []Diagnostic
	byParent := map[string]map[string]map[string]any{}
	for _, a := range items {
		id, _ := a["id"].(string)
		kind := strings.ReplaceAll(firstString(a["type"], ""), "_", "-")
		status, _ := a["status"].(string)
		path, _ := a["path"].(string)
		if id == "" || kind == "" || status == "" || path == "" {
			out = append(out, Diagnostic{Error, "artifacts-registry", ".product/artifacts.json", "Artifact entry is missing id, type, status, or path.", "Regenerate the artifacts registry."})
			continue
		}
		if _, ok := s.Text[filepath.ToSlash(path)]; !ok {
			if _, jsonOK := s.JSON[filepath.ToSlash(path)]; !jsonOK {
				out = append(out, Diagnostic{Error, "artifacts-registry", path, id + " path does not exist.", "Fix the registry path."})
			}
		}
		for _, parent := range stringSlice(a["parentIds"]) {
			if byParent[parent] == nil {
				byParent[parent] = map[string]map[string]any{}
			}
			byParent[parent][kind] = a
		}
	}
	for _, uc := range items {
		if strings.ReplaceAll(firstString(uc["type"], ""), "_", "-") != "use-case" {
			continue
		}
		id, _ := uc["id"].(string)
		children := byParent[id]
		sequence := []struct{ child, parent, rule string }{{"design", "specification", "design requires an approved Specification"}, {"implementation-plan", "design", "implementation plan requires approved Design"}, {"execution-graph", "implementation-plan", "execution graph requires approved Implementation Plan"}, {"taskset", "execution-graph", "tasks require approved Execution Graph"}}
		for _, gate := range sequence {
			child := children[gate.child]
			if child == nil {
				continue
			}
			childStatus, _ := child["status"].(string)
			if !needsParent(childStatus) {
				continue
			}
			parent := children[gate.parent]
			parentStatus := "missing"
			if parent != nil {
				parentStatus, _ = parent["status"].(string)
			}
			if parent == nil || !feeds(parentStatus) {
				path, _ := child["path"].(string)
				out = append(out, Diagnostic{Error, "approval-gates", path, gate.rule + ".", "Keep the child draft until its parent is approved."})
			}
		}
	}
	return out
}
func title(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}
