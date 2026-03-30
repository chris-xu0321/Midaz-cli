// Package skills provides embedded skill files for the seer-q CLI.
// Skills are discovered via the */SKILL.md tree pattern.
// Each skill's metadata lives in its own SKILL.md frontmatter.
package skills

import "embed"

//go:embed */SKILL.md
var FS embed.FS
