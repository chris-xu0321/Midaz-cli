# Midaz CLI (`seer-q`)

Query CLI for the [Seer](https://github.com/SparkssL/Seer) market intelligence system. Retrieves structured, evidence-backed market analysis from the Seer API.

## Why seer-q?

- **Structured market intelligence** — topics, threads, claims, snapshots, global regime verdicts
- **Agent-native** — 3 skills for Claude Code and other AI agents via `npx skills add`
- **JSON envelope output** — machine-readable with `view_url` links, exit codes, error hints
- **Single binary** — Go, cross-platform, zero runtime dependencies

## Installation

### From npm (recommended)

```bash
npm config set @midaz:registry https://npm.pkg.github.com
npm login --registry=https://npm.pkg.github.com  # GitHub PAT with read:packages scope
GITHUB_TOKEN=<your-pat> npm install -g @midaz/seer-cli
```

`GITHUB_TOKEN` is needed to download the binary from GitHub Releases during install.

### From source

```bash
git clone https://github.com/chris-xu0321/Midaz-cli.git
cd Midaz-cli && make install
```

## Quick Start

### Human Users

```bash
seer-q search "AI regulation"       # Fuzzy search topics, threads, assets
seer-q market                        # Global regime + all topics
seer-q topic <id>                    # Topic detail + threads
seer-q thread <id>                   # Thread detail + claims + market links
seer-q snapshot                      # Latest global regime snapshot
```

All commands return JSON envelopes. Use `--format pretty` for indented output or `--raw` for the raw API response.

### AI Agents

```bash
npx skills add chris-xu0321/Midaz-cli --all -y
```

Skills provide structured guidance for querying the Seer API. See [target compatibility](docs/target-compatibility.md) for supported platforms.

## Skills

| Skill | Description |
|-------|-------------|
| `seer-shared` | Response format, config, common rules |
| `seer-market` | Search, browse, and analyze topics, threads, claims |
| `seer-api-explorer` | Discover commands via schema introspection |

## Development

```bash
make build       # Build seer-q binary
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

Private — see repository access controls.
