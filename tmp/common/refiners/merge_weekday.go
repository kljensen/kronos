package refiners

import (
	"regexp"

	kronos "github.com/kljensen/kronos"
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
			textBetween := context.Text[currentResult.Index+len(currentResult.Text) : nextResult.Index]

			if r.shouldMergeResults(textBetween, currentResult, nextResult) {
				mergedResult := r.mergeResults(textBetween, currentResult, nextResult)
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
) bool {
	// Merge when:
	// 1. Current result is weekday-only
	// 2. Current result has no certain hour
	// 3. Next result has a certain day
	// 4. Text between is just optional comma and whitespace
	weekdayThenNormalDate := currentResult.Start.IsOnlyWeekdayComponent() &&
		!currentResult.Start.IsCertain(kronos.ComponentHour) &&
		nextResult.Start.IsCertain(kronos.ComponentDay)

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
) *kronos.ParsingResult {
	// Create a new result based on nextResult but with updated index and text
	newResult := &kronos.ParsingResult{
		Reference: nextResult.Reference,
		Index:     currentResult.Index,
		Text:      currentResult.Text + textBetween + nextResult.Text,
		Start:     nextResult.Start.Clone(),
		End:       nextResult.End,
		RefDate:   nextResult.RefDate,
	}

	// Assign the weekday from the current result
	newResult.Start.Assign(kronos.ComponentWeekday, currentResult.Start.Get(kronos.ComponentWeekday))
	if newResult.End != nil {
		newResult.End.Assign(kronos.ComponentWeekday, currentResult.Start.Get(kronos.ComponentWeekday))
	}

	return newResult
}
