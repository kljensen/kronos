package data

import (
	"time"

	"github.com/kljensen/kronos"
)

// DefaultTimezoneAbbrMap is a map of common timezone abbreviations to their offsets in minutes.
// Values can be either:
// - int: a fixed offset in minutes
// - AmbiguousTimezoneMap: for timezones that observe DST
var DefaultTimezoneAbbrMap = kronos.TimezoneAbbrMap{
	// UTC/GMT
	"UTC": 0,
	"GMT": 0,
	"Z":   0,

	// North American Timezones
	"EST": -300,
	"EDT": -240,
	"ET": &kronos.AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -240, // EDT = UTC-4
		TimezoneOffsetNonDst:    -300, // EST = UTC-5
		DstStart: func(year int) time.Time {
			// DST starts 2nd Sunday of March at 2 AM
			return getNthWeekdayOfMonth(year, time.March, time.Sunday, 2, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends 1st Sunday of November at 2 AM
			return getNthWeekdayOfMonth(year, time.November, time.Sunday, 1, 2)
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
	"CET": &kronos.AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: 120, // CEST = UTC+2
		TimezoneOffsetNonDst:    60,  // CET = UTC+1
		DstStart: func(year int) time.Time {
			// DST starts last Sunday of March at 2 AM
			return getLastWeekdayOfMonth(year, time.March, time.Sunday, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends last Sunday of October at 3 AM
			return getLastWeekdayOfMonth(year, time.October, time.Sunday, 3)
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

// ApproximationWords is a list of words that indicate approximate time expressions
var ApproximationWords = []string{
	"about",
	"around",
	"roughly",
	"approximately",
	"approx",
	"circa",
}

// EmptyDuration represents an explicit empty duration.
// This is defined as zero day, second, and millisecond.
var EmptyDuration = kronos.Duration{
	kronos.TimeunitDay:         0,
	kronos.TimeunitSecond:      0,
	kronos.TimeunitMillisecond: 0,
}

// Helper functions for DST calculations

// getNthWeekdayOfMonth returns the date of the nth occurrence of a given weekday
// in a given month and year.
func getNthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int, hour int) time.Time {
	dayOfMonth := 0
	count := 0

	// Iterate through days of the month
	for count < n {
		dayOfMonth++
		date := time.Date(year, time.Month(month), dayOfMonth, hour, 0, 0, 0, time.UTC)

		// Check if this day is the target weekday
		if int(date.Weekday()) == int(weekday) {
			count++
		}
	}

	return time.Date(year, time.Month(month), dayOfMonth, hour, 0, 0, 0, time.UTC)
}

// getLastWeekdayOfMonth returns the date of the last occurrence of a given weekday
// in a given month and year.
func getLastWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, hour int) time.Time {
	// Start with the first day of the next month
	nextMonth := time.Date(year, time.Month(month)+1, 1, 12, 0, 0, 0, time.UTC)

	// Convert weekdays to 1-indexed (Monday=1, Sunday=7)
	targetWeekday := int(weekday)
	if targetWeekday == 0 {
		targetWeekday = 7
	}

	firstWeekdayNextMonth := int(nextMonth.Weekday())
	if firstWeekdayNextMonth == 0 {
		firstWeekdayNextMonth = 7
	}

	// Calculate how many days to go back
	var dayDiff int
	switch {
	case firstWeekdayNextMonth == targetWeekday:
		dayDiff = 7
	case firstWeekdayNextMonth < targetWeekday:
		dayDiff = 7 + firstWeekdayNextMonth - targetWeekday
	default:
		dayDiff = firstWeekdayNextMonth - targetWeekday
	}

	// Go back to find the last occurrence
	result := nextMonth.AddDate(0, 0, -dayDiff)
	return time.Date(result.Year(), result.Month(), result.Day(), hour, 0, 0, 0, time.UTC)
}
