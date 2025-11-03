package refiners

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
)

// UnlikelyFormatFilter filters out unlikely or impossible date/time results.
// It removes results that are just numbers, have invalid dates, or unrealistic times.
type UnlikelyFormatFilter struct {
	strictMode bool
}

// NewUnlikelyFormatFilter creates a new UnlikelyFormatFilter
func NewUnlikelyFormatFilter(strictMode bool) *UnlikelyFormatFilter {
	return &UnlikelyFormatFilter{
		strictMode: strictMode,
	}
}

// Refine filters out unlikely results
func (f *UnlikelyFormatFilter) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	filtered := make([]*kronos.ParsingResult, 0, len(results))

	for _, result := range results {
		if f.isValid(context, result) {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

func (f *UnlikelyFormatFilter) isValid(context *kronos.ParsingContext, result *kronos.ParsingResult) bool {
	// Remove results that are just numbers or dots
	if regexp.MustCompile(`^\d*(\.\d*)?$`).MatchString(strings.ReplaceAll(result.Text, " ", "")) {
		if context.Option.Debug {
			context.DebugLog("Removing unlikely result '%s'", result.Text)
		}
		return false
	}

	// Check if start date is valid
	if !result.Start.IsValidDate() {
		if context.Option.Debug {
			context.DebugLog("Removing invalid result: %s (%s)", result, result.Start)
		}
		return false
	}

	// Check if end date is valid
	if result.End != nil && !result.End.IsValidDate() {
		if context.Option.Debug {
			context.DebugLog("Removing invalid result: %s (%s)", result, result.End)
		}
		return false
	}

	// Apply strict mode checks
	if f.strictMode {
		return f.isStrictModeValid(context, result)
	}

	return true
}

func (f *UnlikelyFormatFilter) isStrictModeValid(context *kronos.ParsingContext, result *kronos.ParsingResult) bool {
	// In strict mode, remove weekday-only components
	if result.Start.IsOnlyWeekdayComponent() {
		if context.Option.Debug {
			context.DebugLog("(Strict) Removing weekday only component: %s (%s)", result, result.Start)
		}
		return false
	}

	return true
}
