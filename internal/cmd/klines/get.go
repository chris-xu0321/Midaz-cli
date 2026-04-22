package klines

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
)

func runGet(f *cmdutil.Factory, opts *cmdutil.RunOpts, assetID string) error {
	return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
		Path:      "/api/klines/" + url.PathEscape(assetID),
		Normalize: cmdutil.NormalizePassthrough,
	})
}
