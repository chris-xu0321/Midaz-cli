# Release Gate Checklist

Last updated: 2026-04-02

## Automated (must all pass)

### Go tests
- [ ] `go test ./...` - all tests pass

### Smoke test
- [ ] `bash test/smoke-test.sh` - all commands return `ok:true`
- [ ] API reachable at `https://www.midaz.xyz` (production) or `localhost:4000` (local)

### Skills distribution
- [ ] `bash test/skills-dist-test.sh` - artifact is complete
- [ ] Skills present, frontmatter valid, no leaked files (embed.go excluded)

### npm package (when releasing CLI)
- [ ] `npm pack --dry-run` - package contains `scripts/install.js`, `scripts/run.js`, and `CHANGELOG.md`
- [ ] GitHub Release has binary archives for all 6 platform targets
- [ ] `npm install -g @midaz/cli` downloads and runs binary successfully

## Manual (per release)

### CLI independence
- [ ] CLI installs and runs without skills being installed
- [ ] `seer-q version`, `seer-q doctor`, `seer-q health` work with no skills
- [ ] No Bash, Python, or Node.js dependency for end users

### Setup command
- [ ] `seer-q setup all --yes` installs skills to all targets
- [ ] `seer-q setup auto --yes` with no agent dirs returns empty result with hint
- [ ] `seer-q setup claude --yes --force` overwrites existing skill files
- [ ] `seer-q setup auto --dry-run` works without `--yes`
- [ ] `seer-q setup auto` (no `--yes`) fails with confirmation_required error

### Installer scripts
- [ ] `bash install.sh` end-to-end on macOS/Linux
- [ ] `install.ps1` end-to-end on Windows
- [ ] `bash install.sh --agent claude` installs only Claude skills

### Skill independence (legacy)
- [ ] `npx skills add SparkssL/Midaz-cli -y -g` still works

### Cross-platform (when releasing CLI)
- [ ] Windows: `seer-q version` returns correct OS/arch
- [ ] macOS: verify binary runs (if available)
- [ ] Linux: verify binary runs (if available)

### Agent compatibility
- [ ] Claude Code: skills installed via `seer-q setup claude --yes`
- [ ] Claude Code: `seer-q search` callable from agent context
- [ ] Codex: skills installed via `seer-q setup codex --yes`
- [ ] Codex: `seer-q search` callable from agent context
