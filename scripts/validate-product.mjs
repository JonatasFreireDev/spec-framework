#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const cwd = process.cwd();
const args = process.argv.slice(2);

function argValue(name, fallback = "") {
  const index = args.indexOf(name);
  return index === -1 ? fallback : args[index + 1] ?? fallback;
}

function withoutOwnArgs(values) {
  const stripped = [];
  for (let index = 0; index < values.length; index += 1) {
    const value = values[index];
    if (value === "--product-root" || value === "--framework-root") {
      index += 1;
      continue;
    }
    stripped.push(value);
  }
  return stripped;
}

function findFrameworkRoot() {
  const explicit = argValue("--framework-root", "");
  if (explicit) return path.resolve(cwd, explicit);
  const installed = path.resolve(cwd, ".spec-framework");
  if (fs.existsSync(installed)) return installed;
  return path.resolve(scriptDir, "..");
}

const frameworkRoot = findFrameworkRoot();
const productRoot = path.resolve(cwd, argValue("--product-root", "product"));
const validator = path.join(frameworkRoot, "validators", "framework-validator.mjs");

if (!fs.existsSync(validator)) {
  console.error(`Validator not found: ${validator}`);
  process.exit(1);
}

const result = spawnSync(
  process.execPath,
  [
    validator,
    "--product-root",
    path.relative(cwd, productRoot) || ".",
    "--framework-root",
    path.relative(cwd, frameworkRoot) || ".",
    ...withoutOwnArgs(args),
  ],
  {
    cwd,
    stdio: "inherit",
    windowsHide: true,
  }
);

process.exitCode = result.status ?? 1;
