package en

import (
	"github.com/kljensen/kronos/internal/en/parsers"
)

// ENTimeExpressionParser is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENTimeExpressionParser = parsers.ENTimeExpressionParser

// NewENTimeExpressionParser is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENTimeExpressionParser(strictMode bool) *ENTimeExpressionParser {
	return parsers.NewENTimeExpressionParser(strictMode)
}
