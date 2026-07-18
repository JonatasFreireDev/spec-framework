package dispatcher

import (
	"fmt"
	"os"
	"path/filepath"
)

const skillTemplate = `---
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

Before resolving a specialized skill, read the runtime's AGENTS.framework.md common agent rules. Then run spec-framework skill path <skill-name> from the repository root to resolve the selected contract. Read the returned versioned SKILL.md completely, then follow it. The CLI resolves the version pinned by the product manifest from the external user cache. Direct diagnostic CLI commands remain available; dispatch never grants approval or mutation authority.

## Native questions

Definition and planning skills use the harness-native structured question tool when it is available. In this harness, map the canonical native_user_question capability to %s. Inspect repository and CLI evidence before asking, and do not replace an available question tool with silent assumptions or a question buried in prose.
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
	var questionTool string
	switch agent {
	case "codex":
		root = filepath.Join(home, ".agents", "skills", "spec-framework")
		questionTool = "request_user_input"
	case "cursor":
		root = filepath.Join(home, ".cursor", "skills", "spec-framework")
		questionTool = "Cursor's native user-question tool when exposed"
	case "claude":
		root = filepath.Join(home, ".claude", "skills", "spec-framework")
		questionTool = "AskUserQuestion"
	default:
		return "", fmt.Errorf("unsupported agent %q", agent)
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(root, "SKILL.md")
	if err := os.WriteFile(path, []byte(fmt.Sprintf(skillTemplate, questionTool)), 0644); err != nil {
		return "", err
	}
	if agent == "codex" {
		legacyRoot := filepath.Join(home, ".codex", "skills", "spec-framework")
		if err := os.RemoveAll(legacyRoot); err != nil {
			return "", fmt.Errorf("remove legacy Codex dispatcher: %w", err)
		}
	}
	return path, nil
}
