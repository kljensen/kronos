package en

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
)

var (
	patternFollowBetween = regexp.MustCompile(`^\s*$`)
)

// ENMergeRelativeFollowByDateRefiner merges a relative date/time that follows an absolute date.
// Examples:
//   - "2 weeks before 2020-02-13"
//   - "2 days after next Friday"
type ENMergeRelativeFollowByDateRefiner struct{}

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

		textBetween := context.Text()[startIdx:endIdx]
		if !patternFollowBetween.MatchString(textBetween) {
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
		duration := ParseDuration(current.Text())
		if hasImpliedEarlierReferenceDate(current) {
			duration = kronos.ReverseDuration(duration)
		}

		// Create new reference from next result's date
		newRef := context.Reference().FromDate(next.Start().Date())
		components := kronos.CreateRelativeFromReference(newRef, duration)

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
