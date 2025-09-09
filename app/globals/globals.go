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

const (
	HAuthorization HttpHeader = "Authorization"
)

const Csrf = "_csrf"

const (
	// key to authenticated user info in gin context
	GAuth = "authInfo"
)

// -----------------------------------------------------------------------------

var (
	ErrorInvalidCredential = errors.New("invalid credentials")
)