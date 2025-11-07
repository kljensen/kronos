package kronos

import "time"

// This file provides access to internal types for the internal/ packages.
// These types are intentionally unexported to external users but needed by
// internal implementation packages.
//
// INTERNAL USE ONLY - These exports are implementation details that may
// change without notice. External code should not use these types.

// InternalParsingContext is an alias for parsingContext, allowing internal packages
// to use the unexported type.
type InternalParsingContext = parsingContext

// InternalParsingComponents is an alias for parsingComponents, allowing internal packages
// to use the unexported type.
type InternalParsingComponents = parsingComponents

// InternalParsingResult is an alias for parsingResult, allowing internal packages
// to use the unexported type.
type InternalParsingResult = parsingResult

// InternalReferenceWithTimezone is an alias for referenceWithTimezone, allowing internal packages
// to use the unexported type.
type InternalReferenceWithTimezone = referenceWithTimezone

// InternalParsingOption is an alias for parsingOption, allowing internal packages
// to use the unexported type.
type InternalParsingOption = parsingOption

// InternalParsingReference is an alias for parsingReference, allowing internal packages
// to use the unexported type.
type InternalParsingReference = parsingReference

// InternalParsingResultWithBoundary is an alias for parsingResultWithBoundary, allowing internal packages
// to use the unexported type.
type InternalParsingResultWithBoundary = parsingResultWithBoundary

// Constructor functions for internal packages

// InternalNewParsingContext creates a parsing context (internal use only).
func InternalNewParsingContext(text string, refDate any, option *parsingOption) *parsingContext {
	return newParsingContext(text, refDate, option)
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

// InternalMergeDateTimeComponent merges date and time components (internal use only).
func InternalMergeDateTimeComponent(dateComp, timeComp *parsingComponents) *parsingComponents {
	return mergeDateTimeComponent(dateComp, timeComp)
}

// InternalDeterminePeriodFromDuration determines period from duration (internal use only).
func InternalDeterminePeriodFromDuration(duration Duration) Period {
	return determinePeriodFromDuration(duration)
}

// Weekday calculation helpers for internal packages

// weekdayTo1Indexed converts time.Weekday (0=Sunday) to 1-indexed (1=Monday, 7=Sunday).
func weekdayTo1Indexed(weekday time.Weekday) int {
	w := int(weekday)
	if w == 0 {
		return 7
	}
	return w
}

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
	targetWeekday := weekdayTo1Indexed(weekday)
	firstWeekdayNextMonth := weekdayTo1Indexed(nextMonth.Weekday())

	// Calculate how many days to go back
	dayDiff := firstWeekdayNextMonth - targetWeekday
	if dayDiff <= 0 {
		dayDiff += 7
	}

	// Go back to find the last occurrence
	result := nextMonth.AddDate(0, 0, -dayDiff)
	return time.Date(result.Year(), result.Month(), result.Day(), hour, 0, 0, 0, time.UTC)
}
