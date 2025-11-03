package en

import (
	"regexp"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
)

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
	return strings.Join(keys, "|")
}

// Patterns
var (
	MonthPattern        = MatchAnyPattern(MonthDictionary)
	FullMonthPattern    = MatchAnyPattern(FullMonthNameDictionary)
	WeekdayPattern      = MatchAnyPattern(WeekdayDictionary)
	IntegerWordPattern  = MatchAnyPattern(IntegerWordDictionary)
	OrdinalWordPattern  = MatchAnyPattern(OrdinalWordDictionary)
	OrdinalNumberPattern = `(?:` + OrdinalWordPattern + `|[0-9]{1,2}(?:st|nd|rd|th)?)`
	YearPattern         = `(?:[1-9][0-9]{0,3}\s{0,2}(?:BE|AD|BC|BCE|CE)|[1-2][0-9]{3}|[5-9][0-9]|2[0-5])`
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

	// Check integer words
	if val, ok := IntegerWordDictionary[lower]; ok {
		return float64(val)
	}

	// Check special cases
	if lower == "a" || lower == "an" || lower == "the" {
		return 1
	}
	if strings.Contains(lower, "few") {
		return 3
	}
	if strings.Contains(lower, "half") {
		return 0.5
	}
	if strings.Contains(lower, "couple") {
		return 2
	}
	if strings.Contains(lower, "several") {
		return 7
	}

	// Try parsing as number
	val, _ := strconv.ParseFloat(lower, 64)
	return val
}
