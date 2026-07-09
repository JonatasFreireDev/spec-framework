#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";

const root = process.cwd();
const args = new Set(process.argv.slice(2));
const writeReport = args.has("--write-report");

const allowedStatuses = new Set([
  "draft",
  "proposed",
  "approved",
  "in_progress",
  "implemented",
  "validated",
  "released",
  "deprecated",
  "superseded",
]);

const requiredUseCaseFiles = [
  "context.md",
  "use-case.md",
  "specification.md",
  "design.md",
  "implementation-plan.md",
  "execution-graph.json",
  "tasks.md",
  "tests.md",
  "analytics.md",
  "audit.md",
];

const results = [];

function rel(filePath) {
  return path.relative(root, filePath).replaceAll(path.sep, "/");
}

function addResult(severity, check, file, message, fix = "") {
  results.push({
    severity,
    check,
    file: file ? rel(file) : "",
    message,
    fix,
  });
}

function walk(dir, output = []) {
  if (!fs.existsSync(dir)) return output;
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (entry.name === ".git") continue;
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walk(full, output);
    } else {
      output.push(full);
    }
  }
  return output;
}

function readText(file) {
  return fs.readFileSync(file, "utf8").replace(/^\uFEFF/, "");
}

function parseJsonFile(file) {
  try {
    return { ok: true, value: JSON.parse(readText(file)) };
  } catch (error) {
    return { ok: false, error };
  }
}

function findUseCaseDirs() {
  const dirs = [];
  const domainsDir = path.join(root, "domains");
  if (!fs.existsSync(domainsDir)) return dirs;

  function visit(dir) {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const full = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        if (path.basename(path.dirname(full)) === "use-cases") {
          dirs.push(full);
        }
        visit(full);
      }
    }
  }

  visit(domainsDir);
  return dirs;
}

function parseContextMeta(text) {
  const match = text.match(/```yaml\s+([\s\S]*?)```/);
  if (!match) return null;
  const meta = {};
  for (const rawLine of match[1].split(/\r?\n/)) {
    const line = rawLine.trim();
    const simple = line.match(/^([A-Za-z0-9_-]+):\s*(.+?)\s*$/);
    if (simple) meta[simple[1]] = simple[2].replace(/^["']|["']$/g, "");
  }
  return meta;
}

function pathExistsMaybeProductPrefix(value) {
  const direct = path.join(root, value);
  if (fs.existsSync(direct)) return { exists: true, normalized: value };
  if (value.startsWith("product/")) {
    const stripped = value.slice("product/".length);
    if (fs.existsSync(path.join(root, stripped))) {
      return { exists: true, normalized: stripped, stalePrefix: true };
    }
  }
  return { exists: false, normalized: value };
}

function validateUseCaseBundles() {
  for (const dir of findUseCaseDirs()) {
    for (const fileName of requiredUseCaseFiles) {
      const file = path.join(dir, fileName);
      if (!fs.existsSync(file)) {
        addResult(
          "error",
          "use-case-bundle",
          dir,
          `Missing required use-case artifact: ${fileName}`,
          `Create ${fileName} or explain why the use case is not executable yet.`
        );
      }
    }
  }
}

function validateExecutionGraphs() {
  for (const file of walk(path.join(root, "domains")).filter((item) =>
    item.endsWith("execution-graph.json")
  )) {
    const parsed = parseJsonFile(file);
    if (!parsed.ok) {
      addResult("error", "execution-graph", file, `Invalid JSON: ${parsed.error.message}`, "Fix JSON syntax.");
      continue;
    }

    const graph = parsed.value;
    if (!graph.id) addResult("error", "execution-graph", file, "Graph is missing id.", "Add graph id.");
    if (!graph.sourceSpecification) {
      addResult("error", "execution-graph", file, "Graph is missing sourceSpecification.", "Link graph to a specification id.");
    }
    if (!Array.isArray(graph.nodes)) {
      addResult("error", "execution-graph", file, "Graph nodes must be an array.", "Add nodes array.");
      continue;
    }

    const ids = new Set();
    for (const node of graph.nodes) {
      if (!node.id) {
        addResult("error", "execution-graph", file, "A graph node is missing id.", "Add node id.");
        continue;
      }
      if (ids.has(node.id)) {
        addResult("error", "execution-graph", file, `Duplicate graph node id: ${node.id}`, "Use unique task ids.");
      }
      ids.add(node.id);
      for (const field of ["title", "type", "dependsOn"]) {
        if (!(field in node)) {
          addResult("error", "execution-graph", file, `Node ${node.id} is missing ${field}.`, `Add ${field}.`);
        }
      }
      if (!("ownerSkill" in node)) {
        addResult("warning", "execution-graph", file, `Node ${node.id} is missing ownerSkill.`, "Assign the responsible skill.");
      }
      if (!("writeScope" in node)) {
        addResult("warning", "execution-graph", file, `Node ${node.id} is missing writeScope.`, "Declare write scope.");
      }
      if (!("acceptanceChecks" in node)) {
        addResult("warning", "execution-graph", file, `Node ${node.id} is missing acceptanceChecks.`, "Add acceptanceChecks.");
      }
      if (Array.isArray(node.dependsOn)) {
        for (const dependency of node.dependsOn) {
          if (!ids.has(dependency) && !graph.nodes.some((candidate) => candidate.id === dependency)) {
            addResult(
              "error",
              "execution-graph",
              file,
              `Node ${node.id} depends on missing node ${dependency}.`,
              "Fix dependsOn or add the missing node."
            );
          }
        }
      }
    }
  }
}

function validateContexts() {
  for (const file of walk(path.join(root, "domains")).filter((item) => path.basename(item) === "context.md")) {
    const meta = parseContextMeta(readText(file));
    if (!meta) {
      addResult("error", "context", file, "context.md is missing yaml metadata block.", "Add a yaml metadata block.");
      continue;
    }
    for (const field of ["id", "type", "name", "status", "owner_skill"]) {
      if (!meta[field]) addResult("error", "context", file, `Missing context field: ${field}.`, `Add ${field}.`);
    }
    if (meta.status && !allowedStatuses.has(meta.status)) {
      addResult("error", "context", file, `Invalid status: ${meta.status}.`, "Use a framework-approved status.");
    }
    if (!readText(file).includes("## Handoff")) {
      addResult("warning", "context", file, "Missing Handoff section.", "Add next skill and required reading.");
    }
  }
}

function validateProductPrefixLinks(files) {
  const staleProductPath = /(^|[^.\w-])product\/(domains|knowledge|audits|foundation|design|engineering|releases|skills|FRAMEWORK\.md)/;
  for (const file of files.filter((item) => /\.(md|json)$/.test(item))) {
    const fileRel = rel(file);
    if (fileRel === "FRAMEWORK.md" || fileRel === "audits/framework-validation-report.md") continue;
    const text = readText(file);
    if (staleProductPath.test(text)) {
      addResult(
        "warning",
        "paths",
        file,
        "Found product/ path prefix outside FRAMEWORK.md.",
        "Use repository-root-relative paths unless quoting FRAMEWORK.md."
      );
    }
  }
}

function validateDecisionsIndex() {
  const file = path.join(root, ".product", "decisions.json");
  const parsed = parseJsonFile(file);
  if (!parsed.ok) {
    addResult("error", "decisions-index", file, `Invalid decisions index JSON: ${parsed.error.message}`, "Fix JSON syntax.");
    return;
  }
  const decisions = parsed.value.decisions;
  if (!Array.isArray(decisions)) {
    addResult("error", "decisions-index", file, "decisions must be an array.", "Add decisions array.");
    return;
  }
  for (const decision of decisions) {
    if (!decision.id) addResult("error", "decisions-index", file, "Decision entry missing id.", "Add id.");
    if (!decision.status) addResult("error", "decisions-index", file, `Decision ${decision.id ?? "(unknown)"} missing status.`, "Add status.");
    if (decision.path) {
      const pathCheck = pathExistsMaybeProductPrefix(decision.path);
      if (!pathCheck.exists) {
        addResult("error", "decisions-index", file, `Decision ${decision.id} path does not exist: ${decision.path}`, "Fix path.");
      } else if (pathCheck.stalePrefix) {
        addResult("warning", "decisions-index", file, `Decision ${decision.id} path uses stale product/ prefix.`, `Use ${pathCheck.normalized}.`);
      }
    }
    for (const artifact of decision.affectedArtifacts ?? []) {
      const pathCheck = pathExistsMaybeProductPrefix(artifact);
      if (!pathCheck.exists) {
        addResult("error", "decisions-index", file, `Affected artifact does not exist: ${artifact}`, "Fix or remove affected artifact path.");
      } else if (pathCheck.stalePrefix) {
        addResult("warning", "decisions-index", file, `Affected artifact uses stale product/ prefix: ${artifact}`, `Use ${pathCheck.normalized}.`);
      }
    }
  }
}

function validateMermaidAndTemplates(files) {
  for (const file of files.filter((item) => item.endsWith(".md"))) {
    const text = readText(file);
    if (/```mermaid[\s\S]*?flowchart/.test(text) && !text.includes("classDef done")) {
      addResult("warning", "mermaid", file, "Flowchart is missing Mermaid progress classes.", "Add done/current/pending/blocked classDef.");
    }
  }

  const templateDir = path.join(root, "knowledge", "templates");
  for (const file of walk(templateDir).filter((item) => item.endsWith(".md"))) {
    const text = readText(file);
    if (!/Snapshot|Executive Snapshot/.test(text)) {
      addResult("warning", "templates", file, "Template missing Snapshot section.", "Add Snapshot or Executive Snapshot.");
    }
  }
}

function normalizeMarkdownLinkTarget(target) {
  let clean = target.trim();
  if (!clean || clean.startsWith("#")) return null;
  if (/^(https?:|mailto:|tel:|javascript:)/i.test(clean)) return null;
  if (clean.startsWith("<") && clean.endsWith(">")) clean = clean.slice(1, -1);
  clean = clean.split(/\s+/)[0];
  clean = clean.split("#")[0];
  if (!clean || clean.includes("[") || clean.includes("]")) return null;
  try {
    clean = decodeURI(clean);
  } catch {
    // Keep the raw target if it is not URI-encoded.
  }
  return clean;
}

function validateMarkdownLinks(files) {
  const templateDir = rel(path.join(root, "knowledge", "templates"));
  const linkPattern = /(?<!!)\[[^\]\n]+\]\(([^)\n]+)\)/g;
  for (const file of files.filter((item) => item.endsWith(".md"))) {
    const fileRel = rel(file);
    if (fileRel.startsWith(`${templateDir}/`)) continue;

    const text = readText(file);
    for (const match of text.matchAll(linkPattern)) {
      const target = normalizeMarkdownLinkTarget(match[1]);
      if (!target) continue;

      const resolved = path.resolve(path.dirname(file), target);
      if (!resolved.startsWith(root)) {
        addResult("warning", "links", file, `Markdown link points outside repository: ${target}`, "Keep local documentation links inside the repository.");
        continue;
      }
      if (!fs.existsSync(resolved)) {
        addResult("error", "links", file, `Broken Markdown link: ${target}`, "Fix the path or remove the link.");
      }
    }
  }
}

function severityCounts() {
  return {
    error: results.filter((item) => item.severity === "error").length,
    warning: results.filter((item) => item.severity === "warning").length,
    note: results.filter((item) => item.severity === "note").length,
  };
}

function verdict() {
  const counts = severityCounts();
  if (counts.error > 0) return "🔴 blocked";
  if (counts.warning > 0) return "🟡 ready_with_notes";
  return "✅ ready";
}

function markdownEscape(value) {
  return String(value ?? "").replaceAll("|", "\\|").replace(/\r?\n/g, " ");
}

function generateReport() {
  const counts = severityCounts();
  const now = new Date().toISOString().slice(0, 10);
  const rows = results.length
    ? results
        .map(
          (item) =>
            `| ${item.severity === "error" ? "🔴" : item.severity === "warning" ? "🟡" : "🔵"} ${item.severity} | ${markdownEscape(item.check)} | ${markdownEscape(item.file)} | ${markdownEscape(item.message)} | ${markdownEscape(item.fix)} |`
        )
        .join("\n")
    : "| ✅ ready | framework | repository | No findings. | None |";

  return `# Framework Validation Report

## 🧭 Executive Snapshot

| Field | Value |
| --- | --- |
| Date | ${now} |
| Validator | \`engineering/validators/framework-validator.mjs\` |
| Verdict | ${verdict()} |
| Errors | ${counts.error} |
| Warnings | ${counts.warning} |
| Notes | ${counts.note} |

## 🗺️ Validation Flow

\`\`\`mermaid
flowchart LR
  A["Scan repository"] --> B["Validate context.md"]
  B --> C["Validate use-case bundles"]
  C --> D["Validate execution graphs"]
  D --> E["Validate decisions index"]
  E --> F["Validate links and Mermaid"]
  F --> G["Report verdict"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A,B,C,D,E,F done;
  class G current;
\`\`\`

## 🚦 Check Summary

| Check | Status |
| --- | --- |
| Context metadata | ${results.some((item) => item.check === "context" && item.severity === "error") ? "🔴 has errors" : "✅ no errors"} |
| Use-case bundles | ${results.some((item) => item.check === "use-case-bundle" && item.severity === "error") ? "🔴 has errors" : "✅ no errors"} |
| Execution graph JSON and dependencies | ${results.some((item) => item.check === "execution-graph" && item.severity === "error") ? "🔴 has errors" : "✅ no errors"} |
| Decisions index | ${results.some((item) => item.check === "decisions-index") ? "🟡 findings" : "✅ no findings"} |
| Mermaid visual standard | ${results.some((item) => item.check === "mermaid") ? "🟡 findings" : "✅ no findings"} |
| Markdown links | ${results.some((item) => item.check === "links" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "links") ? "🟡 findings" : "✅ no findings"} |
| Template snapshots | ${results.some((item) => item.check === "templates") ? "🟡 findings" : "✅ no findings"} |

## 🔎 Findings

| Severity | Check | File | Finding | Suggested Fix |
| --- | --- | --- | --- | --- |
${rows}

## 🏁 Result

| Field | Value |
| --- | --- |
| Verdict | ${verdict()} |
| Required next step | ${counts.error > 0 ? "Fix blocking errors and re-run validator." : counts.warning > 0 ? "Review warnings, fix stale metadata, and re-run validator." : "Proceed to next framework step."} |
`;
}

const allFiles = walk(root);
validateUseCaseBundles();
validateExecutionGraphs();
validateContexts();
validateProductPrefixLinks(allFiles);
validateDecisionsIndex();
validateMermaidAndTemplates(allFiles);
validateMarkdownLinks(allFiles);

const report = generateReport();
if (writeReport) {
  const reportPath = path.join(root, "audits", "framework-validation-report.md");
  fs.writeFileSync(reportPath, report, "utf8");
  console.log(`Wrote ${rel(reportPath)}`);
}

const counts = severityCounts();
console.log(`Verdict: ${verdict()} (${counts.error} errors, ${counts.warning} warnings, ${counts.note} notes)`);
if (!writeReport && results.length) {
  console.log(results.map((item) => `${item.severity.toUpperCase()} ${item.check} ${item.file}: ${item.message}`).join("\n"));
}

process.exitCode = counts.error > 0 ? 1 : 0;
