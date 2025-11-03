package en

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENMonthNameLittleEndianParser parses "DD Month YYYY" format (little-endian)
// Examples: "10 August 2012", "3rd Feb 82", "31st March, 2016"
// Also handles ranges: "10 to 15 August 2012"
type ENMonthNameLittleEndianParser struct {
	*common.AbstractParserWithWordBoundary
}

// NewENMonthNameLittleEndianParser creates a new ENMonthNameLittleEndianParser
func NewENMonthNameLittleEndianParser() *ENMonthNameLittleEndianParser {
	parser := &ENMonthNameLittleEndianParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENMonthNameLittleEndianParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Pattern: (on)? DAY (to|-)? DAY? (-|/|of) MONTH (-|/|,)? YEAR?
	pattern := `(?i)(?:on\s{0,3})?` +
		`(` + OrdinalNumberPattern + `)` +
		`(?:` +
		`\s{0,3}(?:to|\-|until|through|till)?\s{0,3}` +
		`(` + OrdinalNumberPattern + `)` +
		`)?` +
		`(?:-|/|\s{0,3}(?:of)?\s{0,3})` +
		`(` + MonthPattern + `)` +
		`(?:` +
		`(?:-|/|,?\s{0,3})` +
		`(` + YearPattern + `)` +
		`)?` +
		``

	return regexp.MustCompile(pattern)
}

func (p *ENMonthNameLittleEndianParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 4 {
		return nil
	}

	day := ParseOrdinalNumber(match[1])

	// Validate day - if > 31, it might be a year-like number (e.g., "96 Aug")
	if day > 31 {
		return nil
	}

	monthName := strings.ToLower(match[3])
	month := MonthDictionary[monthName]

	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Assign(kronos.ComponentDay, day).
		AddTag("parser/ENMonthNameLittleEndianParser")

	// Handle year if present
	if len(match) > 4 && match[4] != "" {
		year := ParseYear(match[4])
		components.Assign(kronos.ComponentYear, year)
	} else {
		// Find closest year to reference using preference setting
		year := kronos.FindYearClosestToRefWithPreference(context.RefDate(), day, month, context.Option().Preference)
		components.Imply(kronos.ComponentYear, year)
	}

	// Handle date range (e.g., "10 to 15 August")
	if match[2] != "" {
		endDay := ParseOrdinalNumber(match[2])

		// Create result with end date
		endComponents := components.Clone()
		endComponents.Assign(kronos.ComponentDay, endDay)

		result := context.CreateParsingResult(0, match[0], components, endComponents)
		return result
	}

	return components
}
