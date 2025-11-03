package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENMonthNameLittleEndianParser_SingleExpression(t *testing.T) {
	parser := NewENMonthNameLittleEndianParser()
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedText  string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "10 August 2012",
			expectedText:  "10 August 2012",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "3rd Feb 82",
			expectedText:  "3rd Feb 82",
			expectedDay:   3,
			expectedMonth: 2,
			expectedYear:  1982,
		},
		{
			text:          "Sun 15Sep",
			expectedText:  "15Sep",
			expectedDay:   15,
			expectedMonth: 9,
			expectedYear:  2013,
		},
		{
			text:          "SUN 15SEP",
			expectedText:  "15SEP",
			expectedDay:   15,
			expectedMonth: 9,
			expectedYear:  2013,
		},
		{
			text:          "The Deadline is 10 August",
			expectedText:  "10 August",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "31st March, 2016",
			expectedText:  "31st March, 2016",
			expectedDay:   31,
			expectedMonth: 3,
			expectedYear:  2016,
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

func TestENMonthNameLittleEndianParser_WithWeekday(t *testing.T) {
	parser := NewENMonthNameLittleEndianParser()
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedText  string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "The Deadline is Tuesday, 10 January",
			expectedText:  "10 January",
			expectedDay:   10,
			expectedMonth: 1,
			expectedYear:  2013,
		},
		{
			text:          "The Deadline is Tue, 10 January",
			expectedText:  "10 January",
			expectedDay:   10,
			expectedMonth: 1,
			expectedYear:  2013,
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

func TestENMonthNameLittleEndianParser_OrdinalNumbers(t *testing.T) {
	parser := NewENMonthNameLittleEndianParser()
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "1st Jan 2020",
			expectedDay:   1,
			expectedMonth: 1,
			expectedYear:  2020,
		},
		{
			text:          "2nd Feb 2020",
			expectedDay:   2,
			expectedMonth: 2,
			expectedYear:  2020,
		},
		{
			text:          "3rd Mar 2020",
			expectedDay:   3,
			expectedMonth: 3,
			expectedYear:  2020,
		},
		{
			text:          "21st April 2020",
			expectedDay:   21,
			expectedMonth: 4,
			expectedYear:  2020,
		},
		{
			text:          "22nd May 2020",
			expectedDay:   22,
			expectedMonth: 5,
			expectedYear:  2020,
		},
		{
			text:          "23rd June 2020",
			expectedDay:   23,
			expectedMonth: 6,
			expectedYear:  2020,
		},
		{
			text:          "31st December 2020",
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

func TestENMonthNameLittleEndianParser_InvalidDates(t *testing.T) {
	parser := NewENMonthNameLittleEndianParser()
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	// Test cases that should NOT parse (day > 31)
	tests := []string{
		"96 Aug", // 96 is too large for a day
		"50 January 2020",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse(text, refDate, nil)
			

			assert.Empty(t, results, "Should not parse: %s", text)
		})
	}
}
