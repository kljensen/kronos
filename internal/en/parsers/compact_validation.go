package parsers

import "time"

// isValidDate checks if year, month, day form a valid date
func isValidDate(year, month, day int) bool {
	if year < 1000 || year > 9999 {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}
	if day < 1 || day > 31 {
		return false
	}

	// Check if the date is actually valid (e.g., no Feb 30)
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.Year() == year && int(t.Month()) == month && t.Day() == day
}

// isValidMonthDay checks if month and day form a valid date
func isValidMonthDay(month, day int) bool {
	if month < 1 || month > 12 {
		return false
	}
	if day < 1 || day > 31 {
		return false
	}

	// Use a leap year to validate day ranges (2020 was a leap year)
	t := time.Date(2020, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return int(t.Month()) == month && t.Day() == day
}

// isValidTime checks if hour, minute, second form a valid time
func isValidTime(hour, minute, second int) bool {
	return hour >= 0 && hour <= 23 &&
		minute >= 0 && minute <= 59 &&
		second >= 0 && second <= 59
}

// convertTwoDigitYear converts a 2-digit year to a 4-digit year
// Rule: 00-69 → 2000-2069, 70-99 → 1970-1999
func convertTwoDigitYear(year int) int {
	if year >= 0 && year <= 69 {
		return 2000 + year
	}
	return 1900 + year
}
