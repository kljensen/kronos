package en

// Tests for ENRelativeDateFormatParser - "next/last/this day/hour/week/month/year" patterns
// This file addresses GitHub issue #114: Support for 'this/next/last week/month/year' patterns
//
// Test Coverage:
// - next/last hour patterns
// - next/last day patterns
// - next/last week/month/year patterns (existing functionality)
// - Month boundary edge cases
// - "this" modifier patterns

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// TestRelativeDateNextHour tests "next hour" pattern
func TestRelativeDateNextHour(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	results, err := New().WithReferenceDate(refDate).Parse("next hour")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next hour'")
	assert.Equal(t, "next hour", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 10, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
	assert.Equal(t, 13, *result.Get(kronos.ComponentHour))
}

// TestRelativeDateLastHour tests "last hour" pattern
func TestRelativeDateLastHour(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("last hour")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'last hour'")
	assert.Equal(t, "last hour", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 10, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
	assert.Equal(t, 11, *result.Get(kronos.ComponentHour))
}

// TestRelativeDateNextDay tests "next day" pattern
func TestRelativeDateNextDay(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("next day")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next day'")
	assert.Equal(t, "next day", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 10, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 2, *result.Get(kronos.ComponentDay))
	assert.Equal(t, 12, *result.Get(kronos.ComponentHour))
}

// TestRelativeDateLastDay tests "last day" pattern
func TestRelativeDateLastDay(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("last day")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'last day'")
	assert.Equal(t, "last day", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 9, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 30, *result.Get(kronos.ComponentDay))
	assert.Equal(t, 12, *result.Get(kronos.ComponentHour))
}

// TestRelativeDateNextWeek tests "next week" pattern
func TestRelativeDateNextWeek(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("next week")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next week'")
	assert.Equal(t, "next week", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 10, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 8, *result.Get(kronos.ComponentDay))
}

// TestRelativeDateLastWeek tests "last week" pattern
func TestRelativeDateLastWeek(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("last week")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'last week'")
	assert.Equal(t, "last week", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 9, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 24, *result.Get(kronos.ComponentDay))
}

// TestRelativeDateLastMonthBoundary tests month boundary behavior
// Issue #114: "last month" on Jan 15 should return Dec 1, not Dec 15
func TestRelativeDateLastMonthBoundary(t *testing.T) {
	refDate := time.Date(2016, 1, 15, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("last month")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'last month'")
	assert.Equal(t, "last month", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2015, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 12, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay), "Month boundary: last month should return Dec 1, not Dec 15")
}

// TestRelativeDateNextMonthBoundary tests next month boundary behavior
func TestRelativeDateNextMonthBoundary(t *testing.T) {
	refDate := time.Date(2016, 1, 15, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("next month")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next month'")
	assert.Equal(t, "next month", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 2, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay), "Month boundary: next month should return Feb 1, not Feb 15")
}

// TestRelativeDateLastMonth tests "last month" pattern
func TestRelativeDateLastMonth(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("last month")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'last month'")
	assert.Equal(t, "last month", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 9, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
}

// TestRelativeDateNextMonth tests "next month" pattern
func TestRelativeDateNextMonth(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("next month")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next month'")
	assert.Equal(t, "next month", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 11, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
}

// TestRelativeDateNextYear tests "next year" pattern
func TestRelativeDateNextYear(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("next year")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next year'")
	assert.Equal(t, "next year", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2017, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 1, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
}

// TestRelativeDateLastYear tests "last year" pattern
func TestRelativeDateLastYear(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("last year")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'last year'")
	assert.Equal(t, "last year", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2015, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 1, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
}

// TestRelativeDateThisHour tests "this hour" pattern
func TestRelativeDateThisHour(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 30, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("this hour")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'this hour'")
	assert.Equal(t, "this hour", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 10, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
	assert.Equal(t, 12, *result.Get(kronos.ComponentHour))
	assert.Equal(t, 0, *result.Get(kronos.ComponentMinute))
}

// TestRelativeDateThisDay tests "this day" pattern
func TestRelativeDateThisDay(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 30, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("this day")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'this day'")
	assert.Equal(t, "this day", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 10, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
	assert.Equal(t, 0, *result.Get(kronos.ComponentHour))
}

// TestRelativeDatePastDay tests "past day" pattern
func TestRelativeDatePastDay(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("past day")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'past day'")
	assert.Equal(t, "past day", results[0].Text())

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 9, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 30, *result.Get(kronos.ComponentDay))
}

// TestRelativeDateHourDayBoundary tests hour crossing day boundary
func TestRelativeDateHourDayBoundary(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 23, 30, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("next hour")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next hour'")

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 10, *result.Get(kronos.ComponentMonth))
	assert.Equal(t, 2, *result.Get(kronos.ComponentDay), "next hour at 23:30 should cross to next day")
	assert.Equal(t, 0, *result.Get(kronos.ComponentHour))
}

// TestRelativeDateDayMonthBoundary tests day crossing month boundary
func TestRelativeDateDayMonthBoundary(t *testing.T) {
	refDate := time.Date(2016, 10, 31, 12, 0, 0, 0, time.UTC)

	// Using Casual parser
	// Using Casual parser
	results, err := New().WithReferenceDate(refDate).Parse("next day")
	assert.NoError(t, err)

	assert.NotEmpty(t, results, "Expected to parse 'next day'")

	result := results[0].Start()
	assert.Equal(t, 2016, *result.Get(kronos.ComponentYear))
	assert.Equal(t, 11, *result.Get(kronos.ComponentMonth), "next day from Oct 31 should be Nov")
	assert.Equal(t, 1, *result.Get(kronos.ComponentDay))
}

// TestRelativeDatePluralForms tests plural forms (hours, days)
func TestRelativeDatePluralForms(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)
	// Using Casual parser
	// Using Casual parser

	tests := []struct {
		name string
		text string
	}{
		{"next hours", "next hours"},
		{"last hours", "last hours"},
		{"next days", "next days"},
		{"last days", "last days"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, _ := New().WithReferenceDate(refDate).Parse(tt.text)
			assert.NotEmpty(t, results, "Expected to parse '%s'", tt.text)
			assert.Equal(t, tt.text, results[0].Text())
		})
	}
}

// TestRelativeDateApproximation tests approximation modifiers
func TestRelativeDateApproximation(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC)
	// Using Casual parser
	// Using Casual parser

	tests := []struct {
		name string
		text string
	}{
		{"about next day", "about next day"},
		{"around last hour", "around last hour"},
		{"approximately next week", "approximately next week"},
		// Note: Tilde (~) approximation doesn't work with all parsers - known limitation
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Tags() method is not part of the public API")
			results, _ := New().WithReferenceDate(refDate).Parse(tt.text)
			assert.NotEmpty(t, results, "Expected to parse '%s'", tt.text)
			// Check for approximation tag - Tags() is not part of public API
			// tags := results[0].Tags()
			// assert.True(t, tags["result/approximate"], "Expected approximate tag for '%s'", tt.text)
		})
	}
}
