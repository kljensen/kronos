package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToTimezoneOffset_NumericOffset(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)

	// Test positive offset
	offset := 300
	result := toTimezoneOffset(offset, instant, nil)
	assert.NotNil(t, result)
	assert.Equal(t, 300, *result)

	// Test negative offset
	offset = -480
	result = toTimezoneOffset(offset, instant, nil)
	assert.NotNil(t, result)
	assert.Equal(t, -480, *result)

	// Test zero offset
	offset = 0
	result = toTimezoneOffset(offset, instant, nil)
	assert.NotNil(t, result)
	assert.Equal(t, 0, *result)
}

func TestToTimezoneOffset_StringTimezone(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		tz       string
		expected int
	}{
		{"UTC", "UTC", 0},
		{"GMT", "GMT", 0},
		{"EST", "EST", -300},
		{"EDT", "EDT", -240},
		{"PST", "PST", -480},
		{"PDT", "PDT", -420},
		{"JST", "JST", 540},
		{"AEST", "AEST", 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTimezoneOffset(tt.tz, instant, nil)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestToTimezoneOffset_WithOverrides(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)

	overrides := TimezoneAbbrMap{
		"CUSTOM": 123,
		"TEST":   -456,
	}

	// Test custom timezone
	result := toTimezoneOffset("CUSTOM", instant, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, 123, *result)

	// Test override taking precedence over default
	overrides["UTC"] = 999
	result = toTimezoneOffset("UTC", instant, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, 999, *result)
}

func TestToTimezoneOffset_AmbiguousTimezone(t *testing.T) {
	// Create an ambiguous timezone (like ET - Eastern Time)
	et := AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -240, // EDT = UTC-4
		TimezoneOffsetNonDst:    -300, // EST = UTC-5
		DstStart: func(year int) time.Time {
			// DST starts 2nd Sunday of March at 2 AM
			return getNthWeekdayOfMonth(year, MonthMarch, WeekdaySunday, 2, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends 1st Sunday of November at 2 AM
			return getNthWeekdayOfMonth(year, MonthNovember, WeekdaySunday, 1, 2)
		},
	}

	overrides := TimezoneAbbrMap{
		"ET": et,
	}

	// Test during DST (June 1, 2020)
	instantDuringDst := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
	result := toTimezoneOffset("ET", instantDuringDst, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -240, *result, "Should use DST offset during summer")

	// Test during non-DST (January 1, 2020)
	instantNonDst := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	result = toTimezoneOffset("ET", instantNonDst, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -300, *result, "Should use non-DST offset during winter")

	// Test right at DST start (should be non-DST)
	instantAtDstStart := time.Date(2020, 3, 8, 1, 0, 0, 0, time.UTC)
	result = toTimezoneOffset("ET", instantAtDstStart, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -300, *result, "Should use non-DST offset before DST starts")

	// Test right after DST start (should be DST)
	instantAfterDstStart := time.Date(2020, 3, 8, 8, 0, 0, 0, time.UTC)
	result = toTimezoneOffset("ET", instantAfterDstStart, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -240, *result, "Should use DST offset after DST starts")
}

func TestToTimezoneOffset_NilTimezone(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
	result := toTimezoneOffset(nil, instant, nil)
	assert.Nil(t, result)
}

func TestToTimezoneOffset_UnknownTimezone(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
	result := toTimezoneOffset("UNKNOWN_TZ", instant, nil)
	assert.Nil(t, result)
}

func TestToTimezoneOffset_AmbiguousWithZeroInstant(t *testing.T) {
	et := AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -240,
		TimezoneOffsetNonDst:    -300,
		DstStart:                func(year int) time.Time { return time.Time{} },
		DstEnd:                  func(year int) time.Time { return time.Time{} },
	}

	overrides := TimezoneAbbrMap{
		"ET": et,
	}

	// Zero instant should return nil for ambiguous timezone
	result := toTimezoneOffset("ET", time.Time{}, overrides)
	assert.Nil(t, result)
}

func TestgetNthWeekdayOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    Month
		weekday  Weekday
		n        int
		hour     int
		expected time.Time
	}{
		{
			name:     "2nd Sunday of March 2020",
			year:     2020,
			month:    MonthMarch,
			weekday:  WeekdaySunday,
			n:        2,
			hour:     2,
			expected: time.Date(2020, 3, 8, 2, 0, 0, 0, time.UTC),
		},
		{
			name:     "1st Monday of November 2020",
			year:     2020,
			month:    MonthNovember,
			weekday:  WeekdayMonday,
			n:        1,
			hour:     2,
			expected: time.Date(2020, 11, 2, 2, 0, 0, 0, time.UTC),
		},
		{
			name:     "4th Thursday of November 2020 (Thanksgiving)",
			year:     2020,
			month:    MonthNovember,
			weekday:  WeekdayThursday,
			n:        4,
			hour:     0,
			expected: time.Date(2020, 11, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "1st Friday of January 2021",
			year:     2021,
			month:    MonthJanuary,
			weekday:  WeekdayFriday,
			n:        1,
			hour:     12,
			expected: time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "3rd Tuesday of July 2020",
			year:     2020,
			month:    MonthJuly,
			weekday:  WeekdayTuesday,
			n:        3,
			hour:     6,
			expected: time.Date(2020, 7, 21, 6, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNthWeekdayOfMonth(tt.year, tt.month, tt.weekday, tt.n, tt.hour)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestgetLastWeekdayOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    Month
		weekday  Weekday
		hour     int
		expected time.Time
	}{
		{
			name:     "Last Sunday of October 2020",
			year:     2020,
			month:    MonthOctober,
			weekday:  WeekdaySunday,
			hour:     3,
			expected: time.Date(2020, 10, 25, 3, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last Monday of December 2020",
			year:     2020,
			month:    MonthDecember,
			weekday:  WeekdayMonday,
			hour:     0,
			expected: time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last Friday of February 2020 (leap year)",
			year:     2020,
			month:    MonthFebruary,
			weekday:  WeekdayFriday,
			hour:     12,
			expected: time.Date(2020, 2, 28, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last Saturday of April 2021",
			year:     2021,
			month:    MonthApril,
			weekday:  WeekdaySaturday,
			hour:     18,
			expected: time.Date(2021, 4, 24, 18, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last Wednesday of June 2020",
			year:     2020,
			month:    MonthJune,
			weekday:  WeekdayWednesday,
			hour:     9,
			expected: time.Date(2020, 6, 24, 9, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLastWeekdayOfMonth(tt.year, tt.month, tt.weekday, tt.hour)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultTimezoneAbbrMap(t *testing.T) {
	// Test that common timezones are present
	assert.NotNil(t, DefaultTimezoneAbbrMap)
	assert.Equal(t, 0, DefaultTimezoneAbbrMap["UTC"])
	assert.Equal(t, 0, DefaultTimezoneAbbrMap["GMT"])
	assert.Equal(t, -300, DefaultTimezoneAbbrMap["EST"])
	assert.Equal(t, -480, DefaultTimezoneAbbrMap["PST"])
	assert.Equal(t, 540, DefaultTimezoneAbbrMap["JST"])
}

// TestTimezoneAbbreviations_Comprehensive tests comprehensive timezone abbreviations
// based on Python dateparser test suite
func TestTimezoneAbbreviations_Comprehensive(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		tz       string
		expected int // offset in minutes
	}{
		// UTC/GMT variants
		{"UTC", "UTC", 0},
		{"GMT", "GMT", 0},
		{"Z (Zulu)", "Z", 0},

		// North American - Eastern
		{"EST - Eastern Standard Time", "EST", -300},
		{"EDT - Eastern Daylight Time", "EDT", -240},

		// North American - Central
		{"CST - Central Standard Time", "CST", -360},
		{"CDT - Central Daylight Time", "CDT", -300},

		// North American - Mountain
		{"MST - Mountain Standard Time", "MST", -420},
		{"MDT - Mountain Daylight Time", "MDT", -360},

		// North American - Pacific
		{"PST - Pacific Standard Time", "PST", -480},
		{"PDT - Pacific Daylight Time", "PDT", -420},

		// North American - Alaska
		{"AKST - Alaska Standard Time", "AKST", -540},
		{"AKDT - Alaska Daylight Time", "AKDT", -480},

		// North American - Hawaii
		{"HST - Hawaii Standard Time", "HST", -600},
		{"HAST - Hawaii-Aleutian Standard Time", "HAST", -600},
		{"HADT - Hawaii-Aleutian Daylight Time", "HADT", -540},

		// European - Western
		{"WET - Western European Time", "WET", 0},
		{"WEST - Western European Summer Time", "WEST", 60},
		{"BST - British Summer Time", "BST", 60},
		{"IST - Irish Standard Time (used as BST equiv)", "IST", 60},

		// European - Central
		{"CEST - Central European Summer Time", "CEST", 120},

		// European - Eastern
		{"EET - Eastern European Time", "EET", 120},
		{"EEST - Eastern European Summer Time", "EEST", 180},

		// European - Moscow
		{"MSK - Moscow Standard Time", "MSK", 180},

		// Asian - East Asia
		{"JST - Japan Standard Time", "JST", 540},
		{"KST - Korean Standard Time", "KST", 540},
		{"HKT - Hong Kong Time", "HKT", 480},
		{"SGT - Singapore Time", "SGT", 480},
		{"CST_CHINA - China Standard Time", "CST_CHINA", 480},

		// Asian - South Asia
		{"IST_INDIA - India Standard Time", "IST_INDIA", 330},
		{"PKT - Pakistan Standard Time", "PKT", 300},

		// Asian - Southeast Asia
		{"WIB - Western Indonesia Time", "WIB", 420},
		{"WITA - Central Indonesia Time", "WITA", 480},
		{"WIT - Eastern Indonesia Time", "WIT", 540},

		// Australian
		{"AEST - Australian Eastern Standard Time", "AEST", 600},
		{"AEDT - Australian Eastern Daylight Time", "AEDT", 660},
		{"ACST - Australian Central Standard Time", "ACST", 570},
		{"ACDT - Australian Central Daylight Time", "ACDT", 630},
		{"AWST - Australian Western Standard Time", "AWST", 480},

		// New Zealand
		{"NZST - New Zealand Standard Time", "NZST", 720},
		{"NZDT - New Zealand Daylight Time", "NZDT", 780},

		// South American
		{"BRT - Brasilia Time", "BRT", -180},
		{"ART - Argentina Time", "ART", -180},

		// Other
		{"GET - Georgia Eastern Time", "GET", 240},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTimezoneOffset(tt.tz, instant, nil)
			assert.NotNil(t, result, "Timezone %s should be recognized", tt.tz)
			assert.Equal(t, tt.expected, *result, "Timezone %s offset mismatch", tt.tz)
		})
	}
}

// TestToTimezoneOffset_IANATimezones tests IANA timezone database names
func TestToTimezoneOffset_IANATimezones(t *testing.T) {
	tests := []struct {
		name             string
		location         string
		instant          time.Time
		expectedOffsetMi int // offset in minutes
	}{
		{
			name:             "America/New_York - winter",
			location:         "America/New_York",
			instant:          time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: -300, // EST
		},
		{
			name:             "America/New_York - summer",
			location:         "America/New_York",
			instant:          time.Date(2020, 7, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: -240, // EDT
		},
		{
			name:             "America/Los_Angeles - winter",
			location:         "America/Los_Angeles",
			instant:          time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: -480, // PST
		},
		{
			name:             "America/Los_Angeles - summer",
			location:         "America/Los_Angeles",
			instant:          time.Date(2020, 7, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: -420, // PDT
		},
		{
			name:             "Europe/London - winter",
			location:         "Europe/London",
			instant:          time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 0, // GMT
		},
		{
			name:             "Europe/London - summer",
			location:         "Europe/London",
			instant:          time.Date(2020, 7, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 60, // BST
		},
		{
			name:             "Asia/Tokyo - no DST",
			location:         "Asia/Tokyo",
			instant:          time.Date(2020, 7, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 540, // JST
		},
		{
			name:             "Australia/Sydney - summer",
			location:         "Australia/Sydney",
			instant:          time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 660, // AEDT
		},
		{
			name:             "Australia/Sydney - winter",
			location:         "Australia/Sydney",
			instant:          time.Date(2020, 7, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 600, // AEST
		},
		{
			name:             "Europe/Paris - winter",
			location:         "Europe/Paris",
			instant:          time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 60, // CET
		},
		{
			name:             "Europe/Paris - summer",
			location:         "Europe/Paris",
			instant:          time.Date(2020, 7, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 120, // CEST
		},
		{
			name:             "Asia/Kolkata - no DST",
			location:         "Asia/Kolkata",
			instant:          time.Date(2020, 7, 15, 12, 0, 0, 0, time.UTC),
			expectedOffsetMi: 330, // IST (India)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTimezoneOffset(tt.location, tt.instant, nil)
			assert.NotNil(t, result, "Location %s should be recognized", tt.location)
			assert.Equal(t, tt.expectedOffsetMi, *result, "Location %s offset mismatch", tt.location)
		})
	}
}

// TestToTimezoneOffset_EdgeCases tests edge cases and unusual scenarios
func TestToTimezoneOffset_EdgeCases(t *testing.T) {
	t.Run("negative numeric offset", func(t *testing.T) {
		instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
		offset := -480
		result := toTimezoneOffset(offset, instant, nil)
		assert.NotNil(t, result)
		assert.Equal(t, -480, *result)
	})

	t.Run("positive numeric offset", func(t *testing.T) {
		instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
		offset := 330
		result := toTimezoneOffset(offset, instant, nil)
		assert.NotNil(t, result)
		assert.Equal(t, 330, *result)
	})

	t.Run("nil timezone", func(t *testing.T) {
		instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
		result := toTimezoneOffset(nil, instant, nil)
		assert.Nil(t, result)
	})

	t.Run("unknown string timezone", func(t *testing.T) {
		instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
		result := toTimezoneOffset("INVALID_TZ", instant, nil)
		assert.Nil(t, result)
	})

	t.Run("case-sensitive timezone", func(t *testing.T) {
		instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
		// Lowercase should not match
		result := toTimezoneOffset("utc", instant, nil)
		assert.Nil(t, result)
	})
}

// TestToTimezoneOffset_OffsetFormats tests various offset format interpretations
func TestToTimezoneOffset_OffsetFormats(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		offset   int
		expected int
	}{
		{"zero offset", 0, 0},
		{"positive hours only", 300, 300},
		{"negative hours only", -300, -300},
		{"half hour offset positive", 330, 330},
		{"half hour offset negative", -330, -330},
		{"quarter hour offset", 345, 345},
		{"large positive offset", 720, 720},
		{"large negative offset", -720, -720},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTimezoneOffset(tt.offset, instant, nil)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

// TestToTimezoneOffset_DSTTransitions tests DST transition edge cases
func TestToTimezoneOffset_DSTTransitions(t *testing.T) {
	tests := []struct {
		name     string
		tz       string
		instant  time.Time
		expected int
	}{
		// ET transitions in 2020
		{
			name:     "ET before DST 2020",
			tz:       "ET",
			instant:  time.Date(2020, 3, 7, 12, 0, 0, 0, time.UTC),
			expected: -300, // EST
		},
		{
			name:     "ET after DST start 2020",
			tz:       "ET",
			instant:  time.Date(2020, 3, 9, 12, 0, 0, 0, time.UTC),
			expected: -240, // EDT
		},
		{
			name:     "ET before DST end 2020",
			tz:       "ET",
			instant:  time.Date(2020, 10, 31, 12, 0, 0, 0, time.UTC),
			expected: -240, // EDT
		},
		{
			name:     "ET after DST end 2020",
			tz:       "ET",
			instant:  time.Date(2020, 11, 2, 12, 0, 0, 0, time.UTC),
			expected: -300, // EST
		},
		// CET transitions in 2020
		{
			name:     "CET before DST 2020",
			tz:       "CET",
			instant:  time.Date(2020, 3, 28, 12, 0, 0, 0, time.UTC),
			expected: 60, // CET
		},
		{
			name:     "CET after DST start 2020",
			tz:       "CET",
			instant:  time.Date(2020, 3, 30, 12, 0, 0, 0, time.UTC),
			expected: 120, // CEST
		},
		{
			name:     "CET before DST end 2020",
			tz:       "CET",
			instant:  time.Date(2020, 10, 24, 12, 0, 0, 0, time.UTC),
			expected: 120, // CEST
		},
		{
			name:     "CET after DST end 2020",
			tz:       "CET",
			instant:  time.Date(2020, 10, 26, 12, 0, 0, 0, time.UTC),
			expected: 60, // CET
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTimezoneOffset(tt.tz, tt.instant, nil)
			assert.NotNil(t, result, "Timezone %s should be recognized", tt.tz)
			assert.Equal(t, tt.expected, *result, "Timezone %s offset mismatch at %s", tt.tz, tt.instant)
		})
	}
}

// TestToTimezoneOffset_OverridePrecedence tests that overrides take precedence
func TestToTimezoneOffset_OverridePrecedence(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)

	t.Run("override existing timezone", func(t *testing.T) {
		overrides := TimezoneAbbrMap{
			"UTC": 123, // Override UTC to non-zero
		}
		result := toTimezoneOffset("UTC", instant, overrides)
		assert.NotNil(t, result)
		assert.Equal(t, 123, *result, "Override should take precedence")
	})

	t.Run("override with new timezone", func(t *testing.T) {
		overrides := TimezoneAbbrMap{
			"CUSTOM": 456,
		}
		result := toTimezoneOffset("CUSTOM", instant, overrides)
		assert.NotNil(t, result)
		assert.Equal(t, 456, *result, "Custom timezone should work")
	})

	t.Run("default still works without override", func(t *testing.T) {
		overrides := TimezoneAbbrMap{
			"CUSTOM": 456,
		}
		result := toTimezoneOffset("EST", instant, overrides)
		assert.NotNil(t, result)
		assert.Equal(t, -300, *result, "Default timezone should still work")
	})
}

// TestToTimezoneOffset_AmbiguousDSTBoundaries tests exact DST transition boundaries
func TestToTimezoneOffset_AmbiguousDSTBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		tz       string
		instant  time.Time
		expected int
	}{
		// Test exact DST start for ET in 2021 (March 14, 2021 at 2 AM)
		{
			name:     "ET at exact DST start 2021",
			tz:       "ET",
			instant:  time.Date(2021, 3, 14, 2, 0, 0, 0, time.UTC),
			expected: -300, // Should still be EST before transition
		},
		{
			name:     "ET after DST start 2021",
			tz:       "ET",
			instant:  time.Date(2021, 3, 14, 3, 0, 0, 0, time.UTC),
			expected: -240, // Should be EDT after transition
		},
		// Test exact DST end for ET in 2021 (November 7, 2021 at 2 AM)
		{
			name:     "ET before DST end 2021",
			tz:       "ET",
			instant:  time.Date(2021, 11, 7, 1, 0, 0, 0, time.UTC),
			expected: -240, // Should still be EDT before transition
		},
		{
			name:     "ET after DST end 2021",
			tz:       "ET",
			instant:  time.Date(2021, 11, 7, 8, 0, 0, 0, time.UTC),
			expected: -300, // Should be EST after transition
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTimezoneOffset(tt.tz, tt.instant, nil)
			assert.NotNil(t, result, "Timezone %s should be recognized", tt.tz)
			assert.Equal(t, tt.expected, *result, "Timezone %s offset mismatch at boundary %s", tt.tz, tt.instant)
		})
	}
}
