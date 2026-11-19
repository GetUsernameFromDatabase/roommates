package services

import (
	"context"
	"slices"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
)

// NotificationSetting is whether the member gets push notifications of a kind (g.NotifyKinds) from a household
type NotificationSetting struct {
	Kind string
	On   bool
}

// ListNotificationSettings is every kind in g.NotifyKinds order, a kind is on unless the member muted it
func (s *DefaultService) ListNotificationSettings(ctx context.Context, actor Membership) ([]NotificationSetting, error) {
	muted, err := s.db.ListNotificationMutes(ctx, dbqueries.ListNotificationMutesParams{
		HouseholdID: actor.HouseholdID, UserID: actor.UserID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]NotificationSetting, 0, len(g.NotifyKinds))
	for _, k := range g.NotifyKinds {
		out = append(out, NotificationSetting{Kind: k, On: !slices.Contains(muted, k)})
	}
	return out, nil
}

// SetNotification turns a kind on or off for the member in the household, g.ErrorNotFound for an unknown kind
func (s *DefaultService) SetNotification(ctx context.Context, actor Membership, kind string, on bool) error {
	if !slices.Contains(g.NotifyKinds, kind) {
		return g.ErrorNotFound
	}
	if on {
		return s.db.UnmuteNotification(ctx, dbqueries.UnmuteNotificationParams{
			HouseholdID: actor.HouseholdID, UserID: actor.UserID, Kind: kind,
		})
	}
	return s.db.MuteNotification(ctx, dbqueries.MuteNotificationParams{
		HouseholdID: actor.HouseholdID, UserID: actor.UserID, Kind: kind,
	})
}
