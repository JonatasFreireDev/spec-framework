package designsystem

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitAndInspect(t *testing.T) {
	root := t.TempDir()
	if _, err := Init(root, "generate"); err != nil {
		t.Fatal(err)
	}
	i, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if i.ID != "DSYS-001" || i.Version != "0.1.0" || len(i.Blockers) != 0 {
		t.Fatalf("unexpected inspection: %+v", i)
	}
}

func TestTokenAliasValidation(t *testing.T) {
	root := t.TempDir()
	if _, err := Init(root, "generate"); err != nil {
		t.Fatal(err)
	}
	doc := TokenDocument{SchemaVersion: 1, System: "DSYS-001", Version: "0.1.0", Tokens: map[string]any{
		"color": map[string]any{
			"primitive": map[string]any{"blue": map[string]any{"value": "#00f", "type": "color"}},
			"semantic":  map[string]any{"action": map[string]any{"value": "{color.primitive.blue}", "type": "color"}},
		},
	}}
	data, _ := json.Marshal(doc)
	if err := os.WriteFile(filepath.Join(root, "design", "system", "tokens", "tokens.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	i, err := Inspect(root)
	if err != nil || i.Tokens != 2 || len(i.Blockers) != 0 {
		t.Fatalf("unexpected validation: %+v %v", i, err)
	}
}

func TestBrokenAliasBlocks(t *testing.T) {
	root := t.TempDir()
	if _, err := Init(root, "generate"); err != nil {
		t.Fatal(err)
	}
	doc := TokenDocument{SchemaVersion: 1, System: "DSYS-001", Version: "0.1.0", Tokens: map[string]any{"x": map[string]any{"value": "{missing}", "type": "color"}}}
	data, _ := json.Marshal(doc)
	_ = os.WriteFile(filepath.Join(root, "design", "system", "tokens", "tokens.json"), data, 0o644)
	i, _ := Inspect(root)
	if len(i.Blockers) == 0 {
		t.Fatal("expected broken alias blocker")
	}
}

func TestMigrateDryRun(t *testing.T) {
	root := t.TempDir()
	items, err := Migrate(root, true)
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected dry run: %v %v", items, err)
	}
	if _, err := os.Stat(filepath.Join(root, "design", "system")); !os.IsNotExist(err) {
		t.Fatal("dry run mutated the product")
	}
}
