package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENCasualTimeParser_ThisMorning(t *testing.T) {
	parser := NewENCasualTimeParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("The Deadline was this morning ", refDate, nil)
	

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "this morning", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
	assert.Equal(t, 6, *result.Start().Get(kronos.ComponentHour))
}

func TestENCasualTimeParser_ThisAfternoon(t *testing.T) {
	parser := NewENCasualTimeParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("The Deadline is this afternoon ", refDate, nil)
	

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "this afternoon", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
	assert.Equal(t, 15, *result.Start().Get(kronos.ComponentHour))
}

func TestENCasualTimeParser_ThisEvening(t *testing.T) {
	parser := NewENCasualTimeParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("The Deadline is this evening ", refDate, nil)
	

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "this evening", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
	assert.Equal(t, 20, *result.Start().Get(kronos.ComponentHour))
}

func TestENCasualTimeParser_Night(t *testing.T) {
	parser := NewENCasualTimeParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("See you night", refDate, nil)
	

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "night", result.Text())
	assert.Equal(t, 20, *result.Start().Get(kronos.ComponentHour))
}

func TestENCasualTimeParser_Noon(t *testing.T) {
	parser := NewENCasualTimeParser()
	refDate := time.Date(2012, 8, 10, 8, 0, 0, 0, time.UTC)

	tests := []string{"noon", "midday"}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("Meet at "+text, refDate, nil)
			

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, text, result.Text())
			assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, 12, *result.Start().Get(kronos.ComponentHour))
			assert.Equal(t, 0, *result.Start().Get(kronos.ComponentMinute))
		})
	}
}

func TestENCasualTimeParser_Midnight(t *testing.T) {
	parser := NewENCasualTimeParser()

	tests := []struct {
		name        string
		refDate     time.Time
		expectedDay int
	}{
		{
			name:        "Midnight in afternoon",
			refDate:     time.Date(2012, 8, 10, 14, 0, 0, 0, time.UTC),
			expectedDay: 11, // Next day's midnight
		},
		{
			name:        "Midnight in early morning",
			refDate:     time.Date(2012, 8, 10, 1, 0, 0, 0, time.UTC),
			expectedDay: 10, // Same day's midnight
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("Meet at midnight", tt.refDate, nil)
			

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, "midnight", result.Text())
			assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay))
			assert.Equal(t, 0, *result.Start().Get(kronos.ComponentHour))
			assert.Equal(t, 0, *result.Start().Get(kronos.ComponentMinute))
		})
	}
}

func TestENCasualTimeParser_WithoutThis(t *testing.T) {
	parser := NewENCasualTimeParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		text         string
		expectedHour int
	}{
		{"See you morning", 6},
		{"See you afternoon", 15},
		{"See you evening", 20},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse(tt.text, refDate, nil)
			

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, tt.expectedHour, *result.Start().Get(kronos.ComponentHour))
		})
	}
}
