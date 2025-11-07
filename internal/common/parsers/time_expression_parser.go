// Package common provides shared utilities and parsers for date/time parsing.
//
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// atoiSafe converts a string to an integer, returning 0 and false on error.
// Returns the value and true on success.
func atoiSafe(s string) (int, bool) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return val, true
}

// Time parsing capture group constants
const (
	TimeHourGroup          = 2
	TimeMinuteGroup        = 3
	TimeSecondGroup        = 4
	TimeFractionalSecGroup = 5
	TimeAMPMGroup          = 6
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

// AbstractTimeExpressionParser is an abstract base for parsing time expressions.
// It supports 12-hour format (3pm, 3:30pm), 24-hour format (15:30),
// meridiem handling, and optional "at" keyword.
type AbstractTimeExpressionParser struct {
	// Required fields
	primaryPrefix  func() string
	followingPhase func() string
	strictMode     bool

	// Optional customization
	primaryPatternLeftBoundary func() string
	primarySuffix              func() string
	followingSuffix            func() string
	patternFlags               func() string

	// Additional extraction logic
	extractPrimaryTimeComponentsHook   func(*kronos.ParsingContext, []string, *kronos.ParsingComponents) bool
	extractFollowingTimeComponentsHook func(*kronos.ParsingContext, []string, *kronos.ParsingResult, *kronos.ParsingComponents) bool
	checkAndReturnWithoutFollowingHook func(*kronos.ParsingResult) *kronos.ParsingResult

	// Cached patterns
	cachedPrimaryPrefix        string
	cachedPrimarySuffix        string
	cachedPrimaryTimePattern   *regexp.Regexp
	cachedFollowingPhase       string
	cachedFollowingSuffix      string
	cachedFollowingTimePattern *regexp.Regexp
}

// NewAbstractTimeExpressionParser creates a new time expression parser.
func NewAbstractTimeExpressionParser(
	primaryPrefix func() string,
	followingPhase func() string,
	strictMode bool,
) *AbstractTimeExpressionParser {
	return &AbstractTimeExpressionParser{
		primaryPrefix:  primaryPrefix,
		followingPhase: followingPhase,
		strictMode:     strictMode,
		primaryPatternLeftBoundary: func() string {
			return `(^|\s|T|\b)`
		},
		primarySuffix: func() string {
			// Go regexp doesn't support lookaheads (?! and ?=)
			// We use word boundary \b or whitespace/end
			return `(?:\s|$|\b)`
		},
		followingSuffix: func() string {
			// Go regexp doesn't support lookaheads (?! and ?=)
			// We use word boundary \b or whitespace/end
			return `(?:\s|$|\b)`
		},
		patternFlags: func() string {
			return "i"
		},
	}
}

// SetPrimarySuffix allows customization of the primary suffix pattern
func (p *AbstractTimeExpressionParser) SetPrimarySuffix(fn func() string) {
	p.primarySuffix = fn
}

// SetFollowingSuffix allows customization of the following suffix pattern
func (p *AbstractTimeExpressionParser) SetFollowingSuffix(fn func() string) {
	p.followingSuffix = fn
}

// SetExtractPrimaryTimeComponentsHook allows additional extraction logic
func (p *AbstractTimeExpressionParser) SetExtractPrimaryTimeComponentsHook(
	fn func(*kronos.ParsingContext, []string, *kronos.ParsingComponents) bool,
) {
	p.extractPrimaryTimeComponentsHook = fn
}

// SetExtractFollowingTimeComponentsHook allows additional extraction logic for following time
func (p *AbstractTimeExpressionParser) SetExtractFollowingTimeComponentsHook(
	fn func(*kronos.ParsingContext, []string, *kronos.ParsingResult, *kronos.ParsingComponents) bool,
) {
	p.extractFollowingTimeComponentsHook = fn
}

// SetCheckAndReturnWithoutFollowingHook allows customization of validation logic
func (p *AbstractTimeExpressionParser) SetCheckAndReturnWithoutFollowingHook(
	fn func(*kronos.ParsingResult) *kronos.ParsingResult,
) {
	p.checkAndReturnWithoutFollowingHook = fn
}

// buildPrimaryTimePattern constructs the primary time pattern
func buildPrimaryTimePattern(leftBoundary, primaryPrefix, primarySuffix, flags string) *regexp.Regexp {
	pattern := leftBoundary +
		primaryPrefix +
		`(\d{1,4})` + // Hour
		`(?:` +
		`(?:\.|:|：)` +
		`(\d{1,2})` + // Minute
		`(?:` +
		`(?::|：)` +
		`(\d{2})` + // Second
		`(?:\.(\d{1,9}))?` + // Fractional seconds (up to nanoseconds)
		`)?` +
		`)?` +
		`(?:\s*(a\.m\.|p\.m\.|am?|pm?))?` + // AM/PM
		primarySuffix

	flagPrefix := ""
	if flags == "i" {
		flagPrefix = "(?i)"
	}
	return regexp.MustCompile(flagPrefix + pattern)
}

// buildFollowingTimePattern constructs the following time pattern
func buildFollowingTimePattern(followingPhase, followingSuffix string) *regexp.Regexp {
	pattern := `^(` + followingPhase + `)` +
		`(\d{1,4})` + // Hour
		`(?:` +
		`(?:\.|:|：)` +
		`(\d{1,2})` + // Minute
		`(?:` +
		`(?:\.|:|：)` +
		`(\d{1,2})` + // Second
		`(?:\.(\d{1,9}))?` + // Fractional seconds (up to nanoseconds)
		`)?` +
		`)?` +
		`(?:\s*(a\.m\.|p\.m\.|am?|pm?))?` + // AM/PM
		followingSuffix

	return regexp.MustCompile("(?i)" + pattern)
}

// Pattern returns the primary time pattern
func (p *AbstractTimeExpressionParser) Pattern(context *kronos.ParsingContext) *regexp.Regexp {
	return p.getPrimaryTimePatternThroughCache()
}

func (p *AbstractTimeExpressionParser) getPrimaryTimePatternThroughCache() *regexp.Regexp {
	primaryPrefix := p.primaryPrefix()
	primarySuffix := p.primarySuffix()

	if p.cachedPrimaryPrefix == primaryPrefix && p.cachedPrimarySuffix == primarySuffix && p.cachedPrimaryTimePattern != nil {
		return p.cachedPrimaryTimePattern
	}

	p.cachedPrimaryTimePattern = buildPrimaryTimePattern(
		p.primaryPatternLeftBoundary(),
		primaryPrefix,
		primarySuffix,
		p.patternFlags(),
	)
	p.cachedPrimaryPrefix = primaryPrefix
	p.cachedPrimarySuffix = primarySuffix
	return p.cachedPrimaryTimePattern
}

func (p *AbstractTimeExpressionParser) getFollowingTimePatternThroughCache() *regexp.Regexp {
	followingPhase := p.followingPhase()
	followingSuffix := p.followingSuffix()

	if p.cachedFollowingPhase == followingPhase && p.cachedFollowingSuffix == followingSuffix && p.cachedFollowingTimePattern != nil {
		return p.cachedFollowingTimePattern
	}

	p.cachedFollowingTimePattern = buildFollowingTimePattern(followingPhase, followingSuffix)
	p.cachedFollowingPhase = followingPhase
	p.cachedFollowingSuffix = followingSuffix
	return p.cachedFollowingTimePattern
}

// Extract parses time expression from the match
func (p *AbstractTimeExpressionParser) Extract(context *kronos.ParsingContext, match []string) interface{} {
	startComponents := p.ExtractPrimaryTimeComponents(context, match, false)
	if startComponents == nil {
		// If the match seems like a year (e.g., "2013.12:..."),
		// skip the year part and try matching again
		if len(match[0]) >= 4 && regexp.MustCompile(`^\d{4}`).MatchString(match[0]) {
			// Return nil to skip this match
			return nil
		}
		return nil
	}

	// Calculate index and text
	// The match[1] is the left boundary capture group
	// The index should be relative to the match start (will be adjusted by chrono.go)
	// We remove the left boundary from the text
	index := len(match[1])
	text := match[0][len(match[1]):]
	// Don't trim yet - we may need to preserve spaces for range patterns
	result := context.CreateParsingResult(index, text, startComponents, nil)

	// Look for following time pattern (for ranges like "10:00 - 21:45")
	// Try to find it in the remaining part of the match or immediately after
	// Find the match position in the context text
	textIndex := strings.Index(context.Text(), match[0])
	if textIndex < 0 {
		return p.checkAndReturnWithoutFollowingPattern(result)
	}
	remainingText := context.Text()[textIndex+len(match[0]):]
	followingPattern := p.getFollowingTimePatternThroughCache()
	followingMatch := followingPattern.FindStringSubmatch(remainingText)

	// Pattern "456-12", "2022-12" should not be time without proper context
	if regexp.MustCompile(`^\d{3,4}`).MatchString(text) && followingMatch != nil {
		// e.g., "2022-12"
		if regexp.MustCompile(`^\s*([+-])\s*\d{2,4}$`).MatchString(followingMatch[0]) {
			return nil
		}
		// e.g., "2022-12:01..."
		if regexp.MustCompile(`^\s*([+-])\s*\d{2}\W\d{2}`).MatchString(followingMatch[0]) {
			return nil
		}
	}

	if followingMatch == nil || regexp.MustCompile(`^\s*([+-])\s*\d{3,4}$`).MatchString(followingMatch[0]) {
		// No following pattern - trim the text now
		result = context.CreateParsingResult(index, strings.TrimRight(text, " \t"), startComponents, nil)
		return p.checkAndReturnWithoutFollowingPattern(result)
	}

	// Check if the following match is followed by a slash (date pattern like "15/15")
	// If so, reject it as it's not a time range but part of a date expression
	followingMatchEnd := textIndex + len(match[0]) + len(followingMatch[0])
	if followingMatchEnd < len(context.Text()) && context.Text()[followingMatchEnd] == '/' {
		// No following pattern - trim the text now
		result = context.CreateParsingResult(index, strings.TrimRight(text, " \t"), startComponents, nil)
		return p.checkAndReturnWithoutFollowingPattern(result)
	}

	endComponents := p.ExtractFollowingTimeComponents(context, followingMatch, result)
	if endComponents != nil {
		// Create a new result with the extended text and end components
		// Trim trailing whitespace from the combined text
		newText := strings.TrimRight(text+followingMatch[0], " \t")
		result = context.CreateParsingResult(index, newText, startComponents, endComponents)
	}

	return p.checkAndReturnWithFollowingPattern(result)
}

// isDigit checks if a byte is a digit character (0-9).
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// looksLikeDecimalRange checks if the match appears to be part of a decimal range
// rather than a time expression. For example, "10.1 - 10.12" should not parse "10.12" as a time.
func looksLikeDecimalRange(context *kronos.ParsingContext, match []string) bool {
	matchStartIndex := strings.Index(context.Text(), match[0])
	if matchStartIndex < 0 {
		return false
	}

	// Look back for pattern like "X.Y - "
	if matchStartIndex > 0 && match[TimeMinuteGroup] != "" {
		lookback := context.Text()[:matchStartIndex]
		if regexp.MustCompile(`\d+\.\d+\s*[-–]\s*$`).MatchString(lookback) {
			return true
		}
	}

	// Look ahead for pattern like " - X.Y" where Y is single digit
	matchEndIndex := matchStartIndex + len(match[0])
	if matchEndIndex < len(context.Text()) && match[TimeMinuteGroup] == "" && match[TimeAMPMGroup] == "" {
		lookahead := context.Text()[matchEndIndex:]
		if regexp.MustCompile(`^\s*[-–]\s*\d+\.\d`).MatchString(lookahead) {
			return true
		}
	}

	return false
}

// isPartOfLargerNumber checks if match is embedded in a larger number sequence.
// For example, "2012-1400" should not match "12-1400" as a time, and "20" from "2020" should be rejected.
func isPartOfLargerNumber(context *kronos.ParsingContext, match []string) bool {
	matchStartInText := strings.Index(context.Text(), match[0])
	if matchStartInText < 0 {
		return false
	}

	// Check what precedes the match
	if matchStartInText > 0 {
		prevChar := context.Text()[matchStartInText-1]
		// Directly preceded by a digit (e.g., "2012-1400" matching "12-1400")
		if isDigit(prevChar) {
			return true
		}
		// Preceded by decimal point with digit (e.g., "10.1" matching "1")
		if prevChar == '.' && matchStartInText > 1 {
			if isDigit(context.Text()[matchStartInText-2]) {
				return true
			}
		}
	}

	// Check what follows (only for hour-only matches)
	if match[TimeMinuteGroup] == "" {
		matchEndInText := matchStartInText + len(match[0])
		if matchEndInText < len(context.Text()) {
			if isDigit(context.Text()[matchEndInText]) {
				return true
			}
		}
	}

	return false
}

// parseMeridiem parses AM/PM indicator and adjusts hour accordingly.
// Returns the adjusted hour, meridiem value, and success status.
// For AM: hour 12 becomes 0 (midnight). For PM: hours 1-11 add 12.
func parseMeridiem(ampmStr string, hour int) (int, *helpers.Meridiem, bool) {
	if hour > maxHour12Format {
		return 0, nil, false
	}

	ampm := strings.ToLower(string(ampmStr[0]))
	if ampm == "a" {
		m := helpers.MeridiemAM
		if hour == maxHour12Format {
			hour = 0
		}
		return hour, &m, true
	}

	if ampm == "p" {
		m := helpers.MeridiemPM
		if hour != maxHour12Format {
			hour += maxHour12Format
		}
		return hour, &m, true
	}

	return hour, nil, true
}

// ExtractPrimaryTimeComponents extracts time components from the primary match
func (p *AbstractTimeExpressionParser) ExtractPrimaryTimeComponents(
	context *kronos.ParsingContext,
	match []string,
	strict bool,
) *kronos.ParsingComponents {
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
	hour, ok := atoiSafe(match[TimeHourGroup])
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
		minute, ok = atoiSafe(match[TimeMinuteGroup])
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
		second, ok := atoiSafe(match[TimeSecondGroup])
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
		nanos, ok := atoiSafe(fracStr)
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
	context *kronos.ParsingContext,
	match []string,
	result *kronos.ParsingResult,
) *kronos.ParsingComponents {
	const noMeridiem = -1

	components := context.CreateParsingComponents(nil)
	resultStart, hasStart := helpers.AsParsingComponents(result.Start())

	// Parse seconds
	if match[TimeSecondGroup] != "" {
		second, ok := atoiSafe(match[TimeSecondGroup])
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
		nanos, ok := atoiSafe(fracStr)
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

	hour, ok := atoiSafe(match[TimeHourGroup])
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
		minute, ok = atoiSafe(match[TimeMinuteGroup])
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

// Validation patterns (compiled once for efficiency)
var (
	singleDigitPattern       = regexp.MustCompile(`^\d$`)
	threeOrMoreDigitsPattern = regexp.MustCompile(`^\d\d\d+$`)
	endsWithSingleAPPattern  = regexp.MustCompile(`\d[apAP]$`)
	plainNumberPattern       = regexp.MustCompile(`^\d+(?:\.\d+)?$`)
	endingNumbersPattern     = regexp.MustCompile(`[^\d:.]([\d.]+)$`)
	hasDigitPattern          = regexp.MustCompile(`\d`)
	twoDigitDecimalPattern   = regexp.MustCompile(`\d(\.\d{2})+$`)
)

func (p *AbstractTimeExpressionParser) checkAndReturnWithoutFollowingPattern(result *kronos.ParsingResult) *kronos.ParsingResult {
	// Use hook if provided
	if p.checkAndReturnWithoutFollowingHook != nil {
		return p.checkAndReturnWithoutFollowingHook(result)
	}

	text := strings.TrimSpace(result.Text())

	// Reject single digits (e.g., "1")
	if singleDigitPattern.MatchString(text) {
		return nil
	}

	// Reject three or more digits without separators (e.g., "203", "2014")
	if threeOrMoreDigitsPattern.MatchString(text) {
		return nil
	}

	// Reject single letter AM/PM suffix (e.g., "1a", "123p")
	if endsWithSingleAPPattern.MatchString(text) {
		return nil
	}

	// In strict mode, reject standalone numbers
	if p.strictMode && plainNumberPattern.MatchString(text) {
		return nil
	}

	// Validate trailing numbers after prefix (e.g., "at 10")
	if endingNumbers := endingNumbersPattern.FindStringSubmatch(text); endingNumbers != nil {
		if hasDigitPattern.MatchString(endingNumbers[1]) {
			nums := endingNumbers[1]

			if p.strictMode {
				return nil
			}

			// Reject single-digit decimals (e.g., "at 1.2")
			if strings.Contains(nums, ".") && !twoDigitDecimalPattern.MatchString(nums) {
				return nil
			}

			// Reject hours above 24
			if val, ok := atoiSafe(nums); ok && val > maxHour24Format {
				return nil
			}
		}
	}

	return result
}

var (
	startsWithDigitDashPattern = regexp.MustCompile(`^\d+\s*-`)
	pureNumberRangePattern     = regexp.MustCompile(`^\d+-\d+$`)
	numberRangePattern         = regexp.MustCompile(`[^\d:.]([\d.]+)\s*-\s*([\d.]+)$`)
)

func (p *AbstractTimeExpressionParser) checkAndReturnWithFollowingPattern(result *kronos.ParsingResult) *kronos.ParsingResult {
	text := strings.TrimSpace(result.Text())

	// In strict mode, reject simple number ranges (e.g., "7-730", "10 - 20")
	if p.strictMode && startsWithDigitDashPattern.MatchString(text) {
		return nil
	}

	// Reject pure number ranges without context (e.g., "12-14")
	if pureNumberRangePattern.MatchString(text) {
		return nil
	}

	// Validate number ranges after prefix (e.g., "at 10 - 20")
	if rangeMatch := numberRangePattern.FindStringSubmatch(text); rangeMatch != nil {
		if p.strictMode {
			return nil
		}

		startNum, endNum := rangeMatch[1], rangeMatch[2]

		// Reject single-digit decimals (e.g., "at 1.2 - 2.3")
		if strings.Contains(startNum, ".") && !twoDigitDecimalPattern.MatchString(startNum) {
			return nil
		}
		if strings.Contains(endNum, ".") && !twoDigitDecimalPattern.MatchString(endNum) {
			return nil
		}

		// Reject hours above 24
		if startVal, ok := atoiSafe(startNum); ok && startVal > maxHour24Format {
			return nil
		}
		if endVal, ok := atoiSafe(endNum); ok && endVal > maxHour24Format {
			return nil
		}
	}

	return result
}
