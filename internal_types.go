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

// InternalGetNthWeekdayOfMonth returns the date of the nth occurrence of a given weekday
// in a given month and year (internal use only).
func InternalGetNthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int, hour int) time.Time {
	dayOfMonth := 0
	count := 0

	// Iterate through days of the month
	for count < n {
		dayOfMonth++
		date := time.Date(year, time.Month(month), dayOfMonth, hour, 0, 0, 0, time.UTC)

		// Check if this day is the target weekday
		if int(date.Weekday()) == int(weekday) {
			count++
		}
	}

	return time.Date(year, time.Month(month), dayOfMonth, hour, 0, 0, 0, time.UTC)
}

// InternalGetLastWeekdayOfMonth returns the date of the last occurrence of a given weekday
// in a given month and year (internal use only).
func InternalGetLastWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, hour int) time.Time {
	// Start with the first day of the next month
	nextMonth := time.Date(year, time.Month(month)+1, 1, 12, 0, 0, 0, time.UTC)

	// Convert weekdays to 1-indexed (Monday=1, Sunday=7)
	targetWeekday := int(weekday)
	if targetWeekday == 0 {
		targetWeekday = 7
	}

	firstWeekdayNextMonth := int(nextMonth.Weekday())
	if firstWeekdayNextMonth == 0 {
		firstWeekdayNextMonth = 7
	}

	// Calculate how many days to go back
	var dayDiff int
	switch {
	case firstWeekdayNextMonth == targetWeekday:
		dayDiff = 7
	case firstWeekdayNextMonth < targetWeekday:
		dayDiff = 7 + firstWeekdayNextMonth - targetWeekday
	default:
		dayDiff = firstWeekdayNextMonth - targetWeekday
	}

	// Go back to find the last occurrence
	result := nextMonth.AddDate(0, 0, -dayDiff)
	return time.Date(result.Year(), result.Month(), result.Day(), hour, 0, 0, 0, time.UTC)
}
