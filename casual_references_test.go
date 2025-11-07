package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 45, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := now(reference)

	// All date and time components should be certain
	assert.True(t, component.IsCertain(ComponentYear))
	assert.True(t, component.IsCertain(ComponentMonth))
	assert.True(t, component.IsCertain(ComponentDay))
	assert.True(t, component.IsCertain(ComponentHour))
	assert.True(t, component.IsCertain(ComponentMinute))
	assert.True(t, component.IsCertain(ComponentSecond))
	assert.True(t, component.IsCertain(ComponentTimezoneOffset))

	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 15, *component.Get(ComponentDay))
	assert.Equal(t, 14, *component.Get(ComponentHour))
	assert.Equal(t, 30, *component.Get(ComponentMinute))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/now"])
}

func TestToday(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 45, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := today(reference)

	// Date components should be certain
	assert.True(t, component.IsCertain(ComponentYear))
	assert.True(t, component.IsCertain(ComponentMonth))
	assert.True(t, component.IsCertain(ComponentDay))

	// Time components should be implied
	assert.False(t, component.IsCertain(ComponentHour))
	assert.False(t, component.IsCertain(ComponentMinute))
	assert.False(t, component.IsCertain(ComponentSecond))

	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 15, *component.Get(ComponentDay))

	// Meridiem should be deleted
	assert.Nil(t, component.Get(ComponentMeridiem))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/today"])
}

func TestYesterday(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := yesterday(reference)

	assert.True(t, component.IsCertain(ComponentYear))
	assert.True(t, component.IsCertain(ComponentMonth))
	assert.True(t, component.IsCertain(ComponentDay))

	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 14, *component.Get(ComponentDay))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/yesterday"])
}

func TestTomorrow(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := tomorrow(reference)

	assert.True(t, component.IsCertain(ComponentYear))
	assert.True(t, component.IsCertain(ComponentMonth))
	assert.True(t, component.IsCertain(ComponentDay))

	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 16, *component.Get(ComponentDay))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/tomorrow"])
}

func TestTheDayAfter(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	// 3 days after
	component := theDayAfter(reference, 3)

	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 18, *component.Get(ComponentDay))
}

func TestTheDayBefore(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	// 3 days before
	component := theDayBefore(reference, 3)

	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 12, *component.Get(ComponentDay))
}

func TestTonight(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := tonight(reference)

	// Date should be certain
	assert.True(t, component.IsCertain(ComponentYear))
	assert.True(t, component.IsCertain(ComponentMonth))
	assert.True(t, component.IsCertain(ComponentDay))

	// Time should be implied (22:00 / 10 PM by default)
	assert.False(t, component.IsCertain(ComponentHour))
	assert.Equal(t, 22, *component.Get(ComponentHour))
	assert.Equal(t, int(1), *component.Get(ComponentMeridiem))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/tonight"])
}

func TestTonightWithHour(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := tonightWithHour(reference, 20)

	assert.Equal(t, 20, *component.Get(ComponentHour))
}

func TestLastNight_EarlyMorning(t *testing.T) {
	// Before 6 AM - should refer to previous day's night
	refDate := time.Date(2020, 6, 15, 3, 0, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := lastNight(reference)

	// Should be June 14
	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 14, *component.Get(ComponentDay))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/lastNight"])
}

func TestLastNight_AfterMorning(t *testing.T) {
	// After 6 AM - should refer to same day's night
	refDate := time.Date(2020, 6, 15, 14, 0, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := lastNight(reference)

	// Should be June 15
	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 15, *component.Get(ComponentDay))
}

func TestEvening(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := evening(reference)

	// Time should be implied (20:00 / 8 PM by default)
	assert.False(t, component.IsCertain(ComponentHour))
	assert.Equal(t, 20, *component.Get(ComponentHour))
	assert.Equal(t, int(1), *component.Get(ComponentMeridiem))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/evening"])
}

func TestYesterdayEvening(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := yesterdayEvening(reference)

	// Date should be yesterday
	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 14, *component.Get(ComponentDay))

	// Time should be implied evening
	assert.Equal(t, 20, *component.Get(ComponentHour))
	assert.Equal(t, int(1), *component.Get(ComponentMeridiem))

	// Check tags
	tags := component.Tags()
	assert.True(t, tags["casualReference/yesterday"])
	assert.True(t, tags["casualReference/evening"])
}

func TestMidnight_EarlyMorning(t *testing.T) {
	// Before 2 AM - refers to current midnight
	refDate := time.Date(2020, 6, 15, 1, 0, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := midnight(reference)

	// Should be same day (June 15)
	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 15, *component.Get(ComponentDay))
	assert.Equal(t, 0, *component.Get(ComponentHour))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/midnight"])
}

func TestMidnight_AfterMorning(t *testing.T) {
	// After 2 AM - refers to coming midnight (next day)
	refDate := time.Date(2020, 6, 15, 14, 0, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := midnight(reference)

	// Should be next day (June 16)
	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 6, *component.Get(ComponentMonth))
	assert.Equal(t, 16, *component.Get(ComponentDay))
	assert.Equal(t, 0, *component.Get(ComponentHour))
}

func TestMorning(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := morning(reference)

	// Time should be implied (6:00 AM by default)
	assert.False(t, component.IsCertain(ComponentHour))
	assert.Equal(t, 6, *component.Get(ComponentHour))
	assert.Equal(t, int(0), *component.Get(ComponentMeridiem))
	assert.Equal(t, 0, *component.Get(ComponentMinute))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/morning"])
}

func TestAfternoon(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := afternoon(reference)

	// Time should be implied (15:00 / 3 PM by default)
	assert.False(t, component.IsCertain(ComponentHour))
	assert.Equal(t, 15, *component.Get(ComponentHour))
	assert.Equal(t, int(1), *component.Get(ComponentMeridiem))
	assert.Equal(t, 0, *component.Get(ComponentMinute))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/afternoon"])
}

func TestNoon(t *testing.T) {
	refDate := time.Date(2020, 6, 15, 14, 30, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := noon(reference)

	// Hour should be certain (12:00)
	assert.True(t, component.IsCertain(ComponentHour))
	assert.Equal(t, 12, *component.Get(ComponentHour))
	assert.Equal(t, 0, *component.Get(ComponentMinute))

	// Check tag
	tags := component.Tags()
	assert.True(t, tags["casualReference/noon"])
}

func TestMonthBoundaries(t *testing.T) {
	// Test yesterday crossing month boundary
	refDate := time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := yesterday(reference)

	assert.Equal(t, 2020, *component.Get(ComponentYear))
	assert.Equal(t, 5, *component.Get(ComponentMonth))
	assert.Equal(t, 31, *component.Get(ComponentDay))
}

func TestYearBoundaries(t *testing.T) {
	// Test yesterday crossing year boundary
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	reference := newReferenceWithTimezone(refDate, nil)

	component := yesterday(reference)

	assert.Equal(t, 2019, *component.Get(ComponentYear))
	assert.Equal(t, 12, *component.Get(ComponentMonth))
	assert.Equal(t, 31, *component.Get(ComponentDay))

	// Test tomorrow crossing year boundary
	refDate2 := time.Date(2020, 12, 31, 12, 0, 0, 0, time.UTC)
	reference2 := newReferenceWithTimezone(refDate2, nil)

	component2 := tomorrow(reference2)

	assert.Equal(t, 2021, *component2.Get(ComponentYear))
	assert.Equal(t, 1, *component2.Get(ComponentMonth))
	assert.Equal(t, 1, *component2.Get(ComponentDay))
}
