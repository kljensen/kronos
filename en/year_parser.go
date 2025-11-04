package en

import (
	"regexp"
	"strconv"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENYearParser parses standalone year expressions
// Examples: "2020", "1999", "2025"
type ENYearParser struct {
	*common.AbstractParserWithWordBoundary
}

// NewENYearParser creates a new ENYearParser
func NewENYearParser() *ENYearParser {
	parser := &ENYearParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		// Use custom boundary that excludes date separators (-, /, @, .)
		// This prevents matching years within date expressions like "2020-01-15" or "2020@01@15"
		// Allow common sentence boundaries: start, whitespace, some punctuation
		// Note: Period is excluded from leading boundary to avoid matching in "2020..01"
		// Must be a CAPTURING group for AbstractParserWithWordBoundary
		func() string {
			return `(^|[\s,;:!?\(\)])`
		},
	)

	return parser
}

func (p *ENYearParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Match 4-digit years (1000-2999)
	// This is more restrictive than YEAR_PATTERN to avoid matching other numbers
	// Note: AbstractParserWithWordBoundary will prepend (^|[\s,;:!?()]) for the left boundary
	// Right boundary: must be followed by common punctuation, whitespace, or end of string
	// This prevents matching years within date expressions like "2020-01-15" or "2020@01@15"
	// For periods, we also capture the next character to check for malformed dates
	pattern := `([12][0-9]{3})(\s|,|;|:|\?|!|\)|$|\.([^.\d\-/]|$))`
	return regexp.MustCompile(pattern)
}

func (p *ENYearParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 2 {
		return nil
	}

	year, err := strconv.Atoi(match[1])
	if err != nil {
		return nil
	}

	// The match[0] is the full match including trailing characters
	// We need to extract just the year part (4 digits)
	adjustedText := match[1] // Just the year, no trailing chars

	// Create components for the year
	// Date is set to January 1st of the year
	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentYear, year).
		Imply(kronos.ComponentMonth, 1).
		Imply(kronos.ComponentDay, 1).
		SetPeriod(kronos.PeriodYear)

	// Return a ParsingResultWithBoundary with the adjusted text
	// This excludes the trailing character
	return &kronos.ParsingResultWithBoundary{
		Components:         components,
		AdjustedText:       adjustedText,
		BoundaryLen:        0,    // Will be set by AbstractParserWithWordBoundary
		IncludeBoundaryIdx: true, // Will be overridden by AbstractParserWithWordBoundary
	}
}
