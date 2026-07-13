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
	if _, err := os.Stat(filepath.Join(target, "product")); err == nil {
		return Result{}, fmt.Errorf("target already contains product/: %s", target)
	} else if !os.IsNotExist(err) {
		return Result{}, err
	}
	point, err := ParseStartingPoint(opts.StartingPoint)
	if err != nil {
		return Result{}, err
	}
	plan, err := buildInitializationPlan(point)
	if err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return Result{}, err
	}
	stagingRoot, err := os.MkdirTemp(target, ".spec-framework-init-")
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(stagingRoot)
	if err := writeInitializationPlan(filepath.Join(stagingRoot, "product"), plan); err != nil {
		return Result{}, err
	}
	result, err := Upgrade(Options{Target: stagingRoot, Version: opts.Version, Agents: opts.Agents, StartingPoint: point, Force: true})
	if err != nil {
		return Result{}, err
	}
	version := opts.Version
	if version == "" {
		version = "dev"
	}
	if err := writeStarterGuides(stagingRoot, version, opts.Agents, point); err != nil {
		return Result{}, err
	}
	result.StartingPoint = point
	for _, action := range plan.Actions {
		switch action {
		case "create-import-run":
			runID, err := sourceimport.CreateRun(filepath.Join(stagingRoot, "product"), opts.Sources)
			if err != nil {
				return Result{}, err
			}
			result.ImportID = runID
			if err := updateImportManifest(stagingRoot, runID); err != nil {
				return Result{}, err
			}
			if err := pinImportGuide(stagingRoot, runID); err != nil {
				return Result{}, err
			}
		}
	}
	if err := commitStagedProduct(filepath.Join(stagingRoot, "product"), filepath.Join(target, "product")); err != nil {
		return Result{}, err
	}
	result.Target = target
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
