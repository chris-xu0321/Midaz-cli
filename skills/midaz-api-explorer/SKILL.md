---
name: midaz-api-explorer
version: 0.7.0
description: Discover and explore Midaz CLI commands via schema introspection — use when existing skills don't cover the user's need
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz API Explorer

> Read [midaz-shared](../midaz-shared/SKILL.md) for response format and common rules.

Fallback skill. Use when the user's need isn't obviously covered by `midaz-market`, `midaz-account`, or `midaz-desk`.

## When to Use

- User asks about a command you don't recognize.
- User needs data that the domain skills don't document.
- You want to discover what commands are available.
- You need the exact input/output contract for a command.

## Discovery Flow

### 1. Check domain skills first

Before running `schema`, ask: does `midaz-market`, `midaz-account`, or `midaz-desk` already cover this? If yes, use that skill.

### 2. List every command

```
midaz schema
```

Returns every registered command with description, positional args, flags, and the endpoints it hits. Use this index to pick the right command.

### 3. Inspect one command

```
midaz schema <command>
```

Returns the full contract for that command: arguments, flags, and endpoint mapping.

### 4. Call with --raw

```
midaz <command> [args] [flags] --raw
```

Bypasses the envelope and prints the raw API response. Handy for:

- Inspecting full response structure when the envelope normalizer hides fields.
- Debugging unexpected output.
- Exploring fields not documented in skills.

### 5. Interpret & synthesize

1. Parse the JSON (`.data` payload when not using `--raw`).
2. Identify the relevant fields for the user's question.
3. Synthesize into natural language.
4. Surface any `view_url` as a markdown link.

## Auth & Subscription Gotchas

- Exit code 6 → the command needs auth. Run `midaz auth status`; if it also fails, `midaz auth login`.
- Exit code 7 → the command needs an allowed subscription. Run `midaz subscription status`; surface the trial/active state to the user before running `midaz subscription start --yes`.
- Exit code 2 with `confirmation_required` → the command is a write; add `--yes`.

## Example

User: "what stages does the pipeline use for processing?"

1. Not covered by the market skill → use API explorer.
2. `midaz schema decisions` → see `--stage` flag.
3. `midaz decisions --limit 5` → real stage names in output.
4. Summarize the processing stages found.

User: "is there a way to list all my PATs?"

1. `midaz schema` → find `auth` command.
2. `midaz schema auth` → surface `keys {list,create,revoke}`.
3. `midaz auth keys list` → summarize active keys with labels and last-used dates.
