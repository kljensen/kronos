package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENCompactFormatParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENCompactFormatParser = parsers.ENCompactFormatParser

// NewENCompactFormatParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENCompactFormatParser() *ENCompactFormatParser {
	return parsers.NewENCompactFormatParser()
}
