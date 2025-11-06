//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/experimental"
	"github.com/kljensen/kronos/internal/helpers"
	"github.com/stretchr/testify/assert"
)

// TestMicrosecondPrecisionInTimestamps tests parsing timestamps with microsecond precision
// as documented in issue #83 and based on Python's dateparser test cases
func TestMicrosecondPrecisionInTimestamps(t *testing.T) {
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		input          string
		expectedYear   int
		expectedMonth  int
		expectedDay    int
		expectedHour   int
		expectedMinute int
		expectedSecond int
		expectedNanos  int
	}{
		{
			name:           "1 digit fractional seconds (.1)",
			input:          "1/1/16 9:02:43.1",
			expectedYear:   2016,
			expectedMonth:  1,
			expectedDay:    1,
			expectedHour:   9,
			expectedMinute: 2,
			expectedSecond: 43,
			expectedNanos:  100000000, // 0.1 seconds = 100 milliseconds = 100000000 nanoseconds
		},
		{
			name:           "2 digits fractional seconds (.12)",
			input:          "2012-01-21 13:11:23.12",
			expectedYear:   2012,
			expectedMonth:  1,
			expectedDay:    21,
			expectedHour:   13,
			expectedMinute: 11,
			expectedSecond: 23,
			expectedNanos:  120000000, // 0.12 seconds = 120 milliseconds
		},
		{
			name:           "3 digits fractional seconds (.678)",
			input:          "21 January 2012 13:11:23.678",
			expectedYear:   2012,
			expectedMonth:  1,
			expectedDay:    21,
			expectedHour:   13,
			expectedMinute: 11,
			expectedSecond: 23,
			expectedNanos:  678000000, // 0.678 seconds = 678 milliseconds
		},
		{
			name:           "6 digits fractional seconds (.123456)",
			input:          "2012-01-21 13:11:23.123456",
			expectedYear:   2012,
			expectedMonth:  1,
			expectedDay:    21,
			expectedHour:   13,
			expectedMinute: 11,
			expectedSecond: 23,
			expectedNanos:  123456000, // 0.123456 seconds = 123456 microseconds
		},
		{
			name:           "9 digits fractional seconds (.123456789)",
			input:          "2012-01-21 13:11:23.123456789",
			expectedYear:   2012,
			expectedMonth:  1,
			expectedDay:    21,
			expectedHour:   13,
			expectedMinute: 11,
			expectedSecond: 23,
			expectedNanos:  123456789, // Full nanosecond precision
		},
		{
			name:           "ISO format with 4 digits (.4876)",
			input:          "2016-05-07T23:45:00.4876",
			expectedYear:   2016,
			expectedMonth:  5,
			expectedDay:    7,
			expectedHour:   23,
			expectedMinute: 45,
			expectedSecond: 0,
			expectedNanos:  487600000, // 0.4876 seconds = 487.6 milliseconds
		},
		{
			name:           "Time with AM and fractional seconds",
			input:          "9:02:43.1 AM",
			expectedYear:   2012,
			expectedMonth:  8,
			expectedDay:    10,
			expectedHour:   9,
			expectedMinute: 2,
			expectedSecond: 43,
			expectedNanos:  100000000,
		},
		{
			name:           "Time with PM and fractional seconds",
			input:          "1:30:15.5 PM",
			expectedYear:   2012,
			expectedMonth:  8,
			expectedDay:    10,
			expectedHour:   13,
			expectedMinute: 30,
			expectedSecond: 15,
			expectedNanos:  500000000,
		},
		{
			name:           "24-hour format with microseconds",
			input:          "15:04:05.123456",
			expectedYear:   2012,
			expectedMonth:  8,
			expectedDay:    10,
			expectedHour:   15,
			expectedMinute: 4,
			expectedSecond: 5,
			expectedNanos:  123456000,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			chrono := experimental.EnglishCasualChrono()
			results := chrono.Parse(tt.input, refDate, nil)

			assert.NotEmpty(t, results, "Should parse: %s", tt.input)
			if len(results) == 0 {
				return
			}

			result := results[0]
			start := result.Start()

			if tt.expectedYear > 0 {
				assert.Equal(t, tt.expectedYear, *start.Get(kronos.ComponentYear), "Year mismatch")
			}
			if tt.expectedMonth > 0 {
				assert.Equal(t, tt.expectedMonth, *start.Get(kronos.ComponentMonth), "Month mismatch")
			}
			if tt.expectedDay > 0 {
				assert.Equal(t, tt.expectedDay, *start.Get(kronos.ComponentDay), "Day mismatch")
			}

			assert.Equal(t, tt.expectedHour, *start.Get(kronos.ComponentHour), "Hour mismatch")
			assert.Equal(t, tt.expectedMinute, *start.Get(kronos.ComponentMinute), "Minute mismatch")
			assert.Equal(t, tt.expectedSecond, *start.Get(kronos.ComponentSecond), "Second mismatch")

			// Check nanoseconds
			actualDate := result.Date()
			assert.Equal(t, tt.expectedNanos, actualDate.Nanosecond(), "Nanosecond mismatch")
		})
	}
}

// TestRelativeTimeWithMicroseconds tests relative time expressions using microseconds
func TestRelativeTimeWithMicroseconds(t *testing.T) {
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name          string
		input         string
		expectedNanos int
	}{
		{
			name:          "2 microseconds ago",
			input:         "2 microseconds ago",
			expectedNanos: -2000, // -2 microseconds in nanoseconds
		},
		{
			name:          "500 nanoseconds ago",
			input:         "500 nanoseconds ago",
			expectedNanos: -500,
		},
		{
			name:          "2.5 microseconds ago",
			input:         "2.5 microseconds ago",
			expectedNanos: -2500, // -2.5 microseconds in nanoseconds
		},
		{
			name:          "1 millisecond ago",
			input:         "1 millisecond ago",
			expectedNanos: -1000000, // -1 millisecond in nanoseconds
		},
		{
			name:          "250 milliseconds ago",
			input:         "250 milliseconds ago",
			expectedNanos: -250000000,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			chrono := experimental.EnglishCasualChrono()
			results := chrono.Parse(tt.input, refDate, nil)

			assert.NotEmpty(t, results, "Should parse: %s", tt.input)
			if len(results) == 0 {
				return
			}

			result := results[0]
			actualDate := result.Date()
			expectedDate := refDate.Add(time.Duration(tt.expectedNanos) * time.Nanosecond)

			assert.Equal(t, expectedDate, actualDate, "Date mismatch for: %s", tt.input)
		})
	}
}

// TestMicrosecondInDurationCalculations tests microsecond precision in duration calculations
func TestMicrosecondInDurationCalculations(t *testing.T) {
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name             string
		duration         kronos.Duration
		expectedDuration time.Duration
	}{
		{
			name: "100 microseconds",
			duration: kronos.Duration{
				kronos.TimeunitMicrosecond: 100,
			},
			expectedDuration: 100 * time.Microsecond,
		},
		{
			name: "1000 nanoseconds",
			duration: kronos.Duration{
				kronos.TimeunitNanosecond: 1000,
			},
			expectedDuration: 1000 * time.Nanosecond,
		},
		{
			name: "1.5 milliseconds",
			duration: kronos.Duration{
				kronos.TimeunitMillisecond: 1.5,
			},
			expectedDuration: 1500 * time.Microsecond,
		},
		{
			name: "Mixed: 1 millisecond + 500 microseconds",
			duration: kronos.Duration{
				kronos.TimeunitMillisecond: 1,
				kronos.TimeunitMicrosecond: 500,
			},
			expectedDuration: 1500 * time.Microsecond,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helpers.AddDuration(refDate, tt.duration)
			assert.NoError(t, err)
			expected := refDate.Add(tt.expectedDuration)

			assert.Equal(t, expected, result, "Duration calculation mismatch")
		})
	}
}

// TestTrailingZerosInFractionalSeconds tests that trailing zeros are handled correctly
func TestTrailingZerosInFractionalSeconds(t *testing.T) {
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name          string
		input1        string // With trailing zeros
		input2        string // Without trailing zeros
		expectedNanos int
	}{
		{
			name:          "0.1 vs 0.100",
			input1:        "12:00:00.100",
			input2:        "12:00:00.1",
			expectedNanos: 100000000,
		},
		{
			name:          "0.123 vs 0.123000",
			input1:        "12:00:00.123000",
			input2:        "12:00:00.123",
			expectedNanos: 123000000,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			chrono := experimental.EnglishCasualChrono()

			results1 := chrono.Parse(tt.input1, refDate, nil)
			assert.NotEmpty(t, results1)

			results2 := chrono.Parse(tt.input2, refDate, nil)
			assert.NotEmpty(t, results2)

			date1 := results1[0].Date()
			date2 := results2[0].Date()

			assert.Equal(t, tt.expectedNanos, date1.Nanosecond(), "Input1 nanoseconds mismatch")
			assert.Equal(t, tt.expectedNanos, date2.Nanosecond(), "Input2 nanoseconds mismatch")
			assert.Equal(t, date1, date2, "Both inputs should parse to the same time")
		})
	}
}
