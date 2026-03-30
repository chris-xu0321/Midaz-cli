#!/usr/bin/env bash
set -euo pipefail

# Publish @midaz/seer-cli to npm registry.
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
REGISTRY="$(node -p "require('$CLI_DIR/package.json').publishConfig.registry")"

echo "=== seer-q npm release v${VERSION} ==="
echo "  registry: ${REGISTRY}"

# Guard: refuse to publish to public npm
if [[ "$REGISTRY" == "https://registry.npmjs.org" || -z "$REGISTRY" ]]; then
  echo "ERROR: publishConfig.registry must point to a private registry"
  echo "  Current: ${REGISTRY:-<empty>}"
  echo "  Set it in package.json -> publishConfig.registry"
  exit 1
fi

# Step 1: Cross-compile (goreleaser creates GitHub Release with binaries)
echo ""
echo "=== Building cross-platform binaries ==="
cd "$CLI_DIR"
GORELEASER_CURRENT_TAG="v${VERSION}" goreleaser release --clean

# Step 2: Publish npm package
echo ""
echo "=== Publishing @midaz/seer-cli ==="
cd "$CLI_DIR"
npm publish --access restricted $DRY_RUN

echo ""
echo "=== Done: @midaz/seer-cli@${VERSION} ==="
echo "  Install: npm install -g @midaz/seer-cli"
