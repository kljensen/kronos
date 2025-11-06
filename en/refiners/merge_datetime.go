package refiners

import (
	"github.com/kljensen/kronos/internal/en/refiners"
)

// ENMergeDateTimeRefiner is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENMergeDateTimeRefiner = refiners.ENMergeDateTimeRefiner

// NewENMergeDateTimeRefiner is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENMergeDateTimeRefiner() *ENMergeDateTimeRefiner {
	return refiners.NewENMergeDateTimeRefiner()
}
