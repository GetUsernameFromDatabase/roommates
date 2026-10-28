package globals

//lint:file-ignore ST1003 goal is to ignore underscore issues

type HTTPHeaderKey string

// constants for headers
// if there becomes a need to have more, could be nice to copy from
// https://pkg.go.dev/github.com/go-http-utils/headers#pkg-constants

// HTMX header
const (
	HHK_HXRequest  HTTPHeaderKey = "HX-Request"
	HHK_HXRedirect HTTPHeaderKey = "HX-Redirect"
)

const (
	HHK_Authorization  HTTPHeaderKey = "Authorization"
	HHK_AcceptLanguage HTTPHeaderKey = "Accept-Language"
)

// https://gin-gonic.com/en/docs/middleware/security-headers/#what-each-header-protects-against

func (c HTTPHeaderKey) String() string {
	return string(c)
}
