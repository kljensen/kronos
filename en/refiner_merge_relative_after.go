package en

import (
	"github.com/kljensen/kronos/internal/en/refiners"
)

// ENMergeRelativeAfterDateRefiner is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENMergeRelativeAfterDateRefiner = refiners.ENMergeRelativeAfterDateRefiner

// NewENMergeRelativeAfterDateRefiner is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENMergeRelativeAfterDateRefiner() *ENMergeRelativeAfterDateRefiner {
	return refiners.NewENMergeRelativeAfterDateRefiner()
}
