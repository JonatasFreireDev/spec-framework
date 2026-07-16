package sourceimport

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func materializableSources(productRoot, runID string) (map[string]bool, bool, error) {
	runRoot := filepath.Join(productRoot, "knowledge", "imports", "runs", runID)
	if _, err := os.Stat(filepath.Join(runRoot, "inventory", "index.json")); err == nil {
		sources, err := scalableSources(productRoot, runID)
		if err != nil {
			return nil, true, err
		}
		known := map[string]bool{}
		for _, source := range sources {
			known[source.Path] = true
		}
		return known, true, nil
	}
	data, err := os.ReadFile(filepath.Join(runRoot, "inventory.json"))
	if err != nil {
		return nil, false, err
	}
	var inventory Inventory
	if err := json.Unmarshal(trimBOM(data), &inventory); err != nil {
		return nil, false, err
	}
	known := map[string]bool{}
	for _, source := range inventory.Sources {
		known[source.Path] = true
	}
	return known, false, nil
}

func requireReviewedChunks(productRoot, runID string) error {
	entries, err := os.ReadDir(filepath.Join(productRoot, "knowledge", "imports", "runs", runID, "chunks"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		var chunk Chunk
		if err := readJSONFile(filepath.Join(productRoot, "knowledge", "imports", "runs", runID, "chunks", entry.Name()), &chunk); err != nil {
			return err
		}
		if chunk.Status != "reviewed" && chunk.Status != "excluded" {
			return errors.New("all scalable import chunks must be reviewed or excluded before materialization")
		}
	}
	return nil
}
