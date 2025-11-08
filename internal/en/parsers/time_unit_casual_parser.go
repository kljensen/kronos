//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal"
	endata "github.com/kljensen/kronos/internal/en/data"
	"github.com/kljensen/kronos/internal/helpers"
	"github.com/kljensen/kronos/internal/parsing"
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
		func(context *kronos.InternalParsingContext) *regexp.Regexp {
			timeUnitPattern := endata.TimeUnitPattern
			if !parser.allowAbbreviations {
				timeUnitPattern = endata.TimeUnitNoAbbrPattern
			}

			pattern := approximationPattern +
				`(this|last|past|next|after|\+|-)\s*` +
				`((?:(?:an?\s+)?(?:half|dozen|several|couple|few|one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|an|a|the|[0-9]+(?:[.,][0-9]+)?)\s*(?:an?\s+)?(?:of\s+)?)?` +
				timeUnitPattern +
				`(?:(?:\s*,?\s*)(?:(?:an?\s+)?(?:half|dozen|several|couple|few|one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|an|a|the|[0-9]+(?:[.,][0-9]+)?)\s*(?:an?\s+)?(?:of\s+)?)?` +
				timeUnitPattern + `)*)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.InternalParsingContext, match []string) any {
			if len(match) < 3 {
				return nil
			}

			// Check if approximation words were used by examining the full match
			fullMatch := match[0]
			_, isApproximate := helpers.StripApproximationWords(fullMatch)

			prefix := strings.ToLower(match[1])
			duration := endata.ParseDuration(match[2])
			if endata.IsEmptyDuration(duration) {
				return nil
			}

			// Reverse duration for "last", "past", and "-"
			switch prefix {
			case "last", "past", "-":
				duration = internal.ReverseDuration(duration)
			}

			// Create relative result from reference
			components := helpers.CreateRelativeFromReference(context.Reference(), duration, internal.EmptyDuration)
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
