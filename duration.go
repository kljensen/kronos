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

	// Process years (cascade fractional part to months)
	if val, exists := working[TimeunitYear]; exists {
		floor := int(val)
		date = addYears(date, floor)
		remainder := val - float64(floor)
		if remainder > 0 {
			working[TimeunitMonth] = working[TimeunitMonth] + remainder*12
		}
	}

	// Process quarters (convert to months)
	if val, exists := working[TimeunitQuarter]; exists {
		floor := int(val)
		date = addMonths(date, floor*3)
	}

	// Process months (cascade fractional part to weeks)
	if val, exists := working[TimeunitMonth]; exists {
		floor := int(val)
		date = addMonths(date, floor)
		remainder := val - float64(floor)
		if remainder > 0 {
			working[TimeunitWeek] = working[TimeunitWeek] + remainder*4
		}
	}

	// Process weeks (cascade fractional part to days)
	if val, exists := working[TimeunitWeek]; exists {
		floor := int(val)
		date = date.AddDate(0, 0, floor*7)
		remainder := val - float64(floor)
		if remainder > 0 {
			// Round to nearest day for week fractions
			working[TimeunitDay] = working[TimeunitDay] + float64(int(remainder*7+0.5))
		}
	}

	// Process days (cascade fractional part to hours)
	if val, exists := working[TimeunitDay]; exists {
		floor := int(val)
		date = date.AddDate(0, 0, floor)
		remainder := val - float64(floor)
		if remainder > 0 {
			// Round to nearest hour for day fractions
			working[TimeunitHour] = working[TimeunitHour] + float64(int(remainder*24+0.5))
		}
	}

	// Process hours (cascade fractional part to minutes)
	if val, exists := working[TimeunitHour]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Hour)
		remainder := val - float64(floor)
		if remainder > 0 {
			// Round to nearest minute for hour fractions
			working[TimeunitMinute] = working[TimeunitMinute] + float64(int(remainder*60+0.5))
		}
	}

	// Process minutes (cascade fractional part to seconds)
	if val, exists := working[TimeunitMinute]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Minute)
		remainder := val - float64(floor)
		if remainder > 0 {
			// Round to nearest second for minute fractions
			working[TimeunitSecond] = working[TimeunitSecond] + float64(int(remainder*60+0.5))
		}
	}

	// Process seconds (cascade fractional part to milliseconds)
	if val, exists := working[TimeunitSecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Second)
		remainder := val - float64(floor)
		if remainder > 0 {
			// Round to nearest millisecond for second fractions
			working[TimeunitMillisecond] = working[TimeunitMillisecond] + float64(int(remainder*1000+0.5))
		}
	}

	// Process milliseconds
	if val, exists := working[TimeunitMillisecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Millisecond)
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
