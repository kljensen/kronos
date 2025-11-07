package kronos

// Test helper functions for creating casual reference components.
// These functions are only used by tests to verify parsing behavior.
// They are kept in the main package to avoid export/import complexity.

// now returns a ParsingComponents representing the current moment.
// Both date and time components are certain, and it includes timezone offset.
func now(reference *referenceWithTimezone) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	assignSimilarDate(component, targetDate)
	assignSimilarTime(component, targetDate)
	component.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	component.AddTag("casualReference/now")
	component.SetPeriod(PeriodTime)

	return component
}

// today returns a ParsingComponents representing today's date.
// Date components are certain, time components are implied.
func today(reference *referenceWithTimezone) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	assignSimilarDate(component, targetDate)
	implySimilarTime(component, targetDate)
	component.Delete(ComponentMeridiem)
	component.AddTag("casualReference/today")
	component.SetPeriod(PeriodDay)

	return component
}

// yesterday returns a ParsingComponents representing yesterday's date.
// Date components are certain, time components are implied.
func yesterday(reference *referenceWithTimezone) *parsingComponents {
	component := theDayBefore(reference, 1)
	component.AddTag("casualReference/yesterday")
	component.SetPeriod(PeriodDay)
	return component
}

// tomorrow returns a ParsingComponents representing tomorrow's date.
// Date components are certain, time components are implied.
func tomorrow(reference *referenceWithTimezone) *parsingComponents {
	component := theDayAfter(reference, 1)
	component.AddTag("casualReference/tomorrow")
	component.SetPeriod(PeriodDay)
	return component
}

// theDayBefore returns a ParsingComponents representing n days before the reference date.
// Date components are certain, time components are implied.
func theDayBefore(reference *referenceWithTimezone, nDays int) *parsingComponents {
	return theDayAfter(reference, -nDays)
}

// theDayAfter returns a ParsingComponents representing n days after the reference date.
// Date components are certain, time components are implied.
func theDayAfter(reference *referenceWithTimezone, nDays int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	newDate := targetDate.AddDate(0, 0, nDays)

	assignSimilarDate(component, newDate)
	implySimilarTime(component, newDate)
	component.Delete(ComponentMeridiem)
	component.SetPeriod(PeriodDay)

	return component
}

// tonight returns a ParsingComponents representing tonight.
// Date components are certain, time is implied (default: 10 PM / 22:00).
func tonight(reference *referenceWithTimezone) *parsingComponents {
	return tonightWithHour(reference, 22)
}

// tonightWithHour returns a ParsingComponents representing tonight with a specific implied hour.
func tonightWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	assignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, 1) // PM
	component.AddTag("casualReference/tonight")
	component.SetPeriod(PeriodDay)

	return component
}

// lastNight returns a ParsingComponents representing last night.
// If the reference time is before 6 AM, it refers to the previous night.
// Otherwise, it refers to the night of the current day.
func lastNight(reference *referenceWithTimezone) *parsingComponents {
	return lastNightWithHour(reference, 0)
}

// lastNightWithHour returns a ParsingComponents representing last night with a specific implied hour.
func lastNightWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	// If it's very early morning (before 6 AM), "last night" refers to yesterday
	if targetDate.Hour() < 6 {
		targetDate = targetDate.AddDate(0, 0, -1)
	}

	assignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/lastNight")
	component.SetPeriod(PeriodDay)

	return component
}

// evening returns a ParsingComponents representing evening time.
// Time is implied (default: 8 PM / 20:00).
func evening(reference *referenceWithTimezone) *parsingComponents {
	return eveningWithHour(reference, 20)
}

// eveningWithHour returns a ParsingComponents representing evening with a specific implied hour.
func eveningWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, 1) // PM
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/evening")
	component.SetPeriod(PeriodTime)

	return component
}

// yesterdayEvening returns a ParsingComponents representing yesterday evening.
func yesterdayEvening(reference *referenceWithTimezone) *parsingComponents {
	return yesterdayEveningWithHour(reference, 20)
}

// yesterdayEveningWithHour returns a ParsingComponents representing yesterday evening with a specific hour.
func yesterdayEveningWithHour(reference *referenceWithTimezone, implyHour int) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	targetDate = targetDate.AddDate(0, 0, -1)

	assignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, 1) // PM
	component.AddTag("casualReference/yesterday")
	component.AddTag("casualReference/evening")
	component.SetPeriod(PeriodTime)

	return component
}

// midnight returns a ParsingComponents representing midnight.
// If the reference time is after 2 AM, it refers to the coming midnight (next day).
// Otherwise, it refers to the current midnight.
func midnight(reference *referenceWithTimezone) *parsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	// Unless it's very early morning (0-2 AM), assume midnight refers to the coming midnight
	if targetDate.Hour() > 2 {
		duration := Duration{TimeunitDay: 1}
		newDate, err := addDuration(targetDate, duration)
		if err != nil {
			// Duration calculation failed - return nil
			return nil
		}
		implySimilarDate(component, newDate)
	}

	component.Assign(ComponentHour, 0)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/midnight")
	component.SetPeriod(PeriodTime)

	return component
}

// morning returns a ParsingComponents representing morning time.
// Time is implied (default: 6 AM / 06:00).
func morning(reference *referenceWithTimezone) *parsingComponents {
	return morningWithHour(reference, 6)
}

// morningWithHour returns a ParsingComponents representing morning with a specific implied hour.
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

// afternoon returns a ParsingComponents representing afternoon time.
// Time is implied (default: 3 PM / 15:00).
func afternoon(reference *referenceWithTimezone) *parsingComponents {
	return afternoonWithHour(reference, 15)
}

// afternoonWithHour returns a ParsingComponents representing afternoon with a specific implied hour.
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

// noon returns a ParsingComponents representing noon (12:00 PM).
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
