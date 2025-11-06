//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"regexp"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// MergeWeekdayComponentRefiner merges weekday-only results with adjacent date results.
// Examples:
//   - "Sunday 12/7/2014" -> merges weekday into date
//   - "Tuesday, January 13, 2012" -> merges weekday into date
type MergeWeekdayComponentRefiner struct{}

// NewMergeWeekdayComponentRefiner creates a new MergeWeekdayComponentRefiner
func NewMergeWeekdayComponentRefiner() *MergeWeekdayComponentRefiner {
	return &MergeWeekdayComponentRefiner{}
}

// Refine merges weekday components with adjacent date components
func (r *MergeWeekdayComponentRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.ParsingResult, 0, len(results))
	i := 0

	for i < len(results) {
		currentResult := results[i]

		// Check if we can merge with the next result
		if i+1 < len(results) {
			nextResult := results[i+1]
			start := currentResult.Index() + len(currentResult.Text())
			end := nextResult.Index()
			textBetween, okRange := helpers.SafeSlice(context.Text(), start, end)
			if !okRange {
				merged = append(merged, currentResult)
				i++
				continue
			}

			if r.shouldMergeResults(textBetween, currentResult, nextResult, context) {
				mergedResult := r.mergeResults(textBetween, currentResult, nextResult, context)
				merged = append(merged, mergedResult)
				i += 2 // Skip both current and next
				continue
			}
		}

		// No merge, just add current result
		merged = append(merged, currentResult)
		i++
	}

	return merged
}

func (r *MergeWeekdayComponentRefiner) shouldMergeResults(
	textBetween string,
	currentResult *kronos.ParsingResult,
	nextResult *kronos.ParsingResult,
	context *kronos.ParsingContext,
) bool {
	// Merge when:
	// 1. Current result is weekday-only
	// 2. Current result has no certain hour
	// 3. Next result has a certain day
	// 4. Text between is just optional comma and whitespace
	currentStart, okCurrent := kronos.XAsParsingComponents(currentResult.Start())
	nextStart, okNext := kronos.XAsParsingComponents(nextResult.Start())
	if !okCurrent || !okNext {
		return false
	}

	weekdayThenNormalDate := currentStart.IsOnlyWeekdayComponent() &&
		!currentStart.IsCertain(kronos.ComponentHour) &&
		nextStart.IsCertain(kronos.ComponentDay)

	if !weekdayThenNormalDate {
		return false
	}

	// Check that text between is just comma and/or whitespace
	return regexp.MustCompile(`^,?\s*$`).MatchString(textBetween)
}

func (r *MergeWeekdayComponentRefiner) mergeResults(
	textBetween string,
	currentResult *kronos.ParsingResult,
	nextResult *kronos.ParsingResult,
	context *kronos.ParsingContext,
) *kronos.ParsingResult {
	// Get start components
	currentStart, okCurrent := kronos.XAsParsingComponents(currentResult.Start())
	nextStart, okNext := kronos.XAsParsingComponents(nextResult.Start())
	if !okCurrent || !okNext {
		return currentResult
	}

	// Clone the next result's start components
	newStart := nextStart.Clone()

	// Assign the weekday from the current result
	weekdayVal := currentStart.Get(kronos.ComponentWeekday)
	if weekdayVal != nil {
		newStart.Assign(kronos.ComponentWeekday, *weekdayVal)
	}

	// Handle end components if present
	var newEnd *kronos.ParsingComponents
	if nextResult.End() != nil {
		if nextEnd, ok := kronos.XAsParsingComponents(nextResult.End()); ok {
			newEnd = nextEnd.Clone()
			if weekdayVal != nil {
				newEnd.Assign(kronos.ComponentWeekday, *weekdayVal)
			}
		}
	}

	// Create a new result with correct index and text
	resultIndex := currentResult.Index()
	resultText := currentResult.Text() + textBetween + nextResult.Text()
	newResult := context.CreateParsingResult(resultIndex, resultText, newStart, newEnd)

	return newResult
}
