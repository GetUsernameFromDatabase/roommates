package pages

import (
	"context"
	"roommates/internal/utils"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// sections of a household, the bottom bar on phones and the tab row on wide screens
const (
	SectionShopping = "shopping"
	SectionExpenses = "expenses"
	SectionNotes    = "notes"
	SectionFiles    = "files"
	SectionChat     = "chat"
	SectionChores   = "chores"
	SectionCalendar = "calendar"
	SectionPeople   = "people"
	SectionActivity = "activity"
	SectionSettings = "settings"
)

// HouseholdNav is where in a household a page is
type HouseholdNav struct {
	HouseholdID pgtype.UUID
	Section     string
}

// pageMainClass is the single column of the signed in pages, 16px gutters on phones,
// the chat is wider so both sides of the conversation have room, and on phones
// only clears the tab bar so it gets most of the screen, with the safe area
// only in the installed app like [householdNav]
const pageMainClass = "mx-auto flex flex-1 w-full max-w-lg flex-col gap-8 px-4 py-6 sm:p-10 md:has-[#chat]:max-w-3xl lg:max-w-3xl lg:has-[#cal-month]:max-w-5xl has-[#chat]:min-h-0 max-lg:has-[#chat]:pt-2 max-lg:has-[#chat]:pb-18 [@media(display-mode:standalone)]:max-lg:has-[#chat]:pb-[calc(4.5rem+env(safe-area-inset-bottom))]"

// sheetClass is the one raised surface of a page, the thing used most on it
const sheetClass = "rounded-2xl bg-card p-4 text-card-foreground ring-1 ring-foreground/10 sm:p-5"

// the sections only in [householdMenu] on phones
const desktopOnly = "max-lg:hidden"

// navRowClass is a row of the nav menus and the sidebar footer, the current one is a pill
const navRowClass = "flex min-h-11 w-full items-center gap-3 rounded-full px-3 font-semibold outline-none hover:bg-secondary focus-visible:ring-3 focus-visible:ring-ring/50 aria-[current=page]:bg-secondary [&[aria-current=page]>svg]:text-primary"

// menuPanelClass is the panel a nav menu opens, it grows out of its button
const menuPanelClass = "absolute z-30 mt-2 rounded-3xl bg-card p-2 ring-1 ring-foreground/10 shadow-lg group-open:animate-in group-open:fade-in-0 group-open:zoom-in-95 motion-reduce:animate-none"

// toastClass is [toast], live-resync.js reads it from <body> for its request-failed toast
const toastClass = "toast fixed right-[max(1rem,env(safe-area-inset-right))] bottom-[calc(6rem+env(safe-area-inset-bottom))] left-[max(1rem,env(safe-area-inset-left))] z-30 mx-auto flex w-fit max-w-sm items-center gap-2.5 rounded-3xl bg-foreground py-3 pr-5 pl-4 text-sm font-medium text-background shadow-lg lg:right-6 lg:bottom-6 lg:left-auto lg:mx-0"

var dirtySaveScript = templ.NewOnceHandle()

// householdName is the name of one of the viewer's households, "" when it is not one of them
func householdName(ctx context.Context, id pgtype.UUID) string {
	for _, h := range utils.GetHouseholds(ctx) {
		if h.ID == id {
			return h.Name
		}
	}
	return ""
}

// pageTitle is "Section · Household" inside a household, the page's own title elsewhere
func pageTitle(ctx context.Context, title string, nav *HouseholdNav) string {
	if nav == nil {
		return title
	}
	name := householdName(ctx, nav.HouseholdID)
	if name == "" || name == title {
		return title
	}
	return title + " · " + name
}
