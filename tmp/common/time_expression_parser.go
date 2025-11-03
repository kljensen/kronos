package common

import (
	"regexp"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
)

// Time parsing capture group constants
const (
	TimeHourGroup       = 2
	TimeMinuteGroup     = 3
	TimeSecondGroup     = 4
	TimeMillisecondGroup = 5
	TimeAMPMGroup       = 6
)

// AbstractTimeExpressionParser is an abstract base for parsing time expressions.
// It supports 12-hour format (3pm, 3:30pm), 24-hour format (15:30),
// meridiem handling, and optional "at" keyword.
type AbstractTimeExpressionParser struct {
	// Required fields
	primaryPrefix   func() string
	followingPhase  func() string
	strictMode      bool

	// Optional customization
	primaryPatternLeftBoundary func() string
	primarySuffix              func() string
	followingSuffix            func() string
	patternFlags               func() string

	// Additional extraction logic
	extractPrimaryTimeComponentsHook  func(*kronos.ParsingContext, []string, *kronos.ParsingComponents) bool
	extractFollowingTimeComponentsHook func(*kronos.ParsingContext, []string, *kronos.ParsingResult, *kronos.ParsingComponents) bool
	checkAndReturnWithoutFollowingHook func(*kronos.ParsingResult) *kronos.ParsingResult
	checkAndReturnWithFollowingHook    func(*kronos.ParsingResult) *kronos.ParsingResult

	// Cached patterns
	cachedPrimaryPrefix       string
	cachedPrimarySuffix       string
	cachedPrimaryTimePattern  *regexp.Regexp
	cachedFollowingPhase      string
	cachedFollowingSuffix     string
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
			return `(?!/)(?=\W|$)`
		},
		followingSuffix: func() string {
			return `(?!/)(?=\W|$)`
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

// SetExtractPrimaryTimeComponentsHook allows additional extraction logic
func (p *AbstractTimeExpressionParser) SetExtractPrimaryTimeComponentsHook(
	fn func(*kronos.ParsingContext, []string, *kronos.ParsingComponents) bool,
) {
	p.extractPrimaryTimeComponentsHook = fn
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
		`(?:\.(\d{1,6}))?` + // Millisecond
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
		`(?:\.(\d{1,6}))?` + // Millisecond
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
	index := 0
	text := match[0][len(match[1]):]
	result := context.CreateParsingResult(index, text, startComponents, nil)

	// Look for following time pattern (for ranges like "10:00 - 21:45")
	// Try to find it in the remaining part of the match or immediately after
	// Note: Without context.Index, we search in context.Text() starting after match[0]
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
		return p.checkAndReturnWithoutFollowingPattern(result)
	}

	endComponents := p.ExtractFollowingTimeComponents(context, followingMatch, result)
	if endComponents != nil {
		// Create a new result with the extended text and end components
		newText := text + followingMatch[0]
		startComponents := result.Start().(*kronos.ParsingComponents)
		result = context.CreateParsingResult(index, newText, startComponents, endComponents)
	}

	return p.checkAndReturnWithFollowingPattern(result)
}

// ExtractPrimaryTimeComponents extracts time components from the primary match
func (p *AbstractTimeExpressionParser) ExtractPrimaryTimeComponents(
	context *kronos.ParsingContext,
	match []string,
	strict bool,
) *kronos.ParsingComponents {
	components := context.CreateParsingComponents(nil)
	minute := 0
	var meridiem *kronos.Meridiem

	// Parse hour
	hour, _ := strconv.Atoi(match[TimeHourGroup])
	if hour > 100 {
		// When time is like '2019', it is more likely a year
		if len(match[TimeHourGroup]) == 4 && match[TimeMinuteGroup] == "" && match[TimeAMPMGroup] == "" {
			return nil
		}

		if p.strictMode || match[TimeMinuteGroup] != "" {
			return nil
		}

		minute = hour % 100
		hour = hour / 100
	}

	if hour > 24 {
		return nil
	}

	// Parse minute
	if match[TimeMinuteGroup] != "" {
		if len(match[TimeMinuteGroup]) == 1 && match[TimeAMPMGroup] == "" {
			// Skip single digit minute e.g., "at 1.1 xx"
			return nil
		}
		minute, _ = strconv.Atoi(match[TimeMinuteGroup])
	}

	if minute >= 60 {
		return nil
	}

	if hour > 12 {
		m := kronos.MeridiemPM
		meridiem = &m
	}

	// Parse AM/PM
	if match[TimeAMPMGroup] != "" {
		if hour > 12 {
			return nil
		}
		ampm := strings.ToLower(string(match[TimeAMPMGroup][0]))
		if ampm == "a" {
			m := kronos.MeridiemAM
			meridiem = &m
			if hour == 12 {
				hour = 0
			}
		}
		if ampm == "p" {
			m := kronos.MeridiemPM
			meridiem = &m
			if hour != 12 {
				hour += 12
			}
		}
	}

	components.Assign(kronos.ComponentHour, hour)
	components.Assign(kronos.ComponentMinute, minute)

	if meridiem != nil {
		components.Assign(kronos.ComponentMeridiem, int(*meridiem))
	} else {
		if hour < 12 {
			components.Imply(kronos.ComponentMeridiem, int(kronos.MeridiemAM))
		} else {
			components.Imply(kronos.ComponentMeridiem, int(kronos.MeridiemPM))
		}
	}

	// Parse milliseconds
	if match[TimeMillisecondGroup] != "" {
		// Take first 3 digits
		msStr := match[TimeMillisecondGroup]
		if len(msStr) > 3 {
			msStr = msStr[:3]
		}
		millisecond, _ := strconv.Atoi(msStr)
		if millisecond >= 1000 {
			return nil
		}
		components.Assign(kronos.ComponentMillisecond, millisecond)
	}

	// Parse seconds
	if match[TimeSecondGroup] != "" {
		second, _ := strconv.Atoi(match[TimeSecondGroup])
		if second >= 60 {
			return nil
		}
		components.Assign(kronos.ComponentSecond, second)
	}

	// Call hook if provided
	if p.extractPrimaryTimeComponentsHook != nil {
		if !p.extractPrimaryTimeComponentsHook(context, match, components) {
			return nil
		}
	}

	return components
}

// ExtractFollowingTimeComponents extracts time components from the following match (for time ranges)
func (p *AbstractTimeExpressionParser) ExtractFollowingTimeComponents(
	context *kronos.ParsingContext,
	match []string,
	result *kronos.ParsingResult,
) *kronos.ParsingComponents {
	components := context.CreateParsingComponents(nil)

	// Parse milliseconds
	if match[TimeMillisecondGroup] != "" {
		msStr := match[TimeMillisecondGroup]
		if len(msStr) > 3 {
			msStr = msStr[:3]
		}
		millisecond, _ := strconv.Atoi(msStr)
		if millisecond >= 1000 {
			return nil
		}
		components.Assign(kronos.ComponentMillisecond, millisecond)
	}

	// Parse seconds
	if match[TimeSecondGroup] != "" {
		second, _ := strconv.Atoi(match[TimeSecondGroup])
		if second >= 60 {
			return nil
		}
		components.Assign(kronos.ComponentSecond, second)
	}

	hour, _ := strconv.Atoi(match[TimeHourGroup])
	minute := 0
	meridiem := -1

	// Parse minute
	if match[TimeMinuteGroup] != "" {
		minute, _ = strconv.Atoi(match[TimeMinuteGroup])
	} else if hour > 100 {
		minute = hour % 100
		hour = hour / 100
	}

	if minute >= 60 || hour > 24 {
		return nil
	}

	if hour >= 12 {
		meridiem = int(kronos.MeridiemPM)
	}

	// Parse AM/PM
	if match[TimeAMPMGroup] != "" {
		if hour > 12 {
			return nil
		}

		ampm := strings.ToLower(string(match[TimeAMPMGroup][0]))
		if ampm == "a" {
			meridiem = int(kronos.MeridiemAM)
			if hour == 12 {
				hour = 0
				resultStart := result.Start().(*kronos.ParsingComponents)
				if !resultStart.IsCertain(kronos.ComponentDay) {
					dayVal := resultStart.Get(kronos.ComponentDay)
					if dayVal != nil {
						components.Imply(kronos.ComponentDay, *dayVal+1)
					}
				}
			}
		}

		if ampm == "p" {
			meridiem = int(kronos.MeridiemPM)
			if hour != 12 {
				hour += 12
			}
		}

		resultStart := result.Start().(*kronos.ParsingComponents)
		if !resultStart.IsCertain(kronos.ComponentMeridiem) {
			if meridiem == int(kronos.MeridiemAM) {
				resultStart.Imply(kronos.ComponentMeridiem, int(kronos.MeridiemAM))
				hourVal := resultStart.Get(kronos.ComponentHour)
				if hourVal != nil && *hourVal == 12 {
					resultStart.Assign(kronos.ComponentHour, 0)
				}
			} else {
				resultStart.Imply(kronos.ComponentMeridiem, int(kronos.MeridiemPM))
				hourVal := resultStart.Get(kronos.ComponentHour)
				if hourVal != nil && *hourVal != 12 {
					resultStart.Assign(kronos.ComponentHour, *hourVal+12)
				}
			}
		}
	}

	components.Assign(kronos.ComponentHour, hour)
	components.Assign(kronos.ComponentMinute, minute)

	if meridiem >= 0 {
		components.Assign(kronos.ComponentMeridiem, meridiem)
	} else {
		resultStart := result.Start().(*kronos.ParsingComponents)
		startAtPM := resultStart.IsCertain(kronos.ComponentMeridiem)
		if startAtPM {
			startHourVal := resultStart.Get(kronos.ComponentHour)
			if startHourVal != nil && *startHourVal > 12 {
				startAtPM = true
			} else {
				startAtPM = false
			}
		}

		if startAtPM {
			startHourVal := resultStart.Get(kronos.ComponentHour)
			if startHourVal != nil && *startHourVal-12 > hour {
				// 10pm - 1 (am)
				components.Imply(kronos.ComponentMeridiem, int(kronos.MeridiemAM))
			} else if hour <= 12 {
				components.Assign(kronos.ComponentHour, hour+12)
				components.Assign(kronos.ComponentMeridiem, int(kronos.MeridiemPM))
			}
		} else if hour > 12 {
			components.Imply(kronos.ComponentMeridiem, int(kronos.MeridiemPM))
		} else if hour <= 12 {
			components.Imply(kronos.ComponentMeridiem, int(kronos.MeridiemAM))
		}
	}

	if components.Date().Before(result.Start().Date()) {
		dayVal := components.Get(kronos.ComponentDay)
		if dayVal != nil {
			components.Imply(kronos.ComponentDay, *dayVal+1)
		}
	}

	return components
}

func (p *AbstractTimeExpressionParser) checkAndReturnWithoutFollowingPattern(result *kronos.ParsingResult) *kronos.ParsingResult {
	// Use hook if provided
	if p.checkAndReturnWithoutFollowingHook != nil {
		return p.checkAndReturnWithoutFollowingHook(result)
	}

	// Single digit (e.g., "1") should not be counted as time expression
	if regexp.MustCompile(`^\d$`).MatchString(result.Text()) {
		return nil
	}

	// Three or more digits (e.g., "203", "2014") should not be counted as time expression
	if regexp.MustCompile(`^\d\d\d+$`).MatchString(result.Text()) {
		return nil
	}

	// Instead of "am/pm", it ends with "a" or "p" (e.g., "1a", "123p"), this seems unlikely
	if regexp.MustCompile(`\d[apAP]$`).MatchString(result.Text()) {
		return nil
	}

	// If it ends only with numbers or dots
	endingWithNumbers := regexp.MustCompile(`[^\d:.]([\d.]+)$`).FindStringSubmatch(result.Text())
	if endingWithNumbers != nil {
		endingNumbers := endingWithNumbers[1]

		// In strict mode (e.g., "at 1" or "at 1.2"), this should not be accepted
		if p.strictMode {
			return nil
		}

		// If it ends only with dot single digit, e.g., "at 1.2"
		if strings.Contains(endingNumbers, ".") && !regexp.MustCompile(`\d(\.\d{2})+$`).MatchString(endingNumbers) {
			return nil
		}

		// If it ends only with numbers above 24, e.g., "at 25"
		endingNumberVal, _ := strconv.Atoi(endingNumbers)
		if endingNumberVal > 24 {
			return nil
		}
	}

	return result
}

func (p *AbstractTimeExpressionParser) checkAndReturnWithFollowingPattern(result *kronos.ParsingResult) *kronos.ParsingResult {
	if regexp.MustCompile(`^\d+-\d+$`).MatchString(result.Text()) {
		return nil
	}

	// If it ends only with numbers or dots
	endingWithNumbers := regexp.MustCompile(`[^\d:.]([\d.]+)\s*-\s*([\d.]+)$`).FindStringSubmatch(result.Text())
	if endingWithNumbers != nil {
		// In strict mode (e.g., "at 1-3" or "at 1.2 - 2.3"), this should not be accepted
		if p.strictMode {
			return nil
		}

		startingNumbers := endingWithNumbers[1]
		endingNumbers := endingWithNumbers[2]

		// If it ends only with dot single digit, e.g., "at 1.2"
		if strings.Contains(endingNumbers, ".") && !regexp.MustCompile(`\d(\.\d{2})+$`).MatchString(endingNumbers) {
			return nil
		}

		// If it ends only with numbers above 24, e.g., "at 25"
		endingNumberVal, _ := strconv.Atoi(endingNumbers)
		startingNumberVal, _ := strconv.Atoi(startingNumbers)
		if endingNumberVal > 24 || startingNumberVal > 24 {
			return nil
		}
	}

	return result
}
