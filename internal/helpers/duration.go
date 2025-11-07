package helpers

import (
	"fmt"
	"math"
	"time"

	"github.com/kljensen/kronos"
)

// Bounds constants for date arithmetic.
const (
	MinYear            = 1
	MaxYear            = 9999
	MaxYearsDuration   = 10000
	MaxMonthsDuration  = 120000
	MaxDaysDuration    = 3650000
	MaxHoursDuration   = 87600000
	MaxMinutesDuration = 5256000000
)

// AddDuration returns the date after adding the given duration to ref.
// It handles fractional durations by cascading remainders to smaller units.
// For example, 1.5 months becomes 1 month + 2 weeks.
// Returns an error if the duration or resulting date is out of bounds.
func AddDuration(ref time.Time, duration kronos.Duration) (time.Time, error) {
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
	working := make(kronos.Duration)
	for k, v := range duration {
		working[k] = v
	}

	// Process decades (convert to years)
	if val, exists := working[kronos.TimeunitDecade]; exists {
		const yearsPerDecade = 10
		working[kronos.TimeunitYear] += val * yearsPerDecade
	}

	// Process years (cascade fractional part to months)
	if val, exists := working[kronos.TimeunitYear]; exists {
		if err := checkFloatToIntOverflow(val, kronos.TimeunitYear); err != nil {
			return time.Time{}, err
		}
		floor := int(val)
		date = date.AddDate(floor, 0, 0)
		if err := validateDate(date); err != nil {
			return time.Time{}, err
		}
		remainder := val - float64(floor)
		if remainder != 0 {
			working[kronos.TimeunitMonth] += remainder * 12
		}
	}

	// Process quarters (convert to months)
	if val, exists := working[kronos.TimeunitQuarter]; exists {
		if err := checkFloatToIntOverflow(val, kronos.TimeunitQuarter); err != nil {
			return time.Time{}, err
		}
		floor := int(val)
		date = date.AddDate(0, floor*3, 0)
		if err := validateDate(date); err != nil {
			return time.Time{}, err
		}
	}

	// Process months (cascade fractional part to weeks)
	if val, exists := working[kronos.TimeunitMonth]; exists {
		if err := checkFloatToIntOverflow(val, kronos.TimeunitMonth); err != nil {
			return time.Time{}, err
		}
		floor := int(val)
		date = date.AddDate(0, floor, 0)
		if err := validateDate(date); err != nil {
			return time.Time{}, err
		}
		remainder := val - float64(floor)
		if remainder != 0 {
			working[kronos.TimeunitWeek] += remainder * 4
		}
	}

	// Process weeks (cascade fractional part to days)
	if val, exists := working[kronos.TimeunitWeek]; exists {
		floor := int(val)
		date = date.AddDate(0, 0, floor*7)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional weeks to days
			// Round to nearest day for both positive and negative values
			days := remainder * 7
			const roundingOffset = 0.5
			if days > 0 {
				working[kronos.TimeunitDay] += float64(int(days + roundingOffset))
			} else if days < 0 {
				working[kronos.TimeunitDay] += float64(int(days - roundingOffset))
			}
		}
	}

	// Process days (cascade fractional part to hours)
	if val, exists := working[kronos.TimeunitDay]; exists {
		floor := int(val)
		date = date.AddDate(0, 0, floor)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional days to hours
			// Round to nearest hour for both positive and negative values
			hours := remainder * 24
			const roundingOffset = 0.5
			if hours > 0 {
				working[kronos.TimeunitHour] += float64(int(hours + roundingOffset))
			} else if hours < 0 {
				working[kronos.TimeunitHour] += float64(int(hours - roundingOffset))
			}
		}
	}

	// Process hours (cascade fractional part to minutes)
	if val, exists := working[kronos.TimeunitHour]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Hour)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional hours to minutes
			// For positive values: round to nearest minute
			// For negative values: preserve sign and round to nearest minute
			minutes := remainder * 60
			const roundingOffset = 0.5
			if minutes > 0 {
				working[kronos.TimeunitMinute] += float64(int(minutes + roundingOffset))
			} else if minutes < 0 {
				working[kronos.TimeunitMinute] += float64(int(minutes - roundingOffset))
			}
		}
	}

	// Process minutes (cascade fractional part to seconds)
	if val, exists := working[kronos.TimeunitMinute]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Minute)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional minutes to seconds
			// For positive values: round to nearest second
			// For negative values: preserve sign and round to nearest second
			seconds := remainder * 60
			const roundingOffset = 0.5
			if seconds > 0 {
				working[kronos.TimeunitSecond] += float64(int(seconds + roundingOffset))
			} else if seconds < 0 {
				working[kronos.TimeunitSecond] += float64(int(seconds - roundingOffset))
			}
		}
	}

	// Process seconds (cascade fractional part to milliseconds)
	if val, exists := working[kronos.TimeunitSecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Second)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional seconds to milliseconds
			// For positive values: round to nearest millisecond
			// For negative values: preserve sign and round to nearest millisecond
			milliseconds := remainder * 1000
			const roundingOffset = 0.5
			if milliseconds > 0 {
				working[kronos.TimeunitMillisecond] += float64(int(milliseconds + roundingOffset))
			} else if milliseconds < 0 {
				working[kronos.TimeunitMillisecond] += float64(int(milliseconds - roundingOffset))
			}
		}
	}

	// Process milliseconds (cascade fractional part to microseconds)
	if val, exists := working[kronos.TimeunitMillisecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Millisecond)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional milliseconds to microseconds
			microseconds := remainder * 1000
			const roundingOffset = 0.5
			if microseconds > 0 {
				working[kronos.TimeunitMicrosecond] += float64(int(microseconds + roundingOffset))
			} else if microseconds < 0 {
				working[kronos.TimeunitMicrosecond] += float64(int(microseconds - roundingOffset))
			}
		}
	}

	// Process microseconds (cascade fractional part to nanoseconds)
	if val, exists := working[kronos.TimeunitMicrosecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Microsecond)
		remainder := val - float64(floor)
		if remainder != 0 {
			// Convert fractional microseconds to nanoseconds
			nanoseconds := remainder * 1000
			const roundingOffset = 0.5
			if nanoseconds > 0 {
				working[kronos.TimeunitNanosecond] += float64(int(nanoseconds + roundingOffset))
			} else if nanoseconds < 0 {
				working[kronos.TimeunitNanosecond] += float64(int(nanoseconds - roundingOffset))
			}
		}
	}

	// Process nanoseconds
	if val, exists := working[kronos.TimeunitNanosecond]; exists {
		floor := int(val)
		date = date.Add(time.Duration(floor) * time.Nanosecond)
	}

	// Final validation of result date
	if err := validateDate(date); err != nil {
		return time.Time{}, err
	}

	return date, nil
}

// ReverseDuration returns the reversed duration (e.g., back into the past instead of future).
// All values in the duration are negated.
func ReverseDuration(duration kronos.Duration) kronos.Duration {
	reversed := make(kronos.Duration)
	for key, val := range duration {
		reversed[key] = -val
	}
	return reversed
}

// Helper functions

func validateDate(t time.Time) error {
	year := t.Year()
	if year < MinYear || year > MaxYear {
		return fmt.Errorf("date year %d is outside valid range [%d, %d]", year, MinYear, MaxYear)
	}
	return nil
}

func validateDuration(d kronos.Duration) error {
	for unit, value := range d {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("duration contains invalid value for %s: %f", unit, value)
		}

		absValue := math.Abs(value)
		switch unit {
		case kronos.TimeunitYear, kronos.TimeunitDecade:
			if absValue > MaxYearsDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, MaxYearsDuration)
			}
		case kronos.TimeunitMonth, kronos.TimeunitQuarter:
			if absValue > MaxMonthsDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, MaxMonthsDuration)
			}
		case kronos.TimeunitWeek, kronos.TimeunitDay:
			if absValue > MaxDaysDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, MaxDaysDuration)
			}
		case kronos.TimeunitHour:
			if absValue > MaxHoursDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, MaxHoursDuration)
			}
		case kronos.TimeunitMinute:
			if absValue > MaxMinutesDuration {
				return fmt.Errorf("duration %s value %f exceeds maximum %d", unit, value, MaxMinutesDuration)
			}
		}
	}
	return nil
}

func checkCascadingOverflow(d kronos.Duration) error {
	totalYears := d[kronos.TimeunitYear] + d[kronos.TimeunitDecade]*10
	totalMonths := totalYears*12 + d[kronos.TimeunitMonth] + d[kronos.TimeunitQuarter]*3
	totalDays := d[kronos.TimeunitDay] + d[kronos.TimeunitWeek]*7

	if math.Abs(totalYears) > MaxYearsDuration {
		return fmt.Errorf("cascading year duration %f exceeds maximum %d", totalYears, MaxYearsDuration)
	}
	if math.Abs(totalMonths) > MaxMonthsDuration {
		return fmt.Errorf("cascading month duration %f exceeds maximum %d", totalMonths, MaxMonthsDuration)
	}
	if math.Abs(totalDays) > MaxDaysDuration {
		return fmt.Errorf("cascading day duration %f exceeds maximum %d", totalDays, MaxDaysDuration)
	}

	return nil
}

func checkFloatToIntOverflow(value float64, unit kronos.Timeunit) error {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return fmt.Errorf("duration %s value %f would overflow int32 conversion", unit, value)
	}
	return nil
}
