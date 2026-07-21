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

func TestReviewFindingValidationRequiresNormalizedEvidence(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".product", "reviews", "findings")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	write := func(name string, value any) {
		data, _ := json.Marshal(value)
		if err := os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("one.json", map[string]any{"id": "RF-1", "source": "github", "reference": "https://example.test/1", "severity": "warning", "description": "missing test", "status": "open", "scope": "src/a.go", "evidence": "comment", "owner": "qa"})
	write("duplicate.json", map[string]any{"id": "RF-1", "source": "github", "reference": "https://example.test/2", "severity": "warning", "description": "missing test", "status": "open", "scope": "src/b.go", "evidence": "comment", "owner": "qa"})
	write("invalid.json", map[string]any{"id": "RF-2"})
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	duplicate, invalid := false, false
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Check != "review-findings" {
			continue
		}
		duplicate = duplicate || strings.Contains(diagnostic.Message, "Duplicate")
		invalid = invalid || strings.Contains(diagnostic.Message, "required")
	}
	if !duplicate || !invalid {
		t.Fatalf("duplicate=%v invalid=%v diagnostics=%+v", duplicate, invalid, result.Diagnostics)
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

func TestSkillDiscoveryContractIsMechanicallyRequired(t *testing.T) {
	root := t.TempDir()
	skill := filepath.Join(root, "framework", "skills", "feature", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(skill), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skill, []byte("# Feature\n\n## Workflow\n"), 0644); err != nil {
		t.Fatal(err)
	}
	diagnostics := validateSkillDiscoveryContracts(Snapshot{FrameworkRoot: root})
	found := 0
	for _, diagnostic := range diagnostics {
		if diagnostic.Check == "skill-discovery-contract" && diagnostic.File == "framework/skills/feature/SKILL.md" {
			found++
		}
	}
	if found != 2 {
		t.Fatalf("expected missing section and reference diagnostics, got %+v", diagnostics)
	}
	content := "# Feature\n\n## Discovery and challenge\n\nFollow [contract](../discovery-and-challenge.md).\n"
	if err := os.WriteFile(skill, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range validateSkillDiscoveryContracts(Snapshot{FrameworkRoot: root}) {
		if diagnostic.File == "framework/skills/feature/SKILL.md" {
			t.Fatalf("valid contract rejected: %+v", diagnostic)
		}
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

func TestApprovedArtifactsRequireCanonicalTemplateSections(t *testing.T) {
	s := Snapshot{
		Text: map[string]string{"tasks/TK-001.md": "# Task\n\n## Objective\nDo the thing.\n"},
		JSON: map[string]any{".product/artifacts.json": map[string]any{"artifacts": []any{map[string]any{
			"id": "TK-001", "type": "task", "status": "approved", "path": "tasks/TK-001.md",
		}}}},
	}
	diagnostics := validateTemplateConformance(s)
	if len(diagnostics) != 1 || diagnostics[0].Severity != Error || diagnostics[0].Check != "template-conformance" {
		t.Fatalf("expected one blocking template diagnostic, got %+v", diagnostics)
	}
}

func TestExecutionGraphRejectsLegacyTaskDependencyField(t *testing.T) {
	s := Snapshot{Text: map[string]string{}, JSON: map[string]any{"execution-graph.json": map[string]any{
		"status": "draft",
		"nodes":  []any{map[string]any{"id": "TK-1", "path": "tasks/TK-1.md", "dependsOn": []any{}, "dependencies": []any{"TK-0"}, "writeScope": []any{"src"}}},
	}}}
	for _, diagnostic := range validateGraph("execution-graph.json", s.JSON["execution-graph.json"], s) {
		if strings.Contains(diagnostic.Message, "non-canonical dependencies") {
			return
		}
	}
	t.Fatal("expected non-canonical dependency diagnostic")
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

func TestLegacySpecificationShellRemainsCompatibleBeforeV2OptIn(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	shell := "# Contract\nStatus: draft\n\nThis contract will receive stable REQ/AC identifiers during the next approved content evolution.\n\n## Scope\n\n- Preserve approved behavior.\n"
	texts := map[string]string{
		base + "/context.md":              "---\nrigor_tier: L\n---\n",
		base + "/technical-discovery.md":  "# Technical Discovery\n",
		base + "/engineering-proposal.md": "# Engineering Proposal\n",
		base + "/engineering-review.md":   "# Engineering Review\n",
	}
	for _, name := range []string{"product", "behavior", "ux", "api", "data", "security", "quality", "observability", "rollout"} {
		texts[base+"/contracts/"+name+".md"] = shell
	}

	diagnostics := validateDeliveryClosure(Snapshot{Text: texts})
	if len(diagnostics) != 0 {
		t.Fatalf("legacy contract without specification_contract_version must remain compatible, got %+v", diagnostics)
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

func TestContextEvolutionRelationsRequireListValues(t *testing.T) {
	text := "```yaml\nid: UC-1\ntype: use-case\nname: Example\nstatus: draft\nowner_skill: Use Case AI\nslug: u\nrelations:\n  extends: UC-0\n  reuses: not-a-list\n  impacts:\n    - services/api\n```\n"
	find := false
	for _, diagnostic := range validateContextFull("domains/d/goals/g/features/f/use-cases/u/context.md", text) {
		if diagnostic.Check == "context-relations" {
			find = true
		}
	}
	if !find {
		t.Fatal("expected invalid evolution relation diagnostic")
	}
}

func TestContextEvolutionRelationsRemainOptional(t *testing.T) {
	text := "```yaml\nid: UC-1\ntype: use-case\nname: Example\nstatus: draft\nowner_skill: Use Case AI\nslug: u\n```\n\n## Handoff\n"
	for _, diagnostic := range validateContextFull("domains/d/goals/g/features/f/use-cases/u/context.md", text) {
		if diagnostic.Check == "context-relations" {
			t.Fatalf("legacy context should remain compatible: %+v", diagnostic)
		}
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

func TestEventsFixtureHasOnlyExpectedSpecificationMigrationStaleness(t *testing.T) {
	frameworkRoot := filepath.Clean(filepath.Join("..", ".."))
	productRoot := filepath.Join(frameworkRoot, "examples", "events")
	result, err := Validate(context.Background(), productRoot, frameworkRoot)
	if err != nil {
		t.Fatal(err)
	}
	if result.Errors != 0 {
		t.Fatalf("events fixture has blocking findings: %+v", result.Diagnostics)
	}
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Severity == Warning && diagnostic.Check != "staleness" {
			t.Fatalf("events fixture has an unexpected warning: %+v", diagnostic)
		}
		if diagnostic.Check == "specification-depth" || diagnostic.Check == "specification-traceability" || diagnostic.Check == "specification-depth-migration" {
			t.Fatalf("events fixture did not complete Specification v2 migration: %+v", diagnostic)
		}
	}
}

func TestSpecificationV2RejectsGenericShells(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	root := specificationV2Root("proposed")
	shell := "# Contract\n\n| Status | proposed |\n\nThis contract will receive stable REQ/AC identifiers during the next approved content evolution.\n"
	s := Snapshot{FrameworkRoot: filepath.Clean(filepath.Join("..", "..")), Text: map[string]string{
		base + "/context.md":            "---\nrigor_tier: S\nspecification_contract_version: 2\n---\n",
		base + "/specification.md":      root,
		base + "/contracts/behavior.md": shell,
		base + "/contracts/quality.md":  shell,
	}}
	diagnostics := validateSpecificationDepth(s)
	foundBlocking := false
	for _, diagnostic := range diagnostics {
		if diagnostic.Check == "specification-depth" && diagnostic.Severity == Error {
			foundBlocking = true
		}
	}
	if !foundBlocking {
		t.Fatalf("expected v2 generic shells to block proposed Specification, got %+v", diagnostics)
	}
}

func TestLegacySpecificationDoesNotRequireV2Registry(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	s := Snapshot{FrameworkRoot: t.TempDir(), Text: map[string]string{
		base + "/context.md":       "---\nrigor_tier: S\n---\n",
		base + "/specification.md": "# Legacy Specification\n",
	}}
	diagnostics := validateSpecificationDepth(s)
	if len(diagnostics) != 1 || diagnostics[0].Check != "specification-depth-migration" || diagnostics[0].Severity != Note {
		t.Fatalf("legacy Specification should need no v2 registry, got %+v", diagnostics)
	}
}

func TestSpecificationV2AcceptsConcernSpecificTierSContracts(t *testing.T) {
	base := "domains/d/goals/g/features/f/use-cases/u"
	behavior := "# Behavior Contract\n\n| Status | proposed |\n\n## Preconditions And Triggers\nConcrete trigger from [use case](../use-case.md).\n\n## State Transitions\nConcrete transition.\n\n## Alternate Error And Edge Flows\nConcrete failure and recovery.\n\n## Invariants\nConcrete invariant.\n\n## Requirements\n\n| ID | Requirement | Source | Acceptance criteria | Dependencies |\n| --- | --- | --- | --- | --- |\n| REQ-001 | Deterministic behavior. | [Use case](../use-case.md) | AC-001 | None |\n"
	quality := "# Quality Contract\n\n| Status | proposed |\n\n## Quality Risks\nConcrete risk.\n\n## Acceptance Traceability\nConcrete mapping.\n\n## Test Levels And Environments\nConcrete level and environment.\n\n## Evidence And Exit Conditions\nConcrete evidence condition.\n\n## Requirements\n\n| ID | Requirement | Source | Acceptance criteria | Dependencies |\n| --- | --- | --- | --- | --- |\n| REQ-002 | Verification is repeatable. | [Use case](../use-case.md) | AC-002 | REQ-001 |\n"
	s := Snapshot{FrameworkRoot: filepath.Clean(filepath.Join("..", "..")), Text: map[string]string{
		base + "/context.md":            "---\nrigor_tier: S\nspecification_contract_version: 2\n---\n",
		base + "/specification.md":      specificationV2Root("proposed"),
		base + "/contracts/behavior.md": behavior,
		base + "/contracts/quality.md":  quality,
	}}
	if diagnostics := validateSpecificationDepth(s); len(diagnostics) != 0 {
		t.Fatalf("expected complete Tier S v2 bundle, got %+v", diagnostics)
	}
}

func specificationV2Root(status string) string {
	return "# Specification\n\n| Status | " + status + " |\n\n## Evidence And Boundary\nConcrete evidence from [use case](use-case.md).\n\n## Cross-Contract Synthesis\nConcrete synthesis.\n\n## Traceability Summary\nREQ-001 maps to AC-001.\n\n## Adversarial Review\nNo material gap found.\n\n## Open Questions And Decisions\n\n| Question/Decision | Owner | Blocks |\n| --- | --- | --- |\n| None | None | None |\n"
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

func TestProposedTestsMustPinConfiguredQualitySystem(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"engineering/context.md":                      "---\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\n---\n",
		"engineering/engineering-system.md":           "| Field | Value |\n| --- | --- |\n| ID | ENGSYS-001 |\n| Status | draft |\n| Version | 1.0.0 |\n",
		"engineering/engineering-system.yaml":         "schema_version: 1\nid: ENGSYS-001\nstatus: draft\nversion: 1.0.0\norigin_mode: generate\nscope: product\nareas:\n  quality:\n    contract: quality/quality-system.md\n    maturity: baseline\n    evidence: []\n",
		"engineering/quality/quality-system.md":       "| Field | Value |\n| --- | --- |\n| Engineering System | ENGSYS-001 @ 1.0.0 |\n| Status | draft |\n\n| Area | Policy | Evidence | Maturity |\n| --- | --- | --- | --- |\n| Behavioral | strategy | none | baseline |\n| Accessibility | strategy | none | baseline |\n| Security and privacy | strategy | none | baseline |\n| Performance and reliability | model | none | baseline |\n| Observability | model | none | baseline |\n",
		"engineering/quality/quality-model.md":        "# Model\n",
		"engineering/quality/test-strategy.md":        "# Strategy\n",
		"engineering/quality/quality-system.yaml":     "schema_version: 1\nengineering_system: ENGSYS-001\nversion: 1.0.0\nstatus: draft\nareas:\n  behavioral: {maturity: baseline, policy: test-strategy.md}\n  accessibility: {maturity: baseline, policy: test-strategy.md}\n  security_privacy: {maturity: baseline, policy: test-strategy.md, delegated_gate: security-review}\n  performance_reliability: {maturity: baseline, policy: quality-model.md}\n  observability: {maturity: baseline, policy: quality-model.md}\ngate_source: knowledge/conventions/gates.md\nenvironments: [ci]\ntest_data_classes: [synthetic]\nplatforms: [server]\nexceptions:\n  require_owner: true\n  require_residual_risk: true\n  require_expiry_or_review: true\n",
		"domains/d/use-cases/u/tests.md":              "| Field | Value |\n| --- | --- |\n| Status | proposed |\n| Engineering System | ENGSYS-OLD @ 0.9.0 |\n| Quality policy | engineering/quality/quality-system.md |\n| Applicable risks | behavior |\n| Environments | CI |\n| Test data | synthetic |\n| Platforms | server |\n| Deviations or exceptions | None |\n",
		"domains/d/use-cases/u/contracts/behavior.md": "| ID | Requirement | Acceptance criteria |\n| --- | --- | --- |\n| REQ-1 | works | AC-001 |\n",
		"domains/d/use-cases/u/qa-evidence.md":        "| Field | Value |\n| --- | --- |\n| Status | approved |\n| Engineering System | ENGSYS-001 @ 1.0.0 |\n\n| Check | Evidence | Result | Notes |\n| --- | --- | --- | --- |\n| Quality System conformity | policy | blocked | missing |\n| Environment and test data policy | CI | blocked | missing |\n| Flaky test and exception policy | scan | blocked | missing |\n",
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
	found := 0
	for _, diagnostic := range validateEngineeringSystem(snapshot) {
		if diagnostic.Check == "quality-system-consumer" {
			found++
		}
	}
	if found != 2 {
		t.Fatalf("expected pin and structural mapping diagnostics, got %d", found)
	}
	qaFound := 0
	for _, diagnostic := range validateEngineeringSystem(snapshot) {
		if diagnostic.Check == "quality-system-qa" {
			qaFound++
		}
	}
	if qaFound != 3 {
		t.Fatalf("expected three QA policy diagnostics, got %d", qaFound)
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

func TestQualityTraceabilityRequiresStructuredCompleteRows(t *testing.T) {
	text := "## Acceptance Traceability\n\n| Acceptance Criterion | Risk | Validation Method | Test Level | Evidence | Owner |\n| --- | --- | --- | --- | --- | --- |\n| [AC-001](specification.md) | high | automated | integration | test\\|log | qa |\n"
	rows := markdownTableRows(text, "Acceptance Traceability")
	if len(rows) != 1 || !containsExactID(rows[0]["acceptance criterion"], "AC-001") || rows[0]["validation method"] != "automated" || rows[0]["evidence"] != "test|log" {
		t.Fatalf("rows=%v", rows)
	}
	if qualityCheckPassed("| Quality System conformity | policy | N/A | none |\n", "Quality System conformity") {
		t.Fatal("N/A satisfied a required Quality System QA check")
	}
	if !qualityCheckPassed("| Quality System conformity | policy | passed | verified |\n", "Quality System conformity") {
		t.Fatal("passed result did not satisfy Quality System QA check")
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

func TestDomainModelingWarnings(t *testing.T) {
	snapshot := Snapshot{Text: map[string]string{
		"context.md":                   "```yaml\nname: FocusFlow\nslug: focusflow\n```\n",
		"domains/focusflow/context.md": "```yaml\nstatus: approved\n```\n",
		"domains/focusflow/domain.md":  "## Owns\n\n- Authentication, login, and tasks.\n",
	}}
	diagnostics := validateDomainModeling(snapshot)
	found := map[string]bool{}
	for _, diagnostic := range diagnostics {
		found[diagnostic.Check] = true
		if diagnostic.Severity != Warning {
			t.Fatalf("diagnostic %#v is not a warning", diagnostic)
		}
	}
	for _, check := range []string{"domain-product-name", "domain-missing-boundaries", "domain-chain-incomplete", "domain-monolith"} {
		if !found[check] {
			t.Errorf("missing %s warning: %#v", check, diagnostics)
		}
	}
}

func TestDomainAuthenticationWarningOnlyUsesOwnershipSection(t *testing.T) {
	if domainOwnsAuthentication("### Owns\n\n- Tasks.\n\n### Does Not Own\n\n- Authentication.\n") {
		t.Fatal("authentication listed as non-ownership must not trigger domain-monolith")
	}
	if !domainOwnsAuthentication("## Owns\n\n- Authentication and login.\n\n## Does Not Own\n\n- Payments.\n") {
		t.Fatal("authentication ownership must trigger domain-monolith")
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

func TestBlocksBrokenMarkdownSectionLink(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "index.md"), []byte("[Missing](target.md#missing)\n"), 0644)
	_ = os.WriteFile(filepath.Join(root, "target.md"), []byte("# Present\n"), 0644)
	result, err := Validate(context.Background(), root, root)
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Check == "links" && strings.Contains(diagnostic.Message, "section") {
			return
		}
	}
	t.Fatalf("expected broken section diagnostic: %+v", result)
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
