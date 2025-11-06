package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENCasualDateParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENCasualDateParser = parsers.ENCasualDateParser

// NewENCasualDateParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENCasualDateParser() *ENCasualDateParser {
	return parsers.NewENCasualDateParser()
}
