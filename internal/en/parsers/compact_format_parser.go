//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/parsing"
)

// ENCompactFormatParser handles dates/times without separators
// Supports formats like:
// - 20200315 (YYYYMMDD)
// - 200315 (YYMMDD)
// - 20200315143045 (YYYYMMDDHHmmss)
// - 1430 (HHmm)
// - 143045 (HHmmss)
type ENCompactFormatParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENCompactFormatParser creates a new compact format parser
func NewENCompactFormatParser() *ENCompactFormatParser {
	parser := &ENCompactFormatParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		// Use custom boundary that excludes hyphens and digits
		// This prevents matching numbers that are part of hyphenated expressions like "2012-14"
		// Must be a CAPTURING group for AbstractParserWithWordBoundary
		func() string {
			return `(^|[^\w-])`
		},
	)

	return parser
}

func (p *ENCompactFormatParser) innerPattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	// Match sequences of 4-14 digits
	// We want to match compact date/time formats but not:
	// - Single/double/triple digits
	// - Very long numbers (15+ digits like unix timestamps)
	// - Numbers followed by hyphens (like "2012-" in "2012-14")
	// Capture the trailing character to properly handle word boundaries
	return regexp.MustCompile(`(\d{4}|\d{6}|\d{8}|\d{10}|\d{12}|\d{14})([^\w-]|$)`)
}

func (p *ENCompactFormatParser) innerExtract(context *kronos.InternalParsingContext, match []string) any {
	if len(match) < 3 {
		return nil
	}

	digitStr := match[1]
	trailingChar := match[2]

	// Remove the trailing character from the match text
	adjustedText := match[0]
	if len(trailingChar) > 0 && len(adjustedText) > 0 {
		adjustedText = adjustedText[:len(adjustedText)-len(trailingChar)]
	}

	var components *kronos.InternalParsingComponents

	switch len(digitStr) {
	case 4:
		// Could be MMDD or HHmm
		components = p.tryParse4Digits(digitStr, context)
	case 6:
		// Could be YYMMDD or HHmmss
		components = p.tryParse6Digits(digitStr, context)
	case 8:
		// YYYYMMDD
		components = p.tryParse8Digits(digitStr, context)
	case 10:
		// YYYYMMDDHHmm (no seconds)
		components = p.tryParse10Digits(digitStr, context)
	case 12:
		// YYYYMMDDHHmmss
		components = p.tryParse12Digits(digitStr, context)
	case 14:
		// YYYYMMDDHHmmss with milliseconds
		components = p.tryParse14Digits(digitStr, context)
	default:
		return nil
	}

	if components == nil {
		return nil
	}

	// Return a ParsingResultWithBoundary with the adjusted text
	// This excludes the trailing character
	return &kronos.InternalParsingResultWithBoundary{
		Components:         components,
		AdjustedText:       adjustedText,
		BoundaryLen:        0,    // Will be set by AbstractParserWithWordBoundary
		IncludeBoundaryIdx: true, // Will be overridden by AbstractParserWithWordBoundary
	}
}
