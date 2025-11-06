package kronos

import "time"

// Internal helper functions exported with X prefix for use by internal/ packages only.
// These should NOT be used by external code.

// XSafeSlice is an internal helper for internal/ packages.
func XSafeSlice(text string, start, end int) (string, bool) {
	return safeSlice(text, start, end)
}

// XAsParsingComponents is an internal helper for internal/ packages.
func XAsParsingComponents(pc ParsedComponents) (*ParsingComponents, bool) {
	return asParsingComponents(pc)
}

// XImplySimilarDate is an internal helper for internal/ packages.
func XImplySimilarDate(components *ParsingComponents, date time.Time) {
	implySimilarDate(components, date)
}

// XAssignSimilarDate is an internal helper for internal/ packages.
func XAssignSimilarDate(components *ParsingComponents, date time.Time) {
	assignSimilarDate(components, date)
}

// XAssignSimilarTime is an internal helper for internal/ packages.
func XAssignSimilarTime(components *ParsingComponents, date time.Time) {
	assignSimilarTime(components, date)
}

// XNewParsingContext is an internal helper for internal/ packages.
func XNewParsingContext(text string, refDate interface{}, option *ParsingOption) *ParsingContext {
	return newParsingContext(text, refDate, option)
}

// XNewParsingComponents is an internal helper for internal/ packages.
func XNewParsingComponents(reference *ReferenceWithTimezone, knownComponents map[Component]int) *ParsingComponents {
	return newParsingComponents(reference, knownComponents)
}

// XFindMostLikelyADYear is an internal helper for internal/ packages.
func XFindMostLikelyADYear(rawYear int) int {
	return findMostLikelyADYear(rawYear)
}

// XFindYearClosestToRefWithPreference is an internal helper for internal/ packages.
func XFindYearClosestToRefWithPreference(refDate time.Time, day, month int, preference DatePreference) int {
	return findYearClosestToRefWithPreference(refDate, day, month, preference)
}

// XToday is an internal helper for internal/ packages.
func XToday(reference *ReferenceWithTimezone) *ParsingComponents {
	return today(reference)
}

// XTomorrow is an internal helper for internal/ packages.
func XTomorrow(reference *ReferenceWithTimezone) *ParsingComponents {
	return tomorrow(reference)
}

// XYesterday is an internal helper for internal/ packages.
func XYesterday(reference *ReferenceWithTimezone) *ParsingComponents {
	return yesterday(reference)
}

// XNow is an internal helper for internal/ packages.
func XNow(reference *ReferenceWithTimezone) *ParsingComponents {
	return now(reference)
}

// XTheDayAfter is an internal helper for internal/ packages.
func XTheDayAfter(reference *ReferenceWithTimezone, nDays int) *ParsingComponents {
	return theDayAfter(reference, nDays)
}

// XTheDayBefore is an internal helper for internal/ packages.
func XTheDayBefore(reference *ReferenceWithTimezone, nDays int) *ParsingComponents {
	return theDayBefore(reference, nDays)
}

// XMidnight is an internal helper for internal/ packages.
func XMidnight(reference *ReferenceWithTimezone) *ParsingComponents {
	return midnight(reference)
}

// XNoon is an internal helper for internal/ packages.
func XNoon(reference *ReferenceWithTimezone) *ParsingComponents {
	return noon(reference)
}

// XMorning is an internal helper for internal/ packages.
func XMorning(reference *ReferenceWithTimezone) *ParsingComponents {
	return morning(reference)
}

// XAfternoon is an internal helper for internal/ packages.
func XAfternoon(reference *ReferenceWithTimezone) *ParsingComponents {
	return afternoon(reference)
}

// XEvening is an internal helper for internal/ packages.
func XEvening(reference *ReferenceWithTimezone) *ParsingComponents {
	return evening(reference)
}

// XGetDaysToWeekday is an internal helper for internal/ packages.
func XGetDaysToWeekday(refDate time.Time, targetWeekday Weekday, modifier *string) int {
	return getDaysToWeekday(refDate, targetWeekday, modifier)
}

// XGetLastWeekday is an internal helper for internal/ packages.
func XGetLastWeekday(refDate time.Time, targetWeekday Weekday) time.Time {
	return getLastWeekday(refDate, targetWeekday)
}

// XGetNextWeekday is an internal helper for internal/ packages.
func XGetNextWeekday(refDate time.Time, targetWeekday Weekday) time.Time {
	return getNextWeekday(refDate, targetWeekday)
}

// XGetThisWeekday is an internal helper for internal/ packages.
func XGetThisWeekday(refDate time.Time, targetWeekday Weekday, forward bool) time.Time {
	return getThisWeekday(refDate, targetWeekday, forward)
}

// XCreateRelativeFromReference is an internal helper for internal/ packages.
func XCreateRelativeFromReference(reference *ReferenceWithTimezone, duration Duration) *ParsingComponents {
	return createRelativeFromReference(reference, duration)
}

// XAddDuration is an internal helper for internal/ packages.
func XAddDuration(ref time.Time, duration Duration) (time.Time, error) {
	return addDuration(ref, duration)
}

// XReverseDuration is an internal helper for internal/ packages.
func XReverseDuration(duration Duration) Duration {
	return reverseDuration(duration)
}

// XMergeDateTimeComponent is an internal helper for internal/ packages.
func XMergeDateTimeComponent(dateComp, timeComp *ParsingComponents) *ParsingComponents {
	return mergeDateTimeComponent(dateComp, timeComp)
}

// XNewParsingResult is an internal helper for internal/ packages.
func XNewParsingResult(reference *ReferenceWithTimezone, index int, text string, start, end *ParsingComponents) *ParsingResult {
	return newParsingResult(reference, index, text, start, end)
}

// XToTimezoneOffset is an internal helper for internal/ packages.
func XToTimezoneOffset(tz interface{}, instant time.Time, overrides TimezoneAbbrMap) *int {
	return toTimezoneOffset(tz, instant, overrides)
}

// XGetLastWeekdayOfMonth is an internal helper for internal/ packages.
func XGetLastWeekdayOfMonth(year int, month Month, weekday Weekday, hour int) time.Time {
	return getLastWeekdayOfMonth(year, month, weekday, hour)
}

// XGetNthWeekdayOfMonth is an internal helper for internal/ packages.
func XGetNthWeekdayOfMonth(year int, month Month, weekday Weekday, n int, hour int) time.Time {
	return getNthWeekdayOfMonth(year, month, weekday, n, hour)
}

// XStripApproximationWords is an internal helper for internal/ packages.
func XStripApproximationWords(input string) (string, bool) {
	return stripApproximationWords(input)
}

// XMergeDateTimeResult is an internal helper for internal/ packages.
func XMergeDateTimeResult(dateResult, timeResult *ParsingResult) *ParsingResult {
	return mergeDateTimeResult(dateResult, timeResult)
}

// XAfternoonWithHour is an internal helper for internal/ packages.
func XAfternoonWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	return afternoonWithHour(reference, implyHour)
}

// XEveningWithHour is an internal helper for internal/ packages.
func XEveningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	return eveningWithHour(reference, implyHour)
}

// XMorningWithHour is an internal helper for internal/ packages.
func XMorningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	return morningWithHour(reference, implyHour)
}

// XTonightWithHour is an internal helper for internal/ packages.
func XTonightWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	return tonightWithHour(reference, implyHour)
}
