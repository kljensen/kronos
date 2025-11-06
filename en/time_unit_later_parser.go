package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENTimeUnitLaterFormatParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENTimeUnitLaterFormatParser = parsers.ENTimeUnitLaterFormatParser

// NewENTimeUnitLaterFormatParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENTimeUnitLaterFormatParser(strictMode bool) *ENTimeUnitLaterFormatParser {
	return parsers.NewENTimeUnitLaterFormatParser(strictMode)
}
