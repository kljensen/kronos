package refiners

import (
	"regexp"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
)

var (
	timezoneOffsetPattern = regexp.MustCompile(`(?i)^\s*(?:\(?(?:GMT|UTC)\s?)?([+-])(\d{1,2})(?::?(\d{2}))?` + `\)?`)
)

const (
	timezoneOffsetSignGroup         = 1
	timezoneOffsetHourOffsetGroup   = 2
	timezoneOffsetMinuteOffsetGroup = 3
)

// ExtractTimezoneOffsetRefiner extracts timezone offsets from text following a parsed result.
// Examples: "+0900", "-05:00", "GMT+8", "UTC-5"
type ExtractTimezoneOffsetRefiner struct{}

// NewExtractTimezoneOffsetRefiner creates a new ExtractTimezoneOffsetRefiner
func NewExtractTimezoneOffsetRefiner() *ExtractTimezoneOffsetRefiner {
	return &ExtractTimezoneOffsetRefiner{}
}

// Refine extracts timezone offsets and adds them to results
func (r *ExtractTimezoneOffsetRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	for _, result := range results {
		if result.Start.IsCertain(kronos.ComponentTimezoneOffset) {
			continue
		}

		suffix := context.Text[result.Index+len(result.Text):]
		match := timezoneOffsetPattern.FindStringSubmatch(suffix)
		if match == nil {
			continue
		}

		if context.Option.Debug {
			context.DebugLog("Extracting timezone: '%s' into: %s", match[0], result)
		}

		hourOffset, _ := strconv.Atoi(match[timezoneOffsetHourOffsetGroup])
		minuteOffset := 0
		if match[timezoneOffsetMinuteOffsetGroup] != "" {
			minuteOffset, _ = strconv.Atoi(match[timezoneOffsetMinuteOffsetGroup])
		}

		timezoneOffset := hourOffset*60 + minuteOffset

		// No timezones have offsets greater than 14 hours, so disregard this match
		if timezoneOffset > 14*60 {
			continue
		}

		if match[timezoneOffsetSignGroup] == "-" {
			timezoneOffset = -timezoneOffset
		}

		if result.End != nil {
			result.End.Assign(kronos.ComponentTimezoneOffset, timezoneOffset)
		}

		result.Start.Assign(kronos.ComponentTimezoneOffset, timezoneOffset)
		result.Text += match[0]
	}

	return results
}
