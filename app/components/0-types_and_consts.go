// will have all types and constants just to have it more easily renamed in .go files
package components

import (
	"roommates/globals"
	"roommates/rdb"

	"github.com/a-h/templ"
)

// html classes for label element
// will override default `uk-form-label`
type LabelClass string

// used to get desired element from templ component
//
// useful to have related elements in one place
type ElementType string

const (
	// meant to make an element visible, modal or commannd
	EOpener ElementType = "open"
	// thing meant to be opened by EOpener
	EModal ElementType = "modal"
)

// get icon from https://lucide.dev/icons/?focus
type Icon string

// id for house form element
const (
	HfId              = "house-form"
	HfSearchResultsId = "HouseRoomatesInputSearchResults"
	HfRoomateInputId  = "houseForm-roommates-input"
	HfModalId         = "house-modal"
)

type SPageWrapper struct {
	AuthInfo *rdb.UserSessionValue
	PathURL  string
}

const IdRootLayout = "root-layout"

// attributes for htmx elements
var (
	// hx boost, swapping innerHTML of IdRootLayout
	HtmxPageSwapAttributes = templ.Attributes{
		"hx-boost":  "true",
		"hx-target": "#" + IdRootLayout,
		"hx-swap":   "innerHTML",
	}
	// swaps out house form element on click
	// be sure to add house_id when editing house
	GetHtmxHouseForm = templ.Attributes{
		"hx-get":     globals.RHtmxHouseForm,
		"hx-target":  "#" + HfId,
		"hx-swap":    "outerHTML",
		"hx-trigger": "click",
	}
)
