package engineeringsystem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

var allowedTriggers = map[string]bool{
	"new_dependency":               true,
	"external_integration":         true,
	"data_ownership_change":        true,
	"migration":                    true,
	"architecture_boundary_change": true,
	"deployment_change":            true,
	"security_boundary_change":     true,
	"operational_change":           true,
}

var allowedMaturity = map[string]bool{
	"baseline": true,
	"mapped":   true,
	"governed": true,
	"verified": true,
	"operated": true,
}

type Area struct {
	Name     string `json:"name"`
	Contract string `json:"contract"`
	Maturity string `json:"maturity"`
	Evidence int    `json:"evidence"`
}

type Inspection struct {
	ID               string   `json:"id"`
	Status           string   `json:"status"`
	Version          string   `json:"version"`
	OriginMode       string   `json:"originMode"`
	Scope            string   `json:"scope"`
	Areas            []Area   `json:"areas"`
	Decisions        int      `json:"decisions"`
	Standards        int      `json:"standards"`
	FitnessFunctions int      `json:"fitnessFunctions"`
	Blockers         []string `json:"blockers,omitempty"`
}

type contextDocument struct {
	ID                  string   `yaml:"id"`
	Status              string   `yaml:"status"`
	Version             string   `yaml:"version"`
	OriginMode          string   `yaml:"origin_mode"`
	EngineeringTriggers []string `yaml:"engineering_triggers"`
}

type catalogDocument struct {
	SchemaVersion    int                    `yaml:"schema_version"`
	ID               string                 `yaml:"id"`
	Status           string                 `yaml:"status"`
	Version          string                 `yaml:"version"`
	OriginMode       string                 `yaml:"origin_mode"`
	Scope            string                 `yaml:"scope"`
	Areas            map[string]catalogArea `yaml:"areas"`
	Decisions        []any                  `yaml:"decisions"`
	Standards        []any                  `yaml:"standards"`
	FitnessFunctions []any                  `yaml:"fitness_functions"`
}

type catalogArea struct {
	Contract string   `yaml:"contract"`
	Maturity string   `yaml:"maturity"`
	Evidence []string `yaml:"evidence"`
}

func Inspect(root string) (Inspection, error) {
	dir := filepath.Join(root, "engineering")
	contextData, err := os.ReadFile(filepath.Join(dir, "context.md"))
	if err != nil {
		return Inspection{}, err
	}
	var context contextDocument
	if err := yaml.Unmarshal([]byte(yamlPayload(string(contextData))), &context); err != nil {
		return Inspection{}, fmt.Errorf("engineering/context.md has invalid YAML metadata: %w", err)
	}
	catalogData, err := os.ReadFile(filepath.Join(dir, "engineering-system.yaml"))
	if err != nil {
		return Inspection{}, err
	}
	var catalog catalogDocument
	if err := yaml.Unmarshal(catalogData, &catalog); err != nil {
		return Inspection{}, fmt.Errorf("engineering-system.yaml is invalid YAML: %w", err)
	}
	i := Inspection{
		ID:               context.ID,
		Status:           context.Status,
		Version:          context.Version,
		OriginMode:       context.OriginMode,
		Scope:            catalog.Scope,
		Decisions:        len(catalog.Decisions),
		Standards:        len(catalog.Standards),
		FitnessFunctions: len(catalog.FitnessFunctions),
	}
	if catalog.SchemaVersion != 1 {
		i.Blockers = append(i.Blockers, "catalog schema_version must be 1")
	}
	if !regexp.MustCompile(`^ENGSYS-[A-Z0-9-]+$`).MatchString(i.ID) {
		i.Blockers = append(i.Blockers, "context engineering system id is invalid")
	}
	if !regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`).MatchString(i.Version) {
		i.Blockers = append(i.Blockers, "context semantic version is invalid")
	}
	if !oneOf(i.OriginMode, "generate", "evolve", "adopt") {
		i.Blockers = append(i.Blockers, "context origin mode is invalid")
	}
	if i.Scope == "" {
		i.Blockers = append(i.Blockers, "catalog scope is missing")
	}
	for field, values := range map[string][2]string{
		"id":          {context.ID, catalog.ID},
		"status":      {context.Status, catalog.Status},
		"version":     {context.Version, catalog.Version},
		"origin_mode": {context.OriginMode, catalog.OriginMode},
	} {
		if values[1] == "" || values[0] != values[1] {
			i.Blockers = append(i.Blockers, fmt.Sprintf("context and catalog %s do not match", field))
		}
	}
	if len(catalog.Areas) == 0 {
		i.Blockers = append(i.Blockers, "catalog areas are missing")
	}
	for name, source := range catalog.Areas {
		area := Area{Name: name, Contract: source.Contract, Maturity: source.Maturity, Evidence: len(source.Evidence)}
		i.Areas = append(i.Areas, area)
		if area.Contract == "" || area.Maturity == "" {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s is missing contract or maturity", name))
			continue
		}
		if !allowedMaturity[area.Maturity] {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s has invalid maturity %s", name, area.Maturity))
		}
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(area.Contract))); err != nil {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s contract %s is missing", name, area.Contract))
		}
		if area.Maturity != "baseline" && area.Evidence == 0 {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s maturity %s requires evidence", name, area.Maturity))
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "engineering-system.md")); err != nil {
		i.Blockers = append(i.Blockers, "engineering-system.md is missing")
	}
	sort.Slice(i.Areas, func(left, right int) bool { return i.Areas[left].Name < i.Areas[right].Name })
	i.Blockers = unique(i.Blockers)
	sort.Strings(i.Blockers)
	return i, nil
}

func Validate(root string) (Inspection, error) { return Inspect(root) }

func Migrate(root string, dryRun bool) ([]string, error) {
	path := filepath.Join(root, "engineering", "engineering-system.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("engineering-system.yaml is invalid YAML: %w", err)
	}
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("engineering-system.yaml must be a YAML mapping")
	}
	mapping := document.Content[0]
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == "schema_version" {
			if mapping.Content[index+1].Value != "1" {
				return nil, fmt.Errorf("unsupported schema_version %s", mapping.Content[index+1].Value)
			}
			return []string{"Engineering System catalog already uses schema_version 1"}, nil
		}
	}
	change := "ADD engineering/engineering-system.yaml schema_version: 1"
	if dryRun {
		return []string{change}, nil
	}
	key := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "schema_version"}
	value := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"}
	mapping.Content = append([]*yaml.Node{key, value}, mapping.Content...)
	updated, err := yaml.Marshal(&document)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return nil, err
	}
	return []string{change}, nil
}

func Triggers(text string) (valid, invalid []string) {
	var context contextDocument
	if err := yaml.Unmarshal([]byte(yamlPayload(text)), &context); err != nil {
		return nil, []string{"invalid_yaml"}
	}
	seen := map[string]bool{}
	for _, value := range context.EngineeringTriggers {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		if allowedTriggers[value] {
			valid = append(valid, value)
		} else {
			invalid = append(invalid, value)
		}
	}
	sort.Strings(valid)
	sort.Strings(invalid)
	return valid, invalid
}

func AllowedTriggers() []string {
	var out []string
	for trigger := range allowedTriggers {
		out = append(out, trigger)
	}
	sort.Strings(out)
	return out
}

func yamlPayload(text string) string {
	text = strings.TrimPrefix(text, "\ufeff")
	trimmed := strings.TrimSpace(text)
	if strings.HasPrefix(trimmed, "---") {
		lines := strings.Split(trimmed, "\n")
		for index := 1; index < len(lines); index++ {
			if strings.TrimSpace(lines[index]) == "---" {
				return strings.Join(lines[1:index], "\n")
			}
		}
	}
	lower := strings.ToLower(text)
	if start := strings.Index(lower, "```yaml"); start >= 0 {
		body := text[start+len("```yaml"):]
		if end := strings.Index(body, "```"); end >= 0 {
			return body[:end]
		}
	}
	return text
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
		if item != "" && !seen[item] {
			seen[item] = true
			out = append(out, item)
		}
	}
	return out
}
