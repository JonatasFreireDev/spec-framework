package adapters

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Definition struct {
	ID             string   `json:"id"`
	Provider       string   `json:"provider"`
	Package        string   `json:"package"`
	Modes          []string `json:"modes"`
	Runtime        string   `json:"runtime"`
	InstallCommand []string `json:"installCommand"`
	UpdateCommand  []string `json:"updateCommand"`
}

type Status struct {
	Definition
	Installed bool     `json:"installed"`
	Paths     []string `json:"paths,omitempty"`
	RuntimeOK bool     `json:"runtimeReady"`
	NpxPath   string   `json:"npxPath,omitempty"`
}

type Doctor struct {
	Status
	LatestVersion string   `json:"latestVersion,omitempty"`
	Checks        []string `json:"checks"`
	Blockers      []string `json:"blockers,omitempty"`
}

func Registry() []Definition {
	return []Definition{{ID: "impeccable", Provider: "pbakaus/impeccable", Package: "impeccable", Modes: []string{"generate", "evolve"}, Runtime: "node+npx", InstallCommand: []string{"skills", "install"}, UpdateCommand: []string{"skills", "update"}}}
}

func Lookup(id string) (Definition, error) {
	for _, item := range Registry() {
		if item.ID == id {
			return item, nil
		}
	}
	return Definition{}, fmt.Errorf("unknown adapter %q", id)
}

func Inspect(root, id string) (Status, error) {
	definition, err := Lookup(id)
	if err != nil {
		return Status{}, err
	}
	npx, _ := findNpx()
	status := Status{Definition: definition, RuntimeOK: npx != "", NpxPath: npx}
	patterns := []string{
		filepath.Join(root, ".agents", "skills", id, "SKILL.md"),
		filepath.Join(root, ".claude", "skills", id, "SKILL.md"),
		filepath.Join(root, ".cursor", "skills", id, "SKILL.md"),
		filepath.Join(root, ".github", "skills", id, "SKILL.md"),
		filepath.Join(root, ".gemini", "skills", id, "SKILL.md"),
	}
	for _, path := range patterns {
		if _, err := os.Stat(path); err == nil {
			rel, _ := filepath.Rel(root, path)
			status.Paths = append(status.Paths, filepath.ToSlash(rel))
		}
	}
	sort.Strings(status.Paths)
	status.Installed = len(status.Paths) > 0
	return status, nil
}

func Diagnose(root, id string, checkLatest bool) (Doctor, error) {
	status, err := Inspect(root, id)
	if err != nil {
		return Doctor{}, err
	}
	d := Doctor{Status: status}
	if status.RuntimeOK {
		d.Checks = append(d.Checks, "npx runtime found: "+status.NpxPath)
	} else {
		d.Blockers = append(d.Blockers, "Node.js/npx is not available")
	}
	if status.Installed {
		d.Checks = append(d.Checks, fmt.Sprintf("adapter skill found in %d harness path(s)", len(status.Paths)))
	} else {
		d.Blockers = append(d.Blockers, "adapter is not installed in a supported project harness path")
	}
	if checkLatest && status.RuntimeOK {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		npm, npmErr := findNpm()
		if npmErr != nil {
			d.Blockers = append(d.Blockers, npmErr.Error())
			return d, nil
		}
		cmd := exec.CommandContext(ctx, npm, "view", status.Package, "version")
		var output bytes.Buffer
		cmd.Stdout, cmd.Stderr = &output, &output
		if err := cmd.Run(); err != nil {
			d.Blockers = append(d.Blockers, "could not query latest provider version: "+strings.TrimSpace(output.String()))
		} else {
			d.LatestVersion = strings.TrimSpace(output.String())
			d.Checks = append(d.Checks, "latest npm version: "+d.LatestVersion)
		}
	}
	return d, nil
}

func ProviderArgv(id, action, version string) ([]string, error) {
	definition, err := Lookup(id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(version) == "" {
		return nil, fmt.Errorf("%s requires an explicit --version", action)
	}
	var suffix []string
	switch action {
	case "install":
		suffix = definition.InstallCommand
	case "update":
		suffix = definition.UpdateCommand
	default:
		return nil, fmt.Errorf("unsupported adapter action %q", action)
	}
	return append([]string{definition.Package + "@" + version}, suffix...), nil
}

func ResolveVersion(id, requested string) (string, error) {
	definition, err := Lookup(id)
	if err != nil {
		return "", err
	}
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return "", errors.New("an explicit version or latest is required")
	}
	if requested != "latest" {
		if !regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`).MatchString(requested) {
			return "", fmt.Errorf("invalid provider version %q; use semantic version or latest", requested)
		}
		return requested, nil
	}
	npm, err := findNpm()
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, npm, "view", definition.Package, "version")
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("resolve latest %s version: %s", id, strings.TrimSpace(output.String()))
	}
	resolved := strings.TrimSpace(output.String())
	if !regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`).MatchString(resolved) {
		return "", fmt.Errorf("provider returned invalid version %q", resolved)
	}
	return resolved, nil
}

func Execute(root, id, action, version string, stdout, stderr io.Writer) error {
	npx, err := findNpx()
	if err != nil {
		return err
	}
	argv, err := ProviderArgv(id, action, version)
	if err != nil {
		return err
	}
	cmd := exec.Command(npx, argv...)
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func findNpx() (string, error) {
	names := []string{"npx"}
	if runtime.GOOS == "windows" {
		names = []string{"npx.cmd", "npx.exe", "npx"}
	}
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("Node.js/npx runtime not found")
}

func findNpm() (string, error) {
	names := []string{"npm"}
	if runtime.GOOS == "windows" {
		names = []string{"npm.cmd", "npm.exe", "npm"}
	}
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("npm runtime not found")
}
