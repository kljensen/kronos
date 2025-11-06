// Package parsing contains internal parsing result types and structures.
// These types should not be used directly by external code. Instead, use the
// public Result and Components interfaces from the main kronos package.
//
// Note: This package cannot import the main kronos package to avoid import cycles.
// Therefore, it uses basic types (string, int) where the main package uses custom types
// (Component, Period). The type system ensures correctness via the type aliases in the
// main package.
package parsing

import (
	"time"
)

// Component represents a date/time component (mirrored from main package to avoid cycles).
type Component string

// Period represents the granularity of a parsed date/time (mirrored from main package).
type Period int

// ParsingOption holds parsing configuration (forward declaration to avoid cycles).
// The actual definition is in the main package.
type ParsingOption struct {
	ForwardDate bool
	Preference  int // DatePreference
	DateOrder   int
	Timezones   map[string]interface{}
	Debug       func(string)
}

// Settings holds comprehensive parsing settings (forward declaration to avoid cycles).
// The actual definition is in the main package.
type Settings struct {
	// Simplified structure - actual fields are in main package
	_placeholder byte //nolint:unused // Placeholder to prevent empty struct, actual fields are in main package
}

// ReferenceWithTimezone represents a reference date/time with an optional timezone offset.
// It is used as the reference point for parsing relative dates and times.
//
// NOTE: Fields are exported to allow methods in the main package to access them via type aliases.
// External code should not access these fields directly.
type ReferenceWithTimezone struct {
	Instant        time.Time
	TimezoneOffset *int
}

// ParsingComponents represents a collection of parsed date/time components.
// Components are stored as either "known" (directly parsed) or "implied" (inferred).
//
// NOTE: This type has been moved to internal/parsing. External code should use
// the Components interface or the type alias from the main kronos package.
//
// Fields are exported to allow methods in the main package to access them via type aliases.
// External code should not access these fields directly.
type ParsingComponents struct {
	KnownValues   map[Component]int
	ImpliedValues map[Component]int
	Reference     *ReferenceWithTimezone
	Tags          map[string]bool
	Period        Period
}

// ParsingResult represents a parsed result containing date/time information.
//
// NOTE: This type has been moved to internal/parsing. External code should use
// the Result interface or the type alias from the main kronos package.
//
// Fields are exported to allow methods in the main package to access them via type aliases.
// External code should not access these fields directly.
type ParsingResult struct {
	Reference *ReferenceWithTimezone
	RefDate   time.Time
	Index     int
	Text      string
	Start     *ParsingComponents
	End       *ParsingComponents
}

// ParsingResultWithBoundary wraps ParsingComponents with boundary information.
// This is used internally to communicate the adjusted text (without boundary) to chrono.go
// when parsers using AbstractParserWithWordBoundary return ParsingComponents.
//
// Fields are exported to allow usage from main package.
type ParsingResultWithBoundary struct {
	Components         *ParsingComponents
	AdjustedText       string
	BoundaryLen        int
	IncludeBoundaryIdx bool // If true, index points past boundary; if false, index points at boundary start
}

// ParsingContext holds the context for parsing operations.
// It contains the text to parse, options, and reference information.
//
// NOTE: This type has been moved to internal/parsing. External code should use
// the type alias from the main kronos package.
//
// Fields are exported to allow methods in the main package to access them via type aliases.
// External code should not access these fields directly.
type ParsingContext struct {
	Text      string
	Option    ParsingOption
	Reference *ReferenceWithTimezone
	RefDate   time.Time
	Settings  *Settings
}
