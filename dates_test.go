package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssignSimilarDate(t *testing.T) {
	reference := NewReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := NewParsingComponents(reference, nil)

	targetDate := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	AssignSimilarDate(components, targetDate)

	assert.True(t, components.IsCertain(ComponentYear))
	assert.True(t, components.IsCertain(ComponentMonth))
	assert.True(t, components.IsCertain(ComponentDay))

	assert.Equal(t, 2021, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
}

func TestAssignSimilarTime(t *testing.T) {
	reference := NewReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := NewParsingComponents(reference, nil)

	// Test AM time
	targetDate := time.Date(2021, 6, 15, 8, 30, 45, 123000000, time.UTC)
	AssignSimilarTime(components, targetDate)

	assert.True(t, components.IsCertain(ComponentHour))
	assert.True(t, components.IsCertain(ComponentMinute))
	assert.True(t, components.IsCertain(ComponentSecond))
	assert.True(t, components.IsCertain(ComponentMillisecond))
	assert.True(t, components.IsCertain(ComponentMeridiem))

	assert.Equal(t, 8, *components.Get(ComponentHour))
	assert.Equal(t, 30, *components.Get(ComponentMinute))
	assert.Equal(t, 45, *components.Get(ComponentSecond))
	assert.Equal(t, 123, *components.Get(ComponentMillisecond))
	assert.Equal(t, int(MeridiemAM), *components.Get(ComponentMeridiem))
}

func TestAssignSimilarTime_PM(t *testing.T) {
	reference := NewReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := NewParsingComponents(reference, nil)

	// Test PM time
	targetDate := time.Date(2021, 6, 15, 14, 30, 0, 0, time.UTC)
	AssignSimilarTime(components, targetDate)

	assert.Equal(t, 14, *components.Get(ComponentHour))
	assert.Equal(t, int(MeridiemPM), *components.Get(ComponentMeridiem))
}

func TestImplySimilarDate(t *testing.T) {
	reference := NewReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := NewParsingComponents(reference, nil)

	// Set year as certain first
	components.Assign(ComponentYear, 2025)

	targetDate := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	ImplySimilarDate(components, targetDate)

	// Year should remain 2025 (certain value not overridden by imply)
	assert.True(t, components.IsCertain(ComponentYear))
	assert.Equal(t, 2025, *components.Get(ComponentYear))

	// Month and day should be implied
	assert.False(t, components.IsCertain(ComponentMonth))
	assert.False(t, components.IsCertain(ComponentDay))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
}

func TestImplySimilarTime(t *testing.T) {
	reference := NewReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := NewParsingComponents(reference, nil)

	// Set hour as certain first
	components.Assign(ComponentHour, 18)

	targetDate := time.Date(2021, 6, 15, 10, 30, 45, 0, time.UTC)
	ImplySimilarTime(components, targetDate)

	// Hour should remain 18 (certain value not overridden by imply)
	assert.True(t, components.IsCertain(ComponentHour))
	assert.Equal(t, 18, *components.Get(ComponentHour))

	// Minute and second should be implied
	assert.False(t, components.IsCertain(ComponentMinute))
	assert.False(t, components.IsCertain(ComponentSecond))
	assert.Equal(t, 30, *components.Get(ComponentMinute))
	assert.Equal(t, 45, *components.Get(ComponentSecond))
}

func TestFindMostLikelyADYear(t *testing.T) {
	tests := []struct {
		name     string
		rawYear  int
		expected int
	}{
		{
			name:     "4-digit year unchanged",
			rawYear:  2020,
			expected: 2020,
		},
		{
			name:     "3-digit year unchanged",
			rawYear:  999,
			expected: 999,
		},
		{
			name:     "year 0 -> 2000",
			rawYear:  0,
			expected: 2000,
		},
		{
			name:     "year 20 -> 2020s",
			rawYear:  20,
			expected: 2020,
		},
		{
			name:     "year 25 -> 2020s",
			rawYear:  25,
			expected: 2025,
		},
		{
			name:     "year 99 -> 1999",
			rawYear:  99,
			expected: 1999,
		},
		{
			name:     "year 50 -> 1950",
			rawYear:  50,
			expected: 1950,
		},
		{
			name:     "year 80 -> 1980",
			rawYear:  80,
			expected: 1980,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindMostLikelyADYear(tt.rawYear)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindYearClosestToRef(t *testing.T) {
	tests := []struct {
		name     string
		refDate  time.Time
		day      int
		month    int
		expected int
	}{
		{
			name:     "exact match with reference year",
			refDate:  time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC),
			day:      15,
			month:    6,
			expected: 2020,
		},
		{
			name:     "future date in same year",
			refDate:  time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
			day:      20,
			month:    3,
			expected: 2020,
		},
		{
			name:     "past date in same year",
			refDate:  time.Date(2020, 12, 15, 12, 0, 0, 0, time.UTC),
			day:      20,
			month:    3,
			expected: 2021, // March 20, 2021 is closer than March 20, 2020
		},
		{
			name:     "date closer to previous year",
			refDate:  time.Date(2020, 1, 5, 12, 0, 0, 0, time.UTC),
			day:      25,
			month:    12,
			expected: 2019,
		},
		{
			name:     "date closer to next year",
			refDate:  time.Date(2020, 12, 25, 12, 0, 0, 0, time.UTC),
			day:      5,
			month:    1,
			expected: 2021,
		},
		{
			name:     "February 29 on leap year ref",
			refDate:  time.Date(2020, 3, 1, 12, 0, 0, 0, time.UTC),
			day:      29,
			month:    2,
			expected: 2020,
		},
		{
			name:     "mid-year reference",
			refDate:  time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC),
			day:      1,
			month:    1,
			expected: 2020,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindYearClosestToRef(tt.refDate, tt.day, tt.month)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssignOverridesImply(t *testing.T) {
	reference := NewReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := NewParsingComponents(reference, nil)

	targetDate := time.Date(2021, 6, 15, 14, 30, 0, 0, time.UTC)

	// First imply a date
	ImplySimilarDate(components, targetDate)
	assert.False(t, components.IsCertain(ComponentDay))
	assert.Equal(t, 15, *components.Get(ComponentDay))

	// Now assign a different date
	targetDate2 := time.Date(2022, 8, 20, 0, 0, 0, 0, time.UTC)
	AssignSimilarDate(components, targetDate2)

	// Day should now be certain and have the new value
	assert.True(t, components.IsCertain(ComponentDay))
	assert.Equal(t, 20, *components.Get(ComponentDay))
	assert.Equal(t, 8, *components.Get(ComponentMonth))
	assert.Equal(t, 2022, *components.Get(ComponentYear))
}
