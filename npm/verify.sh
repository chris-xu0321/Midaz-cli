#!/usr/bin/env bash
# Validates the npm packaging pipeline output after build.sh has run.
# Usage: bash verify.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_DIR="$(dirname "$SCRIPT_DIR")"
NPM_DIST="$CLI_DIR/dist/npm"
VERSION="$(node -p "require('$SCRIPT_DIR/package.json').version")"

EXPECTED_PLATFORMS=("darwin-arm64" "darwin-x64" "linux-arm64" "linux-x64" "win32-arm64" "win32-x64")
ERRORS=0

echo "=== Verifying seer-q v${VERSION} npm packages ==="

# 1. Check all 6 platform package dirs exist with binary + package.json
for plat in "${EXPECTED_PLATFORMS[@]}"; do
  pkg_dir="$NPM_DIST/@midaz/seer-cli-${plat}"
  if [[ "$plat" == win32-* ]]; then bin="bin/seer-q.exe"; else bin="bin/seer-q"; fi

  if [ ! -f "$pkg_dir/package.json" ]; then
    echo "FAIL: missing $pkg_dir/package.json"
    ERRORS=$((ERRORS + 1))
  fi
  if [ ! -f "$pkg_dir/$bin" ]; then
    echo "FAIL: missing $pkg_dir/$bin"
    ERRORS=$((ERRORS + 1))
  fi

  # Verify platform package version matches
  pkg_ver="$(node -p "require('$pkg_dir/package.json').version" 2>/dev/null || echo MISSING)"
  if [ "$pkg_ver" != "$VERSION" ]; then
    echo "FAIL: $plat version mismatch: expected $VERSION, got $pkg_ver"
    ERRORS=$((ERRORS + 1))
  fi
done

# 2. Check meta package has optionalDependencies for all 6 platforms
META="$NPM_DIST/meta/package.json"
if [ ! -f "$META" ]; then
  echo "FAIL: missing $META"
  ERRORS=$((ERRORS + 1))
else
  meta_ver="$(node -p "require('$META').version")"
  if [ "$meta_ver" != "$VERSION" ]; then
    echo "FAIL: meta version mismatch: expected $VERSION, got $meta_ver"
    ERRORS=$((ERRORS + 1))
  fi

  for plat in "${EXPECTED_PLATFORMS[@]}"; do
    dep_ver="$(node -p "require('$META').optionalDependencies?.['@midaz/seer-cli-${plat}'] ?? 'MISSING'")"
    if [ "$dep_ver" != "$VERSION" ]; then
      echo "FAIL: meta optionalDependencies['@midaz/seer-cli-${plat}'] = $dep_ver (expected $VERSION)"
      ERRORS=$((ERRORS + 1))
    fi
  done
fi

# 3. Check run.js exists in meta
if [ ! -f "$NPM_DIST/meta/scripts/run.js" ]; then
  echo "FAIL: missing meta/scripts/run.js"
  ERRORS=$((ERRORS + 1))
fi

# 4. npm pack --dry-run all 7 packages
for plat in "${EXPECTED_PLATFORMS[@]}"; do
  pkg_dir="$NPM_DIST/@midaz/seer-cli-${plat}"
  if ! (cd "$pkg_dir" && npm pack --dry-run > /dev/null 2>&1); then
    echo "FAIL: npm pack failed for @midaz/seer-cli-${plat}"
    ERRORS=$((ERRORS + 1))
  fi
done
if ! (cd "$NPM_DIST/meta" && npm pack --dry-run > /dev/null 2>&1); then
  echo "FAIL: npm pack failed for @midaz/seer-cli (meta)"
  ERRORS=$((ERRORS + 1))
fi

echo ""
if [ "$ERRORS" -eq 0 ]; then
  echo "PASS: all 7 packages verified for v${VERSION}"
else
  echo "FAIL: $ERRORS error(s) found"
  exit 1
fi
