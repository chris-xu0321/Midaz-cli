// Package agent provides embedded agent resource files for the seer-q CLI.
// These are the source-of-truth skill and command files that get installed
// into a workspace's .claude/ directory via "seer-q agent install".
package agent

import _ "embed"

//go:embed cmd/seer.md
var SeerCmdMD []byte

//go:embed skills/seer-market.md
var SeerSkillMD []byte
