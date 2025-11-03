package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateRelativeFromReference_DateOnly(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitDay: 3,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Date components should be certain
	assert.True(t, components.IsCertain(ComponentYear))
	assert.True(t, components.IsCertain(ComponentMonth))
	assert.True(t, components.IsCertain(ComponentDay))
	assert.True(t, components.IsCertain(ComponentWeekday))

	// Time components should be implied
	assert.False(t, components.IsCertain(ComponentHour))
	assert.False(t, components.IsCertain(ComponentMinute))
	assert.False(t, components.IsCertain(ComponentSecond))

	// Should be 3 days after reference (June 18, 2020)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 18, *components.Get(ComponentDay))

	// Check tags
	tags := components.Tags()
	assert.True(t, tags["result/relativeDate"])
	assert.False(t, tags["result/relativeDateAndTime"])
}

func TestCreateRelativeFromReference_WithTime(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitHour: 2,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Both date and time components should be certain
	assert.True(t, components.IsCertain(ComponentYear))
	assert.True(t, components.IsCertain(ComponentMonth))
	assert.True(t, components.IsCertain(ComponentDay))
	assert.True(t, components.IsCertain(ComponentHour))
	assert.True(t, components.IsCertain(ComponentMinute))
	assert.True(t, components.IsCertain(ComponentSecond))
	assert.True(t, components.IsCertain(ComponentTimezoneOffset))

	// Should be 2 hours after reference
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
	assert.Equal(t, 16, *components.Get(ComponentHour))
	assert.Equal(t, 30, *components.Get(ComponentMinute))

	// Check tags
	tags := components.Tags()
	assert.True(t, tags["result/relativeDate"])
	assert.True(t, tags["result/relativeDateAndTime"])
}

func TestCreateRelativeFromReference_WeekDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitWeek: 2,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Date components should be certain
	assert.True(t, components.IsCertain(ComponentYear))
	assert.True(t, components.IsCertain(ComponentMonth))
	assert.True(t, components.IsCertain(ComponentDay))

	// Weekday should be implied for week duration
	assert.False(t, components.IsCertain(ComponentWeekday))

	// Should be 2 weeks (14 days) after reference (June 29, 2020)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 29, *components.Get(ComponentDay))

	// Check tags
	tags := components.Tags()
	assert.True(t, tags["result/relativeDate"])
	assert.False(t, tags["result/relativeDateAndTime"])
}

func TestCreateRelativeFromReference_MonthDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitMonth: 3,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Month and year should be certain
	assert.True(t, components.IsCertain(ComponentMonth))
	assert.True(t, components.IsCertain(ComponentYear))

	// Day should be implied
	assert.False(t, components.IsCertain(ComponentDay))

	// Should be 3 months after reference (September 15, 2020)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 9, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))

	// Check tags
	tags := components.Tags()
	assert.True(t, tags["result/relativeDate"])
	assert.False(t, tags["result/relativeDateAndTime"])
}

func TestCreateRelativeFromReference_YearDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitYear: 1,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Year should be certain
	assert.True(t, components.IsCertain(ComponentYear))

	// Month and day should be implied
	assert.False(t, components.IsCertain(ComponentMonth))
	assert.False(t, components.IsCertain(ComponentDay))

	// Should be 1 year after reference (June 15, 2021)
	assert.Equal(t, 2021, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
}

func TestCreateRelativeFromReference_QuarterDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitQuarter: 1,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Year should be certain
	assert.True(t, components.IsCertain(ComponentYear))

	// Month and day should be implied
	assert.False(t, components.IsCertain(ComponentMonth))
	assert.False(t, components.IsCertain(ComponentDay))

	// Should be 1 quarter (3 months) after reference (September 15, 2020)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 9, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
}

func TestCreateRelativeFromReference_NegativeDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitDay: -5,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Should be 5 days before reference (June 10, 2020)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 10, *components.Get(ComponentDay))
}

func TestCreateRelativeFromReference_MixedDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	duration := Duration{
		TimeunitDay:    2,
		TimeunitHour:   3,
		TimeunitMinute: 15,
	}

	components := CreateRelativeFromReference(reference, duration)

	// Both date and time should be certain (has time components)
	assert.True(t, components.IsCertain(ComponentYear))
	assert.True(t, components.IsCertain(ComponentDay))
	assert.True(t, components.IsCertain(ComponentHour))
	assert.True(t, components.IsCertain(ComponentMinute))

	// Should be 2 days, 3 hours, 15 minutes after reference
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 17, *components.Get(ComponentDay))
	assert.Equal(t, 17, *components.Get(ComponentHour))
	assert.Equal(t, 45, *components.Get(ComponentMinute))

	// Check tags
	tags := components.Tags()
	assert.True(t, tags["result/relativeDate"])
	assert.True(t, tags["result/relativeDateAndTime"])
}

func TestCreateRelativeFromReference_EmptyDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	components := CreateRelativeFromReference(reference, EmptyDuration)

	// Should be same as reference date (no change)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
}

func TestCreateRelativeFromReference_NilDuration(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	components := CreateRelativeFromReference(reference, nil)

	// Should handle nil duration gracefully (treat as empty)
	assert.NotNil(t, components)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
}

func TestAddDurationAsImplied(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	components := NewParsingComponents(reference, nil)

	// Set some initial certain values
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 10)

	// Add 5 days as implied
	duration := Duration{
		TimeunitDay: 5,
	}
	components.AddDurationAsImplied(duration)

	// Certain values should remain certain
	assert.True(t, components.IsCertain(ComponentYear))
	assert.True(t, components.IsCertain(ComponentMonth))
	assert.True(t, components.IsCertain(ComponentDay))

	// But the values should not change (because they're already certain)
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 10, *components.Get(ComponentDay))
}

func TestAddDurationAsImplied_ImpliedValues(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := NewReferenceWithTimezone(refDate, nil)

	components := NewParsingComponents(reference, nil)

	// Don't set day as certain - it will be implied

	// Add 5 days as implied
	duration := Duration{
		TimeunitDay: 5,
	}
	components.AddDurationAsImplied(duration)

	// Day should still be implied, but updated
	assert.False(t, components.IsCertain(ComponentDay))

	// The date should be advanced by 5 days from the reference
	date := components.Date()
	expected := refDate.AddDate(0, 0, 5)

	// Compare year, month, day
	assert.Equal(t, expected.Year(), date.Year())
	assert.Equal(t, expected.Month(), date.Month())
	assert.Equal(t, expected.Day(), date.Day())
}
