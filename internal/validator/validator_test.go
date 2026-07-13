package validator

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/engineeringsystem"
)

func TestImportValidationDetectsChangedSourceAndDuplicateTargets(t *testing.T) {
	root := t.TempDir()
	run := filepath.Join(root, "knowledge", "imports", "runs", "IMPORT-001")
	sourceRel := "knowledge/imports/sources/epic.md"
	source := filepath.Join(root, filepath.FromSlash(sourceRel))
	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(run, 0755); err != nil {
		t.Fatal(err)
	}
	original := []byte("original")
	sum := sha256.Sum256(original)
	if err := os.WriteFile(source, []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	write := func(name string, value any) {
		data, _ := json.Marshal(value)
		if name == "mapping.json" {
			data = append([]byte{0xef, 0xbb, 0xbf}, data...)
		}
		if err := os.WriteFile(filepath.Join(run, name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("inventory.json", map[string]any{"schema_version": 1, "import_id": "IMPORT-001", "sources": []any{map[string]any{"path": sourceRel, "sha256": fmt.Sprintf("%x", sum[:])}}})
	write("import-plan.json", map[string]any{"materialization_approved": false})
	write("mapping.json", map[string]any{"mappings": []any{map[string]any{"id": "MAP-1", "selected": true, "target": "domains/a/domain.md", "source_documents": []any{sourceRel}}, map[string]any{"id": "MAP-2", "selected": true, "target": "domains/a/domain.md", "source_documents": []any{sourceRel}}}})
	for _, name := range []string{"conflicts.md", "import-report.md"} {
		if err := os.WriteFile(filepath.Join(run, name), []byte("# Report"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	changed, duplicate := false, false
	for _, d := range result.Diagnostics {
		if d.Check == "imports" && strings.Contains(d.Message, "Source changed") {
			changed = true
		}
		if d.Check == "imports" && strings.Contains(d.Message, "Multiple selected mappings") {
			duplicate = true
		}
	}
	if !changed || !duplicate {
		t.Fatalf("changed=%v duplicate=%v diagnostics=%+v", changed, duplicate, result.Diagnostics)
	}
}

func TestDeliveryClosureRejectsLegacyAndUnknownHandoffs(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "framework", "skills", "known"), 0755)
	s := Snapshot{Root: root, FrameworkRoot: root, Text: map[string]string{"framework/skills/a/SKILL.md": "## Handoff\nNext: 05-old.md\n", "framework/skills/b/SKILL.md": "## Handoff\nNext: missing-skill.\n"}}
	d := validateDeliveryClosure(s)
	legacy, unknown := false, false
	for _, x := range d {
		if strings.Contains(x.Message, "Legacy numbered") {
			legacy = true
		}
		if strings.Contains(x.Message, "Unknown next skill") {
			unknown = true
		}
	}
	if !legacy || !unknown {
		t.Fatalf("legacy=%v unknown=%v diagnostics=%+v", legacy, unknown, d)
	}
}

func TestValidatedTaskRequiresMatchingDiffHashes(t *testing.T) {
	text := "# Task\n\n| Field | Value |\n| --- | --- |\n| Status | validated |\n| Branch | feature/x |\n| Base commit | abcdef1 |\n| Diff hash | hash-a |\n| Changed paths | src/x.go |\n| Test status | passed |\n| Commits | abcdef2 |\n| Code paths | src/x.go |\n| Code Review diff hash | hash-a |\n| QA diff hash | hash-b |\n"
	s := Snapshot{Text: map[string]string{"domains/x/use-cases/u/tasks/TK-1.md": text}}
	d := validateDeliveryClosure(s)
	found := false
	for _, x := range d {
		if x.Check == "diff-staleness" {
			found = true
		}
	}
	if !found {
		t.Fatalf("diagnostics=%+v", d)
	}
}

func TestTierLRequiresEngineeringProposalAndReview(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"domains/d/goals/g/features/f/use-cases/u/context.md": "---\nrigor_tier: L\n---\n",
	}}
	diagnostics := validateUseCaseBundles(s)
	wants := map[string]bool{"engineering-proposal.md": false, "engineering-review.md": false}
	for _, diagnostic := range diagnostics {
		for name := range wants {
			if strings.HasSuffix(diagnostic.File, name) {
				wants[name] = true
			}
		}
	}
	for name, found := range wants {
		if !found {
			t.Fatalf("missing diagnostic for %s: %+v", name, diagnostics)
		}
	}
}

func TestTierLAdvancedPlanRequiresPassedEngineeringReview(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	s := Snapshot{Text: map[string]string{
		base + "/context.md":              "---\nrigor_tier: L\n---\n",
		base + "/implementation-plan.md":  "| Field | Value |\n| --- | --- |\n| Status | `proposed` |\n",
		base + "/technical-discovery.md":  "| Field | Value |\n| --- | --- |\n| Status | `approved` |\n| Verdict | Not required |\n",
		base + "/engineering-review.md":   "| Field | Value |\n| --- | --- |\n| Status | `draft` |\n| Verdict | `blocked` |\n",
		base + "/engineering-proposal.md": "| Field | Value |\n| --- | --- |\n| Status | `draft` |\n",
	}}
	found := false
	for _, diagnostic := range validateDeliveryClosure(s) {
		if diagnostic.Check == "engineering-review-gate" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected engineering-review-gate diagnostic")
	}
}

func TestRegistryNormalizesInlineCodeTableValues(t *testing.T) {
	path := "domains/d/goals/g/features/f/use-cases/u/engineering-proposal.md"
	s := Snapshot{Text: map[string]string{
		"domains/d/goals/g/features/f/use-cases/u/context.md": "---\nid: UC-1\nrigor_tier: L\n---\n",
		path: "| Field | Value |\n| --- | --- |\n| ID | `ENGPROP-1` |\n| Status | `draft` |\n| Owner skill | `engineering-proposal` |\n| Level | `L1` |\n| Priority | `P0` |\n| Rationale | Inherited. |\n",
	}}
	artifacts := buildRegistry(s)
	if len(artifacts) != 1 {
		t.Fatalf("artifacts=%+v", artifacts)
	}
	artifact := artifacts[0]
	delivery, _ := artifact["delivery"].(map[string]any)
	if artifact["id"] != "ENGPROP-1" || artifact["status"] != "draft" || artifact["ownerSkill"] != "engineering-proposal" || delivery["level"] != "L1" || delivery["priority"] != "P0" {
		t.Fatalf("artifact=%+v", artifact)
	}
}

func TestFoundationWithoutRegistryMetadataIsRejected(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"foundation/problem/problem.md": "# Problem\n\n## Status\n\nDraft.\n",
	}, JSON: map[string]any{
		".product/artifacts.json": map[string]any{"artifacts": []any{map[string]any{"id": "DOMAIN-1", "type": "domain", "status": "draft", "path": "domains/d/context.md"}}},
	}}
	found := false
	for _, diagnostic := range validateRegistryAndApprovalGates(s) {
		if diagnostic.Check == "foundation-registry" && diagnostic.File == "foundation/problem/problem.md" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected missing Foundation registry diagnostic")
	}
}

func TestExistingFeatureBriefWithoutRegistryMetadataIsRejected(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"foundation/feature-brief.md": "# Feature Brief\n",
	}, JSON: map[string]any{
		".product/framework.json": map[string]any{"starting_point": "existing-feature"},
		".product/artifacts.json": map[string]any{"artifacts": []any{map[string]any{"id": "DOMAIN-1", "type": "domain", "status": "draft", "path": "domains/d/context.md"}}},
	}}
	for _, diagnostic := range validateRegistryAndApprovalGates(s) {
		if diagnostic.Check == "foundation-registry" && diagnostic.File == "foundation/feature-brief.md" {
			return
		}
	}
	t.Fatal("expected missing Feature Brief registry diagnostic")
}

func TestExistingImplementationAssessmentWithoutRegistryMetadataIsRejected(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"knowledge/assessments/implementation-assessment.md": "# Implementation Assessment\n",
	}, JSON: map[string]any{
		".product/framework.json": map[string]any{"starting_point": "existing-implementation"},
		".product/artifacts.json": map[string]any{"artifacts": []any{map[string]any{"id": "DOMAIN-1", "type": "domain", "status": "draft", "path": "domains/d/context.md"}}},
	}}
	for _, diagnostic := range validateRegistryAndApprovalGates(s) {
		if diagnostic.Check == "starting-point-registry" && diagnostic.File == "knowledge/assessments/implementation-assessment.md" {
			return
		}
	}
	t.Fatal("expected missing Implementation Assessment registry diagnostic")
}

func TestExistingProductBaselineWithoutRegistryMetadataIsRejected(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"foundation/product-baseline.md":  "# Product Baseline\n",
		"foundation/strategy/strategy.md": "| ID | STRATEGY-1 |\n| Type | strategy |\n| Status | draft |\n",
	}, JSON: map[string]any{
		".product/framework.json": map[string]any{"starting_point": "existing-product"},
		".product/artifacts.json": map[string]any{"artifacts": []any{map[string]any{"id": "STRATEGY-1", "type": "strategy", "status": "draft", "path": "foundation/strategy/strategy.md"}}},
	}}
	for _, diagnostic := range validateRegistryAndApprovalGates(s) {
		if diagnostic.Check == "foundation-registry" && diagnostic.File == "foundation/product-baseline.md" {
			return
		}
	}
	t.Fatal("expected missing Product Baseline registry diagnostic")
}

func TestTierMRegistryDoesNotRequireEngineeringReview(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	s := Snapshot{Text: map[string]string{
		base + "/context.md":             "---\nid: UC-1\ntype: use-case\nstatus: approved\nrigor_tier: M\n---\n",
		base + "/technical-discovery.md": "| Field | Value |\n| --- | --- |\n| ID | TD-1 |\n| Status | approved |\n",
		base + "/implementation-plan.md": "| Field | Value |\n| --- | --- |\n| ID | PLAN-1 |\n| Status | proposed |\n",
	}}
	for _, diagnostic := range validateRegistryAndApprovalGates(s) {
		if diagnostic.Check == "approval-gates" && strings.Contains(diagnostic.Message, "Engineering Review") {
			t.Fatalf("Tier M should remain compatible: %+v", diagnostic)
		}
	}
}

func TestTierMTriggerRequiresEngineeringProposalAndReview(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	s := Snapshot{Text: map[string]string{
		base + "/context.md": "---\nrigor_tier: M\nengineering_triggers:\n  - new_dependency\n---\n",
	}}
	diagnostics := validateUseCaseBundles(s)
	wants := map[string]bool{"engineering-proposal.md": false, "engineering-review.md": false}
	for _, diagnostic := range diagnostics {
		for name := range wants {
			if strings.HasSuffix(diagnostic.File, name) {
				wants[name] = true
			}
		}
	}
	for name, found := range wants {
		if !found {
			t.Fatalf("trigger did not require %s: %+v", name, diagnostics)
		}
	}
}

func TestUnknownEngineeringTriggerIsRejected(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	s := Snapshot{Text: map[string]string{
		base + "/context.md": "---\nrigor_tier: M\nengineering_triggers:\n  - invented_trigger\n---\n",
	}}
	found := false
	for _, diagnostic := range validateDeliveryClosure(s) {
		if diagnostic.Check == "engineering-trigger" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected unknown engineering trigger diagnostic")
	}
}

func TestStructuredNotApplicableRequiresRationale(t *testing.T) {
	s := Snapshot{Text: map[string]string{
		"domains/d/use-cases/u/design.md": "| Field | Value |\n| --- | --- |\n| Status | not_applicable |\n| Rationale | TBD |\n",
	}}
	found := false
	for _, diagnostic := range validateDeliveryClosure(s) {
		if diagnostic.Check == "not-applicable" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected not_applicable rationale diagnostic")
	}
}

func TestEventsFixtureRemainsReady(t *testing.T) {
	frameworkRoot := filepath.Clean(filepath.Join("..", ".."))
	productRoot := filepath.Join(frameworkRoot, "examples", "events")
	result, err := Validate(context.Background(), productRoot, frameworkRoot)
	if err != nil {
		t.Fatal(err)
	}
	if result.Errors != 0 || result.Warnings != 0 || result.Notes != 0 {
		t.Fatalf("events fixture is not ready: %+v", result.Diagnostics)
	}
}

func TestEngineeringProposalMustPinCurrentSystem(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"engineering/context.md":                        "---\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering/engineering-system.md":             "| Field | Value |\n| --- | --- |\n| ID | ENGSYS-001 |\n| Status | draft |\n",
		"engineering/engineering-system.yaml":           "schema_version: 1\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: product\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: baseline\n    evidence: []\n",
		"engineering/architecture/modules.md":           "# Modules\n",
		"domains/d/use-cases/u/engineering-proposal.md": "| Field | Value |\n| --- | --- |\n| ID | ENGPROP-1 |\n| Status | draft |\n| Engineering System | ENGSYS-OLD @ 0.9.0 |\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := Scan(context.Background(), root, root, 1)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, diagnostic := range validateEngineeringSystem(snapshot) {
		if diagnostic.Check == "engineering-system-consumer" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected invalid Engineering System pin diagnostic")
	}
}

func TestPassedEngineeringReviewMustMatchProposalHash(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"engineering/context.md":                        "---\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering/engineering-system.md":             "# System\n",
		"engineering/engineering-system.yaml":           "schema_version: 1\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: not-configured\nareas:\n  modules:\n    contract: architecture/modules.md\n    maturity: baseline\n    evidence: []\n",
		"engineering/architecture/modules.md":           "# Modules\n",
		"domains/d/use-cases/u/engineering-proposal.md": "# Current Proposal\n\n| Field | Value |\n| --- | --- |\n| Engineering System | Not configured |\n",
		"domains/d/use-cases/u/engineering-review.md":   "| Field | Value |\n| --- | --- |\n| Verdict | passed |\n| Proposal hash | deadbeef |\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := Scan(context.Background(), root, root, 1)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, diagnostic := range validateEngineeringSystem(snapshot) {
		if diagnostic.Check == "engineering-review-staleness" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected stale Engineering Review diagnostic")
	}
}

func TestEngineeringSystemCompositeHashParity(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"engineering/context.md":              "status: approved\n",
		"engineering/engineering-system.md":   "system\n",
		"engineering/engineering-system.yaml": "status: approved\n",
		"engineering/architecture/modules.md": "modules\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := Scan(context.Background(), root, root, 1)
	if err != nil {
		t.Fatal(err)
	}
	want, err := engineeringsystem.CompositeHash(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := artifactHash(snapshot, "engineering/engineering-system.md", snapshot.Text["engineering/engineering-system.md"])
	if got != want {
		t.Fatalf("validator hash=%s engineering hash=%s", got, want)
	}
}

func TestDiagnosticsAreDeterministic(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "domains", "a"), 0755)
	_ = os.WriteFile(filepath.Join(root, "domains", "a", "context.md"), []byte("status: draft\n"), 0644)
	first, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	second, _ := Validate(context.Background(), root, root)
	if len(first.Diagnostics) == 0 || len(first.Diagnostics) != len(second.Diagnostics) {
		t.Fatalf("diagnostics=%v", first.Diagnostics)
	}
	for i := range first.Diagnostics {
		if first.Diagnostics[i] != second.Diagnostics[i] {
			t.Fatal("unstable diagnostics")
		}
	}
}

func TestBlocksBrokenMarkdownLink(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "index.md"), []byte("[Missing](missing.md)\n"), 0644)
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Check == "links" {
			found = true
		}
	}
	if !found {
		t.Fatalf("%+v", result)
	}
}

func TestRequiresMatchingApprovalRecord(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".product"), 0755)
	artifact := "# Artifact\n"
	_ = os.WriteFile(filepath.Join(root, "artifact.md"), []byte(artifact), 0644)
	registry := map[string]any{"artifacts": []any{map[string]any{"id": "ART-1", "status": "approved", "path": "artifact.md"}}}
	data, _ := json.Marshal(registry)
	_ = os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0644)
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range result.Diagnostics {
		if d.Check == "approval-records" {
			found = true
		}
	}
	if !found {
		t.Fatalf("%+v", result)
	}
}
