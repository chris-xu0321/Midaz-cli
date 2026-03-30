#!/usr/bin/env bash
set -euo pipefail

# Environment checks
if [ "${BASH_VERSINFO[0]}" -lt 4 ]; then
  echo "ERROR: Bash 4+ required (for associative arrays). Current: $BASH_VERSION"
  exit 1
fi
for cmd in node sed tar unzip mktemp; do
  if ! command -v "$cmd" > /dev/null 2>&1; then
    echo "ERROR: required command not found: $cmd"
    exit 1
  fi
done

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_DIR="$(dirname "$SCRIPT_DIR")"
DIST="$CLI_DIR/dist"
NPM_DIST="$DIST/npm"

# Read version and registry from meta package.json (single source of truth)
VERSION="$(node -p "require('$SCRIPT_DIR/package.json').version")"
REGISTRY="$(node -p "require('$SCRIPT_DIR/package.json').publishConfig.registry")"

echo "=== Packaging seer-q v${VERSION} ==="

# Platform map: npm_key -> goreleaser archive suffix
# goreleaser uses goos-goarch naming (controlled by .goreleaser.yml name_template)
# npm uses process.platform-process.arch naming
declare -A PLATFORMS=(
  ["darwin-arm64"]="darwin-arm64"
  ["darwin-x64"]="darwin-amd64"
  ["linux-arm64"]="linux-arm64"
  ["linux-x64"]="linux-amd64"
  ["win32-arm64"]="windows-arm64"
  ["win32-x64"]="windows-amd64"
)

rm -rf "$NPM_DIST"
OPT_DEPS=""

for npm_key in "${!PLATFORMS[@]}"; do
  go_key="${PLATFORMS[$npm_key]}"
  os="${npm_key%-*}"
  arch="${npm_key#*-}"

  # Determine archive name and extension
  if [[ "$os" == "win32" ]]; then
    archive="$DIST/seer-q-${VERSION}-${go_key}.zip"
    bin_ext=".exe"
  else
    archive="$DIST/seer-q-${VERSION}-${go_key}.tar.gz"
    bin_ext=""
  fi

  # Verify archive exists
  if [ ! -f "$archive" ]; then
    echo "ERROR: archive not found: $archive"
    echo "  Expected from goreleaser name_template: seer-q-VERSION-OS-ARCH"
    echo "  Run 'make release' first"
    exit 1
  fi

  # Create platform package directory
  pkg_dir="$NPM_DIST/@midaz/seer-cli-${npm_key}"
  mkdir -p "$pkg_dir/bin"

  # Extract binary from archive
  tmpdir="$(mktemp -d)"
  if [[ "$os" == "win32" ]]; then
    unzip -q "$archive" -d "$tmpdir"
  else
    tar -xzf "$archive" -C "$tmpdir"
  fi
  cp "$tmpdir/seer-q${bin_ext}" "$pkg_dir/bin/seer-q${bin_ext}"
  rm -rf "$tmpdir"

  # Generate platform package.json from template
  sed -e "s/{{OS}}/$os/g" \
      -e "s/{{ARCH}}/$arch/g" \
      -e "s/{{VERSION}}/$VERSION/g" \
      -e "s|{{REGISTRY}}|$REGISTRY|g" \
      "$SCRIPT_DIR/platform-template/package.json.tmpl" \
      > "$pkg_dir/package.json"

  # Accumulate optionalDependencies entry
  OPT_DEPS+="\"@midaz/seer-cli-${npm_key}\": \"${VERSION}\","

  echo "  packed: @midaz/seer-cli-${npm_key}@${VERSION}"
done

# Inject optionalDependencies into meta package.json for publish
# Work on a copy in dist so we don't mutate the source file
META_PUBLISH="$NPM_DIST/meta"
mkdir -p "$META_PUBLISH/scripts"
cp "$SCRIPT_DIR/package.json" "$META_PUBLISH/package.json"
cp "$SCRIPT_DIR/scripts/run.js" "$META_PUBLISH/scripts/run.js"

# Remove trailing comma, wrap in object
OPT_DEPS="{${OPT_DEPS%,}}"
node -e "
  const fs = require('fs');
  const pkg = JSON.parse(fs.readFileSync('$META_PUBLISH/package.json', 'utf8'));
  pkg.optionalDependencies = ${OPT_DEPS};
  fs.writeFileSync('$META_PUBLISH/package.json', JSON.stringify(pkg, null, 2) + '\n');
"

echo "  packed: @midaz/seer-cli@${VERSION} (meta)"
echo ""
echo "Platform packages ready in $NPM_DIST"
