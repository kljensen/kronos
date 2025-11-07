package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPipeline_TimezoneConversion_BasicUTCToNY tests basic timezone conversion from UTC to America/New_York
func TestPipeline_TimezoneConversion_BasicUTCToNY(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	// Create a test result with UTC time: 2020-06-15 18:00:00 UTC
	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 18)
	components.Assign(ComponentMinute, 0)
	components.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "June 15, 2020 at 6:00 PM", components, nil)
	results := []*parsingResult{result}

	// Apply timezone conversion
	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)
	assert.Len(t, converted, 1)

	// Verify the components are updated correctly
	// 18:00 UTC = 14:00 EDT (UTC-4 in summer)
	start := converted[0].Start()
	assert.Equal(t, 2020, *start.Get(ComponentYear))
	assert.Equal(t, 6, *start.Get(ComponentMonth))
	assert.Equal(t, 15, *start.Get(ComponentDay))
	assert.Equal(t, 14, *start.Get(ComponentHour))
	assert.Equal(t, 0, *start.Get(ComponentMinute))
	assert.Equal(t, 0, *start.Get(ComponentSecond))
}

func TestPipeline_TimezoneConversion_StartDateReflectsTargetTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)
	startComponents := newParsingComponents(ref, nil)
	startComponents.Assign(ComponentYear, 2020)
	startComponents.Assign(ComponentMonth, 6)
	startComponents.Assign(ComponentDay, 15)
	startComponents.Assign(ComponentHour, 18)
	startComponents.Assign(ComponentMinute, 0)
	startComponents.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "June 15, 2020 at 6:00 PM", startComponents, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	targetLoc, _ := time.LoadLocation("America/New_York")
	expected := time.Date(2020, 6, 15, 14, 0, 0, 0, targetLoc)

	assert.True(t, converted[0].Date().Equal(expected), "converted date %v does not match expected %v", converted[0].Date(), expected)
	assert.Equal(t, targetLoc.String(), converted[0].Date().Location().String())
	assert.True(t, converted[0].RefDate().Equal(expected), "converted ref date %v does not match expected %v", converted[0].RefDate(), expected)
}

// TestPipeline_TimezoneConversion_UTCToTokyo tests conversion to Asia/Tokyo
func TestPipeline_TimezoneConversion_UTCToTokyo(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "Asia/Tokyo"
	pipeline := newPipeline(nil, settings)

	// Create a test result with UTC time: 2020-06-15 18:00:00 UTC
	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 18)
	components.Assign(ComponentMinute, 0)
	components.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "June 15, 2020 at 6:00 PM", components, nil)
	results := []*parsingResult{result}

	// Apply timezone conversion
	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)
	assert.Len(t, converted, 1)

	// Verify the components are updated correctly
	// 18:00 UTC = 03:00 JST next day (UTC+9)
	start := converted[0].Start()
	assert.Equal(t, 2020, *start.Get(ComponentYear))
	assert.Equal(t, 6, *start.Get(ComponentMonth))
	assert.Equal(t, 16, *start.Get(ComponentDay), "Date should advance to next day")
	assert.Equal(t, 3, *start.Get(ComponentHour))
	assert.Equal(t, 0, *start.Get(ComponentMinute))
	assert.Equal(t, 0, *start.Get(ComponentSecond))
}

// TestPipeline_TimezoneConversion_DateBoundaryChange tests when timezone conversion crosses date boundary
func TestPipeline_TimezoneConversion_DateBoundaryChange(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	// Create a test result with UTC time: 2020-06-15 23:00:00 UTC
	// This should become 2020-06-15 19:00:00 EDT (same day)
	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 23)
	components.Assign(ComponentMinute, 0)
	components.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "June 15, 2020 at 11:00 PM", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	assert.Equal(t, 2020, *start.Get(ComponentYear))
	assert.Equal(t, 6, *start.Get(ComponentMonth))
	assert.Equal(t, 15, *start.Get(ComponentDay))
	assert.Equal(t, 19, *start.Get(ComponentHour))

	// Test crossing to previous day
	components2 := newParsingComponents(ref, nil)
	components2.Assign(ComponentYear, 2020)
	components2.Assign(ComponentMonth, 6)
	components2.Assign(ComponentDay, 16)
	components2.Assign(ComponentHour, 3)
	components2.Assign(ComponentMinute, 0)
	components2.Assign(ComponentSecond, 0)

	result2 := newParsingResult(ref, 0, "June 16, 2020 at 3:00 AM", components2, nil)
	results2 := []*parsingResult{result2}

	converted2, err := pipeline.applyTimezoneConversion(results2)
	assert.NoError(t, err)

	start2 := converted2[0].Start()
	assert.Equal(t, 2020, *start2.Get(ComponentYear))
	assert.Equal(t, 6, *start2.Get(ComponentMonth))
	assert.Equal(t, 15, *start2.Get(ComponentDay), "Date should go back to previous day")
	assert.Equal(t, 23, *start2.Get(ComponentHour))
}

// TestPipeline_TimezoneConversion_DSTBoundary tests DST boundary handling
func TestPipeline_TimezoneConversion_DSTBoundary(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 3, 8, 12, 0, 0, 0, utc), nil)

	// Test before DST starts (EST: UTC-5)
	components1 := newParsingComponents(ref, nil)
	components1.Assign(ComponentYear, 2020)
	components1.Assign(ComponentMonth, 3)
	components1.Assign(ComponentDay, 8)
	components1.Assign(ComponentHour, 6) // 6:00 UTC
	components1.Assign(ComponentMinute, 0)
	components1.Assign(ComponentSecond, 0)

	result1 := newParsingResult(ref, 0, "March 8, 2020 at 6:00 AM", components1, nil)
	results1 := []*parsingResult{result1}

	converted1, err := pipeline.applyTimezoneConversion(results1)
	assert.NoError(t, err)

	start1 := converted1[0].Start()
	assert.Equal(t, 1, *start1.Get(ComponentHour), "Should be 1:00 AM EST (UTC-5)")

	// Test after DST starts (EDT: UTC-4)
	components2 := newParsingComponents(ref, nil)
	components2.Assign(ComponentYear, 2020)
	components2.Assign(ComponentMonth, 3)
	components2.Assign(ComponentDay, 8)
	components2.Assign(ComponentHour, 8) // 8:00 UTC
	components2.Assign(ComponentMinute, 0)
	components2.Assign(ComponentSecond, 0)

	result2 := newParsingResult(ref, 0, "March 8, 2020 at 8:00 AM", components2, nil)
	results2 := []*parsingResult{result2}

	converted2, err := pipeline.applyTimezoneConversion(results2)
	assert.NoError(t, err)

	start2 := converted2[0].Start()
	assert.Equal(t, 4, *start2.Get(ComponentHour), "Should be 4:00 AM EDT (UTC-4)")
}

// TestPipeline_TimezoneConversion_RangeResults tests timezone conversion for date ranges
func TestPipeline_TimezoneConversion_RangeResults(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	// Create a range result
	startComponents := newParsingComponents(ref, nil)
	startComponents.Assign(ComponentYear, 2020)
	startComponents.Assign(ComponentMonth, 6)
	startComponents.Assign(ComponentDay, 15)
	startComponents.Assign(ComponentHour, 18)
	startComponents.Assign(ComponentMinute, 0)
	startComponents.Assign(ComponentSecond, 0)

	endComponents := newParsingComponents(ref, nil)
	endComponents.Assign(ComponentYear, 2020)
	endComponents.Assign(ComponentMonth, 6)
	endComponents.Assign(ComponentDay, 15)
	endComponents.Assign(ComponentHour, 22)
	endComponents.Assign(ComponentMinute, 0)
	endComponents.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "June 15, 2020 from 6:00 PM to 10:00 PM", startComponents, endComponents)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)
	assert.Len(t, converted, 1)

	// Verify start is converted
	start := converted[0].Start()
	assert.Equal(t, 14, *start.Get(ComponentHour), "Start: 18:00 UTC = 14:00 EDT")

	// Verify end is converted
	end := converted[0].End()
	assert.NotNil(t, end)
	assert.Equal(t, 18, *end.Get(ComponentHour), "End: 22:00 UTC = 18:00 EDT")
}

// TestPipeline_TimezoneConversion_Midnight tests midnight handling
func TestPipeline_TimezoneConversion_midnight(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/Los_Angeles"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	// UTC midnight
	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 16)
	components.Assign(ComponentHour, 0)
	components.Assign(ComponentMinute, 0)
	components.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "June 16, 2020 at midnight", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	// 00:00 UTC on June 16 = 17:00 PDT on June 15 (UTC-7 in summer)
	assert.Equal(t, 2020, *start.Get(ComponentYear))
	assert.Equal(t, 6, *start.Get(ComponentMonth))
	assert.Equal(t, 15, *start.Get(ComponentDay), "Should be previous day")
	assert.Equal(t, 17, *start.Get(ComponentHour))
}

// TestPipeline_TimezoneConversion_Noon tests noon handling
func TestPipeline_TimezoneConversion_noon(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "Asia/Shanghai"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	// UTC noon
	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 12)
	components.Assign(ComponentMinute, 0)
	components.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "June 15, 2020 at noon", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	// 12:00 UTC = 20:00 CST (UTC+8)
	assert.Equal(t, 2020, *start.Get(ComponentYear))
	assert.Equal(t, 6, *start.Get(ComponentMonth))
	assert.Equal(t, 15, *start.Get(ComponentDay))
	assert.Equal(t, 20, *start.Get(ComponentHour))
}

// TestPipeline_TimezoneConversion_YearBoundary tests year boundary changes
func TestPipeline_TimezoneConversion_YearBoundary(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "Pacific/Auckland"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2019, 12, 31, 12, 0, 0, 0, utc), nil)

	// UTC time on Dec 31, 2019 at 11:00 PM
	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2019)
	components.Assign(ComponentMonth, 12)
	components.Assign(ComponentDay, 31)
	components.Assign(ComponentHour, 23)
	components.Assign(ComponentMinute, 0)
	components.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "December 31, 2019 at 11:00 PM", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	// 23:00 UTC on Dec 31, 2019 = 12:00 NZDT on Jan 1, 2020 (UTC+13 in summer)
	assert.Equal(t, 2020, *start.Get(ComponentYear), "Year should advance to 2020")
	assert.Equal(t, 1, *start.Get(ComponentMonth), "Month should be January")
	assert.Equal(t, 1, *start.Get(ComponentDay), "Day should be 1st")
	assert.Equal(t, 12, *start.Get(ComponentHour))
}

// TestPipeline_TimezoneConversion_YearBoundaryBackward tests backward year boundary
func TestPipeline_TimezoneConversion_YearBoundaryBackward(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/Los_Angeles"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 1, 1, 12, 0, 0, 0, utc), nil)

	// UTC time on Jan 1, 2020 at 02:00 AM
	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 1)
	components.Assign(ComponentDay, 1)
	components.Assign(ComponentHour, 2)
	components.Assign(ComponentMinute, 0)
	components.Assign(ComponentSecond, 0)

	result := newParsingResult(ref, 0, "January 1, 2020 at 2:00 AM", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	// 02:00 UTC on Jan 1, 2020 = 18:00 PST on Dec 31, 2019 (UTC-8 in winter)
	assert.Equal(t, 2019, *start.Get(ComponentYear), "Year should go back to 2019")
	assert.Equal(t, 12, *start.Get(ComponentMonth), "Month should be December")
	assert.Equal(t, 31, *start.Get(ComponentDay), "Day should be 31st")
	assert.Equal(t, 18, *start.Get(ComponentHour))
}

// TestPipeline_TimezoneConversion_PreserveCertainty tests that certain/implied status is preserved
func TestPipeline_TimezoneConversion_PreserveCertainty(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	// Only hour is certain (explicitly parsed), everything else is implied
	components.Assign(ComponentHour, 18)

	result := newParsingResult(ref, 0, "6:00 PM", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	// Hour should still be certain
	assert.True(t, start.IsCertain(ComponentHour))
	// But other components should remain implied
	assert.False(t, start.IsCertain(ComponentYear))
	assert.False(t, start.IsCertain(ComponentMonth))
	assert.False(t, start.IsCertain(ComponentDay))
}

// TestPipeline_TimezoneConversion_InvalidTimezone tests error handling for invalid timezone
func TestPipeline_TimezoneConversion_InvalidTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "Invalid/Timezone"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 18)

	result := newParsingResult(ref, 0, "June 15, 2020 at 6:00 PM", components, nil)
	results := []*parsingResult{result}

	_, err := pipeline.applyTimezoneConversion(results)
	assert.Error(t, err, "Should return error for invalid timezone")
}

// TestPipeline_TimezoneConversion_EmptyResults tests handling of empty results
func TestPipeline_TimezoneConversion_EmptyResults(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	results := []*parsingResult{}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)
	assert.Len(t, converted, 0)
}

// TestPipeline_TimezoneConversion_NilStartComponents tests handling of nil start components
func TestPipeline_TimezoneConversion_NilStartComponents(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	result := newParsingResult(ref, 0, "test", nil, nil)
	results := []*parsingResult{result}

	// Should not panic with nil start components
	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)
	assert.Len(t, converted, 1)
}

// TestPipeline_TimezoneConversion_WithMilliseconds tests subsecond precision
func TestPipeline_TimezoneConversion_WithMilliseconds(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 18)
	components.Assign(ComponentMinute, 30)
	components.Assign(ComponentSecond, 45)
	components.Assign(ComponentMillisecond, 123)
	components.Assign(ComponentMicrosecond, 456)
	components.Assign(ComponentNanosecond, 789)

	result := newParsingResult(ref, 0, "June 15, 2020 at 6:30:45.123 PM", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	assert.Equal(t, 14, *start.Get(ComponentHour))
	assert.Equal(t, 30, *start.Get(ComponentMinute))
	assert.Equal(t, 45, *start.Get(ComponentSecond))
	assert.Equal(t, 123, *start.Get(ComponentMillisecond), "Milliseconds should be preserved")
}

// TestPipeline_TimezoneConversion_MeridiemUpdate tests that meridiem is updated correctly
func TestPipeline_TimezoneConversion_MeridiemUpdate(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	// 11:00 AM UTC should become 7:00 AM EDT (AM to AM)
	components1 := newParsingComponents(ref, nil)
	components1.Assign(ComponentYear, 2020)
	components1.Assign(ComponentMonth, 6)
	components1.Assign(ComponentDay, 15)
	components1.Assign(ComponentHour, 11)
	components1.Assign(ComponentMeridiem, int(0))

	result1 := newParsingResult(ref, 0, "June 15, 2020 at 11:00 AM", components1, nil)
	results1 := []*parsingResult{result1}

	converted1, err := pipeline.applyTimezoneConversion(results1)
	assert.NoError(t, err)

	start1 := converted1[0].Start()
	assert.Equal(t, 7, *start1.Get(ComponentHour))
	assert.Equal(t, int(0), *start1.Get(ComponentMeridiem))

	// 1:00 PM UTC (13:00) should become 9:00 AM EDT (PM to AM - crosses meridiem)
	components2 := newParsingComponents(ref, nil)
	components2.Assign(ComponentYear, 2020)
	components2.Assign(ComponentMonth, 6)
	components2.Assign(ComponentDay, 15)
	components2.Assign(ComponentHour, 13)
	components2.Assign(ComponentMeridiem, int(1))

	result2 := newParsingResult(ref, 0, "June 15, 2020 at 1:00 PM", components2, nil)
	results2 := []*parsingResult{result2}

	converted2, err := pipeline.applyTimezoneConversion(results2)
	assert.NoError(t, err)

	start2 := converted2[0].Start()
	assert.Equal(t, 9, *start2.Get(ComponentHour))
	assert.Equal(t, int(0), *start2.Get(ComponentMeridiem), "Meridiem should change from PM to AM")
}

// TestPipeline_TimezoneConversion_TimezoneOffset tests that timezone offset is updated
func TestPipeline_TimezoneConversion_TimezoneOffset(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := newPipeline(nil, settings)

	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 18)
	components.Assign(ComponentTimezoneOffset, 0) // UTC offset

	result := newParsingResult(ref, 0, "June 15, 2020 at 6:00 PM UTC", components, nil)
	results := []*parsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	assert.NoError(t, err)

	start := converted[0].Start()
	// EDT is UTC-4, which is -240 minutes
	assert.Equal(t, -240, *start.Get(ComponentTimezoneOffset), "Timezone offset should be updated to EDT (-240 minutes)")
}

// TestUpdateComponentsFromDate_AllComponents tests the helper function with all components
func TestUpdateComponentsFromDate_AllComponents(t *testing.T) {
	utc := time.UTC
	ref := newReferenceWithTimezone(time.Date(2020, 6, 15, 12, 0, 0, 0, utc), nil)

	components := newParsingComponents(ref, nil)
	components.Assign(ComponentYear, 2020)
	components.Assign(ComponentMonth, 6)
	components.Assign(ComponentDay, 15)
	components.Assign(ComponentHour, 18)
	components.Assign(ComponentMinute, 30)
	components.Assign(ComponentSecond, 45)
	components.Assign(ComponentMillisecond, 123)
	components.Assign(ComponentMicrosecond, 456)
	components.Assign(ComponentNanosecond, 789)

	// Create a new date in a different timezone
	nyLoc, _ := time.LoadLocation("America/New_York")
	newDate := time.Date(2020, 6, 15, 14, 30, 45, 123000000, nyLoc)

	updateComponentsFromDate(components, newDate)

	// Verify all components are updated
	assert.Equal(t, 2020, *components.Get(ComponentYear))
	assert.Equal(t, 6, *components.Get(ComponentMonth))
	assert.Equal(t, 15, *components.Get(ComponentDay))
	assert.Equal(t, 14, *components.Get(ComponentHour))
	assert.Equal(t, 30, *components.Get(ComponentMinute))
	assert.Equal(t, 45, *components.Get(ComponentSecond))
	assert.Equal(t, 123, *components.Get(ComponentMillisecond))
	assert.Equal(t, 0, *components.Get(ComponentMicrosecond))
	assert.Equal(t, 0, *components.Get(ComponentNanosecond))
}
