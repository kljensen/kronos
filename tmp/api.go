// Package kronos provides natural language date/time parsing for Go.
// It is a port of the JavaScript library Chrono (https://github.com/wanasit/chrono).
//
// Basic usage:
//
//	import "github.com/markusmobius/go-chrono"
//	import "time"
//
//	date := kronos.ParseDate("tomorrow at 3pm", time.Now(), kronos.ParsingOption{})
//	results := kronos.Parse("Monday to Friday", time.Now(), kronos.ParsingOption{})
//
// The package supports both casual and strict parsing modes through
// the Casual and Strict configurations.
package kronos

import (
	"time"

	"github.com/markusmobius/go-chrono/en"
)

var (
	// Casual is a Chrono instance configured for parsing casual English.
	// This is an alias for en.Casual.
	Casual = en.Casual

	// Strict is a Chrono instance configured for parsing strict English.
	// This is an alias for en.Strict.
	Strict = en.Strict
)

// Parse parses the text using casual English configuration and returns all parsed results.
// This is a convenience function that calls en.Parse.
//
// Parameters:
//   - text: The text to parse
//   - ref: Reference date/time for relative dates (e.g., "tomorrow" is relative to this)
//   - option: Parsing options (can be empty struct for defaults)
//
// Example:
//
//	results := kronos.Parse("Meet me tomorrow at 3pm", time.Now(), kronos.ParsingOption{})
//	for _, result := range results {
//	    fmt.Println(result.Start.Date())
//	}
func Parse(text string, ref time.Time, option ParsingOption) []*ParsingResult {
	return en.Parse(text, ref, option)
}

// ParseDate parses the text using casual English configuration and returns the first parsed date.
// Returns zero time if no date is found.
// This is a convenience function that calls en.ParseDate.
//
// Parameters:
//   - text: The text to parse
//   - ref: Reference date/time for relative dates
//   - option: Parsing options (can be empty struct for defaults)
//
// Example:
//
//	date := kronos.ParseDate("tomorrow at 3pm", time.Now(), kronos.ParsingOption{})
//	if !date.IsZero() {
//	    fmt.Println("Meeting at:", date)
//	}
func ParseDate(text string, ref time.Time, option ParsingOption) time.Time {
	return en.ParseDate(text, ref, option)
}

// New creates a new Chrono instance with the given configuration.
// If you need custom parsing behavior, create a Configuration with
// your desired parsers and refiners, then pass it to this function.
//
// Example:
//
//	config := &kronos.Configuration{
//	    Parsers: []kronos.Parser{...},
//	    Refiners: []kronos.Refiner{...},
//	}
//	chrono := kronos.New(config)
//	results := chrono.Parse("some text", time.Now(), kronos.ParsingOption{})
func New(config *Configuration) *Chrono {
	return NewChrono(config)
}
