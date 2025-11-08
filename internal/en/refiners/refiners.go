//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal"
	commonrefiners "github.com/kljensen/kronos/internal/common/refiners"
	endata "github.com/kljensen/kronos/internal/en/data"
	"github.com/kljensen/kronos/internal/helpers"
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

// ENExtractYearSuffixRefiner extracts year suffixes from dates.
// Example: "Dec 12, 2020" - pulls the year suffix
type ENExtractYearSuffixRefiner struct{}

// NewENExtractYearSuffixRefiner creates a new ENExtractYearSuffixRefiner
func NewENExtractYearSuffixRefiner() *ENExtractYearSuffixRefiner {
	return &ENExtractYearSuffixRefiner{}
}

const (
	// YEAR_PATTERN matches year patterns including BE, AD, BC, BCE, CE suffixes
	yearPattern = `(?:[1-9][0-9]{0,3}\s{0,2}(?:BE|AD|BC|BCE|CE)|[1-2][0-9]{3}|[5-9][0-9]|2[0-5])`
)

var yearSuffixPattern = regexp.MustCompile(`^\s*(` + yearPattern + `)`)

// Refine extracts year suffixes from dates in parsing results
func (r *ENExtractYearSuffixRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	for i, result := range results {
		resultStart, okStart := helpers.AsParsingComponents(result.Start())
		if !okStart {
			continue
		}
		if !resultStart.IsDateWithUnknownYear() {
			continue
		}

		start := result.Index() + len(result.Text())
		suffix, okSlice := helpers.SafeSlice(context.Text(), start, len(context.Text()))
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
		var resultEnd *kronos.InternalParsingComponents
		if result.End() != nil {
			if endComponents, okEnd := helpers.AsParsingComponents(result.End()); okEnd {
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
		year, err := strconv.Atoi(strings.TrimSpace(cleaned))
		if err != nil {
			return 0
		}
		return year - 543
	}

	// Before Christ / Before Common Era
	if regexp.MustCompile(`(?i)BCE?`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*BCE?`).ReplaceAllString(match, "")
		year, err := strconv.Atoi(strings.TrimSpace(cleaned))
		if err != nil {
			return 0
		}
		return -year
	}

	// Anno Domini / Common Era
	if regexp.MustCompile(`(?i)(AD|CE)`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*(AD|CE)`).ReplaceAllString(match, "")
		year, err := strconv.Atoi(strings.TrimSpace(cleaned))
		if err != nil {
			return 0
		}
		return year
	}

	// Regular year number
	year, err := strconv.Atoi(strings.TrimSpace(match))
	if err != nil {
		return 0
	}
	return helpers.FindMostLikelyADYear(year)
}

// ENMergeRelativeAfterDateRefiner merges a relative date/time that comes after an absolute date.
// Examples:
//   - "2020-02-13 +2 weeks"
//   - "next tuesday +10 days"
type ENMergeRelativeAfterDateRefiner struct{}

// NewENMergeRelativeAfterDateRefiner creates a new ENMergeRelativeAfterDateRefiner
func NewENMergeRelativeAfterDateRefiner() *ENMergeRelativeAfterDateRefiner {
	return &ENMergeRelativeAfterDateRefiner{}
}

var (
	patternAfterBetween      = regexp.MustCompile(`^\s*$`)
	patternPositiveFollowing = regexp.MustCompile(`^[+-]`)
	patternNegativeFollowing = regexp.MustCompile(`^-`)
)

func isPositiveFollowingReference(result *kronos.InternalParsingResult) bool {
	return patternPositiveFollowing.MatchString(result.Text())
}

func isNegativeFollowingReference(result *kronos.InternalParsingResult) bool {
	return patternNegativeFollowing.MatchString(result.Text())
}

// Refine merges relative date expressions that come after absolute dates
func (r *ENMergeRelativeAfterDateRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.InternalParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]

		// Check if dates are adjacent
		startIdx := current.Index() + len(current.Text())
		endIdx := next.Index()

		// Skip if results overlap or are not in order
		if startIdx > endIdx {
			merged = append(merged, current)
			current = next
			continue
		}

		textBetween, ok := helpers.SafeSlice(context.Text(), startIdx, endIdx)
		if !ok || !patternAfterBetween.MatchString(textBetween) {
			merged = append(merged, current)
			current = next
			continue
		}

		// Check if next has +/- prefix
		if !isPositiveFollowingReference(next) && !isNegativeFollowingReference(next) {
			merged = append(merged, current)
			current = next
			continue
		}

		// Merge the results
		duration := endata.ParseDuration(strings.TrimPrefix(next.Text(), "+"))
		if isNegativeFollowingReference(next) {
			duration = internal.ReverseDuration(duration)
		}

		// Create new reference from current result's date
		newRef := kronos.InternalNewReferenceWithTimezone(current.Start().Date(), nil)
		components := helpers.CreateRelativeFromReference(newRef, duration, internal.EmptyDuration)

		// Create merged result
		resultIndex := current.Index()
		resultText := current.Text() + textBetween + next.Text()
		result := context.CreateParsingResult(resultIndex, resultText, components, nil)

		current = result
	}

	if current != nil {
		merged = append(merged, current)
	}

	return merged
}

// ENMergeRelativeFollowByDateRefiner merges a relative date/time that follows an absolute date.
// Examples:
//   - "2 weeks before 2020-02-13"
//   - "2 days after next Friday"
type ENMergeRelativeFollowByDateRefiner struct{}

// NewENMergeRelativeFollowByDateRefiner creates a new ENMergeRelativeFollowByDateRefiner
func NewENMergeRelativeFollowByDateRefiner() *ENMergeRelativeFollowByDateRefiner {
	return &ENMergeRelativeFollowByDateRefiner{}
}

var patternFollowBetween = regexp.MustCompile(`^\s*$`)

func hasImpliedEarlierReferenceDate(result *kronos.InternalParsingResult) bool {
	text := strings.ToLower(result.Text())
	return strings.HasSuffix(text, " before") || strings.HasSuffix(text, " from")
}

func hasImpliedLaterReferenceDate(result *kronos.InternalParsingResult) bool {
	text := strings.ToLower(result.Text())
	return strings.HasSuffix(text, " after") || strings.HasSuffix(text, " since")
}

// Refine merges relative date expressions that follow absolute dates
func (r *ENMergeRelativeFollowByDateRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.InternalParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]

		// Check if dates are adjacent
		startIdx := current.Index() + len(current.Text())
		endIdx := next.Index()

		// Skip if results overlap or are not in order
		if startIdx > endIdx {
			merged = append(merged, current)
			current = next
			continue
		}

		textBetween, ok := helpers.SafeSlice(context.Text(), startIdx, endIdx)
		if !ok || !patternFollowBetween.MatchString(textBetween) {
			merged = append(merged, current)
			current = next
			continue
		}

		// Check if current has implied earlier/later reference
		if !hasImpliedEarlierReferenceDate(current) && !hasImpliedLaterReferenceDate(current) {
			merged = append(merged, current)
			current = next
			continue
		}

		// Check if next implies an absolute date
		dayVal := next.Start().Get(kronos.ComponentDay)
		monthVal := next.Start().Get(kronos.ComponentMonth)
		yearVal := next.Start().Get(kronos.ComponentYear)
		if dayVal == nil || *dayVal == 0 || monthVal == nil || *monthVal == 0 || yearVal == nil || *yearVal == 0 {
			merged = append(merged, current)
			current = next
			continue
		}

		// Merge the results
		duration := endata.ParseDuration(current.Text())
		if hasImpliedEarlierReferenceDate(current) {
			duration = internal.ReverseDuration(duration)
		}

		// Create new reference from next result's date
		newRef := kronos.InternalNewReferenceWithTimezone(next.Start().Date(), nil)
		components := helpers.CreateRelativeFromReference(newRef, duration, internal.EmptyDuration)

		// Create merged result
		resultIndex := current.Index()
		resultText := current.Text() + textBetween + next.Text()
		result := context.CreateParsingResult(resultIndex, resultText, components, nil)

		current = result
	}

	if current != nil {
		merged = append(merged, current)
	}

	return merged
}

// ENUnlikelyFormatFilter filters out unlikely English date formats.
type ENUnlikelyFormatFilter struct{}

// NewENUnlikelyFormatFilter creates a new ENUnlikelyFormatFilter
func NewENUnlikelyFormatFilter() *ENUnlikelyFormatFilter {
	return &ENUnlikelyFormatFilter{}
}

var mayContextPattern = regexp.MustCompile(`(?i)\b(in)$`)

// Refine filters out unlikely English date formats from parsing results
func (f *ENUnlikelyFormatFilter) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	filtered := make([]*kronos.InternalParsingResult, 0, len(results))

	for _, result := range results {
		if !f.isValid(context, result) {
			continue
		}
		filtered = append(filtered, result)
	}

	return filtered
}

func (f *ENUnlikelyFormatFilter) isValid(context *kronos.InternalParsingContext, result *kronos.InternalParsingResult) bool {
	text := strings.TrimSpace(result.Text())

	// If the result consists of the whole text, it's likely valid
	if text == strings.TrimSpace(context.Text()) {
		return true
	}

	// "may" is a month name but also a modal verb
	// Check if the text before "may" follows allowed patterns
	if strings.ToLower(text) == "may" {
		textBefore := strings.TrimSpace(context.Text()[:result.Index()])
		if !mayContextPattern.MatchString(textBefore) {
			return false
		}
	}

	// "the second" could refer to the ordinal number or timeunit
	if strings.HasSuffix(strings.ToLower(text), "the second") {
		textAfter := strings.TrimSpace(context.Text()[result.Index()+len(result.Text()):])
		if len(textAfter) > 0 {
			return false
		}
	}

	return true
}
