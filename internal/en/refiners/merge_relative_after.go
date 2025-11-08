//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal"
	endata "github.com/kljensen/kronos/internal/en/data"
	"github.com/kljensen/kronos/internal/helpers"
)

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
