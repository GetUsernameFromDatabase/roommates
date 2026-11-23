package pages

import (
	"context"
	"io/fs"

	"roommates/internal/locales"
	"roommates/internal/utils"
	"roommates/web/components/select"
)

// languageItems names each language in itself, so it can be found whichever one the page is in,
// built at render time because translations load after package vars are set
func languageItems() []selectcomp.ItemData {
	items := make([]selectcomp.ItemData, 0, len(locales.All))
	for _, lang := range locales.All {
		name := utils.T(locales.WithLocale(context.Background(), lang.String()), locales.LK_LanguageName)
		// locales_test.go checks every language has one, a missing flag only leaves the name
		flag, _ := fs.ReadFile(locales.Content, lang.String()+"/flag.svg")
		items = append(items, selectcomp.ItemData{Value: lang.String(), Label: languageLabel(lang.String(), name, string(flag))})
	}
	return items
}
