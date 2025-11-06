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
var NewParsingContext = kronos.NewParsingContext
