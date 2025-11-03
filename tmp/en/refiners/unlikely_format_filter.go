package refiners

import (
	"regexp"
	"strings"

	. "github.com/markusmobius/go-chrono"
)

var (
	mayContextPattern = regexp.MustCompile(`(?i)\b(in)$`)
)

// ENUnlikelyFormatFilter filters out unlikely English date formats.
type ENUnlikelyFormatFilter struct{}

func NewENUnlikelyFormatFilter() *ENUnlikelyFormatFilter {
	return &ENUnlikelyFormatFilter{}
}

func (f *ENUnlikelyFormatFilter) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	filtered := make([]*ParsingResult, 0, len(results))

	for _, result := range results {
		if !f.isValid(context, result) {
			continue
		}
		filtered = append(filtered, result)
	}

	return filtered
}

func (f *ENUnlikelyFormatFilter) isValid(context *ParsingContext, result *ParsingResult) bool {
	text := strings.TrimSpace(result.Text)

	// If the result consists of the whole text, it's likely valid
	if text == strings.TrimSpace(context.Text) {
		return true
	}

	// "may" is a month name but also a modal verb
	// Check if the text before "may" follows allowed patterns
	if strings.ToLower(text) == "may" {
		textBefore := strings.TrimSpace(context.Text[:result.Index])
		if !mayContextPattern.MatchString(textBefore) {
			return false
		}
	}

	// "the second" could refer to the ordinal number or timeunit
	if strings.HasSuffix(strings.ToLower(text), "the second") {
		textAfter := strings.TrimSpace(context.Text[result.Index+len(result.Text):])
		if len(textAfter) > 0 {
			return false
		}
	}

	return true
}
