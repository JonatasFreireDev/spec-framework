#!/usr/bin/env node
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import crypto from "node:crypto";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..", "..");
const validatorScript = path.join(repoRoot, "engineering", "validators", "framework-validator.mjs");
const moveScript = path.join(repoRoot, "engineering", "move-artifact.mjs");
const initProductScript = path.join(repoRoot, "scripts", "init-product.mjs");

const tests = [];

function test(name, fn) {
  tests.push({ name, fn });
}

function mkdir(dir) {
  fs.mkdirSync(dir, { recursive: true });
}

function write(root, relativePath, content) {
  const file = path.join(root, relativePath);
  mkdir(path.dirname(file));
  fs.writeFileSync(file, content.replace(/\r\n?/g, "\n"), "utf8");
}

function copy(root, source, relativePath) {
  const target = path.join(root, relativePath);
  mkdir(path.dirname(target));
  fs.copyFileSync(source, target);
}

function withFixture(name, fn) {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), `spec-framework-${name}-`));
  try {
    scaffoldBase(root);
    fn(root);
  } finally {
    fs.rmSync(root, { recursive: true, force: true });
  }
}

function scaffoldBase(root) {
  write(root, ".product/ids.json", JSON.stringify({
    policy: "slug-scoped",
    deprecated_counters: true,
  }, null, 2));
  write(root, ".product/decisions.json", JSON.stringify({ decisions: [] }, null, 2));
  write(root, ".product/derivations.json", JSON.stringify({ derivations: [] }, null, 2));
  write(root, ".product/artifacts.json", JSON.stringify({ artifacts: [] }, null, 2));
  write(root, "knowledge/templates/minimal-template.md", "# Minimal Template\n\n## Snapshot\n\n| Field | Value |\n| --- | --- |\n| Purpose | Test fixture. |\n");
  mkdir(path.join(root, "audits", "readiness"));
  copy(root, moveScript, "engineering/move-artifact.mjs");
  write(root, ".codex/skills/code-runner/SKILL.md", "---\nname: code-runner\n---\n\n# Code Runner\n");
  write(root, ".codex/skills/task-generator/SKILL.md", "---\nname: task-generator\n---\n\n# Task Generator\n");
}

function runNode(script, cwd, args = []) {
  return spawnSync(process.execPath, [script, ...args], {
    cwd,
    encoding: "utf8",
    windowsHide: true,
  });
}

function runValidator(cwd, args = []) {
  return runNode(validatorScript, cwd, args);
}

function output(result) {
  return `${result.stdout ?? ""}${result.stderr ?? ""}`;
}

function normalizedHash(text) {
  const normalized = text
    .replace(/\r\n?/g, "\n")
    .split("\n")
    .map((line) => line.replace(/[ \t]+$/g, ""))
    .join("\n");
  return crypto.createHash("sha256").update(normalized, "utf8").digest("hex");
}

function writeApprovedArtifact(root, artifact) {
  write(root, artifact.path, artifact.content);
  write(root, `.product/history/approval-${artifact.id}-${artifact.status}.json`, JSON.stringify({
    artifact_id: artifact.id,
    path: artifact.path,
    content_hash: normalizedHash(artifact.content),
    status_granted: artifact.status,
    approved_by: "test-human",
    approved_at: "2026-07-10T00:00:00.000Z",
    notes: "test fixture",
  }, null, 2));
}

function scaffoldTierSUseCase(root, nodes) {
  const dir = "domains/test/goals/goal/features/feature/use-cases/use-case";
  write(root, `${dir}/context.md`, `# Context

\`\`\`yaml
id: UC-TEST
type: use-case
name: Test Use Case
status: draft
slug: use-case
owner_skill: use-case
delivery:
  level: L1
  priority: P1
  depends_on: []
  rationale: Test fixture.
rigor_tier: S
\`\`\`

## Handoff

| Field | Value |
| --- | --- |
| Next skill | task-generator |
`);
  write(root, `${dir}/use-case.md`, `# Use Case

| Field | Value |
| --- | --- |
| ID | UC-TEST |
| Status | draft |
| Delivery Level | L1 |
| Priority | P1 |
| Rationale | Test fixture. |
`);
  write(root, `${dir}/specification.md`, `# Specification

| Field | Value |
| --- | --- |
| ID | SPEC-TEST |
| Status | draft |
| Delivery Level | L1 |
| Priority | P1 |
| Rationale | Test fixture. |
`);
  write(root, `${dir}/tests.md`, `# Tests

| Field | Value |
| --- | --- |
| ID | TEST-TEST |
| Status | draft |
`);
  const graph = {
    id: "GRAPH-TEST",
    status: "draft",
    sourceSpecification: "SPEC-TEST",
    delivery: { level: "L1", priority: "P1", rationale: "Test fixture." },
    nodes,
  };
  write(root, `${dir}/execution-graph.json`, JSON.stringify(graph, null, 2));
  write(root, `${dir}/tasks.md`, `# Tasks

Generated index from execution-graph.json.

${nodes.map((node) => `- [${node.id}](${node.path})`).join("\n")}
`);
  for (const node of nodes) {
    write(root, `${dir}/${node.path}`, `# Task: ${node.title}

## Snapshot

| Field | Value |
| --- | --- |
| ID | ${node.id} |
| Status | draft |
| Source graph | GRAPH-TEST |
| Source specification | SPEC-TEST |
| Source node | ${node.id} |
| Owner skill | task-generator |
| Next skill | code-runner |

## Delivery

| Field | Value |
| --- | --- |
| Level | L1 |
| Priority | P1 |
| Depends on | ${(node.dependsOn ?? []).join(", ") || "none"} |
| Rationale | Test fixture. |

## Task Contract

| Field | Value |
| --- | --- |
| Title | ${node.title} |
| Type | ${node.type} |
| Depends on | ${(node.dependsOn ?? []).join(", ") || "none"} |
| Source sections | Specification |
| Write scope | ${(node.writeScope ?? []).join(", ") || "missing"} |
| Graph node status | pending |
`);
  }
}

test("validator blocks approved artifacts without approval records", () => {
  withFixture("approval-record", (root) => {
    write(root, "approved-audit.md", "# Approved Audit\n\n| Status | approved |\n");
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "AUD-TEST",
          type: "audit",
          status: "approved",
          path: "approved-audit.md",
          documents: { canonical: "approved-audit.md" },
        },
      ],
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /approval-records/);
    assert.match(output(result), /no matching approval record exists/);
  });
});

test("validator blocks proposed descendants when their source hash is stale", () => {
  withFixture("staleness", (root) => {
    const source = "# Source\n\nCurrent source content.\n";
    write(root, "source.md", source);
    write(root, "derived.md", "# Derived\n\nGenerated from source.\n");
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "SRC-TEST",
          type: "audit",
          status: "draft",
          path: "source.md",
          documents: { canonical: "source.md" },
        },
        {
          id: "DER-TEST",
          type: "audit",
          status: "proposed",
          path: "derived.md",
          documents: { canonical: "derived.md" },
        },
      ],
    }, null, 2));
    write(root, ".product/derivations.json", JSON.stringify({
      derivations: [
        {
          artifact_id: "DER-TEST",
          path: "derived.md",
          derived_from: [
            {
              artifact_id: "SRC-TEST",
              path: "source.md",
              content_hash: normalizedHash(`${source}\nchanged elsewhere`),
            },
          ],
        },
      ],
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /staleness/);
    assert.match(output(result), /DER-TEST is stale/);
  });
});

test("validator reports writeScope overlap between parallel nodes as a phase A warning", () => {
  withFixture("write-scope", (root) => {
    scaffoldTierSUseCase(root, [
      {
        id: "TK-A",
        path: "tasks/TK-A.md",
        title: "Task A",
        type: "backend",
        dependsOn: [],
        writeScope: ["src"],
      },
      {
        id: "TK-B",
        path: "tasks/TK-B.md",
        title: "Task B",
        type: "backend",
        dependsOn: [],
        writeScope: ["src/foo.ts"],
      },
    ]);

    const result = runValidator(root, ["--write-registry"]);

    assert.equal(result.status, 0, output(result));
    assert.match(output(result), /WARNING write-scope/);
    assert.match(output(result), /Parallel nodes TK-A and TK-B have overlapping writeScope/);
  });
});

test("validator blocks task handoffs that reference missing skills", () => {
  withFixture("skill-reference", (root) => {
    scaffoldTierSUseCase(root, [
      {
        id: "TK-A",
        path: "tasks/TK-A.md",
        title: "Task A",
        type: "backend",
        dependsOn: [],
        writeScope: ["src/task-a"],
      },
    ]);
    const taskFile = path.join(root, "domains/test/goals/goal/features/feature/use-cases/use-case/tasks/TK-A.md");
    fs.writeFileSync(
      taskFile,
      fs.readFileSync(taskFile, "utf8").replace("| Next skill | code-runner |", "| Next skill | missing-skill |"),
      "utf8"
    );

    const result = runValidator(root, ["--write-registry"]);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /skill-reference/);
    assert.match(output(result), /Next skill references missing skill missing-skill/);
  });
});

test("validator blocks approved QA evidence with placeholder gate output", () => {
  withFixture("qa-evidence", (root) => {
    const qa = `# QA Evidence

| Field | Value |
| --- | --- |
| ID | QA-TEST |
| Status | approved |

## Gate Evidence

| Field | Value |
| --- | --- |
| Test command | N/A until validation |
| Gate logs | N/A until validation |
| CI URL | N/A |
| Screenshots | N/A |
| Environment | N/A |
| Limitations | N/A |

## QA Verdict

| Field | Value |
| --- | --- |
| Verdict | blocked |
`;
    write(root, "qa-evidence.md", qa);
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "QA-TEST",
          type: "qa-evidence",
          status: "approved",
          path: "qa-evidence.md",
          documents: { canonical: "qa-evidence.md" },
        },
      ],
    }, null, 2));
    write(root, ".product/history/approval-QA-TEST-approved.json", JSON.stringify({
      artifact_id: "QA-TEST",
      path: "qa-evidence.md",
      content_hash: normalizedHash(qa),
      status_granted: "approved",
      approved_by: "test-human",
      approved_at: "2026-07-10T00:00:00.000Z",
      notes: "test fixture",
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /qa-evidence/);
    assert.match(output(result), /no real gate output or limitation is recorded/);
  });
});

test("validator blocks approved findings without route and owner", () => {
  withFixture("failure-routing", (root) => {
    const audit = `# Audit

| Field | Value |
| --- | --- |
| ID | AUD-ROUTE |
| Status | approved |

## Findings

| Severity | Finding | Evidence | Required Fix |
| --- | --- | --- | --- |
| blocker | Escaped defect | test.log | Fix root cause |
`;
    write(root, "audit.md", audit);
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "AUD-ROUTE",
          type: "audit",
          status: "approved",
          path: "audit.md",
          documents: { canonical: "audit.md" },
        },
      ],
    }, null, 2));
    write(root, ".product/history/approval-AUD-ROUTE-approved.json", JSON.stringify({
      artifact_id: "AUD-ROUTE",
      path: "audit.md",
      content_hash: normalizedHash(audit),
      status_granted: "approved",
      approved_by: "test-human",
      approved_at: "2026-07-10T00:00:00.000Z",
      notes: "test fixture",
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /failure-routing/);
    assert.match(output(result), /missing Route/);
    assert.match(output(result), /missing Owner/);
  });
});

test("validator blocks validated use cases without approved Code Review", () => {
  withFixture("code-review-gate", (root) => {
    const artifacts = [
      { id: "UC-VAL", type: "use-case", status: "validated", path: "uc.md", content: "# Use Case\n\n| Status | validated |\n" },
      { id: "TEST-VAL", type: "tests", status: "approved", path: "tests.md", content: "# Tests\n\n| ID | TEST-VAL |\n| Status | approved |\n" },
      { id: "QA-VAL", type: "qa-evidence", status: "approved", path: "qa.md", content: "# QA Evidence\n\n| ID | QA-VAL |\n| Status | approved |\n\n## Gate Evidence\n\n| Field | Value |\n| --- | --- |\n| Test command | npm test |\n| Gate logs | tests.log |\n| Environment | CI |\n| Limitations | N/A |\n\n## QA Verdict\n\n| Field | Value |\n| --- | --- |\n| Verdict | passed |\n" },
      { id: "SEC-VAL", type: "security-review", status: "approved", path: "security.md", content: "# Security Review\n\n| ID | SEC-VAL |\n| Status | approved |\n" },
      { id: "AUD-VAL", type: "audit", status: "approved", path: "audit.md", content: "# Audit\n\n| ID | AUD-VAL |\n| Status | approved |\n" },
    ];
    for (const artifact of artifacts) writeApprovedArtifact(root, artifact);
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: artifacts.map((artifact) => ({
        id: artifact.id,
        type: artifact.type,
        status: artifact.status,
        path: artifact.path,
        parentIds: artifact.id === "UC-VAL" ? [] : ["UC-VAL"],
        documents: { canonical: artifact.path },
      })),
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /validation-gates/);
    assert.match(output(result), /code-review\.md/);
  });
});

test("validator blocks approved Code Review required fixes without route and owner", () => {
  withFixture("code-review-quality", (root) => {
    const review = `# Code Review

| Field | Value |
| --- | --- |
| ID | CR-TEST |
| Status | approved |

## Findings

| Severity | Finding | Evidence | Required Fix |
| --- | --- | --- | --- |
| required_fix | Missing edge handling | src/app.ts:10 | Handle empty state |

## Review Verdict

| Field | Value |
| --- | --- |
| Verdict | passed |
| Completeness passed | yes |
| Adherence passed | yes |
| Quality passed | yes |
`;
    writeApprovedArtifact(root, { id: "CR-TEST", type: "code-review", status: "approved", path: "code-review.md", content: review });
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "CR-TEST",
          type: "code-review",
          status: "approved",
          path: "code-review.md",
          documents: { canonical: "code-review.md" },
        },
      ],
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /failure-routing/);
    assert.match(output(result), /Missing edge handling/);
  });
});

test("validator blocks implemented tasks without traceable commit references", () => {
  withFixture("commit-reference", (root) => {
    write(root, "src/app.ts", "export const value = 1;\n");
    const task = `# Task

| Field | Value |
| --- | --- |
| ID | TK-COMMIT |
| Status | implemented |

## Implementation Links

| Field | Value |
| --- | --- |
| Branch | feature/test |
| Commits | not-a-hash |
| PR | N/A until validation |
| Code paths | src/app.ts |
`;
    writeApprovedArtifact(root, { id: "TK-COMMIT", type: "task", status: "implemented", path: "task.md", content: task });
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "TK-COMMIT",
          type: "task",
          status: "implemented",
          path: "task.md",
          documents: { canonical: "task.md" },
        },
      ],
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /code-evidence/);
    assert.match(output(result), /traceable commit hash or commit URL/);
  });
});

test("validator blocks validated tasks without traceable PR references", () => {
  withFixture("pr-reference", (root) => {
    write(root, "src/app.ts", "export const value = 1;\n");
    const task = `# Task

| Field | Value |
| --- | --- |
| ID | TK-PR |
| Status | validated |

## Implementation Links

| Field | Value |
| --- | --- |
| Branch | feature/test |
| Commits | abc1234 |
| PR | review surface pending |
| Code paths | src/app.ts |

## Validation Evidence

| Field | Value |
| --- | --- |
| Test status | passed |
| Gate logs | test.log |
| CI URL | N/A |
| Screenshots | N/A |
| QA evidence | qa.md |
`;
    writeApprovedArtifact(root, { id: "TK-PR", type: "task", status: "validated", path: "task.md", content: task });
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "TK-PR",
          type: "task",
          status: "validated",
          path: "task.md",
          documents: { canonical: "task.md" },
        },
      ],
    }, null, 2));

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /code-evidence/);
    assert.match(output(result), /not a traceable PR reference/);
  });
});

test("validator blocks broken relative Markdown links", () => {
  withFixture("markdown-link", (root) => {
    write(root, "docs/index.md", "# Docs\n\n[Missing](missing.md)\n");

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /links/);
    assert.match(output(result), /Broken Markdown link: missing\.md/);
  });
});

test("validator checks real Markdown links inside templates", () => {
  withFixture("template-link", (root) => {
    write(root, "knowledge/templates/template-with-link.md", "# Template\n\n## Snapshot\n\n[Broken](missing-template-target.md)\n");

    const result = runValidator(root);

    assert.notEqual(result.status, 0, output(result));
    assert.match(output(result), /knowledge\/templates\/template-with-link\.md/);
    assert.match(output(result), /Broken Markdown link: missing-template-target\.md/);
  });
});

test("move-artifact moves folders, rewrites Markdown links and JSON paths, and reports free-text mentions", () => {
  withFixture("move", (root) => {
    write(root, "domains/old/use-case/file.md", "# Target\n\nMoved content.\n");
    write(root, "docs/link.md", "# Link\n\n[Target](../domains/old/use-case/file.md#target)\n");
    write(root, "docs/free-text.md", "Review domains/old/use-case before release.\n");
    write(root, ".product/artifacts.json", JSON.stringify({
      artifacts: [
        {
          id: "UC-MOVE",
          path: "domains/old/use-case",
          documents: { canonical: "domains/old/use-case/file.md" },
        },
      ],
    }, null, 2));

    const result = runNode(moveScript, root, [
      "--from",
      "domains/old/use-case",
      "--to",
      "domains/new/use-case",
    ]);

    assert.equal(result.status, 0, output(result));
    assert.equal(fs.existsSync(path.join(root, "domains/old/use-case")), false);
    assert.equal(fs.existsSync(path.join(root, "domains/new/use-case/file.md")), true);
    assert.match(fs.readFileSync(path.join(root, "docs/link.md"), "utf8"), /\.\.\/domains\/new\/use-case\/file\.md#target/);
    assert.match(fs.readFileSync(path.join(root, ".product/artifacts.json"), "utf8"), /domains\/new\/use-case\/file\.md/);
    assert.match(output(result), /Rewritten files: 2/);
    assert.match(output(result), /Free-text mentions requiring review: 1/);
    assert.match(output(result), /docs\/free-text\.md:1/);
  });
});

test("init-product creates a product repo with installed framework assets", () => {
  const parent = fs.mkdtempSync(path.join(os.tmpdir(), "spec-framework-init-"));
  const target = path.join(parent, "new-product");
  try {
    const result = runNode(initProductScript, repoRoot, ["--target", target]);

    assert.equal(result.status, 0, output(result));
    assert.equal(fs.existsSync(path.join(target, "product", ".product", "framework.json")), true);
    assert.equal(fs.existsSync(path.join(target, ".spec-framework", "FRAMEWORK.md")), true);
    assert.equal(fs.existsSync(path.join(target, ".spec-framework", "validators", "framework-validator.mjs")), true);
    assert.equal(fs.existsSync(path.join(target, ".spec-framework", "tools", "move-artifact.mjs")), true);
    assert.equal(fs.existsSync(path.join(target, ".spec-framework", "skills", "code-runner", "SKILL.md")), true);
    assert.equal(fs.existsSync(path.join(target, ".codex", "skills", "code-runner", "SKILL.md")), true);

    const manifest = JSON.parse(fs.readFileSync(path.join(target, ".spec-framework", "manifest.json"), "utf8"));
    assert.equal(manifest.product_root, "product");
    assert.equal(manifest.installed_assets.validators, true);

    const installedValidator = path.join(target, ".spec-framework", "validators", "framework-validator.mjs");
    const validation = runNode(installedValidator, target, ["--product-root", "product", "--framework-root", ".spec-framework"]);
    assert.equal(validation.status, 0, output(validation));
    assert.match(output(validation), /ready/);
  } finally {
    fs.rmSync(parent, { recursive: true, force: true });
  }
});

let passed = 0;
for (const item of tests) {
  try {
    item.fn();
    passed += 1;
    console.log(`ok - ${item.name}`);
  } catch (error) {
    console.error(`not ok - ${item.name}`);
    console.error(error.stack || error.message);
    process.exitCode = 1;
  }
}

console.log(`${passed}/${tests.length} tests passed`);
