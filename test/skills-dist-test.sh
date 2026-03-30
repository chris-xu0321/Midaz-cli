#!/usr/bin/env bash
# Test: skills distribution artifact is complete and valid.
# Usage: bash skills-dist-test.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_DIR="$(dirname "$SCRIPT_DIR")"
ERRORS=0

echo "=== Skills Distribution Test ==="

# Publish to temp directory
TMPDIR=$(mktemp -d)
bash "$CLI_DIR/npm/publish-skills.sh" "$TMPDIR" 0.0.0-test 2>&1 | tail -5

# Check expected skills
for skill in seer-shared seer-market seer-api-explorer; do
  if [ -f "$TMPDIR/skills/$skill/SKILL.md" ]; then
    echo "PASS: $skill/SKILL.md exists"
  else
    echo "FAIL: $skill/SKILL.md missing"
    ERRORS=$((ERRORS + 1))
  fi
done

# Check command wrapper (advisory — target-specific, not required)
if [ -f "$TMPDIR/commands/claude/seer.md" ]; then
  echo "PASS: commands/claude/seer.md exists"
else
  echo "WARN: commands/claude/seer.md missing (optional Claude asset)"
fi

# Check no Go files leaked
GO_COUNT=$(find "$TMPDIR" -name "*.go" | wc -l)
if [ "$GO_COUNT" -eq 0 ]; then
  echo "PASS: no .go files in artifact"
else
  echo "FAIL: $GO_COUNT .go files found"
  ERRORS=$((ERRORS + 1))
fi

# Check frontmatter consistency
for md in "$TMPDIR"/skills/*/SKILL.md; do
  name=$(grep "^name:" "$md" | head -1 | awk '{print $2}')
  dir=$(basename "$(dirname "$md")")
  if [ "$name" = "$dir" ]; then
    echo "PASS: $dir frontmatter name matches directory"
  else
    echo "FAIL: $dir frontmatter name '$name' != directory '$dir'"
    ERRORS=$((ERRORS + 1))
  fi
done

rm -rf "$TMPDIR"
echo ""
if [ "$ERRORS" -eq 0 ]; then
  echo "=== ALL PASSED ==="
else
  echo "=== $ERRORS FAILED ==="
  exit 1
fi
