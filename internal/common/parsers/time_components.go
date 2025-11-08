// Package parsers provides shared utilities and parsers for date/time parsing.
//
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// Time validation constants
const (
	maxHour24Format      = 24
	maxHour12Format      = 12
	maxMinute            = 60
	maxSecond            = 60
	hourCombinedFormat   = 100 // Format like "1430" for 14:30
	fourDigitYearLength  = 4
	maxNanoDigits        = 9
	singleDigitThreshold = 1
)

// ExtractPrimaryTimeComponents extracts time components from the primary match
func (p *AbstractTimeExpressionParser) ExtractPrimaryTimeComponents(
	context *kronos.InternalParsingContext,
	match []string,
	strict bool,
) *kronos.InternalParsingComponents {
	// Reject decimal ranges and embedded numbers early
	if looksLikeDecimalRange(context, match) {
		return nil
	}

	if isPartOfLargerNumber(context, match) {
		return nil
	}

	components := context.CreateParsingComponents(nil)
	minute := 0
	var meridiem *helpers.Meridiem

	// Parse hour
	hour, ok := helpers.AtoiSafe(match[TimeHourGroup])
	if !ok {
		return nil
	}

	// Handle combined hour-minute format (e.g., "1430" for 14:30)
	if hour > hourCombinedFormat {
		// When time is like '2019', it is more likely a year
		if len(match[TimeHourGroup]) == fourDigitYearLength && match[TimeMinuteGroup] == "" && match[TimeAMPMGroup] == "" {
			return nil
		}

		if p.strictMode || match[TimeMinuteGroup] != "" {
			return nil
		}

		minute = hour % hourCombinedFormat
		hour /= hourCombinedFormat
	}

	if hour > maxHour24Format {
		return nil
	}

	// Parse minute
	if match[TimeMinuteGroup] != "" {
		if len(match[TimeMinuteGroup]) == singleDigitThreshold && match[TimeAMPMGroup] == "" {
			// Skip single digit minute e.g., "at 1.1 xx"
			return nil
		}
		minute, ok = helpers.AtoiSafe(match[TimeMinuteGroup])
		if !ok {
			return nil
		}
	}

	if minute >= maxMinute {
		return nil
	}

	// Infer PM for 24-hour format
	if hour > maxHour12Format {
		m := helpers.MeridiemPM
		meridiem = &m
	}

	// Parse AM/PM
	if match[TimeAMPMGroup] != "" {
		var ok bool
		hour, meridiem, ok = parseMeridiem(match[TimeAMPMGroup], hour)
		if !ok {
			return nil
		}
	}

	components.Assign(kronos.ComponentHour, hour)
	components.Assign(kronos.ComponentMinute, minute)

	if meridiem != nil {
		components.Assign(kronos.ComponentMeridiem, int(*meridiem))
	} else {
		if hour < maxHour12Format {
			components.Imply(kronos.ComponentMeridiem, int(helpers.MeridiemAM))
		} else {
			components.Imply(kronos.ComponentMeridiem, int(helpers.MeridiemPM))
		}
	}

	// Parse seconds
	if match[TimeSecondGroup] != "" {
		second, ok := helpers.AtoiSafe(match[TimeSecondGroup])
		if !ok {
			return nil
		}
		if second >= maxSecond {
			return nil
		}
		components.Assign(kronos.ComponentSecond, second)
	}

	// Parse fractional seconds (up to nanoseconds)
	if match[TimeFractionalSecGroup] != "" {
		fracStr := match[TimeFractionalSecGroup]
		// Pad or truncate to 9 digits (nanoseconds)
		for len(fracStr) < maxNanoDigits {
			fracStr += "0"
		}
		if len(fracStr) > maxNanoDigits {
			fracStr = fracStr[:maxNanoDigits]
		}
		nanos, ok := helpers.AtoiSafe(fracStr)
		if !ok {
			return nil
		}

		// Store as milliseconds, microseconds, and nanoseconds for compatibility
		millisecond := nanos / 1000000
		remainingNanos := nanos % 1000000
		microsecond := remainingNanos / 1000
		nanosecond := remainingNanos % 1000

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

	// Call hook if provided
	if p.extractPrimaryTimeComponentsHook != nil {
		if !p.extractPrimaryTimeComponentsHook(context, match, components) {
			return nil
		}
	}

	// Set period to time-level since we're parsing time components
	components.SetPeriod(kronos.PeriodTime)

	return components
}

// ExtractFollowingTimeComponents extracts time components from the following match (for time ranges)
func (p *AbstractTimeExpressionParser) ExtractFollowingTimeComponents(
	context *kronos.InternalParsingContext,
	match []string,
	result *kronos.InternalParsingResult,
) *kronos.InternalParsingComponents {
	const noMeridiem = -1

	components := context.CreateParsingComponents(nil)
	resultStart, hasStart := helpers.AsParsingComponents(result.Start())

	// Parse seconds
	if match[TimeSecondGroup] != "" {
		second, ok := helpers.AtoiSafe(match[TimeSecondGroup])
		if !ok {
			return nil
		}
		if second >= maxSecond {
			return nil
		}
		components.Assign(kronos.ComponentSecond, second)
	}

	// Parse fractional seconds (up to nanoseconds)
	if match[TimeFractionalSecGroup] != "" {
		fracStr := match[TimeFractionalSecGroup]
		// Pad or truncate to 9 digits (nanoseconds)
		for len(fracStr) < maxNanoDigits {
			fracStr += "0"
		}
		if len(fracStr) > maxNanoDigits {
			fracStr = fracStr[:maxNanoDigits]
		}
		nanos, ok := helpers.AtoiSafe(fracStr)
		if !ok {
			return nil
		}

		// Store as milliseconds, microseconds, and nanoseconds for compatibility
		millisecond := nanos / 1000000
		remainingNanos := nanos % 1000000
		microsecond := remainingNanos / 1000
		nanosecond := remainingNanos % 1000

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

	// Check if this looks like part of a decimal range
	matchPos := strings.LastIndex(context.Text(), match[0])
	if matchPos > 0 {
		lookback := context.Text()[:matchPos]
		if regexp.MustCompile(`\d+\.\d+\s*[-–]\s*$`).MatchString(lookback) {
			return nil
		}
	}

	hour, ok := helpers.AtoiSafe(match[TimeHourGroup])
	if !ok {
		return nil
	}
	minute := 0
	meridiem := noMeridiem

	// Parse minute
	if match[TimeMinuteGroup] != "" {
		if len(match[TimeMinuteGroup]) == singleDigitThreshold && match[TimeAMPMGroup] == "" {
			// Skip single digit minute in following time e.g., "10 - 10.1"
			return nil
		}
		minute, ok = helpers.AtoiSafe(match[TimeMinuteGroup])
		if !ok {
			return nil
		}
	} else if hour > hourCombinedFormat {
		minute = hour % hourCombinedFormat
		hour /= hourCombinedFormat
	}

	if minute >= maxMinute || hour > maxHour24Format {
		return nil
	}

	if hour >= maxHour12Format {
		meridiem = int(helpers.MeridiemPM)
	}

	// Parse AM/PM
	if match[TimeAMPMGroup] != "" {
		if hour > maxHour12Format {
			return nil
		}

		ampm := strings.ToLower(string(match[TimeAMPMGroup][0]))
		if ampm == "a" {
			meridiem = int(helpers.MeridiemAM)
			if hour == maxHour12Format {
				hour = 0
				// Crossing midnight - advance day
				if hasStart && !resultStart.IsCertain(kronos.ComponentDay) {
					if dayVal := resultStart.Get(kronos.ComponentDay); dayVal != nil {
						components.Imply(kronos.ComponentDay, *dayVal+1)
					}
				}
			}
		}

		if ampm == "p" {
			meridiem = int(helpers.MeridiemPM)
			if hour != maxHour12Format {
				hour += maxHour12Format
			}
		}

		// Backfill meridiem to start time if not certain
		if hasStart && !resultStart.IsCertain(kronos.ComponentMeridiem) {
			if meridiem == int(helpers.MeridiemAM) {
				resultStart.Imply(kronos.ComponentMeridiem, int(helpers.MeridiemAM))
				if hourVal := resultStart.Get(kronos.ComponentHour); hourVal != nil && *hourVal == maxHour12Format {
					resultStart.Assign(kronos.ComponentHour, 0)
				}
			} else {
				resultStart.Imply(kronos.ComponentMeridiem, int(helpers.MeridiemPM))
				if hourVal := resultStart.Get(kronos.ComponentHour); hourVal != nil && *hourVal != maxHour12Format {
					resultStart.Assign(kronos.ComponentHour, *hourVal+maxHour12Format)
				}
			}
		}
	}

	components.Assign(kronos.ComponentHour, hour)
	components.Assign(kronos.ComponentMinute, minute)

	// Assign or imply meridiem
	if meridiem >= 0 {
		components.Assign(kronos.ComponentMeridiem, meridiem)
	} else {
		// Infer meridiem from start time if available
		startAtPM := false
		if hasStart && resultStart.IsCertain(kronos.ComponentMeridiem) {
			if startHourVal := resultStart.Get(kronos.ComponentHour); startHourVal != nil && *startHourVal > maxHour12Format {
				startAtPM = true
			}
		}

		switch {
		case startAtPM:
			if startHourVal := resultStart.Get(kronos.ComponentHour); startHourVal != nil && *startHourVal-maxHour12Format > hour {
				// e.g., "10pm - 1" means 1am next day
				components.Imply(kronos.ComponentMeridiem, int(helpers.MeridiemAM))
			} else if hour <= maxHour12Format {
				components.Assign(kronos.ComponentHour, hour+maxHour12Format)
				components.Assign(kronos.ComponentMeridiem, int(helpers.MeridiemPM))
			}
		case hour > maxHour12Format:
			components.Imply(kronos.ComponentMeridiem, int(helpers.MeridiemPM))
		case hour <= maxHour12Format:
			components.Imply(kronos.ComponentMeridiem, int(helpers.MeridiemAM))
		}
	}

	if components.Date().Before(result.Start().Date()) {
		dayVal := components.Get(kronos.ComponentDay)
		if dayVal != nil {
			components.Imply(kronos.ComponentDay, *dayVal+1)
		}
	}

	// Call hook if provided
	if p.extractFollowingTimeComponentsHook != nil {
		if !p.extractFollowingTimeComponentsHook(context, match, result, components) {
			return nil
		}
	}

	// Add the same parser tag as the start components
	startTags := result.Start().Tags()
	for tag := range startTags {
		components.AddTag(tag)
	}

	// Set period to time-level since we're parsing time components
	components.SetPeriod(kronos.PeriodTime)

	return components
}
