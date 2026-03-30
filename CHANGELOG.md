# Changelog

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
