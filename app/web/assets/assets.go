// Package assets [embed.FS]
package assets

import (
	"embed"
)

//go:embed all:public favicon/favicon.ico favicon/*.png pwa
var PublicAssetsFS embed.FS
