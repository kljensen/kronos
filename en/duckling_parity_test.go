//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

// Tests ported from: github.com/facebook/duckling - Duckling/Time/EN/Corpus.hs
// License: BSD-3-Clause
// Source: https://github.com/facebook/duckling
//
// Duckling is a Haskell library that parses natural language into structured data.
// It uses corpus-based testing with both positive and negative examples.
//
// Reference date for all tests: Tuesday, February 12, 2013 at 4:30:00 AM

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/en/parsers"
	"github.com/stretchr/testify/assert"
)

// Standard reference date used by Duckling corpus
var ducklingRefDate = time.Date(2013, 2, 12, 4, 30, 0, 0, time.UTC)

// intPtr is defined in time_expression_test.go and reused here

// TestDucklingParity_RelativeTime tests relative time expressions
// Ported from: Duckling/Time/EN/Corpus.hs
func TestDucklingParity_RelativeTime(t *testing.T) {
	tests := []ducklingTestCase{
		{
			name:          "now",
			text:          "now",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   12,
			expectedHour:  intPtr(4),
			expectedMin:   intPtr(30),
		},
		{
			name:          "right now",
			text:          "right now",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   12,
			expectedHour:  intPtr(4),
			expectedMin:   intPtr(30),
		},
		{
			name:          "just now",
			text:          "just now",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   12,
			expectedHour:  intPtr(4),
			expectedMin:   intPtr(30),
		},
		{
			name:          "today",
			text:          "today",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   12,
		},
		{
			name:          "at this time",
			text:          "at this time",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   12,
			expectedHour:  intPtr(4),
			expectedMin:   intPtr(30),
			skip:          true,
			skipReason:    "SKIP: 'at this time' expression not implemented - see issue #112",
		},
		{
			name:          "yesterday",
			text:          "yesterday",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   11,
		},
		{
			name:          "tomorrow",
			text:          "tomorrow",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDucklingTest(t, tt)
		})
	}
}

// TestDucklingParity_DayOfWeek tests day of week expressions
// Ported from: Duckling/Time/EN/Corpus.hs
func TestDucklingParity_DayOfWeek(t *testing.T) {
	tests := []ducklingTestCase{
		{
			name:          "monday",
			text:          "monday",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   18, // Next Monday from Tuesday Feb 12
			skip:          true,
			skipReason:    "SKIP: Bare weekday defaults to 'this' not 'next' - behavior difference vs Duckling",
		},
		{
			name:          "mon.",
			text:          "mon.",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   18,
			skip:          true,
			skipReason:    "SKIP: Bare weekday defaults to 'this' not 'next' - behavior difference vs Duckling",
		},
		{
			name:          "this monday",
			text:          "this monday",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   18,
		},
		{
			name:          "next tuesday",
			text:          "next tuesday",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   19,
		},
		{
			name:          "last sunday",
			text:          "last sunday",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   10, // Last Sunday before Tuesday Feb 12
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDucklingTest(t, tt)
		})
	}
}

// TestDucklingParity_SpecificDates tests specific date expressions
// Ported from: Duckling/Time/EN/Corpus.hs
func TestDucklingParity_SpecificDates(t *testing.T) {
	tests := []ducklingTestCase{
		{
			name:          "2/15",
			text:          "2/15",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   15,
			skip:          true,
			skipReason:    "SKIP: Ambiguous slash dates without year may parse differently - needs investigation",
		},
		{
			name:          "march 3 2015",
			text:          "march 3 2015",
			expectedYear:  2015,
			expectedMonth: 3,
			expectedDay:   3,
			skip:          true,
			skipReason:    "SKIP: 'month day year' format parsing issue - likely parser ordering",
		},
		{
			name:          "the 1st of march",
			text:          "the 1st of march",
			expectedYear:  2013,
			expectedMonth: 3,
			expectedDay:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDucklingTest(t, tt)
		})
	}
}

// TestDucklingParity_TimeOfDay tests time expressions
// Ported from: Duckling/Time/EN/Corpus.hs
func TestDucklingParity_TimeOfDay(t *testing.T) {
	tests := []ducklingTestCase{
		{
			name:          "at 3am",
			text:          "at 3am",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   13, // Forward to next occurrence
			expectedHour:  intPtr(3),
			expectedMin:   intPtr(0),
			skip:          true,
			skipReason:    "SKIP: Time before reference defaults to same day not next day - behavior difference",
		},
		{
			name:          "3pm",
			text:          "3pm",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   12,
			expectedHour:  intPtr(15),
			expectedMin:   intPtr(0),
		},
		{
			name:          "at midnight",
			text:          "at midnight",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   13,
			expectedHour:  intPtr(0),
			expectedMin:   intPtr(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDucklingTest(t, tt)
		})
	}
}

// TestDucklingParity_RelativeDurations tests relative duration expressions
// Ported from: Duckling/Time/EN/Corpus.hs
func TestDucklingParity_RelativeDurations(t *testing.T) {
	tests := []ducklingTestCase{
		{
			name:          "in 2 minutes",
			text:          "in 2 minutes",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   12,
			expectedHour:  intPtr(4),
			expectedMin:   intPtr(32),
			skip:          true,
			skipReason:    "SKIP: 'in X minutes' with minute precision not implemented - see issue #112",
		},
		{
			name:          "in a week",
			text:          "in a week",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   19,
			skip:          true,
			skipReason:    "SKIP: 'in a week' (singular with article) not implemented - see issue #112",
		},
		{
			name:          "7 days ago",
			text:          "7 days ago",
			expectedYear:  2013,
			expectedMonth: 2,
			expectedDay:   5,
			skip:          true,
			skipReason:    "SKIP: 'X days ago' returns incorrect date - needs investigation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDucklingTest(t, tt)
		})
	}
}

// TestDucklingParity_Intervals tests date range/interval expressions
// Ported from: Duckling/Time/EN/Corpus.hs
func TestDucklingParity_Intervals(t *testing.T) {
	t.Skip("SKIP: Date range/interval parsing requires range refiners - see issue #112")

	// These tests require parsing date ranges like:
	// - "July 13-15" → start: July 13, end: July 15
	// - "this morning" → start: morning begin, end: morning end
	// - "this week-end" → start: Friday evening, end: Sunday night
	//
	// Kronos would need refiners to handle range expressions properly.
}

// TestDucklingParity_Holidays tests holiday recognition
// Ported from: Duckling/Time/EN/Corpus.hs
func TestDucklingParity_Holidays(t *testing.T) {
	t.Skip("SKIP: Holiday recognition not implemented - would need dedicated holiday parser")

	// Duckling recognizes holidays like:
	// - "thanksgiving" → November 28, 2013
	// - "christmas" → December 25, 2013
	// - "easter" → March 31, 2013
	//
	// This would require a dedicated holiday parser with year-specific calculations.
}

// TestDucklingParity_NegativeExamples tests inputs that should NOT parse
// Ported from: Duckling/Time/EN/Corpus.hs (negative corpus)
func TestDucklingParity_NegativeExamples(t *testing.T) {
	negativeTests := []string{
		"laughing out loud",
		"1 adult",
		"we are separated",
		"three twenty", // Should not parse as time without context
	}

	for _, text := range negativeTests {
		t.Run(text, func(t *testing.T) {
			chrono := createDucklingChrono()
			results := chrono.Parse(text, ducklingRefDate, nil)

			// These should not parse as dates
			assert.Empty(t, results, "Should not parse: %s", text)
		})
	}
}

// TestDucklingParity_GrainPrecision tests grain-based precision
// Duckling tracks the "grain" or precision of parsed dates (year, month, day, hour, etc.)
func TestDucklingParity_GrainPrecision(t *testing.T) {
	t.Skip("SKIP: Grain-based precision tracking would need IsCertain() assertions for each component")

	// Duckling tracks grain/precision levels:
	// - "2015" → year grain only
	// - "March 2015" → month grain
	// - "March 3, 2015" → day grain
	// - "March 3, 2015 at 3pm" → hour grain
	//
	// In Kronos, this maps to IsCertain() for components.
	// A full test suite would verify certainty for each expression.
}

// createDucklingChrono creates a Chrono instance for Duckling parity tests
func createDucklingChrono() *kronos.Chrono {
	config := &kronos.Configuration{
		Parsers: []kronos.Parser{
			parsers.NewENCasualDateParser(),
			parsers.NewENCasualTimeParser(),
			parsers.NewENWeekdayParser(),
			parsers.NewENTimeExpressionParser(false),
			parsers.NewENMonthNameParser(),
			parsers.NewENSlashMonthFormatParser(),
		},
	}
	return kronos.NewChrono(config)
}

// ducklingTestCase is the common test case structure for Duckling parity tests
type ducklingTestCase struct {
	name          string
	text          string
	expectedYear  int
	expectedMonth int
	expectedDay   int
	expectedHour  *int
	expectedMin   *int
	skip          bool
	skipReason    string
}

// runDucklingTest runs a single Duckling parity test case
func runDucklingTest(t *testing.T, tt ducklingTestCase) {
	if tt.skip {
		t.Skip(tt.skipReason)
	}

	chrono := createDucklingChrono()
	results := chrono.Parse(tt.text, ducklingRefDate, nil)

	assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
	if len(results) == 0 {
		return
	}

	result := results[0]
	assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
	assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
	assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")

	if tt.expectedHour != nil {
		assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch")
	}
	if tt.expectedMin != nil {
		assert.Equal(t, *tt.expectedMin, *result.Start().Get(kronos.ComponentMinute), "Minute mismatch")
	}
}
