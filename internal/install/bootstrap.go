package install

import (
	"encoding/json"
	"fmt"
	"strings"

	framework "github.com/JonatasFreireDev/spec-framework"
)

type bootstrapCatalog struct {
	SchemaVersion int                                    `json:"schema_version"`
	Profiles      map[string]declarativeBootstrapProfile `json:"profiles"`
}

type declarativeBootstrapProfile struct {
	Title       string                     `json:"title"`
	Simple      string                     `json:"simple_explanation"`
	WhenToUse   string                     `json:"when_to_use"`
	FirstAction string                     `json:"first_action"`
	Steps       []declarativeBootstrapStep `json:"steps"`
	Rules       []string                   `json:"rules"`
	After       string                     `json:"after"`
}

type declarativeBootstrapStep struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Goal   string   `json:"goal"`
	Read   []string `json:"read"`
	Write  []string `json:"write"`
	Prompt string   `json:"prompt"`
	Gate   string   `json:"gate"`
	Next   string   `json:"next"`
}

func declarativeBootstrapFor(startingPoint string) (string, error) {
	data, err := framework.Assets.ReadFile("framework/init/bootstrap.json")
	if err != nil {
		return "", err
	}
	var catalog bootstrapCatalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return "", err
	}
	if catalog.SchemaVersion != 1 {
		return "", fmt.Errorf("unsupported bootstrap schema %d", catalog.SchemaVersion)
	}
	profile, ok := catalog.Profiles[startingPoint]
	if !ok {
		return "", fmt.Errorf("unknown bootstrap profile %q", startingPoint)
	}
	if profile.Title == "" || profile.Simple == "" || len(profile.Steps) == 0 {
		return "", fmt.Errorf("incomplete bootstrap profile %q", startingPoint)
	}
	return renderBootstrap(startingPoint, profile), nil
}

func renderBootstrap(startingPoint string, profile declarativeBootstrapProfile) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Product Bootstrap\n\nStarting point: **%s**\n\n", startingPoint)
	fmt.Fprintf(&b, "## What this means\n\n%s\n\n", profile.Simple)
	fmt.Fprintf(&b, "## When to use this path\n\n%s\n\n", profile.WhenToUse)
	fmt.Fprintf(&b, "## Your first action\n\n%s\n\n", profile.FirstAction)
	b.WriteString("## Guided steps\n\n")
	for index, step := range profile.Steps {
		fmt.Fprintf(&b, "### %d. %s\n\n", index+1, step.Title)
		fmt.Fprintf(&b, "**Goal:** %s\n\n", step.Goal)
		if len(step.Read) > 0 {
			b.WriteString("**Agent reads:**\n\n")
			for _, path := range step.Read {
				fmt.Fprintf(&b, "- `%s`\n", path)
			}
			b.WriteString("\n")
		}
		if len(step.Write) > 0 {
			b.WriteString("**Agent may propose or fill:**\n\n")
			for _, path := range step.Write {
				fmt.Fprintf(&b, "- `%s`\n", path)
			}
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "**Prompt:**\n\n> %s\n\n", strings.ReplaceAll(step.Prompt, "\n", "\n> "))
		fmt.Fprintf(&b, "**Gate:** %s\n\n", step.Gate)
		if step.Next != "" {
			fmt.Fprintf(&b, "**Next:** %s\n\n", step.Next)
		}
	}
	b.WriteString("## Rules for the agent\n\n")
	for _, rule := range profile.Rules {
		fmt.Fprintf(&b, "- %s\n", rule)
	}
	b.WriteString("\n## After this path\n\n")
	b.WriteString(profile.After)
	b.WriteString("\n")
	return b.String()
}
