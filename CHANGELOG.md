# Changelog

## 0.6.1 — 2026-04-16

### Added
- **Radar pins**: `midaz workspace radar {pin,unpin,pins}` wrap the new entity-pin API (`POST`/`DELETE /api/ws/radar/pin`, `GET /api/ws/radar/pins`). Distinct from free-form radar lines — pins carry provenance and participate in L4 refresh.

### Fixed
- `midaz snapshot` now hits `/api/global` (the legacy `/api/global/snapshot` + `/api/global/snapshots` routes were removed upstream). The `--history` / `--limit` flags are dropped — the new API has no history variant.

### Changed
- 402 responses now include a hint pointing at `midaz subscription status` / `midaz subscription start` instead of a bare error.
- `client.Delete()` now accepts a JSON body (needed by radar `unpin`).

## 0.6.0 — 2026-04-15

### Added
- **Full rebrand**: binary renamed `seer-q` → `midaz`; skills `seer-*` → `midaz-*`; config dir `~/.config/seer/` → `~/.config/midaz/`; env vars `SEER_*` → `MIDAZ_*`. Legacy names still work for one release (deprecation notice).
- **Auth subsystem**: `midaz auth {login,logout,status,whoami,keys}`. Browser-relay login via `/cli-auth` (requires Seer frontend companion page); `--paste` fallback for headless; `--token` / `MIDAZ_TOKEN` for CI. Credentials in `~/.config/midaz/auth.json` (mode 0600), profile-shaped for future multi-profile.
- **Workspace**: `midaz workspace {get,settings,view,share,radar,playbook,telegram}` — mirrors `/workspace/settings` in the web app.
- **Onboarding**: `midaz onboard {status,generate,complete}` covering both AI-generated and direct paths.
- **Invitations**: `midaz invite redeem <code> --yes`.
- **Subscription**: `midaz subscription {status,start,portal}` — Stripe Checkout and Customer Portal URLs auto-open.
- **Intel**: `midaz intel {list,push,rm}` for private notes / sources (subscription-gated).
- **Assets**: `midaz assets {list,get,thesis}` for per-asset thesis links + evidence.
- **Delta**: `midaz delta [--hours N]` surfaces new claims + theses + topics.
- **Thesis rename**: `midaz theses`, `midaz thesis <id>` hit `/api/theses*`. `threads`/`thread` still work as hidden deprecated aliases.
- **New exit codes**: 6 (auth required / invalid), 7 (subscription required / workspace paused).
- **New skills**: `midaz-account` (auth/onboarding/invite/subscription), `midaz-workspace` (settings/radar/playbook/telegram/intel/assets). Existing skills rewritten and rebranded.

### Changed
- HTTP client now supports POST/PATCH/DELETE with JSON bodies and auto-attaches `Authorization: Bearer` when the user is logged in.
- `setup` now installs the 5 new `midaz-*` skills and cleans up any stale `seer-*` skill directories it finds.
- `doctor` adds an `auth` check that prints masked credentials.

## 0.5.0 — 2026-04-02

### Added
- `seer-q setup` command — install embedded skills to agent directories (claude, codex, auto, all)
- `install.sh` — zero-dependency installer for macOS/Linux (downloads binary + installs skills)
- `install.ps1` — zero-dependency installer for Windows
- Skills embedded in binary via `go:embed` (`skills/embed.go`)
- `--yes`, `--force`, `--dry-run`, `--skill-dir` flags on setup command

### Changed
- Primary install method is now `curl | sh` / `irm | iex` (no npm required)
- npm install remains supported as secondary option
- Updated all documentation to reflect new install flow
- `test/skills-dist-test.sh` updated to allow `embed.go` in skills/

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
- `seer-q agent install/uninstall/doctor` bridge commands — use `npx skills add SparkssL/Midaz-cli --all -y` instead
- Embedded skills/agent Go packages (`agent/`, `skills/embed.go`)
- `npm/skills-repo-template/` (skills distributed from this repo directly)

### Changed
- `/seer` command wrapper moved from `agent/cmd/seer.md` to `commands/claude/seer.md`
- Skills install references updated from `SparkssL/seer-skills` to `SparkssL/Midaz-cli`
- Consistency tests now read skills from disk instead of embedded FS

## 0.2.0 — 2026-03-30

Initial release from standalone repository. Bootstrapped from [SparkssL/Seer](https://github.com/SparkssL/Seer) `apps/cli/` subtree.

### Features

- 16 query commands: `search`, `market`, `topics`, `topic`, `threads`, `thread`, `claims`, `sources`, `snapshot`, `usage`, `decisions`, `health`, `version`, `doctor`, `config`, `schema`
- Agent compatibility bridge (`agent install/uninstall/doctor`) — deprecated, use skill installer
- Cross-platform npm distribution: meta package + 6 platform packages (`@midaz/cli-*`)
- 3 embedded skills: `seer-shared`, `seer-market`, `seer-api-explorer`
- goreleaser-based multi-platform builds (darwin/linux/windows, amd64/arm64)
- 23 golden JSON contract tests
- Test infrastructure: smoke test, skills distribution test, npm install test
- JSON envelope response format with exit codes and error hints
