package pages

import (
	"roommates/internal/services"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCalendarGrid(t *testing.T) {
	day := func(m time.Month, d int) time.Time { return time.Date(2026, m, d, 0, 0, 0, 0, time.UTC) }
	for _, c := range []struct {
		month    time.Time
		from, to time.Time
	}{
		// starts on a Sunday, ends on a Monday
		{day(11, 1), day(10, 26), day(12, 7)},
		// starts on a Tuesday, the next month on a Friday
		{day(12, 1), day(11, 30), day(12, 1).AddDate(0, 1, 3)},
		// February 2027 is exactly four weeks from Monday
		{time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 3, 1, 0, 0, 0, 0, time.UTC)},
	} {
		from, to := CalendarGrid(c.month)
		if !from.Equal(c.from) || !to.Equal(c.to) {
			t.Errorf("%s: got %s–%s, want %s–%s", c.month.Format("2006-01"), from, to, c.from, c.to)
		}
		if n := len(calendarWeeks(c.month)); n*7 != int(to.Sub(from).Hours()/24) {
			t.Errorf("%s: %d weeks for %s–%s", c.month.Format("2006-01"), n, from, to)
		}
	}
}

func TestCalendarAgendaMonths(t *testing.T) {
	item := func(m time.Month, d int) services.CalendarItem {
		on := time.Date(2026, m, d, 0, 0, 0, 0, time.UTC)
		return services.CalendarItem{Kind: services.CalendarKindEvent, Date: on, Title: "", Location: "", AllDay: true, Clock: "", EndClock: "",
			ID: pgtype.UUID{}, On: on, AssigneeName: "", AmountCents: 0, Overdue: false}
	}
	items := []services.CalendarItem{item(11, 28), item(11, 30), item(12, 2)}
	headings := func(last time.Time) []bool {
		var got []bool
		for _, m := range calendarAgendaMonths(items, last) {
			got = append(got, m.Heading)
		}
		return got
	}
	// the first block names every month
	if got := headings(time.Time{}); len(got) != 2 || !got[0] || !got[1] {
		t.Errorf("first block: %v", got)
	}
	// a block after one that showed November does not name it again
	if got := headings(time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)); len(got) != 2 || got[0] || !got[1] {
		t.Errorf("after November: %v", got)
	}
	// after a block that ended in October, November is named
	if got := headings(time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)); len(got) != 2 || !got[0] {
		t.Errorf("after October: %v", got)
	}
}
