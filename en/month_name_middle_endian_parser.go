package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENMonthNameMiddleEndianParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENMonthNameMiddleEndianParser = parsers.ENMonthNameMiddleEndianParser

// NewENMonthNameMiddleEndianParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENMonthNameMiddleEndianParser(shouldSkipYearLikeDate bool) *ENMonthNameMiddleEndianParser {
	return parsers.NewENMonthNameMiddleEndianParser(shouldSkipYearLikeDate)
}
