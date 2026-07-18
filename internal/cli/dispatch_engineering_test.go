package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/dispatch"
)

func TestDispatchEngineeringAssignAndReturn(t *testing.T) {
	root := t.TempDir()
	work := "WORK-ENG-001"
	handoff := filepath.Join(root, ".product", "workspaces", work, "engineering-handoff.json")
	if err := os.MkdirAll(filepath.Dir(handoff), 0755); err != nil {
		t.Fatal(err)
	}
	contract := []byte(`{"schema_version":1,"execution":{"mode":"delegated","context_policy":"minimal","max_parallel":1,"fallback":"sequential"},"routes":[{"skill":"technical-landscape","phase":1,"depends_on":[],"write_scope":["engineering/catalog"],"status":"pending"}]}`)
	if err := os.WriteFile(handoff, contract, 0644); err != nil {
		t.Fatal(err)
	}
	var out, errout bytes.Buffer
	code := runDispatch([]string{"assign", "--product-root", root, "--work", work, "--task", ".product/workspaces/" + work + "/engineering-handoff.json", "--role", "technical-landscape", "--agent", "landscape-1", "--yes"}, &out, &errout)
	if code != 0 {
		t.Fatalf("assign code=%d out=%q err=%q", code, out.String(), errout.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != 2 || fields[0] != "ASSIGNED" {
		t.Fatalf("unexpected assign output %q", out.String())
	}
	id := fields[1]
	outputPath := filepath.Join(root, "engineering", "catalog", "catalog.yaml")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		t.Fatal(err)
	}
	content := []byte("schema_version: 1\n")
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	out.Reset()
	errout.Reset()
	code = runDispatch([]string{"return", "--product-root", root, "--work", work, "--id", id, "--agent", "landscape-1", "--summary", "mapped", "--evidence", "catalog", "--output-hashes", "engineering/catalog/catalog.yaml=" + hex.EncodeToString(sum[:]), "--decision-candidates", "DEC-ENG-001", "--yes"}, &out, &errout)
	if code != 0 {
		t.Fatalf("return code=%d out=%q err=%q", code, out.String(), errout.String())
	}
	data, err := os.ReadFile(filepath.Join(root, ".product", "workspaces", work, "dispatches", id+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var envelope dispatch.Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Status != "returned" || len(envelope.DecisionCandidates) != 1 || envelope.DecisionCandidates[0] != "DEC-ENG-001" {
		t.Fatalf("returned envelope=%+v", envelope)
	}
}
