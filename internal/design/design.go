package design

import (
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
)

var imageExtensions = map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true, ".svg": true, ".pdf": true}

func Init(productRoot, useCase, mode string) (string, error) {
	if !oneOf(mode, "generate", "evolve", "adopt") {
		return "", fmt.Errorf("invalid design origin mode %q", mode)
	}
	uc, rel, err := resolveUseCase(productRoot, useCase)
	if err != nil {
		return "", err
	}
	slug := filepath.Base(uc)
	dir := filepath.Join(productRoot, "design", "use-cases", slug)
	for _, child := range []string{"wireframes", "mockups", "prototype", "evidence"} {
		if err := os.MkdirAll(filepath.Join(dir, child), 0o755); err != nil {
			return "", err
		}
	}
	m := UseCaseManifest{SchemaVersion: 1, UseCase: rel, OriginMode: mode, Maturity: "contract", FidelityPolicy: defaultFidelity(mode), Sources: []string{}, NonProduction: true}
	path := filepath.Join(dir, "manifest.json")
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("design manifest already exists: %s", filepath.ToSlash(path))
	}
	return filepath.ToSlash(path), writeJSON(path, m)
}

func ImportImages(productRoot, useCase, source, authority, sourceID string, copyAssets bool) (SourceManifest, string, error) {
	if !oneOf(authority, "visual_canonical", "reference", "inspiration") {
		return SourceManifest{}, "", fmt.Errorf("invalid source authority %q", authority)
	}
	uc, _, err := resolveUseCase(productRoot, useCase)
	if err != nil {
		return SourceManifest{}, "", err
	}
	if sourceID == "" {
		sourceID, err = nextSourceID(filepath.Join(productRoot, "design", "sources"))
		if err != nil {
			return SourceManifest{}, "", err
		}
	}
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return SourceManifest{}, "", err
	}
	files, err := visualFiles(sourceAbs)
	if err != nil {
		return SourceManifest{}, "", err
	}
	if len(files) == 0 {
		return SourceManifest{}, "", errors.New("source contains no supported visual files")
	}
	destination := filepath.Join(productRoot, "design", "sources", sourceID)
	assets := filepath.Join(destination, "assets")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		return SourceManifest{}, "", err
	}
	h := sha256.New()
	manifest := SourceManifest{SchemaVersion: 1, ID: sourceID, Type: "images", Authority: authority, Location: filepath.ToSlash(filepath.Join("design", "sources", sourceID, "assets")), Adapter: "builtin-images"}
	for i, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return SourceManifest{}, "", err
		}
		rel, _ := filepath.Rel(sourceAbs, file)
		if info, statErr := os.Stat(sourceAbs); statErr == nil && !info.IsDir() {
			rel = filepath.Base(file)
		}
		rel = filepath.Clean(rel)
		h.Write([]byte(filepath.ToSlash(rel)))
		h.Write(data)
		target := filepath.Join(assets, rel)
		if copyAssets {
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return SourceManifest{}, "", err
			}
			if err := os.WriteFile(target, data, 0o644); err != nil {
				return SourceManifest{}, "", err
			}
		}
		manifest.Screens = append(manifest.Screens, Screen{ID: fmt.Sprintf("SCREEN-%03d", i+1), Name: screenName(rel), Asset: filepath.ToSlash(filepath.Join("assets", rel))})
	}
	manifest.Version = Version{Kind: "sha256", Value: hex.EncodeToString(h.Sum(nil))}
	manifestPath := filepath.Join(destination, "manifest.json")
	if err := writeJSON(manifestPath, manifest); err != nil {
		return SourceManifest{}, "", err
	}
	if err := attachSource(productRoot, filepath.Base(uc), sourceID, authority); err != nil {
		return SourceManifest{}, "", err
	}
	return manifest, filepath.ToSlash(manifestPath), nil
}

func Inspect(productRoot, useCase string) (Inspection, error) {
	uc, _, err := resolveUseCase(productRoot, useCase)
	if err != nil {
		return Inspection{}, err
	}
	path := filepath.Join(productRoot, "design", "use-cases", filepath.Base(uc), "manifest.json")
	var m UseCaseManifest
	if err := readJSON(path, &m); err != nil {
		return Inspection{}, err
	}
	i := Inspection{UseCase: m.UseCase, OriginMode: m.OriginMode, Maturity: m.Maturity, FidelityPolicy: m.FidelityPolicy, Sources: m.Sources, Mappings: len(m.Mappings)}
	for _, id := range m.Sources {
		var source SourceManifest
		if err := readJSON(filepath.Join(productRoot, "design", "sources", id, "manifest.json"), &source); err != nil {
			i.Blockers = append(i.Blockers, "missing source manifest: "+id)
			continue
		}
		i.Screens += len(source.Screens)
		if source.Authority == "visual_canonical" && source.Version.Value == "" {
			i.Blockers = append(i.Blockers, "canonical source is unversioned: "+id)
		}
	}
	for _, mapping := range m.Mappings {
		if oneOf(mapping.Coverage, "missing", "conflict") {
			i.Blockers = append(i.Blockers, mapping.Requirement+" is "+mapping.Coverage)
		}
	}
	sort.Strings(i.Blockers)
	return i, nil
}

func UpdateMappings(productRoot, useCase string, mappings []Mapping) (string, error) {
	uc, _, err := resolveUseCase(productRoot, useCase)
	if err != nil {
		return "", err
	}
	path := filepath.Join(productRoot, "design", "use-cases", filepath.Base(uc), "manifest.json")
	var m UseCaseManifest
	if err := readJSON(path, &m); err != nil {
		return "", err
	}
	for _, item := range mappings {
		if item.Requirement == "" || !oneOf(item.Coverage, "covered", "partial", "missing", "conflict", "not-verifiable", "not-applicable") {
			return "", fmt.Errorf("invalid mapping for requirement %q", item.Requirement)
		}
	}
	m.Mappings = mappings
	return filepath.ToSlash(path), writeJSON(path, m)
}

func Verify(productRoot, useCase string) (Inspection, error) { return Inspect(productRoot, useCase) }

func attachSource(root, slug, sourceID, authority string) error {
	path := filepath.Join(root, "design", "use-cases", slug, "manifest.json")
	var m UseCaseManifest
	if err := readJSON(path, &m); err != nil {
		return fmt.Errorf("initialize design before importing: %w", err)
	}
	if !contains(m.Sources, sourceID) {
		m.Sources = append(m.Sources, sourceID)
		sort.Strings(m.Sources)
	}
	if authority == "visual_canonical" {
		m.OriginMode = "adopt"
		m.FidelityPolicy = "strict"
		m.Maturity = "mockup"
	}
	return writeJSON(path, m)
}

func resolveUseCase(root, value string) (string, string, error) {
	if value == "" {
		return "", "", errors.New("use case is required")
	}
	p := value
	if !filepath.IsAbs(p) {
		p = filepath.Join(root, filepath.FromSlash(value))
	}
	p, _ = filepath.Abs(p)
	r, _ := filepath.Abs(root)
	rel, err := filepath.Rel(r, p)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", "", errors.New("use case path escapes product root")
	}
	info, err := os.Stat(p)
	if err != nil || !info.IsDir() {
		return "", "", fmt.Errorf("use case directory not found: %s", filepath.ToSlash(value))
	}
	if _, err := os.Stat(filepath.Join(p, "specification.md")); err != nil {
		return "", "", errors.New("use case must contain specification.md")
	}
	return p, filepath.ToSlash(rel), nil
}

func visualFiles(source string) ([]string, error) {
	info, err := os.Stat(source)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if !imageExtensions[strings.ToLower(filepath.Ext(source))] {
			return nil, fmt.Errorf("unsupported visual file: %s", source)
		}
		return []string{source}, nil
	}
	var out []string
	err = filepath.WalkDir(source, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && imageExtensions[strings.ToLower(filepath.Ext(path))] {
			out = append(out, path)
		}
		return err
	})
	sort.Strings(out)
	return out, err
}

func nextSourceID(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return "DSRC-001", nil
	}
	if err != nil {
		return "", err
	}
	used := map[string]bool{}
	for _, e := range entries {
		used[e.Name()] = true
	}
	for n := 1; ; n++ {
		id := fmt.Sprintf("DSRC-%03d", n)
		if !used[id] {
			return id, nil
		}
	}
}

func HashFile(path string) (string, error) {
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

func readJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func oneOf(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func defaultFidelity(mode string) string {
	if mode == "adopt" {
		return "strict"
	}
	if mode == "evolve" {
		return "balanced"
	}
	return "exploratory"
}

func screenName(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strings.ReplaceAll(strings.ReplaceAll(name, "-", " "), "_", " ")
	return strings.Title(name) //nolint:staticcheck
}
