package refiners

import (
	"regexp"
	"strings"

	. "github.com/markusmobius/go-chrono"
	"github.com/markusmobius/go-chrono/en"
)

var (
	yearSuffixPattern = regexp.MustCompile(`^\s*(` + en.YEAR_PATTERN + `)`)
)

// ENExtractYearSuffixRefiner extracts year suffixes from dates.
// Example: "Dec 12, 2020" - pulls the year suffix
type ENExtractYearSuffixRefiner struct{}

func NewENExtractYearSuffixRefiner() *ENExtractYearSuffixRefiner {
	return &ENExtractYearSuffixRefiner{}
}

func (r *ENExtractYearSuffixRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	for _, result := range results {
		if !result.Start.IsDateWithUnknownYear() {
			continue
		}

		suffix := context.Text[result.Index+len(result.Text):]
		match := yearSuffixPattern.FindStringSubmatch(suffix)
		if match == nil {
			continue
		}

		// If the suffix match is just a short number, don't assume it's a year
		if len(strings.TrimSpace(match[0])) <= 3 {
			continue
		}

		year := en.ParseYear(match[1])
		if result.End != nil {
			result.End.Assign(ComponentYear, year)
		}
		result.Start.Assign(ComponentYear, year)
		result.Text += match[0]
	}

	return results
}
