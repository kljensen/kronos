// Package en provides English language support for Chrono.
// It includes parsers, refiners, and configurations for parsing English dates and times.
package en

import (
	"fmt"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/common/parsers"
	"github.com/kljensen/kronos/internal/common/refiners"
	enparsers "github.com/kljensen/kronos/internal/en/parsers"
	enrefiners "github.com/kljensen/kronos/internal/en/refiners"
)

// New creates a new ParserBuilder using the casual English configuration.
// This is the recommended way to parse English dates with the builder pattern.
//
// Example:
//
//	parser := en.New().
//	    WithReferenceDate(time.Now()).
//	    PreferPast()
//	results, err := parser.Parse("last Monday")
func New() *kronos.ParserBuilder {
	return kronos.New(englishCasualChrono())
}

// NewStrict creates a new ParserBuilder using strict English configuration.
// In strict mode, only formal date/time patterns are recognized.
//
// Example:
//
//	parser := en.NewStrict().WithReferenceDate(time.Now())
//	results, err := parser.Parse("2020-03-15")
func NewStrict() *kronos.ParserBuilder {
	return kronos.New(englishStrictChrono())
}

// NewGB creates a new ParserBuilder using UK-style English configuration.
// It uses little-endian date format (day/month/year) and casual expressions.
//
// Example:
//
//	parser := en.NewGB().WithReferenceDate(time.Now())
//	results, err := parser.Parse("15/03/2020")
func NewGB() *kronos.ParserBuilder {
	return kronos.New(englishGBChrono()).DateOrder(kronos.DateOrderDMY)
}

// includeCommonConfiguration adds common parsers and refiners to a configuration.
func includeCommonConfiguration(config *kronos.Configuration, strictMode bool) *kronos.Configuration {
	// Add ISO format parser at the beginning
	config.Parsers = append([]kronos.Parser{parsers.NewISOFormatParser()}, config.Parsers...)

	// Add common refiners at the beginning
	config.Refiners = append([]kronos.Refiner{
		refiners.NewMergeWeekdayComponentRefiner(),
		refiners.NewExtractTimezoneOffsetRefiner(),
		refiners.NewOverlapRemovalRefiner(),
	}, config.Refiners...)

	// Add common refiners at the end
	config.Refiners = append(config.Refiners,
		refiners.NewExtractTimezoneAbbrRefiner(nil),
		refiners.NewOverlapRemovalRefiner(),
		refiners.NewDatePreferenceRefiner(), // Apply date preferences before ForwardDateRefiner
		refiners.NewForwardDateRefiner(),    // ForwardDate option overrides preferences
		refiners.NewUnlikelyFormatFilter(strictMode),
	)

	return config
}

// createCasualConfiguration creates a casual English configuration.
func createCasualConfiguration(littleEndian bool) *kronos.Configuration {
	config := createConfiguration(false, littleEndian)

	// Add unlikely format filter for additional filtering
	config.Refiners = append(config.Refiners, enrefiners.NewENUnlikelyFormatFilter())

	return config
}

// createConfiguration creates a standard English configuration.
func createConfiguration(strictMode, littleEndian bool) *kronos.Configuration {
	config := &kronos.Configuration{
		Parsers: []kronos.Parser{
			parsers.NewSlashDateFormatParser(littleEndian),
			enparsers.NewENTimeUnitWithinFormatParser(strictMode),
			enparsers.NewENMonthNameLittleEndianParser(),
			enparsers.NewENMonthNameMiddleEndianParser(littleEndian), // shouldSkipYearLikeDate
			enparsers.NewENWeekdayParser(),
			enparsers.NewENSlashMonthFormatParser(),
			enparsers.NewENTimeExpressionParser(strictMode),
			enparsers.NewENTimeUnitAgoFormatParser(strictMode),
			enparsers.NewENTimeUnitLaterFormatParser(strictMode),
			enparsers.NewENCasualDateParser(),
			enparsers.NewENCasualTimeParser(),
			enparsers.NewENMonthNameParser(),
			enparsers.NewENRelativeDateFormatParser(),
			enparsers.NewENTimeUnitCasualRelativeFormatParser(true),
			enparsers.NewENYearParser(),          // Add year-only parser before compact format
			enparsers.NewENCompactFormatParser(), // Add compact format parser last as catch-all
		},
		Refiners: []kronos.Refiner{
			enrefiners.NewENMergeDateTimeRefiner(),
		},
	}

	// Apply common configuration
	config = includeCommonConfiguration(config, strictMode)

	// Add year/month/day parser at the beginning
	config.Parsers = append([]kronos.Parser{enparsers.NewENYearMonthDayParser(strictMode)}, config.Parsers...)

	// Add relative date refiners at the beginning
	config.Refiners = append([]kronos.Refiner{
		enrefiners.NewENMergeRelativeFollowByDateRefiner(),
		enrefiners.NewENMergeRelativeAfterDateRefiner(),
		refiners.NewOverlapRemovalRefiner(),
	}, config.Refiners...)

	// Re-apply date time refiner after timezone refinement
	config.Refiners = append(config.Refiners, enrefiners.NewENMergeDateTimeRefiner())

	// Extract year after merging date and time
	config.Refiners = append(config.Refiners, enrefiners.NewENExtractYearSuffixRefiner())

	// Keep date range refiner at the end
	config.Refiners = append(config.Refiners, enrefiners.NewENMergeDateRangeRefiner())

	return config
}

// Internal convenience functions for the builder API
func englishCasualChrono() *kronos.Chrono {
	return kronos.NewChrono(createCasualConfiguration(false))
}

func englishStrictChrono() *kronos.Chrono {
	return kronos.NewChrono(createConfiguration(true, false))
}

func englishGBChrono() *kronos.Chrono {
	return kronos.NewChrono(createCasualConfiguration(true))
}

// ParseSimple is a convenience function that parses text using the new builder API.
// It uses casual English configuration and the current time as reference.
//
// Example:
//
//	results, err := en.ParseSimple("tomorrow at 3pm")
func ParseSimple(text string) ([]kronos.Result, error) {
	results, err := New().Parse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse text: %w", err)
	}
	return results, nil
}

// ParseDateSimple is a convenience function that parses text and returns the first date.
// It uses casual English configuration and the current time as reference.
//
// Example:
//
//	date, err := en.ParseDateSimple("tomorrow at 3pm")
func ParseDateSimple(text string) (*time.Time, error) {
	date, err := New().ParseDate(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	return date, nil
}
