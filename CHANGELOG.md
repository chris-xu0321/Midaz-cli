# Changelog

## Unreleased

### Added
- New `midaz-onboard` skill — interactive trader onboarding ritual synced
  from Seer; produces profile, radar, and playbook and commits them via
  the existing `midaz onboard complete` / `midaz desk radar set` /
  `midaz desk playbook set` verbs.

### Fixed
- `midaz-account` skill: radar cap corrected from `≤5` to `≤12` items to
  match `internal/cmd/desk/radar/add.go`.

### Breaking
- Removed the npm distribution channel. `npm install -g @midaz/cli` is no
  longer published; the `@midaz/cli` package, its postinstall downloader,
  `scripts/run.js`, `scripts/install.js`, `npm/publish.sh`, and
  `package.json` are all gone. Install via the curl (`install.sh`) or
  PowerShell (`install.ps1`) one-liners instead. Version is now sourced
  from the git tag only (goreleaser resolves it in CI).
- Removed deprecated `seer-q` npm shim, `SEER_*` env var fallbacks, and the
  `~/.config/seer/` legacy config read. The grace window from v0.6.0 is over —
  use `midaz` and `MIDAZ_*` env vars. Existing users with config under
  `~/.config/seer/config.json` must move it to `~/.config/midaz/config.json`
  (or re-run `midaz config set`).

## 0.7.2 — 2026-04-22

Companion release for Seer's post-L4 drift: the `/api/topics*` routes were
removed in favor of a canonical driver ontology, and several new endpoints
landed (klines, asset timelines, desk preferences, full pipeline refresh,
per-run usage). This release closes that gap.

### Breaking
- **Removed commands** whose backing endpoints no longer exist on Seer:
  - `midaz topics`, `midaz topic <id>` — `/api/topics*` was dropped in
    Seer `a5c4c59` ("drop Topic legacy"). Use `midaz drivers` /
    `midaz driver <id>` instead.
  - `midaz assets thesis <asset> <thesis>` — `/api/assets/:id/theses/:thesis_id`
    was removed. Use `midaz assets timeline <asset>` for per-asset event
    history, or `midaz thesis <id>` for the thesis itself.
- **`midaz desk radar add --topic` removed**: replaced by `--driver`. Radar
  lines now use the `driver:` prefix (was `topic:`).

### Added
- **Driver ontology commands** (canonical replacement for topics):
  - `midaz drivers` → `GET /api/drivers`
  - `midaz driver <id>` → `GET /api/drivers/:id`
  - `midaz driver-links` → `GET /api/driver-links`
- **Klines** (price history): `midaz klines` and `midaz klines <asset_id>`
  → `GET /api/klines`, `GET /api/klines/:assetId`.
- **Assets timeline**: `midaz assets timeline <asset_id> [--limit N]` →
  `GET /api/assets/:id/timeline`. Replaces the removed `assets thesis`.
- **Desk preferences**: `midaz desk preferences get` /
  `midaz desk preferences set --language <code> --yes` →
  `PATCH /api/desk/preferences`. Supported languages: `en`, `zh-CN`, `ja`,
  `ko`, `es`, `fr`.
- **Desk full refresh**: `midaz desk refresh --yes` →
  `POST /api/desk/refresh`. Triggers a full pipeline refresh (market
  rebuild + personal desk), distinct from `desk regenerate` which only
  replays L4.
- **Usage by-run**: `midaz usage by-run <run_id>` →
  `GET /api/usage/by-run/:runId` for per-pipeline-run token-cost breakdown.

### Changed
- `midaz desk radar add --driver` resolves labels against `/api/drivers/:id`
  using the `name` field (was `/api/topics/:id` with `name`).
- Skill metadata bumped to `0.7.2`; `midaz-market` and `midaz-desk`
  rewritten to document drivers, klines, asset timeline, preferences,
  and the refresh/regenerate/reonboard trio.

## 0.7.1 — 2026-04-16

### Added
- **Manual desk refresh verbs**: `midaz desk regenerate` (L4Cause.manual, `POST /api/desk/personal-desk/regenerate`) and `midaz desk reonboard` (L4Cause.personal_input, `POST /api/desk/onboard` with current radar+playbook re-submitted). Both owner-only; both require an active subscription (regenerate directly, reonboard transitively via the settings round-trip). Mirrors the "Regenerate personal desk" and "Run setup" buttons in the web Desk Preferences panel.

### Changed
- Trimmed server-internal jargon out of user-facing `--help` text across `auth status`, `invite redeem`, `onboard`, and `desk radar pin/unpin` (stale `workspace` → `desk`, endpoint paths, DB terminology).

## 0.7.0 — 2026-04-16

Companion release for Seer's `workspaces → desks` rename. Backend dropped its
`/api/ws*` and `/api/workspaces/:slug/view` mounts with no compat layer, so
this release is required for any CLI built against earlier Seer.

### Breaking
- **Command rename**: `midaz workspace` → `midaz desk` (and all subcommands:
  `desk get/settings/view/share/radar/playbook/telegram`). The `workspace`
  command is removed — scripts must update.
- **Skill rename**: `midaz-workspace` → `midaz-desk`. The embedded skill
  directory and on-disk install path both move; reinstall via
  `midaz skills install --yes`.
- **API path swap**: every endpoint hit by the CLI under `/api/ws*` now
  targets `/api/desk*`. Affects `desk` subtree, `onboard`, and
  `subscription status`.
- **`auth.json` field rename**: `workspace_id`/`workspace_slug` →
  `desk_id`/`desk_slug`. Existing files are read transparently via a
  one-time migration; on next `auth login` (or any save) the file is
  rewritten with the new keys. `midaz auth whoami` output also switches
  to `desk_id`/`desk_slug`.
- **`claims --thread` flag removed**: backend dropped the `?thread_id`
  query parameter (Seer commit `cc2c6ad`); the CLI flag had been a
  silent no-op. Use `midaz thesis <id>` to read the embedded `claims[]`
  array on a thesis.

### Changed
- `desk view` short help now says "Personal market read" (mirrors
  Seer's `market_view → market_read` rename).
- Frontend URL hint in `midaz auth login --paste` points at
  `/desk/settings` (the web SPA route was renamed in lockstep).
- `/api/app/me` parser accepts both `desk` and `workspace` top-level
  fields during the rename window; prefers `desk`.
- Loopback CLI auth callback accepts both `desk_*` and `workspace_*`
  payload keys (the `/cli-auth` page dual-emits during transition).

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
