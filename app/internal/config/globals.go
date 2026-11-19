// Package config constants (mostly keys) used for accessing data in a centralized manner
//
// This package should be as independent as possible to avoid circular dependencies
package config

import "time"

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
	// takes ?token=, the link sent to a new address
	RouteConfirmEmail = "/confirm-email"
	// POST theme, mode and next
	RouteAppearance = "/appearance"
	// POST language and next
	RouteLanguage = "/language"
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

// recurring_expenses.interval values, the table CHECKs these
const (
	IntervalWeek  = "week"
	IntervalMonth = "month"
	IntervalYear  = "year"
)

// Intervals are the repeat intervals in the order the form shows them
var Intervals = []string{IntervalWeek, IntervalMonth, IntervalYear}

// push notification kinds a member can turn off per household, notification_mutes.kind CHECKs these
const (
	NotifyChat      = "chat"
	NotifyMoney     = "money"
	NotifyRecurring = "recurring"
	NotifyHousehold = "household"
)

// NotifyKinds are the notification kinds in the order the settings show them
var NotifyKinds = []string{NotifyChat, NotifyMoney, NotifyRecurring, NotifyHousehold}

// HouseholdPurgeAfter is how long a deleted household stays restorable (read-only) before it is purged
const HouseholdPurgeAfter = 30 * 24 * time.Hour

// activity.grp values, the table CHECKs these
const (
	ActivityMoney    = "money"
	ActivityShopping = "shopping"
	ActivityNotes    = "notes"
	ActivityMembers  = "members"
)

// ActivityRetention is how long activity of each group is kept
var ActivityRetention = map[string]time.Duration{
	ActivityMoney:    365 * 24 * time.Hour,
	ActivityMembers:  365 * 24 * time.Hour,
	ActivityNotes:    90 * 24 * time.Hour,
	ActivityShopping: 30 * 24 * time.Hour,
}

// activity.kind values
const (
	ActExpenseAdded         = "expense_added"
	ActExpenseEdited        = "expense_edited"
	ActExpenseDeleted       = "expense_deleted"
	ActSharePaid            = "share_paid"
	ActShareUnpaid          = "share_unpaid"
	ActSettlementRequested  = "settlement_requested"
	ActSettlementConfirmed  = "settlement_confirmed"
	ActSettlementDeclined   = "settlement_declined"
	ActRecurringAdded       = "recurring_added"
	ActRecurringEdited      = "recurring_edited"
	ActRecurringDeleted     = "recurring_deleted"
	ActRecurringSkipped     = "recurring_skipped"
	ActListCreated          = "list_created"
	ActListRenamed          = "list_renamed"
	ActListDeleted          = "list_deleted"
	ActItemBought           = "item_bought"
	ActItemsCleared         = "items_cleared"
	ActNoteCreated          = "note_created"
	ActNoteEdited           = "note_edited"
	ActNoteDeleted          = "note_deleted"
	ActFileUploaded         = "file_uploaded"
	ActFileRemoved          = "file_removed"
	ActMemberJoined         = "member_joined"
	ActMemberLeft           = "member_left"
	ActMemberRemoved        = "member_removed"
	ActMemberPromoted       = "member_promoted"
	ActOwnershipTransferred = "ownership_transferred"
	ActSteppedDown          = "stepped_down"
	ActHouseholdRenamed     = "household_renamed"
	ActCurrencyChanged      = "currency_changed"
	ActHouseholdDeleted     = "household_deleted"
	ActHouseholdRestored    = "household_restored"
)

// ActivityGroup is the group (and so the retention) of each activity kind
var ActivityGroup = map[string]string{
	ActExpenseAdded: ActivityMoney, ActExpenseEdited: ActivityMoney, ActExpenseDeleted: ActivityMoney,
	ActSharePaid: ActivityMoney, ActShareUnpaid: ActivityMoney,
	ActSettlementRequested: ActivityMoney, ActSettlementConfirmed: ActivityMoney, ActSettlementDeclined: ActivityMoney,
	ActRecurringAdded: ActivityMoney, ActRecurringEdited: ActivityMoney, ActRecurringDeleted: ActivityMoney, ActRecurringSkipped: ActivityMoney,
	ActListCreated: ActivityShopping, ActListRenamed: ActivityShopping, ActListDeleted: ActivityShopping,
	ActItemBought: ActivityShopping, ActItemsCleared: ActivityShopping,
	ActNoteCreated: ActivityNotes, ActNoteEdited: ActivityNotes, ActNoteDeleted: ActivityNotes,
	ActFileUploaded: ActivityNotes, ActFileRemoved: ActivityNotes,
	ActMemberJoined: ActivityMembers, ActMemberLeft: ActivityMembers, ActMemberRemoved: ActivityMembers,
	ActMemberPromoted: ActivityMembers, ActOwnershipTransferred: ActivityMembers, ActSteppedDown: ActivityMembers,
	ActHouseholdRenamed: ActivityMembers, ActCurrencyChanged: ActivityMembers,
	ActHouseholdDeleted: ActivityMembers, ActHouseholdRestored: ActivityMembers,
}

// ExpenseCategories are the rows of the expense_categories table, which expenses and shopping lists reference
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
	ParamRecurringID = "recurringID"
)
