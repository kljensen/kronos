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
	StrictParsing bool     // Validate strictly - reject ambiguous dates
	Normalize     bool     // Unicode normalization before parsing
	SkipTokens    []string // Words to ignore during parsing
	RequireParts  []string // Required components (e.g., "year", "month", "day")

	// Relative dates
	RelativeBase *time.Time // Base time for relative date calculations

	// Period tracking
	ReturnTimeAsPeriod bool // Track parsing granularity (year, month, day, time)

	// Parser control
	EnabledParsers []string // Which parsers to use (empty = all)
	ParserOrder    []string // Order of parser execution (empty = default)

	// Performance
	MaxParsers int           // Maximum number of parsers to run (0 = unlimited)
	Timeout    time.Duration // Parsing timeout (0 = no timeout)

	// Compatibility
	ForwardDate bool // Legacy: parse only forward dates
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
		Normalize:           true,
		SkipTokens:          []string{},
		RequireParts:        []string{},
		RelativeBase:        nil,
		ReturnTimeAsPeriod:  false,
		EnabledParsers:      []string{},
		ParserOrder:         []string{},
		MaxParsers:          0,
		Timeout:             0,
		ForwardDate:         false,
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

	// Validate RequireParts
	validParts := map[string]bool{
		"year":   true,
		"month":  true,
		"day":    true,
		"hour":   true,
		"minute": true,
		"second": true,
	}
	for _, part := range s.RequireParts {
		if !validParts[part] {
			return fmt.Errorf("invalid required part: %s (valid: year, month, day, hour, minute, second)", part)
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

// ToParsingOption converts Settings to the legacy ParsingOption format
// for backward compatibility with existing parsers.
func (s Settings) ToParsingOption(timezones TimezoneAbbrMap) ParsingOption {
	return ParsingOption{
		ForwardDate: s.ForwardDate,
		Preference:  s.PreferDatesFrom,
		DateOrder:   s.DateOrder,
		Timezones:   timezones,
		Debug:       nil,
	}
}

// ApplySettings creates a new ParsingContext with settings applied.
// This allows settings to influence the parsing context.
func ApplySettings(text string, refDate time.Time, settings Settings) (*ParsingContext, error) {
	// Validate settings first
	if err := ValidateSettings(settings); err != nil {
		return nil, err
	}

	// Apply normalization if enabled
	if settings.Normalize {
		text = SanitizeInput(text)
	}

	// Apply skip tokens
	for _, token := range settings.SkipTokens {
		// Simple token removal - could be enhanced with word boundary matching
		text = removeToken(text, token)
	}

	// Determine reference date
	ref := refDate
	if settings.RelativeBase != nil {
		ref = *settings.RelativeBase
	}

	// Create parsing option from settings
	opt := settings.ToParsingOption(nil)

	// Create context
	ctx := NewParsingContext(text, ref, &opt)

	return ctx, nil
}

// removeToken removes all occurrences of a token from text.
// This is a simple implementation that preserves spaces around tokens.
//
//nolint:gofumpt // Function formatting is correct
func removeToken(text string, token string) string {
	// Simple approach: iterate through the text and skip tokens
	result := ""
	i := 0
	for i < len(text) {
		if i+len(token) <= len(text) && text[i:i+len(token)] == token {
			// Check if this is a word boundary (not part of a larger word)
			isWordBoundary := true
			if i > 0 && text[i-1] != ' ' {
				isWordBoundary = false
			}
			if i+len(token) < len(text) && text[i+len(token)] != ' ' {
				isWordBoundary = false
			}

			if isWordBoundary {
				// Skip the token but preserve the surrounding space structure
				i += len(token)
				// If there's a trailing space, include one space in the result
				if i < len(text) && text[i] == ' ' {
					result += " "
					i++
				}
			} else {
				result += string(text[i])
				i++
			}
		} else {
			result += string(text[i])
			i++
		}
	}
	return result
}
