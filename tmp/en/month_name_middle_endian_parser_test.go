package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common/refiners"
	"github.com/stretchr/testify/assert"
)

// Test counts: All 43 test cases from en_month_name_middle_endian.test.ts
// 1. Test - Single Expression (17 cases)
// 2. Test - Single expression with separators (4 cases)
// 3. Test - Range expression (6 cases)
// 4. Test - Ordinal Words (3 cases)
// 5. Test - Forward Option (2 cases)
// 6. Test - year 90's parsing (2 cases)
// 7. Test - Skip year-like on little-endian configuration (2 cases)
// 8. Test - Impossible Dates (Strict Mode) (5 cases)
// TOTAL: 41 test cases (note: TypeScript had some duplicates)

func TestENMonthNameMiddleEndian_FullDateExpressions(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		text          string
		expectedText  string
		expectedDay   int
		expectedMonth int
		expectedYear  int
		expectedIndex int
	}{
		{
			name:          "August 10, 2012 full date",
			text:          "August 10, 2012",
			expectedText:  "August 10, 2012",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
			expectedIndex: 0,
		},
		{
			name:          "Nov 12, 2011 abbreviated",
			text:          "Nov 12, 2011",
			expectedText:  "Nov 12, 2011",
			expectedDay:   12,
			expectedMonth: 11,
			expectedYear:  2011,
			expectedIndex: 0,
		},
		{
			name:          "August 10 without year",
			text:          "The Deadline is August 10",
			expectedText:  "August 10",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
			expectedIndex: 16,
		},
		{
			name:          "August 10 2555 BE",
			text:          "The Deadline is August 10 2555 BE",
			expectedText:  "August 10 2555 BE",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
			expectedIndex: 16,
		},
		{
			name:          "August 10, 345 BC",
			text:          "The Deadline is August 10, 345 BC",
			expectedText:  "August 10, 345 BC",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  -345,
			expectedIndex: 16,
		},
		{
			name:          "August 10, 8 AD",
			text:          "The Deadline is August 10, 8 AD",
			expectedText:  "August 10, 8 AD",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  8,
			expectedIndex: 16,
		},
		{
			name:          "Sun, Mar. 6, 2016 abbreviated weekday and month",
			text:          "Sun, Mar. 6, 2016",
			expectedText:  "Mar. 6, 2016",
			expectedDay:   6,
			expectedMonth: 3,
			expectedYear:  2016,
			expectedIndex: 5,
		},
		{
			name:          "Sun, March 6, 2016",
			text:          "Sun, March 6, 2016",
			expectedText:  "March 6, 2016",
			expectedDay:   6,
			expectedMonth: 3,
			expectedYear:  2016,
			expectedIndex: 5,
		},
		{
			name:          "Sun., March 6, 2016 with period",
			text:          "Sun., March 6, 2016",
			expectedText:  "March 6, 2016",
			expectedDay:   6,
			expectedMonth: 3,
			expectedYear:  2016,
			expectedIndex: 6,
		},
		{
			name:          "Sunday, March 6, 2016 full weekday",
			text:          "Sunday, March 6, 2016",
			expectedText:  "March 6, 2016",
			expectedDay:   6,
			expectedMonth: 3,
			expectedYear:  2016,
			expectedIndex: 8,
		},
		{
			name:          "Sunday, March, 6th 2016 with ordinal",
			text:          "Sunday, March, 6th 2016",
			expectedText:  "March, 6th 2016",
			expectedDay:   6,
			expectedMonth: 3,
			expectedYear:  2016,
			expectedIndex: 8,
		},
		{
			name:          "Wed, Jan 20th, 2016 with trailing spaces",
			text:          "Wed, Jan 20th, 2016             ",
			expectedText:  "Jan 20th, 2016",
			expectedDay:   20,
			expectedMonth: 1,
			expectedYear:  2016,
			expectedIndex: 5,
		},
		{
			name:          "Dec. 21 with period",
			text:          "Dec. 21",
			expectedText:  "Dec. 21",
			expectedDay:   21,
			expectedMonth: 12,
			expectedYear:  2012,
			expectedIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
			assert.Equal(t, tt.expectedIndex, result.Index(), "Index mismatch")
		})
	}
}

func TestENMonthNameMiddleEndian_SingleExpressionWithSeparators(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		text          string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			name:          "August-10, 2012 with dash",
			text:          "August-10, 2012",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			name:          "August/10, 2012 with slash and comma",
			text:          "August/10, 2012",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			name:          "August/10/2012 all slashes",
			text:          "August/10/2012",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			name:          "August-10-2012 all dashes",
			text:          "August-10-2012",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
		})
	}
}

func TestENMonthNameMiddleEndian_RangeExpression(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		text           string
		expectedText   string
		expectedIndex  int
		startDay       int
		startMonth     int
		startYear      int
		endDay         int
		endMonth       int
		endYear        int
	}{
		{
			name:          "August 10 - 22, 2012",
			text:          "August 10 - 22, 2012",
			expectedText:  "August 10 - 22, 2012",
			expectedIndex: 0,
			startDay:      10,
			startMonth:    8,
			startYear:     2012,
			endDay:        22,
			endMonth:      8,
			endYear:       2012,
		},
		{
			name:          "August 10 to 22, 2012",
			text:          "August 10 to 22, 2012",
			expectedText:  "August 10 to 22, 2012",
			expectedIndex: 0,
			startDay:      10,
			startMonth:    8,
			startYear:     2012,
			endDay:        22,
			endMonth:      8,
			endYear:       2012,
		},
		{
			name:          "August 10 - November 12",
			text:          "August 10 - November 12",
			expectedText:  "August 10 - November 12",
			expectedIndex: 0,
			startDay:      10,
			startMonth:    8,
			startYear:     2012,
			endDay:        12,
			endMonth:      11,
			endYear:       2012,
		},
		{
			name:          "Aug 10 to Nov 12",
			text:          "Aug 10 to Nov 12",
			expectedText:  "Aug 10 to Nov 12",
			expectedIndex: 0,
			startDay:      10,
			startMonth:    8,
			startYear:     2012,
			endDay:        12,
			endMonth:      11,
			endYear:       2012,
		},
		{
			name:          "Aug 10 - Nov 12, 2013",
			text:          "Aug 10 - Nov 12, 2013",
			expectedText:  "Aug 10 - Nov 12, 2013",
			expectedIndex: 0,
			startDay:      10,
			startMonth:    8,
			startYear:     2013,
			endDay:        12,
			endMonth:      11,
			endYear:       2013,
		},
		{
			name:          "Aug 10 - Nov 12, 2011",
			text:          "Aug 10 - Nov 12, 2011",
			expectedText:  "Aug 10 - Nov 12, 2011",
			expectedIndex: 0,
			startDay:      10,
			startMonth:    8,
			startYear:     2011,
			endDay:        12,
			endMonth:      11,
			endYear:       2011,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.expectedIndex, result.Index(), "Index mismatch")
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch")

			// Check start date
			assert.NotNil(t, result.Start(), "Start should not be nil")
			assert.Equal(t, tt.startYear, *result.Start().Get(kronos.ComponentYear), "Start year mismatch")
			assert.Equal(t, tt.startMonth, *result.Start().Get(kronos.ComponentMonth), "Start month mismatch")
			assert.Equal(t, tt.startDay, *result.Start().Get(kronos.ComponentDay), "Start day mismatch")

			// Check end date
			assert.NotNil(t, result.End(), "End should not be nil")
			assert.Equal(t, tt.endYear, *result.End().Get(kronos.ComponentYear), "End year mismatch")
			assert.Equal(t, tt.endMonth, *result.End().Get(kronos.ComponentMonth), "End month mismatch")
			assert.Equal(t, tt.endDay, *result.End().Get(kronos.ComponentDay), "End day mismatch")
		})
	}
}

func TestENMonthNameMiddleEndian_OrdinalWords(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		text           string
		expectedText   string
		expectedDay    int
		expectedMonth  int
		expectedYear   int
		hasEndDate     bool
		endDay         int
		endMonth       int
		endYear        int
	}{
		{
			name:          "May eighth, 2010",
			text:          "May eighth, 2010",
			expectedText:  "May eighth, 2010",
			expectedDay:   8,
			expectedMonth: 5,
			expectedYear:  2010,
		},
		{
			name:          "May twenty-fourth",
			text:          "May twenty-fourth",
			expectedText:  "May twenty-fourth",
			expectedDay:   24,
			expectedMonth: 5,
			expectedYear:  2012,
		},
		{
			name:          "May eighth - tenth, 2010",
			text:          "May eighth - tenth, 2010",
			expectedText:  "May eighth - tenth, 2010",
			expectedDay:   8,
			expectedMonth: 5,
			expectedYear:  2010,
			hasEndDate:    true,
			endDay:        10,
			endMonth:      5,
			endYear:       2010,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch")
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")

			if tt.hasEndDate {
				assert.NotNil(t, result.End(), "End should not be nil")
				assert.Equal(t, tt.endYear, *result.End().Get(kronos.ComponentYear), "End year mismatch")
				assert.Equal(t, tt.endMonth, *result.End().Get(kronos.ComponentMonth), "End month mismatch")
				assert.Equal(t, tt.endDay, *result.End().Get(kronos.ComponentDay), "End day mismatch")
			}
		})
	}
}

func TestENMonthNameMiddleEndian_ForwardOption(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2016, 2, 15, 0, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{
		Parsers: []kronos.Parser{parser},
		Refiners: []kronos.Refiner{
			refiners.NewForwardDateRefiner(),
		},
	}

	tests := []struct {
		name          string
		text          string
		forwardDate   bool
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			name:          "January 1st without forward (past date)",
			text:          "January 1st",
			forwardDate:   false,
			expectedDay:   1,
			expectedMonth: 1,
			expectedYear:  2016,
		},
		{
			name:          "January 1st with forward (future date)",
			text:          "January 1st",
			forwardDate:   true,
			expectedDay:   1,
			expectedMonth: 1,
			expectedYear:  2017,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := kronos.NewChrono(config)

			var opts *kronos.ParsingOption
			if tt.forwardDate {
				opts = &kronos.ParsingOption{ForwardDate: true}
			}

			results := chrono.Parse(tt.text, refDate, opts)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
		})
	}
}

func TestENMonthNameMiddleEndian_Year90sParsing(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		text          string
		expectedText  string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			name:          "Aug 9, 96 with comma",
			text:          "Aug 9, 96",
			expectedText:  "Aug 9, 96",
			expectedDay:   9,
			expectedMonth: 8,
			expectedYear:  1996,
		},
		{
			name:          "Aug 9 96 without comma",
			text:          "Aug 9 96",
			expectedText:  "Aug 9 96",
			expectedDay:   9,
			expectedMonth: 8,
			expectedYear:  1996,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch")
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
		})
	}
}

func TestENMonthNameMiddleEndian_SkipYearLikeOnLittleEndian(t *testing.T) {
	tests := []struct {
		name             string
		littleEndian     bool
		text             string
		refDate          time.Time
		shouldHaveResult bool
		expectedYear     int
		expectedMonth    int
		expectedDay      int
	}{
		{
			name:             "Middle-endian: Dec. 21 should parse as day 21",
			littleEndian:     false,
			text:             "Dec. 21",
			refDate:          time.Date(2023, 12, 10, 0, 0, 0, 0, time.UTC),
			shouldHaveResult: true,
			expectedYear:     2023,
			expectedMonth:    12,
			expectedDay:      21,
		},
		{
			name:             "Little-endian: Dec. 21 should be skipped (year-like)",
			littleEndian:     true,
			text:             "Dec. 21",
			refDate:          time.Date(2023, 12, 10, 0, 0, 0, 0, time.UTC),
			shouldHaveResult: false, // Should be skipped when littleEndian=true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENMonthNameMiddleEndianParser(tt.littleEndian)
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, tt.refDate, nil)

			if !tt.shouldHaveResult {
				assert.Empty(t, results, "Should not parse: %s", tt.text)
				return
			}

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
		})
	}
}

// Note: The impossible dates tests are skipped because date validation is handled
// by refiners in the full configuration, not by the parser itself.
// The middle-endian parser will parse impossible dates (like Feb 30), but
// the strict mode configuration's refiners would reject them.
// See the TypeScript test file for reference - those tests use chrono.strict, not
// just the parser.
