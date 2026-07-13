package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	framework "github.com/JonatasFreireDev/spec-framework"
	"github.com/JonatasFreireDev/spec-framework/internal/dispatcher"
	"github.com/JonatasFreireDev/spec-framework/internal/runtimeassets"
	"github.com/JonatasFreireDev/spec-framework/internal/sourceimport"
)

type Agent string

const (
	Codex  Agent = "codex"
	Cursor Agent = "cursor"
	Claude Agent = "claude"
)

var agentRoots = map[Agent]string{Codex: "codex", Cursor: "cursor", Claude: "claude"}

type Options struct {
	Target, Version string
	Agents          []Agent
	StartingPoint   string
	Sources         []string
	Force           bool
}
type Result struct {
	Target, SpecRoot string
	SkillCount       int
	StartingPoint    string
	ImportID         string
}

var StartingPoints = []string{"new-product", "existing-product", "existing-documents", "existing-feature", "existing-implementation", "audit-only"}

func ParseStartingPoint(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "new-product", nil
	}
	for _, candidate := range StartingPoints {
		if value == candidate {
			return value, nil
		}
	}
	return "", fmt.Errorf("unsupported starting point %q", value)
}

func ParseAgents(value string) ([]Agent, error) {
	seen := map[Agent]bool{}
	var out []Agent
	for _, raw := range strings.Split(value, ",") {
		a := Agent(strings.ToLower(strings.TrimSpace(raw)))
		if a == "" {
			continue
		}
		if _, ok := agentRoots[a]; !ok {
			return nil, fmt.Errorf("unsupported agent %q", raw)
		}
		if !seen[a] {
			seen[a] = true
			out = append(out, a)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("select at least one agent: codex, cursor, or claude")
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out, nil
}

func InstalledAgents(target string) ([]Agent, error) {
	data, err := os.ReadFile(filepath.Join(target, "product", ".product", "framework.json"))
	if err != nil {
		return nil, fmt.Errorf("read installed agents: %w", err)
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("parse installed manifest: %w", err)
	}
	raw, _ := value["agents"].([]any)
	var names []string
	for _, item := range raw {
		if name, ok := item.(string); ok {
			names = append(names, name)
		}
	}
	return ParseAgents(strings.Join(names, ","))
}

func Init(opts Options) (Result, error) {
	target, err := filepath.Abs(opts.Target)
	if err != nil {
		return Result{}, err
	}
	if _, err := os.Stat(filepath.Join(target, "product")); err == nil && !opts.Force {
		return Result{}, fmt.Errorf("target already contains product/: %s", target)
	}
	if err := copyTree("starter/product", filepath.Join(target, "product")); err != nil {
		return Result{}, err
	}
	point, err := ParseStartingPoint(opts.StartingPoint)
	if err != nil {
		return Result{}, err
	}
	result, err := Upgrade(Options{Target: target, Version: opts.Version, Agents: opts.Agents, StartingPoint: point, Force: true})
	if err != nil {
		return Result{}, err
	}
	if err := configureStartingPointProduct(target, point); err != nil {
		return Result{}, err
	}
	version := opts.Version
	if version == "" {
		version = "dev"
	}
	if err := writeStarterGuides(target, version, opts.Agents, point); err != nil {
		return Result{}, err
	}
	result.StartingPoint = point
	if point == "existing-documents" {
		runID, err := sourceimport.CreateRun(filepath.Join(target, "product"), opts.Sources)
		if err != nil {
			return Result{}, err
		}
		result.ImportID = runID
		if err := updateImportManifest(target, runID); err != nil {
			return Result{}, err
		}
		if err := pinImportGuide(target, runID); err != nil {
			return Result{}, err
		}
	}
	return result, nil
}

func pinImportGuide(target, runID string) error {
	path := filepath.Join(target, "product", "BOOTSTRAP.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	updated := strings.NewReplacer("<IMPORT-NNN>", runID, "<latest-run>", runID).Replace(string(data))
	return writeFile(path, []byte(updated), 0644)
}

func configureStartingPointProduct(target, point string) error {
	switch point {
	case "existing-feature":
		return configureExistingFeatureProduct(target)
	case "existing-product":
		return configureExistingProduct(target)
	case "existing-implementation":
		return configureExistingImplementationProduct(target)
	default:
		return nil
	}
}

func configureExistingProduct(target string) error {
	productRoot := filepath.Join(target, "product")
	template, err := fs.ReadFile(framework.Assets, "framework/template/product-baseline-template.md")
	if err != nil {
		return err
	}
	baseline := strings.NewReplacer(
		"[PRODUCT-BASELINE-XXX]", "PRODUCT-BASELINE-TBD",
		"[draft | proposed | approved]", "draft",
	).Replace(string(template))
	if err := writeFile(filepath.Join(productRoot, "foundation", "product-baseline.md"), []byte(baseline), 0644); err != nil {
		return err
	}
	guidance := map[string][2]string{
		filepath.Join(productRoot, "context.md"): {
			"Do not create domains or features until `foundation/problem/problem.md`, `foundation/vision/vision.md`, and `foundation/strategy/strategy.md` contain product-specific content.",
			"For `existing-product`, approve `foundation/product-baseline.md` from repository and operational evidence, then approve `foundation/strategy/strategy.md` before creating delivery work.",
		},
		filepath.Join(productRoot, "foundation", "README.md"): {
			"Fill `problem/problem.md`, then `vision/vision.md`, then `strategy/strategy.md`.",
			"For `existing-product`, complete and individually approve `product-baseline.md`, then `strategy/strategy.md`. Escalate to the full Foundation path when audience, delivered value, or product direction is uncertain.",
		},
	}
	for path, replacement := range guidance {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		updated := strings.Replace(string(data), replacement[0], replacement[1], 1)
		if updated == string(data) {
			return fmt.Errorf("existing-product guidance anchor not found: %s", path)
		}
		if err := writeFile(path, []byte(updated), 0644); err != nil {
			return err
		}
	}
	registryPath := filepath.Join(productRoot, ".product", "artifacts.json")
	var registry struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		return err
	}
	consolidated := map[string]bool{"problem": true, "vision": true, "product-principles": true, "north-star": true}
	filtered := make([]map[string]any, 0, len(registry.Artifacts))
	for _, artifact := range registry.Artifacts {
		kind, _ := artifact["type"].(string)
		if consolidated[kind] {
			continue
		}
		if kind == "strategy" {
			artifact["parentIds"] = []string{"PRODUCT-BASELINE-TBD"}
		}
		filtered = append(filtered, artifact)
	}
	registry.Artifacts = append([]map[string]any{{
		"id": "PRODUCT-BASELINE-TBD", "type": "product-baseline", "status": "draft",
		"path": "foundation/product-baseline.md", "parentIds": []string{},
	}}, filtered...)
	return writeJSON(registryPath, registry)
}

func configureExistingFeatureProduct(target string) error {
	productRoot := filepath.Join(target, "product")
	briefPath := filepath.Join(productRoot, "foundation", "feature-brief.md")
	brief, err := fs.ReadFile(framework.Assets, "framework/template/feature-brief-template.md")
	if err != nil {
		return err
	}
	brief = []byte(strings.NewReplacer(
		"# Feature Brief: [feature name]", "# Feature Brief",
		"[FBR-XXX]", "FEATURE-BRIEF-TBD",
		"[draft | proposed | approved]", "draft",
		"[FT-XXX or product-relative feature path]", "FT-TEMPLATE",
	).Replace(string(brief)))
	if err := writeFile(briefPath, brief, 0644); err != nil {
		return err
	}
	guidance := map[string][2]string{
		filepath.Join(productRoot, "context.md"): {
			"Do not create domains or features until `foundation/problem/problem.md`, `foundation/vision/vision.md`, and `foundation/strategy/strategy.md` contain product-specific content.",
			"For `existing-feature`, complete and approve `foundation/feature-brief.md` before creating a workspace. Escalate to full Foundation when the work requires product-wide direction.",
		},
		filepath.Join(productRoot, "foundation", "README.md"): {
			"Fill `problem/problem.md`, then `vision/vision.md`, then `strategy/strategy.md`.",
			"For `existing-feature`, complete and individually approve `feature-brief.md`. Use the full Problem -> Vision -> Strategy path only when the scope is broad or uncertain.",
		},
	}
	for path, replacement := range guidance {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		updated := strings.Replace(string(data), replacement[0], replacement[1], 1)
		if updated == string(data) {
			return fmt.Errorf("existing-feature guidance anchor not found: %s", path)
		}
		if err := writeFile(path, []byte(updated), 0644); err != nil {
			return err
		}
	}
	registryPath := filepath.Join(productRoot, ".product", "artifacts.json")
	var registry struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		return err
	}
	fullFoundation := map[string]bool{"problem": true, "vision": true, "product-principles": true, "north-star": true, "strategy": true}
	filtered := make([]map[string]any, 0, len(registry.Artifacts)+1)
	for _, artifact := range registry.Artifacts {
		kind, _ := artifact["type"].(string)
		if fullFoundation[kind] {
			continue
		}
		if kind == "feature" {
			parents, _ := artifact["parentIds"].([]any)
			parents = append(parents, "FEATURE-BRIEF-TBD")
			artifact["parentIds"] = parents
		}
		filtered = append(filtered, artifact)
	}
	registry.Artifacts = append([]map[string]any{{
		"id": "FEATURE-BRIEF-TBD", "type": "feature-brief", "status": "draft",
		"path": "foundation/feature-brief.md", "parentIds": []string{}, "targetFeature": "FT-TEMPLATE",
	}}, filtered...)
	return writeJSON(registryPath, registry)
}

func configureExistingImplementationProduct(target string) error {
	productRoot := filepath.Join(target, "product")
	template, err := fs.ReadFile(framework.Assets, "framework/template/implementation-assessment-template.md")
	if err != nil {
		return err
	}
	assessment := strings.NewReplacer(
		"[IMPL-ASSESS-XXX]", "IMPL-ASSESS-TBD",
		"[draft | proposed | approved]", "draft",
	).Replace(string(template))
	assessmentPath := filepath.Join(productRoot, "knowledge", "assessments", "implementation-assessment.md")
	if err := writeFile(assessmentPath, []byte(assessment), 0644); err != nil {
		return err
	}
	contextPath := filepath.Join(productRoot, "context.md")
	contextData, err := os.ReadFile(contextPath)
	if err != nil {
		return err
	}
	oldRule := "Do not create domains or features until `foundation/problem/problem.md`, `foundation/vision/vision.md`, and `foundation/strategy/strategy.md` contain product-specific content."
	newRule := "For `existing-implementation`, approve `knowledge/assessments/implementation-assessment.md`, then complete the full Foundation path before creating delivery work. Observed code is evidence, not approved product intent."
	updated := strings.Replace(string(contextData), oldRule, newRule, 1)
	if updated == string(contextData) {
		return fmt.Errorf("existing-implementation guidance anchor not found: %s", contextPath)
	}
	if err := writeFile(contextPath, []byte(updated), 0644); err != nil {
		return err
	}
	foundationReadme := filepath.Join(productRoot, "foundation", "README.md")
	readmeData, err := os.ReadFile(foundationReadme)
	if err != nil {
		return err
	}
	oldNext := "Fill `problem/problem.md`, then `vision/vision.md`, then `strategy/strategy.md`."
	newNext := "For `existing-implementation`, first approve `../knowledge/assessments/implementation-assessment.md`. Then fill `problem/problem.md`, `vision/vision.md`, and `strategy/strategy.md` from reviewed evidence without treating observed behavior as intent."
	updatedReadme := strings.Replace(string(readmeData), oldNext, newNext, 1)
	if updatedReadme == string(readmeData) {
		return fmt.Errorf("existing-implementation Foundation guidance anchor not found: %s", foundationReadme)
	}
	if err := writeFile(foundationReadme, []byte(updatedReadme), 0644); err != nil {
		return err
	}
	registryPath := filepath.Join(productRoot, ".product", "artifacts.json")
	var registry struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		return err
	}
	for _, artifact := range registry.Artifacts {
		if artifact["type"] == "problem" {
			parents, _ := artifact["parentIds"].([]any)
			artifact["parentIds"] = append(parents, "IMPL-ASSESS-TBD")
		}
	}
	registry.Artifacts = append([]map[string]any{{
		"id": "IMPL-ASSESS-TBD", "type": "implementation-assessment", "status": "draft",
		"path": "knowledge/assessments/implementation-assessment.md", "parentIds": []string{},
	}}, registry.Artifacts...)
	return writeJSON(registryPath, registry)
}

func Upgrade(opts Options) (Result, error) {
	target, err := filepath.Abs(opts.Target)
	if err != nil {
		return Result{}, err
	}
	if _, err := os.Stat(filepath.Join(target, "product")); err != nil {
		return Result{}, fmt.Errorf("target does not contain product/: %s", target)
	}
	productManifest := filepath.Join(target, "product", ".product", "framework.json")
	if strings.TrimSpace(opts.StartingPoint) == "" {
		opts.StartingPoint = installedStartingPoint(productManifest)
	}
	version := opts.Version
	if version == "" {
		version = "dev"
	}
	agents := make([]string, len(opts.Agents))
	for i, a := range opts.Agents {
		agents[i] = string(a)
	}
	point, err := ParseStartingPoint(opts.StartingPoint)
	if err != nil {
		return Result{}, err
	}
	manifest := map[string]any{
		"schema_version": 3, "framework": "spec-framework", "version": version, "product_root": ".",
		"agents": agents, "starting_point": point,
		"activation": map[string]any{"mode": "manifest-only"},
		"runtime":    map[string]any{"source": "user-cache", "channel": "stable"},
	}
	if err := writeJSON(productManifest, manifest); err != nil {
		return Result{}, err
	}
	spec, err := runtimeassets.Ensure(version)
	if err != nil {
		return Result{}, err
	}
	for _, agent := range opts.Agents {
		if _, err := dispatcher.Install(string(agent)); err != nil {
			return Result{}, fmt.Errorf("install %s dispatcher: %w", agent, err)
		}
	}
	return Result{Target: target, SpecRoot: spec, SkillCount: 0, StartingPoint: point}, nil
}

func installedStartingPoint(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "new-product"
	}
	var value map[string]any
	if json.Unmarshal(data, &value) != nil {
		return "new-product"
	}
	point, _ := value["starting_point"].(string)
	if point == "" {
		return "new-product"
	}
	return point
}

func updateImportManifest(target, runID string) error {
	for _, path := range []string{filepath.Join(target, "product", ".product", "framework.json")} {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var value map[string]any
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		value["import"] = map[string]any{"latest_run": runID, "sources_path": "knowledge/imports/sources"}
		if err := writeJSON(path, value); err != nil {
			return err
		}
	}
	return nil
}

func copyTree(source, target string) error {
	return fs.WalkDir(framework.Assets, source, func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(source, name)
		dest := target
		if rel != "." {
			dest = filepath.Join(target, rel)
		}
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		data, err := framework.Assets.ReadFile(name)
		if err != nil {
			return err
		}
		return writeFile(dest, data, 0644)
	})
}
func writeFile(name string, data []byte, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	tmp := name + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, name)
}
func writeJSON(name string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return writeFile(name, append(data, '\n'), 0644)
}
