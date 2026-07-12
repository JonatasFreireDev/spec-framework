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
	intro := "Begin with the product foundation and follow the canonical gates in order."
	if startingPoint == "existing-documents" {
		intro = "Begin by reviewing the generated import inventory and plan. Source documents are evidence, not approved product artifacts."
	}
	if startingPoint == "existing-product" {
		intro = "Begin by mapping existing product context and decisions into the canonical structure before creating downstream artifacts."
	}
	if startingPoint == "existing-feature" {
		intro = "Begin by validating the feature's Domain, User Goal, scope, and parent approvals before generating Use Cases."
	}
	if startingPoint == "existing-implementation" {
		intro = "Begin with a reverse audit of code, tests, decisions, and documentary gaps. Do not infer approval from implementation."
	}
	if startingPoint == "audit-only" {
		intro = "Begin with read-only gap, conflict, dependency, impact, and consistency audits."
	}
	return fmt.Sprintf(`# Product Bootstrap

Starting point: **%s**

%s

Use this checklist in order. Do not generate downstream artifacts from incomplete or unapproved parents.

## 1. Repository setup

- [ ] Initialize Git and add the first baseline commit.
- [ ] Confirm the versioned spec-framework CLI is available on PATH.
- [ ] Confirm spec-framework validate runs locally.
- [ ] Confirm product/.product/framework.json pins the expected framework version.

## 2. Product identity

- [ ] Replace PRODUCT-TBD and TBD Product in product/context.md.
- [ ] Update product/.product/state.json with the product name and date.
- [ ] Record the intended audience and product owner.

## 3. Problem

- [ ] Complete product/foundation/problem/problem.md with a specific user pain.
- [ ] Add evidence, research, interviews, and opportunities where available.
- [ ] Review the Problem before proceeding to Vision.

## 4. Vision and strategy

- [ ] Complete Vision, principles, north star, and their context.
- [ ] Complete positioning, personas, metrics, roadmap, and Strategy context.
- [ ] Keep downstream artifacts draft until their parent gate is approved.

## 5. Engineering gates

- [ ] Replace commands marked TBD by product adopter in product/knowledge/conventions/gates.md.
- [ ] Complete the security baseline for the chosen stack.
- [ ] Remove gates that do not apply and explain why.

## 6. First product slice

- [ ] Copy product/domains/_template-domain/ to a stable domain slug.
- [ ] Replace template IDs, names, slugs, parents, delivery metadata, and handoffs.
- [ ] Continue Domain -> User Goal -> Domain Evolution -> selected Feature -> Use Case.
- [ ] Use spec-framework work --feature <id-or-path>, then status/next for navigation.
- [ ] Create modular Specification contracts according to rigor.
- [ ] Generate Design, Technical Discovery, resolve the Architecture Gate, then applicable Engineering Proposal and Engineering Review, Plan, Graph, and Tasks.
- [ ] Run spec-framework gates before Code Runner and claim graph tasks when coordinating agents.

## 7. Validate readiness

~~~text
spec-framework validate
spec-framework validate --write-registry --write-report
~~~

Structural ready means paths and contracts are valid. Product readiness additionally requires relevant content, approvals, decisions, gates, and evidence.
`, startingPoint, intro)
}
