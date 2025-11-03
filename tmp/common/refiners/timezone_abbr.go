package refiners

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
)

var (
	timezoneNamePattern = regexp.MustCompile(`(?i)^\s*,?\s*\(?([A-Z]{2,4})\)?(?=\W|$)`)
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
	timezoneOverrides := make(map[string]int)
	for k, v := range r.timezoneOverrides {
		timezoneOverrides[k] = v
	}
	if context.Option.Timezones != nil {
		for k, v := range context.Option.Timezones {
			timezoneOverrides[k] = v
		}
	}

	for _, result := range results {
		suffix := context.Text[result.Index+len(result.Text):]
		match := timezoneNamePattern.FindStringSubmatch(suffix)
		if match == nil {
			continue
		}

		timezoneAbbr := strings.ToUpper(match[1])

		// Determine the reference date for timezone lookup
		refDate := result.Start.Date()
		if refDate.IsZero() && result.RefDate != nil {
			refDate = *result.RefDate
		}
		if refDate.IsZero() {
			refDate = context.Reference.Instant
		}

		// Look up timezone offset
		extractedTimezoneOffset, err := kronos.ToTimezoneOffset(timezoneAbbr, refDate, timezoneOverrides)
		if err != nil {
			continue
		}

		if context.Option.Debug {
			context.DebugLog("Extracting timezone: '%s' into: %d for: %s", timezoneAbbr, extractedTimezoneOffset, result.Start)
		}

		currentTimezoneOffset := result.Start.Get(kronos.ComponentTimezoneOffset)
		if currentTimezoneOffset != 0 && extractedTimezoneOffset != currentTimezoneOffset {
			// We may already have extracted the timezone offset e.g., "11 am GMT+0900 (JST)"
			// - if they are equal, we also want to take the abbreviation text into result
			// - if they are not equal, we trust the offset more
			if result.Start.IsCertain(kronos.ComponentTimezoneOffset) {
				continue
			}

			// This is often because it's relative time with inferred timezone (e.g., "in 1 hour", "tomorrow")
			// Then, we want to double-check the abbr case (e.g., "GET" not "get")
			if timezoneAbbr != match[1] {
				continue
			}
		}

		if result.Start.IsOnlyDate() {
			// If the time is not explicitly mentioned,
			// Then, we also want to double-check the abbr case (e.g., "GET" not "get")
			if timezoneAbbr != match[1] {
				continue
			}
		}

		result.Text += match[0]

		if !result.Start.IsCertain(kronos.ComponentTimezoneOffset) {
			result.Start.Assign(kronos.ComponentTimezoneOffset, extractedTimezoneOffset)
		}

		if result.End != nil && !result.End.IsCertain(kronos.ComponentTimezoneOffset) {
			result.End.Assign(kronos.ComponentTimezoneOffset, extractedTimezoneOffset)
		}
	}

	return results
}
