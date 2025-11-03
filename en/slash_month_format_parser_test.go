package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENSlashMonthFormatParser(t *testing.T) {
	parser := NewENSlashMonthFormatParser()
	refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		text          string
		expectedMonth int
		expectedYear  int
		expectedDay   int
		shouldParse   bool
	}{
		{
			text:          "11/2005",
			expectedMonth: 11,
			expectedYear:  2005,
			expectedDay:   1,
			shouldParse:   true,
		},
		{
			text:          "06/2005",
			expectedMonth: 6,
			expectedYear:  2005,
			expectedDay:   1,
			shouldParse:   true,
		},
		{
			text:          "1/2020",
			expectedMonth: 1,
			expectedYear:  2020,
			expectedDay:   1,
			shouldParse:   true,
		},
		{
			text:          "12/1999",
			expectedMonth: 12,
			expectedYear:  1999,
			expectedDay:   1,
			shouldParse:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refDate, nil)

			if tt.shouldParse {
				assert.NotEmpty(t, results, "Expected to parse: %s", tt.text)
				if len(results) == 0 {
					return
				}

				result := results[0]
				assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
				assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
				assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")
				assert.False(t, result.Start().IsCertain(kronos.ComponentDay), "Day should be implied")
			} else {
				assert.Empty(t, results, "Should not parse: %s", tt.text)
			}
		})
	}
}
