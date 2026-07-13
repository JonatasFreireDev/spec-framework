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
	invData, err := os.ReadFile(filepath.Join(runRoot, "inventory.json"))
	if err != nil {
		return nil, err
	}
	var inventory Inventory
	if err := json.Unmarshal(trimBOM(invData), &inventory); err != nil {
		return nil, fmt.Errorf("parse inventory: %w", err)
	}
	knownSources := map[string]bool{}
	for _, source := range inventory.Sources {
		knownSources[source.Path] = true
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
			if !knownSources[source] {
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
		if err := os.WriteFile(dest, []byte(mapping.DraftContent), 0644); err != nil {
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
		sum := sha256.Sum256([]byte(mapping.DraftContent))
		hashes[filepath.ToSlash(mapping.Target)] = hex.EncodeToString(sum[:])
	}
	plan["materialized_hashes"] = hashes
	plan["status"] = "materialized"
	if err := writeJSON(planPath, plan); err != nil {
		rollback()
		return nil, err
	}
	return created, nil
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
