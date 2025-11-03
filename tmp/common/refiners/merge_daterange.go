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
	if current.End() != nil || next.End() != nil {
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
	fromStart := fromResult.Start().(*ParsingComponents)
	toStart := toResult.Start().(*ParsingComponents)

	// Imply similar components between from and to
	if !fromStart.IsOnlyWeekdayComponent() && !toStart.IsOnlyWeekdayComponent() {
		// All possible components to check
		allComponents := []Component{
			ComponentYear, ComponentMonth, ComponentDay, ComponentWeekday,
			ComponentHour, ComponentMinute, ComponentSecond, ComponentMillisecond,
			ComponentMeridiem, ComponentTimezoneOffset,
		}

		// Copy certain components from toResult to fromResult
		for _, comp := range allComponents {
			if toStart.IsCertain(comp) && !fromStart.IsCertain(comp) {
				val := toStart.Get(comp)
				if val != nil {
					fromStart.Imply(comp, *val)
				}
			}
		}

		// Copy certain components from fromResult to toResult
		for _, comp := range allComponents {
			if fromStart.IsCertain(comp) && !toStart.IsCertain(comp) {
				val := fromStart.Get(comp)
				if val != nil {
					toStart.Imply(comp, *val)
				}
			}
		}
	}

	// Handle reversed dates
	if fromStart.Date().After(toStart.Date()) {
		fromDate := fromStart.Date()
		toDate := toStart.Date()

		// If toResult is only weekday, try adding 7 days
		if toStart.IsOnlyWeekdayComponent() {
			nextWeek := toDate.Add(7 * 24 * time.Hour)
			if nextWeek.After(fromDate) {
				toDate = nextWeek
				toStart.Imply(ComponentDay, toDate.Day())
				toStart.Imply(ComponentMonth, int(toDate.Month()))
				toStart.Imply(ComponentYear, toDate.Year())
			}
		} else if fromStart.IsOnlyWeekdayComponent() {
			// If fromResult is only weekday, try subtracting 7 days
			prevWeek := fromDate.Add(-7 * 24 * time.Hour)
			if prevWeek.Before(toDate) {
				fromDate = prevWeek
				fromStart.Imply(ComponentDay, fromDate.Day())
				fromStart.Imply(ComponentMonth, int(fromDate.Month()))
				fromStart.Imply(ComponentYear, fromDate.Year())
			}
		} else if toStart.IsDateWithUnknownYear() {
			// Try adding a year to toResult
			nextYear := toDate.AddDate(1, 0, 0)
			if nextYear.After(fromDate) {
				toDate = nextYear
				toStart.Imply(ComponentYear, toDate.Year())
			}
		} else if fromStart.IsDateWithUnknownYear() {
			// Try subtracting a year from fromResult
			prevYear := fromDate.AddDate(-1, 0, 0)
			if prevYear.Before(toDate) {
				fromDate = prevYear
				fromStart.Imply(ComponentYear, fromDate.Year())
			}
		} else {
			// Swap if still reversed
			fromResult, toResult = toResult, fromResult
			fromStart, toStart = toStart, fromStart
		}
	}

	// Create the range result by creating a new ParsingResult
	// Calculate index and text based on order
	var resultIndex int
	var resultText string
	if fromResult.Index() < toResult.Index() {
		resultIndex = fromResult.Index()
		resultText = fromResult.Text() + textBetween + toResult.Text()
	} else {
		resultIndex = toResult.Index()
		resultText = toResult.Text() + textBetween + fromResult.Text()
	}

	result := context.CreateParsingResult(resultIndex, resultText, fromStart, toStart)
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
		textBetween := context.Text()[current.Index()+len(current.Text()) : next.Index()]

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
