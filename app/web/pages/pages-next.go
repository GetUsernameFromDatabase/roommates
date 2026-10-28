package pages

import (
	"net/url"
	g "roommates/internal/config"

	"github.com/a-h/templ"
)

// withNext keeps ?next= when switching between the login and register pages
func withNext(route, next string) templ.SafeURL {
	if next == "" {
		return templ.SafeURL(route)
	}
	return templ.SafeURL(route + "?" + url.Values{g.QueryNext: {next}}.Encode())
}
