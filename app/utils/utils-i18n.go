package utils

import (
	"context"
	"roommates/locales"

	"github.com/invopop/ctxi18n/i18n"
)

// wrapper of i18n.T
//   - if default is empty ("") then will use key as the default value
func T(ctx context.Context, key locales.LK, defaultValue string, args ...any) string {
	if defaultValue == "" {
		defaultValue = string(key)
	}
	defaultArg := i18n.Default(defaultValue)
	return i18n.T(ctx, string(key), append(args, defaultArg)...)
}

// wrapper of i18n.N
//
// i18nMap i18n.M is used since without it %!(EXTRA ...) is added on plural rule
// which does not need arguments like %s or %d
func N(ctx context.Context, key locales.LK, count int, i18nMap i18n.M, args ...any) string {
	return i18n.N(ctx, string(key), count, append(args, i18nMap)...)
}
