#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";
import {
  copyDir,
  frameworkRepo,
  installFrameworkAssets,
  rel,
} from "./framework-assets.mjs";

const args = process.argv.slice(2);

function argValue(name, fallback = "") {
  const index = args.indexOf(name);
  return index === -1 ? fallback : args[index + 1] ?? fallback;
}

const targetArg = argValue("--target", args[0] ?? "");
const force = args.includes("--force");
const mirrorCodexSkills = !args.includes("--no-codex-skills");

function usage() {
  console.log("Usage: node scripts/init-product.mjs --target <path> [--force] [--no-codex-skills]");
}

if (!targetArg) {
  usage();
  process.exit(1);
}

const target = path.resolve(process.cwd(), targetArg);
if (fs.existsSync(target) && fs.readdirSync(target).length > 0 && !force) {
  console.error(`Target is not empty: ${target}`);
  console.error("Use --force to merge starter files into an existing directory.");
  process.exit(1);
}

fs.mkdirSync(target, { recursive: true });

copyDir(path.join(frameworkRepo, "starter"), target);

const { specRoot } = installFrameworkAssets(target, { mirrorCodexSkills });

console.log(`Initialized Spec Framework product at ${target}`);
console.log(`- Product root: ${rel(path.join(target, "product"), target)}`);
console.log(`- Framework assets: ${rel(specRoot, target)}`);
console.log(`- Codex skills mirror: ${mirrorCodexSkills ? ".codex/skills" : "disabled"}`);
console.log("Next: fill product/foundation and run:");
console.log("node .spec-framework/tools/validate-product.mjs");
