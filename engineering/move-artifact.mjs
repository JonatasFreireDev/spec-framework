#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";

const root = process.cwd();
const args = process.argv.slice(2);

function argValue(name) {
  const index = args.indexOf(name);
  return index === -1 ? "" : args[index + 1] ?? "";
}

const fromArg = argValue("--from");
const toArg = argValue("--to");
const dryRun = args.includes("--dry-run");

function usage() {
  console.log("Usage: node engineering/move-artifact.mjs --from <old-path> --to <new-path> [--dry-run]");
}

function rel(filePath) {
  return path.relative(root, filePath).replaceAll(path.sep, "/");
}

function isInsideRoot(filePath) {
  const resolvedRoot = path.resolve(root);
  const resolved = path.resolve(filePath);
  return resolved === resolvedRoot || resolved.startsWith(`${resolvedRoot}${path.sep}`);
}

function walk(dir, output = []) {
  if (!fs.existsSync(dir)) return output;
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (entry.name === ".git") continue;
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) walk(full, output);
    else output.push(full);
  }
  return output;
}

function readText(file) {
  return fs.readFileSync(file, "utf8").replace(/^\uFEFF/, "");
}

function writeText(file, text) {
  fs.writeFileSync(file, text.replace(/\r\n?/g, "\n"), "utf8");
}

function normalizeMarkdownLinkTarget(target) {
  let clean = target.trim();
  if (!clean || clean.startsWith("#")) return null;
  if (/^(https?:|mailto:|tel:|javascript:)/i.test(clean)) return null;
  if (clean.startsWith("<") && clean.endsWith(">")) clean = clean.slice(1, -1);
  const hashIndex = clean.indexOf("#");
  const anchor = hashIndex === -1 ? "" : clean.slice(hashIndex);
  const pathPart = hashIndex === -1 ? clean : clean.slice(0, hashIndex);
  if (!pathPart || pathPart.includes("[") || pathPart.includes("]")) return null;
  try {
    return { pathPart: decodeURI(pathPart), anchor };
  } catch {
    return { pathPart, anchor };
  }
}

function markdownRelativeTarget(fromFile, targetAbs, anchor) {
  let target = path.relative(path.dirname(fromFile), targetAbs).replaceAll(path.sep, "/");
  if (!target.startsWith(".")) target = `./${target}`;
  return `${target}${anchor}`;
}

function rewriteMarkdownLinks(file, oldAbs, newAbs) {
  const text = readText(file);
  const linkPattern = /(?<!!)\[([^\]\n]+)\]\(([^)\n]+)\)/g;
  let changed = false;
  const next = text.replace(linkPattern, (full, label, target) => {
    const parsed = normalizeMarkdownLinkTarget(target);
    if (!parsed) return full;
    const resolved = path.resolve(path.dirname(file), parsed.pathPart);
    if (resolved !== oldAbs && !resolved.startsWith(`${oldAbs}${path.sep}`)) return full;
    const suffix = path.relative(oldAbs, resolved);
    const movedTarget = suffix ? path.join(newAbs, suffix) : newAbs;
    changed = true;
    return `[${label}](${markdownRelativeTarget(file, movedTarget, parsed.anchor)})`;
  });
  if (changed && !dryRun) writeText(file, next);
  return changed;
}

function rewriteJsonValue(value, oldRel, newRel) {
  if (Array.isArray(value)) {
    let changed = false;
    const next = value.map((item) => {
      const result = rewriteJsonValue(item, oldRel, newRel);
      changed = changed || result.changed;
      return result.value;
    });
    return { value: next, changed };
  }
  if (value && typeof value === "object") {
    let changed = false;
    const next = Object.fromEntries(Object.entries(value).map(([key, item]) => {
      const result = rewriteJsonValue(item, oldRel, newRel);
      changed = changed || result.changed;
      return [key, result.value];
    }));
    return { value: next, changed };
  }
  if (typeof value !== "string") return { value, changed: false };
  if (value === oldRel) return { value: newRel, changed: true };
  if (value.startsWith(`${oldRel}/`)) return { value: `${newRel}${value.slice(oldRel.length)}`, changed: true };
  return { value, changed: false };
}

function rewriteJsonPaths(file, oldRel, newRel) {
  const text = readText(file);
  let parsed;
  try {
    parsed = JSON.parse(text);
  } catch {
    return false;
  }
  const result = rewriteJsonValue(parsed, oldRel, newRel);
  if (!result.changed) return false;
  if (!dryRun) writeText(file, `${JSON.stringify(result.value, null, 2)}\n`);
  return true;
}

function freeTextMentions(file, oldRel) {
  const text = readText(file);
  if (!text.includes(oldRel)) return [];
  const mentions = [];
  for (const [index, line] of text.split(/\r?\n/).entries()) {
    if (!line.includes(oldRel)) continue;
    if (/\[[^\]]+\]\([^)]+\)/.test(line)) continue;
    mentions.push(`${rel(file)}:${index + 1}: ${line.trim()}`);
  }
  return mentions;
}

if (!fromArg || !toArg) {
  usage();
  process.exit(1);
}

const fromAbs = path.resolve(root, fromArg);
const toAbs = path.resolve(root, toArg);
if (!isInsideRoot(fromAbs) || !isInsideRoot(toAbs)) {
  console.error("Both --from and --to must stay inside the repository.");
  process.exit(1);
}
if (!fs.existsSync(fromAbs)) {
  console.error(`Source does not exist: ${fromArg}`);
  process.exit(1);
}
if (fs.existsSync(toAbs)) {
  console.error(`Target already exists: ${toArg}`);
  process.exit(1);
}

const oldRel = rel(fromAbs);
const newRel = rel(toAbs);
const filesBeforeMove = walk(root).filter((file) => /\.(md|json)$/.test(file));

if (!dryRun) {
  fs.mkdirSync(path.dirname(toAbs), { recursive: true });
  fs.renameSync(fromAbs, toAbs);
}

const rewrites = [];
const mentionReports = [];
for (const file of filesBeforeMove) {
  const currentFile = file === fromAbs || file.startsWith(`${fromAbs}${path.sep}`)
    ? path.join(toAbs, path.relative(fromAbs, file))
    : file;
  if (!fs.existsSync(currentFile)) continue;

  if (currentFile.endsWith(".md") && rewriteMarkdownLinks(currentFile, fromAbs, toAbs)) {
    rewrites.push(`${rel(currentFile)} markdown-links`);
  }
  if (currentFile.endsWith(".json") && rewriteJsonPaths(currentFile, oldRel, newRel)) {
    rewrites.push(`${rel(currentFile)} json-paths`);
  }
  mentionReports.push(...freeTextMentions(currentFile, oldRel));
}

console.log(`${dryRun ? "Dry run" : "Moved"}: ${oldRel} -> ${newRel}`);
console.log(`Rewritten files: ${rewrites.length}`);
for (const item of rewrites) console.log(`- ${item}`);
console.log(`Free-text mentions requiring review: ${mentionReports.length}`);
for (const item of mentionReports) console.log(`- ${item}`);
