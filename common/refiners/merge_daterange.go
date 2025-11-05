package refiners

import (
	"regexp"
	"time"

	"github.com/kljensen/kronos"
)

// AbstractMergeDateRangeRefiner merges two date results into a date range.
type AbstractMergeDateRangeRefiner struct {
	BaseMergingRefiner
	PatternBetweenFunc func() *regexp.Regexp
}

// ShouldMergeResults determines if two results should be merged into a date range.
func (r *AbstractMergeDateRangeRefiner) ShouldMergeResults(textBetween string, current, next *kronos.ParsingResult, context *kronos.ParsingContext) bool {
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

// MergeResults merges two date results into a single date range result.
func (r *AbstractMergeDateRangeRefiner) MergeResults(textBetween string, fromResult, toResult *kronos.ParsingResult, context *kronos.ParsingContext) *kronos.ParsingResult {
	fromStart, okFrom := kronos.AsParsingComponents(fromResult.Start())
	toStart, okTo := kronos.AsParsingComponents(toResult.Start())
	if !okFrom || !okTo {
		return fromResult
	}

	// Imply similar components between from and to
	if !fromStart.IsOnlyWeekdayComponent() && !toStart.IsOnlyWeekdayComponent() {
		// All possible components to check
		allComponents := []kronos.Component{
			kronos.ComponentYear, kronos.ComponentMonth, kronos.ComponentDay, kronos.ComponentWeekday,
			kronos.ComponentHour, kronos.ComponentMinute, kronos.ComponentSecond, kronos.ComponentMillisecond,
			kronos.ComponentMicrosecond, kronos.ComponentNanosecond,
			kronos.ComponentMeridiem, kronos.ComponentTimezoneOffset,
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

		switch {
		case toStart.IsOnlyWeekdayComponent():
			// If toResult is only weekday, try adding 7 days
			nextWeek := toDate.Add(7 * 24 * time.Hour)
			if nextWeek.After(fromDate) {
				toDate = nextWeek
				toStart.Imply(kronos.ComponentDay, toDate.Day())
				toStart.Imply(kronos.ComponentMonth, int(toDate.Month()))
				toStart.Imply(kronos.ComponentYear, toDate.Year())
			}
		case fromStart.IsOnlyWeekdayComponent():
			// If fromResult is only weekday, try subtracting 7 days
			prevWeek := fromDate.Add(-7 * 24 * time.Hour)
			if prevWeek.Before(toDate) {
				fromDate = prevWeek
				fromStart.Imply(kronos.ComponentDay, fromDate.Day())
				fromStart.Imply(kronos.ComponentMonth, int(fromDate.Month()))
				fromStart.Imply(kronos.ComponentYear, fromDate.Year())
			}
		case toStart.IsDateWithUnknownYear():
			// Try adding a year to toResult
			nextYear := toDate.AddDate(1, 0, 0)
			if nextYear.After(fromDate) {
				toDate = nextYear
				toStart.Imply(kronos.ComponentYear, toDate.Year())
			}
		case fromStart.IsDateWithUnknownYear():
			// Try subtracting a year from fromResult
			prevYear := fromDate.AddDate(-1, 0, 0)
			if prevYear.Before(toDate) {
				fromDate = prevYear
				fromStart.Imply(kronos.ComponentYear, fromDate.Year())
			}
		default:
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

// Refine processes results to merge date ranges based on the pattern between them.
func (r *AbstractMergeDateRangeRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
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
