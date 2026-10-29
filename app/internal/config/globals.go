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
	// POST theme, mode and next
	RouteAppearance = "/appearance"
	// the per-device settings, appearance and notifications
	RouteSettings = "/settings"
	// POST the "avatar" file under [RouteSettings] (at /settings or a household's), or delete it
	RouteAvatar       = "/avatar"
	RouteAvatarDelete = "/avatar/delete"
	// GET /users/:userID/avatar is a user's picture
	RouteUsers = "/users"
	// POST the JSON of a browser's PushSubscription
	RoutePushSubscribe = "/push/subscribe"
	// POST {"endpoint": ...}
	RoutePushUnsubscribe = "/push/unsubscribe"
	// the web app manifest, makes the app installable
	RouteManifest = "/manifest.webmanifest"
	// the service worker shows it when a page cannot be loaded
	RouteOffline = "/offline"
)

// QueryNext is the query and form key of the path to return to after signing in
const QueryNext = "next"

// memberships.role values, the table CHECKs these
const (
	RoleOwner  = "owner"
	RoleMember = "member"
)

// notes.visibility values, the table CHECKs these
const (
	NoteVisibilityHousehold = "household"
	NoteVisibilityPrivate   = "private"
)

// ExpenseCategories are the expenses.category values, the table CHECKs these
var ExpenseCategories = []string{"groceries", "utilities", "rent", "household", "other"}

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
	ParamListID      = "listID"
	ParamItemID      = "itemID"
	ParamNoteID      = "noteID"
	ParamFileID      = "fileID"
	ParamExpenseID   = "expenseID"
	ParamMessageID   = "messageID"
	ParamSettleID    = "settlementID"
)
