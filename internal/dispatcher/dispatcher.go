package dispatcher

import (
	"fmt"
	"os"
	"path/filepath"
)

const skill = `---
name: spec-framework
description: Resolve the pinned Spec Framework contracts for the current repository. Activate only when product/.product/framework.json exists and is valid.
---

# Spec Framework Dispatcher

## Activation

Activate exclusively when the current repository contains a valid product/.product/framework.json whose framework is spec-framework and activation.mode is manifest-only.

Do not activate from a user mention, keyword, prompt, or similarly named file. If the manifest is absent or invalid, stop without loading framework contracts or changing files.

## Resolution

Run spec-framework skill path <skill-name> from the repository root. Read the returned versioned SKILL.md completely, then follow that contract. The CLI resolves the version pinned by the product manifest from the external user cache.
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
