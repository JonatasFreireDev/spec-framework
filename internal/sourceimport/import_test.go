package sourceimport

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateRunInventoriesMultipleSourcesWithoutMaterializingArtifacts(t *testing.T) {
	root := t.TempDir()
	sources := filepath.Join(t.TempDir(), "docs")
	if err := os.MkdirAll(sources, 0755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"payments.md": "# Payments", "events.md": "# Events"} {
		if err := os.WriteFile(filepath.Join(sources, name), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	runID, err := CreateRun(root, []string{sources})
	if err != nil {
		t.Fatal(err)
	}
	if runID != "IMPORT-001" {
		t.Fatalf("run=%s", runID)
	}
	data, err := os.ReadFile(filepath.Join(root, "knowledge", "imports", "runs", runID, "inventory.json"))
	if err != nil {
		t.Fatal(err)
	}
	var inv Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		t.Fatal(err)
	}
	if len(inv.Sources) != 2 {
		t.Fatalf("sources=%d", len(inv.Sources))
	}
	if _, err := os.Stat(filepath.Join(root, "domains")); !os.IsNotExist(err) {
		t.Fatal("import unexpectedly materialized domains")
	}
	traceData, err := os.ReadFile(filepath.Join(root, "knowledge", "imports", "runs", runID, "traceability.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trace Traceability
	if err := json.Unmarshal(traceData, &trace); err != nil {
		t.Fatal(err)
	}
	if trace.Status != "unreviewed" || len(trace.Sources) != 2 || trace.Sources[0].ReviewStatus != "unreviewed" {
		t.Fatalf("unexpected traceability: %+v", trace)
	}
}

func TestCreateRunRequiresSources(t *testing.T) {
	if _, err := CreateRun(t.TempDir(), nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestMaterializeRequiresExplicitApprovalAndCreatesDraft(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(t.TempDir(), "epic.md")
	if err := os.WriteFile(source, []byte("# Epic"), 0644); err != nil {
		t.Fatal(err)
	}
	runID, err := CreateRun(root, []string{source})
	if err != nil {
		t.Fatal(err)
	}
	invData, _ := os.ReadFile(filepath.Join(root, "knowledge", "imports", "runs", runID, "inventory.json"))
	var inv Inventory
	_ = json.Unmarshal(invData, &inv)
	mapping := MappingFile{SchemaVersion: 1, ImportID: runID, Mappings: []Mapping{{ID: "MAP-001", Target: "domains/payments/domain.md", ArtifactType: "domain", Selected: true, SourceDocuments: []string{inv.Sources[0].Path}, DraftContent: "---\nstatus: draft\nsource_documents:\n  - " + inv.Sources[0].Path + "\n---\n# Payments\n"}}}
	if err := writeJSON(filepath.Join(root, "knowledge", "imports", "runs", runID, "mapping.json"), mapping); err != nil {
		t.Fatal(err)
	}
	if _, err := Materialize(root, runID, ""); err == nil {
		t.Fatal("expected approval identity error")
	}
	created, err := Materialize(root, runID, "Jonatas")
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 {
		t.Fatalf("created=%v", created)
	}
	if _, err := os.Stat(filepath.Join(root, "domains", "payments", "domain.md")); err != nil {
		t.Fatal(err)
	}
	planData, err := os.ReadFile(filepath.Join(root, "knowledge", "imports", "runs", runID, "import-plan.json"))
	if err != nil {
		t.Fatal(err)
	}
	var plan map[string]any
	if err := json.Unmarshal(planData, &plan); err != nil {
		t.Fatal(err)
	}
	hashes, _ := plan["materialized_hashes"].(map[string]any)
	if hashes["domains/payments/domain.md"] == "" {
		t.Fatalf("materialized draft hash missing: %v", plan)
	}
	traceData, err := os.ReadFile(filepath.Join(root, "knowledge", "imports", "runs", runID, "traceability.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trace Traceability
	if err := json.Unmarshal(traceData, &trace); err != nil {
		t.Fatal(err)
	}
	if trace.Status != "materialized_as_draft" || len(trace.Sources[0].MaterializedPaths) != 1 || trace.Sources[0].MaterializedPaths[0] != "domains/payments/domain.md" {
		t.Fatalf("materialization not reflected in traceability: %+v", trace)
	}
	if _, err := Materialize(root, runID, "Jonatas"); err == nil {
		t.Fatal("expected overwrite protection")
	}
}

func TestMaterializeAcceptsUTF8BOMJSON(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(t.TempDir(), "epic.md")
	if err := os.WriteFile(source, []byte("# Epic"), 0644); err != nil {
		t.Fatal(err)
	}
	runID, err := CreateRun(root, []string{source})
	if err != nil {
		t.Fatal(err)
	}
	runRoot := filepath.Join(root, "knowledge", "imports", "runs", runID)
	invData, _ := os.ReadFile(filepath.Join(runRoot, "inventory.json"))
	var inv Inventory
	_ = json.Unmarshal(invData, &inv)
	mapping := MappingFile{SchemaVersion: 1, ImportID: runID, Mappings: []Mapping{{ID: "MAP-001", Target: "domains/bom/domain.md", Selected: true, SourceDocuments: []string{inv.Sources[0].Path}, DraftContent: "status: draft\nsource_documents:\n  - " + inv.Sources[0].Path + "\n"}}}
	data, _ := json.Marshal(mapping)
	data = append([]byte{0xef, 0xbb, 0xbf}, data...)
	if err := os.WriteFile(filepath.Join(runRoot, "mapping.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Materialize(root, runID, "Windows User"); err != nil {
		t.Fatal(err)
	}
}

func TestMaterializeRejectsEscapingAndDuplicateTargets(t *testing.T) {
	for name, mappings := range map[string][]Mapping{
		"escape":    {{ID: "MAP-001", Target: "../outside.md", Selected: true, SourceDocuments: []string{"knowledge/imports/sources/a.md"}, DraftContent: "status: draft\nsource_documents: []"}},
		"duplicate": {{ID: "MAP-001", Target: "domains/a/domain.md", Selected: true, SourceDocuments: []string{"a"}, DraftContent: "status: draft\nsource_documents: []"}, {ID: "MAP-002", Target: "domains/a/domain.md", Selected: true, SourceDocuments: []string{"a"}, DraftContent: "status: draft\nsource_documents: []"}},
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			runID := "IMPORT-001"
			runRoot := filepath.Join(root, "knowledge", "imports", "runs", runID)
			if err := os.MkdirAll(runRoot, 0755); err != nil {
				t.Fatal(err)
			}
			if err := writeJSON(filepath.Join(runRoot, "mapping.json"), MappingFile{SchemaVersion: 1, ImportID: runID, Mappings: mappings}); err != nil {
				t.Fatal(err)
			}
			if _, err := Materialize(root, runID, "Jonatas"); err == nil {
				t.Fatal("expected rejection")
			}
		})
	}
}
