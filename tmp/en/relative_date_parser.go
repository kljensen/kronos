package en

import (
	"regexp"
	"strings"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENRelativeDateFormatParser parses expressions like:
// "next Tuesday", "last Friday", "this Monday", "this week", "next month"
type ENRelativeDateFormatParser struct {
	*common.AbstractParserWithWordBoundary
}

// TimeUnitRelativeDictionary maps time unit words to Timeunit values
var TimeUnitRelativeDictionary = map[string]kronos.Timeunit{
	"week":    kronos.TimeunitWeek,
	"weeks":   kronos.TimeunitWeek,
	"month":   kronos.TimeunitMonth,
	"months":  kronos.TimeunitMonth,
	"quarter": kronos.TimeunitQuarter,
	"year":    kronos.TimeunitYear,
	"years":   kronos.TimeunitYear,
}

// NewENRelativeDateFormatParser creates a new parser for relative date expressions
func NewENRelativeDateFormatParser() *ENRelativeDateFormatParser {
	parser := &ENRelativeDateFormatParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		func(context *kronos.ParsingContext) *regexp.Regexp {
			pattern := `(this|last|past|next|after\s*this)\s*(` + MatchAnyPattern(TimeUnitRelativeDictionary) + `)(?=\s*)?(?=\W|$)`
			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
			if len(match) < 3 {
				return nil
			}

			modifier := strings.ToLower(match[1])
			unitWord := strings.ToLower(match[2])
			timeunit, ok := TimeUnitRelativeDictionary[unitWord]
			if !ok {
				return nil
			}

			// Handle "next" and "after this"
			if modifier == "next" || strings.HasPrefix(modifier, "after") {
				duration := kronos.Duration{timeunit: 1}
				return kronos.CreateRelativeFromReference(context.Reference(), duration)
			}

			// Handle "last" and "past"
			if modifier == "last" || modifier == "past" {
				duration := kronos.Duration{timeunit: -1}
				return kronos.CreateRelativeFromReference(context.Reference(), duration)
			}

			// Handle "this" - set to beginning of current period
			components := context.CreateParsingComponents(nil)
			refDate := context.Reference().Instant()

			switch timeunit {
			case kronos.TimeunitWeek:
				// Start of this week (Sunday)
				date := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, refDate.Location())
				date = date.AddDate(0, 0, -int(date.Weekday()))
				components.Imply(kronos.ComponentDay, date.Day())
				components.Imply(kronos.ComponentMonth, int(date.Month()))
				components.Imply(kronos.ComponentYear, date.Year())

			case kronos.TimeunitMonth:
				// Start of this month
				date := time.Date(refDate.Year(), refDate.Month(), 1, 0, 0, 0, 0, refDate.Location())
				components.Imply(kronos.ComponentDay, date.Day())
				components.Assign(kronos.ComponentYear, date.Year())
				components.Assign(kronos.ComponentMonth, int(date.Month()))

			case kronos.TimeunitYear:
				// Start of this year
				date := time.Date(refDate.Year(), 1, 1, 0, 0, 0, 0, refDate.Location())
				components.Imply(kronos.ComponentDay, date.Day())
				components.Imply(kronos.ComponentMonth, int(date.Month()))
				components.Assign(kronos.ComponentYear, date.Year())
			}

			return components
		},
		nil,
	)

	return parser
}
