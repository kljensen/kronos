package en

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENTimeUnitCasualRelativeFormatParser parses expressions like:
// "this week", "next month", "last year", "past week", "+3 days", "-2 weeks"
type ENTimeUnitCasualRelativeFormatParser struct {
	*common.AbstractParserWithWordBoundary
	allowAbbreviations bool
}

// NewENTimeUnitCasualRelativeFormatParser creates a new parser for casual relative time expressions
func NewENTimeUnitCasualRelativeFormatParser(allowAbbreviations bool) *ENTimeUnitCasualRelativeFormatParser {
	parser := &ENTimeUnitCasualRelativeFormatParser{
		allowAbbreviations: allowAbbreviations,
	}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		func(context *kronos.ParsingContext) *regexp.Regexp {
			timeUnitPattern := TimeUnitPattern
			if !parser.allowAbbreviations {
				timeUnitPattern = TimeUnitNoAbbrPattern
			}

			pattern := `(this|last|past|next|after|\+|-)\s*` +
				`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
				timeUnitPattern +
				`(?:\s+(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
				timeUnitPattern + `)*)(?=\W|$)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
			if len(match) < 3 {
				return nil
			}

			prefix := strings.ToLower(match[1])
			duration := ParseDuration(match[2])
			if duration.IsEmpty() {
				return nil
			}

			// Reverse duration for "last", "past", and "-"
			switch prefix {
			case "last", "past", "-":
				duration = duration.Reverse()
			}

			// Create relative result from reference
			components := kronos.CreateRelativeFromReference(context.Reference, duration)
			if components != nil {
				components.AddTag("result/relativeDate")
				if duration.Hour != 0 || duration.Minute != 0 || duration.Second != 0 {
					components.AddTag("result/relativeDateAndTime")
				}
			}
			return components
		},
		nil,
	)

	return parser
}
