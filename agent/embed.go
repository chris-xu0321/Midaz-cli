// Package agent provides embedded Claude-specific command files for the seer-q CLI.
// Skills have moved to the skills package (apps/cli/skills/).
// This package only embeds the Claude command wrapper (cmd/seer.md).
package agent

import _ "embed"

//go:embed cmd/seer.md
var SeerCmdMD []byte
