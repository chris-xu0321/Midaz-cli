#!/usr/bin/env bash
set -euo pipefail

# Environment checks
for cmd in node goreleaser npm; do
  if ! command -v "$cmd" > /dev/null 2>&1; then
    echo "ERROR: required command not found: $cmd"
    exit 1
  fi
done

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_DIR="$(dirname "$SCRIPT_DIR")"
NPM_DIST="$CLI_DIR/dist/npm"
DRY_RUN="${1:-}"  # pass "--dry-run" for testing

# Read version from package.json (single version authority)
VERSION="$(node -p "require('$SCRIPT_DIR/package.json').version")"
REGISTRY="$(node -p "require('$SCRIPT_DIR/package.json').publishConfig.registry")"

echo "=== seer-q npm release v${VERSION} ==="
echo "  registry: ${REGISTRY}"

# Guard: refuse to publish to public npm
if [[ "$REGISTRY" == "https://registry.npmjs.org" || -z "$REGISTRY" ]]; then
  echo "ERROR: publishConfig.registry must point to a private registry"
  echo "  Current: ${REGISTRY:-<empty>}"
  echo "  Set it in npm/package.json -> publishConfig.registry"
  exit 1
fi

# Step 1: Cross-compile
echo ""
echo "=== Building cross-platform binaries ==="
cd "$CLI_DIR"
GORELEASER_CURRENT_TAG="v${VERSION}" goreleaser release --clean --skip=publish

# Step 2: Package
echo ""
echo "=== Packaging npm platform packages ==="
bash "$SCRIPT_DIR/build.sh"

# Step 2.5: Verify
echo ""
echo "=== Verifying packages ==="
bash "$SCRIPT_DIR/verify.sh"

# Step 3: Publish platform packages (must go first)
echo ""
echo "=== Publishing platform packages ==="
for pkg_dir in "$NPM_DIST"/@midaz/seer-cli-*/; do
  pkg_name="$(node -p "require('${pkg_dir}package.json').name")"
  echo "  publishing ${pkg_name}..."
  cd "$pkg_dir"
  npm publish --access restricted $DRY_RUN
done

# Step 4: Publish meta package (after platform packages are available)
echo ""
echo "=== Publishing meta package ==="
cd "$NPM_DIST/meta"
npm publish --access restricted $DRY_RUN

echo ""
echo "=== Done: @midaz/seer-cli@${VERSION} ==="
echo "  Install: npm install -g @midaz/seer-cli"
