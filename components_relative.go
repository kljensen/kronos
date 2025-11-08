package kronos

// ============================================================================
// Period determination
// ============================================================================

// DeterminePeriodFromDuration determines the granularity/period based on a duration.
// The period represents the finest time unit present in the duration.
// This follows the pattern from Python's dateparser.
func determinePeriodFromDuration(duration Duration) Period {
	if duration == nil {
		return PeriodDay // Default
	}

	// Check from finest to coarsest granularity
	// Time components (hour, minute, second) indicate time-level precision
	for _, timeunit := range []Timeunit{TimeunitSecond, TimeunitMinute, TimeunitHour} {
		if _, exists := duration[timeunit]; exists {
			return PeriodTime
		}
	}

	// Day indicates day-level precision
	if _, exists := duration[TimeunitDay]; exists {
		return PeriodDay
	}

	// Week indicates week-level precision
	if _, exists := duration[TimeunitWeek]; exists {
		return PeriodWeek
	}

	// Month indicates month-level precision
	if _, exists := duration[TimeunitMonth]; exists {
		return PeriodMonth
	}

	// Year, decade, or quarter indicate year-level precision
	for _, timeunit := range []Timeunit{TimeunitYear, TimeunitDecade, TimeunitQuarter} {
		if _, exists := duration[timeunit]; exists {
			return PeriodYear
		}
	}

	// Default to day if no specific duration is found
	return PeriodDay
}

// DeterminePeriodFromComponents determines the granularity/period based on which
// components are certain (explicitly mentioned). The period represents the finest
// granularity of date/time information that was directly parsed.
func determinePeriodFromComponents(pc *parsingComponents) Period {
	if pc == nil {
		return PeriodUnknown
	}

	// If any time components (hour, minute, second) are certain, it's time-level
	if pc.IsCertain(ComponentHour) || pc.IsCertain(ComponentMinute) || pc.IsCertain(ComponentSecond) {
		return PeriodTime
	}

	// If day is certain, it's day-level
	if pc.IsCertain(ComponentDay) {
		return PeriodDay
	}

	// If weekday is certain (without day/month), it's week-level
	if pc.IsCertain(ComponentWeekday) && !pc.IsCertain(ComponentDay) {
		return PeriodWeek
	}

	// If month is certain (without day), it's month-level
	if pc.IsCertain(ComponentMonth) {
		return PeriodMonth
	}

	// If only year is certain, it's year-level
	if pc.IsCertain(ComponentYear) {
		return PeriodYear
	}

	// Default to unknown if nothing is certain
	return PeriodUnknown
}

// ============================================================================
// Relative date creation
// ============================================================================

// CreateRelativeFromReference creates a ParsingComponents from a duration relative to the reference.
// It handles date-only durations (implies time) and time durations (assigns both date and time).
// This is used for parsing relative expressions like "in 3 days", "2 hours ago", etc.
// Returns nil if the duration calculation fails (e.g., overflow).
func createRelativeFromReference(reference *referenceWithTimezone, duration Duration) *parsingComponents {
	if duration == nil {
		duration = Duration{}
	}

	date, err := addDuration(reference.GetDateWithAdjustedTimezone(), duration)
	if err != nil {
		// Duration calculation failed - return nil to indicate invalid result
		return nil
	}

	components := newParsingComponents(reference, nil)
	components.AddTag("result/relativeDate")

	// Determine and set the period based on the duration
	period := determinePeriodFromDuration(duration)
	components.SetPeriod(period)

	// Check if duration contains time components
	hasTimeComponents := false
	for _, timeunit := range []Timeunit{TimeunitHour, TimeunitMinute, TimeunitSecond, TimeunitMillisecond} {
		if _, exists := duration[timeunit]; exists {
			hasTimeComponents = true
			break
		}
	}

	if hasTimeComponents {
		// Duration includes time - assign both date and time as certain
		components.AddTag("result/relativeDateAndTime")
		components.AssignSimilarTime(date)
		components.AssignSimilarDate(date)
		components.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	} else {
		// Duration is date-only - imply time components
		components.ImplySimilarTime(date)
		components.Imply(ComponentTimezoneOffset, reference.GetTimezoneOffset())

		// Handle different date granularities
		if _, hasDayDuration := duration[TimeunitDay]; hasDayDuration {
			// Day duration - assign day, month, year and weekday
			components.Assign(ComponentDay, date.Day())
			components.Assign(ComponentMonth, int(date.Month()))
			components.Assign(ComponentYear, date.Year())
			components.Assign(ComponentWeekday, int(date.Weekday()))
		} else if _, hasWeekDuration := duration[TimeunitWeek]; hasWeekDuration {
			// Week duration - assign day, month, year and imply weekday
			components.Assign(ComponentDay, date.Day())
			components.Assign(ComponentMonth, int(date.Month()))
			components.Assign(ComponentYear, date.Year())
			components.Imply(ComponentWeekday, int(date.Weekday()))
		} else {
			// Month/year duration - imply day
			components.Imply(ComponentDay, date.Day())

			if _, hasMonthDuration := duration[TimeunitMonth]; hasMonthDuration {
				// Month duration - assign month and year
				components.Assign(ComponentMonth, int(date.Month()))
				components.Assign(ComponentYear, date.Year())
			} else {
				// Imply month
				components.Imply(ComponentMonth, int(date.Month()))

				if _, hasYearDuration := duration[TimeunitYear]; hasYearDuration {
					// Year duration - assign year
					components.Assign(ComponentYear, date.Year())
				} else if _, hasQuarterDuration := duration[TimeunitQuarter]; hasQuarterDuration {
					// Quarter duration - assign year
					components.Assign(ComponentYear, date.Year())
				} else {
					// Imply year
					components.Imply(ComponentYear, date.Year())
				}
			}
		}
	}

	return components
}

// AddDurationAsImplied adds the duration to the current components and implies the result.
// This is useful for modifying existing parsing components with a relative offset.
// Returns nil if the duration calculation fails (e.g., overflow).
func (pc *parsingComponents) AddDurationAsImplied(duration Duration) *parsingComponents {
	// Get the current date from this component
	currentDate := pc.Date()

	// Add the duration
	newDate, err := addDuration(currentDate, duration)
	if err != nil {
		// Duration calculation failed - return nil to indicate invalid result
		return nil
	}

	// Imply the new date components
	pc.ImplySimilarDate(newDate)
	pc.ImplySimilarTime(newDate)

	return pc
}
