// Package en provides English language support for Chrono.
// It includes parsers, refiners, and configurations for parsing English dates and times.
package en

import (
	"fmt"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/experimental"
)

// New creates a new ParserBuilder using the casual English configuration.
// This is the recommended way to parse English dates with the builder pattern.
//
// Example:
//
//	parser := en.New().
//	    WithReferenceDate(time.Now()).
//	    PreferPast()
//	results, err := parser.Parse("last Monday")
func New() *kronos.ParserBuilder {
	return kronos.New(experimental.EnglishCasualChrono())
}

// NewStrict creates a new ParserBuilder using strict English configuration.
// In strict mode, only formal date/time patterns are recognized.
//
// Example:
//
//	parser := en.NewStrict().WithReferenceDate(time.Now())
//	results, err := parser.Parse("2020-03-15")
func NewStrict() *kronos.ParserBuilder {
	return kronos.New(experimental.EnglishStrictChrono())
}

// NewGB creates a new ParserBuilder using UK-style English configuration.
// It uses little-endian date format (day/month/year) and casual expressions.
//
// Example:
//
//	parser := en.NewGB().WithReferenceDate(time.Now())
//	results, err := parser.Parse("15/03/2020")
func NewGB() *kronos.ParserBuilder {
	return kronos.New(experimental.EnglishGBChrono()).DateOrder(kronos.DateOrderDMY)
}

// ParseSimple is a convenience function that parses text using the new builder API.
// It uses casual English configuration and the current time as reference.
//
// Example:
//
//	results, err := en.ParseSimple("tomorrow at 3pm")
func ParseSimple(text string) ([]kronos.Result, error) {
	results, err := New().Parse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse text: %w", err)
	}
	return results, nil
}

// ParseDateSimple is a convenience function that parses text and returns the first date.
// It uses casual English configuration and the current time as reference.
//
// Example:
//
//	date, err := en.ParseDateSimple("tomorrow at 3pm")
func ParseDateSimple(text string) (*time.Time, error) {
	date, err := New().ParseDate(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	return date, nil
}
