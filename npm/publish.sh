#!/usr/bin/env bash
set -euo pipefail

# Publish @midaz/cli to the public npm registry.
# Binaries must be available on GitHub Releases (via goreleaser).
# The npm package uses a postinstall hook to download the binary.

for cmd in node goreleaser npm; do
  if ! command -v "$cmd" > /dev/null 2>&1; then
    echo "ERROR: required command not found: $cmd"
    exit 1
  fi
done

CLI_DIR="$(cd "$(dirname "$0")/.." && pwd)"
DRY_RUN="${1:-}"  # pass "--dry-run" for testing

VERSION="$(node -p "require('$CLI_DIR/package.json').version")"

echo "=== seer-q npm release v${VERSION} ==="

# Step 1: Cross-compile (goreleaser creates GitHub Release with binaries)
echo ""
echo "=== Building cross-platform binaries ==="
cd "$CLI_DIR"
GORELEASER_CURRENT_TAG="v${VERSION}" goreleaser release --clean

# Step 2: Publish npm package
echo ""
echo "=== Publishing @midaz/cli ==="
cd "$CLI_DIR"
npm publish --access public $DRY_RUN

echo ""
echo "=== Done: @midaz/cli@${VERSION} ==="
echo "  Install: npm install -g @midaz/cli"
