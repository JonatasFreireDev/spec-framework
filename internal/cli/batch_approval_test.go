package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func TestApproveBatchPreviewAndApplyByID(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := workflow.Registry{Artifacts: []workflow.Artifact{{ID: "PRB-001", Type: "problem", Status: "draft", Path: "foundation/problem/problem.md"}}}
	data, _ := json.Marshal(registry)
	if err := os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "foundation", "problem", "problem.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("status: draft\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var out, errout bytes.Buffer
	code := New("test").Run([]string{"approve-batch", "--product-root", root, "--ids", "PRB-001"}, &out, &errout)
	if code != 0 || !strings.Contains(out.String(), "APPROVE PRB-001") {
		t.Fatalf("preview=%d out=%s err=%s", code, out.String(), errout.String())
	}
	out.Reset()
	errout.Reset()
	code = New("test").Run([]string{"approve-batch", "--product-root", root, "--ids", "PRB-001", "--approved-by", "Owner", "--yes"}, &out, &errout)
	if code != 0 || !strings.Contains(out.String(), "APPROVED PRB-001") {
		t.Fatalf("apply=%d out=%s err=%s", code, out.String(), errout.String())
	}
}
