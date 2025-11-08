// Package parsers provides shared utilities and parsers for date/time parsing.
//
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
)

// Time parsing capture group constants
const (
	TimeHourGroup          = 2
	TimeMinuteGroup        = 3
	TimeSecondGroup        = 4
	TimeFractionalSecGroup = 5
	TimeAMPMGroup          = 6
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
	extractPrimaryTimeComponentsHook   func(*kronos.InternalParsingContext, []string, *kronos.InternalParsingComponents) bool
	extractFollowingTimeComponentsHook func(*kronos.InternalParsingContext, []string, *kronos.InternalParsingResult, *kronos.InternalParsingComponents) bool
	checkAndReturnWithoutFollowingHook func(*kronos.InternalParsingResult) *kronos.InternalParsingResult

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
	fn func(*kronos.InternalParsingContext, []string, *kronos.InternalParsingComponents) bool,
) {
	p.extractPrimaryTimeComponentsHook = fn
}

// SetExtractFollowingTimeComponentsHook allows additional extraction logic for following time
func (p *AbstractTimeExpressionParser) SetExtractFollowingTimeComponentsHook(
	fn func(*kronos.InternalParsingContext, []string, *kronos.InternalParsingResult, *kronos.InternalParsingComponents) bool,
) {
	p.extractFollowingTimeComponentsHook = fn
}

// SetCheckAndReturnWithoutFollowingHook allows customization of validation logic
func (p *AbstractTimeExpressionParser) SetCheckAndReturnWithoutFollowingHook(
	fn func(*kronos.InternalParsingResult) *kronos.InternalParsingResult,
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
func (p *AbstractTimeExpressionParser) Pattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	return p.getPrimaryTimePatternThroughCache()
}

// Extract parses time expression from the match
func (p *AbstractTimeExpressionParser) Extract(context *kronos.InternalParsingContext, match []string) any {
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
