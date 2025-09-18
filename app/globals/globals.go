package globals

import "errors"

// Key for cookies
type CookieKey string

const (
	CSessionToken CookieKey = "session_token"
)

// Key for http header
type HttpHeader string

// if there becomes a need to have more, could be nice to copy from
// https://pkg.go.dev/github.com/go-http-utils/headers#pkg-constants
// constants for headers
const (
	//
	HAuthorization HttpHeader = "Authorization"
	HHXRequest     HttpHeader = "HX-Request"
	HHXRedirect    HttpHeader = "HX-Redirect"
)

const Csrf = "_csrf"

// constants for gin context keys
const (
	// key to authenticated user info in gin context
	GAuth = "authInfo" // do not see a need to add ContextKey type to this
)

// constants for routes, see routes.go
const (
	//
	RHouses    = "/houses"
	RLogin     = "/login"
	RMessaging = "/messaging"
	RNotes     = "/notes"
	RPayments  = "/payments"
	RProfile   = "/profile"
	RRegister  = "/register"
)

// -----------------------------------------------------------------------------

var (
	ErrorInvalidCredential    = errors.New("invalid credentials")
	ErrorAccountAlreadyExists = errors.New("account already exists")
	ErrorHtmxRequired         = errors.New("htmx required")
)
