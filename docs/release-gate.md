# Release Gate Checklist

Last updated: 2026-03-30

## Automated (must all pass)

### Go tests
- [ ] `go test ./...` — all 35+ tests pass

### Smoke test
- [ ] `bash test/smoke-test.sh` — all commands return ok:true
  - Requires: API running at localhost:4000

### Skills distribution
- [ ] `bash test/skills-dist-test.sh` — artifact is complete
  - Skills present, frontmatter valid, no leaked files

### npm packages (when releasing CLI)
- [ ] `bash npm/verify.sh` — all 7 packages valid
  - Requires: goreleaser build completed
- [ ] `bash test/npm-install-test.sh` — install + bin shim works

## Manual (per release)

### CLI independence
- [ ] CLI installs and runs without skills being installed
- [ ] `seer-q version`, `seer-q doctor`, `seer-q health` work with no skills
- [ ] No Bash or Python dependency for end users

### Skill independence
- [ ] Skills install without requiring CLI source repo access
- [ ] Skills work when installed from the private skills repo
- [ ] New skills can be added to the skills repo without CLI rebuild

### Cross-platform (when releasing CLI)
- [ ] Windows: `seer-q version` returns correct OS/arch
- [ ] macOS: verify binary runs (if available)
- [ ] Linux: verify binary runs (if available)

### Agent compatibility
- [ ] Claude Code: skills discovered from `.claude/skills/*/SKILL.md`
- [ ] Claude Code: `seer-q search` callable from agent context
- [ ] Codex: PASS or DEFER (if deferred: document rationale in release notes, track in target-compatibility.md)

### Bridge status
- [ ] `seer-q agent install claude` still works (compatibility bridge, deprecated)
- [ ] Bridge removal criteria tracked in `docs/bridge-removal.md`
