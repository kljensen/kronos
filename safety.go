package kronos

// AsParsingComponents attempts to cast the given ParsedComponents interface to
// a *ParsingComponents. It returns the concrete value and true when the cast
// succeeds, or nil and false otherwise.
func AsParsingComponents(pc ParsedComponents) (*ParsingComponents, bool) {
	if pc == nil {
		return nil, false
	}

	components, ok := pc.(*ParsingComponents)
	if !ok || components == nil {
		return nil, false
	}

	return components, true
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
