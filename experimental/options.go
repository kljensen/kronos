package experimental

import (
	"time"

	"github.com/kljensen/kronos"
)

// Option is a function that modifies parser settings.
// This allows advanced configuration without cluttering the main builder API.
type Option func(*kronos.Settings)

// WithParserOrder sets the order in which parsers are executed.
// Parser names should match those registered in the global registry.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithParserOrder("iso8601", "en_casual_date"))
func WithParserOrder(parserNames ...string) Option {
	return func(s *kronos.Settings) {
		s.ParserOrder = parserNames
	}
}

// WithEnabledParsers specifies which parsers to use.
// If not set, all parsers from the configuration are used.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithEnabledParsers("iso8601", "en_casual_date"))
func WithEnabledParsers(parserNames ...string) Option {
	return func(s *kronos.Settings) {
		s.EnabledParsers = parserNames
	}
}

// DisableDefaultParsers removes all default parsers.
// Use with WithEnabledParsers to specify a custom parser set.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.DisableDefaultParsers()).
//	    WithOption(experimental.WithEnabledParsers("iso8601"))
func DisableDefaultParsers() Option {
	return func(s *kronos.Settings) {
		s.EnabledParsers = []string{}
	}
}

// WithMaxParsers limits the number of parsers to execute.
// This can improve performance when you only need the first few matches.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithMaxParsers(3))
func WithMaxParsers(max int) Option {
	return func(s *kronos.Settings) {
		s.MaxParsers = max
	}
}

// WithTimeout sets a parsing timeout.
// If parsing takes longer than this duration, it will be aborted.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithTimeout(5 * time.Second))
func WithTimeout(timeout time.Duration) Option {
	return func(s *kronos.Settings) {
		s.Timeout = timeout
	}
}

// WithSkipTokens specifies words to ignore during parsing.
// These tokens are removed from the input text before parsing.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithSkipTokens("at", "on"))
func WithSkipTokens(tokens ...string) Option {
	return func(s *kronos.Settings) {
		s.SkipTokens = tokens
	}
}

// WithRequiredParts specifies which date/time components must be present.
// Valid parts: "year", "month", "day", "hour", "minute", "second"
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithRequiredParts("year", "month", "day"))
func WithRequiredParts(parts ...string) Option {
	return func(s *kronos.Settings) {
		s.RequireParts = parts
	}
}

// WithNormalization enables or disables Unicode normalization.
// Normalization is enabled by default.
//
// Example:
//
//	builder := kronos.New().
//	    WithOption(experimental.WithNormalization(false))
func WithNormalization(normalize bool) Option {
	return func(s *kronos.Settings) {
		s.Normalize = normalize
	}
}

// WithRelativeBase sets a custom base time for relative date calculations.
// If not set, the reference date passed to Parse is used.
//
// Example:
//
//	baseTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
//	builder := kronos.New().
//	    WithOption(experimental.WithRelativeBase(baseTime))
func WithRelativeBase(base time.Time) Option {
	return func(s *kronos.Settings) {
		s.RelativeBase = &base
	}
}

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
