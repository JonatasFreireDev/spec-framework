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
)

type Agent string

const (
	Codex  Agent = "codex"
	Cursor Agent = "cursor"
	Claude Agent = "claude"
)

var agentRoots = map[Agent]string{Codex: ".agents/skills", Cursor: ".cursor/skills", Claude: ".claude/skills"}

type Options struct {
	Target, Version string
	Agents          []Agent
	Force           bool
}
type Result struct {
	Target, SpecRoot string
	SkillCount       int
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

func Init(opts Options) (Result, error) {
	target, err := filepath.Abs(opts.Target)
	if err != nil {
		return Result{}, err
	}
	if entries, err := os.ReadDir(target); err == nil && len(entries) > 0 && !opts.Force {
		return Result{}, fmt.Errorf("target is not empty: %s", target)
	}
	if err := copyTree("starter", target); err != nil {
		return Result{}, err
	}
	result, err := Upgrade(Options{Target: target, Version: opts.Version, Agents: opts.Agents, Force: true})
	if err != nil {
		return Result{}, err
	}
	version := opts.Version
	if version == "" {
		version = "dev"
	}
	if err := writeStarterGuides(target, version, opts.Agents); err != nil {
		return Result{}, err
	}
	return result, nil
}

func Upgrade(opts Options) (Result, error) {
	target, err := filepath.Abs(opts.Target)
	if err != nil {
		return Result{}, err
	}
	if _, err := os.Stat(filepath.Join(target, "product")); err != nil {
		return Result{}, fmt.Errorf("target does not contain product/: %s", target)
	}
	spec := filepath.Join(target, ".spec-framework")
	for source, dest := range map[string]string{
		"FRAMEWORK.md": "FRAMEWORK.md", "framework/AGENTS.framework.md": "AGENTS.framework.md",
		"framework/decisions": "decisions", "framework/skills": "skills", "framework/template": "templates",
	} {
		if err := copyTree(source, filepath.Join(spec, dest)); err != nil {
			return Result{}, err
		}
	}
	count := 0
	for _, agent := range opts.Agents {
		root := filepath.Join(target, filepath.FromSlash(agentRoots[agent]))
		if err := os.RemoveAll(root); err != nil {
			return Result{}, err
		}
		if err := copyTree("framework/skills", root); err != nil {
			return Result{}, err
		}
		if agent != Codex {
			entries, _ := os.ReadDir(root)
			for _, skill := range entries {
				_ = os.RemoveAll(filepath.Join(root, skill.Name(), "agents"))
			}
		}
		count++
	}
	version := opts.Version
	if version == "" {
		version = "dev"
	}
	agents := make([]string, len(opts.Agents))
	for i, a := range opts.Agents {
		agents[i] = string(a)
	}
	manifest := map[string]any{"schema_version": 1, "version": version, "product_root": "product", "agents": agents, "installed_assets": map[string]bool{"framework_document": true, "framework_agent": true, "decisions": true, "skills": true, "templates": true}}
	if err := writeJSON(filepath.Join(spec, "manifest.json"), manifest); err != nil {
		return Result{}, err
	}
	productManifest := filepath.Join(target, "product", ".product", "framework.json")
	if data, err := os.ReadFile(productManifest); err == nil {
		var value map[string]any
		if json.Unmarshal(data, &value) == nil {
			value["version"] = version
			value["framework_assets_path"] = "../.spec-framework"
			value["product_root"] = "."
			value["agents"] = agents
			value["installed_assets"] = manifest["installed_assets"]
			_ = writeJSON(productManifest, value)
		}
	}
	workflow := productWorkflow(version)
	if err := writeFile(filepath.Join(target, ".github", "workflows", "framework-validation.yml"), []byte(workflow), 0644); err != nil {
		return Result{}, err
	}
	return Result{Target: target, SpecRoot: spec, SkillCount: count}, nil
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
