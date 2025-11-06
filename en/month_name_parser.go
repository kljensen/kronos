package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENMonthNameParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENMonthNameParser = parsers.ENMonthNameParser

// NewENMonthNameParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENMonthNameParser() *ENMonthNameParser {
	return parsers.NewENMonthNameParser()
}
