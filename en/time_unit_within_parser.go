package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENTimeUnitWithinFormatParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENTimeUnitWithinFormatParser = parsers.ENTimeUnitWithinFormatParser

// NewENTimeUnitWithinFormatParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENTimeUnitWithinFormatParser(strictMode bool) *ENTimeUnitWithinFormatParser {
	return parsers.NewENTimeUnitWithinFormatParser(strictMode)
}
