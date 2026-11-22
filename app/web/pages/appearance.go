package pages

import (
	"roommates/internal/locales"

	"github.com/a-h/templ"
)

var themeNames = map[string]locales.LK{
	"home":  locales.LK_AppearanceThemeHome,
	"plain": locales.LK_AppearanceThemePlain,
}

var modeNames = map[string]locales.LK{
	"system": locales.LK_AppearanceModeSystem,
	"light":  locales.LK_AppearanceModeLight,
	"dark":   locales.LK_AppearanceModeDark,
}

var appearancePreviewScript = templ.NewOnceHandle()
