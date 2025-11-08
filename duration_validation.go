package kronos

import (
	"fmt"
	"math"
	"time"
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
)

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
