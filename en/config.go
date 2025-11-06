//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
	commonrefiners "github.com/kljensen/kronos/common/refiners"
	"github.com/kljensen/kronos/internal/en/parsers"
	"github.com/kljensen/kronos/internal/en/refiners"
)

// includeCommonConfiguration adds common parsers and refiners to a configuration.
func includeCommonConfiguration(config *kronos.Configuration, strictMode bool) *kronos.Configuration {
	// Add ISO format parser at the beginning
	config.Parsers = append([]kronos.Parser{common.NewISOFormatParser()}, config.Parsers...)

	// Add common refiners at the beginning
	config.Refiners = append([]kronos.Refiner{
		commonrefiners.NewMergeWeekdayComponentRefiner(),
		commonrefiners.NewExtractTimezoneOffsetRefiner(),
		commonrefiners.NewOverlapRemovalRefiner(),
	}, config.Refiners...)

	// Add common refiners at the end
	config.Refiners = append(config.Refiners,
		commonrefiners.NewExtractTimezoneAbbrRefiner(nil),
		commonrefiners.NewOverlapRemovalRefiner(),
		commonrefiners.NewDatePreferenceRefiner(), // Apply date preferences before ForwardDateRefiner
		commonrefiners.NewForwardDateRefiner(),    // ForwardDate option overrides preferences
		commonrefiners.NewUnlikelyFormatFilter(strictMode),
	)

	return config
}

// CreateCasualConfiguration creates a casual English configuration.
// This is now an alias for CreateConfiguration since casual parsers are included by default.
//
// Deprecated: Use en.New() to create a parser instance instead of accessing configurations directly.
// This function is maintained for backward compatibility.
func CreateCasualConfiguration(littleEndian bool) *kronos.Configuration {
	config := CreateConfiguration(false, littleEndian)

	// Add unlikely format filter for additional filtering
	config.Refiners = append(config.Refiners, refiners.NewENUnlikelyFormatFilter())

	return config
}

// CreateConfiguration creates a standard English configuration.
//
// Deprecated: Use en.New() to create a parser instance instead of accessing configurations directly.
// This function is maintained for backward compatibility.
func CreateConfiguration(strictMode, littleEndian bool) *kronos.Configuration {
	config := &kronos.Configuration{
		Parsers: []kronos.Parser{
			common.NewSlashDateFormatParser(littleEndian),
			parsers.NewENTimeUnitWithinFormatParser(strictMode),
			parsers.NewENMonthNameLittleEndianParser(),
			parsers.NewENMonthNameMiddleEndianParser(littleEndian), // shouldSkipYearLikeDate
			parsers.NewENWeekdayParser(),
			parsers.NewENSlashMonthFormatParser(),
			parsers.NewENTimeExpressionParser(strictMode),
			parsers.NewENTimeUnitAgoFormatParser(strictMode),
			parsers.NewENTimeUnitLaterFormatParser(strictMode),
			parsers.NewENCasualDateParser(),
			parsers.NewENCasualTimeParser(),
			parsers.NewENMonthNameParser(),
			parsers.NewENRelativeDateFormatParser(),
			parsers.NewENTimeUnitCasualRelativeFormatParser(true),
			parsers.NewENYearParser(),          // Add year-only parser before compact format
			parsers.NewENCompactFormatParser(), // Add compact format parser last as catch-all
		},
		Refiners: []kronos.Refiner{
			refiners.NewENMergeDateTimeRefiner(),
		},
	}

	// Apply common configuration
	config = includeCommonConfiguration(config, strictMode)

	// Add year/month/day parser at the beginning
	config.Parsers = append([]kronos.Parser{parsers.NewENYearMonthDayParser(strictMode)}, config.Parsers...)

	// Add relative date refiners at the beginning
	config.Refiners = append([]kronos.Refiner{
		refiners.NewENMergeRelativeFollowByDateRefiner(),
		refiners.NewENMergeRelativeAfterDateRefiner(),
		commonrefiners.NewOverlapRemovalRefiner(),
	}, config.Refiners...)

	// Re-apply date time refiner after timezone refinement
	config.Refiners = append(config.Refiners, refiners.NewENMergeDateTimeRefiner())

	// Extract year after merging date and time
	config.Refiners = append(config.Refiners, refiners.NewENExtractYearSuffixRefiner())

	// Keep date range refiner at the end
	config.Refiners = append(config.Refiners, refiners.NewENMergeDateRangeRefiner())

	return config
}
