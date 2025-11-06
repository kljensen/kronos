package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENTimeUnitCasualRelativeFormatParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENTimeUnitCasualRelativeFormatParser = parsers.ENTimeUnitCasualRelativeFormatParser

// NewENTimeUnitCasualRelativeFormatParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENTimeUnitCasualRelativeFormatParser(allowAbbreviations bool) *ENTimeUnitCasualRelativeFormatParser {
	return parsers.NewENTimeUnitCasualRelativeFormatParser(allowAbbreviations)
}
