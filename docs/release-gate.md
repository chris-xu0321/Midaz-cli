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

### GitHub Release (when releasing CLI)
- [ ] GitHub Release has binary archives for all 6 platform targets
- [ ] `install.sh` / `install.ps1` fetch and run the new binary successfully

## Manual (per release)

### CLI independence
- [ ] CLI installs and runs without skills being installed
- [ ] `midaz version`, `midaz doctor`, `midaz health` work with no skills
- [ ] No Bash, Python, or Node.js dependency for end users

### Setup command
- [ ] `midaz setup all --yes` installs skills to all targets
- [ ] `midaz setup auto --yes` with no agent dirs returns empty result with hint
- [ ] `midaz setup claude --yes --force` overwrites existing skill files
- [ ] `midaz setup auto --dry-run` works without `--yes`
- [ ] `midaz setup auto` (no `--yes`) fails with confirmation_required error

### Installer scripts
- [ ] `bash install.sh` end-to-end on macOS/Linux
- [ ] `install.ps1` end-to-end on Windows
- [ ] `bash install.sh --agent claude` installs only Claude skills

### Cross-platform (when releasing CLI)
- [ ] Windows: `midaz version` returns correct OS/arch
- [ ] macOS: verify binary runs (if available)
- [ ] Linux: verify binary runs (if available)

### Agent compatibility
- [ ] Claude Code: skills installed via `midaz setup claude --yes`
- [ ] Claude Code: `midaz search` callable from agent context
- [ ] Codex: skills installed via `midaz setup codex --yes`
- [ ] Codex: `midaz search` callable from agent context
