package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/cli"
)

func TestGoCLIInitValidateUpgradeAndMove(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "product")
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex,cursor,claude", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stderr=%s", code, stderr.String())
	}
	for _, path := range []string{".agents/skills/code-runner/SKILL.md", ".cursor/skills/code-runner/SKILL.md", ".claude/skills/code-runner/SKILL.md", ".spec-framework/manifest.json"} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); err != nil {
			t.Fatal(err)
		}
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate", "--product-root", filepath.Join(target, "product"), "--framework-root", filepath.Join(target, ".spec-framework")}, &stdout, &stderr); code != 0 {
		t.Fatalf("absolute validate=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chdir(target); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate"}, &stdout, &stderr); code != 0 {
		t.Fatalf("validate=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	moveSource := filepath.Join(target, "product", "move-source")
	if err = os.MkdirAll(moveSource, 0755); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(moveSource, "context.md"), []byte("# Move\n"), 0644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"move", "--from", "product/move-source", "--to", "product/moved"}, &stdout, &stderr); code != 0 {
		t.Fatalf("move=%d stderr=%s", code, stderr.String())
	}
	if _, err = os.Stat(filepath.Join(target, "product", "moved", "context.md")); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"upgrade", "--target", target, "--agents", "codex,cursor,claude", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("upgrade=%d stderr=%s", code, stderr.String())
	}
}

func TestCLIExistingDocumentsMaterialization(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "repo")
	source := filepath.Join(root, "epic.md")
	if err := os.WriteFile(source, []byte("# Payments"), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex", "--starting-point", "existing-documents", "--sources", source, "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stderr=%s", code, stderr.String())
	}
	runRoot := filepath.Join(target, "product", "knowledge", "imports", "runs", "IMPORT-001")
	invData, _ := os.ReadFile(filepath.Join(runRoot, "inventory.json"))
	var inv map[string]any
	if err := json.Unmarshal(invData, &inv); err != nil {
		t.Fatal(err)
	}
	sources := inv["sources"].([]any)
	sourceRel := sources[0].(map[string]any)["path"].(string)
	mapping := map[string]any{"schema_version": 1, "import_id": "IMPORT-001", "mappings": []any{map[string]any{"id": "MAP-001", "target": "domains/payments/domain.md", "artifact_type": "domain", "selected": true, "source_documents": []string{sourceRel}, "draft_content": "---\nstatus: draft\nsource_documents:\n  - " + sourceRel + "\n---\n# Payments\n"}}}
	data, _ := json.MarshalIndent(mapping, "", "  ")
	if err := os.WriteFile(filepath.Join(runRoot, "mapping.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	productRoot := filepath.Join(target, "product")
	if code := app.Run([]string{"import", "materialize", "--product-root", productRoot, "--run", "IMPORT-001", "--approved-by", "Product Owner", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("materialize=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(filepath.Join(productRoot, "domains", "payments", "domain.md")); err != nil {
		t.Fatal(err)
	}
}

func TestCLIWorkspaceApprovalGatesAndGraph(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "repo")
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d %s", code, stderr.String())
	}
	product := filepath.Join(target, "product")
	approve := func(path string) {
		stdout.Reset()
		stderr.Reset()
		if code := app.Run([]string{"approve", "--product-root", product, "--artifact", path, "--grant", "approved", "--approved-by", "Test Owner", "--yes"}, &stdout, &stderr); code != 0 {
			t.Fatalf("approve %s=%d %s", path, code, stderr.String())
		}
	}
	approve("domains/_template-domain/context.md")
	approve("domains/_template-domain/goals/_template-goal/context.md")
	approve("domains/_template-domain/goals/_template-goal/features/_template-feature/context.md")
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"work", "--product-root", product, "--feature", "FT-TEMPLATE", "--created-by", "Test Owner"}, &stdout, &stderr); code != 0 {
		t.Fatalf("work=%d %s", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"status", "--product-root", product, "--work", "WORK-001"}, &stdout, &stderr); code != 0 {
		t.Fatalf("status=%d %s %s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate", "--product-root", product, "--framework-root", filepath.Join(target, ".spec-framework")}, &stdout, &stderr); code != 0 {
		t.Fatalf("validate after approvals=%d %s %s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"gates", "--product-root", product}, &stdout, &stderr); code == 0 {
		t.Fatal("placeholder gates should block")
	}
	graph := "domains/_template-domain/goals/_template-goal/features/_template-feature/use-cases/_template-use-case/execution-graph.json"
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"graph", "ready", "--product-root", product, "--graph", graph}, &stdout, &stderr); code != 0 {
		t.Fatalf("ready=%d %s", code, stderr.String())
	}
	if code := app.Run([]string{"graph", "claim", "--product-root", product, "--graph", graph, "--task", "TASK-TEMPLATE-001", "--agent", "codex"}, &stdout, &stderr); code != 0 {
		t.Fatalf("claim=%d %s", code, stderr.String())
	}
}
