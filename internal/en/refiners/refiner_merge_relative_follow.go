//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"github.com/kljensen/kronos/internal/en/data"
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
)

var patternFollowBetween = regexp.MustCompile(`^\s*$`)

// ENMergeRelativeFollowByDateRefiner merges a relative date/time that follows an absolute date.
// Examples:
//   - "2 weeks before 2020-02-13"
//   - "2 days after next Friday"
type ENMergeRelativeFollowByDateRefiner struct{}

// NewENMergeRelativeFollowByDateRefiner creates a new ENMergeRelativeFollowByDateRefiner
func NewENMergeRelativeFollowByDateRefiner() *ENMergeRelativeFollowByDateRefiner {
	return &ENMergeRelativeFollowByDateRefiner{}
}

func hasImpliedEarlierReferenceDate(result *kronos.ParsingResult) bool {
	text := strings.ToLower(result.Text())
	return strings.HasSuffix(text, " before") || strings.HasSuffix(text, " from")
}

func hasImpliedLaterReferenceDate(result *kronos.ParsingResult) bool {
	text := strings.ToLower(result.Text())
	return strings.HasSuffix(text, " after") || strings.HasSuffix(text, " since")
}

// Refine merges relative date expressions that follow absolute dates
func (r *ENMergeRelativeFollowByDateRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.ParsingResult, 0, len(results))
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

		textBetween, ok := kronos.XSafeSlice(context.Text(), startIdx, endIdx)
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
		duration := data.ParseDuration(current.Text())
		if hasImpliedEarlierReferenceDate(current) {
			duration = kronos.XReverseDuration(duration)
		}

		// Create new reference from next result's date
		newRef := context.Reference().FromDate(next.Start().Date())
		components := kronos.XCreateRelativeFromReference(newRef, duration)

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
