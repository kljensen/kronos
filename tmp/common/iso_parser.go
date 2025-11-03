package common

import (
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
	*AbstractParserWithWordBoundary
}

// ISO 8601 pattern
// http://www.w3.org/TR/NOTE-datetime
// Note: Go doesn't support lookahead (?=...), so we capture the trailing character instead
var isoPattern = regexp.MustCompile(
	`(?i)([0-9]{4})\-([0-9]{1,2})\-([0-9]{1,2})` +
		`(?:T` +
		`([0-9]{1,2}):([0-9]{1,2})` +
		`(?::([0-9]{1,2})(?:\.(\d{1,4}))?)?` +
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

	parser.AbstractParserWithWordBoundary = NewAbstractParserWithWordBoundary(
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
	year, _ := strconv.Atoi(match[isoYearGroup])
	month, _ := strconv.Atoi(match[isoMonthGroup])
	day, _ := strconv.Atoi(match[isoDayGroup])

	components := context.CreateParsingComponents(map[kronos.Component]int{
		kronos.ComponentYear:  year,
		kronos.ComponentMonth: month,
		kronos.ComponentDay:   day,
	})

	// Parse time components if present
	if match[isoHourGroup] != "" {
		hour, _ := strconv.Atoi(match[isoHourGroup])
		minute, _ := strconv.Atoi(match[isoMinuteGroup])

		components.Assign(kronos.ComponentHour, hour)
		components.Assign(kronos.ComponentMinute, minute)

		// Parse seconds if present
		if match[isoSecondGroup] != "" {
			second, _ := strconv.Atoi(match[isoSecondGroup])
			components.Assign(kronos.ComponentSecond, second)
		}

		// Parse milliseconds if present
		if match[isoMillisecondGroup] != "" {
			// The millisecond group can be 1-4 digits, normalize to milliseconds
			msStr := match[isoMillisecondGroup]
			// Pad or truncate to 3 digits
			for len(msStr) < 3 {
				msStr += "0"
			}
			if len(msStr) > 3 {
				msStr = msStr[:3]
			}
			millisecond, _ := strconv.Atoi(msStr)
			components.Assign(kronos.ComponentMillisecond, millisecond)
		}

		// Parse timezone if present
		if match[isoTZDGroup] != "" {
			offset := 0
			if match[isoTZDHourGroup] != "" {
				hourOffset, _ := strconv.Atoi(match[isoTZDHourGroup])
				minuteOffset := 0
				if match[isoTZDMinuteGroup] != "" {
					minuteOffset, _ = strconv.Atoi(match[isoTZDMinuteGroup])
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

	return components.AddTag("parser/ISOFormatParser")
}
