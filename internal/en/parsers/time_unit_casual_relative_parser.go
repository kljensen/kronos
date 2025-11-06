package parsers

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/parsing"
	"github.com/kljensen/kronos/internal/en/data"
)

// ENTimeUnitCasualRelativeFormatParser parses expressions like:
// "this week", "next month", "last year", "past week", "+3 days", "-2 weeks"
type ENTimeUnitCasualRelativeFormatParser struct {
	*parsing.AbstractParserWithWordBoundary
	allowAbbreviations bool
}

// NewENTimeUnitCasualRelativeFormatParser creates a new parser for casual relative time expressions
func NewENTimeUnitCasualRelativeFormatParser(allowAbbreviations bool) *ENTimeUnitCasualRelativeFormatParser {
	parser := &ENTimeUnitCasualRelativeFormatParser{
		allowAbbreviations: allowAbbreviations,
	}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		func(context *kronos.ParsingContext) *regexp.Regexp {
			timeUnitPattern := data.TimeUnitPattern
			if !parser.allowAbbreviations {
				timeUnitPattern = data.TimeUnitNoAbbrPattern
			}

			// Add optional approximation words at the beginning
			// Tilde is handled separately because it's a symbol, not a word
			approximationPattern := `(?:~\s*|(?:about|around|roughly|approximately|approx|circa)\s+)?`
			pattern := approximationPattern +
				`(this|last|past|next|after|\+|-)\s*` +
				`((?:(?:an?\s+)?(?:half|dozen|several|couple|few|one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|an|a|the|[0-9]+(?:[.,][0-9]+)?)\s*(?:an?\s+)?(?:of\s+)?)?` +
				timeUnitPattern +
				`(?:(?:\s*,?\s*)(?:(?:an?\s+)?(?:half|dozen|several|couple|few|one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|an|a|the|[0-9]+(?:[.,][0-9]+)?)\s*(?:an?\s+)?(?:of\s+)?)?` +
				timeUnitPattern + `)*)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
			if len(match) < 3 {
				return nil
			}

			// Check if approximation words were used by examining the full match
			fullMatch := match[0]
			_, isApproximate := kronos.StripApproximationWords(fullMatch)

			prefix := strings.ToLower(match[1])
			duration := data.ParseDuration(match[2])
			if data.IsEmptyDuration(duration) {
				return nil
			}

			// Reverse duration for "last", "past", and "-"
			switch prefix {
			case "last", "past", "-":
				duration = kronos.ReverseDuration(duration)
			}

			// Create relative result from reference
			components := kronos.CreateRelativeFromReference(context.Reference(), duration)
			if components != nil {
				components.AddTag("result/relativeDate")
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
