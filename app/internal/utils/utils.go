// Package utils a collection of helper functions used throughout the codebase
package utils

import "github.com/jackc/pgx/v5/pgtype"

// ParseUUID parses a UUID from a path param, session value etc.
func ParseUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}
