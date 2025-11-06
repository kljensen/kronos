package types

// Timeunit represents a unit of time for calculations and operations.
// This is an internal type used for duration calculations and parsing.
type Timeunit string

// Time unit constants for internal use
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

// Meridiem represents the AM/PM indicator for internal parsing.
type Meridiem int

// Meridiem constants for internal use
const (
	MeridiemAM Meridiem = 0
	MeridiemPM Meridiem = 1
)

// Common time constants for internal calculations
const (
	HoursPerDay            = 24
	MinutesPerHour         = 60
	SecondsPerMinute       = 60
	MillisecondsPerSecond  = 1000
	MicrosecondsPerMS      = 1000
	MicrosecondsPerSecond  = 1000000
	NanosecondsPerMicro    = 1000
	NanosecondsPerMS       = 1000000
	SecondsPerHour         = 3600
	MinutesPerDay          = 1440
	DaysPerWeek            = 7
	MonthsPerYear          = 12
	MonthsPerQuarter       = 3
	WeeksPerMonthApprox    = 4
	YearLookAheadThreshold = 20 // For 2-digit year conversion
)
