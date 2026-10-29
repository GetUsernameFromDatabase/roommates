package pages

import (
	"context"
	"fmt"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// householdURL is /h/:householdID followed by parts
func householdURL(id pgtype.UUID, parts ...string) templ.SafeURL {
	u := g.RouteHousehold + "/" + id.String()
	for _, p := range parts {
		u += "/" + url.PathEscape(p)
	}
	return templ.SafeURL(u)
}

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

func inviteURL(code string) templ.SafeURL {
	return templ.SafeURL(g.RouteInvite + "?" + url.Values{"code": {code}}.Encode())
}

// where on the people page the error of an action shows, next to that action
const (
	PeopleAtMembers     = "members"
	PeopleAtInvites     = "invites"
	PeopleAtInviteEmail = "invite-email"
	PeopleAtLeave       = "leave"
)

// formatDate shows only the date, the server time zone (UTC in Docker) is not the user's
func formatDate(t pgtype.Timestamptz) string {
	return t.Time.Format(time.DateOnly)
}

// formatSize is a byte count in B, KB or MB
func formatSize(n int64) string {
	switch {
	case n < 1<<10:
		return fmt.Sprintf("%d B", n)
	case n < 1<<20:
		return fmt.Sprintf("%.0f KB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	}
}

// RoleName is the translated memberships.role
func RoleName(ctx context.Context, role string) string {
	if role == g.RoleOwner {
		return utils.T(ctx, locales.LK_HouseholdRoleOwner)
	}
	return utils.T(ctx, locales.LK_HouseholdRoleMember)
}

// shoppingListURL is /h/:householdID/lists/:listID followed by parts
func shoppingListURL(actor services.Membership, l services.ShoppingList, parts ...string) templ.SafeURL {
	return householdURL(actor.HouseholdID, append([]string{"lists", l.List.ID.String()}, parts...)...)
}

// shoppingSwap posts a form with htmx and swaps in only #shopping-list of the response
func shoppingSwap(actor services.Membership, l services.ShoppingList, parts ...string) templ.Attributes {
	return templ.Attributes{
		"hx-post":   string(shoppingListURL(actor, l, parts...)),
		"hx-target": "#shopping-list",
		"hx-select": "#shopping-list",
		"hx-swap":   "outerHTML",
	}
}
