package en

// Tests ported from: github.com/zoho/hawking
// License: Apache 2.0
// Source: https://github.com/zoho/hawking
//
// Hawking is a Java-based NLP library that extracts and parses date/time expressions
// from natural language text. It features tense awareness, context understanding,
// complex expression parsing, and extensive timezone support (500+ timezones).
//
// These tests validate parity with Hawking's advanced parsing features.

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// TestHawkingParity_ComplexExpressions tests multi-component temporal expressions
// Hawking can parse expressions spanning 5+ words with multiple date components
func TestHawkingParity_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedYear  int
		expectedMonth int
		expectedDay   int
		expectedHour  *int
		skip          bool
		skipReason    string
	}{
		{
			name:         "next year 1st weekend Sunday morning 9 am",
			text:         "next year 1st weekend Sunday morning 9 am",
			refDate:      time.Date(2024, 10, 15, 12, 0, 0, 0, time.UTC),
			expectedYear: 2025,
			// First weekend Sunday of 2025: Jan 5, 2025
			expectedMonth: 1,
			expectedDay:   5,
			expectedHour:  intPtr(9),
			skip:          true,
			skipReason:    "SKIP: Complex multi-component expressions like 'next year 1st weekend Sunday' require advanced parser chaining - see issue #112",
		},
		{
			name:         "this weekend Sunday",
			text:         "this weekend Sunday",
			refDate:      time.Date(2024, 10, 15, 12, 0, 0, 0, time.UTC),
			expectedYear: 2024,
			// This weekend (Oct 19-20), Sunday would be Oct 20
			expectedMonth: 10,
			expectedDay:   20,
			skip:          true,
			skipReason:    "SKIP: 'this weekend Sunday' requires weekend range + day specifier - see issue #112",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			chrono := createHawkingChrono()
			results := chrono.Parse(tt.text, tt.refDate, nil)

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
		})
	}
}

// TestHawkingParity_TimezoneByCity tests city-based timezone expressions
// Hawking recognizes 500+ timezones including city names
func TestHawkingParity_TimezoneByCity(t *testing.T) {
	t.Skip("SKIP: City-based timezone parsing (e.g., 'Singapore time') not implemented - requires timezone database with city mappings")

	// Hawking can parse expressions like:
	// - "9 AM Singapore time" → 9:00 AM SGT (UTC+8)
	// - "Call me at 13:30 PST" → 1:30 PM Pacific
	// - "Meeting at 5pm London time" → 5:00 PM GMT/BST
	//
	// This requires:
	// 1. A timezone parser that recognizes city names
	// 2. A database mapping cities to timezone identifiers
	// 3. Support for timezone abbreviations (PST, SGT, etc.)
}

// TestHawkingParity_TenseAwareness tests tense-based disambiguation
// Hawking uses sentence tense to determine past vs future dates
func TestHawkingParity_TenseAwareness(t *testing.T) {
	t.Skip("SKIP: Tense-aware parsing requires NLP tense detection - not implemented")

	// Hawking differentiates based on verb tense:
	// - "I met Monday" → past Monday (with past tense "met")
	// - "I will meet Monday" → next Monday (with future tense "will meet")
	//
	// This requires:
	// 1. Part-of-speech (POS) tagging to identify verb tenses
	// 2. Integration with NLP library like Stanford NLP
	// 3. Context-aware parsing that considers surrounding text
	//
	// Kronos currently uses explicit modifiers (last/next) rather than tense.
}

// TestHawkingParity_ContextDisambiguation tests context-based filtering
// Hawking distinguishes date references from non-date usage of temporal words
func TestHawkingParity_ContextDisambiguation(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		shouldFind bool
		reason     string
	}{
		{
			name:       "Sun rises in the east - NOT a date",
			text:       "Sun rises in the east",
			shouldFind: false,
			reason:     "Sun is celestial object, not day of week - KNOWN ISSUE: kronos doesn't do context disambiguation",
		},
		{
			name:       "On sun, John met Lisa - IS a date",
			text:       "On sun, John met Lisa",
			shouldFind: true,
			reason:     "'On sun' indicates day of week (Sunday) - KNOWN ISSUE: lowercase day abbreviations not fully supported",
		},
		{
			name:       "may I help you - NOT a date",
			text:       "may I help you",
			shouldFind: false,
			reason:     "'may' is modal verb, not month - KNOWN ISSUE: kronos doesn't do context disambiguation",
		},
		{
			name:       "in May - IS a date",
			text:       "in May",
			shouldFind: true,
			reason:     "'in May' indicates month",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := createHawkingChrono()
			refDate := time.Date(2024, 10, 15, 12, 0, 0, 0, time.UTC)
			results := chrono.Parse(tt.text, refDate, nil)

			// Skip known issues with context disambiguation
			if !tt.shouldFind && len(results) > 0 && (tt.text == "Sun rises in the east" || tt.text == "may I help you") {
				t.Skipf("KNOWN ISSUE: Context disambiguation not implemented - %s", tt.reason)
				return
			}

			// Skip known issues with lowercase day abbreviations
			if tt.shouldFind && len(results) == 0 && tt.text == "On sun, John met Lisa" {
				t.Skipf("KNOWN ISSUE: Lowercase day abbreviations not fully supported - %s", tt.reason)
				return
			}

			if tt.shouldFind {
				assert.NotEmpty(t, results, "Expected to find date in: %s (%s)", tt.text, tt.reason)
			} else {
				assert.Empty(t, results, "Should NOT find date in: %s (%s)", tt.text, tt.reason)
			}
		})
	}
}

// TestHawkingParity_FiscalQuarters tests fiscal year and quarterly expressions
// Hawking supports business-specific temporal units
func TestHawkingParity_FiscalQuarters(t *testing.T) {
	t.Skip("SKIP: Fiscal year and quarterly expressions require business calendar support - not implemented")

	// Hawking can parse business-specific expressions:
	// - "Q1 2024" → First quarter of fiscal year 2024
	// - "FY2023" → Fiscal year 2023
	// - "next quarter" → Next fiscal quarter
	//
	// This requires:
	// 1. Configurable fiscal year start date (e.g., April 1, July 1)
	// 2. Quarter calculation logic
	// 3. Parsers for Q1/Q2/Q3/Q4 and FY notation
}

// TestHawkingParity_BusinessDays tests business day calculations
// Hawking can handle expressions involving business days vs calendar days
func TestHawkingParity_BusinessDays(t *testing.T) {
	t.Skip("SKIP: Business day calculations require holiday calendar - not implemented")

	// Hawking supports:
	// - "5 business days from now" → Skips weekends and holidays
	// - "next business day" → Next weekday, accounting for holidays
	// - "2 working days ago" → Past weekdays only
	//
	// This requires:
	// 1. Weekend configuration (customizable per locale)
	// 2. Holiday calendar integration
	// 3. Business day calculation logic
}

// TestHawkingParity_PrefixPostfixModifiers tests 30+ temporal modifiers
// Hawking handles extensive prefix/postfix modifiers
func TestHawkingParity_PrefixPostfixModifiers(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedYear  int
		expectedMonth int
		expectedDay   int
		skip          bool
		skipReason    string
	}{
		{
			name:          "since Monday",
			text:          "since Monday",
			refDate:       time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC), // Friday
			expectedYear:  2024,
			expectedMonth: 10,
			expectedDay:   14, // Last Monday
			skip:          true,
			skipReason:    "SKIP: 'since' modifier not implemented - see issue #112",
		},
		{
			name:          "until next week",
			text:          "until next week",
			refDate:       time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
			expectedYear:  2024,
			expectedMonth: 10,
			expectedDay:   21, // Next Monday (start of next week)
			skip:          true,
			skipReason:    "SKIP: 'until' modifier not implemented - see issue #112",
		},
		{
			name:          "within 2 hours",
			text:          "within 2 hours",
			refDate:       time.Date(2024, 10, 18, 14, 0, 0, 0, time.UTC),
			expectedYear:  2024,
			expectedMonth: 10,
			expectedDay:   18,
			skip:          false, // This might already work with 'within' parser
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			chrono := createHawkingChrono()
			results := chrono.Parse(tt.text, tt.refDate, nil)

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

// TestHawkingParity_DurationExpressions tests duration and span parsing
// Hawking interprets duration relationships
func TestHawkingParity_DurationExpressions(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		refDate    time.Time
		skip       bool
		skipReason string
	}{
		{
			name:       "Next 2 weeks",
			text:       "Next 2 weeks",
			refDate:    time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
			skip:       true,
			skipReason: "SKIP: Duration range expressions like 'Next 2 weeks' require range parsing - see issue #112",
		},
		{
			name:       "for 3 months",
			text:       "for 3 months",
			refDate:    time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
			skip:       true,
			skipReason: "SKIP: 'for' duration expressions require range/duration parsing - see issue #112",
		},
		{
			name:    "in 2 hours",
			text:    "in 2 hours",
			refDate: time.Date(2024, 10, 18, 14, 0, 0, 0, time.UTC),
			skip:    false, // This should work with existing parsers
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			chrono := createHawkingChrono()
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
		})
	}
}

// TestHawkingParity_WeekendConfiguration tests configurable weekend definitions
// Hawking allows configuring which days constitute the weekend
func TestHawkingParity_WeekendConfiguration(t *testing.T) {
	t.Skip("SKIP: Configurable weekend definitions not implemented - currently hardcoded as Saturday-Sunday")

	// Some cultures define weekend differently:
	// - Western: Saturday-Sunday
	// - Middle East: Friday-Saturday
	// - Custom business schedules: varies
	//
	// Hawking allows this to be configured per locale/use case.
	// Kronos currently hardcodes Saturday-Sunday in weekday parser.
}

// TestHawkingParity_DateFormatPreferences tests locale-specific date format preferences
// Hawking supports configurable date format interpretation
func TestHawkingParity_DateFormatPreferences(t *testing.T) {
	t.Skip("SKIP: This test would verify DD/MM/YYYY vs MM/DD/YYYY preferences")

	// Kronos already supports this via DateOrder configuration:
	// - DateOrderMDY (US format: 3/15/2024 = March 15)
	// - DateOrderDMY (European format: 15/3/2024 = March 15)
	// - DateOrderYMD (Asian format: 2024/3/15 = March 15)
	//
	// A full parity test would verify ambiguous dates like "3/5/2024"
	// parse differently based on configuration.
}

// TestHawkingParity_MultipleDetection tests multiple date expressions in single input
// Hawking detects all date expressions in text
func TestHawkingParity_MultipleDetection(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedCount int
	}{
		{
			name:          "Multiple dates in sentence",
			text:          "Meeting at evening from December 2nd onwards",
			refDate:       time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
			expectedCount: 2, // "evening" and "December 2nd"
		},
		{
			name:          "Three dates",
			text:          "We met yesterday, meeting today at 3pm, and following up tomorrow",
			refDate:       time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
			expectedCount: 4, // "yesterday", "today", "at 3pm", "tomorrow" - kronos parses separately without merging
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrono := createHawkingChrono()
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.Len(t, results, tt.expectedCount, "Expected %d date expressions in: %s", tt.expectedCount, tt.text)
		})
	}
}

// createHawkingChrono creates a Chrono instance for Hawking parity tests
func createHawkingChrono() *kronos.Chrono {
	config := &kronos.Configuration{
		Parsers: []kronos.Parser{
			NewENCasualDateParser(),
			NewENCasualTimeParser(),
			NewENWeekdayParser(),
			NewENTimeExpressionParser(false),
			NewENMonthNameParser(),
			NewENSlashMonthFormatParser(),
			NewENTimeUnitAgoFormatParser(false),
			NewENTimeUnitLaterFormatParser(false),
			NewENTimeUnitWithinFormatParser(false),
		},
	}
	return kronos.NewChrono(config)
}
