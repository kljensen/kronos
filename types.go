package kronos

import (
	"time"
)

// Component represents a date/time component that can be parsed.
// Components are used as keys in maps, so they are string constants.
type Component string

// Date/time component constants
const (
	// ComponentYear represents the year component
	ComponentYear           Component = "year"
	ComponentMonth          Component = "month"
	ComponentDay            Component = "day"
	ComponentWeekday        Component = "weekday"
	ComponentHour           Component = "hour"
	ComponentMinute         Component = "minute"
	ComponentSecond         Component = "second"
	ComponentMillisecond    Component = "millisecond"
	ComponentMicrosecond    Component = "microsecond"
	ComponentNanosecond     Component = "nanosecond"
	ComponentMeridiem       Component = "meridiem"
	ComponentTimezoneOffset Component = "timezoneOffset"
)

// Timeunit represents a unit of time for calculations and operations.
// This type is used as a key in the Duration map type for specifying time durations.
//
// Note: While Timeunit is part of the public API (used by Duration), it is primarily
// an implementation detail. Future versions may move this to a more restricted scope
// while maintaining backward compatibility for Duration operations.
type Timeunit string

// Time unit constants for use with Duration type.
// These constants are required for working with Duration maps.
//
// Example:
//
//	duration := kronos.Duration{
//	    kronos.TimeunitDay: 5,
//	    kronos.TimeunitHour: 3,
//	}
const (
	TimeunitYear        Timeunit = "year"
	TimeunitMonth       Timeunit = "month"
	TimeunitWeek        Timeunit = "week"
	TimeunitDay         Timeunit = "day"
	TimeunitHour        Timeunit = "hour"
	TimeunitMinute      Timeunit = "minute"
	TimeunitSecond      Timeunit = "second"
	TimeunitMillisecond Timeunit = "millisecond"
	TimeunitMicrosecond Timeunit = "microsecond"
	TimeunitNanosecond  Timeunit = "nanosecond"
	TimeunitQuarter     Timeunit = "quarter"
	TimeunitDecade      Timeunit = "decade"
)

// DatePreference specifies how ambiguous dates (with missing components) should be resolved.
// This controls whether dates like "March 15" (without year) or "10:00" (without date)
// should be interpreted as being in the past, future, or current period relative to the reference time.
type DatePreference int

const (
	// PreferCurrentPeriod (default) chooses the date in the current period.
	// For dates with missing year: chooses the current year.
	// For times with missing date: chooses the current day.
	// Example: If reference is Feb 15, 2015 15:30
	//   - "March" → March 2015 (current year)
	//   - "10:00" → Feb 15, 2015 10:00 (current day, even if past)
	PreferCurrentPeriod DatePreference = iota

	// PreferPast chooses dates in the past.
	// For dates with missing year: if the date would be in the future, use previous year.
	// For times with missing date: if the time would be in the future, use previous day.
	// Example: If reference is Feb 15, 2015 15:30
	//   - "March" → March 2014 (last March, in the past)
	//   - "10:00" → Feb 15, 2015 10:00 (earlier today)
	//   - "18:00" → Feb 14, 2015 18:00 (yesterday, since 18:00 today hasn't happened yet)
	PreferPast

	// PreferFuture chooses dates in the future.
	// For dates with missing year: if the date would be in the past, use next year.
	// For times with missing date: if the time would be in the past, use next day.
	// Example: If reference is Feb 15, 2015 15:30
	//   - "March" → March 2015 (next March, in the future)
	//   - "10:00" → Feb 16, 2015 10:00 (tomorrow morning, since 10:00 already passed today)
	//   - "18:00" → Feb 15, 2015 18:00 (later today)
	PreferFuture
)

// Period represents the granularity of a parsed date/time expression.
// It indicates the finest level of precision explicitly mentioned in the input.
// For example:
//   - "2020" has year-level precision
//   - "March 2020" has month-level precision
//   - "yesterday" has day-level precision
//   - "10:30" has time-level precision
type Period int

const (
	// PeriodUnknown indicates the period could not be determined.
	PeriodUnknown Period = iota
	// PeriodYear represents year-level precision (e.g., "2020", "last year").
	PeriodYear
	// PeriodMonth represents month-level precision (e.g., "March", "3 months ago").
	PeriodMonth
	// PeriodWeek represents week-level precision (e.g., "last week", "2 weeks ago").
	PeriodWeek
	// PeriodDay represents day-level precision (e.g., "yesterday", "March 15").
	PeriodDay
	// PeriodTime represents time-level precision (e.g., "10:30", "2 hours ago").
	PeriodTime
)

// String returns the string representation of the Period.
func (p Period) String() string {
	switch p {
	case PeriodYear:
		return "year"
	case PeriodMonth:
		return "month"
	case PeriodWeek:
		return "week"
	case PeriodDay:
		return "day"
	case PeriodTime:
		return "time"
	default:
		return "unknown"
	}
}

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

// ParsingOption contains configuration options for parsing.
//
// Deprecated: Direct use of ParsingOption is discouraged. Use the builder pattern
// for common options and WithOption for advanced options:
//
//	// Common usage:
//	parser := kronos.New(en.Casual).
//	    WithDateOrder(kronos.DateOrderDMY).
//	    PreferFuture()
//
//	// Advanced usage with WithOption:
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.TimezoneOverrides = customTimezones
//	    })
type parsingOption struct {
	// ForwardDate indicates whether to parse only forward dates
	// (results should be after the reference date).
	// This affects date/time implication (e.g. weekday or time mentioning).
	ForwardDate bool

	// Preference specifies how ambiguous dates should be resolved.
	// Controls whether dates like "March 15" (without year) or "10:00" (without date)
	// should be interpreted as past, future, or current period.
	// Default is PreferCurrentPeriod.
	Preference DatePreference

	// DateOrder specifies the order of date components in ambiguous formats.
	// Use DateOrderMDY for US format (12/31/2020), DateOrderDMY for European format (31/12/2020),
	// or DateOrderYMD for ISO format (2020/12/31).
	// Default is DateOrderMDY.
	DateOrder DateOrder

	// Timezones provides additional timezone keywords for parsers to recognize.
	// Any value provided will override the default handling of that value.
	Timezones TimezoneAbbrMap

	// Debug is an internal debug event handler.
	Debug DebugHandler
}

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
type AmbiguousTimezoneMap struct {
	// TimezoneOffsetDuringDst is the offset in minutes during DST.
	TimezoneOffsetDuringDst int

	// TimezoneOffsetNonDst is the offset in minutes when DST is not in effect.
	TimezoneOffsetNonDst int

	// DstStart returns the start date of DST for the given year.
	DstStart func(year int) time.Time

	// DstEnd returns the end date of DST for the given year.
	DstEnd func(year int) time.Time
}

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
type TimezoneAbbrMap map[string]interface{}

// ParsingReference contains reference information for parsing dates/times.
//
// Deprecated: Direct use of ParsingReference is discouraged. Use the builder pattern
// with WithReferenceDate instead:
//
//	parser := kronos.New(en.Casual).
//	    WithReferenceDate(time.Now())
//
// For advanced timezone reference configuration, use the experimental package.
type parsingReference struct {
	// Instant is the reference date/time when the input is written or mentioned.
	// This affects date/time implication (e.g. weekday or time mentioning).
	// If nil, the current time is used.
	Instant *time.Time

	// Timezone is the reference timezone where the input is written or mentioned.
	// Date/time implication will account for the difference between input timezone
	// and the current system timezone.
	// Can be either a timezone name (string) or offset in minutes (int).
	Timezone interface{}
}

// ParsedComponents represents a collection of parsed date/time components.
// Each component has three levels of certainty:
//   - Certain (or Known): The component is directly mentioned and parsed.
//   - Implied: The component is not directly mentioned, but implied by other information.
//   - Unknown: The component is not mentioned at all.
type ParsedComponents interface {
	// IsCertain returns true if the component is certain (directly mentioned).
	IsCertain(component Component) bool

	// Get returns the component value for either Certain or Implied values.
	// Returns nil if the component is Unknown.
	Get(component Component) *int

	// Date returns a time.Time object constructed from the components.
	Date() time.Time

	// Tags returns debugging tags for the parsed component.
	Tags() map[string]bool
}

// ParsedResult represents a parsed result or final output.
// Each result object represents a date/time (or date/time range) mentioned in the input.
type ParsedResult interface {
	// RefDate returns the reference date used for parsing.
	RefDate() time.Time

	// Index returns the position in the input text where this result was found.
	Index() int

	// Text returns the matched text from the input.
	Text() string

	// Start returns the starting date/time components.
	Start() ParsedComponents

	// End returns the ending date/time components for a range, or nil for a single date/time.
	End() ParsedComponents

	// Date returns a time.Time object created from the start components.
	Date() time.Time

	// Tags returns combined debugging tags from start and end components.
	Tags() map[string]bool
}
