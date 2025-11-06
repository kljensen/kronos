// Package en provides English language support for Chrono.
// It includes parsers, refiners, and configurations for parsing English dates and times.
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"fmt"
	"time"

	"github.com/kljensen/kronos"
)

var (
	// Casual is a Chrono instance configured for parsing casual English.
	// It recognizes informal expressions like "today", "tomorrow", "next week", etc.
	//
	// Deprecated: Use New() to create a parser instance with the builder pattern API.
	// Migration example:
	//   Before: results := en.Casual.Parse("tomorrow", now, nil)
	//   After:  results, _ := en.New().Parse("tomorrow")
	Casual = kronos.NewChrono(CreateCasualConfiguration(false))

	// Strict is a Chrono instance configured for parsing strict English.
	// It only recognizes formal date/time patterns and avoids casual expressions.
	//
	// Deprecated: Use NewStrict() to create a parser instance with the builder pattern API.
	// Migration example:
	//   Before: results := en.Strict.Parse("2020-03-15", now, nil)
	//   After:  results, _ := en.NewStrict().Parse("2020-03-15")
	Strict = kronos.NewChrono(CreateConfiguration(true, false))

	// GB is a Chrono instance configured for parsing UK-style English.
	// It uses little-endian date format (day/month/year) and casual expressions.
	//
	// Deprecated: Use NewGB() to create a parser instance with the builder pattern API.
	// Migration example:
	//   Before: results := en.GB.Parse("15/03/2020", now, nil)
	//   After:  results, _ := en.NewGB().Parse("15/03/2020")
	GB = kronos.NewChrono(CreateCasualConfiguration(true))
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
	return kronos.New(Casual)
}

// NewStrict creates a new ParserBuilder using strict English configuration.
// In strict mode, only formal date/time patterns are recognized.
//
// Example:
//
//	parser := en.NewStrict().WithReferenceDate(time.Now())
//	results, err := parser.Parse("2020-03-15")
func NewStrict() *kronos.ParserBuilder {
	return kronos.New(Strict)
}

// NewGB creates a new ParserBuilder using UK-style English configuration.
// It uses little-endian date format (day/month/year) and casual expressions.
//
// Example:
//
//	parser := en.NewGB().WithReferenceDate(time.Now())
//	results, err := parser.Parse("15/03/2020")
func NewGB() *kronos.ParserBuilder {
	return kronos.New(GB).DateOrder(kronos.DateOrderDMY)
}

// StrictParser is deprecated: use NewStrict() instead.
// This function is maintained for backward compatibility.
func StrictParser() *kronos.ParserBuilder {
	return NewStrict()
}

// GBParser is deprecated: use NewGB() instead.
// This function is maintained for backward compatibility.
func GBParser() *kronos.ParserBuilder {
	return NewGB()
}

// Parse parses the text and returns all parsed results.
// It uses the casual English configuration.
// This is the legacy API - for the new builder pattern, use New() instead.
func Parse(text string, ref time.Time, option *kronos.ParsingOption) []*kronos.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate parses the text and returns the first parsed date.
// It uses the casual English configuration.
// Returns nil if no date is found.
// This is the legacy API - for the new builder pattern, use New().ParseDate() instead.
func ParseDate(text string, ref time.Time, option *kronos.ParsingOption) *time.Time {
	return Casual.ParseDate(text, ref, option)
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
