#!/usr/bin/env bash
# Smoke test: validates seer-q against a running API.
# Usage: bash smoke-test.sh [binary-path] [api-url]
# Hermetic: uses temp config, explicit API URL, forces JSON output.
set -euo pipefail

BIN="${1:-seer-q}"
API_URL="${2:-http://localhost:4000}"

# Hermetic environment: temp config with valid empty JSON, explicit API URL
export SEER_CONFIG_PATH="$(mktemp)"
echo '{}' > "$SEER_CONFIG_PATH"
export SEER_API_URL="$API_URL"
export SEER_FORMAT="json"
trap 'rm -f "$SEER_CONFIG_PATH"' EXIT

ERRORS=0

check() {
  local name="$1"; shift
  if output=$("$@" --format json 2>&1); then
    # Buffer all stdin before parsing (handles large responses)
    ok=$(echo "$output" | node -e "let b='';process.stdin.on('data',d=>b+=d);process.stdin.on('end',()=>{try{console.log(JSON.parse(b).ok)}catch{console.log(false)}})" 2>/dev/null)
    if [ "$ok" = "true" ]; then
      echo "PASS: $name"
    else
      echo "FAIL: $name (ok != true)"
      ERRORS=$((ERRORS + 1))
    fi
  else
    echo "FAIL: $name (exit code $?)"
    ERRORS=$((ERRORS + 1))
  fi
}

echo "=== Seer CLI Smoke Test ==="
echo "  binary: $BIN"
echo "  api_url: $API_URL"
echo "  config: $SEER_CONFIG_PATH (isolated)"
echo ""

# Core commands (no API required)
check "version"       "$BIN" version
check "schema"        "$BIN" schema
check "config list"   "$BIN" config list

# API commands (require running API)
check "health"        "$BIN" health
check "doctor"        "$BIN" doctor
check "market"        "$BIN" market
check "topics"        "$BIN" topics
check "search"        "$BIN" search "test"
check "snapshot"      "$BIN" snapshot
check "claims"        "$BIN" claims
check "sources"       "$BIN" sources
check "usage"         "$BIN" usage

echo ""
if [ "$ERRORS" -eq 0 ]; then
  echo "=== ALL PASSED ==="
else
  echo "=== $ERRORS FAILED ==="
  exit 1
fi
