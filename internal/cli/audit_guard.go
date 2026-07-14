package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func auditOnlyMutation(args []string) (bool, string) {
	if len(args) == 0 || !auditOnlyActive(productRootArgument(args)) {
		return false, ""
	}
	command := args[0]
	subcommand := ""
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		subcommand = args[1]
	}
	mutating := false
	switch command {
	case "move", "import", "approve":
		mutating = true
	case "approve-batch":
		mutating = boolArgumentEnabled(args, "--yes") && !boolArgumentEnabled(args, "--dry-run")
	case "validate":
		mutating = boolArgumentEnabled(args, "--write-registry") || boolArgumentEnabled(args, "--write-report")
	case "work":
		mutating = argumentValue(args, "--feature") != ""
	case "approve-stage":
		mutating = boolArgumentEnabled(args, "--yes")
	case "graph":
		mutating = subcommand != "" && subcommand != "ready"
	case "design":
		mutating = map[string]bool{"init": true, "import": true, "register": true, "map": true}[subcommand] ||
			(subcommand == "migrate" && !boolArgumentEnabled(args, "--dry-run")) ||
			((subcommand == "audit" || subcommand == "verify-fidelity") && boolArgumentEnabled(args, "--write-report"))
	case "design-system", "engineering-system":
		mutating = subcommand == "init" || (subcommand == "migrate" && !boolArgumentEnabled(args, "--dry-run"))
	case "decisions":
		mutating = boolArgumentEnabled(args, "--yes") || boolArgumentEnabled(args, "--interactive")
	case "migrate":
		mutating = !boolArgumentEnabled(args, "--dry-run")
	case "adapters":
		mutating = subcommand == "install" || subcommand == "update"
	case "handoff", "checkpoint", "lease", "commands", "integrate":
		mutating = true
	case "schedule":
		mutating = true
	case "runtime":
		mutating = !boolArgumentEnabled(args, "--dry-run")
	}
	if !mutating {
		return false, ""
	}
	label := command
	if subcommand != "" {
		label += " " + subcommand
	}
	return true, fmt.Sprintf("audit-only blocks product mutation: %s; use read-only inspection or explicitly transition the starting point", label)
}

func auditOnlyActive(root string) bool {
	data, err := os.ReadFile(filepath.Join(root, ".product", "framework.json"))
	if err != nil {
		return false
	}
	var manifest map[string]any
	if json.Unmarshal(data, &manifest) != nil {
		return false
	}
	point, _ := manifest["starting_point"].(string)
	return point == "audit-only"
}

func productRootArgument(args []string) string {
	root := argumentValue(args, "--product-root")
	if root != "" {
		return absolutePath(root)
	}
	if len(args) > 0 && args[0] == "adapters" {
		if repositoryRoot := argumentValue(args, "--root"); repositoryRoot != "" {
			return filepath.Join(absolutePath(repositoryRoot), "product")
		}
	}
	if len(args) > 0 && args[0] == "migrate" {
		if target := argumentValue(args, "--target"); target != "" {
			return filepath.Join(absolutePath(target), "product")
		}
	}
	return absolutePath("product")
}

func absolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, path)
}

func argumentValue(args []string, name string) string {
	for i, arg := range args {
		if strings.HasPrefix(arg, name+"=") {
			return strings.TrimSpace(strings.TrimPrefix(arg, name+"="))
		}
		if arg == name && i+1 < len(args) {
			return strings.TrimSpace(args[i+1])
		}
	}
	return ""
}

func boolArgumentEnabled(args []string, name string) bool {
	for _, arg := range args {
		if arg == name {
			return true
		}
		if strings.HasPrefix(arg, name+"=") {
			value := strings.TrimPrefix(arg, name+"=")
			enabled, err := strconv.ParseBool(value)
			return err != nil || enabled
		}
	}
	return false
}
