package kronos

import (
	"testing"
	"time"

	"github.com/kljensen/kronos/internal/chrono"
	"github.com/stretchr/testify/assert"
)

func TestParserBuilder_Basic(t *testing.T) {
	// Create a simple configuration for testing
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	// Test New
	builder := New(c)
	assert.NotNil(t, builder)
	assert.NotNil(t, builder.config)
	assert.NotNil(t, builder.settings)
}

func TestParserBuilder_Chainable(t *testing.T) {
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	// Test that methods are chainable
	refDate := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)
	builder := New(c).
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
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	builder := New(c).DateOrder(DateOrderDMY)
	assert.Equal(t, DateOrderDMY, builder.settings.DateOrder)

	builder = New(c).DateOrder(DateOrderMDY)
	assert.Equal(t, DateOrderMDY, builder.settings.DateOrder)

	builder = New(c).DateOrder(DateOrderYMD)
	assert.Equal(t, DateOrderYMD, builder.settings.DateOrder)
}

func TestParserBuilder_PreferenceChaining(t *testing.T) {
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	// Test PreferPast
	builder := New(c).PreferPast()
	assert.Equal(t, PreferPast, builder.settings.PreferDatesFrom)

	// Test PreferFuture
	builder = New(c).PreferFuture()
	assert.Equal(t, PreferFuture, builder.settings.PreferDatesFrom)

	// Test PreferCurrentPeriod
	builder = New(c).PreferCurrentPeriod()
	assert.Equal(t, PreferCurrentPeriod, builder.settings.PreferDatesFrom)
}

func TestParserBuilder_Timezone(t *testing.T) {
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	builder := New(c).
		Timezone("America/New_York")

	assert.Equal(t, "America/New_York", builder.settings.Timezone)
}

func TestParserBuilder_StrictCasual(t *testing.T) {
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	// Default is casual
	builder := New(c)
	assert.False(t, builder.settings.StrictParsing)

	// Test Strict
	builder = New(c).Strict()
	assert.True(t, builder.settings.StrictParsing)

	// Test Casual
	builder = New(c).Casual()
	assert.False(t, builder.settings.StrictParsing)
}

func TestParserBuilder_ForwardDate(t *testing.T) {
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	builder := New(c).PreferFuture()
	assert.Equal(t, PreferFuture, builder.settings.PreferDatesFrom)
}

func TestPackageLevelParse(t *testing.T) {
	config := &chrono.Configuration{
		Parsers:  []any{},
		Refiners: []any{},
	}
	c := chrono.NewChrono(config)

	// Test Parse via builder - should not panic even with no parsers
	results, err := New(c).Parse("test")
	assert.NoError(t, err)
	assert.NotNil(t, results)

	// Test ParseDate via builder - should not panic even with no parsers
	date, err := New(c).ParseDate("test")
	assert.NoError(t, err)
	assert.Nil(t, date) // No parsers, no results
}
