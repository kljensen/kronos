//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"regexp"
	"strconv"

	kronos "github.com/kljensen/kronos"
)

var timezoneOffsetPattern = regexp.MustCompile(`(?i)^\s*(?:\(?(?:GMT|UTC)\s?)?([+-])(\d{1,2})(?::?(\d{2}))?` + `\)?`)

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
	for i, result := range results {
		resultStart, ok := kronos.XAsParsingComponents(result.Start())
		if !ok {
			continue
		}
		if resultStart.IsCertain(kronos.ComponentTimezoneOffset) {
			continue
		}

		// Calculate the position after the result text
		suffixStart := result.Index() + len(result.Text())
		// Check if we're beyond the end of the text
		suffix, okSlice := kronos.XSafeSlice(context.Text(), suffixStart, len(context.Text()))
		if !okSlice {
			continue
		}
		match := timezoneOffsetPattern.FindStringSubmatch(suffix)
		if match == nil {
			continue
		}

		if context.Option().Debug != nil {
			context.Debug(func() {
				// Log: Extracting timezone offset
			})
		}

		hourOffset, err := strconv.Atoi(match[timezoneOffsetHourOffsetGroup])
		if err != nil {
			continue
		}
		minuteOffset := 0
		if match[timezoneOffsetMinuteOffsetGroup] != "" {
			minuteOffset, err = strconv.Atoi(match[timezoneOffsetMinuteOffsetGroup])
			if err != nil {
				continue
			}
		}

		timezoneOffset := hourOffset*60 + minuteOffset

		// No timezones have offsets greater than 14 hours, so disregard this match
		if timezoneOffset > 14*60 {
			continue
		}

		if match[timezoneOffsetSignGroup] == "-" {
			timezoneOffset = -timezoneOffset
		}

		var resultEnd *kronos.ParsingComponents
		if result.End() != nil {
			if endComponents, okEnd := kronos.XAsParsingComponents(result.End()); okEnd {
				resultEnd = endComponents
				resultEnd.Assign(kronos.ComponentTimezoneOffset, timezoneOffset)
			}
		}

		resultStart.Assign(kronos.ComponentTimezoneOffset, timezoneOffset)

		// Create new result with updated text
		// The match includes leading whitespace from the regex, so append it directly
		newText := result.Text() + match[0]
		newResult := context.CreateParsingResult(result.Index(), newText, resultStart, resultEnd)
		results[i] = newResult
	}

	return results
}
