package kronos

import (
	"testing"
	"time"
)

// TestComponentConstants verifies all Component constants are defined.
func TestComponentConstants(t *testing.T) {
	components := []Component{
		ComponentYear,
		ComponentMonth,
		ComponentDay,
		ComponentWeekday,
		ComponentHour,
		ComponentMinute,
		ComponentSecond,
		ComponentMillisecond,
		ComponentMeridiem,
		ComponentTimezoneOffset,
	}

	expected := []string{
		"year",
		"month",
		"day",
		"weekday",
		"hour",
		"minute",
		"second",
		"millisecond",
		"meridiem",
		"timezoneOffset",
	}

	if len(components) != len(expected) {
		t.Errorf("Expected %d components, got %d", len(expected), len(components))
	}

	for i, component := range components {
		if string(component) != expected[i] {
			t.Errorf("Component[%d]: expected %q, got %q", i, expected[i], string(component))
		}
	}
}

// TestTimeunitConstants verifies all Timeunit constants are defined.
func TestTimeunitConstants(t *testing.T) {
	timeunits := []Timeunit{
		TimeunitYear,
		TimeunitMonth,
		TimeunitWeek,
		TimeunitDay,
		TimeunitHour,
		TimeunitMinute,
		TimeunitSecond,
		TimeunitMillisecond,
		TimeunitQuarter,
	}

	expected := []string{
		"year",
		"month",
		"week",
		"day",
		"hour",
		"minute",
		"second",
		"millisecond",
		"quarter",
	}

	if len(timeunits) != len(expected) {
		t.Errorf("Expected %d timeunits, got %d", len(expected), len(timeunits))
	}

	for i, timeunit := range timeunits {
		if string(timeunit) != expected[i] {
			t.Errorf("Timeunit[%d]: expected %q, got %q", i, expected[i], string(timeunit))
		}
	}
}

// NOTE: Tests for Meridiem, Weekday, and Month constants have been removed.
// These types are now internal only. Users should use time.Weekday and time.Month
// from the standard library instead.

// NOTE: parsingOption tests removed as it's an internal implementation detail.
// Settings tests are in settings_test.go

// TestAmbiguousTimezoneMap verifies AmbiguousTimezoneMap structure.
func TestAmbiguousTimezoneMap(t *testing.T) {
	atz := AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -240,
		TimezoneOffsetNonDst:    -300,
		DstStart: func(year int) time.Time {
			return time.Date(year, 3, 10, 2, 0, 0, 0, time.UTC)
		},
		DstEnd: func(year int) time.Time {
			return time.Date(year, 11, 3, 2, 0, 0, 0, time.UTC)
		},
	}

	if atz.TimezoneOffsetDuringDst != -240 {
		t.Errorf("Expected DST offset -240, got %d", atz.TimezoneOffsetDuringDst)
	}

	if atz.TimezoneOffsetNonDst != -300 {
		t.Errorf("Expected non-DST offset -300, got %d", atz.TimezoneOffsetNonDst)
	}

	dstStart := atz.DstStart(2024)
	expectedStart := time.Date(2024, 3, 10, 2, 0, 0, 0, time.UTC)
	if !dstStart.Equal(expectedStart) {
		t.Errorf("Expected DST start %v, got %v", expectedStart, dstStart)
	}

	dstEnd := atz.DstEnd(2024)
	expectedEnd := time.Date(2024, 11, 3, 2, 0, 0, 0, time.UTC)
	if !dstEnd.Equal(expectedEnd) {
		t.Errorf("Expected DST end %v, got %v", expectedEnd, dstEnd)
	}
}

// TestTimezoneAbbrMap verifies TimezoneAbbrMap can store both simple and ambiguous timezones.
func TestTimezoneAbbrMap(t *testing.T) {
	tzMap := make(TimezoneAbbrMap)

	// Simple timezone
	tzMap["UTC"] = 0
	tzMap["EST"] = -300

	// Ambiguous timezone
	tzMap["EDT"] = AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -240,
		TimezoneOffsetNonDst:    -300,
		DstStart: func(year int) time.Time {
			return time.Date(year, 3, 10, 2, 0, 0, 0, time.UTC)
		},
		DstEnd: func(year int) time.Time {
			return time.Date(year, 11, 3, 2, 0, 0, 0, time.UTC)
		},
	}

	if offset, ok := tzMap["UTC"].(int); !ok || offset != 0 {
		t.Errorf("Expected UTC offset 0, got %v", tzMap["UTC"])
	}

	if offset, ok := tzMap["EST"].(int); !ok || offset != -300 {
		t.Errorf("Expected EST offset -300, got %v", tzMap["EST"])
	}

	if atz, ok := tzMap["EDT"].(AmbiguousTimezoneMap); !ok {
		t.Errorf("Expected EDT to be AmbiguousTimezoneMap, got %T", tzMap["EDT"])
	} else if atz.TimezoneOffsetDuringDst != -240 {
		t.Errorf("Expected EDT DST offset -240, got %d", atz.TimezoneOffsetDuringDst)
	}
}

// TestparsingReference verifies parsingReference structure.
func TestParsingReference(t *testing.T) {
	now := time.Now()
	ref := parsingReference{
		Instant:  &now,
		Timezone: "America/New_York",
	}

	if ref.Instant == nil {
		t.Fatal("Instant should not be nil")
	}

	if !ref.Instant.Equal(now) {
		t.Errorf("Expected instant %v, got %v", now, ref.Instant)
	}

	if tz, ok := ref.Timezone.(string); !ok || tz != "America/New_York" {
		t.Errorf("Expected timezone 'America/New_York', got %v", ref.Timezone)
	}
}

// TestparsingReferenceWithOffset verifies parsingReference can use offset.
func TestParsingReferenceWithOffset(t *testing.T) {
	now := time.Now()
	ref := parsingReference{
		Timezone: -300, // EST offset in minutes
	}
	_ = now // Reference time not needed for this test

	if offset, ok := ref.Timezone.(int); !ok || offset != -300 {
		t.Errorf("Expected timezone offset -300, got %v", ref.Timezone)
	}
}

// TestComponentAsMapKey verifies Component can be used as a map key.
func TestComponentAsMapKey(t *testing.T) {
	components := make(map[Component]int)
	components[ComponentYear] = 2024
	components[ComponentMonth] = 11
	components[ComponentDay] = 2

	if components[ComponentYear] != 2024 {
		t.Errorf("Expected year 2024, got %d", components[ComponentYear])
	}

	if components[ComponentMonth] != 11 {
		t.Errorf("Expected month 11, got %d", components[ComponentMonth])
	}

	if components[ComponentDay] != 2 {
		t.Errorf("Expected day 2, got %d", components[ComponentDay])
	}
}

// TestTimeunitAsMapKey verifies Timeunit can be used as a map key.
func TestTimeunitAsMapKey(t *testing.T) {
	offsets := make(map[Timeunit]int)
	offsets[TimeunitDay] = 1
	offsets[TimeunitHour] = 2
	offsets[TimeunitMinute] = 30

	if offsets[TimeunitDay] != 1 {
		t.Errorf("Expected day offset 1, got %d", offsets[TimeunitDay])
	}

	if offsets[TimeunitHour] != 2 {
		t.Errorf("Expected hour offset 2, got %d", offsets[TimeunitHour])
	}

	if offsets[TimeunitMinute] != 30 {
		t.Errorf("Expected minute offset 30, got %d", offsets[TimeunitMinute])
	}
}
