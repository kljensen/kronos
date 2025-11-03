package en

import (
	"regexp"
	"strconv"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENSlashMonthFormatParser parses MM/YYYY format
// Examples: 11/2005, 06/2005
type ENSlashMonthFormatParser struct {
	*common.AbstractParserWithWordBoundary
}

// NewENSlashMonthFormatParser creates a new ENSlashMonthFormatParser
func NewENSlashMonthFormatParser() *ENSlashMonthFormatParser {
	parser := &ENSlashMonthFormatParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENSlashMonthFormatParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Pattern: MM/YYYY (month 1-12, year 4 digits)
	pattern := `(?i)([0-9]|0[1-9]|1[012])/([0-9]{4})`
	return regexp.MustCompile(pattern)
}

func (p *ENSlashMonthFormatParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 3 {
		return nil
	}

	month, _ := strconv.Atoi(match[1])
	year, _ := strconv.Atoi(match[2])

	return context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Assign(kronos.ComponentYear, year).
		Imply(kronos.ComponentDay, 1).
		SetPeriod(kronos.PeriodMonth)
}
