package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// TestNoonMidnightKeywords tests the noon and midnight keyword support
// as requested in issue #82
func TestNoonMidnightKeywords(t *testing.T) {
	tests := []struct {
		name             string
		text             string
		refDate          time.Time
		expectedText     string
		expectedYear     int
		expectedMonth    int
		expectedDay      int
		expectedHour     int
		expectedMinute   int
		expectedSecond   int
		expectedMillisec int
	}{
		// Test cases from Python dateparser (issue #82)
		{
			name:           "November 19, 2014 at noon",
			text:           "November 19, 2014 at noon",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "November 19, 2014 at noon",
			expectedYear:   2014,
			expectedMonth:  11,
			expectedDay:    19,
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "December 13, 2014 at midnight",
			text:           "December 13, 2014 at midnight",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "December 13, 2014 at midnight",
			expectedYear:   2014,
			expectedMonth:  12,
			expectedDay:    13,
			expectedHour:   0,
			expectedMinute: 0,
			expectedSecond: 0,
		},

		// Standalone keywords
		{
			name:           "at noon",
			text:           "at noon",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "noon",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    1,
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "at midnight (afternoon ref)",
			text:           "at midnight",
			refDate:        time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
			expectedText:   "midnight",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    2, // Next day's midnight when ref is afternoon
			expectedHour:   0,
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "at midnight (early morning ref)",
			text:           "at midnight",
			refDate:        time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC),
			expectedText:   "midnight",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    1, // Same day when ref is early morning
			expectedHour:   0,
			expectedMinute: 0,
			expectedSecond: 0,
		},

		// Combined with relative dates
		{
			name:           "noon tomorrow",
			text:           "noon tomorrow",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "noon tomorrow",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    2,
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "midnight tonight",
			text:           "midnight tonight",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "midnight tonight",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    1,
			expectedHour:   0, // Midnight should be 00:00, not 12:00
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "yesterday at noon",
			text:           "yesterday at noon",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "yesterday at noon",
			expectedYear:   2022,
			expectedMonth:  12,
			expectedDay:    31,
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},

		// Combined with weekdays
		{
			name:           "next Tuesday at midnight",
			text:           "next Tuesday at midnight",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC), // Sunday
			expectedText:   "next Tuesday at midnight",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    3, // Tuesday
			expectedHour:   0,
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "next Friday at noon",
			text:           "next Friday at noon",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC), // Sunday
			expectedText:   "next Friday at noon",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    6, // Friday
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},

		// Case variations
		{
			name:           "NOON (uppercase)",
			text:           "Meet at NOON",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "NOON",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    1,
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "Midnight (capitalized)",
			text:           "Meet at Midnight",
			refDate:        time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
			expectedText:   "Midnight",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    2,
			expectedHour:   0,
			expectedMinute: 0,
			expectedSecond: 0,
		},

		// Without "at" prefix
		{
			name:           "noon without at",
			text:           "Meet noon",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "noon",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    1,
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},
		{
			name:           "midnight without at",
			text:           "Meet midnight",
			refDate:        time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
			expectedText:   "midnight",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    2,
			expectedHour:   0,
			expectedMinute: 0,
			expectedSecond: 0,
		},

		// Alternative form: midday
		{
			name:           "midday (synonym for noon)",
			text:           "Meet at midday",
			refDate:        time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedText:   "midday",
			expectedYear:   2023,
			expectedMonth:  1,
			expectedDay:    1,
			expectedHour:   12,
			expectedMinute: 0,
			expectedSecond: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, _ := New().WithReferenceDate(tt.refDate).Parse(tt.text)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.expectedText, result.Text(), "Matched text mismatch")

			start := result.Start()
			assert.Equal(t, tt.expectedYear, *start.Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *start.Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *start.Get(kronos.ComponentDay), "Day mismatch")
			assert.Equal(t, tt.expectedHour, *start.Get(kronos.ComponentHour), "Hour mismatch")
			assert.Equal(t, tt.expectedMinute, *start.Get(kronos.ComponentMinute), "Minute mismatch")

			// Verify the full datetime
			expectedDate := time.Date(tt.expectedYear, time.Month(tt.expectedMonth), tt.expectedDay,
				tt.expectedHour, tt.expectedMinute, tt.expectedSecond, tt.expectedMillisec*1000000, time.UTC)
			assert.Equal(t, expectedDate, result.Start().Date(), "Full date mismatch")
		})
	}
}

// TestNoonMidnightEdgeCases tests edge cases and potential issues
func TestNoonMidnightEdgeCases(t *testing.T) {
	t.Run("midnight should be hour 0, not hour 12", func(t *testing.T) {
		refDate := time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC)
		results, _ := New().WithReferenceDate(refDate).Parse("midnight")

		assert.NotEmpty(t, results)
		result := results[0]
		hour := result.Start().Get(kronos.ComponentHour)
		assert.NotNil(t, hour)
		assert.Equal(t, 0, *hour, "Midnight should be hour 0, not 12")
	})

	t.Run("noon should be hour 12", func(t *testing.T) {
		refDate := time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC)
		results, _ := New().WithReferenceDate(refDate).Parse("noon")

		assert.NotEmpty(t, results)
		result := results[0]
		hour := result.Start().Get(kronos.ComponentHour)
		assert.NotNil(t, hour)
		assert.Equal(t, 12, *hour, "Noon should be hour 12")
	})

	t.Run("midnight tonight should not be converted to 12 PM", func(t *testing.T) {
		// This is the critical test for the bug fix
		refDate := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
		results, _ := New().WithReferenceDate(refDate).Parse("midnight tonight")

		assert.NotEmpty(t, results)
		result := results[0]
		hour := result.Start().Get(kronos.ComponentHour)
		assert.NotNil(t, hour)
		assert.Equal(t, 0, *hour, "Midnight tonight should be hour 0, not 12")

		// Verify the date is correct too
		expectedDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, expectedDate, result.Start().Date())
	})

	t.Run("midnight with explicit date should work", func(t *testing.T) {
		refDate := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
		results, _ := New().WithReferenceDate(refDate).Parse("at midnight on August 12")

		assert.NotEmpty(t, results)
		result := results[0]
		assert.Equal(t, 0, *result.Start().Get(kronos.ComponentHour))
		assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
		assert.Equal(t, 12, *result.Start().Get(kronos.ComponentDay))
	})
}

// TestNoonMidnightWithTimezones tests noon/midnight with timezone expressions
func TestNoonMidnightWithTimezones(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		refDate      time.Time
		expectedHour int
		expectedTZ   *int
	}{
		{
			name:         "noon EST",
			text:         "noon EST",
			refDate:      time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedHour: 12,
			expectedTZ:   intPtr(-300), // EST is UTC-5
		},
		{
			name:         "midnight PST",
			text:         "midnight PST",
			refDate:      time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedHour: 0,
			expectedTZ:   intPtr(-480), // PST is UTC-8
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, _ := New().WithReferenceDate(tt.refDate).Parse(tt.text)

			assert.NotEmpty(t, results)
			result := results[0]
			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))

			if tt.expectedTZ != nil {
				tz := result.Start().Get(kronos.ComponentTimezoneOffset)
				assert.NotNil(t, tz, "Expected timezone to be parsed")
				assert.Equal(t, *tt.expectedTZ, *tz, "Timezone offset mismatch")
			}
		})
	}
}
