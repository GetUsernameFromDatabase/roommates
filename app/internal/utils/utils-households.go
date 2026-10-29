package utils

import (
	"context"
	"roommates/internal/db/dbqueries"
)

type householdsKey struct{}

// WithHouseholds stores the viewer's households, for the nav
func WithHouseholds(ctx context.Context, h []dbqueries.ListHouseholdsOfUserRow) context.Context {
	return context.WithValue(ctx, householdsKey{}, h)
}

func GetHouseholds(ctx context.Context) []dbqueries.ListHouseholdsOfUserRow {
	h, _ := ctx.Value(householdsKey{}).([]dbqueries.ListHouseholdsOfUserRow)
	return h
}
