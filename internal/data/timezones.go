package data

import "time"

// AmbiguousTimezoneMap defines a timezone that has different offsets
// depending on whether daylight saving time (DST) is in effect.
type AmbiguousTimezoneMap struct {
	// TimezoneOffsetDuringDst is the offset in minutes during DST.
	TimezoneOffsetDuringDst int

	// TimezoneOffsetNonDst is the offset in minutes when DST is not in effect.
	TimezoneOffsetNonDst int

	// DstStart returns the start date of DST for the given year.
	DstStart func(year int) time.Time

	// DstEnd returns the end date of DST for the given year.
	DstEnd func(year int) time.Time
}

// TimezoneAbbrMap maps timezone abbreviations to their offsets.
// Values can be either a simple offset (in minutes) or an AmbiguousTimezoneMap
// for timezones that observe DST.
type TimezoneAbbrMap map[string]any

// DefaultTimezoneAbbrMap contains the standard timezone abbreviations and their offsets.
// External code should use this through the public API rather than accessing directly.
var DefaultTimezoneAbbrMap = TimezoneAbbrMap{
	// UTC/GMT
	"UTC": 0,
	"GMT": 0,
	"Z":   0,

	// North American Timezones
	"EST": -300,
	"EDT": -240,
	"ET": &AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -240, // EDT = UTC-4
		TimezoneOffsetNonDst:    -300, // EST = UTC-5
		DstStart: func(year int) time.Time {
			// DST starts 2nd Sunday of March at 2 AM
			return GetNthWeekdayOfMonth(year, time.March, time.Sunday, 2, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends 1st Sunday of November at 2 AM
			return GetNthWeekdayOfMonth(year, time.November, time.Sunday, 1, 2)
		},
	},
	"CST":  -360,
	"CDT":  -300,
	"MST":  -420,
	"MDT":  -360,
	"PST":  -480,
	"PDT":  -420,
	"AKST": -540,
	"AKDT": -480,
	"HST":  -600,
	"HAST": -600,
	"HADT": -540,

	// European Timezones
	"BST":  60,
	"IST":  60,
	"WET":  0,
	"WEST": 60,
	"CET": &AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: 120, // CEST = UTC+2
		TimezoneOffsetNonDst:    60,  // CET = UTC+1
		DstStart: func(year int) time.Time {
			// DST starts last Sunday of March at 2 AM
			return GetLastWeekdayOfMonth(year, time.March, time.Sunday, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends last Sunday of October at 3 AM
			return GetLastWeekdayOfMonth(year, time.October, time.Sunday, 3)
		},
	},
	"CEST": 120,
	"EET":  120,
	"EEST": 180,
	"MSK":  180,

	// Asian Timezones
	"JST":       540,
	"KST":       540,
	"HKT":       480,
	"SGT":       480,
	"CST_CHINA": 480, // China Standard Time
	"IST_INDIA": 330, // India Standard Time
	"PKT":       300,
	"WIB":       420,
	"WITA":      480,
	"WIT":       540,

	// Australian Timezones
	"AEST": 600,
	"AEDT": 660,
	"ACST": 570,
	"ACDT": 630,
	"AWST": 480,

	// Other Timezones
	"NZST": 720,
	"NZDT": 780,
	"BRT":  -180,
	"ART":  -180,
	"GET":  240, // Georgia Eastern Time (UTC+4)
}

// GetNthWeekdayOfMonth returns the Nth occurrence of a weekday in a given month.
// For example, the 2nd Sunday of March 2024 at 2 AM.
func GetNthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int, hour int) time.Time {
	// Start at the first day of the month
	firstDay := time.Date(year, month, 1, hour, 0, 0, 0, time.UTC)

	// Find the first occurrence of the target weekday
	firstWeekday := firstDay
	for firstWeekday.Weekday() != weekday {
		firstWeekday = firstWeekday.AddDate(0, 0, 1)
	}

	// Add (n-1) weeks to get the Nth occurrence
	return firstWeekday.AddDate(0, 0, (n-1)*7)
}

// GetLastWeekdayOfMonth returns the last occurrence of a weekday in a given month.
// For example, the last Sunday of October 2024 at 3 AM.
func GetLastWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, hour int) time.Time {
	// Start at the last day of the month
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	lastDay := time.Date(nextYear, nextMonth, 1, hour, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	// Walk backwards to find the last occurrence of the target weekday
	for lastDay.Weekday() != weekday {
		lastDay = lastDay.AddDate(0, 0, -1)
	}

	return lastDay
}
