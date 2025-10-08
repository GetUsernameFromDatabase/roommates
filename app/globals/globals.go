package globals

import "errors"

// Key for cookies
type CookieKey string

const (
	CSessionToken CookieKey = "session_token"
)

// Key for http header
type HttpHeader string

// constants for headers
// if there becomes a need to have more, could be nice to copy from
// https://pkg.go.dev/github.com/go-http-utils/headers#pkg-constants
const (
	HAuthorization HttpHeader = "Authorization"
	HHXRequest     HttpHeader = "HX-Request"
	HHXRedirect    HttpHeader = "HX-Redirect"
)

const Csrf = "_csrf"

// key for gin.Context
type RequestContextKey string

// constants for gin context keys
const (
	// key to authenticated user info in gin context
	GAuth RequestContextKey = "authInfo" // do not see a need to add ContextKey type to this
)

// constants for routes, see routes.go
const (
	RHouses    = "/houses"
	RLogin     = "/login"
	RMessaging = "/messaging"
	RNotes     = "/notes"
	RPayments  = "/payments"
	RProfile   = "/profile"
	RRegister  = "/register"
	RUser      = "/user"

	RHouseID = RHouses + "/:id"
	RUserID  = RUser + "/:id"
	RNoteID  = RNotes + "/:id"

	RHxRoomateSearch = RHouses + "/roomate-search"
	RHxHouseForm     = RHouses + "/house-form"

	RHxHouseResidentsBadge = RHouseID + "/residents-badge"
	RHxNoteForm            = RHouseID + "/note-form"

	RHxNoteInHouseAccordion = RNoteID + "/view-house-accordion"
)

// -----------------------------------------------------------------------------

var (
	ErrorInvalidCredential    = errors.New("invalid credentials")
	ErrorAccountAlreadyExists = errors.New("account already exists")
	ErrorHxRequired           = errors.New("htmx required")
	ErrorNotAllowedToModify   = errors.New("not allowed to modify")
	ErrorInvalidID            = errors.New("invalid id")
)
