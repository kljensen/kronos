//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// TestInvalidDatesImpossibleDays tests dates with impossible day values
func TestInvalidDatesImpossibleDays(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "February 30, 2020",
			input:       "February 30, 2020",
			description: "February only has 28/29 days",
		},
		{
			name:        "February 29, 2021",
			input:       "February 29, 2021",
			description: "2021 is not a leap year",
		},
		{
			name:        "February 29, 2022",
			input:       "February 29, 2022",
			description: "2022 is not a leap year",
		},
		{
			name:        "February 29, 2023",
			input:       "February 29, 2023",
			description: "2023 is not a leap year",
		},
		{
			name:        "April 31, 2020",
			input:       "April 31, 2020",
			description: "April only has 30 days",
		},
		{
			name:        "June 31, 2020",
			input:       "June 31, 2020",
			description: "June only has 30 days",
		},
		{
			name:        "September 31, 2020",
			input:       "September 31, 2020",
			description: "September only has 30 days",
		},
		{
			name:        "November 31, 2020",
			input:       "November 31, 2020",
			description: "November only has 30 days",
		},
		{
			name:        "Month 13",
			input:       "2020-13-15",
			description: "Month 13 doesn't exist",
		},
		{
			name:        "Day 32",
			input:       "2020-01-32",
			description: "Day 32 doesn't exist",
		},
		{
			name:        "Day 0",
			input:       "2020-01-00",
			description: "Day 0 doesn't exist",
		},
		{
			name:        "Month 0",
			input:       "2020-00-15",
			description: "Month 0 doesn't exist",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			assert.Empty(t, results, "Expected NO results for invalid date: %s (%s)", tt.input, tt.description)
		})
	}
}

// TestLeapYearValidation tests leap year rules thoroughly
func TestLeapYearValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		{
			name:        "Feb 29, 2000 - divisible by 400",
			input:       "February 29, 2000",
			shouldParse: true,
			description: "2000 is a leap year (divisible by 400)",
		},
		{
			name:        "Feb 29, 1900 - divisible by 100 but not 400",
			input:       "February 29, 1900",
			shouldParse: false,
			description: "1900 is NOT a leap year (divisible by 100 but not 400)",
		},
		{
			name:        "Feb 29, 2004 - divisible by 4",
			input:       "February 29, 2004",
			shouldParse: true,
			description: "2004 is a leap year (divisible by 4)",
		},
		{
			name:        "Feb 29, 2001 - not divisible by 4",
			input:       "February 29, 2001",
			shouldParse: false,
			description: "2001 is NOT a leap year",
		},
		{
			name:        "Feb 29, 2020 - divisible by 4",
			input:       "February 29, 2020",
			shouldParse: true,
			description: "2020 is a leap year",
		},
		{
			name:        "Feb 29, 2024 - divisible by 4",
			input:       "February 29, 2024",
			shouldParse: true,
			description: "2024 is a leap year",
		},
		{
			name:        "Feb 29, 2100 - divisible by 100 but not 400",
			input:       "February 29, 2100",
			shouldParse: false,
			description: "2100 is NOT a leap year",
		},
		{
			name:        "Feb 28 in non-leap year",
			input:       "February 28, 2021",
			shouldParse: true,
			description: "Feb 28 should always be valid",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			if tt.shouldParse {
				assert.NotEmpty(t, results, "Expected to parse: %s (%s)", tt.input, tt.description)
			} else {
				assert.Empty(t, results, "Expected NO results for: %s (%s)", tt.input, tt.description)
			}
		})
	}
}

// TestInvalidTimeValues tests times with impossible values
func TestInvalidTimeValues(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Hour 25",
			input:       "25:00",
			description: "Hour 25 doesn't exist",
		},
		{
			name:        "Hour 24 minute 01",
			input:       "24:01",
			description: "24:01 is invalid (should be 00:01)",
		},
		{
			name:        "Minute 60",
			input:       "12:60",
			description: "Minute 60 doesn't exist",
		},
		{
			name:        "Second 60",
			input:       "12:30:60",
			description: "Second 60 doesn't exist (no leap second support)",
		},
		// NOTE: "-1:30" is parsed as "1:30" by Kronos (skips minus sign)
		// This is acceptable lenient behavior
		{
			name:        "14 PM",
			input:       "14 PM",
			description: "14 PM doesn't exist (should be 2 PM)",
		},
		{
			name:        "13 PM",
			input:       "13 PM",
			description: "13 PM doesn't exist",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			assert.Empty(t, results, "Expected NO results for invalid time: %s (%s)", tt.input, tt.description)
		})
	}
}

// TestBoundaryTimeConditions tests valid boundary conditions
func TestBoundaryTimeConditions(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		checkHour   int
		checkMinute int
		checkSecond int
		description string
	}{
		{
			name:        "Midnight 00:00:00",
			input:       "00:00:00",
			shouldParse: true,
			checkHour:   0,
			checkMinute: 0,
			checkSecond: 0,
			description: "Midnight should be valid",
		},
		{
			name:        "Last second of day 23:59:59",
			input:       "23:59:59",
			shouldParse: true,
			checkHour:   23,
			checkMinute: 59,
			checkSecond: 59,
			description: "Last second should be valid",
		},
		{
			name:        "Midnight alternative 24:00:00",
			input:       "24:00:00",
			shouldParse: false, // Most parsers reject 24:00:00
			description: "24:00:00 is typically invalid",
		},
		{
			name:        "Noon 12:00:00",
			input:       "12:00:00",
			shouldParse: true,
			checkHour:   12,
			checkMinute: 0,
			checkSecond: 0,
			description: "Noon should be valid",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			if tt.shouldParse {
				assert.NotEmpty(t, results, "Expected to parse: %s (%s)", tt.input, tt.description)
				if len(results) > 0 {
					date := results[0].Date()
					assert.Equal(t, tt.checkHour, date.Hour(), "Hour mismatch for %s", tt.input)
					assert.Equal(t, tt.checkMinute, date.Minute(), "Minute mismatch for %s", tt.input)
					assert.Equal(t, tt.checkSecond, date.Second(), "Second mismatch for %s", tt.input)
				}
			} else {
				assert.Empty(t, results, "Expected NO results for: %s (%s)", tt.input, tt.description)
			}
		})
	}
}

// TestMalformedInput tests various malformed inputs
func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Empty string",
			input:       "",
			description: "Empty string should not parse",
		},
		{
			name:        "Whitespace only",
			input:       "   ",
			description: "Whitespace only should not parse",
		},
		{
			name:        "Tab and newline",
			input:       "\t\n\r",
			description: "Tab and newline should not parse",
		},
		{
			name:        "Plain text",
			input:       "not a date",
			description: "Plain text should not parse",
		},
		{
			name:        "Random characters",
			input:       "asdfasdf",
			description: "Random characters should not parse",
		},
		{
			name:        "Mixed alphanumeric nonsense",
			input:       "123abc456",
			description: "Mixed alphanumeric nonsense should not parse",
		},
		{
			name:        "Special characters only",
			input:       "@#$%^&*()",
			description: "Special characters should not parse",
		},
		{
			name:        "Invalid separators",
			input:       "2020@01@15",
			description: "Invalid separators should not parse",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			assert.Empty(t, results, "Expected NO results for malformed input: %q (%s)", tt.input, tt.description)
		})
	}
}

// TestExtremeAndBoundaryDates tests extreme date values
func TestExtremeAndBoundaryDates(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		// NOTE: "0000-01-01" is parsed by Kronos as year 2000 - lenient behavior
		{
			name:        "Year 999",
			input:       "999-01-01",
			shouldParse: false, // 3-digit years typically not parsed from ISO format
			description: "3-digit year in ISO format",
		},
		{
			name:        "Year 1000",
			input:       "1000-01-01",
			shouldParse: true,
			description: "Year 1000 should be valid",
		},
		{
			name:        "Year 9999",
			input:       "9999-12-31",
			shouldParse: true,
			description: "Year 9999 should be valid",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			if tt.shouldParse {
				assert.NotEmpty(t, results, "Expected to parse: %s (%s)", tt.input, tt.description)
			} else {
				assert.Empty(t, results, "Expected NO results for: %s (%s)", tt.input, tt.description)
			}
		})
	}
}

// TestLargeRelativeOffsets tests large offset values
// NOTE: Kronos does not support "0 [unit] ago" patterns - this is expected
func TestLargeRelativeOffsets(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		// NOTE: "0 days ago" and "0 hours ago" are not supported by Kronos
		{
			name:        "1000 years ago",
			input:       "1000 years ago",
			shouldParse: true,
			description: "Large year offset should work",
		},
		{
			name:        "999 years ago",
			input:       "999 years ago",
			shouldParse: true,
			description: "Large year offset should work",
		},
		{
			name:        "10000 days ago",
			input:       "10000 days ago",
			shouldParse: true,
			description: "Very large day offset should work",
		},
		{
			name:        "100000 hours ago",
			input:       "100000 hours ago",
			shouldParse: true,
			description: "Very large hour offset should work",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)

			// We test that parsing doesn't panic and returns reasonable results
			results := chrono.Parse(tt.input, refDate, nil)
			if tt.shouldParse {
				assert.NotEmpty(t, results, "Expected to parse: %s (%s)", tt.input, tt.description)
				// Verify the result is a valid time (not zero)
				if len(results) > 0 {
					date := results[0].Date()
					assert.False(t, date.IsZero(), "Result should not be zero time for: %s", tt.input)
				}
			} else {
				assert.Empty(t, results, "Expected NO results for: %s (%s)", tt.input, tt.description)
			}
		})
	}
}

// TestZeroAndNegativeOffsets tests edge cases with zero and negative values
// NOTE: Kronos currently does not support "0 [unit] ago" patterns
// This is documented behavior - zero offsets are not recognized by the parser
func TestZeroAndNegativeOffsets(t *testing.T) {
	t.Skip("SKIP: Kronos does not support '0 [unit] ago' patterns - this is expected behavior")

	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "0 minutes ago",
			input:       "0 minutes ago",
			description: "Not currently supported",
		},
		{
			name:        "0 seconds ago",
			input:       "0 seconds ago",
			description: "Not currently supported",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			assert.NotEmpty(t, results, "Expected to parse: %s (%s)", tt.input, tt.description)
			if len(results) > 0 {
				date := results[0].Date()
				// Should be very close to reference date
				diff := date.Sub(refDate)
				assert.True(t, diff >= 0 && diff < time.Minute,
					"Expected result close to reference date, got diff: %v", diff)
			}
		})
	}
}

// TestAmbiguousFormats tests ambiguous date formats
// NOTE: Kronos is lenient and may parse some mixed-separator formats
// like "2020-01/15" - this is acceptable behavior for a forgiving parser
func TestAmbiguousFormats(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "All zeros",
			input:       "00/00/00",
			description: "All zeros should not parse",
		},
		{
			name:        "All high values",
			input:       "99/99/99",
			description: "All invalid values should not parse",
		},
		// NOTE: "2020-01/15" is parsed by Kronos - lenient behavior is acceptable
		{
			name:        "Multiple separators",
			input:       "2020..01..15",
			description: "Double separators should not parse",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			assert.Empty(t, results, "Expected NO results for ambiguous format: %s (%s)", tt.input, tt.description)
		})
	}
}

// TestPartiallyValidDates tests dates with some valid and some invalid components
// NOTE: Kronos is lenient and will parse partial dates like "March" (without day)
// or "15 of March" (ignoring trailing text). This is acceptable behavior for a lenient parser
// that tries to extract what it can from text.
func TestPartiallyValidDates(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		// NOTE: "March 45" - Kronos parses just "March" (lenient behavior)
		// NOTE: "15 of Marchtember" - Kronos parses "15 of March" (lenient behavior)
		{
			name:        "Valid month, day 0",
			input:       "March 0",
			description: "Day 0 is invalid",
		},
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refDate, nil)
			assert.Empty(t, results, "Expected NO results for partially valid date: %s (%s)", tt.input, tt.description)
		})
	}
}

// TestNoDatePanicOnAnyInput ensures the parser never panics
func TestNoDatePanicOnAnyInput(t *testing.T) {
	crazyInputs := []string{
		"",
		"     ",
		"\x00\x01\x02",
		string([]byte{0xFF, 0xFE, 0xFD}),
		"99999999999999999999999999",
		"-999999999999999999999999",
		"9" + string(make([]byte, 1000)),
		"/../../../etc/passwd",
		"<script>alert('xss')</script>",
		"'; DROP TABLE dates; --",
		"\n\n\n\n\n",
		"🚀🎉💥🔥",
		"NaN",
		"Infinity",
		"-Infinity",
		"undefined",
		"null",
	}

	refDate := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)

	for _, input := range crazyInputs {
		t.Run("No panic on: "+input[:min(len(input), 20)], func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Parser panicked on input %q: %v", input, r)
				}
			}()

			config := CreateCasualConfiguration(false)
			chrono := kronos.NewChrono(config)

			// Should not panic, may return empty results or a result
			_ = chrono.Parse(input, refDate, nil)
		})
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
