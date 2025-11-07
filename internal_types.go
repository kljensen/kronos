package kronos

import "time"

// This file provides access to internal types for the internal/ packages.
// These types are intentionally unexported to external users but needed by
// internal implementation packages.
//
// INTERNAL USE ONLY - These exports are implementation details that may
// change without notice. External code should not use these types.

// Type aliases for internal packages to use unexported types
type (
	InternalParsingContext            = parsingContext
	InternalParsingComponents         = parsingComponents
	InternalParsingResult             = parsingResult
	InternalReferenceWithTimezone     = referenceWithTimezone
	InternalParsingOption             = parsingOption
	InternalParsingReference          = parsingReference
	InternalParsingResultWithBoundary = parsingResultWithBoundary
)

// Constructor functions for internal packages

// InternalNewParsingContext creates a parsing context (internal use only).
func InternalNewParsingContext(text string, refDate interface{}, option *parsingOption) *parsingContext {
	return newParsingContextFromOption(text, refDate, option)
}

// InternalNewParsingComponents creates parsing components (internal use only).
func InternalNewParsingComponents(reference *referenceWithTimezone, knownComponents map[Component]int) *parsingComponents {
	return newParsingComponents(reference, knownComponents)
}

// InternalNewParsingResult creates a parsing result (internal use only).
func InternalNewParsingResult(reference *referenceWithTimezone, index int, text string, start, end *parsingComponents) *parsingResult {
	return newParsingResult(reference, index, text, start, end)
}

// InternalNewReferenceWithTimezone creates a reference with timezone (internal use only).
func InternalNewReferenceWithTimezone(instant time.Time, timezoneOffset *int) *referenceWithTimezone {
	return newReferenceWithTimezone(instant, timezoneOffset)
}

// InternalMergeDateTimeResult merges date and time results (internal use only).
func InternalMergeDateTimeResult(dateResult, timeResult *parsingResult) *parsingResult {
	return mergeDateTimeResult(dateResult, timeResult)
}

// InternalValidateSettings validates parsing settings (internal use only).
func InternalValidateSettings(s Settings) error {
	return validateSettings(s)
}

// InternalParseWithSettings parses text with custom settings (internal use only).
func InternalParseWithSettings(text string, refDate time.Time, settings Settings, config *Configuration) ([]*parsingResult, error) {
	return parseWithSettings(text, refDate, settings, config)
}

// Weekday calculation helpers for internal packages

// InternalGetNthWeekdayOfMonth wraps the unexported weekday calculation (internal use only).
func InternalGetNthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int, hour int) time.Time {
	return getNthWeekdayOfMonth(year, month, weekday, n, hour)
}

// InternalGetLastWeekdayOfMonth wraps the unexported weekday calculation (internal use only).
func InternalGetLastWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, hour int) time.Time {
	return getLastWeekdayOfMonth(year, month, weekday, hour)
}
