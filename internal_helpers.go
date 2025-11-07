package kronos

// Internal helper functions for accessing private fields.
// These are used by internal/ and experimental packages.

// These functions are NOT part of the public API and should only be used by
// internal/ and experimental/ packages. They are capitalized to be accessible
// from internal packages but are prefixed with Internal to indicate they are
// not intended for external use.

func InternalNewParsingContext(text string, refDate interface{}, option *ParsingOption) *ParsingContext {
	return newParsingContext(text, refDate, option)
}

func InternalNewParsingComponents(reference *ReferenceWithTimezone, knownComponents map[Component]int) *ParsingComponents {
	return newParsingComponents(reference, knownComponents)
}

func InternalNewParsingResult(reference *ReferenceWithTimezone, index int, text string, start, end *ParsingComponents) *ParsingResult {
	return newParsingResult(reference, index, text, start, end)
}

func InternalMergeDateTimeResult(dateResult, timeResult *ParsingResult) *ParsingResult {
	return mergeDateTimeResult(dateResult, timeResult)
}

func InternalNewParserRegistry() *ParserRegistry {
	return newParserRegistry()
}

func InternalGlobalRegistry() *ParserRegistry {
	return globalRegistry
}

func InternalRegister(name string, info ParserInfo, factory ParserFactory) {
	internalRegister(name, info, factory)
}
