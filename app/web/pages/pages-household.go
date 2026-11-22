package pages

import (
	"context"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/utils"
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// sseTrigger is the hx-trigger of server-sent events named topics,
// and of live-resync when the page may have missed some (live-resync.js)
func sseTrigger(topics []string) string {
	return "sse:" + strings.Join(topics, ", sse:") + ", live-resync from:body"
}

// unreadLabel caps the count so the nav keeps its width
func unreadLabel(n int64) string {
	if n > 99 {
		return "99+"
	}
	return strconv.FormatInt(n, 10)
}

func inviteURL(code string) string {
	return routes.Invite + "?" + url.Values{"code": {code}}.Encode()
}

// where on the people page the error of an action shows, next to that action
const (
	PeopleAtMembers     = "members"
	PeopleAtInvites     = "invites"
	PeopleAtInviteEmail = "invite-email"
	PeopleAtLeave       = "leave"
	PeopleAtStepDown    = "step-down"
)

type membershipKey struct{}

// WithMembership keeps the membership of a household page for the layout, see [MembershipFrom]
func WithMembership(ctx context.Context, m services.Membership) context.Context {
	return context.WithValue(ctx, membershipKey{}, m)
}

// MembershipFrom is the membership middleware.NewHouseholdMember stored, ok is false outside a household
func MembershipFrom(ctx context.Context) (services.Membership, bool) {
	m, ok := ctx.Value(membershipKey{}).(services.Membership)
	return m, ok
}

// ownerCount is how many of members are owners
func ownerCount(members []dbqueries.ListActiveMembersRow) int {
	n := 0
	for _, m := range members {
		if m.Role == g.RoleOwner {
			n++
		}
	}
	return n
}

// formatDate shows only the date, the server time zone (UTC in Docker) is not the user's
func formatDate(ctx context.Context, t pgtype.Timestamptz) string {
	return utils.FormatDate(ctx, t.Time)
}

// formatSize is a byte count in B, KB or MB
func formatSize(ctx context.Context, n int64) string {
	return utils.FormatSize(ctx, n)
}

// RoleName is the translated memberships.role
func RoleName(ctx context.Context, role string) string {
	if role == g.RoleOwner {
		return utils.T(ctx, locales.LK_HouseholdRoleOwner)
	}
	return utils.T(ctx, locales.LK_HouseholdRoleMember)
}

// shoppingSwap posts a form to action with htmx and swaps in only #shopping-list of the response,
// as a view transition so a bought item slides down into the bag
func shoppingSwap(action string) templ.Attributes {
	return templ.Attributes{
		"hx-post":   action,
		"hx-target": "#shopping-list",
		"hx-select": "#shopping-list",
		"hx-swap":   "outerHTML transition:true",
	}
}

// itemTransitionName lets the browser follow an item from the pad to the bag
func itemTransitionName(itemID string) templ.SafeCSS {
	return templ.SafeCSS("view-transition-name:item-" + itemID)
}

// hasUnlinked says whether a bought item has no expense yet
func hasUnlinked(bought []dbqueries.ListShoppingItemsRow) bool {
	return slices.ContainsFunc(bought, func(it dbqueries.ListShoppingItemsRow) bool { return !it.ExpenseID.Valid })
}

// splitBought parts the items of a list, the query sorts the bought ones last
func splitBought(items []dbqueries.ListShoppingItemsRow) (open, bought []dbqueries.ListShoppingItemsRow) {
	i := slices.IndexFunc(items, func(it dbqueries.ListShoppingItemsRow) bool { return it.CheckedAt.Valid })
	if i < 0 {
		return items, nil
	}
	return items[:i], items[i:]
}
