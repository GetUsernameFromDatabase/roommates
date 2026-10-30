package services

import (
	"context"
	"crypto/rand"
	"roommates/internal/locales"
	"roommates/internal/push"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestPushPreview(t *testing.T) {
	long := strings.Repeat("õ", 200)
	for body, want := range map[string]string{
		"Who took my milk?":       "Who took my milk?",
		"a\r\nsd\n\n  asd\tdasd ": "a sd asd dasd",
		long:                      strings.Repeat("õ", pushPreviewLength-1) + "…",
		strings.Repeat("🥛", 81):   strings.Repeat("🥛", pushPreviewLength-1) + "…",
	} {
		if got := pushPreview(body); got != want {
			t.Errorf("pushPreview(%q) = %q, want %q", body, got, want)
		}
	}
}

func TestValidPushSubscription(t *testing.T) {
	keys := push.Keys{P256dh: "BNcR", Auth: "tBHI"}
	for endpoint, want := range map[string]bool{
		"https://fcm.googleapis.com/fcm/send/abc":            true,
		"https://updates.push.services.mozilla.com/wpush/v2": true,
		"https://web.push.apple.com/QGuQ":                    true,
		"https://wns2-db5p.notify.windows.com/w/?token=x":    true,
		"http://fcm.googleapis.com/fcm/send/abc":             false,
		"https://fcm.googleapis.com:8443/fcm/send/abc":       false,
		"https://user@fcm.googleapis.com/fcm/send/abc":       false,
		"https://127.0.0.1/admin":                            false,
		"https://evilpush.services.mozilla.com.example/x":    false,
		"https://notify.windows.com.example.com/x":           false,
	} {
		if got := validPushSubscription(push.Subscription{Endpoint: endpoint, Keys: keys}); got != want {
			t.Errorf("%s: got %v, want %v", endpoint, got, want)
		}
	}
	if validPushSubscription(push.Subscription{Endpoint: "https://fcm.googleapis.com/x", Keys: push.Keys{P256dh: "", Auth: "a"}}) {
		t.Error("missing key was accepted")
	}
}

type sentPush struct {
	endpoint string
	msg      push.Message
}

// fakePusher collects notifications, endpoints in gone answer [push.ErrGone]
type fakePusher struct {
	mu   sync.Mutex
	sent []sentPush
	gone map[string]bool
}

func newFakePusher() *fakePusher {
	return &fakePusher{mu: sync.Mutex{}, sent: nil, gone: map[string]bool{}}
}

func (p *fakePusher) Send(_ context.Context, sub push.Subscription, msg push.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.gone[sub.Endpoint] {
		return push.ErrGone
	}
	p.sent = append(p.sent, sentPush{endpoint: sub.Endpoint, msg: msg})
	return nil
}

func (p *fakePusher) PublicKey() string { return "test" }

// takePushes waits for the notifications being sent and returns them
func takePushes(s *DefaultService) []sentPush {
	s.pushing.Wait()
	p := s.pusher.(*fakePusher)
	p.mu.Lock()
	defer p.mu.Unlock()
	sent := p.sent
	p.sent = nil
	return sent
}

func subscribeTestDevice(t *testing.T, s *DefaultService, userID pgtype.UUID, lang string) string {
	t.Helper()
	endpoint := "https://fcm.googleapis.com/fcm/send/" + rand.Text()
	sub := push.Subscription{Endpoint: endpoint, Keys: push.Keys{P256dh: "BNcR", Auth: "tBHI"}}
	if err := s.SubscribePush(t.Context(), userID, sub, lang); err != nil {
		t.Fatal(err)
	}
	return endpoint
}

func wantPushes(t *testing.T, step string, got []sentPush, endpoints ...string) {
	t.Helper()
	if len(got) != len(endpoints) {
		t.Fatalf("%s: got %d notifications %v, want %d", step, len(got), got, len(endpoints))
	}
	for i, e := range endpoints {
		if got[i].endpoint != e {
			t.Fatalf("%s: notification %d went to %s, want %s", step, i, got[i].endpoint, e)
		}
	}
}

func TestPushNotifications(t *testing.T) {
	s := newTestService(t)
	locales.LoadLocales()
	ctx := t.Context()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	ownerDevice := subscribeTestDevice(t, s, ownerM.UserID, "en")
	memberDevice := subscribeTestDevice(t, s, memberM.UserID, "et")
	subscribeTestDevice(t, s, strangerM.UserID, "en")

	if err := s.PostMessage(ctx, memberM, "Who took my milk?", pgtype.UUID{}, nil); err != nil {
		t.Fatal(err)
	}
	got := takePushes(s)
	wantPushes(t, "chat", got, ownerDevice)
	if got[0].msg.Body != "Who took my milk?" || !strings.HasSuffix(got[0].msg.URL, "/chat") {
		t.Fatalf("chat: got %+v", got[0].msg)
	}

	if _, err := s.CreateExpense(ctx, ownerM, ExpenseInput{
		Title:      "Rent",
		Amount:     "10",
		Category:   "rent",
		IncurredOn: "2026-10-01",
		Custom:     false,
		Members:    []pgtype.UUID{ownerM.UserID, memberM.UserID},
		Shares:     nil,
		ListID:     pgtype.UUID{},
		ItemID:     pgtype.UUID{},
	}); err != nil {
		t.Fatal(err)
	}
	got = takePushes(s)
	wantPushes(t, "expense", got, memberDevice)
	// in the language the member subscribed with
	if !strings.Contains(got[0].msg.Body, "sinu osa on 5,00") {
		t.Fatalf("expense: got %+v", got[0].msg)
	}

	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "5"); err != nil {
		t.Fatal(err)
	}
	wantPushes(t, "settlement", takePushes(s), ownerDevice)

	s.pusher.(*fakePusher).gone[ownerDevice] = true
	if err := s.PostMessage(ctx, memberM, "hello?", pgtype.UUID{}, nil); err != nil {
		t.Fatal(err)
	}
	wantPushes(t, "gone", takePushes(s))
	subs, err := s.db.ListPushSubscriptionsOfUsers(ctx, []pgtype.UUID{ownerM.UserID})
	if err != nil || len(subs) != 0 {
		t.Fatalf("gone subscription was kept: %v %v", subs, err)
	}

	if err := s.UnsubscribePush(ctx, memberM.UserID, memberDevice); err != nil {
		t.Fatal(err)
	}
	if err := s.PostMessage(ctx, ownerM, "bye", pgtype.UUID{}, nil); err != nil {
		t.Fatal(err)
	}
	wantPushes(t, "unsubscribed", takePushes(s))
}
