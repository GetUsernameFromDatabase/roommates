package handlers

import (
	"io"
	"roommates/internal/events"
	"time"

	"github.com/gin-gonic/gin"
)

// eventsPing keeps idle streams open through proxies
const eventsPing = 25 * time.Second

// Events streams the topics of the household's changes as server-sent events,
// the data is always "1" since an event without data is not dispatched by browsers
//
// the stream ends when the user leaves or is removed from the household
func (h *DefaultHandler) Events(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	topics, cancel := h.hub.Subscribe(actor.HouseholdID)
	defer cancel()

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("X-Accel-Buffering", "no")
	// an id makes a browser that reconnects by itself send Last-Event-ID,
	// the changes in between were missed, so the page refetches everything
	start := "id: 1\n\n"
	if ctx.GetHeader("Last-Event-ID") != "" {
		for _, t := range events.Topics {
			start += "event: " + t + "\ndata: 1\n\n"
		}
	}
	if _, err := io.WriteString(ctx.Writer, start); err != nil {
		return
	}
	ctx.Writer.Flush()
	ping := time.NewTicker(eventsPing)
	defer ping.Stop()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Request.Context().Done():
			return false
		case <-ping.C:
			_, err := io.WriteString(w, ": ping\n\n")
			return err == nil
		case t := <-topics:
			if t == events.TopicPeople {
				if _, err := h.service.GetMembership(ctx, actor.HouseholdID, actor.UserID); err != nil {
					return false
				}
			}
			_, err := io.WriteString(w, "event: "+t+"\ndata: 1\n\n")
			return err == nil
		}
	})
}
