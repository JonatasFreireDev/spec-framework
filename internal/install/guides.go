package install

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeStarterGuides(target, version string, agents []Agent, startingPoint string) error {
	names := make([]string, len(agents))
	for i, agent := range agents {
		names[i] = string(agent)
	}
	selected := strings.Join(names, ", ")
	bootstrap := bootstrapFor(startingPoint)
	header := fmt.Sprintf("Framework version: **%s**\n\nConfigured agents: **%s**\n\n", version, selected)
	return writeFile(filepath.Join(target, "product", "BOOTSTRAP.md"), []byte(header+bootstrap), 0644)
}

func bootstrapFor(startingPoint string) string {
	rendered, err := declarativeBootstrapFor(startingPoint)
	if err != nil {
		return fmt.Sprintf("# Product Bootstrap\n\nStarting point: **%s**\n\nBootstrap profile unavailable: %v\n", startingPoint, err)
	}
	return rendered
}
