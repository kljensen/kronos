package helpers

import (
	"time"

	"github.com/kljensen/kronos"
)

// AssignSimilarDate assigns (force updates) the parsing components to the same day as the target.
// This sets year, month, and day as certain (known) values.
func AssignSimilarDate(components *kronos.ParsingComponents, date time.Time) {
	components.Assign(kronos.ComponentDay, date.Day())
	components.Assign(kronos.ComponentMonth, int(date.Month()))
	components.Assign(kronos.ComponentYear, date.Year())
}

// AssignSimilarTime assigns (force updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as certain (known) values.
func AssignSimilarTime(components *kronos.ParsingComponents, date time.Time) {
	components.Assign(kronos.ComponentHour, date.Hour())
	components.Assign(kronos.ComponentMinute, date.Minute())
	components.Assign(kronos.ComponentSecond, date.Second())

	// Break down nanoseconds into milliseconds, microseconds, and nanoseconds
	totalNanos := date.Nanosecond()
	millisecond := totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond := remainingNanos / 1000
	nanosecond := remainingNanos % 1000

	components.Assign(kronos.ComponentMillisecond, millisecond)
	if microsecond > 0 {
		components.Assign(kronos.ComponentMicrosecond, microsecond)
	}
	if nanosecond > 0 {
		components.Assign(kronos.ComponentNanosecond, nanosecond)
	}

	// Set meridiem based on hour
	if date.Hour() < 12 {
		components.Assign(kronos.ComponentMeridiem, 0) // AM
	} else {
		components.Assign(kronos.ComponentMeridiem, 1) // PM
	}
}

// ImplySimilarDate implies (weakly updates) the parsing components to the same day as the target.
// This sets year, month, and day as implied values (only if not already certain).
func ImplySimilarDate(components *kronos.ParsingComponents, date time.Time) {
	components.Imply(kronos.ComponentDay, date.Day())
	components.Imply(kronos.ComponentMonth, int(date.Month()))
	components.Imply(kronos.ComponentYear, date.Year())
}

// ImplySimilarTime implies (weakly updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as implied values (only if not already certain).
func ImplySimilarTime(components *kronos.ParsingComponents, date time.Time) {
	components.Imply(kronos.ComponentHour, date.Hour())
	components.Imply(kronos.ComponentMinute, date.Minute())
	components.Imply(kronos.ComponentSecond, date.Second())

	// Break down nanoseconds into milliseconds, microseconds, and nanoseconds
	totalNanos := date.Nanosecond()
	millisecond := totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond := remainingNanos / 1000
	nanosecond := remainingNanos % 1000

	components.Imply(kronos.ComponentMillisecond, millisecond)
	if microsecond > 0 {
		components.Imply(kronos.ComponentMicrosecond, microsecond)
	}
	if nanosecond > 0 {
		components.Imply(kronos.ComponentNanosecond, nanosecond)
	}

	// Set meridiem based on hour
	if date.Hour() < 12 {
		components.Imply(kronos.ComponentMeridiem, 0) // AM
	} else {
		components.Imply(kronos.ComponentMeridiem, 1) // PM
	}
}

// MergeDateTimeComponent merges date and time components.
func MergeDateTimeComponent(dateComp, timeComp *kronos.ParsingComponents) *kronos.ParsingComponents {
	result := dateComp.Clone()

	// Merge time components
	hourVal := timeComp.Get(kronos.ComponentHour)
	minuteVal := timeComp.Get(kronos.ComponentMinute)
	secondVal := timeComp.Get(kronos.ComponentSecond)
	millisecondVal := timeComp.Get(kronos.ComponentMillisecond)
	microsecondVal := timeComp.Get(kronos.ComponentMicrosecond)
	nanosecondVal := timeComp.Get(kronos.ComponentNanosecond)

	if timeComp.IsCertain(kronos.ComponentHour) {
		if hourVal != nil {
			result.Assign(kronos.ComponentHour, *hourVal)
		}
		if minuteVal != nil {
			result.Assign(kronos.ComponentMinute, *minuteVal)
		}

		if timeComp.IsCertain(kronos.ComponentSecond) {
			if secondVal != nil {
				result.Assign(kronos.ComponentSecond, *secondVal)
			}
			if millisecondVal != nil {
				if timeComp.IsCertain(kronos.ComponentMillisecond) {
					result.Assign(kronos.ComponentMillisecond, *millisecondVal)
				} else {
					result.Imply(kronos.ComponentMillisecond, *millisecondVal)
				}
			}
			if microsecondVal != nil {
				if timeComp.IsCertain(kronos.ComponentMicrosecond) {
					result.Assign(kronos.ComponentMicrosecond, *microsecondVal)
				} else {
					result.Imply(kronos.ComponentMicrosecond, *microsecondVal)
				}
			}
			if nanosecondVal != nil {
				if timeComp.IsCertain(kronos.ComponentNanosecond) {
					result.Assign(kronos.ComponentNanosecond, *nanosecondVal)
				} else {
					result.Imply(kronos.ComponentNanosecond, *nanosecondVal)
				}
			}
		} else {
			if secondVal != nil {
				result.Imply(kronos.ComponentSecond, *secondVal)
			}
			if millisecondVal != nil {
				result.Imply(kronos.ComponentMillisecond, *millisecondVal)
			}
			if microsecondVal != nil {
				result.Imply(kronos.ComponentMicrosecond, *microsecondVal)
			}
			if nanosecondVal != nil {
				result.Imply(kronos.ComponentNanosecond, *nanosecondVal)
			}
		}
	} else {
		if hourVal != nil {
			result.Imply(kronos.ComponentHour, *hourVal)
		}
		if minuteVal != nil {
			result.Imply(kronos.ComponentMinute, *minuteVal)
		}
		if secondVal != nil {
			result.Imply(kronos.ComponentSecond, *secondVal)
		}
		if millisecondVal != nil {
			result.Imply(kronos.ComponentMillisecond, *millisecondVal)
		}
		if microsecondVal != nil {
			result.Imply(kronos.ComponentMicrosecond, *microsecondVal)
		}
		if nanosecondVal != nil {
			result.Imply(kronos.ComponentNanosecond, *nanosecondVal)
		}
	}

	// Merge timezone
	if timeComp.IsCertain(kronos.ComponentTimezoneOffset) {
		if tzVal := timeComp.Get(kronos.ComponentTimezoneOffset); tzVal != nil {
			result.Assign(kronos.ComponentTimezoneOffset, *tzVal)
		}
	}

	// Merge meridiem
	timeMeridiem := timeComp.Get(kronos.ComponentMeridiem)
	resultMeridiem := result.Get(kronos.ComponentMeridiem)
	if timeComp.IsCertain(kronos.ComponentMeridiem) {
		if timeMeridiem != nil {
			result.Assign(kronos.ComponentMeridiem, *timeMeridiem)
		}
	} else if timeMeridiem != nil && *timeMeridiem != 0 && resultMeridiem != nil && *resultMeridiem == 0 {
		result.Imply(kronos.ComponentMeridiem, *timeMeridiem)
	}

	// Apply PM meridiem adjustment
	// Note: hour 0 (midnight) should not be converted even with PM meridiem
	// Only hours 1-11 should be adjusted for PM (becoming 13-23)
	resultMeridiem = result.Get(kronos.ComponentMeridiem)
	resultHour := result.Get(kronos.ComponentHour)
	if resultMeridiem != nil && *resultMeridiem == 1 && resultHour != nil && *resultHour > 0 && *resultHour < 12 { // PM
		if timeComp.IsCertain(kronos.ComponentHour) {
			result.Assign(kronos.ComponentHour, *resultHour+12)
		} else {
			result.Imply(kronos.ComponentHour, *resultHour+12)
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
	if timePeriod == kronos.PeriodTime || datePeriod == kronos.PeriodTime {
		result.SetPeriod(kronos.PeriodTime)
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
