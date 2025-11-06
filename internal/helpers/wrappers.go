package helpers

import "github.com/kljensen/kronos"

// AsParsingComponents attempts to cast the given ParsedComponents interface to
// a *ParsingComponents. It returns the concrete value and true when the cast
// succeeds, or nil and false otherwise.
func AsParsingComponents(pc kronos.ParsedComponents) (*kronos.ParsingComponents, bool) {
	if pc == nil {
		return nil, false
	}

	components, ok := pc.(*kronos.ParsingComponents)
	if !ok || components == nil {
		return nil, false
	}

	return components, true
}

// NewParsingComponents creates a new ParsingComponents with the given reference.
// It initializes implied values based on the reference date.
//
// This is a wrapper around the internal implementation in the kronos package.
// We must use the X function because it needs access to private fields.
func NewParsingComponents(reference *kronos.ReferenceWithTimezone, knownComponents map[kronos.Component]int) *kronos.ParsingComponents {
	return kronos.XNewParsingComponents(reference, knownComponents)
}

// MergeDateTimeResult merges a date-only result with a time-only result.
//
// This is a wrapper around the internal implementation in the kronos package.
// We must use the X function because it needs access to private fields.
func MergeDateTimeResultWrapper(dateResult, timeResult *kronos.ParsingResult) *kronos.ParsingResult {
	return kronos.XMergeDateTimeResult(dateResult, timeResult)
}
