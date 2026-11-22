package pages

import (
	"roommates/internal/services"

	"github.com/a-h/templ"
)

// SettingsMember is what the settings show only inside a household
type SettingsMember struct {
	Actor       services.Membership
	Household   services.HouseholdSettings
	ShowPayment bool
	// push kinds switched on or off in this household
	Notifications []services.NotificationSetting
}

var avatarCropScript = templ.NewOnceHandle()
