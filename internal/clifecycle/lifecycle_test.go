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
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchiveName(t *testing.T) {
	for _, test := range []struct{ goos, arch, want string }{
		{"windows", "amd64", "spec-framework_1.2.3_windows_amd64.zip"},
		{"linux", "arm64", "spec-framework_1.2.3_linux_arm64.tar.gz"},
		{"darwin", "amd64", "spec-framework_1.2.3_darwin_amd64.tar.gz"},
	} {
		got, err := archiveName("1.2.3", test.goos, test.arch)
		if err != nil || got != test.want {
			t.Fatalf("got=%q err=%v want=%q", got, err, test.want)
		}
	}
	if _, err := archiveName("1.2.3", "plan9", "amd64"); err == nil {
		t.Fatal("unsupported OS accepted")
	}
}

func TestCheckRejectsUnsafeExplicitVersion(t *testing.T) {
	manager := testManager(t, "https://example.invalid")
	if _, err := manager.Check(context.Background(), "../../malicious"); err == nil {
		t.Fatal("unsafe version accepted")
	}
}

func TestCheckResolvesLatestRelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/latest" {
			http.NotFound(response, request)
			return
		}
		fmt.Fprint(response, `{"tag_name":"v1.2.3"}`)
	}))
	defer server.Close()
	manager := testManager(t, server.URL)
	release, err := manager.Check(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if release.Latest != "1.2.3" || !release.UpdateAvailable || !strings.HasSuffix(release.URL, "/v1.2.3/spec-framework_1.2.3_windows_amd64.zip") {
		t.Fatalf("release=%+v", release)
	}
}

func TestUpdateVerifiesChecksumAndAtomicallyReplacesBinary(t *testing.T) {
	archive := zipBinary(t, []byte("new binary"))
	sum := sha256.Sum256(archive)
	archiveName := "spec-framework_1.2.3_windows_amd64.zip"
	server := releaseServer(t, archiveName, archive, hex.EncodeToString(sum[:]))
	defer server.Close()
	manager := testManager(t, server.URL)
	result, err := manager.Update(context.Background(), "1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Updated {
		t.Fatal("update did not report replacement")
	}
	data, err := os.ReadFile(manager.Executable)
	if err != nil || string(data) != "new binary" {
		t.Fatalf("binary=%q err=%v", data, err)
	}
	if _, err := os.Stat(manager.Executable + ".old"); !os.IsNotExist(err) {
		t.Fatalf("backup remained: %v", err)
	}
}

func TestChecksumFailurePreservesCurrentBinary(t *testing.T) {
	archive := zipBinary(t, []byte("new binary"))
	server := releaseServer(t, "spec-framework_1.2.3_windows_amd64.zip", archive, strings.Repeat("0", 64))
	defer server.Close()
	manager := testManager(t, server.URL)
	if _, err := manager.Update(context.Background(), "1.2.3"); err == nil {
		t.Fatal("bad checksum accepted")
	}
	data, err := os.ReadFile(manager.Executable)
	if err != nil || string(data) != "old binary" {
		t.Fatalf("current binary changed: %q %v", data, err)
	}
}

func TestCandidateSmokeFailurePreservesCurrentBinary(t *testing.T) {
	archive := zipBinary(t, []byte("new binary"))
	sum := sha256.Sum256(archive)
	server := releaseServer(t, "spec-framework_1.2.3_windows_amd64.zip", archive, hex.EncodeToString(sum[:]))
	defer server.Close()
	manager := testManager(t, server.URL)
	manager.ValidateCandidate = func(context.Context, string, string) error { return fmt.Errorf("smoke failed") }
	if _, err := manager.Update(context.Background(), "1.2.3"); err == nil {
		t.Fatal("failed candidate accepted")
	}
	data, err := os.ReadFile(manager.Executable)
	if err != nil || string(data) != "old binary" {
		t.Fatalf("current binary changed: %q %v", data, err)
	}
}

func TestTarGzExtraction(t *testing.T) {
	archive := tarBinary(t, []byte("unix binary"))
	path, err := extractBinary("spec-framework_1.2.3_linux_amd64.tar.gz", archive, t.TempDir(), "linux")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "unix binary" {
		t.Fatalf("binary=%q err=%v", data, err)
	}
}

func TestPendingUpdateRecoveryRestoresBackup(t *testing.T) {
	root := t.TempDir()
	current := filepath.Join(root, "spec-framework.exe")
	backup := current + ".old"
	candidate := current + ".new"
	if err := os.WriteFile(backup, []byte("old"), 0755); err != nil {
		t.Fatal(err)
	}
	state, _ := json.Marshal(map[string]string{"current": current, "candidate": candidate, "backup": backup})
	if err := os.WriteFile(current+".update.json", state, 0644); err != nil {
		t.Fatal(err)
	}
	manager := Manager{Executable: current}
	if err := manager.recoverPendingUpdate(); err == nil || !strings.Contains(err.Error(), "restored") {
		t.Fatalf("recovery err=%v", err)
	}
	data, err := os.ReadFile(current)
	if err != nil || string(data) != "old" {
		t.Fatalf("restored=%q err=%v", data, err)
	}
}

func TestUninstallAndPurgeStayWithinOwnedPaths(t *testing.T) {
	for _, purge := range []bool{false, true} {
		t.Run(fmt.Sprintf("purge=%v", purge), func(t *testing.T) {
			root := t.TempDir()
			executable := filepath.Join(root, "bin", "spec-framework")
			manifest := filepath.Join(filepath.Dir(executable), "install.json")
			cache := filepath.Join(root, "spec-framework")
			agentHome := filepath.Join(root, "agents")
			product := filepath.Join(root, "repo", "product", "owned.md")
			for _, path := range []string{executable, filepath.Join(cache, "versions", "1", ".complete"), filepath.Join(agentHome, ".agents", "skills", "spec-framework", "SKILL.md"), filepath.Join(agentHome, ".codex", "skills", "spec-framework", "SKILL.md"), product} {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("owned"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			manifestData, err := json.Marshal(map[string]any{"schema_version": 1, "managed_by": "spec-framework-installer", "executable": executable})
			if err != nil || os.WriteFile(manifest, manifestData, 0644) != nil {
				t.Fatalf("write install manifest: %v", err)
			}
			manager := Manager{Executable: executable, GOOS: "linux", CacheRoot: cache, AgentHome: agentHome}
			if plan := manager.PlanUninstall(purge); !plan.Managed || plan.Manifest != manifest {
				t.Fatalf("managed installation not recognized: %+v", plan)
			}
			if _, err := manager.Uninstall(purge); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(executable); !os.IsNotExist(err) {
				t.Fatalf("binary remained: %v", err)
			}
			if _, err := os.Stat(manifest); !os.IsNotExist(err) {
				t.Fatalf("manifest remained: %v", err)
			}
			_, cacheErr := os.Stat(cache)
			_, dispatcherErr := os.Stat(filepath.Join(agentHome, ".agents", "skills", "spec-framework"))
			_, legacyDispatcherErr := os.Stat(filepath.Join(agentHome, ".codex", "skills", "spec-framework"))
			if purge && (!os.IsNotExist(cacheErr) || !os.IsNotExist(dispatcherErr) || !os.IsNotExist(legacyDispatcherErr)) {
				t.Fatalf("purge left cache=%v dispatcher=%v legacy=%v", cacheErr, dispatcherErr, legacyDispatcherErr)
			}
			if !purge && (cacheErr != nil || dispatcherErr != nil || legacyDispatcherErr != nil) {
				t.Fatalf("standard uninstall removed cache=%v dispatcher=%v legacy=%v", cacheErr, dispatcherErr, legacyDispatcherErr)
			}
			if data, err := os.ReadFile(product); err != nil || string(data) != "owned" {
				t.Fatalf("product changed: %q %v", data, err)
			}
		})
	}
}

func TestUpdateManagedManifestPreservesOwnershipAndRefreshesVersion(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "spec-framework")
	manifest := filepath.Join(filepath.Dir(executable), "install.json")
	data := []byte(`{"schema_version":1,"managed_by":"spec-framework-installer","version":"1.0.0","executable":"` + filepath.ToSlash(executable) + `"}`)
	if err := os.WriteFile(manifest, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := updateManagedManifest(executable, "1.1.0"); err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	updated, err := os.ReadFile(manifest)
	if err != nil || json.Unmarshal(updated, &record) != nil {
		t.Fatalf("read refreshed manifest: %v", err)
	}
	if record["version"] != "1.1.0" || record["managed_by"] != "spec-framework-installer" || record["updated_at"] == nil {
		t.Fatalf("unexpected refreshed manifest: %#v", record)
	}
}

func testManager(t *testing.T, server string) Manager {
	t.Helper()
	executable := filepath.Join(t.TempDir(), "spec-framework.exe")
	if err := os.WriteFile(executable, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}
	return Manager{CurrentVersion: "1.0.0", Executable: executable, GOOS: "windows", GOARCH: "amd64", APIBase: server, ReleaseBase: server, Client: http.DefaultClient, ValidateCandidate: func(context.Context, string, string) error { return nil }}
}

func releaseServer(t *testing.T, archiveName string, archive []byte, checksum string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch filepath.Base(request.URL.Path) {
		case archiveName:
			response.Write(archive)
		case "checksums.txt":
			fmt.Fprintf(response, "%s  %s\n", checksum, archiveName)
		default:
			http.NotFound(response, request)
		}
	}))
}

func zipBinary(t *testing.T, data []byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	file, err := writer.Create("spec-framework.exe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func tarBinary(t *testing.T, data []byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	tarWriter := tar.NewWriter(gzipWriter)
	if err := tarWriter.WriteHeader(&tar.Header{Name: "spec-framework", Mode: 0755, Size: int64(len(data)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := tarWriter.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
