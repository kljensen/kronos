//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// tryParse4Digits attempts to parse 4-digit strings
// Could be: MMDD or HHmm
// Prioritize time if hours >= 13 (clearly not a month)
// Otherwise prioritize date if month is valid
func (p *ENCompactFormatParser) tryParse4Digits(s string, ctx *kronos.InternalParsingContext) *kronos.InternalParsingComponents {
	first2, ok := helpers.AtoiSafe(s[0:2])
	if !ok {
		return nil
	}
	last2, ok := helpers.AtoiSafe(s[2:4])
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
func (p *ENCompactFormatParser) tryParse6Digits(s string, ctx *kronos.InternalParsingContext) *kronos.InternalParsingComponents {
	first2, ok := helpers.AtoiSafe(s[0:2])
	if !ok {
		return nil
	}
	middle2, ok := helpers.AtoiSafe(s[2:4])
	if !ok {
		return nil
	}
	last2, ok := helpers.AtoiSafe(s[4:6])
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
func (p *ENCompactFormatParser) tryParse8Digits(s string, ctx *kronos.InternalParsingContext) *kronos.InternalParsingComponents {
	year, ok := helpers.AtoiSafe(s[0:4])
	if !ok {
		return nil
	}
	month, ok := helpers.AtoiSafe(s[4:6])
	if !ok {
		return nil
	}
	day, ok := helpers.AtoiSafe(s[6:8])
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

// tryParse10Digits attempts to parse 10-digit strings (YYYYMMDDHHmm)
func (p *ENCompactFormatParser) tryParse10Digits(s string, ctx *kronos.InternalParsingContext) *kronos.InternalParsingComponents {
	return p.tryParseDateTime(s, 4, 6, 8, 10, -1, -1, ctx)
}

// tryParse12Digits attempts to parse 12-digit strings (YYYYMMDDHHmmss)
func (p *ENCompactFormatParser) tryParse12Digits(s string, ctx *kronos.InternalParsingContext) *kronos.InternalParsingComponents {
	return p.tryParseDateTime(s, 4, 6, 8, 10, 12, -1, ctx)
}

// tryParse14Digits attempts to parse 14-digit strings (YYYYMMDDHHmmss)
func (p *ENCompactFormatParser) tryParse14Digits(s string, ctx *kronos.InternalParsingContext) *kronos.InternalParsingComponents {
	return p.tryParseDateTime(s, 4, 6, 8, 10, 12, 14, ctx)
}

// tryParseDateTime is a helper that parses datetime strings with year, month, day, and optional hour, minute, second
// Pass -1 for minute or second positions to imply 0 values
func (p *ENCompactFormatParser) tryParseDateTime(s string, yearLen, monthPos, dayPos, hourPos, minutePos, secondPos int, ctx *kronos.InternalParsingContext) *kronos.InternalParsingComponents {
	year, ok := helpers.AtoiSafe(s[0:yearLen])
	if !ok {
		return nil
	}
	month, ok := helpers.AtoiSafe(s[yearLen:monthPos])
	if !ok {
		return nil
	}
	day, ok := helpers.AtoiSafe(s[monthPos:dayPos])
	if !ok {
		return nil
	}
	hour, ok := helpers.AtoiSafe(s[dayPos:hourPos])
	if !ok {
		return nil
	}

	var minute, second int
	if minutePos > 0 {
		minute, ok = helpers.AtoiSafe(s[hourPos:minutePos])
		if !ok {
			return nil
		}
	}
	if secondPos > 0 {
		second, ok = helpers.AtoiSafe(s[minutePos:secondPos])
		if !ok {
			return nil
		}
	}

	if !isValidDate(year, month, day) || !isValidTime(hour, minute, second) {
		return nil
	}

	components := ctx.CreateParsingComponents(nil)
	components.Assign(kronos.ComponentYear, year)
	components.Assign(kronos.ComponentMonth, month)
	components.Assign(kronos.ComponentDay, day)
	components.Assign(kronos.ComponentHour, hour)

	if minutePos > 0 {
		components.Assign(kronos.ComponentMinute, minute)
	} else {
		components.Imply(kronos.ComponentMinute, 0)
	}

	if secondPos > 0 {
		components.Assign(kronos.ComponentSecond, second)
	} else {
		components.Imply(kronos.ComponentSecond, 0)
	}

	components.AddTag("parser/ENCompactFormatParser/datetime")
	return components
}
