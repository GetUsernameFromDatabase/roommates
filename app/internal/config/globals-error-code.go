package config

//lint:file-ignore ST1003 goal is to ignore underscore issues

import "errors"

type ErrorCode string

const (
	EC_MissingLocale ErrorCode = "MISSING_LOCALE"

	EC_InternalError ErrorCode = "INTERNAL_ERROR"
)

func (c ErrorCode) String() string {
	return string(c)
}

var (
	ErrorInvalidCredential    = errors.New("invalid credentials")
	ErrorAccountAlreadyExists = errors.New("account already exists")
	// also used when the thing exists but the user is not allowed to know that
	ErrorNotFound      = errors.New("not found")
	ErrorForbidden     = errors.New("forbidden")
	ErrorInviteInvalid = errors.New("invite is invalid or expired")
	// the only owner cannot leave while other members remain
	ErrorLastOwner = errors.New("last owner cannot leave")
	// unknown, expired or already used email token
	ErrorTokenInvalid     = errors.New("token is invalid or expired")
	ErrorEmailNotVerified = errors.New("email is not verified")
	// a form value that fails validation in a service
	ErrorInvalidInput = errors.New("invalid input")
)
