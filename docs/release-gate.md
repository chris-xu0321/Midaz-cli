# Release Gate Checklist

Last updated: 2026-03-30

## Automated (must all pass)

### Go tests
- [ ] `go test ./...` - all 35+ tests pass

### Smoke test
- [ ] `bash test/smoke-test.sh` - all commands return `ok:true`
- [ ] API reachable at `https://www.midaz.xyz` (production) or `localhost:4000` (local)

### Skills distribution
- [ ] `bash test/skills-dist-test.sh` - artifact is complete
- [ ] Skills present, frontmatter valid, no leaked files

### npm package (when releasing CLI)
- [ ] `npm pack --dry-run` - package contains `scripts/install.js`, `scripts/run.js`, and `CHANGELOG.md`
- [ ] GitHub Release has binary archives for all 6 platform targets
- [ ] `npm install -g @midaz/cli` downloads and runs binary successfully

## Manual (per release)

### CLI independence
- [ ] CLI installs and runs without skills being installed
- [ ] `seer-q version`, `seer-q doctor`, `seer-q health` work with no skills
- [ ] No Bash or Python dependency for end users

### Skill independence
- [ ] `npx skills add chris-xu0321/Midaz-cli -y -g` installs all skills
- [ ] Skills work when installed via `npx skills add chris-xu0321/Midaz-cli -y -g`

### Cross-platform (when releasing CLI)
- [ ] Windows: `seer-q version` returns correct OS/arch
- [ ] macOS: verify binary runs (if available)
- [ ] Linux: verify binary runs (if available)

### Agent compatibility
- [ ] Claude Code: skills installed via `npx skills add chris-xu0321/Midaz-cli -y -g`
- [ ] Claude Code: `seer-q search` callable from agent context
- [ ] Codex: skills installed via `npx skills add chris-xu0321/Midaz-cli -y -g`
- [ ] Codex: `seer-q search` callable from agent context
