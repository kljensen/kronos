//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

// Tests ported from: github.com/olebedev/when - rules/en/en_test.go
// License: Apache 2.0
// Source: https://github.com/olebedev/when
//
// The 'when' library is a natural language date/time parser with pluggable rules.
// These tests validate parity with 'when' parsing behavior for English expressions.
//
// Reference date for tests: January 6, 2016 00:00:00 UTC (Wednesday)

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// Helper to create pointer to int (if not already defined)
// intPtr is used throughout test files

// TestWhenParity_AfternoonParsing tests parsing of afternoon time expressions
// Ported from: github.com/olebedev/when - rules/en/en_test.go
func TestWhenParity_AfternoonParsing(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedText  string
		expectedYear  int
		expectedMonth int
		expectedDay   int
		expectedHour  *int
		skip          bool
		skipReason    string
	}{
		{
			name:          "tonight at 11:10 pm",
			text:          "tonight at 11:10 pm",
			refDate:       time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC),
			expectedText:  "tonight at 11:10 pm",
			expectedYear:  2016,
			expectedMonth: 1,
			expectedDay:   6,
			expectedHour:  intPtr(23),
			skip:          true,
			skipReason:    "SKIP: Combined 'tonight at TIME' expressions need parser merging/refiner - see issue #112",
		},
		{
			name:          "at Friday afternoon",
			text:          "at Friday afternoon",
			refDate:       time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC),
			expectedText:  "Friday afternoon",
			expectedYear:  2016,
			expectedMonth: 1,
			expectedDay:   8, // Next Friday from Wednesday Jan 6
			expectedHour:  intPtr(15),
			skip:          true,
			skipReason:    "SKIP: Combined weekday + time of day (afternoon) needs refiner - see issue #112",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			chrono := createFullChrono()
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch")

			if result.Start() != nil {
				assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
				assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
				assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
				if tt.expectedHour != nil {
					assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch")
				}
			}
		})
	}
}

// TestWhenParity_NextDayWithTime tests next day/hour patterns with specific times
// Ported from: github.com/olebedev/when - rules/en/en_test.go
func TestWhenParity_NextDayWithTime(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		refDate        time.Time
		expectedText   string
		expectedYear   int
		expectedMonth  int
		expectedDay    int
		expectedHour   *int
		expectedMinute *int
		skip           bool
		skipReason     string
	}{
		{
			name:           "next tuesday at 14:00",
			text:           "in next tuesday at 14:00",
			refDate:        time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC),
			expectedText:   "next tuesday at 14:00",
			expectedYear:   2016,
			expectedMonth:  1,
			expectedDay:    12, // Next Tuesday from Wednesday Jan 6
			expectedHour:   intPtr(14),
			expectedMinute: intPtr(0),
			skip:           true,
			skipReason:     "SKIP: Combined 'next WEEKDAY at TIME' needs refiner to merge date+time - see issue #112",
		},
		{
			name:           "next tuesday at 2p",
			text:           "in next tuesday at 2p",
			refDate:        time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC),
			expectedText:   "next tuesday at 2p",
			expectedYear:   2016,
			expectedMonth:  1,
			expectedDay:    12,
			expectedHour:   intPtr(14),
			expectedMinute: intPtr(0),
			skip:           true,
			skipReason:     "SKIP: Combined 'next WEEKDAY at TIME' needs refiner to merge date+time - see issue #112",
		},
		{
			name:           "next wednesday at 2:25 p.m.",
			text:           "in next wednesday at 2:25 p.m.",
			refDate:        time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC),
			expectedText:   "next wednesday at 2:25 p.m.",
			expectedYear:   2016,
			expectedMonth:  1,
			expectedDay:    13, // Next Wednesday from Wednesday Jan 6
			expectedHour:   intPtr(14),
			expectedMinute: intPtr(25),
			skip:           true,
			skipReason:     "SKIP: Combined 'next WEEKDAY at TIME' needs refiner to merge date+time - see issue #112",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			chrono := createFullChrono()
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch")

			if result.Start() != nil {
				assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
				assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
				assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
				if tt.expectedHour != nil {
					assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch")
				}
				if tt.expectedMinute != nil {
					assert.Equal(t, *tt.expectedMinute, *result.Start().Get(kronos.ComponentMinute), "Minute mismatch")
				}
			}
		})
	}
}

// TestWhenParity_PastReferences tests past temporal references
// Ported from: github.com/olebedev/when - rules/en/en_test.go
func TestWhenParity_PastReferences(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedText  string
		expectedYear  int
		expectedMonth int
		expectedDay   int
		expectedHour  *int
		skip          bool
		skipReason    string
	}{
		{
			name:          "11 am past tuesday",
			text:          "at 11 am past tuesday",
			refDate:       time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC),
			expectedText:  "11 am past tuesday",
			expectedYear:  2016,
			expectedMonth: 1,
			expectedDay:   5, // Last Tuesday from Wednesday Jan 6
			expectedHour:  intPtr(11),
			skip:          true,
			skipReason:    "SKIP: Combined 'TIME past WEEKDAY' needs refiner to merge time+date - see issue #112",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			chrono := createFullChrono()
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch")

			if result.Start() != nil {
				assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
				assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
				assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
				if tt.expectedHour != nil {
					assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch")
				}
			}
		})
	}
}

// TestWhenParity_DistanceClustering tests distance-based rule clustering
// This feature in 'when' library groups nearby temporal patterns
func TestWhenParity_DistanceClustering(t *testing.T) {
	t.Skip("SKIP: Distance-based clustering of matches not implemented in kronos - this is a 'when' library specific feature")

	// The 'when' library has a 'distance' parameter (default: 5) that determines
	// which rule matches should be clustered together. This affects how closely
	// positioned temporal expressions are grouped into single results.
	//
	// Example: "tomorrow at 3pm" might be parsed as two separate matches
	// ("tomorrow" and "3pm") that get clustered based on their proximity.
	//
	// Kronos uses a different approach with parsers and refiners, so this
	// specific clustering behavior may differ.
}

// createFullChrono creates a Chrono instance with full English parser suite
func createFullChrono() *kronos.Chrono {
	config := &kronos.Configuration{
		Parsers: []kronos.Parser{
			NewENCasualDateParser(),
			NewENCasualTimeParser(),
			NewENWeekdayParser(),
			NewENTimeExpressionParser(false),
			NewENMonthNameParser(),
		},
	}
	return kronos.NewChrono(config)
}
