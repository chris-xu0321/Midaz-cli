// Package skills embeds all skill directories for bundling into the midaz binary.
package skills

import "embed"

//go:embed all:midaz-shared all:midaz-market all:midaz-api-explorer all:midaz-account all:midaz-workspace
var FS embed.FS
