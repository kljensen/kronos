//nolint:staticcheck // SA1019: Must use deprecated types during transition
package experimental

import (
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/data"
	"github.com/kljensen/kronos/internal/helpers"
)

// ParsingContext provides context during the parsing process.
// Re-exported from the main package for the experimental API.
type ParsingContext = kronos.ParsingContext

// ParsingResult represents a parsed date/time result.
// Re-exported from the main package for the experimental API.
type ParsingResult = kronos.ParsingResult

// ParsingComponents represents parsed date/time components.
// Re-exported from the main package for the experimental API.
type ParsingComponents = kronos.ParsingComponents

// ReferenceWithTimezone represents a reference date/time with an optional timezone offset.
// Re-exported from the main package for the experimental API.
type ReferenceWithTimezone = kronos.ReferenceWithTimezone

// ParsingOption configures parsing behavior.
// Re-exported from the main package for the experimental API.
type ParsingOption = kronos.ParsingOption

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

// Factory functions for creating parsing types
// These are re-exported for use by custom parser authors.

// NewParsingContext creates a new parsing context.
// Re-exported from the helpers package for the experimental API.
var NewParsingContext = kronos.XNewParsingContext

// NewParsingComponents creates a new ParsingComponents instance.
// Re-exported from the helpers package for the experimental API.
var NewParsingComponents = helpers.NewParsingComponents

// NewParsingResult creates a new ParsingResult instance.
// Re-exported from the main package for the experimental API.
var NewParsingResult = kronos.XNewParsingResult

// Date construction helpers
// These functions create ParsingComponents for common relative dates.

// Today returns components representing today's date.
// Date components are certain, time components are implied.
var Today = helpers.Today

// Tomorrow returns components representing tomorrow's date.
// Date components are certain, time components are implied.
var Tomorrow = helpers.Tomorrow

// Yesterday returns components representing yesterday's date.
// Date components are certain, time components are implied.
var Yesterday = helpers.Yesterday

// Now returns components representing the current moment.
// Both date and time components are certain, and it includes timezone offset.
var Now = helpers.Now

// Midnight returns components representing midnight.
// If the reference time is after 2 AM, it refers to the coming midnight (next day).
// Otherwise, it refers to the current midnight.
var Midnight = helpers.Midnight

// Noon returns components representing noon (12:00 PM).
var Noon = helpers.Noon

// Morning returns components representing morning time.
// Time is implied (default: 6 AM / 06:00).
var Morning = helpers.Morning

// Afternoon returns components representing afternoon time.
// Time is implied (default: 3 PM / 15:00).
var Afternoon = helpers.Afternoon

// Evening returns components representing evening time.
// Time is implied (default: 8 PM / 20:00).
var Evening = helpers.Evening

// Component helpers
// These functions manipulate ParsingComponents.

// AssignSimilarDate assigns (force updates) the parsing components to the same day as the target.
// This sets year, month, and day as certain (known) values.
var AssignSimilarDate = helpers.AssignSimilarDate

// AssignSimilarTime assigns (force updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as certain (known) values.
var AssignSimilarTime = helpers.AssignSimilarTime

// MergeDateTimeComponent merges date and time components.
var MergeDateTimeComponent = helpers.MergeDateTimeComponent

// Date math helpers
// These functions perform date calculations.

// GetLastWeekday returns the date of the last occurrence of the target weekday
// before the reference date (not including the reference date itself).
var GetLastWeekday = helpers.GetLastWeekday

// GetNextWeekday returns the date of the next occurrence of the target weekday
// after the reference date (not including the reference date itself).
var GetNextWeekday = helpers.GetNextWeekday

// GetThisWeekday returns the date of the target weekday in "this" week.
// The forward parameter controls the direction:
// - If forward=true: looks forward from refDate (inclusive)
// - If forward=false: looks backward from refDate (inclusive)
var GetThisWeekday = helpers.GetThisWeekday

// AddDuration returns the date after adding the given duration to ref.
// It handles fractional durations by cascading remainders to smaller units.
// Returns an error if the duration or resulting date is out of bounds.
var AddDuration = helpers.AddDuration

// ReverseDuration returns the reversed duration (e.g., back into the past instead of future).
// All values in the duration are negated.
var ReverseDuration = helpers.ReverseDuration

// FindMostLikelyADYear converts a 2-digit year to a 4-digit year.
// Years 0-99 are mapped to 1900-2099 range based on current year and threshold.
var FindMostLikelyADYear = helpers.FindMostLikelyADYear

// FindYearClosestToRefWithPreference finds the year based on date preference settings.
// This allows control over whether ambiguous dates should be resolved to past, future, or current period.
var FindYearClosestToRefWithPreference = helpers.FindYearClosestToRefWithPreference

// GetDaysToWeekday returns the number of days from refDate to the target weekday
// based on the modifier ("this", "next", "last", or nil for closest).
var GetDaysToWeekday = helpers.GetDaysToWeekday

// GetLastWeekdayOfMonth returns the date of the last occurrence of a given weekday
// in a given month and year.
var GetLastWeekdayOfMonth = helpers.GetLastWeekdayOfMonth

// GetNthWeekdayOfMonth returns the date of the nth occurrence of a given weekday
// in a given month and year.
var GetNthWeekdayOfMonth = helpers.GetNthWeekdayOfMonth

// Data constants
// These are commonly used constants for parsing.

// DefaultTimezoneAbbrMap is a map of common timezone abbreviations to their offsets in minutes.
var DefaultTimezoneAbbrMap = data.DefaultTimezoneAbbrMap

// ApproximationWords is a list of words that indicate approximate time expressions.
var ApproximationWords = data.ApproximationWords

// EmptyDuration represents an explicit empty duration.
var EmptyDuration = data.EmptyDuration

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

// ParsingReference contains reference information for parsing.
// Re-exported from the main package for the experimental API.
type ParsingReference = kronos.ParsingReference

// DefaultSettings returns sensible default settings.
// This is re-exported from the main package for backward compatibility.
var DefaultSettings = kronos.DefaultSettings

// ValidateSettings checks settings for consistency and validity.
// This is re-exported from the main package for backward compatibility.
var ValidateSettings = kronos.ValidateSettings
