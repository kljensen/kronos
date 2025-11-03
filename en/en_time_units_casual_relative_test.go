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

// TestExtendedNumberWords tests extended number words like "a", "an", "few", "several", "couple", "half"
func TestExtendedNumberWords(t *testing.T) {
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
		expectedSecond *int
	}{
		{
			name:           "a minute ago",
			text:           "last a minute",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last a minute",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(11),
			expectedMinute: ptrInt(59),
		},
		{
			name:           "an hour ago",
			text:           "last an hour",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last an hour",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(11),
			expectedMinute: ptrInt(0),
		},
		{
			name:          "a day ago",
			text:          "last a day",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "last a day",
			expectedYear:  2016,
			expectedMonth: 9,
			expectedDay:   30,
			expectedHour:  ptrInt(12),
		},
		{
			name:          "a week ago",
			text:          "last a week",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "last a week",
			expectedYear:  2016,
			expectedMonth: 9,
			expectedDay:   24, // Sept 24
			expectedHour:  ptrInt(12),
		},
		{
			name:          "a couple of days ago",
			text:          "last a couple of days",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "last a couple of days",
			expectedYear:  2016,
			expectedMonth: 9,
			expectedDay:   29, // Sept 29 (Oct 1 - 2 days)
			expectedHour:  ptrInt(12),
		},
		{
			name:          "couple days ago (without of)",
			text:          "last couple days",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "last couple days",
			expectedYear:  2016,
			expectedMonth: 9,
			expectedDay:   29, // Sept 29
			expectedHour:  ptrInt(12),
		},
		{
			name:           "a few hours ago",
			text:           "last a few hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last a few hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(9),
			expectedMinute: ptrInt(0), // 12 - 3 = 9
		},
		{
			name:           "few hours ago (without a)",
			text:           "last few hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last few hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(9),
			expectedMinute: ptrInt(0),
		},
		{
			name:          "several weeks ago",
			text:          "last several weeks",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "last several weeks",
			expectedYear:  2016,
			expectedMonth: 8,
			expectedDay:   13, // Oct 1 - 49 days (7 weeks) = Aug 13
			expectedHour:  ptrInt(12),
		},
		{
			name:           "half an hour ago",
			text:           "last half an hour",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last half an hour",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(11),
			expectedMinute: ptrInt(30), // 12:00 - 0.5h = 11:30
		},
		{
			name:           "half a day ago",
			text:           "last half a day",
			refDate:        time.Date(2016, 10, 2, 12, 0, 0, 0, time.UTC),
			expectedText:   "last half a day",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    2,
			expectedHour:   ptrInt(0), // Oct 2 12:00 - 12h = Oct 2 00:00
			expectedMinute: ptrInt(0),
		},
		{
			name:          "a dozen hours ago",
			text:          "last a dozen hours",
			refDate:       time.Date(2016, 10, 2, 12, 0, 0, 0, time.UTC),
			expectedText:  "last a dozen hours",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   2,
			expectedHour:  ptrInt(0), // Oct 2 12:00 - 12h = Oct 2 00:00
		},
		// Test with "next" prefix
		{
			name:           "next an hour",
			text:           "next an hour",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "next an hour",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(13),
			expectedMinute: ptrInt(0),
		},
		{
			name:          "next a couple of days",
			text:          "next a couple of days",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:  "next a couple of days",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   3,
			expectedHour:  ptrInt(12),
		},
		{
			name:           "next half an hour",
			text:           "next half an hour",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "next half an hour",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(30),
		},
		// Test with + prefix
		{
			name:           "+a minute",
			text:           "+a minute",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "+a minute",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(1),
		},
		{
			name:           "+half an hour",
			text:           "+half an hour",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "+half an hour",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(30),
		},
		// Test with - prefix
		{
			name:           "-an hour",
			text:           "-an hour",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "-an hour",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(11),
			expectedMinute: ptrInt(0),
		},
		{
			name:          "-a couple days",
			text:          "-a couple days",
			refDate:       time.Date(2016, 10, 3, 12, 0, 0, 0, time.UTC),
			expectedText:  "-a couple days",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   1,
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
			if tt.expectedMinute != nil {
				assert.Equal(t, *tt.expectedMinute, *result.Start().Get(kronos.ComponentMinute), "Minute mismatch for: %s", tt.text)
			}
			if tt.expectedSecond != nil {
				assert.Equal(t, *tt.expectedSecond, *result.Start().Get(kronos.ComponentSecond), "Second mismatch for: %s", tt.text)
			}
		})
	}
}

// TestFractionalTimeUnits tests fractional time unit expressions
func TestFractionalTimeUnits(t *testing.T) {
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
		expectedSecond *int
	}{
		{
			name:           "2.5 hours with next",
			text:           "next 2.5 hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "next 2.5 hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(14),
			expectedMinute: ptrInt(30),
		},
		{
			name:           "1.5 days with last",
			text:           "last 1.5 days",
			refDate:        time.Date(2016, 10, 2, 12, 0, 0, 0, time.UTC),
			expectedText:   "last 1.5 days",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(0), // -1.5 days = -1 day -12 hours = Oct 1 00:00
			expectedMinute: ptrInt(0),
		},
		{
			name:           "0.5 weeks from now",
			text:           "next 0.5 weeks",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "next 0.5 weeks",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    5, // 0.5 weeks = 3.5 days, rounded to 4 days
			expectedHour:   ptrInt(12),
		},
		{
			name:           "3.25 minutes ago",
			text:           "last 3.25 minutes",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last 3.25 minutes",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(11),
			expectedMinute: ptrInt(56),
			expectedSecond: ptrInt(45),
		},
		{
			name:           "in 2.5 hours",
			text:           "next 2.5 hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "next 2.5 hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(14),
			expectedMinute: ptrInt(30),
		},
		{
			name:           "10.75 minutes ago",
			text:           "last 10.75 minutes",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last 10.75 minutes",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(11),
			expectedMinute: ptrInt(49),
			expectedSecond: ptrInt(15),
		},
		{
			name:           "2,5 hours with last (comma separator)",
			text:           "last 2,5 hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last 2,5 hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(9),
			expectedMinute: ptrInt(30),
		},
		{
			name:           "1,5 days with past (comma separator)",
			text:           "past 1,5 days",
			refDate:        time.Date(2016, 10, 2, 12, 0, 0, 0, time.UTC),
			expectedText:   "past 1,5 days",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(0), // -1.5 days = -1 day -12 hours = Oct 1 00:00
			expectedMinute: ptrInt(0),
		},
		{
			name:           "+0.5 hours",
			text:           "+0.5 hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "+0.5 hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(12),
			expectedMinute: ptrInt(30),
		},
		{
			name:           "+1.5 days 2.5 hours",
			text:           "+1.5 days 2.5 hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "+1.5 days 2.5 hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    3,
			expectedHour:   ptrInt(2),
			expectedMinute: ptrInt(30),
		},
		{
			name:           "0.001 seconds with last",
			text:           "last 0.001 seconds",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "last 0.001 seconds",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(11),
			expectedMinute: ptrInt(59),
			expectedSecond: ptrInt(59), // -0.001 seconds = -1 millisecond
		},
		{
			name:           "-2.5 hours",
			text:           "-2.5 hours",
			refDate:        time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedText:   "-2.5 hours",
			expectedYear:   2016,
			expectedMonth:  10,
			expectedDay:    1,
			expectedHour:   ptrInt(9),
			expectedMinute: ptrInt(30),
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
			if tt.expectedSecond != nil {
				assert.Equal(t, *tt.expectedSecond, *result.Start().Get(kronos.ComponentSecond), "Second mismatch for: %s", tt.text)
			}
		})
	}
}
