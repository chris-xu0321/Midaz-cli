# Midaz CLI (`seer-q`)

Query CLI for the [Seer](https://github.com/SparkssL/Seer) market intelligence system. It retrieves structured, evidence-backed market analysis from the Seer API.

## Why seer-q?

- **Structured market intelligence** - topics, threads, claims, snapshots, and global regime verdicts
- **Agent-native** - 3 skills installed directly from this repo via `npx skills add chris-xu0321/Midaz-cli -y -g`
- **JSON envelope output** - machine-readable responses with `view_url` links, exit codes, and error hints
- **Single binary** - Go, cross-platform, zero runtime dependencies

## Installation

### From npm (recommended)

```bash
npm install -g @midaz/cli
```

### From source

```bash
git clone https://github.com/chris-xu0321/Midaz-cli.git
cd Midaz-cli
make install
```

## Quick Start

### Human Users

```bash
seer-q search "AI regulation"       # Fuzzy search topics, threads, assets
seer-q market                       # Global regime + all topics
seer-q topic <id>                   # Topic detail + threads
seer-q thread <id>                  # Thread detail + claims + market links
seer-q snapshot                     # Latest global regime snapshot
```

All commands return JSON envelopes. Use `--format pretty` for indented output or `--raw` for the raw API response.

### AI Agents

Install all skills:

```bash
npm install -g @midaz/cli
npx skills add chris-xu0321/Midaz-cli -y -g
```

Skills provide structured guidance for querying the Seer API. See [target compatibility](docs/target-compatibility.md) for platform-specific notes.

## Skills

| Skill | Description |
|-------|-------------|
| `seer-shared` | Response format, config, and common rules |
| `seer-market` | Search, browse, and analyze topics, threads, and claims |
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

- [CLI Reference](docs/cli-reference.md) - full command documentation
- [Target Compatibility](docs/target-compatibility.md) - agent platform setup
- [Release Gate](docs/release-gate.md) - QA checklist
- [Changelog](CHANGELOG.md)

## License

[MIT](LICENSE)
