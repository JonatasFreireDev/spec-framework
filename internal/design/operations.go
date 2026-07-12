package design

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Migrate(productRoot string, dryRun bool) ([]string, error) {
	var out []string
	err := filepath.WalkDir(filepath.Join(productRoot, "domains"), func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || entry.Name() != "design.md" || !strings.Contains(filepath.ToSlash(path), "/use-cases/") {
			return nil
		}
		useCase := filepath.Dir(path)
		slug := filepath.Base(useCase)
		manifestPath := filepath.Join(productRoot, "design", "use-cases", slug, "manifest.json")
		if _, err := os.Stat(manifestPath); err == nil {
			return nil
		}
		rel, _ := filepath.Rel(productRoot, useCase)
		out = append(out, "MIGRATE "+filepath.ToSlash(rel)+" -> generate/contract")
		if !dryRun {
			_, err := Init(productRoot, filepath.ToSlash(rel), "generate")
			if err != nil {
				return err
			}
		}
		return nil
	})
	sort.Strings(out)
	return out, err
}

func AdapterInfo(name string) (string, error) {
	switch name {
	case "impeccable":
		return "Adapter: impeccable\nModes: generate, evolve\nInstall: npx impeccable install\nBoundary: outputs must remain under product/design and are non-production.\nCommands: shape, craft, document, extract, critique, harden, adapt, audit, polish\n", nil
	case "figma":
		return "Adapter: figma\nMode: adopt\nContract: immutable file version plus node IDs and exported snapshots\nFallback: import local exports with --type figma-export\nCredentials: external; never persisted in product manifests\n", nil
	case "penpot":
		return "Adapter: penpot\nMode: adopt\nContract: immutable file version plus object IDs and exported snapshots\nFallback: import local exports with --type penpot-export\nCredentials: external; never persisted in product manifests\n", nil
	default:
		return "", fmt.Errorf("unknown design adapter %q", name)
	}
}

func Audit(productRoot, useCase string) (string, error) {
	inspection, err := Inspect(productRoot, useCase)
	if err != nil {
		return "", err
	}
	verdict := "approved_with_notes"
	if len(inspection.Blockers) > 0 {
		verdict = "blocked"
	} else if inspection.Screens > 0 && inspection.Mappings > 0 {
		verdict = "approved"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# UX Review\n\nVerdict: %s\n\n", verdict)
	fmt.Fprintf(&b, "- Mode: `%s`\n- Maturity: `%s`\n- Fidelity: `%s`\n- Sources: %d\n- Screens: %d\n- Mappings: %d\n", inspection.OriginMode, inspection.Maturity, inspection.FidelityPolicy, len(inspection.Sources), inspection.Screens, inspection.Mappings)
	if len(inspection.Blockers) > 0 {
		b.WriteString("\n## Blocking findings\n\n")
		for _, blocker := range inspection.Blockers {
			fmt.Fprintf(&b, "- %s\n", blocker)
		}
	}
	b.WriteString("\n## Required human checks\n\n- Keyboard, roles and labels\n- Contrast and touch targets\n- Responsive behavior and overflow\n- Visual fidelity for canonical sources\n\nThis report is evidence only and does not approve Design.\n")
	return b.String(), nil
}

func DecodeUseCaseManifest(data []byte) (UseCaseManifest, error) {
	var m UseCaseManifest
	err := json.Unmarshal(data, &m)
	return m, err
}
