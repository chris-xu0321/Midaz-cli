---
name: seer-api-explorer
version: 1.0.0
description: Discover seer-q commands you don't know about via schema introspection
metadata: {"requires":{"bins":["seer-q"]}}
---

# Seer API Explorer

> Use this only when `seer-market` doesn't cover what you need.

```bash
seer-q schema              # List all commands
seer-q schema <command>    # Inspect a command's args, flags, endpoints
seer-q <command> --raw     # Get raw API response (no envelope)
```

This is your fallback for discovering commands not documented in the market skill.
