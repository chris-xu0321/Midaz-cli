# Target Compatibility

Last updated: 2026-03-30

## Support Matrix

| Target | Status | Install Method |
|--------|--------|----------------|
| Claude Code | Supported | `npx skills add SparkssL/seer-skills --all -y` |
| Codex | Planned | Blocked on Codex environment availability for testing |

## Claude Code

### Install

```bash
# Step 1: CLI
npm install -g @midaz/seer-cli

# Step 2: Skills
npx skills add SparkssL/seer-skills --all -y
```

Requires GitHub access to the private skills repository.

Skills are discovered from `.claude/skills/*/SKILL.md`. YAML frontmatter provides metadata.

**Fallback** (if skill installer does not support private repos):

```bash
git clone git@github.com:SparkssL/seer-skills.git /tmp/seer-skills
cp -r /tmp/seer-skills/skills/* .claude/skills/
```

### Optional: Claude command wrapper

The `/seer` slash command shortcut is available in the skills repo:

```bash
cp /tmp/seer-skills/commands/claude/seer.md .claude/cmd/seer.md
```

### Verify

```bash
seer-q doctor
seer-q search "test"
```

## Codex (Planned)

The Seer skill tree (`skills/*/SKILL.md`) uses standard YAML frontmatter and markdown content. Codex compatibility is planned but blocked on:

- No Codex environment available for testing
- Codex skill directory convention not yet verified
- `seer-q` binary availability on Codex PATH not tested

The skill format is architecturally compatible. End-to-end testing is required before marking as supported.

## Adding a New Target

1. Verify `seer-q` runs on the target's supported platforms
2. Determine the target's skill directory convention
3. Install skills from the Seer skill source per that convention
4. Verify the target discovers and uses the skill content
5. Update this matrix with tested results
