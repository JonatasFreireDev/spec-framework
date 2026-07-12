package engineeringsystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

func Inspect(root string) (Inspection, error) {
	dir := filepath.Join(root, "engineering")
	context, err := os.ReadFile(filepath.Join(dir, "context.md"))
	if err != nil {
		return Inspection{}, err
	}
	catalogPath := filepath.Join(dir, "engineering-system.yaml")
	catalog, err := os.ReadFile(catalogPath)
	if err != nil {
		return Inspection{}, err
	}
	i := Inspection{
		ID:         field(string(context), "id"),
		Status:     field(string(context), "status"),
		Version:    field(string(context), "version"),
		OriginMode: field(string(context), "origin_mode"),
		Scope:      field(string(catalog), "scope"),
	}
	i.Areas, i.Decisions, i.Standards, i.FitnessFunctions, i.Blockers = parseCatalog(string(catalog), dir)
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
	if _, err := os.Stat(filepath.Join(dir, "engineering-system.md")); err != nil {
		i.Blockers = append(i.Blockers, "engineering-system.md is missing")
	}
	sort.Slice(i.Areas, func(left, right int) bool { return i.Areas[left].Name < i.Areas[right].Name })
	i.Blockers = unique(i.Blockers)
	sort.Strings(i.Blockers)
	return i, nil
}

func Validate(root string) (Inspection, error) { return Inspect(root) }

func Triggers(text string) (valid, invalid []string) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	in := false
	indent := -1
	seen := map[string]bool{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		currentIndent := len(line) - len(strings.TrimLeft(line, " \t"))
		if strings.HasPrefix(strings.ToLower(trimmed), "engineering_triggers:") {
			in = true
			indent = currentIndent
			continue
		}
		if !in {
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if currentIndent <= indent || !strings.HasPrefix(trimmed, "-") {
			break
		}
		value := strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "-")), "`\"'")
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

func parseCatalog(text, root string) ([]Area, int, int, int, []string) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	section := ""
	areaName := ""
	area := Area{}
	var areas []Area
	counts := map[string]int{}
	var blockers []string
	flush := func() {
		if areaName == "" {
			return
		}
		area.Name = areaName
		areas = append(areas, area)
		areaName = ""
		area = Area{}
	}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent == 0 && strings.HasSuffix(trimmed, ":") {
			flush()
			section = strings.TrimSuffix(trimmed, ":")
			continue
		}
		if section == "areas" && indent == 2 && strings.HasSuffix(trimmed, ":") {
			flush()
			areaName = strings.TrimSuffix(trimmed, ":")
			continue
		}
		if section == "areas" && areaName != "" && indent >= 4 {
			key, value := pair(trimmed)
			switch key {
			case "contract":
				area.Contract = value
			case "maturity":
				area.Maturity = value
			case "evidence":
				if value != "[]" && value != "" {
					area.Evidence++
				}
			}
			continue
		}
		if (section == "decisions" || section == "standards" || section == "fitness_functions") && strings.HasPrefix(trimmed, "-") {
			counts[section]++
		}
	}
	flush()
	if len(areas) == 0 {
		blockers = append(blockers, "catalog areas are missing")
	}
	for _, item := range areas {
		if item.Contract == "" || item.Maturity == "" {
			blockers = append(blockers, fmt.Sprintf("area %s is missing contract or maturity", item.Name))
			continue
		}
		if !allowedMaturity[item.Maturity] {
			blockers = append(blockers, fmt.Sprintf("area %s has invalid maturity %s", item.Name, item.Maturity))
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(item.Contract))); err != nil {
			blockers = append(blockers, fmt.Sprintf("area %s contract %s is missing", item.Name, item.Contract))
		}
		if item.Maturity != "baseline" && item.Evidence == 0 {
			blockers = append(blockers, fmt.Sprintf("area %s maturity %s requires evidence", item.Name, item.Maturity))
		}
	}
	return areas, counts["decisions"], counts["standards"], counts["fitness_functions"], blockers
}

func pair(line string) (string, string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), "`\"'")
}

func field(text, name string) string {
	re := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(name) + `:\s*([^\r\n#]+)`)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return ""
	}
	return strings.Trim(strings.TrimSpace(match[1]), "`\"'")
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

var ErrNotConfigured = errors.New("Engineering System is not configured")
