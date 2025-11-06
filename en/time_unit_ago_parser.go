package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENTimeUnitAgoFormatParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENTimeUnitAgoFormatParser = parsers.ENTimeUnitAgoFormatParser

// NewENTimeUnitAgoFormatParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENTimeUnitAgoFormatParser(strictMode bool) *ENTimeUnitAgoFormatParser {
	return parsers.NewENTimeUnitAgoFormatParser(strictMode)
}
