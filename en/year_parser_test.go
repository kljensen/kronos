package en

import (
	"testing"
	"time"

	"github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENYearParser(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	parser := NewENYearParser()

	tests := []struct {
		input         string
		expectedYear  int
		expectedMonth int
		expectedDay   int
		shouldMatch   bool
	}{
		{"2020", 2020, 1, 1, true},
		{"1999", 1999, 1, 1, true},
		{"2025", 2025, 1, 1, true},
		{"1900", 1900, 1, 1, true},
		{"2999", 2999, 1, 1, true},
		{"999", 0, 0, 0, false},   // 3-digit year should not match
		{"20", 0, 0, 0, false},    // 2-digit year should not match
		{"12345", 0, 0, 0, false}, // 5-digit should not match
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.input, refTime, nil)

			if !tt.shouldMatch {
				assert.Empty(t, results, "Expected no match for %q", tt.input)
				return
			}

			assert.NotEmpty(t, results, "Expected match for %q", tt.input)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.expectedYear, *result.Start().Get(kronos.ComponentYear), "Year mismatch")
			assert.Equal(t, tt.expectedMonth, *result.Start().Get(kronos.ComponentMonth), "Month mismatch")
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay), "Day mismatch")

			// Check period
			components, ok := kronos.asParsingComponents(result.Start())
			assert.True(t, ok, "Expected ParsingComponents")
			period := components.Period()
			assert.Equal(t, kronos.PeriodYear, period, "Expected PeriodYear")
		})
	}
}

func TestENYearParserWithContext(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	parser := NewENYearParser()

	tests := []struct {
		text string
		year int
	}{
		{"The year 2020 was interesting", 2020},
		{"In 1999 we had Y2K fears", 1999},
		{"Looking forward to 2025!", 2025},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, refTime, nil)

			assert.NotEmpty(t, results, "Expected match for %q", tt.text)
			if len(results) == 0 {
				return
			}

			result := results[0]
			assert.Equal(t, tt.year, *result.Start().Get(kronos.ComponentYear), "Year mismatch for %q", tt.text)
		})
	}
}
