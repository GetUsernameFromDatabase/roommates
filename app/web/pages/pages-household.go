package pages

import (
	"context"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/utils"
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

func inviteURL(code string) templ.SafeURL {
	return templ.SafeURL(g.RouteInvite + "?" + url.Values{"code": {code}}.Encode())
}

// formatDate shows only the date, the server time zone (UTC in Docker) is not the user's
func formatDate(t pgtype.Timestamptz) string {
	return t.Time.Format(time.DateOnly)
}

// RoleName is the translated memberships.role
func RoleName(ctx context.Context, role string) string {
	if role == g.RoleOwner {
		return utils.T(ctx, locales.LK_HouseholdRoleOwner)
	}
	return utils.T(ctx, locales.LK_HouseholdRoleMember)
}
