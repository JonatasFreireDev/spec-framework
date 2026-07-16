package extensions

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Manifest declares an optional adapter. It grants no authority by itself.
type Manifest struct {
	ID           string   `json:"id"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

// Discover reads only versioned manifests supplied by the framework or an
// adopter. Discovery grants no capability and never executes extension code.
func Discover(dir string) ([]Manifest, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []Manifest{}, nil
	}
	if err != nil {
		return nil, err
	}
	var manifests []Manifest
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, err
		}
		if err := manifest.Validate(); err != nil {
			return nil, err
		}
		manifests = append(manifests, manifest)
	}
	sort.Slice(manifests, func(i, j int) bool { return manifests[i].ID < manifests[j].ID })
	return manifests, nil
}

// EnabledCapability requires an explicit product-owned enablement record.
// A manifest alone is only declarative metadata.
func EnabledCapability(productRoot string, manifest Manifest, capability string) (bool, error) {
	if err := manifest.Validate(); err != nil {
		return false, err
	}
	declared := false
	for _, item := range manifest.Capabilities {
		declared = declared || item == capability
	}
	if !declared {
		return false, nil
	}
	path := filepath.Join(productRoot, ".product", "extensions", manifest.ID+".json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var enabled struct {
		Version      string   `json:"version"`
		Capabilities []string `json:"capabilities"`
	}
	if err := json.Unmarshal(data, &enabled); err != nil {
		return false, err
	}
	if enabled.Version != manifest.Version {
		return false, nil
	}
	for _, item := range enabled.Capabilities {
		if item == capability {
			return true, nil
		}
	}
	return false, nil
}

func (m Manifest) Validate() error {
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.Version) == "" {
		return errors.New("extension id and version are required")
	}
	allowed := map[string]bool{"artifacts.read": true, "reviews.import": true, "runtime.dispatch": true, "events.read": true}
	for _, capability := range m.Capabilities {
		if !allowed[capability] {
			return errors.New("unsupported extension capability: " + capability)
		}
	}
	return nil
}
