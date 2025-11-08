package component

import "time"

// BreakdownNanoseconds splits nanoseconds into milliseconds, microseconds, and nanoseconds.
func BreakdownNanoseconds(totalNanos int) (millisecond, microsecond, nanosecond int) {
	millisecond = totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond = remainingNanos / 1000
	nanosecond = remainingNanos % 1000
	return
}

// HourToMeridiem converts an hour (0-23) to meridiem value (0=AM, 1=PM).
func HourToMeridiem(hour int) int {
	if hour < 12 {
		return 0 // AM
	}
	return 1 // PM
}

// GetValueOrDefault returns the value pointed to by ptr, or defaultVal if ptr is nil.
func GetValueOrDefault(ptr *int, defaultVal int) int {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// CreateDateInLocation creates a time.Time from component values in the specified location.
// Component values should be provided as pointers (nil means use default).
func CreateDateInLocation(
	year, month, day, hour, minute, second, millisecond, microsecond, nanosecond *int,
	location *time.Location,
) time.Time {
	const (
		defaultYear  = 2000
		defaultMonth = 1
		defaultDay   = 1
	)

	y := GetValueOrDefault(year, defaultYear)
	m := GetValueOrDefault(month, defaultMonth)
	d := GetValueOrDefault(day, defaultDay)
	h := GetValueOrDefault(hour, 0)
	min := GetValueOrDefault(minute, 0)
	sec := GetValueOrDefault(second, 0)
	ms := GetValueOrDefault(millisecond, 0)
	us := GetValueOrDefault(microsecond, 0)
	ns := GetValueOrDefault(nanosecond, 0)

	// Calculate total nanoseconds from milliseconds, microseconds, and nanoseconds
	totalNanos := (ms * 1000000) + (us * 1000) + ns

	return time.Date(y, time.Month(m), d, h, min, sec, totalNanos, location)
}
