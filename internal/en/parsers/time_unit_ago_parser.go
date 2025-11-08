//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal"
	endata "github.com/kljensen/kronos/internal/en/data"
	"github.com/kljensen/kronos/internal/helpers"
	"github.com/kljensen/kronos/internal/parsing"
)

// ENTimeUnitAgoFormatParser parses expressions like:
// "3 days ago", "2 hours ago", "5 minutes before", "15 minutes earlier"
type ENTimeUnitAgoFormatParser struct {
	*parsing.AbstractParserWithWordBoundary

	strictMode bool
}

// NewENTimeUnitAgoFormatParser creates a new parser for "time ago" expressions
func NewENTimeUnitAgoFormatParser(strictMode bool) *ENTimeUnitAgoFormatParser {
	parser := &ENTimeUnitAgoFormatParser{
		strictMode: strictMode,
	}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		func(context *kronos.InternalParsingContext) *regexp.Regexp {
			timeUnitPattern := endata.TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = endata.TimeUnitNoAbbrPattern
			}

			pattern := approximationPattern +
				`((?:(?:[0-9]+(?:[.,][0-9]+)?|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern +
				`(?:` + unitSeparator + `(?:(?:[0-9]+(?:[.,][0-9]+)?|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern + `)*)\s{0,5}(?:ago|before|earlier)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.InternalParsingContext, match []string) any {
			if len(match) < 2 {
				return nil
			}

			// Check if approximation words were used by examining the full match
			fullMatch := match[0]
			_, isApproximate := helpers.StripApproximationWords(fullMatch)

			duration := endata.ParseDuration(match[1])
			if endata.IsEmptyDuration(duration) {
				return nil
			}

			// Reverse the duration (go backwards in time)
			reversedDuration := internal.ReverseDuration(duration)

			// Create relative result from reference
			components := helpers.CreateRelativeFromReference(context.Reference(), reversedDuration, internal.EmptyDuration)
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
