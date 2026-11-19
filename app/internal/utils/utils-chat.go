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

type pendingSettlementsKey struct{}

// WithPendingSettlements stores how many payments wait for the viewer to confirm them, for the nav
func WithPendingSettlements(ctx context.Context, n int64) context.Context {
	return context.WithValue(ctx, pendingSettlementsKey{}, n)
}

func GetPendingSettlements(ctx context.Context) int64 {
	n, _ := ctx.Value(pendingSettlementsKey{}).(int64)
	return n
}
