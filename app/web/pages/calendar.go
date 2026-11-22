package pages

import "github.com/a-h/templ"

// calendarSwap are the attributes of a link that shows another day or month of the grid in place
func calendarSwap(url string) templ.Attributes {
	return templ.Attributes{"hx-get": url, "hx-target": "#cal-month", "hx-select": "#cal-month", "hx-swap": "outerHTML", "hx-replace-url": "true"}
}
