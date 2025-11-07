package helpers

import (
	"time"

	"github.com/kljensen/kronos"
)

// ToTimezoneOffset converts various timezone representations to an offset in minutes.
// It supports:
// - Numeric offsets (int) - returned as-is
// - Timezone names (string) - looked up in overrides map or DefaultTimezoneAbbrMap
// - Ambiguous timezones with DST - determined based on the instant time
//
// Returns nil if the timezone cannot be resolved.
func ToTimezoneOffset(tz interface{}, instant time.Time, overrides kronos.TimezoneAbbrMap, defaultMap kronos.TimezoneAbbrMap) *int {
	if tz == nil {
		return nil
	}

	// If it's already a numeric offset, return it
	if offset, ok := tz.(int); ok {
		return &offset
	}

	// If it's a string, look it up in the maps
	if name, ok := tz.(string); ok {
		// First check overrides
		if overrides != nil {
			if val, exists := overrides[name]; exists {
				return resolveTimezoneValue(val, instant)
			}
		}

		// Then check default map
		if val, exists := defaultMap[name]; exists {
			return resolveTimezoneValue(val, instant)
		}

		// Try to load as a location name (e.g., "America/New_York")
		loc, err := time.LoadLocation(name)
		if err == nil {
			_, offset := instant.In(loc).Zone()
			offsetMinutes := offset / 60
			return &offsetMinutes
		}
	}

	return nil
}

// FromInput creates a ReferenceWithTimezone from either a ParsingReference or a time.Time.
// It also handles timezone conversion using the provided timezoneOverrides.
func FromInput(input interface{}, timezoneOverrides kronos.TimezoneAbbrMap, defaultMap kronos.TimezoneAbbrMap, newRefFn func(time.Time, *int) *kronos.InternalReferenceWithTimezone) *kronos.InternalReferenceWithTimezone {
	if input == nil {
		return newRefFn(time.Time{}, nil)
	}

	switch v := input.(type) {
	case time.Time:
		return newRefFn(v, nil)
	case kronos.InternalParsingReference:
		instant := time.Now()
		if v.Instant != nil {
			instant = *v.Instant
		}

		var timezoneOffset *int
		if v.Timezone != nil {
			timezoneOffset = ToTimezoneOffset(v.Timezone, instant, timezoneOverrides, defaultMap)
		}

		return newRefFn(instant, timezoneOffset)
	default:
		return newRefFn(time.Time{}, nil)
	}
}

// Helper functions

// resolveTimezoneValue resolves a timezone value from the map to an offset.
// Handles both simple integer offsets and AmbiguousTimezoneMap entries.
func resolveTimezoneValue(val interface{}, instant time.Time) *int {
	// Simple integer offset
	if offset, ok := val.(int); ok {
		return &offset
	}

	// Ambiguous timezone with DST (value type)
	if ambiguous, ok := val.(kronos.AmbiguousTimezoneMap); ok {
		return resolveAmbiguousTimezone(ambiguous, instant)
	}

	// Ambiguous timezone with DST (pointer type)
	if ambiguous, ok := val.(*kronos.AmbiguousTimezoneMap); ok {
		return resolveAmbiguousTimezone(*ambiguous, instant)
	}

	return nil
}

// resolveAmbiguousTimezone resolves an ambiguous timezone to its offset based on the instant.
func resolveAmbiguousTimezone(ambiguous kronos.AmbiguousTimezoneMap, instant time.Time) *int {
	// Without a valid instant, we can't determine DST status
	if instant.IsZero() {
		return nil
	}

	year := instant.Year()
	dstStart := ambiguous.DstStart(year)
	dstEnd := ambiguous.DstEnd(year)

	// Check if instant is during DST period
	if instant.After(dstStart) && !instant.After(dstEnd) {
		offset := ambiguous.TimezoneOffsetDuringDst
		return &offset
	}

	// Not during DST
	offset := ambiguous.TimezoneOffsetNonDst
	return &offset
}
