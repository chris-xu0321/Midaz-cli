# Midaz CLI (`midaz`)

Query CLI for the [Midaz](https://www.midaz.xyz) Interactive Cognitive Trading Map. It retrieves structured, evidence-backed market intelligence from the Midaz API and supports full desk management — auth, onboarding, subscription, radar/playbook, private intel, and more.

## Why `midaz`?

- **Structured market intelligence** — drivers, theses, claims, snapshots, assets, klines, and global regime verdicts
- **Agent-native** — 5 skills bundled in the binary, installed on demand via `midaz skills install`
- **JSON envelope output** — machine-readable responses with `view_url` links, exit codes, and error hints
- **Single binary** — Go, cross-platform, zero runtime dependencies

## Getting started

### Installation

Install the `midaz` binary. The installer places only the binary on your machine — skills and auth are explicit follow-up steps.

macOS / Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.ps1 | iex
```

Once installed, verify the CLI is set up:

```bash
midaz version
```

### Login

Most write endpoints and your desk require an account:

```bash
midaz auth login
```

### Install agent skills (optional)

If you run Claude Code, Codex, or another agent that consumes skill bundles, install them with:

```bash
midaz skills install --yes
```

Targets: `auto` (detected, default), `claude`, `codex`, `all`, or `--skill-dir <path>` for a custom directory. Use `--dry-run` to preview without writing. See [target compatibility](docs/target-compatibility.md) for platform notes.

## Quick Start

### Human users

```bash
midaz search "AI regulation"   # Fuzzy search drivers, theses, assets
midaz market                   # Global regime + drivers + thesis memberships
midaz drivers                  # All active drivers (world-layer)
midaz driver <id>              # Driver detail + thread members + asset contributions
midaz thesis <id>              # Thesis detail + claims + market links
midaz snapshot                 # Latest global regime snapshot
midaz auth login               # Sign in
midaz onboard                  # Complete desk onboarding
midaz desk                     # Manage radar, playbook, preferences, sharing, Telegram
```

All commands return JSON envelopes. Use `--format pretty` for indented output or `--raw` for the raw API response.

### AI agents

Install the binary with the one-liner above, then install skills:

```bash
midaz skills install --yes
```

Inside Claude Code or Codex, the skills self-register under `~/.claude/skills` or `~/.codex/skills` and guide the agent through Midaz commands.

## Skills

| Skill | Description |
|-------|-------------|
| `midaz-shared` | Shared concepts — auth model, response format, global flags, safety rules |
| `midaz-market` | Search, browse, and analyze drivers, theses, claims, assets, klines, deltas, and regime |
| `midaz-account` | Authenticate, redeem invitations, complete onboarding, and manage subscription |
| `midaz-desk` | Manage radar, playbook, sharing, Telegram alerts, private intel, asset tracking |
| `midaz-onboard` | Interactive 5-round trader onboarding ritual — produces profile + radar + playbook, syncs to the desk |
| `midaz-api-explorer` | Discover commands via schema introspection — fallback when other skills don't fit |

## Development

```bash
make build       # Build midaz binary
make test        # Run all Go tests
make qa          # Tests + skills validation + smoke test
make release     # Cross-platform build via goreleaser
make install     # Install to /usr/local/bin (or PREFIX)
```

## Links

- [CLI Reference](docs/cli-reference.md) — full command documentation
- [Target Compatibility](docs/target-compatibility.md) — agent platform setup
- [Release Gate](docs/release-gate.md) — QA checklist
- [Changelog](CHANGELOG.md)

## License

[MIT](LICENSE)
