// Package assets [embed.FS]
package assets

import (
	"embed"
)

//go:embed all:public favicon/favicon.ico
var PublicAssetsFS embed.FS
