//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/parsing"
	"github.com/kljensen/kronos/internal/en/data"
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

func (p *ENMonthNameParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Pattern: (in)? MONTH (,|-|of)? YEAR?
	pattern := `(?i)((?:in)\s*)?` +
		`(` + data.MonthPattern + `)` +
		`(?:` +
		`(?:\s*(?:,|-|of))?\s*(` + data.YearPattern + `)` +
		`)?`

	return regexp.MustCompile(pattern)
}

func (p *ENMonthNameParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
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
		year := kronos.XFindYearClosestToRefWithPreference(context.RefDate(), 1, month, context.Option().Preference)
		components.Imply(kronos.ComponentYear, year)
	}

	// Set period to month-level
	components.SetPeriod(kronos.PeriodMonth)

	return components
}
