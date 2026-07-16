package extensions

import (
	"errors"
	"strings"
)

// Manifest declares an optional adapter. It grants no authority by itself.
type Manifest struct {
	ID           string   `json:"id"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
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
