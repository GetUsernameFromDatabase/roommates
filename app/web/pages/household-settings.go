package pages

import (
	"fmt"
	"strings"
)

// timeZoneLabel is an IANA name as people read it: "America/New_York" is "America / New York"
func timeZoneLabel(zone string) string {
	return strings.ReplaceAll(strings.ReplaceAll(zone, "/", " / "), "_", " ")
}

func reminderHourLabel(hour int) string {
	return fmt.Sprintf("%02d:00", hour)
}
