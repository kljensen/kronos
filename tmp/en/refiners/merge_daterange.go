package refiners

import (
	"regexp"

	"github.com/kljensen/kronos/common/refiners"
)

// ENMergeDateRangeRefiner merges before and after results.
// Examples:
//   - "2020-02-13 to 2020-02-15"
//   - "Wednesday - Friday"
type ENMergeDateRangeRefiner struct {
	refiners.AbstractMergeDateRangeRefiner
}

func NewENMergeDateRangeRefiner() *ENMergeDateRangeRefiner {
	r := &ENMergeDateRangeRefiner{}
	r.PatternBetweenFunc = r.PatternBetween
	return r
}

func (r *ENMergeDateRangeRefiner) PatternBetween() *regexp.Regexp {
	return regexp.MustCompile(`(?i)^\s*(to|-|–|until|through|till)\s*$`)
}
