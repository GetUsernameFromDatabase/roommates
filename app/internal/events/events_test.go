package events

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestHub(t *testing.T) {
	h := NewHub()
	flat := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	other := pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	a, cancelA := h.Subscribe(flat)
	b, cancelB := h.Subscribe(other)
	defer cancelB()

	h.Publish(flat, TopicChat)
	if got := <-a; got != TopicChat {
		t.Errorf("got %q, want %q", got, TopicChat)
	}
	if len(b) != 0 {
		t.Error("another household got the event")
	}

	// a subscriber that stops reading must not block the others
	for range subscriberBuffer + 5 {
		h.Publish(flat, TopicShopping)
	}
	if len(a) != subscriberBuffer {
		t.Errorf("buffered %d, want %d", len(a), subscriberBuffer)
	}

	cancelA()
	if _, ok := h.subs[flat]; ok {
		t.Error("household without subscribers is kept")
	}
	h.Publish(flat, TopicChat)

	h.publishAll()
	if len(b) != len(Topics) {
		t.Errorf("after publishAll got %d topics, want %d", len(b), len(Topics))
	}
}
