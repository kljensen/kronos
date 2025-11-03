package refiners

import (
	"time"

	kronos "github.com/kljensen/kronos"
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
func (r *ForwardDateRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	if !context.Option.ForwardDate {
		return results
	}

	for _, result := range results {
		refDate := context.Reference.GetDateWithAdjustedTimezone()

		// Handle time-only results
		if result.Start.IsOnlyTime() && context.Reference.Instant.After(result.Start.Date()) {
			refFollowingDay := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, refDate.Location())
			refFollowingDay = refFollowingDay.AddDate(0, 0, 1)

			kronos.ImplySimilarDate(result.Start, refFollowingDay)
			if context.Option.Debug {
				context.DebugLog("ForwardDateRefiner adjusted %s time from ref date (%s) to following day (%s)",
					result, refDate, refFollowingDay)
			}

			if result.End != nil && result.End.IsOnlyTime() {
				kronos.ImplySimilarDate(result.End, refFollowingDay)
				if result.Start.Date().After(result.End.Date()) {
					refFollowingDay = refFollowingDay.AddDate(0, 0, 1)
					kronos.ImplySimilarDate(result.End, refFollowingDay)
				}
			}
		}

		// Handle weekday-only results
		if result.Start.IsOnlyWeekdayComponent() && refDate.After(result.Start.Date()) {
			daysToAdd := int(result.Start.Get(kronos.ComponentWeekday)) - int(refDate.Weekday())
			if daysToAdd <= 0 {
				daysToAdd += 7
			}

			adjustedDate := kronos.AddDuration(refDate, kronos.Duration{Day: daysToAdd})
			kronos.ImplySimilarDate(result.Start, adjustedDate)

			if context.Option.Debug {
				context.DebugLog("ForwardDateRefiner adjusted %s weekday (%s)", result, result.Start)
			}

			if result.End != nil && result.End.IsOnlyWeekdayComponent() {
				daysToAdd = int(result.End.Get(kronos.ComponentWeekday)) - int(adjustedDate.Weekday())
				if daysToAdd <= 0 {
					daysToAdd += 7
				}
				adjustedDate = kronos.AddDuration(adjustedDate, kronos.Duration{Day: daysToAdd})
				kronos.ImplySimilarDate(result.End, adjustedDate)

				if context.Option.Debug {
					context.DebugLog("ForwardDateRefiner adjusted %s weekday (%s)", result, result.End)
				}
			}
		}

		// Handle dates with unknown year
		if result.Start.IsDateWithUnknownYear() && refDate.After(result.Start.Date()) {
			for i := 0; i < 3 && refDate.After(result.Start.Date()); i++ {
				result.Start.Imply(kronos.ComponentYear, result.Start.Get(kronos.ComponentYear)+1)
				if context.Option.Debug {
					context.DebugLog("ForwardDateRefiner adjusted %s year (%s)", result, result.Start)
				}

				if result.End != nil && !result.End.IsCertain(kronos.ComponentYear) {
					result.End.Imply(kronos.ComponentYear, result.End.Get(kronos.ComponentYear)+1)
					if context.Option.Debug {
						context.DebugLog("ForwardDateRefiner adjusted %s year (%s)", result, result.End)
					}
				}
			}
		}
	}

	return results
}
