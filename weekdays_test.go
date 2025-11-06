package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNextWeekday(t *testing.T) {
	tests := []struct {
		name          string
		refDate       time.Time
		targetWeekday time.Weekday
		expectedDays  int
		expectedDate  time.Time
	}{
		{
			name:          "Monday to next Monday",
			refDate:       time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC), // Monday
			targetWeekday: time.Monday,
			expectedDays:  7,
			expectedDate:  time.Date(2020, 1, 13, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Monday to next Tuesday",
			refDate:       time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC), // Monday
			targetWeekday: time.Tuesday,
			expectedDays:  1,
			expectedDate:  time.Date(2020, 1, 7, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Friday to next Monday",
			refDate:       time.Date(2020, 1, 10, 12, 0, 0, 0, time.UTC), // Friday
			targetWeekday: time.Monday,
			expectedDays:  3,
			expectedDate:  time.Date(2020, 1, 13, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Saturday to next Sunday",
			refDate:       time.Date(2020, 1, 11, 12, 0, 0, 0, time.UTC), // Saturday
			targetWeekday: time.Sunday,
			expectedDays:  1,
			expectedDate:  time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Sunday to next Saturday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: time.Saturday,
			expectedDays:  6,
			expectedDate:  time.Date(2020, 1, 18, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNextWeekday(tt.refDate, tt.targetWeekday)
			assert.Equal(t, tt.expectedDate, result)
		})
	}
}

func TestGetLastWeekday(t *testing.T) {
	tests := []struct {
		name          string
		refDate       time.Time
		targetWeekday time.Weekday
		expectedDate  time.Time
	}{
		{
			name:          "Monday to last Monday",
			refDate:       time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC), // Monday
			targetWeekday: time.Monday,
			expectedDate:  time.Date(2019, 12, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Friday to last Monday",
			refDate:       time.Date(2020, 1, 10, 12, 0, 0, 0, time.UTC), // Friday
			targetWeekday: time.Monday,
			expectedDate:  time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Tuesday to last Friday",
			refDate:       time.Date(2020, 1, 7, 12, 0, 0, 0, time.UTC), // Tuesday
			targetWeekday: time.Friday,
			expectedDate:  time.Date(2020, 1, 3, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Sunday to last Sunday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: time.Sunday,
			expectedDate:  time.Date(2020, 1, 5, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLastWeekday(tt.refDate, tt.targetWeekday)
			assert.Equal(t, tt.expectedDate, result)
		})
	}
}

func TestGetThisWeekday(t *testing.T) {
	tests := []struct {
		name          string
		refDate       time.Time
		targetWeekday time.Weekday
		forward       bool
		expectedDate  time.Time
	}{
		{
			name:          "Wednesday to this Monday (forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Monday,
			forward:       true,
			expectedDate:  time.Date(2020, 1, 13, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Monday (backward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Monday,
			forward:       false,
			expectedDate:  time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Wednesday (forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Wednesday,
			forward:       true,
			expectedDate:  time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Wednesday (backward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Wednesday,
			forward:       false,
			expectedDate:  time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Friday (forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Friday,
			forward:       true,
			expectedDate:  time.Date(2020, 1, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Friday (backward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Friday,
			forward:       false,
			expectedDate:  time.Date(2020, 1, 3, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThisWeekday(tt.refDate, tt.targetWeekday, tt.forward)
			assert.Equal(t, tt.expectedDate, result)
		})
	}
}

func TestGetDaysToWeekday(t *testing.T) {
	tests := []struct {
		name          string
		refDate       time.Time
		targetWeekday time.Weekday
		modifier      *string
		expectedDays  int
	}{
		{
			name:          "this Monday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Monday,
			modifier:      stringPtr("this"),
			expectedDays:  5, // Forward to next Monday
		},
		{
			name:          "last Monday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Monday,
			modifier:      stringPtr("last"),
			expectedDays:  -2, // Back to previous Monday
		},
		{
			name:          "next Monday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Monday,
			modifier:      stringPtr("next"),
			expectedDays:  5, // Next Monday (since Monday < Wednesday)
		},
		{
			name:          "next Friday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Friday,
			modifier:      stringPtr("next"),
			expectedDays:  9, // Next Friday (Friday > Wednesday, so +7)
		},
		{
			name:          "next Sunday from Sunday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: time.Sunday,
			modifier:      stringPtr("next"),
			expectedDays:  7,
		},
		{
			name:          "next Monday from Sunday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: time.Monday,
			modifier:      stringPtr("next"),
			expectedDays:  1,
		},
		{
			name:          "next Saturday from Saturday",
			refDate:       time.Date(2020, 1, 11, 12, 0, 0, 0, time.UTC), // Saturday
			targetWeekday: time.Saturday,
			modifier:      stringPtr("next"),
			expectedDays:  7,
		},
		{
			name:          "next Sunday from Saturday",
			refDate:       time.Date(2020, 1, 11, 12, 0, 0, 0, time.UTC), // Saturday
			targetWeekday: time.Sunday,
			modifier:      stringPtr("next"),
			expectedDays:  8,
		},
		{
			name:          "closest Monday from Wednesday (should go back)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Monday,
			modifier:      nil,
			expectedDays:  -2, // Backward 2 days is closer than forward 5
		},
		{
			name:          "closest Friday from Wednesday (should go forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: time.Friday,
			modifier:      nil,
			expectedDays:  2, // Forward 2 days is closer than backward 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDaysToWeekday(tt.refDate, tt.targetWeekday, tt.modifier)
			assert.Equal(t, tt.expectedDays, result)
		})
	}
}

func TestGetDaysForwardToWeekday(t *testing.T) {
	// Monday Jan 6, 2020
	refDate := time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC)

	assert.Equal(t, 0, getDaysForwardToWeekday(refDate, time.Monday))
	assert.Equal(t, 1, getDaysForwardToWeekday(refDate, time.Tuesday))
	assert.Equal(t, 2, getDaysForwardToWeekday(refDate, time.Wednesday))
	assert.Equal(t, 3, getDaysForwardToWeekday(refDate, time.Thursday))
	assert.Equal(t, 4, getDaysForwardToWeekday(refDate, time.Friday))
	assert.Equal(t, 5, getDaysForwardToWeekday(refDate, time.Saturday))
	assert.Equal(t, 6, getDaysForwardToWeekday(refDate, time.Sunday))
}

func TestGetBackwardDaysToWeekday(t *testing.T) {
	// Monday Jan 6, 2020
	refDate := time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC)

	assert.Equal(t, -7, getBackwardDaysToWeekday(refDate, time.Monday))
	assert.Equal(t, -6, getBackwardDaysToWeekday(refDate, time.Tuesday))
	assert.Equal(t, -5, getBackwardDaysToWeekday(refDate, time.Wednesday))
	assert.Equal(t, -4, getBackwardDaysToWeekday(refDate, time.Thursday))
	assert.Equal(t, -3, getBackwardDaysToWeekday(refDate, time.Friday))
	assert.Equal(t, -2, getBackwardDaysToWeekday(refDate, time.Saturday))
	assert.Equal(t, -1, getBackwardDaysToWeekday(refDate, time.Sunday))
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
