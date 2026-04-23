# Target Compatibility

Last updated: 2026-04-02

## Support Matrix

| Target | Status | Install Method |
|--------|--------|----------------|
| Claude Code | Supported | `midaz setup claude --yes` |
| Codex | Supported | `midaz setup codex --yes` |

Skills are embedded in the `midaz` binary and installed via `midaz setup`.

## Claude Code

### Install

```bash
curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh
midaz setup claude --yes
```

Skills are written to `~/.claude/skills/`. If existing symlinks point to `~/.agents/skills/`, writes are transparently resolved to the symlink target.

### Verify

```bash
midaz doctor
midaz search "test"
```

## Codex

### Install

```bash
curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh
midaz setup codex --yes
```

Skills are written to `~/.codex/skills/`.

### Verify

```bash
midaz version
midaz doctor
```

## Adding a New Target

1. Add the target to `resolveTargets()` in `internal/cmd/setup/setup.go`.
2. Verify `midaz` runs on the target's supported platforms.
3. Run `midaz setup <target> --yes` and verify skills are discovered by the agent.
4. Update this matrix with tested results.
