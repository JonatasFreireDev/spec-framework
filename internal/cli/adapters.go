package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/adapters"
)

func runAdapters(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "adapters requires list, status, doctor, install, or update")
		return 2
	}
	action := args[0]
	flags := flag.NewFlagSet("adapters "+action, flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "repository root")
	version := flags.String("version", "", "explicit provider CLI version")
	yes := flags.Bool("yes", false, "confirm external mutation")
	checkLatest := flags.Bool("check-latest", false, "query npm for the latest version")
	asJSON := flags.Bool("json", false, "JSON output")
	raw := args[1:]
	id := ""
	if len(raw) > 0 && !strings.HasPrefix(raw[0], "-") {
		id = raw[0]
		raw = raw[1:]
	}
	if err := flags.Parse(raw); err != nil {
		return 2
	}
	if id == "" && flags.NArg() > 0 {
		id = flags.Arg(0)
	}
	cwd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(cwd, p)
	}
	switch action {
	case "list":
		items := adapters.Registry()
		if *asJSON {
			data, _ := json.MarshalIndent(items, "", "  ")
			fmt.Fprintln(stdout, string(data))
		} else {
			for _, item := range items {
				fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\n", item.ID, item.Provider, item.Runtime, strings.Join(item.Modes, ","))
			}
		}
	case "status":
		if id == "" {
			fmt.Fprintln(stderr, "adapters status requires an adapter id")
			return 2
		}
		status, err := adapters.Inspect(p, id)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if *asJSON {
			data, _ := json.MarshalIndent(status, "", "  ")
			fmt.Fprintln(stdout, string(data))
		} else {
			fmt.Fprintf(stdout, "Adapter: %s\n- Provider: %s\n- Installed: %t\n- Runtime ready: %t\n", status.ID, status.Provider, status.Installed, status.RuntimeOK)
			for _, path := range status.Paths {
				fmt.Fprintln(stdout, "- Path:", path)
			}
		}
	case "doctor":
		if id == "" {
			fmt.Fprintln(stderr, "adapters doctor requires an adapter id")
			return 2
		}
		doctor, err := adapters.Diagnose(p, id, *checkLatest)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if *asJSON {
			data, _ := json.MarshalIndent(doctor, "", "  ")
			fmt.Fprintln(stdout, string(data))
		} else {
			fmt.Fprintf(stdout, "Adapter doctor: %s\n", doctor.ID)
			for _, check := range doctor.Checks {
				fmt.Fprintln(stdout, "OK", check)
			}
			for _, blocker := range doctor.Blockers {
				fmt.Fprintln(stdout, "BLOCKED", blocker)
			}
		}
		if len(doctor.Blockers) > 0 {
			return 1
		}
	case "install", "update":
		if id == "" {
			fmt.Fprintf(stderr, "adapters %s requires an adapter id\n", action)
			return 2
		}
		status, err := adapters.Inspect(p, id)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		argv, err := adapters.ProviderArgv(id, action, *version)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprintf(stdout, "Adapter mutation preview\n- Action: %s\n- Provider: %s\n- Package: %s@%s\n- Runtime: %s\n- Working directory: %s\n- Command: %s %s\n", action, status.Provider, status.Package, *version, status.NpxPath, filepath.Clean(p), status.NpxPath, strings.Join(argv, " "))
		if !*yes {
			fmt.Fprintln(stdout, "Re-run with --yes to execute the external provider command.")
			return 0
		}
		if !status.RuntimeOK {
			fmt.Fprintln(stderr, "Node.js/npx runtime is not ready")
			return 1
		}
		var providerOut, providerErr bytes.Buffer
		if err := adapters.Execute(p, id, action, *version, &providerOut, &providerErr); err != nil {
			fmt.Fprint(stdout, providerOut.String())
			fmt.Fprint(stderr, providerErr.String())
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprint(stdout, providerOut.String())
		if providerErr.Len() > 0 {
			fmt.Fprint(stderr, providerErr.String())
		}
		after, _ := adapters.Inspect(p, id)
		fmt.Fprintf(stdout, "Adapter %s complete. Installed: %t\n", action, after.Installed)
	default:
		fmt.Fprintln(stderr, "unknown adapters command", action)
		return 2
	}
	return 0
}
