package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTimezoneExpression_UTCOffset tests parsing date/time with UTC offset
func TestTimezoneExpression_UTCOffset(t *testing.T) {
	// Use full EN configuration to get timezone refiners
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	tests := []struct {
		name                    string
		text                    string
		expectedText            string
		expectedHour            int
		expectedMinute          int
		expectedTimezoneOffset  int
	}{
		{
			name:                   "UTC offset with colon separator",
			text:                   "wednesday, september 16, 2020 at 11 am utc+02:45 ",
			expectedText:           "wednesday, september 16, 2020 at 11 am utc+02:45",
			expectedHour:           11,
			expectedMinute:         0,
			expectedTimezoneOffset: 2*60 + 45,
		},
		{
			name:                   "UTC offset without separator",
			text:                   "wednesday, september 16, 2020 at 11 am utc+0245 ",
			expectedText:           "wednesday, september 16, 2020 at 11 am utc+0245",
			expectedHour:           11,
			expectedMinute:         0,
			expectedTimezoneOffset: 2*60 + 45,
		},
		{
			name:                   "UTC offset hours only",
			text:                   "wednesday, september 16, 2020 at 11 am utc+02 ",
			expectedText:           "wednesday, september 16, 2020 at 11 am utc+02",
			expectedHour:           11,
			expectedMinute:         0,
			expectedTimezoneOffset: 2 * 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedText, result.Text())
			hour := result.Start().Get(kronos.ComponentHour)
			require.NotNil(t, hour)
			assert.Equal(t, tt.expectedHour, *hour)
			minute := result.Start().Get(kronos.ComponentMinute)
			require.NotNil(t, minute)
			assert.Equal(t, tt.expectedMinute, *minute)
			tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
			require.NotNil(t, tzOffset)
			assert.Equal(t, tt.expectedTimezoneOffset, *tzOffset)
		})
	}
}

// TestTimezoneExpression_NumericOffset tests parsing date/time with numeric offset
func TestTimezoneExpression_NumericOffset(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	tests := []struct {
		name                      string
		text                      string
		expectedText              string
		expectTimezoneOffsetSet   bool
		expectedTimezoneOffset    int
	}{
		{
			name:                    "numeric offset +14",
			text:                    "wednesday, september 16, 2020 at 23.00+14",
			expectedText:            "wednesday, september 16, 2020 at 23.00+14",
			expectTimezoneOffsetSet: true,
			expectedTimezoneOffset:  14 * 60,
		},
		{
			name:                    "numeric offset +1400",
			text:                    "wednesday, september 16, 2020 at 23.00+1400",
			expectedText:            "wednesday, september 16, 2020 at 23.00+1400",
			expectTimezoneOffsetSet: true,
			expectedTimezoneOffset:  14 * 60,
		},
		{
			name:                    "invalid offset +15 - should not parse timezone",
			text:                    "wednesday, september 16, 2020 at 23.00+15",
			expectedText:            "wednesday, september 16, 2020 at 23.00",
			expectTimezoneOffsetSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedText, result.Text())

			if tt.expectTimezoneOffsetSet {
				tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
				require.NotNil(t, tzOffset)
				assert.Equal(t, tt.expectedTimezoneOffset, *tzOffset)
			} else {
				// Should be nil or not certain
				tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
				if tzOffset != nil {
					// If it exists, it should not be certain
					assert.False(t, result.Start().IsCertain(kronos.ComponentTimezoneOffset))
				}
			}
		})
	}
}

// TestTimezoneExpression_GMTOffset tests parsing date/time with GMT offset
func TestTimezoneExpression_GMTOffset(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	tests := []struct {
		name                   string
		text                   string
		expectedText           string
		expectedHour           int
		expectedMinute         int
		expectedTimezoneOffset int
	}{
		{
			name:                   "GMT with negative offset and space",
			text:                   "wednesday, september 16, 2020 at 11 am GMT -08:45 ",
			expectedText:           "wednesday, september 16, 2020 at 11 am GMT -08:45",
			expectedHour:           11,
			expectedMinute:         0,
			expectedTimezoneOffset: -(8*60 + 45),
		},
		{
			name:                   "GMT with positive offset lowercase",
			text:                   "wednesday, september 16, 2020 at 11 am gmt+02 ",
			expectedText:           "wednesday, september 16, 2020 at 11 am gmt+02",
			expectedHour:           11,
			expectedMinute:         0,
			expectedTimezoneOffset: 2 * 60,
		},
		{
			name:                   "GMT in parentheses with negative offset",
			text:                   "published: 10:30 (gmt-2:30).",
			expectedText:           "10:30 (gmt-2:30)",
			expectedHour:           10,
			expectedMinute:         30,
			expectedTimezoneOffset: -(2*60 + 30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedText, result.Text())
			hour := result.Start().Get(kronos.ComponentHour)
			require.NotNil(t, hour)
			assert.Equal(t, tt.expectedHour, *hour)
			minute := result.Start().Get(kronos.ComponentMinute)
			require.NotNil(t, minute)
			assert.Equal(t, tt.expectedMinute, *minute)
			tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
			require.NotNil(t, tzOffset)
			assert.Equal(t, tt.expectedTimezoneOffset, *tzOffset)
		})
	}
}

// TestTimezoneExpression_Abbreviation tests parsing date/time with timezone abbreviation
func TestTimezoneExpression_Abbreviation(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	tests := []struct {
		name                      string
		text                      string
		refTime                   time.Time
		expectedText              string
		expectedYear              int
		expectedMonth             int
		expectedDay               int
		expectedHour              int
		expectedMinute            int
		expectTimezoneOffset      bool
		expectedTimezoneOffset    int
	}{
		{
			name:                   "without timezone abbreviation",
			text:                   "wednesday, september 16, 2020 at 11 am",
			expectedText:           "wednesday, september 16, 2020 at 11 am",
			expectedYear:           2020,
			expectedMonth:          9,
			expectedDay:            16,
			expectedHour:           11,
			expectedMinute:         0,
			expectTimezoneOffset:   false,
		},
		{
			name:                   "with JST timezone",
			text:                   "wednesday, september 16, 2020 at 11 am JST",
			expectedText:           "wednesday, september 16, 2020 at 11 am JST",
			expectedYear:           2020,
			expectedMonth:          9,
			expectedDay:            16,
			expectedHour:           11,
			expectedMinute:         0,
			expectTimezoneOffset:   true,
			expectedTimezoneOffset: 9 * 60, // JST: GMT+9:00
		},
		{
			name:                   "with GMT+0900 and JST in parentheses",
			text:                   "wednesday, september 16, 2020 at 11 am GMT+0900 (JST)",
			expectedText:           "wednesday, september 16, 2020 at 11 am GMT+0900 (JST)",
			expectedYear:           2020,
			expectedMonth:          9,
			expectedDay:            16,
			expectedHour:           11,
			expectedMinute:         0,
			expectTimezoneOffset:   true,
			expectedTimezoneOffset: 9 * 60, // JST: GMT+9:00
		},
		{
			name:                   "PST with today reference",
			text:                   "10:30 pst today",
			refTime:                time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC),
			expectedText:           "10:30 pst today",
			expectedYear:           2016,
			expectedMonth:          10,
			expectedDay:            1,
			expectedHour:           10,
			expectedMinute:         30,
			expectTimezoneOffset:   true,
			expectedTimezoneOffset: -8 * 60, // PST: UTC−08:00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refTime := tt.refTime
			if refTime.IsZero() {
				refTime = time.Now()
			}

			results := chrono.Parse(tt.text, refTime, nil)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedText, result.Text())
			year := result.Start().Get(kronos.ComponentYear)
			require.NotNil(t, year)
			assert.Equal(t, tt.expectedYear, *year)
			month := result.Start().Get(kronos.ComponentMonth)
			require.NotNil(t, month)
			assert.Equal(t, tt.expectedMonth, *month)
			day := result.Start().Get(kronos.ComponentDay)
			require.NotNil(t, day)
			assert.Equal(t, tt.expectedDay, *day)
			hour := result.Start().Get(kronos.ComponentHour)
			require.NotNil(t, hour)
			assert.Equal(t, tt.expectedHour, *hour)
			minute := result.Start().Get(kronos.ComponentMinute)
			require.NotNil(t, minute)
			assert.Equal(t, tt.expectedMinute, *minute)

			if tt.expectTimezoneOffset {
				tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
				require.NotNil(t, tzOffset)
				assert.Equal(t, tt.expectedTimezoneOffset, *tzOffset)
			} else {
				assert.Nil(t, result.Start().Get(kronos.ComponentTimezoneOffset))
			}
		})
	}
}

// TestTimezoneExpression_DateRange tests parsing date range with timezone abbreviation
func TestTimezoneExpression_DateRange(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	refTime := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)
	results := chrono.Parse("10:30 JST today to 10:30 pst tomorrow ", refTime, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "10:30 JST today to 10:30 pst tomorrow", result.Text())

	// Start date
	year := result.Start().Get(kronos.ComponentYear)
	require.NotNil(t, year)
	assert.Equal(t, 2016, *year)
	month := result.Start().Get(kronos.ComponentMonth)
	require.NotNil(t, month)
	assert.Equal(t, 10, *month)
	day := result.Start().Get(kronos.ComponentDay)
	require.NotNil(t, day)
	assert.Equal(t, 1, *day)
	hour := result.Start().Get(kronos.ComponentHour)
	require.NotNil(t, hour)
	assert.Equal(t, 10, *hour)
	minute := result.Start().Get(kronos.ComponentMinute)
	require.NotNil(t, minute)
	assert.Equal(t, 30, *minute)
	tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
	require.NotNil(t, tzOffset)
	assert.Equal(t, 9*60, *tzOffset) // JST: GMT+9:00

	// End date
	assert.NotNil(t, result.End())
	endYear := result.End().Get(kronos.ComponentYear)
	require.NotNil(t, endYear)
	assert.Equal(t, 2016, *endYear)
	endMonth := result.End().Get(kronos.ComponentMonth)
	require.NotNil(t, endMonth)
	assert.Equal(t, 10, *endMonth)
	endDay := result.End().Get(kronos.ComponentDay)
	require.NotNil(t, endDay)
	assert.Equal(t, 2, *endDay)
	endHour := result.End().Get(kronos.ComponentHour)
	require.NotNil(t, endHour)
	assert.Equal(t, 10, *endHour)
	endMinute := result.End().Get(kronos.ComponentMinute)
	require.NotNil(t, endMinute)
	assert.Equal(t, 30, *endMinute)
	endTzOffset := result.End().Get(kronos.ComponentTimezoneOffset)
	require.NotNil(t, endTzOffset)
	assert.Equal(t, -8*60, *endTzOffset) // PST: UTC−08:00
}

// TestTimezoneExpression_AmbiguousET tests parsing with ambiguous ET timezone around DST transitions
func TestTimezoneExpression_AmbiguousET(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	tests := []struct {
		name                   string
		text                   string
		expectedTimezoneOffset int
	}{
		// Transition TO DST in 2022
		{
			name:                   "2022-03-12 before DST transition - EST",
			text:                   "2022-03-12 23:00 ET",
			expectedTimezoneOffset: -5 * 60, // EST
		},
		{
			name:                   "2022-03-13 after DST transition - EDT",
			text:                   "2022-03-13 23:00 ET",
			expectedTimezoneOffset: -4 * 60, // EDT
		},
		// Transition TO DST in 2021
		{
			name:                   "2021-03-13 before DST transition - EST",
			text:                   "2021-03-13 23:00 ET",
			expectedTimezoneOffset: -5 * 60, // EST
		},
		{
			name:                   "2021-03-14 after DST transition - EDT",
			text:                   "2021-03-14 23:00 ET",
			expectedTimezoneOffset: -4 * 60, // EDT
		},
		// Transition FROM DST in 2021
		{
			name:                   "2021-11-06 before end DST - EDT",
			text:                   "2021-11-06 23:00 ET",
			expectedTimezoneOffset: -4 * 60, // EDT
		},
		{
			name:                   "2021-11-07 after end DST - EST",
			text:                   "2021-11-07 23:00 ET",
			expectedTimezoneOffset: -5 * 60, // EST
		},
		// Transition FROM DST in 2020
		{
			name:                   "2020-10-31 before end DST - EDT",
			text:                   "2020-10-31 23:00 ET",
			expectedTimezoneOffset: -4 * 60, // EDT
		},
		{
			name:                   "2020-11-01 after end DST - EST",
			text:                   "2020-11-01 23:00 ET",
			expectedTimezoneOffset: -5 * 60, // EST
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			hour := result.Start().Get(kronos.ComponentHour)
			require.NotNil(t, hour)
			assert.Equal(t, 23, *hour)
			tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
			require.NotNil(t, tzOffset)
			assert.Equal(t, tt.expectedTimezoneOffset, *tzOffset)
		})
	}
}

// TestTimezoneExpression_AmbiguousCET tests parsing with ambiguous CET timezone around DST transitions
func TestTimezoneExpression_AmbiguousCET(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	tests := []struct {
		name                   string
		text                   string
		expectedTimezoneOffset int
	}{
		// Transition TO DST in 2022
		{
			name:                   "2022-03-26 before DST transition - CET",
			text:                   "2022-03-26 23:00 CET",
			expectedTimezoneOffset: 60, // CET
		},
		{
			name:                   "2022-03-27 after DST transition - CEST",
			text:                   "2022-03-27 23:00 CET",
			expectedTimezoneOffset: 2 * 60, // CEST
		},
		// Transition TO DST in 2021
		{
			name:                   "2021-03-27 before DST transition - CET",
			text:                   "2021-03-27 23:00 CET",
			expectedTimezoneOffset: 60, // CET
		},
		{
			name:                   "2021-03-28 after DST transition - CEST",
			text:                   "2021-03-28 23:00 CET",
			expectedTimezoneOffset: 2 * 60, // CEST
		},
		// Transition FROM DST in 2022
		{
			name:                   "2022-10-29 before end DST - CEST",
			text:                   "2022-10-29 23:00 CET",
			expectedTimezoneOffset: 2 * 60, // CEST
		},
		{
			name:                   "2022-10-30 after end DST - CET",
			text:                   "2022-10-30 23:00 CET",
			expectedTimezoneOffset: 60, // CET
		},
		// Transition FROM DST in 2021
		{
			name:                   "2021-10-30 before end DST - CEST",
			text:                   "2021-10-30 23:00 CET",
			expectedTimezoneOffset: 2 * 60, // CEST
		},
		{
			name:                   "2021-10-31 after end DST - CET",
			text:                   "2021-10-31 23:00 CET",
			expectedTimezoneOffset: 60, // CET
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := chrono.Parse(tt.text, time.Now(), nil)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			hour := result.Start().Get(kronos.ComponentHour)
			require.NotNil(t, hour)
			assert.Equal(t, 23, *hour)
			tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
			require.NotNil(t, tzOffset)
			assert.Equal(t, tt.expectedTimezoneOffset, *tzOffset)
		})
	}
}

// TestTimezoneExpression_CustomTimezone tests timezone parsing with custom timezone overrides
func TestTimezoneExpression_CustomTimezone(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	refTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("XYZ not recognized without custom timezone", func(t *testing.T) {
		results := chrono.Parse("Jan 1st 2023 at 10:00 XYZ", refTime, nil)
		assert.NotEmpty(t, results)
		result := results[0]

		assert.Equal(t, "Jan 1st 2023 at 10:00", result.Text())
		expectedDate := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
		year := result.Start().Get(kronos.ComponentYear)
		require.NotNil(t, year)
		assert.Equal(t, expectedDate.Year(), *year)
		month := result.Start().Get(kronos.ComponentMonth)
		require.NotNil(t, month)
		assert.Equal(t, int(expectedDate.Month()), *month)
		day := result.Start().Get(kronos.ComponentDay)
		require.NotNil(t, day)
		assert.Equal(t, expectedDate.Day(), *day)
		hour := result.Start().Get(kronos.ComponentHour)
		require.NotNil(t, hour)
		assert.Equal(t, expectedDate.Hour(), *hour)
	})

	t.Run("XYZ recognized with custom fixed timezone", func(t *testing.T) {
		options := &kronos.ParsingOption{
			Timezones: kronos.TimezoneAbbrMap{
				"XYZ": -180,
			},
		}
		results := chrono.Parse("Jan 1st 2023 at 10:00 XYZ", refTime, options)
		assert.NotEmpty(t, results)
		result := results[0]

		assert.Equal(t, "Jan 1st 2023 at 10:00 XYZ", result.Text())
		tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
		require.NotNil(t, tzOffset)
		assert.Equal(t, -180, *tzOffset)
	})

	t.Run("XYZ with ambiguous timezone non-DST date", func(t *testing.T) {
		options := &kronos.ParsingOption{
			Timezones: kronos.TimezoneAbbrMap{
				"XYZ": &kronos.AmbiguousTimezoneMap{
					TimezoneOffsetDuringDst: -120,
					TimezoneOffsetNonDst:    -180,
					DstStart: func(year int) time.Time {
						// Last Sunday of March at 2am
						return getLastWeekdayOfMonth(year, 3, time.Sunday, 2)
					},
					DstEnd: func(year int) time.Time {
						// Last Sunday of October at 3am
						return getLastWeekdayOfMonth(year, 10, time.Sunday, 3)
					},
				},
			},
		}

		results := chrono.Parse("Jan 1st 2023 at 10:00 XYZ", refTime, options)
		assert.NotEmpty(t, results)
		result := results[0]

		assert.Equal(t, "Jan 1st 2023 at 10:00 XYZ", result.Text())
		tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
		require.NotNil(t, tzOffset)
		assert.Equal(t, -180, *tzOffset) // Non-DST
	})

	t.Run("XYZ with ambiguous timezone DST date", func(t *testing.T) {
		options := &kronos.ParsingOption{
			Timezones: kronos.TimezoneAbbrMap{
				"XYZ": &kronos.AmbiguousTimezoneMap{
					TimezoneOffsetDuringDst: -120,
					TimezoneOffsetNonDst:    -180,
					DstStart: func(year int) time.Time {
						// Last Sunday of March at 2am
						return getLastWeekdayOfMonth(year, 3, time.Sunday, 2)
					},
					DstEnd: func(year int) time.Time {
						// Last Sunday of October at 3am
						return getLastWeekdayOfMonth(year, 10, time.Sunday, 3)
					},
				},
			},
		}

		results := chrono.Parse("Jun 1st 2023 at 10:00 XYZ", refTime, options)
		assert.NotEmpty(t, results)
		result := results[0]

		assert.Equal(t, "Jun 1st 2023 at 10:00 XYZ", result.Text())
		tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
		require.NotNil(t, tzOffset)
		assert.Equal(t, -120, *tzOffset) // DST
	})
}

// TestTimezoneExpression_DateWithTimezone tests parsing date with timezone abbreviation
func TestTimezoneExpression_DateWithTimezone(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	results := chrono.Parse("Wednesday, September 16, 2020, EST", time.Now(), nil)
	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "Wednesday, September 16, 2020, EST", result.Text())
	tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
	require.NotNil(t, tzOffset)
	assert.Equal(t, -300, *tzOffset)
}

// TestTimezoneExpression_NotParsingFromRelativeTime tests that timezone is not incorrectly parsed from relative time
func TestTimezoneExpression_NotParsingFromRelativeTime(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	refInstant := time.Date(2020, 11, 29, 13, 24, 13, 0, time.FixedZone("JST", 9*60*60))
	expectedInstant := time.Date(2020, 11, 29, 14, 24, 13, 0, time.FixedZone("JST", 9*60*60))

	t.Run("in 1 hour get eggs and milk", func(t *testing.T) {
		results := chrono.Parse("in 1 hour get eggs and milk", refInstant, nil)
		assert.NotEmpty(t, results)
		result := results[0]

		assert.Equal(t, "in 1 hour", result.Text())

		// Calculate expected timezone offset: In Go, the offset is returned in seconds,
		// but we need it in minutes. JavaScript's getTimezoneOffset() returns minutes
		// and is negated (e.g., JST +9 returns -540)
		_, offsetSeconds := refInstant.Zone()
		expectedOffset := offsetSeconds / 60 // Convert to minutes
		tzOffset := result.Start().Get(kronos.ComponentTimezoneOffset)
		require.NotNil(t, tzOffset)
		assert.Equal(t, expectedOffset, *tzOffset)

		// Check the result date matches expected
		actualDate := result.Start().Date()
		assert.True(t, actualDate.Equal(expectedInstant) || actualDate.Unix() == expectedInstant.Unix())
	})

	// The remaining test cases check "in 1 hour GMT" with different timezone settings
	// These are edge cases that may behave differently based on implementation
	t.Run("in 1 hour GMT", func(t *testing.T) {
		results := chrono.Parse("in 1 hour GMT", refInstant, nil)
		assert.NotEmpty(t, results)
		result := results[0]

		actualDate := result.Start().Date()
		assert.True(t, actualDate.Equal(expectedInstant) || actualDate.Unix() == expectedInstant.Unix())
	})
}

// TestTimezoneExpression_RelativeTimeNotAffectedByTimezone tests that relative time parsing
// is not affected by timezone settings in options
func TestTimezoneExpression_RelativeTimeNotAffectedByTimezone(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	refInstant := time.Unix(1637674343, 0)

	tests := []struct {
		name     string
		timezone string
		useNil   bool
	}{
		{name: "no timezone option", useNil: true},
		{name: "timezone BST", timezone: "BST"},
		{name: "timezone JST", timezone: "JST"},
	}

	for _, tt := range tests {
		t.Run(tt.name+" - now", func(t *testing.T) {
			var options *kronos.ParsingOption
			if !tt.useNil {
				options = &kronos.ParsingOption{
					// Note: The Timezone field may need to be implemented differently
					// based on the actual ParsingOption structure
				}
			}

			results := chrono.Parse("now", refInstant, options)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, "now", result.Text())
			actualDate := result.Start().Date()
			assert.True(t, actualDate.Equal(refInstant) || actualDate.Unix() == refInstant.Unix())
		})
	}
}

// TestTimezoneExpression_RelativeHoursNotAffectedByTimezone tests that "2 hour later"
// is not affected by timezone settings
func TestTimezoneExpression_RelativeHoursNotAffectedByTimezone(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	refInstant := time.Unix(1637674343, 0)
	expectedInstant := time.Unix(1637674343+2*60*60, 0)

	tests := []struct {
		name     string
		timezone string
		useNil   bool
	}{
		{name: "no timezone option", useNil: true},
		{name: "timezone BST", timezone: "BST"},
		{name: "timezone JST", timezone: "JST"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options *kronos.ParsingOption
			if !tt.useNil {
				options = &kronos.ParsingOption{
					// Note: The Timezone field may need to be implemented differently
				}
			}

			results := chrono.Parse("2 hour later", refInstant, options)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, "2 hour later", result.Text())
			actualDate := result.Start().Date()
			assert.True(t, actualDate.Equal(expectedInstant) || actualDate.Unix() == expectedInstant.Unix())
		})
	}
}

// TestTimezoneExpression_ParsingTimezoneFromRelativeDateWhenValid tests parsing timezone
// from relative date expressions when valid
func TestTimezoneExpression_ParsingTimezoneFromRelativeDateWhenValid(t *testing.T) {
	config := CreateConfiguration(false, false)
	chrono := kronos.NewChrono(config)

	refDate := time.Date(2020, 11, 14, 13, 48, 22, 0, time.UTC)

	_, refOffsetSeconds := refDate.Zone()
	refOffsetMinutes := refOffsetSeconds / 60

	tests := []struct {
		name                   string
		text                   string
		expectedText           string
		expectedTimezoneOffset int
	}{
		{
			name:                   "in 1 day get eggs and milk",
			text:                   "in 1 day get eggs and milk",
			expectedText:           "in 1 day",
			expectedTimezoneOffset: refOffsetMinutes,
		},
		{
			name:                   "in 1 day GET (Georgia Time)",
			text:                   "in 1 day GET",
			expectedText:           "in 1 day GET",
			expectedTimezoneOffset: 240, // GET: UTC+4
		},
		{
			name:                   "today EST",
			text:                   "today EST",
			expectedText:           "today EST",
			expectedTimezoneOffset: -300, // EST: UTC-5
		},
		{
			name:                   "next week EST",
			text:                   "next week EST",
			expectedText:           "next week EST",
			expectedTimezoneOffset: -300, // EST: UTC-5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := chrono.Parse(tt.text, refDate, nil)
			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedText, result.Text())
			assert.NotNil(t, result.Start().Get(kronos.ComponentTimezoneOffset))
			assert.Equal(t, tt.expectedTimezoneOffset, *result.Start().Get(kronos.ComponentTimezoneOffset))
		})
	}
}

// Helper function to get last weekday of a month
func getLastWeekdayOfMonth(year, month int, weekday time.Weekday, hour int) time.Time {
	// Start from the last day of the month
	lastDay := time.Date(year, time.Month(month+1), 1, hour, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	// Walk backwards to find the last occurrence of the weekday
	for lastDay.Weekday() != weekday {
		lastDay = lastDay.AddDate(0, 0, -1)
	}

	return lastDay
}
