// will have all types and constants just to have it more easily renamed in .go files
package components

import (
	"roommates/rdb"

	"github.com/a-h/templ"
)

// html classes for label element
// will override default `uk-form-label`
type LabelClass string

// used to get desired element from templ component
//
// useful to have related elements in one place
//  NB: not to be used anymore, use htmx to replace modal-here and then open with hyperscript
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
	HfSearchResultsId = "HouseRoommatesInputSearchResults"
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
	AtrHtmxReplaceMeOnRevealed = templ.Attributes{
		"hx-trigger": "revealed",
		"hx-swap": "outerHTML",
	}
	// hx boost, swapping innerHTML of IdRootLayout
	AtrHtmxPageSwap = templ.Attributes{
		"hx-boost":  "true",
		"hx-target": "#" + IdRootLayout,
		"hx-swap":   "innerHTML",
	}
	// swap inner of #modal-here with modal then open with hyperscript
	AtrHtmxSwapModal = templ.Attributes{
		"hx-target": "#modal-here",
		"hx-swap": "innerHTML",
		"_": HSOpenModal,
	}
)

// hyperscript constants
const (
	// open modal after htmx load 
	HSOpenModal = "on htmx:afterOnLoad wait 10ms then call UIkit.modal('#modal-here > *').show()"
)

// programmatic UIkit elements
const (
	uiKitHRISR = "UIkit.dropdown('#" + HfSearchResultsId + "')"
)
