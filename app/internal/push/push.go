// Package push sends web push notifications to browsers that subscribed
//
// https://developer.mozilla.org/en-US/docs/Web/API/Push_API
package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/logger"
	"slices"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
)

var log = logger.Main.With().Str("package", "push").Logger()

// ErrGone means the subscription cannot be used anymore and can be deleted
var ErrGone = errors.New("push subscription gone")

type Keys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

// Subscription is what PushSubscription.toJSON() gives in the browser
type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     Keys   `json:"keys"`
}

// Message is read by the service worker, web/assets/pwa/sw.js
type Message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	// page opened when the notification is clicked
	URL string `json:"url"`
	// a notification replaces the shown one with the same tag
	Tag string `json:"tag"`
}

type Pusher interface {
	Send(ctx context.Context, sub Subscription, msg Message) error
	// PublicKey is the applicationServerKey of subscribe(), empty when pushes are not sent
	PublicKey() string
}

// New returns a pusher with the VAPID keys from env, see [g.EVK_VAPIDPublicKey]
//
// in development without keys, notifications are logged instead
func New() Pusher {
	if g.EVK_VAPIDPrivateKey.Get("") == "" && g.IsDevelopment() {
		log.Warn().Msg("VAPID_PRIVATE_KEY not set, push notifications are written to the log")
		return LogPusher{}
	}
	return &WebPusher{
		publicKey:  g.EVK_VAPIDPublicKey.MustGet(),
		privateKey: g.EVK_VAPIDPrivateKey.MustGet(),
		subject:    g.EVK_VAPIDSubject.Get(g.GetAppBaseURL()),
		client:     &http.Client{Transport: nil, CheckRedirect: nil, Jar: nil, Timeout: 30 * time.Second},
	}
}

// LogPusher writes notifications to the log, for development
type LogPusher struct{}

func (LogPusher) Send(_ context.Context, sub Subscription, msg Message) error {
	log.Info().Str("endpoint", sub.Endpoint).Str("title", msg.Title).Str("body", msg.Body).Str("url", msg.URL).Msg("push not sent")
	return nil
}

func (LogPusher) PublicKey() string { return "" }

type WebPusher struct {
	publicKey  string
	privateKey string
	subject    string
	client     *http.Client
}

func (p *WebPusher) PublicKey() string { return p.publicKey }

func (p *WebPusher) Send(ctx context.Context, sub Subscription, msg Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	s := &webpush.Subscription{Endpoint: sub.Endpoint, Keys: webpush.Keys{Auth: sub.Keys.Auth, P256dh: sub.Keys.P256dh}}
	res, err := webpush.SendNotificationWithContext(ctx, payload, s, &webpush.Options{
		HTTPClient:      p.client,
		Subscriber:      p.subject,
		VAPIDPublicKey:  p.publicKey,
		VAPIDPrivateKey: p.privateKey,
		// a phone that is off for a day misses it, by then the app shows it anyway
		TTL:     int((24 * time.Hour).Seconds()),
		Urgency: webpush.UrgencyNormal,
		// Mozilla limits the body to 4096 bytes once base64 encoded (about 3 KB),
		// webpush-go pads every message to RecordSize, messages stay under 1 KB
		RecordSize:      2048,
		Topic:           "",
		VapidExpiration: time.Time{},
	})
	if err != nil {
		return err
	}
	// push services explain a rejection in the body
	reason, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
	_ = res.Body.Close()
	switch {
	// these never work again for this subscription, e.g. 403 after the VAPID keys changed,
	// 413 and 429 are about the message or the rate instead
	case slices.Contains([]int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusGone}, res.StatusCode):
		return ErrGone
	case res.StatusCode >= 300:
		return fmt.Errorf("push service answered %s: %s", res.Status, reason)
	}
	return nil
}
