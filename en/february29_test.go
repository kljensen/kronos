package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestFebruary29SmartLeapYearSelection(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		refDate      time.Time
		preference   kronos.DatePreference
		expectedYear int
	}{
		{
			name:         "February 29 from 2023 (non-leap), prefer past",
			input:        "February 29",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferPast,
			expectedYear: 2020,
		},
		{
			name:         "February 29 from 2023 (non-leap), prefer future",
			input:        "February 29",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferFuture,
			expectedYear: 2024,
		},
		{
			name:         "February 29 from 2024 (leap year), prefer current",
			input:        "February 29",
			refDate:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferCurrentPeriod,
			expectedYear: 2024,
		},
		{
			name:         "Feb 29 from 2023, prefer current",
			input:        "Feb 29",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferCurrentPeriod,
			expectedYear: 2024,
		},
		{
			name:         "29th of February from 2023, prefer past",
			input:        "29th of February",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferPast,
			expectedYear: 2020,
		},
		{
			name:         "February 29th from 2023, prefer future",
			input:        "February 29th",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferFuture,
			expectedYear: 2024,
		},
		// Edge cases with century boundaries
		{
			name:         "Feb 29 from 2100 (non-leap century), prefer past",
			input:        "February 29",
			refDate:      time.Date(2100, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferPast,
			expectedYear: 2096,
		},
		{
			name:         "Feb 29 from 2100 (non-leap century), prefer future",
			input:        "February 29",
			refDate:      time.Date(2100, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferFuture,
			expectedYear: 2104,
		},
		{
			name:         "Feb 29 from 2000 (leap century), prefer current",
			input:        "February 29",
			refDate:      time.Date(2000, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferCurrentPeriod,
			expectedYear: 2000,
		},
		// Slash format
		{
			name:         "2/29 from 2023, prefer future",
			input:        "2/29",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferFuture,
			expectedYear: 2024,
		},
		{
			name:         "02/29 from 2023, prefer past",
			input:        "02/29",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferPast,
			expectedYear: 2020,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := New().WithReferenceDate(tt.refDate)
			switch tt.preference {
			case kronos.PreferPast:
				builder = builder.PreferPast()
			case kronos.PreferFuture:
				builder = builder.PreferFuture()
			case kronos.PreferCurrentPeriod:
				builder = builder.PreferCurrentPeriod()
			}

			results, err := builder.Parse(tt.input)
			assert.NoError(t, err)

			assert.NotEmpty(t, results, "Should parse %s", tt.input)
			if len(results) > 0 {
				result := results[0]
				date := result.Start().Date()

				assert.Equal(t, tt.expectedYear, date.Year(), "Year should be %d for input %q", tt.expectedYear, tt.input)
				assert.Equal(t, 2, int(date.Month()), "Month should be February")
				assert.Equal(t, 29, date.Day(), "Day should be 29")
			}
		})
	}
}

func TestFebruary29WithTime(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		refDate      time.Time
		preference   kronos.DatePreference
		expectedYear int
		expectedHour int
	}{
		{
			name:         "February 29 at 3:30 PM from 2023, prefer future",
			input:        "February 29 at 3:30 PM",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferFuture,
			expectedYear: 2024,
			expectedHour: 15,
		},
		{
			name:         "Feb 29 10:00 from 2023, prefer past",
			input:        "Feb 29 10:00",
			refDate:      time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference:   kronos.PreferPast,
			expectedYear: 2020,
			expectedHour: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := New().WithReferenceDate(tt.refDate)
			switch tt.preference {
			case kronos.PreferPast:
				builder = builder.PreferPast()
			case kronos.PreferFuture:
				builder = builder.PreferFuture()
			case kronos.PreferCurrentPeriod:
				builder = builder.PreferCurrentPeriod()
			}

			results, err := builder.Parse(tt.input)
			assert.NoError(t, err)

			assert.NotEmpty(t, results, "Should parse %s", tt.input)
			if len(results) > 0 {
				result := results[0]
				date := result.Start().Date()

				assert.Equal(t, tt.expectedYear, date.Year(), "Year should be %d", tt.expectedYear)
				assert.Equal(t, 2, int(date.Month()), "Month should be February")
				assert.Equal(t, 29, date.Day(), "Day should be 29")
				assert.Equal(t, tt.expectedHour, date.Hour(), "Hour should be %d", tt.expectedHour)
			}
		})
	}
}
