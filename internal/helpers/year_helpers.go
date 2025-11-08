package helpers

import (
	"time"

	"github.com/kljensen/kronos"
)

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
func FindMostLikelyADYear(rawYear int) int {
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
func FindYearClosestToRefWithPreference(refDate time.Time, day, month int, preference kronos.DatePreference) int {
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
	case kronos.PreferPast:
		// Choose the most recent date that is in the past
		for i := len(candidates) - 1; i >= 0; i-- {
			if candidates[i].Before(refDate) || candidates[i].Equal(refDate) {
				return candidates[i].Year()
			}
		}
		// If all candidates are in the future, return the earliest one
		return candidates[0].Year()

	case kronos.PreferFuture:
		// Choose the nearest date that is in the future
		for i := range len(candidates) {
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

// isLeapYear checks if a year is a leap year according to the Gregorian calendar rules.
func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

// findPreviousLeapYear finds the most recent leap year before or at baseYear.
func findPreviousLeapYear(baseYear int) int {
	const lowerBound = 1900
	for year := baseYear; year >= lowerBound; year-- {
		if isLeapYear(year) {
			return year
		}
	}
	return baseYear
}

// findNextLeapYear finds the next leap year after or at baseYear.
func findNextLeapYear(baseYear int) int {
	const upperBound = 9999
	for year := baseYear; year <= upperBound; year++ {
		if isLeapYear(year) {
			return year
		}
	}
	return baseYear
}

// findNearestLeapYear finds the nearest leap year based on date preference.
func findNearestLeapYear(baseYear int, preference kronos.DatePreference) int {
	// If the base year is already a leap year, use it for all preferences
	if isLeapYear(baseYear) {
		return baseYear
	}

	switch preference {
	case kronos.PreferPast:
		return findPreviousLeapYear(baseYear)

	case kronos.PreferFuture:
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
