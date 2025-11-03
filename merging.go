package kronos

import "time"

// MergeDateTimeResult merges a date-only result with a time-only result.
func MergeDateTimeResult(dateResult, timeResult *ParsingResult) *ParsingResult {
	result := dateResult.Clone()
	beginDate := dateResult.Start().(*ParsingComponents)
	beginTime := timeResult.Start().(*ParsingComponents)

	result.start = MergeDateTimeComponent(beginDate, beginTime)

	if dateResult.End() != nil || timeResult.End() != nil {
		var endDate, endTime *ParsingComponents
		if dateResult.End() == nil {
			endDate = dateResult.Start().(*ParsingComponents)
		} else {
			endDate = dateResult.End().(*ParsingComponents)
		}
		if timeResult.End() == nil {
			endTime = timeResult.Start().(*ParsingComponents)
		} else {
			endTime = timeResult.End().(*ParsingComponents)
		}

		endDateTime := MergeDateTimeComponent(endDate, endTime)

		// If date has no end and the merged end time is before start time,
		// the end should be on the next day
		if dateResult.End() == nil && endDateTime.Date().Before(result.Start().Date()) {
			nextDay := endDateTime.Date().Add(24 * time.Hour)
			if endDateTime.IsCertain(ComponentDay) {
				AssignSimilarDate(endDateTime, nextDay)
			} else {
				ImplySimilarDate(endDateTime, nextDay)
			}
		}

		result.end = endDateTime
	}

	return result
}

// MergeDateTimeComponent merges date and time components.
func MergeDateTimeComponent(dateComp, timeComp *ParsingComponents) *ParsingComponents {
	result := dateComp.Clone()

	// Merge time components
	if timeComp.IsCertain(ComponentHour) {
		result.Assign(ComponentHour, *timeComp.Get(ComponentHour))
		result.Assign(ComponentMinute, *timeComp.Get(ComponentMinute))

		if timeComp.IsCertain(ComponentSecond) {
			result.Assign(ComponentSecond, *timeComp.Get(ComponentSecond))
			if timeComp.IsCertain(ComponentMillisecond) {
				result.Assign(ComponentMillisecond, *timeComp.Get(ComponentMillisecond))
			} else {
				result.Imply(ComponentMillisecond, *timeComp.Get(ComponentMillisecond))
			}
		} else {
			result.Imply(ComponentSecond, *timeComp.Get(ComponentSecond))
			result.Imply(ComponentMillisecond, *timeComp.Get(ComponentMillisecond))
		}
	} else {
		result.Imply(ComponentHour, *timeComp.Get(ComponentHour))
		result.Imply(ComponentMinute, *timeComp.Get(ComponentMinute))
		result.Imply(ComponentSecond, *timeComp.Get(ComponentSecond))
		result.Imply(ComponentMillisecond, *timeComp.Get(ComponentMillisecond))
	}

	// Merge timezone
	if timeComp.IsCertain(ComponentTimezoneOffset) {
		result.Assign(ComponentTimezoneOffset, *timeComp.Get(ComponentTimezoneOffset))
	}

	// Merge meridiem
	if timeComp.IsCertain(ComponentMeridiem) {
		result.Assign(ComponentMeridiem, *timeComp.Get(ComponentMeridiem))
	} else if timeComp.Get(ComponentMeridiem) != nil && *timeComp.Get(ComponentMeridiem) != 0 && result.Get(ComponentMeridiem) != nil && *result.Get(ComponentMeridiem) == 0 {
		result.Imply(ComponentMeridiem, *timeComp.Get(ComponentMeridiem))
	}

	// Apply PM meridiem adjustment
	if result.Get(ComponentMeridiem) != nil && *result.Get(ComponentMeridiem) == int(MeridiemPM) && result.Get(ComponentHour) != nil && *result.Get(ComponentHour) < 12 {
		if timeComp.IsCertain(ComponentHour) {
			result.Assign(ComponentHour, *result.Get(ComponentHour)+12)
		} else {
			result.Imply(ComponentHour, *result.Get(ComponentHour)+12)
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
