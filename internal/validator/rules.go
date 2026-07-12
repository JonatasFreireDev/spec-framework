package validator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/engineeringsystem"
)

func validateImportRuns(s Snapshot) []Diagnostic {
	var out []Diagnostic
	for rel := range s.Text {
		if !strings.HasPrefix(rel, "knowledge/imports/runs/IMPORT-") || !strings.HasSuffix(rel, "/inventory.json") {
			continue
		}
		value, parsed := s.JSON[rel]
		if !parsed {
			out = append(out, Diagnostic{Error, "imports", rel, "Inventory is not valid JSON.", "Regenerate the import inventory."})
			continue
		}
		obj, ok := value.(map[string]any)
		if !ok {
			out = append(out, Diagnostic{Error, "imports", rel, "Inventory is not a JSON object.", "Regenerate the import inventory."})
			continue
		}
		runDir := filepath.ToSlash(filepath.Dir(rel))
		for _, required := range []string{"inventory.json", "import-plan.json", "mapping.json", "conflicts.md", "import-report.md"} {
			path := runDir + "/" + required
			if _, exists := s.Text[path]; !exists {
				out = append(out, Diagnostic{Error, "imports", path, "Import run file is missing.", "Regenerate or repair the import run before materialization."})
			}
		}
		sources, _ := obj["sources"].([]any)
		knownSources := map[string]bool{}
		for _, raw := range sources {
			source, _ := raw.(map[string]any)
			path, _ := source["path"].(string)
			expected, _ := source["sha256"].(string)
			if path == "" || expected == "" {
				out = append(out, Diagnostic{Error, "imports", rel, "Source entry lacks path or sha256.", "Regenerate the inventory."})
				continue
			}
			knownSources[path] = true
			data, err := os.ReadFile(filepath.Join(s.Root, filepath.FromSlash(path)))
			if err != nil {
				out = append(out, Diagnostic{Error, "imports", path, "Imported source is missing.", "Restore the source or create a new import run."})
				continue
			}
			sum := sha256.Sum256(data)
			if hex.EncodeToString(sum[:]) != expected {
				out = append(out, Diagnostic{Warning, "imports", path, "Source changed after this import inventory was created.", "Create a new import run and review affected mappings."})
			}
		}
		plan, _ := s.JSON[runDir+"/import-plan.json"].(map[string]any)
		if plan == nil {
			out = append(out, Diagnostic{Error, "imports", runDir + "/import-plan.json", "Import plan is not valid JSON.", "Repair the plan before materialization."})
		}
		if approved, _ := plan["materialization_approved"].(bool); approved {
			approvedBy, _ := plan["materialization_approved_by"].(string)
			approvedAt, _ := plan["materialization_approved_at"].(string)
			if strings.TrimSpace(approvedBy) == "" || strings.TrimSpace(approvedAt) == "" {
				out = append(out, Diagnostic{Error, "imports", runDir + "/import-plan.json", "Materialization approval lacks approver or timestamp.", "Record explicit human approval evidence."})
			}
		}
		mapping, mappingOK := s.JSON[runDir+"/mapping.json"].(map[string]any)
		if !mappingOK {
			out = append(out, Diagnostic{Error, "imports", runDir + "/mapping.json", "Import mapping is not valid JSON.", "Repair the mapping before materialization."})
			continue
		}
		if importID, _ := mapping["import_id"].(string); importID != filepath.Base(filepath.FromSlash(runDir)) {
			out = append(out, Diagnostic{Error, "imports", runDir + "/mapping.json", "Mapping import_id does not match its run directory.", "Use the enclosing IMPORT-NNN id."})
		}
		mappings, _ := mapping["mappings"].([]any)
		targets := map[string]bool{}
		for _, raw := range mappings {
			m, _ := raw.(map[string]any)
			selected, _ := m["selected"].(bool)
			if !selected {
				continue
			}
			target, _ := m["target"].(string)
			id := fmt.Sprint(m["id"])
			if target == "" {
				out = append(out, Diagnostic{Error, "imports", runDir + "/mapping.json", "Selected mapping " + id + " has no target.", "Add a product-relative target."})
				continue
			}
			clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(target)))
			if clean == ".." || strings.HasPrefix(clean, "../") || filepath.IsAbs(filepath.FromSlash(target)) {
				out = append(out, Diagnostic{Error, "imports", runDir + "/mapping.json", "Selected mapping " + id + " escapes the product root.", "Use a safe product-relative target."})
				continue
			}
			key := strings.ToLower(clean)
			if targets[key] {
				out = append(out, Diagnostic{Error, "imports", runDir + "/mapping.json", "Multiple selected mappings target " + clean + ".", "Reconcile duplicate targets before materialization."})
			}
			targets[key] = true
			refs, _ := m["source_documents"].([]any)
			if len(refs) == 0 {
				out = append(out, Diagnostic{Error, "imports", runDir + "/mapping.json", "Selected mapping " + id + " has no source_documents.", "Add source-level traceability."})
			}
			for _, ref := range refs {
				source := fmt.Sprint(ref)
				if !knownSources[source] {
					out = append(out, Diagnostic{Error, "imports", runDir + "/mapping.json", "Selected mapping " + id + " references an uninventoried source: " + source + ".", "Use a source path from inventory.json."})
				}
			}
			if approved, _ := plan["materialization_approved"].(bool); approved {
				if _, exists := s.Text[clean]; !exists {
					out = append(out, Diagnostic{Error, "imports", clean, "Approved selected mapping was not materialized.", "Materialize the approved run or correct the mapping."})
				}
			}
		}
	}
	return out
}

func validateDeliveryClosure(s Snapshot) []Diagnostic {
	var out []Diagnostic
	for rel, text := range s.Text {
		if !strings.HasPrefix(rel, "domains/") || !strings.HasSuffix(rel, ".md") {
			continue
		}
		fields, meta := tableFields(text), metadata(text)
		status := strings.ToLower(first(fields["status"], meta["status"]))
		if status == "not_applicable" && !validatorMeaningful(first(fields["rationale"], meta["rationale"])) {
			out = append(out, Diagnostic{Error, "not-applicable", rel, "Structured not_applicable status requires a non-placeholder rationale.", "Add a Rationale field explaining why the artifact does not apply."})
		}
	}
	known := map[string]bool{"END": true}
	frameworkSkillTexts := map[string]string{}
	for _, root := range []string{filepath.Join(s.FrameworkRoot, "skills"), filepath.Join(s.FrameworkRoot, "framework", "skills")} {
		entries, _ := os.ReadDir(root)
		for _, e := range entries {
			if e.IsDir() {
				known[e.Name()] = true
				path := filepath.Join(root, e.Name(), "SKILL.md")
				if data, err := os.ReadFile(path); err == nil {
					rel, _ := filepath.Rel(s.FrameworkRoot, path)
					frameworkSkillTexts[filepath.ToSlash(rel)] = string(data)
				}
			}
		}
	}
	nextPattern := regexp.MustCompile(`(?mi)^Next:\s*([^\r\n]+)`)
	canonicalPattern := regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	allSkillTexts := map[string]string{}
	for rel, text := range frameworkSkillTexts {
		allSkillTexts[rel] = text
	}
	for rel, text := range s.Text {
		if strings.HasSuffix(rel, "SKILL.md") {
			allSkillTexts[rel] = text
		}
	}
	for rel, text := range allSkillTexts {
		if !strings.HasSuffix(rel, "SKILL.md") {
			continue
		}
		for _, match := range nextPattern.FindAllStringSubmatch(text, -1) {
			raw := strings.TrimSpace(strings.TrimSuffix(match[1], "."))
			first := strings.FieldsFunc(raw, func(r rune) bool { return r == ' ' || r == ',' })[0]
			if regexp.MustCompile(`^\d+-`).MatchString(first) || strings.HasSuffix(first, ".md") {
				out = append(out, Diagnostic{Error, "skill-handoff", rel, "Legacy numbered handoff: " + first + ".", "Use the canonical skill folder name."})
				continue
			}
			if canonicalPattern.MatchString(first) && !known[first] && first != "human" {
				out = append(out, Diagnostic{Error, "skill-handoff", rel, "Unknown next skill: " + first + ".", "Reference an installed canonical skill name."})
			}
		}
	}
	for rel, text := range s.Text {
		if !strings.HasSuffix(rel, "/context.md") || !strings.Contains(rel, "/use-cases/") {
			continue
		}
		tier := strings.ToUpper(metadata(text)["rigor_tier"])
		if tier == "" {
			tier = strings.ToUpper(strings.Trim(tableFields(text)["tier"], "` []"))
		}
		if tier == "N/A" {
			continue
		}
		base := filepath.ToSlash(filepath.Dir(rel))
		triggers, invalidTriggers := engineeringsystem.Triggers(text)
		for _, trigger := range invalidTriggers {
			if trigger == "invalid_yaml" {
				out = append(out, Diagnostic{Error, "engineering-trigger", rel, "Use-case context has invalid YAML metadata.", "Repair the context YAML before evaluating engineering triggers."})
				continue
			}
			out = append(out, Diagnostic{Error, "engineering-trigger", rel, "Unknown engineering trigger " + trigger + ".", "Use a trigger listed by spec-framework engineering-system triggers."})
		}
		engineeringApplies := tier == "L" || len(triggers) > 0
		traceSeverity := Warning
		if graph, ok := s.JSON[base+"/execution-graph.json"].(map[string]any); ok {
			status := strings.ToLower(firstString(graph["status"], "draft"))
			if status != "" && status != "draft" {
				traceSeverity = Error
			}
		}
		required := []string{"contracts/behavior.md", "contracts/quality.md"}
		if tier == "M" || tier == "L" {
			required = append(required, "contracts/product.md", "contracts/ux.md", "contracts/api.md", "contracts/data.md", "contracts/rollout.md", "technical-discovery.md")
		}
		if tier == "L" {
			required = append(required, "contracts/security.md", "contracts/observability.md")
		}
		if engineeringApplies {
			required = append(required, "engineering-proposal.md", "engineering-review.md")
		}
		for _, name := range required {
			path := base + "/" + name
			if _, ok := s.Text[path]; !ok {
				out = append(out, Diagnostic{Warning, "delivery-closure", path, "Rigor " + tier + " contract is missing.", "Create it or mark the contract Not applicable with rationale during migration."})
			}
		}
		plan := s.Text[base+"/implementation-plan.md"]
		planStatus := strings.ToLower(strings.Trim(tableFields(plan)["status"], "` []"))
		if needsParent(planStatus) {
			discovery := s.Text[base+"/technical-discovery.md"]
			discoveryStatus := strings.ToLower(strings.Trim(tableFields(discovery)["status"], "` []"))
			if !feeds(discoveryStatus) {
				out = append(out, Diagnostic{Error, "approval-gates", base + "/implementation-plan.md", "Implementation Plan requires approved Technical Discovery.", "Approve Technical Discovery before advancing the plan."})
			}
			verdict, _, rationale := validatorArchitectureGate(discovery)
			resolved := (verdict == "Not required" && validatorMeaningful(rationale)) || (verdict == "Decision required" && regexp.MustCompile(`DEC-\d+`).MatchString(discovery))
			if !resolved {
				out = append(out, Diagnostic{Error, "architecture-gate", base + "/technical-discovery.md", "Architecture Gate is unresolved.", "Reference an approved DEC-* or record Not required with rationale."})
			}
			if engineeringApplies {
				review := s.Text[base+"/engineering-review.md"]
				fields := tableFields(review)
				reviewStatus := strings.ToLower(strings.Trim(fields["status"], "` []"))
				verdict := strings.ToLower(strings.Trim(fields["verdict"], "` []"))
				if !feeds(reviewStatus) || verdict != "passed" {
					out = append(out, Diagnostic{Error, "engineering-review-gate", base + "/implementation-plan.md", "Implementation Plan requires a passed approved Engineering Review for this delivery.", "Keep the plan draft until Engineering Review passes against the current proposal."})
				}
				proposal := s.Text[base+"/engineering-proposal.md"]
				if fields["proposal hash"] != Hash(proposal) {
					out = append(out, Diagnostic{Error, "engineering-review-staleness", base + "/engineering-review.md", "Engineering Review does not match the current proposal hash.", "Re-run Engineering Review before advancing the plan."})
				}
			}
		}
		reqPattern := regexp.MustCompile(`\bREQ-\d+\b`)
		acPattern := regexp.MustCompile(`\bAC-\d+\b`)
		testPattern := regexp.MustCompile(`\bTEST-\d+\b`)
		contracts := ""
		for path, body := range s.Text {
			if strings.HasPrefix(path, base+"/contracts/") {
				contracts += "\n" + body
			}
		}
		tasks := ""
		for path, body := range s.Text {
			if strings.HasPrefix(path, base+"/tasks/") {
				tasks += "\n" + body
			}
		}
		quality := s.Text[base+"/contracts/quality.md"]
		for _, id := range uniqueStrings(reqPattern.FindAllString(contracts, -1)) {
			if !strings.Contains(tasks, id) {
				out = append(out, Diagnostic{traceSeverity, "traceability", base, "Requirement " + id + " has no task mapping.", "Reference it from at least one task."})
			}
		}
		for _, id := range uniqueStrings(acPattern.FindAllString(contracts, -1)) {
			if !strings.Contains(tasks, id) {
				out = append(out, Diagnostic{traceSeverity, "traceability", base, "Acceptance criterion " + id + " has no task mapping.", "Reference it from a task."})
			}
			if !strings.Contains(quality, id) {
				out = append(out, Diagnostic{traceSeverity, "traceability", base + "/contracts/quality.md", "Acceptance criterion " + id + " has no quality mapping.", "Map it to a TEST-* or evidence method."})
			}
		}
		if len(acPattern.FindAllString(contracts, -1)) > 0 && len(testPattern.FindAllString(quality, -1)) == 0 {
			out = append(out, Diagnostic{traceSeverity, "traceability", base + "/contracts/quality.md", "Acceptance criteria have no TEST-* identifiers.", "Add stable test ids or explicit non-test evidence."})
		}
	}
	gatesText := s.Text["knowledge/conventions/gates.md"]
	tbdGates := strings.Contains(strings.ToUpper(gatesText), "TBD")
	for rel, text := range s.Text {
		if !strings.Contains(rel, "/tasks/") || !strings.HasSuffix(rel, ".md") {
			continue
		}
		fields := tableFields(text)
		status := strings.Trim(fields["status"], "` []")
		if status == "" {
			status = metadata(text)["status"]
		}
		if status == "implemented" || status == "validated" || status == "released" {
			for _, field := range []string{"branch", "base commit", "diff hash", "changed paths", "test status"} {
				v := strings.ToLower(fields[field])
				if v == "" || strings.Contains(v, "n/a until") || strings.Contains(v, "pending") {
					out = append(out, Diagnostic{Error, "working-tree-evidence", rel, "Task " + status + " lacks " + field + ".", "Record immutable working-tree evidence before implementation status."})
				}
			}
		}
		if tbdGates && (status == "in_progress" || status == "implemented" || status == "validated" || status == "released") {
			out = append(out, Diagnostic{Error, "gate-readiness", rel, "Task advanced while applicable gate commands remain TBD.", "Configure gates or mark them N/A with rationale before Code Runner."})
		}
		if status == "validated" || status == "released" {
			for _, field := range []string{"commits", "code paths", "code review diff hash", "qa diff hash"} {
				v := strings.ToLower(fields[field])
				if v == "" || strings.Contains(v, "pending") || strings.Contains(v, "n/a") {
					out = append(out, Diagnostic{Error, "validation-evidence", rel, "Task " + status + " lacks " + field + ".", "Commit only after Code Review and QA approve the same diff hash."})
				}
			}
			diff := strings.Trim(fields["diff hash"], "` ")
			review := strings.Trim(fields["code review diff hash"], "` ")
			qa := strings.Trim(fields["qa diff hash"], "` ")
			if diff != "" && (review != diff || qa != diff) {
				out = append(out, Diagnostic{Error, "diff-staleness", rel, "Code Review and QA did not approve the current diff hash.", "Re-run both independent gates on the current working-tree snapshot before commit."})
			}
		}
	}
	return out
}

func uniqueStrings(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			out = append(out, item)
		}
	}
	sort.Strings(out)
	return out
}

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
		out[key] = strings.Trim(strings.TrimSpace(m[2]), "`")
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
		triggers, _ := engineeringsystem.Triggers(contextText)
		if tier == "L" || len(triggers) > 0 {
			required = append(required, "engineering-proposal.md", "engineering-review.md")
		}
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
			for _, field := range []string{"branch", "base commit", "diff hash", "changed paths", "test status"} {
				if placeholder(fields[field]) {
					out = append(out, Diagnostic{Error, "code-evidence", path, "Task is " + status + " but has no concrete " + field + ".", "Record immutable working-tree evidence."})
				}
			}
		}
		if kind == "task" && (status == "validated" || status == "released") {
			for _, field := range []string{"commits", "code paths", "test status", "gate logs", "qa evidence", "code review diff hash", "qa diff hash"} {
				if placeholder(fields[field]) {
					out = append(out, Diagnostic{Error, "code-evidence", path, "Validated task has no concrete " + field + ".", "Record passing validation evidence."})
				}
			}
			if commit := fields["commits"]; commit != "" && !regexp.MustCompile(`(?i)([0-9a-f]{7,40}|https?://\S+/commit/\S+)`).MatchString(commit) {
				out = append(out, Diagnostic{Error, "code-evidence", path, "Commits is not a traceable commit hash or commit URL.", "Record a commit hash or URL."})
			}
			if pr := fields["pr"]; !placeholder(pr) && !regexp.MustCompile(`(?i)(https?://\S+/(pull|merge_requests)/\d+|#\d+)`).MatchString(pr) {
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
	records := map[string][]map[string]any{}
	for path, value := range s.JSON {
		if strings.HasPrefix(path, ".product/history/approval-") {
			if rec, ok := value.(map[string]any); ok {
				id, _ := rec["artifact_id"].(string)
				records[id] = append(records[id], rec)
			}
		}
	}
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
		decisionType, _ := d["type"].(string)
		if decisionType == "" {
			if scope, _ := d["scope"].(string); strings.Contains(scope, "architecture") {
				decisionType = "architecture"
			}
		}
		if decisionType != "" && !map[string]bool{"product": true, "architecture": true, "security": true, "data": true, "delivery": true}[decisionType] {
			out = append(out, Diagnostic{Error, "decisions-index", ".product/decisions.json", id + " has unsupported type " + decisionType + ".", "Use a canonical decision type."})
		}
		status, _ := d["status"].(string)
		if status == "approved" && path != "" {
			text := s.Text[filepath.ToSlash(path)]
			matched := false
			for _, rec := range records[id] {
				if rec["path"] == path && rec["status_granted"] == "approved" && rec["content_hash"] == Hash(text) {
					matched = true
				}
			}
			if !matched {
				out = append(out, Diagnostic{Error, "decision-approval", path, id + " is approved but has no current hash-matching approval record.", "Re-approve the current decision through the human approval flow."})
			}
		}
		if effects, exists := d["workflowEffects"]; exists {
			obj, ok := effects.(map[string]any)
			if !ok {
				out = append(out, Diagnostic{Error, "decision-effects", ".product/decisions.json", id + " workflowEffects must be an object.", "Use structured effect arrays."})
			} else {
				for _, field := range []string{"requiredTaskTypes", "requiredGates", "requiredEvidence", "requiredWriteScopes", "sharedResources"} {
					if v, present := obj[field]; present {
						if _, ok := v.([]any); !ok {
							out = append(out, Diagnostic{Error, "decision-effects", ".product/decisions.json", id + " " + field + " must be an array.", "Use an array, including [] when empty."})
						}
					}
				}
			}
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
			if strings.HasSuffix(path, "/technical-discovery.md") {
				d := known[id]
				if d != nil && firstString(d["status"], "") != "approved" {
					out = append(out, Diagnostic{Error, "architecture-gate", path, id + " is not approved.", "Approve the decision before advancing Technical Discovery."})
				}
				if d != nil && !decisionApplies(d, path) {
					out = append(out, Diagnostic{Error, "decision-scope", path, id + " does not apply to this use-case scope.", "Add the affected path to decision scope/affectedArtifacts or reference the correct decision."})
				}
			}
		}
	}
	return dedupeDiagnostics(out)
}
func decisionApplies(d map[string]any, path string) bool {
	var scopes []string
	scopes = append(scopes, stringSlice(d["scope"])...)
	if s, ok := d["scope"].(string); ok {
		scopes = append(scopes, strings.Split(s, "/")...)
	}
	scopes = append(scopes, stringSlice(d["affectedArtifacts"])...)
	if len(scopes) == 0 {
		return false
	}
	for _, s := range scopes {
		s = filepath.ToSlash(strings.TrimSpace(s))
		if s == "product" || s == "architecture" || s == "security" || s == "data" || s == "delivery" || strings.HasPrefix(path, strings.TrimSuffix(s, "/")) || strings.HasPrefix(s, filepath.ToSlash(filepath.Dir(path))) {
			return true
		}
	}
	return false
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
	required := map[string]bool{"domain": true, "goal": true, "feature": true, "use-case": true, "specification": true, "technical-discovery": true, "engineering-proposal": true, "engineering-review": true, "implementation-plan": true, "execution-graph": true, "taskset": true, "task": true}
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

func validatorArchitectureGate(text string) (string, string, string) {
	fields := tableFields(text)
	if fields["verdict"] != "" {
		return fields["verdict"], fields["decision"], fields["rationale"]
	}
	pattern := regexp.MustCompile(`(?mi)^\|\s*(Decision required|Not required)\s*\|\s*([^|]*)\|\s*([^|]*)\|$`)
	match := pattern.FindStringSubmatch(text)
	if len(match) == 4 {
		return strings.TrimSpace(match[1]), strings.TrimSpace(match[2]), strings.TrimSpace(match[3])
	}
	return "", "", ""
}

func validatorMeaningful(value string) bool {
	value = strings.ToLower(strings.Trim(strings.TrimSpace(value), "`[]"))
	return value != "" && value != "n/a" && value != "none" && value != "tbd" && value != "pending" && !strings.Contains(value, "placeholder")
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
		ucPath, _ := uc["path"].(string)
		contextText := s.Text[filepath.ToSlash(ucPath)]
		tierL := strings.ToUpper(metadata(contextText)["rigor_tier"]) == "L"
		triggers, _ := engineeringsystem.Triggers(contextText)
		engineeringApplies := tierL || len(triggers) > 0
		sequence := []struct{ child, parent, rule string }{{"design", "specification", "design requires an approved Specification"}, {"technical-discovery", "design", "technical discovery requires approved Design"}, {"engineering-proposal", "technical-discovery", "engineering proposal requires approved Technical Discovery"}, {"engineering-review", "engineering-proposal", "engineering review requires an approved Engineering Proposal"}, {"implementation-plan", "engineering-review", "implementation plan requires approved Engineering Review when applicable"}, {"execution-graph", "implementation-plan", "execution graph requires approved Implementation Plan"}, {"taskset", "execution-graph", "tasks require approved Execution Graph"}}
		for _, gate := range sequence {
			if !engineeringApplies && (gate.child == "engineering-proposal" || gate.child == "engineering-review" || gate.parent == "engineering-review") {
				continue
			}
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
	ucIDPattern := regexp.MustCompile(`\bUC-[A-Z0-9-]+\b`)
	for _, feature := range items {
		if strings.ReplaceAll(firstString(feature["type"], ""), "_", "-") != "feature" {
			continue
		}
		id := firstString(feature["id"], "")
		path := firstString(feature["path"], "")
		status := strings.ToLower(firstString(feature["status"], "draft"))
		severity := Warning
		if status != "draft" {
			severity = Error
		}
		declared := map[string]bool{}
		for _, x := range ucIDPattern.FindAllString(s.Text[path], -1) {
			declared[x] = true
		}
		actual := map[string]bool{}
		for _, child := range items {
			for _, parent := range stringSlice(child["parentIds"]) {
				if parent == id && strings.ReplaceAll(firstString(child["type"], ""), "_", "-") == "use-case" {
					actual[firstString(child["id"], "")] = true
				}
			}
		}
		for x := range declared {
			if !actual[x] {
				out = append(out, Diagnostic{severity, "feature-coverage", path, "Feature declares use case " + x + " but no canonical Use Case exists.", "Create the Use Case or remove it from the feature scope."})
			}
		}
		for x := range actual {
			if !declared[x] {
				out = append(out, Diagnostic{severity, "feature-coverage", path, "Canonical use case " + x + " is not listed by the Feature.", "Add it to the Feature coverage table."})
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
