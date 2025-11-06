package parsers

import (
	"github.com/kljensen/kronos/internal/parsing"
	"regexp"
	"strconv"

	kronos "github.com/kljensen/kronos"
)

// ISOFormatParser parses ISO 8601 date/time formats.
// Supports:
// - YYYY-MM-DD
// - YYYY-MM-DDThh:mm
// - YYYY-MM-DDThh:mm:ss
// - YYYY-MM-DDThh:mm:ss.sss
// - With optional timezone (Z, +hh:mm, -hh:mm)
type ISOFormatParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// ISO 8601 pattern
// http://www.w3.org/TR/NOTE-datetime
// Note: Go doesn't support lookahead (?=...), so we capture the trailing character instead
var isoPattern = regexp.MustCompile(
	`(?i)([0-9]{4})\-([0-9]{1,2})\-([0-9]{1,2})` +
		`(?:T` +
		`([0-9]{1,2}):([0-9]{1,2})` +
		`(?::([0-9]{1,2})(?:\.(\d{1,9}))?)?` +
		`(Z|([+-]\d{2}):?(\d{2})?)?` +
		`)?` +
		`(\W|$)`)

const (
	isoYearGroup        = 1
	isoMonthGroup       = 2
	isoDayGroup         = 3
	isoHourGroup        = 4
	isoMinuteGroup      = 5
	isoSecondGroup      = 6
	isoMillisecondGroup = 7
	isoTZDGroup         = 8
	isoTZDHourGroup     = 9
	isoTZDMinuteGroup   = 10
	isoTrailingGroup    = 11 // The trailing \W or $ character
)

// NewISOFormatParser creates a new ISO format parser.
func NewISOFormatParser() *ISOFormatParser {
	parser := &ISOFormatParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // Use default word boundary
	)

	return parser
}

func (p *ISOFormatParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	return isoPattern
}

func (p *ISOFormatParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	// Parse year, month, day
	year, err := strconv.Atoi(match[isoYearGroup])
	if err != nil {
		return nil
	}
	month, err := strconv.Atoi(match[isoMonthGroup])
	if err != nil {
		return nil
	}
	day, err := strconv.Atoi(match[isoDayGroup])
	if err != nil {
		return nil
	}

	components := context.CreateParsingComponents(map[kronos.Component]int{
		kronos.ComponentYear:  year,
		kronos.ComponentMonth: month,
		kronos.ComponentDay:   day,
	})

	// Trim the trailing character from the match text (captured by (\W|$))
	// We need to return an adjusted match that excludes the trailing character
	trailingChar := match[isoTrailingGroup]
	adjustedText := match[0]
	if len(trailingChar) > 0 && len(adjustedText) > 0 {
		// Remove the trailing character from the text
		adjustedText = adjustedText[:len(adjustedText)-len(trailingChar)]
	}

	// Parse time components if present
	if match[isoHourGroup] != "" {
		hour, err := strconv.Atoi(match[isoHourGroup])
		if err != nil {
			return nil
		}
		minute, err := strconv.Atoi(match[isoMinuteGroup])
		if err != nil {
			return nil
		}

		components.Assign(kronos.ComponentHour, hour)
		components.Assign(kronos.ComponentMinute, minute)

		// Parse seconds if present
		if match[isoSecondGroup] != "" {
			second, err := strconv.Atoi(match[isoSecondGroup])
			if err != nil {
				return nil
			}
			components.Assign(kronos.ComponentSecond, second)
		}

		// Parse fractional seconds if present (up to nanoseconds)
		if match[isoMillisecondGroup] != "" {
			fracStr := match[isoMillisecondGroup]
			// Pad or truncate to 9 digits (nanoseconds)
			for len(fracStr) < 9 {
				fracStr += "0"
			}
			if len(fracStr) > 9 {
				fracStr = fracStr[:9]
			}
			nanos, err := strconv.Atoi(fracStr)
			if err != nil {
				return nil
			}

			// Store as milliseconds, microseconds, and nanoseconds for compatibility
			millisecond := nanos / kronos.NanosecondsPerMS
			remainingNanos := nanos % kronos.NanosecondsPerMS
			microsecond := remainingNanos / kronos.NanosecondsPerMicro
			nanosecond := remainingNanos % kronos.NanosecondsPerMicro

			if millisecond > 0 {
				components.Assign(kronos.ComponentMillisecond, millisecond)
			}
			if microsecond > 0 {
				components.Assign(kronos.ComponentMicrosecond, microsecond)
			}
			if nanosecond > 0 {
				components.Assign(kronos.ComponentNanosecond, nanosecond)
			}
		}

		// Parse timezone if present
		if match[isoTZDGroup] != "" {
			offset := 0
			if match[isoTZDHourGroup] != "" {
				hourOffset, err := strconv.Atoi(match[isoTZDHourGroup])
				if err != nil {
					return nil
				}
				minuteOffset := 0
				if match[isoTZDMinuteGroup] != "" {
					minuteOffset, err = strconv.Atoi(match[isoTZDMinuteGroup])
					if err != nil {
						return nil
					}
				}

				offset = hourOffset * 60
				if offset < 0 {
					offset -= minuteOffset
				} else {
					offset += minuteOffset
				}
			}
			components.Assign(kronos.ComponentTimezoneOffset, offset)
		}
	}

	components.AddTag("parser/ISOFormatParser")

	// The parsing.AbstractParserWithWordBoundary will wrap this in ParsingResultWithBoundary
	// But we need to trim the trailing character from the text
	// Return a ParsingResultWithBoundary with the adjusted text
	return &kronos.ParsingResultWithBoundary{
		Components:         components,
		AdjustedText:       adjustedText,
		BoundaryLen:        0,    // Will be set by parsing.AbstractParserWithWordBoundary
		IncludeBoundaryIdx: true, // Will be overridden by parsing.AbstractParserWithWordBoundary
	}
}
