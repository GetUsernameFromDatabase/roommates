package pages

import (
	g "roommates/internal/config"
	"roommates/internal/locales"
)

// notificationLabels are the name and the description of a kind of notification
func notificationLabels(kind string) (title, description locales.LK) {
	switch kind {
	case g.NotifyChat:
		return locales.LK_NotificationsChat, locales.LK_NotificationsChatDescription
	case g.NotifyMoney:
		return locales.LK_NotificationsMoney, locales.LK_NotificationsMoneyDescription
	case g.NotifyRecurring:
		return locales.LK_NotificationsRecurring, locales.LK_NotificationsRecurringDescription
	case g.NotifyChores:
		return locales.LK_NotificationsChores, locales.LK_NotificationsChoresDescription
	case g.NotifyCalendar:
		return locales.LK_NotificationsCalendar, locales.LK_NotificationsCalendarDescription
	default:
		return locales.LK_NotificationsHousehold, locales.LK_NotificationsHouseholdDescription
	}
}
