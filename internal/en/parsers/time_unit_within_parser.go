package parsers

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/parsing"
	"github.com/kljensen/kronos/internal/en/data"
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
		func(context *kronos.ParsingContext) *regexp.Regexp {
			timeUnitPattern := data.TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = data.TimeUnitNoAbbrPattern
			}

			// With optional prefix if forwardDate is enabled, required prefix otherwise
			// Pattern supports: "1 year 2 months", "1 year, 2 months", "1 year and 2 months", "1 year, 2 months and 3 days"
			// Separator between units: space, comma+space, or " and "
			unitSeparator := `(?:\s*,\s*|\s+and\s+|\s+)`
			// Word numbers pattern includes 1-19 and tens (20, 30, ..., 90)
			wordNumbers := `half|dozen|several|couple|few|ninety|eighty|seventy|sixty|fifty|forty|thirty|twenty|nineteen|eighteen|seventeen|sixteen|fifteen|fourteen|thirteen|twelve|eleven|ten|nine|eight|seven|six|five|four|three|two|one|a|an|the`
			var pattern string
			option := context.Option()
			if option.ForwardDate {
				pattern = `(?:(?:within|in|for)\s*)?` +
					`(?:(?:about|around|roughly|approximately|just)\s*(?:~\s*)?)?` +
					`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
					timeUnitPattern +
					`(?:` + unitSeparator + `(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
					timeUnitPattern + `)*)(?:\s|$|\b)`
			} else {
				pattern = `(?:within|in|for)\s*` +
					`(?:(?:about|around|roughly|approximately|just)\s*(?:~\s*)?)?` +
					`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
					timeUnitPattern +
					`(?:` + unitSeparator + `(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
					timeUnitPattern + `)*)(?:\s|$|\b)`
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

			duration := data.ParseDuration(match[1])
			if data.IsEmptyDuration(duration) {
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
