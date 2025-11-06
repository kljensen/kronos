//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

// TestMergingRelativeDates tests merging of relative date expressions
// like "2 weeks after yesterday" or "2 days after next Friday".
func TestMergingRelativeDates(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		refDate     time.Time
		expectedTxt string
		expectYear  int
		expectMonth int
		expectDay   int
		expectWday  *int
		certainDay  bool
		certainMon  bool
		certainYear bool
		expectHour  int
	}{
		{
			name:        "2 weeks after yesterday",
			text:        "2 weeks after yesterday",
			refDate:     time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			expectedTxt: "2 weeks after yesterday",
			expectYear:  2022,
			expectMonth: 2,
			expectDay:   15,
			expectWday:  ptr(2), // Tuesday
			certainDay:  true,
			certainMon:  true,
			certainYear: true,
			expectHour:  0,
		},
		{
			name:        "2 months before 02/02",
			text:        "2 months before 02/02",
			refDate:     time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			expectedTxt: "2 months before 02/02",
			expectYear:  2021,
			expectMonth: 12,
			expectDay:   2,
			expectWday:  nil,
			certainDay:  false,
			certainMon:  true,
			certainYear: true,
			expectHour:  12,
		},
		{
			name:        "2 days after next Friday",
			text:        "2 days after next Friday",
			refDate:     time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			expectedTxt: "2 days after next Friday",
			expectYear:  2022,
			expectMonth: 2,
			expectDay:   13,
			expectWday:  nil,
			certainDay:  true,
			certainMon:  true,
			certainYear: true,
			expectHour:  12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasualConfiguration(true)
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.Len(t, results, 1, "Expected exactly 1 result")
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.expectedTxt, result.Text(), "Text mismatch")

			start := result.Start()
			assert.NotNil(t, start, "Start should not be nil")
			if start == nil {
				return
			}

			// Check date components
			assert.Equal(t, tt.expectYear, *start.Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectMonth, *start.Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectDay, *start.Get(kronos.ComponentDay), "Day mismatch")

			if tt.expectWday != nil {
				assert.Equal(t, *tt.expectWday, *start.Get(kronos.ComponentWeekday), "Weekday mismatch")
			}

			// Check certainty
			assert.Equal(t, tt.certainDay, start.IsCertain(kronos.ComponentDay), "Day certainty mismatch")
			assert.Equal(t, tt.certainMon, start.IsCertain(kronos.ComponentMonth), "Month certainty mismatch")
			assert.Equal(t, tt.certainYear, start.IsCertain(kronos.ComponentYear), "Year certainty mismatch")

			// Check time
			assert.Equal(t, tt.expectHour, *start.Get(kronos.ComponentHour), "Hour mismatch")
		})
	}
}

// ptr is a helper function to get a pointer to an int
func ptr(i int) *int {
	return &i
}
