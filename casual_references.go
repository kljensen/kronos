package kronos

// Now returns a ParsingComponents representing the current moment.
// Both date and time components are certain, and it includes timezone offset.
func Now(reference *ReferenceWithTimezone) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	AssignSimilarDate(component, targetDate)
	AssignSimilarTime(component, targetDate)
	component.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	component.AddTag("casualReference/now")
	component.SetPeriod(PeriodTime)

	return component
}

// Today returns a ParsingComponents representing today's date.
// Date components are certain, time components are implied.
func Today(reference *ReferenceWithTimezone) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	AssignSimilarDate(component, targetDate)
	ImplySimilarTime(component, targetDate)
	component.Delete(ComponentMeridiem)
	component.AddTag("casualReference/today")
	component.SetPeriod(PeriodDay)

	return component
}

// Yesterday returns a ParsingComponents representing yesterday's date.
// Date components are certain, time components are implied.
func Yesterday(reference *ReferenceWithTimezone) *ParsingComponents {
	component := TheDayBefore(reference, 1)
	component.AddTag("casualReference/yesterday")
	component.SetPeriod(PeriodDay)
	return component
}

// Tomorrow returns a ParsingComponents representing tomorrow's date.
// Date components are certain, time components are implied.
func Tomorrow(reference *ReferenceWithTimezone) *ParsingComponents {
	component := TheDayAfter(reference, 1)
	component.AddTag("casualReference/tomorrow")
	component.SetPeriod(PeriodDay)
	return component
}

// TheDayBefore returns a ParsingComponents representing n days before the reference date.
// Date components are certain, time components are implied.
func TheDayBefore(reference *ReferenceWithTimezone, nDays int) *ParsingComponents {
	return TheDayAfter(reference, -nDays)
}

// TheDayAfter returns a ParsingComponents representing n days after the reference date.
// Date components are certain, time components are implied.
func TheDayAfter(reference *ReferenceWithTimezone, nDays int) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	newDate := targetDate.AddDate(0, 0, nDays)

	AssignSimilarDate(component, newDate)
	ImplySimilarTime(component, newDate)
	component.Delete(ComponentMeridiem)
	component.SetPeriod(PeriodDay)

	return component
}

// Tonight returns a ParsingComponents representing tonight.
// Date components are certain, time is implied (default: 10 PM / 22:00).
func Tonight(reference *ReferenceWithTimezone) *ParsingComponents {
	return TonightWithHour(reference, 22)
}

// TonightWithHour returns a ParsingComponents representing tonight with a specific implied hour.
func TonightWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	AssignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.AddTag("casualReference/tonight")
	component.SetPeriod(PeriodDay)

	return component
}

// LastNight returns a ParsingComponents representing last night.
// If the reference time is before 6 AM, it refers to the previous night.
// Otherwise, it refers to the night of the current day.
func LastNight(reference *ReferenceWithTimezone) *ParsingComponents {
	return LastNightWithHour(reference, 0)
}

// LastNightWithHour returns a ParsingComponents representing last night with a specific implied hour.
func LastNightWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	// If it's very early morning (before 6 AM), "last night" refers to yesterday
	if targetDate.Hour() < 6 {
		targetDate = targetDate.AddDate(0, 0, -1)
	}

	AssignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/lastNight")
	component.SetPeriod(PeriodDay)

	return component
}

// Evening returns a ParsingComponents representing evening time.
// Time is implied (default: 8 PM / 20:00).
func Evening(reference *ReferenceWithTimezone) *ParsingComponents {
	return EveningWithHour(reference, 20)
}

// EveningWithHour returns a ParsingComponents representing evening with a specific implied hour.
func EveningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	component := NewParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/evening")
	component.SetPeriod(PeriodTime)

	return component
}

// YesterdayEvening returns a ParsingComponents representing yesterday evening.
func YesterdayEvening(reference *ReferenceWithTimezone) *ParsingComponents {
	return YesterdayEveningWithHour(reference, 20)
}

// YesterdayEveningWithHour returns a ParsingComponents representing yesterday evening with a specific hour.
func YesterdayEveningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	targetDate = targetDate.AddDate(0, 0, -1)

	AssignSimilarDate(component, targetDate)
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
func Midnight(reference *ReferenceWithTimezone) *ParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	// Unless it's very early morning (0-2 AM), assume midnight refers to the coming midnight
	if targetDate.Hour() > 2 {
		duration := Duration{TimeunitDay: 1}
		newDate := AddDuration(targetDate, duration)
		ImplySimilarDate(component, newDate)
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
func Morning(reference *ReferenceWithTimezone) *ParsingComponents {
	return MorningWithHour(reference, 6)
}

// MorningWithHour returns a ParsingComponents representing morning with a specific implied hour.
func MorningWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	component := NewParsingComponents(reference, nil)

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
func Afternoon(reference *ReferenceWithTimezone) *ParsingComponents {
	return AfternoonWithHour(reference, 15)
}

// AfternoonWithHour returns a ParsingComponents representing afternoon with a specific implied hour.
func AfternoonWithHour(reference *ReferenceWithTimezone, implyHour int) *ParsingComponents {
	component := NewParsingComponents(reference, nil)

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
func Noon(reference *ReferenceWithTimezone) *ParsingComponents {
	component := NewParsingComponents(reference, nil)

	component.Imply(ComponentMeridiem, int(MeridiemPM))
	component.Assign(ComponentHour, 12)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/noon")
	component.SetPeriod(PeriodTime)

	return component
}
