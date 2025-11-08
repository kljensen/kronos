package data

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/kljensen/kronos"
)

// YearPattern matches year patterns including BE, AD, BC, BCE, CE suffixes
// Note: Longer alternatives (BCE) must come before shorter ones (BC, AD)
const YearPattern = `(?:[1-9][0-9]{0,3}\s{0,2}(?:BCE|BE|CE|AD|BC)|[1-2][0-9]{3}|[5-9][0-9]|2[0-5])`

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
var WeekdayDictionary = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"sun":       time.Sunday,
	"sun.":      time.Sunday,
	"monday":    time.Monday,
	"mon":       time.Monday,
	"mon.":      time.Monday,
	"tuesday":   time.Tuesday,
	"tue":       time.Tuesday,
	"tue.":      time.Tuesday,
	"wednesday": time.Wednesday,
	"wed":       time.Wednesday,
	"wed.":      time.Wednesday,
	"thursday":  time.Thursday,
	"thurs":     time.Thursday,
	"thurs.":    time.Thursday,
	"thur":      time.Thursday,
	"thur.":     time.Thursday,
	"thu":       time.Thursday,
	"thu.":      time.Thursday,
	"friday":    time.Friday,
	"fri":       time.Friday,
	"fri.":      time.Friday,
	"saturday":  time.Saturday,
	"sat":       time.Saturday,
	"sat.":      time.Saturday,
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
	// Numbers 13-19
	"thirteen":  13.0,
	"fourteen":  14.0,
	"fifteen":   15.0,
	"sixteen":   16.0,
	"seventeen": 17.0,
	"eighteen":  18.0,
	"nineteen":  19.0,
	// Tens (20-90)
	"twenty":  20.0,
	"thirty":  30.0,
	"forty":   40.0,
	"fifty":   50.0,
	"sixty":   60.0,
	"seventy": 70.0,
	"eighty":  80.0,
	"ninety":  90.0,
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
	"ns":           kronos.TimeunitNanosecond,
	"nano":         kronos.TimeunitNanosecond,
	"nanos":        kronos.TimeunitNanosecond,
	"nanosecond":   kronos.TimeunitNanosecond,
	"nanoseconds":  kronos.TimeunitNanosecond,
	"us":           kronos.TimeunitMicrosecond,
	"micro":        kronos.TimeunitMicrosecond,
	"micros":       kronos.TimeunitMicrosecond,
	"microsecond":  kronos.TimeunitMicrosecond,
	"microseconds": kronos.TimeunitMicrosecond,
	"ms":           kronos.TimeunitMillisecond,
	"milli":        kronos.TimeunitMillisecond,
	"millis":       kronos.TimeunitMillisecond,
	"millisecond":  kronos.TimeunitMillisecond,
	"milliseconds": kronos.TimeunitMillisecond,
	"s":            kronos.TimeunitSecond,
	"sec":          kronos.TimeunitSecond,
	"second":       kronos.TimeunitSecond,
	"seconds":      kronos.TimeunitSecond,
	"m":            kronos.TimeunitMinute,
	"min":          kronos.TimeunitMinute,
	"mins":         kronos.TimeunitMinute,
	"minute":       kronos.TimeunitMinute,
	"minutes":      kronos.TimeunitMinute,
	"h":            kronos.TimeunitHour,
	"hr":           kronos.TimeunitHour,
	"hrs":          kronos.TimeunitHour,
	"hour":         kronos.TimeunitHour,
	"hours":        kronos.TimeunitHour,
	"d":            kronos.TimeunitDay,
	"day":          kronos.TimeunitDay,
	"days":         kronos.TimeunitDay,
	"w":            kronos.TimeunitWeek,
	"week":         kronos.TimeunitWeek,
	"weeks":        kronos.TimeunitWeek,
	"mo":           kronos.TimeunitMonth,
	"mon":          kronos.TimeunitMonth,
	"mos":          kronos.TimeunitMonth,
	"month":        kronos.TimeunitMonth,
	"months":       kronos.TimeunitMonth,
	"qtr":          kronos.TimeunitQuarter,
	"quarter":      kronos.TimeunitQuarter,
	"quarters":     kronos.TimeunitQuarter,
	"y":            kronos.TimeunitYear,
	"yr":           kronos.TimeunitYear,
	"year":         kronos.TimeunitYear,
	"years":        kronos.TimeunitYear,
	"decade":       kronos.TimeunitDecade,
	"decades":      kronos.TimeunitDecade,
}

// MatchAnyPattern creates a pattern that matches any key from a map
func MatchAnyPattern(dict any) string {
	var keys []string
	switch d := dict.(type) {
	case map[string]int:
		for k := range d {
			keys = append(keys, regexp.QuoteMeta(k))
		}
	case map[string]time.Weekday:
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
)

// Time-related patterns
var (
	// NumberPattern matches integers and float numbers
	NumberPattern = `(?:\d+(?:\.\d*)?|\.\d+)`

	// TimeUnitNoAbbrPattern matches full time unit names only (no abbreviations)
	TimeUnitNoAbbrPattern = `(?:nanoseconds|nanosecond|microseconds|microsecond|milliseconds|millisecond|decades|decade|seconds|second|minutes|minute|hours|hour|days|day|weeks|week|months|month|quarters|quarter|years|year)`

	// TimeUnitPattern matches time units including abbreviations
	// Order matters: longer alternatives come first to ensure correct matching
	TimeUnitPattern = `(?:nanoseconds|nanosecond|microseconds|microsecond|milliseconds|millisecond|decades|decade|seconds|second|minutes|minute|hours|hour|days|day|weeks|week|months|month|quarters|quarter|years|year|nanos|nano|micros|micro|millis|milli|mins|min|hrs|hr|sec|mos|mon|qtr|yr|mo|ns|us|ms|s|m|h|d|w|y)`
)
