#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_DIR="$(dirname "$SCRIPT_DIR")"
SKILLS_SRC="$CLI_DIR/skills"
TARGET="${1:?Usage: publish-skills.sh <target-repo-dir> <version>}"
VERSION="${2:?Usage: publish-skills.sh <target-repo-dir> <version>}"

if [ ! -d "$SKILLS_SRC" ]; then
  echo "ERROR: skills source not found: $SKILLS_SRC"
  exit 1
fi

echo "=== Publishing seer-skills v${VERSION} ==="

SKILLS_DEST="$TARGET/skills"
rm -rf "$SKILLS_DEST"
mkdir -p "$SKILLS_DEST"

COUNT=0
for dir in "$SKILLS_SRC"/*/; do
  skill="$(basename "$dir")"
  [ -f "$dir/SKILL.md" ] || continue

  # Copy entire skill directory (SKILL.md, references/, templates/, etc.)
  cp -r "$dir" "$SKILLS_DEST/$skill"

  echo "  synced: $skill"
  COUNT=$((COUNT + 1))
done

if [ "$COUNT" -eq 0 ]; then
  echo "ERROR: no skills found in $SKILLS_SRC"
  exit 1
fi

# Clean non-distributable files
find "$SKILLS_DEST" \( -name "*.go" -o -name "*.swp" -o -name "*~" \
  -o -name ".DS_Store" -o -name "Thumbs.db" -o -name "*.bak" \) -delete

# Verify published skills have valid frontmatter
VERIFY_ERRORS=0
for skill_dir in "$SKILLS_DEST"/*/; do
  md="$skill_dir/SKILL.md"
  [ -f "$md" ] || continue
  if ! head -1 "$md" | grep -q "^---"; then
    echo "VERIFY FAIL: $md missing frontmatter block"
    VERIFY_ERRORS=$((VERIFY_ERRORS + 1))
  fi
  if ! grep -q "^name:" "$md"; then
    echo "VERIFY FAIL: $md missing name field"
    VERIFY_ERRORS=$((VERIFY_ERRORS + 1))
  fi
  if ! grep -q "^version:" "$md"; then
    echo "VERIFY FAIL: $md missing version field"
    VERIFY_ERRORS=$((VERIFY_ERRORS + 1))
  fi
done
if [ "$VERIFY_ERRORS" -gt 0 ]; then
  echo "ERROR: $VERIFY_ERRORS frontmatter validation failures"
  exit 1
fi

# Include Claude command wrapper as a target-specific asset
# Skill installers ignore files outside skills/; manual install uses this
COMMANDS_DEST="$TARGET/commands/claude"
if [ -f "$CLI_DIR/agent/cmd/seer.md" ]; then
  mkdir -p "$COMMANDS_DEST"
  cp "$CLI_DIR/agent/cmd/seer.md" "$COMMANDS_DEST/seer.md"
  echo "  synced: commands/claude/seer.md"
fi

echo ""
echo "=== $COUNT skills synced to $SKILLS_DEST ==="
echo ""
echo "Next steps:"
echo "  cd $TARGET"
echo "  git add -A && git commit -m 'skills: v$VERSION'"
echo "  git tag v$VERSION && git push origin main --tags"
