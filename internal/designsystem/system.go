package designsystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Token struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}

type TokenDocument struct {
	SchemaVersion int            `json:"schemaVersion"`
	System        string         `json:"system"`
	Version       string         `json:"version"`
	Tokens        map[string]any `json:"tokens"`
}

type Inspection struct {
	ID         string   `json:"id"`
	Status     string   `json:"status"`
	Version    string   `json:"version"`
	OriginMode string   `json:"originMode"`
	Tokens     int      `json:"tokens"`
	Components int      `json:"components"`
	Patterns   int      `json:"patterns"`
	Sources    int      `json:"sources"`
	Blockers   []string `json:"blockers,omitempty"`
}

type Impact struct {
	ID        string   `json:"id"`
	Version   string   `json:"version"`
	Consumers []string `json:"consumers"`
	Blockers  []string `json:"blockers,omitempty"`
}

func Init(root, mode string) (string, error) {
	if !oneOf(mode, "generate", "evolve", "adopt") {
		return "", fmt.Errorf("invalid Design System mode %q", mode)
	}
	dir := filepath.Join(root, "design", "system")
	if _, err := os.Stat(filepath.Join(dir, "context.md")); err == nil {
		return "", errors.New("Design System already exists")
	}
	for _, child := range []string{"foundations", "tokens", "components", "patterns", "sources", "evidence"} {
		if err := os.MkdirAll(filepath.Join(dir, child), 0o755); err != nil {
			return "", err
		}
	}
	context := fmt.Sprintf("```yaml\nid: DSYS-001\ntype: design-system\nname: Product Design System\nstatus: draft\nowner_skill: design-system\nslug: system\norigin_mode: %s\nversion: 0.1.0\ndelivery:\n  level: L0\n  priority: P1\n  depends_on:\n    - foundation/vision\n    - foundation/strategy\n  rationale: Shared foundations for consistent product interfaces.\ndocuments:\n  canonical: design-system.md\n  tokens: tokens/tokens.json\n  themes: tokens/themes.json\nsources: []\ndecisions: []\nopen_questions: []\n```\n\n# Context: Product Design System\n", mode)
	if err := os.WriteFile(filepath.Join(dir, "context.md"), []byte(context), 0o644); err != nil {
		return "", err
	}
	canonical := "# Design System: Product\n\n## Snapshot\n\n| Field | Value |\n| --- | --- |\n| ID | `DSYS-001` |\n| Status | `draft` |\n| Version | `0.1.0` |\n| Origin | `" + mode + "` |\n\n## Purpose And Principles\n\nDraft.\n\n## Tokens And Themes\n\n- [Tokens](tokens/tokens.json)\n- [Themes](tokens/themes.json)\n\n## Approval\n\nNo approval has been granted.\n"
	if err := os.WriteFile(filepath.Join(dir, "design-system.md"), []byte(canonical), 0o644); err != nil {
		return "", err
	}
	if err := writeJSON(filepath.Join(dir, "tokens", "tokens.json"), TokenDocument{SchemaVersion: 1, System: "DSYS-001", Version: "0.1.0", Tokens: map[string]any{}}); err != nil {
		return "", err
	}
	if err := writeJSON(filepath.Join(dir, "tokens", "themes.json"), map[string]any{"schemaVersion": 1, "system": "DSYS-001", "version": "0.1.0", "themes": map[string]any{}}); err != nil {
		return "", err
	}
	return filepath.ToSlash(dir), nil
}

func Inspect(root string) (Inspection, error) {
	dir := filepath.Join(root, "design", "system")
	context, err := os.ReadFile(filepath.Join(dir, "context.md"))
	if err != nil {
		return Inspection{}, err
	}
	i := Inspection{ID: field(string(context), "id"), Status: field(string(context), "status"), Version: field(string(context), "version"), OriginMode: field(string(context), "origin_mode")}
	var doc TokenDocument
	data, err := os.ReadFile(filepath.Join(dir, "tokens", "tokens.json"))
	if err != nil {
		i.Blockers = append(i.Blockers, "tokens/tokens.json is missing")
	} else if err := json.Unmarshal(data, &doc); err != nil {
		i.Blockers = append(i.Blockers, "tokens/tokens.json is invalid JSON")
	} else {
		flat := map[string]Token{}
		flatten("", doc.Tokens, flat)
		i.Tokens = len(flat)
		i.Blockers = append(i.Blockers, validateTokens(doc, flat)...)
	}
	i.Components = markdownCount(filepath.Join(dir, "components"))
	i.Patterns = markdownCount(filepath.Join(dir, "patterns"))
	i.Sources = manifestCount(filepath.Join(dir, "sources"))
	if i.ID == "" || i.Version == "" || !oneOf(i.OriginMode, "generate", "evolve", "adopt") {
		i.Blockers = append(i.Blockers, "context identity, version, or origin mode is invalid")
	}
	sort.Strings(i.Blockers)
	return i, nil
}

func Validate(root string) (Inspection, error) { return Inspect(root) }

func ImpactReport(root, id string) (Impact, error) {
	i, err := Inspect(root)
	if err != nil {
		return Impact{}, err
	}
	if id != "" && id != i.ID {
		return Impact{}, fmt.Errorf("Design System %s is not declared", id)
	}
	report := Impact{ID: i.ID, Version: i.Version, Blockers: i.Blockers}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || strings.Contains(filepath.ToSlash(path), "/design/system/") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), i.ID) {
			rel, _ := filepath.Rel(root, path)
			report.Consumers = append(report.Consumers, filepath.ToSlash(rel))
		}
		return nil
	})
	sort.Strings(report.Consumers)
	return report, err
}

func Migrate(root string, dryRun bool) ([]string, error) {
	path := filepath.Join(root, "design", "system", "context.md")
	if _, err := os.Stat(path); err == nil {
		return []string{"Design System already exists"}, nil
	}
	if dryRun {
		return []string{"CREATE design/system as generate/draft DSYS-001"}, nil
	}
	created, err := Init(root, "generate")
	if err != nil {
		return nil, err
	}
	return []string{"CREATED " + created}, nil
}

func validateTokens(doc TokenDocument, flat map[string]Token) []string {
	var out []string
	if doc.SchemaVersion != 1 || !regexp.MustCompile(`^DSYS-[0-9]{3,}$`).MatchString(doc.System) || !semver(doc.Version) {
		out = append(out, "token document schema, system id, or semantic version is invalid")
	}
	alias := regexp.MustCompile(`^\{([^}]+)\}$`)
	deps := map[string]string{}
	for path, token := range flat {
		if token.Type == "" {
			out = append(out, "token "+path+" is missing type")
		}
		if value, ok := token.Value.(string); ok {
			if match := alias.FindStringSubmatch(value); len(match) == 2 {
				deps[path] = match[1]
				if _, exists := flat[match[1]]; !exists {
					out = append(out, "token "+path+" references missing token "+match[1])
				}
			}
		}
	}
	for start := range deps {
		seen := map[string]bool{}
		for at := start; deps[at] != ""; at = deps[at] {
			if seen[at] {
				out = append(out, "token alias cycle includes "+at)
				break
			}
			seen[at] = true
		}
	}
	return unique(out)
}

func flatten(prefix string, value map[string]any, out map[string]Token) {
	for key, raw := range value {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		object, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if tokenValue, exists := object["value"]; exists {
			out[path] = Token{Value: tokenValue, Type: fmt.Sprint(object["type"])}
			continue
		}
		flatten(path, object, out)
	}
}

func field(text, name string) string {
	re := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(name) + `:\s*([^\r\n#]+)`)
	m := re.FindStringSubmatch(text)
	if len(m) != 2 {
		return ""
	}
	return strings.Trim(strings.TrimSpace(m[1]), "`\"'")
}

func markdownCount(dir string) int {
	entries, _ := os.ReadDir(dir)
	n := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") && strings.ToLower(entry.Name()) != "readme.md" {
			n++
		}
	}
	return n
}

func manifestCount(dir string) int {
	entries, _ := os.ReadDir(dir)
	n := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			n++
		}
	}
	return n
}

func semver(value string) bool {
	return regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`).MatchString(value)
}
func oneOf(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}
func unique(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			out = append(out, item)
		}
	}
	return out
}
func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
