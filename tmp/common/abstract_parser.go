package common

import (
	"regexp"
	"unicode"

	kronos "github.com/kljensen/kronos"
)

// AbstractParserWithWordBoundary is a base parser that checks for word boundaries
// before applying the inner pattern and extraction.
// It wraps the inner pattern with word boundary checks and adjusts the match
// before passing it to the inner extraction logic.
type AbstractParserWithWordBoundary struct {
	// innerPattern should return the core regex pattern (without boundary checks)
	innerPattern func(context *kronos.ParsingContext) *regexp.Regexp

	// innerExtract should extract components from the match
	innerExtract func(context *kronos.ParsingContext, match []string) interface{}

	// patternLeftBoundary returns the left boundary pattern
	patternLeftBoundary func() string

	// Cache for pattern
	cachedInnerPattern *regexp.Regexp
	cachedPattern      *regexp.Regexp
}

// NewAbstractParserWithWordBoundary creates a new AbstractParserWithWordBoundary.
func NewAbstractParserWithWordBoundary(
	innerPattern func(context *kronos.ParsingContext) *regexp.Regexp,
	innerExtract func(context *kronos.ParsingContext, match []string) interface{},
	patternLeftBoundary func() string,
) *AbstractParserWithWordBoundary {
	if patternLeftBoundary == nil {
		patternLeftBoundary = func() string {
			return `(\W|^)`
		}
	}

	return &AbstractParserWithWordBoundary{
		innerPattern:        innerPattern,
		innerExtract:        innerExtract,
		patternLeftBoundary: patternLeftBoundary,
	}
}

// Pattern implements Parser.Pattern.
// It wraps the inner pattern with word boundary checks.
func (p *AbstractParserWithWordBoundary) Pattern(context *kronos.ParsingContext) *regexp.Regexp {
	currentInnerPattern := p.innerPattern(context)

	// Check if we need to rebuild the pattern
	if p.cachedInnerPattern != nil && p.cachedInnerPattern == currentInnerPattern {
		return p.cachedPattern
	}

	// Rebuild the pattern with boundaries
	p.cachedInnerPattern = currentInnerPattern
	boundaryPrefix := p.patternLeftBoundary()

	// Combine boundary + inner pattern
	patternStr := boundaryPrefix + currentInnerPattern.String()

	// Preserve the flags from the inner pattern
	flags := ""
	if isRegexpCaseInsensitive(currentInnerPattern) {
		flags = "(?i)"
	}

	p.cachedPattern = regexp.MustCompile(flags + patternStr)
	return p.cachedPattern
}

// Extract implements Parser.Extract.
// It validates word boundaries and adjusts the match before calling innerExtract.
func (p *AbstractParserWithWordBoundary) Extract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 2 {
		return p.innerExtract(context, match)
	}

	// The first capture group is the boundary (e.g., non-word char or start of string)
	header := match[1]

	// Adjust the match array by removing the header
	// This makes the inner pattern's capture groups align correctly
	adjustedMatch := make([]string, len(match)-1)
	adjustedMatch[0] = match[0][len(header):] // Full match without header

	// Shift remaining groups down by 1
	for i := 2; i < len(match); i++ {
		adjustedMatch[i-1] = match[i]
	}

	return p.innerExtract(context, adjustedMatch)
}

// Helper function to check if a regexp has case-insensitive flag
func isRegexpCaseInsensitive(re *regexp.Regexp) bool {
	// Go's regexp doesn't expose flags directly, so we check the string representation
	str := re.String()
	return len(str) >= 4 && str[:4] == "(?i)"
}

// isWordBoundary checks if there's a word boundary at the given position
func isWordBoundary(text string, pos int) bool {
	if pos <= 0 || pos >= len(text) {
		return true
	}

	prevChar := rune(text[pos-1])
	currChar := rune(text[pos])

	prevIsWord := unicode.IsLetter(prevChar) || unicode.IsDigit(prevChar) || prevChar == '_'
	currIsWord := unicode.IsLetter(currChar) || unicode.IsDigit(currChar) || currChar == '_'

	return prevIsWord != currIsWord
}
