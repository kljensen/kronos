package en

import (
	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
	"github.com/kljensen/kronos/common/refiners"
	enrefiners "github.com/kljensen/kronos/en/refiners"
)

// includeCommonConfiguration adds common parsers and refiners to a configuration.
func includeCommonConfiguration(config *kronos.Configuration, strictMode bool) *kronos.Configuration {
	// Add ISO format parser at the beginning
	config.Parsers = append([]kronos.Parser{common.NewISOFormatParser()}, config.Parsers...)

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
		refiners.NewForwardDateRefiner(),     // ForwardDate option overrides preferences
		refiners.NewUnlikelyFormatFilter(strictMode),
	)

	return config
}

// CreateCasualConfiguration creates a casual English configuration.
// This is now an alias for CreateConfiguration since casual parsers are included by default.
func CreateCasualConfiguration(littleEndian bool) *kronos.Configuration {
	config := CreateConfiguration(false, littleEndian)

	// Add unlikely format filter for additional filtering
	config.Refiners = append(config.Refiners, enrefiners.NewENUnlikelyFormatFilter())

	return config
}

// CreateConfiguration creates a standard English configuration.
func CreateConfiguration(strictMode, littleEndian bool) *kronos.Configuration {
	config := &kronos.Configuration{
		Parsers: []kronos.Parser{
			common.NewSlashDateFormatParser(littleEndian),
			NewENTimeUnitWithinFormatParser(strictMode),
			NewENMonthNameLittleEndianParser(),
			NewENMonthNameMiddleEndianParser(littleEndian), // shouldSkipYearLikeDate
			NewENWeekdayParser(),
			NewENSlashMonthFormatParser(),
			NewENTimeExpressionParser(strictMode),
			NewENTimeUnitAgoFormatParser(strictMode),
			NewENTimeUnitLaterFormatParser(strictMode),
			NewENCasualDateParser(),
			NewENCasualTimeParser(),
			NewENMonthNameParser(),
			NewENRelativeDateFormatParser(),
			NewENTimeUnitCasualRelativeFormatParser(true),
			NewENCompactFormatParser(), // Add compact format parser last as catch-all
		},
		Refiners: []kronos.Refiner{
			enrefiners.NewENMergeDateTimeRefiner(),
		},
	}

	// Apply common configuration
	config = includeCommonConfiguration(config, strictMode)

	// Add year/month/day parser at the beginning
	config.Parsers = append([]kronos.Parser{NewENYearMonthDayParser(strictMode)}, config.Parsers...)

	// Add relative date refiners at the beginning
	config.Refiners = append([]kronos.Refiner{
		NewENMergeRelativeFollowByDateRefiner(),
		NewENMergeRelativeAfterDateRefiner(),
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
