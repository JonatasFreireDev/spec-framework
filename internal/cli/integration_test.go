package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/cli"
	"github.com/JonatasFreireDev/spec-framework/internal/sourceimport"
)

func reviewInitialImportChunk(t *testing.T, productRoot string) {
	t.Helper()
	chunk, err := sourceimport.Resume(productRoot, "IMPORT-001", "CHUNK-0001", "test-importer")
	if err != nil {
		t.Fatal(err)
	}
	if err := sourceimport.RecordChunkReview(productRoot, "IMPORT-001", chunk.ID, "test-importer", sourceimport.ChunkReview{SourceEvidence: map[string][]sourceimport.Evidence{"SRC-000001": {{Locator: "line 1", Claim: "fixture reviewed"}}}}); err != nil {
		t.Fatal(err)
	}
}

func TestGoCLIInitValidateUpgradeAndMove(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	parent := t.TempDir()
	target := filepath.Join(parent, "product")
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex,cursor,claude", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stderr=%s", code, stderr.String())
	}
	for _, path := range []string{"product/.product/framework.json", "product/BOOTSTRAP.md", "product/tools/check-links.py"} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); err != nil {
			t.Fatal(err)
		}
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate"}, &stdout, &stderr); code == 0 {
		t.Fatal("validate outside the product repository should not activate")
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
		t.Fatalf("absolute validate=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate"}, &stdout, &stderr); code != 0 {
		t.Fatalf("validate=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"work", "--feature", "FT-TEMPLATE", "--created-by", "Test Owner"}, &stdout, &stderr); code != 0 {
		t.Fatalf("work=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"guide", "--work", "WORK-001"}, &stdout, &stderr); code != 1 {
		t.Fatalf("blocked guide=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("Feature scope: domains/_template-domain/goals/_template-goal/features/_template-feature/context.md")) || !bytes.Contains(stdout.Bytes(), []byte("Skill: feature")) {
		t.Fatalf("guide omitted verified route context: %s", stdout.String())
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

func TestCLIStructuredEngineeringBaselineValidatesWithoutNormalization(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "repo")
	productRoot := filepath.Join(target, "product")
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex", "--no-code-roots", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}

	// Simulate Technical Landscape materialization using the shipped indexed
	// catalog and entity contracts.
	files := map[string]string{
		"engineering/catalog/catalog.yaml":                    "schema_version: 1\nowner_skill: technical-landscape\nentities:\n  systems:\n    SYS-MF-001: systems/meetfriends.yaml\n  applications: {}\n  components: {}\n  repositories: {}\n  data_stores: {}\n  interfaces: {}\n  deployments: {}\nrelations: []\n",
		"engineering/catalog/systems/meetfriends.yaml":        "schema_version: 1\nid: SYS-MF-001\ntype: system\nstatus: draft\nname: MeetFriends\nevidence: []\n",
		"engineering/architecture/topology.yaml":              "schema_version: 1\nowner_skill: technical-landscape\nsystems: [SYS-MF-001]\napplications: []\ncomponents: []\nrepositories: []\ndata_stores: []\ninterfaces: []\ndeployments: []\nrelations: []\n",
		"engineering/standards/standards.yaml":                "schema_version: 1\nowner_skill: engineering-standards\nprofiles:\n  PROFILE-PRODUCT-DEFAULT: profiles/product-default.yaml\nstandards:\n  STD-API-001: catalog/api.yaml\nexceptions: {}\n",
		"engineering/standards/profiles/product-default.yaml": "schema_version: 1\nid: PROFILE-PRODUCT-DEFAULT\nversion: 0.1.0\nstatus: draft\nextends: []\nstandards: [STD-API-001]\n",
		"engineering/standards/catalog/api.yaml":              "schema_version: 1\nid: STD-API-001\nversion: 1.0.0\nstatus: draft\ncategory: api\nlevel: required\nrules:\n  - id: STD-API-001-R01\n    requirement: Publish an API schema\n    verification: [schema]\n",
		"engineering/operations/operations.yaml":              "schema_version: 1\nowner_skill: operations-baseline\nenvironments:\n  ENV-PRODUCTION: environments/production.yaml\ndeployments: {}\nrunbooks: {}\n",
		"engineering/operations/environments/production.yaml": "schema_version: 1\nid: ENV-PRODUCTION\nstatus: draft\npurpose: production\n",
	}
	for name, body := range files {
		path := filepath.Join(productRoot, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate", "--product-root", productRoot}, &stdout, &stderr); code != 0 {
		t.Fatalf("validate=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}

	embedded := "schema_version: 1\nowner_skill: technical-landscape\nentities:\n  systems:\n    SYS-MF-001:\n      name: MeetFriends\n      evidence: []\nrelations: []\n"
	if err := os.WriteFile(filepath.Join(productRoot, "engineering", "catalog", "catalog.yaml"), []byte(embedded), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate", "--product-root", productRoot}, &stdout, &stderr); code == 0 {
		t.Fatalf("embedded catalog unexpectedly validated: stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	for _, expected := range []string{"systems id SYS-MF-001", "expected a relative YAML file path", "not an embedded entity"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("CLI diagnostic omitted %q: stdout=%s", expected, stdout.String())
		}
	}
}

func TestCLIAgentLedCodeRootDiscovery(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	root := t.TempDir()
	target := filepath.Join(root, "repo")
	if err := os.MkdirAll(filepath.Join(target, "web"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "web", "package.json"), []byte(`{"name":"web"}`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex", "--code-roots", "web:frontend", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	manifestPath := filepath.Join(target, "product", ".product", "framework.json")
	manifest, err := os.ReadFile(manifestPath)
	if err != nil || !bytes.Contains(manifest, []byte(`"mode": "agent-declared"`)) || !bytes.Contains(manifest, []byte(`"role": "frontend"`)) {
		t.Fatalf("agent discovery missing from manifest: %s %v", manifest, err)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("agent-declared (confirmed)")) {
		t.Fatalf("init did not report discovery authority: %s", stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"upgrade", "--target", target, "--code-roots", "web:web-client", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("upgrade=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	manifest, err = os.ReadFile(manifestPath)
	if err != nil || !bytes.Contains(manifest, []byte(`"role": "web-client"`)) || !bytes.Contains(stdout.Bytes(), []byte("agent-declared (confirmed)")) {
		t.Fatalf("upgrade did not replace confirmed roots: %s stdout=%s err=%v", manifest, stdout.String(), err)
	}

	noCode := filepath.Join(root, "no-code")
	if err := os.MkdirAll(filepath.Join(noCode, "tooling"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(noCode, "tooling", "package.json"), []byte(`{"private":true}`), 0644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"init", "--target", noCode, "--agents", "codex", "--no-code-roots", "--yes"}, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "agent-confirmed-none (confirmed)") {
		t.Fatalf("confirmed no-code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	manifest, err = os.ReadFile(filepath.Join(noCode, "product", ".product", "framework.json"))
	if err != nil || bytes.Contains(manifest, []byte(`"path": "tooling"`)) {
		t.Fatalf("CLI candidate overrode agent no-code decision: %s %v", manifest, err)
	}
}

func TestCLIFoundationApprovalsCreateTraceableHistory(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "repo")
	productRoot := filepath.Join(target, "product")
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stderr=%s", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate", "--product-root", productRoot, "--write-registry"}, &stdout, &stderr); code != 0 {
		t.Fatalf("registry regeneration=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	artifacts := []string{
		"foundation/problem/problem.md",
		"foundation/vision/vision.md",
		"foundation/vision/principles.md",
		"foundation/vision/north-star.md",
		"foundation/strategy/strategy.md",
	}
	for _, artifact := range artifacts {
		stdout.Reset()
		stderr.Reset()
		code := app.Run([]string{"approve", "--product-root", productRoot, "--artifact", artifact, "--grant", "approved", "--approved-by", "Test Owner", "--yes"}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("approve %s=%d stdout=%s stderr=%s", artifact, code, stdout.String(), stderr.String())
		}
	}
	history, err := filepath.Glob(filepath.Join(productRoot, ".product", "history", "approval-*.json"))
	if err != nil || len(history) != len(artifacts) {
		t.Fatalf("history=%v err=%v", history, err)
	}
	for _, contextPath := range []string{"foundation/problem/context.md", "foundation/vision/context.md", "foundation/strategy/context.md"} {
		data, err := os.ReadFile(filepath.Join(productRoot, filepath.FromSlash(contextPath)))
		if err != nil || !bytes.Contains(data, []byte("status: approved")) {
			t.Fatalf("context %s not synchronized: %s err=%v", contextPath, data, err)
		}
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"validate", "--product-root", productRoot}, &stdout, &stderr); code != 0 {
		t.Fatalf("validate=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
}

func TestCLIExistingDocumentsMaterialization(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
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
	reviewInitialImportChunk(t, productRoot)
	if code := app.Run([]string{"import", "materialize", "--product-root", productRoot, "--run", "IMPORT-001", "--approved-by", "Product Owner", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("materialize=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(filepath.Join(productRoot, "domains", "payments", "domain.md")); err != nil {
		t.Fatal(err)
	}
}

func TestCLIExistingDocumentsNormalizeApproveAndWork(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	root := t.TempDir()
	target := filepath.Join(root, "repo")
	source := filepath.Join(root, "brief.md")
	if err := os.WriteFile(source, []byte("# FocusFlow"), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	app := cli.New("integration")
	if code := app.Run([]string{"init", "--target", target, "--agents", "codex", "--starting-point", "existing-documents", "--sources", source, "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("init=%d stderr=%s", code, stderr.String())
	}
	productRoot := filepath.Join(target, "product")
	runRoot := filepath.Join(productRoot, "knowledge", "imports", "runs", "IMPORT-001")
	invData, err := os.ReadFile(filepath.Join(runRoot, "inventory.json"))
	if err != nil {
		t.Fatal(err)
	}
	var inv map[string]any
	if err := json.Unmarshal(invData, &inv); err != nil {
		t.Fatal(err)
	}
	sourceRel := inv["sources"].([]any)[0].(map[string]any)["path"].(string)
	reviewInitialImportChunk(t, productRoot)
	featurePath := "domains/imported/goals/imported/features/imported/context.md"
	draft := "---\nid: FT-IMPORT-001\ntype: feature\nname: Imported Feature\nstatus: draft\nowner_skill: feature\nslug: imported\nrigor_tier: S\nsource_documents:\n  - " + sourceRel + "\n---\n\n# Imported Feature\n"
	mapping := map[string]any{"schema_version": 1, "import_id": "IMPORT-001", "mappings": []any{map[string]any{"id": "MAP-001", "target": featurePath, "artifact_type": "feature", "selected": true, "source_documents": []string{sourceRel}, "draft_content": draft}}}
	data, _ := json.MarshalIndent(mapping, "", "  ")
	if err := os.WriteFile(filepath.Join(runRoot, "mapping.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"import", "materialize", "--product-root", productRoot, "--run", "IMPORT-001", "--approved-by", "Product Owner", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("materialize=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	registryPath := filepath.Join(productRoot, ".product", "artifacts.json")
	registryData, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	var registry map[string]any
	if err := json.Unmarshal(registryData, &registry); err != nil {
		t.Fatal(err)
	}
	artifacts := registry["artifacts"].([]any)
	artifacts = append(artifacts, map[string]any{"id": "FT-IMPORT-001", "type": "feature", "status": "draft", "path": featurePath, "parentIds": []string{}})
	registry["artifacts"] = artifacts
	registryData, _ = json.MarshalIndent(registry, "", "  ")
	if err := os.WriteFile(registryPath, append(registryData, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"approve", "--product-root", productRoot, "--artifact", featurePath, "--grant", "approved", "--approved-by", "Product Owner", "--yes"}, &stdout, &stderr); code == 0 {
		t.Fatal("import-draft approval unexpectedly succeeded")
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"template", "audit", "--product-root", productRoot, "--framework-root", filepath.Join(target, ".spec-framework"), "--artifact", featurePath}, &stdout, &stderr); code != 0 {
		t.Fatalf("template audit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"template", "normalize", "--product-root", productRoot, "--framework-root", filepath.Join(target, ".spec-framework"), "--artifact", featurePath, "--skill", "feature"}, &stdout, &stderr); code != 0 {
		t.Fatalf("template normalize=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"approve", "--product-root", productRoot, "--artifact", featurePath, "--grant", "approved", "--approved-by", "Product Owner", "--yes"}, &stdout, &stderr); code != 0 {
		t.Fatalf("approve normalized=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"work", "--product-root", productRoot, "--feature", "FT-IMPORT-001", "--created-by", "Product Owner"}, &stdout, &stderr); code != 0 {
		t.Fatalf("work=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
}

func TestCLIWorkspaceApprovalGatesAndGraph(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
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
	if code := app.Run([]string{"validate", "--product-root", product}, &stdout, &stderr); code != 0 {
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
