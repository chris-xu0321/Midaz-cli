// Package skills embeds all skill directories for bundling into the seer-q binary.
package skills

import "embed"

//go:embed all:seer-shared all:seer-market all:seer-api-explorer
var FS embed.FS
