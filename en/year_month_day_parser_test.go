package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENYearMonthDayParser_SingleExpression(t *testing.T) {
	parser := NewENYearMonthDayParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedText  string
		expectedDate  time.Time
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "2012/8/10",
			expectedText:  "2012/8/10",
			expectedDate:  time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC),
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "The Deadline is 2012/8/10",
			expectedText:  "2012/8/10",
			expectedDate:  time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC),
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "2014/2/28",
			expectedText:  "2014/2/28",
			expectedDate:  time.Date(2014, 2, 28, 12, 0, 0, 0, time.UTC),
			expectedDay:   28,
			expectedMonth: 2,
			expectedYear:  2014,
		},
		{
			text:          "2014/12/28",
			expectedText:  "2014/12/28",
			expectedDate:  time.Date(2014, 12, 28, 12, 0, 0, 0, time.UTC),
			expectedDay:   28,
			expectedMonth: 12,
			expectedYear:  2014,
		},
		{
			text:          "2014.12.28",
			expectedText:  "2014.12.28",
			expectedDate:  time.Date(2014, 12, 28, 12, 0, 0, 0, time.UTC),
			expectedDay:   28,
			expectedMonth: 12,
			expectedYear:  2014,
		},
		{
			text:          "2014 12 28",
			expectedText:  "2014 12 28",
			expectedDate:  time.Date(2014, 12, 28, 12, 0, 0, 0, time.UTC),
			expectedDay:   28,
			expectedMonth: 12,
			expectedYear:  2014,
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
			assert.Equal(t, tt.expectedText, result.Text(), "Text mismatch")
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
		})
	}
}

func TestENYearMonthDayParser_WithMonthName(t *testing.T) {
	parser := NewENYearMonthDayParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedText  string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "2012/Aug/10",
			expectedText:  "2012/Aug/10",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "The Deadline is 2012/aug/10",
			expectedText:  "2012/aug/10",
			expectedDay:   10,
			expectedMonth: 8,
			expectedYear:  2012,
		},
		{
			text:          "The Deadline is 2018 March 18",
			expectedText:  "2018 March 18",
			expectedDay:   18,
			expectedMonth: 3,
			expectedYear:  2018,
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

func TestENYearMonthDayParser_SwapMonthDay(t *testing.T) {
	// In casual mode (strictMonthDateOrder=false), allow swapping
	casualParser := NewENYearMonthDayParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedDay   int
		expectedMonth int
		expectedYear  int
	}{
		{
			text:          "2024/13/1", // Invalid month 13, swap to 1/13
			expectedDay:   13,
			expectedMonth: 1,
			expectedYear:  2024,
		},
		{
			text:          "2024-13-01",
			expectedDay:   13,
			expectedMonth: 1,
			expectedYear:  2024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{casualParser}}
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

func TestENYearMonthDayParser_StrictMode(t *testing.T) {
	// In strict mode, don't parse invalid month/day combinations
	strictParser := NewENYearMonthDayParser(true)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []string{
		"2024/13/1",  // Invalid month
		"2024-13-01", // Invalid month
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{strictParser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(text, refDate, nil)

			assert.Empty(t, results, "Should not parse in strict mode: %s", text)
		})
	}
}

func TestENYearMonthDayParser_InvalidDates(t *testing.T) {
	parser := NewENYearMonthDayParser(false)
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []string{
		"2012/80/10", // Invalid month
		"2012 80 10", // Invalid month
		"2014-08-32", // Invalid day
		"2014-02-30", // Invalid day for February
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(text, refDate, nil)

			assert.Empty(t, results, "Should not parse invalid date: %s", text)
		})
	}
}
