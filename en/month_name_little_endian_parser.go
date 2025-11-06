package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENMonthNameLittleEndianParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENMonthNameLittleEndianParser = parsers.ENMonthNameLittleEndianParser

// NewENMonthNameLittleEndianParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENMonthNameLittleEndianParser() *ENMonthNameLittleEndianParser {
	return parsers.NewENMonthNameLittleEndianParser()
}
