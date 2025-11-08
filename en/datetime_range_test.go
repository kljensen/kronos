package en

// TestDateTimeRangePreservation tests that time ranges are preserved when
// combined with date expressions.
//
// This test file addresses GitHub issue #150: "Time ranges within date
// expressions are not preserved". When a time range like "8:30 PM - 10:45 PM"
// is combined with a date expression, both start and end times should be
// preserved.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDateTimeRangeP0CoreCases tests the core time range preservation scenarios.
// These are the most critical cases (P0) that must work.
func TestDateTimeRangeP0CoreCases(t *testing.T) {
	// Reference date for all tests
	refDate := time.Date(2024, 11, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		text      string
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "Date with PM time range",
			text:      "Wed, Nov 5, 2025 8:30 PM - 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
		{
			name:      "Date with AM time range",
			text:      "Wed, Nov 5, 2025 8:30 AM - 10:45 AM",
			wantStart: time.Date(2025, 11, 5, 8, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 10, 45, 0, 0, time.UTC),
		},
		{
			name:      "Date with mixed AM/PM time range",
			text:      "Wed, Nov 5, 2025 11:30 AM - 2:45 PM",
			wantStart: time.Date(2025, 11, 5, 11, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 14, 45, 0, 0, time.UTC),
		},
		{
			name:      "Date with 24-hour time range",
			text:      "Wed, Nov 5, 2025 14:30 - 16:45",
			wantStart: time.Date(2025, 11, 5, 14, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 16, 45, 0, 0, time.UTC),
		},
		{
			name:      "ISO date with time range",
			text:      "2025-11-05 8:30 PM - 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := New().WithReferenceDate(refDate).Parse(tt.text)
			require.NoError(t, err, "Parse should not return error")
			require.NotEmpty(t, results, "Should parse at least one result")

			result := results[0]

			// Verify start time
			startDate := result.Start().Date()
			assert.Equal(t, tt.wantStart.Year(), startDate.Year())
			assert.Equal(t, tt.wantStart.Month(), startDate.Month())
			assert.Equal(t, tt.wantStart.Day(), startDate.Day())
			assert.Equal(t, tt.wantStart.Hour(), startDate.Hour())
			assert.Equal(t, tt.wantStart.Minute(), startDate.Minute())

			// Verify end time is present
			require.NotNil(t, result.End(), "End time should be present for time range")

			// Verify end time values
			endDate := result.End().Date()
			assert.Equal(t, tt.wantEnd.Year(), endDate.Year())
			assert.Equal(t, tt.wantEnd.Month(), endDate.Month())
			assert.Equal(t, tt.wantEnd.Day(), endDate.Day())
			assert.Equal(t, tt.wantEnd.Hour(), endDate.Hour())
			assert.Equal(t, tt.wantEnd.Minute(), endDate.Minute())
		})
	}
}

// TestDateTimeRangeP1SeparatorVariations tests different separator styles.
// These are important but not critical (P1).
func TestDateTimeRangeP1SeparatorVariations(t *testing.T) {
	refDate := time.Date(2024, 11, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		text      string
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "Hyphen separator with spaces",
			text:      "Nov 5, 2025 8:30 PM - 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
		{
			name:      "Hyphen separator without spaces",
			text:      "Nov 5, 2025 8:30 PM-10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
		{
			name:      "To separator",
			text:      "Nov 5, 2025 8:30 PM to 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
		{
			name:      "Double hyphen separator",
			text:      "Nov 5, 2025 8:30 PM -- 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
		{
			name:      "Triple hyphen separator",
			text:      "Nov 5, 2025 8:30 PM --- 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
		{
			name:      "En-dash separator",
			text:      "Nov 5, 2025 8:30 PM – 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
		{
			name:      "Em-dash separator",
			text:      "Nov 5, 2025 8:30 PM — 10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := New().WithReferenceDate(refDate).Parse(tt.text)
			require.NoError(t, err, "Parse should not return error")
			require.NotEmpty(t, results, "Should parse at least one result")

			result := results[0]

			// Verify start time
			startDate := result.Start().Date()
			assert.Equal(t, tt.wantStart.Hour(), startDate.Hour())
			assert.Equal(t, tt.wantStart.Minute(), startDate.Minute())

			// Verify end time is present
			require.NotNil(t, result.End(), "End time should be present for time range")

			// Verify end time values
			endDate := result.End().Date()
			assert.Equal(t, tt.wantEnd.Hour(), endDate.Hour())
			assert.Equal(t, tt.wantEnd.Minute(), endDate.Minute())
		})
	}
}

// TestDateTimeRangeP2EdgeCases tests edge cases and complex scenarios.
// These are lower priority but still important (P2).
func TestDateTimeRangeP2EdgeCases(t *testing.T) {
	refDate := time.Date(2024, 11, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		text      string
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "Overnight time range (crosses midnight)",
			text:      "Nov 5, 2025 11:30 PM - 1:30 AM",
			wantStart: time.Date(2025, 11, 5, 23, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 6, 1, 30, 0, 0, time.UTC),
		},
		{
			name:      "Time range with seconds",
			text:      "Nov 5, 2025 8:30:15 PM - 10:45:30 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 15, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 30, 0, time.UTC),
		},
		{
			name:      "Multiple spaces in separator",
			text:      "Nov 5, 2025 8:30 PM  -  10:45 PM",
			wantStart: time.Date(2025, 11, 5, 20, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 11, 5, 22, 45, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := New().WithReferenceDate(refDate).Parse(tt.text)
			require.NoError(t, err, "Parse should not return error")
			require.NotEmpty(t, results, "Should parse at least one result")

			result := results[0]

			// Verify start time
			startDate := result.Start().Date()
			assert.Equal(t, tt.wantStart.Hour(), startDate.Hour())
			assert.Equal(t, tt.wantStart.Minute(), startDate.Minute())
			assert.Equal(t, tt.wantStart.Second(), startDate.Second())

			// Verify end time is present
			require.NotNil(t, result.End(), "End time should be present for time range")

			// Verify end time values
			endDate := result.End().Date()
			assert.Equal(t, tt.wantEnd.Hour(), endDate.Hour())
			assert.Equal(t, tt.wantEnd.Minute(), endDate.Minute())
			assert.Equal(t, tt.wantEnd.Second(), endDate.Second())

			// For overnight ranges, verify day increment
			if tt.wantEnd.Day() > tt.wantStart.Day() {
				assert.Equal(t, tt.wantEnd.Day(), endDate.Day(),
					"End day should be incremented for overnight range")
			}
		})
	}
}
