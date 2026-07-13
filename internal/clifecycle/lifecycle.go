package clifecycle

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	defaultAPIBase     = "https://api.github.com/repos/JonatasFreireDev/spec-framework/releases"
	defaultReleaseBase = "https://github.com/JonatasFreireDev/spec-framework/releases/download"
)

type Manager struct {
	CurrentVersion    string
	Executable        string
	GOOS              string
	GOARCH            string
	APIBase           string
	ReleaseBase       string
	CacheRoot         string
	AgentHome         string
	Client            *http.Client
	ValidateCandidate func(context.Context, string, string) error
}

type Release struct {
	Current         string
	Latest          string
	Archive         string
	URL             string
	ChecksumsURL    string
	UpdateAvailable bool
}

type UpdateResult struct {
	Release   Release
	Updated   bool
	Scheduled bool
}

type UninstallPlan struct {
	Executable  string
	InstallDir  string
	Manifest    string
	Managed     bool
	CacheRoot   string
	Dispatchers []string
	Purge       bool
}

func Default(version string) (Manager, error) {
	executable, err := os.Executable()
	if err != nil {
		return Manager{}, err
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return Manager{}, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Manager{}, err
	}
	if override := strings.TrimSpace(os.Getenv("SPEC_FRAMEWORK_CACHE")); override != "" {
		cache = override
	} else {
		cache = filepath.Join(cache, "spec-framework")
	}
	if override := strings.TrimSpace(os.Getenv("SPEC_FRAMEWORK_AGENT_HOME")); override != "" {
		home = override
	}
	return Manager{
		CurrentVersion: normalizeVersion(version), Executable: executable,
		GOOS: runtime.GOOS, GOARCH: runtime.GOARCH,
		APIBase: defaultAPIBase, ReleaseBase: defaultReleaseBase,
		CacheRoot: cache, AgentHome: home,
		Client: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (manager Manager) Check(ctx context.Context, requested string) (Release, error) {
	if err := manager.recoverPendingUpdate(); err != nil {
		return Release{}, err
	}
	version := normalizeVersion(requested)
	if version == "" {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(manager.APIBase, "/")+"/latest", nil)
		if err != nil {
			return Release{}, err
		}
		manager.setRequestHeaders(request)
		response, err := manager.client().Do(request)
		if err != nil {
			return Release{}, fmt.Errorf("resolve latest release: %w", err)
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return Release{}, fmt.Errorf("resolve latest release: HTTP %d", response.StatusCode)
		}
		var payload struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload); err != nil {
			return Release{}, fmt.Errorf("parse latest release: %w", err)
		}
		version = normalizeVersion(payload.TagName)
	}
	if version == "" {
		return Release{}, errors.New("release version is empty")
	}
	if !regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z.-]*$`).MatchString(version) || strings.Contains(version, "..") {
		return Release{}, fmt.Errorf("invalid release version %q", version)
	}
	archive, err := archiveName(version, manager.GOOS, manager.GOARCH)
	if err != nil {
		return Release{}, err
	}
	base := strings.TrimRight(manager.ReleaseBase, "/") + "/v" + version
	return Release{
		Current: manager.CurrentVersion, Latest: version, Archive: archive,
		URL: base + "/" + archive, ChecksumsURL: base + "/checksums.txt",
		UpdateAvailable: normalizeVersion(manager.CurrentVersion) != version,
	}, nil
}

func (manager Manager) Update(ctx context.Context, requested string) (UpdateResult, error) {
	release, err := manager.Check(ctx, requested)
	if err != nil {
		return UpdateResult{}, err
	}
	if !release.UpdateAvailable {
		return UpdateResult{Release: release}, nil
	}
	if manager.CurrentVersion == "dev" {
		return UpdateResult{}, errors.New("cannot replace a development binary; install a released CLI first")
	}
	archive, err := manager.download(ctx, release.URL, 256<<20)
	if err != nil {
		return UpdateResult{}, err
	}
	checksums, err := manager.download(ctx, release.ChecksumsURL, 4<<20)
	if err != nil {
		return UpdateResult{}, err
	}
	if err := verifyChecksum(release.Archive, archive, checksums); err != nil {
		return UpdateResult{}, err
	}
	installDir := filepath.Dir(manager.Executable)
	staging, err := os.MkdirTemp(installDir, ".spec-framework-update-")
	if err != nil {
		return UpdateResult{}, err
	}
	defer os.RemoveAll(staging)
	newBinary, err := extractBinary(release.Archive, archive, staging, manager.GOOS)
	if err != nil {
		return UpdateResult{}, err
	}
	if info, err := os.Stat(newBinary); err != nil || info.Size() == 0 {
		return UpdateResult{}, errors.New("release archive contains an empty or missing CLI binary")
	}
	if manager.GOOS != "windows" {
		if err := os.Chmod(newBinary, 0755); err != nil {
			return UpdateResult{}, err
		}
	}
	if err := manager.validateCandidate(ctx, newBinary, release.Latest); err != nil {
		return UpdateResult{}, err
	}
	if manager.GOOS == "windows" && isCurrentExecutable(manager.Executable) {
		if err := scheduleWindowsUpdate(manager.Executable, newBinary, release.Latest); err != nil {
			return UpdateResult{}, err
		}
		return UpdateResult{Release: release, Updated: true, Scheduled: true}, nil
	}
	if err := replaceExecutable(manager.Executable, newBinary, manager.GOOS); err != nil {
		return UpdateResult{}, err
	}
	if err := updateManagedManifest(manager.Executable, release.Latest); err != nil {
		return UpdateResult{}, fmt.Errorf("CLI updated but install manifest could not be refreshed: %w", err)
	}
	return UpdateResult{Release: release, Updated: true}, nil
}

func (manager Manager) PlanUninstall(purge bool) UninstallPlan {
	manifest := filepath.Join(filepath.Dir(manager.Executable), "install.json")
	plan := UninstallPlan{Executable: manager.Executable, InstallDir: filepath.Dir(manager.Executable), Manifest: manifest, Purge: purge}
	if data, err := os.ReadFile(manifest); err == nil {
		var record struct {
			SchemaVersion int    `json:"schema_version"`
			ManagedBy     string `json:"managed_by"`
			Executable    string `json:"executable"`
		}
		if json.Unmarshal(data, &record) == nil && record.SchemaVersion == 1 && record.ManagedBy == "spec-framework-installer" && filepath.Clean(record.Executable) == filepath.Clean(manager.Executable) {
			plan.Managed = true
		}
	}
	if purge {
		plan.CacheRoot = manager.CacheRoot
		for _, harness := range []string{".codex", ".cursor", ".claude"} {
			plan.Dispatchers = append(plan.Dispatchers, filepath.Join(manager.AgentHome, harness, "skills", "spec-framework"))
		}
	}
	return plan
}

func (manager Manager) Uninstall(purge bool) (UninstallPlan, error) {
	plan := manager.PlanUninstall(purge)
	if purge {
		if plan.CacheRoot != "" && filepath.Base(filepath.Clean(plan.CacheRoot)) != "spec-framework" {
			return plan, fmt.Errorf("refusing unsafe cache root %s", plan.CacheRoot)
		}
		for _, path := range plan.Dispatchers {
			if filepath.Base(path) != "spec-framework" {
				return plan, fmt.Errorf("refusing unsafe dispatcher path %s", path)
			}
		}
	}
	if manager.GOOS == "windows" {
		return plan, scheduleWindowsUninstall(plan)
	}
	paths := []string{manager.Executable}
	if plan.Managed {
		paths = append(paths, plan.Manifest)
	}
	if purge {
		paths = append(paths, plan.Dispatchers...)
		if plan.CacheRoot != "" {
			paths = append(paths, plan.CacheRoot)
		}
	}
	if err := transactionalRemove(paths); err != nil {
		return plan, err
	}
	return plan, nil
}

func (manager Manager) client() *http.Client {
	if manager.Client != nil {
		return manager.Client
	}
	return &http.Client{Timeout: 30 * time.Second}
}

func (manager Manager) download(ctx context.Context, url string, limit int64) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	manager.setRequestHeaders(request)
	response, err := manager.client().Do(request)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: HTTP %d", url, response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("download %s exceeds size limit", url)
	}
	return data, nil
}

func (manager Manager) setRequestHeaders(request *http.Request) {
	request.Header.Set("User-Agent", "spec-framework/"+manager.CurrentVersion)
	if token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
}

func archiveName(version, goos, goarch string) (string, error) {
	if goarch != "amd64" && goarch != "arm64" {
		return "", fmt.Errorf("unsupported architecture %q", goarch)
	}
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	} else if goos != "linux" && goos != "darwin" {
		return "", fmt.Errorf("unsupported operating system %q", goos)
	}
	return fmt.Sprintf("spec-framework_%s_%s_%s%s", version, goos, goarch, ext), nil
}

func normalizeVersion(version string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(version), "v"))
}

func verifyChecksum(name string, archive, checksums []byte) error {
	want := ""
	for _, line := range strings.Split(string(checksums), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.TrimPrefix(fields[len(fields)-1], "*") == name {
			want = strings.ToLower(fields[0])
			break
		}
	}
	if want == "" {
		return fmt.Errorf("checksum for %s not found", name)
	}
	sum := sha256.Sum256(archive)
	if hex.EncodeToString(sum[:]) != want {
		return fmt.Errorf("checksum verification failed for %s", name)
	}
	return nil
}

func extractBinary(archiveName string, data []byte, target, goos string) (string, error) {
	binaryName := "spec-framework"
	if goos == "windows" {
		binaryName += ".exe"
	}
	dest := filepath.Join(target, binaryName)
	if strings.HasSuffix(archiveName, ".zip") {
		reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return "", err
		}
		for _, file := range reader.File {
			if filepath.Base(file.Name) != binaryName {
				continue
			}
			input, err := file.Open()
			if err != nil {
				return "", err
			}
			err = copyLimitedFile(dest, input, 256<<20)
			input.Close()
			if err != nil {
				return "", err
			}
			return dest, nil
		}
		return "", fmt.Errorf("archive does not contain %s", binaryName)
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == binaryName {
			if err := copyLimitedFile(dest, tarReader, 256<<20); err != nil {
				return "", err
			}
			return dest, nil
		}
	}
	return "", fmt.Errorf("archive does not contain %s", binaryName)
}

func copyLimitedFile(path string, input io.Reader, limit int64) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	written, copyErr := io.Copy(file, io.LimitReader(input, limit+1))
	closeErr := file.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if written > limit {
		return errors.New("extracted binary exceeds size limit")
	}
	return nil
}

func replaceExecutable(current, candidate, goos string) error {
	backup := current + ".old"
	_ = os.Remove(backup)
	if goos != "windows" {
		if err := copyPath(current, backup); err != nil {
			return fmt.Errorf("prepare executable backup: %w", err)
		}
		if err := os.Rename(candidate, current); err != nil {
			_ = os.Remove(current)
			_ = os.Rename(backup, current)
			return fmt.Errorf("install updated executable: %w", err)
		}
		_ = os.Remove(backup)
		return nil
	}
	if err := os.Rename(current, backup); err != nil {
		return fmt.Errorf("prepare executable replacement: %w", err)
	}
	if err := os.Rename(candidate, current); err != nil {
		_ = os.Rename(backup, current)
		return fmt.Errorf("install updated executable: %w", err)
	}
	_ = os.Remove(backup)
	return nil
}

func (manager Manager) validateCandidate(ctx context.Context, path, version string) error {
	if manager.ValidateCandidate != nil {
		return manager.ValidateCandidate(ctx, path, version)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, path, "version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("candidate CLI smoke failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	if strings.TrimSpace(string(output)) != "spec-framework "+version {
		return fmt.Errorf("candidate CLI reports %q, expected %q", strings.TrimSpace(string(output)), "spec-framework "+version)
	}
	return nil
}

func copyPath(source, target string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return err
	}
	output, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		_ = os.Remove(target)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(target)
		return closeErr
	}
	return nil
}

func isCurrentExecutable(path string) bool {
	current, err := os.Executable()
	if err != nil {
		return false
	}
	current, currentErr := filepath.EvalSymlinks(current)
	path, pathErr := filepath.EvalSymlinks(path)
	if currentErr != nil || pathErr != nil {
		return false
	}
	return strings.EqualFold(filepath.Clean(current), filepath.Clean(path))
}

func scheduleWindowsUpdate(current, candidate, version string) error {
	staged := current + ".new"
	statePath := current + ".update.json"
	data, err := os.ReadFile(candidate)
	if err != nil {
		return err
	}
	if err := os.WriteFile(staged, data, 0755); err != nil {
		return fmt.Errorf("stage Windows executable replacement: %w", err)
	}
	state, _ := json.MarshalIndent(map[string]string{"operation": "update", "status": "scheduled", "current": current, "candidate": staged, "backup": current + ".old"}, "", "  ")
	if err := os.WriteFile(statePath, append(state, '\n'), 0644); err != nil {
		_ = os.Remove(staged)
		return err
	}
	script, err := os.CreateTemp("", "spec-framework-update-*.ps1")
	if err != nil {
		_ = os.Remove(staged)
		_ = os.Remove(statePath)
		return err
	}
	currentQuoted := strings.ReplaceAll(current, "'", "''")
	stagedQuoted := strings.ReplaceAll(staged, "'", "''")
	backupQuoted := strings.ReplaceAll(current+".old", "'", "''")
	stateQuoted := strings.ReplaceAll(statePath, "'", "''")
	manifestQuoted := strings.ReplaceAll(filepath.Join(filepath.Dir(current), "install.json"), "'", "''")
	versionQuoted := strings.ReplaceAll(version, "'", "''")
	fmt.Fprintf(script, `Start-Sleep -Milliseconds 750
$current='%s'; $staged='%s'; $backup='%s'
Remove-Item -LiteralPath $backup -Force -ErrorAction SilentlyContinue
try {
  Move-Item -LiteralPath $current -Destination $backup -Force
  Move-Item -LiteralPath $staged -Destination $current -Force
  $manifest='%s'
  if (Test-Path -LiteralPath $manifest) {
    $record=Get-Content -LiteralPath $manifest -Raw | ConvertFrom-Json
    if ($record.schema_version -eq 1 -and $record.managed_by -eq 'spec-framework-installer') {
      $record.version='%s'; $record.updated_at=[DateTime]::UtcNow.ToString('o')
      $record | ConvertTo-Json | Set-Content -LiteralPath $manifest -Encoding UTF8
    }
  }
  Remove-Item -LiteralPath $backup -Force -ErrorAction SilentlyContinue
  Remove-Item -LiteralPath '%s' -Force -ErrorAction SilentlyContinue
} catch {
  if ((Test-Path -LiteralPath $backup) -and -not (Test-Path -LiteralPath $current)) { Move-Item -LiteralPath $backup -Destination $current -Force }
  Remove-Item -LiteralPath $staged -Force -ErrorAction SilentlyContinue
}
Remove-Item -LiteralPath $PSCommandPath -Force -ErrorAction SilentlyContinue
`, currentQuoted, stagedQuoted, backupQuoted, manifestQuoted, versionQuoted, stateQuoted)
	if err := script.Close(); err != nil {
		_ = os.Remove(staged)
		return err
	}
	command := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", script.Name())
	if err := command.Start(); err != nil {
		_ = os.Remove(staged)
		_ = os.Remove(statePath)
		return err
	}
	return nil
}

func updateManagedManifest(executable, version string) error {
	path := filepath.Join(filepath.Dir(executable), "install.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var record map[string]any
	if err := json.Unmarshal(data, &record); err != nil {
		return nil
	}
	if record["schema_version"] != float64(1) || record["managed_by"] != "spec-framework-installer" {
		return nil
	}
	record["version"] = version
	record["updated_at"] = time.Now().UTC().Format(time.RFC3339Nano)
	updated, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, append(updated, '\n'), 0644); err != nil {
		return err
	}
	if err := os.Rename(temporary, path); err != nil {
		_ = os.Remove(temporary)
		return err
	}
	return nil
}

func (manager Manager) recoverPendingUpdate() error {
	statePath := manager.Executable + ".update.json"
	data, err := os.ReadFile(statePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var state struct{ Current, Candidate, Backup string }
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("invalid pending update state: %w", err)
	}
	if state.Current != manager.Executable || state.Candidate == "" || state.Backup == "" {
		return errors.New("unsafe pending update state")
	}
	_, currentErr := os.Stat(state.Current)
	_, candidateErr := os.Stat(state.Candidate)
	_, backupErr := os.Stat(state.Backup)
	if currentErr == nil && os.IsNotExist(candidateErr) {
		_ = os.Remove(state.Backup)
		return os.Remove(statePath)
	}
	if os.IsNotExist(currentErr) && backupErr == nil {
		if err := os.Rename(state.Backup, state.Current); err != nil {
			return fmt.Errorf("recover previous CLI update: %w", err)
		}
		_ = os.Remove(state.Candidate)
		_ = os.Remove(statePath)
		return errors.New("previous CLI update failed and the prior binary was restored")
	}
	if currentErr == nil && candidateErr == nil {
		_ = os.Remove(state.Candidate)
		_ = os.Remove(state.Backup)
		_ = os.Remove(statePath)
		return errors.New("previous Windows CLI update did not complete; staged files were cleaned")
	}
	return errors.New("previous CLI update is incomplete and requires manual inspection")
}

func sameOrParent(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func transactionalRemove(paths []string) error {
	paths = collapsePaths(paths)
	type movedPath struct{ original, staged string }
	var moved []movedPath
	rollback := func() {
		for i := len(moved) - 1; i >= 0; i-- {
			_ = os.Rename(moved[i].staged, moved[i].original)
		}
	}
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		} else if err != nil {
			rollback()
			return err
		}
		staged := path + ".uninstalling"
		_ = os.RemoveAll(staged)
		if err := os.Rename(path, staged); err != nil {
			rollback()
			return fmt.Errorf("stage uninstall path %s: %w", path, err)
		}
		moved = append(moved, movedPath{path, staged})
	}
	for _, path := range moved {
		if err := os.RemoveAll(path.staged); err != nil {
			return fmt.Errorf("remove staged uninstall path %s: %w", path.staged, err)
		}
	}
	return nil
}

func collapsePaths(paths []string) []string {
	var out []string
	for _, candidate := range paths {
		if candidate == "" {
			continue
		}
		covered := false
		for _, parent := range paths {
			if parent != candidate && sameOrParent(parent, candidate) {
				covered = true
				break
			}
		}
		if !covered {
			out = append(out, candidate)
		}
	}
	return out
}

func scheduleWindowsUninstall(plan UninstallPlan) error {
	renamed := plan.Executable + ".uninstalling"
	_ = os.Remove(renamed)
	if err := os.Rename(plan.Executable, renamed); err != nil {
		return fmt.Errorf("prepare uninstall: %w", err)
	}
	script, err := os.CreateTemp("", "spec-framework-uninstall-*.ps1")
	if err != nil {
		_ = os.Rename(renamed, plan.Executable)
		return err
	}
	paths := []string{renamed}
	if plan.Managed {
		paths = append(paths, plan.Manifest)
	}
	if plan.Purge && plan.CacheRoot != "" {
		paths = append(paths, plan.CacheRoot)
	}
	if plan.Purge {
		paths = append(paths, plan.Dispatchers...)
	}
	for _, path := range paths {
		quoted := strings.ReplaceAll(path, "'", "''")
		fmt.Fprintf(script, "Start-Sleep -Milliseconds 750\nRemove-Item -LiteralPath '%s' -Recurse -Force -ErrorAction SilentlyContinue\n", quoted)
	}
	install := strings.ReplaceAll(plan.InstallDir, "'", "''")
	fmt.Fprintf(script, "$p=[Environment]::GetEnvironmentVariable('Path','User'); $n=(($p -split ';') | Where-Object { $_ -and $_ -ne '%s' }) -join ';'; [Environment]::SetEnvironmentVariable('Path',$n,'User')\nRemove-Item -LiteralPath $PSCommandPath -Force -ErrorAction SilentlyContinue\n", install)
	if err := script.Close(); err != nil {
		return err
	}
	command := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", script.Name())
	if err := command.Start(); err != nil {
		_ = os.Rename(renamed, plan.Executable)
		return err
	}
	return nil
}
