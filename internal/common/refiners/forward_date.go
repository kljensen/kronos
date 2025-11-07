//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// ForwardDateRefiner enforces the 'forwardDate' option on results.
// When there are missing components (e.g., "March 12-13" without year, or "Thursday"),
// it adjusts the result to be in the future rather than the past.
type ForwardDateRefiner struct{}

// NewForwardDateRefiner creates a new ForwardDateRefiner
func NewForwardDateRefiner() *ForwardDateRefiner {
	return &ForwardDateRefiner{}
}

// Refine adjusts dates to be in the future when forwardDate option is enabled
func (r *ForwardDateRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	if !context.Option().ForwardDate {
		return results
	}

	for _, result := range results {
		refDate := context.Reference().GetDateWithAdjustedTimezone()
		refInstant := context.Reference().Instant()

		// Handle time-only results
		resultStart, okStart := helpers.AsParsingComponents(result.Start())
		if !okStart {
			continue
		}
		if resultStart.IsOnlyTime() && refInstant.After(resultStart.Date()) {
			refFollowingDay := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, refDate.Location())
			refFollowingDay = refFollowingDay.AddDate(0, 0, 1)

			helpers.ImplySimilarDate(resultStart, refFollowingDay)
			if context.Option().Debug != nil {
				context.Debug(func() {
					// Log: ForwardDateRefiner adjusted time from ref date to following day
				})
			}

			if result.End() != nil {
				if resultEnd, okEnd := helpers.AsParsingComponents(result.End()); okEnd {
					if resultEnd.IsOnlyTime() {
						helpers.ImplySimilarDate(resultEnd, refFollowingDay)
						if resultStart.Date().After(resultEnd.Date()) {
							refFollowingDay = refFollowingDay.AddDate(0, 0, 1)
							helpers.ImplySimilarDate(resultEnd, refFollowingDay)
						}
					}
				}
			}
		}

		// Handle weekday-only results
		if resultStart.IsOnlyWeekdayComponent() && refDate.After(resultStart.Date()) {
			weekdayVal := resultStart.Get(kronos.ComponentWeekday)
			if weekdayVal != nil {
				daysToAdd := *weekdayVal - int(refDate.Weekday())
				if daysToAdd <= 0 {
					daysToAdd += 7
				}

				adjustedDate, err := helpers.AddDuration(refDate, kronos.Duration{kronos.TimeunitDay: float64(daysToAdd)})
				if err != nil {
					// Duration calculation failed - skip this adjustment
					continue
				}
				helpers.ImplySimilarDate(resultStart, adjustedDate)

				if context.Option().Debug != nil {
					context.Debug(func() {
						// Log: ForwardDateRefiner adjusted weekday
					})
				}

				if result.End() != nil {
					if resultEnd, okEnd := helpers.AsParsingComponents(result.End()); okEnd {
						if resultEnd.IsOnlyWeekdayComponent() {
							endWeekdayVal := resultEnd.Get(kronos.ComponentWeekday)
							if endWeekdayVal != nil {
								daysToAdd = *endWeekdayVal - int(adjustedDate.Weekday())
								if daysToAdd <= 0 {
									daysToAdd += 7
								}
								adjustedDate, err = helpers.AddDuration(adjustedDate, kronos.Duration{kronos.TimeunitDay: float64(daysToAdd)})
								if err != nil {
									// Duration calculation failed - skip this adjustment
									continue
								}
								helpers.ImplySimilarDate(resultEnd, adjustedDate)

								if context.Option().Debug != nil {
									context.Debug(func() {
										// Log: ForwardDateRefiner adjusted end weekday
									})
								}
							}
						}
					}
				}
			}
		}

		// Handle dates with unknown year
		if resultStart.IsDateWithUnknownYear() && refDate.After(resultStart.Date()) {
			for i := 0; i < 3 && refDate.After(resultStart.Date()); i++ {
				yearVal := resultStart.Get(kronos.ComponentYear)
				if yearVal != nil {
					resultStart.Imply(kronos.ComponentYear, *yearVal+1)
					if context.Option().Debug != nil {
						context.Debug(func() {
							// Log: ForwardDateRefiner adjusted year
						})
					}

					if result.End() != nil {
						if resultEnd, okEnd := helpers.AsParsingComponents(result.End()); okEnd {
							if !resultEnd.IsCertain(kronos.ComponentYear) {
								endYearVal := resultEnd.Get(kronos.ComponentYear)
								if endYearVal != nil {
									resultEnd.Imply(kronos.ComponentYear, *endYearVal+1)
									if context.Option().Debug != nil {
										context.Debug(func() {
											// Log: ForwardDateRefiner adjusted end year
										})
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return results
}
