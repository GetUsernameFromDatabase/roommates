package pages

import (
	"slices"
	"strconv"
	"time"

	g "roommates/internal/config"
	"roommates/internal/locales"
)

// scheduleNone is the schedule control's value for an event that does not repeat
const scheduleNone = "none"

// schedulePick is what the [schedulePicker] shows: a weekdays schedule, "every N unit" and with
// None also no repeat, which is a Unit of ""
type schedulePick struct {
	Unit     string
	Every    int
	Weekdays []time.Weekday
	None     bool
}

// value is the schedule control's value, the unit is the select's for "every"
func (p schedulePick) value() string {
	switch {
	case p.None && p.Unit == "":
		return scheduleNone
	case p.Unit == g.ChoreUnitWeekdays:
		return g.ChoreUnitWeekdays
	}
	return choreEveryValue
}

func (p schedulePick) options() []segmentOption {
	if p.None {
		return append([]segmentOption{{Value: scheduleNone, Label: locales.LK_CalendarRepeatNone}}, choreScheduleOptions...)
	}
	return choreScheduleOptions
}

func (p schedulePick) picked(d time.Weekday) bool {
	return slices.Contains(p.Weekdays, d)
}

// scheduleEvery is the number input's value, 1 until one is typed
func scheduleEvery(n int) string {
	if n < 1 {
		return "1"
	}
	return strconv.Itoa(n)
}
