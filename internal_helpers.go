package kronos

// Internal helper functions exported with X prefix for use by internal/ packages only.
// These should NOT be used by external code.
//
// NOTE: Most X-prefixed helpers have been retired and moved to internal/helpers.
// Only a minimal set of functions that require private field access remain here.

// XNewParsingContext is an internal helper for internal/ and experimental packages.
// It creates a new ParsingContext. This must remain in the root package because
// it needs access to private fields.
func XNewParsingContext(text string, refDate interface{}, option *ParsingOption) *ParsingContext {
	return newParsingContext(text, refDate, option)
}

// XNewParsingComponents is an internal helper for internal/ and experimental packages.
// It creates a new ParsingComponents. This must remain in the root package because
// it needs access to private fields.
func XNewParsingComponents(reference *ReferenceWithTimezone, knownComponents map[Component]int) *ParsingComponents {
	return newParsingComponents(reference, knownComponents)
}

// XNewParsingResult is an internal helper for internal/ and experimental packages.
// It creates a new ParsingResult. This must remain in the root package because
// it needs access to private fields.
func XNewParsingResult(reference *ReferenceWithTimezone, index int, text string, start, end *ParsingComponents) *ParsingResult {
	return newParsingResult(reference, index, text, start, end)
}

// XMergeDateTimeResult is an internal helper for internal/ and experimental packages.
// It merges date and time results. This must remain in the root package because
// it needs access to private fields.
func XMergeDateTimeResult(dateResult, timeResult *ParsingResult) *ParsingResult {
	return mergeDateTimeResult(dateResult, timeResult)
}
