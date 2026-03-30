#!/usr/bin/env bash
# Test: validate skill files are complete and well-formed.
# Usage: bash skills-dist-test.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_DIR="$(dirname "$SCRIPT_DIR")"
ERRORS=0

echo "=== Skills Validation Test ==="

# Check expected skills exist
for skill in seer-shared seer-market seer-api-explorer; do
  if [ -f "$CLI_DIR/skills/$skill/SKILL.md" ]; then
    echo "PASS: $skill/SKILL.md exists"
  else
    echo "FAIL: $skill/SKILL.md missing"
    ERRORS=$((ERRORS + 1))
  fi
done

# Check no Go files leaked into skills/
GO_COUNT=$(find "$CLI_DIR/skills" -name "*.go" | wc -l)
if [ "$GO_COUNT" -eq 0 ]; then
  echo "PASS: no .go files in skills/"
else
  echo "FAIL: $GO_COUNT .go files found in skills/"
  ERRORS=$((ERRORS + 1))
fi

# Check frontmatter consistency
for md in "$CLI_DIR"/skills/*/SKILL.md; do
  name=$(grep "^name:" "$md" | head -1 | awk '{print $2}')
  dir=$(basename "$(dirname "$md")")
  if [ "$name" = "$dir" ]; then
    echo "PASS: $dir frontmatter name matches directory"
  else
    echo "FAIL: $dir frontmatter name '$name' != directory '$dir'"
    ERRORS=$((ERRORS + 1))
  fi
done

echo ""
if [ "$ERRORS" -eq 0 ]; then
  echo "=== ALL PASSED ==="
else
  echo "=== $ERRORS FAILED ==="
  exit 1
fi
