package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENCompactFormatParser(t *testing.T) {
	parser := NewENCompactFormatParser()
	refDate := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		text          string
		expectedYear  *int
		expectedMonth *int
		expectedDay   *int
		expectedHour  *int
		expectedMin   *int
		expectedSec   *int
		shouldParse   bool
	}{
		// 8-digit dates (YYYYMMDD)
		{
			name:          "8-digit date YYYYMMDD",
			text:          "20200315",
			expectedYear:  intPtr(2020),
			expectedMonth: intPtr(3),
			expectedDay:   intPtr(15),
			shouldParse:   true,
		},
		{
			name:          "8-digit date with leading text",
			text:          "Report from 20211026 shows",
			expectedYear:  intPtr(2021),
			expectedMonth: intPtr(10),
			expectedDay:   intPtr(26),
			shouldParse:   true,
		},
		{
			name:          "8-digit date at end",
			text:          "Due date: 20220131",
			expectedYear:  intPtr(2022),
			expectedMonth: intPtr(1),
			expectedDay:   intPtr(31),
			shouldParse:   true,
		},
		{
			name:        "Invalid 8-digit date (month 13)",
			text:        "20211332",
			shouldParse: false,
		},
		{
			name:        "Invalid 8-digit date (day 32)",
			text:        "20210132",
			shouldParse: false,
		},
		{
			name:        "Invalid 8-digit date (Feb 30)",
			text:        "20210230",
			shouldParse: false,
		},

		// 6-digit dates (YYMMDD)
		{
			name:          "6-digit date YYMMDD (2021)",
			text:          "211026",
			expectedYear:  intPtr(2021),
			expectedMonth: intPtr(10),
			expectedDay:   intPtr(26),
			shouldParse:   true,
		},
		{
			name:          "6-digit date YYMMDD (2000)",
			text:          "000101",
			expectedYear:  intPtr(2000),
			expectedMonth: intPtr(1),
			expectedDay:   intPtr(1),
			shouldParse:   true,
		},
		{
			name:          "6-digit date YYMMDD (2069)",
			text:          "691231",
			expectedYear:  intPtr(2069),
			expectedMonth: intPtr(12),
			expectedDay:   intPtr(31),
			shouldParse:   true,
		},
		{
			name:          "6-digit date YYMMDD (1970)",
			text:          "700101",
			expectedYear:  intPtr(1970),
			expectedMonth: intPtr(1),
			expectedDay:   intPtr(1),
			shouldParse:   true,
		},
		{
			name:          "6-digit date YYMMDD (1999)",
			text:          "991231",
			expectedYear:  intPtr(1999),
			expectedMonth: intPtr(12),
			expectedDay:   intPtr(31),
			shouldParse:   true,
		},

		// 14-digit datetime (YYYYMMDDHHmmss)
		{
			name:          "14-digit datetime",
			text:          "20211026141200",
			expectedYear:  intPtr(2021),
			expectedMonth: intPtr(10),
			expectedDay:   intPtr(26),
			expectedHour:  intPtr(14),
			expectedMin:   intPtr(12),
			expectedSec:   intPtr(0),
			shouldParse:   true,
		},
		{
			name:          "14-digit datetime with seconds",
			text:          "20200315143045",
			expectedYear:  intPtr(2020),
			expectedMonth: intPtr(3),
			expectedDay:   intPtr(15),
			expectedHour:  intPtr(14),
			expectedMin:   intPtr(30),
			expectedSec:   intPtr(45),
			shouldParse:   true,
		},
		{
			name:          "14-digit datetime midnight",
			text:          "20230101000000",
			expectedYear:  intPtr(2023),
			expectedMonth: intPtr(1),
			expectedDay:   intPtr(1),
			expectedHour:  intPtr(0),
			expectedMin:   intPtr(0),
			expectedSec:   intPtr(0),
			shouldParse:   true,
		},
		{
			name:          "14-digit datetime end of day",
			text:          "20231231235959",
			expectedYear:  intPtr(2023),
			expectedMonth: intPtr(12),
			expectedDay:   intPtr(31),
			expectedHour:  intPtr(23),
			expectedMin:   intPtr(59),
			expectedSec:   intPtr(59),
			shouldParse:   true,
		},

		// 12-digit datetime (YYYYMMDDHHmm)
		{
			name:          "12-digit datetime",
			text:          "202110261412",
			expectedYear:  intPtr(2021),
			expectedMonth: intPtr(10),
			expectedDay:   intPtr(26),
			expectedHour:  intPtr(14),
			expectedMin:   intPtr(12),
			shouldParse:   true,
		},

		// 4-digit time (HHmm)
		{
			name:        "4-digit time HHmm",
			text:        "1430",
			expectedHour: intPtr(14),
			expectedMin:  intPtr(30),
			shouldParse: true,
		},
		{
			name:        "4-digit time midnight",
			text:        "0000",
			expectedHour: intPtr(0),
			expectedMin:  intPtr(0),
			shouldParse: true,
		},
		{
			name:        "4-digit time end of day",
			text:        "2359",
			expectedHour: intPtr(23),
			expectedMin:  intPtr(59),
			shouldParse: true,
		},
		{
			name:        "Invalid 4-digit time (hour 24)",
			text:        "2400",
			shouldParse: false,
		},
		{
			name:        "Invalid 4-digit time (minute 60)",
			text:        "1360",
			shouldParse: false,
		},

		// 6-digit time (HHmmss)
		{
			name:        "6-digit time HHmmss",
			text:        "143045",
			expectedHour: intPtr(14),
			expectedMin:  intPtr(30),
			expectedSec:  intPtr(45),
			shouldParse: true,
		},
		{
			name:        "6-digit time with zero seconds",
			text:        "143000",
			expectedHour: intPtr(14),
			expectedMin:  intPtr(30),
			expectedSec:  intPtr(0),
			shouldParse: true,
		},
		{
			name:        "Invalid 6-digit time (second 60)",
			text:        "143060",
			shouldParse: false,
		},

		// 4-digit date (MMDD)
		{
			name:          "4-digit date MMDD",
			text:          "0315",
			expectedMonth: intPtr(3),
			expectedDay:   intPtr(15),
			shouldParse:   true,
		},
		{
			name:          "4-digit date MMDD December",
			text:          "1225",
			expectedMonth: intPtr(12),
			expectedDay:   intPtr(25),
			shouldParse:   true,
		},
		{
			name:        "4-digit ambiguous 1332 (could be time 13:32)",
			text:        "1332",
			expectedHour: intPtr(13),
			expectedMin:  intPtr(32),
			shouldParse: true,
		},
		{
			name:        "4-digit ambiguous 0132 (could be date 01-32 or time 01:32)",
			text:        "0132",
			expectedHour: intPtr(1),
			expectedMin:  intPtr(32),
			shouldParse: true,
		},

		// Edge cases
		{
			name:          "Leap year Feb 29",
			text:          "20200229",
			expectedYear:  intPtr(2020),
			expectedMonth: intPtr(2),
			expectedDay:   intPtr(29),
			shouldParse:   true,
		},
		{
			name:        "Non-leap year Feb 29",
			text:        "20210229",
			shouldParse: false,
		},
		{
			name:        "Too short (3 digits)",
			text:        "123",
			shouldParse: false,
		},
		{
			name:        "Too long (15 digits)",
			text:        "123456789012345",
			shouldParse: false,
		},

		// Word boundary tests
		{
			name:          "In sentence with word boundaries",
			text:          "The meeting on 20211026 was productive",
			expectedYear:  intPtr(2021),
			expectedMonth: intPtr(10),
			expectedDay:   intPtr(26),
			shouldParse:   true,
		},
		{
			name:        "Part of longer number (no word boundary)",
			text:        "ID123456789",
			shouldParse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			if tt.shouldParse {
				assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
				if len(results) == 0 {
					return
				}

				result := results[0]
				start := result.Start()

				if tt.expectedYear != nil {
					assert.Equal(t, *tt.expectedYear, *start.Get(kronos.ComponentYear), "Year mismatch")
				}
				if tt.expectedMonth != nil {
					assert.Equal(t, *tt.expectedMonth, *start.Get(kronos.ComponentMonth), "Month mismatch")
				}
				if tt.expectedDay != nil {
					assert.Equal(t, *tt.expectedDay, *start.Get(kronos.ComponentDay), "Day mismatch")
				}
				if tt.expectedHour != nil {
					assert.Equal(t, *tt.expectedHour, *start.Get(kronos.ComponentHour), "Hour mismatch")
				}
				if tt.expectedMin != nil {
					assert.Equal(t, *tt.expectedMin, *start.Get(kronos.ComponentMinute), "Minute mismatch")
				}
				if tt.expectedSec != nil {
					assert.Equal(t, *tt.expectedSec, *start.Get(kronos.ComponentSecond), "Second mismatch")
				}
			} else {
				assert.Empty(t, results, "Should not parse: %s", tt.text)
			}
		})
	}
}

func TestENCompactFormatParser_Integration(t *testing.T) {
	// Test with full EN configuration
	refDate := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		text        string
		expectCount int
		checkFunc   func(*testing.T, []*kronos.ParsingResult)
	}{
		{
			name:        "Multiple compact dates in text",
			text:        "Events on 20210315 and 20211026",
			expectCount: 2,
			checkFunc: func(t *testing.T, results []*kronos.ParsingResult) {
				assert.Equal(t, 2021, *results[0].Start().Get(kronos.ComponentYear))
				assert.Equal(t, 3, *results[0].Start().Get(kronos.ComponentMonth))
				assert.Equal(t, 2021, *results[1].Start().Get(kronos.ComponentYear))
				assert.Equal(t, 10, *results[1].Start().Get(kronos.ComponentMonth))
			},
		},
		{
			name:        "Compact date and time separately",
			text:        "Date 20210315 time 1430",
			expectCount: 2,
			checkFunc: func(t *testing.T, results []*kronos.ParsingResult) {
				// First result should be the date
				assert.Equal(t, 2021, *results[0].Start().Get(kronos.ComponentYear))
				assert.Equal(t, 3, *results[0].Start().Get(kronos.ComponentMonth))
				// Second result should be the time
				assert.Equal(t, 14, *results[1].Start().Get(kronos.ComponentHour))
				assert.Equal(t, 30, *results[1].Start().Get(kronos.ComponentMinute))
			},
		},
		{
			name:        "Compact datetime format",
			text:        "Timestamp: 20211026141200",
			expectCount: 1,
			checkFunc: func(t *testing.T, results []*kronos.ParsingResult) {
				assert.Equal(t, 2021, *results[0].Start().Get(kronos.ComponentYear))
				assert.Equal(t, 10, *results[0].Start().Get(kronos.ComponentMonth))
				assert.Equal(t, 26, *results[0].Start().Get(kronos.ComponentDay))
				assert.Equal(t, 14, *results[0].Start().Get(kronos.ComponentHour))
				assert.Equal(t, 12, *results[0].Start().Get(kronos.ComponentMinute))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewENCompactFormatParser()
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			assert.Len(t, results, tt.expectCount, "Result count mismatch")
			if len(results) == tt.expectCount && tt.checkFunc != nil {
				tt.checkFunc(t, results)
			}
		})
	}
}
