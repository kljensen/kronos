package kronos

import (
	"regexp"
	"strings"
	"time"
)

// Package helpers contains internal helper functions for the kronos package.
//
// NOTE ON APPARENT DUPLICATION:
// Many functions in this file have similar implementations in internal/helpers.
// This is intentional and necessary to avoid import cycles:
//   - internal/helpers provides functions for use by internal parsers/refiners
//   - This file provides functions for use by main package code
//   - Go's import cycle restrictions prevent consolidation into a single location
//
// The implementations are kept in sync manually. When modifying these functions,
// ensure corresponding changes are made in internal/helpers if applicable.

// ============================================================================
// Type casting helpers
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
// Date/time assignment and implication helpers
// ============================================================================

// assignSimilarDate assigns (force updates) the parsing components to the same day as the target.
// This sets year, month, and day as certain (known) values.
func assignSimilarDate(components *parsingComponents, date time.Time) {
	components.Assign(ComponentDay, date.Day())
	components.Assign(ComponentMonth, int(date.Month()))
	components.Assign(ComponentYear, date.Year())
}

// assignSimilarTime assigns (force updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as certain (known) values.
func assignSimilarTime(components *parsingComponents, date time.Time) {
	components.Assign(ComponentHour, date.Hour())
	components.Assign(ComponentMinute, date.Minute())
	components.Assign(ComponentSecond, date.Second())

	// Break down nanoseconds into milliseconds, microseconds, and nanoseconds
	totalNanos := date.Nanosecond()
	millisecond := totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond := remainingNanos / 1000
	nanosecond := remainingNanos % 1000

	components.Assign(ComponentMillisecond, millisecond)
	if microsecond > 0 {
		components.Assign(ComponentMicrosecond, microsecond)
	}
	if nanosecond > 0 {
		components.Assign(ComponentNanosecond, nanosecond)
	}

	// Set meridiem based on hour
	if date.Hour() < 12 {
		components.Assign(ComponentMeridiem, 0) // AM
	} else {
		components.Assign(ComponentMeridiem, 1) // PM
	}
}

// implySimilarDate implies (weakly updates) the parsing components to the same day as the target.
// This sets year, month, and day as implied values (only if not already certain).
func implySimilarDate(components *parsingComponents, date time.Time) {
	components.Imply(ComponentDay, date.Day())
	components.Imply(ComponentMonth, int(date.Month()))
	components.Imply(ComponentYear, date.Year())
}

// implySimilarTime implies (weakly updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as implied values (only if not already certain).
func implySimilarTime(components *parsingComponents, date time.Time) {
	components.Imply(ComponentHour, date.Hour())
	components.Imply(ComponentMinute, date.Minute())
	components.Imply(ComponentSecond, date.Second())

	// Break down nanoseconds into milliseconds, microseconds, and nanoseconds
	totalNanos := date.Nanosecond()
	millisecond := totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond := remainingNanos / 1000
	nanosecond := remainingNanos % 1000

	components.Imply(ComponentMillisecond, millisecond)
	if microsecond > 0 {
		components.Imply(ComponentMicrosecond, microsecond)
	}
	if nanosecond > 0 {
		components.Imply(ComponentNanosecond, nanosecond)
	}

	// Set meridiem based on hour
	if date.Hour() < 12 {
		components.Imply(ComponentMeridiem, 0) // AM
	} else {
		components.Imply(ComponentMeridiem, 1) // PM
	}
}
