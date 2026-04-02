# Target Compatibility

Last updated: 2026-04-02

## Support Matrix

| Target | Status | Install Method |
|--------|--------|----------------|
| Claude Code | Supported | `seer-q setup claude --yes` |
| Codex | Supported | `seer-q setup codex --yes` |

Skills are embedded in the `seer-q` binary and installed via `seer-q setup`. The legacy `npx skills add SparkssL/Midaz-cli -y -g` method also works for users with npm.

## Claude Code

### Install

```bash
# One-line (recommended):
curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh

# Or manually:
npm install -g @midaz/cli
seer-q setup claude --yes
```

Skills are written to `~/.claude/skills/`. If existing symlinks point to `~/.agents/skills/`, writes are transparently resolved to the symlink target.

### Verify

```bash
seer-q doctor
seer-q search "test"
```

## Codex

### Install

```bash
# One-line (recommended):
curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh

# Or manually:
npm install -g @midaz/cli
seer-q setup codex --yes
```

Skills are written to `~/.codex/skills/`.

### Verify

```bash
seer-q version
seer-q doctor
```

## Adding a New Target

1. Add the target to `resolveTargets()` in `internal/cmd/setup/setup.go`.
2. Verify `seer-q` runs on the target's supported platforms.
3. Run `seer-q setup <target> --yes` and verify skills are discovered by the agent.
4. Update this matrix with tested results.
