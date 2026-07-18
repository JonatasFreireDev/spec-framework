package install

import (
	"fmt"
	"os"
	pathpkg "path"
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

const (
	CodeRootDiscoveryAgentDeclared      = "agent-declared"
	CodeRootDiscoveryAgentConfirmedNone = "agent-confirmed-none"
	CodeRootDiscoveryCLIFallback        = "cli-fallback"
)

// resolveCodeRoots makes agent declarations authoritative. The CLI heuristic
// remains a compatible fallback, but its result is explicitly unconfirmed.
func resolveCodeRoots(target string, explicit []runtimeassets.CodeRoot, mode string) ([]runtimeassets.CodeRoot, runtimeassets.CodeRootDiscovery, error) {
	mode = strings.TrimSpace(mode)
	switch mode {
	case CodeRootDiscoveryAgentDeclared:
		roots, err := normalizeCodeRoots(explicit)
		if err != nil {
			return nil, runtimeassets.CodeRootDiscovery{}, err
		}
		if len(roots) == 0 {
			return nil, runtimeassets.CodeRootDiscovery{}, fmt.Errorf("agent-declared code-root discovery requires at least one --code-roots entry")
		}
		for _, root := range roots {
			info, statErr := os.Stat(filepath.Join(target, filepath.FromSlash(root.Path)))
			if statErr != nil || !info.IsDir() {
				return nil, runtimeassets.CodeRootDiscovery{}, fmt.Errorf("agent-declared code root %q does not exist as a directory", root.Path)
			}
		}
		return roots, discoveryForMode(mode), nil
	case CodeRootDiscoveryAgentConfirmedNone:
		return []runtimeassets.CodeRoot{}, discoveryForMode(mode), nil
	case "", CodeRootDiscoveryCLIFallback:
		roots, err := detectCodeRoots(target)
		return roots, discoveryForMode(CodeRootDiscoveryCLIFallback), err
	default:
		return nil, runtimeassets.CodeRootDiscovery{}, fmt.Errorf("unsupported code-root discovery mode %q", mode)
	}
}

func discoveryForMode(mode string) runtimeassets.CodeRootDiscovery {
	status := "confirmed"
	if mode == CodeRootDiscoveryCLIFallback || mode == "" {
		mode = CodeRootDiscoveryCLIFallback
		status = "needs-agent-review"
	}
	return runtimeassets.CodeRootDiscovery{Mode: mode, Status: status}
}

func normalizeCodeRoots(explicit []runtimeassets.CodeRoot) ([]runtimeassets.CodeRoot, error) {
	byPath := map[string]runtimeassets.CodeRoot{}
	for _, root := range explicit {
		rawPath := strings.TrimSpace(root.Path)
		slashPath := strings.ReplaceAll(rawPath, "\\", "/")
		path := pathpkg.Clean(slashPath)
		role := strings.ToLower(strings.TrimSpace(root.Role))
		if rawPath == "" || path == "." && slashPath != "." || role == "" || path == "product" || path == ".." || strings.HasPrefix(path, "../") || strings.HasPrefix(slashPath, "/") || filepath.IsAbs(rawPath) || filepath.VolumeName(rawPath) != "" {
			return nil, fmt.Errorf("invalid code root %q; use a repository-relative sibling path with a role", root.Path)
		}
		if previous, exists := byPath[path]; exists {
			return nil, fmt.Errorf("duplicate code root %q has roles %q and %q; declare one semantic owner", path, previous.Role, role)
		}
		byPath[path] = runtimeassets.CodeRoot{Path: path, Role: role}
	}
	return sortedCodeRoots(byPath), nil
}

// detectCodeRoots inventories immediate repository siblings only. It produces
// candidates for compatibility; the agent-owned pre-init inventory is the
// canonical semantic discovery path.
func detectCodeRoots(target string) ([]runtimeassets.CodeRoot, error) {
	byPath := map[string]runtimeassets.CodeRoot{}
	if role := detectCodeRole(target); role != "" {
		byPath["."] = runtimeassets.CodeRoot{Path: ".", Role: role}
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
	return sortedCodeRoots(byPath), nil
}

func sortedCodeRoots(byPath map[string]runtimeassets.CodeRoot) []runtimeassets.CodeRoot {
	roots := make([]runtimeassets.CodeRoot, 0, len(byPath))
	for _, root := range byPath {
		roots = append(roots, root)
	}
	sort.Slice(roots, func(i, j int) bool { return roots[i].Path < roots[j].Path })
	return roots
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

func writeCodeDiscovery(productRoot string, roots []runtimeassets.CodeRoot, discovery runtimeassets.CodeRootDiscovery) error {
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
	b.WriteString("` |\n| Discovery mode | `")
	b.WriteString(discovery.Mode)
	b.WriteString("` |\n| Discovery status | `")
	b.WriteString(discovery.Status)
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
