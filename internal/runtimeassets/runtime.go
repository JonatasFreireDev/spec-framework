package runtimeassets

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	framework "github.com/JonatasFreireDev/spec-framework"
)

type Manifest struct {
	SchemaVersion  int        `json:"schema_version"`
	Framework      string     `json:"framework"`
	Version        string     `json:"version"`
	ProductRoot    string     `json:"product_root"`
	StartingPoint  string     `json:"starting_point"`
	Agents         []string   `json:"agents"`
	CodeRoots      []CodeRoot `json:"code_roots,omitempty"`
	BaselinePolicy struct {
		PreSpecification string `json:"pre_specification,omitempty"`
	} `json:"baseline_policy,omitempty"`
	Activation struct {
		Mode string `json:"mode"`
	} `json:"activation"`
	Runtime struct {
		Source  string `json:"source"`
		Channel string `json:"channel"`
	} `json:"runtime"`
}

// CodeRoot is an implementation area located alongside product/. It records
// the semantic role rather than coupling the framework to one language.
type CodeRoot struct {
	Path string `json:"path"`
	Role string `json:"role"`
}

func CacheRoot() (string, error) {
	if override := strings.TrimSpace(os.Getenv("SPEC_FRAMEWORK_CACHE")); override != "" {
		return filepath.Abs(override)
	}
	root, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "spec-framework"), nil
}

func VersionRoot(version string) (string, error) {
	root, err := CacheRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "versions", safeVersion(version)), nil
}

func Ensure(version string) (string, error) {
	root, err := VersionRoot(version)
	if err != nil {
		return "", err
	}
	marker := filepath.Join(root, ".complete")
	if _, err := os.Stat(marker); err == nil {
		return root, nil
	}
	versions := filepath.Dir(root)
	if err := os.MkdirAll(versions, 0755); err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp(versions, ".install-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	assets := map[string]string{
		"FRAMEWORK.md":                      "FRAMEWORK.md",
		"docs/artifact-registry-modules.md": "docs/artifact-registry-modules.md",
		"docs/execution-runtime.md":         "docs/execution-runtime.md",
		"docs/engineering-systems.md":       "docs/engineering-systems.md",
		"docs/lifecycle-and-approvals.md":   "docs/lifecycle-and-approvals.md",
		"framework/AGENTS.framework.md":     "AGENTS.framework.md",
		"framework/delivery-closure.md":     "delivery-closure.md",
		"framework/init":                    "init",
		"framework/skills":                  "skills",
		"examples/events":                   "examples/events",
	}
	for source, dest := range assets {
		if err := copyTree(source, filepath.Join(tmp, dest)); err != nil {
			return "", err
		}
	}
	if err := os.WriteFile(filepath.Join(tmp, ".complete"), []byte(version+"\n"), 0644); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, root); err != nil {
		if _, statErr := os.Stat(marker); statErr == nil {
			return root, nil
		}
		return "", err
	}
	return root, nil
}

func Discover(start string) (string, Manifest, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", Manifest{}, err
	}
	if info, err := os.Stat(current); err == nil && !info.IsDir() {
		current = filepath.Dir(current)
	}
	for {
		path := filepath.Join(current, "product", ".product", "framework.json")
		if data, err := os.ReadFile(path); err == nil {
			var manifest Manifest
			if err := json.Unmarshal(data, &manifest); err != nil {
				return "", Manifest{}, fmt.Errorf("invalid Spec Framework manifest %s: %w", path, err)
			}
			if manifest.Framework != "spec-framework" || manifest.Activation.Mode != "manifest-only" || strings.TrimSpace(manifest.Version) == "" {
				return "", Manifest{}, fmt.Errorf("manifest %s is not an active Spec Framework product", path)
			}
			return current, manifest, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", Manifest{}, fmt.Errorf("no Spec Framework product found; expected product/.product/framework.json")
}

func Resolve(start string) (repoRoot, productRoot, frameworkRoot string, manifest Manifest, err error) {
	repoRoot, manifest, err = Discover(start)
	if err != nil {
		return
	}
	productRoot = filepath.Join(repoRoot, "product")
	frameworkRoot, err = Ensure(manifest.Version)
	return
}

func safeVersion(version string) string {
	version = strings.TrimSpace(strings.TrimPrefix(version, "v"))
	return strings.NewReplacer("/", "_", "\\", "_", ":", "_").Replace(version)
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
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		return os.WriteFile(dest, data, 0644)
	})
}
