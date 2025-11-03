// Package en provides English language support for Chrono.
// It includes parsers, refiners, and configurations for parsing English dates and times.
package en

import (
	"time"

	. "github.com/kljensen/kronos"
)

var (
	// Casual is a Chrono instance configured for parsing casual English.
	// It recognizes informal expressions like "today", "tomorrow", "next week", etc.
	Casual = NewChrono(CreateCasualConfiguration(false))

	// Strict is a Chrono instance configured for parsing strict English.
	// It only recognizes formal date/time patterns and avoids casual expressions.
	Strict = NewChrono(CreateConfiguration(true, false))

	// GB is a Chrono instance configured for parsing UK-style English.
	// It uses little-endian date format (day/month/year) and casual expressions.
	GB = NewChrono(CreateCasualConfiguration(true))
)

// Parse parses the text and returns all parsed results.
// It uses the casual English configuration.
func Parse(text string, ref time.Time, option *ParsingOption) []*ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate parses the text and returns the first parsed date.
// It uses the casual English configuration.
// Returns zero time if no date is found.
func ParseDate(text string, ref time.Time, option *ParsingOption) *time.Time {
	return Casual.ParseDate(text, ref, option)
}
