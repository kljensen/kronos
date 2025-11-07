package kronos

import (
	"testing"
	"time"
)

// TestReferenceWithTimezone tests basic ReferenceWithTimezone functionality
func TestReferenceWithTimezone(t *testing.T) {
	t.Run("NewReferenceWithTimezone with instant and offset", func(t *testing.T) {
		instant := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
		offset := -300 // EST offset
		ref := newReferenceWithTimezone(instant, &offset)

		if !ref.Instant().Equal(instant) {
			t.Errorf("Expected instant %v, got %v", instant, ref.Instant())
		}

		if ref.timezoneOffset == nil || *ref.timezoneOffset != offset {
			t.Errorf("Expected offset %d, got %v", offset, ref.timezoneOffset)
		}
	})

	t.Run("NewReferenceWithTimezone with zero instant uses current time", func(t *testing.T) {
		before := time.Now()
		ref := newReferenceWithTimezone(time.Time{}, nil)
		after := time.Now()

		if ref.Instant().Before(before) || ref.Instant().After(after) {
			t.Errorf("Expected instant to be current time, got %v", ref.Instant())
		}
	})

	t.Run("newReferenceWithTimezone creates reference without timezone offset", func(t *testing.T) {
		date := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
		ref := newReferenceWithTimezone(date, nil)

		if !ref.Instant().Equal(date) {
			t.Errorf("Expected instant %v, got %v", date, ref.Instant())
		}

		if ref.timezoneOffset != nil {
			t.Errorf("Expected nil timezone offset, got %v", ref.timezoneOffset)
		}
	})
}

func TestFromInput(t *testing.T) {
	t.Run("FromInput with time.Time", func(t *testing.T) {
		instant := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
		ref := fromInput(instant, nil)

		if !ref.Instant().Equal(instant) {
			t.Errorf("Expected instant %v, got %v", instant, ref.Instant())
		}

		if ref.timezoneOffset != nil {
			t.Errorf("Expected nil timezone offset, got %v", ref.timezoneOffset)
		}
	})

	t.Run("FromInput with parsingReference", func(t *testing.T) {
		instant := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
		offset := -300
		parsingRef := parsingReference{
			Instant:  &instant,
			Timezone: offset,
		}

		ref := fromInput(parsingRef, nil)

		if !ref.Instant().Equal(instant) {
			t.Errorf("Expected instant %v, got %v", instant, ref.Instant())
		}

		if ref.timezoneOffset == nil || *ref.timezoneOffset != offset {
			t.Errorf("Expected offset %d, got %v", offset, ref.timezoneOffset)
		}
	})

	t.Run("FromInput with nil", func(t *testing.T) {
		before := time.Now()
		ref := fromInput(nil, nil)
		after := time.Now()

		if ref.Instant().Before(before) || ref.Instant().After(after) {
			t.Errorf("Expected instant to be current time, got %v", ref.Instant())
		}
	})
}

func TestGetTimezoneOffset(t *testing.T) {
	t.Run("GetTimezoneOffset with explicit offset", func(t *testing.T) {
		offset := -300
		ref := newReferenceWithTimezone(time.Now(), &offset)

		if ref.GetTimezoneOffset() != offset {
			t.Errorf("Expected offset %d, got %d", offset, ref.GetTimezoneOffset())
		}
	})

	t.Run("GetTimezoneOffset without explicit offset uses system", func(t *testing.T) {
		instant := time.Now()
		ref := newReferenceWithTimezone(instant, nil)

		_, systemOffset := instant.Zone()
		expectedOffset := systemOffset / 60

		if ref.GetTimezoneOffset() != expectedOffset {
			t.Errorf("Expected system offset %d, got %d", expectedOffset, ref.GetTimezoneOffset())
		}
	})
}

func TestParsingComponents(t *testing.T) {
	ref := newReferenceWithTimezone(time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC), nil)

	t.Run("NewParsingComponents sets default implied values", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)

		// Check implied values from reference
		if pc.Get(ComponentYear) == nil || *pc.Get(ComponentYear) != 2024 {
			t.Errorf("Expected year 2024 to be implied")
		}

		if pc.Get(ComponentMonth) == nil || *pc.Get(ComponentMonth) != 11 {
			t.Errorf("Expected month 11 to be implied")
		}

		if pc.Get(ComponentDay) == nil || *pc.Get(ComponentDay) != 2 {
			t.Errorf("Expected day 2 to be implied")
		}

		if pc.Get(ComponentHour) == nil || *pc.Get(ComponentHour) != 12 {
			t.Errorf("Expected hour 12 to be implied")
		}

		if pc.Get(ComponentMinute) == nil || *pc.Get(ComponentMinute) != 0 {
			t.Errorf("Expected minute 0 to be implied")
		}

		if pc.Get(ComponentSecond) == nil || *pc.Get(ComponentSecond) != 0 {
			t.Errorf("Expected second 0 to be implied")
		}
	})

	t.Run("IsCertain returns true only for known values", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)

		if pc.IsCertain(ComponentYear) {
			t.Errorf("Year should not be certain initially")
		}

		pc.Assign(ComponentYear, 2025)

		if !pc.IsCertain(ComponentYear) {
			t.Errorf("Year should be certain after assignment")
		}
	})

	t.Run("Assign sets known value and removes implied", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)

		// Initially implied
		if pc.IsCertain(ComponentYear) {
			t.Errorf("Year should not be certain initially")
		}

		pc.Assign(ComponentYear, 2025)

		if !pc.IsCertain(ComponentYear) {
			t.Errorf("Year should be certain after assignment")
		}

		if val := pc.Get(ComponentYear); val == nil || *val != 2025 {
			t.Errorf("Expected year 2025, got %v", val)
		}
	})

	t.Run("Imply does not override known values", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)

		pc.Assign(ComponentYear, 2025)
		pc.Imply(ComponentYear, 2026)

		if val := pc.Get(ComponentYear); val == nil || *val != 2025 {
			t.Errorf("Expected year 2025 (known value should not be overridden), got %v", val)
		}
	})

	t.Run("Delete removes both known and implied values", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)

		pc.Assign(ComponentYear, 2025)
		pc.Assign(ComponentMonth, 12)

		pc.Delete(ComponentYear, ComponentMonth)

		if pc.Get(ComponentYear) != nil {
			t.Errorf("Year should be deleted")
		}

		if pc.Get(ComponentMonth) != nil {
			t.Errorf("Month should be deleted")
		}
	})

	t.Run("Clone creates independent copy", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentYear, 2025)
		pc.AddTag("test-tag")

		clone := pc.Clone()

		// Modify original
		pc.Assign(ComponentYear, 2026)
		pc.AddTag("another-tag")

		// Clone should not be affected
		if val := clone.Get(ComponentYear); val == nil || *val != 2025 {
			t.Errorf("Expected cloned year 2025, got %v", val)
		}

		tags := clone.Tags()
		if !tags["test-tag"] {
			t.Errorf("Expected test-tag in clone")
		}

		if tags["another-tag"] {
			t.Errorf("Clone should not have another-tag")
		}
	})

	t.Run("IsOnlyDate returns true when no time components are certain", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentYear, 2025)
		pc.Assign(ComponentMonth, 12)
		pc.Assign(ComponentDay, 25)

		if !pc.IsOnlyDate() {
			t.Errorf("Expected IsOnlyDate to be true")
		}

		pc.Assign(ComponentHour, 14)

		if pc.IsOnlyDate() {
			t.Errorf("Expected IsOnlyDate to be false after setting hour")
		}
	})

	t.Run("IsOnlyTime returns true when no date components are certain", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentHour, 14)
		pc.Assign(ComponentMinute, 30)

		if !pc.IsOnlyTime() {
			t.Errorf("Expected IsOnlyTime to be true")
		}

		pc.Assign(ComponentDay, 25)

		if pc.IsOnlyTime() {
			t.Errorf("Expected IsOnlyTime to be false after setting day")
		}
	})

	t.Run("IsOnlyWeekdayComponent", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentWeekday, int(time.Monday))

		if !pc.IsOnlyWeekdayComponent() {
			t.Errorf("Expected IsOnlyWeekdayComponent to be true")
		}

		pc.Assign(ComponentDay, 2)

		if pc.IsOnlyWeekdayComponent() {
			t.Errorf("Expected IsOnlyWeekdayComponent to be false after setting day")
		}
	})

	t.Run("IsDateWithUnknownYear", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentMonth, 12)

		if !pc.IsDateWithUnknownYear() {
			t.Errorf("Expected IsDateWithUnknownYear to be true")
		}

		pc.Assign(ComponentYear, 2025)

		if pc.IsDateWithUnknownYear() {
			t.Errorf("Expected IsDateWithUnknownYear to be false after setting year")
		}
	})

	t.Run("IsValidDate validates component values", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentYear, 2024)
		pc.Assign(ComponentMonth, 2)
		pc.Assign(ComponentDay, 29) // Valid in leap year

		if !pc.IsValidDate() {
			t.Errorf("Expected date to be valid")
		}

		pc.Assign(ComponentDay, 30) // Invalid for February

		if pc.IsValidDate() {
			t.Errorf("Expected date to be invalid")
		}
	})

	t.Run("Date constructs time.Time from components", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentYear, 2024)
		pc.Assign(ComponentMonth, 11)
		pc.Assign(ComponentDay, 2)
		pc.Assign(ComponentHour, 14)
		pc.Assign(ComponentMinute, 30)
		pc.Assign(ComponentSecond, 45)

		date := pc.Date()

		if date.Year() != 2024 {
			t.Errorf("Expected year 2024, got %d", date.Year())
		}

		if date.Month() != 11 {
			t.Errorf("Expected month 11, got %d", date.Month())
		}

		if date.Day() != 2 {
			t.Errorf("Expected day 2, got %d", date.Day())
		}

		if date.Hour() != 14 {
			t.Errorf("Expected hour 14, got %d", date.Hour())
		}

		if date.Minute() != 30 {
			t.Errorf("Expected minute 30, got %d", date.Minute())
		}

		if date.Second() != 45 {
			t.Errorf("Expected second 45, got %d", date.Second())
		}
	})

	t.Run("AddTag and Tags", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)

		pc.AddTag("tag1")
		pc.AddTag("tag2")

		tags := pc.Tags()

		if !tags["tag1"] {
			t.Errorf("Expected tag1 to be present")
		}

		if !tags["tag2"] {
			t.Errorf("Expected tag2 to be present")
		}

		if len(tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(tags))
		}
	})

	t.Run("String returns debug representation", func(t *testing.T) {
		pc := newParsingComponents(ref, nil)
		pc.Assign(ComponentYear, 2024)
		pc.AddTag("test")

		str := pc.String()

		if str == "" {
			t.Errorf("Expected non-empty string representation")
		}
	})
}

func TestParsingResult(t *testing.T) {
	ref := newReferenceWithTimezone(time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC), nil)

	t.Run("NewParsingResult creates result", func(t *testing.T) {
		start := newParsingComponents(ref, nil)
		start.Assign(ComponentYear, 2024)

		result := newParsingResult(ref, 0, "Nov 2, 2024", start, nil)

		if result.Index() != 0 {
			t.Errorf("Expected index 0, got %d", result.Index())
		}

		if result.Text() != "Nov 2, 2024" {
			t.Errorf("Expected text 'Nov 2, 2024', got %s", result.Text())
		}

		if result.Start() == nil {
			t.Errorf("Expected start to be non-nil")
		}

		if result.End() != nil {
			t.Errorf("Expected end to be nil")
		}
	})

	t.Run("NewParsingResult with nil start creates default", func(t *testing.T) {
		result := newParsingResult(ref, 5, "test", nil, nil)

		if result.Start() == nil {
			t.Errorf("Expected start to be created by default")
		}
	})

	t.Run("Clone creates independent copy", func(t *testing.T) {
		start := newParsingComponents(ref, nil)
		start.Assign(ComponentYear, 2024)
		start.AddTag("start-tag")

		result := newParsingResult(ref, 0, "test", start, nil)
		result.AddTag("result-tag")

		clone := result.Clone()

		// Modify original
		result.Start().(*parsingComponents).Assign(ComponentYear, 2025)
		result.AddTag("new-tag")

		// Clone should not be affected
		clonedStart := clone.Start().(*parsingComponents)
		if val := clonedStart.Get(ComponentYear); val == nil || *val != 2024 {
			t.Errorf("Expected cloned year 2024, got %v", val)
		}

		tags := clone.Tags()
		if tags["new-tag"] {
			t.Errorf("Clone should not have new-tag")
		}
	})

	t.Run("Date delegates to start.Date()", func(t *testing.T) {
		start := newParsingComponents(ref, nil)
		start.Assign(ComponentYear, 2024)
		start.Assign(ComponentMonth, 11)
		start.Assign(ComponentDay, 2)

		result := newParsingResult(ref, 0, "test", start, nil)
		date := result.Date()

		if date.Year() != 2024 || date.Month() != 11 || date.Day() != 2 {
			t.Errorf("Expected date 2024-11-02, got %v", date)
		}
	})

	t.Run("AddTag adds to both start and end", func(t *testing.T) {
		start := newParsingComponents(ref, nil)
		end := newParsingComponents(ref, nil)

		result := newParsingResult(ref, 0, "test", start, end)
		result.AddTag("test-tag")

		if !start.Tags()["test-tag"] {
			t.Errorf("Expected test-tag in start")
		}

		if !end.Tags()["test-tag"] {
			t.Errorf("Expected test-tag in end")
		}
	})

	t.Run("Tags combines start and end tags", func(t *testing.T) {
		start := newParsingComponents(ref, nil)
		start.AddTag("start-tag")

		end := newParsingComponents(ref, nil)
		end.AddTag("end-tag")

		result := newParsingResult(ref, 0, "test", start, end)

		tags := result.Tags()

		if !tags["start-tag"] {
			t.Errorf("Expected start-tag in combined tags")
		}

		if !tags["end-tag"] {
			t.Errorf("Expected end-tag in combined tags")
		}

		if len(tags) != 2 {
			t.Errorf("Expected 2 combined tags, got %d", len(tags))
		}
	})

	t.Run("String returns debug representation", func(t *testing.T) {
		start := newParsingComponents(ref, nil)
		result := newParsingResult(ref, 5, "Nov 2", start, nil)
		result.AddTag("test")

		str := result.String()

		if str == "" {
			t.Errorf("Expected non-empty string representation")
		}
	})

	t.Run("RefDate returns reference date", func(t *testing.T) {
		result := newParsingResult(ref, 0, "test", nil, nil)

		if !result.RefDate().Equal(ref.Instant()) {
			t.Errorf("Expected RefDate to equal reference instant")
		}
	})
}
