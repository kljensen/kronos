package helpers

import (
	"regexp"
	"strings"
)

// ApproximationWords is a list of words that indicate approximate time expressions
var ApproximationWords = []string{
	"about",
	"around",
	"roughly",
	"approximately",
	"approx",
	"circa",
}

// SafeSlice returns the substring of text between start (inclusive) and end
// (exclusive) while clamping the requested range to valid bounds. The returned
// boolean is false when start is beyond the end of the string, indicating that
// the requested slice could not be produced safely.
func SafeSlice(text string, start, end int) (string, bool) {
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

// StripApproximationWords removes approximation modifiers from the input text
// and returns the cleaned text along with a flag indicating if approximation was present.
// The approximation words are matched case-insensitively and removed with their trailing whitespace.
func StripApproximationWords(input string) (cleaned string, isApproximate bool) {
	cleaned = input
	isApproximate = false

	// Check for tilde (~) symbol as approximation marker
	tildePattern := regexp.MustCompile(`(?i)~\s*`)
	if tildePattern.MatchString(cleaned) {
		isApproximate = true
		cleaned = tildePattern.ReplaceAllString(cleaned, "")
	}

	// Check for approximation words
	for _, word := range ApproximationWords {
		// Use word boundaries to avoid matching words like "about" in "roundabout"
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\s+`)
		if pattern.MatchString(cleaned) {
			isApproximate = true
			cleaned = pattern.ReplaceAllString(cleaned, "")
		}
	}

	return strings.TrimSpace(cleaned), isApproximate
}
