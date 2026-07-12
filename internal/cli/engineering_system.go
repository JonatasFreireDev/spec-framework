package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/engineeringsystem"
)

func runEngineeringSystem(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "engineering-system requires inspect, validate, or triggers")
		return 2
	}
	command := args[0]
	flags := flag.NewFlagSet("engineering-system "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	asJSON := flags.Bool("json", false, "JSON output")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if command == "triggers" {
		items := engineeringsystem.AllowedTriggers()
		if *asJSON {
			data, _ := json.MarshalIndent(map[string]any{"triggers": items}, "", "  ")
			fmt.Fprintln(stdout, string(data))
		} else {
			fmt.Fprintln(stdout, strings.Join(items, "\n"))
		}
		return 0
	}
	cwd, _ := os.Getwd()
	productRoot := *root
	if !filepath.IsAbs(productRoot) {
		productRoot = filepath.Join(cwd, productRoot)
	}
	if command != "inspect" && command != "validate" {
		fmt.Fprintln(stderr, "unknown engineering-system command", command)
		return 2
	}
	inspection, err := engineeringsystem.Inspect(productRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *asJSON {
		data, _ := json.MarshalIndent(inspection, "", "  ")
		fmt.Fprintln(stdout, string(data))
	} else {
		fmt.Fprintf(stdout, "Engineering System: %s\n- Status: %s\n- Version: %s\n- Origin: %s\n- Scope: %s\n- Areas: %d\n- Decisions: %d\n- Standards: %d\n- Fitness functions: %d\n", inspection.ID, inspection.Status, inspection.Version, inspection.OriginMode, inspection.Scope, len(inspection.Areas), inspection.Decisions, inspection.Standards, inspection.FitnessFunctions)
		for _, area := range inspection.Areas {
			fmt.Fprintf(stdout, "- %s: %s (%s, evidence=%d)\n", area.Name, area.Maturity, area.Contract, area.Evidence)
		}
		for _, blocker := range inspection.Blockers {
			fmt.Fprintln(stdout, "BLOCKED:", blocker)
		}
	}
	if command == "validate" && len(inspection.Blockers) > 0 {
		return 1
	}
	return 0
}
