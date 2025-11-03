package en

import (
	"testing"
	"time"

	"github.com/kljensen/kronos"
)

// TestPeriodIntegrationWithParser tests period tracking with the actual parsing system
func TestPeriodIntegrationWithParser(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		input          string
		expectedPeriod kronos.Period
		description    string
	}{
		// Time-level expressions
		{"10:30", kronos.PeriodTime, "absolute time"},
		{"3pm", kronos.PeriodTime, "time with meridiem"},
		{"2 hours ago", kronos.PeriodTime, "relative time"},
		{"now", kronos.PeriodTime, "now keyword"},
		{"noon", kronos.PeriodTime, "noon"},
		{"midnight", kronos.PeriodTime, "midnight"},
		{"this morning", kronos.PeriodTime, "morning"},
		{"this afternoon", kronos.PeriodTime, "afternoon"},
		{"this evening", kronos.PeriodTime, "evening"},

		// Day-level expressions
		{"yesterday", kronos.PeriodDay, "yesterday"},
		{"tomorrow", kronos.PeriodDay, "tomorrow"},
		{"today", kronos.PeriodDay, "today"},
		{"3 days ago", kronos.PeriodDay, "days ago"},
		{"March 15", kronos.PeriodDay, "month and day"},
		{"March 15, 2020", kronos.PeriodDay, "full date"},
		{"2020-03-15", kronos.PeriodDay, "ISO date"},
		{"tonight", kronos.PeriodDay, "tonight"},
		{"last night", kronos.PeriodDay, "last night"},

		// Week-level expressions
		{"last week", kronos.PeriodWeek, "last week"},
		{"next week", kronos.PeriodWeek, "next week"},
		{"this week", kronos.PeriodWeek, "this week"},
		{"2 weeks ago", kronos.PeriodWeek, "weeks ago"},
		{"Monday", kronos.PeriodWeek, "weekday"},
		{"next Friday", kronos.PeriodWeek, "next weekday"},

		// Month-level expressions
		{"March", kronos.PeriodMonth, "month name only"},
		{"March 2020", kronos.PeriodMonth, "month and year"},
		{"2 months ago", kronos.PeriodMonth, "months ago"},
		{"11/2005", kronos.PeriodMonth, "MM/YYYY format"},
		{"last month", kronos.PeriodMonth, "last month"},

		// Year-level expressions
		{"2020", kronos.PeriodYear, "year only"},
		{"last year", kronos.PeriodYear, "last year"},
		{"next year", kronos.PeriodYear, "next year"},
		{"1 year ago", kronos.PeriodYear, "year ago"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Parse using the English parser
			results := Parse(tt.input, refTime, nil)

			if len(results) == 0 {
				t.Fatalf("No results parsed from %q", tt.input)
			}

			result := results[0]
			components, ok := kronos.AsParsingComponents(result.Start())
			if !ok {
				t.Fatalf("Could not convert result to ParsingComponents")
			}

			period := components.Period()
			if period != tt.expectedPeriod {
				t.Errorf("Input %q: period = %v (%s), want %v (%s)",
					tt.input, period, period.String(), tt.expectedPeriod, tt.expectedPeriod.String())
			}
		})
	}
}

// TestPeriodWithCombinedExpressions tests period for combined date+time expressions
func TestPeriodWithCombinedExpressions(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		input          string
		expectedPeriod kronos.Period
		description    string
	}{
		{"yesterday at 11:30", kronos.PeriodTime, "day + time should be time-level"},
		{"March 15 at 3pm", kronos.PeriodTime, "date + time should be time-level"},
		{"tomorrow morning", kronos.PeriodTime, "tomorrow + time of day should be time-level"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			results := Parse(tt.input, refTime, nil)

			if len(results) == 0 {
				t.Fatalf("No results parsed from %q", tt.input)
			}

			result := results[0]
			components, ok := kronos.AsParsingComponents(result.Start())
			if !ok {
				t.Fatalf("Could not convert result to ParsingComponents")
			}

			period := components.Period()
			if period != tt.expectedPeriod {
				t.Errorf("Input %q: period = %v (%s), want %v (%s)",
					tt.input, period, period.String(), tt.expectedPeriod, tt.expectedPeriod.String())
			}
		})
	}
}

// TestPeriodConsistencyAcrossParsers tests that similar expressions have consistent periods
func TestPeriodConsistencyAcrossParsers(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)

	// All these should be day-level
	dayExpressions := []string{
		"March 15",
		"15th March",
		"03/15/2020",
		"2020-03-15",
		"yesterday",
		"tomorrow",
		"today",
	}

	for _, input := range dayExpressions {
		t.Run(input, func(t *testing.T) {
			results := Parse(input, refTime, nil)

			if len(results) == 0 {
				t.Skipf("Parser doesn't support %q", input)
			}

			result := results[0]
			components, ok := kronos.AsParsingComponents(result.Start())
			if !ok {
				t.Fatalf("Could not convert result to ParsingComponents")
			}

			period := components.Period()
			if period != kronos.PeriodDay {
				t.Errorf("Input %q: period = %v (%s), expected all day expressions to be day-level",
					input, period, period.String())
			}
		})
	}

	// All these should be month-level
	monthExpressions := []string{
		"March",
		"March 2020",
		"11/2005",
	}

	for _, input := range monthExpressions {
		t.Run(input, func(t *testing.T) {
			results := Parse(input, refTime, nil)

			if len(results) == 0 {
				t.Skipf("Parser doesn't support %q", input)
			}

			result := results[0]
			components, ok := kronos.AsParsingComponents(result.Start())
			if !ok {
				t.Fatalf("Could not convert result to ParsingComponents")
			}

			period := components.Period()
			if period != kronos.PeriodMonth {
				t.Errorf("Input %q: period = %v (%s), expected all month expressions to be month-level",
					input, period, period.String())
			}
		})
	}
}

// TestPeriodPreservationInMerging tests that period is preserved during merging operations
func TestPeriodPreservationInMerging(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)

	// Test expressions that might get merged (e.g., "yesterday" + "at 3pm")
	// The period should be updated to the finest granularity after merging
	tests := []struct {
		input          string
		expectedPeriod kronos.Period
	}{
		{"yesterday at 3pm", kronos.PeriodTime}, // Merged: day + time = time
		{"March 15 at 10:30", kronos.PeriodTime}, // Merged: date + time = time
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			results := Parse(tt.input, refTime, nil)

			if len(results) == 0 {
				t.Skipf("No results for %q", tt.input)
			}

			result := results[0]
			components, ok := kronos.AsParsingComponents(result.Start())
			if !ok {
				t.Fatalf("Could not convert result to ParsingComponents")
			}

			period := components.Period()
			if period != tt.expectedPeriod {
				t.Errorf("Merged expression %q: period = %v (%s), want %v (%s)",
					tt.input, period, period.String(), tt.expectedPeriod, tt.expectedPeriod.String())
			}
		})
	}
}
