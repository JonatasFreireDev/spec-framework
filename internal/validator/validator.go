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

	"github.com/JonatasFreireDev/spec-framework/internal/designsystem"
	"github.com/JonatasFreireDev/spec-framework/internal/engineeringsystem"
)

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Note    Severity = "note"
)

type Diagnostic struct {
	Severity Severity `json:"severity"`
	Check    string   `json:"check"`
	File     string   `json:"file"`
	Message  string   `json:"message"`
	Fix      string   `json:"fix"`
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
	return validate(ctx, root, frameworkRoot, false, "", "", "")
}

// ValidateStrict promotes delivery warnings on approved (or later) artifacts
// to errors at approval and delivery boundaries.
func ValidateStrict(ctx context.Context, root, frameworkRoot string) (Result, error) {
	return validate(ctx, root, frameworkRoot, true, "", "", "")
}

// ValidateCandidate is read-only and validates an artifact as if its proposed
// content and status had already been written.
func ValidateCandidate(ctx context.Context, root, frameworkRoot, artifactPath, content, status string) (Result, error) {
	return validate(ctx, root, frameworkRoot, true, artifactPath, content, status)
}

// AuditTemplate returns only structural/provenance findings for one registered
// artifact, making it suitable for an actionable CLI audit command.
func AuditTemplate(ctx context.Context, root, frameworkRoot, artifactPath string) ([]Diagnostic, error) {
	snap, err := Scan(ctx, root, frameworkRoot, 0)
	if err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(root, artifactPath)
	if err != nil {
		return nil, err
	}
	rel = filepath.ToSlash(rel)
	if registry, ok := snap.JSON[".product/artifacts.json"].(map[string]any); ok {
		if items, ok := registry["artifacts"].([]any); ok {
			for _, raw := range items {
				if item, ok := raw.(map[string]any); ok && filepath.ToSlash(fmt.Sprint(item["path"])) == rel {
					item["status"] = "approved"
				}
			}
		}
	}
	var findings []Diagnostic
	for _, diagnostic := range validateTemplateConformance(snap) {
		if filepath.ToSlash(diagnostic.File) == rel && diagnostic.Check == "template-conformance" {
			findings = append(findings, diagnostic)
		}
	}
	return findings, nil
}

func validate(ctx context.Context, root, frameworkRoot string, strict bool, candidatePath, candidateContent, candidateStatus string) (Result, error) {
	snap, err := Scan(ctx, root, frameworkRoot, 0)
	if err != nil {
		return Result{}, err
	}
	if candidatePath != "" {
		rel, err := filepath.Rel(root, candidatePath)
		if err != nil {
			return Result{}, err
		}
		rel = filepath.ToSlash(rel)
		snap.Text[rel] = candidateContent
		if registry, ok := snap.JSON[".product/artifacts.json"].(map[string]any); ok {
			if items, ok := registry["artifacts"].([]any); ok {
				for _, raw := range items {
					if item, ok := raw.(map[string]any); ok && filepath.ToSlash(fmt.Sprint(item["path"])) == rel {
						item["status"] = candidateStatus
					}
				}
			}
		}
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
	d = append(d, validateDeliveryClosure(snap)...)
	d = append(d, validateSkillDiscoveryContracts(snap)...)
	d = append(d, validateDomainModeling(snap)...)
	d = append(d, validateDesignArtifacts(snap)...)
	d = append(d, validateDesignSystem(snap)...)
	d = append(d, validateEngineeringSystem(snap)...)
	d = append(d, validateTemplateConformance(snap)...)
	if strict {
		d = promoteApprovedWarnings(d, snap)
	}
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

func validateEngineeringSystem(s Snapshot) []Diagnostic {
	var out []Diagnostic
	inspection := engineeringsystem.Inspection{}
	configured := false
	if _, exists := s.Text["engineering/context.md"]; exists {
		var err error
		inspection, err = engineeringsystem.Inspect(s.Root)
		if err != nil {
			out = append(out, Diagnostic{Error, "engineering-system", "engineering/context.md", err.Error(), "Restore or initialize the Engineering System contract."})
		} else {
			for _, blocker := range inspection.Blockers {
				out = append(out, Diagnostic{Error, "engineering-system", "engineering/engineering-system.yaml", blocker, "Repair the Engineering System catalog, contract paths, maturity, or evidence."})
			}
			configured = inspection.Scope != "not-configured"
		}
	}
	for rel, proposal := range s.Text {
		if !strings.HasPrefix(rel, "domains/") || !strings.HasSuffix(rel, "/engineering-proposal.md") {
			continue
		}
		fields := tableFields(proposal)
		pin := strings.TrimSpace(fields["engineering system"])
		if !configured {
			if !strings.EqualFold(pin, "Not configured") {
				out = append(out, Diagnostic{Error, "engineering-system-consumer", rel, "Engineering Proposal must declare Not configured while no Engineering System is configured.", "Use Not configured or establish the shared Engineering System first."})
			}
		} else {
			expected := inspection.ID + " @ " + inspection.Version
			if pin != expected {
				out = append(out, Diagnostic{Error, "engineering-system-consumer", rel, "Engineering Proposal does not pin the declared Engineering System id and version.", "Set Engineering System to " + expected + "."})
			}
			status := markdownStatus(proposal)
			if requiresApproval(status) && (inspection.Status != "approved" || !currentApproval(s, inspection.ID, "engineering/engineering-system.md", inspection.Status)) {
				out = append(out, Diagnostic{Error, "engineering-system-consumer", rel, "Advanced Engineering Proposal consumes an Engineering System without current approval evidence.", "Approve the current Engineering System or keep the proposal draft."})
			}
		}
	}
	if configured && inspection.QualitySystem {
		expected := inspection.ID + " @ " + inspection.Version
		knownExceptions := map[string]bool{}
		for _, id := range inspection.QualityExceptions {
			knownExceptions[id] = true
		}
		knownDimensions := map[string][]string{
			"environments": inspection.QualityEnvironments,
			"test data":    inspection.QualityTestDataClasses,
			"platforms":    inspection.QualityPlatforms,
		}
		for rel, tests := range s.Text {
			if !strings.HasPrefix(rel, "domains/") || !strings.HasSuffix(rel, "/tests.md") {
				continue
			}
			status := markdownStatus(tests)
			if status == "draft" {
				continue
			}
			fields := tableFields(tests)
			pin := strings.TrimSpace(fields["engineering system"])
			if pin != expected {
				out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests do not pin the configured Engineering Quality System.", "Set Engineering System to " + expected + "."})
			}
			if !strings.Contains(strings.ToLower(tests), "deviations or exceptions") {
				out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests do not declare Quality System deviations or exceptions.", "Record None or reference a governed quality exception."})
			}
			for _, field := range []string{"quality policy", "applicable risks"} {
				if !validatorMeaningful(fields[field]) {
					out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests have no meaningful " + field + " value.", "Apply the configured Quality System policy explicitly."})
				}
			}
			if !strings.Contains(filepath.ToSlash(strings.ToLower(fields["quality policy"])), "engineering/quality/quality-system.md") {
				out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests do not reference the canonical Engineering Quality System policy.", "Set Quality policy to engineering/quality/quality-system.md."})
			}
			for field, allowed := range knownDimensions {
				declared := declaredValues(fields[field])
				if len(declared) == 0 {
					out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests have no declared " + field + " values.", "Select values configured by the Engineering Quality System."})
					continue
				}
				allowedSet := map[string]bool{}
				for _, value := range allowed {
					allowedSet[value] = true
				}
				for _, value := range declared {
					if !allowedSet[value] {
						out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests reference unconfigured " + field + " value " + value + ".", "Use a value from quality-system.yaml."})
					}
				}
			}
			deviations := strings.ToLower(strings.Trim(strings.TrimSpace(fields["deviations or exceptions"]), "`[]"))
			exceptions := uniqueStrings(regexp.MustCompile(`\bQEX-[A-Z0-9-]+\b`).FindAllString(fields["deviations or exceptions"], -1))
			if deviations != "none" && len(exceptions) == 0 {
				out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests must declare None or at least one QEX-* exception.", "Use None when no deviation applies."})
			}
			for _, exception := range exceptions {
				if !knownExceptions[exception] {
					out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests reference unavailable quality exception " + exception + ".", "Use an open, unexpired, in-scope exception or remove the reference."})
					continue
				}
				scope := filepath.ToSlash(inspection.QualityExceptionScopes[exception])
				if scope != "product" && !strings.HasPrefix(rel, strings.TrimSuffix(scope, "/")+"/") {
					out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Quality exception " + exception + " does not apply to this use case.", "Use an exception whose scope contains this tests.md path."})
				}
			}
			base := filepath.ToSlash(filepath.Dir(rel))
			contracts := ""
			for path, body := range s.Text {
				if strings.HasPrefix(path, base+"/contracts/") {
					contracts += "\n" + body
				}
			}
			traceability := markdownTableRows(tests, "Acceptance Traceability")
			for _, acceptance := range uniqueStrings(regexp.MustCompile(`\bAC-[A-Z0-9-]+\b`).FindAllString(contracts, -1)) {
				var mapped map[string]string
				for _, row := range traceability {
					if containsExactID(row["acceptance criterion"], acceptance) {
						mapped = row
						break
					}
				}
				if mapped == nil {
					out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Tests do not structurally map acceptance criterion " + acceptance + ".", "Add it to the Acceptance Traceability table."})
					continue
				}
				for _, column := range []string{"risk", "validation method", "test level", "evidence", "owner"} {
					if !validatorMeaningful(mapped[column]) {
						out = append(out, Diagnostic{Error, "quality-system-consumer", rel, "Acceptance criterion " + acceptance + " has no meaningful " + column + ".", "Complete every required Acceptance Traceability column."})
					}
				}
			}
		}
		for rel, qa := range s.Text {
			if !strings.HasPrefix(rel, "domains/") || !strings.HasSuffix(rel, "/qa-evidence.md") || !requiresApproval(markdownStatus(qa)) {
				continue
			}
			fields := tableFields(qa)
			if strings.TrimSpace(fields["engineering system"]) != expected {
				out = append(out, Diagnostic{Error, "quality-system-qa", rel, "Approved QA evidence does not pin the configured Engineering Quality System.", "Set Engineering System to " + expected + " and re-run QA."})
			}
			for _, check := range []string{"Quality System conformity", "Environment and test data policy", "Flaky test and exception policy"} {
				if !qualityCheckPassed(qa, check) {
					out = append(out, Diagnostic{Error, "quality-system-qa", rel, "Approved QA evidence has no passed " + check + " check.", "Run and record the check against the pinned policy."})
				}
			}
		}
	}
	for rel, review := range s.Text {
		if !strings.HasPrefix(rel, "domains/") || !strings.HasSuffix(rel, "/engineering-review.md") {
			continue
		}
		fields := tableFields(review)
		if strings.ToLower(fields["verdict"]) != "passed" {
			continue
		}
		proposalPath := filepath.ToSlash(filepath.Join(filepath.Dir(rel), "engineering-proposal.md"))
		proposal, exists := s.Text[proposalPath]
		if !exists || fields["proposal hash"] != Hash(proposal) {
			out = append(out, Diagnostic{Error, "engineering-review-staleness", rel, "Passed Engineering Review does not match the current Engineering Proposal hash.", "Re-run Engineering Review and record the current SHA-256 proposal hash."})
		}
	}
	return out
}

func qualityCheckPassed(text, label string) bool {
	pattern := regexp.MustCompile(`(?mi)^\|\s*` + regexp.QuoteMeta(label) + `\s*\|[^|]*\|\s*passed\s*\|`)
	return pattern.MatchString(text)
}

func declaredValues(value string) []string {
	value = strings.Trim(strings.TrimSpace(value), "`[]")
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ';' })
	seen := map[string]bool{}
	var out []string
	for _, part := range parts {
		part = strings.ToLower(strings.TrimSpace(part))
		if part != "" && !seen[part] {
			seen[part] = true
			out = append(out, part)
		}
	}
	return out
}

func markdownTableRows(text, heading string) []map[string]string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	active := false
	var headers []string
	var rows []map[string]string
	separatorPattern := regexp.MustCompile(`^:?-{3,}:?$`)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if active {
				break
			}
			active = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")), heading)
			continue
		}
		if !active || !strings.HasPrefix(trimmed, "|") || !strings.HasSuffix(trimmed, "|") {
			continue
		}
		cells := splitMarkdownRow(strings.Trim(trimmed, "|"))
		for index := range cells {
			cells[index] = strings.Trim(strings.TrimSpace(cells[index]), "`[]")
		}
		separator := true
		for _, cell := range cells {
			separator = separator && separatorPattern.MatchString(cell)
		}
		if separator {
			continue
		}
		if headers == nil {
			for _, cell := range cells {
				headers = append(headers, strings.ToLower(cell))
			}
			continue
		}
		if len(cells) != len(headers) {
			continue
		}
		row := map[string]string{}
		for index, header := range headers {
			row[header] = cells[index]
		}
		rows = append(rows, row)
	}
	return rows
}

func splitMarkdownRow(line string) []string {
	var cells []string
	var current strings.Builder
	escaped := false
	for _, r := range line {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '|' {
			cells = append(cells, current.String())
			current.Reset()
			continue
		}
		current.WriteRune(r)
	}
	if escaped {
		current.WriteRune('\\')
	}
	cells = append(cells, current.String())
	return cells
}

func containsExactID(value, expected string) bool {
	for _, id := range regexp.MustCompile(`\bAC-[A-Z0-9-]+\b`).FindAllString(value, -1) {
		if id == expected {
			return true
		}
	}
	return false
}

func currentApproval(s Snapshot, id, path, status string) bool {
	text, exists := s.Text[filepath.ToSlash(path)]
	if !exists || !requiresApproval(status) {
		return false
	}
	expected := artifactHash(s, filepath.ToSlash(path), text)
	for rel, value := range s.JSON {
		if !strings.HasPrefix(rel, ".product/history/approval-") {
			continue
		}
		record, _ := value.(map[string]any)
		if record["artifact_id"] == id && filepath.ToSlash(fmt.Sprint(record["path"])) == filepath.ToSlash(path) && record["status_granted"] == status && record["content_hash"] == expected {
			return true
		}
	}
	return false
}

func validateDesignSystem(s Snapshot) []Diagnostic {
	contextPath := "design/system/context.md"
	if _, exists := s.Text[contextPath]; !exists {
		return nil
	}
	var out []Diagnostic
	out = append(out, validateContext(contextPath, s.Text[contextPath])...)
	inspection, err := designsystem.Inspect(s.Root)
	if err != nil {
		return append(out, Diagnostic{Error, "design-system", contextPath, err.Error(), "Restore or initialize the Design System."})
	}
	for _, blocker := range inspection.Blockers {
		out = append(out, Diagnostic{Error, "design-system", "design/system/tokens/tokens.json", blocker, "Repair the Design System contract or tokens."})
	}
	for rel, text := range s.Text {
		if !strings.HasSuffix(rel, "/design.md") || !strings.HasPrefix(rel, "domains/") || strings.Contains(strings.ToLower(text), "not applicable") {
			continue
		}
		status := markdownStatus(text)
		if !requiresApproval(status) && status != "proposed" {
			continue
		}
		if !strings.Contains(text, "design_system:") || !strings.Contains(text, "id: "+inspection.ID) || !strings.Contains(text, "version: "+inspection.Version) {
			out = append(out, Diagnostic{Error, "design-system-consumer", rel, "Proposed-or-later Design must pin the declared Design System id and version", "Add design_system id/path/version and consumed tokens/components/patterns."})
		}
		if status != "draft" && inspection.Status != "approved" {
			out = append(out, Diagnostic{Error, "design-system-consumer", rel, "Design consumes a Design System that is not approved", "Approve the system with a current approval record or keep Design draft."})
		}
	}
	return out
}

func markdownStatus(text string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?im)^\s*[-|]?\s*status\s*[:|]\s*` + "`?" + `([a-z_ -]+)`),
		regexp.MustCompile(`(?im)^\s*status:\s*([a-z_]+)`),
	}
	for _, pattern := range patterns {
		if match := pattern.FindStringSubmatch(text); len(match) == 2 {
			return strings.TrimSpace(strings.Trim(match[1], "`| "))
		}
	}
	return "draft"
}

func validateDesignArtifacts(s Snapshot) []Diagnostic {
	var out []Diagnostic
	sources := map[string]map[string]any{}
	for rel, value := range s.JSON {
		if !strings.HasPrefix(rel, "design/sources/") || !strings.HasSuffix(rel, "/manifest.json") {
			continue
		}
		m, ok := value.(map[string]any)
		if !ok {
			out = append(out, Diagnostic{Error, "design-source", rel, "Invalid source manifest JSON", "Write a valid manifest object."})
			continue
		}
		id := fmt.Sprint(m["id"])
		if !regexp.MustCompile(`^DSRC-[0-9]{3,}$`).MatchString(id) {
			out = append(out, Diagnostic{Error, "design-source", rel, "Invalid or missing source id", "Use DSRC-NNN."})
		}
		if !validString(m["authority"], "visual_canonical", "reference", "inspiration") {
			out = append(out, Diagnostic{Error, "design-source", rel, "Invalid source authority", "Use visual_canonical, reference, or inspiration."})
		}
		version, _ := m["version"].(map[string]any)
		if strings.TrimSpace(fmt.Sprint(version["value"])) == "" {
			out = append(out, Diagnostic{Error, "design-source", rel, "Source version is missing", "Record an immutable version or SHA-256."})
		}
		screens, ok := m["screens"].([]any)
		if !ok {
			out = append(out, Diagnostic{Error, "design-source", rel, "screens must be an array", "Add a screen inventory."})
		}
		seen := map[string]bool{}
		for _, raw := range screens {
			screen, _ := raw.(map[string]any)
			screenID := fmt.Sprint(screen["id"])
			if screenID == "" || seen[screenID] {
				out = append(out, Diagnostic{Error, "design-source", rel, "Screen IDs must be present and unique", "Assign stable SCREEN-NNN ids."})
			}
			seen[screenID] = true
			if asset := strings.TrimSpace(fmt.Sprint(screen["asset"])); asset != "" {
				assetPath := filepath.ToSlash(filepath.Join(filepath.Dir(rel), asset))
				if _, exists := s.Text[assetPath]; !exists {
					out = append(out, Diagnostic{Error, "design-source", rel, "Missing visual asset: " + asset, "Restore the snapshot or update the manifest."})
				}
			}
		}
		sources[id] = m
	}
	for rel, value := range s.JSON {
		if !strings.HasPrefix(rel, "design/use-cases/") || !strings.HasSuffix(rel, "/manifest.json") {
			continue
		}
		m, ok := value.(map[string]any)
		if !ok {
			continue
		}
		if !validString(m["originMode"], "generate", "evolve", "adopt") {
			out = append(out, Diagnostic{Error, "design-contract", rel, "Invalid originMode", "Use generate, evolve, or adopt."})
		}
		if !validString(m["maturity"], "contract", "wireframe", "mockup", "prototype") {
			out = append(out, Diagnostic{Error, "design-contract", rel, "Invalid maturity", "Use contract, wireframe, mockup, or prototype."})
		}
		if !validString(m["fidelityPolicy"], "strict", "balanced", "exploratory") {
			out = append(out, Diagnostic{Error, "design-contract", rel, "Invalid fidelityPolicy", "Use strict, balanced, or exploratory."})
		}
		for _, id := range stringAnySlice(m["sources"]) {
			if _, exists := sources[id]; !exists {
				out = append(out, Diagnostic{Error, "design-contract", rel, "Missing source manifest: " + id, "Import or restore the source."})
			}
		}
		for _, raw := range anySlice(m["mappings"]) {
			mapping, _ := raw.(map[string]any)
			if strings.TrimSpace(fmt.Sprint(mapping["requirement"])) == "" {
				out = append(out, Diagnostic{Error, "design-coverage", rel, "Mapping is missing requirement", "Reference a stable REQ-* id."})
			}
			if !validString(mapping["coverage"], "covered", "partial", "missing", "conflict", "not-verifiable", "not-applicable") {
				out = append(out, Diagnostic{Error, "design-coverage", rel, "Invalid mapping coverage", "Use a supported coverage state."})
			}
		}
	}
	return out
}

func validString(value any, options ...string) bool {
	s := fmt.Sprint(value)
	for _, option := range options {
		if s == option {
			return true
		}
	}
	return false
}

func anySlice(value any) []any {
	items, _ := value.([]any)
	return items
}

func stringAnySlice(value any) []string {
	var out []string
	for _, item := range anySlice(value) {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
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
	graphStatus := strings.ToLower(fmt.Sprint(object["status"]))
	requireTaskFiles := graphStatus == "materialized" || graphStatus == "approved" || graphStatus == "in_progress" || graphStatus == "implemented" || graphStatus == "validated" || graphStatus == "released"
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
			if _, ok := snap.Text[full]; !ok && requireTaskFiles {
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
			if target == "" || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") || strings.Contains(target, "<") {
				continue
			}
			fragment := ""
			if i := strings.Index(target, "#"); i >= 0 {
				fragment = target[i+1:]
				target = target[:i]
			}
			if target == "" {
				if fragment != "" && !markdownAnchors(text)[fragment] {
					out = append(out, Diagnostic{Error, "links", rel, "Broken Markdown section link: #" + fragment, "Create the target heading or update the link."})
				}
				continue
			}
			decoded, err := url.PathUnescape(strings.Trim(target, "<>"))
			if err == nil {
				target = decoded
			}
			candidate := filepath.Clean(filepath.Join(s.Root, filepath.Dir(filepath.FromSlash(rel)), filepath.FromSlash(target)))
			if _, err := os.Stat(candidate); err != nil {
				out = append(out, Diagnostic{Error, "links", rel, "Broken Markdown link: " + target, "Create the target or update the link."})
			} else if fragment != "" {
				targetText, readErr := os.ReadFile(candidate)
				if readErr != nil || !markdownAnchors(string(targetText))[fragment] {
					out = append(out, Diagnostic{Error, "links", rel, "Broken Markdown section link: " + target + "#" + fragment, "Create the target heading or update the link."})
				}
			}
		}
	}
	return out
}

var markdownHeading = regexp.MustCompile(`(?m)^#{1,6}\s+(.+?)\s*#*\s*$`)

func markdownAnchors(text string) map[string]bool {
	anchors := map[string]bool{}
	for _, match := range markdownHeading.FindAllStringSubmatch(text, -1) {
		value := strings.ToLower(strings.TrimSpace(match[1]))
		value = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(value, "")
		value = regexp.MustCompile(`[^\pL\pN\s-]`).ReplaceAllString(value, "")
		value = regexp.MustCompile(`[\s-]+`).ReplaceAllString(value, "-")
		anchors[strings.Trim(value, "-")] = true
	}
	return anchors
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
		expected := artifactHash(s, filepath.ToSlash(path), text)
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

func artifactHash(s Snapshot, path, text string) string {
	if filepath.Base(path) != "engineering-system.md" {
		return Hash(text)
	}
	var paths []string
	for candidate := range s.Text {
		if strings.HasPrefix(candidate, "engineering/") {
			paths = append(paths, candidate)
		}
	}
	sort.Strings(paths)
	var content strings.Builder
	for _, candidate := range paths {
		content.WriteString(candidate)
		content.WriteByte('\n')
		content.WriteString(normalizedText(s.Text[candidate]))
		content.WriteByte('\n')
	}
	return Hash(content.String())
}

func normalizedText(text string) string {
	text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(text, "\n")
	for index, line := range lines {
		lines[index] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
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
