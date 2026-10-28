// Package config constants (mostly keys) used for accessing data in a centralized manner
//
// This package should be as independent as possible to avoid circular dependencies
package config

// route paths for gin
const (
	RouteLogin      = "/login"
	RouteLogout     = "/logout"
	RouteRegister   = "/register"
	RouteHouseholds = "/households"
	// takes ?code=
	RouteInvite = "/invite"
	// prefix of everything scoped to one household, followed by /:householdID
	RouteHousehold = "/h"
	// takes ?token=
	RouteVerifyEmail       = "/verify-email"
	RouteVerifyEmailResend = "/verify-email/resend"
	RouteForgotPassword    = "/forgot-password"
	// takes ?token=
	RouteResetPassword = "/reset-password"
)

// QueryNext is the query and form key of the path to return to after signing in
const QueryNext = "next"

// memberships.role values, the table CHECKs these
const (
	RoleOwner  = "owner"
	RoleMember = "member"
)

// keys of values set on [gin.Context] by middleware
const (
	// services.Membership set by middleware.NewHouseholdMember
	GinKeyMembership = "membership"
)

// path params
const (
	ParamHouseholdID = "householdID"
	ParamInviteID    = "inviteID"
	ParamUserID      = "userID"
)
