package design

import (
	"os"
	"path/filepath"
	"testing"
)

func fixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	uc := filepath.Join(root, "domains", "events", "goals", "join", "features", "qr", "use-cases", "check-in")
	if err := os.MkdirAll(uc, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uc, "specification.md"), []byte("# Specification\nREQ-UX-001\nAC-001\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rel, _ := filepath.Rel(root, uc)
	return root, filepath.ToSlash(rel)
}

func TestInitImportMapInspect(t *testing.T) {
	root, uc := fixture(t)
	if _, err := Init(root, uc, "generate"); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(t.TempDir(), "mobile-default.png")
	if err := os.WriteFile(source, []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest, _, err := ImportImages(root, uc, source, "visual_canonical", "DSRC-001", true)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Version.Kind != "sha256" || len(manifest.Screens) != 1 {
		t.Fatalf("unexpected source manifest: %+v", manifest)
	}
	_, err = UpdateMappings(root, uc, []Mapping{{Requirement: "REQ-UX-001", Criterion: "AC-001", Screen: "SCREEN-001", State: "default", Coverage: "covered"}})
	if err != nil {
		t.Fatal(err)
	}
	inspection, err := Inspect(root, uc)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.OriginMode != "adopt" || inspection.Maturity != "mockup" || inspection.Mappings != 1 || len(inspection.Blockers) != 0 {
		t.Fatalf("unexpected inspection: %+v", inspection)
	}
}

func TestMissingMappingBlocks(t *testing.T) {
	root, uc := fixture(t)
	if _, err := Init(root, uc, "generate"); err != nil {
		t.Fatal(err)
	}
	if _, err := UpdateMappings(root, uc, []Mapping{{Requirement: "REQ-UX-001", Coverage: "missing"}}); err != nil {
		t.Fatal(err)
	}
	inspection, err := Inspect(root, uc)
	if err != nil {
		t.Fatal(err)
	}
	if len(inspection.Blockers) != 1 {
		t.Fatalf("expected blocker, got %+v", inspection)
	}
}
