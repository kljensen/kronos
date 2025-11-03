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
	result := ToTimezoneOffset(offset, instant, nil)
	assert.NotNil(t, result)
	assert.Equal(t, 300, *result)

	// Test negative offset
	offset = -480
	result = ToTimezoneOffset(offset, instant, nil)
	assert.NotNil(t, result)
	assert.Equal(t, -480, *result)

	// Test zero offset
	offset = 0
	result = ToTimezoneOffset(offset, instant, nil)
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
			result := ToTimezoneOffset(tt.tz, instant, nil)
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
	result := ToTimezoneOffset("CUSTOM", instant, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, 123, *result)

	// Test override taking precedence over default
	overrides["UTC"] = 999
	result = ToTimezoneOffset("UTC", instant, overrides)
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
			return GetNthWeekdayOfMonth(year, MonthMarch, WeekdaySunday, 2, 2)
		},
		DstEnd: func(year int) time.Time {
			// DST ends 1st Sunday of November at 2 AM
			return GetNthWeekdayOfMonth(year, MonthNovember, WeekdaySunday, 1, 2)
		},
	}

	overrides := TimezoneAbbrMap{
		"ET": et,
	}

	// Test during DST (June 1, 2020)
	instantDuringDst := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
	result := ToTimezoneOffset("ET", instantDuringDst, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -240, *result, "Should use DST offset during summer")

	// Test during non-DST (January 1, 2020)
	instantNonDst := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	result = ToTimezoneOffset("ET", instantNonDst, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -300, *result, "Should use non-DST offset during winter")

	// Test right at DST start (should be non-DST)
	instantAtDstStart := time.Date(2020, 3, 8, 1, 0, 0, 0, time.UTC)
	result = ToTimezoneOffset("ET", instantAtDstStart, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -300, *result, "Should use non-DST offset before DST starts")

	// Test right after DST start (should be DST)
	instantAfterDstStart := time.Date(2020, 3, 8, 8, 0, 0, 0, time.UTC)
	result = ToTimezoneOffset("ET", instantAfterDstStart, overrides)
	assert.NotNil(t, result)
	assert.Equal(t, -240, *result, "Should use DST offset after DST starts")
}

func TestToTimezoneOffset_NilTimezone(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
	result := ToTimezoneOffset(nil, instant, nil)
	assert.Nil(t, result)
}

func TestToTimezoneOffset_UnknownTimezone(t *testing.T) {
	instant := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
	result := ToTimezoneOffset("UNKNOWN_TZ", instant, nil)
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
	result := ToTimezoneOffset("ET", time.Time{}, overrides)
	assert.Nil(t, result)
}

func TestGetNthWeekdayOfMonth(t *testing.T) {
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
			result := GetNthWeekdayOfMonth(tt.year, tt.month, tt.weekday, tt.n, tt.hour)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetLastWeekdayOfMonth(t *testing.T) {
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
			result := GetLastWeekdayOfMonth(tt.year, tt.month, tt.weekday, tt.hour)
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
