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
