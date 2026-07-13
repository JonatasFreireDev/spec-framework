package validator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteReportUsesAtomicCanonicalPaths(t *testing.T) {
	root := t.TempDir()
	paths, err := WriteReport(root, Result{})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Fatalf("paths=%v", paths)
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Fatal(err)
		}
	}
}
func TestWriteRegistrySortsArtifacts(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "docs"), 0755)
	_ = os.WriteFile(filepath.Join(root, "docs", "b.md"), []byte("| ID | B-1 |\n| Type | test |\n| Status | draft |\n"), 0644)
	_ = os.WriteFile(filepath.Join(root, "docs", "a.md"), []byte("| ID | A-1 |\n| Type | test |\n| Status | draft |\n"), 0644)
	path, err := WriteRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	var value struct {
		Artifacts []map[string]any `json:"artifacts"`
	}
	if json.Unmarshal(data, &value) != nil || len(value.Artifacts) != 2 {
		t.Fatal(string(data))
	}
	if value.Artifacts[0]["id"] != "A-1" {
		t.Fatalf("%v", value.Artifacts)
	}
}

func TestBuildRegistryIncludesFoundationParents(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"foundation/problem/problem.md":   "| ID | `PROBLEM-1` |\n| Type | `problem` |\n| Status | `draft` |\n",
		"foundation/vision/vision.md":     "| ID | `VISION-1` |\n| Type | `vision` |\n| Parent IDs | `PROBLEM-1` |\n| Status | `draft` |\n",
		"foundation/vision/principles.md": "| ID | `PRINCIPLES-1` |\n| Type | `product-principles` |\n| Parent IDs | `VISION-1` |\n| Status | `draft` |\n",
		"foundation/vision/north-star.md": "| ID | `NORTH-STAR-1` |\n| Type | `north-star` |\n| Parent IDs | `VISION-1` |\n| Status | `draft` |\n",
		"foundation/strategy/strategy.md": "| ID | `STRATEGY-1` |\n| Type | `strategy` |\n| Parent IDs | `VISION-1, PRINCIPLES-1, NORTH-STAR-1` |\n| Status | `draft` |\n",
	}}
	artifacts := buildRegistry(s)
	if len(artifacts) != 5 {
		t.Fatalf("artifacts=%+v", artifacts)
	}
	for _, artifact := range artifacts {
		if artifact["id"] != "STRATEGY-1" {
			continue
		}
		parents, _ := artifact["parentIds"].([]string)
		if len(parents) != 3 || parents[0] != "VISION-1" || parents[1] != "PRINCIPLES-1" || parents[2] != "NORTH-STAR-1" {
			t.Fatalf("strategy parents=%v", parents)
		}
		return
	}
	t.Fatal("strategy missing from registry")
}

func TestBuildRegistryUsesOnlyFeatureBriefForExistingFeature(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"foundation/problem/problem.md":           "| ID | PROBLEM-1 |\n| Type | problem |\n| Status | draft |\n",
		"foundation/feature-brief.md":             "| ID | FEATURE-BRIEF-1 |\n| Type | feature-brief |\n| Status | draft |\n| Target Feature | FT-1 |\n",
		"domains/d/goals/g/features/f/feature.md": "| ID | FT-1 |\n| Status | draft |\n",
		"domains/d/goals/g/features/f/context.md": "```yaml\nid: FT-1\ntype: feature\nparents:\n  - GOAL-1\nchildren:\n  - UC-1\ndepends_on:\n  - ../goal.md\ndecisions:\n  - DEC-1\ndelivery:\n  depends_on:\n    - GOAL-1\n```\n",
	}, JSON: map[string]any{
		".product/framework.json": map[string]any{"starting_point": "existing-feature"},
	}}
	artifacts := buildRegistry(s)
	if len(artifacts) != 2 {
		t.Fatalf("artifacts=%+v", artifacts)
	}
	for _, artifact := range artifacts {
		if artifact["id"] == "FEATURE-BRIEF-1" && artifact["targetFeature"] != "FT-1" {
			t.Fatalf("Feature Brief target was not preserved: %+v", artifact)
		}
		if artifact["id"] == "FT-1" {
			parents, _ := artifact["parentIds"].([]string)
			if len(parents) != 2 || parents[0] != "GOAL-1" || parents[1] != "FEATURE-BRIEF-1" {
				t.Fatalf("feature parents=%v", parents)
			}
			children, _ := artifact["childIds"].([]string)
			dependsOn, _ := artifact["dependsOn"].([]string)
			decisions, _ := artifact["decisions"].([]string)
			if len(children) != 1 || children[0] != "UC-1" || len(dependsOn) != 1 || len(decisions) != 1 {
				t.Fatalf("context relationships were not preserved: %+v", artifact)
			}
			return
		}
	}
	t.Fatal("feature missing from registry")
}

func TestBuildRegistryLinksExistingImplementationAssessmentToProblem(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"knowledge/assessments/implementation-assessment.md": "| ID | IMPL-ASSESS-1 |\n| Type | implementation-assessment |\n| Status | draft |\n",
		"foundation/problem/problem.md":                      "| ID | PROBLEM-1 |\n| Type | problem |\n| Status | draft |\n",
	}, JSON: map[string]any{
		".product/framework.json": map[string]any{"starting_point": "existing-implementation"},
	}}
	artifacts := buildRegistry(s)
	if len(artifacts) != 2 {
		t.Fatalf("artifacts=%+v", artifacts)
	}
	for _, artifact := range artifacts {
		if artifact["id"] == "PROBLEM-1" {
			parents, _ := artifact["parentIds"].([]string)
			if len(parents) != 1 || parents[0] != "IMPL-ASSESS-1" {
				t.Fatalf("problem parents=%v", parents)
			}
			return
		}
	}
	t.Fatal("Problem missing from registry")
}

func TestBuildRegistryUsesProductBaselineAndStrategyForExistingProduct(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"foundation/product-baseline.md":  "| ID | PRODUCT-BASELINE-1 |\n| Type | product-baseline |\n| Status | draft |\n",
		"foundation/problem/problem.md":   "| ID | PROBLEM-1 |\n| Type | problem |\n| Status | draft |\n",
		"foundation/vision/vision.md":     "| ID | VISION-1 |\n| Type | vision |\n| Status | draft |\n",
		"foundation/strategy/strategy.md": "| ID | STRATEGY-1 |\n| Type | strategy |\n| Status | draft |\n",
	}, JSON: map[string]any{
		".product/framework.json": map[string]any{"starting_point": "existing-product"},
	}}
	artifacts := buildRegistry(s)
	if len(artifacts) != 2 {
		t.Fatalf("artifacts=%+v", artifacts)
	}
	for _, artifact := range artifacts {
		if artifact["id"] == "STRATEGY-1" {
			parents, _ := artifact["parentIds"].([]string)
			if len(parents) != 1 || parents[0] != "PRODUCT-BASELINE-1" {
				t.Fatalf("strategy parents=%v", parents)
			}
			return
		}
	}
	t.Fatal("Strategy missing from registry")
}
