package en

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENMonthNameMiddleEndianParser parses "Month DD, YYYY" format (middle-endian, US-style)
// Examples: "August 10, 2012", "January 13", "Dec 12, 2020"
// Also handles ranges: "January 12 - 15, 2012"
type ENMonthNameMiddleEndianParser struct {
	*common.AbstractParserWithWordBoundary
	shouldSkipYearLikeDate bool
}

// NewENMonthNameMiddleEndianParser creates a new ENMonthNameMiddleEndianParser
func NewENMonthNameMiddleEndianParser(shouldSkipYearLikeDate bool) *ENMonthNameMiddleEndianParser {
	parser := &ENMonthNameMiddleEndianParser{
		shouldSkipYearLikeDate: shouldSkipYearLikeDate,
	}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENMonthNameMiddleEndianParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Pattern: MONTH (-|/|,) DAY (to|-)? DAY? (-|/|,)? YEAR?
	pattern := `(?i)(` + MonthPattern + `)` +
		`(?:-|/|\s*,?\s*)` +
		`(` + OrdinalNumberPattern + `)\s*` +
		`(?:` +
		`(?:to|\-)\s*` +
		`(` + OrdinalNumberPattern + `)\s*` +
		`)?` +
		`(?:` +
		`(?:-|/|\s*,\s*|\s+)` +
		`(` + YearPattern + `)` +
		`)?`

	return regexp.MustCompile(pattern)
}

func (p *ENMonthNameMiddleEndianParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 3 {
		return nil
	}

	monthName := strings.ToLower(match[1])
	month := MonthDictionary[monthName]

	day := ParseOrdinalNumber(match[2])

	// Validate day
	if day > 31 {
		return nil
	}

	// Skip year-like dates if configured (e.g., "January 21" where 21 looks like a year)
	if p.shouldSkipYearLikeDate {
		// No range, no year, and day looks like a year (20-25)
		if match[3] == "" && (len(match) < 5 || match[4] == "") {
			yearLikePattern := regexp.MustCompile(`^2[0-5]$`)
			if yearLikePattern.MatchString(match[2]) {
				return nil
			}
		}
	}

	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Assign(kronos.ComponentDay, day).
		AddTag("parser/ENMonthNameMiddleEndianParser")

	// Handle year if present
	if len(match) > 4 && match[4] != "" {
		year := ParseYear(match[4])
		components.Assign(kronos.ComponentYear, year)
	} else {
		// Find closest year to reference
		year := kronos.FindYearClosestToRef(context.RefDate(), day, month)
		components.Imply(kronos.ComponentYear, year)
	}

	// Handle date range (e.g., "January 12 - 15, 2012")
	if match[3] != "" {
		endDay := ParseOrdinalNumber(match[3])

		// Create result with end date
		endComponents := components.Clone()
		endComponents.Assign(kronos.ComponentDay, endDay)

		result := context.CreateParsingResult(0, match[0], components, endComponents)
		return result
	}

	return components
}
