package kronos

import "time"

// GetNextWeekday returns the date of the next occurrence of the target weekday
// after the reference date (not including the reference date itself).
//
// For example, if refDate is Monday and targetWeekday is Monday, it returns
// the following Monday (7 days later).
func getNextWeekday(refDate time.Time, targetWeekday Weekday) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	// Calculate days forward to next occurrence
	daysForward := target - refWeekday
	if daysForward <= 0 {
		daysForward += 7
	}

	return refDate.AddDate(0, 0, daysForward)
}

// GetLastWeekday returns the date of the last occurrence of the target weekday
// before the reference date (not including the reference date itself).
//
// For example, if refDate is Monday and targetWeekday is Monday, it returns
// the previous Monday (7 days earlier).
func getLastWeekday(refDate time.Time, targetWeekday Weekday) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	// Calculate days backward to last occurrence
	daysBackward := target - refWeekday
	if daysBackward >= 0 {
		daysBackward -= 7
	}

	return refDate.AddDate(0, 0, daysBackward)
}

// GetThisWeekday returns the date of the target weekday in "this" week.
//
// The forward parameter controls the direction:
// - If forward=true: looks forward from refDate (inclusive)
// - If forward=false: looks backward from refDate (inclusive)
//
// Examples with refDate = Wednesday:
// - getThisWeekday(Wednesday, Monday, forward=true) -> next Monday (5 days forward)
// - getThisWeekday(Wednesday, Monday, forward=false) -> previous Monday (2 days back)
// - getThisWeekday(Wednesday, Wednesday, forward=true) -> same Wednesday (0 days)
// - getThisWeekday(Wednesday, Wednesday, forward=false) -> same Wednesday (0 days)
func getThisWeekday(refDate time.Time, targetWeekday Weekday, forward bool) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	if refWeekday == target {
		// Same weekday - return the reference date
		return refDate
	}

	if forward {
		// Look forward
		daysForward := target - refWeekday
		if daysForward < 0 {
			daysForward += 7
		}
		return refDate.AddDate(0, 0, daysForward)
	} else {
		// Look backward
		daysBackward := target - refWeekday
		if daysBackward > 0 {
			daysBackward -= 7
		}
		return refDate.AddDate(0, 0, daysBackward)
	}
}

// GetDaysToWeekday returns the number of days from refDate to the target weekday
// based on the modifier.
//
// Modifiers:
// - "this": Returns the target weekday in the current week (looking forward)
// - "next": Returns the target weekday in the next occurrence
// - "last": Returns the target weekday in the last occurrence (negative value)
// - nil/empty: Returns the closest weekday (forward or backward)
func getDaysToWeekday(refDate time.Time, targetWeekday Weekday, modifier *string) int {
	refWeekday := Weekday(refDate.Weekday())

	if modifier == nil {
		return getDaysToWeekdayClosest(refDate, targetWeekday)
	}

	switch *modifier {
	case "this":
		return getDaysForwardToWeekday(refDate, targetWeekday)
	case "last":
		return getBackwardDaysToWeekday(refDate, targetWeekday)
	case "next":
		// Special handling for "next" based on the current weekday
		if refWeekday == WeekdaySunday {
			if targetWeekday == WeekdaySunday {
				return 7
			}
			return int(targetWeekday)
		}

		if refWeekday == WeekdaySaturday {
			if targetWeekday == WeekdaySaturday {
				return 7
			}
			if targetWeekday == WeekdaySunday {
				return 8
			}
			return 1 + int(targetWeekday)
		}

		// For weekdays (Mon-Fri)
		if targetWeekday < refWeekday && targetWeekday != WeekdaySunday {
			return getDaysForwardToWeekday(refDate, targetWeekday)
		}
		return getDaysForwardToWeekday(refDate, targetWeekday) + 7
	}

	return getDaysToWeekdayClosest(refDate, targetWeekday)
}

// getDaysForwardToWeekday returns the number of days forward to reach the target weekday.
// Returns 0 if already on the target weekday.
func getDaysForwardToWeekday(refDate time.Time, targetWeekday Weekday) int {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	forwardCount := target - refWeekday
	if forwardCount < 0 {
		forwardCount += 7
	}

	return forwardCount
}

// getBackwardDaysToWeekday returns the number of days backward to reach the target weekday.
// Returns a negative number (or 0 if already on the target weekday).
func getBackwardDaysToWeekday(refDate time.Time, targetWeekday Weekday) int {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	backwardCount := target - refWeekday
	if backwardCount >= 0 {
		backwardCount -= 7
	}

	return backwardCount
}

// getDaysToWeekdayClosest returns the number of days to the closest occurrence
// of the target weekday (either forward or backward).
func getDaysToWeekdayClosest(refDate time.Time, targetWeekday Weekday) int {
	backward := getBackwardDaysToWeekday(refDate, targetWeekday)
	forward := getDaysForwardToWeekday(refDate, targetWeekday)

	// Choose the direction with fewer days
	if forward < -backward {
		return forward
	}
	return backward
}
