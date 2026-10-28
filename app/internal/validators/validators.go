// Package validators makes a singleton of [govalidator]
// and holds other validation functions
//
// [govalidator]: https://pkg.go.dev/github.com/go-playground/validator/v10#pkg-overview
package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var Main = validator.New()

// LogValidationErrors logs [validator.ValidationErrors] errors
//
// made to be used with [validator.Validate.Struct]
func LogValidationErrors(log zerolog.Logger, err error) {
	// TODO: maybe check for [validator.InvalidValidationError]
	// https://github.com/go-playground/validator#error-return-value
	for _, err := range err.(validator.ValidationErrors) {
		log.Err(err).
			Str("field", err.Field()).
			Any("value", err.Value()).
			Str("tag", err.Tag()).
			Msg("")
	}
}

// -----------------------------------------------------------------------------

// ValidateCookieHashBlockPair see https://github.com/gorilla/securecookie#examples
//   - hashKey must be a key with 32 or 64 bytes
//   - blockKey valid lengths are 16 (will select AES-128), 24 (AES-192), or 32 (AES-256) bytes
func ValidateCookieHashBlockPair(hashKey, blockKey []byte) error {
	if len(hashKey) != 32 && len(hashKey) != 64 {
		return fmt.Errorf("hash key must be 32 or 64 bytes, got %d", len(hashKey))
	}
	switch len(blockKey) {
	case 16, 24, 32:
	default:
		return fmt.Errorf("block key must be 16, 24 or 32 bytes, got %d", len(blockKey))
	}
	return nil
}
