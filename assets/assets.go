package assets

import "embed"

// ShimTemplates is a bunch of shim templates that can be used to create shims.
//go:embed shims
var ShimTemplates embed.FS
