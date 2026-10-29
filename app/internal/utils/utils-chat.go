package utils

import "context"

type unreadChatKey struct{}

// WithUnreadChat stores how many chat messages the viewer has not read, for the nav
func WithUnreadChat(ctx context.Context, n int64) context.Context {
	return context.WithValue(ctx, unreadChatKey{}, n)
}

func GetUnreadChat(ctx context.Context) int64 {
	n, _ := ctx.Value(unreadChatKey{}).(int64)
	return n
}
