package kronos

import "time"

// Duration represents a directed time duration as a set of values by timeunits.
// Positive values mean the duration goes into the future.
// Duration supports fractional values (e.g., 1.5 months).
type Duration map[Timeunit]float64

// EmptyDuration represents an explicit empty duration.
// This is defined as zero day, second, and millisecond.
var EmptyDuration = Duration{
	TimeunitDay:         0,
	TimeunitSecond:      0,
	TimeunitMillisecond: 0,
}

// AddDuration returns the date after adding the given duration to ref.
// It handles fractional durations by cascading remainders to smaller units.
// For example, 1.5 months becomes 1 month + 2 weeks.
func AddDuration(ref time.Time, duration Duration) time.Time {
	date := ref

	// Create a working copy to handle fractional cascading
	working := make(Duration)
	for k, v := range duration {
		working[k] = v
	}

	// Process decades (convert to years)
	if val, exists := working[TimeunitDecade]; exists {
		const yearsPerDecade = 10
		working[TimeunitYear] = working[TimeunitYear] + val*yearsPerDecade
	}

	// Process years (cascade fractional part to months)
	if val, exists := working[TimeunitYear]; exists {
		floor := int(val)
		date = addYears(date, floor)
		remainder := val - float64(floor)
		if remainder > 0 {
			working[TimeunitMonth] = working[TimeunitMonth] + remainder*MonthsPerYear
		}
	}

	// Process quarters (convert to months)
	if val, exists := working[TimeunitQuarter]; exists {
		floor := int(val)
		date = addMonths(date, floor*MonthsPerQuarter)
	}

	// Process months (cascade fractional part to weeks)
	if val, exists := working[TimeunitMonth]; exists {
		floor := int(val)
		date = addMonths(date, floor)
		remainder := val - float64(floor)
		if remainder > 0 {
			working[TimeunitWeek] = working[TimeunitWeek] + remainder*WeeksPerMonthApprox
		}
	}

	// Process weeks (cascade fractional part to days)
	if val, exists := working[TimeunitWeek]; exists {
		floor := int(val)
		date = date.AddDate(0, 0, floor*DaysPerWeek)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional weeks to days
			// Round to nearest day for both positive and negative values
			days := remainder * DaysPerWeek
			const roundingOffset = 0.5
			if days > 0 {
				working[TimeunitDay] = working[TimeunitDay] + float64(int(days+roundingOffset))
			} else if days < 0 {
				working[TimeunitDay] = working[TimeunitDay] + float64(int(days-roundingOffset))
			}
		}
	}

	// Process days (cascade fractional part to hours)
	if val, exists := working[TimeunitDay]; exists {
		floor := int(val)
		date = date.AddDate(0, 0, floor)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional days to hours
			// Round to nearest hour for both positive and negative values
			hours := remainder * HoursPerDay
			const roundingOffset = 0.5
			if hours > 0 {
				working[TimeunitHour] = working[TimeunitHour] + float64(int(hours+roundingOffset))
			} else if hours < 0 {
				working[TimeunitHour] = working[TimeunitHour] + float64(int(hours-roundingOffset))
			}
		}
	}

	// Process hours (cascade fractional part to minutes)
	if val, exists := working[TimeunitHour]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Hour)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional hours to minutes
			// For positive values: round to nearest minute
			// For negative values: preserve sign and round to nearest minute
			minutes := remainder * MinutesPerHour
			const roundingOffset = 0.5
			if minutes > 0 {
				working[TimeunitMinute] = working[TimeunitMinute] + float64(int(minutes+roundingOffset))
			} else if minutes < 0 {
				working[TimeunitMinute] = working[TimeunitMinute] + float64(int(minutes-roundingOffset))
			}
		}
	}

	// Process minutes (cascade fractional part to seconds)
	if val, exists := working[TimeunitMinute]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Minute)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional minutes to seconds
			// For positive values: round to nearest second
			// For negative values: preserve sign and round to nearest second
			seconds := remainder * SecondsPerMinute
			const roundingOffset = 0.5
			if seconds > 0 {
				working[TimeunitSecond] = working[TimeunitSecond] + float64(int(seconds+roundingOffset))
			} else if seconds < 0 {
				working[TimeunitSecond] = working[TimeunitSecond] + float64(int(seconds-roundingOffset))
			}
		}
	}

	// Process seconds (cascade fractional part to milliseconds)
	if val, exists := working[TimeunitSecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Second)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional seconds to milliseconds
			// For positive values: round to nearest millisecond
			// For negative values: preserve sign and round to nearest millisecond
			milliseconds := remainder * MillisecondsPerSecond
			const roundingOffset = 0.5
			if milliseconds > 0 {
				working[TimeunitMillisecond] = working[TimeunitMillisecond] + float64(int(milliseconds+roundingOffset))
			} else if milliseconds < 0 {
				working[TimeunitMillisecond] = working[TimeunitMillisecond] + float64(int(milliseconds-roundingOffset))
			}
		}
	}

	// Process milliseconds (cascade fractional part to microseconds)
	if val, exists := working[TimeunitMillisecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Millisecond)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional milliseconds to microseconds
			microseconds := remainder * MicrosecondsPerMS
			const roundingOffset = 0.5
			if microseconds > 0 {
				working[TimeunitMicrosecond] = working[TimeunitMicrosecond] + float64(int(microseconds+roundingOffset))
			} else if microseconds < 0 {
				working[TimeunitMicrosecond] = working[TimeunitMicrosecond] + float64(int(microseconds-roundingOffset))
			}
		}
	}

	// Process microseconds (cascade fractional part to nanoseconds)
	if val, exists := working[TimeunitMicrosecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Microsecond)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional microseconds to nanoseconds
			nanoseconds := remainder * NanosecondsPerMicro
			const roundingOffset = 0.5
			if nanoseconds > 0 {
				working[TimeunitNanosecond] = working[TimeunitNanosecond] + float64(int(nanoseconds+roundingOffset))
			} else if nanoseconds < 0 {
				working[TimeunitNanosecond] = working[TimeunitNanosecond] + float64(int(nanoseconds-roundingOffset))
			}
		}
	}

	// Process nanoseconds
	if val, exists := working[TimeunitNanosecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Nanosecond)
	}

	return date
}

// ReverseDuration returns the reversed duration (e.g., back into the past instead of future).
// All values in the duration are negated.
func ReverseDuration(duration Duration) Duration {
	reversed := make(Duration)
	for key, val := range duration {
		reversed[key] = -val
	}
	return reversed
}

// addYears adds the specified number of years to the date.
// It handles month overflow correctly (e.g., Jan 31 + 1 year = Jan 31 next year).
func addYears(date time.Time, years int) time.Time {
	return date.AddDate(years, 0, 0)
}

// addMonths adds the specified number of months to the date.
// It handles varying month lengths correctly (e.g., Jan 31 + 1 month = Feb 28/29).
func addMonths(date time.Time, months int) time.Time {
	return date.AddDate(0, months, 0)
}
