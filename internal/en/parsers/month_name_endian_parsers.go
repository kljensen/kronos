//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/en/data"
	"github.com/kljensen/kronos/internal/helpers"
	"github.com/kljensen/kronos/internal/parsing"
)

// ENMonthNameLittleEndianParser parses "DD Month YYYY" format (little-endian)
// Examples: "10 August 2012", "3rd Feb 82", "31st March, 2016"
// Also handles ranges: "10 to 15 August 2012"
type ENMonthNameLittleEndianParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENMonthNameLittleEndianParser creates a new ENMonthNameLittleEndianParser
func NewENMonthNameLittleEndianParser() *ENMonthNameLittleEndianParser {
	parser := &ENMonthNameLittleEndianParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENMonthNameLittleEndianParser) innerPattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	// Pattern: (on)? DAY (to|-)? DAY? (-|/|of) MONTH (-|/|,)? YEAR?
	pattern := `(?i)(?:on\s{0,3})?` +
		`(` + data.OrdinalNumberPattern + `)` +
		`(?:` +
		`\s{0,3}(?:to|\-|until|through|till)?\s{0,3}` +
		`(` + data.OrdinalNumberPattern + `)` +
		`)?` +
		`(?:-|/|\s{0,3}(?:of)?\s{0,3})` +
		`(` + data.MonthPattern + `)` +
		`(?:` +
		`(?:-|/|,?\s{0,3})` +
		`(` + data.YearPattern + `)` +
		`)?` +
		``

	return regexp.MustCompile(pattern)
}

func (p *ENMonthNameLittleEndianParser) innerExtract(context *kronos.InternalParsingContext, match []string) interface{} {
	if len(match) < 4 {
		return nil
	}

	day := data.ParseOrdinalNumber(match[1])

	// Validate day - if > 31, it might be a year-like number (e.g., "96 Aug")
	if day > 31 {
		return nil
	}

	monthName := strings.ToLower(match[3])
	month := data.MonthDictionary[monthName]

	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Assign(kronos.ComponentDay, day).
		AddTag("parser/ENMonthNameLittleEndianParser").
		SetPeriod(kronos.PeriodDay) // Day-level period for specific dates

	// Handle year if present
	if len(match) > 4 && match[4] != "" {
		year := data.ParseYear(match[4])
		components.Assign(kronos.ComponentYear, year)
	} else {
		// Find closest year to reference using preference setting
		year := helpers.FindYearClosestToRefWithPreference(context.RefDate(), day, month, context.Option().PreferDatesFrom)
		components.Imply(kronos.ComponentYear, year)
	}

	// Handle date range (e.g., "10 to 15 August")
	if match[2] != "" {
		endDay := data.ParseOrdinalNumber(match[2])

		// Create result with end date
		endComponents := components.Clone()
		endComponents.Assign(kronos.ComponentDay, endDay)

		result := context.CreateParsingResult(0, match[0], components, endComponents)
		return result
	}

	return components
}

// ENMonthNameMiddleEndianParser parses "Month DD, YYYY" format (middle-endian, US-style)
// Examples: "August 10, 2012", "January 13", "Dec 12, 2020"
// Also handles ranges: "January 12 - 15, 2012"
type ENMonthNameMiddleEndianParser struct {
	*parsing.AbstractParserWithWordBoundary

	shouldSkipYearLikeDate bool
}

// NewENMonthNameMiddleEndianParser creates a new ENMonthNameMiddleEndianParser
func NewENMonthNameMiddleEndianParser(shouldSkipYearLikeDate bool) *ENMonthNameMiddleEndianParser {
	parser := &ENMonthNameMiddleEndianParser{
		shouldSkipYearLikeDate: shouldSkipYearLikeDate,
	}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENMonthNameMiddleEndianParser) innerPattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	// Pattern: MONTH (-|/|,) DAY (to|-)? (MONTH)? DAY? (-|/|,|space)? YEAR?
	// Supports:
	// - "August 10" (month + day)
	// - "August 10, 2012" (month + day + year with comma)
	// - "August 10 2012" (month + day + year with space)
	// - "August 10 2555 BE" (month + day + year with suffix)
	// - "August 10 - 22, 2012" (month + day + range within same month + year)
	// - "August 10 - November 12" (month + day + range to different month + day)
	pattern := `(?i)(` + data.MonthPattern + `)` +
		`(?:-|/|\s*,?\s*)` +
		`(` + data.OrdinalNumberPattern + `)` +
		`(?:` +
		`\s*(?:to|\-)\s*` +
		`(?:(` + data.MonthPattern + `)(?:-|/|\s*,?\s*))?` + // optional second month for ranges
		`(` + data.OrdinalNumberPattern + `)` +
		`)?` +
		`(?:` +
		`(?:\s*(?:-|/|,)\s*|\s+)` + // punctuation separator or space
		`(` + data.YearPattern + `)` +
		`)?`

	return regexp.MustCompile(pattern)
}

func (p *ENMonthNameMiddleEndianParser) innerExtract(context *kronos.InternalParsingContext, match []string) interface{} {
	if len(match) < 3 {
		return nil
	}

	// Capture groups after boundary adjustment:
	// match[0] = full match text (without boundary)
	// match[1] = first month
	// match[2] = first day
	// match[3] = second month (optional, for ranges like "Aug 10 - Nov 12")
	// match[4] = second day (optional, for ranges)
	// match[5] = year (optional)

	monthName := strings.ToLower(match[1])
	month := data.MonthDictionary[monthName]

	day := data.ParseOrdinalNumber(match[2])

	// Validate day
	if day > 31 {
		return nil
	}

	// Find where the day portion ends in the original text
	// This helps us check if the day is part of a longer number (like "20" from "2012")
	monthPos := strings.Index(strings.ToLower(context.Text()), strings.ToLower(match[1]))
	if monthPos >= 0 {
		textAfterMonth := context.Text()[monthPos+len(match[1]):]
		dayPattern := regexp.MustCompile(`(?i)(?:-|/|\s*,?\s*)(` + regexp.QuoteMeta(match[2]) + `)`)
		dayMatch := dayPattern.FindStringSubmatchIndex(textAfterMonth)
		if len(dayMatch) >= 4 {
			// dayMatch[3] is the end of the captured day group
			dayEnd := dayMatch[3]
			if dayEnd < len(textAfterMonth) && textAfterMonth[dayEnd] >= '0' && textAfterMonth[dayEnd] <= '9' {
				// Day is followed by another digit, likely part of a year (e.g., "20" from "2012")
				return nil
			}
		}
	}

	// Reject if the match is followed by a colon (time indicator like ":00")
	// This prevents false matches like "May 12:00" being parsed as "May 12"
	matchEnd := strings.Index(context.Text(), match[0]) + len(match[0])
	if matchEnd < len(context.Text()) && context.Text()[matchEnd] == ':' {
		return nil
	}

	// Skip year-like dates if configured (e.g., "January 21" where 21 looks like a year)
	if p.shouldSkipYearLikeDate {
		// No range, no year, and day looks like a year (20-25)
		if (len(match) < 5 || match[4] == "") && (len(match) < 6 || match[5] == "") {
			yearLikePattern := regexp.MustCompile(`^2[0-5]$`)
			if yearLikePattern.MatchString(match[2]) {
				return nil
			}
		}
	}

	components := context.CreateParsingComponents(nil).
		Assign(kronos.ComponentMonth, month).
		Assign(kronos.ComponentDay, day).
		AddTag("parser/ENMonthNameMiddleEndianParser").
		SetPeriod(kronos.PeriodDay)

	// Handle year if present
	if len(match) > 5 && match[5] != "" {
		year := data.ParseYear(match[5])
		components.Assign(kronos.ComponentYear, year)
	} else {
		// Find closest year to reference using preference setting
		year := helpers.FindYearClosestToRefWithPreference(context.RefDate(), day, month, context.Option().PreferDatesFrom)
		components.Imply(kronos.ComponentYear, year)
	}

	// Handle date range (e.g., "August 10 - 22" or "August 10 - November 12")
	if len(match) > 4 && match[4] != "" {
		endDay := data.ParseOrdinalNumber(match[4])
		endMonth := month // default to same month

		// Check if there's a different month for the end date
		if match[3] != "" {
			endMonthName := strings.ToLower(match[3])
			endMonth = data.MonthDictionary[endMonthName]
		}

		// Create result with end date
		endComponents := components.Clone()
		endComponents.Assign(kronos.ComponentDay, endDay)
		endComponents.Assign(kronos.ComponentMonth, endMonth)

		result := context.CreateParsingResult(0, match[0], components, endComponents)
		return result
	}

	// Return a ParsingResult instead of bare ParsingComponents
	// This is necessary for AbstractParserWithWordBoundary to properly handle the boundary offset
	result := context.CreateParsingResult(0, match[0], components)
	return result
}
