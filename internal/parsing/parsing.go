// Package parsing contains internal parsing implementation types for kronos.
// These types were previously exposed in the main kronos package but have been moved
// to internal to reduce the public API surface.
//
// The parsing package provides:
//   - Context - holds parsing context (text, options, reference)
//   - Components - represents parsed date/time components with certainty levels
//   - Result - represents a parsed result with date/time information
//   - ReferenceWithTimezone - represents a reference date/time with timezone
//
// NOTE: This package cannot import the main kronos package to avoid import cycles.
// Therefore, it uses mirrored types where necessary. The main package provides
// type aliases and wrapper functions to maintain API compatibility.
//
// External code should not use types from this package directly. Instead, use:
//   - The ParsedComponents and ParsedResult interfaces from kronos
//   - The ParsingComponents and ParsingResult type aliases from kronos (deprecated)
//
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsing

import (
	"regexp"
	"time"

	kronos "github.com/kljensen/kronos"
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

// AbstractParserWithWordBoundary is a base parser that checks for word boundaries
// before applying the inner pattern and extraction.
// It wraps the inner pattern with word boundary checks and adjusts the match
// before passing it to the inner extraction logic.
type AbstractParserWithWordBoundary struct {
	// innerPattern should return the core regex pattern (without boundary checks)
	innerPattern func(context *kronos.InternalParsingContext) *regexp.Regexp

	// innerExtract should extract components from the match
	innerExtract func(context *kronos.InternalParsingContext, match []string) interface{}

	// patternLeftBoundary returns the left boundary pattern
	patternLeftBoundary func() string

	// Cache for pattern
	cachedInnerPattern *regexp.Regexp
	cachedPattern      *regexp.Regexp
}

// NewAbstractParserWithWordBoundary creates a new AbstractParserWithWordBoundary.
func NewAbstractParserWithWordBoundary(
	innerPattern func(context *kronos.InternalParsingContext) *regexp.Regexp,
	innerExtract func(context *kronos.InternalParsingContext, match []string) interface{},
	patternLeftBoundary func() string,
) *AbstractParserWithWordBoundary {
	if patternLeftBoundary == nil {
		patternLeftBoundary = func() string {
			return `(\W|^)`
		}
	}

	return &AbstractParserWithWordBoundary{
		innerPattern:        innerPattern,
		innerExtract:        innerExtract,
		patternLeftBoundary: patternLeftBoundary,
	}
}

// Pattern implements Parser.Pattern.
// It wraps the inner pattern with word boundary checks.
func (p *AbstractParserWithWordBoundary) Pattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	currentInnerPattern := p.innerPattern(context)

	// Check if we need to rebuild the pattern
	if p.cachedInnerPattern != nil && p.cachedInnerPattern == currentInnerPattern {
		return p.cachedPattern
	}

	// Rebuild the pattern with boundaries
	p.cachedInnerPattern = currentInnerPattern
	boundaryPrefix := p.patternLeftBoundary()

	// Combine boundary + inner pattern
	patternStr := boundaryPrefix + currentInnerPattern.String()

	// Preserve the flags from the inner pattern
	flags := ""
	if isRegexpCaseInsensitive(currentInnerPattern) {
		flags = "(?i)"
	}

	p.cachedPattern = regexp.MustCompile(flags + patternStr)
	return p.cachedPattern
}

// Extract implements Parser.Extract.
// It validates word boundaries and adjusts the match before calling innerExtract.
func (p *AbstractParserWithWordBoundary) Extract(context *kronos.InternalParsingContext, match []string) interface{} {
	if len(match) < 2 {
		return p.innerExtract(context, match)
	}

	// The first capture group is the boundary (e.g., non-word char or start of string)
	header := match[1]
	headerLen := len(header)

	// Adjust the match array by removing the header
	// This makes the inner pattern's capture groups align correctly
	adjustedMatch := make([]string, len(match)-1)
	adjustedMatch[0] = match[0][headerLen:] // Full match without header

	// Shift remaining groups down by 1
	for i := 2; i < len(match); i++ {
		adjustedMatch[i-1] = match[i]
	}

	result := p.innerExtract(context, adjustedMatch)

	// Wrap the result to communicate boundary information to chrono.go
	switch v := result.(type) {
	case *kronos.InternalParsingResult:
		// Index should point past the boundary (skip the boundary characters)
		// headerLen is the length of the boundary to skip
		v.SetIndex(headerLen)
		return v
	case *kronos.InternalParsingResultWithBoundary:
		// Parser returned a ParsingResultWithBoundary with custom adjusted text
		// Update it with the correct boundary length
		v.BoundaryLen = headerLen
		v.IncludeBoundaryIdx = true
		return v
	case *kronos.InternalParsingComponents:
		// Wrap in a struct that provides both the components and the adjusted text
		// The text should NOT include the boundary character
		// Trim trailing whitespace that may have been matched by (?:\s|$|\b)
		adjustedText := adjustedMatch[0]
		for len(adjustedText) > 0 {
			lastChar := adjustedText[len(adjustedText)-1]
			if lastChar == ' ' || lastChar == '\t' || lastChar == '\n' || lastChar == '\r' {
				adjustedText = adjustedText[:len(adjustedText)-1]
			} else {
				break
			}
		}
		return &kronos.InternalParsingResultWithBoundary{
			Components:         v,
			AdjustedText:       adjustedText, // Exclude boundary and trailing space from text
			BoundaryLen:        headerLen,    // Length of boundary to skip
			IncludeBoundaryIdx: true,         // Index should point past the boundary
		}
	default:
		return result
	}
}

// Helper function to check if a regexp has case-insensitive flag
func isRegexpCaseInsensitive(re *regexp.Regexp) bool {
	// Go's regexp doesn't expose flags directly, so we check the string representation
	str := re.String()
	return len(str) >= 4 && str[:4] == "(?i)"
}
