//nolint:staticcheck // SA1019: Must use deprecated types during transition
package experimental

import (
	"github.com/kljensen/kronos"
)

// Option is a function that modifies parser settings.
// This allows advanced configuration without cluttering the main builder API.
type Option func(*kronos.Settings)

// WithDayPreference sets how to interpret day when it's ambiguous.
// Use DayPreferCurrent (default), DayPreferFirst, or DayPreferLast.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithDayPreference(kronos.DayPreferFirst))
func WithDayPreference(pref kronos.DayPreference) Option {
	return func(s *kronos.Settings) {
		s.PreferDayOfMonth = pref
	}
}

// WithTimezoneAware enables returning timezone information in results.
// This is an experimental feature.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithTimezoneAware(true))
func WithTimezoneAware(aware bool) Option {
	return func(s *kronos.Settings) {
		s.ReturnTimezoneAware = aware
	}
}

// WithTimeAsPeriod enables tracking parsing granularity.
// This tracks whether the parser found year, month, day, or time precision.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithTimeAsPeriod(true))
func WithTimeAsPeriod(asPeriod bool) Option {
	return func(s *kronos.Settings) {
		s.ReturnTimeAsPeriod = asPeriod
	}
}

// WithTimezoneOverrides sets custom timezone abbreviation mappings.
// This allows you to override the default timezone abbreviations with custom ones.
// Values can be either a simple offset in minutes (int) or an AmbiguousTimezoneMap
// for timezones that observe DST.
//
// Example:
//
//	import "github.com/kljensen/kronos/experimental"
//
//	customTimezones := kronos.TimezoneAbbrMap{
//	    "CUSTOM": 123,  // Custom timezone at UTC+2:03
//	    "TEST": -456,   // Custom timezone at UTC-7:36
//	}
//	builder := kronos.New().
//	    WithOption(experimental.WithTimezoneOverrides(customTimezones))
func WithTimezoneOverrides(timezones kronos.TimezoneAbbrMap) Option {
	return func(s *kronos.Settings) {
		s.TimezoneOverrides = timezones
	}
}

// WithDebugHandler sets a debug callback for parsing events.
// This is useful for understanding how the parser works and debugging issues.
// The handler receives debug messages during parsing.
//
// Example:
//
//	import "log"
//	import "github.com/kljensen/kronos/experimental"
//
//	builder := kronos.New().
//	    WithOption(experimental.WithDebugHandler(func(msg string) {
//	        log.Printf("DEBUG: %s", msg)
//	    }))
func WithDebugHandler(handler kronos.DebugHandler) Option {
	return func(s *kronos.Settings) {
		s.DebugHandler = handler
	}
}
