package kronos

import "time"

// AssignSimilarDate assigns (force updates) the parsing components to the same day as the target.
// This sets year, month, and day as certain (known) values.
func AssignSimilarDate(components *ParsingComponents, date time.Time) {
	components.Assign(ComponentDay, date.Day())
	components.Assign(ComponentMonth, int(date.Month()))
	components.Assign(ComponentYear, date.Year())
}

// AssignSimilarTime assigns (force updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, and meridiem as certain (known) values.
func AssignSimilarTime(components *ParsingComponents, date time.Time) {
	components.Assign(ComponentHour, date.Hour())
	components.Assign(ComponentMinute, date.Minute())
	components.Assign(ComponentSecond, date.Second())
	components.Assign(ComponentMillisecond, date.Nanosecond()/NanosecondsPerMS)

	// Set meridiem based on hour
	if date.Hour() < HoursPerDay/2 {
		components.Assign(ComponentMeridiem, int(MeridiemAM))
	} else {
		components.Assign(ComponentMeridiem, int(MeridiemPM))
	}
}

// ImplySimilarDate implies (weakly updates) the parsing components to the same day as the target.
// This sets year, month, and day as implied values (only if not already certain).
func ImplySimilarDate(components *ParsingComponents, date time.Time) {
	components.Imply(ComponentDay, date.Day())
	components.Imply(ComponentMonth, int(date.Month()))
	components.Imply(ComponentYear, date.Year())
}

// ImplySimilarTime implies (weakly updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, and meridiem as implied values (only if not already certain).
func ImplySimilarTime(components *ParsingComponents, date time.Time) {
	components.Imply(ComponentHour, date.Hour())
	components.Imply(ComponentMinute, date.Minute())
	components.Imply(ComponentSecond, date.Second())
	components.Imply(ComponentMillisecond, date.Nanosecond()/NanosecondsPerMS)

	// Set meridiem based on hour
	if date.Hour() < HoursPerDay/2 {
		components.Imply(ComponentMeridiem, int(MeridiemAM))
	} else {
		components.Imply(ComponentMeridiem, int(MeridiemPM))
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
	if year2000s <= currentYear+YearLookAheadThreshold {
		return year2000s
	}

	// Otherwise, use 1900s
	const year1900Base = 1900
	return year1900Base + rawYear
}

// FindYearClosestToRef finds the year (past or future) that is closest to the reference date
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
func FindYearClosestToRef(refDate time.Time, day, month int) int {
	const defaultImpliedHour = 12 // Use noon for comparison

	refYear := refDate.Year()

	// Try the reference year and adjacent years
	candidates := []time.Time{
		time.Date(refYear-1, time.Month(month), day, defaultImpliedHour, 0, 0, 0, time.Local),
		time.Date(refYear, time.Month(month), day, defaultImpliedHour, 0, 0, 0, time.Local),
		time.Date(refYear+1, time.Month(month), day, defaultImpliedHour, 0, 0, 0, time.Local),
	}

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
