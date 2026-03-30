# Seer Skills

Private skill source for the Seer market intelligence CLI.

Requires the `seer-q` CLI installed first (`npm install -g @midaz/seer-cli`).

## Available Skills

| Skill | Description |
|-------|-------------|
| seer-shared | Shared concepts, response format, config, and common rules |
| seer-market | Search, browse, and analyze topics, threads, claims |
| seer-api-explorer | Discover and explore API commands via schema introspection |

## Installation

### Claude Code

```bash
npx skills add SparkssL/seer-skills --all -y
```

Requires GitHub access to the private skills repository.

### Other Ecosystems

See [target-compatibility.md](https://github.com/SparkssL/seer-cli/blob/main/docs/target-compatibility.md) in the CLI repository.

## Updating

```bash
npx skills add SparkssL/seer-skills --all -y
```
