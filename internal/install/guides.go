package install

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeStarterGuides(target, version string, agents []Agent, startingPoint string) error {
	names := make([]string, len(agents))
	for i, agent := range agents {
		names[i] = string(agent)
	}
	selected := strings.Join(names, ", ")
	bootstrap := bootstrapFor(startingPoint)
	header := fmt.Sprintf("Framework version: **%s**\n\nConfigured agents: **%s**\n\n", version, selected)
	return writeFile(filepath.Join(target, "product", "BOOTSTRAP.md"), []byte(header+bootstrap), 0644)
}

func bootstrapFor(startingPoint string) string {
	profile := bootstrapProfileFor(startingPoint)
	validationCommand := "spec-framework validate --product-root product --write-registry --write-report"
	workspaceGuidance := `After a stable Domain, User Goal, and Feature exist in the registry:

~~~text
spec-framework work --feature <id-or-path> --created-by "Your Name"
spec-framework guide --work WORK-001
spec-framework dashboard --work WORK-001
~~~`
	if startingPoint == "audit-only" {
		validationCommand = "spec-framework validate --product-root product"
		workspaceGuidance = "Audit-only does not create a WORK-NNN workspace. Inspect existing state and report findings in terminal output."
	}
	return fmt.Sprintf(`# Product Bootstrap

Starting point: **%s**

## Start here

| Question | Answer |
| --- | --- |
| Where am I? | **%s** |
| What is ready? | The framework structure and runtime are installed. Product content is still draft. |
| What should I do next? | **%s** |
| Recommended working style | **%s** |

> spec-framework validate can pass while the product is still empty. Structural validity means files and contracts are well-formed; it does not mean product decisions are approved or implementation-ready.

## First session

1. Confirm the repository baseline and framework pin.

   ~~~text
   git status
   spec-framework version
   %s
   ~~~

2. Read the current evidence before writing product truth.

   - product/context.md
   - product/.product/framework.json
   - product/audits/framework-validation-report.md
   - For imported documents: the latest run under product/knowledge/imports/runs/

3. Choose the amount of discovery, not a weaker approval standard.

   | Style | Use when | Expected depth |
   | --- | --- | --- |
   | **Lean** | A brief, feature, or implementation already makes the intended outcome clear. | Concise, scoped Foundation contracts with evidence and explicit boundaries. |
   | **Full** | The product, audience, problem, or strategic bet is still uncertain. | Research, alternatives, personas, metrics, roadmap, and broader discovery. |

   %s

4. Complete only the next artifact. Do not change its status by hand.

   %s

5. %s

   %s

   %s

## Prompts for the agent

Use these prompts as a starting point. The agent must read the cited evidence first, propose changes in the correct template, show unresolved questions, and stop before approval or materialization.

%s

## Foundation path

%s

Parent approvals are enforced mechanically. If a command is blocked, read the reported parent instead of editing statuses manually.

## Before the first workspace

status, next, guide, and dashboard require a WORK-NNN workspace. Not having one during Foundation is expected.

## Before modeling domains

For every starting point that creates or revises domains, read the pinned framework runtime's examples/events/ before the first domain change. It is the canonical modeling reference: a domain is a coherent business area, not the product name, a sidebar section, or a container for every capability. Use its domains/events/domain.md and domains/README.md to model explicit ownership, Does Not Own boundaries, cross-domain dependencies, and one walking-skeleton chain: Domain -> User Goal -> Feature -> Use Case. In audit-only mode, use the same reference to assess existing domain boundaries without changing them.

Do not put authentication into an unrelated business domain; model a users/identity boundary or record why the product genuinely has one identity domain. Do not stop at domain.md: create the first goal, feature, and use case before creating a workspace.

%s

## Engineering readiness

- Replace TBD commands in product/knowledge/conventions/gates.md before implementation.
- Complete the security baseline for the chosen stack.
- Keep downstream artifacts draft until their parent gates have current approval evidence.
- Run spec-framework gates before Code Runner.

## Check progress

~~~text
%s
~~~

Then ask:

1. Which artifacts are still draft or proposed?
2. Which approved statuses have matching current history records?
3. What is the first blocked parent?
4. Is the chosen Foundation scope lean, full, or feature-scoped?
`, startingPoint, profile.location, profile.nextAction, profile.style, validationCommand, profile.scopeRule, profile.artifactAction, profile.approvalIntro, profile.approvalGuidance, profile.approvalRule, profile.agentPrompts, profile.foundationPath, workspaceGuidance, validationCommand)
}

type bootstrapProfile struct {
	location         string
	nextAction       string
	style            string
	scopeRule        string
	artifactAction   string
	approvalIntro    string
	approvalGuidance string
	approvalRule     string
	foundationPath   string
	agentPrompts     string
}

const fullFoundationScope = "Lean keeps Problem, Vision, Principles, North Star, and Strategy proportional to the available evidence. It reduces depth, not approval integrity."

const fullFoundationApproval = `~~~text
   spec-framework approve --product-root product --artifact foundation/problem/problem.md --approved-by "Your Name"
   spec-framework approve --product-root product --artifact foundation/problem/problem.md --approved-by "Your Name" --yes
   ~~~`

const artifactApprovalIntro = "Preview the approval, review the exact artifact and hash, then apply it explicitly."

const artifactApprovalRule = "Approval is valid only when the CLI writes a matching record under product/.product/history/. A Markdown status edit alone is not approval."

const fullFoundationPath = `Approve one artifact at a time in this order:

| Step | Artifact | What it establishes |
| --- | --- | --- |
| 1 | foundation/problem/problem.md | The evidenced pain or opportunity. |
| 2 | foundation/vision/vision.md | The intended outcome and boundaries. |
| 3 | foundation/vision/principles.md | Rules and trade-offs for decisions. |
| 4 | foundation/vision/north-star.md | Value signal and guardrails. |
| 5 | foundation/strategy/strategy.md | The scoped delivery bet. |`

func bootstrapProfileFor(startingPoint string) bootstrapProfile {
	profiles := map[string]bootstrapProfile{
		"new-product": {
			location:         "L0 Foundation: product identity, then Problem",
			nextAction:       "Replace PRODUCT-TBD in product/context.md, then complete foundation/problem/problem.md",
			style:            "Full until the problem and audience are clear",
			scopeRule:        fullFoundationScope,
			artifactAction:   "Start with `product/context.md`, then write an evidence-based `product/foundation/problem/problem.md`.",
			approvalIntro:    artifactApprovalIntro,
			approvalGuidance: fullFoundationApproval,
			approvalRule:     artifactApprovalRule,
			foundationPath:   fullFoundationPath,
		},
		"existing-product": {
			location:       "Product Baseline, before future Strategy",
			nextAction:     "Complete foundation/product-baseline.md from repository, runtime, user, and operational evidence",
			style:          "Lean when existing evidence is reliable; Full where decisions are unclear",
			scopeRule:      "Existing-product consolidates current Problem, Vision, Principles, and North Star evidence into Product Baseline, while keeping future Strategy separate. Escalate to full Foundation when audience, value, or direction is uncertain.",
			artifactAction: "Complete `product/foundation/product-baseline.md` from code and operating evidence; label inferred intent and unknowns explicitly.",
			approvalIntro:  artifactApprovalIntro,
			approvalGuidance: `~~~text
   spec-framework approve --product-root product --artifact foundation/product-baseline.md --approved-by "Your Name"
   spec-framework approve --product-root product --artifact foundation/product-baseline.md --approved-by "Your Name" --yes
			   ~~~`,
			approvalRule: artifactApprovalRule,
			foundationPath: `Approve the two artifacts individually in this order:

| Step | Artifact | What it establishes |
| --- | --- | --- |
| 1 | foundation/product-baseline.md | The evidenced product that exists today. |
| 2 | foundation/strategy/strategy.md | The future bets, trade-offs, priorities, and metrics. |`,
		},
		"existing-documents": {
			location:       "Latest import run, before canonical product artifacts",
			nextAction:     "Review inventory, conflicts, and selected mappings, then explicitly materialize them as drafts",
			style:          "Lean when the sources already contain a coherent brief",
			scopeRule:      "Existing-documents uses the latest import run as its entry contract. Materialization approves selected draft writes, not the resulting product artifacts.",
			artifactAction: "Review `product/knowledge/imports/runs/<latest-run>/inventory.json`, `traceability.json`, `conflicts.md`, and `mapping.json`. Ask the Artifact Importer agent to read every source and record covered claims and unmapped gaps before materialization. Sources are evidence, not approved product truth.",
			approvalIntro:  "Review every selected mapping, target, source reference, conflict, and draft body; then authorize draft materialization explicitly.",
			approvalGuidance: `~~~text
   spec-framework import materialize --product-root product --run <IMPORT-NNN> --approved-by "Your Name" --yes
   ~~~`,
			approvalRule:   "Materialization records who authorized the selected draft writes in import-plan.json. It does not create product approval history.",
			foundationPath: "After materialization, route each draft through its normal owner and parent gates. Use the full Foundation path unless the human explicitly selects another supported starting-point contract.",
		},
		"existing-feature": {
			location:       "Feature Brief, before the first workspace",
			nextAction:     "Complete and approve foundation/feature-brief.md for this bounded feature",
			style:          "Lean and feature-scoped",
			scopeRule:      "Existing-feature replaces the full product Foundation package with one Feature Brief. Escalate to full Foundation when product direction or scope is broad or uncertain.",
			artifactAction: "Complete `product/foundation/feature-brief.md`; do not invent a product-wide strategy. After approval, model only the bounded feature's domain slice using the pinned runtime's examples/events/ reference.",
			approvalIntro:  artifactApprovalIntro,
			approvalGuidance: `~~~text
   spec-framework approve --product-root product --artifact foundation/feature-brief.md --approved-by "Your Name"
   spec-framework approve --product-root product --artifact foundation/feature-brief.md --approved-by "Your Name" --yes
			   ~~~`,
			approvalRule:   artifactApprovalRule,
			foundationPath: "Approve foundation/feature-brief.md before creating the first WORK-NNN workspace. Its approval covers only the bounded feature described there.",
		},
		"existing-implementation": {
			location:       "Implementation Assessment, before canonical Foundation",
			nextAction:     "Complete and approve knowledge/assessments/implementation-assessment.md from observed evidence",
			style:          "Lean where implementation evidence is strong",
			scopeRule:      "Existing-implementation adds an Implementation Assessment before the full Foundation path. Observed code is evidence, not approved product intent.",
			artifactAction: "Complete `product/knowledge/assessments/implementation-assessment.md` without changing application code or inventing product decisions.",
			approvalIntro:  artifactApprovalIntro,
			approvalGuidance: `~~~text
   spec-framework approve --product-root product --artifact knowledge/assessments/implementation-assessment.md --approved-by "Your Name"
   spec-framework approve --product-root product --artifact knowledge/assessments/implementation-assessment.md --approved-by "Your Name" --yes
			   ~~~`,
			approvalRule:   artifactApprovalRule,
			foundationPath: "Approve knowledge/assessments/implementation-assessment.md first. Then continue with the full Foundation path:\n\n" + fullFoundationPath,
		},
		"audit-only": {
			location:         "Read-only audit entry point",
			nextAction:       "Run validation and inspect gaps without advancing artifact statuses",
			style:            "Audit only; do not manufacture Foundation approvals",
			scopeRule:        "Audit-only inspects structural and evidentiary gaps without creating or approving product truth.",
			artifactAction:   "Keep product artifacts unchanged unless the human explicitly converts an audit finding into scoped Foundation work.",
			approvalIntro:    "Keep this session read-only.",
			approvalGuidance: "Do not run approval commands in audit-only mode.",
			approvalRule:     "Audit findings are evidence and do not alter product approval history.",
			foundationPath:   "No Foundation path advances during an audit-only session. Report gaps and stop before product mutation.",
		},
	}
	if profile, ok := profiles[startingPoint]; ok {
		profile.agentPrompts = agentPromptsFor(startingPoint)
		return profile
	}
	profile := profiles["new-product"]
	profile.agentPrompts = agentPromptsFor("new-product")
	return profile
}

func agentPromptsFor(startingPoint string) string {
	switch startingPoint {
	case "new-product":
		return `### Product identity
> Read ` + "`product/context.md`" + `. Ask me for the product name, owner, audience, and first business area. Propose only the smallest edits needed to replace ` + "`TBD`" + ` placeholders. Do not invent facts.

### Problem
> Read ` + "`product/context.md`" + `, the Problem context, and any evidence I provide. Ask focused questions about audience, pain, frequency, urgency, alternatives, and evidence. Draft ` + "`foundation/problem/problem.md`" + ` using the template, label assumptions, and list unanswered questions. Stop before approval.

### Foundation sequence
> After I confirm the previous artifact, read its current content and parent approval state. Propose the next artifact only: Vision, Product Principles, North Star, or Strategy. Preserve traceability to the parent, call out contradictions, and never change status or create approval records.`
	case "existing-product":
		return `### Product Baseline
> Read the repository, tests, runtime configuration, operational evidence, and ` + "`product/context.md`" + `. Propose ` + "`foundation/product-baseline.md`" + ` with observed facts separated from inferred intent and unknowns. Cite evidence paths. Do not treat existing code as approved product scope.

### Strategy
> Read the approved Product Baseline and propose ` + "`foundation/strategy/strategy.md`" + ` with future bets, trade-offs, metrics, and roadmap. Ask before making assumptions about direction. Stop before approval.`
	case "existing-documents":
		return `### Source review and traceability
> Read every source in ` + "`product/knowledge/imports/sources/`" + ` and the current ` + "`inventory.json`" + `. For each source, update ` + "`traceability.json`" + ` with review status, section-level evidence, extracted claims, candidate ids, mapped targets, and gaps. Do not silently discard content.

### Draft proposals
> Compare the traced claims with existing product artifacts and templates. Propose ` + "`mapping.json`" + ` entries for useful drafts, preserving ` + "`source_documents`" + ` references. Record conflicts and ambiguous ownership in ` + "`conflicts.md`" + `. Do not materialize or approve anything.

### Human review
> Summarize what is covered, what is unmapped, which sources conflict, and which mappings are safe to materialize as drafts. Ask for explicit selection before running the materialization command.`
	case "existing-feature":
		return `### Feature Brief
> Read the request, existing product context, relevant decisions, and nearby artifacts. Ask questions about the user, outcome, scope, non-goals, constraints, success signal, and delivery level. Draft ` + "`foundation/feature-brief.md`" + ` for this bounded feature only. Flag when the request actually needs full Foundation. Stop before approval.

### Bounded delivery slice
> After the Feature Brief is approved, read the pinned ` + "`examples/events/`" + ` reference and choose one Domain -> User Goal -> Feature -> Use Case walking skeleton that fits the brief. Create only the smallest draft slice and preserve parent traceability.`
	case "existing-implementation":
		return `### Implementation Assessment
> Inspect the code, configuration, database schema, integrations, tests, and operational evidence. Fill ` + "`knowledge/assessments/implementation-assessment.md`" + ` with observed behavior, architecture, risks, and unknowns. Separate evidence from inferred product intent. Do not modify application code or create approval records.

### Derive product truth carefully
> Use the approved assessment as evidence for the full Foundation. Propose Problem, Vision, Principles, North Star, and Strategy one artifact at a time, asking questions where the implementation does not prove user intent.`
	case "audit-only":
		return `### Read-only audit
> Read ` + "`product/BOOTSTRAP.md`" + `, the manifest, validation output, contexts, import runs, decisions, and existing artifacts. Report structural gaps, stale evidence, missing traceability, conflicts, and approval inconsistencies. Do not edit product files, materialize imports, approve artifacts, or create workspaces.

### Human handoff
> Summarize the first safe remediation, the owning skill, required evidence, and the exact command or human decision needed. Keep the session read-only.`
	default:
		return agentPromptsFor("new-product")
	}
}
