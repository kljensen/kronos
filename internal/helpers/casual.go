package helpers

import "github.com/kljensen/kronos"

// Now returns components representing the current moment.
// Both date and time components are certain, and it includes timezone offset.
func Now(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	component.AssignSimilarDate(targetDate)
	component.AssignSimilarTime(targetDate)
	component.Assign(kronos.ComponentTimezoneOffset, reference.GetTimezoneOffset())
	component.AddTag("casualReference/now")
	component.SetPeriod(kronos.PeriodTime)

	return component
}

// Today returns components representing today's date.
// Date components are certain, time components are implied.
func Today(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	component.AssignSimilarDate(targetDate)
	component.ImplySimilarTime(targetDate)
	component.Delete(kronos.ComponentMeridiem)
	component.AddTag("casualReference/today")
	component.SetPeriod(kronos.PeriodDay)

	return component
}

// Yesterday returns components representing yesterday's date.
// Date components are certain, time components are implied.
func Yesterday(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	component := TheDayBefore(reference, 1)
	component.AddTag("casualReference/yesterday")
	component.SetPeriod(kronos.PeriodDay)
	return component
}

// Tomorrow returns components representing tomorrow's date.
// Date components are certain, time components are implied.
func Tomorrow(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	component := TheDayAfter(reference, 1)
	component.AddTag("casualReference/tomorrow")
	component.SetPeriod(kronos.PeriodDay)
	return component
}

// TheDayBefore returns components representing n days before the reference date.
// Date components are certain, time components are implied.
func TheDayBefore(reference *kronos.InternalReferenceWithTimezone, nDays int) *kronos.InternalParsingComponents {
	return TheDayAfter(reference, -nDays)
}

// TheDayAfter returns components representing n days after the reference date.
// Date components are certain, time components are implied.
func TheDayAfter(reference *kronos.InternalReferenceWithTimezone, nDays int) *kronos.InternalParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	newDate := targetDate.AddDate(0, 0, nDays)

	component.AssignSimilarDate(newDate)
	component.ImplySimilarTime(newDate)
	component.Delete(kronos.ComponentMeridiem)
	component.SetPeriod(kronos.PeriodDay)

	return component
}

// Midnight returns components representing midnight.
// If the reference time is after 2 AM, it refers to the coming midnight (next day).
// Otherwise, it refers to the current midnight.
func Midnight(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	// Unless it's very early morning (0-2 AM), assume midnight refers to the coming midnight
	if targetDate.Hour() > 2 {
		duration := kronos.Duration{kronos.TimeunitDay: 1}
		newDate, err := AddDuration(targetDate, duration)
		if err != nil {
			// Duration calculation failed - return nil
			return nil
		}
		component.ImplySimilarDate(newDate)
	}

	component.Assign(kronos.ComponentHour, 0)
	component.Imply(kronos.ComponentMinute, 0)
	component.Imply(kronos.ComponentSecond, 0)
	component.Imply(kronos.ComponentMillisecond, 0)
	component.AddTag("casualReference/midnight")
	component.SetPeriod(kronos.PeriodTime)

	return component
}

// Noon returns components representing noon (12:00 PM).
func Noon(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	component := NewParsingComponents(reference, nil)

	component.Imply(kronos.ComponentMeridiem, 1) // PM
	component.Assign(kronos.ComponentHour, 12)
	component.Imply(kronos.ComponentMinute, 0)
	component.Imply(kronos.ComponentSecond, 0)
	component.Imply(kronos.ComponentMillisecond, 0)
	component.AddTag("casualReference/noon")
	component.SetPeriod(kronos.PeriodTime)

	return component
}

// Morning returns components representing morning time.
// Time is implied (default: 6 AM / 06:00).
func Morning(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	return MorningWithHour(reference, 6)
}

// MorningWithHour returns components representing morning with a specific implied hour.
func MorningWithHour(reference *kronos.InternalReferenceWithTimezone, implyHour int) *kronos.InternalParsingComponents {
	component := NewParsingComponents(reference, nil)

	component.Imply(kronos.ComponentMeridiem, 0) // AM
	component.Imply(kronos.ComponentHour, implyHour)
	component.Imply(kronos.ComponentMinute, 0)
	component.Imply(kronos.ComponentSecond, 0)
	component.Imply(kronos.ComponentMillisecond, 0)
	component.AddTag("casualReference/morning")
	component.SetPeriod(kronos.PeriodTime)

	return component
}

// Afternoon returns components representing afternoon time.
// Time is implied (default: 3 PM / 15:00).
func Afternoon(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	return AfternoonWithHour(reference, 15)
}

// AfternoonWithHour returns components representing afternoon with a specific implied hour.
func AfternoonWithHour(reference *kronos.InternalReferenceWithTimezone, implyHour int) *kronos.InternalParsingComponents {
	component := NewParsingComponents(reference, nil)

	component.Imply(kronos.ComponentMeridiem, 1) // PM
	component.Imply(kronos.ComponentHour, implyHour)
	component.Imply(kronos.ComponentMinute, 0)
	component.Imply(kronos.ComponentSecond, 0)
	component.Imply(kronos.ComponentMillisecond, 0)
	component.AddTag("casualReference/afternoon")
	component.SetPeriod(kronos.PeriodTime)

	return component
}

// Evening returns components representing evening time.
// Time is implied (default: 8 PM / 20:00).
func Evening(reference *kronos.InternalReferenceWithTimezone) *kronos.InternalParsingComponents {
	return EveningWithHour(reference, 20)
}

// EveningWithHour returns components representing evening with a specific implied hour.
func EveningWithHour(reference *kronos.InternalReferenceWithTimezone, implyHour int) *kronos.InternalParsingComponents {
	component := NewParsingComponents(reference, nil)

	component.Imply(kronos.ComponentMeridiem, 1) // PM
	component.Imply(kronos.ComponentHour, implyHour)
	component.AddTag("casualReference/evening")
	component.SetPeriod(kronos.PeriodTime)

	return component
}

// TonightWithHour returns components representing tonight with a specific implied hour.
func TonightWithHour(reference *kronos.InternalReferenceWithTimezone, implyHour int) *kronos.InternalParsingComponents {
	targetDate := reference.GetDateWithAdjustedTimezone()
	component := NewParsingComponents(reference, nil)

	component.AssignSimilarDate(targetDate)
	component.Imply(kronos.ComponentHour, implyHour)
	component.Imply(kronos.ComponentMeridiem, 1) // PM
	component.AddTag("casualReference/tonight")
	component.SetPeriod(kronos.PeriodDay)

	return component
}
