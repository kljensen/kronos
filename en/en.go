// Package en provides English language support for Chrono.
// It includes parsers, refiners, and configurations for parsing English dates and times.
package en

import (
	"time"

	"github.com/kljensen/kronos"
)

var (
	// Casual is a Chrono instance configured for parsing casual English.
	// It recognizes informal expressions like "today", "tomorrow", "next week", etc.
	Casual = kronos.NewChrono(CreateCasualConfiguration(false))

	// Strict is a Chrono instance configured for parsing strict English.
	// It only recognizes formal date/time patterns and avoids casual expressions.
	// Deprecated: Use StrictParser() for the builder pattern API.
	Strict = kronos.NewChrono(CreateConfiguration(true, false))

	// GB is a Chrono instance configured for parsing UK-style English.
	// It uses little-endian date format (day/month/year) and casual expressions.
	// Deprecated: Use GBParser() for the builder pattern API.
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

// StrictParser creates a new ParserBuilder using strict English configuration.
// In strict mode, only formal date/time patterns are recognized.
//
// Example:
//
//	parser := en.StrictParser().WithReferenceDate(time.Now())
//	results, err := parser.Parse("2020-03-15")
func StrictParser() *kronos.ParserBuilder {
	return kronos.New(Strict)
}

// GBParser creates a new ParserBuilder using UK-style English configuration.
// It uses little-endian date format (day/month/year) and casual expressions.
//
// Example:
//
//	parser := en.GBParser().WithReferenceDate(time.Now())
//	results, err := parser.Parse("15/03/2020")
func GBParser() *kronos.ParserBuilder {
	return kronos.New(GB)
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
	return New().Parse(text)
}

// ParseDateSimple is a convenience function that parses text and returns the first date.
// It uses casual English configuration and the current time as reference.
//
// Example:
//
//	date, err := en.ParseDateSimple("tomorrow at 3pm")
func ParseDateSimple(text string) (*time.Time, error) {
	return New().ParseDate(text)
}
