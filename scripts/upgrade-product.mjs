#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";
import {
  installFrameworkAssets,
  rel,
} from "./framework-assets.mjs";

const args = process.argv.slice(2);

function argValue(name, fallback = "") {
  const index = args.indexOf(name);
  return index === -1 ? fallback : args[index + 1] ?? fallback;
}

const targetArg = argValue("--target", args[0] ?? ".");
const mirrorCodexSkills = !args.includes("--no-codex-skills");

function usage() {
  console.log("Usage: node scripts/upgrade-product.mjs [--target <path>] [--no-codex-skills]");
}

if (args.includes("--help")) {
  usage();
  process.exit(0);
}

const target = path.resolve(process.cwd(), targetArg);
const productRoot = path.join(target, "product");
const specRoot = path.join(target, ".spec-framework");

if (!fs.existsSync(productRoot) || !fs.existsSync(specRoot)) {
  console.error(`Target does not look like an initialized Spec Framework product: ${target}`);
  console.error("Expected both product/ and .spec-framework/.");
  process.exit(1);
}

const result = installFrameworkAssets(target, { mirrorCodexSkills });

console.log(`Upgraded Spec Framework assets at ${target}`);
console.log(`- Product root preserved: ${rel(productRoot, target)}`);
console.log(`- Framework assets updated: ${rel(result.specRoot, target)}`);
console.log(`- Version: ${result.version}`);
console.log(`- Codex skills mirror: ${mirrorCodexSkills ? ".codex/skills" : "disabled"}`);
console.log("Next: run:");
console.log("node .spec-framework/tools/validate-product.mjs");
