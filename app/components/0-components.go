// will have all types and constants just to have it more easily renamed in .go files
package components

import (
	"encoding/json"
	"roommates/rdb"

	"github.com/a-h/templ"
)

// html classes for label element
// will override default `uk-form-label`
type LabelClass string

// get icon from https://lucide.dev/icons/?focus
type Icon string

type SPageWrapper struct {
	AuthInfo *rdb.UserSessionValue
	PathURL  string
}

// used to get desired element from templ component
//
// useful to have related elements in one place
//
//	NB: not to be used anymore, use htmx to replace modal-here and then open with hyperscript
type ElementType string

// TODO: get rid of this
const (
	// meant to make an element visible, modal or commannd
	EOpener ElementType = "open"
	// thing meant to be opened by EOpener
	EModal ElementType = "modal"
)

// ---| vars ---

// attributes for htmx elements
var (
	AtrHxReplaceMeOnRevealed = templ.Attributes{
		"hx-trigger": "revealed",
		"hx-swap":    "outerHTML",
	}
	// hx boost, swapping innerHTML of IdRootLayout
	AtrHxPageSwap = templ.Attributes{
		"hx-boost":  "true",
		"hx-target": "#" + IdRootLayout,
		"hx-swap":   "innerHTML",
	}
	// swap inner of #modal-here with modal then open with hyperscript
	AtrHxSwapModal = templ.Attributes{
		"hx-target": "#modal-here",
		"hx-swap":   "innerHTML",
		"_":         HSOpenModal,
	}
)

// --- vars |---

// ---| consts ---

const IdRootLayout = "root-layout"

// id for house form element
const (
	HfId              = "house-form"
	HfSearchResultsId = "HouseRoommatesInputSearchResults"
	HfRoomateInputId  = "houseForm-roommates-input"
)

// id for house note element
const (
	HnId = "house-note"
)

// hyperscript constants
const (
	// open modal after htmx load
	HSOpenModal = "on htmx:afterOnLoad wait 10ms then call UIkit.modal('#modal-here > *').show()"
	// this is a solution to the problem caused by HSOpenModal putting modal outside of #modal-here
	// and essentially endlessly duplicating the modal in the dom on modal open
	//
	// I do not want to change HSOpenModal as that is a solution to visual problems
	// caused by following htmx uikit modal example https://htmx.org/examples/modal-uikit/
	//  - no need to deal with closing it manually as uikit handles that
	//  - will not work the same way as examples are on franken-ui
	//
	// wait 1s is used to allow animations to run their course, would like to get rid of it
	// but could not find a quick way to react to visibility change
	HSRemoveModalWhenHidden = "on mutation of @class if not me.classList.contains('uk-open') then wait 1s then remove me"
)

// programmatic UIkit elements
const (
	uiKitHRISR = "UIkit.dropdown('#" + HfSearchResultsId + "')"
)

// --- consts |---

// ---| funcs ---

func FormSwapOuterHxAttributes(id string) templ.Attributes {
	return templ.Attributes{
		"hx-swap":   "outerHTML",
		"hx-target": "#" + id,
	}
}

// json marshals the map and converts it into string
func HxValsData(data map[string]string) string {
	marshalled, _ := json.Marshal(data)
	return string(marshalled)
}

// --- funcs |---
