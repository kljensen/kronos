package en

import (
	"regexp"
	"strings"

	. "github.com/kljensen/kronos"
)

var (
	patternAfterBetween      = regexp.MustCompile(`^\s*$`)
	patternPositiveFollowing = regexp.MustCompile(`^[+-]`)
	patternNegativeFollowing = regexp.MustCompile(`^-`)
)

// ENMergeRelativeAfterDateRefiner merges a relative date/time that comes after an absolute date.
// Examples:
//   - "2020-02-13 +2 weeks"
//   - "next tuesday +10 days"
type ENMergeRelativeAfterDateRefiner struct{}

func NewENMergeRelativeAfterDateRefiner() *ENMergeRelativeAfterDateRefiner {
	return &ENMergeRelativeAfterDateRefiner{}
}

func isPositiveFollowingReference(result *ParsingResult) bool {
	return patternPositiveFollowing.MatchString(result.Text())
}

func isNegativeFollowingReference(result *ParsingResult) bool {
	return patternNegativeFollowing.MatchString(result.Text())
}

func (r *ENMergeRelativeAfterDateRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*ParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]

		// Check if dates are adjacent
		textBetween := context.Text()[current.Index()+len(current.Text()) : next.Index()]
		if !patternAfterBetween.MatchString(textBetween) {
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
		duration := ParseDuration(strings.TrimPrefix(next.Text(), "+"))
		if isNegativeFollowingReference(next) {
			duration = ReverseDuration(duration)
		}

		// Create new reference from current result's date
		newRef := context.Reference().FromDate(current.Start().Date())
		components := CreateRelativeFromReference(newRef, duration)

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
