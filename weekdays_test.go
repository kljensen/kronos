package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNextWeekday(t *testing.T) {
	tests := []struct {
		name           string
		refDate        time.Time
		targetWeekday  Weekday
		expectedDays   int
		expectedDate   time.Time
	}{
		{
			name:          "Monday to next Monday",
			refDate:       time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC), // Monday
			targetWeekday: WeekdayMonday,
			expectedDays:  7,
			expectedDate:  time.Date(2020, 1, 13, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Monday to next Tuesday",
			refDate:       time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC), // Monday
			targetWeekday: WeekdayTuesday,
			expectedDays:  1,
			expectedDate:  time.Date(2020, 1, 7, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Friday to next Monday",
			refDate:       time.Date(2020, 1, 10, 12, 0, 0, 0, time.UTC), // Friday
			targetWeekday: WeekdayMonday,
			expectedDays:  3,
			expectedDate:  time.Date(2020, 1, 13, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Saturday to next Sunday",
			refDate:       time.Date(2020, 1, 11, 12, 0, 0, 0, time.UTC), // Saturday
			targetWeekday: WeekdaySunday,
			expectedDays:  1,
			expectedDate:  time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Sunday to next Saturday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: WeekdaySaturday,
			expectedDays:  6,
			expectedDate:  time.Date(2020, 1, 18, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNextWeekday(tt.refDate, tt.targetWeekday)
			assert.Equal(t, tt.expectedDate, result)
		})
	}
}

func TestGetLastWeekday(t *testing.T) {
	tests := []struct {
		name          string
		refDate       time.Time
		targetWeekday Weekday
		expectedDate  time.Time
	}{
		{
			name:          "Monday to last Monday",
			refDate:       time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC), // Monday
			targetWeekday: WeekdayMonday,
			expectedDate:  time.Date(2019, 12, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Friday to last Monday",
			refDate:       time.Date(2020, 1, 10, 12, 0, 0, 0, time.UTC), // Friday
			targetWeekday: WeekdayMonday,
			expectedDate:  time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Tuesday to last Friday",
			refDate:       time.Date(2020, 1, 7, 12, 0, 0, 0, time.UTC), // Tuesday
			targetWeekday: WeekdayFriday,
			expectedDate:  time.Date(2020, 1, 3, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Sunday to last Sunday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: WeekdaySunday,
			expectedDate:  time.Date(2020, 1, 5, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLastWeekday(tt.refDate, tt.targetWeekday)
			assert.Equal(t, tt.expectedDate, result)
		})
	}
}

func TestGetThisWeekday(t *testing.T) {
	tests := []struct {
		name          string
		refDate       time.Time
		targetWeekday Weekday
		forward       bool
		expectedDate  time.Time
	}{
		{
			name:          "Wednesday to this Monday (forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayMonday,
			forward:       true,
			expectedDate:  time.Date(2020, 1, 13, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Monday (backward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayMonday,
			forward:       false,
			expectedDate:  time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Wednesday (forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayWednesday,
			forward:       true,
			expectedDate:  time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Wednesday (backward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayWednesday,
			forward:       false,
			expectedDate:  time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Friday (forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayFriday,
			forward:       true,
			expectedDate:  time.Date(2020, 1, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name:          "Wednesday to this Friday (backward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayFriday,
			forward:       false,
			expectedDate:  time.Date(2020, 1, 3, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetThisWeekday(tt.refDate, tt.targetWeekday, tt.forward)
			assert.Equal(t, tt.expectedDate, result)
		})
	}
}

func TestGetDaysToWeekday(t *testing.T) {
	tests := []struct {
		name          string
		refDate       time.Time
		targetWeekday Weekday
		modifier      *string
		expectedDays  int
	}{
		{
			name:          "this Monday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayMonday,
			modifier:      stringPtr("this"),
			expectedDays:  5, // Forward to next Monday
		},
		{
			name:          "last Monday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayMonday,
			modifier:      stringPtr("last"),
			expectedDays:  -2, // Back to previous Monday
		},
		{
			name:          "next Monday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayMonday,
			modifier:      stringPtr("next"),
			expectedDays:  5, // Next Monday (since Monday < Wednesday)
		},
		{
			name:          "next Friday from Wednesday",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayFriday,
			modifier:      stringPtr("next"),
			expectedDays:  9, // Next Friday (Friday > Wednesday, so +7)
		},
		{
			name:          "next Sunday from Sunday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: WeekdaySunday,
			modifier:      stringPtr("next"),
			expectedDays:  7,
		},
		{
			name:          "next Monday from Sunday",
			refDate:       time.Date(2020, 1, 12, 12, 0, 0, 0, time.UTC), // Sunday
			targetWeekday: WeekdayMonday,
			modifier:      stringPtr("next"),
			expectedDays:  1,
		},
		{
			name:          "next Saturday from Saturday",
			refDate:       time.Date(2020, 1, 11, 12, 0, 0, 0, time.UTC), // Saturday
			targetWeekday: WeekdaySaturday,
			modifier:      stringPtr("next"),
			expectedDays:  7,
		},
		{
			name:          "next Sunday from Saturday",
			refDate:       time.Date(2020, 1, 11, 12, 0, 0, 0, time.UTC), // Saturday
			targetWeekday: WeekdaySunday,
			modifier:      stringPtr("next"),
			expectedDays:  8,
		},
		{
			name:          "closest Monday from Wednesday (should go back)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayMonday,
			modifier:      nil,
			expectedDays:  -2, // Backward 2 days is closer than forward 5
		},
		{
			name:          "closest Friday from Wednesday (should go forward)",
			refDate:       time.Date(2020, 1, 8, 12, 0, 0, 0, time.UTC), // Wednesday
			targetWeekday: WeekdayFriday,
			modifier:      nil,
			expectedDays:  2, // Forward 2 days is closer than backward 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDaysToWeekday(tt.refDate, tt.targetWeekday, tt.modifier)
			assert.Equal(t, tt.expectedDays, result)
		})
	}
}

func TestGetDaysForwardToWeekday(t *testing.T) {
	// Monday Jan 6, 2020
	refDate := time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC)

	assert.Equal(t, 0, getDaysForwardToWeekday(refDate, WeekdayMonday))
	assert.Equal(t, 1, getDaysForwardToWeekday(refDate, WeekdayTuesday))
	assert.Equal(t, 2, getDaysForwardToWeekday(refDate, WeekdayWednesday))
	assert.Equal(t, 3, getDaysForwardToWeekday(refDate, WeekdayThursday))
	assert.Equal(t, 4, getDaysForwardToWeekday(refDate, WeekdayFriday))
	assert.Equal(t, 5, getDaysForwardToWeekday(refDate, WeekdaySaturday))
	assert.Equal(t, 6, getDaysForwardToWeekday(refDate, WeekdaySunday))
}

func TestGetBackwardDaysToWeekday(t *testing.T) {
	// Monday Jan 6, 2020
	refDate := time.Date(2020, 1, 6, 12, 0, 0, 0, time.UTC)

	assert.Equal(t, -7, getBackwardDaysToWeekday(refDate, WeekdayMonday))
	assert.Equal(t, -6, getBackwardDaysToWeekday(refDate, WeekdayTuesday))
	assert.Equal(t, -5, getBackwardDaysToWeekday(refDate, WeekdayWednesday))
	assert.Equal(t, -4, getBackwardDaysToWeekday(refDate, WeekdayThursday))
	assert.Equal(t, -3, getBackwardDaysToWeekday(refDate, WeekdayFriday))
	assert.Equal(t, -2, getBackwardDaysToWeekday(refDate, WeekdaySaturday))
	assert.Equal(t, -1, getBackwardDaysToWeekday(refDate, WeekdaySunday))
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
