package kronos

import (
	"fmt"
	"maps"
	"math"
	"time"
)

// Timeunit represents a unit of time for calculations and operations.
// This type is used as a key in the Duration map type for specifying time durations.
//
// Note: While Timeunit is part of the public API (used by Duration), it is primarily
// an implementation detail. Future versions may move this to a more restricted scope
// while maintaining backward compatibility for Duration operations.
type Timeunit string

// Time unit constants for use with Duration type.
// These constants are required for working with Duration maps.
//
// Example:
//
//	duration := kronos.Duration{
//	    kronos.TimeunitDay: 5,
//	    kronos.TimeunitHour: 3,
//	}
const (
	// Large time units
	TimeunitDecade  Timeunit = "decade"
	TimeunitYear    Timeunit = "year"
	TimeunitQuarter Timeunit = "quarter"
	TimeunitMonth   Timeunit = "month"

	// Medium time units
	TimeunitWeek Timeunit = "week"
	TimeunitDay  Timeunit = "day"

	// Small time units
	TimeunitHour        Timeunit = "hour"
	TimeunitMinute      Timeunit = "minute"
	TimeunitSecond      Timeunit = "second"
	TimeunitMillisecond Timeunit = "millisecond"
	TimeunitMicrosecond Timeunit = "microsecond"
	TimeunitNanosecond  Timeunit = "nanosecond"
)

// Bounds constants for date arithmetic.
const (
	minYear            = 1
	maxYear            = 9999
	maxYearsDuration   = 10000
	maxMonthsDuration  = 120000
	maxDaysDuration    = 3650000
	maxHoursDuration   = 87600000
	maxMinutesDuration = 5256000000

	// roundingOffset is used for rounding fractional time units
	roundingOffset = 0.5
)

// Duration represents a directed time duration as a set of values by timeunits.
// Positive values mean the duration goes into the future.
// Duration supports fractional values (e.g., 1.5 months).
type Duration map[Timeunit]float64

// validateDate checks if a date is within valid bounds.
func validateDate(t time.Time) error {
	year := t.Year()
	if year < minYear || year > maxYear {
		return fmt.Errorf("date year %d is outside valid range [%d, %d]", year, minYear, maxYear)
	}
	return nil
}

// validateDuration checks if duration values are within reasonable bounds.
func validateDuration(d Duration) error {
	for unit, value := range d {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("duration contains invalid value for %s: %f", unit, value)
		}

		absValue := math.Abs(value)
		switch unit {
		case TimeunitYear, TimeunitDecade:
			if absValue > maxYearsDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, maxYearsDuration)
			}
		case TimeunitMonth, TimeunitQuarter:
			if absValue > maxMonthsDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, maxMonthsDuration)
			}
		case TimeunitWeek, TimeunitDay:
			if absValue > maxDaysDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, maxDaysDuration)
			}
		case TimeunitHour:
			if absValue > maxHoursDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, maxHoursDuration)
			}
		case TimeunitMinute:
			if absValue > maxMinutesDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, maxMinutesDuration)
			}
		}
	}
	return nil
}

// checkCascadingOverflow checks if cascading operations will cause overflow.
func checkCascadingOverflow(d Duration) error {
	totalYears := d[TimeunitYear] + d[TimeunitDecade]*10
	totalMonths := totalYears*12 + d[TimeunitMonth] + d[TimeunitQuarter]*3
	totalDays := d[TimeunitDay] + d[TimeunitWeek]*7

	if math.Abs(totalYears) > maxYearsDuration {
		return fmt.Errorf("cascading year duration %f exceeds maximum %d", totalYears, maxYearsDuration)
	}
	if math.Abs(totalMonths) > maxMonthsDuration {
		return fmt.Errorf("cascading month duration %f exceeds maximum %d", totalMonths, maxMonthsDuration)
	}
	if math.Abs(totalDays) > maxDaysDuration {
		return fmt.Errorf("cascading day duration %f exceeds maximum %d", totalDays, maxDaysDuration)
	}

	return nil
}

// checkFloatToIntOverflow checks if converting a float to int would overflow.
func checkFloatToIntOverflow(value float64, unit Timeunit) error {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return fmt.Errorf("duration %s value %f would overflow int32 conversion", unit, value)
	}
	return nil
}

// InternalAddDuration is exported for use by internal packages only.
// External code should not use this function directly.
//
// It returns the date after adding the given duration to ref.
// It handles fractional durations by cascading remainders to smaller units.
// For example, 1.5 months becomes 1 month + 2 weeks.
// Returns an error if the duration or resulting date is out of bounds.
func InternalAddDuration(ref time.Time, duration Duration) (time.Time, error) {
	return addDuration(ref, duration)
}

// addDuration is the internal implementation.
func addDuration(ref time.Time, duration Duration) (time.Time, error) {
	// Validate input date
	if err := validateDate(ref); err != nil {
		return time.Time{}, err
	}

	// Validate duration values
	if err := validateDuration(duration); err != nil {
		return time.Time{}, err
	}

	// Check for cascading overflow
	if err := checkCascadingOverflow(duration); err != nil {
		return time.Time{}, err
	}

	date := ref

	// Create a working copy to handle fractional cascading
	working := make(Duration)
	maps.Copy(working, duration)

	// Process decades (convert to years)
	if val, exists := working[TimeunitDecade]; exists {
		const yearsPerDecade = 10
		working[TimeunitYear] += val * yearsPerDecade
	}

	// Process years (cascade fractional part to months)
	if val, exists := working[TimeunitYear]; exists {
		if err := checkFloatToIntOverflow(val, TimeunitYear); err != nil {
			return time.Time{}, err
		}
		floor := int(val)
		date = date.AddDate(floor, 0, 0)
		if err := validateDate(date); err != nil {
			return time.Time{}, err
		}
		remainder := val - float64(floor)
		if remainder != 0 {
			working[TimeunitMonth] += remainder * 12
		}
	}

	// Process quarters (convert to months)
	if val, exists := working[TimeunitQuarter]; exists {
		if err := checkFloatToIntOverflow(val, TimeunitQuarter); err != nil {
			return time.Time{}, err
		}
		floor := int(val)
		date = date.AddDate(0, floor*3, 0)
		if err := validateDate(date); err != nil {
			return time.Time{}, err
		}
	}

	// Process months (cascade fractional part to weeks)
	if val, exists := working[TimeunitMonth]; exists {
		if err := checkFloatToIntOverflow(val, TimeunitMonth); err != nil {
			return time.Time{}, err
		}
		floor := int(val)
		date = date.AddDate(0, floor, 0)
		if err := validateDate(date); err != nil {
			return time.Time{}, err
		}
		remainder := val - float64(floor)
		if remainder != 0 {
			working[TimeunitWeek] += remainder * 4
		}
	}

	// Process weeks (cascade fractional part to days)
	if val, exists := working[TimeunitWeek]; exists {
		floor := int(val)
		date = date.AddDate(0, 0, floor*7)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional weeks to days
			// Round to nearest day for both positive and negative values
			days := remainder * 7
			if days > 0 {
				working[TimeunitDay] += float64(int(days + roundingOffset))
			} else if days < 0 {
				working[TimeunitDay] += float64(int(days - roundingOffset))
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
			hours := remainder * 24
			if hours > 0 {
				working[TimeunitHour] += float64(int(hours + roundingOffset))
			} else if hours < 0 {
				working[TimeunitHour] += float64(int(hours - roundingOffset))
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
			minutes := remainder * 60
			if minutes > 0 {
				working[TimeunitMinute] += float64(int(minutes + roundingOffset))
			} else if minutes < 0 {
				working[TimeunitMinute] += float64(int(minutes - roundingOffset))
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
			seconds := remainder * 60
			if seconds > 0 {
				working[TimeunitSecond] += float64(int(seconds + roundingOffset))
			} else if seconds < 0 {
				working[TimeunitSecond] += float64(int(seconds - roundingOffset))
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
			milliseconds := remainder * 1000
			if milliseconds > 0 {
				working[TimeunitMillisecond] += float64(int(milliseconds + roundingOffset))
			} else if milliseconds < 0 {
				working[TimeunitMillisecond] += float64(int(milliseconds - roundingOffset))
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
			microseconds := remainder * 1000
			if microseconds > 0 {
				working[TimeunitMicrosecond] += float64(int(microseconds + roundingOffset))
			} else if microseconds < 0 {
				working[TimeunitMicrosecond] += float64(int(microseconds - roundingOffset))
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
			nanoseconds := remainder * 1000
			if nanoseconds > 0 {
				working[TimeunitNanosecond] += float64(int(nanoseconds + roundingOffset))
			} else if nanoseconds < 0 {
				working[TimeunitNanosecond] += float64(int(nanoseconds - roundingOffset))
			}
		}
	}

	// Process nanoseconds
	if val, exists := working[TimeunitNanosecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Nanosecond)
	}

	// Final validation of result date
	if err := validateDate(date); err != nil {
		return time.Time{}, err
	}

	return date, nil
}

// InternalReverseDuration is exported for use by internal packages only.
// External code should not use this function directly.
//
// It returns the reversed duration (e.g., back into the past instead of future).
// All values in the duration are negated.
func InternalReverseDuration(duration Duration) Duration {
	return reverseDuration(duration)
}

// reverseDuration is the internal implementation.
func reverseDuration(duration Duration) Duration {
	reversed := make(Duration)
	for key, val := range duration {
		reversed[key] = -val
	}
	return reversed
}
