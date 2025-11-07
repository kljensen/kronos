package internal

import (
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// DefaultTimezoneAbbrMap is a map of common timezone abbreviations to their offsets in minutes.
// Values can be either:
// - int: a fixed offset in minutes
// - AmbiguousTimezoneMap: for timezones that observe DST
//
//nolint:staticcheck // SA1019 Internal package legitimately uses deprecated types
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
			return helpers.GetNthWeekdayOfMonth(year, time.March, time.Sunday, 2, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends 1st Sunday of November at 2 AM
			return helpers.GetNthWeekdayOfMonth(year, time.November, time.Sunday, 1, 2)
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
			return helpers.GetLastWeekdayOfMonth(year, time.March, time.Sunday, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends last Sunday of October at 3 AM
			return helpers.GetLastWeekdayOfMonth(year, time.October, time.Sunday, 3)
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
