package kronos

import "time"

// AmbiguousTimezoneMap defines a timezone that has different offsets
// depending on whether daylight saving time (DST) is in effect.
//
// Deprecated: Direct use of AmbiguousTimezoneMap is discouraged. Use the builder pattern
// with WithOption instead:
//
//	customTimezones := kronos.TimezoneAbbrMap{
//	    "ET": &kronos.AmbiguousTimezoneMap{
//	        TimezoneOffsetDuringDst: -240,
//	        TimezoneOffsetNonDst: -300,
//	        DstStart: func(year int) time.Time { ... },
//	        DstEnd: func(year int) time.Time { ... },
//	    },
//	}
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.TimezoneOverrides = customTimezones
//	    })
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
//
// Deprecated: Direct use of TimezoneAbbrMap is discouraged. Use the builder pattern
// with WithOption instead:
//
//	customTimezones := kronos.TimezoneAbbrMap{
//	    "CUSTOM": 123,  // UTC+2:03
//	    "TEST": -456,   // UTC-7:36
//	}
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.TimezoneOverrides = customTimezones
//	    })
type TimezoneAbbrMap map[string]interface{}

// DebugHandler is a function that handles debug events.
// It receives a debug message for logging or analysis.
//
// Deprecated: Direct use of DebugHandler is discouraged. Use the builder pattern
// with WithOption instead:
//
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.DebugHandler = func(msg string) {
//	            log.Printf("DEBUG: %s", msg)
//	        }
//	    })
type DebugHandler func(message string)

// defaultTimezoneAbbrMap is an internal copy of timezone abbreviations for use within the root package.
// External code should use the types and functions provided by the public API.
var defaultTimezoneAbbrMap = TimezoneAbbrMap{
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
	"CET": &AmbiguousTimezoneMap{
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

// ToTimezoneOffset converts various timezone representations to an offset in minutes.
// It supports:
// - Numeric offsets (int) - returned as-is
// - Timezone names (string) - looked up in overrides map or default timezone map
// - Ambiguous timezones with DST - determined based on the instant time
//
// Returns nil if the timezone cannot be resolved.
func toTimezoneOffset(tz interface{}, instant time.Time, overrides TimezoneAbbrMap) *int {
	if tz == nil {
		return nil
	}

	// If it's already a numeric offset, return it
	if offset, ok := tz.(int); ok {
		return &offset
	}

	// If it's a string, look it up in the maps
	if name, ok := tz.(string); ok {
		// First check overrides
		if overrides != nil {
			if val, exists := overrides[name]; exists {
				return resolveTimezoneValue(val, instant)
			}
		}

		// Then check default map
		if val, exists := defaultTimezoneAbbrMap[name]; exists {
			return resolveTimezoneValue(val, instant)
		}

		// Try to load as a location name (e.g., "America/New_York")
		loc, err := time.LoadLocation(name)
		if err == nil {
			_, offset := instant.In(loc).Zone()
			offsetMinutes := offset / 60
			return &offsetMinutes
		}
	}

	return nil
}

// resolveTimezoneValue resolves a timezone value from the map to an offset.
// Handles both simple integer offsets and AmbiguousTimezoneMap entries.
func resolveTimezoneValue(val interface{}, instant time.Time) *int {
	// Simple integer offset
	if offset, ok := val.(int); ok {
		return &offset
	}

	// Ambiguous timezone with DST (value type)
	if ambiguous, ok := val.(AmbiguousTimezoneMap); ok {
		return resolveAmbiguousTimezone(ambiguous, instant)
	}

	// Ambiguous timezone with DST (pointer type)
	if ambiguous, ok := val.(*AmbiguousTimezoneMap); ok {
		return resolveAmbiguousTimezone(*ambiguous, instant)
	}

	return nil
}

// resolveAmbiguousTimezone resolves an ambiguous timezone to its offset based on the instant.
func resolveAmbiguousTimezone(ambiguous AmbiguousTimezoneMap, instant time.Time) *int {
	// Without a valid instant, we can't determine DST status
	if instant.IsZero() {
		return nil
	}

	year := instant.Year()
	dstStart := ambiguous.DstStart(year)
	dstEnd := ambiguous.DstEnd(year)

	// Check if instant is during DST period
	if instant.After(dstStart) && !instant.After(dstEnd) {
		offset := ambiguous.TimezoneOffsetDuringDst
		return &offset
	}

	// Not during DST
	offset := ambiguous.TimezoneOffsetNonDst
	return &offset
}

// GetNthWeekdayOfMonth returns the date of the nth occurrence of a given weekday
// in a given month and year.
//
// For example: 2nd Sunday of March 2020, or 1st Monday of November 2020.
//
// Parameters:
// - year: The year
// - month: The month (1-12)
// - weekday: The target weekday (0=Sunday, 6=Saturday)
// - n: The occurrence (1=first, 2=second, 3=third, 4=fourth)
// - hour: The hour of day for the returned time
//
//nolint:gofumpt // Function formatting is correct
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

// GetLastWeekdayOfMonth returns the date of the last occurrence of a given weekday
// in a given month and year.
//
// For example: Last Sunday of October 2020.
//
// Parameters:
// - year: The year
// - month: The month (1-12)
// - weekday: The target weekday (0=Sunday, 6=Saturday)
// - hour: The hour of day for the returned time
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
