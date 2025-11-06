//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common/refiners"
	"github.com/stretchr/testify/assert"
)

// TestENTimeExpression_ParsingTextOffset tests that the parser correctly identifies
// the index and text of time expressions with various offsets
func TestENTimeExpression_ParsingTextOffset(t *testing.T) {
	parser := NewENTimeExpressionParser(false)
	refDate := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		text          string
		expectedIndex int
		expectedText  string
	}{
		{
			name:          "time with leading spaces",
			text:          "  11 AM ",
			expectedIndex: 0, // After sanitization, leading spaces are trimmed
			expectedText:  "11 AM",
		},
		{
			name:          "time after year with 'at' keyword",
			text:          "2020 at  11 AM ",
			expectedIndex: 5,
			expectedText:  "at 11 AM", // After sanitization, multiple spaces are collapsed to single space
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedIndex, result.Index())
			assert.Equal(t, tt.expectedText, result.Text())
		})
	}
}

// TestENTimeExpression_Basic tests basic time expression parsing with hours, minutes, seconds
func TestENTimeExpression_Basic(t *testing.T) {
	parser := NewENTimeExpressionParser(false)
	refDate := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("20:32:13", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "20:32:13", result.Text())
	assert.Equal(t, 20, *result.Start().Get(kronos.ComponentHour))
	assert.Equal(t, 32, *result.Start().Get(kronos.ComponentMinute))
	assert.Equal(t, 13, *result.Start().Get(kronos.ComponentSecond))
	assert.Equal(t, int(kronos.MeridiemPM), *result.Start().Get(kronos.ComponentMeridiem))

	assert.Contains(t, result.Tags(), "parser/ENTimeExpressionParser")
	assert.Contains(t, result.Start().Tags(), "parser/ENTimeExpressionParser")
}

// TestENTimeExpression_WithClues tests time expressions with contextual clues like "at night", "in the afternoon"
func TestENTimeExpression_WithClues(t *testing.T) {
	parser := NewENTimeExpressionParser(false)
	refDate := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		text             string
		expectedHour     int
		expectedMeridiem int
	}{
		{
			name:             "1 at night",
			text:             "1 at night",
			expectedHour:     1,
			expectedMeridiem: int(kronos.MeridiemAM),
		},
		{
			name:             "1 in the afternoon",
			text:             "1 in the afternoon",
			expectedHour:     13,
			expectedMeridiem: int(kronos.MeridiemPM),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))
			assert.Equal(t, tt.expectedMeridiem, *result.Start().Get(kronos.ComponentMeridiem))
		})
	}
}

// TestENTimeExpression_AfterDate tests time expressions following date patterns
func TestENTimeExpression_AfterDate(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		text          string
		expectedYear  int
		expectedMonth int
		expectedDay   int
		expectedHour  int
		expectedMin   int
	}{
		{
			name:          "date space time",
			text:          "05/31/2024 14:15",
			expectedYear:  2024,
			expectedMonth: 5,
			expectedDay:   31,
			expectedHour:  14,
			expectedMin:   15,
		},
		{
			name:          "date colon time",
			text:          "05/31/2024:14:15",
			expectedYear:  2024,
			expectedMonth: 5,
			expectedDay:   31,
			expectedHour:  14,
			expectedMin:   15,
		},
		{
			name:          "date dash time",
			text:          "05/31/2024-14:15",
			expectedYear:  2024,
			expectedMonth: 5,
			expectedDay:   31,
			expectedHour:  14,
			expectedMin:   15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use full EN configuration to get date-time merging
			config := CreateConfiguration(false, false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))
			assert.Equal(t, tt.expectedMin, *result.Start().Get(kronos.ComponentMinute))
			assert.Equal(t, int(kronos.MeridiemPM), *result.Start().Get(kronos.ComponentMeridiem))

			assert.Contains(t, result.Tags(), "parser/ENTimeExpressionParser")
			assert.Contains(t, result.Start().Tags(), "parser/ENTimeExpressionParser")
		})
	}
}

// TestENTimeExpression_BeforeDate tests time expressions before date patterns
func TestENTimeExpression_BeforeDate(t *testing.T) {
	refDate := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		text             string
		expectedYear     int
		expectedMonth    int
		expectedDay      int
		expectedHour     int
		expectedMin      int
		expectedMeridiem int
	}{
		{
			name:             "time before slash date",
			text:             "14:15 05/31/2024",
			expectedYear:     2024,
			expectedMonth:    5,
			expectedDay:      31,
			expectedHour:     14,
			expectedMin:      15,
			expectedMeridiem: int(kronos.MeridiemPM),
		},
		{
			name:             "AM time with comma and month",
			text:             "8:23 AM, Jul 9",
			expectedYear:     2016,
			expectedMonth:    7,
			expectedDay:      9,
			expectedHour:     8,
			expectedMin:      23,
			expectedMeridiem: int(kronos.MeridiemAM),
		},
		{
			name:             "AM time with bullet separator",
			text:             "8:23 AM ∙ Jul 9",
			expectedYear:     2016,
			expectedMonth:    7,
			expectedDay:      9,
			expectedHour:     8,
			expectedMin:      23,
			expectedMeridiem: int(kronos.MeridiemAM),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use full EN configuration to get date-time merging
			config := CreateConfiguration(false, false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))
			assert.Equal(t, tt.expectedMin, *result.Start().Get(kronos.ComponentMinute))
			assert.Equal(t, tt.expectedMeridiem, *result.Start().Get(kronos.ComponentMeridiem))
		})
	}
}

// TestENTimeExpression_TimeRange tests time range expressions with various separators
func TestENTimeExpression_TimeRange(t *testing.T) {
	parser := NewENTimeExpressionParser(false)
	_ = time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC) // default refDate

	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		startHour     int
		startMin      int
		startSec      int
		startMeridiem int
		endHour       int
		endMin        int
		endSec        int
		endMeridiem   int
	}{
		{
			name:          "time range with dash",
			text:          "10:00:00 - 21:45:00",
			refDate:       time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC),
			startHour:     10,
			startMin:      0,
			startSec:      0,
			startMeridiem: int(kronos.MeridiemAM),
			endHour:       21,
			endMin:        45,
			endSec:        0,
			endMeridiem:   int(kronos.MeridiemPM),
		},
		{
			name:          "time range with until",
			text:          "10:00:00 until 21:45:00",
			refDate:       time.Date(2016, 10, 1, 11, 0, 0, 0, time.UTC),
			startHour:     10,
			startMin:      0,
			startSec:      0,
			startMeridiem: int(kronos.MeridiemAM),
			endHour:       21,
			endMin:        45,
			endSec:        0,
			endMeridiem:   int(kronos.MeridiemPM),
		},
		{
			name:          "time range with till",
			text:          "10:00:00 till 21:45:00",
			refDate:       time.Date(2016, 10, 1, 11, 0, 0, 0, time.UTC),
			startHour:     10,
			startMin:      0,
			startSec:      0,
			startMeridiem: int(kronos.MeridiemAM),
			endHour:       21,
			endMin:        45,
			endSec:        0,
			endMeridiem:   int(kronos.MeridiemPM),
		},
		{
			name:          "time range with through",
			text:          "10:00:00 through 21:45:00",
			refDate:       time.Date(2016, 10, 1, 11, 0, 0, 0, time.UTC),
			startHour:     10,
			startMin:      0,
			startSec:      0,
			startMeridiem: int(kronos.MeridiemAM),
			endHour:       21,
			endMin:        45,
			endSec:        0,
			endMeridiem:   int(kronos.MeridiemPM),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			assert.Contains(t, result.Tags(), "parser/ENTimeExpressionParser")

			assert.Equal(t, tt.startHour, *result.Start().Get(kronos.ComponentHour))
			assert.Equal(t, tt.startMin, *result.Start().Get(kronos.ComponentMinute))
			assert.Equal(t, tt.startSec, *result.Start().Get(kronos.ComponentSecond))
			assert.Equal(t, tt.startMeridiem, *result.Start().Get(kronos.ComponentMeridiem))
			assert.Contains(t, result.Start().Tags(), "parser/ENTimeExpressionParser")

			assert.NotNil(t, result.End())
			assert.Equal(t, tt.endHour, *result.End().Get(kronos.ComponentHour))
			assert.Equal(t, tt.endMin, *result.End().Get(kronos.ComponentMinute))
			assert.Equal(t, tt.endSec, *result.End().Get(kronos.ComponentSecond))
			assert.Equal(t, tt.endMeridiem, *result.End().Get(kronos.ComponentMeridiem))
			assert.Contains(t, result.End().Tags(), "parser/ENTimeExpressionParser")
		})
	}
}

// TestENTimeExpression_NonRange tests that invalid range patterns are not parsed as ranges
func TestENTimeExpression_NonRange(t *testing.T) {
	parser := NewENTimeExpressionParser(false)
	refDate := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("10:00:00 - 15/15", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "10:00:00", result.Text())
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentHour))
	assert.Equal(t, 0, *result.Start().Get(kronos.ComponentMinute))
	assert.Equal(t, 0, *result.Start().Get(kronos.ComponentSecond))
	assert.Equal(t, int(kronos.MeridiemAM), *result.Start().Get(kronos.ComponentMeridiem))
}

// TestENTimeExpression_CasualTimeNumber tests casual time expressions with time of day
func TestENTimeExpression_CasualTimeNumber(t *testing.T) {
	parser := NewENTimeExpressionParser(false)
	refDate := time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		text             string
		expectedYear     int
		expectedMonth    int
		expectedDay      int
		expectedHour     int
		expectedMin      *int
		expectedMeridiem *int
	}{
		{
			name:          "11 at night",
			text:          "11 at night",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   1,
			expectedHour:  23,
		},
		{
			name:          "11 tonight",
			text:          "11 tonight",
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   1,
			expectedHour:  23,
		},
		{
			name:             "6 in the morning",
			text:             "6 in the morning",
			expectedYear:     2016,
			expectedMonth:    10,
			expectedDay:      1,
			expectedHour:     6,
			expectedMin:      intPtr(0),
			expectedMeridiem: intPtr(int(kronos.MeridiemAM)),
		},
		{
			name:             "6 in the afternoon",
			text:             "6 in the afternoon",
			expectedYear:     2016,
			expectedMonth:    10,
			expectedDay:      1,
			expectedHour:     18,
			expectedMin:      intPtr(0),
			expectedMeridiem: intPtr(int(kronos.MeridiemPM)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))

			if tt.expectedMin != nil {
				assert.Equal(t, *tt.expectedMin, *result.Start().Get(kronos.ComponentMinute))
			}
			if tt.expectedMeridiem != nil {
				assert.Equal(t, *tt.expectedMeridiem, *result.Start().Get(kronos.ComponentMeridiem))
			}
		})
	}
}

// TestENTimeExpression_TimeRangeMeridiemHandling tests meridiem propagation in time ranges
func TestENTimeExpression_TimeRangeMeridiemHandling(t *testing.T) {
	parser := NewENTimeExpressionParser(false)
	_ = time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC) // default refDate

	tests := []struct {
		name                 string
		text                 string
		refDate              time.Time
		startYear            int
		startMonth           int
		startDay             int
		startHour            int
		startMeridiem        *int
		startMeridiemCertain *bool
		endYear              int
		endMonth             int
		endDay               int
		endHour              int
		endMeridiem          *int
		endMeridiemCertain   *bool
	}{
		{
			name:       "10 - 11 at night",
			text:       "10 - 11 at night",
			refDate:    time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC),
			startYear:  2016,
			startMonth: 10,
			startDay:   1,
			startHour:  22,
			endYear:    2016,
			endMonth:   10,
			endDay:     1,
			endHour:    23,
		},
		{
			name:          "8pm - 11",
			text:          "8pm - 11",
			refDate:       time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC),
			startYear:     2016,
			startMonth:    10,
			startDay:      1,
			startHour:     20,
			startMeridiem: intPtr(int(kronos.MeridiemPM)),
			endYear:       2016,
			endMonth:      10,
			endDay:        1,
			endHour:       23,
			endMeridiem:   intPtr(int(kronos.MeridiemPM)),
		},
		{
			name:          "8 - 11pm",
			text:          "8 - 11pm",
			refDate:       time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC),
			startYear:     2016,
			startMonth:    10,
			startDay:      1,
			startHour:     20,
			startMeridiem: intPtr(int(kronos.MeridiemPM)),
			endYear:       2016,
			endMonth:      10,
			endDay:        1,
			endHour:       23,
			endMeridiem:   intPtr(int(kronos.MeridiemPM)),
		},
		{
			name:          "7 - 8",
			text:          "7 - 8",
			refDate:       time.Date(2016, 10, 1, 8, 0, 0, 0, time.UTC),
			startYear:     2016,
			startMonth:    10,
			startDay:      1,
			startHour:     7,
			startMeridiem: intPtr(int(kronos.MeridiemAM)),
			endYear:       2016,
			endMonth:      10,
			endDay:        1,
			endHour:       8,
			endMeridiem:   intPtr(int(kronos.MeridiemAM)),
		},
		{
			name:                 "1pm-3",
			text:                 "1pm-3",
			refDate:              time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
			startYear:            2012,
			startMonth:           8,
			startDay:             10,
			startHour:            13,
			startMeridiem:        intPtr(int(kronos.MeridiemPM)),
			startMeridiemCertain: boolPtr(true),
			endYear:              2012,
			endMonth:             8,
			endDay:               10,
			endHour:              15,
			endMeridiem:          intPtr(int(kronos.MeridiemPM)),
			endMeridiemCertain:   boolPtr(true),
		},
		{
			name:                 "1am-3",
			text:                 "1am-3",
			refDate:              time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
			startYear:            2012,
			startMonth:           8,
			startDay:             10,
			startHour:            1,
			startMeridiem:        intPtr(int(kronos.MeridiemAM)),
			startMeridiemCertain: boolPtr(true),
			endYear:              2012,
			endMonth:             8,
			endDay:               10,
			endHour:              3,
			endMeridiem:          intPtr(int(kronos.MeridiemAM)),
			endMeridiemCertain:   boolPtr(false),
		},
		{
			name:                 "11pm-3 (crosses midnight)",
			text:                 "11pm-3",
			refDate:              time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
			startYear:            2012,
			startMonth:           8,
			startDay:             10,
			startHour:            23,
			startMeridiem:        intPtr(int(kronos.MeridiemPM)),
			startMeridiemCertain: boolPtr(true),
			endYear:              2012,
			endMonth:             8,
			endDay:               11,
			endHour:              3,
			endMeridiem:          intPtr(int(kronos.MeridiemAM)),
			endMeridiemCertain:   boolPtr(false),
		},
		{
			name:               "12-3am",
			text:               "12-3am",
			refDate:            time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
			startYear:          2012,
			startMonth:         8,
			startDay:           10,
			startHour:          0,
			startMeridiem:      intPtr(int(kronos.MeridiemAM)),
			endYear:            2012,
			endMonth:           8,
			endDay:             10,
			endHour:            3,
			endMeridiem:        intPtr(int(kronos.MeridiemAM)),
			endMeridiemCertain: boolPtr(true),
		},
		{
			name:               "12-3pm",
			text:               "12-3pm",
			refDate:            time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
			startYear:          2012,
			startMonth:         8,
			startDay:           10,
			startHour:          12,
			startMeridiem:      intPtr(int(kronos.MeridiemPM)),
			endYear:            2012,
			endMonth:           8,
			endDay:             10,
			endHour:            15,
			endMeridiem:        intPtr(int(kronos.MeridiemPM)),
			endMeridiemCertain: boolPtr(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.text, result.Text())
			assert.Equal(t, tt.startYear, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, tt.startMonth, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.startDay, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, tt.startHour, *result.Start().Get(kronos.ComponentHour))

			if tt.startMeridiem != nil {
				assert.Equal(t, *tt.startMeridiem, *result.Start().Get(kronos.ComponentMeridiem))
			}
			if tt.startMeridiemCertain != nil {
				assert.Equal(t, *tt.startMeridiemCertain, result.Start().IsCertain(kronos.ComponentMeridiem))
			}

			assert.NotNil(t, result.End())
			assert.Equal(t, tt.endYear, *result.End().Get(kronos.ComponentYear))
			assert.Equal(t, tt.endMonth, *result.End().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.endDay, *result.End().Get(kronos.ComponentDay))
			assert.Equal(t, tt.endHour, *result.End().Get(kronos.ComponentHour))

			if tt.endMeridiem != nil {
				assert.Equal(t, *tt.endMeridiem, *result.End().Get(kronos.ComponentMeridiem))
			}
			if tt.endMeridiemCertain != nil {
				assert.Equal(t, *tt.endMeridiemCertain, result.End().IsCertain(kronos.ComponentMeridiem))
			}
		})
	}
}

// TestENTimeExpression_TimeRangeNextDay tests time ranges that span to the next day
func TestENTimeExpression_TimeRangeNextDay(t *testing.T) {
	refDate := time.Date(2017, 7, 7, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		text          string
		startDay      int
		startMonth    int
		startYear     int
		startHour     int
		startMeridiem int
	}{
		{
			name:          "Dec 31 10pm - 1am",
			text:          "December 31, 2022 10:00 pm - 1:00 am",
			startDay:      31,
			startMonth:    12,
			startYear:     2022,
			startHour:     22,
			startMeridiem: int(kronos.MeridiemPM),
		},
		{
			name:          "Dec 31 10pm - midnight",
			text:          "December 31, 2022 10:00 pm - 12:00 am",
			startDay:      31,
			startMonth:    12,
			startYear:     2022,
			startHour:     22,
			startMeridiem: int(kronos.MeridiemPM),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use full EN configuration to get date-time merging
			config := CreateConfiguration(false, false)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.startDay, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, tt.startMonth, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.startYear, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, tt.startHour, *result.Start().Get(kronos.ComponentHour))
			assert.Equal(t, tt.startMeridiem, *result.Start().Get(kronos.ComponentMeridiem))
		})
	}
}

// TestENTimeExpression_CasualPositive tests parsing casual time expressions
func TestENTimeExpression_CasualPositive(t *testing.T) {
	parser := NewENTimeExpressionParser(false)

	tests := []struct {
		name         string
		text         string
		expectedText string
		expectedHour int
		expectedMin  *int
	}{
		{
			name:         "at 1",
			text:         "at 1",
			expectedText: "at 1",
			expectedHour: 1,
		},
		{
			name:         "at 12",
			text:         "at 12",
			expectedText: "at 12",
			expectedHour: 12,
		},
		{
			name:         "at 12.30",
			text:         "at 12.30",
			expectedText: "at 12.30",
			expectedHour: 12,
			expectedMin:  intPtr(30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, time.Time{}, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedText, result.Text())
			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))

			if tt.expectedMin != nil {
				assert.Equal(t, *tt.expectedMin, *result.Start().Get(kronos.ComponentMinute))
			}
		})
	}
}

// TestENTimeExpression_NegativeYearLike tests that year-like patterns are not parsed as times
func TestENTimeExpression_NegativeYearLike(t *testing.T) {
	parser := NewENTimeExpressionParser(false)

	tests := []string{
		"2020",
		"2020  ",
		"2019 to 2020",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(text, time.Time{}, nil)

			assert.Empty(t, results, "Should not parse: %s", text)
		})
	}
}

// TestENTimeExpression_NegativeAtSomeNumbers tests that "at" with non-time numbers is not parsed
func TestENTimeExpression_NegativeAtSomeNumbers(t *testing.T) {
	parser := NewENTimeExpressionParser(false)

	tests := []string{
		"I'm at 101,194 points!",
		"I'm at 101 points!",
		"I'm at 10.1",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(text, time.Time{}, nil)

			assert.Empty(t, results, "Should not parse: %s", text)
		})
	}
}

// TestENTimeExpression_NegativeAtRanges tests that invalid "at" ranges are not parsed
func TestENTimeExpression_NegativeAtRanges(t *testing.T) {
	parser := NewENTimeExpressionParser(false)

	tests := []string{
		"I'm at 10.1 - 10.12",
		"I'm at 10 - 10.1",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(text, time.Time{}, nil)

			assert.Empty(t, results, "Should not parse: %s", text)
		})
	}
}

// TestENTimeExpression_StrictNegative tests strict mode rejects casual patterns
func TestENTimeExpression_StrictNegative(t *testing.T) {
	parser := NewENTimeExpressionParser(true) // Strict mode

	tests := []string{
		"I'm at 101,194 points!",
		"I'm at 101 points!",
		"I'm at 10.1",
		"I'm at 10",
		"2020",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(text, time.Time{}, nil)

			assert.Empty(t, results, "Strict mode should not parse: %s", text)
		})
	}
}

// TestENTimeExpression_StrictNegativeRanges tests strict mode rejects casual range patterns
func TestENTimeExpression_StrictNegativeRanges(t *testing.T) {
	parser := NewENTimeExpressionParser(true) // Strict mode

	tests := []string{
		"I'm at 10.1 - 10.12",
		"I'm at 10 - 10.1",
		"I'm at 10 - 20",
		"7-730",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(text, time.Time{}, nil)

			assert.Empty(t, results, "Strict mode should not parse: %s", text)
		})
	}
}

// TestENTimeExpression_ForwardDateFlag tests the forwardDate option
func TestENTimeExpression_ForwardDateFlag(t *testing.T) {
	parser := NewENTimeExpressionParser(false)

	tests := []struct {
		name          string
		text          string
		refDate       time.Time
		expectedYear  int
		expectedMonth int
		expectedDay   int
		expectedHour  int
		endExpected   bool
		endYear       *int
		endMonth      *int
		endDay        *int
		endHour       *int
	}{
		{
			name:          "1am after reference time",
			text:          "1am",
			refDate:       time.Date(2022, 5, 26, 1, 57, 0, 0, time.UTC),
			expectedYear:  2022,
			expectedMonth: 5,
			expectedDay:   27,
			expectedHour:  1,
		},
		{
			name:          "11am after reference time",
			text:          "11am",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   2,
			expectedHour:  11,
		},
		{
			name:          "11am to 1am range",
			text:          "  11am to 1am  ",
			refDate:       time.Date(2016, 10, 1, 12, 0, 0, 0, time.UTC),
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   2,
			expectedHour:  11,
			endExpected:   true,
			endYear:       intPtr(2016),
			endMonth:      intPtr(10),
			endDay:        intPtr(3),
			endHour:       intPtr(1),
		},
		{
			name:          "10am to 12pm same day",
			text:          "  10am to 12pm  ",
			refDate:       time.Date(2016, 10, 1, 11, 0, 0, 0, time.UTC),
			expectedYear:  2016,
			expectedMonth: 10,
			expectedDay:   2,
			expectedHour:  10,
			endExpected:   true,
			endYear:       intPtr(2016),
			endMonth:      intPtr(10),
			endDay:        intPtr(2),
			endHour:       intPtr(12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{
				Parsers:  []kronos.Parser{parser},
				Refiners: []kronos.Refiner{refiners.NewForwardDateRefiner()},
			}
			chrono := kronos.NewChrono(config)
			option := &kronos.ParsingOption{ForwardDate: true}
			results := chrono.Parse(tt.text, tt.refDate, option)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))

			if tt.endExpected {
				assert.NotNil(t, result.End())
				if tt.endYear != nil {
					assert.Equal(t, *tt.endYear, *result.End().Get(kronos.ComponentYear))
				}
				if tt.endMonth != nil {
					assert.Equal(t, *tt.endMonth, *result.End().Get(kronos.ComponentMonth))
				}
				if tt.endDay != nil {
					assert.Equal(t, *tt.endDay, *result.End().Get(kronos.ComponentDay))
				}
				if tt.endHour != nil {
					assert.Equal(t, *tt.endHour, *result.End().Get(kronos.ComponentHour))
				}
			}
		})
	}
}

// TestENTimeExpression_ForwardDateFlagTimezone tests forwardDate with timezone
func TestENTimeExpression_ForwardDateFlagTimezone(t *testing.T) {
	parser := NewENTimeExpressionParser(false)

	// CDT is UTC-5, so Wed May 26 2022 01:57:00 GMT-0500 is May 26 at 01:57 CDT
	instant := time.Date(2022, 5, 26, 6, 57, 0, 0, time.UTC) // UTC equivalent
	timezone := "CDT"

	ref := kronos.ParsingReference{
		Instant:  &instant,
		Timezone: &timezone,
	}

	config := &kronos.Configuration{
		Parsers:  []kronos.Parser{parser},
		Refiners: []kronos.Refiner{refiners.NewForwardDateRefiner()},
	}
	chrono := kronos.NewChrono(config)
	option := &kronos.ParsingOption{ForwardDate: true}
	results := chrono.Parse("1am", ref, option)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, 2022, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 5, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 27, *result.Start().Get(kronos.ComponentDay))
	assert.Equal(t, 1, *result.Start().Get(kronos.ComponentHour))
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
