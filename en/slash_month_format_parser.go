package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENSlashMonthFormatParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENSlashMonthFormatParser = parsers.ENSlashMonthFormatParser

// NewENSlashMonthFormatParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENSlashMonthFormatParser() *ENSlashMonthFormatParser {
	return parsers.NewENSlashMonthFormatParser()
}
