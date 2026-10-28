// Package locales handles i18n
package locales

import (
	"roommates/logger"

	"github.com/invopop/ctxi18n"
	"github.com/invopop/ctxi18n/i18n"
)

func LoadLocales() {
	ctxi18n.DefaultLocale = i18n.Code(Default)
	if err := ctxi18n.Load(Content); err != nil {
		logger.Main.Fatal().Err(err).Msg("error loading i18n content")
	}
}
