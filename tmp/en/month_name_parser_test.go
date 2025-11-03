package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENMonthNameParser_MonthYear(t *testing.T) {
	parser := NewENMonthNameParser()
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedMonth int
		expectedYear  int
		expectedDay   int
		yearCertain   bool
	}{
		{
			text:          "September 2012",
			expectedMonth: 9,
			expectedYear:  2012,
			expectedDay:   1,
			yearCertain:   true,
		},
		{
			text:          "Sept 2012",
			expectedMonth: 9,
			expectedYear:  2012,
			expectedDay:   1,
			yearCertain:   true,
		},
		{
			text:          "Sep 2012",
			expectedMonth: 9,
			expectedYear:  2012,
			expectedDay:   1,
			yearCertain:   true,
		},
		{
			text:          "Sep. 2012",
			expectedMonth: 9,
			expectedYear:  2012,
			expectedDay:   1,
			yearCertain:   true,
		},
		{
			text:          "Sep-2012",
			expectedMonth: 9,
			expectedYear:  2012,
			expectedDay:   1,
			yearCertain:   true,
		},
		{
			text:          "in June of 2022",
			expectedMonth: 6,
			expectedYear:  2022,
			expectedDay:   1,
			yearCertain:   true,
		},
		{
			text:          "Statement of comprehensive income for the year ended Dec. 2021",
			expectedMonth: 12,
			expectedYear:  2021,
			expectedDay:   1,
			yearCertain:   true,
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
			assert.Equal(t, tt.yearCertain, result.Start().IsCertain(kronos.ComponentYear), "Year certainty mismatch")
			assert.True(t, result.Start().IsCertain(kronos.ComponentMonth), "Month should be certain")
			assert.False(t, result.Start().IsCertain(kronos.ComponentDay), "Day should be implied")
		})
	}
}

func TestENMonthNameParser_MonthOnly(t *testing.T) {
	parser := NewENMonthNameParser()

	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedMonth int
		expectedYear  int
		expectedDay   int
	}{
		{
			name:          "In January (from November)",
			text:          "In January",
			refDate:       time.Date(2020, 11, 22, 0, 0, 0, 0, time.UTC),
			expectedMonth: 1,
			expectedYear:  2021, // Next year since we're in November
			expectedDay:   1,
		},
		{
			name:          "in Jan (from November)",
			text:          "in Jan",
			refDate:       time.Date(2020, 11, 22, 0, 0, 0, 0, time.UTC),
			expectedMonth: 1,
			expectedYear:  2021,
			expectedDay:   1,
		},
		{
			name:          "May (from November)",
			text:          "May",
			refDate:       time.Date(2020, 11, 22, 0, 0, 0, 0, time.UTC),
			expectedMonth: 5,
			expectedYear:  2021,
			expectedDay:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse(tt.text, tt.refDate, nil)
			

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Contains(t, result.Text(), tt.text, "Text mismatch")
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
			assert.False(t, result.Start().IsCertain(kronos.ComponentYear), "Year should be implied")
			assert.True(t, result.Start().IsCertain(kronos.ComponentMonth), "Month should be certain")
			assert.False(t, result.Start().IsCertain(kronos.ComponentDay), "Day should be implied")
		})
	}
}
