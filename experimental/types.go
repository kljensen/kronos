package experimental

import (
	"github.com/kljensen/kronos"
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

// NewParsingContext creates a new parsing context.
// Re-exported from the main package for the experimental API.
var NewParsingContext = kronos.XNewParsingContext

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
