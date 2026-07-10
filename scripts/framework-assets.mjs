import fs from "node:fs";
import path from "node:path";
import { execSync } from "node:child_process";
import { fileURLToPath } from "node:url";

export const frameworkRepo = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");

export function rel(filePath, from = frameworkRepo) {
  return path.relative(from, filePath).replaceAll(path.sep, "/");
}

export function copyDir(source, target) {
  if (!fs.existsSync(source)) return;
  fs.mkdirSync(path.dirname(target), { recursive: true });
  fs.cpSync(source, target, {
    recursive: true,
    force: true,
    filter: (item) => !item.split(path.sep).includes(".git"),
  });
}

export function copyFile(source, target) {
  fs.mkdirSync(path.dirname(target), { recursive: true });
  fs.copyFileSync(source, target);
}

export function gitVersion() {
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

export function readJson(file) {
  return JSON.parse(fs.readFileSync(file, "utf8"));
}

export function writeJson(file, value) {
  fs.writeFileSync(file, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

export function productWorkflow() {
  return `name: Framework Validation

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
        run: |
          node --check .spec-framework/validators/framework-validator.mjs
          node --check .spec-framework/tools/validate-product.mjs

      - name: Framework validator
        run: node .spec-framework/tools/validate-product.mjs
`;
}

export function installedAssets() {
  return {
    framework_document: true,
    decisions: true,
    skills: true,
    templates: true,
    validators: true,
    tools: true,
    ci: true,
  };
}

export function installFrameworkAssets(target, options = {}) {
  const mirrorCodexSkills = options.mirrorCodexSkills !== false;
  const specRoot = path.join(target, ".spec-framework");

  copyFile(path.join(frameworkRepo, "FRAMEWORK.md"), path.join(specRoot, "FRAMEWORK.md"));
  copyDir(path.join(frameworkRepo, "engineering", "decisions"), path.join(specRoot, "decisions"));
  copyDir(path.join(frameworkRepo, ".codex", "skills"), path.join(specRoot, "skills"));
  copyDir(path.join(frameworkRepo, "knowledge", "templates"), path.join(specRoot, "templates"));
  copyDir(path.join(frameworkRepo, "engineering", "validators"), path.join(specRoot, "validators"));
  copyFile(path.join(frameworkRepo, "engineering", "move-artifact.mjs"), path.join(specRoot, "tools", "move-artifact.mjs"));
  copyFile(path.join(frameworkRepo, "scripts", "validate-product.mjs"), path.join(specRoot, "tools", "validate-product.mjs"));
  copyFile(path.join(frameworkRepo, "AGENTS.md"), path.join(specRoot, "AGENTS.framework.md"));

  if (mirrorCodexSkills) {
    copyDir(path.join(frameworkRepo, ".codex", "skills"), path.join(target, ".codex", "skills"));
  }

  fs.mkdirSync(path.join(target, ".github", "workflows"), { recursive: true });
  fs.writeFileSync(path.join(target, ".github", "workflows", "framework-validation.yml"), productWorkflow(), "utf8");

  const version = gitVersion();
  const assets = installedAssets();

  const manifestFile = path.join(specRoot, "manifest.json");
  if (fs.existsSync(manifestFile)) {
    const manifest = readJson(manifestFile);
    manifest.version = version;
    manifest.installed_assets = assets;
    writeJson(manifestFile, manifest);
  }

  const productFrameworkFile = path.join(target, "product", ".product", "framework.json");
  if (fs.existsSync(productFrameworkFile)) {
    const productFramework = readJson(productFrameworkFile);
    productFramework.version = version;
    productFramework.framework_assets_path = "../.spec-framework";
    productFramework.product_root = ".";
    productFramework.installed_assets = assets;
    writeJson(productFrameworkFile, productFramework);
  }

  return { specRoot, version, mirrorCodexSkills };
}
