package kronos

import "time"

// Internal helper functions and type aliases for accessing internal types.
// These are used by internal/ packages only.
//
// WARNING: These are NOT part of the public API. External users should NOT
// use these - they expose implementation details that may change without notice.

// Type aliases to expose internal types to internal/ packages.
// These types are exported to allow internal packages to use them,
// but external code should use the public interfaces instead (Result, Components, etc.)
type (
	InternalParsingContext            = parsingContext
	InternalParsingComponents         = parsingComponents
	InternalParsingResult             = parsingResult
	InternalReferenceWithTimezone     = referenceWithTimezone
	InternalParsingOption             = parsingOption
	InternalParsingReference          = parsingReference
	InternalParsingResultWithBoundary = parsingResultWithBoundary
)

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
