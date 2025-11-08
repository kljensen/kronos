package kronos

import (
	"time"
)

// mergeDateTimeResult merges a date-only result with a time-only result.
func mergeDateTimeResult(dateResult, timeResult *parsingResult) *parsingResult {
	result := dateResult.Clone()

	beginDate := dateResult.start
	beginTime := timeResult.start

	result.start = mergeDateTimeComponent(beginDate, beginTime)

	if dateResult.end != nil || timeResult.end != nil {
		var endDate, endTime *parsingComponents
		if dateResult.end == nil {
			endDate = dateResult.start
		} else {
			endDate = dateResult.end
		}
		if timeResult.end == nil {
			endTime = timeResult.start
		} else {
			endTime = timeResult.end
		}

		endDateTime := mergeDateTimeComponent(endDate, endTime)

		// If date has no end and the merged end time is before start time,
		// the end should be on the next day
		if dateResult.end == nil && endDateTime.Date().Before(result.start.Date()) {
			nextDay := endDateTime.Date().Add(24 * time.Hour)
			if endDateTime.IsCertain(ComponentDay) {
				endDateTime.AssignSimilarDate(nextDay)
			} else {
				endDateTime.ImplySimilarDate(nextDay)
			}
		}

		result.end = endDateTime
	}

	return result
}

// assignOrImplyComponent assigns or implies a component value based on source certainty.
func assignOrImplyComponent(result, source *parsingComponents, component Component) {
	val := source.Get(component)
	if val == nil {
		return
	}
	if source.IsCertain(component) {
		result.Assign(component, *val)
	} else {
		result.Imply(component, *val)
	}
}

// mergeComponents merges multiple components from source to result.
// Uses assignOrImplyComponent to preserve certainty information.
func mergeComponents(result, source *parsingComponents, components ...Component) {
	for _, component := range components {
		assignOrImplyComponent(result, source, component)
	}
}

// mergeDateTimeComponent merges date and time components.
func mergeDateTimeComponent(dateComp, timeComp *parsingComponents) *parsingComponents {
	result := dateComp.Clone()

	// Merge all time components, preserving certainty
	mergeComponents(result, timeComp,
		ComponentHour,
		ComponentMinute,
		ComponentSecond,
		ComponentMillisecond,
		ComponentMicrosecond,
		ComponentNanosecond,
	)

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
	if resultMeridiem != nil && *resultMeridiem == 1 && resultHour != nil && *resultHour > 0 && *resultHour < 12 { // PM
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
