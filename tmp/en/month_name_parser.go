package en

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENMonthNameParser parses standalone month names with optional year
// Examples: "January", "January, 2012", "in June of 2022", "Sep 2012"
type ENMonthNameParser struct {
	*common.AbstractParserWithWordBoundary
}

// NewENMonthNameParser creates a new ENMonthNameParser
func NewENMonthNameParser() *ENMonthNameParser {
	parser := &ENMonthNameParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENMonthNameParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Pattern: (in)? MONTH (,|-|of)? YEAR?
	// Simplified lookahead pattern from: (?=[^\s\w]|\s+[^0-9]|\s+$|$)
	pattern := `(?i)((?:in)\s*)?` +
		`(` + MonthPattern + `)` +
		`\s*` +
		`(?:` +
		`(?:,|-|of)?\s*(` + YearPattern + `)?` +
		`)?` +
		`(?:\s|$|\b)`

	return regexp.MustCompile(pattern)
}

func (p *ENMonthNameParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 3 {
		return nil
	}

	monthName := strings.ToLower(match[2])

	// Skip unlikely short words unless they're full month names
	if len(match[0]) <= 3 {
		if _, ok := FullMonthNameDictionary[monthName]; !ok {
			return nil
		}
	}

	month := MonthDictionary[monthName]

	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Imply(kronos.ComponentDay, 1).
		AddTag("parser/ENMonthNameParser")

	// Handle year if present
	if len(match) > 3 && match[3] != "" {
		year := ParseYear(match[3])
		components.Assign(kronos.ComponentYear, year)
	} else {
		// Find closest year to reference
		year := kronos.FindYearClosestToRef(context.RefDate(), 1, month)
		components.Imply(kronos.ComponentYear, year)
	}

	return components
}
