package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENYearParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENYearParser = parsers.ENYearParser

// NewENYearParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENYearParser() *ENYearParser {
	return parsers.NewENYearParser()
}
