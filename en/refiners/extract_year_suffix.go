package refiners

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/kljensen/kronos"
)

const (
	// YEAR_PATTERN matches year patterns including BE, AD, BC, BCE, CE suffixes
	yearPattern = `(?:[1-9][0-9]{0,3}\s{0,2}(?:BE|AD|BC|BCE|CE)|[1-2][0-9]{3}|[5-9][0-9]|2[0-5])`
)

var (
	yearSuffixPattern = regexp.MustCompile(`^\s*(` + yearPattern + `)`)
)

// ENExtractYearSuffixRefiner extracts year suffixes from dates.
// Example: "Dec 12, 2020" - pulls the year suffix
type ENExtractYearSuffixRefiner struct{}

func NewENExtractYearSuffixRefiner() *ENExtractYearSuffixRefiner {
	return &ENExtractYearSuffixRefiner{}
}

func (r *ENExtractYearSuffixRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	for i, result := range results {
		resultStart, okStart := kronos.AsParsingComponents(result.Start())
		if !okStart {
			continue
		}
		if !resultStart.IsDateWithUnknownYear() {
			continue
		}

		start := result.Index() + len(result.Text())
		suffix, okSlice := kronos.SafeSlice(context.Text(), start, len(context.Text()))
		if !okSlice {
			continue
		}
		match := yearSuffixPattern.FindStringSubmatch(suffix)
		if match == nil {
			continue
		}

		// If the suffix match is just a short number, don't assume it's a year
		if len(strings.TrimSpace(match[0])) <= 3 {
			continue
		}

		year := parseYear(match[1])
		var resultEnd *kronos.ParsingComponents
		if result.End() != nil {
			if endComponents, okEnd := kronos.AsParsingComponents(result.End()); okEnd {
				resultEnd = endComponents
				resultEnd.Assign(kronos.ComponentYear, year)
			}
		}
		resultStart.Assign(kronos.ComponentYear, year)

		// Create new result with updated text
		newText := result.Text() + match[0]
		newResult := context.CreateParsingResult(result.Index(), newText, resultStart, resultEnd)
		results[i] = newResult
	}

	return results
}

// parseYear parses a year pattern (handles BE, AD, BC, BCE, CE)
func parseYear(match string) int {
	// Buddhist Era
	if regexp.MustCompile(`(?i)BE`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*BE`).ReplaceAllString(match, "")
		year, _ := strconv.Atoi(strings.TrimSpace(cleaned))
		return year - 543
	}

	// Before Christ / Before Common Era
	if regexp.MustCompile(`(?i)BCE?`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*BCE?`).ReplaceAllString(match, "")
		year, _ := strconv.Atoi(strings.TrimSpace(cleaned))
		return -year
	}

	// Anno Domini / Common Era
	if regexp.MustCompile(`(?i)(AD|CE)`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*(AD|CE)`).ReplaceAllString(match, "")
		year, _ := strconv.Atoi(strings.TrimSpace(cleaned))
		return year
	}

	// Regular year number
	year, _ := strconv.Atoi(strings.TrimSpace(match))
	return kronos.FindMostLikelyADYear(year)
}
