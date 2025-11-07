package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParserBuilder_Basic(t *testing.T) {
	// Create a simple configuration for testing
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	// Test New
	builder := New(chrono)
	assert.NotNil(t, builder)
	assert.NotNil(t, builder.config)
	assert.NotNil(t, builder.settings)
}

func TestParserBuilder_Chainable(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	// Test that methods are chainable
	refDate := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)
	builder := New(chrono).
		WithReferenceDate(refDate).
		Strict().
		PreferPast().
		Timezone("UTC")

	assert.NotNil(t, builder)
	assert.Equal(t, refDate, builder.refDate)
	assert.True(t, builder.settings.StrictParsing)
	assert.Equal(t, PreferPast, builder.settings.PreferDatesFrom)
	assert.Equal(t, "UTC", builder.settings.Timezone)
}

func TestParserBuilder_DateOrder(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	builder := New(chrono).DateOrder(DateOrderDMY)
	assert.Equal(t, DateOrderDMY, builder.settings.DateOrder)

	builder = New(chrono).DateOrder(DateOrderMDY)
	assert.Equal(t, DateOrderMDY, builder.settings.DateOrder)

	builder = New(chrono).DateOrder(DateOrderYMD)
	assert.Equal(t, DateOrderYMD, builder.settings.DateOrder)
}

func TestParserBuilder_PreferenceChaining(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	// Test PreferPast
	builder := New(chrono).PreferPast()
	assert.Equal(t, PreferPast, builder.settings.PreferDatesFrom)

	// Test PreferFuture
	builder = New(chrono).PreferFuture()
	assert.Equal(t, PreferFuture, builder.settings.PreferDatesFrom)

	// Test PreferCurrentPeriod
	builder = New(chrono).PreferCurrentPeriod()
	assert.Equal(t, PreferCurrentPeriod, builder.settings.PreferDatesFrom)
}

func TestParserBuilder_Timezone(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	builder := New(chrono).
		Timezone("America/New_York")

	assert.Equal(t, "America/New_York", builder.settings.Timezone)
}

func TestParserBuilder_StrictCasual(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	// Default is casual
	builder := New(chrono)
	assert.False(t, builder.settings.StrictParsing)

	// Test Strict
	builder = New(chrono).Strict()
	assert.True(t, builder.settings.StrictParsing)

	// Test Casual
	builder = New(chrono).Casual()
	assert.False(t, builder.settings.StrictParsing)
}

func TestParserBuilder_ForwardDate(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	builder := New(chrono).PreferFuture()
	assert.Equal(t, PreferFuture, builder.settings.PreferDatesFrom)
}

func TestPackageLevelParse(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	chrono := NewChrono(config)

	// Test Parse via builder - should not panic even with no parsers
	results, err := New(chrono).Parse("test")
	assert.NoError(t, err)
	assert.NotNil(t, results)

	// Test ParseDate via builder - should not panic even with no parsers
	date, err := New(chrono).ParseDate("test")
	assert.NoError(t, err)
	assert.Nil(t, date) // No parsers, no results
}
