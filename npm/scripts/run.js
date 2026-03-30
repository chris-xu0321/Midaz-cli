#!/usr/bin/env node
const { execFileSync } = require("child_process");
const path = require("path");

const platformKey = `${process.platform}-${process.arch === "x64" ? "x64" : process.arch}`;
const pkgName = `@midaz/seer-cli-${platformKey}`;
const ext = process.platform === "win32" ? ".exe" : "";

let binPath;
try {
  const pkgDir = path.dirname(require.resolve(`${pkgName}/package.json`));
  binPath = path.join(pkgDir, "bin", "seer-q" + ext);
} catch {
  console.error(`Unsupported platform: ${process.platform}-${process.arch}`);
  console.error(`The platform package ${pkgName} was not installed.`);
  console.error(`Supported: darwin-arm64, darwin-x64, linux-arm64, linux-x64, win32-arm64, win32-x64`);
  process.exit(1);
}

try {
  execFileSync(binPath, process.argv.slice(2), { stdio: "inherit" });
} catch (e) {
  process.exit(e.status || 1);
}
