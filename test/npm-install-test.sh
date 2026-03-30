#!/usr/bin/env bash
# Test: npm install from local tarballs + run via bin shim.
# Usage: bash npm-install-test.sh
# Requires: npm packages built in dist/npm/
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_DIR="$(dirname "$SCRIPT_DIR")"
NPM_DIST="$CLI_DIR/dist/npm"

if [ ! -d "$NPM_DIST/meta" ]; then
  echo "ERROR: npm packages not built. Run 'bash npm/build.sh' first."
  exit 1
fi

echo "=== npm Install Test ==="

# Detect current platform
PLAT="$(node -p "process.platform + '-' + (process.arch === 'x64' ? 'x64' : process.arch)")"
echo "  platform: $PLAT"

PLAT_DIR="$NPM_DIST/@midaz/seer-cli-${PLAT}"
if [ ! -d "$PLAT_DIR" ]; then
  echo "SKIP: platform package not found for $PLAT (cross-platform build required)"
  exit 0
fi

# Pack tarballs
TMPDIR=$(mktemp -d)
(cd "$PLAT_DIR" && npm pack --pack-destination "$TMPDIR" > /dev/null 2>&1)
(cd "$NPM_DIST/meta" && npm pack --pack-destination "$TMPDIR" > /dev/null 2>&1)

# Install in isolated directory
mkdir -p "$TMPDIR/test"
(cd "$TMPDIR/test" && npm init -y > /dev/null 2>&1)
(cd "$TMPDIR/test" && npm install "$TMPDIR"/midaz-seer-cli-${PLAT}-*.tgz > /dev/null 2>&1)
(cd "$TMPDIR/test" && npm install "$TMPDIR"/midaz-seer-cli-[0-9]*.tgz > /dev/null 2>&1)

# Test via bin shim (not direct node invocation)
# This validates the "bin" entry in package.json is wired correctly
BIN="$(cd "$TMPDIR/test" && npm bin)/seer-q"
if [ ! -f "$BIN" ] && [ ! -f "${BIN}.cmd" ]; then
  echo "FAIL: seer-q bin shim not created by npm install"
  rm -rf "$TMPDIR"
  exit 1
fi
echo "  bin shim: $BIN"

VERSION=$(cd "$TMPDIR/test" && "$(npm bin)/seer-q" version --format json 2>&1 | node -e "process.stdin.on('data',d=>console.log(JSON.parse(d).data?.version??'FAIL'))")
echo "  installed version: $VERSION"

if [ "$VERSION" = "FAIL" ] || [ -z "$VERSION" ]; then
  echo "FAIL: bin shim did not return valid version"
  rm -rf "$TMPDIR"
  exit 1
fi

echo "PASS: npm install + bin shim works"
rm -rf "$TMPDIR"
