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
		func(context *kronos.ParsingContext) *regexp.Regexp {
			timeUnitPattern := endata.TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = endata.TimeUnitNoAbbrPattern
			}

			// Add optional approximation words at the beginning
			// Tilde is handled separately because it's a symbol, not a word
			approximationPattern := `(?:~\s*|(?:about|around|roughly|approximately|approx|circa)\s+)?`
			// Pattern supports: "1 year 2 months", "1 year, 2 months", "1 year and 2 months", "1 year, 2 months and 3 days"
			// Separator between units: space, comma+space, or " and "
			unitSeparator := `(?:\s*,\s*|\s+and\s+|\s+)`
			// Word numbers pattern includes 1-19 and tens (20, 30, ..., 90)
			wordNumbers := `half|dozen|several|couple|few|ninety|eighty|seventy|sixty|fifty|forty|thirty|twenty|nineteen|eighteen|seventeen|sixteen|fifteen|fourteen|thirteen|twelve|eleven|ten|nine|eight|seven|six|five|four|three|two|one|a|an|the`
			pattern := approximationPattern +
				`((?:(?:[0-9]+(?:[.,][0-9]+)?|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern +
				`(?:` + unitSeparator + `(?:(?:[0-9]+(?:[.,][0-9]+)?|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern + `)*)\s{0,5}(?:ago|before|earlier)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
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
			reversedDuration := helpers.ReverseDuration(duration)

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
