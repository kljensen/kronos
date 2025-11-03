package en

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
)

// YEAR_PATTERN matches year patterns including BE, AD, BC, BCE, CE suffixes
// Note: Longer alternatives (BCE) must come before shorter ones (BC, AD)
const YEAR_PATTERN = `(?:[1-9][0-9]{0,3}\s{0,2}(?:BCE|BE|CE|AD|BC)|[1-2][0-9]{3}|[5-9][0-9]|2[0-5])`

// MonthDictionary maps month names to month numbers (1-12)
var MonthDictionary = map[string]int{
	// Full names
	"january":   1,
	"february":  2,
	"march":     3,
	"april":     4,
	"may":       5,
	"june":      6,
	"july":      7,
	"august":    8,
	"september": 9,
	"october":   10,
	"november":  11,
	"december":  12,
	// Abbreviated names
	"jan":   1,
	"jan.":  1,
	"feb":   2,
	"feb.":  2,
	"mar":   3,
	"mar.":  3,
	"apr":   4,
	"apr.":  4,
	"jun":   6,
	"jun.":  6,
	"jul":   7,
	"jul.":  7,
	"aug":   8,
	"aug.":  8,
	"sep":   9,
	"sep.":  9,
	"sept":  9,
	"sept.": 9,
	"oct":   10,
	"oct.":  10,
	"nov":   11,
	"nov.":  11,
	"dec":   12,
	"dec.":  12,
}

// FullMonthNameDictionary contains only full month names
var FullMonthNameDictionary = map[string]int{
	"january":   1,
	"february":  2,
	"march":     3,
	"april":     4,
	"may":       5,
	"june":      6,
	"july":      7,
	"august":    8,
	"september": 9,
	"october":   10,
	"november":  11,
	"december":  12,
}

// WeekdayDictionary maps weekday names to weekday numbers (0-6)
var WeekdayDictionary = map[string]kronos.Weekday{
	"sunday":    kronos.WeekdaySunday,
	"sun":       kronos.WeekdaySunday,
	"sun.":      kronos.WeekdaySunday,
	"monday":    kronos.WeekdayMonday,
	"mon":       kronos.WeekdayMonday,
	"mon.":      kronos.WeekdayMonday,
	"tuesday":   kronos.WeekdayTuesday,
	"tue":       kronos.WeekdayTuesday,
	"tue.":      kronos.WeekdayTuesday,
	"wednesday": kronos.WeekdayWednesday,
	"wed":       kronos.WeekdayWednesday,
	"wed.":      kronos.WeekdayWednesday,
	"thursday":  kronos.WeekdayThursday,
	"thurs":     kronos.WeekdayThursday,
	"thurs.":    kronos.WeekdayThursday,
	"thur":      kronos.WeekdayThursday,
	"thur.":     kronos.WeekdayThursday,
	"thu":       kronos.WeekdayThursday,
	"thu.":      kronos.WeekdayThursday,
	"friday":    kronos.WeekdayFriday,
	"fri":       kronos.WeekdayFriday,
	"fri.":      kronos.WeekdayFriday,
	"saturday":  kronos.WeekdaySaturday,
	"sat":       kronos.WeekdaySaturday,
	"sat.":      kronos.WeekdaySaturday,
}

// IntegerWordDictionary maps word numbers to integers
var IntegerWordDictionary = map[string]int{
	"one":    1,
	"two":    2,
	"three":  3,
	"four":   4,
	"five":   5,
	"six":    6,
	"seven":  7,
	"eight":  8,
	"nine":   9,
	"ten":    10,
	"eleven": 11,
	"twelve": 12,
}

// NumberWordDictionary maps all number words to their float values
// This includes regular numbers, informal quantifiers, and fractional values
var NumberWordDictionary = map[string]float64{
	// Regular number words (1-12)
	"one":    1.0,
	"two":    2.0,
	"three":  3.0,
	"four":   4.0,
	"five":   5.0,
	"six":    6.0,
	"seven":  7.0,
	"eight":  8.0,
	"nine":   9.0,
	"ten":    10.0,
	"eleven": 11.0,
	"twelve": 12.0,
	// Articles that indicate singular
	"a":  1.0,
	"an": 1.0,
	// Special quantifiers
	"couple":  2.0, // "a couple" typically means 2
	"few":     3.0, // "a few" typically means 3-5, using 3 as standard
	"several": 7.0, // "several" typically means 3-7, using 7 for backward compatibility
	"dozen":   12.0,
	// Fractional values
	"half": 0.5,
}

// OrdinalWordDictionary maps ordinal words to integers
var OrdinalWordDictionary = map[string]int{
	"first":          1,
	"second":         2,
	"third":          3,
	"fourth":         4,
	"fifth":          5,
	"sixth":          6,
	"seventh":        7,
	"eighth":         8,
	"ninth":          9,
	"tenth":          10,
	"eleventh":       11,
	"twelfth":        12,
	"thirteenth":     13,
	"fourteenth":     14,
	"fifteenth":      15,
	"sixteenth":      16,
	"seventeenth":    17,
	"eighteenth":     18,
	"nineteenth":     19,
	"twentieth":      20,
	"twenty first":   21,
	"twenty-first":   21,
	"twenty second":  22,
	"twenty-second":  22,
	"twenty third":   23,
	"twenty-third":   23,
	"twenty fourth":  24,
	"twenty-fourth":  24,
	"twenty fifth":   25,
	"twenty-fifth":   25,
	"twenty sixth":   26,
	"twenty-sixth":   26,
	"twenty seventh": 27,
	"twenty-seventh": 27,
	"twenty eighth":  28,
	"twenty-eighth":  28,
	"twenty ninth":   29,
	"twenty-ninth":   29,
	"thirtieth":      30,
	"thirty first":   31,
	"thirty-first":   31,
}

// TimeUnitDictionary maps time unit words to Timeunit values
var TimeUnitDictionary = map[string]kronos.Timeunit{
	"s":        kronos.TimeunitSecond,
	"sec":      kronos.TimeunitSecond,
	"second":   kronos.TimeunitSecond,
	"seconds":  kronos.TimeunitSecond,
	"m":        kronos.TimeunitMinute,
	"min":      kronos.TimeunitMinute,
	"mins":     kronos.TimeunitMinute,
	"minute":   kronos.TimeunitMinute,
	"minutes":  kronos.TimeunitMinute,
	"h":        kronos.TimeunitHour,
	"hr":       kronos.TimeunitHour,
	"hrs":      kronos.TimeunitHour,
	"hour":     kronos.TimeunitHour,
	"hours":    kronos.TimeunitHour,
	"d":        kronos.TimeunitDay,
	"day":      kronos.TimeunitDay,
	"days":     kronos.TimeunitDay,
	"w":        kronos.TimeunitWeek,
	"week":     kronos.TimeunitWeek,
	"weeks":    kronos.TimeunitWeek,
	"mo":       kronos.TimeunitMonth,
	"mon":      kronos.TimeunitMonth,
	"mos":      kronos.TimeunitMonth,
	"month":    kronos.TimeunitMonth,
	"months":   kronos.TimeunitMonth,
	"qtr":      kronos.TimeunitQuarter,
	"quarter":  kronos.TimeunitQuarter,
	"quarters": kronos.TimeunitQuarter,
	"y":        kronos.TimeunitYear,
	"yr":       kronos.TimeunitYear,
	"year":     kronos.TimeunitYear,
	"years":    kronos.TimeunitYear,
	"decade":   kronos.TimeunitDecade,
	"decades":  kronos.TimeunitDecade,
}

// Pattern builders

// MatchAnyPattern creates a pattern that matches any key from a map
func MatchAnyPattern(dict interface{}) string {
	var keys []string
	switch d := dict.(type) {
	case map[string]int:
		for k := range d {
			keys = append(keys, regexp.QuoteMeta(k))
		}
	case map[string]kronos.Weekday:
		for k := range d {
			keys = append(keys, regexp.QuoteMeta(k))
		}
	case map[string]kronos.Timeunit:
		for k := range d {
			keys = append(keys, regexp.QuoteMeta(k))
		}
	}
	// Sort keys for deterministic pattern generation
	// Longer strings first to ensure greedy matching (e.g., "september" before "sep")
	sort.SliceStable(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) > len(keys[j])
		}
		return keys[i] < keys[j]
	})
	return strings.Join(keys, "|")
}

// Patterns
var (
	MonthPattern         = MatchAnyPattern(MonthDictionary)
	FullMonthPattern     = MatchAnyPattern(FullMonthNameDictionary)
	WeekdayPattern       = MatchAnyPattern(WeekdayDictionary)
	IntegerWordPattern   = MatchAnyPattern(IntegerWordDictionary)
	OrdinalWordPattern   = MatchAnyPattern(OrdinalWordDictionary)
	OrdinalNumberPattern = `(?:` + OrdinalWordPattern + `|[0-9]{1,2}(?:st|nd|rd|th)?)`
	YearPattern          = YEAR_PATTERN
)

// ParseOrdinalNumber parses an ordinal number pattern (e.g., "1st", "2nd", "first", "second")
func ParseOrdinalNumber(match string) int {
	lower := strings.ToLower(match)

	// Check if it's a word
	if val, ok := OrdinalWordDictionary[lower]; ok {
		return val
	}

	// Remove ordinal suffix (st, nd, rd, th)
	cleaned := regexp.MustCompile(`(?i)(st|nd|rd|th)$`).ReplaceAllString(lower, "")
	val, _ := strconv.Atoi(cleaned)
	return val
}

// ParseYear parses a year pattern (handles BE, AD, BC, BCE, CE)
func ParseYear(match string) int {
	// Buddhist Era
	if regexp.MustCompile(`(?i)BE`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*BE`).ReplaceAllString(match, "")
		year, _ := strconv.Atoi(strings.TrimSpace(cleaned))
		return year - 543
	}

	// Before Christ / Before Common Era
	if regexp.MustCompile(`(?i)BCE?`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*BCE?`).ReplaceAllString(match, "")
		year, _ := strconv.Atoi(strings.TrimSpace(cleaned))
		return -year
	}

	// Anno Domini / Common Era
	if regexp.MustCompile(`(?i)(AD|CE)`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*(AD|CE)`).ReplaceAllString(match, "")
		year, _ := strconv.Atoi(strings.TrimSpace(cleaned))
		return year
	}

	// Regular year number
	year, _ := strconv.Atoi(strings.TrimSpace(match))
	return kronos.FindMostLikelyADYear(year)
}

// ParseNumberPattern parses number-like patterns including words
func ParseNumberPattern(match string) float64 {
	lower := strings.ToLower(strings.TrimSpace(match))

	// Check number word dictionary (includes all word numbers and quantifiers)
	if val, ok := NumberWordDictionary[lower]; ok {
		return val
	}

	// Special handling for "the" - only treat as 1 in certain contexts
	// For now, skip "the" as it's not typically used as a quantity
	if lower == "the" {
		return 1
	}

	// Normalize comma to dot for decimal separator (European format support)
	normalized := strings.ReplaceAll(lower, ",", ".")

	// Try parsing as number
	val, _ := strconv.ParseFloat(normalized, 64)
	return val
}

// Time-related patterns

var (
	// NumberPattern matches integers and float numbers
	NumberPattern = `(?:\d+(?:\.\d*)?|\.\d+)`

	// TimeUnitNoAbbrPattern matches full time unit names only (no abbreviations)
	TimeUnitNoAbbrPattern = `(?:decades|decade|seconds|second|minutes|minute|hours|hour|days|day|weeks|week|months|month|quarters|quarter|years|year)`

	// TimeUnitPattern matches time units including abbreviations
	// Order matters: longer alternatives come first to ensure correct matching
	TimeUnitPattern = `(?:decades|decade|seconds|second|minutes|minute|hours|hour|days|day|weeks|week|months|month|quarters|quarter|years|year|mins|min|hrs|hr|sec|mos|mon|qtr|yr|mo|s|m|h|d|w|y)`
)

// ParseDuration parses a duration expression like "3 days", "5 hours 30 minutes", "half an hour"
func ParseDuration(text string) kronos.Duration {
	result := make(kronos.Duration)

	// Pattern for matching time units
	// Supports formats like:
	// - "3 days"
	// - "two weeks"
	// - "5 hours 30 minutes"
	// - "a week"
	// - "half an hour"
	// - "1d 5h 30m"
	// - "2.5 hours" (decimal with dot)
	// - "2,5 hours" (decimal with comma, European format)
	// - "a couple of days" -> captures "couple"
	// - "a few hours" -> captures "few"
	// - "several weeks" -> captures "several"
	// Word numbers: one, two, three, four, five, six, seven, eight, nine, ten, eleven, twelve
	// Extended quantifiers: a, an, couple, few, several, half, dozen
	// Note: Using flexible pattern to support consecutive units like "2hr5min"
	// Also supports optional "of" separator (e.g., "couple of days")
	// Pattern breakdown:
	// 1. Optional leading "a/an" (not captured)
	// 2. Number/quantifier word (captured)
	// 3. Optional "a/an" again for "half an hour"
	// 4. Optional "of"
	// 5. Time unit (captured)
	pattern := regexp.MustCompile(`(?i)(?:an?\s+)?(half|dozen|several|couple|few|one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|an|a|the|[0-9]+(?:[.,][0-9]+)?)\s*(?:an?\s+)?(?:of\s+)?(` + TimeUnitPattern + `)`)
	matchesIdx := pattern.FindAllStringSubmatchIndex(text, -1)

	for _, matchIdx := range matchesIdx {
		// matchIdx has: [fullStart, fullEnd, numStart, numEnd, unitStart, unitEnd]
		if len(matchIdx) < 6 {
			continue
		}

		fullStart := matchIdx[0]

		// Get the unit match (group 2)
		unitStart, unitEnd := matchIdx[4], matchIdx[5]
		unit := strings.ToLower(text[unitStart:unitEnd])

		// Skip if match is preceded by a letter (part of a larger word like "them")
		// But allow if the match starts with a digit (e.g., "2hr5min" where "5min" follows "hr")
		if fullStart > 0 {
			prevChar := text[fullStart-1]
			matchStartsWithDigit := (matchIdx[2] >= 0 && matchIdx[2] == fullStart &&
				len(text[matchIdx[2]:matchIdx[3]]) > 0 &&
				text[matchIdx[2]] >= '0' && text[matchIdx[2]] <= '9')

			if !matchStartsWithDigit &&
				((prevChar >= 'a' && prevChar <= 'z') || (prevChar >= 'A' && prevChar <= 'Z')) {
				continue
			}
		}

		// Skip if the unit is immediately followed by a letter (part of a larger word)
		if unitEnd < len(text) {
			nextChar := text[unitEnd]
			if (nextChar >= 'a' && nextChar <= 'z') || (nextChar >= 'A' && nextChar <= 'Z') {
				continue
			}
		}

		// Get the number part (group 1, if any)
		var numStr string
		if matchIdx[2] >= 0 {
			numStr = strings.TrimSpace(text[matchIdx[2]:matchIdx[3]])
		}

		// Single-letter time units (s, m, h, d, w, y) should have an explicit number
		// This prevents matching "am" as "a" + "m" or "them" as "the" + "m"
		if len(unit) == 1 {
			// Only allow if we have an actual digit
			if numStr == "" || (numStr != "" && !strings.ContainsAny(numStr, "0123456789")) {
				continue
			}
		}

		var num float64
		if numStr == "" {
			num = 1
		} else {
			num = ParseNumberPattern(numStr)
		}

		// Map to Timeunit and add to duration
		if timeunit, ok := TimeUnitDictionary[unit]; ok {
			// Accumulate values for same timeunit
			if existing, exists := result[timeunit]; exists {
				result[timeunit] = existing + num
			} else {
				result[timeunit] = num
			}
		}
	}

	return result
}

// IsEmpty returns true if the duration has no non-zero values
func IsEmptyDuration(d kronos.Duration) bool {
	if len(d) == 0 {
		return true
	}
	for _, v := range d {
		if v != 0 {
			return false
		}
	}
	return true
}
