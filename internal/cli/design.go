package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/design"
)

func runDesign(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "design requires init, import, register, inspect, map, verify, migrate, adapter, audit, or verify-fidelity")
		return 2
	}
	command := args[0]
	flags := flag.NewFlagSet("design "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	useCase := flags.String("use-case", "", "use-case product-relative path")
	mode := flags.String("mode", "generate", "generate, evolve, or adopt")
	sourceType := flags.String("type", "images", "source type")
	source := flags.String("source", "", "source file, directory, or URL")
	authority := flags.String("authority", "reference", "visual-canonical, reference, or inspiration")
	sourceID := flags.String("source-id", "", "optional DSRC-NNN id")
	copyAssets := flags.Bool("copy", true, "copy local assets into the product")
	mappingsFile := flags.String("mappings", "", "JSON file containing mapping array")
	jsonOutput := flags.Bool("json", false, "JSON output")
	adapter := flags.String("adapter", "", "adapter name")
	dryRun := flags.Bool("dry-run", false, "preview migration without writing")
	writeReport := flags.Bool("write-report", false, "write UX review evidence under product/design")
	maturity := flags.String("maturity", "wireframe", "target visual maturity")
	version := flags.String("version", "", "immutable external source version")
	nodes := flags.String("nodes", "", "comma-separated Figma node IDs or Penpot object IDs")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	auth := strings.ReplaceAll(*authority, "-", "_")
	switch command {
	case "init":
		path, err := design.Init(p, *useCase, *mode)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Initialized", path)
	case "import":
		if *sourceType != "images" && *sourceType != "figma-export" && *sourceType != "penpot-export" {
			fmt.Fprintln(stderr, "local import currently supports images, figma-export, or penpot-export")
			return 2
		}
		if *source == "" {
			fmt.Fprintln(stderr, "design import requires --source")
			return 2
		}
		manifest, path, err := design.ImportImages(p, *useCase, *source, auth, *sourceID, *copyAssets)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "Imported %s (%d screens)\n- Manifest: %s\n- Version: %s:%s\n", manifest.ID, len(manifest.Screens), path, manifest.Version.Kind, manifest.Version.Value)
	case "register":
		metadata := map[string]string{}
		if strings.TrimSpace(*nodes) != "" {
			metadata["selection"] = *nodes
		}
		manifest, path, err := design.RegisterRemote(p, *useCase, *sourceType, *source, *version, auth, *sourceID, metadata)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "Registered %s source %s\n- Manifest: %s\n- Version: %s\n", manifest.Type, manifest.ID, path, manifest.Version.Value)
	case "inspect", "verify":
		result, err := design.Inspect(p, *useCase)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		writeDesignInspection(result, *jsonOutput, stdout)
		if len(result.Blockers) > 0 {
			return 1
		}
	case "map":
		if *mappingsFile == "" {
			fmt.Fprintln(stderr, "design map requires --mappings <json-file>")
			return 2
		}
		data, err := os.ReadFile(*mappingsFile)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		var mappings []design.Mapping
		if err := json.Unmarshal(data, &mappings); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		path, err := design.UpdateMappings(p, *useCase, mappings)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "Mapped %d requirements in %s\n", len(mappings), path)
	case "migrate":
		result, err := design.Migrate(p, *dryRun)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		for _, item := range result {
			fmt.Fprintln(stdout, item)
		}
	case "adapter":
		if *adapter == "" {
			fmt.Fprintln(stderr, "design adapter requires --adapter impeccable|figma|penpot")
			return 2
		}
		if *adapter == "impeccable" && *useCase != "" {
			plan, err := design.ImpeccablePlan(p, *useCase, *maturity)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			fmt.Fprint(stdout, design.EncodeAdapterPlan(plan))
			return 0
		}
		info, err := design.AdapterInfo(*adapter)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprint(stdout, info)
	case "audit", "verify-fidelity":
		report, err := design.Audit(p, *useCase)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprint(stdout, report)
		if *writeReport {
			path, err := design.WriteAudit(p, *useCase, report)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			fmt.Fprintln(stdout, "Evidence:", path)
		}
	default:
		fmt.Fprintln(stderr, "unknown design command", command)
		return 2
	}
	return 0
}

func writeDesignInspection(result design.Inspection, asJSON bool, output io.Writer) {
	if asJSON {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprintln(output, string(data))
		return
	}
	fmt.Fprintf(output, "Design: %s\n- Mode: %s\n- Maturity: %s\n- Fidelity: %s\n- Sources: %d\n- Screens: %d\n- Mappings: %d\n", result.UseCase, result.OriginMode, result.Maturity, result.FidelityPolicy, len(result.Sources), result.Screens, result.Mappings)
	for _, blocker := range result.Blockers {
		fmt.Fprintln(output, "BLOCKED:", blocker)
	}
}
