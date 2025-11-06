package refiners

import (
	"github.com/kljensen/kronos/internal/en/refiners"
)

// ENUnlikelyFormatFilter is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENUnlikelyFormatFilter = refiners.ENUnlikelyFormatFilter

// NewENUnlikelyFormatFilter is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENUnlikelyFormatFilter() *ENUnlikelyFormatFilter {
	return refiners.NewENUnlikelyFormatFilter()
}
