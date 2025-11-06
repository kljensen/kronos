package kronos

// Now returns a ParsingComponents representing the current moment.
// Both date and time components are certain, and it includes timezone offset.
func now(reference *ReferenceWithTimezone) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	assignSimilarDate(component, targetDate)
	assignSimilarTime(component, targetDate)
	component.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	component.AddTag("casualReference/now")
	component.SetPeriod(PeriodTime)

	return component
}

// Today returns a ParsingComponents representing today's date.
// Date components are certain, time components are implied.
func today(reference *ReferenceWithTimezone) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	assignSimilarDate(component, targetDate)
	implySimilarTime(component, targetDate)
	component.Delete(ComponentMeridiem)
	component.AddTag("casualReference/today")
	component.SetPeriod(PeriodDay)

	return component
}

// Yesterday returns a ParsingComponents representing yesterday's date.
// Date components are certain, time components are implied.
func yesterday(reference *ReferenceWithTimezone) *ParsingComponents {
	component := theDayBefore(reference, 1)
	component.AddTag("casualReference/yesterday")
	component.SetPeriod(PeriodDay)
	return component
}

// Tomorrow returns a ParsingComponents representing tomorrow's date.
// Date components are certain, time components are implied.
func tomorrow(reference *ReferenceWithTimezone) *ParsingComponents {
	component := theDayAfter(reference, 1)
	component.AddTag("casualReference/tomorrow")
	component.SetPeriod(PeriodDay)
	return component
}

// TheDayBefore returns a ParsingComponents representing n days before the reference date.
// Date components are certain, time components are implied.
func theDayBefore(reference *ReferenceWithTimezone, nDays int) *ParsingComponents {
	return theDayAfter(reference, -nDays)
}

// TheDayAfter returns a ParsingComponents representing n days after the reference date.
// Date components are certain, time components are implied.
func theDayAfter(reference *ReferenceWithTimezone, nDays int) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	newDate := targetDate.AddDate(0, 0, nDays)

	assignSimilarDate(component, newDate)
	implySimilarTime(component, newDate)
	component.Delete(ComponentMeridiem)
	component.SetPeriod(PeriodDay)

	return component
}

// Tonight returns a ParsingComponents representing tonight.
// Date components are certain, time is implied (default: 10 PM / 22:00).
func tonight(reference *ReferenceWithTimezone) *ParsingComponents {
	return tonightWithHour(reference, 22)
}

// TonightWithHour returns a ParsingComponents representing tonight with a specific implied hour.
func tonightWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	assignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.AddTag("casualReference/tonight")
	component.SetPeriod(PeriodDay)

	return component
}

// LastNight returns a ParsingComponents representing last night.
// If the reference time is before 6 AM, it refers to the previous night.
// Otherwise, it refers to the night of the current day.
func lastNight(reference *ReferenceWithTimezone) *ParsingComponents {
	return lastNightWithHour(reference, 0)
}

// LastNightWithHour returns a ParsingComponents representing last night with a specific implied hour.
func lastNightWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
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

// Evening returns a ParsingComponents representing evening time.
// Time is implied (default: 8 PM / 20:00).
func evening(reference *ReferenceWithTimezone) *ParsingComponents {
	return eveningWithHour(reference, 20)
}

// EveningWithHour returns a ParsingComponents representing evening with a specific implied hour.
func eveningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/evening")
	component.SetPeriod(PeriodTime)

	return component
}

// YesterdayEvening returns a ParsingComponents representing yesterday evening.
func yesterdayEvening(reference *ReferenceWithTimezone) *ParsingComponents {
	return yesterdayEveningWithHour(reference, 20)
}

// YesterdayEveningWithHour returns a ParsingComponents representing yesterday evening with a specific hour.
func yesterdayEveningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := newParsingComponents(reference, nil)

	targetDate = targetDate.AddDate(0, 0, -1)

	assignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.AddTag("casualReference/yesterday")
	component.AddTag("casualReference/evening")
	component.SetPeriod(PeriodTime)

	return component
}

// Midnight returns a ParsingComponents representing midnight.
// If the reference time is after 2 AM, it refers to the coming midnight (next day).
// Otherwise, it refers to the current midnight.
func midnight(reference *ReferenceWithTimezone) *ParsingComponents {
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

// Morning returns a ParsingComponents representing morning time.
// Time is implied (default: 6 AM / 06:00).
func morning(reference *ReferenceWithTimezone) *ParsingComponents {
	return morningWithHour(reference, 6)
}

// MorningWithHour returns a ParsingComponents representing morning with a specific implied hour.
func morningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, int(MeridiemAM))
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/morning")
	component.SetPeriod(PeriodTime)

	return component
}

// Afternoon returns a ParsingComponents representing afternoon time.
// Time is implied (default: 3 PM / 15:00).
func afternoon(reference *ReferenceWithTimezone) *ParsingComponents {
	return afternoonWithHour(reference, 15)
}

// AfternoonWithHour returns a ParsingComponents representing afternoon with a specific implied hour.
func afternoonWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/afternoon")
	component.SetPeriod(PeriodTime)

	return component
}

// Noon returns a ParsingComponents representing noon (12:00 PM).
func noon(reference *ReferenceWithTimezone) *ParsingComponents {
	component := newParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.Assign(ComponentHour, 12)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/noon")
	component.SetPeriod(PeriodTime)

	return component
}
