package kronos

import (
	"regexp"
	"strings"
	"time"
)

// Test helper functions are only used by tests to verify parsing behavior.
// These cannot import internal/helpers due to import cycles (internal/helpers imports kronos).
// Therefore we keep implementations here that are test-only.
//
// NOTE: The real implementations are in internal/helpers/*.go and are used by the actual parsers.
// This duplication is necessary due to Go's import cycle restrictions. Tests in package kronos
// cannot import internal/helpers because internal/helpers imports kronos.

// ============================================================================
// Casual reference helpers (test-only versions)
// ============================================================================

func now(reference *referenceWithTimezone) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	component.AssignSimilarDate(targetDate)
	component.AssignSimilarTime(targetDate)
	component.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	component.AddTag("casualReference/now")
	component.SetPeriod(PeriodTime)

	return component
}

func today(reference *referenceWithTimezone) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	component.AssignSimilarDate(targetDate)
	component.ImplySimilarTime(targetDate)
	component.Delete(ComponentMeridiem)
	component.AddTag("casualReference/today")
	component.SetPeriod(PeriodDay)

	return component
}

func yesterday(reference *referenceWithTimezone) *parsingComponents {
	component := theDayBefore(reference, 1)
	component.AddTag("casualReference/yesterday")
	component.SetPeriod(PeriodDay)
	return component
}

func tomorrow(reference *referenceWithTimezone) *parsingComponents {
	component := theDayAfter(reference, 1)
	component.AddTag("casualReference/tomorrow")
	component.SetPeriod(PeriodDay)
	return component
}

func theDayBefore(reference *referenceWithTimezone, nDays int) *parsingComponents {
	return theDayAfter(reference, -nDays)
}

func theDayAfter(reference *referenceWithTimezone, nDays int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	newDate := targetDate.AddDate(0, 0, nDays)

	component.AssignSimilarDate(newDate)
	component.ImplySimilarTime(newDate)
	component.Delete(ComponentMeridiem)
	component.SetPeriod(PeriodDay)

	return component
}

func tonight(reference *referenceWithTimezone) *parsingComponents {
	return tonightWithHour(reference, 22)
}

func tonightWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	component.AssignSimilarDate(targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, 1) // PM
	component.AddTag("casualReference/tonight")
	component.SetPeriod(PeriodDay)

	return component
}

func lastNight(reference *referenceWithTimezone) *parsingComponents {
	return lastNightWithHour(reference, 0)
}

func lastNightWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	// If it's very early morning (before 6 AM), "last night" refers to yesterday
	if targetDate.Hour() < 6 {
		targetDate = targetDate.AddDate(0, 0, -1)
	}

	component.AssignSimilarDate(targetDate)
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/lastNight")
	component.SetPeriod(PeriodDay)

	return component
}

func evening(reference *referenceWithTimezone) *parsingComponents {
	return eveningWithHour(reference, 20)
}

func eveningWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, 1) // PM
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/evening")
	component.SetPeriod(PeriodTime)

	return component
}

func yesterdayEvening(reference *referenceWithTimezone) *parsingComponents {
	return yesterdayEveningWithHour(reference, 20)
}

func yesterdayEveningWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	targetDate = targetDate.AddDate(0, 0, -1)

	component.AssignSimilarDate(targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, 1) // PM
	component.AddTag("casualReference/yesterday")
	component.AddTag("casualReference/evening")
	component.SetPeriod(PeriodTime)

	return component
}

func midnight(reference *referenceWithTimezone) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	// Unless it's very early morning (0-2 AM), assume midnight refers to the coming midnight
	if targetDate.Hour() > 2 {
		duration := Duration{TimeunitDay: 1}
		newDate, err := addDuration(targetDate, duration)
		if err != nil {
			return nil
		}
		component.ImplySimilarDate(newDate)
	}

	component.Assign(ComponentHour, 0)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/midnight")
	component.SetPeriod(PeriodTime)

	return component
}

func morning(reference *referenceWithTimezone) *parsingComponents {
	return morningWithHour(reference, 6)
}

func morningWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, 0) // AM
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/morning")
	component.SetPeriod(PeriodTime)

	return component
}

func afternoon(reference *referenceWithTimezone) *parsingComponents {
	return afternoonWithHour(reference, 15)
}

func afternoonWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, 1) // PM
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/afternoon")
	component.SetPeriod(PeriodTime)

	return component
}

func noon(reference *referenceWithTimezone) *parsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, 1) // PM
	component.Assign(ComponentHour, 12)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/noon")
	component.SetPeriod(PeriodTime)

	return component
}

// ============================================================================
// Weekday calculation helpers (test-only versions)
// ============================================================================

func getNextWeekday(refDate time.Time, targetWeekday time.Weekday) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	daysForward := target - refWeekday
	if daysForward <= 0 {
		daysForward += 7
	}

	return refDate.AddDate(0, 0, daysForward)
}

func getLastWeekday(refDate time.Time, targetWeekday time.Weekday) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	daysBackward := target - refWeekday
	if daysBackward >= 0 {
		daysBackward -= 7
	}

	return refDate.AddDate(0, 0, daysBackward)
}

func getThisWeekday(refDate time.Time, targetWeekday time.Weekday, forward bool) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	if refWeekday == target {
		return refDate
	}

	if forward {
		daysForward := target - refWeekday
		if daysForward < 0 {
			daysForward += 7
		}
		return refDate.AddDate(0, 0, daysForward)
	} else {
		daysBackward := target - refWeekday
		if daysBackward > 0 {
			daysBackward -= 7
		}
		return refDate.AddDate(0, 0, daysBackward)
	}
}

func getDaysToWeekday(refDate time.Time, targetWeekday time.Weekday, modifier *string) int {
	refWeekday := time.Weekday(refDate.Weekday())

	if modifier == nil {
		return getDaysToWeekdayClosest(refDate, targetWeekday)
	}

	switch *modifier {
	case "this":
		return getDaysForwardToWeekday(refDate, targetWeekday)
	case "last":
		return getBackwardDaysToWeekday(refDate, targetWeekday)
	case "next":
		if refWeekday == time.Sunday {
			if targetWeekday == time.Sunday {
				return 7
			}
			return int(targetWeekday)
		}

		if refWeekday == time.Saturday {
			if targetWeekday == time.Saturday {
				return 7
			}
			if targetWeekday == time.Sunday {
				return 8
			}
			return 1 + int(targetWeekday)
		}

		if targetWeekday < refWeekday && targetWeekday != time.Sunday {
			return getDaysForwardToWeekday(refDate, targetWeekday)
		}
		return getDaysForwardToWeekday(refDate, targetWeekday) + 7
	}

	return getDaysToWeekdayClosest(refDate, targetWeekday)
}

func getDaysForwardToWeekday(refDate time.Time, targetWeekday time.Weekday) int {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	forwardCount := target - refWeekday
	if forwardCount < 0 {
		forwardCount += 7
	}

	return forwardCount
}

func getBackwardDaysToWeekday(refDate time.Time, targetWeekday time.Weekday) int {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	backwardCount := target - refWeekday
	if backwardCount >= 0 {
		backwardCount -= 7
	}

	return backwardCount
}

func getDaysToWeekdayClosest(refDate time.Time, targetWeekday time.Weekday) int {
	backward := getBackwardDaysToWeekday(refDate, targetWeekday)
	forward := getDaysForwardToWeekday(refDate, targetWeekday)

	if forward < -backward {
		return forward
	}
	return backward
}

// ============================================================================
// Year calculation helpers (test-only versions)
// ============================================================================

func findMostLikelyADYear(rawYear int) int {
	const twoDigitThreshold = 100

	if rawYear >= twoDigitThreshold {
		return rawYear
	}

	currentYear := time.Now().Year()
	const year2000Base = 2000
	year2000s := year2000Base + rawYear

	if year2000s <= currentYear+20 {
		return year2000s
	}

	const year1900Base = 1900
	return year1900Base + rawYear
}

func findYearClosestToRef(refDate time.Time, day, month int) int {
	return findYearClosestToRefWithPreference(refDate, day, month, PreferCurrentPeriod)
}

func findYearClosestToRefWithPreference(refDate time.Time, day, month int, preference DatePreference) int {
	const defaultImpliedHour = 12

	refYear := refDate.Year()
	location := refDate.Location()

	if month == 2 && day == 29 {
		return findNearestLeapYear(refYear, preference)
	}

	candidates := []time.Time{
		time.Date(refYear-1, time.Month(month), day, defaultImpliedHour, 0, 0, 0, location),
		time.Date(refYear, time.Month(month), day, defaultImpliedHour, 0, 0, 0, location),
		time.Date(refYear+1, time.Month(month), day, defaultImpliedHour, 0, 0, 0, location),
	}

	switch preference {
	case PreferPast:
		for i := len(candidates) - 1; i >= 0; i-- {
			if candidates[i].Before(refDate) || candidates[i].Equal(refDate) {
				return candidates[i].Year()
			}
		}
		return candidates[0].Year()

	case PreferFuture:
		for i := 0; i < len(candidates); i++ {
			if candidates[i].After(refDate) || candidates[i].Equal(refDate) {
				return candidates[i].Year()
			}
		}
		return candidates[len(candidates)-1].Year()

	default: // PreferCurrentPeriod
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

func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

func findPreviousLeapYear(baseYear int) int {
	const lowerBound = 1900
	for year := baseYear; year >= lowerBound; year-- {
		if isLeapYear(year) {
			return year
		}
	}
	return baseYear
}

func findNextLeapYear(baseYear int) int {
	const upperBound = 9999
	for year := baseYear; year <= upperBound; year++ {
		if isLeapYear(year) {
			return year
		}
	}
	return baseYear
}

func findNearestLeapYear(baseYear int, preference DatePreference) int {
	if isLeapYear(baseYear) {
		return baseYear
	}

	switch preference {
	case PreferPast:
		return findPreviousLeapYear(baseYear)

	case PreferFuture:
		return findNextLeapYear(baseYear)

	default: // PreferCurrentPeriod
		next := findNextLeapYear(baseYear)
		prev := findPreviousLeapYear(baseYear)

		distToNext := next - baseYear
		distToPrev := baseYear - prev

		if distToNext <= distToPrev {
			return next
		}
		return prev
	}
}

// ============================================================================
// Approximation helpers (test-only)
// ============================================================================

var approximationWords = []string{
	"about",
	"around",
	"roughly",
	"approximately",
	"approx",
	"circa",
}

func stripApproximationWords(input string) (cleaned string, isApproximate bool) {
	cleaned = input
	isApproximate = false

	tildePattern := regexp.MustCompile(`(?i)~\s*`)
	if tildePattern.MatchString(cleaned) {
		isApproximate = true
		cleaned = tildePattern.ReplaceAllString(cleaned, "")
	}

	for _, word := range approximationWords {
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\s+`)
		if pattern.MatchString(cleaned) {
			isApproximate = true
			cleaned = pattern.ReplaceAllString(cleaned, "")
		}
	}

	return strings.TrimSpace(cleaned), isApproximate
}
