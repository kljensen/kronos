//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"regexp"
	"strings"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// OverlapRemovalRefiner removes overlapping parse results.
// When two results overlap, it keeps the longer/more specific one.
type OverlapRemovalRefiner struct{}

// NewOverlapRemovalRefiner creates a new OverlapRemovalRefiner
func NewOverlapRemovalRefiner() *OverlapRemovalRefiner {
	return &OverlapRemovalRefiner{}
}

// Refine removes overlapping results
func (r *OverlapRemovalRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	if len(results) < 2 {
		return results
	}

	filteredResults := make([]*kronos.InternalParsingResult, 0, len(results))
	prevResult := results[0]

	for i := 1; i < len(results); i++ {
		result := results[i]

		// Check if results overlap
		if result.Index() >= prevResult.Index()+len(prevResult.Text()) {
			// No overlap, keep previous result and move to current
			filteredResults = append(filteredResults, prevResult)
			prevResult = result
			continue
		}

		// Results overlap - keep the longer one
		var kept *kronos.InternalParsingResult
		if len(result.Text()) > len(prevResult.Text()) {
			kept = result
		} else {
			kept = prevResult
		}

		if context.Option().DebugHandler != nil {
			context.Debug(func() {
				// Log: OverlapRemovalRefiner removing result
			})
		}

		prevResult = kept
	}

	// Add the last result
	if prevResult != nil {
		filteredResults = append(filteredResults, prevResult)
	}

	return filteredResults
}

// UnlikelyFormatFilter filters out unlikely or impossible date/time results.
// It removes results that are just numbers, have invalid dates, or unrealistic times.
type UnlikelyFormatFilter struct {
	strictMode bool
}

// NewUnlikelyFormatFilter creates a new UnlikelyFormatFilter
func NewUnlikelyFormatFilter(strictMode bool) *UnlikelyFormatFilter {
	return &UnlikelyFormatFilter{
		strictMode: strictMode,
	}
}

// Refine filters out unlikely results
func (f *UnlikelyFormatFilter) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	filtered := make([]*kronos.InternalParsingResult, 0, len(results))

	for _, result := range results {
		if f.isValid(context, result) {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

func (f *UnlikelyFormatFilter) isValid(context *kronos.InternalParsingContext, result *kronos.InternalParsingResult) bool {
	// Remove results that are just numbers or dots
	// Exception: Allow 4-digit years (1000-2999) which are valid year-only expressions
	textWithoutSpaces := strings.ReplaceAll(result.Text(), " ", "")
	if regexp.MustCompile(`^\d*(\.\d*)?$`).MatchString(textWithoutSpaces) {
		// Check if it's a 4-digit year
		if regexp.MustCompile(`^[12]\d{3}$`).MatchString(textWithoutSpaces) {
			// Check if the result has a year-level period
			if resultStart, okStart := helpers.AsParsingComponents(result.Start()); okStart {
				if resultStart.Period() == kronos.PeriodYear {
					// This is a valid year-only expression, don't filter it
					return true
				}
			}
		}
		if context.Option().DebugHandler != nil {
			context.Debug(func() {
				// Log: Removing unlikely result
			})
		}
		return false
	}

	// Check if start date is valid
	resultStart, okStart := helpers.AsParsingComponents(result.Start())
	if !okStart {
		return false
	}
	if !resultStart.IsValidDate() {
		if context.Option().DebugHandler != nil {
			context.Debug(func() {
				// Log: Removing invalid result
			})
		}
		return false
	}

	// Check if end date is valid
	if result.End() != nil {
		if resultEnd, okEnd := helpers.AsParsingComponents(result.End()); okEnd {
			if !resultEnd.IsValidDate() {
				if context.Option().DebugHandler != nil {
					context.Debug(func() {
						// Log: Removing invalid result with invalid end date
					})
				}
				return false
			}
		} else {
			if context.Option().DebugHandler != nil {
				context.Debug(func() {
					// Log: Removing invalid result with invalid end date
				})
			}
			return false
		}
	}

	// Apply strict mode checks
	if f.strictMode {
		return f.isStrictModeValid(context, result)
	}

	return true
}

func (f *UnlikelyFormatFilter) isStrictModeValid(context *kronos.InternalParsingContext, result *kronos.InternalParsingResult) bool {
	// In strict mode, remove weekday-only components
	resultStart, okStart := helpers.AsParsingComponents(result.Start())
	if !okStart {
		return false
	}
	if resultStart.IsOnlyWeekdayComponent() {
		if context.Option().DebugHandler != nil {
			context.Debug(func() {
				// Log: (Strict) Removing weekday only component
			})
		}
		return false
	}

	return true
}

// DatePreferenceRefiner adjusts ambiguous dates based on the Preference setting.
// This handles:
// - Time-only dates (e.g., "10:00") - adjusts the implied date component
// - Dates with missing year (e.g., "March 15") - handled by FindYearClosestToRefWithPreference
//
// The refiner respects the following preferences:
// - PreferCurrentPeriod: Keep dates in the current period (default behavior, no adjustment needed)
// - PreferPast: Adjust dates to be in the past
// - PreferFuture: Adjust dates to be in the future
type DatePreferenceRefiner struct{}

// NewDatePreferenceRefiner creates a new DatePreferenceRefiner
func NewDatePreferenceRefiner() *DatePreferenceRefiner {
	return &DatePreferenceRefiner{}
}

// Refine adjusts dates based on the preference setting
func (r *DatePreferenceRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	preference := context.Option().PreferDatesFrom

	// PreferCurrentPeriod is the default - no adjustment needed
	if preference == kronos.PreferCurrentPeriod {
		return results
	}

	refDate := context.Reference().GetDateWithAdjustedTimezone()
	refInstant := context.Reference().Instant()

	for _, result := range results {
		resultStart, okStart := helpers.AsParsingComponents(result.Start())
		if !okStart {
			continue
		}

		// Handle time-only results (e.g., "10:00", "3pm")
		// These have only time components certain, no date components
		if resultStart.IsOnlyTime() {
			r.adjustTimeOnlyComponents(resultStart, refDate, refInstant, preference, context)

			// Also adjust end time if it's a time range
			if result.End() != nil {
				if resultEnd, okEnd := helpers.AsParsingComponents(result.End()); okEnd {
					if resultEnd.IsOnlyTime() {
						r.adjustTimeOnlyComponents(resultEnd, refDate, refInstant, preference, context)

						// Ensure end time is after start time
						// If end is before start, it means the range crosses midnight
						if resultStart.Date().After(resultEnd.Date()) {
							// Move end time to next day
							nextDay := resultEnd.Date().AddDate(0, 0, 1)
							helpers.ImplySimilarDate(resultEnd, nextDay)
						}
					}
				}
			}
		}

		// Note: Dates with missing year are already handled by FindYearClosestToRefWithPreference
		// in the parsers themselves, so we don't need to adjust them here.
	}

	return results
}

// adjustTimeOnlyComponents adjusts time-only components based on preference
func (r *DatePreferenceRefiner) adjustTimeOnlyComponents(
	components *kronos.InternalParsingComponents,
	refDate time.Time,
	refInstant time.Time,
	preference kronos.DatePreference,
	context *kronos.InternalParsingContext,
) {
	// Get the constructed date with the current implied date
	constructedDate := components.Date()

	switch preference {
	case kronos.PreferPast:
		// If the constructed time is after the reference instant, move it to the previous day
		// This ensures the time is in the past
		if constructedDate.After(refInstant) {
			refPreviousDay := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, refDate.Location())
			refPreviousDay = refPreviousDay.AddDate(0, 0, -1)
			helpers.ImplySimilarDate(components, refPreviousDay)

			if context.Option().DebugHandler != nil {
				context.Debug(func() {
					// Log: DatePreferenceRefiner (PreferPast) adjusted time to previous day
				})
			}
		}

	case kronos.PreferFuture:
		// If the constructed time is before or equal to the reference instant, move it to the next day
		// This ensures the time is in the future
		if constructedDate.Before(refInstant) || constructedDate.Equal(refInstant) {
			refFollowingDay := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, refDate.Location())
			refFollowingDay = refFollowingDay.AddDate(0, 0, 1)
			helpers.ImplySimilarDate(components, refFollowingDay)

			if context.Option().DebugHandler != nil {
				context.Debug(func() {
					// Log: DatePreferenceRefiner (PreferFuture) adjusted time to following day
				})
			}
		}
	}
}
