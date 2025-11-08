package data

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// ParseOrdinalNumber parses an ordinal number pattern (e.g., "1st", "2nd", "first", "second")
func ParseOrdinalNumber(match string) int {
	lower := strings.ToLower(match)

	// Check if it's a word
	if val, ok := OrdinalWordDictionary[lower]; ok {
		return val
	}

	// Remove ordinal suffix (st, nd, rd, th)
	cleaned := regexp.MustCompile(`(?i)(st|nd|rd|th)$`).ReplaceAllString(lower, "")
	val, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0
	}
	return val
}

// ParseYear parses a year pattern (handles BE, AD, BC, BCE, CE)
func ParseYear(match string) int {
	// Buddhist Era
	if regexp.MustCompile(`(?i)BE`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*BE`).ReplaceAllString(match, "")
		year, err := strconv.Atoi(strings.TrimSpace(cleaned))
		if err != nil {
			return 0
		}
		return year - 543
	}

	// Before Christ / Before Common Era
	if regexp.MustCompile(`(?i)BCE?`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*BCE?`).ReplaceAllString(match, "")
		year, err := strconv.Atoi(strings.TrimSpace(cleaned))
		if err != nil {
			return 0
		}
		return -year
	}

	// Anno Domini / Common Era
	if regexp.MustCompile(`(?i)(AD|CE)`).MatchString(match) {
		cleaned := regexp.MustCompile(`(?i)\s*(AD|CE)`).ReplaceAllString(match, "")
		year, err := strconv.Atoi(strings.TrimSpace(cleaned))
		if err != nil {
			return 0
		}
		return year
	}

	// Regular year number
	year, err := strconv.Atoi(strings.TrimSpace(match))
	if err != nil {
		return 0
	}
	return helpers.FindMostLikelyADYear(year)
}

// ParseNumberPattern parses number-like patterns including words
func ParseNumberPattern(match string) float64 {
	lower := strings.ToLower(strings.TrimSpace(match))

	// Check number word dictionary (includes all word numbers and quantifiers)
	if val, ok := NumberWordDictionary[lower]; ok {
		return val
	}

	// Special handling for "the" - only treat as 1 in certain contexts
	// For now, skip "the" as it's not typically used as a quantity
	if lower == "the" {
		return 1
	}

	// Normalize comma to dot for decimal separator (European format support)
	normalized := strings.ReplaceAll(lower, ",", ".")

	// Try parsing as number
	val, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0
	}
	return val
}

// ParseDuration parses a duration expression like "3 days", "5 hours 30 minutes", "half an hour"
func ParseDuration(text string) kronos.Duration {
	result := make(kronos.Duration)

	// Normalize the text to handle commas and "and" connectors
	// Replace patterns like "1 year, 2 months" or "1 year and 2 months" with "1 year 2 months"
	// BUT preserve commas in decimal numbers like "2,5 hours" (European format)
	normalized := text

	// Replace comma+space followed by a number or word (not after a digit immediately before comma)
	// This preserves "2,5 hours" but replaces "2 hours, 3 minutes"
	// Pattern: match comma followed by space and then a non-digit word or number with letter
	// We need to preserve commas that are between digits (decimal separators)
	normalized = regexp.MustCompile(`([a-zA-Z])\s*,\s*([0-9])`).ReplaceAllString(normalized, "$1 $2")

	// Replace " and " between time units with a single space
	// This pattern looks for "and" surrounded by spaces/word boundaries
	normalized = regexp.MustCompile(`\s+and\s+`).ReplaceAllString(normalized, " ")

	// Pattern for matching time units
	// Supports formats like:
	// - "3 days"
	// - "two weeks"
	// - "5 hours 30 minutes"
	// - "a week"
	// - "half an hour"
	// - "1d 5h 30m"
	// - "2.5 hours" (decimal with dot)
	// - "2,5 hours" (decimal with comma, European format)
	// - "a couple of days" -> captures "couple"
	// - "a few hours" -> captures "few"
	// - "several weeks" -> captures "several"
	// Word numbers: one through nineteen, and tens (twenty through ninety)
	// Extended quantifiers: a, an, couple, few, several, half, dozen
	// Note: Using flexible pattern to support consecutive units like "2hr5min"
	// Also supports optional "of" separator (e.g., "couple of days")
	// Pattern breakdown:
	// 1. Optional leading "a/an" (not captured)
	// 2. Number/quantifier word (captured)
	// 3. Optional "a/an" again for "half an hour"
	// 4. Optional "of"
	// 5. Time unit (captured)
	pattern := regexp.MustCompile(`(?i)(?:an?\s+)?(half|dozen|several|couple|few|ninety|eighty|seventy|sixty|fifty|forty|thirty|twenty|nineteen|eighteen|seventeen|sixteen|fifteen|fourteen|thirteen|twelve|eleven|ten|nine|eight|seven|six|five|four|three|two|one|an|a|the|[0-9]+(?:[.,][0-9]+)?)\s*(?:an?\s+)?(?:of\s+)?(` + TimeUnitPattern + `)`)
	matchesIdx := pattern.FindAllStringSubmatchIndex(normalized, -1)

	for _, matchIdx := range matchesIdx {
		// matchIdx has: [fullStart, fullEnd, numStart, numEnd, unitStart, unitEnd]
		if len(matchIdx) < 6 {
			continue
		}

		fullStart := matchIdx[0]

		// Get the unit match (group 2)
		unitStart, unitEnd := matchIdx[4], matchIdx[5]
		unit := strings.ToLower(normalized[unitStart:unitEnd])

		// Skip if match is preceded by a letter (part of a larger word like "them")
		// But allow if the match starts with a digit (e.g., "2hr5min" where "5min" follows "hr")
		if fullStart > 0 {
			prevChar := normalized[fullStart-1]
			matchStartsWithDigit := (matchIdx[2] >= 0 && matchIdx[2] == fullStart &&
				len(normalized[matchIdx[2]:matchIdx[3]]) > 0 &&
				normalized[matchIdx[2]] >= '0' && normalized[matchIdx[2]] <= '9')

			if !matchStartsWithDigit &&
				((prevChar >= 'a' && prevChar <= 'z') || (prevChar >= 'A' && prevChar <= 'Z')) {
				continue
			}
		}

		// Skip if the unit is immediately followed by a letter (part of a larger word)
		if unitEnd < len(normalized) {
			nextChar := normalized[unitEnd]
			if (nextChar >= 'a' && nextChar <= 'z') || (nextChar >= 'A' && nextChar <= 'Z') {
				continue
			}
		}

		// Get the number part (group 1, if any)
		var numStr string
		if matchIdx[2] >= 0 {
			numStr = strings.TrimSpace(normalized[matchIdx[2]:matchIdx[3]])
		}

		// Single-letter time units (s, m, h, d, w, y) should have an explicit number
		// This prevents matching "am" as "a" + "m" or "them" as "the" + "m"
		if len(unit) == 1 {
			// Only allow if we have an actual digit
			if numStr == "" || (numStr != "" && !strings.ContainsAny(numStr, "0123456789")) {
				continue
			}
		}

		var num float64
		if numStr == "" {
			num = 1
		} else {
			num = ParseNumberPattern(numStr)
		}

		// Map to Timeunit and add to duration
		if timeunit, ok := TimeUnitDictionary[unit]; ok {
			// Accumulate values for same timeunit
			if existing, exists := result[timeunit]; exists {
				result[timeunit] = existing + num
			} else {
				result[timeunit] = num
			}
		}
	}

	return result
}

// IsEmptyDuration returns true if the duration has no non-zero values
func IsEmptyDuration(d kronos.Duration) bool {
	if len(d) == 0 {
		return true
	}
	for _, v := range d {
		if v != 0 {
			return false
		}
	}
	return true
}
