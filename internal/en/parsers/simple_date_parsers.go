//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/en/data"
	"github.com/kljensen/kronos/internal/helpers"
	"github.com/kljensen/kronos/internal/parsing"
)

// ENMonthNameParser parses standalone month names with optional year
// Examples: "January", "January, 2012", "in June of 2022", "Sep 2012"
type ENMonthNameParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENMonthNameParser creates a new ENMonthNameParser
func NewENMonthNameParser() *ENMonthNameParser {
	parser := &ENMonthNameParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENMonthNameParser) innerPattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	// Pattern: (in)? MONTH (,|-|of)? YEAR?
	pattern := `(?i)((?:in)\s*)?` +
		`(` + data.MonthPattern + `)` +
		`(?:` +
		`(?:\s*(?:,|-|of))?\s*(` + data.YearPattern + `)` +
		`)?`

	return regexp.MustCompile(pattern)
}

func (p *ENMonthNameParser) innerExtract(context *kronos.InternalParsingContext, match []string) interface{} {
	if len(match) < 3 {
		return nil
	}

	monthName := strings.ToLower(match[2])

	// Skip unlikely short words unless they're full month names
	if len(match[0]) <= 3 {
		if _, ok := data.FullMonthNameDictionary[monthName]; !ok {
			return nil
		}
	}

	month := data.MonthDictionary[monthName]

	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Imply(kronos.ComponentDay, 1).
		AddTag("parser/ENMonthNameParser")

	// Handle year if present
	if len(match) > 3 && match[3] != "" {
		year := data.ParseYear(match[3])
		components.Assign(kronos.ComponentYear, year)
	} else {
		// Find closest year to reference using preference setting
		year := helpers.FindYearClosestToRefWithPreference(context.RefDate(), 1, month, context.Option().Preference)
		components.Imply(kronos.ComponentYear, year)
	}

	// Set period to month-level
	components.SetPeriod(kronos.PeriodMonth)

	return components
}

// ENSlashMonthFormatParser parses MM/YYYY format
// Examples: 11/2005, 06/2005
type ENSlashMonthFormatParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENSlashMonthFormatParser creates a new ENSlashMonthFormatParser
func NewENSlashMonthFormatParser() *ENSlashMonthFormatParser {
	parser := &ENSlashMonthFormatParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENSlashMonthFormatParser) innerPattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	// Pattern: MM/YYYY (month 1-12, year 4 digits)
	pattern := `(?i)([0-9]|0[1-9]|1[012])/([0-9]{4})`
	return regexp.MustCompile(pattern)
}

func (p *ENSlashMonthFormatParser) innerExtract(context *kronos.InternalParsingContext, match []string) interface{} {
	if len(match) < 3 {
		return nil
	}

	month, err := strconv.Atoi(match[1])
	if err != nil {
		return nil
	}
	year, err := strconv.Atoi(match[2])
	if err != nil {
		return nil
	}

	return context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Assign(kronos.ComponentYear, year).
		Imply(kronos.ComponentDay, 1).
		SetPeriod(kronos.PeriodMonth)
}

// ENYearParser parses standalone year expressions
// Examples: "2020", "1999", "2025"
type ENYearParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENYearParser creates a new ENYearParser
func NewENYearParser() *ENYearParser {
	parser := &ENYearParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
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

func (p *ENYearParser) innerPattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	// Match 4-digit years (1000-2999)
	// This is more restrictive than data.YearPattern to avoid matching other numbers
	// Note: AbstractParserWithWordBoundary will prepend (^|[\s,;:!?()]) for the left boundary
	// Right boundary: must be followed by common punctuation, whitespace, or end of string
	// This prevents matching years within date expressions like "2020-01-15" or "2020@01@15"
	// For periods, we also capture the next character to check for malformed dates
	pattern := `([12][0-9]{3})(\s|,|;|:|\?|!|\)|$|\.([^.\d\-/]|$))`
	return regexp.MustCompile(pattern)
}

func (p *ENYearParser) innerExtract(context *kronos.InternalParsingContext, match []string) interface{} {
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
	return &kronos.InternalParsingResultWithBoundary{
		Components:         components,
		AdjustedText:       adjustedText,
		BoundaryLen:        0,    // Will be set by AbstractParserWithWordBoundary
		IncludeBoundaryIdx: true, // Will be overridden by AbstractParserWithWordBoundary
	}
}
