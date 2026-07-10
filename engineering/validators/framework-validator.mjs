#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";

const root = process.cwd();
const args = new Set(process.argv.slice(2));
const writeReport = args.has("--write-report");
const writeRegistry = args.has("--write-registry");

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

const requiredUseCaseFilesByTier = {
  S: ["context.md", "use-case.md", "specification.md", "tasks.md", "tests.md"],
  M: ["context.md", "use-case.md", "specification.md", "design.md", "implementation-plan.md", "execution-graph.json", "tasks.md", "tests.md"],
  L: ["context.md", "use-case.md", "specification.md", "design.md", "implementation-plan.md", "execution-graph.json", "tasks.md", "tests.md", "analytics.md", "audit.md", "qa-evidence.md", "security-review.md"],
  "N/A": requiredUseCaseFiles,
};

const allowedRigorTiers = new Set(["S", "M", "L", "N/A"]);
const tierLTriggerPatterns = [
  { name: "auth", pattern: /\bauth(?:entication|enticated)?\b|\blogin\b|\bsession\b/i },
  { name: "permissions", pattern: /\bpermission(?:s)?\b|\bauthori[sz]ation\b|\brole(?:s)?\b/i },
  { name: "payment", pattern: /\bpayment(?:s)?\b|\bpaid\b|\bbilling\b|\bticketing\b/i },
  { name: "PII", pattern: /\bPII\b|\bpersonal data\b|\bLGPD\b|\bprivacy\b|\bsensitive\b/i },
  { name: "upload/UGC", pattern: /\bupload\b|\bUGC\b|\buser-generated\b/i },
  { name: "public surface", pattern: /\bpublic surface\b|\bpublic endpoint\b|\bpublic page\b|\bpublic\b/i },
  { name: "RLS/policies", pattern: /\bRLS\b|\bpolic(?:y|ies)\b|\bmigration\b/i },
];

const allowedDeliveryLevels = new Set(["L0", "L1", "L2", "L3", "L4", "L5", "N/A"]);
const allowedPriorities = new Set(["P0", "P1", "P2", "P3", "N/A"]);
const deliveryRequiredTypes = new Set([
  "domain",
  "goal",
  "feature",
  "use-case",
  "specification",
  "implementation-plan",
  "execution-graph",
  "taskset",
  "task",
]);

const results = [];
let generatedRegistry = null;

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

function normalizeArtifactContent(text) {
  return text
    .replace(/\r\n?/g, "\n")
    .split("\n")
    .map((line) => line.replace(/[ \t]+$/g, ""))
    .join("\n");
}

function artifactContentHash(file) {
  return crypto
    .createHash("sha256")
    .update(normalizeArtifactContent(readText(file)), "utf8")
    .digest("hex");
}

function artifactCurrentHash(artifact) {
  if (!artifact?.path) return "";
  const file = path.join(root, artifact.path);
  if (!fs.existsSync(file)) return "";
  return artifactContentHash(file);
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

function parseContextRigorTier(text) {
  const field = parseMarkdownField(text, "Tier");
  const source = field || text.match(/^rigor_tier:\s*(.+?)\s*$/im)?.[1] || "";
  const value = source.trim().replace(/^["'`]+|["'`]+$/g, "").toUpperCase();
  return value === "N/A" || ["S", "M", "L"].includes(value) ? value : "";
}

function parseYamlList(text, key) {
  const lines = text.split(/\r?\n/);
  const values = [];
  const start = lines.findIndex((line) => line.trim() === `${key}:`);
  if (start === -1) return values;

  for (let index = start + 1; index < lines.length; index += 1) {
    const line = lines[index];
    if (/^\S/.test(line) && !line.trim().startsWith("- ")) break;
    const item = line.trim().match(/^-\s+(.+?)\s*$/);
    if (item) values.push(item[1].replace(/^["']|["']$/g, ""));
  }
  return values;
}

function parseYamlDelivery(text) {
  const lines = text.split(/\r?\n/);
  const delivery = {};
  const start = lines.findIndex((line) => line.trim() === "delivery:");
  if (start === -1) return delivery;

  for (let index = start + 1; index < lines.length; index += 1) {
    const line = lines[index];
    if (/^\S/.test(line)) break;
    const pair = line.trim().match(/^([A-Za-z0-9_-]+):\s*(.+?)\s*$/);
    if (pair) delivery[pair[1]] = pair[2].replace(/^["']|["']$/g, "");
  }
  delivery.depends_on = parseYamlList(text.slice(text.indexOf("delivery:")), "depends_on");
  return delivery;
}

function normalizeDeliveryLevel(value) {
  const match = String(value ?? "").trim().match(/^(L[0-5]|N\/A)\b/i);
  return match ? match[1].toUpperCase() : "";
}

function normalizePriority(value) {
  const match = String(value ?? "").trim().match(/^(P[0-3]|N\/A)\b/i);
  return match ? match[1].toUpperCase() : "";
}

function findContextFiles() {
  return walk(path.join(root, "domains")).filter((item) => path.basename(item) === "context.md");
}

function normalizeArtifactType(type) {
  return (type ?? "unknown").replaceAll("_", "-");
}

function normalizeOwnerSkill(ownerSkill) {
  return (ownerSkill ?? "")
    .replace(/^\d+-/, "")
    .replace(/\.md$/, "");
}

function parseMarkdownSectionItems(text, heading) {
  const lines = text.split(/\r?\n/);
  const start = lines.findIndex((line) => line.trim() === `## ${heading}`);
  if (start === -1) return [];
  const items = [];
  for (let index = start + 1; index < lines.length; index += 1) {
    const line = lines[index];
    if (line.startsWith("## ")) break;
    const item = line.trim().match(/^-\s+(.+?)\s*$/);
    if (item) items.push(item[1]);
  }
  return items;
}

function parseLeadingIds(items) {
  return items
    .map((item) => item.split(/\s+-\s+/)[0]?.trim())
    .filter((item) => item && item !== "N/A" && /^[A-Z]+-[A-Za-z0-9.:-]+$/.test(item));
}

function parseDecisionIds(items) {
  return [...new Set(
    items
      .flatMap((item) => [...String(item).matchAll(/\bDEC-\d+\b/g)].map((match) => match[0]))
      .filter(Boolean)
  )];
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

function inferArtifactDocuments(contextFile, meta) {
  const dir = path.dirname(contextFile);
  const type = normalizeArtifactType(meta.type);
  const documents = {
    context: rel(contextFile),
  };

  const candidatesByType = {
    domain: { canonical: "domain.md", decisions: "decisions.md" },
    goal: { canonical: "goal.md", journeys: "journeys.md" },
    feature: { canonical: "feature.md", decisions: "decisions.md" },
    "use-case": {
      canonical: "use-case.md",
      specification: "specification.md",
      design: "design.md",
      implementationPlan: "implementation-plan.md",
      executionGraph: "execution-graph.json",
      tasks: "tasks.md",
      tests: "tests.md",
      qaEvidence: "qa-evidence.md",
      codeReview: "code-review.md",
      securityReview: "security-review.md",
      analytics: "analytics.md",
      audit: "audit.md",
      readme: "README.md",
    },
  };

  const candidates = candidatesByType[type] ?? {};
  for (const [name, fileName] of Object.entries(candidates)) {
    const candidate = path.join(dir, fileName);
    if (fs.existsSync(candidate)) documents[name] = rel(candidate);
  }
  return documents;
}

function parseMarkdownField(text, field) {
  const escaped = field.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const table = text.match(new RegExp(`\\|\\s*${escaped}\\s*\\|\\s*(.+?)\\s*\\|`, "i"));
  if (table) {
    const value = table[1].trim().replace(/^`|`$/g, "");
    if (field === "ID" && /^(Scenario|Task|Artifact|Decision|Metric|Persona|Opportunity)$/i.test(value)) return "";
    return value;
  }
  const bullet = text.match(new RegExp(`^-\\s+${escaped}:\\s*(.+?)\\s*$`, "im"));
  if (bullet) return bullet[1].trim();
  return "";
}

function parseDelimitedValueList(value) {
  const clean = String(value ?? "").trim();
  if (!clean || /^none$/i.test(clean) || /^N\/A$/i.test(clean)) return [];
  return clean
    .replace(/^\[|\]$/g, "")
    .split(/\s*,\s*|\s*;\s*/)
    .map((item) => item.trim().replace(/^`|`$/g, ""))
    .filter(Boolean);
}

function normalizeWriteScopePath(value) {
  return String(value ?? "")
    .trim()
    .replace(/\\/g, "/")
    .replace(/^\.\//, "")
    .replace(/\/+$/g, "");
}

function isPlaceholderScope(value) {
  const clean = normalizeWriteScopePath(value);
  return !clean || /^\[.+\]$/.test(clean) || /path\/or\/module|TBD|placeholder/i.test(clean);
}

function writeScopesOverlap(left, right) {
  const a = normalizeWriteScopePath(left);
  const b = normalizeWriteScopePath(right);
  if (!a || !b) return false;
  return a === b || a.startsWith(`${b}/`) || b.startsWith(`${a}/`);
}

function graphDependencyMap(nodes) {
  return new Map(nodes.map((node) => [node.id, Array.isArray(node.dependsOn) ? node.dependsOn : []]));
}

function hasDependencyPath(fromId, toId, dependencies, visited = new Set()) {
  if (fromId === toId) return true;
  if (visited.has(fromId)) return false;
  visited.add(fromId);
  for (const dependency of dependencies.get(fromId) ?? []) {
    if (dependency === toId || hasDependencyPath(dependency, toId, dependencies, visited)) return true;
  }
  return false;
}

function nodesAreParallel(left, right, dependencies) {
  return !hasDependencyPath(left.id, right.id, dependencies) && !hasDependencyPath(right.id, left.id, dependencies);
}

function normalizedStringArray(value) {
  return Array.isArray(value)
    ? value.map((item) => String(item ?? "").trim()).filter(Boolean)
    : [];
}

function knownSkillNames() {
  const skillsDir = path.join(root, ".codex", "skills");
  if (!fs.existsSync(skillsDir)) return new Set();
  return new Set(
    fs.readdirSync(skillsDir, { withFileTypes: true })
      .filter((entry) => entry.isDirectory() && fs.existsSync(path.join(skillsDir, entry.name, "SKILL.md")))
      .map((entry) => entry.name)
  );
}

function normalizeSkillReference(value) {
  return String(value ?? "")
    .trim()
    .replace(/^`|`$/g, "")
    .replace(/\s+AI$/i, "")
    .toLowerCase()
    .replace(/\s+/g, "-");
}

function validateSkillReference(file, field, value, skills) {
  const skill = normalizeSkillReference(value);
  if (!skill || isPlaceholderValue(skill) || /^n\/a$/i.test(skill) || skill === "none") return;
  if (!skills.has(skill)) {
    addResult(
      "error",
      "skill-reference",
      file,
      `${field} references missing skill ${value}.`,
      `Create .codex/skills/${skill}/SKILL.md or update ${field} to an existing skill.`
    );
  }
}

function validateWriteScopeSafety(file, graph) {
  const nodes = Array.isArray(graph.nodes) ? graph.nodes : [];
  const dependencies = graphDependencyMap(nodes);

  for (const node of nodes) {
    const writeScope = normalizedStringArray(node.writeScope);
    if (!Array.isArray(node.writeScope) || writeScope.length === 0 || writeScope.some(isPlaceholderScope)) {
      addResult(
        "warning",
        "write-scope",
        file,
        `Node ${node.id ?? "(missing id)"} is missing a concrete writeScope.`,
        "Declare the paths or modules this task may create or change. FDR-003 keeps this as warning during Phase A."
      );
    }
    if ("sharedResources" in node && !Array.isArray(node.sharedResources)) {
      addResult("warning", "write-scope", file, `Node ${node.id ?? "(missing id)"} sharedResources must be an array.`, "Set sharedResources to an array of stable resource names or remove it.");
    }
  }

  for (let leftIndex = 0; leftIndex < nodes.length; leftIndex += 1) {
    for (let rightIndex = leftIndex + 1; rightIndex < nodes.length; rightIndex += 1) {
      const left = nodes[leftIndex];
      const right = nodes[rightIndex];
      if (!left.id || !right.id || !nodesAreParallel(left, right, dependencies)) continue;

      const leftScopes = normalizedStringArray(left.writeScope).filter((item) => !isPlaceholderScope(item));
      const rightScopes = normalizedStringArray(right.writeScope).filter((item) => !isPlaceholderScope(item));
      for (const leftScope of leftScopes) {
        for (const rightScope of rightScopes) {
          if (writeScopesOverlap(leftScope, rightScope)) {
            addResult(
              "warning",
              "write-scope",
              file,
              `Parallel nodes ${left.id} and ${right.id} have overlapping writeScope: ${leftScope} <-> ${rightScope}.`,
              "Add a dependency, split the scope, or merge the tasks. FDR-003 keeps this as warning during Phase A."
            );
          }
        }
      }

      const leftResources = normalizedStringArray(left.sharedResources);
      const rightResources = normalizedStringArray(right.sharedResources);
      for (const resource of leftResources) {
        if (rightResources.includes(resource)) {
          addResult(
            "warning",
            "write-scope",
            file,
            `Parallel nodes ${left.id} and ${right.id} share resource ${resource}.`,
            "Serialize tasks with a dependency or assign the shared resource to one node. FDR-003 keeps this as warning during Phase A."
          );
        }
      }
    }
  }
}

function isPlaceholderValue(value) {
  const clean = String(value ?? "").trim();
  if (!clean) return true;
  if (/^`?N\/A/i.test(clean)) return true;
  if (/^`?pending`?$/i.test(clean)) return true;
  if (/until implementation|until validation|to be produced|link or notes|path or URL/i.test(clean)) return true;
  if (/^\[.+\]$/.test(clean)) return true;
  return false;
}

function statusRequiresImplementationEvidence(status) {
  return ["implemented", "validated", "released"].includes(status);
}

function statusRequiresConcreteValidationEvidence(status) {
  return ["validated", "released"].includes(status);
}

function looksLikeCommitReference(value) {
  const clean = String(value ?? "").trim();
  return /\b[a-f0-9]{7,40}\b/i.test(clean) || /^https?:\/\/\S+\/commit\/[a-f0-9]{7,40}\b/i.test(clean);
}

function looksLikePrReference(value) {
  const clean = String(value ?? "").trim();
  return /^#?\d+$/.test(clean) || /^PR-\d+$/i.test(clean) || /^https?:\/\/\S+\/pull\/\d+\b/i.test(clean);
}

function firstKnownId(text, prefixes) {
  for (const prefix of prefixes) {
    const match = text.match(new RegExp(`\\b${prefix}-[A-Z0-9.]+\\b`));
    if (match) return match[0];
  }
  return "";
}

function artifactFromDocument(parentArtifact, documentKey, documentPath) {
  const fullPath = path.join(root, documentPath);
  if (!fs.existsSync(fullPath)) return null;

  if (documentKey === "executionGraph") {
    const parsed = parseJsonFile(fullPath);
    const graph = parsed.ok ? parsed.value : {};
    return {
      id: graph.id || `${parentArtifact.id}:execution-graph`,
      type: "execution-graph",
      name: graph.id || "Execution Graph",
      status: graph.status || "unknown",
      ownerSkill: "execution-graph",
      path: documentPath,
      parentIds: [parentArtifact.id],
      childIds: Array.isArray(graph.nodes) ? graph.nodes.map((node) => node.id).filter(Boolean) : [],
      dependsOn: [],
      decisions: parseDecisionIds([
        ...(graph.delivery?.depends_on ?? []),
        ...(graph.delivery?.dependsOn ?? []),
      ]),
      delivery: {
        ...(graph.delivery ?? {}),
        level: normalizeDeliveryLevel(graph.delivery?.level),
        priority: normalizePriority(graph.delivery?.priority),
        depends_on: graph.delivery?.depends_on ?? graph.delivery?.dependsOn ?? [],
      },
      documents: {
        canonical: documentPath,
      },
    };
  }

  const text = readText(fullPath);
  const config = {
    specification: { type: "specification", prefixes: ["SPEC"], ownerSkill: "specification" },
    design: { type: "design", prefixes: ["DES", "DESIGN"], ownerSkill: "ux-ui" },
    implementationPlan: { type: "implementation-plan", prefixes: ["PLAN"], ownerSkill: "implementation-planner" },
    tasks: { type: "taskset", prefixes: ["TASKSET"], ownerSkill: "task-generator" },
    tests: { type: "tests", prefixes: ["TEST"], ownerSkill: "qa" },
    qaEvidence: { type: "qa-evidence", prefixes: ["QA"], ownerSkill: "qa" },
    codeReview: { type: "code-review", prefixes: ["CR"], ownerSkill: "code-review" },
    securityReview: { type: "security-review", prefixes: ["SEC"], ownerSkill: "security-review" },
    analytics: { type: "analytics", prefixes: ["ANA"], ownerSkill: "documentation-writer" },
    audit: { type: "audit", prefixes: ["AUD"], ownerSkill: "audit-orchestrator" },
  }[documentKey];
  if (!config) return null;

  const heading = text.match(/^#\s+(.+?)\s*$/m)?.[1]?.trim() ?? config.type;
  return {
    id: parseMarkdownField(text, "ID") || firstKnownId(text, config.prefixes) || `${parentArtifact.id}:${documentKey}`,
    type: config.type,
    name: heading.replace(/^(Specification|Design|Implementation Plan|Tasks|Tests|QA Evidence|Code Review|Security Review|Analytics|Audit):\s*/i, ""),
    status: parseMarkdownField(text, "Status") || "unknown",
    ownerSkill: config.ownerSkill,
    path: documentPath,
    parentIds: [parentArtifact.id],
    childIds: [],
    dependsOn: [],
    decisions: parseDecisionIds(text.split(/\r?\n/)),
    delivery: {
      level: normalizeDeliveryLevel(parseMarkdownField(text, "Delivery Level") || parseMarkdownField(text, "Level")),
      priority: normalizePriority(parseMarkdownField(text, "Priority")),
      depends_on: parseMarkdownSectionItems(text, "Dependencies").map((item) => item.split(" - ")[0].replace(/^`|`$/g, "")),
      rationale: parseMarkdownField(text, "Rationale"),
    },
    documents: {
      canonical: documentPath,
    },
  };
}

function taskArtifactFromFile(useCaseArtifact, graphArtifact, documentPath, node = {}) {
  const fullPath = path.join(root, documentPath);
  if (!fs.existsSync(fullPath)) return null;

  const text = readText(fullPath);
  const id = parseMarkdownField(text, "ID") || firstKnownId(text, ["TK", "TASK"]) || node.id;
  const heading = text.match(/^#\s+(.+?)\s*$/m)?.[1]?.trim() ?? id;

  return {
    id,
    type: "task",
    name: heading.replace(/^Task:\s*/i, ""),
    status: parseMarkdownField(text, "Status") || "unknown",
    ownerSkill: normalizeOwnerSkill(parseMarkdownField(text, "Owner skill") || node.ownerSkill || "task-generator"),
    path: documentPath,
    parentIds: graphArtifact?.id ? [graphArtifact.id] : [useCaseArtifact.id],
    childIds: [],
    dependsOn: parseDelimitedValueList(parseMarkdownField(text, "Depends on")),
    decisions: parseDecisionIds(text.split(/\r?\n/)),
    delivery: {
      level: normalizeDeliveryLevel(parseMarkdownField(text, "Delivery Level") || parseMarkdownField(text, "Level")),
      priority: normalizePriority(parseMarkdownField(text, "Priority")),
      depends_on: parseDelimitedValueList(parseMarkdownField(text, "Depends on")),
      rationale: parseMarkdownField(text, "Rationale"),
    },
    documents: {
      canonical: documentPath,
    },
  };
}

function buildArtifactsRegistry() {
  const artifacts = [];
  const contextFiles = findContextFiles();

  for (const contextFile of contextFiles) {
    const text = readText(contextFile);
    const meta = parseContextMeta(text);
    if (!meta?.id) continue;
    const type = normalizeArtifactType(meta.type);
    const yamlParents = parseYamlList(text, "parents");
    const yamlChildren = parseYamlList(text, "children");
    const yamlDecisions = parseYamlList(text, "decisions");
    const markdownParents = parseLeadingIds(parseMarkdownSectionItems(text, "Parent Artifacts"));
    const markdownChildren = parseLeadingIds(parseMarkdownSectionItems(text, "Child Artifacts"));
    const markdownDecisions = parseDecisionIds(parseMarkdownSectionItems(text, "Decisions"));

    const artifact = {
      id: meta.id,
      type,
      name: meta.name ?? "",
      status: meta.status ?? "unknown",
      ownerSkill: normalizeOwnerSkill(meta.owner_skill),
      path: rel(contextFile),
      parentIds: [...new Set([...yamlParents, ...markdownParents])],
      childIds: [...new Set([...yamlChildren, ...markdownChildren])],
      dependsOn: parseYamlList(text, "depends_on"),
      decisions: [...new Set([...yamlDecisions, ...markdownDecisions])],
      delivery: parseYamlDelivery(text),
      rigorTier: parseContextRigorTier(text),
      documents: inferArtifactDocuments(contextFile, meta),
    };
    artifacts.push(artifact);

    if (type === "use-case") {
      const documentArtifacts = [];
      for (const key of ["specification", "design", "implementationPlan", "executionGraph", "tasks", "tests", "qaEvidence", "codeReview", "securityReview", "analytics", "audit"]) {
        const documentPath = artifact.documents[key];
        const documentArtifact = documentPath ? artifactFromDocument(artifact, key, documentPath) : null;
        if (documentArtifact) documentArtifacts.push(documentArtifact);
      }
      artifacts.push(...documentArtifacts);

      const graphPath = artifact.documents.executionGraph;
      const graphArtifact = documentArtifacts.find((item) => item.type === "execution-graph");
      const graphFile = graphPath ? path.join(root, graphPath) : null;
      const parsedGraph = graphFile && fs.existsSync(graphFile) ? parseJsonFile(graphFile) : { ok: false };
      if (parsedGraph.ok && Array.isArray(parsedGraph.value.nodes)) {
        for (const node of parsedGraph.value.nodes) {
          if (!node.path) continue;
          const taskPath = rel(path.resolve(path.dirname(graphFile), node.path));
          const taskArtifact = taskArtifactFromFile(artifact, graphArtifact, taskPath, node);
          if (taskArtifact) artifacts.push(taskArtifact);
        }
      }
    }
  }

  const decisionFile = path.join(root, ".product", "decisions.json");
  const parsedDecisions = parseJsonFile(decisionFile);
  if (parsedDecisions.ok && Array.isArray(parsedDecisions.value.decisions)) {
    for (const decision of parsedDecisions.value.decisions) {
      artifacts.push({
        id: decision.id,
        type: "decision",
        name: decision.title ?? decision.id,
        status: decision.status ?? "unknown",
        ownerSkill: "product-historian",
        path: decision.path,
        parentIds: [],
        childIds: [],
        dependsOn: [],
        decisions: [],
        delivery: {},
        documents: {
          canonical: decision.path,
        },
        affectedArtifacts: decision.affectedArtifacts ?? [],
      });
    }
  }

  artifacts.sort((a, b) => a.id.localeCompare(b.id));
  return {
    generatedAt: new Date().toISOString(),
    generator: "engineering/validators/framework-validator.mjs",
    artifacts,
  };
}

function writeArtifactsRegistry() {
  generatedRegistry = buildArtifactsRegistry();
  const file = path.join(root, ".product", "artifacts.json");
  fs.writeFileSync(file, `${JSON.stringify(generatedRegistry, null, 2)}\n`, "utf8");
  console.log(`Wrote ${rel(file)}`);
}

function currentArtifactsRegistry() {
  return generatedRegistry ?? parseJsonFile(path.join(root, ".product", "artifacts.json")).value;
}

function currentArtifacts() {
  const registry = currentArtifactsRegistry();
  return Array.isArray(registry?.artifacts) ? registry.artifacts : [];
}

function statusCanFeedDownstream(status) {
  return ["approved", "implemented", "validated", "released"].includes(status);
}

function statusRequiresApprovedParent(status) {
  return status === "proposed" || statusCanFeedDownstream(status);
}

function statusRequiresApprovalRecord(status) {
  return statusCanFeedDownstream(status);
}

function isDraftLike(status) {
  return status === "draft";
}

function validateUseCaseBundles() {
  for (const dir of findUseCaseDirs()) {
    const contextFile = path.join(dir, "context.md");
    const contextText = fs.existsSync(contextFile) ? readText(contextFile) : "";
    const tier = parseContextRigorTier(contextText) || "L";
    const requiredFiles = requiredUseCaseFilesByTier[tier] ?? requiredUseCaseFilesByTier.L;
    for (const fileName of requiredFiles) {
      const file = path.join(dir, fileName);
      if (!fs.existsSync(file)) {
        addResult(
          "error",
          "use-case-bundle",
          dir,
          `Missing required Tier ${tier} use-case artifact: ${fileName}`,
          `Create ${fileName}, lower the tier with approval, or mark the use case as N/A if it is only structural.`
        );
      }
    }
  }
}

function validateRigorTiers() {
  for (const dir of findUseCaseDirs()) {
    const contextFile = path.join(dir, "context.md");
    if (!fs.existsSync(contextFile)) continue;

    const contextText = readText(contextFile);
    const meta = parseContextMeta(contextText);
    const tier = parseContextRigorTier(contextText);
    const fileText = walk(dir)
      .filter((file) => /\.(md|json)$/.test(file))
      .map((file) => readText(file))
      .join("\n");

    if (!tier) {
      addResult("error", "rigor-tier", contextFile, `${meta?.id ?? rel(dir)} is missing rigor_tier.`, "Add rigor_tier: S, M, L, or N/A to context.md.");
      continue;
    }

    if (!allowedRigorTiers.has(tier)) {
      addResult("error", "rigor-tier", contextFile, `${meta?.id ?? rel(dir)} has invalid rigor_tier: ${tier}.`, "Use S, M, L, or N/A.");
      continue;
    }

    if (tier === "N/A") continue;

    const triggers = tierLTriggerPatterns
      .filter((trigger) => trigger.pattern.test(fileText))
      .map((trigger) => trigger.name);
    if (triggers.length > 0 && tier !== "L") {
      addResult(
        "error",
        "rigor-tier",
        contextFile,
        `${meta?.id ?? rel(dir)} is Tier ${tier}, but Tier L trigger(s) were detected: ${[...new Set(triggers)].join(", ")}.`,
        "Raise rigor_tier to L or remove/clarify the trigger with human approval."
      );
    }

    if (tier === "S") {
      for (const optionalFile of ["design.md", "analytics.md", "audit.md"]) {
        const file = path.join(dir, optionalFile);
        if (fs.existsSync(file) && !/Not applicable/i.test(readText(file))) {
          addResult("error", "rigor-tier", file, `Tier S optional artifact ${optionalFile} exists but is not marked Not applicable.`, "Mark it Not applicable or raise the use case tier.");
        }
      }
    }
  }
}

function validateDeliveryMetadata() {
  for (const artifact of currentArtifacts().filter((item) => deliveryRequiredTypes.has(item.type))) {
    const delivery = artifact.delivery ?? {};
    const level = normalizeDeliveryLevel(delivery.level);
    const priority = normalizePriority(delivery.priority);
    const file = artifact.path ? path.join(root, artifact.path) : null;

    if (!level) {
      addResult("warning", "delivery", file, `${artifact.id} is missing delivery.level.`, "Add delivery.level using L0-L5 or N/A for placeholders.");
    } else if (!allowedDeliveryLevels.has(level)) {
      addResult("warning", "delivery", file, `${artifact.id} has invalid delivery.level: ${delivery.level}.`, "Use L0, L1, L2, L3, L4, L5, or N/A.");
    }

    if (!priority) {
      addResult("warning", "delivery", file, `${artifact.id} is missing delivery.priority.`, "Add delivery.priority using P0-P3 or N/A for placeholders.");
    } else if (!allowedPriorities.has(priority)) {
      addResult("warning", "delivery", file, `${artifact.id} has invalid delivery.priority: ${delivery.priority}.`, "Use P0, P1, P2, P3, or N/A.");
    }

    const dependsOn = delivery.depends_on ?? delivery.dependsOn ?? [];
    if (!Array.isArray(dependsOn)) {
      addResult("warning", "delivery", file, `${artifact.id} delivery dependencies must be a list.`, "Use delivery.depends_on as an array/list.");
    }

    if (!String(delivery.rationale ?? "").trim()) {
      addResult("warning", "delivery", file, `${artifact.id} is missing delivery.rationale.`, "Explain why this level and priority were assigned.");
    }
  }
}

function decisionIndexById() {
  const parsed = parseJsonFile(path.join(root, ".product", "decisions.json"));
  const decisions = Array.isArray(parsed.value?.decisions) ? parsed.value.decisions : [];
  return new Map(decisions.map((decision) => [decision.id, decision]));
}

function artifactDecisionReferences(artifact) {
  const deliveryDependsOn = artifact.delivery?.depends_on ?? artifact.delivery?.dependsOn ?? [];
  return [...new Set([
    ...parseDecisionIds(artifact.decisions ?? []),
    ...parseDecisionIds(artifact.dependsOn ?? []),
    ...parseDecisionIds(Array.isArray(deliveryDependsOn) ? deliveryDependsOn : []),
  ])];
}

function validateDecisionReferences() {
  const decisionsById = decisionIndexById();

  for (const artifact of currentArtifacts().filter((item) => item.type !== "decision")) {
    const file = artifact.path ? path.join(root, artifact.path) : null;
    const deliveryDependsOn = artifact.delivery?.depends_on ?? artifact.delivery?.dependsOn ?? [];
    const deliveryDecisionRefs = parseDecisionIds(Array.isArray(deliveryDependsOn) ? deliveryDependsOn : []);

    for (const decisionId of artifactDecisionReferences(artifact)) {
      const decision = decisionsById.get(decisionId);
      if (!decision) {
        addResult("error", "decision-references", file, `${artifact.id} references ${decisionId}, but it is missing from .product/decisions.json.`, "Add the decision to .product/decisions.json or remove the reference.");
        continue;
      }

      if (decision.path && !fs.existsSync(path.join(root, decision.path))) {
        addResult("error", "decision-references", file, `${artifact.id} references ${decisionId}, but its decision file is missing.`, "Fix the decision path or create the decision record.");
      }

      if (deliveryDecisionRefs.includes(decisionId) && decision.status !== "approved") {
        addResult("warning", "decision-references", file, `${artifact.id} depends on ${decisionId}, but decision status is ${decision.status}.`, "Approve the decision or remove it from delivery dependencies.");
      }
    }
  }
}

function addGateResult(file, child, parent, rule, fix) {
  addResult(
    "error",
    "approval-gates",
    file,
    `${child?.id ?? "Downstream artifact"} is ${child?.status ?? "missing"}, but ${parent?.id ?? "required parent"} is ${parent?.status ?? "missing"}: ${rule}`,
    fix
  );
}

function artifactByType(artifacts, parentId, type) {
  return artifacts.find((artifact) => artifact.type === type && artifact.parentIds?.includes(parentId));
}

function artifactById() {
  return new Map(currentArtifacts().map((artifact) => [artifact.id, artifact]));
}

function isPlaceholderArtifact(artifact) {
  return artifact?.id?.includes("EXAMPLE") || artifact?.delivery?.level === "N/A";
}

function validateSourceField(artifact, field, expectedId, rule) {
  if (!artifact?.path || !expectedId) return;
  const file = path.join(root, artifact.path);
  if (!fs.existsSync(file) || !artifact.path.endsWith(".md")) return;
  const actual = parseMarkdownField(readText(file), field);
  if (!actual && isPlaceholderArtifact(artifact)) return;
  if (!actual) {
    addResult("warning", "traceability", file, `${artifact.id} is missing ${field}.`, `Set ${field} to ${expectedId}.`);
  } else if (actual !== expectedId) {
    addResult("error", "traceability", file, `${artifact.id} ${field} is ${actual}, expected ${expectedId}: ${rule}`, `Update ${field} to ${expectedId}.`);
  }
}

function validateTraceability() {
  const artifacts = currentArtifacts();
  const byId = artifactById();

  for (const artifact of artifacts) {
    const file = artifact.path ? path.join(root, artifact.path) : null;
    for (const parentId of artifact.parentIds ?? []) {
      const parent = byId.get(parentId);
      if (parent && !(parent.childIds ?? []).includes(artifact.id)) {
        addResult("warning", "traceability", file, `${artifact.id} points to parent ${parentId}, but the parent does not list it as a child.`, `Add ${artifact.id} to ${parentId} children or remove the parent link.`);
      }
    }

    for (const childId of artifact.childIds ?? []) {
      if (childId.includes("..") || /^(TK|TASK)-/.test(childId)) continue;
      const child = byId.get(childId);
      if (child && !(child.parentIds ?? []).includes(artifact.id)) {
        addResult("warning", "traceability", file, `${artifact.id} lists child ${childId}, but the child does not point back to it.`, `Add ${artifact.id} to ${childId} parents or remove the child link.`);
      }
    }
  }

  for (const useCase of artifacts.filter((artifact) => artifact.type === "use-case")) {
    const spec = artifactByType(artifacts, useCase.id, "specification");
    const design = artifactByType(artifacts, useCase.id, "design");
    const plan = artifactByType(artifacts, useCase.id, "implementation-plan");
    const graph = artifactByType(artifacts, useCase.id, "execution-graph");
    const tasks = artifactByType(artifacts, useCase.id, "taskset");

    validateSourceField(spec, "Source use case", useCase.id, "Specification must trace to its use case.");
    validateSourceField(design, "Source specification", spec?.id, "Design must trace to the use-case Specification.");
    validateSourceField(plan, "Source specification", spec?.id, "Implementation Plan must trace to the use-case Specification.");
    validateSourceField(tasks, "Source graph", graph?.id, "Tasks must trace to the use-case Execution Graph.");
    if (tasks && parseMarkdownField(readText(path.join(root, tasks.path)), "Source specification")) {
      validateSourceField(tasks, "Source specification", spec?.id, "Tasks must trace to the use-case Specification.");
    }

    if (graph?.path) {
      const graphFile = path.join(root, graph.path);
      const parsed = parseJsonFile(graphFile);
      if (parsed.ok) {
        if (parsed.value.sourceSpecification !== spec?.id) {
          addResult("error", "traceability", graphFile, `${graph.id} sourceSpecification is ${parsed.value.sourceSpecification ?? "missing"}, expected ${spec?.id}.`, `Set sourceSpecification to ${spec?.id}.`);
        }
        if (parsed.value.sourceImplementationPlan && parsed.value.sourceImplementationPlan !== plan?.id) {
          addResult("error", "traceability", graphFile, `${graph.id} sourceImplementationPlan is ${parsed.value.sourceImplementationPlan}, expected ${plan?.id}.`, `Set sourceImplementationPlan to ${plan?.id}.`);
        }
      }
    }
  }
}

function validateStatusPolicy() {
  const byId = artifactById();
  for (const artifact of currentArtifacts()) {
    if (artifact.type === "decision" || isPlaceholderArtifact(artifact)) continue;
    const file = artifact.path ? path.join(root, artifact.path) : null;
    for (const parentId of artifact.parentIds ?? []) {
      const parent = byId.get(parentId);
      if (!parent || isPlaceholderArtifact(parent)) continue;
      if (statusCanFeedDownstream(artifact.status) && !statusCanFeedDownstream(parent.status)) {
        addResult(
          "error",
          "status-policy",
          file,
          `${artifact.id} is ${artifact.status}, but parent ${parent.id} is ${parent.status}.`,
          `Approve or validate ${parent.id} before advancing ${artifact.id}.`
        );
      }
    }
  }
}

function validateApprovalGates() {
  const artifacts = currentArtifacts();
  const useCases = artifacts.filter((artifact) => artifact.type === "use-case");

  for (const useCase of useCases) {
    const spec = artifactByType(artifacts, useCase.id, "specification");
    const design = artifactByType(artifacts, useCase.id, "design");
    const plan = artifactByType(artifacts, useCase.id, "implementation-plan");
    const graph = artifactByType(artifacts, useCase.id, "execution-graph");
    const tasks = artifactByType(artifacts, useCase.id, "taskset");

    if (design && !spec && statusRequiresApprovedParent(design.status)) {
      addGateResult(path.join(root, design.path), design, spec, "design requires an existing Specification.", "Create specification.md before design.md.");
    } else if (design && !statusCanFeedDownstream(spec?.status) && statusRequiresApprovedParent(design.status)) {
      addGateResult(path.join(root, design.path), design, spec, "design should not move beyond draft before Specification approval.", "Approve the Specification or keep Design as draft.");
    }

    if (plan && !design && statusRequiresApprovedParent(plan.status)) {
      addGateResult(path.join(root, plan.path), plan, design, "implementation plan requires design.md.", "Create design.md or mark Design as Not applicable.");
    } else if (plan && !statusCanFeedDownstream(design?.status) && statusRequiresApprovedParent(plan.status)) {
      addGateResult(path.join(root, plan.path), plan, design, "implementation plan should not move beyond draft before Design approval.", "Approve Design, mark Design Not applicable, or keep Implementation Plan as draft.");
    }

    if (graph && !plan && statusRequiresApprovedParent(graph.status)) {
      addGateResult(path.join(root, graph.path), graph, plan, "execution graph requires implementation-plan.md.", "Create implementation-plan.md before execution-graph.json.");
    } else if (graph && !statusCanFeedDownstream(plan?.status) && statusRequiresApprovedParent(graph.status)) {
      addGateResult(path.join(root, graph.path), graph, plan, "execution graph should not move beyond draft before Implementation Plan approval.", "Approve Implementation Plan or keep Execution Graph as draft.");
    }

    if (tasks && !graph && statusRequiresApprovedParent(tasks.status)) {
      addGateResult(path.join(root, tasks.path), tasks, graph, "tasks require execution-graph.json.", "Create and validate execution-graph.json before tasks.md.");
    } else if (tasks && !statusCanFeedDownstream(graph?.status) && statusRequiresApprovedParent(tasks.status)) {
      addGateResult(path.join(root, tasks.path), tasks, graph, "tasks should not move beyond draft before Execution Graph approval.", "Approve Execution Graph or keep Tasks as draft.");
    }
  }
}

function approvalRecordFiles() {
  const historyDir = path.join(root, ".product", "history");
  if (!fs.existsSync(historyDir)) return [];
  return walk(historyDir).filter((file) => path.basename(file).startsWith("approval-") && file.endsWith(".json"));
}

function approvalRecordsByArtifact() {
  const records = new Map();
  for (const file of approvalRecordFiles()) {
    const parsed = parseJsonFile(file);
    if (!parsed.ok) {
      addResult("error", "approval-records", file, `Invalid approval record JSON: ${parsed.error.message}`, "Ask a human to fix or replace the approval record.");
      continue;
    }

    const record = parsed.value;
    for (const field of ["artifact_id", "path", "content_hash", "status_granted", "approved_by", "approved_at", "notes"]) {
      if (!record[field]) {
        addResult("error", "approval-records", file, `Approval record is missing ${field}.`, "Ask a human to recreate the approval record from the template.");
      }
    }

    if (record.status_granted && !statusRequiresApprovalRecord(record.status_granted)) {
      addResult("error", "approval-records", file, `Approval record grants non-approval status ${record.status_granted}.`, "Use approved, in_progress, implemented, validated, or released.");
    }

    if (record.path && !fs.existsSync(path.join(root, record.path))) {
      addResult("error", "approval-records", file, `Approval record path does not exist: ${record.path}`, "Ask a human to fix the path or supersede the record.");
    }

    const key = record.artifact_id || rel(file);
    if (!records.has(key)) records.set(key, []);
    records.get(key).push({ ...record, recordPath: rel(file) });
  }
  return records;
}

function validateApprovalRecords() {
  const recordsByArtifact = approvalRecordsByArtifact();
  for (const artifact of currentArtifacts()) {
    if (!statusRequiresApprovalRecord(artifact.status) || isPlaceholderArtifact(artifact)) continue;
    const file = artifact.path ? path.join(root, artifact.path) : null;
    if (!file || !fs.existsSync(file)) continue;
    const expectedHash = artifactContentHash(file);
    const records = recordsByArtifact.get(artifact.id) ?? [];
    const matchingRecord = records.find((record) =>
      record.path === artifact.path &&
      record.status_granted === artifact.status &&
      record.content_hash === expectedHash
    );

    if (!matchingRecord) {
      const staleRecord = records.find((record) => record.path === artifact.path && record.status_granted === artifact.status);
      addResult(
        "error",
        "approval-records",
        file,
        staleRecord
          ? `${artifact.id} is ${artifact.status}, but approval record ${staleRecord.recordPath} hash does not match current content.`
          : `${artifact.id} is ${artifact.status}, but no matching approval record exists in .product/history/.`,
        "Do not auto-fix approval records. Ask the approving human to create a matching record."
      );
    }
  }
}

function loadDerivations() {
  const file = path.join(root, ".product", "derivations.json");
  if (!fs.existsSync(file)) {
    addResult("warning", "staleness", file, "Derivations index is missing.", "Generate .product/derivations.json for staleness checks.");
    return [];
  }

  const parsed = parseJsonFile(file);
  if (!parsed.ok) {
    addResult("error", "staleness", file, `Invalid derivations JSON: ${parsed.error.message}`, "Fix JSON syntax.");
    return [];
  }

  if (!Array.isArray(parsed.value.derivations)) {
    addResult("error", "staleness", file, "derivations must be an array.", "Add a derivations array.");
    return [];
  }

  return parsed.value.derivations;
}

function validateDerivationSchema(entry, index, file) {
  for (const field of ["artifact_id", "path", "derived_from"]) {
    if (!(field in entry)) {
      addResult("error", "staleness", file, `Derivation entry ${index + 1} is missing ${field}.`, `Add ${field}.`);
    }
  }
  if (!Array.isArray(entry.derived_from)) {
    addResult("error", "staleness", file, `Derivation entry ${index + 1} derived_from must be an array.`, "Use derived_from as an array.");
    return;
  }
  for (const [sourceIndex, source] of entry.derived_from.entries()) {
    for (const field of ["artifact_id", "path", "content_hash"]) {
      if (!(field in source)) {
        addResult("error", "staleness", file, `Derivation entry ${index + 1} source ${sourceIndex + 1} is missing ${field}.`, `Add ${field}.`);
      }
    }
  }
}

function validateStaleness() {
  const file = path.join(root, ".product", "derivations.json");
  const derivations = loadDerivations();
  const artifactsById = artifactById();

  for (const [index, entry] of derivations.entries()) {
    validateDerivationSchema(entry, index, file);
    const artifact = artifactsById.get(entry.artifact_id);
    const artifactFile = entry.path ? path.join(root, entry.path) : file;

    if (!artifact) {
      addResult("warning", "staleness", file, `Derivation entry references unknown artifact ${entry.artifact_id}.`, "Regenerate derivations or fix the artifact id.");
      continue;
    }
    if (entry.path !== artifact.path) {
      addResult("warning", "staleness", artifactFile, `${entry.artifact_id} derivation path is ${entry.path}, but registry path is ${artifact.path}.`, "Update .product/derivations.json.");
    }
    if (!fs.existsSync(artifactFile)) {
      addResult("error", "staleness", file, `Derived artifact path does not exist: ${entry.path}`, "Fix the path or remove the derivation entry.");
      continue;
    }

    for (const source of Array.isArray(entry.derived_from) ? entry.derived_from : []) {
      const sourceFile = source.path ? path.join(root, source.path) : null;
      if (!sourceFile || !fs.existsSync(sourceFile)) {
        addResult("error", "staleness", artifactFile, `${entry.artifact_id} source path does not exist: ${source.path ?? "missing"}`, "Fix the source path or regenerate the artifact.");
        continue;
      }

      const currentHash = artifactContentHash(sourceFile);
      if (currentHash !== source.content_hash) {
        addResult(
          statusRequiresApprovedParent(artifact.status) ? "error" : "warning",
          "staleness",
          artifactFile,
          `${entry.artifact_id} is stale because source ${source.artifact_id} changed since derivation.`,
          "Regenerate or re-approve the derived artifact, then update the derivation baseline."
        );
      }
    }
  }
}

function statusRequiresValidationEvidence(status) {
  return ["validated", "released"].includes(status);
}

function useCaseRequiresValidationEvidence(useCase, parts) {
  return [
    useCase,
    parts.spec,
    parts.design,
    parts.plan,
    parts.graph,
    parts.tasks,
    parts.tests,
    parts.qaEvidence,
    parts.codeReview,
    parts.securityReview,
    parts.audit,
  ].some((artifact) => statusRequiresValidationEvidence(artifact?.status));
}

function addValidationGateResult(file, child, requiredArtifact, rule, fix) {
  addResult(
    "error",
    "validation-gates",
    file,
    `${child?.id ?? "Artifact"} is ${child?.status ?? "missing"}, but ${requiredArtifact ?? "required evidence"} is not approved: ${rule}`,
    fix
  );
}

function validateValidationGates() {
  const artifacts = currentArtifacts();
  const useCases = artifacts.filter((artifact) => artifact.type === "use-case");

  for (const useCase of useCases) {
    const parts = {
      spec: artifactByType(artifacts, useCase.id, "specification"),
      design: artifactByType(artifacts, useCase.id, "design"),
      plan: artifactByType(artifacts, useCase.id, "implementation-plan"),
      graph: artifactByType(artifacts, useCase.id, "execution-graph"),
      tasks: artifactByType(artifacts, useCase.id, "taskset"),
      tests: artifactByType(artifacts, useCase.id, "tests"),
      qaEvidence: artifactByType(artifacts, useCase.id, "qa-evidence"),
      codeReview: artifactByType(artifacts, useCase.id, "code-review"),
      securityReview: artifactByType(artifacts, useCase.id, "security-review"),
      audit: artifactByType(artifacts, useCase.id, "audit"),
    };

    if (!useCaseRequiresValidationEvidence(useCase, parts)) continue;

    const target = [useCase, parts.audit, parts.tasks, parts.graph, parts.plan, parts.design, parts.spec]
      .find((artifact) => statusRequiresValidationEvidence(artifact?.status)) ?? useCase;
    const targetFile = target.path ? path.join(root, target.path) : null;

    if (!statusCanFeedDownstream(parts.tests?.status)) {
      addValidationGateResult(
        targetFile,
        target,
        parts.tests ? `${parts.tests.id} (${parts.tests.status})` : "tests.md",
        "validation requires approved tests before an artifact can be validated or released.",
        "Approve tests.md or keep the target artifact below validated."
      );
    }

    if (!statusCanFeedDownstream(parts.qaEvidence?.status)) {
      addValidationGateResult(
        targetFile,
        target,
        parts.qaEvidence ? `${parts.qaEvidence.id} (${parts.qaEvidence.status})` : "qa-evidence.md",
        "validation requires approved QA evidence covering acceptance criteria, tasks, security controls, and residual risks.",
        "Create and approve qa-evidence.md or keep the target artifact below validated."
      );
    }

    if (!statusCanFeedDownstream(parts.codeReview?.status)) {
      addValidationGateResult(
        targetFile,
        target,
        parts.codeReview ? `${parts.codeReview.id} (${parts.codeReview.status})` : "code-review.md",
        "validation requires approved Code Review covering completeness, adherence, and quality.",
        "Create and approve code-review.md or keep the target artifact below validated."
      );
    }

    if (!statusCanFeedDownstream(parts.securityReview?.status)) {
      addValidationGateResult(
        targetFile,
        target,
        parts.securityReview ? `${parts.securityReview.id} (${parts.securityReview.status})` : "security-review.md",
        "validation and release require approved Security Review for executable use cases.",
        "Create and approve security-review.md or keep the target artifact below validated."
      );
    }

    if (!statusCanFeedDownstream(parts.audit?.status)) {
      addValidationGateResult(
        targetFile,
        target,
        parts.audit ? `${parts.audit.id} (${parts.audit.status})` : "audit.md",
        "release-grade validation requires an approved audit with no blocking findings.",
        "Approve audit.md or keep the target artifact below validated."
      );
    }
  }
}

function validateCodeEvidenceGates() {
  for (const task of currentArtifacts().filter((artifact) => artifact.type === "task")) {
    if (!statusRequiresImplementationEvidence(task.status)) continue;

    const file = task.path ? path.join(root, task.path) : null;
    if (!file || !fs.existsSync(file)) continue;
    const text = readText(file);

    const branch = parseMarkdownField(text, "Branch");
    const commits = parseMarkdownField(text, "Commits");
    const codePaths = parseMarkdownField(text, "Code paths");

    if (isPlaceholderValue(branch)) {
      addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but Branch is missing or placeholder.`, "Add the implementation branch to the task file.");
    }
    if (isPlaceholderValue(commits)) {
      addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but Commits is missing or placeholder.`, "Add one or more implementation commit hashes.");
    } else {
      const commitRefs = parseDelimitedValueList(commits);
      const refsToCheck = commitRefs.length > 0 ? commitRefs : [commits];
      if (!refsToCheck.some(looksLikeCommitReference)) {
        addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but Commits does not contain a traceable commit hash or commit URL.`, "Add at least one 7-40 character git hash or commit URL.");
      }
    }
    if (isPlaceholderValue(codePaths)) {
      addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but Code paths is missing or placeholder.`, "Add repository-relative code paths touched by this task.");
    }

    for (const codePath of parseDelimitedValueList(codePaths)) {
      if (isPlaceholderValue(codePath) || /^https?:\/\//i.test(codePath)) continue;
      const resolved = path.resolve(root, codePath);
      if (!resolved.startsWith(root)) {
        addResult("error", "code-evidence", file, `${task.id} code path points outside repository: ${codePath}`, "Use repository-relative code paths inside the monorepo.");
      } else if (!fs.existsSync(resolved)) {
        addResult("warning", "code-evidence", file, `${task.id} code path does not exist yet: ${codePath}`, "Confirm the path or keep the task below implemented until code exists.");
      }
    }

    if (!statusRequiresConcreteValidationEvidence(task.status)) continue;

    const pr = parseMarkdownField(text, "PR");
    const testStatus = parseMarkdownField(text, "Test status");
    const evidenceFields = [
      parseMarkdownField(text, "Gate logs"),
      parseMarkdownField(text, "CI URL"),
      parseMarkdownField(text, "Screenshots"),
      parseMarkdownField(text, "QA evidence"),
    ];

    if (isPlaceholderValue(pr)) {
      addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but PR is missing or placeholder.`, "Add the PR URL, PR id, or equivalent review surface.");
    } else if (!looksLikePrReference(pr)) {
      addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but PR is not a traceable PR reference.`, "Use a PR number, #number, PR-number, or pull request URL.");
    }
    if (!/^(passed|approved|success|succeeded)$/i.test(String(testStatus ?? "").trim())) {
      addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but Test status is not passing.`, "Set Test status to passed/approved/success only after QA evidence exists.");
    }
    if (!evidenceFields.some((value) => !isPlaceholderValue(value))) {
      addResult("error", "code-evidence", file, `${task.id} is ${task.status}, but no concrete validation evidence is linked.`, "Add gate logs, CI URL, screenshots, or QA evidence.");
    }
  }
}

function validateQaEvidenceQuality() {
  for (const qaEvidence of currentArtifacts().filter((artifact) => artifact.type === "qa-evidence")) {
    if (!statusCanFeedDownstream(qaEvidence.status)) continue;

    const file = qaEvidence.path ? path.join(root, qaEvidence.path) : null;
    if (!file || !fs.existsSync(file)) continue;
    const text = readText(file);

    const gateLogs = parseMarkdownField(text, "Gate logs");
    const testCommand = parseMarkdownField(text, "Test command");
    const environment = parseMarkdownField(text, "Environment");
    const verdictValue = parseMarkdownField(text, "Verdict");
    const limitations = parseMarkdownField(text, "Limitations");
    const hasConcreteEvidence = [gateLogs, limitations].some((value) => !isPlaceholderValue(value));

    if (isPlaceholderValue(testCommand)) {
      addResult("error", "qa-evidence", file, `${qaEvidence.id} is ${qaEvidence.status}, but Test command is missing or placeholder.`, "Record the actual gate command, method, CI job, or manual verification method.");
    }
    if (isPlaceholderValue(environment)) {
      addResult("error", "qa-evidence", file, `${qaEvidence.id} is ${qaEvidence.status}, but Environment is missing or placeholder.`, "Record local, CI, staging, production, or another concrete verification environment.");
    }
    if (!hasConcreteEvidence) {
      addResult("error", "qa-evidence", file, `${qaEvidence.id} is ${qaEvidence.status}, but no real gate output or limitation is recorded.`, "Add gate logs, captured output, CI URL, or an explicit limitation explaining why the gate could not run.");
    }
    if (!/^(passed|passed_with_notes)$/i.test(String(verdictValue).trim())) {
      addResult("error", "qa-evidence", file, `${qaEvidence.id} is ${qaEvidence.status}, but Verdict is not passing.`, "Use passed or passed_with_notes only after independent QA evidence exists.");
    }
  }
}

function validateCodeReviewQuality() {
  for (const codeReview of currentArtifacts().filter((artifact) => artifact.type === "code-review")) {
    if (!statusCanFeedDownstream(codeReview.status)) continue;

    const file = codeReview.path ? path.join(root, codeReview.path) : null;
    if (!file || !fs.existsSync(file)) continue;
    const text = readText(file);
    const verdictValue = parseMarkdownField(text, "Verdict");
    const completeness = parseMarkdownField(text, "Completeness passed");
    const adherence = parseMarkdownField(text, "Adherence passed");
    const quality = parseMarkdownField(text, "Quality passed");

    if (!/^(passed|passed_with_notes)$/i.test(String(verdictValue).trim())) {
      addResult("error", "code-review", file, `${codeReview.id} is ${codeReview.status}, but Verdict is not passing.`, "Use passed or passed_with_notes only after blockers and required fixes are resolved.");
    }
    for (const [field, value] of [["Completeness passed", completeness], ["Adherence passed", adherence], ["Quality passed", quality]]) {
      if (!/^yes$/i.test(String(value).trim())) {
        addResult("error", "code-review", file, `${codeReview.id} is ${codeReview.status}, but ${field} is not yes.`, `Set ${field} to yes only after review evidence supports it.`);
      }
    }
  }
}

function parseMarkdownTableRows(text) {
  const lines = text.split(/\r?\n/);
  const tables = [];
  for (let index = 0; index < lines.length; index += 1) {
    if (!lines[index].trim().startsWith("|")) continue;
    const header = splitMarkdownTableRow(lines[index]);
    const separator = lines[index + 1] ? splitMarkdownTableRow(lines[index + 1]) : [];
    if (!header.length || !separator.every((cell) => /^:?-{3,}:?$/.test(cell.trim()))) continue;

    const rows = [];
    let cursor = index + 2;
    while (cursor < lines.length && lines[cursor].trim().startsWith("|")) {
      rows.push(splitMarkdownTableRow(lines[cursor]));
      cursor += 1;
    }
    tables.push({ header, rows });
    index = cursor;
  }
  return tables;
}

function splitMarkdownTableRow(line) {
  return line
    .trim()
    .replace(/^\|/, "")
    .replace(/\|$/, "")
    .split("|")
    .map((cell) => cell.trim().replace(/^`|`$/g, ""));
}

function isBlockingFindingSeverity(value) {
  return /(^|\s|\/)(blocker|high)(\s|\/|$)/i.test(String(value ?? "")) || String(value ?? "").includes("🔴");
}

function validateBlockingFindingRoutes() {
  const routedArtifactTypes = new Set(["qa-evidence", "code-review", "security-review", "audit"]);
  for (const artifact of currentArtifacts().filter((item) => routedArtifactTypes.has(item.type))) {
    if (!statusCanFeedDownstream(artifact.status)) continue;
    const file = artifact.path ? path.join(root, artifact.path) : null;
    if (!file || !fs.existsSync(file)) continue;

    for (const table of parseMarkdownTableRows(readText(file))) {
      const normalizedHeader = table.header.map((cell) => cell.toLowerCase());
      const severityIndex = normalizedHeader.indexOf("severity");
      const findingIndex = normalizedHeader.indexOf("finding");
      if (severityIndex === -1 || findingIndex === -1) continue;

      const routeIndex = normalizedHeader.indexOf("route");
      const ownerIndex = normalizedHeader.indexOf("owner");
      for (const row of table.rows) {
        const finding = row[findingIndex] || "(unnamed finding)";
        const severity = String(row[severityIndex] ?? "").trim();
        const requiresRoute = isBlockingFindingSeverity(severity) || /^required_fix$/i.test(severity);
        if (!requiresRoute) continue;
        if (routeIndex === -1 || isPlaceholderValue(row[routeIndex])) {
          addResult("error", "failure-routing", file, `${artifact.id} blocking or required finding "${finding}" is missing Route.`, "Route blockers and required fixes with FDR-006: bug-fixer, code-runner, qa, or product-historian.");
        }
        if (ownerIndex === -1 || isPlaceholderValue(row[ownerIndex])) {
          addResult("error", "failure-routing", file, `${artifact.id} blocking or required finding "${finding}" is missing Owner.`, "Assign the finding to a skill or human owner.");
        }
      }
    }
  }
}

function validateExecutionGraphs() {
  const skills = knownSkillNames();
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
    validateWriteScopeSafety(file, graph);

    const graphDir = path.dirname(file);
    const tasksIndexFile = path.join(graphDir, "tasks.md");
    const tasksIndexText = fs.existsSync(tasksIndexFile) ? readText(tasksIndexFile) : "";
    if (!tasksIndexText.includes("Generated index")) {
      addResult("error", "tasks-index", tasksIndexFile, "tasks.md must be marked as a generated index.", "Regenerate tasks.md from execution-graph.json and tasks/*.md.");
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
      for (const field of ["id", "path", "dependsOn"]) {
        if (!(field in node)) {
          addResult("error", "execution-graph", file, `Node ${node.id} is missing ${field}.`, `Add ${field}.`);
        }
      }

      if (node.path) {
        if (typeof node.path !== "string" || !node.path.startsWith("tasks/") || !node.path.endsWith(".md")) {
          addResult("error", "execution-graph", file, `Node ${node.id} path must point to tasks/<task-id>.md.`, "Set path to the canonical task file.");
        }

        const taskFile = path.resolve(graphDir, node.path);
        if (!taskFile.startsWith(graphDir)) {
          addResult("error", "execution-graph", file, `Node ${node.id} path points outside the use-case directory.`, "Keep task files inside the use-case tasks/ directory.");
        } else if (!fs.existsSync(taskFile)) {
          addResult("error", "execution-graph", file, `Node ${node.id} task file does not exist: ${node.path}`, "Create the task file from knowledge/templates/task-template.md.");
        } else {
          const taskText = readText(taskFile);
          const taskId = parseMarkdownField(taskText, "ID");
          const taskStatus = parseMarkdownField(taskText, "Status");
          const taskSourceGraph = parseMarkdownField(taskText, "Source graph");
          const taskSourceNode = parseMarkdownField(taskText, "Source node");
          const taskTitle = parseMarkdownField(taskText, "Title");
          const taskType = parseMarkdownField(taskText, "Type");
          const taskOwnerSkill = parseMarkdownField(taskText, "Owner skill");
          const taskNextSkill = parseMarkdownField(taskText, "Next skill");
          const taskRequiredNextSkill = parseMarkdownField(taskText, "Required next skill");

          if (taskId !== node.id) {
            addResult("error", "task-file", taskFile, `Task file ID is ${taskId || "missing"}, expected ${node.id}.`, `Set ID to ${node.id}.`);
          }
          if (!taskStatus || !allowedStatuses.has(taskStatus)) {
            addResult("error", "task-file", taskFile, `Task ${node.id} has invalid status: ${taskStatus || "missing"}.`, "Use a framework-approved artifact status.");
          }
          if (taskSourceGraph !== graph.id) {
            addResult("error", "task-file", taskFile, `Task ${node.id} Source graph is ${taskSourceGraph || "missing"}, expected ${graph.id}.`, `Set Source graph to ${graph.id}.`);
          }
          if (taskSourceNode !== node.id) {
            addResult("error", "task-file", taskFile, `Task ${node.id} Source node is ${taskSourceNode || "missing"}, expected ${node.id}.`, `Set Source node to ${node.id}.`);
          }
          if (node.title && taskTitle && node.title !== taskTitle) {
            addResult("error", "task-file", taskFile, `Task ${node.id} title snapshot differs from graph title.`, "Regenerate graph snapshot or update the task file from the approved source.");
          }
          if (node.type && taskType && node.type !== taskType) {
            addResult("error", "task-file", taskFile, `Task ${node.id} type snapshot differs from graph type.`, "Regenerate graph snapshot or update the task file from the approved source.");
          }
          validateSkillReference(taskFile, "Owner skill", taskOwnerSkill || node.ownerSkill, skills);
          validateSkillReference(taskFile, "Next skill", taskNextSkill, skills);
          validateSkillReference(taskFile, "Required next skill", taskRequiredNextSkill, skills);
          if (tasksIndexText && !tasksIndexText.includes(node.path)) {
            addResult("error", "tasks-index", tasksIndexFile, `tasks.md does not link task file ${node.path}.`, "Regenerate tasks.md from the graph and task files.");
          }
        }
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
      } else {
        addResult("error", "execution-graph", file, `Node ${node.id} dependsOn must be an array.`, "Set dependsOn to an array of task ids.");
      }
    }
  }
}

function validateContexts() {
  for (const file of findContextFiles()) {
    const meta = parseContextMeta(readText(file));
    if (!meta) {
      addResult("error", "context", file, "context.md is missing yaml metadata block.", "Add a yaml metadata block.");
      continue;
    }
    for (const field of ["id", "type", "name", "status", "owner_skill"]) {
      if (!meta[field]) addResult("error", "context", file, `Missing context field: ${field}.`, `Add ${field}.`);
    }
    if (!meta.slug) {
      addResult("error", "identity", file, "Missing immutable slug in context metadata.", "Add slug matching the artifact folder name.");
    } else if (meta.slug !== path.basename(path.dirname(file))) {
      addResult("error", "identity", file, `Context slug ${meta.slug} does not match folder ${path.basename(path.dirname(file))}.`, "Keep slug equal to the immutable folder name or move the folder with engineering/move-artifact.mjs.");
    }
    if (meta.status && !allowedStatuses.has(meta.status)) {
      addResult("error", "context", file, `Invalid status: ${meta.status}.`, "Use a framework-approved status.");
    }
    if (!readText(file).includes("## Handoff")) {
      addResult("warning", "context", file, "Missing Handoff section.", "Add next skill and required reading.");
    }
  }
}

function validateIdentityPolicy() {
  const idsFile = path.join(root, ".product", "ids.json");
  if (!fs.existsSync(idsFile)) {
    addResult("error", "identity", idsFile, ".product/ids.json is missing.", "Add identity policy metadata.");
    return;
  }

  const parsed = parseJsonFile(idsFile);
  if (!parsed.ok) {
    addResult("error", "identity", idsFile, `Invalid ids.json: ${parsed.error.message}`, "Fix JSON syntax.");
    return;
  }

  const ids = parsed.value;
  if (ids.policy !== "slug-scoped") {
    addResult("error", "identity", idsFile, `.product/ids.json policy is ${ids.policy ?? "missing"}, expected slug-scoped.`, "Use the DEC-009 identity policy.");
  }
  if (ids.deprecated_counters !== true) {
    addResult("error", "identity", idsFile, ".product/ids.json must mark central counters as deprecated.", "Set deprecated_counters to true.");
  }
  const numericCounters = Object.entries(ids).filter(([, value]) => typeof value === "number");
  if (numericCounters.length > 0) {
    addResult("error", "identity", idsFile, `Central numeric counters remain: ${numericCounters.map(([key]) => key).join(", ")}.`, "Remove central counters and use parent-scoped IDs.");
  }

  const moveTool = path.join(root, "engineering", "move-artifact.mjs");
  if (!fs.existsSync(moveTool)) {
    addResult("error", "identity", moveTool, "Move tooling is missing.", "Create engineering/move-artifact.mjs.");
  }
}

function validateArtifactsRegistry() {
  const file = path.join(root, ".product", "artifacts.json");
  if (!fs.existsSync(file)) {
    addResult("warning", "artifacts-registry", file, "Artifacts registry is missing.", "Run validator with --write-registry.");
    return;
  }

  const parsed = parseJsonFile(file);
  if (!parsed.ok) {
    addResult("error", "artifacts-registry", file, `Invalid artifacts registry JSON: ${parsed.error.message}`, "Fix JSON syntax or regenerate registry.");
    return;
  }

  const registry = parsed.value;
  if (!Array.isArray(registry.artifacts)) {
    addResult("error", "artifacts-registry", file, "Registry artifacts must be an array.", "Regenerate registry.");
    return;
  }

  const ids = new Set();
  for (const artifact of registry.artifacts) {
    if (!artifact.id) {
      addResult("error", "artifacts-registry", file, "Registry artifact missing id.", "Regenerate registry.");
      continue;
    }
    if (ids.has(artifact.id)) {
      addResult("error", "artifacts-registry", file, `Duplicate artifact id: ${artifact.id}`, "Resolve duplicate IDs.");
    }
    ids.add(artifact.id);

    if (artifact.path && !fs.existsSync(path.join(root, artifact.path))) {
      addResult("error", "artifacts-registry", file, `Artifact ${artifact.id} path does not exist: ${artifact.path}`, "Fix path or regenerate registry.");
    }
    for (const documentPath of Object.values(artifact.documents ?? {})) {
      if (documentPath && !fs.existsSync(path.join(root, documentPath))) {
        addResult("error", "artifacts-registry", file, `Artifact ${artifact.id} document path does not exist: ${documentPath}`, "Fix document path or regenerate registry.");
      }
    }
  }

  for (const contextFile of findContextFiles()) {
    const meta = parseContextMeta(readText(contextFile));
    if (meta?.id && !ids.has(meta.id)) {
      addResult("error", "artifacts-registry", file, `Context id missing from registry: ${meta.id}`, "Regenerate registry.");
    }
  }
}

function validateProductPrefixLinks(files) {
  const staleProductPath = /(^|[^.\w-])product\/(domains|knowledge|audits|foundation|design|engineering|releases|skills|FRAMEWORK\.md)/;
  const allowedAdoptionDocs = new Set([
    "FRAMEWORK.md",
    "README.md",
    "framework/adoption.md",
    "starter/README.md",
    "starter/AGENTS.md",
    "starter/.spec-framework/README.md",
    "starter/.spec-framework/decisions/README.md",
    "audits/framework-validation-report.md",
  ]);
  for (const file of files.filter((item) => /\.(md|json)$/.test(item))) {
    const fileRel = rel(file);
    if (allowedAdoptionDocs.has(fileRel)) continue;
    const text = readText(file);
    if (staleProductPath.test(text)) {
      addResult(
        "warning",
        "paths",
        file,
        "Found product/ path prefix outside FRAMEWORK.md.",
        "Use repository-root-relative paths unless documenting product-root adoption."
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

function extractMermaidBlocks(text) {
  const blocks = [];
  const pattern = /```mermaid\s*([\s\S]*?)```/g;
  for (const match of text.matchAll(pattern)) {
    blocks.push(match[1]);
  }
  return blocks;
}

function declaredMermaidNodeIds(block) {
  const ids = new Set();
  for (const match of block.matchAll(/\b([A-Za-z][A-Za-z0-9_]*)\s*(?=\[|\{|\()/g)) {
    ids.add(match[1]);
  }
  return ids;
}

function mermaidClassAssignments(block) {
  const assignments = [];
  const pattern = /^\s*class\s+([A-Za-z0-9_,\s]+)\s+([A-Za-z][A-Za-z0-9_-]*)\s*;?\s*$/gm;
  for (const match of block.matchAll(pattern)) {
    for (const node of match[1].split(",").map((item) => item.trim()).filter(Boolean)) {
      assignments.push({ node, state: match[2] });
    }
  }
  return assignments;
}

function mermaidArtifactBindings(block) {
  const bindings = [];
  const pattern = /%%\s*artifact:\s*([A-Za-z0-9_-]+)\s+node:\s*([A-Za-z][A-Za-z0-9_]*)\s*%%/g;
  for (const match of block.matchAll(pattern)) {
    bindings.push({ artifactId: match[1], node: match[2] });
  }
  return bindings;
}

function artifactStatusToMermaidState(status) {
  if (["approved", "implemented", "validated", "released"].includes(status)) return "done";
  if (["draft", "proposed", "in_progress"].includes(status)) return "current";
  if (["deprecated", "superseded"].includes(status)) return "blocked";
  return "pending";
}

function artifactRegistryById() {
  return new Map(currentArtifacts().map((artifact) => [artifact.id, artifact]));
}

function validateMermaidProgressState(file, block, index, artifactsById) {
  const hasProgressDefinitions = /\bclassDef\s+done\b/.test(block) || /\bclassDef\s+current\b/.test(block);
  if (!hasProgressDefinitions) {
    addResult("warning", "mermaid", file, `Flowchart ${index + 1} is missing Mermaid progress classes.`, "Add done/current/pending/blocked classDef.");
    return;
  }

  const assignments = mermaidClassAssignments(block);
  const allowedStates = new Set(["done", "current", "pending", "blocked"]);
  const declaredIds = declaredMermaidNodeIds(block);
  const nodeStates = new Map(assignments.map((item) => [item.node, item.state]));

  if (!assignments.some((item) => item.state === "current")) {
    addResult("warning", "mermaid-progress", file, `Flowchart ${index + 1} has progress classes but no current node.`, "Assign one node with `class <node> current;`.");
  }

  for (const assignment of assignments) {
    if (!allowedStates.has(assignment.state)) {
      addResult("warning", "mermaid-progress", file, `Flowchart ${index + 1} uses unknown progress state ${assignment.state}.`, "Use only done/current/pending/blocked.");
    }
    if (declaredIds.size > 0 && !declaredIds.has(assignment.node)) {
      addResult("warning", "mermaid-progress", file, `Flowchart ${index + 1} assigns ${assignment.state} to undeclared node ${assignment.node}.`, "Declare the node in the flowchart or fix the class assignment.");
    }
  }

  for (const binding of mermaidArtifactBindings(block)) {
    if (binding.artifactId.includes("XXX")) continue;

    const artifact = artifactsById.get(binding.artifactId);
    if (!artifact) {
      addResult("warning", "mermaid-semantic", file, `Flowchart ${index + 1} references artifact ${binding.artifactId}, but it is not in .product/artifacts.json.`, "Regenerate the registry or fix the artifact binding.");
      continue;
    }

    const visualState = nodeStates.get(binding.node);
    if (!visualState) {
      addResult("warning", "mermaid-semantic", file, `Flowchart ${index + 1} binds ${binding.artifactId} to node ${binding.node}, but that node has no progress class.`, "Assign done/current/pending/blocked to the bound node.");
      continue;
    }

    const expectedState = artifactStatusToMermaidState(artifact.status);
    if (visualState !== expectedState) {
      addResult("warning", "mermaid-semantic", file, `Flowchart ${index + 1} shows ${binding.artifactId} as ${visualState}, but registry status ${artifact.status} maps to ${expectedState}.`, "Update the Mermaid class or the artifact status.");
    }
  }
}

function validateMermaidAndTemplates(files) {
  const artifactsById = artifactRegistryById();
  for (const file of files.filter((item) => item.endsWith(".md"))) {
    const text = readText(file);
    const flowcharts = extractMermaidBlocks(text).filter((block) => /\bflowchart\b/.test(block));
    for (const [index, block] of flowcharts.entries()) {
      validateMermaidProgressState(file, block, index, artifactsById);
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
  if (/^\[.+\]$/.test(clean) || /^(TBD|N\/A)$/i.test(clean)) return null;
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

function stripFencedCodeBlocks(text) {
  return text.replace(/```[\s\S]*?```/g, (block) => "\n".repeat(block.split(/\r?\n/).length - 1));
}

function isTemplateGeneratedArtifactLink(file, target) {
  const fileRel = rel(file);
  if (!fileRel.startsWith("knowledge/templates/")) return false;
  const normalized = normalizeWriteScopePath(target);
  const templateGeneratedTargets = new Set([
    "context.md",
    "use-case.md",
    "specification.md",
    "design.md",
    "implementation-plan.md",
    "execution-graph.json",
    "tasks.md",
    "tests.md",
    "qa-evidence.md",
    "code-review.md",
    "security-review.md",
    "analytics.md",
    "audit.md",
    "../context.md",
    "../specification.md",
    "../implementation-plan.md",
    "../execution-graph.json",
    "../tasks.md",
    "../tests.md",
  ]);
  return templateGeneratedTargets.has(normalized) || /^tasks\/TK-XXX-[0-9]+\.md$/i.test(normalized);
}

function validateMarkdownLinks(files) {
  const linkPattern = /(?<!!)\[[^\]\n]+\]\(([^)\n]+)\)/g;
  for (const file of files.filter((item) => item.endsWith(".md"))) {
    const text = stripFencedCodeBlocks(readText(file));
    for (const match of text.matchAll(linkPattern)) {
      const target = normalizeMarkdownLinkTarget(match[1]);
      if (!target) continue;
      if (isTemplateGeneratedArtifactLink(file, target)) continue;

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

function readinessIcon(status) {
  if (statusCanFeedDownstream(status)) return "✅";
  if (status === "proposed" || status === "in_progress") return "🟡";
  if (status === "draft" || status === "unknown") return "➖";
  return "🔴";
}

function artifactStatusCell(artifact) {
  return artifact ? `${readinessIcon(artifact.status)} ${artifact.status}` : "🔴 missing";
}

function nextReadinessOwner(parts) {
  if (!statusCanFeedDownstream(parts.spec?.status)) return "Specification AI";
  if (!statusCanFeedDownstream(parts.design?.status)) return "UX/UI AI";
  if (!statusCanFeedDownstream(parts.plan?.status)) return "Implementation Planner AI";
  if (!statusCanFeedDownstream(parts.graph?.status)) return "Execution Graph AI";
  if (!statusCanFeedDownstream(parts.tasks?.status)) return "Task AI";
  if (!statusCanFeedDownstream(parts.tests?.status)) return "QA AI";
  if (!statusCanFeedDownstream(parts.qaEvidence?.status)) return "QA AI";
  if (!statusCanFeedDownstream(parts.codeReview?.status)) return "Code Review AI";
  if (!statusCanFeedDownstream(parts.securityReview?.status)) return "Security Review AI";
  return "Audit Orchestrator";
}

function useCaseReadiness(useCase) {
  const artifacts = currentArtifacts();
  const parts = {
    spec: artifactByType(artifacts, useCase.id, "specification"),
    design: artifactByType(artifacts, useCase.id, "design"),
    plan: artifactByType(artifacts, useCase.id, "implementation-plan"),
    graph: artifactByType(artifacts, useCase.id, "execution-graph"),
    tasks: artifactByType(artifacts, useCase.id, "taskset"),
    tests: artifactByType(artifacts, useCase.id, "tests"),
    qaEvidence: artifactByType(artifacts, useCase.id, "qa-evidence"),
    codeReview: artifactByType(artifacts, useCase.id, "code-review"),
    securityReview: artifactByType(artifacts, useCase.id, "security-review"),
    analytics: artifactByType(artifacts, useCase.id, "analytics"),
    audit: artifactByType(artifacts, useCase.id, "audit"),
  };
  const required = [parts.spec, parts.design, parts.plan, parts.graph, parts.tasks];
  const validationRequired = [parts.tests, parts.qaEvidence, parts.codeReview, parts.securityReview, parts.audit];
  const missing = required.filter((item) => !item).length;
  const approved = required.filter((item) => statusCanFeedDownstream(item?.status)).length;
  const score = missing > 0 ? Math.round(((required.length - missing) / required.length) * 100) : Math.round((approved / required.length) * 100);
  const canGenerateTasks = statusCanFeedDownstream(parts.spec?.status) && statusCanFeedDownstream(parts.design?.status) && statusCanFeedDownstream(parts.plan?.status) && parts.graph;
  const validationReady = canGenerateTasks && validationRequired.every((item) => statusCanFeedDownstream(item?.status));
  const verdictLabel = missing > 0 ? "\u{1F534} not_ready" : canGenerateTasks ? "\u{2705} ready" : "\u{1F7E1} in_progress";

  return {
    useCase,
    parts,
    score,
    canGenerateTasks,
    validationReady,
    verdictLabel,
    nextOwner: nextReadinessOwner(parts),
  };
}

function generateReadinessReport() {
  const now = new Date().toISOString().slice(0, 10);
  const useCases = currentArtifacts()
    .filter((artifact) => artifact.type === "use-case" && artifact.delivery?.level !== "N/A")
    .map(useCaseReadiness)
    .sort((a, b) => a.useCase.id.localeCompare(b.useCase.id));
  const readyCount = useCases.filter((item) => item.canGenerateTasks).length;
  const validationReadyCount = useCases.filter((item) => item.validationReady).length;
  const rows = useCases.length
    ? useCases
        .map((item) => {
          const link = item.useCase.path ? `[${item.useCase.id}](../../${item.useCase.path})` : item.useCase.id;
          return `| ${link} | ${markdownEscape(item.useCase.name)} | ${item.verdictLabel} | ${item.score}% | ${artifactStatusCell(item.parts.spec)} | ${artifactStatusCell(item.parts.design)} | ${artifactStatusCell(item.parts.plan)} | ${artifactStatusCell(item.parts.graph)} | ${artifactStatusCell(item.parts.tasks)} | ${artifactStatusCell(item.parts.tests)} | ${artifactStatusCell(item.parts.qaEvidence)} | ${artifactStatusCell(item.parts.codeReview)} | ${artifactStatusCell(item.parts.securityReview)} | ${item.canGenerateTasks ? "yes" : "no"} | ${item.validationReady ? "yes" : "no"} | ${item.nextOwner} |`;
        })
        .join("\n")
    : "| \u{2796} none | No use cases found. | \u{2796} | 0% | \u{2796} | \u{2796} | \u{2796} | \u{2796} | \u{2796} | \u{2796} | \u{2796} | \u{2796} | \u{2796} | no | no | Product Orchestrator |";

  return `# Framework Readiness Matrix

## \u{1F9ED} Executive Snapshot

| Field | Value |
| --- | --- |
| Date | ${now} |
| Auditor | \`engineering/validators/framework-validator.mjs\` |
| Scope | Use cases with real delivery level |
| Use cases checked | ${useCases.length} |
| Ready for task generation | ${readyCount} |
| Ready for validation/release | ${validationReadyCount} |
| Overall verdict | ${readyCount === useCases.length && useCases.length > 0 ? "\u{2705} ready" : "\u{1F7E1} in_progress"} |

## \u{1F5FA}\u{FE0F} Readiness Flow

\`\`\`mermaid
flowchart LR
  U["Use Case"] --> S["Specification"]
  S --> D["Design"]
  D --> P["Implementation Plan"]
  P --> G["Execution Graph"]
  G --> T["Tasks"]
  T --> Q["QA Evidence"]
  Q --> CR["Code Review"]
  CR --> SR["Security Review"]
  SR --> R["Ready For Validation"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class U done;
  class S current;
  class D,P,G,T,Q,CR,SR,R pending;
\`\`\`

## \u{1F6A6} Use Case Matrix

| Use Case | Name | Verdict | Score | Spec | Design | Plan | Graph | Tasks | Tests | QA Evidence | Code Review | Security Review | Can Generate Tasks | Validation Ready | Next Owner |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
${rows}

## \u{1F3C1} Result

| Field | Value |
| --- | --- |
| Current bottleneck | Specification approval and downstream approval gates |
| Recommended next skill | Specification AI |
| Required next step | Review and approve or revise draft/proposed specifications before design, plan, graph, and tasks advance. |
`;
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
| Identity policy | ${results.some((item) => item.check === "identity" && item.severity === "error") ? "\u{1F534} has errors" : results.some((item) => item.check === "identity") ? "\u{1F7E1} findings" : "\u{2705} no findings"} |
| Use-case bundles | ${results.some((item) => item.check === "use-case-bundle" && item.severity === "error") ? "🔴 has errors" : "✅ no errors"} |
| Rigor tiers | ${results.some((item) => item.check === "rigor-tier" && item.severity === "error") ? "\u{1F534} has errors" : results.some((item) => item.check === "rigor-tier") ? "\u{1F7E1} findings" : "\u{2705} no findings"} |
| Approval gates | ${results.some((item) => item.check === "approval-gates" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "approval-gates") ? "🟡 findings" : "✅ no findings"} |
| Approval records | ${results.some((item) => item.check === "approval-records" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "approval-records") ? "🟡 findings" : "✅ no findings"} |
| Derived staleness | ${results.some((item) => item.check === "staleness" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "staleness") ? "🟡 findings" : "✅ no findings"} |
| Validation gates | ${results.some((item) => item.check === "validation-gates" && item.severity === "error") ? "\u{1F534} has errors" : results.some((item) => item.check === "validation-gates") ? "\u{1F7E1} findings" : "\u{2705} no findings"} |
| Code evidence gates | ${results.some((item) => item.check === "code-evidence" && item.severity === "error") ? "\u{1F534} has errors" : results.some((item) => item.check === "code-evidence") ? "\u{1F7E1} findings" : "\u{2705} no findings"} |
| QA evidence quality | ${results.some((item) => item.check === "qa-evidence" && item.severity === "error") ? "\u{1F534} has errors" : results.some((item) => item.check === "qa-evidence") ? "\u{1F7E1} findings" : "\u{2705} no findings"} |
| Code Review quality | ${results.some((item) => item.check === "code-review" && item.severity === "error") ? "\u{1F534} has errors" : results.some((item) => item.check === "code-review") ? "\u{1F7E1} findings" : "\u{2705} no findings"} |
| Failure routing | ${results.some((item) => item.check === "failure-routing" && item.severity === "error") ? "\u{1F534} has errors" : results.some((item) => item.check === "failure-routing") ? "\u{1F7E1} findings" : "\u{2705} no findings"} |
| Traceability | ${results.some((item) => item.check === "traceability" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "traceability") ? "🟡 findings" : "✅ no findings"} |
| Skill references | ${results.some((item) => item.check === "skill-reference" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "skill-reference") ? "🟡 findings" : "✅ no findings"} |
| Status policy | ${results.some((item) => item.check === "status-policy" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "status-policy") ? "🟡 findings" : "✅ no findings"} |
| Delivery metadata | ${results.some((item) => item.check === "delivery" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "delivery") ? "🟡 findings" : "✅ no findings"} |
| Execution graph JSON and dependencies | ${results.some((item) => item.check === "execution-graph" && item.severity === "error") ? "🔴 has errors" : "✅ no errors"} |
| WriteScope safety | ${results.some((item) => item.check === "write-scope" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "write-scope") ? "🟡 findings" : "✅ no findings"} |
| Decisions index | ${results.some((item) => item.check === "decisions-index") ? "🟡 findings" : "✅ no findings"} |
| Decision references | ${results.some((item) => item.check === "decision-references" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "decision-references") ? "🟡 findings" : "✅ no findings"} |
| Artifacts registry | ${results.some((item) => item.check === "artifacts-registry" && item.severity === "error") ? "🔴 has errors" : results.some((item) => item.check === "artifacts-registry") ? "🟡 findings" : "✅ no findings"} |
| Mermaid visual standard | ${results.some((item) => item.check === "mermaid") ? "🟡 findings" : "✅ no findings"} |
| Mermaid progress state | ${results.some((item) => item.check === "mermaid-progress") ? "🟡 findings" : "✅ no findings"} |
| Mermaid semantic state | ${results.some((item) => item.check === "mermaid-semantic") ? "🟡 findings" : "✅ no findings"} |
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
if (writeRegistry) {
  writeArtifactsRegistry();
}
validateUseCaseBundles();
validateRigorTiers();
validateTraceability();
validateStatusPolicy();
validateApprovalGates();
validateApprovalRecords();
validateStaleness();
validateValidationGates();
validateCodeEvidenceGates();
validateQaEvidenceQuality();
validateCodeReviewQuality();
validateBlockingFindingRoutes();
validateDeliveryMetadata();
validateExecutionGraphs();
validateContexts();
validateIdentityPolicy();
validateProductPrefixLinks(allFiles);
validateDecisionsIndex();
validateDecisionReferences();
validateArtifactsRegistry();
validateMermaidAndTemplates(allFiles);
validateMarkdownLinks(allFiles);

const report = generateReport();
if (writeReport) {
  const reportPath = path.join(root, "audits", "framework-validation-report.md");
  fs.writeFileSync(reportPath, report, "utf8");
  console.log(`Wrote ${rel(reportPath)}`);
  const readinessPath = path.join(root, "audits", "readiness", "framework-readiness.md");
  fs.writeFileSync(readinessPath, generateReadinessReport(), "utf8");
  console.log(`Wrote ${rel(readinessPath)}`);
}

const counts = severityCounts();
console.log(`Verdict: ${verdict()} (${counts.error} errors, ${counts.warning} warnings, ${counts.note} notes)`);
if (!writeReport && results.length) {
  console.log(results.map((item) => `${item.severity.toUpperCase()} ${item.check} ${item.file}: ${item.message}`).join("\n"));
}

process.exitCode = counts.error > 0 ? 1 : 0;
