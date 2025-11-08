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

// ENTimeUnitWithinFormatParser parses expressions like:
// "within 3 days", "in 2 hours", "for 5 minutes"
// Creates a date range from now to now + duration
type ENTimeUnitWithinFormatParser struct {
	*parsing.AbstractParserWithWordBoundary

	strictMode bool
}

// NewENTimeUnitWithinFormatParser creates a new parser for "within time" expressions
func NewENTimeUnitWithinFormatParser(strictMode bool) *ENTimeUnitWithinFormatParser {
	parser := &ENTimeUnitWithinFormatParser{
		strictMode: strictMode,
	}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		func(context *kronos.InternalParsingContext) *regexp.Regexp {
			timeUnitPattern := endata.TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = endata.TimeUnitNoAbbrPattern
			}

			var pattern string
			option := context.Option()
			if option.ForwardDate {
				pattern = `(?:(?:within|in|for)\s*)?`
			} else {
				pattern = `(?:within|in|for)\s*`
			}

			pattern += `(?:(?:about|around|roughly|approximately|just)\s*(?:~\s*)?)?` +
				`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern +
				`(?:` + unitSeparator + `(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern + `)*)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.InternalParsingContext, match []string) any {
			if len(match) < 2 {
				return nil
			}

			// Exclude "for the unit" phrases, e.g., "for the year"
			if strings.HasPrefix(strings.ToLower(match[0]), "for") &&
				regexp.MustCompile(`^for\s*the\s*\w+`).MatchString(strings.ToLower(match[0])) {
				return nil
			}

			duration := endata.ParseDuration(match[1])
			if endata.IsEmptyDuration(duration) {
				return nil
			}

			// Create relative result from reference (forward in time)
			components := helpers.CreateRelativeFromReference(context.Reference(), duration, internal.EmptyDuration)
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
