package sourceimport

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const scalableSchemaVersion = 2

type CreateOptions struct {
	Include       []string `json:"include,omitempty"`
	Exclude       []string `json:"exclude,omitempty"`
	MaxFiles      int      `json:"max_files"`
	MaxTotalBytes int64    `json:"max_total_bytes"`
	MaxFileBytes  int64    `json:"max_file_bytes"`
	ChunkSize     int      `json:"chunk_size"`
	BinaryPolicy  string   `json:"binary_policy"`
}

type ScalableSource struct {
	ID           string `json:"id"`
	OriginalPath string `json:"original_path"`
	Path         string `json:"path"`
	Format       string `json:"format"`
	Size         int64  `json:"size"`
	SHA256       string `json:"sha256"`
	Status       string `json:"status"`
	Reason       string `json:"reason,omitempty"`
}
type InventoryIndex struct {
	SchemaVersion int    `json:"schema_version"`
	ImportID      string `json:"import_id"`
	Pages         int    `json:"pages"`
	Sources       int    `json:"sources"`
	TotalBytes    int64  `json:"total_bytes"`
	ConfigHash    string `json:"config_hash"`
}
type Chunk struct {
	SchemaVersion int      `json:"schema_version"`
	ID            string   `json:"id"`
	SourceIDs     []string `json:"source_ids"`
	Status        string   `json:"status"`
	Agent         string   `json:"agent,omitempty"`
	LeaseExpires  string   `json:"lease_expires,omitempty"`
}
type ScalableStatus struct {
	ImportID                                                        string `json:"import_id"`
	Sources, Chunks, Queued, Reviewing, Reviewed, Blocked, Excluded int
}
type ChunkReview struct {
	SourceEvidence map[string][]Evidence `json:"source_evidence"`
	Gaps           map[string][]string   `json:"gaps,omitempty"`
}

func DefaultCreateOptions() CreateOptions {
	return CreateOptions{MaxFiles: 500, MaxTotalBytes: 200 << 20, MaxFileBytes: 10 << 20, ChunkSize: 25, BinaryPolicy: "inventory_only", Exclude: []string{".git/**", "node_modules/**", "**/.git/**", "**/node_modules/**"}}
}
func (o CreateOptions) normalized() (CreateOptions, error) {
	d := DefaultCreateOptions()
	if o.MaxFiles == 0 {
		o.MaxFiles = d.MaxFiles
	}
	if o.MaxTotalBytes == 0 {
		o.MaxTotalBytes = d.MaxTotalBytes
	}
	if o.MaxFileBytes == 0 {
		o.MaxFileBytes = d.MaxFileBytes
	}
	if o.ChunkSize == 0 {
		o.ChunkSize = d.ChunkSize
	}
	if o.BinaryPolicy == "" {
		o.BinaryPolicy = d.BinaryPolicy
	}
	if len(o.Exclude) == 0 {
		o.Exclude = d.Exclude
	}
	if o.MaxFiles < 1 || o.MaxTotalBytes < 1 || o.MaxFileBytes < 1 || o.ChunkSize < 1 {
		return o, errors.New("import limits and chunk size must be positive")
	}
	if o.BinaryPolicy != "inventory_only" && o.BinaryPolicy != "reject" {
		return o, errors.New("binary policy must be inventory_only or reject")
	}
	return o, nil
}

// CreateScalableRun creates a v2 analysis-only run. It checks budgets before
// copying, writes paged JSONL inventory, and never creates product artifacts.
func CreateScalableRun(productRoot string, inputs []string, options CreateOptions) (string, error) {
	options, err := options.normalized()
	if err != nil {
		return "", err
	}
	files, err := expand(inputs)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("no regular source files found")
	}
	var chosen []string
	var total int64
	for _, file := range files {
		rel := filepath.ToSlash(file)
		if matchesAny(rel, options.Exclude) || (len(options.Include) > 0 && !matchesAny(rel, options.Include)) {
			continue
		}
		info, err := os.Stat(file)
		if err != nil {
			return "", err
		}
		if info.Size() > options.MaxFileBytes {
			return "", fmt.Errorf("source exceeds max file bytes: %s", file)
		}
		if len(chosen)+1 > options.MaxFiles || total+info.Size() > options.MaxTotalBytes {
			return "", fmt.Errorf("import budget exceeded before copying %s", file)
		}
		chosen, total = append(chosen, file), total+info.Size()
	}
	if len(chosen) == 0 {
		return "", errors.New("no sources matched import filters")
	}
	base := filepath.Join(productRoot, "knowledge", "imports")
	runID, err := nextRun(filepath.Join(base, "runs"))
	if err != nil {
		return "", err
	}
	runRoot := filepath.Join(base, "runs", runID)
	sourceRoot := filepath.Join(base, "sources", runID)
	for _, dir := range []string{runRoot, sourceRoot, filepath.Join(runRoot, "inventory", "pages"), filepath.Join(runRoot, "chunks")} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}
	if err := writeJSON(filepath.Join(runRoot, "import-config.json"), options); err != nil {
		return "", err
	}
	pageSize := options.ChunkSize
	var page []ScalableSource
	var sources []ScalableSource
	for index, input := range chosen {
		name := fmt.Sprintf("%06d-%s", index+1, filepath.Base(input))
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
		source := ScalableSource{ID: fmt.Sprintf("SRC-%06d", index+1), OriginalPath: input, Path: filepath.ToSlash(rel), Format: strings.TrimPrefix(strings.ToLower(filepath.Ext(input)), "."), Size: info.Size(), SHA256: hash, Status: "queued"}
		if isBinary(input) {
			if options.BinaryPolicy == "reject" {
				return "", fmt.Errorf("binary source rejected: %s", input)
			}
			source.Status, source.Reason = "excluded", "binary inventory only"
		}
		sources, page = append(sources, source), append(page, source)
		if len(page) == pageSize || index == len(chosen)-1 {
			if err := writeJSONL(filepath.Join(runRoot, "inventory", "pages", fmt.Sprintf("PAGE-%04d.jsonl", (index/pageSize)+1)), page); err != nil {
				return "", err
			}
			page = nil
		}
	}
	configData, _ := json.Marshal(options)
	sum := sha256.Sum256(configData)
	index := InventoryIndex{SchemaVersion: scalableSchemaVersion, ImportID: runID, Pages: (len(sources) + pageSize - 1) / pageSize, Sources: len(sources), TotalBytes: total, ConfigHash: hex.EncodeToString(sum[:])}
	if err := writeJSON(filepath.Join(runRoot, "inventory", "index.json"), index); err != nil {
		return "", err
	}
	for offset := 0; offset < len(sources); offset += options.ChunkSize {
		end := offset + options.ChunkSize
		if end > len(sources) {
			end = len(sources)
		}
		ids := make([]string, 0, end-offset)
		status := "queued"
		for _, source := range sources[offset:end] {
			ids = append(ids, source.ID)
			if source.Status == "excluded" {
				status = "excluded"
			}
		}
		if err := writeJSON(filepath.Join(runRoot, "chunks", fmt.Sprintf("CHUNK-%04d.json", offset/options.ChunkSize+1)), Chunk{SchemaVersion: scalableSchemaVersion, ID: fmt.Sprintf("CHUNK-%04d", offset/options.ChunkSize+1), SourceIDs: ids, Status: status}); err != nil {
			return "", err
		}
	}
	if err := writeJSON(filepath.Join(runRoot, "mapping.json"), MappingFile{SchemaVersion: scalableSchemaVersion, ImportID: runID}); err != nil {
		return "", err
	}
	return runID, writeJSON(filepath.Join(runRoot, "import-plan.json"), map[string]any{"schema_version": scalableSchemaVersion, "import_id": runID, "status": "draft", "materialization_approved": false})
}

func ImportStatus(productRoot, runID string) (ScalableStatus, error) {
	runRoot := filepath.Join(productRoot, "knowledge", "imports", "runs", runID)
	indexData, err := os.ReadFile(filepath.Join(runRoot, "inventory", "index.json"))
	if err != nil {
		return ScalableStatus{}, err
	}
	var index InventoryIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		return ScalableStatus{}, err
	}
	status := ScalableStatus{ImportID: runID, Sources: index.Sources}
	entries, err := os.ReadDir(filepath.Join(runRoot, "chunks"))
	if err != nil {
		return status, err
	}
	for _, entry := range entries {
		var chunk Chunk
		if readJSONFile(filepath.Join(runRoot, "chunks", entry.Name()), &chunk) != nil {
			continue
		}
		status.Chunks++
		switch chunk.Status {
		case "queued":
			status.Queued++
		case "reviewing":
			status.Reviewing++
		case "reviewed":
			status.Reviewed++
		case "blocked":
			status.Blocked++
		case "excluded":
			status.Excluded++
		}
	}
	return status, nil
}

// Resume claims one queued/expired chunk for the named importer. It only
// persists operational ownership; the skill still owns the review content.
func Resume(productRoot, runID, chunkID, agent string) (Chunk, error) {
	if strings.TrimSpace(agent) == "" {
		return Chunk{}, errors.New("import resume requires agent")
	}
	dir := filepath.Join(productRoot, "knowledge", "imports", "runs", runID, "chunks")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Chunk{}, err
	}
	for _, entry := range entries {
		var chunk Chunk
		path := filepath.Join(dir, entry.Name())
		if err := readJSONFile(path, &chunk); err != nil {
			return Chunk{}, err
		}
		if chunkID != "" && chunk.ID != chunkID {
			continue
		}
		expired := chunk.LeaseExpires != "" && mustParseTime(chunk.LeaseExpires).Before(time.Now().UTC())
		if chunk.Status != "queued" && !(chunk.Status == "reviewing" && expired) {
			continue
		}
		chunk.Status, chunk.Agent, chunk.LeaseExpires = "reviewing", agent, time.Now().UTC().Add(30*time.Minute).Format(time.RFC3339)
		return chunk, writeJSON(path, chunk)
	}
	return Chunk{}, errors.New("no resumable import chunk")
}

// RecordChunkReview records evidence for every non-excluded source in a leased
// chunk. It is deliberately unable to select mappings or materialize drafts.
func RecordChunkReview(productRoot, runID, chunkID, agent string, review ChunkReview) error {
	if strings.TrimSpace(agent) == "" {
		return errors.New("chunk review requires agent")
	}
	path := filepath.Join(productRoot, "knowledge", "imports", "runs", runID, "chunks", chunkID+".json")
	var chunk Chunk
	if err := readJSONFile(path, &chunk); err != nil {
		return err
	}
	if chunk.Status != "reviewing" || chunk.Agent != agent || mustParseTime(chunk.LeaseExpires).Before(time.Now().UTC()) {
		return errors.New("chunk is not actively leased by this importer")
	}
	sources, err := scalableSources(productRoot, runID)
	if err != nil {
		return err
	}
	sourceByID := map[string]ScalableSource{}
	for _, source := range sources {
		sourceByID[source.ID] = source
	}
	for _, id := range chunk.SourceIDs {
		source, ok := sourceByID[id]
		if !ok {
			return fmt.Errorf("chunk references unknown source %s", id)
		}
		if source.Status == "excluded" {
			continue
		}
		if len(review.SourceEvidence[id]) == 0 {
			return fmt.Errorf("review evidence is required for source %s", id)
		}
	}
	traceDir := filepath.Join(productRoot, "knowledge", "imports", "runs", runID, "traceability")
	if err := os.MkdirAll(traceDir, 0755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(traceDir, chunkID+".json"), review); err != nil {
		return err
	}
	chunk.Status, chunk.Agent, chunk.LeaseExpires = "reviewed", "", ""
	return writeJSON(path, chunk)
}

func scalableSources(productRoot, runID string) ([]ScalableSource, error) {
	runRoot := filepath.Join(productRoot, "knowledge", "imports", "runs", runID)
	data, err := os.ReadFile(filepath.Join(runRoot, "inventory", "index.json"))
	if err != nil {
		return nil, err
	}
	var index InventoryIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}
	var all []ScalableSource
	for page := 1; page <= index.Pages; page++ {
		file, err := os.Open(filepath.Join(runRoot, "inventory", "pages", fmt.Sprintf("PAGE-%04d.jsonl", page)))
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 1024), 4<<20)
		for scanner.Scan() {
			var source ScalableSource
			if err := json.Unmarshal(scanner.Bytes(), &source); err != nil {
				file.Close()
				return nil, err
			}
			all = append(all, source)
		}
		if err := scanner.Err(); err != nil {
			file.Close()
			return nil, err
		}
		file.Close()
	}
	return all, nil
}
func matchesAny(path string, patterns []string) bool {
	path = filepath.ToSlash(path)
	for _, pattern := range patterns {
		pattern = filepath.ToSlash(pattern)
		if ok, _ := filepath.Match(pattern, path); ok {
			return true
		}
		if strings.HasPrefix(pattern, "**/") {
			if ok, _ := filepath.Match(strings.TrimPrefix(pattern, "**/"), filepath.Base(path)); ok {
				return true
			}
		}
		if strings.HasSuffix(pattern, "/**") && strings.Contains(path, strings.TrimSuffix(pattern, "/**")) {
			return true
		}
	}
	return false
}
func isBinary(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".txt", ".json", ".yaml", ".yml", ".csv", ".html", ".xml":
		return false
	}
	return true
}
func writeJSONL(path string, items []ScalableSource) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, item := range items {
		data, _ := json.Marshal(item)
		if _, err := writer.Write(append(data, '\n')); err != nil {
			return err
		}
	}
	return writer.Flush()
}
func readJSONFile(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
func mustParseTime(value string) time.Time {
	parsed, _ := time.Parse(time.RFC3339, value)
	return parsed
}

var _ = io.EOF
var _ = sort.Strings
