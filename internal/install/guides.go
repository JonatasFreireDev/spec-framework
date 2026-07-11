package install

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeStarterGuides(target, version string, agents []Agent) error {
	names := make([]string, len(agents))
	for i, agent := range agents {
		names[i] = string(agent)
	}
	selected := strings.Join(names, ", ")
	flags := strings.Join(names, ",")
	readme := fmt.Sprintf(`# Product Repository

This repository was initialized with Spec Framework %s.

## Start here

1. Read [BOOTSTRAP.md](BOOTSTRAP.md).
2. Replace product identity and starter placeholders under product/.
3. Follow Problem -> Vision -> Strategy before creating the first real Domain.
4. Run spec-framework validate after each documentation step.

## Installed agent integrations

%s

Canonical framework assets live in .spec-framework/. Product-owned artifacts live in product/.

## CLI

~~~text
spec-framework validate
spec-framework validate --write-registry --write-report
spec-framework upgrade --target . --agents %s --yes
spec-framework move --from <old-path> --to <new-path> --dry-run
~~~

A successful validation means the repository is structurally coherent. It does not mean remaining TBD product decisions are complete.
`, version, selected, flags)
	bootstrap := `# Product Bootstrap

Use this checklist in order. Do not generate downstream artifacts from incomplete or unapproved parents.

## 1. Repository setup

- [ ] Initialize Git and add the first baseline commit.
- [ ] Install a versioned spec-framework binary on PATH.
- [ ] Confirm spec-framework validate runs locally.
- [ ] Confirm .github/workflows/framework-validation.yml uses the expected framework version.

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
- [ ] Continue Domain -> User Goal -> Feature -> Use Case -> Specification.
- [ ] Generate Design, Plan, Graph, and Tasks only when their gates permit it.

## 7. Validate readiness

~~~text
spec-framework validate
spec-framework validate --write-registry --write-report
~~~

Structural ready means paths and contracts are valid. Product readiness additionally requires relevant content, approvals, decisions, gates, and evidence.
`
	if err := writeFile(filepath.Join(target, "README.md"), []byte(readme), 0644); err != nil {
		return err
	}
	return writeFile(filepath.Join(target, "BOOTSTRAP.md"), []byte(bootstrap), 0644)
}

func productWorkflow(version string) string {
	if version == "dev" || version == "local" {
		return `name: Framework Validation
on: [pull_request, push]
permissions:
  contents: read
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.26.x"
      - name: Install development Spec Framework CLI
        run: go install github.com/JonatasFreireDev/spec-framework/cmd/spec-framework@master
      - name: Validate framework
        run: spec-framework validate
`
	}
	version = strings.TrimPrefix(version, "v")
	return fmt.Sprintf(`name: Framework Validation
on: [pull_request, push]
permissions:
  contents: read
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install pinned Spec Framework CLI
        env:
          SPEC_FRAMEWORK_VERSION: %q
        run: |
          curl -fsSLO "https://github.com/JonatasFreireDev/spec-framework/releases/download/v${SPEC_FRAMEWORK_VERSION}/spec-framework_${SPEC_FRAMEWORK_VERSION}_linux_amd64.tar.gz"
          curl -fsSLO "https://github.com/JonatasFreireDev/spec-framework/releases/download/v${SPEC_FRAMEWORK_VERSION}/checksums.txt"
          sha256sum --check --ignore-missing checksums.txt
          tar -xzf "spec-framework_${SPEC_FRAMEWORK_VERSION}_linux_amd64.tar.gz"
          sudo install spec-framework /usr/local/bin/spec-framework
      - name: Validate framework
        run: spec-framework validate
`, version)
}
