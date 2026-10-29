// Package events tells open pages of a household that something in it changed
//
// triggers in the database send '<household_id> <topic>' on [Channel] when a
// transaction commits, [Hub.Run] listens and passes the topic to the subscribers
// of that household, the pages then refetch what they show
package events

import (
	"context"
	"roommates/internal/logger"
	"roommates/internal/utils"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var log = logger.Main.With().Str("package", "events").Logger()

// Channel is the Postgres notification channel of the household triggers
const Channel = "household_events"

// topics sent by the triggers in migration 00009, and TopicRead by the chat service
const (
	TopicShopping = "shopping"
	TopicNotes    = "notes"
	TopicFiles    = "files"
	TopicExpenses = "expenses"
	TopicPeople   = "people"
	TopicChat     = "chat"
	// someone read the chat, their badges on other devices go to zero
	TopicRead = "read"
)

// Topics are all topics, sent when events may have been missed
var Topics = []string{TopicShopping, TopicNotes, TopicFiles, TopicExpenses, TopicPeople, TopicChat, TopicRead}

// subscriberBuffer is how many topics wait for a slow subscriber, more are dropped,
// the page refetches on the next one anyway
const subscriberBuffer = 16

type Hub struct {
	mu   sync.Mutex
	subs map[pgtype.UUID]map[chan string]struct{}
}

func NewHub() *Hub {
	return &Hub{mu: sync.Mutex{}, subs: map[pgtype.UUID]map[chan string]struct{}{}}
}

// Subscribe returns the topics of the household's events, cancel must be called
// when they are no longer read
func (h *Hub) Subscribe(householdID pgtype.UUID) (topics <-chan string, cancel func()) {
	ch := make(chan string, subscriberBuffer)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subs[householdID] == nil {
		h.subs[householdID] = map[chan string]struct{}{}
	}
	h.subs[householdID][ch] = struct{}{}
	return ch, func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.subs[householdID], ch)
		if len(h.subs[householdID]) == 0 {
			delete(h.subs, householdID)
		}
	}
}

// Publish passes topic to the household's subscribers without waiting for them
func (h *Hub) Publish(householdID pgtype.UUID, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs[householdID] {
		select {
		case ch <- topic:
		default:
		}
	}
}

// publishAll sends every topic to every subscriber
func (h *Hub) publishAll() {
	h.mu.Lock()
	ids := make([]pgtype.UUID, 0, len(h.subs))
	for id := range h.subs {
		ids = append(ids, id)
	}
	h.mu.Unlock()
	for _, id := range ids {
		for _, t := range Topics {
			h.Publish(id, t)
		}
	}
}

// Run listens on [Channel] with one connection of pool until ctx is done,
// reconnecting with a backoff, after a reconnect every page refetches everything
func (h *Hub) Run(ctx context.Context, pool *pgxpool.Pool) {
	backoff := time.Second
	for first := true; ctx.Err() == nil; first = false {
		err := h.listen(ctx, pool, func() {
			backoff = time.Second
			if !first {
				h.publishAll()
			}
		})
		if ctx.Err() != nil {
			return
		}
		log.Error().Err(err).Dur("retry_in", backoff).Msg("listening for household events failed")
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, 30*time.Second)
	}
}

// listen calls listening once LISTEN has run, then publishes notifications until an error
func (h *Hub) listen(ctx context.Context, pool *pgxpool.Pool, listening func()) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, "LISTEN "+Channel); err != nil {
		return err
	}
	listening()
	for {
		n, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			// LISTEN stays on the connection, pgxpool would hand it out again
			_ = conn.Conn().Close(context.WithoutCancel(ctx))
			return err
		}
		id, topic, ok := strings.Cut(n.Payload, " ")
		householdID, err := utils.ParseUUID(id)
		if !ok || err != nil {
			log.Warn().Str("payload", n.Payload).Msg("bad household event")
			continue
		}
		h.Publish(householdID, topic)
	}
}
