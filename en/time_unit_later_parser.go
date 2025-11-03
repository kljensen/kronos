package en

import (
	"regexp"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENTimeUnitLaterFormatParser parses expressions like:
// "in 3 days", "3 hours later", "5 minutes from now", "2 weeks after"
type ENTimeUnitLaterFormatParser struct {
	*common.AbstractParserWithWordBoundary
	strictMode bool
}

// NewENTimeUnitLaterFormatParser creates a new parser for "time later" expressions
func NewENTimeUnitLaterFormatParser(strictMode bool) *ENTimeUnitLaterFormatParser {
	parser := &ENTimeUnitLaterFormatParser{
		strictMode: strictMode,
	}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		func(context *kronos.ParsingContext) *regexp.Regexp {
			timeUnitPattern := TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = TimeUnitNoAbbrPattern
			}

			// Add optional approximation words at the beginning
			// Tilde is handled separately because it's a symbol, not a word
			approximationPattern := `(?:~\s*|(?:about|around|roughly|approximately|approx|circa)\s+)?`
			pattern := approximationPattern +
				`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
				timeUnitPattern +
				`(?:\s+(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
				timeUnitPattern + `)*)\s{0,5}(?:later|after|from now|henceforth|forward|out)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
			if len(match) < 2 {
				return nil
			}

			// Check if approximation words were used by examining the full match
			fullMatch := match[0]
			_, isApproximate := kronos.StripApproximationWords(fullMatch)

			duration := ParseDuration(match[1])
			if IsEmptyDuration(duration) {
				return nil
			}

			// Create relative result from reference (forward in time)
			components := kronos.CreateRelativeFromReference(context.Reference(), duration)
			if components != nil {
				components.AddTag("result/relativeDate")
				// Add tag for relative date and time if time components are present
				if duration[kronos.TimeunitHour] != 0 || duration[kronos.TimeunitMinute] != 0 || duration[kronos.TimeunitSecond] != 0 {
					components.AddTag("result/relativeDateAndTime")
				}
				// Add approximation tag if approximation words were detected
				if isApproximate {
					components.AddTag("result/approximate")
				}
			}
			return components
		},
		nil,
	)

	return parser
}
