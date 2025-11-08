// Package parser defines the interfaces for implementing custom parsers and refiners in Kronos.
// This package is designed to enable support for multiple languages and custom date formats.
//
// To add support for a new language:
//  1. Implement the Parser interface for each date format you want to recognize
//  2. Implement the Refiner interface for any post-processing logic
//  3. Create a language package (like en, fr, es) that registers your parsers
//
// Example:
//
//	type MyCustomParser struct{}
//
//	func (p *MyCustomParser) Pattern(ctx parser.Context) *regexp.Regexp {
//	    return regexp.MustCompile(`my pattern`)
//	}
//
//	func (p *MyCustomParser) Extract(ctx parser.Context, match []string) any {
//	    // Extract and return components
//	    return ctx.CreateParsingComponents(nil)
//	}
package parser

import (
	"regexp"
	"time"
)

// Parser is the interface for custom date format parsers.
// Implement this interface to add support for new date formats or languages.
//
// Each parser recognizes a specific pattern (via Pattern method) and extracts
// date/time components from matches (via Extract method).
//
// The Extract method can return:
//   - Components (parsed date/time components)
//   - Result (a complete parsing result)
//   - A map of components
//   - nil if extraction fails
type Parser interface {
	// Pattern returns the regular expression pattern for this parser.
	// The pattern is used to find potential matches in the input text.
	// The context provides access to parsing options and reference date.
	Pattern(context Context) *regexp.Regexp

	// Extract is called when the pattern matches.
	// It should analyze the match and return parsed components, a result,
	// a component map, or nil if extraction fails.
	//
	// The match parameter contains the regex match groups, where match[0]
	// is the full matched text and match[1], match[2], etc. are capture groups.
	Extract(context Context, match []string) any
}

// Refiner is the interface for post-processing parsing results.
// Implement this interface to add logic that combines, filters, or enhances results.
//
// Examples of refiners:
//   - Merging separate date and time results into a single datetime
//   - Resolving ambiguous dates based on preferences (past/future)
//   - Filtering unlikely or invalid date formats
//   - Extracting timezone information
type Refiner interface {
	// Refine processes a list of parsing results and returns a refined list.
	// Refiners can:
	//   - Combine multiple results (e.g., merge "Monday" + "3pm" into "Monday 3pm")
	//   - Filter results (e.g., remove unlikely formats)
	//   - Enhance results (e.g., add timezone information)
	//   - Reorder results
	//
	// The context provides access to the original text and parsing options.
	Refine(context Context, results []Result) []Result
}

// Context provides the parsing context including text, options, and reference date.
// This interface is implemented by the internal parsing context and provides
// methods that parsers and refiners need to access.
type Context interface {
	// Text returns the input text being parsed.
	Text() string

	// Reference returns the reference date/time with timezone information.
	// This is used as the basis for relative dates like "tomorrow" or "next week".
	Reference() Reference

	// RefDate returns the reference date as a time.Time.
	// This is a convenience method equivalent to Reference().Instant().
	RefDate() time.Time

	// Option returns the parsing options (forward date, preferences, etc.).
	Option() Option

	// Settings returns the full parsing settings, if available.
	// This may be nil if only basic options were provided.
	Settings() Settings

	// CreateParsingComponents creates new date/time components.
	// Pass nil to create empty components, or a map of component values.
	// Returns components that implement the kronos.Components interface.
	CreateParsingComponents(components any) any

	// CreateParsingResult creates a new parsing result.
	// Supports various signatures:
	//   - (index, text)
	//   - (index, endIndex)
	//   - (index, text, startComponents)
	//   - (index, text, startComponents, endComponents)
	// Returns a Result implementation.
	CreateParsingResult(index int, textOrEndIndex any, args ...any) Result
}

// Reference represents a reference date/time with optional timezone.
// This is used as the basis for parsing relative dates.
type Reference interface {
	// Instant returns the reference date/time.
	Instant() time.Time

	// TimezoneOffset returns the timezone offset in minutes, if known.
	// Returns nil if no timezone information is available.
	TimezoneOffset() *int

	// GetDateWithAdjustedTimezone returns the reference date adjusted for timezone.
	GetDateWithAdjustedTimezone() time.Time
}

// Result represents a parsed date/time result.
// This interface provides access to the matched text, position, and parsed components.
type Result interface {
	// Index returns the position in the input text where this result was found.
	Index() int

	// Text returns the matched text from the input.
	Text() string

	// Start returns the starting date/time components.
	// Returns an object implementing kronos.Components.
	Start() any

	// End returns the ending date/time components for a range, or nil for a single date/time.
	// Returns an object implementing kronos.Components, or nil.
	End() any

	// Date returns a time.Time object created from the start components.
	Date() time.Time

	// Clone creates a deep copy of the result.
	Clone() Result

	// SetIndex sets the position in the input text (used internally by parsers).
	SetIndex(index int)

	// SetStart sets the start component (used by parsers and refiners).
	// The start parameter should be an object implementing kronos.Components.
	SetStart(start any)

	// AddTag adds a debugging tag to the result.
	AddTag(tag string) Result

	// Tags returns debugging tags for this result.
	Tags() map[string]bool

	// RefDate returns the reference date used for parsing.
	RefDate() time.Time
}

// Option represents parsing options.
// This interface provides access to configuration like forward date mode,
// date preferences, and timezone maps.
type Option interface {
	// ForwardDate returns true if only forward dates should be parsed.
	ForwardDate() bool

	// Preference returns the date preference (past, future, or neither).
	Preference() int

	// DateOrder returns the date order (MDY, DMY, or YMD).
	DateOrder() int

	// Timezones returns the timezone abbreviation map.
	Timezones() map[string]any

	// Debug returns the debug handler, if any.
	Debug() func(string)
}

// Settings represents comprehensive parsing settings.
// This is a more complete version of Option that includes all configuration.
// It may be nil if only basic options were provided.
type Settings interface {
	Option
	// Additional methods can be added here as needed
}
