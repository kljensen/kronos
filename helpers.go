package kronos

import (
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// String safety and slicing helpers
// ============================================================================

// asParsingComponents attempts to cast the given Components interface to
// a *parsingComponents. It returns the concrete value and true when the cast
// succeeds, or nil and false otherwise.
func asParsingComponents(pc Components) (*parsingComponents, bool) {
	if pc == nil {
		return nil, false
	}

	components, ok := pc.(*parsingComponents)
	if !ok || components == nil {
		return nil, false
	}

	return components, true
}

// safeSlice returns the substring of text between start (inclusive) and end
// (exclusive) while clamping the requested range to valid bounds. The returned
// boolean is false when start is beyond the end of the string, indicating that
// the requested slice could not be produced safely.
func safeSlice(text string, start, end int) (string, bool) {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	if start > len(text) {
		return "", false
	}
	if end > len(text) {
		end = len(text)
	}

	return text[start:end], true
}

// ============================================================================
// Input sanitization for Unicode normalization
// ============================================================================

// apostropheLookalikes contains Unicode characters that look like apostrophes
// but should be normalized to ASCII apostrophe (U+0027) for parsing.
var apostropheLookalikes = []rune{
	'\u2019', // RIGHT SINGLE QUOTATION MARK
	'\u02bc', // MODIFIER LETTER APOSTROPHE
	'\u02bb', // MODIFIER LETTER TURNED COMMA
	'\u055a', // ARMENIAN APOSTROPHE
	'\ua78c', // LATIN SMALL LETTER SALTILLO
	'\u2032', // PRIME
	'\u2035', // REVERSED PRIME
	'\u02b9', // MODIFIER LETTER PRIME
	'\uff07', // FULLWIDTH APOSTROPHE
}

// zeroWidthChars contains Unicode zero-width and formatting characters
// that should be removed from input text.
// Note: Directional marks (U+200E, U+200F) are handled separately as they
// often act as word separators and should be converted to spaces.
var zeroWidthChars = []rune{
	'\u200b', // ZERO WIDTH SPACE
	'\u200c', // ZERO WIDTH NON-JOINER
	'\u200d', // ZERO WIDTH JOINER
	'\ufeff', // ZERO WIDTH NO-BREAK SPACE (BOM)
	'\u00bb', // RIGHT-POINTING DOUBLE ANGLE QUOTATION MARK
	'\u00b7', // MIDDLE DOT
	'\u064e', // ARABIC FATHA
	'\u064f', // ARABIC DAMMA
}

// multiSpaceRegex matches one or more whitespace characters.
var multiSpaceRegex = regexp.MustCompile(`\s+`)

// sanitizeInput performs comprehensive Unicode normalization on input text
// to handle apostrophe variants, non-breaking spaces, zero-width characters,
// and other Unicode normalization issues that commonly occur in real-world input.
//
// The sanitization process:
// 1. Normalizes apostrophe lookalikes to ASCII apostrophe (')
// 2. Replaces non-breaking spaces (U+00A0) with regular spaces
// 3. Replaces ideographic spaces (U+3000) with regular spaces
// 4. Replaces directional marks with regular spaces (they often act as word separators)
// 5. Removes zero-width characters
// 6. Collapses multiple consecutive whitespace characters to single space
// 7. Trims leading and trailing whitespace
func sanitizeInput(input string) string {
	// 1. Normalize apostrophes to ASCII '
	input = normalizeApostrophes(input)

	// 2. Replace non-breaking space with regular space
	input = strings.ReplaceAll(input, "\u00a0", " ")

	// 3. Replace ideographic space (CJK full-width space) with regular space
	input = strings.ReplaceAll(input, "\u3000", " ")

	// 4. Replace directional marks with spaces (they often separate words)
	input = strings.ReplaceAll(input, "\u200e", " ") // LEFT-TO-RIGHT MARK
	input = strings.ReplaceAll(input, "\u200f", " ") // RIGHT-TO-LEFT MARK

	// 5. Remove zero-width characters
	input = removeZeroWidthChars(input)

	// 6. Collapse multiple spaces to single space
	input = multiSpaceRegex.ReplaceAllString(input, " ")

	// 7. Trim leading/trailing whitespace
	input = strings.TrimSpace(input)

	return input
}

// normalizeApostrophes replaces all apostrophe lookalike Unicode characters
// with the ASCII apostrophe (U+0027).
func normalizeApostrophes(s string) string {
	for _, lookalike := range apostropheLookalikes {
		s = strings.ReplaceAll(s, string(lookalike), "'")
	}
	return s
}

// removeZeroWidthChars removes zero-width and directional formatting characters
// that are invisible but can break parsing.
func removeZeroWidthChars(s string) string {
	for _, zwChar := range zeroWidthChars {
		s = strings.ReplaceAll(s, string(zwChar), "")
	}
	return s
}

// ============================================================================
// Weekday calculation helpers
// ============================================================================

// getNextWeekday returns the date of the next occurrence of the target weekday
// after the reference date (not including the reference date itself).
//
// For example, if refDate is Monday and targetWeekday is Monday, it returns
// the following Monday (7 days later).
func getNextWeekday(refDate time.Time, targetWeekday time.Weekday) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	// Calculate days forward to next occurrence
	daysForward := target - refWeekday
	if daysForward <= 0 {
		daysForward += 7
	}

	return refDate.AddDate(0, 0, daysForward)
}

// getLastWeekday returns the date of the last occurrence of the target weekday
// before the reference date (not including the reference date itself).
//
// For example, if refDate is Monday and targetWeekday is Monday, it returns
// the previous Monday (7 days earlier).
func getLastWeekday(refDate time.Time, targetWeekday time.Weekday) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	// Calculate days backward to last occurrence
	daysBackward := target - refWeekday
	if daysBackward >= 0 {
		daysBackward -= 7
	}

	return refDate.AddDate(0, 0, daysBackward)
}

// getThisWeekday returns the date of the target weekday in "this" week.
//
// The forward parameter controls the direction:
// - If forward=true: looks forward from refDate (inclusive)
// - If forward=false: looks backward from refDate (inclusive)
//
// Examples with refDate = Wednesday:
// - getThisWeekday(Wednesday, Monday, forward=true) -> next Monday (5 days forward)
// - getThisWeekday(Wednesday, Monday, forward=false) -> previous Monday (2 days back)
// - getThisWeekday(Wednesday, Wednesday, forward=true) -> same Wednesday (0 days)
// - getThisWeekday(Wednesday, Wednesday, forward=false) -> same Wednesday (0 days)
func getThisWeekday(refDate time.Time, targetWeekday time.Weekday, forward bool) time.Time {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	if refWeekday == target {
		// Same weekday - return the reference date
		return refDate
	}

	if forward {
		// Look forward
		daysForward := target - refWeekday
		if daysForward < 0 {
			daysForward += 7
		}
		return refDate.AddDate(0, 0, daysForward)
	} else {
		// Look backward
		daysBackward := target - refWeekday
		if daysBackward > 0 {
			daysBackward -= 7
		}
		return refDate.AddDate(0, 0, daysBackward)
	}
}

// getDaysToWeekday returns the number of days from refDate to the target weekday
// based on the modifier.
//
// Modifiers:
// - "this": Returns the target weekday in the current week (looking forward)
// - "next": Returns the target weekday in the next occurrence
// - "last": Returns the target weekday in the last occurrence (negative value)
// - nil/empty: Returns the closest weekday (forward or backward)
func getDaysToWeekday(refDate time.Time, targetWeekday time.Weekday, modifier *string) int {
	refWeekday := time.Weekday(refDate.Weekday())

	if modifier == nil {
		return getDaysToWeekdayClosest(refDate, targetWeekday)
	}

	switch *modifier {
	case "this":
		return getDaysForwardToWeekday(refDate, targetWeekday)
	case "last":
		return getBackwardDaysToWeekday(refDate, targetWeekday)
	case "next":
		// Special handling for "next" based on the current weekday
		if refWeekday == time.Sunday {
			if targetWeekday == time.Sunday {
				return 7
			}
			return int(targetWeekday)
		}

		if refWeekday == time.Saturday {
			if targetWeekday == time.Saturday {
				return 7
			}
			if targetWeekday == time.Sunday {
				return 8
			}
			return 1 + int(targetWeekday)
		}

		// For weekdays (Mon-Fri)
		if targetWeekday < refWeekday && targetWeekday != time.Sunday {
			return getDaysForwardToWeekday(refDate, targetWeekday)
		}
		return getDaysForwardToWeekday(refDate, targetWeekday) + 7
	}

	return getDaysToWeekdayClosest(refDate, targetWeekday)
}

// getDaysForwardToWeekday returns the number of days forward to reach the target weekday.
// Returns 0 if already on the target weekday.
func getDaysForwardToWeekday(refDate time.Time, targetWeekday time.Weekday) int {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	forwardCount := target - refWeekday
	if forwardCount < 0 {
		forwardCount += 7
	}

	return forwardCount
}

// getBackwardDaysToWeekday returns the number of days backward to reach the target weekday.
// Returns a negative number (or 0 if already on the target weekday).
func getBackwardDaysToWeekday(refDate time.Time, targetWeekday time.Weekday) int {
	refWeekday := int(refDate.Weekday())
	target := int(targetWeekday)

	backwardCount := target - refWeekday
	if backwardCount >= 0 {
		backwardCount -= 7
	}

	return backwardCount
}

// getDaysToWeekdayClosest returns the number of days to the closest occurrence
// of the target weekday (either forward or backward).
func getDaysToWeekdayClosest(refDate time.Time, targetWeekday time.Weekday) int {
	backward := getBackwardDaysToWeekday(refDate, targetWeekday)
	forward := getDaysForwardToWeekday(refDate, targetWeekday)

	// Choose the direction with fewer days
	if forward < -backward {
		return forward
	}
	return backward
}
