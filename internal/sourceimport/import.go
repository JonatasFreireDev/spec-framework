package sourceimport

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Source struct {
	Path   string `json:"path"`
	Format string `json:"format"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

type Inventory struct {
	SchemaVersion int      `json:"schema_version"`
	ImportID      string   `json:"import_id"`
	Sources       []Source `json:"sources"`
}

type Mapping struct {
	ID              string   `json:"id"`
	Target          string   `json:"target"`
	ArtifactType    string   `json:"artifact_type"`
	Selected        bool     `json:"selected"`
	SourceDocuments []string `json:"source_documents"`
	DraftContent    string   `json:"draft_content"`
}

type MappingFile struct {
	SchemaVersion int       `json:"schema_version"`
	ImportID      string    `json:"import_id"`
	Mappings      []Mapping `json:"mappings"`
}

// Traceability records the LLM-assisted review of each imported source. It is
// deliberately separate from mapping.json: mappings describe proposed writes,
// while this file records what was read, what was covered, and what remains
// unknown or unmapped.
type Traceability struct {
	SchemaVersion int              `json:"schema_version"`
	ImportID      string           `json:"import_id"`
	Status        string           `json:"status"`
	Sources       []SourceCoverage `json:"sources"`
}

type SourceCoverage struct {
	Path              string     `json:"path"`
	SHA256            string     `json:"sha256"`
	ReviewStatus      string     `json:"review_status"`
	Evidence          []Evidence `json:"evidence"`
	ExtractedClaims   []string   `json:"extracted_claims"`
	CandidateIDs      []string   `json:"candidate_ids"`
	MappedTargets     []string   `json:"mapped_targets"`
	MaterializedPaths []string   `json:"materialized_paths"`
	Gaps              []string   `json:"gaps"`
	Notes             string     `json:"notes"`
}

type Evidence struct {
	Locator string `json:"locator"`
	Claim   string `json:"claim"`
}

// CreateRun copies regular files into the product import area and creates an
// analysis-only run. It deliberately does not infer or materialize product artifacts.
func CreateRun(productRoot string, inputs []string) (string, error) {
	if len(inputs) == 0 {
		return "", fmt.Errorf("existing-documents requires at least one source path")
	}
	base := filepath.Join(productRoot, "knowledge", "imports")
	sourceRoot := filepath.Join(base, "sources")
	if err := os.MkdirAll(sourceRoot, 0755); err != nil {
		return "", err
	}
	files, err := expand(inputs)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no regular source files found")
	}
	runID, err := nextRun(filepath.Join(base, "runs"))
	if err != nil {
		return "", err
	}
	runRoot := filepath.Join(base, "runs", runID)
	if err := os.MkdirAll(runRoot, 0755); err != nil {
		return "", err
	}

	inv := Inventory{SchemaVersion: 1, ImportID: runID}
	for _, input := range files {
		name := uniqueName(sourceRoot, filepath.Base(input))
		dest := filepath.Join(sourceRoot, name)
		if err := copyFile(input, dest); err != nil {
			return "", err
		}
		info, _ := os.Stat(dest)
		hash, err := fileHash(dest)
		if err != nil {
			return "", err
		}
		rel, _ := filepath.Rel(productRoot, dest)
		inv.Sources = append(inv.Sources, Source{Path: filepath.ToSlash(rel), Format: strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), "."), Size: info.Size(), SHA256: hash})
	}
	sort.Slice(inv.Sources, func(i, j int) bool { return inv.Sources[i].Path < inv.Sources[j].Path })
	if err := writeJSON(filepath.Join(runRoot, "inventory.json"), inv); err != nil {
		return "", err
	}
	plan := map[string]any{"schema_version": 1, "import_id": runID, "status": "draft", "sources": inv.Sources, "candidates": []any{}, "open_questions": []string{"Review and classify the inventoried sources."}, "materialization_approved": false}
	if err := writeJSON(filepath.Join(runRoot, "import-plan.json"), plan); err != nil {
		return "", err
	}
	if err := writeJSON(filepath.Join(runRoot, "mapping.json"), map[string]any{"schema_version": 1, "import_id": runID, "mappings": []any{}}); err != nil {
		return "", err
	}
	coverage := make([]SourceCoverage, 0, len(inv.Sources))
	for _, source := range inv.Sources {
		coverage = append(coverage, SourceCoverage{Path: source.Path, SHA256: source.SHA256, ReviewStatus: "unreviewed", Evidence: []Evidence{}, ExtractedClaims: []string{}, CandidateIDs: []string{}, MappedTargets: []string{}, MaterializedPaths: []string{}, Gaps: []string{}})
	}
	traceability := Traceability{SchemaVersion: 1, ImportID: runID, Status: "unreviewed", Sources: coverage}
	if err := writeJSON(filepath.Join(runRoot, "traceability.json"), traceability); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(runRoot, "conflicts.md"), []byte("# Import Conflicts\n\nNo conflicts have been classified yet.\n"), 0644); err != nil {
		return "", err
	}
	report := fmt.Sprintf("# Import Report\n\n| Field | Value |\n| --- | --- |\n| Import | `%s` |\n| Status | `draft` |\n| Sources | `%d` |\n| Candidates | `0` |\n\n## Next step\n\nUse the Artifact Importer to classify sources and propose mappings. No canonical product artifacts were created.\n", runID, len(inv.Sources))
	if err := os.WriteFile(filepath.Join(runRoot, "import-report.md"), []byte(report), 0644); err != nil {
		return "", err
	}
	return runID, nil
}

// Materialize creates only explicitly selected draft mappings. The caller must
// provide the approving human identity; existing files are never overwritten.
func Materialize(productRoot, runID, approvedBy string) ([]string, error) {
	runID = strings.TrimSpace(runID)
	approvedBy = strings.TrimSpace(approvedBy)
	if runID == "" || approvedBy == "" {
		return nil, fmt.Errorf("run id and approved-by are required")
	}
	if !strings.HasPrefix(runID, "IMPORT-") || strings.ContainsAny(runID, `/\\`) {
		return nil, fmt.Errorf("invalid import id %q", runID)
	}
	runRoot := filepath.Join(productRoot, "knowledge", "imports", "runs", runID)
	data, err := os.ReadFile(filepath.Join(runRoot, "mapping.json"))
	if err != nil {
		return nil, err
	}
	var file MappingFile
	if err := json.Unmarshal(trimBOM(data), &file); err != nil {
		return nil, fmt.Errorf("parse mapping: %w", err)
	}
	if file.ImportID != runID {
		return nil, fmt.Errorf("mapping import_id %q does not match %q", file.ImportID, runID)
	}

	tracePath := filepath.Join(runRoot, "traceability.json")
	traceData, traceErr := os.ReadFile(tracePath)
	var traceability Traceability
	if traceErr == nil {
		if err := json.Unmarshal(trimBOM(traceData), &traceability); err != nil {
			return nil, fmt.Errorf("parse traceability: %w", err)
		}
		if traceability.ImportID != runID {
			return nil, fmt.Errorf("traceability import_id %q does not match %q", traceability.ImportID, runID)
		}
	}
	known, scalable, err := materializableSources(productRoot, runID)
	if err != nil {
		return nil, err
	}
	if scalable {
		if err := requireReviewedChunks(productRoot, runID); err != nil {
			return nil, err
		}
	}
	var selected []Mapping
	seen := map[string]bool{}
	for _, mapping := range file.Mappings {
		if !mapping.Selected {
			continue
		}
		if err := validateMapping(mapping); err != nil {
			return nil, fmt.Errorf("mapping %s: %w", mapping.ID, err)
		}
		for _, source := range mapping.SourceDocuments {
			if !known[source] {
				return nil, fmt.Errorf("mapping %s references uninventoried source %s", mapping.ID, source)
			}
		}
		clean := filepath.Clean(filepath.FromSlash(mapping.Target))
		if seen[strings.ToLower(clean)] {
			return nil, fmt.Errorf("duplicate selected target %s", mapping.Target)
		}
		seen[strings.ToLower(clean)] = true
		dest := filepath.Join(productRoot, clean)
		rel, err := filepath.Rel(productRoot, dest)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("mapping %s escapes product root", mapping.ID)
		}
		if _, err := os.Stat(dest); err == nil {
			return nil, fmt.Errorf("target already exists: %s", mapping.Target)
		} else if !os.IsNotExist(err) {
			return nil, err
		}
		selected = append(selected, mapping)
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("no selected mappings to materialize")
	}
	var created []string
	rollback := func() {
		for i := len(created) - 1; i >= 0; i-- {
			_ = os.Remove(filepath.Join(productRoot, filepath.FromSlash(created[i])))
		}
	}
	for _, mapping := range selected {
		dest := filepath.Join(productRoot, filepath.FromSlash(mapping.Target))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			rollback()
			return nil, err
		}
		content := markImportDraft(mapping.DraftContent, runID)
		if err := os.WriteFile(dest, []byte(content), 0644); err != nil {
			rollback()
			return nil, err
		}
		created = append(created, filepath.ToSlash(mapping.Target))
	}
	planPath := filepath.Join(runRoot, "import-plan.json")
	planData, err := os.ReadFile(planPath)
	if err != nil {
		rollback()
		return nil, err
	}
	var plan map[string]any
	if err := json.Unmarshal(trimBOM(planData), &plan); err != nil {
		rollback()
		return nil, err
	}
	plan["materialization_approved"] = true
	plan["materialization_approved_by"] = approvedBy
	plan["materialization_approved_at"] = time.Now().UTC().Format(time.RFC3339)
	plan["materialized_paths"] = created
	hashes := map[string]string{}
	for _, mapping := range selected {
		sum := sha256.Sum256([]byte(markImportDraft(mapping.DraftContent, runID)))
		hashes[filepath.ToSlash(mapping.Target)] = hex.EncodeToString(sum[:])
	}
	plan["materialized_hashes"] = hashes
	plan["status"] = "materialized"
	if err := writeJSON(planPath, plan); err != nil {
		rollback()
		return nil, err
	}
	if traceErr == nil {
		bySource := map[string][]string{}
		for _, mapping := range selected {
			for _, source := range mapping.SourceDocuments {
				bySource[source] = appendUnique(bySource[source], mapping.Target)
			}
		}
		for i := range traceability.Sources {
			source := &traceability.Sources[i]
			for _, target := range bySource[source.Path] {
				source.MappedTargets = appendUnique(source.MappedTargets, target)
				source.MaterializedPaths = appendUnique(source.MaterializedPaths, target)
			}
			if len(source.MaterializedPaths) > 0 {
				source.ReviewStatus = "materialized_as_draft"
			}
		}
		traceability.Status = "materialized_as_draft"
		if err := writeJSON(tracePath, traceability); err != nil {
			return nil, err
		}
	}
	return created, nil
}

// markImportDraft makes the intermediate import format machine-readable while
// preserving the imported body. Normalization skills replace kind with
// skill-normalized and record their owner before approval.
func markImportDraft(content, runID string) string {
	provenance := "provenance:\n  kind: import-draft\n  import_run: " + runID + "\n  normalized_by_skill: \"\"\n"
	if strings.HasPrefix(content, "---\n") {
		if end := strings.Index(content[4:], "\n---"); end >= 0 {
			at := 4 + end
			return content[:at] + "\n" + provenance + content[at:]
		}
	}
	return "---\n" + provenance + "---\n\n" + content
}

// NormalizeProvenance promotes an already template-conformant imported
// artifact without changing its product content or lifecycle status.
func NormalizeProvenance(content, skill string) (string, error) {
	skill = strings.TrimSpace(skill)
	if skill == "" {
		return "", fmt.Errorf("normalizing skill is required")
	}
	if !strings.Contains(content, "kind: import-draft") {
		return "", fmt.Errorf("artifact is not marked as import-draft")
	}
	kindPattern := regexp.MustCompile(`(?m)^(\s*kind:\s*)import-draft\s*$`)
	updated := kindPattern.ReplaceAllString(content, "${1}skill-normalized")
	skillPattern := regexp.MustCompile(`(?m)^(\s*normalized_by_skill:\s*).*$`)
	if !skillPattern.MatchString(updated) {
		return "", fmt.Errorf("artifact provenance is missing normalized_by_skill")
	}
	updated = skillPattern.ReplaceAllString(updated, "${1}"+skill)
	return updated, nil
}

// RecordNormalization preserves the original materialization hash while
// recording the current hash after an owning skill changes the draft.
func RecordNormalization(productRoot, artifactPath, content string) error {
	match := regexp.MustCompile(`(?m)^\s*import_run:\s*([^\s]+)`).FindStringSubmatch(content)
	if len(match) != 2 {
		return fmt.Errorf("normalized artifact is missing provenance.import_run")
	}
	planPath := filepath.Join(productRoot, "knowledge", "imports", "runs", match[1], "import-plan.json")
	data, err := os.ReadFile(planPath)
	if err != nil {
		return err
	}
	var plan map[string]any
	if err := json.Unmarshal(trimBOM(data), &plan); err != nil {
		return err
	}
	rel, err := filepath.Rel(productRoot, artifactPath)
	if err != nil {
		return err
	}
	hashes, _ := plan["normalized_hashes"].(map[string]any)
	if hashes == nil {
		hashes = map[string]any{}
	}
	sum := sha256.Sum256([]byte(content))
	hashes[filepath.ToSlash(rel)] = hex.EncodeToString(sum[:])
	plan["normalized_hashes"] = hashes
	return writeJSON(planPath, plan)
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func validateMapping(mapping Mapping) error {
	if strings.TrimSpace(mapping.ID) == "" {
		return fmt.Errorf("id is required")
	}
	if strings.TrimSpace(mapping.Target) == "" {
		return fmt.Errorf("target is required")
	}
	if len(mapping.SourceDocuments) == 0 {
		return fmt.Errorf("source_documents is required")
	}
	content := strings.TrimSpace(mapping.DraftContent)
	if content == "" {
		return fmt.Errorf("draft_content is required")
	}
	if !strings.Contains(content, "status: draft") && !strings.Contains(content, "| Status | `draft` |") {
		return fmt.Errorf("draft_content must declare draft status")
	}
	if !strings.Contains(content, "source_documents") {
		return fmt.Errorf("draft_content must preserve source_documents")
	}
	return nil
}

func expand(inputs []string) ([]string, error) {
	var out []string
	for _, input := range inputs {
		abs, err := filepath.Abs(strings.TrimSpace(input))
		if err != nil {
			return nil, err
		}
		info, err := os.Stat(abs)
		if err != nil {
			return nil, fmt.Errorf("source %s: %w", input, err)
		}
		if !info.IsDir() {
			out = append(out, abs)
			continue
		}
		err = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				out = append(out, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(out)
	return out, nil
}

func nextRun(root string) (string, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", err
	}
	max := 0
	for _, e := range entries {
		var n int
		if _, err := fmt.Sscanf(e.Name(), "IMPORT-%03d", &n); err == nil && n > max {
			max = n
		}
	}
	return fmt.Sprintf("IMPORT-%03d", max+1), nil
}

func uniqueName(root, name string) string {
	if _, err := os.Stat(filepath.Join(root, name)); os.IsNotExist(err) {
		return name
	}
	ext, base := filepath.Ext(name), strings.TrimSuffix(name, filepath.Ext(name))
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s-%d%s", base, i, ext)
		if _, err := os.Stat(filepath.Join(root, candidate)); os.IsNotExist(err) {
			return candidate
		}
	}
}

func copyFile(source, dest string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func trimBOM(data []byte) []byte {
	return bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf})
}
