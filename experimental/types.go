//nolint:staticcheck // SA1019: Must use deprecated types during transition
package experimental

import (
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal"
)

// Component represents a date/time component (year, month, day, etc.).
// Re-exported from the main package for the experimental API.
type Component = kronos.Component

// Component constants for accessing specific date/time parts.
const (
	ComponentYear           = kronos.ComponentYear
	ComponentMonth          = kronos.ComponentMonth
	ComponentDay            = kronos.ComponentDay
	ComponentWeekday        = kronos.ComponentWeekday
	ComponentHour           = kronos.ComponentHour
	ComponentMinute         = kronos.ComponentMinute
	ComponentSecond         = kronos.ComponentSecond
	ComponentMillisecond    = kronos.ComponentMillisecond
	ComponentMicrosecond    = kronos.ComponentMicrosecond
	ComponentNanosecond     = kronos.ComponentNanosecond
	ComponentMeridiem       = kronos.ComponentMeridiem
	ComponentTimezoneOffset = kronos.ComponentTimezoneOffset
)

// Data constants
// These are commonly used constants for parsing.

// DefaultTimezoneAbbrMap is a map of common timezone abbreviations to their offsets in minutes.
var DefaultTimezoneAbbrMap = internal.DefaultTimezoneAbbrMap

// ApproximationWords is a list of words that indicate approximate time expressions.
var ApproximationWords = internal.ApproximationWords

// EmptyDuration represents an explicit empty duration.
var EmptyDuration = internal.EmptyDuration

// Duration and Timeunit types
// These are re-exported from the main package for convenience.

// Duration represents a time duration with named units.
type Duration = kronos.Duration

// Timeunit represents a unit of time (day, hour, minute, etc.).
type Timeunit = kronos.Timeunit

// Timeunit constants for duration calculations.
const (
	TimeunitYear        = kronos.TimeunitYear
	TimeunitMonth       = kronos.TimeunitMonth
	TimeunitWeek        = kronos.TimeunitWeek
	TimeunitDay         = kronos.TimeunitDay
	TimeunitHour        = kronos.TimeunitHour
	TimeunitMinute      = kronos.TimeunitMinute
	TimeunitSecond      = kronos.TimeunitSecond
	TimeunitMillisecond = kronos.TimeunitMillisecond
	TimeunitMicrosecond = kronos.TimeunitMicrosecond
	TimeunitNanosecond  = kronos.TimeunitNanosecond
	TimeunitQuarter     = kronos.TimeunitQuarter
	TimeunitDecade      = kronos.TimeunitDecade
)

// Period represents the granularity of a parsed date/time.
type Period = kronos.Period

// Period constants.
const (
	PeriodTime  = kronos.PeriodTime
	PeriodDay   = kronos.PeriodDay
	PeriodWeek  = kronos.PeriodWeek
	PeriodMonth = kronos.PeriodMonth
	PeriodYear  = kronos.PeriodYear
)

// Settings contains all configuration for date parsing.
// Re-exported from the main package for the experimental API.
type Settings = kronos.Settings

// DateOrder represents the order of date components in ambiguous formats.
// Re-exported from the main package for the experimental API.
type DateOrder = kronos.DateOrder

// DayPreference specifies how to interpret day when it's ambiguous.
// Re-exported from the main package for the experimental API.
type DayPreference = kronos.DayPreference

// DatePreference specifies whether to prefer past or future dates.
// Re-exported from the main package for the experimental API.
type DatePreference = kronos.DatePreference

// DatePreference constants.
const (
	PreferPast          = kronos.PreferPast
	PreferFuture        = kronos.PreferFuture
	PreferCurrentPeriod = kronos.PreferCurrentPeriod
)

// Weekday type alias for convenience.
type Weekday = time.Weekday

// Weekday constants.
const (
	Sunday    = time.Sunday
	Monday    = time.Monday
	Tuesday   = time.Tuesday
	Wednesday = time.Wednesday
	Thursday  = time.Thursday
	Friday    = time.Friday
	Saturday  = time.Saturday
)

// TimezoneAbbrMap maps timezone abbreviations to their offsets.
// Re-exported from the main package for the experimental API.
type TimezoneAbbrMap = kronos.TimezoneAbbrMap

// AmbiguousTimezoneMap defines a timezone with different DST offsets.
// Re-exported from the main package for the experimental API.
type AmbiguousTimezoneMap = kronos.AmbiguousTimezoneMap

// DebugHandler is a function that handles debug events.
// Re-exported from the main package for the experimental API.
type DebugHandler = kronos.DebugHandler

// DefaultSettings returns sensible default settings.
// This is re-exported from the main package for backward compatibility.
var DefaultSettings = kronos.DefaultSettings

// ValidateSettings checks settings for consistency and validity.
// This is re-exported from the main package for backward compatibility.
var ValidateSettings = kronos.ValidateSettings
