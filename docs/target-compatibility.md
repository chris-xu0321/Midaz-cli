# Target Compatibility

Last updated: 2026-03-30

## Support Matrix

| Target | Status | Install Method |
|--------|--------|----------------|
| Claude Code | Supported | `npx skills add chris-xu0321/Midaz-cli -y -g` |
| Codex | Supported | `npx skills add chris-xu0321/Midaz-cli -y -g` |

Midaz keeps skills in the GitHub repo under `skills/`, and agents install them with `npx skills add chris-xu0321/Midaz-cli -y -g`.

## Claude Code

### Install

```bash
# Step 1: CLI
npm install -g @midaz/cli

# Step 2: Skills
npx skills add chris-xu0321/Midaz-cli -y -g
```

### Verify

```bash
seer-q doctor
seer-q search "test"
```

## Codex

### Install

```bash
# Step 1: CLI
npm install -g @midaz/cli

# Step 2: Skills
npx skills add chris-xu0321/Midaz-cli -y -g
```

### Verify

```bash
seer-q version
seer-q doctor
```

## Adding a New Target

1. Verify `seer-q` runs on the target's supported platforms.
2. Confirm the target works with `npx skills add chris-xu0321/Midaz-cli -y -g`.
3. Verify the target discovers and uses the skill content.
4. Update this matrix with tested results.
