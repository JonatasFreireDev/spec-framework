package install

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/runtimeassets"
)

var codeMarkers = map[string]string{
	"package.json":     "web",
	"go.mod":           "service",
	"Cargo.toml":       "service",
	"pom.xml":          "service",
	"build.gradle":     "service",
	"build.gradle.kts": "service",
}

// discoverCodeRoots inventories immediate repository siblings. Explicit roots
// supplement discovery and let an adopter choose a semantic role such as api,
// web, worker, mobile, infrastructure, or library.
func discoverCodeRoots(target string, explicit []runtimeassets.CodeRoot) ([]runtimeassets.CodeRoot, error) {
	byPath := map[string]runtimeassets.CodeRoot{}
	for _, root := range explicit {
		path := filepath.ToSlash(strings.Trim(strings.TrimSpace(root.Path), "/"))
		role := strings.ToLower(strings.TrimSpace(root.Role))
		if path == "" || role == "" || path == "product" || strings.HasPrefix(path, "../") || filepath.IsAbs(path) {
			return nil, fmt.Errorf("invalid code root %q; use a repository-relative sibling path with a role", root.Path)
		}
		byPath[path] = runtimeassets.CodeRoot{Path: path, Role: role}
	}
	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || entry.Name() == "product" || entry.Name() == "node_modules" || entry.Name() == "vendor" {
			continue
		}
		candidate := filepath.Join(target, entry.Name())
		role := detectCodeRole(candidate)
		if role != "" {
			path := filepath.ToSlash(entry.Name())
			if _, exists := byPath[path]; !exists {
				byPath[path] = runtimeassets.CodeRoot{Path: path, Role: role}
			}
		}
	}
	roots := make([]runtimeassets.CodeRoot, 0, len(byPath))
	for _, root := range byPath {
		roots = append(roots, root)
	}
	sort.Slice(roots, func(i, j int) bool { return roots[i].Path < roots[j].Path })
	return roots, nil
}

func detectCodeRole(root string) string {
	for marker, role := range codeMarkers {
		if _, err := os.Stat(filepath.Join(root, marker)); err == nil {
			return role
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "*.sln")); len(matches) > 0 {
		return "service"
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "*.csproj")); len(matches) > 0 {
		return "service"
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "*.py")); len(matches) > 0 {
		return "service"
	}
	return ""
}

func writeCodeDiscovery(productRoot string, roots []runtimeassets.CodeRoot) error {
	path := filepath.Join(productRoot, "knowledge", "assessments", "product-landscape.md")
	if _, err := os.Stat(path); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("# Product Landscape\n\n## Snapshot\n\n| Field | Value |\n| --- | --- |\n| Status | `draft` |\n| Evidence mode | `")
	if len(roots) == 0 {
		b.WriteString("hypothesis")
	} else {
		b.WriteString("observed-code")
	}
	b.WriteString("` |\n\n## Code Roots\n\n")
	if len(roots) == 0 {
		b.WriteString("No implementation root was detected. Define the intended product surface, language, framework, and official scaffold command before creating code.\n")
	} else {
		b.WriteString("| Path | Role | Discovery status |\n| --- | --- | --- |\n")
		for _, root := range roots {
			fmt.Fprintf(&b, "| `%s/` | `%s` | `pending comprehensive inventory` |\n", root.Path, root.Role)
		}
	}
	b.WriteString("\n## Required Inventory\n\nBefore domain modeling, inspect every declared code root for modules, routes or UI surfaces, data models, integrations, business rules, tests, configuration, design assets, and operational constraints. Record evidence paths and unresolved ownership boundaries; do not reduce a broad codebase to a single feature slice.\n\n## Domain Map\n\nPending comprehensive inventory. List all candidate domains, ownership boundaries, cross-domain workflows, and uncovered areas before creating a delivery Domain, Goal, Feature, or Use Case.\n")
	return writeFile(path, []byte(b.String()), 0644)
}
