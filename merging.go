package kronos

import "time"

// MergeDateTimeResult merges a date-only result with a time-only result.
func mergeDateTimeResult(dateResult, timeResult *ParsingResult) *ParsingResult {
	result := dateResult.Clone()
	beginDate, okDate := asParsingComponents(dateResult.Start())
	beginTime, okTime := asParsingComponents(timeResult.Start())
	if !okDate || !okTime {
		return result
	}

	result.start = mergeDateTimeComponent(beginDate, beginTime)

	if dateResult.End() != nil || timeResult.End() != nil {
		var endDate, endTime *ParsingComponents
		if dateResult.End() == nil {
			endDate, okDate = asParsingComponents(dateResult.Start())
		} else {
			endDate, okDate = asParsingComponents(dateResult.End())
		}
		if timeResult.End() == nil {
			endTime, okTime = asParsingComponents(timeResult.Start())
		} else {
			endTime, okTime = asParsingComponents(timeResult.End())
		}
		if !okDate || !okTime {
			return result
		}

		endDateTime := mergeDateTimeComponent(endDate, endTime)

		// If date has no end and the merged end time is before start time,
		// the end should be on the next day
		if dateResult.End() == nil && endDateTime.Date().Before(result.Start().Date()) {
			nextDay := endDateTime.Date().Add(24 * time.Hour)
			if endDateTime.IsCertain(ComponentDay) {
				assignSimilarDate(endDateTime, nextDay)
			} else {
				implySimilarDate(endDateTime, nextDay)
			}
		}

		result.end = endDateTime
	}

	return result
}

// MergeDateTimeComponent merges date and time components.
func mergeDateTimeComponent(dateComp, timeComp *ParsingComponents) *ParsingComponents {
	result := dateComp.Clone()

	// Merge time components
	hourVal := timeComp.Get(ComponentHour)
	minuteVal := timeComp.Get(ComponentMinute)
	secondVal := timeComp.Get(ComponentSecond)
	millisecondVal := timeComp.Get(ComponentMillisecond)
	microsecondVal := timeComp.Get(ComponentMicrosecond)
	nanosecondVal := timeComp.Get(ComponentNanosecond)

	if timeComp.IsCertain(ComponentHour) {
		if hourVal != nil {
			result.Assign(ComponentHour, *hourVal)
		}
		if minuteVal != nil {
			result.Assign(ComponentMinute, *minuteVal)
		}

		if timeComp.IsCertain(ComponentSecond) {
			if secondVal != nil {
				result.Assign(ComponentSecond, *secondVal)
			}
			if millisecondVal != nil {
				if timeComp.IsCertain(ComponentMillisecond) {
					result.Assign(ComponentMillisecond, *millisecondVal)
				} else {
					result.Imply(ComponentMillisecond, *millisecondVal)
				}
			}
			if microsecondVal != nil {
				if timeComp.IsCertain(ComponentMicrosecond) {
					result.Assign(ComponentMicrosecond, *microsecondVal)
				} else {
					result.Imply(ComponentMicrosecond, *microsecondVal)
				}
			}
			if nanosecondVal != nil {
				if timeComp.IsCertain(ComponentNanosecond) {
					result.Assign(ComponentNanosecond, *nanosecondVal)
				} else {
					result.Imply(ComponentNanosecond, *nanosecondVal)
				}
			}
		} else {
			if secondVal != nil {
				result.Imply(ComponentSecond, *secondVal)
			}
			if millisecondVal != nil {
				result.Imply(ComponentMillisecond, *millisecondVal)
			}
			if microsecondVal != nil {
				result.Imply(ComponentMicrosecond, *microsecondVal)
			}
			if nanosecondVal != nil {
				result.Imply(ComponentNanosecond, *nanosecondVal)
			}
		}
	} else {
		if hourVal != nil {
			result.Imply(ComponentHour, *hourVal)
		}
		if minuteVal != nil {
			result.Imply(ComponentMinute, *minuteVal)
		}
		if secondVal != nil {
			result.Imply(ComponentSecond, *secondVal)
		}
		if millisecondVal != nil {
			result.Imply(ComponentMillisecond, *millisecondVal)
		}
		if microsecondVal != nil {
			result.Imply(ComponentMicrosecond, *microsecondVal)
		}
		if nanosecondVal != nil {
			result.Imply(ComponentNanosecond, *nanosecondVal)
		}
	}

	// Merge timezone
	if timeComp.IsCertain(ComponentTimezoneOffset) {
		if tzVal := timeComp.Get(ComponentTimezoneOffset); tzVal != nil {
			result.Assign(ComponentTimezoneOffset, *tzVal)
		}
	}

	// Merge meridiem
	timeMeridiem := timeComp.Get(ComponentMeridiem)
	resultMeridiem := result.Get(ComponentMeridiem)
	if timeComp.IsCertain(ComponentMeridiem) {
		if timeMeridiem != nil {
			result.Assign(ComponentMeridiem, *timeMeridiem)
		}
	} else if timeMeridiem != nil && *timeMeridiem != 0 && resultMeridiem != nil && *resultMeridiem == 0 {
		result.Imply(ComponentMeridiem, *timeMeridiem)
	}

	// Apply PM meridiem adjustment
	// Note: hour 0 (midnight) should not be converted even with PM meridiem
	// Only hours 1-11 should be adjusted for PM (becoming 13-23)
	resultMeridiem = result.Get(ComponentMeridiem)
	resultHour := result.Get(ComponentHour)
	if resultMeridiem != nil && *resultMeridiem == int(MeridiemPM) && resultHour != nil && *resultHour > 0 && *resultHour < 12 {
		if timeComp.IsCertain(ComponentHour) {
			result.Assign(ComponentHour, *resultHour+12)
		} else {
			result.Imply(ComponentHour, *resultHour+12)
		}
	}

	// Merge tags
	for tag := range dateComp.Tags() {
		result.AddTag(tag)
	}
	for tag := range timeComp.Tags() {
		result.AddTag(tag)
	}

	// Merge period - use the finest granularity
	// If time components are present, period should be PeriodTime
	datePeriod := dateComp.Period()
	timePeriod := timeComp.Period()

	// If either is time-level, the result is time-level
	if timePeriod == PeriodTime || datePeriod == PeriodTime {
		result.SetPeriod(PeriodTime)
	} else {
		// Otherwise, use the finer of the two periods
		if datePeriod > timePeriod {
			result.SetPeriod(datePeriod)
		} else {
			result.SetPeriod(timePeriod)
		}
	}

	return result
}
