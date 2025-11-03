package refiners

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
)

var (
	timezoneNamePattern = regexp.MustCompile(`(?i)^\s*,?\s*\(?([A-Z]{2,4})\)?(?:\W|$)`)
)

// ExtractTimezoneAbbrRefiner extracts timezone abbreviations from text following a parsed result.
// Examples: "UTC", "PST", "JST", "EST"
type ExtractTimezoneAbbrRefiner struct {
	timezoneOverrides map[string]int
}

// NewExtractTimezoneAbbrRefiner creates a new ExtractTimezoneAbbrRefiner
func NewExtractTimezoneAbbrRefiner(timezoneOverrides map[string]int) *ExtractTimezoneAbbrRefiner {
	if timezoneOverrides == nil {
		timezoneOverrides = make(map[string]int)
	}
	return &ExtractTimezoneAbbrRefiner{
		timezoneOverrides: timezoneOverrides,
	}
}

// Refine extracts timezone abbreviations and adds them to results
func (r *ExtractTimezoneAbbrRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	// Merge context timezones with refiner overrides
	timezoneOverrides := make(kronos.TimezoneAbbrMap)
	for k, v := range r.timezoneOverrides {
		timezoneOverrides[k] = v
	}
	if context.Option().Timezones != nil {
		for k, v := range context.Option().Timezones {
			timezoneOverrides[k] = v
		}
	}

	for _, result := range results {
		// Calculate the position after the result text
		suffixStart := result.Index() + len(result.Text())
		// Check if we're beyond the end of the text
		if suffixStart >= len(context.Text()) {
			continue
		}
		suffix := context.Text()[suffixStart:]
		match := timezoneNamePattern.FindStringSubmatch(suffix)
		if match == nil {
			continue
		}

		timezoneAbbr := strings.ToUpper(match[1])

		// Determine the reference date for timezone lookup
		refDate := result.Start().Date()
		if refDate.IsZero() {
			refDate = result.RefDate()
		}
		if refDate.IsZero() {
			refDate = context.Reference().Instant()
		}

		// Look up timezone offset
		extractedTimezoneOffsetPtr := kronos.ToTimezoneOffset(timezoneAbbr, refDate, timezoneOverrides)
		if extractedTimezoneOffsetPtr == nil {
			continue
		}
		extractedTimezoneOffset := *extractedTimezoneOffsetPtr

		if context.Option().Debug != nil {
			context.Debug(func() {
				// Log: Extracting timezone
			})
		}

		resultStart := result.Start().(*kronos.ParsingComponents)
		currentTimezoneOffsetPtr := resultStart.Get(kronos.ComponentTimezoneOffset)
		currentTimezoneOffset := 0
		if currentTimezoneOffsetPtr != nil {
			currentTimezoneOffset = *currentTimezoneOffsetPtr
		}

		if currentTimezoneOffset != 0 && extractedTimezoneOffset != currentTimezoneOffset {
			// We may already have extracted the timezone offset e.g., "11 am GMT+0900 (JST)"
			// - if they are equal, we also want to take the abbreviation text into result
			// - if they are not equal, we trust the offset more
			if resultStart.IsCertain(kronos.ComponentTimezoneOffset) {
				continue
			}

			// This is often because it's relative time with inferred timezone (e.g., "in 1 hour", "tomorrow")
			// Then, we want to double-check the abbr case (e.g., "GET" not "get")
			if timezoneAbbr != match[1] {
				continue
			}
		}

		if resultStart.IsOnlyDate() {
			// If the time is not explicitly mentioned,
			// Then, we also want to double-check the abbr case (e.g., "GET" not "get")
			if timezoneAbbr != match[1] {
				continue
			}
		}

		// Update result with the timezone text (need to create new result since Text is immutable)
		// The regex pattern (?:\W|$) consumes trailing non-word chars, which we want to exclude
		// from the result text (e.g., "pst " should become "pst", but "pst)" should keep the paren)
		matchedText := match[0]
		// Trim trailing whitespace only (not punctuation like parentheses)
		for len(matchedText) > 0 && matchedText[len(matchedText)-1] == ' ' {
			matchedText = matchedText[:len(matchedText)-1]
		}
		newText := result.Text() + matchedText

		if !resultStart.IsCertain(kronos.ComponentTimezoneOffset) {
			resultStart.Assign(kronos.ComponentTimezoneOffset, extractedTimezoneOffset)
		}

		var resultEnd *kronos.ParsingComponents
		if result.End() != nil {
			resultEnd = result.End().(*kronos.ParsingComponents)
			if !resultEnd.IsCertain(kronos.ComponentTimezoneOffset) {
				resultEnd.Assign(kronos.ComponentTimezoneOffset, extractedTimezoneOffset)
			}
		}

		// Create new result with updated text
		newResult := context.CreateParsingResult(result.Index(), newText, resultStart, resultEnd)
		// Replace the result in the slice
		for j, r := range results {
			if r == result {
				results[j] = newResult
				break
			}
		}
	}

	return results
}
