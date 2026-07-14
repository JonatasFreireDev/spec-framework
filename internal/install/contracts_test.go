package install

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
)

func TestEveryStartingPointHasAValidInitializationPlan(t *testing.T) {
	for _, point := range StartingPoints {
		t.Run(point, func(t *testing.T) {
			plan, err := buildInitializationPlan(point)
			if err != nil {
				t.Fatal(err)
			}
			requiredFiles := []string{".product/artifacts.json", ".product/framework.json", "context.md", "tools/check-links.py"}
			if point != "existing-documents" && point != "audit-only" {
				requiredFiles = append(requiredFiles, "domains/_template-domain/context.md")
			}
			for _, required := range requiredFiles {
				if _, ok := plan.Files[required]; !ok {
					t.Fatalf("plan is missing %s", required)
				}
			}
		})
	}
}

func TestStartingPointMaterializationProfiles(t *testing.T) {
	tests := map[string]struct {
		present []string
		absent  []string
	}{
		"new-product": {
			present: []string{"foundation/problem/problem.md", "foundation/vision/vision.md", "domains/_template-domain/context.md", "design/system/context.md", "engineering/context.md"},
		},
		"existing-product": {
			present: []string{"foundation/strategy/strategy.md", "foundation/product-baseline.md", "domains/_template-domain/context.md"},
			absent:  []string{"foundation/problem/problem.md", "foundation/vision/vision.md"},
		},
		"existing-feature": {
			present: []string{"foundation/feature-brief.md", "domains/_template-domain/context.md"},
			absent:  []string{"foundation/problem/problem.md", "foundation/strategy/strategy.md", "design/system/context.md", "engineering/context.md"},
		},
		"existing-documents": {
			present: []string{"context.md"},
			absent:  []string{"foundation/README.md", "foundation/problem/problem.md", "domains/_template-domain/context.md", "design/system/context.md", "engineering/context.md"},
		},
		"existing-implementation": {
			present: []string{"knowledge/assessments/implementation-assessment.md", "foundation/problem/problem.md", "domains/_template-domain/context.md"},
		},
		"audit-only": {
			present: []string{"audits/README.md", "knowledge/conventions/security-baseline.md"},
			absent:  []string{"foundation/README.md", "foundation/problem/problem.md", "knowledge/conventions/gates.md", "domains/_template-domain/context.md", "design/system/context.md", "engineering/context.md"},
		},
	}
	for point, expected := range tests {
		t.Run(point, func(t *testing.T) {
			plan, err := buildInitializationPlan(point)
			if err != nil {
				t.Fatal(err)
			}
			for _, path := range expected.present {
				if _, ok := plan.Files[path]; !ok {
					t.Fatalf("expected %s to be materialized", path)
				}
			}
			for _, path := range expected.absent {
				if _, ok := plan.Files[path]; ok {
					t.Fatalf("did not expect %s to be materialized", path)
				}
			}
			if point == "existing-documents" {
				for _, directory := range []string{"knowledge/imports/sources", "knowledge/imports/runs"} {
					found := false
					for _, planned := range plan.Directories {
						if planned == directory {
							found = true
							break
						}
					}
					if !found {
						t.Fatalf("expected %s directory to be materialized", directory)
					}
				}
			}
		})
	}
}

func TestInitializationPlansAreDeterministic(t *testing.T) {
	for _, point := range StartingPoints {
		first, err := buildInitializationPlan(point)
		if err != nil {
			t.Fatal(err)
		}
		second, err := buildInitializationPlan(point)
		if err != nil {
			t.Fatal(err)
		}
		if planDigest(first) != planDigest(second) {
			t.Fatalf("%s plan is not deterministic", point)
		}
	}
}

func TestStartingPointPlansMatchGoldenContracts(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "entrypoints.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden map[string]struct {
		ArtifactIDs  []string `json:"artifact_ids"`
		SpecialFiles []string `json:"special_files"`
		Directories  []string `json:"directories"`
		Actions      []string `json:"actions"`
	}
	if err := json.Unmarshal(data, &golden); err != nil {
		t.Fatal(err)
	}
	for _, point := range StartingPoints {
		plan, err := buildInitializationPlan(point)
		if err != nil {
			t.Fatal(err)
		}
		expected, ok := golden[point]
		if !ok {
			t.Fatalf("golden contract missing %s", point)
		}
		var registry struct {
			Artifacts []map[string]any `json:"artifacts"`
		}
		if err := json.Unmarshal(plan.Files[".product/artifacts.json"].Data, &registry); err != nil {
			t.Fatal(err)
		}
		var ids []string
		for _, artifact := range registry.Artifacts {
			ids = append(ids, artifact["id"].(string))
		}
		sort.Strings(ids)
		sort.Strings(expected.ArtifactIDs)
		if strings.Join(ids, "\n") != strings.Join(expected.ArtifactIDs, "\n") {
			t.Fatalf("%s artifact ids=%v want=%v", point, ids, expected.ArtifactIDs)
		}
		for _, path := range expected.SpecialFiles {
			if _, ok := plan.Files[path]; !ok {
				t.Fatalf("%s missing special file %s", point, path)
			}
		}
		if strings.Join(plan.Directories, "\n") != strings.Join(expected.Directories, "\n") {
			t.Fatalf("%s directories=%v want=%v", point, plan.Directories, expected.Directories)
		}
		if strings.Join(plan.Actions, "\n") != strings.Join(expected.Actions, "\n") {
			t.Fatalf("%s actions=%v want=%v", point, plan.Actions, expected.Actions)
		}
	}
}

func TestInitializationCatalogDoesNotCopyTheMonolithicStarterRoot(t *testing.T) {
	_, catalog, err := loadInitContract("new-product")
	if err != nil {
		t.Fatal(err)
	}
	for setName, assets := range catalog.Sets {
		for _, asset := range assets {
			if asset.Source == "starter/product" {
				t.Fatalf("asset set %q restores static starter copy", setName)
			}
		}
	}
}

func TestStrictContractRejectsUnknownFieldAndTrailingValue(t *testing.T) {
	for name, contract := range map[string]string{
		"unknown field":  strings.Replace(validFixtureContract(), `"actions":[]`, `"actions":[],"unknown":true`, 1),
		"trailing value": validFixtureContract() + `{}`,
	} {
		t.Run(name, func(t *testing.T) {
			loader := fixtureLoader(contract, validFixtureCatalog(), validFixtureRegistry())
			if _, _, err := loader.load("fixture"); err == nil {
				t.Fatal("invalid JSON contract was accepted")
			}
		})
	}
}

func TestPlannerRejectsAssetCollision(t *testing.T) {
	catalog := `{"schema_version":1,"sets":{"base":[{"source":"assets/doc.md","target":"doc.md"},{"source":"assets/other.md","target":"doc.md"}]}}`
	loader := fixtureLoader(validFixtureContract(), catalog, validFixtureRegistry())
	if _, err := loader.buildPlan("fixture"); err == nil || !strings.Contains(err.Error(), "collision") {
		t.Fatalf("collision error=%v", err)
	}
}

func TestPlannerRejectsAmbiguousPatch(t *testing.T) {
	contract := strings.Replace(validFixtureContract(), `"patches":[]`, `"patches":[{"target":"doc.md","find":"content","replace":"updated"}]`, 1)
	loader := fixtureLoader(contract, validFixtureCatalog(), validFixtureRegistry())
	loader.assets.(fstest.MapFS)["assets/doc.md"] = &fstest.MapFile{Data: []byte("content content")}
	if _, err := loader.buildPlan("fixture"); err == nil || !strings.Contains(err.Error(), "expected one match") {
		t.Fatalf("patch error=%v", err)
	}
}

func TestPlannerRejectsBrokenArtifactRelationship(t *testing.T) {
	registry := `{"artifacts":[{"id":"DOMAIN-TBD","type":"domain","status":"draft","path":"doc.md","parentIds":["MISSING"]}]}`
	loader := fixtureLoader(validFixtureContract(), validFixtureCatalog(), registry)
	if _, err := loader.buildPlan("fixture"); err == nil || !strings.Contains(err.Error(), "unknown parent") {
		t.Fatalf("registry error=%v", err)
	}
}

func TestPlannerMaterializesExplicitEmptyDirectory(t *testing.T) {
	contract := strings.Replace(validFixtureContract(), `"directories":[]`, `"directories":["empty/nested"]`, 1)
	loader := fixtureLoader(contract, validFixtureCatalog(), validFixtureRegistry())
	plan, err := loader.buildPlan("fixture")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := writeInitializationPlan(root, plan); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(root, "empty", "nested"))
	if err != nil || !info.IsDir() {
		t.Fatalf("explicit empty directory was not materialized: %v", err)
	}
}

func TestPlannerRejectsDirectoryTraversalAndFileCollision(t *testing.T) {
	for name, directories := range map[string]string{
		"traversal":      `["../outside"]`,
		"file collision": `["doc.md"]`,
	} {
		t.Run(name, func(t *testing.T) {
			contract := strings.Replace(validFixtureContract(), `"directories":[]`, `"directories":`+directories, 1)
			loader := fixtureLoader(contract, validFixtureCatalog(), validFixtureRegistry())
			if _, err := loader.buildPlan("fixture"); err == nil {
				t.Fatal("invalid explicit directory was accepted")
			}
		})
	}
}

func TestFailedInitializationLeavesNoProductOrStagingTree(t *testing.T) {
	t.Setenv("SPEC_FRAMEWORK_CACHE", filepath.Join(t.TempDir(), "cache"))
	t.Setenv("SPEC_FRAMEWORK_AGENT_HOME", filepath.Join(t.TempDir(), "agents"))
	target := filepath.Join(t.TempDir(), "repo")
	_, err := Init(Options{Target: target, Version: "test", Agents: []Agent{Codex}, StartingPoint: "existing-documents", Sources: []string{filepath.Join(target, "missing.md")}})
	if err == nil {
		t.Fatal("init with a missing import source succeeded")
	}
	if _, statErr := os.Stat(filepath.Join(target, "product")); !os.IsNotExist(statErr) {
		t.Fatalf("failed init left product/: %v", statErr)
	}
	matches, globErr := filepath.Glob(filepath.Join(target, ".spec-framework-init-*"))
	if globErr != nil || len(matches) != 0 {
		t.Fatalf("failed init left staging paths: %v, %v", matches, globErr)
	}
}

func TestInitNeverOverwritesExistingProductWithForce(t *testing.T) {
	target := t.TempDir()
	marker := filepath.Join(target, "product", "owned.md")
	if err := os.MkdirAll(filepath.Dir(marker), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(marker, []byte("owned"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(Options{Target: target, Agents: []Agent{Codex}, Force: true}); err == nil {
		t.Fatal("force init overwrote an existing product")
	}
	data, err := os.ReadFile(marker)
	if err != nil || string(data) != "owned" {
		t.Fatalf("owned marker changed: %q, %v", data, err)
	}
}

func TestSafeProductPathRejectsTraversal(t *testing.T) {
	for _, target := range []string{"../outside", "a/../../outside", "/absolute"} {
		if _, err := safeProductPath("product", target); err == nil {
			t.Fatalf("unsafe target %q accepted", target)
		}
	}
}

func fixtureLoader(contract, catalog, registry string) initContractLoader {
	assets := fstest.MapFS{
		"framework/init/catalog.json":           &fstest.MapFile{Data: []byte(catalog)},
		"framework/init/contracts/fixture.json": &fstest.MapFile{Data: []byte(contract)},
		"assets/doc.md":                         &fstest.MapFile{Data: []byte("content")},
		"assets/other.md":                       &fstest.MapFile{Data: []byte("other")},
		"assets/artifacts.json":                 &fstest.MapFile{Data: []byte(registry)},
	}
	return initContractLoader{assets: assets}
}

func validFixtureCatalog() string {
	return `{"schema_version":1,"sets":{"base":[{"source":"assets/doc.md","target":"doc.md"},{"source":"assets/artifacts.json","target":".product/artifacts.json"}]}}`
}

func validFixtureContract() string {
	return `{"schema_version":1,"id":"fixture","asset_sets":["base"],"directories":[],"files":[],"patches":[],"registry":{},"bootstrap_profile":"fixture","actions":[]}`
}

func validFixtureRegistry() string {
	return `{"artifacts":[{"id":"DOMAIN-TBD","type":"domain","status":"draft","path":"doc.md","parentIds":[]}]}`
}

func planDigest(plan initializationPlan) [32]byte {
	targets := make([]string, 0, len(plan.Files))
	for target := range plan.Files {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	var buffer bytes.Buffer
	for _, directory := range plan.Directories {
		buffer.WriteString("dir:")
		buffer.WriteString(directory)
		buffer.WriteByte(0)
	}
	for _, target := range targets {
		buffer.WriteString(target)
		buffer.WriteByte(0)
		buffer.Write(plan.Files[target].Data)
		buffer.WriteByte(0)
	}
	actions, _ := json.Marshal(plan.Actions)
	buffer.Write(actions)
	return sha256.Sum256(buffer.Bytes())
}

var _ fs.FS = fstest.MapFS{}
