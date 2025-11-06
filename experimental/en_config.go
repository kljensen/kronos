package experimental

import (
	"github.com/kljensen/kronos/internal/common/parsers"
	"github.com/kljensen/kronos/internal/common/refiners"
	enparsers "github.com/kljensen/kronos/internal/en/parsers"
	enrefiners "github.com/kljensen/kronos/internal/en/refiners"
)

// includeCommonConfiguration adds common parsers and refiners to a configuration.
func includeCommonConfiguration(config *Configuration, strictMode bool) *Configuration {
	// Add ISO format parser at the beginning
	config.Parsers = append([]Parser{parsers.NewISOFormatParser()}, config.Parsers...)

	// Add common refiners at the beginning
	config.Refiners = append([]Refiner{
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
func createCasualConfiguration(littleEndian bool) *Configuration {
	config := createConfiguration(false, littleEndian)

	// Add unlikely format filter for additional filtering
	config.Refiners = append(config.Refiners, enrefiners.NewENUnlikelyFormatFilter())

	return config
}

// createConfiguration creates a standard English configuration.
func createConfiguration(strictMode, littleEndian bool) *Configuration {
	config := &Configuration{
		Parsers: []Parser{
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
		Refiners: []Refiner{
			enrefiners.NewENMergeDateTimeRefiner(),
		},
	}

	// Apply common configuration
	config = includeCommonConfiguration(config, strictMode)

	// Add year/month/day parser at the beginning
	config.Parsers = append([]Parser{enparsers.NewENYearMonthDayParser(strictMode)}, config.Parsers...)

	// Add relative date refiners at the beginning
	config.Refiners = append([]Refiner{
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

// EnglishCasualChrono creates a Chrono instance configured for parsing casual English.
// It recognizes informal expressions like "today", "tomorrow", "next week", etc.
func EnglishCasualChrono() *Chrono {
	return NewChrono(createCasualConfiguration(false))
}

// EnglishStrictChrono creates a Chrono instance configured for parsing strict English.
// It only recognizes formal date/time patterns and avoids casual expressions.
func EnglishStrictChrono() *Chrono {
	return NewChrono(createConfiguration(true, false))
}

// EnglishGBChrono creates a Chrono instance configured for parsing UK-style English.
// It uses little-endian date format (day/month/year) and casual expressions.
func EnglishGBChrono() *Chrono {
	return NewChrono(createCasualConfiguration(true))
}
