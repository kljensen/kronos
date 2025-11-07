package helpers

import (
	"github.com/kljensen/kronos"
)

// CreateRelativeFromReference creates a ParsingComponents from a duration relative to the reference.
// It handles date-only durations (implies time) and time durations (assigns both date and time).
// This is used for parsing relative expressions like "in 3 days", "2 hours ago", etc.
// Returns nil if the duration calculation fails (e.g., overflow).
func CreateRelativeFromReference(reference *kronos.InternalReferenceWithTimezone, duration kronos.Duration, emptyDuration kronos.Duration) *kronos.InternalParsingComponents {
	if duration == nil {
		duration = emptyDuration
	}

	date, err := AddDuration(reference.GetDateWithAdjustedTimezone(), duration)
	if err != nil {
		// Duration calculation failed - return nil to indicate invalid result
		return nil
	}

	components := kronos.InternalNewParsingComponents(reference, nil)
	components.AddTag("result/relativeDate")

	// Determine and set the period based on the duration
	period := kronos.InternalDeterminePeriodFromDuration(duration)
	components.SetPeriod(period)

	// Check if duration contains time components
	hasTimeComponents := false
	for _, timeunit := range []kronos.Timeunit{kronos.TimeunitHour, kronos.TimeunitMinute, kronos.TimeunitSecond, kronos.TimeunitMillisecond} {
		if _, exists := duration[timeunit]; exists {
			hasTimeComponents = true
			break
		}
	}

	if hasTimeComponents {
		// Duration includes time - assign both date and time as certain
		components.AddTag("result/relativeDateAndTime")
		components.AssignSimilarTime(date)
		components.AssignSimilarDate(date)
		components.Assign(kronos.ComponentTimezoneOffset, reference.GetTimezoneOffset())
	} else {
		// Duration is date-only - imply time components
		components.ImplySimilarTime(date)
		components.Imply(kronos.ComponentTimezoneOffset, reference.GetTimezoneOffset())

		// Handle different date granularities
		if _, hasDayDuration := duration[kronos.TimeunitDay]; hasDayDuration {
			// Day duration - assign day, month, year and weekday
			components.Assign(kronos.ComponentDay, date.Day())
			components.Assign(kronos.ComponentMonth, int(date.Month()))
			components.Assign(kronos.ComponentYear, date.Year())
			components.Assign(kronos.ComponentWeekday, int(date.Weekday()))
		} else if _, hasWeekDuration := duration[kronos.TimeunitWeek]; hasWeekDuration {
			// Week duration - assign day, month, year and imply weekday
			components.Assign(kronos.ComponentDay, date.Day())
			components.Assign(kronos.ComponentMonth, int(date.Month()))
			components.Assign(kronos.ComponentYear, date.Year())
			components.Imply(kronos.ComponentWeekday, int(date.Weekday()))
		} else {
			// Month/year duration - imply day
			components.Imply(kronos.ComponentDay, date.Day())

			if _, hasMonthDuration := duration[kronos.TimeunitMonth]; hasMonthDuration {
				// Month duration - assign month and year
				components.Assign(kronos.ComponentMonth, int(date.Month()))
				components.Assign(kronos.ComponentYear, date.Year())
			} else {
				// Imply month
				components.Imply(kronos.ComponentMonth, int(date.Month()))

				if _, hasYearDuration := duration[kronos.TimeunitYear]; hasYearDuration {
					// Year duration - assign year
					components.Assign(kronos.ComponentYear, date.Year())
				} else if _, hasQuarterDuration := duration[kronos.TimeunitQuarter]; hasQuarterDuration {
					// Quarter duration - assign year
					components.Assign(kronos.ComponentYear, date.Year())
				} else {
					// Imply year
					components.Imply(kronos.ComponentYear, date.Year())
				}
			}
		}
	}

	return components
}
