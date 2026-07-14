package install

import (
	"strings"
	"testing"
)

func TestBootstrapProfilesGiveEveryStartingPointAnImmediateAction(t *testing.T) {
	for _, startingPoint := range []string{"new-product", "existing-product", "existing-documents", "existing-feature", "existing-implementation", "audit-only"} {
		t.Run(startingPoint, func(t *testing.T) {
			profile := bootstrapProfileFor(startingPoint)
			if profile.location == "" || profile.nextAction == "" || profile.style == "" || profile.scopeRule == "" || profile.artifactAction == "" || profile.approvalIntro == "" || profile.approvalGuidance == "" || profile.approvalRule == "" || profile.foundationPath == "" || profile.agentPrompts == "" {
				t.Fatalf("incomplete profile: %+v", profile)
			}
			bootstrap := bootstrapFor(startingPoint)
			for _, expected := range []string{startingPoint, profile.location, profile.nextAction, profile.style, profile.scopeRule, profile.artifactAction, profile.approvalIntro, profile.approvalGuidance, profile.approvalRule, profile.agentPrompts, profile.foundationPath} {
				if !strings.Contains(bootstrap, expected) {
					t.Errorf("bootstrap does not contain %q", expected)
				}
			}
		})
	}
}

func TestBootstrapExplainsFoundationBeforeWorkspace(t *testing.T) {
	bootstrap := bootstrapFor("existing-feature")
	for _, expected := range []string{
		"Structural validity",
		"replaces the full product Foundation package with one Feature Brief",
		"spec-framework approve --product-root product --artifact foundation/feature-brief.md",
		"A Markdown status edit alone is not approval.",
		"require a WORK-NNN workspace",
		"spec-framework work --feature <id-or-path>",
		"## Before modeling domains",
		"examples/events/",
		"Domain -> User Goal -> Feature -> Use Case",
	} {
		if !strings.Contains(bootstrap, expected) {
			t.Errorf("bootstrap missing %q", expected)
		}
	}
	if strings.Contains(bootstrap, "approve-stage --stage foundation") {
		t.Fatal("bootstrap must not recommend batch Foundation approval")
	}
	if strings.Contains(bootstrap, "foundation/problem/problem.md") {
		t.Fatal("existing-feature bootstrap must not route through full Foundation artifacts")
	}
}

func TestEveryBootstrapIncludesDomainModelingGuidance(t *testing.T) {
	for _, startingPoint := range []string{"new-product", "existing-product", "existing-documents", "existing-feature", "existing-implementation", "audit-only"} {
		bootstrap := bootstrapFor(startingPoint)
		for _, expected := range []string{"Before modeling domains", "every starting point", "examples/events/"} {
			if !strings.Contains(bootstrap, expected) {
				t.Errorf("%s bootstrap missing %q", startingPoint, expected)
			}
		}
	}
}

func TestNewProductBootstrapKeepsFullFoundation(t *testing.T) {
	bootstrap := bootstrapFor("new-product")
	for _, expected := range []string{
		"foundation/problem/problem.md",
		"foundation/vision/vision.md",
		"foundation/strategy/strategy.md",
	} {
		if !strings.Contains(bootstrap, expected) {
			t.Errorf("new-product bootstrap missing %q", expected)
		}
	}
}

func TestExistingImplementationBootstrapStartsWithAssessment(t *testing.T) {
	bootstrap := bootstrapFor("existing-implementation")
	for _, expected := range []string{
		"Implementation Assessment, before canonical Foundation",
		"knowledge/assessments/implementation-assessment.md",
		"foundation/problem/problem.md",
		"foundation/strategy/strategy.md",
	} {
		if !strings.Contains(bootstrap, expected) {
			t.Errorf("existing-implementation bootstrap missing %q", expected)
		}
	}
}

func TestExistingProductBootstrapUsesBaselineAndStrategy(t *testing.T) {
	bootstrap := bootstrapFor("existing-product")
	for _, expected := range []string{
		"Product Baseline, before future Strategy",
		"foundation/product-baseline.md",
		"foundation/strategy/strategy.md",
		"code and operating evidence",
	} {
		if !strings.Contains(bootstrap, expected) {
			t.Errorf("existing-product bootstrap missing %q", expected)
		}
	}
	for _, absent := range []string{"foundation/problem/problem.md", "foundation/vision/vision.md"} {
		if strings.Contains(bootstrap, absent) {
			t.Errorf("existing-product bootstrap should not require consolidated artifact %q", absent)
		}
	}
}

func TestExistingDocumentsBootstrapStartsWithMaterializationGate(t *testing.T) {
	bootstrap := bootstrapFor("existing-documents")
	for _, expected := range []string{
		"Latest import run, before canonical product artifacts",
		"mapping.json",
		"traceability.json",
		"spec-framework import materialize",
		"does not create product approval history",
	} {
		if !strings.Contains(bootstrap, expected) {
			t.Errorf("existing-documents bootstrap missing %q", expected)
		}
	}
	if strings.Contains(bootstrap, "--artifact foundation/problem/problem.md") {
		t.Fatal("existing-documents bootstrap must not skip ahead to Problem approval")
	}
}

func TestAuditOnlyBootstrapUsesNoWriteFlags(t *testing.T) {
	bootstrap := bootstrapFor("audit-only")
	for _, absent := range []string{"--write-registry", "--write-report", "spec-framework approve", "spec-framework work --feature"} {
		if strings.Contains(bootstrap, absent) {
			t.Errorf("audit-only bootstrap contains mutating guidance %q", absent)
		}
	}
	if !strings.Contains(bootstrap, "spec-framework validate --product-root product") || !strings.Contains(bootstrap, "Keep this session read-only") {
		t.Fatal("audit-only bootstrap lacks read-only validation guidance")
	}
}
