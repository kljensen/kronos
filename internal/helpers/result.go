package helpers

import (
	"time"

	"github.com/kljensen/kronos"
)

// CreateRelativeFromReference creates a ParsingComponents from a duration relative to the reference.
// It handles date-only durations (implies time) and time durations (assigns both date and time).
// This is used for parsing relative expressions like "in 3 days", "2 hours ago", etc.
// Returns nil if the duration calculation fails (e.g., overflow).
func CreateRelativeFromReference(reference *kronos.ReferenceWithTimezone, duration kronos.Duration, emptyDuration kronos.Duration) *kronos.ParsingComponents {
	if duration == nil {
		duration = emptyDuration
	}

	date, err := AddDuration(reference.GetDateWithAdjustedTimezone(), duration)
	if err != nil {
		// Duration calculation failed - return nil to indicate invalid result
		return nil
	}

	components := kronos.XNewParsingComponents(reference, nil)
	components.AddTag("result/relativeDate")

	// Determine and set the period based on the duration
	period := determinePeriodFromDuration(duration)
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
		AssignSimilarTime(components, date)
		AssignSimilarDate(components, date)
		components.Assign(kronos.ComponentTimezoneOffset, reference.GetTimezoneOffset())
	} else {
		// Duration is date-only - imply time components
		ImplySimilarTime(components, date)
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

// MergeDateTimeResult merges a date-only result with a time-only result.
// NOTE: This function cannot be fully implemented in this package because it requires
// direct field access to ParsingResult.end which is private. The actual implementation
// remains in the kronos package via kronos.XMergeDateTimeResult.
//
// This is a placeholder for documentation purposes. The real function will be
// exposed through the X-prefixed helper in the main kronos package.
func MergeDateTimeResult(dateResult, timeResult *kronos.ParsingResult) *kronos.ParsingResult {
	// This implementation is a simplified version that works with public APIs only.
	// The full implementation with private field access remains in the kronos package.
	result := dateResult.Clone()
	beginDate, okDate := kronos.XAsParsingComponents(dateResult.Start())
	beginTime, okTime := kronos.XAsParsingComponents(timeResult.Start())
	if !okDate || !okTime {
		return result
	}

	result.SetStart(MergeDateTimeComponent(beginDate, beginTime))

	// For end component merging, we rely on the SetStart approach
	// The full implementation with direct field access is in the main package
	if dateResult.End() != nil || timeResult.End() != nil {
		var endDate, endTime *kronos.ParsingComponents
		if dateResult.End() == nil {
			endDate, okDate = kronos.XAsParsingComponents(dateResult.Start())
		} else {
			endDate, okDate = kronos.XAsParsingComponents(dateResult.End())
		}
		if timeResult.End() == nil {
			endTime, okTime = kronos.XAsParsingComponents(timeResult.Start())
		} else {
			endTime, okTime = kronos.XAsParsingComponents(timeResult.End())
		}
		if !okDate || !okTime {
			return result
		}

		endDateTime := MergeDateTimeComponent(endDate, endTime)

		// If date has no end and the merged end time is before start time,
		// the end should be on the next day
		if dateResult.End() == nil && endDateTime.Date().Before(result.Start().Date()) {
			nextDay := endDateTime.Date().Add(24 * time.Hour)
			if endDateTime.IsCertain(kronos.ComponentDay) {
				AssignSimilarDate(endDateTime, nextDay)
			} else {
				ImplySimilarDate(endDateTime, nextDay)
			}
		}

		// Note: Setting end requires creating a new result since there's no public setter
		// The main package implementation uses direct field access
	}

	return result
}

// Helper functions

// determinePeriodFromDuration determines the granularity/period based on a duration.
// The period represents the finest time unit present in the duration.
// This follows the pattern from Python's dateparser.
func determinePeriodFromDuration(duration kronos.Duration) kronos.Period {
	if duration == nil {
		return kronos.PeriodDay // Default
	}

	// Check from finest to coarsest granularity
	// Time components (hour, minute, second) indicate time-level precision
	for _, timeunit := range []kronos.Timeunit{kronos.TimeunitSecond, kronos.TimeunitMinute, kronos.TimeunitHour} {
		if _, exists := duration[timeunit]; exists {
			return kronos.PeriodTime
		}
	}

	// Day indicates day-level precision
	if _, exists := duration[kronos.TimeunitDay]; exists {
		return kronos.PeriodDay
	}

	// Week indicates week-level precision
	if _, exists := duration[kronos.TimeunitWeek]; exists {
		return kronos.PeriodWeek
	}

	// Month indicates month-level precision
	if _, exists := duration[kronos.TimeunitMonth]; exists {
		return kronos.PeriodMonth
	}

	// Year, decade, or quarter indicate year-level precision
	for _, timeunit := range []kronos.Timeunit{kronos.TimeunitYear, kronos.TimeunitDecade, kronos.TimeunitQuarter} {
		if _, exists := duration[timeunit]; exists {
			return kronos.PeriodYear
		}
	}

	// Default to day if no specific duration is found
	return kronos.PeriodDay
}
