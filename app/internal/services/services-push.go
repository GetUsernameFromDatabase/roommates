package services

import (
	"context"
	"errors"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/push"
	"slices"
	"strings"

	"github.com/invopop/ctxi18n"
	"github.com/jackc/pgx/v5/pgtype"
)

// the server POSTs to the endpoint a browser gives, only push services are allowed
// so it cannot be pointed at the internal network
var (
	pushHosts        = []string{"fcm.googleapis.com", "android.googleapis.com"}
	pushHostSuffixes = []string{".push.services.mozilla.com", ".push.apple.com", ".notify.windows.com"}
)

// pushPreviewLength is how much of a chat message a notification shows, the rest is in the chat
const pushPreviewLength = 80

func validPushSubscription(sub push.Subscription) bool {
	u, err := url.Parse(sub.Endpoint)
	if err != nil || len(sub.Endpoint) > 2048 || u.Scheme != "https" || u.User != nil || u.Port() != "" {
		return false
	}
	host := u.Hostname()
	if !slices.Contains(pushHosts, host) && !slices.ContainsFunc(pushHostSuffixes, func(s string) bool { return strings.HasSuffix(host, s) }) {
		return false
	}
	return sub.Keys.P256dh != "" && sub.Keys.Auth != "" && len(sub.Keys.P256dh) <= 256 && len(sub.Keys.Auth) <= 256
}

// PushPublicKey is the VAPID key browsers subscribe with, empty when pushes are not sent
func (s *DefaultService) PushPublicKey() string {
	return s.pusher.PublicKey()
}

// SubscribePush stores the browser's subscription for userID, notifications come in lang
func (s *DefaultService) SubscribePush(ctx context.Context, userID pgtype.UUID, sub push.Subscription, lang string) error {
	if !validPushSubscription(sub) {
		return g.ErrorInvalidInput
	}
	return s.db.UpsertPushSubscription(ctx, dbqueries.UpsertPushSubscriptionParams{
		Endpoint: sub.Endpoint,
		UserID:   userID,
		P256dh:   sub.Keys.P256dh,
		Auth:     sub.Keys.Auth,
		Lang:     lang,
	})
}

func (s *DefaultService) UnsubscribePush(ctx context.Context, userID pgtype.UUID, endpoint string) error {
	return s.db.DeletePushSubscription(ctx, dbqueries.DeletePushSubscriptionParams{Endpoint: endpoint, UserID: userID})
}

// notifyAsync sends build's message to every device of userIDs, in the language
// each device subscribed with, build gets the names of actor and their household
//
// it runs after the response so push services cannot slow it down, failures are only logged
func (s *DefaultService) notifyAsync(ctx context.Context, actor Membership, userIDs []pgtype.UUID, build func(ctx context.Context, userID pgtype.UUID, actorName, householdName string) push.Message) {
	if len(userIDs) == 0 {
		return
	}
	ctx = context.WithoutCancel(ctx)
	s.pushing.Go(func() {
		subs, err := s.db.ListPushSubscriptionsOfUsers(ctx, userIDs)
		if err != nil || len(subs) == 0 {
			if err != nil {
				log.Error().Err(err).Msg("could not list push subscriptions")
			}
			return
		}
		user, err := s.db.GetUser(ctx, actor.UserID)
		if err != nil {
			log.Error().Err(err).Msg("could not get user for push")
			return
		}
		household, err := s.db.GetHousehold(ctx, actor.HouseholdID)
		if err != nil {
			log.Error().Err(err).Msg("could not get household for push")
			return
		}
		for _, sub := range subs {
			lctx, err := ctxi18n.WithLocale(ctx, sub.Lang)
			if err != nil {
				lctx = ctx
			}
			msg := build(lctx, sub.UserID, user.Username, household.Name)
			err = s.pusher.Send(ctx, push.Subscription{Endpoint: sub.Endpoint, Keys: push.Keys{P256dh: sub.P256dh, Auth: sub.Auth}}, msg)
			if errors.Is(err, push.ErrGone) {
				err = s.db.DeleteGonePushSubscription(ctx, sub.Endpoint)
			}
			if err != nil {
				log.Error().Err(err).Str("endpoint", sub.Endpoint).Msg("could not send push notification")
			}
		}
	})
}

// householdPath is the URL path of a page of the household, like the pages build it
func householdPath(householdID pgtype.UUID, parts ...string) string {
	return g.RouteHousehold + "/" + householdID.String() + "/" + strings.Join(parts, "/")
}

// pushPreview shortens a chat message to one line for a notification
func pushPreview(body string) string {
	body = strings.Join(strings.Fields(body), " ")
	r := []rune(body)
	if len(r) <= pushPreviewLength {
		return body
	}
	return string(r[:pushPreviewLength-1]) + "…"
}
