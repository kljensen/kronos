package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssignSimilarDate(t *testing.T) {
	reference := newReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := newParsingComponents(reference, nil)

	targetDate := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	assignSimilarDate(components, targetDate)

	assert.True(t, components.IsCertain(ComponentYear))
	assert.True(t, components.IsCertain(ComponentMonth))
	assert.True(t, components.IsCertain(ComponentDay))

	assert.Equal(t, 2021, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
}

func TestAssignSimilarTime(t *testing.T) {
	reference := newReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := newParsingComponents(reference, nil)

	// Test AM time
	targetDate := time.Date(2021, 6, 15, 8, 30, 45, 123000000, time.UTC)
	assignSimilarTime(components, targetDate)

	assert.True(t, components.IsCertain(ComponentHour))
	assert.True(t, components.IsCertain(ComponentMinute))
	assert.True(t, components.IsCertain(ComponentSecond))
	assert.True(t, components.IsCertain(ComponentMillisecond))
	assert.True(t, components.IsCertain(ComponentMeridiem))

	assert.Equal(t, 8, *components.Get(ComponentHour))
	assert.Equal(t, 30, *components.Get(ComponentMinute))
	assert.Equal(t, 45, *components.Get(ComponentSecond))
	assert.Equal(t, 123, *components.Get(ComponentMillisecond))
	assert.Equal(t, int(0), *components.Get(ComponentMeridiem))
}

func TestAssignSimilarTime_PM(t *testing.T) {
	reference := newReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := newParsingComponents(reference, nil)

	// Test PM time
	targetDate := time.Date(2021, 6, 15, 14, 30, 0, 0, time.UTC)
	assignSimilarTime(components, targetDate)

	assert.Equal(t, 14, *components.Get(ComponentHour))
	assert.Equal(t, int(1), *components.Get(ComponentMeridiem))
}

func TestImplySimilarDate(t *testing.T) {
	reference := newReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := newParsingComponents(reference, nil)

	// Set year as certain first
	components.Assign(ComponentYear, 2025)

	targetDate := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	implySimilarDate(components, targetDate)

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
	reference := newReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := newParsingComponents(reference, nil)

	// Set hour as certain first
	components.Assign(ComponentHour, 18)

	targetDate := time.Date(2021, 6, 15, 10, 30, 45, 0, time.UTC)
	implySimilarTime(components, targetDate)

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
			result := findMostLikelyADYear(tt.rawYear)
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
			result := findYearClosestToRef(tt.refDate, tt.day, tt.month)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssignOverridesImply(t *testing.T) {
	reference := newReferenceWithTimezone(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	components := newParsingComponents(reference, nil)

	targetDate := time.Date(2021, 6, 15, 14, 30, 0, 0, time.UTC)

	// First imply a date
	implySimilarDate(components, targetDate)
	assert.False(t, components.IsCertain(ComponentDay))
	assert.Equal(t, 15, *components.Get(ComponentDay))

	// Now assign a different date
	targetDate2 := time.Date(2022, 8, 20, 0, 0, 0, 0, time.UTC)
	assignSimilarDate(components, targetDate2)

	// Day should now be certain and have the new value
	assert.True(t, components.IsCertain(ComponentDay))
	assert.Equal(t, 20, *components.Get(ComponentDay))
	assert.Equal(t, 8, *components.Get(ComponentMonth))
	assert.Equal(t, 2022, *components.Get(ComponentYear))
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		name   string
		year   int
		isLeap bool
	}{
		// Divisible by 400 (leap years)
		{name: "2000 divisible by 400", year: 2000, isLeap: true},
		{name: "2400 divisible by 400", year: 2400, isLeap: true},

		// Divisible by 100 but not 400 (NOT leap years)
		{name: "1900 divisible by 100", year: 1900, isLeap: false},
		{name: "2100 divisible by 100", year: 2100, isLeap: false},

		// Divisible by 4 but not 100 (leap years)
		{name: "2004 divisible by 4", year: 2004, isLeap: true},
		{name: "2020 divisible by 4", year: 2020, isLeap: true},
		{name: "2024 divisible by 4", year: 2024, isLeap: true},
		{name: "1896 divisible by 4", year: 1896, isLeap: true},

		// Not divisible by 4 (NOT leap years)
		{name: "2001 not divisible by 4", year: 2001, isLeap: false},
		{name: "2019 not divisible by 4", year: 2019, isLeap: false},
		{name: "2021 not divisible by 4", year: 2021, isLeap: false},
		{name: "2023 not divisible by 4", year: 2023, isLeap: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLeapYear(tt.year)
			assert.Equal(t, tt.isLeap, result, "Year %d leap year detection", tt.year)
		})
	}
}

func TestFindPreviousLeapYear(t *testing.T) {
	tests := []struct {
		name     string
		baseYear int
		expected int
	}{
		{name: "from 2023 non-leap", baseYear: 2023, expected: 2020},
		{name: "from 2024 leap year", baseYear: 2024, expected: 2024},
		{name: "from 2021 non-leap", baseYear: 2021, expected: 2020},
		{name: "from 2020 leap year", baseYear: 2020, expected: 2020},
		{name: "from 2001 non-leap", baseYear: 2001, expected: 2000},
		{name: "from 1901 near 1900", baseYear: 1901, expected: 1901}, // 1900 is NOT leap, but bounded at 1900 returns baseYear
		{name: "from 1900 non-leap", baseYear: 1900, expected: 1900},  // Bounded at 1900, returns baseYear as fallback
		{name: "crossing century", baseYear: 2102, expected: 2096},    // 2100 is NOT a leap year, so skip to 2096
		{name: "crossing century correctly", baseYear: 2099, expected: 2096},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findPreviousLeapYear(tt.baseYear)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindNextLeapYear(t *testing.T) {
	tests := []struct {
		name     string
		baseYear int
		expected int
	}{
		{name: "from 2023 non-leap", baseYear: 2023, expected: 2024},
		{name: "from 2024 leap year", baseYear: 2024, expected: 2024},
		{name: "from 2021 non-leap", baseYear: 2021, expected: 2024},
		{name: "from 2020 leap year", baseYear: 2020, expected: 2020},
		{name: "from 1999 non-leap", baseYear: 1999, expected: 2000},
		{name: "from 2098 non-leap", baseYear: 2098, expected: 2104}, // Skip 2100 (not leap)
		{name: "crossing century", baseYear: 2100, expected: 2104},
		{name: "from 2397 non-leap", baseYear: 2397, expected: 2400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findNextLeapYear(tt.baseYear)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindNearestLeapYear(t *testing.T) {
	tests := []struct {
		name       string
		baseYear   int
		preference DatePreference
		expected   int
	}{
		// PreferPast tests
		{name: "prefer past from 2023", baseYear: 2023, preference: PreferPast, expected: 2020},
		{name: "prefer past from 2024 (leap)", baseYear: 2024, preference: PreferPast, expected: 2024},
		{name: "prefer past from 2021", baseYear: 2021, preference: PreferPast, expected: 2020},

		// PreferFuture tests
		{name: "prefer future from 2023", baseYear: 2023, preference: PreferFuture, expected: 2024},
		{name: "prefer future from 2024 (leap)", baseYear: 2024, preference: PreferFuture, expected: 2024},
		{name: "prefer future from 2021", baseYear: 2021, preference: PreferFuture, expected: 2024},

		// PreferCurrentPeriod tests (prefer nearest, break ties by future)
		{name: "prefer current from 2023", baseYear: 2023, preference: PreferCurrentPeriod, expected: 2024},
		{name: "prefer current from 2024 (leap)", baseYear: 2024, preference: PreferCurrentPeriod, expected: 2024},
		{name: "prefer current from 2022", baseYear: 2022, preference: PreferCurrentPeriod, expected: 2024}, // 2020 is 2 away, 2024 is 2 away, prefer future tie-break gives 2024
		{name: "prefer current from 2021", baseYear: 2021, preference: PreferCurrentPeriod, expected: 2020}, // 2020 is 1 year away, 2024 is 3 years away

		// Edge cases with century years
		{name: "prefer past from 2101 (after non-leap)", baseYear: 2101, preference: PreferPast, expected: 2096},
		{name: "prefer future from 2099 (before non-leap)", baseYear: 2099, preference: PreferFuture, expected: 2104},
		{name: "prefer current from 2100 (non-leap)", baseYear: 2100, preference: PreferCurrentPeriod, expected: 2104}, // 2096 is 4 away, 2104 is 4 away, prefer future = 2104
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findNearestLeapYear(tt.baseYear, tt.preference)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindYearClosestToRefWithPreference_February29(t *testing.T) {
	tests := []struct {
		name       string
		refDate    time.Time
		preference DatePreference
		expected   int
	}{
		// February 29 from non-leap year 2023
		{
			name:       "Feb 29 from 2023, prefer past",
			refDate:    time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferPast,
			expected:   2020,
		},
		{
			name:       "Feb 29 from 2023, prefer future",
			refDate:    time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferFuture,
			expected:   2024,
		},
		{
			name:       "Feb 29 from 2023, prefer current",
			refDate:    time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferCurrentPeriod,
			expected:   2024,
		},

		// February 29 from leap year 2024
		{
			name:       "Feb 29 from 2024, prefer past",
			refDate:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferPast,
			expected:   2024,
		},
		{
			name:       "Feb 29 from 2024, prefer future",
			refDate:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferFuture,
			expected:   2024,
		},
		{
			name:       "Feb 29 from 2024, prefer current",
			refDate:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferCurrentPeriod,
			expected:   2024,
		},

		// February 29 near century boundary (2100 is NOT a leap year)
		{
			name:       "Feb 29 from 2100, prefer past",
			refDate:    time.Date(2100, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferPast,
			expected:   2096,
		},
		{
			name:       "Feb 29 from 2100, prefer future",
			refDate:    time.Date(2100, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferFuture,
			expected:   2104,
		},

		// February 29 from 2000 (leap year divisible by 400)
		{
			name:       "Feb 29 from 2000, prefer current",
			refDate:    time.Date(2000, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferCurrentPeriod,
			expected:   2000,
		},

		// February 29 from 1900 (NOT a leap year)
		{
			name:       "Feb 29 from 1900, prefer past",
			refDate:    time.Date(1900, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferPast,
			expected:   1900, // Bounded at 1900
		},
		{
			name:       "Feb 29 from 1900, prefer future",
			refDate:    time.Date(1900, 3, 1, 12, 0, 0, 0, time.UTC),
			preference: PreferFuture,
			expected:   1904,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findYearClosestToRefWithPreference(tt.refDate, 29, 2, tt.preference)
			assert.Equal(t, tt.expected, result)
		})
	}
}
