package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENRelativeDateFormatParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENRelativeDateFormatParser = parsers.ENRelativeDateFormatParser

// NewENRelativeDateFormatParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENRelativeDateFormatParser() *ENRelativeDateFormatParser {
	return parsers.NewENRelativeDateFormatParser()
}
