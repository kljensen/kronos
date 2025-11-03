package refiners

import (
	"regexp"
	"time"

	. "github.com/kljensen/kronos"
)

// AbstractMergeDateRangeRefiner merges two date results into a date range.
type AbstractMergeDateRangeRefiner struct {
	BaseMergingRefiner
	PatternBetweenFunc func() *regexp.Regexp
}

func (r *AbstractMergeDateRangeRefiner) ShouldMergeResults(textBetween string, current, next *ParsingResult, context *ParsingContext) bool {
	// Both results should not already have an end
	if current.End != nil || next.End != nil {
		return false
	}

	// Check if the text between matches the pattern
	if r.PatternBetweenFunc != nil {
		pattern := r.PatternBetweenFunc()
		return pattern.MatchString(textBetween)
	}

	return false
}

func (r *AbstractMergeDateRangeRefiner) MergeResults(textBetween string, fromResult, toResult *ParsingResult, context *ParsingContext) *ParsingResult {
	// Imply similar components between from and to
	if !fromResult.Start.IsOnlyWeekdayComponent() && !toResult.Start.IsOnlyWeekdayComponent() {
		// Copy certain components from toResult to fromResult
		for _, comp := range toResult.Start.CertainComponents() {
			if !fromResult.Start.IsCertain(comp) {
				fromResult.Start.Imply(comp, toResult.Start.Get(comp))
			}
		}

		// Copy certain components from fromResult to toResult
		for _, comp := range fromResult.Start.CertainComponents() {
			if !toResult.Start.IsCertain(comp) {
				toResult.Start.Imply(comp, fromResult.Start.Get(comp))
			}
		}
	}

	// Handle reversed dates
	if fromResult.Start.Date().After(toResult.Start.Date()) {
		fromDate := fromResult.Start.Date()
		toDate := toResult.Start.Date()

		// If toResult is only weekday, try adding 7 days
		if toResult.Start.IsOnlyWeekdayComponent() {
			nextWeek := toDate.Add(7 * 24 * time.Hour)
			if nextWeek.After(fromDate) {
				toDate = nextWeek
				toResult.Start.Imply(ComponentDay, toDate.Day())
				toResult.Start.Imply(ComponentMonth, int(toDate.Month()))
				toResult.Start.Imply(ComponentYear, toDate.Year())
			}
		} else if fromResult.Start.IsOnlyWeekdayComponent() {
			// If fromResult is only weekday, try subtracting 7 days
			prevWeek := fromDate.Add(-7 * 24 * time.Hour)
			if prevWeek.Before(toDate) {
				fromDate = prevWeek
				fromResult.Start.Imply(ComponentDay, fromDate.Day())
				fromResult.Start.Imply(ComponentMonth, int(fromDate.Month()))
				fromResult.Start.Imply(ComponentYear, fromDate.Year())
			}
		} else if toResult.Start.IsDateWithUnknownYear() {
			// Try adding a year to toResult
			nextYear := toDate.AddDate(1, 0, 0)
			if nextYear.After(fromDate) {
				toDate = nextYear
				toResult.Start.Imply(ComponentYear, toDate.Year())
			}
		} else if fromResult.Start.IsDateWithUnknownYear() {
			// Try subtracting a year from fromResult
			prevYear := fromDate.AddDate(-1, 0, 0)
			if prevYear.Before(toDate) {
				fromDate = prevYear
				fromResult.Start.Imply(ComponentYear, fromDate.Year())
			}
		} else {
			// Swap if still reversed
			fromResult, toResult = toResult, fromResult
		}
	}

	// Create the range result
	result := fromResult.Clone()
	result.Start = fromResult.Start
	result.End = toResult.Start

	// Set index and text
	if fromResult.Index < toResult.Index {
		result.Index = fromResult.Index
		result.Text = fromResult.Text + textBetween + toResult.Text
	} else {
		result.Index = toResult.Index
		result.Text = toResult.Text + textBetween + fromResult.Text
	}

	return result
}

func (r *AbstractMergeDateRangeRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
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
