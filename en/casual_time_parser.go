package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENCasualTimeParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENCasualTimeParser = parsers.ENCasualTimeParser

// NewENCasualTimeParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENCasualTimeParser() *ENCasualTimeParser {
	return parsers.NewENCasualTimeParser()
}
