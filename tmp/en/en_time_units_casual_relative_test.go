package en

// Tests ported from chrono's en_time_units_casual_relative.test.ts
//
// Summary of 26 original test cases:
// - 20 test cases PASSING
// - 6 test cases FAILING (parser limitations)
// - 0 test cases SKIPPED
//
// Test execution summary:
// - Positive time units: 5/6 passing (fails: "next two years" - word numbers not supported)
// - Negative time units: 2/4 passing (fails: "last two weeks" - word numbers, "+2 months, 5 days" - partial parse)
// - Plus sign: 4/4 passing
// - Minus sign: 1/2 passing (fails: "-2hr5min" - compound abbreviation)
// - Without abbreviations: 3/3 passing
// - Negative cases: 4/6 passing (fails: "+am", "+them" - parser incorrectly matches these)
//
// Known parser limitations:
// - Written number words (e.g., "two", "three") not fully supported in all contexts
// - Comma-separated multi-unit expressions partially parse first unit only
// - Compound abbreviations without spaces (e.g., "2hr5min") not supported
// - Some negative cases incorrectly match (e.g., "+am", "+them")

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// TestCasualRelativePositiveTimeUnits tests positive time unit expressions
func TestCasualRelativePositiveTimeUnits(t *testing.T) {
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
	}{
		{
			name:          "next 2 weeks",
			text:          "next 2 weeks",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "next 2 weeks",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   15, // Oct 1 + 14 days
			expectedHour:  ptrInt(12),
		},
		{
			name:          "next 2 days",
			text:          "next 2 days",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "next 2 days",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   3, // Oct 3
			expectedHour:  ptrInt(12),
		},
		{
			name:          "next two years",
			text:          "next two years",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "next two years",
			expectedYear:  2018,
			expectedMonth: 10,
			expectedDay:   1,
			expectedHour:  ptrInt(12),
		},
		{
			name:          "next 2 weeks 3 days",
			text:          "next 2 weeks 3 days",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "next 2 weeks 3 days",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   18, // Oct 1 + 14 + 3 days
			expectedHour:  ptrInt(12),
		},
		{
			name:          "after a year",
			text:          "after a year",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "after a year",
			expectedYear:  2017,
			expectedMonth: 10,
			expectedDay:   1,
			expectedHour:  ptrInt(12),
		},
		{
			name:           "after an hour",
			text:           "after an hour",
			refDate:        time.Date(2016, 10, 1, 15, 0, 0, 0, time.UTC),
			expectedText:   "after an hour",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(16),
			expectedMinute: ptrInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENTimeUnitCasualRelativeFormatParser(true)
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)

			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch for: %s", tt.text)

			if tt.expectedHour != nil {
				assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch for: %s", tt.text)
			}
			if tt.expectedMinute != nil {
				assert.Equal(t, *tt.expectedMinute, *result.Start().Get(kronos.ComponentMinute), "Minute mismatch for: %s", tt.text)
			}
		})
	}
}

// TestCasualRelativeNegativeTimeUnits tests negative time unit expressions (past, last)
func TestCasualRelativeNegativeTimeUnits(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedText  string
		expectedYear  int
		expectedMonth int
		expectedDay   int
		expectedHour  *int
	}{
		{
			name:          "last 2 weeks",
			text:          "last 2 weeks",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "last 2 weeks",
			expectedYear:  2016,
			expectedMonth: 9,
			expectedDay:   17, // Oct 1 - 14 days = Sep 17
			expectedHour:  ptrInt(12),
		},
		{
			name:          "last two weeks",
			text:          "last two weeks",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "last two weeks",
			expectedYear:  2016,
			expectedMonth: 9,
			expectedDay:   17, // Oct 1 - 14 days = Sep 17
			expectedHour:  ptrInt(12),
		},
		{
			name:          "past 2 days",
			text:          "past 2 days",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "past 2 days",
			expectedYear:  2016,
			expectedMonth: 9,
			expectedDay:   29, // Oct 1 - 2 days = Sep 29
			expectedHour:  ptrInt(12),
		},
		{
			name:          "+2 months, 5 days",
			text:          "+2 months, 5 days",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "+2 months, 5 days",
			expectedYear:  2016,
			expectedMonth: 12,
			expectedDay:   6, // Oct 1 + 2 months + 5 days = Dec 6
			expectedHour:  ptrInt(12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENTimeUnitCasualRelativeFormatParser(true)
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)

			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch for: %s", tt.text)

			if tt.expectedHour != nil {
				assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch for: %s", tt.text)
			}
		})
	}
}

// TestCasualRelativePlusSign tests expressions with plus sign
func TestCasualRelativePlusSign(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		refDate        time.Time
		expectedText   string
		expectedDay    int
		expectedHour   *int
		expectedMinute *int
	}{
		{
			name:           "+15 minutes",
			text:           "+15 minutes",
			refDate:        time.Date(2012, 7, 10, 12, 14, 0, 0, time.UTC),
			expectedText:   "+15 minutes",
			expectedDay:    10,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(29),
		},
		{
			name:           "+15min",
			text:           "+15min",
			refDate:        time.Date(2012, 7, 10, 12, 14, 0, 0, time.UTC),
			expectedText:   "+15min",
			expectedDay:    10,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(29),
		},
		{
			name:           "+1 day 2 hour",
			text:           "+1 day 2 hour",
			refDate:        time.Date(2012, 7, 10, 12, 14, 0, 0, time.UTC),
			expectedText:   "+1 day 2 hour",
			expectedDay:    11,
			expectedHour:   ptrInt(14),
			expectedMinute: ptrInt(14),
		},
		{
			name:           "+1m",
			text:           "+1m",
			refDate:        time.Date(2012, 7, 10, 12, 14, 0, 0, time.UTC),
			expectedText:   "+1m",
			expectedDay:    10,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(15),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENTimeUnitCasualRelativeFormatParser(true)
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)

			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch for: %s", tt.text)

			if tt.expectedHour != nil {
				assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch for: %s", tt.text)
			}
			if tt.expectedMinute != nil {
				assert.Equal(t, *tt.expectedMinute, *result.Start().Get(kronos.ComponentMinute), "Minute mismatch for: %s", tt.text)
			}
		})
	}
}

// TestCasualRelativeMinusSign tests expressions with minus sign
func TestCasualRelativeMinusSign(t *testing.T) {
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
	}{
		{
			name:           "-3y",
			text:           "-3y",
			refDate:        time.Date(2015, 7, 10, 12, 14, 0, 0, time.UTC),
			expectedText:   "-3y",
			expectedYear:   2012,
			expectedMonth:  7,
			expectedDay:    10,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(14),
		},
		{
			name:           "-2hr5min",
			text:           "-2hr5min",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "-2hr5min",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(9),
			expectedMinute: ptrInt(55),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENTimeUnitCasualRelativeFormatParser(true)
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)

			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch for: %s", tt.text)

			if tt.expectedHour != nil {
				assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch for: %s", tt.text)
			}
			if tt.expectedMinute != nil {
				assert.Equal(t, *tt.expectedMinute, *result.Start().Get(kronos.ComponentMinute), "Minute mismatch for: %s", tt.text)
			}
		})
	}
}

// TestCasualRelativeWithoutAbbreviations tests parser without abbreviations
func TestCasualRelativeWithoutAbbreviations(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		refDate        time.Time
		shouldParse    bool
		expectedText   string
		expectedYear   int
		expectedMonth  int
		expectedDay    int
		expectedHour   *int
		expectedMinute *int
	}{
		{
			name:        "-3y without abbreviations",
			text:        "-3y",
			refDate:     time.Date(2015, 7, 10, 12, 14, 0, 0, time.UTC),
			shouldParse: false,
		},
		{
			name:        "last 2m without abbreviations",
			text:        "last 2m",
			refDate:     time.Date(2015, 7, 10, 12, 14, 0, 0, time.UTC),
			shouldParse: false,
		},
		{
			name:           "-2 hours 5 minutes without abbreviations",
			text:           "-2 hours 5 minutes",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			shouldParse:    true,
			expectedText:   "-2 hours 5 minutes",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(9),
			expectedMinute: ptrInt(55),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENTimeUnitCasualRelativeFormatParser(false) // No abbreviations
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)

			results := chrono.Parse(tt.text, tt.refDate, nil)

			if !tt.shouldParse {
				assert.Empty(t, results, "Expected NO results for: %s", tt.text)
				return
			}

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.expectedText, "Text mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch for: %s", tt.text)
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch for: %s", tt.text)

			if tt.expectedHour != nil {
				assert.Equal(t, *tt.expectedHour, *result.Start().Get(kronos.ComponentHour), "Hour mismatch for: %s", tt.text)
			}
			if tt.expectedMinute != nil {
				assert.Equal(t, *tt.expectedMinute, *result.Start().Get(kronos.ComponentMinute), "Minute mismatch for: %s", tt.text)
			}
		})
	}
}

// TestCasualRelativeNegativeCases tests that certain patterns should NOT be parsed
func TestCasualRelativeNegativeCases(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		refDate time.Time
	}{
		{
			name:    "3y without prefix",
			text:    "3y",
			refDate: time.Date(2015, 7, 10, 12, 14, 0, 0, time.UTC),
		},
		{
			name:    "1 m with space",
			text:    "1 m",
			refDate: time.Date(2015, 7, 10, 12, 14, 0, 0, time.UTC),
		},
		{
			name:    "the day",
			text:    "the day",
			refDate: time.Date(2015, 7, 10, 12, 14, 0, 0, time.UTC),
		},
		{
			name:    "a day",
			text:    "a day",
			refDate: time.Date(2015, 7, 10, 12, 14, 0, 0, time.UTC),
		},
		{
			name:    "+am",
			text:    "+am",
			refDate: time.Now(),
		},
		{
			name:    "+them",
			text:    "+them",
			refDate: time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENTimeUnitCasualRelativeFormatParser(true)
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)

			results := chrono.Parse(tt.text, tt.refDate, nil)
			assert.Empty(t, results, "Expected NO results for '%s', but got %d", tt.text, len(results))
		})
	}
}
