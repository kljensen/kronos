package kronos

import "time"

// Component represents a date/time component that can be parsed.
// Components are used as keys in maps, so they are string constants.
type Component string

const (
	ComponentYear           Component = "year"
	ComponentMonth          Component = "month"
	ComponentDay            Component = "day"
	ComponentWeekday        Component = "weekday"
	ComponentHour           Component = "hour"
	ComponentMinute         Component = "minute"
	ComponentSecond         Component = "second"
	ComponentMillisecond    Component = "millisecond"
	ComponentMeridiem       Component = "meridiem"
	ComponentTimezoneOffset Component = "timezoneOffset"
)

// Timeunit represents a unit of time for calculations and operations.
type Timeunit string

const (
	TimeunitYear        Timeunit = "year"
	TimeunitMonth       Timeunit = "month"
	TimeunitWeek        Timeunit = "week"
	TimeunitDay         Timeunit = "day"
	TimeunitHour        Timeunit = "hour"
	TimeunitMinute      Timeunit = "minute"
	TimeunitSecond      Timeunit = "second"
	TimeunitMillisecond Timeunit = "millisecond"
	TimeunitQuarter     Timeunit = "quarter"
)

// Meridiem represents the AM/PM indicator.
type Meridiem int

const (
	MeridiemAM Meridiem = 0
	MeridiemPM Meridiem = 1
)

// Common time constants
const (
	HoursPerDay            = 24
	MinutesPerHour         = 60
	SecondsPerMinute       = 60
	MillisecondsPerSecond  = 1000
	NanosecondsPerMS       = 1000000
	SecondsPerHour         = 3600
	MinutesPerDay          = 1440
	DaysPerWeek            = 7
	MonthsPerYear          = 12
	MonthsPerQuarter       = 3
	WeeksPerMonthApprox    = 4
	YearLookAheadThreshold = 20 // For 2-digit year conversion
)

// Weekday represents the day of the week.
type Weekday int

const (
	WeekdaySunday Weekday = iota
	WeekdayMonday
	WeekdayTuesday
	WeekdayWednesday
	WeekdayThursday
	WeekdayFriday
	WeekdaySaturday
)

// Month represents a calendar month.
type Month int

const (
	MonthJanuary Month = iota + 1
	MonthFebruary
	MonthMarch
	MonthApril
	MonthMay
	MonthJune
	MonthJuly
	MonthAugust
	MonthSeptember
	MonthOctober
	MonthNovember
	MonthDecember
)

// DebugHandler is a function that handles debug events.
// It receives a debug message for logging or analysis.
type DebugHandler func(message string)

// ParsingOption contains configuration options for parsing.
type ParsingOption struct {
	// ForwardDate indicates whether to parse only forward dates
	// (results should be after the reference date).
	// This affects date/time implication (e.g. weekday or time mentioning).
	ForwardDate bool

	// Timezones provides additional timezone keywords for parsers to recognize.
	// Any value provided will override the default handling of that value.
	Timezones TimezoneAbbrMap

	// Debug is an internal debug event handler.
	Debug DebugHandler
}

// AmbiguousTimezoneMap defines a timezone that has different offsets
// depending on whether daylight saving time (DST) is in effect.
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
type TimezoneAbbrMap map[string]interface{}

// ParsingReference contains reference information for parsing dates/times.
type ParsingReference struct {
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
