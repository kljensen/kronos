// Package helpers provides internal utility functions for kronos date parsing.
// This package contains helper functions that were previously exported with X-prefixes
// from the main kronos package. These functions are implementation details and should
// not be used by external code.
//
// The helpers package includes:
//   - Casual date reference functions (Today, Tomorrow, Yesterday, Now, etc.)
//   - Date math utilities (year finding, weekday calculations)
//   - Component manipulation (AssignSimilarDate, MergeDateTimeComponent, etc.)
//   - Duration helpers (AddDuration, ReverseDuration)
//   - String utilities (SafeSlice, StripApproximationWords)
//   - Timezone helpers (ToTimezoneOffset, FromInput)
//   - Result helpers (MergeDateTimeResult, CreateRelativeFromReference)
package helpers

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
)

// Meridiem represents the AM/PM indicator for internal parsing.
type Meridiem int

// Meridiem constants for internal use
const (
	MeridiemAM Meridiem = 0
	MeridiemPM Meridiem = 1
)

// Common time constants for internal calculations
const (
	HoursPerDay            = 24
	MinutesPerHour         = 60
	SecondsPerMinute       = 60
	MillisecondsPerSecond  = 1000
	MicrosecondsPerMS      = 1000
	MicrosecondsPerSecond  = 1000000
	NanosecondsPerMicro    = 1000
	NanosecondsPerMS       = 1000000
	SecondsPerHour         = 3600
	MinutesPerDay          = 1440
	DaysPerWeek            = 7
	MonthsPerYear          = 12
	MonthsPerQuarter       = 3
	WeeksPerMonthApprox    = 4
	YearLookAheadThreshold = 20 // For 2-digit year conversion
)

// AsParsingComponents attempts to cast the given Components interface to
// a *ParsingComponents. It returns the concrete value and true when the cast
// succeeds, or nil and false otherwise.
func AsParsingComponents(pc kronos.Components) (*kronos.InternalParsingComponents, bool) {
	if pc == nil {
		return nil, false
	}

	components, ok := pc.(*kronos.InternalParsingComponents)
	if !ok || components == nil {
		return nil, false
	}

	return components, true
}

// NewParsingContext creates a new ParsingContext for parsing operations.
// This uses the internal helper exported from kronos package for internal use.
func NewParsingContext(text string, refDate interface{}, option *kronos.InternalParsingOption) *kronos.InternalParsingContext {
	return kronos.InternalNewParsingContext(text, refDate, option)
}

// NewParsingComponents creates a new ParsingComponents with the given reference.
// It initializes implied values based on the reference date.
//
// This uses the internal helper exported from kronos package for internal use.
func NewParsingComponents(reference *kronos.InternalReferenceWithTimezone, knownComponents map[kronos.Component]int) *kronos.InternalParsingComponents {
	return kronos.InternalNewParsingComponents(reference, knownComponents)
}

// NewParsingResult creates a new ParsingResult with the given parameters.
// This uses the internal helper exported from kronos package for internal use.
func NewParsingResult(reference *kronos.InternalReferenceWithTimezone, index int, text string, start, end *kronos.InternalParsingComponents) *kronos.InternalParsingResult {
	return kronos.InternalNewParsingResult(reference, index, text, start, end)
}

// MergeDateTimeResultWrapper merges a date-only result with a time-only result.
//
// This is a wrapper around the internal implementation in the kronos package.
func MergeDateTimeResultWrapper(dateResult, timeResult *kronos.InternalParsingResult) *kronos.InternalParsingResult {
	return kronos.InternalMergeDateTimeResult(dateResult, timeResult)
}

// ApproximationWords is a list of words that indicate approximate time expressions
var ApproximationWords = []string{
	"about",
	"around",
	"roughly",
	"approximately",
	"approx",
	"circa",
}

// SafeSlice returns the substring of text between start (inclusive) and end
// (exclusive) while clamping the requested range to valid bounds. The returned
// boolean is false when start is beyond the end of the string, indicating that
// the requested slice could not be produced safely.
func SafeSlice(text string, start, end int) (string, bool) {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	if start > len(text) {
		return "", false
	}
	if end > len(text) {
		end = len(text)
	}

	return text[start:end], true
}

// StripApproximationWords removes approximation modifiers from the input text
// and returns the cleaned text along with a flag indicating if approximation was present.
// The approximation words are matched case-insensitively and removed with their trailing whitespace.
func StripApproximationWords(input string) (cleaned string, isApproximate bool) {
	cleaned = input
	isApproximate = false

	// Check for tilde (~) symbol as approximation marker
	tildePattern := regexp.MustCompile(`(?i)~\s*`)
	if tildePattern.MatchString(cleaned) {
		isApproximate = true
		cleaned = tildePattern.ReplaceAllString(cleaned, "")
	}

	// Check for approximation words
	for _, word := range ApproximationWords {
		// Use word boundaries to avoid matching words like "about" in "roundabout"
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\s+`)
		if pattern.MatchString(cleaned) {
			isApproximate = true
			cleaned = pattern.ReplaceAllString(cleaned, "")
		}
	}

	return strings.TrimSpace(cleaned), isApproximate
}
