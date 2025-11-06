package refiners

import (
	"regexp"

	commonrefiners "github.com/kljensen/kronos/internal/common/refiners"
)

// ENMergeDateTimeRefiner merges date-only result and time-only result.
// Examples:
//   - "2020-02-13 at 6pm"
//   - "Tomorrow after 7am"
type ENMergeDateTimeRefiner struct {
	commonrefiners.AbstractMergeDateTimeRefiner
}

// NewENMergeDateTimeRefiner creates a new ENMergeDateTimeRefiner
func NewENMergeDateTimeRefiner() *ENMergeDateTimeRefiner {
	r := &ENMergeDateTimeRefiner{}
	r.PatternBetweenFunc = r.PatternBetween
	return r
}

// PatternBetween returns the regex pattern for matching date-time separators
func (r *ENMergeDateTimeRefiner) PatternBetween() *regexp.Regexp {
	return regexp.MustCompile(`^\s*(T|at|after|before|on|of|,|-|\.|∙|:)?\s*$`)
}
