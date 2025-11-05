package en

import (
	"regexp"
	"strconv"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// atoiSafe converts a string to an integer, returning 0 and false on error.
// Returns the value and true on success.
func atoiSafe(s string) (int, bool) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return val, true
}

// ENCompactFormatParser handles dates/times without separators
// Supports formats like:
// - 20200315 (YYYYMMDD)
// - 200315 (YYMMDD)
// - 20200315143045 (YYYYMMDDHHmmss)
// - 1430 (HHmm)
// - 143045 (HHmmss)
type ENCompactFormatParser struct {
	*common.AbstractParserWithWordBoundary
}

// NewENCompactFormatParser creates a new compact format parser
func NewENCompactFormatParser() *ENCompactFormatParser {
	parser := &ENCompactFormatParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
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

func (p *ENCompactFormatParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	// Match sequences of 4-14 digits
	// We want to match compact date/time formats but not:
	// - Single/double/triple digits
	// - Very long numbers (15+ digits like unix timestamps)
	// - Numbers followed by hyphens (like "2012-" in "2012-14")
	// Capture the trailing character to properly handle word boundaries
	return regexp.MustCompile(`(\d{4}|\d{6}|\d{8}|\d{10}|\d{12}|\d{14})([^\w-]|$)`)
}

func (p *ENCompactFormatParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
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

	var components *kronos.ParsingComponents

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
	return &kronos.ParsingResultWithBoundary{
		Components:         components,
		AdjustedText:       adjustedText,
		BoundaryLen:        0,    // Will be set by AbstractParserWithWordBoundary
		IncludeBoundaryIdx: true, // Will be overridden by AbstractParserWithWordBoundary
	}
}

// tryParse4Digits attempts to parse 4-digit strings
// Could be: MMDD or HHmm
// Prioritize time if hours >= 13 (clearly not a month)
// Otherwise prioritize date if month is valid
func (p *ENCompactFormatParser) tryParse4Digits(s string, ctx *kronos.ParsingContext) *kronos.ParsingComponents {
	first2, ok := atoiSafe(s[0:2])
	if !ok {
		return nil
	}
	last2, ok := atoiSafe(s[2:4])
	if !ok {
		return nil
	}

	// If first 2 digits > 12, it can only be time (hours)
	if first2 > 12 {
		if isValidTime(first2, last2, 0) {
			components := ctx.CreateParsingComponents(nil)
			components.Assign(kronos.ComponentHour, first2)
			components.Assign(kronos.ComponentMinute, last2)
			components.AddTag("parser/ENCompactFormatParser/time")
			return components
		}
		return nil
	}

	// If first 2 digits <= 12, prefer date interpretation
	// Try as date (MMDD)
	if isValidMonthDay(first2, last2) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentMonth, first2)
		components.Assign(kronos.ComponentDay, last2)
		components.AddTag("parser/ENCompactFormatParser/date")
		return components
	}

	// Fall back to time if date doesn't work
	if isValidTime(first2, last2, 0) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentHour, first2)
		components.Assign(kronos.ComponentMinute, last2)
		components.AddTag("parser/ENCompactFormatParser/time")
		return components
	}

	return nil
}

// tryParse6Digits attempts to parse 6-digit strings
// Could be: YYMMDD or HHmmss
// Prioritize date if first 2 digits look like a year (> 23 or starts with 0)
// Otherwise prioritize time
func (p *ENCompactFormatParser) tryParse6Digits(s string, ctx *kronos.ParsingContext) *kronos.ParsingComponents {
	first2, ok := atoiSafe(s[0:2])
	if !ok {
		return nil
	}
	middle2, ok := atoiSafe(s[2:4])
	if !ok {
		return nil
	}
	last2, ok := atoiSafe(s[4:6])
	if !ok {
		return nil
	}

	// If first 2 digits > 23, it's more likely a year (YYMMDD)
	// Also check if middle2 looks like a valid month
	if first2 > 23 || (middle2 >= 1 && middle2 <= 12 && first2 <= 69) {
		fullYear := convertTwoDigitYear(first2)
		if isValidDate(fullYear, middle2, last2) {
			components := ctx.CreateParsingComponents(nil)
			components.Assign(kronos.ComponentYear, fullYear)
			components.Assign(kronos.ComponentMonth, middle2)
			components.Assign(kronos.ComponentDay, last2)
			components.AddTag("parser/ENCompactFormatParser/date")
			return components
		}
	}

	// Try as time (HHmmss)
	if isValidTime(first2, middle2, last2) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentHour, first2)
		components.Assign(kronos.ComponentMinute, middle2)
		components.Assign(kronos.ComponentSecond, last2)
		components.AddTag("parser/ENCompactFormatParser/time")
		return components
	}

	// Try as date (YYMMDD) again if time failed
	fullYear := convertTwoDigitYear(first2)
	if isValidDate(fullYear, middle2, last2) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentYear, fullYear)
		components.Assign(kronos.ComponentMonth, middle2)
		components.Assign(kronos.ComponentDay, last2)
		components.AddTag("parser/ENCompactFormatParser/date")
		return components
	}

	return nil
}

// tryParse8Digits attempts to parse 8-digit strings
// Format: YYYYMMDD
func (p *ENCompactFormatParser) tryParse8Digits(s string, ctx *kronos.ParsingContext) *kronos.ParsingComponents {
	year, ok := atoiSafe(s[0:4])
	if !ok {
		return nil
	}
	month, ok := atoiSafe(s[4:6])
	if !ok {
		return nil
	}
	day, ok := atoiSafe(s[6:8])
	if !ok {
		return nil
	}

	if isValidDate(year, month, day) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentYear, year)
		components.Assign(kronos.ComponentMonth, month)
		components.Assign(kronos.ComponentDay, day)
		components.AddTag("parser/ENCompactFormatParser/date")
		return components
	}

	return nil
}

// tryParse10Digits attempts to parse 10-digit strings
// Format: YYYYMMDDHHmm
func (p *ENCompactFormatParser) tryParse10Digits(s string, ctx *kronos.ParsingContext) *kronos.ParsingComponents {
	year, ok := atoiSafe(s[0:4])
	if !ok {
		return nil
	}
	month, ok := atoiSafe(s[4:6])
	if !ok {
		return nil
	}
	day, ok := atoiSafe(s[6:8])
	if !ok {
		return nil
	}
	hour, ok := atoiSafe(s[8:10])
	if !ok {
		return nil
	}
	minute := 0

	if isValidDate(year, month, day) && isValidTime(hour, minute, 0) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentYear, year)
		components.Assign(kronos.ComponentMonth, month)
		components.Assign(kronos.ComponentDay, day)
		components.Assign(kronos.ComponentHour, hour)
		components.Imply(kronos.ComponentMinute, minute)
		components.AddTag("parser/ENCompactFormatParser/datetime")
		return components
	}

	return nil
}

// tryParse12Digits attempts to parse 12-digit strings
// Format: YYYYMMDDHHmmss
func (p *ENCompactFormatParser) tryParse12Digits(s string, ctx *kronos.ParsingContext) *kronos.ParsingComponents {
	year, ok := atoiSafe(s[0:4])
	if !ok {
		return nil
	}
	month, ok := atoiSafe(s[4:6])
	if !ok {
		return nil
	}
	day, ok := atoiSafe(s[6:8])
	if !ok {
		return nil
	}
	hour, ok := atoiSafe(s[8:10])
	if !ok {
		return nil
	}
	minute, ok := atoiSafe(s[10:12])
	if !ok {
		return nil
	}
	second := 0

	if isValidDate(year, month, day) && isValidTime(hour, minute, second) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentYear, year)
		components.Assign(kronos.ComponentMonth, month)
		components.Assign(kronos.ComponentDay, day)
		components.Assign(kronos.ComponentHour, hour)
		components.Assign(kronos.ComponentMinute, minute)
		components.Imply(kronos.ComponentSecond, second)
		components.AddTag("parser/ENCompactFormatParser/datetime")
		return components
	}

	return nil
}

// tryParse14Digits attempts to parse 14-digit strings
// Format: YYYYMMDDHHmmss
func (p *ENCompactFormatParser) tryParse14Digits(s string, ctx *kronos.ParsingContext) *kronos.ParsingComponents {
	year, ok := atoiSafe(s[0:4])
	if !ok {
		return nil
	}
	month, ok := atoiSafe(s[4:6])
	if !ok {
		return nil
	}
	day, ok := atoiSafe(s[6:8])
	if !ok {
		return nil
	}
	hour, ok := atoiSafe(s[8:10])
	if !ok {
		return nil
	}
	minute, ok := atoiSafe(s[10:12])
	if !ok {
		return nil
	}
	second, ok := atoiSafe(s[12:14])
	if !ok {
		return nil
	}

	if isValidDate(year, month, day) && isValidTime(hour, minute, second) {
		components := ctx.CreateParsingComponents(nil)
		components.Assign(kronos.ComponentYear, year)
		components.Assign(kronos.ComponentMonth, month)
		components.Assign(kronos.ComponentDay, day)
		components.Assign(kronos.ComponentHour, hour)
		components.Assign(kronos.ComponentMinute, minute)
		components.Assign(kronos.ComponentSecond, second)
		components.AddTag("parser/ENCompactFormatParser/datetime")
		return components
	}

	return nil
}

// isValidDate checks if year, month, day form a valid date
func isValidDate(year, month, day int) bool {
	if year < 1000 || year > 9999 {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}
	if day < 1 || day > 31 {
		return false
	}

	// Check if the date is actually valid (e.g., no Feb 30)
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.Year() == year && int(t.Month()) == month && t.Day() == day
}

// isValidMonthDay checks if month and day form a valid date
func isValidMonthDay(month, day int) bool {
	if month < 1 || month > 12 {
		return false
	}
	if day < 1 || day > 31 {
		return false
	}

	// Use a leap year to validate day ranges (2020 was a leap year)
	t := time.Date(2020, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return int(t.Month()) == month && t.Day() == day
}

// isValidTime checks if hour, minute, second form a valid time
func isValidTime(hour, minute, second int) bool {
	return hour >= 0 && hour <= 23 &&
		minute >= 0 && minute <= 59 &&
		second >= 0 && second <= 59
}

// convertTwoDigitYear converts a 2-digit year to a 4-digit year
// Rule: 00-69 → 2000-2069, 70-99 → 1970-1999
func convertTwoDigitYear(year int) int {
	if year >= 0 && year <= 69 {
		return 2000 + year
	}
	return 1900 + year
}
