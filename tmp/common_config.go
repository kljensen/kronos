package kronos

import (
	"github.com/markusmobius/go-chrono/common"
	"github.com/markusmobius/go-chrono/common/refiners"
)

// IncludeCommonConfiguration adds common parsers and refiners to a configuration.
func IncludeCommonConfiguration(config *Configuration, strictMode bool) *Configuration {
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
