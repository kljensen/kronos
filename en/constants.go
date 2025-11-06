package en

import (
	"github.com/kljensen/kronos/internal/en/data"
	kronos "github.com/kljensen/kronos"
)

// Deprecated: The following constants, variables, and functions are deprecated.
// Use en.New() to create a parser instance instead of accessing these directly.
// These are maintained for backward compatibility only.

// YearPattern is deprecated: use en.New() instead.
const YearPattern = data.YearPattern

// Deprecated dictionaries - use en.New() instead
var (
	MonthDictionary         = data.MonthDictionary
	FullMonthNameDictionary = data.FullMonthNameDictionary
	WeekdayDictionary       = data.WeekdayDictionary
	IntegerWordDictionary   = data.IntegerWordDictionary
	NumberWordDictionary    = data.NumberWordDictionary
	OrdinalWordDictionary   = data.OrdinalWordDictionary
	TimeUnitDictionary      = data.TimeUnitDictionary
)

// Deprecated pattern variables - use en.New() instead
var (
	MonthPattern         = data.MonthPattern
	FullMonthPattern     = data.FullMonthPattern
	WeekdayPattern       = data.WeekdayPattern
	IntegerWordPattern   = data.IntegerWordPattern
	OrdinalWordPattern   = data.OrdinalWordPattern
	OrdinalNumberPattern = data.OrdinalNumberPattern
	NumberPattern        = data.NumberPattern
	TimeUnitNoAbbrPattern = data.TimeUnitNoAbbrPattern
	TimeUnitPattern      = data.TimeUnitPattern
)

// Deprecated functions - use en.New() instead

// MatchAnyPattern is deprecated: use en.New() instead.
func MatchAnyPattern(dict interface{}) string {
	return data.MatchAnyPattern(dict)
}

// ParseOrdinalNumber is deprecated: use en.New() instead.
func ParseOrdinalNumber(match string) int {
	return data.ParseOrdinalNumber(match)
}

// ParseYear is deprecated: use en.New() instead.
func ParseYear(match string) int {
	return data.ParseYear(match)
}

// ParseNumberPattern is deprecated: use en.New() instead.
func ParseNumberPattern(match string) float64 {
	return data.ParseNumberPattern(match)
}

// ParseDuration is deprecated: use en.New() instead.
func ParseDuration(text string) kronos.Duration {
	return data.ParseDuration(text)
}

// IsEmptyDuration is deprecated: use en.New() instead.
func IsEmptyDuration(d kronos.Duration) bool {
	return data.IsEmptyDuration(d)
}
