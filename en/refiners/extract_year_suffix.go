package refiners

import (
	"github.com/kljensen/kronos/internal/en/refiners"
)

// ENExtractYearSuffixRefiner is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENExtractYearSuffixRefiner = refiners.ENExtractYearSuffixRefiner

// NewENExtractYearSuffixRefiner is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENExtractYearSuffixRefiner() *ENExtractYearSuffixRefiner {
	return refiners.NewENExtractYearSuffixRefiner()
}
