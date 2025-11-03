package en

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENTimeUnitWithinFormatParser parses expressions like:
// "within 3 days", "in 2 hours", "for 5 minutes"
// Creates a date range from now to now + duration
type ENTimeUnitWithinFormatParser struct {
	*common.AbstractParserWithWordBoundary
	strictMode bool
}

// NewENTimeUnitWithinFormatParser creates a new parser for "within time" expressions
func NewENTimeUnitWithinFormatParser(strictMode bool) *ENTimeUnitWithinFormatParser {
	parser := &ENTimeUnitWithinFormatParser{
		strictMode: strictMode,
	}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		func(context *kronos.ParsingContext) *regexp.Regexp {
			timeUnitPattern := TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = TimeUnitNoAbbrPattern
			}

			// With optional prefix if forwardDate is enabled, required prefix otherwise
			var pattern string
			option := context.Option()
			if option.ForwardDate {
				pattern = `(?:(?:within|in|for)\s*)?` +
					`(?:(?:about|around|roughly|approximately|just)\s*(?:~\s*)?)?` +
					`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
					timeUnitPattern +
					`(?:\s+(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
					timeUnitPattern + `)*)(?=\W|$)`
			} else {
				pattern = `(?:within|in|for)\s*` +
					`(?:(?:about|around|roughly|approximately|just)\s*(?:~\s*)?)?` +
					`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
					timeUnitPattern +
					`(?:\s+(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
					timeUnitPattern + `)*)(?=\W|$)`
			}

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
			if len(match) < 2 {
				return nil
			}

			// Exclude "for the unit" phrases, e.g., "for the year"
			if strings.HasPrefix(strings.ToLower(match[0]), "for") &&
				regexp.MustCompile(`^for\s*the\s*\w+`).MatchString(strings.ToLower(match[0])) {
				return nil
			}

			duration := ParseDuration(match[1])
			if IsEmptyDuration(duration) {
				return nil
			}

			// Create relative result from reference (forward in time)
			components := kronos.CreateRelativeFromReference(context.Reference(), duration)
			if components != nil {
				components.AddTag("result/relativeDate")
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
