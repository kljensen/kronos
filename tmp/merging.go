package kronos

import "time"

// MergeDateTimeResult merges a date-only result with a time-only result.
func MergeDateTimeResult(dateResult, timeResult *ParsingResult) *ParsingResult {
	result := dateResult.Clone()
	beginDate := dateResult.Start
	beginTime := timeResult.Start

	result.Start = MergeDateTimeComponent(beginDate, beginTime)

	if dateResult.End != nil || timeResult.End != nil {
		var endDate, endTime *ParsingComponents
		if dateResult.End == nil {
			endDate = dateResult.Start
		} else {
			endDate = dateResult.End
		}
		if timeResult.End == nil {
			endTime = timeResult.Start
		} else {
			endTime = timeResult.End
		}

		endDateTime := MergeDateTimeComponent(endDate, endTime)

		// If date has no end and the merged end time is before start time,
		// the end should be on the next day
		if dateResult.End == nil && endDateTime.Date().Before(result.Start.Date()) {
			nextDay := endDateTime.Date().Add(24 * time.Hour)
			if endDateTime.IsCertain(ComponentDay) {
				AssignSimilarDate(endDateTime, nextDay)
			} else {
				ImplySimilarDate(endDateTime, nextDay)
			}
		}

		result.End = endDateTime
	}

	return result
}

// MergeDateTimeComponent merges date and time components.
func MergeDateTimeComponent(dateComp, timeComp *ParsingComponents) *ParsingComponents {
	result := dateComp.Clone()

	// Merge time components
	if timeComp.IsCertain(ComponentHour) {
		result.Assign(ComponentHour, timeComp.Get(ComponentHour))
		result.Assign(ComponentMinute, timeComp.Get(ComponentMinute))

		if timeComp.IsCertain(ComponentSecond) {
			result.Assign(ComponentSecond, timeComp.Get(ComponentSecond))
			if timeComp.IsCertain(ComponentMillisecond) {
				result.Assign(ComponentMillisecond, timeComp.Get(ComponentMillisecond))
			} else {
				result.Imply(ComponentMillisecond, timeComp.Get(ComponentMillisecond))
			}
		} else {
			result.Imply(ComponentSecond, timeComp.Get(ComponentSecond))
			result.Imply(ComponentMillisecond, timeComp.Get(ComponentMillisecond))
		}
	} else {
		result.Imply(ComponentHour, timeComp.Get(ComponentHour))
		result.Imply(ComponentMinute, timeComp.Get(ComponentMinute))
		result.Imply(ComponentSecond, timeComp.Get(ComponentSecond))
		result.Imply(ComponentMillisecond, timeComp.Get(ComponentMillisecond))
	}

	// Merge timezone
	if timeComp.IsCertain(ComponentTimezoneOffset) {
		result.Assign(ComponentTimezoneOffset, timeComp.Get(ComponentTimezoneOffset))
	}

	// Merge meridiem
	if timeComp.IsCertain(ComponentMeridiem) {
		result.Assign(ComponentMeridiem, timeComp.Get(ComponentMeridiem))
	} else if timeComp.Get(ComponentMeridiem) != 0 && result.Get(ComponentMeridiem) == 0 {
		result.Imply(ComponentMeridiem, timeComp.Get(ComponentMeridiem))
	}

	// Apply PM meridiem adjustment
	if result.Get(ComponentMeridiem) == MeridiemPM && result.Get(ComponentHour) < 12 {
		if timeComp.IsCertain(ComponentHour) {
			result.Assign(ComponentHour, result.Get(ComponentHour)+12)
		} else {
			result.Imply(ComponentHour, result.Get(ComponentHour)+12)
		}
	}

	// Merge tags
	for tag := range dateComp.Tags() {
		result.AddTag(tag)
	}
	for tag := range timeComp.Tags() {
		result.AddTag(tag)
	}

	return result
}
