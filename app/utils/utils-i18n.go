package utils

import (
	"context"
	"roommates/locales"

	"github.com/invopop/ctxi18n/i18n"
)

// T is a wrapper of [i18n.T]
// which forces key type
func T(ctx context.Context, key locales.LK, args ...any) string {
	return i18n.T(ctx, key.String(), args...)
}

// N is a wrapper of [i18n.N]
// which forces key type
//
// i18nMap ([i18n.M]) is used since without it
//
//	"%!(EXTRA ...)"
//
// is added on plural rule
// which does not need arguments like %s or %d
func N(ctx context.Context, key locales.LK, count int, i18nMap i18n.M, args ...any) string {
	return i18n.N(ctx, key.String(), count, append(args, i18nMap)...)
}

// GetLocaleString gets locale using [i18n.GetLocale],
// turns resulting [i18n.Code] into string
func GetLocaleString(ctx context.Context) string {
	return i18n.GetLocale(ctx).Code().String()
}
