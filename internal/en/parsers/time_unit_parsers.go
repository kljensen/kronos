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
		func(context *kronos.InternalParsingContext, match []string) interface{} {
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

// ENTimeUnitLaterFormatParser parses expressions like:
// "in 3 days", "3 hours later", "5 minutes from now", "2 weeks after"
type ENTimeUnitLaterFormatParser struct {
	*parsing.AbstractParserWithWordBoundary
	strictMode bool
}

// NewENTimeUnitLaterFormatParser creates a new parser for "time later" expressions
func NewENTimeUnitLaterFormatParser(strictMode bool) *ENTimeUnitLaterFormatParser {
	parser := &ENTimeUnitLaterFormatParser{
		strictMode: strictMode,
	}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		func(context *kronos.InternalParsingContext) *regexp.Regexp {
			timeUnitPattern := endata.TimeUnitPattern
			if parser.strictMode {
				timeUnitPattern = endata.TimeUnitNoAbbrPattern
			}

			pattern := approximationPattern +
				`((?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern +
				`(?:` + unitSeparator + `(?:(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+|` + wordNumbers + `)\s*(?:an?\s+)?)?` +
				timeUnitPattern + `)*)\s{0,5}(?:later|after|from now|henceforth|forward|out)(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.InternalParsingContext, match []string) interface{} {
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

			// Create relative result from reference (forward in time)
			components := helpers.CreateRelativeFromReference(context.Reference(), duration, internal.EmptyDuration)
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
		func(context *kronos.InternalParsingContext, match []string) interface{} {
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
		func(context *kronos.InternalParsingContext, match []string) interface{} {
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
				duration = helpers.ReverseDuration(duration)
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
