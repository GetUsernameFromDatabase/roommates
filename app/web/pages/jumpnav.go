package pages

import (
	"roommates/internal/locales"

	"github.com/a-h/templ"
)

// jumpLink is a section of a long page in its [jumpNav], ID is the section's id
type jumpLink struct {
	ID    string
	Label locales.LK
}

var jumpNavScript = templ.NewOnceHandle()
