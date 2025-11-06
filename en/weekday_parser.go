package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENWeekdayParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENWeekdayParser = parsers.ENWeekdayParser

// NewENWeekdayParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENWeekdayParser() *ENWeekdayParser {
	return parsers.NewENWeekdayParser()
}
