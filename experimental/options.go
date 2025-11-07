//nolint:staticcheck // SA1019: Must use deprecated types during transition
package experimental

import (
	"github.com/kljensen/kronos"
)

// Option is a function that modifies parser settings.
// This allows advanced configuration without cluttering the main builder API.
//
// Deprecated: The experimental options are being removed. Use the main ParserBuilder
// API methods instead (PreferPast(), PreferFuture(), DateOrder(), etc.).
type Option func(*kronos.Settings)

// WithDayPreference sets how to interpret day when it's ambiguous.
// Use DayPreferCurrent (default), DayPreferFirst, or DayPreferLast.
//
// Deprecated: This option exposes internal implementation details. If you need
// this level of control, access Settings directly via ParserBuilder.Settings().
func WithDayPreference(pref kronos.DayPreference) Option {
	return func(s *kronos.Settings) {
		s.PreferDayOfMonth = pref
	}
}

// WithTimezoneAware enables returning timezone information in results.
// This is an experimental feature.
//
// Deprecated: This option exposes internal implementation details. If you need
// this level of control, access Settings directly via ParserBuilder.Settings().
func WithTimezoneAware(aware bool) Option {
	return func(s *kronos.Settings) {
		s.ReturnTimezoneAware = aware
	}
}

// WithTimeAsPeriod enables tracking parsing granularity.
// This tracks whether the parser found year, month, day, or time precision.
//
// Deprecated: This option exposes internal implementation details. If you need
// this level of control, access Settings directly via ParserBuilder.Settings().
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
// Deprecated: This option exposes internal implementation details. If you need
// this level of control, access Settings directly via ParserBuilder.Settings().
func WithTimezoneOverrides(timezones kronos.TimezoneAbbrMap) Option {
	return func(s *kronos.Settings) {
		s.TimezoneOverrides = timezones
	}
}

// WithDebugHandler sets a debug callback for parsing events.
// This is useful for understanding how the parser works and debugging issues.
// The handler receives debug messages during parsing.
//
// Deprecated: This option exposes internal implementation details. If you need
// this level of control, access Settings directly via ParserBuilder.Settings().
func WithDebugHandler(handler kronos.DebugHandler) Option {
	return func(s *kronos.Settings) {
		s.DebugHandler = handler
	}
}
