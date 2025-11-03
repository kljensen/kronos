package en

import (
	. "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
	"github.com/kljensen/kronos/common/refiners"
	enrefiners "github.com/kljensen/kronos/en/refiners"
)

// includeCommonConfiguration adds common parsers and refiners to a configuration.
func includeCommonConfiguration(config *Configuration, strictMode bool) *Configuration {
	// Add ISO format parser at the beginning
	config.Parsers = append([]Parser{common.NewISOFormatParser()}, config.Parsers...)

	// Add common refiners at the beginning
	config.Refiners = append([]Refiner{
		refiners.NewMergeWeekdayComponentRefiner(),
		refiners.NewExtractTimezoneOffsetRefiner(),
		refiners.NewOverlapRemovalRefiner(),
	}, config.Refiners...)

	// Add common refiners at the end
	config.Refiners = append(config.Refiners,
		refiners.NewExtractTimezoneAbbrRefiner(),
		refiners.NewOverlapRemovalRefiner(),
		refiners.NewForwardDateRefiner(),
		refiners.NewUnlikelyFormatFilter(strictMode),
	)

	return config
}

// CreateCasualConfiguration creates a casual English configuration.
// This includes parsers for casual language like "today", "tomorrow", etc.
func CreateCasualConfiguration(littleEndian bool) *Configuration {
	config := CreateConfiguration(false, littleEndian)

	// Add casual parsers
	config.Parsers = append(config.Parsers,
		NewENCasualDateParser(),
		NewENCasualTimeParser(),
		NewENMonthNameParser(),
		NewENRelativeDateFormatParser(),
		NewENTimeUnitCasualRelativeFormatParser(true),
	)

	// Add unlikely format filter
	config.Refiners = append(config.Refiners, enrefiners.NewENUnlikelyFormatFilter())

	return config
}

// CreateConfiguration creates a standard English configuration.
func CreateConfiguration(strictMode, littleEndian bool) *Configuration {
	config := &Configuration{
		Parsers: []Parser{
			common.NewSlashDateFormatParser(littleEndian),
			NewENTimeUnitWithinFormatParser(strictMode),
			NewENMonthNameLittleEndianParser(),
			NewENMonthNameMiddleEndianParser(littleEndian), // shouldSkipYearLikeDate
			NewENWeekdayParser(),
			NewENSlashMonthFormatParser(),
			NewENTimeExpressionParser(strictMode),
			NewENTimeUnitAgoFormatParser(strictMode),
			NewENTimeUnitLaterFormatParser(strictMode),
		},
		Refiners: []Refiner{
			enrefiners.NewENMergeDateTimeRefiner(),
		},
	}

	// Apply common configuration
	config = includeCommonConfiguration(config, strictMode)

	// Add year/month/day parser at the beginning
	config.Parsers = append([]Parser{NewENYearMonthDayParser(strictMode)}, config.Parsers...)

	// Add relative date refiners at the beginning
	config.Refiners = append([]Refiner{
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
