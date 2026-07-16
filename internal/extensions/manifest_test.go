package extensions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManifestRejectsUnknownCapability(t *testing.T) {
	if err := (Manifest{ID: "review-import", Version: "0.1.0", Capabilities: []string{"reviews.import"}}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Manifest{ID: "bad", Version: "0.1.0", Capabilities: []string{"artifacts.write"}}).Validate(); err == nil {
		t.Fatal("unknown capability accepted")
	}
}

func TestDiscoveryAndProductEnablementAreSeparate(t *testing.T) {
	root := t.TempDir()
	manifests := filepath.Join(root, "manifests")
	if err := os.MkdirAll(manifests, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifests, "review.json"), []byte(`{"id":"review-import","version":"1","capabilities":["reviews.import"]}`), 0644); err != nil {
		t.Fatal(err)
	}
	found, err := Discover(manifests)
	if err != nil || len(found) != 1 {
		t.Fatalf("found=%+v err=%v", found, err)
	}
	if enabled, err := EnabledCapability(root, found[0], "reviews.import"); err != nil || enabled {
		t.Fatalf("enabled=%v err=%v", enabled, err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".product", "extensions"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".product", "extensions", "review-import.json"), []byte(`{"version":"1","capabilities":["reviews.import"]}`), 0644); err != nil {
		t.Fatal(err)
	}
	if enabled, err := EnabledCapability(root, found[0], "reviews.import"); err != nil || !enabled {
		t.Fatalf("enabled=%v err=%v", enabled, err)
	}
}
