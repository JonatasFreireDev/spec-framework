package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/JonatasFreireDev/spec-framework/internal/designsystem"
)

func runDesignSystem(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "design-system requires init, inspect, validate, or migrate")
		return 2
	}
	command := args[0]
	flags := flag.NewFlagSet("design-system "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	mode := flags.String("mode", "generate", "generate, evolve, or adopt")
	dryRun := flags.Bool("dry-run", false, "preview without writing")
	asJSON := flags.Bool("json", false, "JSON output")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	switch command {
	case "init":
		path, err := designsystem.Init(p, *mode)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Initialized", path)
	case "inspect", "validate":
		i, err := designsystem.Inspect(p)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if *asJSON {
			data, _ := json.MarshalIndent(i, "", "  ")
			fmt.Fprintln(stdout, string(data))
		} else {
			fmt.Fprintf(stdout, "Design System: %s\n- Status: %s\n- Version: %s\n- Origin: %s\n- Tokens: %d\n- Components: %d\n- Patterns: %d\n- Sources: %d\n", i.ID, i.Status, i.Version, i.OriginMode, i.Tokens, i.Components, i.Patterns, i.Sources)
			for _, blocker := range i.Blockers {
				fmt.Fprintln(stdout, "BLOCKED:", blocker)
			}
		}
		if len(i.Blockers) > 0 {
			return 1
		}
	case "migrate":
		items, err := designsystem.Migrate(p, *dryRun)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		for _, item := range items {
			fmt.Fprintln(stdout, item)
		}
	default:
		fmt.Fprintln(stderr, "unknown design-system command", command)
		return 2
	}
	return 0
}
