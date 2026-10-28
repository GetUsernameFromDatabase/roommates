package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"roommates/internal/locales"
	"strings"

	"github.com/invopop/ctxi18n/i18n"
)

// T is a wrapper of [i18n.T]
// which forces key type
func T(ctx context.Context, key locales.LK, args ...any) string {
	s := unescapeJSON(i18n.T(ctx, key.String()))
	if len(args) == 0 {
		return s
	}
	return fmt.Sprintf(s, args...)
}

// N is a wrapper of [i18n.N]
// which forces key type
//
// i18nMap ([i18n.M]) fills %{name} placeholders when there are no args,
// otherwise args fill the fmt verbs
func N(ctx context.Context, key locales.LK, count int, i18nMap i18n.M, args ...any) string {
	s := unescapeJSON(i18n.N(ctx, key.String(), count))
	if len(args) == 0 {
		return i18nMap.Replace(s)
	}
	return fmt.Sprintf(s, args...)
}

// unescapeJSON undoes the escaping ctxi18n leaves in texts: it loads YAML
// through JSON and keeps the raw JSON string, so a newline comes out as `\n`
// and "&" as `&` (i18n.Dict.UnmarshalJSON in ctxi18n v0.9.0)
//
// it runs before args are formatted in, so user input is never unescaped
//
// shortcut: remove once ctxi18n decodes its strings
func unescapeJSON(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var out string
	if err := json.Unmarshal([]byte(`"`+s+`"`), &out); err != nil {
		return s
	}
	return out
}

// GetLocaleString gets locale using [i18n.GetLocale],
// turns resulting [i18n.Code] into string
func GetLocaleString(ctx context.Context) string {
	return i18n.GetLocale(ctx).Code().String()
}
