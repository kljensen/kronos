package en

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestENCasualDateParser_Now(t *testing.T) {
	parser := NewENCasualDateParser()
	refDate := time.Date(2012, 8, 10, 8, 9, 10, 11000000, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("The Deadline is now", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "now", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentHour))
	assert.Equal(t, 9, *result.Start().Get(kronos.ComponentMinute))
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentSecond))
}

func TestENCasualDateParser_Today(t *testing.T) {
	parser := NewENCasualDateParser()
	refDate := time.Date(2012, 8, 10, 14, 12, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("The Deadline is today", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "today", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
}

func TestENCasualDateParser_Tomorrow(t *testing.T) {
	parser := NewENCasualDateParser()

	tests := []struct {
		name        string
		text        string
		refDate     time.Time
		expectedDay int
	}{
		{
			name:        "Tomorrow afternoon",
			text:        "The Deadline is Tomorrow",
			refDate:     time.Date(2012, 8, 10, 17, 10, 0, 0, time.UTC),
			expectedDay: 11,
		},
		{
			name:        "Tomorrow early morning",
			text:        "The Deadline is Tomorrow",
			refDate:     time.Date(2012, 8, 10, 1, 0, 0, 0, time.UTC),
			expectedDay: 11,
		},
		{
			name:        "tmr",
			text:        "See you tmr",
			refDate:     time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC),
			expectedDay: 11,
		},
		{
			name:        "tmrw",
			text:        "See you tmrw",
			refDate:     time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC),
			expectedDay: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
			chrono := kronos.NewChrono(config)
			results := chrono.Parse(tt.text, tt.refDate, nil)

			assert.NotEmpty(t, results)
			result := results[0]

			assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
			assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
			assert.Equal(t, tt.expectedDay, *result.Start().Get(kronos.ComponentDay))
		})
	}
}

func TestENCasualDateParser_Yesterday(t *testing.T) {
	parser := NewENCasualDateParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("The Deadline was yesterday", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "yesterday", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 9, *result.Start().Get(kronos.ComponentDay))
}

func TestENCasualDateParser_LastNight(t *testing.T) {
	parser := NewENCasualDateParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("The Deadline was last night ", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "last night", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 9, *result.Start().Get(kronos.ComponentDay))
	assert.Equal(t, 0, *result.Start().Get(kronos.ComponentHour))
}

func TestENCasualDateParser_Tonight(t *testing.T) {
	parser := NewENCasualDateParser()
	refDate := time.Date(2012, 8, 10, 18, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("See you tonight", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "tonight", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
	// Tonight implies evening hour (22:00)
	assert.Equal(t, 22, *result.Start().Get(kronos.ComponentHour))
}

func TestENCasualDateParser_Overmorrow(t *testing.T) {
	parser := NewENCasualDateParser()
	refDate := time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC)

	config := &kronos.Configuration{Parsers: []kronos.Parser{parser}}
	chrono := kronos.NewChrono(config)
	results := chrono.Parse("See you overmorrow", refDate, nil)

	assert.NotEmpty(t, results)
	result := results[0]

	assert.Equal(t, "overmorrow", result.Text())
	assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
	assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
	assert.Equal(t, 12, *result.Start().Get(kronos.ComponentDay)) // 2 days after ref date
}
