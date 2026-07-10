#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";
import { execSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const frameworkRepo = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
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

function rel(filePath, from = frameworkRepo) {
  return path.relative(from, filePath).replaceAll(path.sep, "/");
}

function copyDir(source, target) {
  if (!fs.existsSync(source)) return;
  fs.mkdirSync(path.dirname(target), { recursive: true });
  fs.cpSync(source, target, {
    recursive: true,
    force: true,
    filter: (item) => !item.split(path.sep).includes(".git"),
  });
}

function copyFile(source, target) {
  fs.mkdirSync(path.dirname(target), { recursive: true });
  fs.copyFileSync(source, target);
}

function gitVersion() {
  try {
    return execSync("git rev-parse --short HEAD", {
      cwd: frameworkRepo,
      encoding: "utf8",
      stdio: ["ignore", "pipe", "ignore"],
    }).trim();
  } catch {
    return "local";
  }
}

function readJson(file) {
  return JSON.parse(fs.readFileSync(file, "utf8"));
}

function writeJson(file, value) {
  fs.writeFileSync(file, `${JSON.stringify(value, null, 2)}\n`, "utf8");
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

const specRoot = path.join(target, ".spec-framework");
copyFile(path.join(frameworkRepo, "FRAMEWORK.md"), path.join(specRoot, "FRAMEWORK.md"));
copyDir(path.join(frameworkRepo, "engineering", "decisions"), path.join(specRoot, "decisions"));
copyDir(path.join(frameworkRepo, ".codex", "skills"), path.join(specRoot, "skills"));
copyDir(path.join(frameworkRepo, "knowledge", "templates"), path.join(specRoot, "templates"));
copyDir(path.join(frameworkRepo, "engineering", "validators"), path.join(specRoot, "validators"));
copyFile(path.join(frameworkRepo, "engineering", "move-artifact.mjs"), path.join(specRoot, "tools", "move-artifact.mjs"));

if (mirrorCodexSkills) {
  copyDir(path.join(frameworkRepo, ".codex", "skills"), path.join(target, ".codex", "skills"));
}

const workflow = `name: Framework Validation

on:
  pull_request:
  push:
    branches:
      - master
      - main

permissions:
  contents: read

jobs:
  validate:
    name: Validate product framework
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "22"

      - name: Syntax check
        run: node --check .spec-framework/validators/framework-validator.mjs

      - name: Framework validator
        run: node .spec-framework/validators/framework-validator.mjs --product-root product --framework-root .spec-framework
`;
copyFile(path.join(frameworkRepo, "AGENTS.md"), path.join(specRoot, "AGENTS.framework.md"));
fs.mkdirSync(path.join(target, ".github", "workflows"), { recursive: true });
fs.writeFileSync(path.join(target, ".github", "workflows", "framework-validation.yml"), workflow, "utf8");

const manifestFile = path.join(specRoot, "manifest.json");
const manifest = readJson(manifestFile);
manifest.version = gitVersion();
manifest.installed_assets = {
  framework_document: true,
  decisions: true,
  skills: true,
  templates: true,
  validators: true,
  tools: true,
  ci: true,
};
writeJson(manifestFile, manifest);

const productFrameworkFile = path.join(target, "product", ".product", "framework.json");
const productFramework = readJson(productFrameworkFile);
productFramework.version = manifest.version;
productFramework.framework_assets_path = "../.spec-framework";
productFramework.product_root = ".";
productFramework.installed_assets = {
  framework_document: true,
  decisions: true,
  skills: true,
  templates: true,
  validators: true,
  tools: true,
  ci: true,
};
writeJson(productFrameworkFile, productFramework);

console.log(`Initialized Spec Framework product at ${target}`);
console.log(`- Product root: ${rel(path.join(target, "product"), target)}`);
console.log(`- Framework assets: ${rel(specRoot, target)}`);
console.log(`- Codex skills mirror: ${mirrorCodexSkills ? ".codex/skills" : "disabled"}`);
console.log("Next: fill product/foundation and run:");
console.log("node .spec-framework/validators/framework-validator.mjs --product-root product --framework-root .spec-framework");
