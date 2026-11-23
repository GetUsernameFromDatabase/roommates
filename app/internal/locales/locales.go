// Package locales handles i18n
package locales

import (
	"context"
	"roommates/internal/logger"

	"github.com/invopop/ctxi18n"
	"github.com/invopop/ctxi18n/i18n"
)

func LoadLocales() {
	ctxi18n.DefaultLocale = i18n.Code(Default)
	if err := ctxi18n.Load(Content); err != nil {
		logger.Main.Fatal().Err(err).Msg("error loading i18n content")
	}
	if ctxi18n.Get(i18n.Code(Default)) == nil {
		logger.Main.Fatal().Str("locale", Default.String()).Msg("default locale is not loaded")
	}
}

// WithLocale sets lang on ctx, an unknown lang falls back to [Default]
func WithLocale(ctx context.Context, lang string) context.Context {
	lctx, err := ctxi18n.WithLocale(ctx, lang)
	if err != nil {
		return ctx // only a missing default locale errors, which LoadLocales refuses at startup
	}
	return lctx
}
