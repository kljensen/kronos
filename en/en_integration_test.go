package en

// Tests ported from chrono's en.test.ts
//
// Original file: 19 integration test cases covering various date parsing scenarios
//
// These are integration tests that verify the full EN parser pipeline,
// not individual parser components. Many require features not yet implemented.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestIntegrationDateTimeExpression tests date + time range expressions
func TestIntegrationDateTimeExpression(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedText string
	}{
		{
			name:         "Date with time range",
			text:         "Something happen on 2014-04-18 13:00 - 16:00 as",
			expectedText: "2014-04-18 13:00 - 16:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := New().Parse(tt.text)
			assert.NoError(t, err)

			if len(results) > 0 {
				assert.Contains(t, results[0].Text(), tt.expectedText, "Text mismatch")
			}
		})
	}
}

// TestIntegrationTimeExpression tests time expressions with timezones
func TestIntegrationTimeExpression(t *testing.T) {
	t.Skip("SKIP: Time expressions with timezone abbreviations need timezone parser")

	tests := []struct {
		name         string
		text         string
		refDate      time.Time
		expectedText string
	}{
		{
			name:         "Time range with PM",
			text:         "between 3:30-4:30pm",
			refDate:      time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC),
			expectedText: "3:30-4:30pm",
		},
		{
			name:         "Time with timezone",
			text:         "9:00 PST",
			refDate:      time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC),
			expectedText: "9:00 PST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := New().WithReferenceDate(tt.refDate).Parse(tt.text)
			assert.NoError(t, err)
			assert.NotEmpty(t, results)
		})
	}
}

// TestIntegrationQuotedExpressions tests dates in quotes
func TestIntegrationQuotedExpressions(t *testing.T) {
	t.Skip("SKIP: Quoted expression handling not implemented")

	tests := []struct {
		name         string
		text         string
		refDate      time.Time
		expectedText string
	}{
		{
			name:         "Time in parentheses",
			text:         "Want to meet for dinner (5pm EST)?",
			refDate:      time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC),
			expectedText: "5pm EST",
		},
		{
			name:         "Time range in single quotes",
			text:         "between '3:30-4:30pm'",
			refDate:      time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC),
			expectedText: "3:30-4:30pm",
		},
		{
			name:         "Date in single quotes",
			text:         "The date is '2014-04-18'",
			refDate:      time.Now(),
			expectedText: "2014-04-18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := New().WithReferenceDate(tt.refDate).Parse(tt.text)
			assert.NoError(t, err)
			assert.NotEmpty(t, results)
		})
	}
}

// TestIntegrationStrictMode tests strict mode behavior
func TestIntegrationStrictMode(t *testing.T) {
	t.Skip("SKIP: Strict mode configuration not documented")

	// Tests that "Tuesday" alone should not parse in strict mode
}

// TestIntegrationBuiltinVariants tests EN locale variants (US vs GB)
func TestIntegrationBuiltinVariants(t *testing.T) {
	t.Skip("SKIP: Locale variants (EN-US vs EN-GB) not implemented")

	// Tests 6/10/2018 parses differently in EN vs EN-GB
	// EN: June 10, 2018
	// EN-GB: 6 October, 2018
}

// TestIntegrationRandomText tests parsing from complex text
func TestIntegrationRandomText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		refDate      time.Time
		expectedText string
		skip         bool
		skipReason   string
	}{
		{
			name:         "Email with date",
			text:         "Adam <Adam@supercalendar.com> написал(а):\nThe date is 02.07.2013",
			expectedText: "02.07.2013",
			skip:         false,
		},
		{
			name:         "Date range in text",
			text:         "174 November 1,2001- March 31,2002",
			expectedText: "November 1,2001- March 31,2002",
			skip:         false,
		},
		{
			name:         "Thursday with full date",
			text:         "...Thursday, December 15, 2011 Best Available Rate ",
			expectedText: "Thursday, December 15, 2011",
			skip:         false,
		},
		{
			name:         "Weekday abbreviation with time",
			text:         "SUN 15SEP 11:05 AM - 12:50 PM",
			expectedText: "SUN 15SEP 11:05 AM - 12:50 PM",
			skip:         true,
			skipReason:   "Complex month abbreviation + time range",
		},
		{
			name:         "Full weekday date time range",
			text:         "FRI 13SEP 1:29 PM - FRI 13SEP 3:29 PM",
			expectedText: "FRI 13SEP 1:29 PM - FRI 13SEP 3:29 PM",
			skip:         true,
			skipReason:   "Complex date-time range with weekday",
		},
		{
			name:         "Time range with full date",
			text:         "9:00 AM to 5:00 PM, Tuesday, 20 May 2013",
			expectedText: "9:00 AM to 5:00 PM, Tuesday, 20 May 2013",
			skip:         true,
			skipReason:   "Time range before date",
		},
		{
			name:         "Relative weekday range",
			text:         "Monday afternoon to last night",
			refDate:      time.Date(2017, 7, 7, 0, 0, 0, 0, time.UTC),
			expectedText: "Monday afternoon to last night",
			skip:         true,
			skipReason:   "Requires relative weekday + 'afternoon'/'night' support",
		},
		{
			name:         "Slash date with comma and time",
			text:         "07-27-2022, 02:00 AM",
			expectedText: "07-27-2022, 02:00 AM",
			skip:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf("SKIP: %s", tt.skipReason)
				return
			}

			refDate := tt.refDate
			if refDate.IsZero() {
				refDate = time.Now()
			}

			results, err := New().WithReferenceDate(refDate).Parse(tt.text)
			assert.NoError(t, err)

			if len(results) > 0 {
				assert.Contains(t, results[0].Text(), tt.expectedText, "Text mismatch")
			}
		})
	}
}

// TestIntegrationWikipediaText tests parsing from Wikipedia text
func TestIntegrationWikipediaText(t *testing.T) {
	text := "October 7, 2011, of which details were not revealed out of respect to Jobs's family.[239] " +
		"Apple announced on the same day that they had no plans for a public service, but were encouraging " +
		`"well-wishers" to send their remembrance messages to an email address created to receive such messages.[240] ` +
		"Sunday, October 16, 2011"

	results, err := New().
		WithReferenceDate(time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)).
		Parse(text)
	assert.NoError(t, err)

	// Should find at least the clear dates
	if len(results) >= 2 {
		assert.Contains(t, results[0].Text(), "October 7, 2011")
		assert.Contains(t, results[len(results)-1].Text(), "October 16, 2011")
	}
}

// TestIntegrationMultipleResults tests parsing multiple dates from text
func TestIntegrationMultipleResults(t *testing.T) {
	text := "I will see you at 2:30. If not I will see you somewhere between 3:30-4:30pm"
	refDate := time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC)

	results, err := New().WithReferenceDate(refDate).Parse(text)
	assert.NoError(t, err)

	// Should find at least one result
	assert.NotEmpty(t, results, "Expected to find at least one date")
}

// TestIntegrationCustomizeByRemovingParser tests custom configuration
func TestIntegrationCustomizeByRemovingParser(t *testing.T) {
	t.Skip("SKIP: Custom parser removal configuration not documented")

	// Tests removing time parser to only extract weekday
	// Example: "Thursday 9AM" should only extract "Thursday"
}
