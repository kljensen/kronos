package en

import (
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENTimeExpressionParser parses English time expressions with keywords:
// at, from, after, before.
// Examples: "at 3pm", "3:30pm", "15:30", "1 at night", "6 in the morning"
type ENTimeExpressionParser struct {
	*common.AbstractTimeExpressionParser
}

// NewENTimeExpressionParser creates a new English time expression parser
func NewENTimeExpressionParser(strictMode bool) *ENTimeExpressionParser {
	parser := &ENTimeExpressionParser{
		AbstractTimeExpressionParser: common.NewAbstractTimeExpressionParser(
			func() string {
				return `(?:(?:at|from)\s*)?` // Optional "at" or "from" prefix
			},
			func() string {
				// En dash (–) is Unicode U+2013, represented as literal character in regexp
				return `\s*(?:\-|–|\~|to|until|through|till|\?)\s*` // Range separator
			},
			strictMode,
		),
	}

	// Set custom primary suffix to handle "o'clock", "at night", "in the morning/afternoon", "tonight"
	parser.SetPrimarySuffix(func() string {
		// Go regexp doesn't support lookaheads (?! and ?=)
		// We use word boundary \b which prevents matching across word boundaries
		return `(?:\s*(?:o\W*clock|at\s*night|tonight|in\s*the\s*(?:morning|afternoon)))?(?:\s|$|\b)`
	})

	// Set custom following suffix to also capture "at night", etc. in ranges
	parser.SetFollowingSuffix(func() string {
		return `(?:\s*(?:o\W*clock|at\s*night|tonight|in\s*the\s*(?:morning|afternoon)))?(?:\s|$|\b)`
	})

	// Helper function to process time clues like "at night", "in the afternoon", "tonight"
	processTimeClues := func(fullMatch string, components *kronos.ParsingComponents) {
		// Handle "at night" or "tonight"
		if strings.Contains(fullMatch, "night") || strings.Contains(fullMatch, "tonight") {
			hourVal := components.Get(kronos.ComponentHour)
			if hourVal != nil {
				hour := *hourVal
				if hour >= 6 && hour < 12 {
					components.Assign(kronos.ComponentHour, hour+12)
					components.Assign(kronos.ComponentMeridiem, int(kronos.MeridiemPM))
				} else if hour < 6 {
					components.Assign(kronos.ComponentMeridiem, int(kronos.MeridiemAM))
				}
			}
		}

		// Handle "in the afternoon"
		if strings.Contains(fullMatch, "afternoon") {
			components.Assign(kronos.ComponentMeridiem, int(kronos.MeridiemPM))
			hourVal := components.Get(kronos.ComponentHour)
			if hourVal != nil {
				hour := *hourVal
				if hour >= 0 && hour <= 6 {
					components.Assign(kronos.ComponentHour, hour+12)
				}
			}
		}

		// Handle "in the morning"
		if strings.Contains(fullMatch, "morning") {
			components.Assign(kronos.ComponentMeridiem, int(kronos.MeridiemAM))
			// Hour stays as-is for morning times
		}
	}

	// Set custom extraction hook to handle "at night", "in the afternoon", etc.
	parser.SetExtractPrimaryTimeComponentsHook(func(
		context *kronos.ParsingContext,
		match []string,
		components *kronos.ParsingComponents,
	) bool {
		processTimeClues(match[0], components)
		// Add parser tag
		components.AddTag("parser/ENTimeExpressionParser")
		return true
	})

	// Set hook for following time components to also handle time clues
	parser.SetExtractFollowingTimeComponentsHook(func(
		context *kronos.ParsingContext,
		match []string,
		result *kronos.ParsingResult,
		components *kronos.ParsingComponents,
	) bool {
		processTimeClues(match[0], components)

		// If the following match has a time clue like "at night", apply it to the start component too
		// (e.g., "10 - 11 at night" means both 10pm and 11pm)
		if strings.Contains(match[0], "night") || strings.Contains(match[0], "afternoon") || strings.Contains(match[0], "morning") {
			startComponents := result.Start().(*kronos.ParsingComponents)
			processTimeClues(match[0], startComponents)
		}

		return true
	})

	return parser
}
