package refiners

import (
	"regexp"

	"github.com/kljensen/kronos"
)

// AbstractMergeDateTimeRefiner merges date-only and time-only results.
type AbstractMergeDateTimeRefiner struct {
	BaseMergingRefiner
	PatternBetweenFunc func() *regexp.Regexp
}

func (r *AbstractMergeDateTimeRefiner) ShouldMergeResults(textBetween string, current, next *kronos.ParsingResult, context *kronos.ParsingContext) bool {
	// Check if one is date-only and the other is time-only
	currentStart, okCurrent := kronos.AsParsingComponents(current.Start())
	nextStart, okNext := kronos.AsParsingComponents(next.Start())
	if !okCurrent || !okNext {
		return false
	}

	isDateTimePair := (currentStart.IsOnlyDate() && nextStart.IsOnlyTime()) ||
		(nextStart.IsOnlyDate() && currentStart.IsOnlyTime())

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

func (r *AbstractMergeDateTimeRefiner) MergeResults(textBetween string, current, next *kronos.ParsingResult, context *kronos.ParsingContext) *kronos.ParsingResult {
	currentStart, okCurrent := kronos.AsParsingComponents(current.Start())
	if !okCurrent {
		return current
	}

	// Calculate index and text for merged result
	var resultIndex int
	var resultText string
	resultIndex = current.Index()
	resultText = current.Text() + textBetween + next.Text()

	// Determine which is date and which is time, then merge
	var result *kronos.ParsingResult
	if currentStart.IsOnlyDate() {
		result = kronos.MergeDateTimeResult(current, next)
	} else {
		result = kronos.MergeDateTimeResult(next, current)
	}

	// Create new result with correct index and text
	startComponents, okStart := kronos.AsParsingComponents(result.Start())
	if !okStart {
		return current
	}
	result = context.CreateParsingResult(resultIndex, resultText, startComponents, nil)

	return result
}

func (r *AbstractMergeDateTimeRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.ParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]
		start := current.Index() + len(current.Text())
		end := next.Index()
		textBetween, okRange := kronos.SafeSlice(context.Text(), start, end)
		if !okRange {
			merged = append(merged, current)
			current = next
			continue
		}

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
