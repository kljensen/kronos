package kronos

import (
	"fmt"
	"time"
)

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

// DayPreference specifies how to interpret day when it's ambiguous.
type DayPreference int

const (
	// DayPreferCurrent prefers the current day of month.
	DayPreferCurrent DayPreference = iota
	// DayPreferFirst prefers the first day of the period.
	DayPreferFirst
	// DayPreferLast prefers the last day of the period.
	DayPreferLast
)

// String returns the string representation of DayPreference.
func (d DayPreference) String() string {
	switch d {
	case DayPreferCurrent:
		return "current"
	case DayPreferFirst:
		return "first"
	case DayPreferLast:
		return "last"
	default:
		return "current"
	}
}

// Settings contains all configuration for date parsing.
// It provides a comprehensive way to customize parsing behavior,
// similar to Python's dateparser settings system.
//
// Deprecated: Direct use of Settings is discouraged in favor of the builder pattern.
// Use the fluent builder API for common options:
//
//	parser := kronos.New(en.Casual).
//	    WithDateOrder(kronos.DateOrderDMY).
//	    WithForwardDate(true).
//	    Strict()
//
// For advanced options, use the experimental package:
//
//	import "github.com/kljensen/kronos/experimental"
//
//	parser := kronos.New(en.Casual).
//	    WithOption(experimental.WithTimezoneOverrides(customTimezones)).
//	    WithOption(experimental.WithDebugHandler(debugFunc))
//
// The Settings struct will remain available for compatibility but new code should
// use the builder pattern.
type Settings struct {
	// Date interpretation
	DateOrder        DateOrder      // Order of date components (MDY, DMY, YMD)
	PreferDatesFrom  DatePreference // Past, Future, CurrentPeriod
	PreferDayOfMonth DayPreference  // Current, First, Last

	// Timezone handling
	Timezone            string // Default timezone name (e.g., "UTC", "America/New_York")
	ToTimezone          string // Convert results to this timezone
	ReturnTimezoneAware bool   // Include timezone information in results

	// Parsing behavior
	StrictParsing bool // Validate strictly - reject ambiguous dates

	// Period tracking
	ReturnTimeAsPeriod bool // Track parsing granularity (year, month, day, time)

	// Advanced configuration
	TimezoneOverrides TimezoneAbbrMap // Custom timezone abbreviations
	DebugHandler      DebugHandler    // Debug callback for parsing events
}

// ToParsingOption converts Settings to ParsingOption for backward compatibility.
// This method is used internally to bridge Settings and ParsingOption.
func (s Settings) ToParsingOption(timezones TimezoneAbbrMap) ParsingOption {
	// Merge default timezones with overrides
	mergedTimezones := timezones
	if s.TimezoneOverrides != nil {
		if mergedTimezones == nil {
			mergedTimezones = make(TimezoneAbbrMap)
		}
		for k, v := range s.TimezoneOverrides {
			mergedTimezones[k] = v
		}
	}

	return ParsingOption{
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
		DateOrder:           DateOrderMDY,
		PreferDatesFrom:     PreferCurrentPeriod,
		PreferDayOfMonth:    DayPreferCurrent,
		Timezone:            "UTC",
		ToTimezone:          "",
		ReturnTimezoneAware: false,
		StrictParsing:       false,
		ReturnTimeAsPeriod:  false,
	}
}

// ValidateSettings checks settings for consistency and validity.
// Returns an error if any setting is invalid.
func ValidateSettings(s Settings) error {
	// Validate timezone
	if s.Timezone != "" {
		if _, err := time.LoadLocation(s.Timezone); err != nil {
			return fmt.Errorf("invalid timezone: %s - %w", s.Timezone, err)
		}
	}

	// Validate ToTimezone
	if s.ToTimezone != "" {
		if _, err := time.LoadLocation(s.ToTimezone); err != nil {
			return fmt.Errorf("invalid to_timezone: %s - %w", s.ToTimezone, err)
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

	// Validate day preference
	if s.PreferDayOfMonth < DayPreferCurrent || s.PreferDayOfMonth > DayPreferLast {
		return fmt.Errorf("invalid day preference: %d", s.PreferDayOfMonth)
	}

	return nil
}

// ApplySettings creates a new ParsingContext with settings applied.
// This allows settings to influence the parsing context.
func applySettings(text string, refDate time.Time, settings Settings) (*ParsingContext, error) {
	// Validate settings first
	if err := ValidateSettings(settings); err != nil {
		return nil, err
	}

	// Apply normalization
	text = sanitizeInput(text)

	// Create parsing option from settings
	opt := ParsingOption{
		Preference: settings.PreferDatesFrom,
		DateOrder:  settings.DateOrder,
		Timezones:  nil,
		Debug:      nil,
	}

	// Create context
	ctx := newParsingContext(text, refDate, &opt)

	return ctx, nil
}
