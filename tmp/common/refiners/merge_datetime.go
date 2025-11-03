package refiners

import (
	"regexp"

	. "github.com/markusmobius/go-chrono"
)

// AbstractMergeDateTimeRefiner merges date-only and time-only results.
type AbstractMergeDateTimeRefiner struct {
	BaseMergingRefiner
	PatternBetweenFunc func() *regexp.Regexp
}

func (r *AbstractMergeDateTimeRefiner) ShouldMergeResults(textBetween string, current, next *ParsingResult, context *ParsingContext) bool {
	// Check if one is date-only and the other is time-only
	isDateTimePair := (current.Start.IsOnlyDate() && next.Start.IsOnlyTime()) ||
		(next.Start.IsOnlyDate() && current.Start.IsOnlyTime())

	if !isDateTimePair {
		return false
	}

	// Check if the text between matches the pattern
	if r.PatternBetweenFunc != nil {
		pattern := r.PatternBetweenFunc()
		return pattern.MatchString(textBetween)
	}

	return false
}

func (r *AbstractMergeDateTimeRefiner) MergeResults(textBetween string, current, next *ParsingResult, context *ParsingContext) *ParsingResult {
	var result *ParsingResult

	// Determine which is date and which is time
	if current.Start.IsOnlyDate() {
		result = MergeDateTimeResult(current, next)
	} else {
		result = MergeDateTimeResult(next, current)
	}

	result.Index = current.Index
	result.Text = current.Text + textBetween + next.Text

	return result
}

func (r *AbstractMergeDateTimeRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*ParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]
		textBetween := context.Text[current.Index+len(current.Text) : next.Index]

		if !r.ShouldMergeResults(textBetween, current, next, context) {
			merged = append(merged, current)
			current = next
		} else {
			current = r.MergeResults(textBetween, current, next, context)
		}
	}

	if current != nil {
		merged = append(merged, current)
	}

	return merged
}
