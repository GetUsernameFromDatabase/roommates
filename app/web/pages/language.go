package pages

import (
	"roommates/internal/locales"
	"roommates/web/components/select"
)

// each in its own language, so it can be found whichever one the page is in
var languageItems = []selectcomp.ItemData{
	{Value: locales.ET.String(), Label: languageLabel(locales.ET.String(), "Eesti")},
	{Value: locales.EN.String(), Label: languageLabel(locales.EN.String(), "English")},
}
