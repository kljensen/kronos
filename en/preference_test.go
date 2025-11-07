//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/experimental"
)

// TestDatePreferenceWithMonth tests month-only expressions with different preferences
func TestDatePreferenceWithMonth(t *testing.T) {
	// Reference: February 15, 2015, 15:30
	refDate := time.Date(2015, 2, 15, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		preference kronos.DatePreference
		expected   time.Time
	}{
		// March tests
		{
			name:       "March with PreferPast should be March 2014",
			input:      "March",
			preference: kronos.PreferPast,
			expected:   time.Date(2014, 3, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "March with PreferFuture should be March 2015",
			input:      "March",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 3, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "March with PreferCurrentPeriod should be March 2015",
			input:      "March",
			preference: kronos.PreferCurrentPeriod,
			expected:   time.Date(2015, 3, 1, 12, 0, 0, 0, time.UTC),
		},
		// August tests
		{
			name:       "August with PreferPast should be August 2014",
			input:      "August",
			preference: kronos.PreferPast,
			expected:   time.Date(2014, 8, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "August with PreferFuture should be August 2015",
			input:      "August",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 8, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "August with PreferCurrentPeriod should be August 2015",
			input:      "August",
			preference: kronos.PreferCurrentPeriod,
			expected:   time.Date(2015, 8, 1, 12, 0, 0, 0, time.UTC),
		},
		// January tests (before reference month)
		{
			name:       "January with PreferPast should be January 2015",
			input:      "January",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "January with PreferFuture should be January 2016",
			input:      "January",
			preference: kronos.PreferFuture,
			expected:   time.Date(2016, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := &kronos.InternalParsingOption{
				Preference: tt.preference,
			}
			chrono := experimental.EnglishCasualChrono()
			results := chrono.Parse(tt.input, refDate, option)
			if len(results) == 0 {
				t.Fatal("Expected at least one result")
			}

			result := results[0].Date()
			if result.Year() != tt.expected.Year() || result.Month() != tt.expected.Month() {
				t.Errorf("Expected %v, got %v", tt.expected.Format("2006-01"), result.Format("2006-01"))
			}
		})
	}
}

// TestDatePreferenceWithMonthAndDay tests date expressions with missing year
func TestDatePreferenceWithMonthAndDay(t *testing.T) {
	// Reference: February 15, 2015, 15:30
	refDate := time.Date(2015, 2, 15, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		preference kronos.DatePreference
		expected   time.Time
	}{
		// March 15 tests (after reference date)
		{
			name:       "March 15 with PreferPast should be March 15, 2014",
			input:      "March 15",
			preference: kronos.PreferPast,
			expected:   time.Date(2014, 3, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "March 15 with PreferFuture should be March 15, 2015",
			input:      "March 15",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 3, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "March 15 with PreferCurrentPeriod should be March 15, 2015",
			input:      "March 15",
			preference: kronos.PreferCurrentPeriod,
			expected:   time.Date(2015, 3, 15, 12, 0, 0, 0, time.UTC),
		},
		// January 20 tests (before reference date but same year)
		{
			name:       "January 20 with PreferPast should be January 20, 2015",
			input:      "January 20",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 1, 20, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "January 20 with PreferFuture should be January 20, 2016",
			input:      "January 20",
			preference: kronos.PreferFuture,
			expected:   time.Date(2016, 1, 20, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "January 20 with PreferCurrentPeriod should be January 20, 2015",
			input:      "January 20",
			preference: kronos.PreferCurrentPeriod,
			expected:   time.Date(2015, 1, 20, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := &kronos.InternalParsingOption{
				Preference: tt.preference,
			}
			chrono := experimental.EnglishCasualChrono()
			results := chrono.Parse(tt.input, refDate, option)
			if len(results) == 0 {
				t.Fatal("Expected at least one result")
			}

			result := results[0].Date()
			if result.Year() != tt.expected.Year() || result.Month() != tt.expected.Month() || result.Day() != tt.expected.Day() {
				t.Errorf("Expected %v, got %v", tt.expected.Format("2006-01-02"), result.Format("2006-01-02"))
			}
		})
	}
}

// TestDatePreferenceWithTimeOnly tests time-only expressions
func TestDatePreferenceWithTimeOnly(t *testing.T) {
	// Reference: February 15, 2015, 15:30
	refDate := time.Date(2015, 2, 15, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		preference kronos.DatePreference
		expected   time.Time
	}{
		// 10:00 (before reference time)
		{
			name:       "10:00 with PreferPast should be Feb 15 10:00",
			input:      "10:00",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 2, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:       "10:00 with PreferFuture should be Feb 16 10:00",
			input:      "10:00",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 2, 16, 10, 0, 0, 0, time.UTC),
		},
		{
			name:       "10:00 with PreferCurrentPeriod should be Feb 15 10:00",
			input:      "10:00",
			preference: kronos.PreferCurrentPeriod,
			expected:   time.Date(2015, 2, 15, 10, 0, 0, 0, time.UTC),
		},
		// 18:00 (after reference time)
		{
			name:       "18:00 with PreferPast should be Feb 14 18:00",
			input:      "18:00",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 2, 14, 18, 0, 0, 0, time.UTC),
		},
		{
			name:       "18:00 with PreferFuture should be Feb 15 18:00",
			input:      "18:00",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 2, 15, 18, 0, 0, 0, time.UTC),
		},
		{
			name:       "18:00 with PreferCurrentPeriod should be Feb 15 18:00",
			input:      "18:00",
			preference: kronos.PreferCurrentPeriod,
			expected:   time.Date(2015, 2, 15, 18, 0, 0, 0, time.UTC),
		},
		// 3pm (before reference time of 15:30, so already in the past)
		{
			name:       "3pm with PreferPast should be Feb 15 15:00",
			input:      "3pm",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 2, 15, 15, 0, 0, 0, time.UTC),
		},
		{
			name:       "3pm with PreferFuture should be Feb 16 15:00",
			input:      "3pm",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 2, 16, 15, 0, 0, 0, time.UTC),
		},
		{
			name:       "3pm with PreferCurrentPeriod should be Feb 15 15:00",
			input:      "3pm",
			preference: kronos.PreferCurrentPeriod,
			expected:   time.Date(2015, 2, 15, 15, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := &kronos.InternalParsingOption{
				Preference: tt.preference,
			}
			chrono := experimental.EnglishCasualChrono()
			results := chrono.Parse(tt.input, refDate, option)
			if len(results) == 0 {
				t.Fatal("Expected at least one result")
			}

			result := results[0].Date()
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected.Format("2006-01-02 15:04"), result.Format("2006-01-02 15:04"))
			}
		})
	}
}

// TestDatePreferenceWithAbsoluteDates ensures that absolute dates ignore preference
func TestDatePreferenceWithAbsoluteDates(t *testing.T) {
	// Reference: February 15, 2015, 15:30
	refDate := time.Date(2015, 2, 15, 15, 30, 0, 0, time.UTC)
	expected := time.Date(2020, 5, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		preference kronos.DatePreference
	}{
		{"Absolute date with PreferPast", "May 10, 2020", kronos.PreferPast},
		{"Absolute date with PreferFuture", "May 10, 2020", kronos.PreferFuture},
		{"Absolute date with PreferCurrentPeriod", "May 10, 2020", kronos.PreferCurrentPeriod},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := &kronos.InternalParsingOption{
				Preference: tt.preference,
			}
			chrono := experimental.EnglishCasualChrono()
			results := chrono.Parse(tt.input, refDate, option)
			if len(results) == 0 {
				t.Fatal("Expected at least one result")
			}

			result := results[0].Date()
			if result.Year() != expected.Year() || result.Month() != expected.Month() || result.Day() != expected.Day() {
				t.Errorf("Expected %v, got %v (preference should not affect absolute dates)",
					expected.Format("2006-01-02"), result.Format("2006-01-02"))
			}
		})
	}
}

// TestDatePreferenceWithLittleEndianFormat tests little-endian date format
func TestDatePreferenceWithLittleEndianFormat(t *testing.T) {
	// Reference: February 15, 2015, 15:30
	refDate := time.Date(2015, 2, 15, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		preference kronos.DatePreference
		expected   time.Time
	}{
		{
			name:       "10 August with PreferPast should be August 10, 2014",
			input:      "10 August",
			preference: kronos.PreferPast,
			expected:   time.Date(2014, 8, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "10 August with PreferFuture should be August 10, 2015",
			input:      "10 August",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 8, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "31 January with PreferPast should be January 31, 2015",
			input:      "31 January",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 1, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "31 January with PreferFuture should be January 31, 2016",
			input:      "31 January",
			preference: kronos.PreferFuture,
			expected:   time.Date(2016, 1, 31, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := &kronos.InternalParsingOption{
				Preference: tt.preference,
			}
			chrono := experimental.EnglishGBChrono()
			results := chrono.Parse(tt.input, refDate, option)
			if len(results) == 0 {
				t.Fatal("Expected at least one result")
			}

			result := results[0].Date()
			if result.Year() != tt.expected.Year() || result.Month() != tt.expected.Month() || result.Day() != tt.expected.Day() {
				t.Errorf("Expected %v, got %v", tt.expected.Format("2006-01-02"), result.Format("2006-01-02"))
			}
		})
	}
}

// TestDatePreferenceDoesNotAffectRelativeDates ensures relative dates are not affected by preference
func TestDatePreferenceDoesNotAffectRelativeDates(t *testing.T) {
	// Reference: February 15, 2015, 15:30
	refDate := time.Date(2015, 2, 15, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		preference kronos.DatePreference
		expected   time.Time
	}{
		{
			name:       "tomorrow with PreferPast should be Feb 16",
			input:      "tomorrow",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 2, 16, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "tomorrow with PreferFuture should be Feb 16",
			input:      "tomorrow",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 2, 16, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "yesterday with PreferPast should be Feb 14",
			input:      "yesterday",
			preference: kronos.PreferPast,
			expected:   time.Date(2015, 2, 14, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "yesterday with PreferFuture should be Feb 14",
			input:      "yesterday",
			preference: kronos.PreferFuture,
			expected:   time.Date(2015, 2, 14, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := &kronos.InternalParsingOption{
				Preference: tt.preference,
			}
			chrono := experimental.EnglishCasualChrono()
			results := chrono.Parse(tt.input, refDate, option)
			if len(results) == 0 {
				t.Fatal("Expected at least one result")
			}

			result := results[0].Date()
			if result.Year() != tt.expected.Year() || result.Month() != tt.expected.Month() || result.Day() != tt.expected.Day() {
				t.Errorf("Expected %v, got %v (preference should not affect relative dates)",
					tt.expected.Format("2006-01-02"), result.Format("2006-01-02"))
			}
		})
	}
}

// TestPreferenceDefault tests that default preference is PreferCurrentPeriod
func TestPreferenceDefault(t *testing.T) {
	// Reference: February 15, 2015, 15:30
	refDate := time.Date(2015, 2, 15, 15, 30, 0, 0, time.UTC)

	// Parse without specifying preference
	chrono := experimental.EnglishCasualChrono()
	results := chrono.Parse("March 15", refDate, nil)

	if len(results) == 0 {
		t.Fatal("Expected at least one result")
	}

	result := results[0].Date()
	expected := time.Date(2015, 3, 15, 12, 0, 0, 0, time.UTC)

	if result.Year() != expected.Year() || result.Month() != expected.Month() || result.Day() != expected.Day() {
		t.Errorf("Expected default to be current period: %v, got %v",
			expected.Format("2006-01-02"), result.Format("2006-01-02"))
	}
}
