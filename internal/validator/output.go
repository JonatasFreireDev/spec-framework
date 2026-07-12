package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func WriteReport(root string, result Result) ([]string, error) {
	date := time.Now().UTC().Format("2006-01-02")
	var b strings.Builder
	fmt.Fprintf(&b, "# Framework Validation Report\n\n## Executive Snapshot\n\n| Field | Value |\n| --- | --- |\n| Date | %s |\n| Validator | `spec-framework validate` |\n| Verdict | %s |\n| Errors | %d |\n| Warnings | %d |\n| Notes | %d |\n\n## Findings\n\n| Severity | Check | File | Finding | Suggested Fix |\n| --- | --- | --- | --- | --- |\n", date, result.Verdict(), result.Errors, result.Warnings, result.Notes)
	if len(result.Diagnostics) == 0 {
		b.WriteString("| ✅ ready | framework | repository | No findings. | None |\n")
	} else {
		for _, d := range result.Diagnostics {
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", escape(string(d.Severity)), escape(d.Check), escape(d.File), escape(d.Message), escape(d.Fix))
		}
	}
	fmt.Fprintf(&b, "\n## Result\n\n| Field | Value |\n| --- | --- |\n| Verdict | %s |\n", result.Verdict())
	report := filepath.Join(root, "audits", "framework-validation-report.md")
	if err := atomicWrite(report, []byte(b.String())); err != nil {
		return nil, err
	}
	readiness := filepath.Join(root, "audits", "readiness", "framework-readiness.md")
	summary := fmt.Sprintf("# Framework Readiness\n\n| Field | Value |\n| --- | --- |\n| Verdict | %s |\n| Errors | %d |\n| Warnings | %d |\n", result.Verdict(), result.Errors, result.Warnings)
	if err := atomicWrite(readiness, []byte(summary)); err != nil {
		return nil, err
	}
	return []string{report, readiness}, nil
}

func WriteRegistry(root string) (string, error) {
	snap, err := Scan(context.Background(), root, root, 0)
	if err != nil {
		return "", err
	}
	artifacts := buildRegistry(snap)
	value := map[string]any{"generatedAt": time.Now().UTC().Format(time.RFC3339Nano), "generator": "spec-framework (Go)", "artifacts": artifacts}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	path := filepath.Join(root, ".product", "artifacts.json")
	if err := atomicWrite(path, append(data, '\n')); err != nil {
		return "", err
	}
	return path, nil
}

func buildRegistry(s Snapshot) []map[string]any {
	byPath := map[string]map[string]any{}
	useCaseIDs := map[string]string{}
	for path, text := range s.Text {
		if strings.HasSuffix(path, "/context.md") && strings.Contains(path, "/use-cases/") {
			if id := metadata(text)["id"]; id != "" {
				useCaseIDs[filepath.ToSlash(filepath.Dir(path))] = id
			}
		}
	}
	for path, text := range s.Text {
		if !strings.HasSuffix(path, ".md") {
			continue
		}
		if filepath.Base(path) == "context.md" {
			continue
		}
		meta := metadata(text)
		fields := tableFields(text)
		companionPath := path[:strings.LastIndex(path, "/")+1] + "context.md"
		companion := metadata(s.Text[companionPath])
		id := first(meta["id"], fields["id"])
		kind := first(meta["type"], fields["type"], inferType(path))
		status := first(meta["status"], fields["status"])
		if id == "" || kind == "" || status == "" {
			continue
		}
		name := first(meta["name"], fields["name"], strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		owner := first(meta["owner_skill"], fields["owner skill"], fields["owner"])
		parents := []string{}
		dir := filepath.ToSlash(filepath.Dir(path))
		for base, parent := range useCaseIDs {
			if strings.HasPrefix(dir, base) && id != parent {
				parents = []string{parent}
				break
			}
		}
		level := first(fields["delivery level"], fields["level"], meta["level"], companion["level"])
		priority := first(fields["priority"], meta["priority"], companion["priority"])
		rationale := first(fields["rationale"], meta["rationale"], companion["rationale"])
		byPath[path] = map[string]any{"id": id, "type": strings.ReplaceAll(kind, "_", "-"), "name": name, "status": status, "ownerSkill": owner, "path": path, "parentIds": parents, "childIds": []string{}, "dependsOn": []string{}, "decisions": []string{}, "delivery": map[string]any{"level": level, "priority": priority, "depends_on": []string{}, "rationale": rationale}, "documents": map[string]string{"canonical": path}}
	}
	for path, value := range s.JSON {
		if !strings.HasSuffix(path, "execution-graph.json") {
			continue
		}
		object, ok := value.(map[string]any)
		if !ok {
			continue
		}
		id, _ := object["id"].(string)
		if id == "" {
			continue
		}
		dir := filepath.ToSlash(filepath.Dir(path))
		parents := []string{}
		if parent := useCaseIDs[dir]; parent != "" {
			parents = []string{parent}
		}
		delivery, _ := object["delivery"].(map[string]any)
		byPath[path] = map[string]any{"id": id, "type": "execution-graph", "name": id, "status": firstString(object["status"], "draft"), "path": path, "parentIds": parents, "childIds": []string{}, "dependsOn": []string{}, "decisions": []string{}, "delivery": delivery, "documents": map[string]string{"canonical": path}}
	}
	var out []map[string]any
	for _, a := range byPath {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		left, _ := out[i]["id"].(string)
		right, _ := out[j]["id"].(string)
		if left == right {
			return out[i]["path"].(string) < out[j]["path"].(string)
		}
		return left < right
	})
	return out
}
func inferType(path string) string {
	base := filepath.Base(path)
	if strings.Contains(filepath.ToSlash(path), "/tasks/") {
		return "task"
	}
	return map[string]string{"domain.md": "domain", "goal.md": "goal", "feature.md": "feature", "use-case.md": "use-case", "specification.md": "specification", "design.md": "design", "technical-discovery.md": "technical-discovery", "implementation-plan.md": "implementation-plan", "tasks.md": "taskset", "tests.md": "tests", "analytics.md": "analytics", "audit.md": "audit", "qa-evidence.md": "qa-evidence", "security-review.md": "security-review", "code-review.md": "code-review"}[base]
}
func first(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
func firstString(value any, fallback string) string {
	if s, ok := value.(string); ok && s != "" {
		return s
	}
	return fallback
}
func escape(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "|", "\\|"), "\n", " ")
}
func atomicWrite(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
