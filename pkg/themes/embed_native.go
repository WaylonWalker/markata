//go:build !js

package themes

import "embed"

//go:embed all:default
var defaultTheme embed.FS
