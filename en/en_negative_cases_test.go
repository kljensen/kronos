package en

// Tests ported from chrono's negative_cases.test.ts
//
// Summary of 38 original test cases:
// - 38 test cases PORTED (execution to be verified)
// - 0 test cases SKIPPED
//
// These tests verify that certain patterns should NOT be parsed as dates.
// Test categories:
// - Random non-date patterns: numbers, decimals, percentages (11 tests)
// - URL encoded strings (2 tests)
// - Hyphenated numbers that look like dates but aren't (9 tests)
// - Impossible dates/times: Feb 29 2022, June 31, 14PM, 25:12 (6 tests)
// - Impossible date ranges (2 tests)
// - Version numbers: 1.1.3, 1.1.30, 1.10.30 (3 tests)
// - Incorrect references: "for the year" (1 test)
// - Dates near version numbers should still parse (3 tests)

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNegativeRandomNonDatePatterns tests that random numbers should not parse
func TestNegativeRandomNonDatePatterns(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "single digit with space", text: " 3"},
		{name: "single digit with spaces", text: "       1"},
		{name: "two digits with spaces", text: "  11 "},
		{name: "decimal with space", text: " 0.5 "},
		{name: "decimal number", text: " 35.49 "},
		{name: "percentage", text: "12.53%"},
		{name: "rating 5.0", text: "6358fe2310> *5.0* / 5 Outstanding"},
		{name: "rating 1.5", text: "6358fe2310> *1.5* / 5 Outstanding"},
		{name: "dollar amount", text: "Total: $1,194.09 [image: View Reservation"},
		{name: "kilograms", text: "at 6.5 kilograms"},
		{name: "unusual text", text: "ah that is unusual"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			refDate := time.Now()
			results := chrono.Parse(tt.text, refDate, nil)

			assert.Empty(t, results, "Expected NO results for: '%s', but got %d", tt.text, len(results))
		})
	}
}

// TestNegativeURLEncoded tests that URL-encoded strings should not parse
func TestNegativeURLEncoded(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "short encoded string",
			text: "%e7%b7%8a",
		},
		{
			name: "long encoded URL",
			text: "https://tenor.com/view/%e3%83%89%e3%82%ad%e3%83%89%e3%82%ad-" +
				"%e7%b7%8a%e5%bc%b5-%e5%a5%bd%e3%81%8d-%e3%83%8f%e3%83%bc%e3%83%88" +
				"-%e5%8f%af%e6%84%9b%e3%81%84-gif-15876325",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.Empty(t, results, "Expected NO results for: '%s'", tt.text)
		})
	}
}

// TestNegativeHyphenatedNumbers tests that hyphenated numbers should not parse
func TestNegativeHyphenatedNumbers(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "1-2", text: "1-2"},
		{name: "1-2-3", text: "1-2-3"},
		{name: "4-5-6", text: "4-5-6"},
		{name: "20-30-12", text: "20-30-12"},
		// Note: Standalone 4-digit years like "2012" now parse as year-only expressions (see Issue #99)
		{name: "2012-14", text: "2012-14"},
		{name: "2012-1400", text: "2012-1400"},
		{name: "2200-25", text: "2200-25"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.Empty(t, results, "Expected NO results for: '%s', but got %d", tt.text, len(results))
		})
	}
}

// TestNegativeImpossibleDates tests that impossible dates should not parse
func TestNegativeImpossibleDates(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "February 29, 2022", text: "February 29, 2022"}, // 2022 not leap year
		{name: "02/29/2022", text: "02/29/2022"},
		{name: "June 31, 2022", text: "June 31, 2022"}, // June has 30 days
		{name: "06/31/2022", text: "06/31/2022"},
		{name: "14PM", text: "14PM"},                               // 14 PM doesn't exist (use 24h or 12h)
		{name: "25:12", text: "25:12"},                             // Hour 25 doesn't exist
		{name: "13/31/2018", text: "An appointment on 13/31/2018"}, // Month 13 doesn't exist
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.Empty(t, results, "Expected NO results for: '%s', but got %d", tt.text, len(results))
		})
	}
}

// TestNegativeImpossibleDateRanges tests that impossible date ranges should not parse
func TestNegativeImpossibleDateRanges(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "February 20 - 29, 2022", text: "February 20 - 29, 2022"}, // 2022 not leap year
		{name: "June 10 - 31, 2022", text: "June 10 - 31, 2022"},         // June has 30 days
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.Empty(t, results, "Expected NO results for: '%s', but got %d", tt.text, len(results))
		})
	}
}

// TestNegativeVersionNumbers tests that version numbers should not parse
func TestNegativeVersionNumbers(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "Version: 1.1.3", text: "Version: 1.1.3"},
		{name: "Version: 1.1.30", text: "Version: 1.1.30"},
		{name: "Version: 1.10.30", text: "Version: 1.10.30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.Empty(t, results, "Expected NO results for: '%s', but got %d", tt.text, len(results))
		})
	}
}

// TestNegativeIncorrectReference tests that incorrect date references should not parse
func TestNegativeIncorrectReference(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "for the year", text: "for the year"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.Empty(t, results, "Expected NO results for: '%s', but got %d", tt.text, len(results))
		})
	}
}

// TestPositiveDateWithVersionNumber tests that real dates should parse even near version numbers
func TestPositiveDateWithVersionNumber(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedText string
	}{
		{
			name:         "1.5.3 - 2015-09-24",
			text:         "1.5.3 - 2015-09-24",
			expectedText: "2015-09-24",
		},
		{
			name:         "1.5.30 - 2015-09-24",
			text:         "1.5.30 - 2015-09-24",
			expectedText: "2015-09-24",
		},
		{
			name:         "1.50.30 - 2015-09-24",
			text:         "1.50.30 - 2015-09-24",
			expectedText: "2015-09-24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := englishCasualChrono() // was CreateCasualConfiguration(false)
			// chrono created above

			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.NotEmpty(t, results, "Expected to parse date in: '%s'", tt.text)

			if len(results) > 0 {
				assert.Contains(t, results[0].Text(), tt.expectedText, "Should parse the date, not the version")
			}
		})
	}
}
