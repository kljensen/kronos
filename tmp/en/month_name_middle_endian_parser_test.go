package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENMonthNameMiddleEndianParser_SingleExpression(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedText  string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "She is getting married soon (July 2017).",
			expectedText:  "July 2017",
			expectedDay:   1,
			expectedMonth: 7,
			expectedYear:  2017,
		},
		{
			text:          "She is leaving in August.",
			expectedText:  "August",
			expectedDay:   1,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "I am arriving sometime in August, 2012, probably.",
			expectedText:  "August, 2012",
			expectedDay:   1,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "August 10, 2012",
			expectedText:  "August 10, 2012",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "Nov 12, 2011",
			expectedText:  "Nov 12, 2011",
			expectedDay:   12,
			expectedMonth: 11,
			expectedYear:  2011,
		},
		{
			text:          "The Deadline is August 10",
			expectedText:  "August 10",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
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

func TestENMonthNameMiddleEndianParser_SpecialYears(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedText  string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "The Deadline is August 10 2555 BE",
			expectedText:  "August 10 2555 BE",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012, // 2555 BE = 2012 AD
		},
		{
			text:          "The Deadline is August 10, 345 BC",
			expectedText:  "August 10, 345 BC",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  -345,
		},
		{
			text:          "The Deadline is August 10, 8 AD",
			expectedText:  "August 10, 8 AD",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
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

func TestENMonthNameMiddleEndianParser_SkipYearLikeDates(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(true)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	// Should skip "January 21" when shouldSkipYearLikeDate is true
	// because 21 could be mistaken for a year
	tests := []string{
		"January 21",
		"February 22",
		"March 23",
		"April 24",
		"May 25",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse(text, refDate, nil)
			

			assert.Empty(t, results, "Should skip year-like date: %s", text)
		})
	}
}

func TestENMonthNameMiddleEndianParser_DontSkipWithYear(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(true)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	// Should NOT skip when year is present, even with shouldSkipYearLikeDate
	tests := []struct {
		text          string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "January 21, 2020",
			expectedDay:   21,
			expectedMonth: 1,
			expectedYear:  2020,
		},
		{
			text:          "February 22 2021",
			expectedDay:   22,
			expectedMonth: 2,
			expectedYear:  2021,
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
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

func TestENMonthNameMiddleEndianParser_OrdinalNumbers(t *testing.T) {
	parser := NewENMonthNameMiddleEndianParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "Jan 1st, 2024",
			expectedDay:   1,
			expectedMonth: 1,
			expectedYear:  2024,
		},
		{
			text:          "February 2nd, 2020",
			expectedDay:   2,
			expectedMonth: 2,
			expectedYear:  2020,
		},
		{
			text:          "March 3rd, 2020",
			expectedDay:   3,
			expectedMonth: 3,
			expectedYear:  2020,
		},
		{
			text:          "April 21st, 2020",
			expectedDay:   21,
			expectedMonth: 4,
			expectedYear:  2020,
		},
		{
			text:          "May 22nd, 2020",
			expectedDay:   22,
			expectedMonth: 5,
			expectedYear:  2020,
		},
		{
			text:          "June 23rd, 2020",
			expectedDay:   23,
			expectedMonth: 6,
			expectedYear:  2020,
		},
		{
			text:          "December 31st, 2020",
			expectedDay:   31,
			expectedMonth: 12,
			expectedYear:  2020,
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
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
