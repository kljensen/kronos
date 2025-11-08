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

// Validation patterns (compiled once for efficiency)
var (
	singleDigitPattern         = regexp.MustCompile(`^\d$`)
	threeOrMoreDigitsPattern   = regexp.MustCompile(`^\d\d\d+$`)
	endsWithSingleAPPattern    = regexp.MustCompile(`\d[apAP]$`)
	plainNumberPattern         = regexp.MustCompile(`^\d+(?:\.\d+)?$`)
	endingNumbersPattern       = regexp.MustCompile(`[^\d:.]([\d.]+)$`)
	hasDigitPattern            = regexp.MustCompile(`\d`)
	twoDigitDecimalPattern     = regexp.MustCompile(`\d(\.\d{2})+$`)
	startsWithDigitDashPattern = regexp.MustCompile(`^\d+\s*-`)
	pureNumberRangePattern     = regexp.MustCompile(`^\d+-\d+$`)
	numberRangePattern         = regexp.MustCompile(`[^\d:.]([\d.]+)\s*-\s*([\d.]+)$`)
)

// isDigit checks if a byte is a digit character (0-9).
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// looksLikeDecimalRange checks if the match appears to be part of a decimal range
// rather than a time expression. For example, "10.1 - 10.12" should not parse "10.12" as a time.
func looksLikeDecimalRange(context *kronos.InternalParsingContext, match []string) bool {
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
func isPartOfLargerNumber(context *kronos.InternalParsingContext, match []string) bool {
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

func (p *AbstractTimeExpressionParser) checkAndReturnWithoutFollowingPattern(result *kronos.InternalParsingResult) *kronos.InternalParsingResult {
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
			if val, ok := helpers.AtoiSafe(nums); ok && val > maxHour24Format {
				return nil
			}
		}
	}

	return result
}

func (p *AbstractTimeExpressionParser) checkAndReturnWithFollowingPattern(result *kronos.InternalParsingResult) *kronos.InternalParsingResult {
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
		if startVal, ok := helpers.AtoiSafe(startNum); ok && startVal > maxHour24Format {
			return nil
		}
		if endVal, ok := helpers.AtoiSafe(endNum); ok && endVal > maxHour24Format {
			return nil
		}
	}

	return result
}
