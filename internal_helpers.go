package kronos

import "time"

// Internal helper functions for accessing private fields.
// These are used by internal/ and experimental packages.

// These functions are NOT part of the public API and should only be used by
// internal/ and experimental/ packages. They are capitalized to be accessible
// from internal packages but are prefixed with Internal to indicate they are
// not intended for external use.

// Type aliases for internal use
type (
	InternalParsingContext            = parsingContext
	InternalParsingComponents         = parsingComponents
	InternalParsingResult             = parsingResult
	InternalReferenceWithTimezone     = referenceWithTimezone
	InternalParsingOption             = parsingOption
	InternalParsingReference          = parsingReference
	InternalParsingResultWithBoundary = parsingResultWithBoundary
)

func InternalNewParsingContext(text string, refDate interface{}, option *parsingOption) *parsingContext {
	return newParsingContext(text, refDate, option)
}

func InternalNewParsingComponents(reference *referenceWithTimezone, knownComponents map[Component]int) *parsingComponents {
	return newParsingComponents(reference, knownComponents)
}

func InternalNewParsingResult(reference *referenceWithTimezone, index int, text string, start, end *parsingComponents) *parsingResult {
	return newParsingResult(reference, index, text, start, end)
}

func InternalMergeDateTimeResult(dateResult, timeResult *parsingResult) *parsingResult {
	return mergeDateTimeResult(dateResult, timeResult)
}

func InternalValidateSettings(s Settings) error {
	return validateSettings(s)
}

func InternalParseWithSettings(text string, refDate time.Time, settings Settings, config *Configuration) ([]*parsingResult, error) {
	return parseWithSettings(text, refDate, settings, config)
}
