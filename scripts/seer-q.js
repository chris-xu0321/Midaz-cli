#!/usr/bin/env node
// Deprecated shim: `seer-q` is now `midaz`. Prints a notice then execs the real binary.
const { execFileSync } = require("child_process");
const path = require("path");

const ext = process.platform === "win32" ? ".exe" : "";
const bin = path.join(__dirname, "..", "bin", "midaz" + ext);

process.stderr.write(
  "note: `seer-q` is deprecated — use `midaz` directly. This shim will be removed in v0.7.\n"
);

try {
  execFileSync(bin, process.argv.slice(2), { stdio: "inherit" });
} catch (e) {
  process.exit(e.status || 1);
}
