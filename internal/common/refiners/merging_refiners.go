//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"regexp"
	"strconv"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// AbstractMergeDateTimeRefiner merges date-only and time-only results.
type AbstractMergeDateTimeRefiner struct {
	BaseMergingRefiner
	PatternBetweenFunc func() *regexp.Regexp
}

// ShouldMergeResults determines if a date-only and time-only result should be merged.
func (r *AbstractMergeDateTimeRefiner) ShouldMergeResults(textBetween string, current, next *kronos.InternalParsingResult, context *kronos.InternalParsingContext) bool {
	// Check if one is date-only and the other is time-only
	currentStart, okCurrent := helpers.AsParsingComponents(current.Start())
	nextStart, okNext := helpers.AsParsingComponents(next.Start())
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

// MergeResults merges a date-only and time-only result into a single date-time result.
func (r *AbstractMergeDateTimeRefiner) MergeResults(textBetween string, current, next *kronos.InternalParsingResult, context *kronos.InternalParsingContext) *kronos.InternalParsingResult {
	currentStart, okCurrent := helpers.AsParsingComponents(current.Start())
	if !okCurrent {
		return current
	}

	// Calculate index and text for merged result
	var resultIndex int
	var resultText string
	resultIndex = current.Index()
	resultText = current.Text() + textBetween + next.Text()

	// Determine which is date and which is time, then merge
	var result *kronos.InternalParsingResult
	if currentStart.IsOnlyDate() {
		result = helpers.MergeDateTimeResult(current, next)
	} else {
		result = helpers.MergeDateTimeResult(next, current)
	}

	// Create new result with correct index and text
	startComponents, okStart := helpers.AsParsingComponents(result.Start())
	if !okStart {
		return current
	}
	result = context.CreateParsingResult(resultIndex, resultText, startComponents, nil)

	return result
}

// Refine processes results to merge date and time components based on the pattern between them.
func (r *AbstractMergeDateTimeRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.InternalParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]
		start := current.Index() + len(current.Text())
		end := next.Index()
		textBetween, okRange := helpers.SafeSlice(context.Text(), start, end)
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
func (r *MergeWeekdayComponentRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.InternalParsingResult, 0, len(results))
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
	currentResult *kronos.InternalParsingResult,
	nextResult *kronos.InternalParsingResult,
	context *kronos.InternalParsingContext,
) bool {
	// Merge when:
	// 1. Current result is weekday-only
	// 2. Current result has no certain hour
	// 3. Next result has a certain day
	// 4. Text between is just optional comma and whitespace
	currentStart, okCurrent := helpers.AsParsingComponents(currentResult.Start())
	nextStart, okNext := helpers.AsParsingComponents(nextResult.Start())
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
	currentResult *kronos.InternalParsingResult,
	nextResult *kronos.InternalParsingResult,
	context *kronos.InternalParsingContext,
) *kronos.InternalParsingResult {
	// Get start components
	currentStart, okCurrent := helpers.AsParsingComponents(currentResult.Start())
	nextStart, okNext := helpers.AsParsingComponents(nextResult.Start())
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
	var newEnd *kronos.InternalParsingComponents
	if nextResult.End() != nil {
		if nextEnd, ok := helpers.AsParsingComponents(nextResult.End()); ok {
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

var timezoneOffsetPattern = regexp.MustCompile(`(?i)^\s*(?:\(?(?:GMT|UTC)\s?)?([+-])(\d{1,2})(?::?(\d{2}))?` + `\)?`)

const (
	timezoneOffsetSignGroup         = 1
	timezoneOffsetHourOffsetGroup   = 2
	timezoneOffsetMinuteOffsetGroup = 3
)

// ExtractTimezoneOffsetRefiner extracts timezone offsets from text following a parsed result.
// Examples: "+0900", "-05:00", "GMT+8", "UTC-5"
type ExtractTimezoneOffsetRefiner struct{}

// NewExtractTimezoneOffsetRefiner creates a new ExtractTimezoneOffsetRefiner
func NewExtractTimezoneOffsetRefiner() *ExtractTimezoneOffsetRefiner {
	return &ExtractTimezoneOffsetRefiner{}
}

// Refine extracts timezone offsets and adds them to results
func (r *ExtractTimezoneOffsetRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	for i, result := range results {
		resultStart, ok := helpers.AsParsingComponents(result.Start())
		if !ok {
			continue
		}
		if resultStart.IsCertain(kronos.ComponentTimezoneOffset) {
			continue
		}

		// Calculate the position after the result text
		suffixStart := result.Index() + len(result.Text())
		// Check if we're beyond the end of the text
		suffix, okSlice := helpers.SafeSlice(context.Text(), suffixStart, len(context.Text()))
		if !okSlice {
			continue
		}
		match := timezoneOffsetPattern.FindStringSubmatch(suffix)
		if match == nil {
			continue
		}

		if context.Option().Debug != nil {
			context.Debug(func() {
				// Log: Extracting timezone offset
			})
		}

		hourOffset, err := strconv.Atoi(match[timezoneOffsetHourOffsetGroup])
		if err != nil {
			continue
		}
		minuteOffset := 0
		if match[timezoneOffsetMinuteOffsetGroup] != "" {
			minuteOffset, err = strconv.Atoi(match[timezoneOffsetMinuteOffsetGroup])
			if err != nil {
				continue
			}
		}

		timezoneOffset := hourOffset*60 + minuteOffset

		// No timezones have offsets greater than 14 hours, so disregard this match
		if timezoneOffset > 14*60 {
			continue
		}

		if match[timezoneOffsetSignGroup] == "-" {
			timezoneOffset = -timezoneOffset
		}

		var resultEnd *kronos.InternalParsingComponents
		if result.End() != nil {
			if endComponents, okEnd := helpers.AsParsingComponents(result.End()); okEnd {
				resultEnd = endComponents
				resultEnd.Assign(kronos.ComponentTimezoneOffset, timezoneOffset)
			}
		}

		resultStart.Assign(kronos.ComponentTimezoneOffset, timezoneOffset)

		// Create new result with updated text
		// The match includes leading whitespace from the regex, so append it directly
		newText := result.Text() + match[0]
		newResult := context.CreateParsingResult(result.Index(), newText, resultStart, resultEnd)
		results[i] = newResult
	}

	return results
}
