# Midaz CLI (`seer-q`)

Query CLI for the [Seer](https://github.com/SparkssL/Seer) market intelligence system. Retrieves structured, evidence-backed market analysis — topics, threads, claims, snapshots, and global regime verdicts — from the Seer API.

## Installation

Packages are hosted on GitHub Packages. Configure npm before installing:

```bash
npm config set @midaz:registry https://npm.pkg.github.com
npm login --registry=https://npm.pkg.github.com  # use GitHub PAT with read:packages scope
npm install -g @midaz/seer-cli
```

Requires Node.js >= 16. Supported platforms: Windows, macOS, Linux (x64, arm64). See [docs/cli-reference.md](docs/cli-reference.md) for full setup details.

## Quick Start

```bash
seer-q search "AI regulation"       # Fuzzy search topics, threads, assets
seer-q market                        # Global regime + all topics
seer-q topic <id>                    # Topic detail + threads
seer-q thread <id>                   # Thread detail + claims + market links
seer-q snapshot                      # Latest global regime snapshot
```

All commands return JSON envelopes. Use `--format pretty` for indented output or `--raw` for the raw API response.

See [docs/cli-reference.md](docs/cli-reference.md) for the full command reference.

## Skills for AI Agents

```bash
npx skills add SparkssL/seer-skills --all -y
```

See [docs/target-compatibility.md](docs/target-compatibility.md) for agent platform setup.

## Development

```bash
make build       # Build seer-q binary
make test        # Run all Go tests
make qa          # Tests + skills-dist-test + smoke-test
make release     # Cross-platform build via goreleaser
make qa-release  # Full QA including npm package verification
```

## License

Private — see repository access controls.
