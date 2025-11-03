package en

import (
	"testing"
	"time"

	"github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
)

func TestNew_BasicUsage(t *testing.T) {
	// Test that New() returns a builder
	builder := New()
	assert.NotNil(t, builder)

	// Test simple parsing
	results, err := builder.Parse("tomorrow")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
}

func TestNew_WithReferenceDate(t *testing.T) {
	refDate := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)

	results, err := New().
		WithReferenceDate(refDate).
		Parse("tomorrow")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	// Tomorrow from March 15 should be March 16
	date := results[0].Date()
	assert.Equal(t, 2020, date.Year())
	assert.Equal(t, time.March, date.Month())
	assert.Equal(t, 16, date.Day())
}

func TestNew_Chainable(t *testing.T) {
	refDate := time.Date(2020, 11, 15, 12, 0, 0, 0, time.UTC)

	// Test that methods are chainable
	builder := New().
		WithReferenceDate(refDate).
		PreferPast().
		Timezone("UTC")

	assert.NotNil(t, builder)

	// Parse "March" which is ambiguous (could be past or future)
	results, err := builder.Parse("March")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	// With PreferPast in November, "March" should be March 2020 (past)
	date := results[0].Date()
	assert.Equal(t, 2020, date.Year())
	assert.Equal(t, time.March, date.Month())
}

func TestNew_PreferFuture(t *testing.T) {
	refDate := time.Date(2020, 11, 15, 12, 0, 0, 0, time.UTC)

	// Test PreferFuture
	results, err := New().
		WithReferenceDate(refDate).
		PreferFuture().
		Parse("March")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	// With PreferFuture in November, "March" should be March 2021 (future)
	date := results[0].Date()
	assert.Equal(t, 2021, date.Year())
	assert.Equal(t, time.March, date.Month())
}

func TestStrictParser(t *testing.T) {
	// Test that StrictParser() returns a builder
	builder := StrictParser()
	assert.NotNil(t, builder)

	// Strict mode should still parse formal dates
	results, err := builder.Parse("2020-03-15")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
}

func TestGBParser(t *testing.T) {
	// Test that GBParser() returns a builder
	builder := GBParser()
	assert.NotNil(t, builder)

	// GB format should parse day/month/year
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	results, err := builder.
		WithReferenceDate(refDate).
		Parse("15/03/2020")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	// Should be March 15, 2020
	date := results[0].Date()
	assert.Equal(t, 2020, date.Year())
	assert.Equal(t, time.March, date.Month())
	assert.Equal(t, 15, date.Day())
}

func TestParseSimple(t *testing.T) {
	// Test ParseSimple convenience function
	results, err := ParseSimple("tomorrow")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
}

func TestParseDateSimple(t *testing.T) {
	// Test ParseDateSimple convenience function
	date, err := ParseDateSimple("tomorrow")
	assert.NoError(t, err)
	assert.NotNil(t, date)

	// Tomorrow should be in the future
	assert.True(t, date.After(time.Now()))
}

func TestParseDateSimple_NoResults(t *testing.T) {
	// Test ParseDateSimple with no results
	date, err := ParseDateSimple("not a date at all xyz")
	assert.NoError(t, err)
	assert.Nil(t, date)
}

func TestNew_DateOrder(t *testing.T) {
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

	// Test MDY (US) format
	results, err := New().
		WithReferenceDate(refDate).
		DateOrder(kronos.DateOrderMDY).
		Parse("03/15/2020")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	date := results[0].Date()
	assert.Equal(t, time.March, date.Month())
	assert.Equal(t, 15, date.Day())
}

func TestBuilder_CompleteExample(t *testing.T) {
	// A complete example showing the builder pattern
	refDate := time.Date(2020, 11, 15, 12, 0, 0, 0, time.UTC)

	date, err := New().
		WithReferenceDate(refDate).
		PreferPast().
		Timezone("UTC").
		ParseDate("last Monday")

	assert.NoError(t, err)
	assert.NotNil(t, date)

	// Last Monday from November 15, 2020 (Sunday) should be November 9
	assert.Equal(t, 2020, date.Year())
	assert.Equal(t, time.November, date.Month())
	assert.Equal(t, 9, date.Day())
}
