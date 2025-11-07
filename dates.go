package kronos

import "time"

// assignSimilarDate assigns (force updates) the parsing components to the same day as the target.
// This sets year, month, and day as certain (known) values.
func assignSimilarDate(components *parsingComponents, date time.Time) {
	components.Assign(ComponentDay, date.Day())
	components.Assign(ComponentMonth, int(date.Month()))
	components.Assign(ComponentYear, date.Year())
}

// AssignSimilarTime assigns (force updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as certain (known) values.
func assignSimilarTime(components *parsingComponents, date time.Time) {
	components.Assign(ComponentHour, date.Hour())
	components.Assign(ComponentMinute, date.Minute())
	components.Assign(ComponentSecond, date.Second())

	// Break down nanoseconds into milliseconds, microseconds, and nanoseconds
	totalNanos := date.Nanosecond()
	millisecond := totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond := remainingNanos / 1000
	nanosecond := remainingNanos % 1000

	components.Assign(ComponentMillisecond, millisecond)
	if microsecond > 0 {
		components.Assign(ComponentMicrosecond, microsecond)
	}
	if nanosecond > 0 {
		components.Assign(ComponentNanosecond, nanosecond)
	}

	// Set meridiem based on hour
	if date.Hour() < 12 {
		components.Assign(ComponentMeridiem, 0) // AM
	} else {
		components.Assign(ComponentMeridiem, 1) // PM
	}
}

// ImplySimilarDate implies (weakly updates) the parsing components to the same day as the target.
// This sets year, month, and day as implied values (only if not already certain).
func implySimilarDate(components *parsingComponents, date time.Time) {
	components.Imply(ComponentDay, date.Day())
	components.Imply(ComponentMonth, int(date.Month()))
	components.Imply(ComponentYear, date.Year())
}

// ImplySimilarTime implies (weakly updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as implied values (only if not already certain).
func implySimilarTime(components *parsingComponents, date time.Time) {
	components.Imply(ComponentHour, date.Hour())
	components.Imply(ComponentMinute, date.Minute())
	components.Imply(ComponentSecond, date.Second())

	// Break down nanoseconds into milliseconds, microseconds, and nanoseconds
	totalNanos := date.Nanosecond()
	millisecond := totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond := remainingNanos / 1000
	nanosecond := remainingNanos % 1000

	components.Imply(ComponentMillisecond, millisecond)
	if microsecond > 0 {
		components.Imply(ComponentMicrosecond, microsecond)
	}
	if nanosecond > 0 {
		components.Imply(ComponentNanosecond, nanosecond)
	}

	// Set meridiem based on hour
	if date.Hour() < 12 {
		components.Imply(ComponentMeridiem, 0) // AM
	} else {
		components.Imply(ComponentMeridiem, 1) // PM
	}
}

// FindMostLikelyADYear converts a 2-digit year to a 4-digit year.
// Years 0-99 are mapped to 1900-2099 range.
//
// Logic:
// - 0-99 maps to 2000-2099 if the year would be <= current year + YearLookAheadThreshold
// - Otherwise maps to 1900-1999
//
// Examples (assuming current year is 2020):
//   - 20 -> 2020
//   - 40 -> 2040 (within 20 years of current)
//   - 50 -> 1950 (would be 2050, which is > 2040, so use 1900s)
//   - 99 -> 1999
func findMostLikelyADYear(rawYear int) int {
	const twoDigitThreshold = 100

	// If it's already a 4-digit year, return as-is
	if rawYear >= twoDigitThreshold {
		return rawYear
	}

	// Get current year for comparison
	currentYear := time.Now().Year()

	// Calculate the 2000s version
	const year2000Base = 2000
	year2000s := year2000Base + rawYear

	// If the year in 2000s would be within threshold of current year, use it
	if year2000s <= currentYear+20 {
		return year2000s
	}

	// Otherwise, use 1900s
	const year1900Base = 1900
	return year1900Base + rawYear
}

// findYearClosestToRef finds the year (past or future) that is closest to the reference date
// for a given day and month.
//
// This is useful when parsing dates without a year (e.g., "March 15").
// The function finds which year makes the date closest to the reference.
//
// Examples:
// - Reference: 2020-01-15, Day: 20, Month: 3 (March 20)
//   - Could be 2019-03-20 (about 10 months ago)
//   - Could be 2020-03-20 (about 2 months ahead)
//   - Could be 2021-03-20 (about 14 months ahead)
//   - Returns 2020 (closest match)
func findYearClosestToRef(refDate time.Time, day, month int) int {
	return findYearClosestToRefWithPreference(refDate, day, month, PreferCurrentPeriod)
}

// FindYearClosestToRefWithPreference finds the year based on date preference settings.
// This allows control over whether ambiguous dates should be resolved to past, future, or current period.
//
// Special handling for February 29:
// When month is 2 and day is 29, this function ensures we select a leap year.
// It adjusts the candidate years to be valid leap years according to the preference.
//
// Parameters:
//   - refDate: The reference date/time for comparison
//   - day: The day of month (1-31)
//   - month: The month (1-12)
//   - preference: How to resolve ambiguous dates (PreferPast, PreferFuture, PreferCurrentPeriod)
//
// Examples with reference date Feb 15, 2015 15:30:
//   - PreferCurrentPeriod: "March 15" → March 15, 2015 (current year)
//   - PreferPast: "March 15" → March 15, 2014 (last year, since March 15, 2015 is in future)
//   - PreferFuture: "March 15" → March 15, 2015 (this year, since it's in future)
//
// Examples with February 29 and reference date March 1, 2023:
//   - PreferPast: "February 29" → February 29, 2020 (previous leap year)
//   - PreferFuture: "February 29" → February 29, 2024 (next leap year)
//   - PreferCurrentPeriod: "February 29" → February 29, 2024 (nearest leap year)
func findYearClosestToRefWithPreference(refDate time.Time, day, month int, preference DatePreference) int {
	const defaultImpliedHour = 12 // Use noon for comparison

	refYear := refDate.Year()
	location := refDate.Location()

	// Special case: February 29 requires leap year selection
	if month == 2 && day == 29 {
		return findNearestLeapYear(refYear, preference)
	}

	// Create candidate dates for the reference year and adjacent years
	candidates := []time.Time{
		time.Date(refYear-1, time.Month(month), day, defaultImpliedHour, 0, 0, 0, location),
		time.Date(refYear, time.Month(month), day, defaultImpliedHour, 0, 0, 0, location),
		time.Date(refYear+1, time.Month(month), day, defaultImpliedHour, 0, 0, 0, location),
	}

	// Handle different preference modes
	switch preference {
	case PreferPast:
		// Choose the most recent date that is in the past
		for i := len(candidates) - 1; i >= 0; i-- {
			if candidates[i].Before(refDate) || candidates[i].Equal(refDate) {
				return candidates[i].Year()
			}
		}
		// If all candidates are in the future, return the earliest one
		return candidates[0].Year()

	case PreferFuture:
		// Choose the nearest date that is in the future
		for i := 0; i < len(candidates); i++ {
			if candidates[i].After(refDate) || candidates[i].Equal(refDate) {
				return candidates[i].Year()
			}
		}
		// If all candidates are in the past, return the latest one
		return candidates[len(candidates)-1].Year()

	default: // PreferCurrentPeriod
		// Find the candidate with the smallest absolute difference from refDate
		var minDiff int64
		closestYear := refYear

		for i, candidate := range candidates {
			diff := candidate.Unix() - refDate.Unix()
			if diff < 0 {
				diff = -diff
			}

			if i == 0 || diff < minDiff {
				minDiff = diff
				closestYear = candidate.Year()
			}
		}

		return closestYear
	}
}

// IsLeapYear checks if a year is a leap year according to the Gregorian calendar rules.
//
// Leap year rules:
// 1. Divisible by 4, AND
// 2. NOT divisible by 100, OR
// 3. Divisible by 400
//
// Examples:
//   - 2000: Leap year (divisible by 400)
//   - 1900: NOT leap year (divisible by 100 but not 400)
//   - 2004: Leap year (divisible by 4, not by 100)
//   - 2001: NOT leap year (not divisible by 4)
func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

// FindPreviousLeapYear finds the most recent leap year before or at baseYear.
// It searches backward from baseYear until it finds a leap year.
// The search is bounded at year 1900 to avoid excessive searching.
//
// Examples:
//   - findPreviousLeapYear(2023) → 2020
//   - findPreviousLeapYear(2024) → 2024
//   - findPreviousLeapYear(1901) → 1900 (bounded)
func findPreviousLeapYear(baseYear int) int {
	const lowerBound = 1900
	for year := baseYear; year >= lowerBound; year-- {
		if isLeapYear(year) {
			return year
		}
	}
	// Fallback: if no leap year found, return base year
	return baseYear
}

// FindNextLeapYear finds the next leap year after or at baseYear.
// It searches forward from baseYear until it finds a leap year.
// The search is bounded at year 9999 to avoid excessive searching.
//
// Examples:
//   - findNextLeapYear(2023) → 2024
//   - findNextLeapYear(2024) → 2024
//   - findNextLeapYear(9997) → 9998 (bounded)
func findNextLeapYear(baseYear int) int {
	const upperBound = 9999
	for year := baseYear; year <= upperBound; year++ {
		if isLeapYear(year) {
			return year
		}
	}
	// Fallback: if no leap year found, return base year
	return baseYear
}

// FindNearestLeapYear finds the nearest leap year based on date preference.
// This is used when parsing February 29 without a year, ensuring we select
// a valid leap year according to the preference setting.
//
// Parameters:
//   - baseYear: The starting year (typically the reference year)
//   - preference: How to resolve ambiguous dates (PreferPast, PreferFuture, PreferCurrentPeriod)
//
// Examples with baseYear 2023 (not a leap year):
//   - PreferPast → 2020 (previous leap year)
//   - PreferFuture → 2024 (next leap year)
//   - PreferCurrentPeriod → 2024 (prefer future when current isn't leap)
func findNearestLeapYear(baseYear int, preference DatePreference) int {
	// If the base year is already a leap year, use it for all preferences
	if isLeapYear(baseYear) {
		return baseYear
	}

	switch preference {
	case PreferPast:
		return findPreviousLeapYear(baseYear)

	case PreferFuture:
		return findNextLeapYear(baseYear)

	default: // PreferCurrentPeriod
		// When current year is not a leap year, prefer the nearest one
		// Break ties by preferring the future
		next := findNextLeapYear(baseYear)
		prev := findPreviousLeapYear(baseYear)

		// Calculate distances
		distToNext := next - baseYear
		distToPrev := baseYear - prev

		// Prefer future if equal distance
		if distToNext <= distToPrev {
			return next
		}
		return prev
	}
}
