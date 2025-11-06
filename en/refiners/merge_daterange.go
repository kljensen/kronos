package refiners

import (
	"github.com/kljensen/kronos/internal/en/refiners"
)

// ENMergeDateRangeRefiner is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENMergeDateRangeRefiner = refiners.ENMergeDateRangeRefiner

// NewENMergeDateRangeRefiner is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENMergeDateRangeRefiner() *ENMergeDateRangeRefiner {
	return refiners.NewENMergeDateRangeRefiner()
}
