package kronos

import "time"

// referenceWithTimezone represents a reference date/time with an optional timezone offset.
// It is used as the reference point for parsing relative dates and times.
type referenceWithTimezone struct {
	instant        time.Time
	timezoneOffset *int
}

// newReferenceWithTimezone creates a new referenceWithTimezone with the given instant and timezone offset.
// If instant is zero, the current time is used.
// If timezoneOffset is nil, the system timezone is used.
func newReferenceWithTimezone(instant time.Time, timezoneOffset *int) *referenceWithTimezone {
	if instant.IsZero() {
		instant = time.Now()
	}
	return &referenceWithTimezone{
		instant:        instant,
		timezoneOffset: timezoneOffset,
	}
}

// fromInput creates a ReferenceWithTimezone from either a ParsingReference or a time.Time.
// It also handles timezone conversion using the provided timezoneOverrides.
func fromInput(input any, timezoneOverrides TimezoneAbbrMap) *referenceWithTimezone {
	if input == nil {
		return newReferenceWithTimezone(time.Time{}, nil)
	}

	switch v := input.(type) {
	case time.Time:
		return newReferenceWithTimezone(v, nil)
	case parsingReference:
		instant := time.Now()
		if v.Instant != nil {
			instant = *v.Instant
		}

		var timezoneOffset *int
		if v.Timezone != nil {
			timezoneOffset = toTimezoneOffset(v.Timezone, instant, timezoneOverrides)
		}

		return newReferenceWithTimezone(instant, timezoneOffset)
	default:
		return newReferenceWithTimezone(time.Time{}, nil)
	}
}

// GetDateWithAdjustedTimezone returns a time.Time with the year, month, day, hour, minute, second
// equal to the reference. The output's instant is NOT the reference's instant when the reference's
// and system's timezone are different.
func (r *referenceWithTimezone) GetDateWithAdjustedTimezone() time.Time {
	date := r.instant
	if r.timezoneOffset != nil {
		adjustment := r.GetSystemTimezoneAdjustmentMinute(r.instant, nil)
		date = date.Add(time.Duration(-adjustment) * time.Minute)
	}
	return date
}

// GetSystemTimezoneAdjustmentMinute returns the number of minutes difference between
// the system's timezone and the reference timezone.
func (r *referenceWithTimezone) GetSystemTimezoneAdjustmentMinute(date time.Time, overrideTimezoneOffset *int) int {
	if date.IsZero() || date.Unix() < 0 {
		// Javascript date timezone calculation got effect when the time epoch < 0
		date = time.Now()
	}

	_, currentOffset := date.Zone()
	currentTimezoneOffset := currentOffset / 60

	targetTimezoneOffset := currentTimezoneOffset
	if overrideTimezoneOffset != nil {
		targetTimezoneOffset = *overrideTimezoneOffset
	} else if r.timezoneOffset != nil {
		targetTimezoneOffset = *r.timezoneOffset
	}

	return currentTimezoneOffset - targetTimezoneOffset
}

// GetTimezoneOffset returns the timezone offset in minutes.
// If no timezone offset is set, it returns the system timezone offset.
func (r *referenceWithTimezone) GetTimezoneOffset() int {
	if r.timezoneOffset != nil {
		return *r.timezoneOffset
	}
	_, offset := r.instant.Zone()
	return offset / 60
}

// Instant returns the reference instant.
func (r *referenceWithTimezone) Instant() time.Time {
	return r.instant
}
