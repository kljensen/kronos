//nolint:staticcheck // SA1019: Must use deprecated types during transition
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
	// Exception: Allow 4-digit years (1000-2999) which are valid year-only expressions
	textWithoutSpaces := strings.ReplaceAll(result.Text(), " ", "")
	if regexp.MustCompile(`^\d*(\.\d*)?$`).MatchString(textWithoutSpaces) {
		// Check if it's a 4-digit year
		if regexp.MustCompile(`^[12]\d{3}$`).MatchString(textWithoutSpaces) {
			// Check if the result has a year-level period
			if resultStart, okStart := kronos.XAsParsingComponents(result.Start()); okStart {
				if resultStart.Period() == kronos.PeriodYear {
					// This is a valid year-only expression, don't filter it
					return true
				}
			}
		}
		if context.Option().Debug != nil {
			context.Debug(func() {
				// Log: Removing unlikely result
			})
		}
		return false
	}

	// Check if start date is valid
	resultStart, okStart := kronos.XAsParsingComponents(result.Start())
	if !okStart {
		return false
	}
	if !resultStart.IsValidDate() {
		if context.Option().Debug != nil {
			context.Debug(func() {
				// Log: Removing invalid result
			})
		}
		return false
	}

	// Check if end date is valid
	if result.End() != nil {
		if resultEnd, okEnd := kronos.XAsParsingComponents(result.End()); okEnd {
			if !resultEnd.IsValidDate() {
				if context.Option().Debug != nil {
					context.Debug(func() {
						// Log: Removing invalid result with invalid end date
					})
				}
				return false
			}
		} else {
			if context.Option().Debug != nil {
				context.Debug(func() {
					// Log: Removing invalid result with invalid end date
				})
			}
			return false
		}
	}

	// Apply strict mode checks
	if f.strictMode {
		return f.isStrictModeValid(context, result)
	}

	return true
}

func (f *UnlikelyFormatFilter) isStrictModeValid(context *kronos.ParsingContext, result *kronos.ParsingResult) bool {
	// In strict mode, remove weekday-only components
	resultStart, okStart := kronos.XAsParsingComponents(result.Start())
	if !okStart {
		return false
	}
	if resultStart.IsOnlyWeekdayComponent() {
		if context.Option().Debug != nil {
			context.Debug(func() {
				// Log: (Strict) Removing weekday only component
			})
		}
		return false
	}

	return true
}
