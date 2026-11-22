// Package config constants (mostly keys) used for accessing data in a centralized manner
//
// This package should be as independent as possible to avoid circular dependencies
package config

import "time"

// route paths and params live in internal/routes

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
	NotifyChores    = "chores"
	NotifyCalendar  = "calendar"
)

// NotifyKinds are the notification kinds in the order the settings show them
var NotifyKinds = []string{NotifyChat, NotifyMoney, NotifyRecurring, NotifyChores, NotifyCalendar, NotifyHousehold}

// events.remind_minutes values, the table CHECKs these; an all-day event takes only
// RemindAtStart (its day) and RemindDayBefore, both at the household's reminder hour
const (
	RemindAtStart   = 0
	Remind15Min     = 15
	RemindHour      = 60
	RemindDayBefore = 1440
)

// RemindTimed and RemindAllDay are the reminder choices in the order the form shows them
var (
	RemindTimed  = []int{RemindAtStart, Remind15Min, RemindHour, RemindDayBefore}
	RemindAllDay = []int{RemindAtStart, RemindDayBefore}
)

// MaxUpcomingEvents caps a household's events whose last occurrence is today or later
const MaxUpcomingEvents = 200

// chores.rotation values, the table CHECKs these
const (
	RotationFixed       = "fixed"
	RotationLeastRecent = "least_recent"
	RotationManual      = "manual"
)

// Rotations are the rotation rules in the order the form shows them
var Rotations = []string{RotationFixed, RotationLeastRecent, RotationManual}

// how a chore repeats: chores.every_unit values (the table CHECKs these), or on weekdays
const (
	ChoreUnitDay      = "day"
	ChoreUnitWeek     = "week"
	ChoreUnitMonth    = "month"
	ChoreUnitWeekdays = "weekdays"
)

// ChoreEveryUnits are the units of an "every N" chore in the order the form shows them
var ChoreEveryUnits = []string{ChoreUnitDay, ChoreUnitWeek, ChoreUnitMonth}

// chore_log.kind values, the table CHECKs these
const (
	ChoreDone    = "done"
	ChoreSkipped = "skipped"
	ChorePassed  = "passed"
)

// HouseholdPurgeAfter is how long a deleted household stays restorable (read-only) before it is purged
const HouseholdPurgeAfter = 30 * 24 * time.Hour

// AccountEraseAfter is how long a deleted account can be restored by signing in before it is erased
const AccountEraseAfter = 30 * 24 * time.Hour

// DataExportValidFor is how long a built data export can be downloaded
const DataExportValidFor = 7 * 24 * time.Hour

// activity.grp values, the table CHECKs these
const (
	ActivityMoney    = "money"
	ActivityShopping = "shopping"
	ActivityNotes    = "notes"
	ActivityMembers  = "members"
	ActivityChores   = "chores"
	ActivityCalendar = "calendar"
)

// ActivityRetention is how long activity of each group is kept
var ActivityRetention = map[string]time.Duration{
	ActivityMoney:    365 * 24 * time.Hour,
	ActivityMembers:  365 * 24 * time.Hour,
	ActivityNotes:    90 * 24 * time.Hour,
	ActivityChores:   90 * 24 * time.Hour,
	ActivityCalendar: 90 * 24 * time.Hour,
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
	ActChoreAdded           = "chore_added"
	ActChoreDeleted         = "chore_deleted"
	ActChoreDone            = "chore_done"
	ActChoreSkipped         = "chore_skipped"
	ActChorePassed          = "chore_passed"
	ActEventAdded           = "event_added"
	ActEventDeleted         = "event_deleted"
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
	ActChoreAdded: ActivityChores, ActChoreDeleted: ActivityChores,
	ActChoreDone: ActivityChores, ActChoreSkipped: ActivityChores, ActChorePassed: ActivityChores,
	ActEventAdded: ActivityCalendar, ActEventDeleted: ActivityCalendar,
}

// ExpenseCategories are the rows of the expense_categories table, which expenses and shopping lists reference
var ExpenseCategories = []string{"groceries", "utilities", "rent", "household", "other"}

// keys of values set on [gin.Context] by middleware
const (
	// services.Membership set by middleware.NewHouseholdMember
	GinKeyMembership = "membership"
)

// Cache-Control values: content-addressed files never change under their URL, a private one
// (a user's own upload) must stay out of shared caches
const (
	CacheImmutable        = "public, max-age=31536000, immutable"
	CachePrivateImmutable = "private, max-age=31536000, immutable"
	CachePrivateDay       = "private, max-age=86400"
)

// HSTSValue is the Strict-Transport-Security header, a year
const HSTSValue = "max-age=31536000"

// Rate limits per client IP, each a count per window
const (
	// each attempt costs a bcrypt hash (~250ms): login, register, new password
	AuthRateLimit = 10
	// each sends an email, so mailboxes cannot be flooded
	EmailRateLimit = 10
	// invite codes are short enough to be worth guessing
	InviteRateLimit = 20
	// RateLimitWindow is the window of the per-IP limits above
	RateLimitWindow = time.Minute
)

// Invite emails per account, the per-IP limit alone is dodged by changing IP
const (
	InviteEmailLimit  = 10
	InviteEmailWindow = time.Hour
)

// HTTP server timeouts; there is no write timeout, it would cut the SSE streams
const (
	ReadHeaderTimeout = 10 * time.Second
	// ReadTimeout leaves room for slow uploads
	ReadTimeout     = 5 * time.Minute
	IdleTimeout     = 120 * time.Second
	ShutdownTimeout = 10 * time.Second
)
