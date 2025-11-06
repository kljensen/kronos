package en

import (
	"github.com/kljensen/kronos/internal/en/refiners"
)

// ENMergeRelativeFollowByDateRefiner is deprecated: use en.New() to create a parser instance.
// This type alias is maintained for backward compatibility.
type ENMergeRelativeFollowByDateRefiner = refiners.ENMergeRelativeFollowByDateRefiner

// NewENMergeRelativeFollowByDateRefiner is deprecated: use en.New() to create a parser instance.
// This constructor is maintained for backward compatibility.
func NewENMergeRelativeFollowByDateRefiner() *ENMergeRelativeFollowByDateRefiner {
	return refiners.NewENMergeRelativeFollowByDateRefiner()
}
