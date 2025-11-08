package kronos

import (
	"fmt"
	"maps"
	"time"
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

// DateOrder represents the order of date components in ambiguous formats.
type DateOrder int

const (
	// DateOrderMDY represents Month-Day-Year order (US format).
	DateOrderMDY DateOrder = iota
	// DateOrderDMY represents Day-Month-Year order (European format).
	DateOrderDMY
	// DateOrderYMD represents Year-Month-Day order (ISO format).
	DateOrderYMD
)

// String returns the string representation of DateOrder.
func (d DateOrder) String() string {
	switch d {
	case DateOrderMDY:
		return "MDY"
	case DateOrderDMY:
		return "DMY"
	case DateOrderYMD:
		return "YMD"
	default:
		return "MDY"
	}
}

// Settings contains all configuration for date parsing.
// It provides a comprehensive way to customize parsing behavior.
//
// Most users should use the fluent builder API for common options:
//
//	parser := en.New().
//	    DateOrder(kronos.DateOrderDMY).
//	    PreferFuture().
//	    Strict()
//
// For advanced options, use WithOption to modify settings directly:
//
//	parser := en.New().
//	    WithOption(func(s *kronos.Settings) {
//	        s.TimezoneOverrides = customTimezones
//	        s.DebugHandler = debugFunc
//	    })
type Settings struct {
	// Date interpretation
	DateOrder       DateOrder      // Order of date components (MDY, DMY, YMD)
	PreferDatesFrom DatePreference // Past, Future, CurrentPeriod

	// Timezone handling
	Timezone string // Default timezone name (e.g., "UTC", "America/New_York")

	// Parsing behavior
	StrictParsing bool // Validate strictly - reject ambiguous dates

	// Advanced configuration
	TimezoneOverrides TimezoneAbbrMap // Custom timezone abbreviations
	DebugHandler      DebugHandler    // Debug callback for parsing events
}

// toparsingOption converts Settings to parsingOption for backward compatibility.
// This method is used internally to bridge Settings and parsingOption.
func (s Settings) toparsingOption(timezones TimezoneAbbrMap) parsingOption {
	// Merge default timezones with overrides
	mergedTimezones := timezones
	if len(s.TimezoneOverrides) > 0 {
		if mergedTimezones == nil {
			mergedTimezones = s.TimezoneOverrides
		} else {
			// Only create new map if we need to merge
			mergedTimezones = make(TimezoneAbbrMap, len(timezones)+len(s.TimezoneOverrides))
			maps.Copy(mergedTimezones, timezones)
			maps.Copy(mergedTimezones, s.TimezoneOverrides)
		}
	}

	return parsingOption{
		ForwardDate: s.PreferDatesFrom == PreferFuture,
		Preference:  s.PreferDatesFrom,
		DateOrder:   s.DateOrder,
		Timezones:   mergedTimezones,
		Debug:       s.DebugHandler,
	}
}

// DefaultSettings returns sensible default settings that maintain
// backward compatibility with existing behavior.
func DefaultSettings() Settings {
	return Settings{
		DateOrder:       DateOrderMDY,
		PreferDatesFrom: PreferCurrentPeriod,
		Timezone:        "UTC",
		StrictParsing:   false,
	}
}

// validateSettings checks settings for consistency and validity.
// Returns an error if any setting is invalid.
func validateSettings(s Settings) error {
	// Validate timezone
	if s.Timezone != "" {
		if _, err := time.LoadLocation(s.Timezone); err != nil {
			return fmt.Errorf("invalid timezone: %s - %w", s.Timezone, err)
		}
	}

	// Validate date order
	if s.DateOrder < DateOrderMDY || s.DateOrder > DateOrderYMD {
		return fmt.Errorf("invalid date order: %d", s.DateOrder)
	}

	// Validate date preference
	if s.PreferDatesFrom < PreferCurrentPeriod || s.PreferDatesFrom > PreferFuture {
		return fmt.Errorf("invalid date preference: %d", s.PreferDatesFrom)
	}

	return nil
}

// ApplySettings creates a new ParsingContext with settings applied.
// This allows settings to influence the parsing context.
func applySettings(text string, refDate time.Time, settings Settings) (*parsingContext, error) {
	// Validate settings first
	if err := validateSettings(settings); err != nil {
		return nil, err
	}

	// Create parsing option from settings
	opt := parsingOption{
		Preference: settings.PreferDatesFrom,
		DateOrder:  settings.DateOrder,
		Timezones:  nil,
		Debug:      nil,
	}

	// Create context (sanitization happens in newParsingContext)
	ctx := newParsingContext(text, refDate, &opt)
	ctx.settings = &settings

	return ctx, nil
}
