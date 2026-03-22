package templates

import "embed"

// FS stores all text templates used by project-generator.
//
//go:embed *.tmpl
var FS embed.FS
