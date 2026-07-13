package dispatcher

import (
	"fmt"
	"os"
	"path/filepath"
)

const skill = `---
name: spec-framework
description: Resolve the pinned Spec Framework contracts for the current repository through Guide-first dispatch. Activate only when product/.product/framework.json exists and is valid.
---

# Spec Framework Dispatcher

## Activation

Activate exclusively when the current repository contains a valid product/.product/framework.json whose framework is spec-framework and activation.mode is manifest-only.

Do not activate from a user mention, keyword, prompt, or similarly named file. If the manifest is absent or invalid, stop without loading framework contracts or changing files.

## Resolution

Use Guide-first dispatch for framework-governed product operations.

Resolve framework-guide first unless one of these verified direct routes exists:

- current-session spec-framework guide, dashboard, status, or next output names the workspace, concrete feature or use-case scope, current gate, and owner skill;
- the human explicitly names both the specialist and the concrete artifact or workspace scope.

A persisted handoff or checkpoint identifies where to resume but is not direct-route evidence by itself. Revalidate it with spec-framework dashboard, status, next, or guide. A skill name, keyword, or remembered chat instruction without concrete scope is only a hint. Resolve framework-guide first. Before following a direct route, validate the manifest, scope, ownership, gate, and staleness against current mechanical state. Return to framework-guide when the route is missing, stale, ambiguous, or conflicting.

Run spec-framework skill path <skill-name> from the repository root to resolve the selected contract. Read the returned versioned SKILL.md completely, then follow it. The CLI resolves the version pinned by the product manifest from the external user cache. Direct diagnostic CLI commands remain available; dispatch never grants approval or mutation authority.
`

func Install(agent string) (string, error) {
	home := os.Getenv("SPEC_FRAMEWORK_AGENT_HOME")
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	var root string
	switch agent {
	case "codex":
		root = filepath.Join(home, ".codex", "skills", "spec-framework")
	case "cursor":
		root = filepath.Join(home, ".cursor", "skills", "spec-framework")
	case "claude":
		root = filepath.Join(home, ".claude", "skills", "spec-framework")
	default:
		return "", fmt.Errorf("unsupported agent %q", agent)
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(root, "SKILL.md")
	if err := os.WriteFile(path, []byte(skill), 0644); err != nil {
		return "", err
	}
	return path, nil
}
