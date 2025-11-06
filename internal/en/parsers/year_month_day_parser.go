package parsers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/parsing"
	"github.com/kljensen/kronos/internal/en/data"
)

// ENYearMonthDayParser parses date formats like YYYY-MM-DD, YYYY/MM/DD, YYYY.MM.DD
// Supports both numeric months and month names (e.g., 2012/Aug/10)
type ENYearMonthDayParser struct {
	*parsing.AbstractParserWithWordBoundary
	strictMonthDateOrder bool
}

// NewENYearMonthDayParser creates a new ENYearMonthDayParser
func NewENYearMonthDayParser(strictMonthDateOrder bool) *ENYearMonthDayParser {
	parser := &ENYearMonthDayParser{
		strictMonthDateOrder: strictMonthDateOrder,
	}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENYearMonthDayParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Pattern: YYYY[-/./ ]MM[-/./ ]DD or YYYY[-/./ ]MONTH[-/./ ]DD
	pattern := `(?i)([0-9]{4})[-.\\/\s]` +
		`(?:(` + data.MonthPattern + `)|([0-9]{1,2}))[-.\\/\s]` +
		`([0-9]{1,2})` +
		``

	return regexp.MustCompile(pattern)
}

func (p *ENYearMonthDayParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 5 {
		return nil
	}

	year, err := strconv.Atoi(match[1])
	if err != nil {
		return nil
	}

	var month, day int

	// Check if month name was matched (group 2) or numeric month (group 3)
	if match[2] != "" {
		// Month name
		month = data.MonthDictionary[strings.ToLower(match[2])]
	} else {
		// Numeric month
		month, err = strconv.Atoi(match[3])
		if err != nil {
			return nil
		}
	}

	day, err = strconv.Atoi(match[4])
	if err != nil {
		return nil
	}

	// Validate and potentially swap month/day if not in strict mode
	if month < 1 || month > 12 {
		if p.strictMonthDateOrder {
			return nil
		}
		// Try swapping if day value could be a valid month
		if day >= 1 && day <= 12 {
			month, day = day, month
		}
	}

	// Validate day
	if day < 1 || day > 31 {
		return nil
	}

	// Validate that the date is actually possible
	// (e.g., reject Feb 30, April 31, etc.)
	// Use Go's time package to validate the date
	// time.Date normalizes invalid dates, so we check if it changed
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return nil
	}

	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentYear, year).
		Assign(kronos.ComponentMonth, month).
		Assign(kronos.ComponentDay, day).
		SetPeriod(kronos.PeriodDay)

	return components
}
