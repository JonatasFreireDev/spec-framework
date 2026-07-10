#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";

const scriptsDir = path.dirname(fileURLToPath(import.meta.url));
const args = process.argv.slice(2);
const command = args[0] ?? "help";
const rest = args.slice(1);

const commands = {
  init: "init-product.mjs",
  validate: "validate-product.mjs",
  upgrade: "upgrade-product.mjs",
};

function usage() {
  console.log(`Usage: spec-framework <command> [options]

Commands:
  init --target <path>       Create a new product repo from starter/.
  validate [options]         Validate an initialized product repo.
  upgrade [--target <path>]  Refresh framework assets in a product repo.
  help                       Show this help.

Examples:
  spec-framework init --target ../my-product
  spec-framework validate
  spec-framework upgrade --target ../my-product
`);
}

if (command === "help" || command === "--help" || command === "-h") {
  usage();
  process.exit(0);
}

const script = commands[command];
if (!script) {
  console.error(`Unknown command: ${command}`);
  usage();
  process.exit(1);
}

const result = spawnSync(process.execPath, [path.join(scriptsDir, script), ...rest], {
  cwd: process.cwd(),
  stdio: "inherit",
  windowsHide: true,
});

process.exitCode = result.status ?? 1;
