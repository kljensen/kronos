package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENYearMonthDayParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENYearMonthDayParser = parsers.ENYearMonthDayParser

// NewENYearMonthDayParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENYearMonthDayParser(strictMonthDateOrder bool) *ENYearMonthDayParser {
	return parsers.NewENYearMonthDayParser(strictMonthDateOrder)
}
