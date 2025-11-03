package refiners

import (
	"regexp"

	. "github.com/markusmobius/go-chrono"
	"github.com/markusmobius/go-chrono/common/refiners"
)

// ENMergeDateTimeRefiner merges date-only result and time-only result.
// Examples:
//   - "2020-02-13 at 6pm"
//   - "Tomorrow after 7am"
type ENMergeDateTimeRefiner struct {
	refiners.AbstractMergeDateTimeRefiner
}

func NewENMergeDateTimeRefiner() *ENMergeDateTimeRefiner {
	r := &ENMergeDateTimeRefiner{}
	r.PatternBetweenFunc = r.PatternBetween
	return r
}

func (r *ENMergeDateTimeRefiner) PatternBetween() *regexp.Regexp {
	return regexp.MustCompile(`^\s*(T|at|after|before|on|of|,|-|\.|∙|:)?\s*$`)
}
