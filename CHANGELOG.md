# Changelog

## 0.4.0 — 2026-03-30

### Changed
- Switched to lark-style npm distribution: single package with `postinstall` binary download from GitHub Releases (replaces 7-package platform model)
- Root `package.json` + `scripts/install.js` + `scripts/run.js` (like lark-cli)
- README rewritten with lark-style installer-facing sections (Why, Install, Quick Start for humans + AI agents)
- `test/skills-dist-test.sh` validates skills directly (no longer depends on publish-skills.sh)

### Removed
- `npm/build.sh`, `npm/platform-template/`, `npm/verify.sh` (platform-package build pipeline)
- `npm/publish-skills.sh` (skills install directly from repo)
- `commands/claude/seer.md` (skills handle agent integration)
- `test/npm-install-test.sh` (platform-package test)

### Added
- `make install` target with PREFIX support
- `.github/workflows/lint.yml` (golangci-lint)

## 0.3.0 — 2026-03-30

### Removed
- `seer-q agent install/uninstall/doctor` bridge commands — use `npx skills add chris-xu0321/Midaz-cli --all -y` instead
- Embedded skills/agent Go packages (`agent/`, `skills/embed.go`)
- `npm/skills-repo-template/` (skills distributed from this repo directly)

### Changed
- `/seer` command wrapper moved from `agent/cmd/seer.md` to `commands/claude/seer.md`
- Skills install references updated from `SparkssL/seer-skills` to `chris-xu0321/Midaz-cli`
- Consistency tests now read skills from disk instead of embedded FS

## 0.2.0 — 2026-03-30

Initial release from standalone repository. Bootstrapped from [SparkssL/Seer](https://github.com/SparkssL/Seer) `apps/cli/` subtree.

### Features

- 16 query commands: `search`, `market`, `topics`, `topic`, `threads`, `thread`, `claims`, `sources`, `snapshot`, `usage`, `decisions`, `health`, `version`, `doctor`, `config`, `schema`
- Agent compatibility bridge (`agent install/uninstall/doctor`) — deprecated, use skill installer
- Cross-platform npm distribution: meta package + 6 platform packages (`@midaz/seer-cli-*`)
- 3 embedded skills: `seer-shared`, `seer-market`, `seer-api-explorer`
- goreleaser-based multi-platform builds (darwin/linux/windows, amd64/arm64)
- 23 golden JSON contract tests
- Test infrastructure: smoke test, skills distribution test, npm install test
- JSON envelope response format with exit codes and error hints
