package helpers

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
