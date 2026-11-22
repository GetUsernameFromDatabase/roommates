// Package utils a collection of helper functions used throughout the codebase
package utils

import (
	"roommates/internal/routes"

	"github.com/jackc/pgx/v5/pgtype"
)

// ParseUUID parses a UUID from a path param, session value etc., in the
// [routes.ShortUUID] form or the standard 36-character one
func ParseUUID(s string) (pgtype.UUID, error) {
	if u, ok := routes.ParseShortUUID(s); ok {
		return u, nil
	}
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}
