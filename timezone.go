package kronos

import (
	"time"

	"github.com/kljensen/kronos/internal/data"
)

// AmbiguousTimezoneMap defines a timezone that has different offsets
// depending on whether daylight saving time (DST) is in effect.
//
// Deprecated: Direct use of AmbiguousTimezoneMap is discouraged. Use the builder pattern
// with WithOption instead:
//
//	customTimezones := kronos.TimezoneAbbrMap{
//	    "ET": &kronos.AmbiguousTimezoneMap{
//	        TimezoneOffsetDuringDst: -240,
//	        TimezoneOffsetNonDst: -300,
//	        DstStart: func(year int) time.Time { ... },
//	        DstEnd: func(year int) time.Time { ... },
//	    },
//	}
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.TimezoneOverrides = customTimezones
//	    })
type AmbiguousTimezoneMap = data.AmbiguousTimezoneMap

// TimezoneAbbrMap maps timezone abbreviations to their offsets.
// Values can be either a simple offset (in minutes) or an AmbiguousTimezoneMap
// for timezones that observe DST.
//
// Deprecated: Direct use of TimezoneAbbrMap is discouraged. Use the builder pattern
// with WithOption instead:
//
//	customTimezones := kronos.TimezoneAbbrMap{
//	    "CUSTOM": 123,  // UTC+2:03
//	    "TEST": -456,   // UTC-7:36
//	}
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.TimezoneOverrides = customTimezones
//	    })
type TimezoneAbbrMap = data.TimezoneAbbrMap

// DebugHandler is a function that handles debug events.
// It receives a debug message for logging or analysis.
//
// Deprecated: Direct use of DebugHandler is discouraged. Use the builder pattern
// with WithOption instead:
//
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.DebugHandler = func(msg string) {
//	            log.Printf("DEBUG: %s", msg)
//	        }
//	    })
type DebugHandler func(message string)

// ToTimezoneOffset converts various timezone representations to an offset in minutes.
// It supports:
// - Numeric offsets (int) - returned as-is
// - Timezone names (string) - looked up in overrides map or default timezone map
// - Ambiguous timezones with DST - determined based on the instant time
//
// Returns nil if the timezone cannot be resolved.
func toTimezoneOffset(tz any, instant time.Time, overrides TimezoneAbbrMap) *int {
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
		if val, exists := data.DefaultTimezoneAbbrMap[name]; exists {
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

// resolveTimezoneValue resolves a timezone value from the map to an offset.
// Handles both simple integer offsets and AmbiguousTimezoneMap entries.
func resolveTimezoneValue(val any, instant time.Time) *int {
	// Simple integer offset
	if offset, ok := val.(int); ok {
		return &offset
	}

	// Ambiguous timezone with DST (value type)
	if ambiguous, ok := val.(data.AmbiguousTimezoneMap); ok {
		return resolveAmbiguousTimezone(ambiguous, instant)
	}

	// Ambiguous timezone with DST (pointer type)
	if ambiguous, ok := val.(*data.AmbiguousTimezoneMap); ok {
		return resolveAmbiguousTimezone(*ambiguous, instant)
	}

	return nil
}

// resolveAmbiguousTimezone resolves an ambiguous timezone to its offset based on the instant.
func resolveAmbiguousTimezone(ambiguous data.AmbiguousTimezoneMap, instant time.Time) *int {
	// Without a valid instant, we can't determine DST status
	if instant.IsZero() {
		return nil
	}

	year := instant.Year()
	dstStart := ambiguous.DstStart(year)
	dstEnd := ambiguous.DstEnd(year)

	// Check if instant is during DST period
	var offset int
	if instant.After(dstStart) && !instant.After(dstEnd) {
		offset = ambiguous.TimezoneOffsetDuringDst
	} else {
		offset = ambiguous.TimezoneOffsetNonDst
	}
	return &offset
}
