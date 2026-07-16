package extensions

import "testing"

func TestManifestRejectsUnknownCapability(t *testing.T) {
	if err := (Manifest{ID: "review-import", Version: "0.1.0", Capabilities: []string{"reviews.import"}}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Manifest{ID: "bad", Version: "0.1.0", Capabilities: []string{"artifacts.write"}}).Validate(); err == nil {
		t.Fatal("unknown capability accepted")
	}
}
