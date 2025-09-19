// will have all types and constants just to have it more easily renamed in .go files
package components

import (
	"roommates/rdb"
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

const (
	HfId              = "house-form"
	HfSearchResultsId = "HouseRoomatesInputSearchResults"
	HfRoomateInputId  = "houseForm-roommates-input"
)

type SPageWrapper struct {
	AuthInfo *rdb.UserSessionValue
	PathURL  string
}

const IdRootLayout = "root-layout"
