package config

//lint:file-ignore ST1003 goal is to ignore underscore issues

import "errors"

type ErrorCode string

const (
	EC_NoToken       ErrorCode = "NO_TOKEN"
	EC_MissingLocale ErrorCode = "MISSING_LOCALE"

	EC_InternalError ErrorCode = "INTERNAL_ERROR"
)

func (c ErrorCode) String() string {
	return string(c)
}

var (
	ErrorInvalidCredential    = errors.New("invalid credentials")
	ErrorAccountAlreadyExists = errors.New("account already exists")
)
