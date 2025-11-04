package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// TestDateOrderDMY_SlashSeparatedDates tests that DateOrder(DMY) works correctly
// for slash-separated dates. This is a regression test for Issue #96.
func TestDateOrderDMY_SlashSeparatedDates(t *testing.T) {
	refDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		input         string
		expectedYear  int
		expectedMonth time.Month
		expectedDay   int
	}{
		{
			name:          "3/5/2024 with DMY should be May 3",
			input:         "3/5/2024",
			expectedYear:  2024,
			expectedMonth: time.May,
			expectedDay:   3,
		},
		{
			name:          "15/3/2024 with DMY should be March 15",
			input:         "15/3/2024",
			expectedYear:  2024,
			expectedMonth: time.March,
			expectedDay:   15,
		},
		{
			name:          "1/12/2024 with DMY should be December 1",
			input:         "1/12/2024",
			expectedYear:  2024,
			expectedMonth: time.December,
			expectedDay:   1,
		},
		{
			name:          "25/12/2024 with DMY should be December 25",
			input:         "25/12/2024",
			expectedYear:  2024,
			expectedMonth: time.December,
			expectedDay:   25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use DMY date order
			date, err := New().
				WithReferenceDate(refDate).
				DateOrder(kronos.DateOrderDMY).
				ParseDate(tt.input)

			assert.NoError(t, err, "Should not error")
			assert.NotNil(t, date, "Should parse date")

			if date != nil {
				assert.Equal(t, tt.expectedYear, date.Year(), "Year mismatch")
				assert.Equal(t, tt.expectedMonth, date.Month(), "Month mismatch")
				assert.Equal(t, tt.expectedDay, date.Day(), "Day mismatch")
			}
		})
	}
}

// TestDateOrderMDY_SlashSeparatedDates tests that DateOrder(MDY) works correctly
// for slash-separated dates (US format).
func TestDateOrderMDY_SlashSeparatedDates(t *testing.T) {
	refDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		input         string
		expectedYear  int
		expectedMonth time.Month
		expectedDay   int
	}{
		{
			name:          "3/5/2024 with MDY should be March 5",
			input:         "3/5/2024",
			expectedYear:  2024,
			expectedMonth: time.March,
			expectedDay:   5,
		},
		{
			name:          "12/25/2024 with MDY should be December 25",
			input:         "12/25/2024",
			expectedYear:  2024,
			expectedMonth: time.December,
			expectedDay:   25,
		},
		{
			name:          "1/15/2024 with MDY should be January 15",
			input:         "1/15/2024",
			expectedYear:  2024,
			expectedMonth: time.January,
			expectedDay:   15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use MDY date order (US format)
			date, err := New().
				WithReferenceDate(refDate).
				DateOrder(kronos.DateOrderMDY).
				ParseDate(tt.input)

			assert.NoError(t, err, "Should not error")
			assert.NotNil(t, date, "Should parse date")

			if date != nil {
				assert.Equal(t, tt.expectedYear, date.Year(), "Year mismatch")
				assert.Equal(t, tt.expectedMonth, date.Month(), "Month mismatch")
				assert.Equal(t, tt.expectedDay, date.Day(), "Day mismatch")
			}
		})
	}
}
