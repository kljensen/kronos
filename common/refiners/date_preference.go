package refiners

import (
	"time"

	kronos "github.com/kljensen/kronos"
)

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
func (r *DatePreferenceRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	preference := context.Option().Preference

	// PreferCurrentPeriod is the default - no adjustment needed
	if preference == kronos.PreferCurrentPeriod {
		return results
	}

	refDate := context.Reference().GetDateWithAdjustedTimezone()
	refInstant := context.Reference().Instant()

	for _, result := range results {
		resultStart, okStart := kronos.AsParsingComponents(result.Start())
		if !okStart {
			continue
		}

		// Handle time-only results (e.g., "10:00", "3pm")
		// These have only time components certain, no date components
		if resultStart.IsOnlyTime() {
			r.adjustTimeOnlyComponents(resultStart, refDate, refInstant, preference, context)

			// Also adjust end time if it's a time range
			if result.End() != nil {
				if resultEnd, okEnd := kronos.AsParsingComponents(result.End()); okEnd {
					if resultEnd.IsOnlyTime() {
						r.adjustTimeOnlyComponents(resultEnd, refDate, refInstant, preference, context)

						// Ensure end time is after start time
						// If end is before start, it means the range crosses midnight
						if resultStart.Date().After(resultEnd.Date()) {
							// Move end time to next day
							nextDay := resultEnd.Date().AddDate(0, 0, 1)
							kronos.ImplySimilarDate(resultEnd, nextDay)
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
	components *kronos.ParsingComponents,
	refDate time.Time,
	refInstant time.Time,
	preference kronos.DatePreference,
	context *kronos.ParsingContext,
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
			kronos.ImplySimilarDate(components, refPreviousDay)

			if context.Option().Debug != nil {
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
			kronos.ImplySimilarDate(components, refFollowingDay)

			if context.Option().Debug != nil {
				context.Debug(func() {
					// Log: DatePreferenceRefiner (PreferFuture) adjusted time to following day
				})
			}
		}
	}
}
