package en

import (
	"regexp"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENTimeUnitAgoFormatParser parses expressions like:
// "3 days ago", "2 hours ago", "5 minutes before", "15 minutes earlier"
type ENTimeUnitAgoFormatParser struct {
	*common.AbstractParserWithWordBoundary
	strictMode bool
}

// NewENTimeUnitAgoFormatParser creates a new parser for "time ago" expressions
func NewENTimeUnitAgoFormatParser(strictMode bool) *ENTimeUnitAgoFormatParser {
	parser := &ENTimeUnitAgoFormatParser{
		strictMode: strictMode,
	}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		func(context *kronos.ParsingContext) *regexp.Regexp {
			timeUnitPattern := TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = TimeUnitNoAbbrPattern
			}

			pattern := `((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
				timeUnitPattern +
				`(?:\s+(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
				timeUnitPattern + `)*)\s{0,5}(?:ago|before|earlier)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
			if len(match) < 2 {
				return nil
			}

			duration := ParseDuration(match[1])
			if IsEmptyDuration(duration) {
				return nil
			}

			// Reverse the duration (go backwards in time)
			reversedDuration := kronos.ReverseDuration(duration)

			// Create relative result from reference
			components := kronos.CreateRelativeFromReference(context.Reference(), reversedDuration)
			if components != nil {
				components.AddTag("result/relativeDate")
				// Add tag for relative date and time if time components are present
				if duration[kronos.TimeunitHour] != 0 || duration[kronos.TimeunitMinute] != 0 || duration[kronos.TimeunitSecond] != 0 {
					components.AddTag("result/relativeDateAndTime")
				}
			}
			return components
		},
		nil,
	)

	return parser
}
