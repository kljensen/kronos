//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"regexp"

	commonrefiners "github.com/kljensen/kronos/internal/common/refiners"
)

// ENMergeDateRangeRefiner merges before and after results.
// Examples:
//   - "2020-02-13 to 2020-02-15"
//   - "Wednesday - Friday"
type ENMergeDateRangeRefiner struct {
	commonrefiners.AbstractMergeDateRangeRefiner
}

// NewENMergeDateRangeRefiner creates a new ENMergeDateRangeRefiner
func NewENMergeDateRangeRefiner() *ENMergeDateRangeRefiner {
	r := &ENMergeDateRangeRefiner{}
	r.PatternBetweenFunc = r.PatternBetween
	return r
}

// PatternBetween returns the regex pattern for matching range separators
func (r *ENMergeDateRangeRefiner) PatternBetween() *regexp.Regexp {
	return regexp.MustCompile(`(?i)^\s*(to|-|–|until|through|till)\s*$`)
}
