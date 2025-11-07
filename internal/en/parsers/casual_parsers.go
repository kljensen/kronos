//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strings"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
	"github.com/kljensen/kronos/internal/parsing"
)

// ENCasualDateParser parses casual date expressions
// Examples: now, today, tonight, tomorrow, tmr, tmrw, yesterday, last night, overmorrow
type ENCasualDateParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENCasualDateParser creates a new ENCasualDateParser
func NewENCasualDateParser() *ENCasualDateParser {
	parser := &ENCasualDateParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENCasualDateParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	pattern := `(?i)(now|today|tonight|tomorrow|overmorrow|tmr|tmrw|yesterday|last\s*night)`
	return regexp.MustCompile(pattern)
}

func (p *ENCasualDateParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 2 {
		return nil
	}

	lowerText := strings.ToLower(match[1])
	lowerText = strings.TrimSpace(lowerText)

	var component *kronos.ParsingComponents

	switch {
	case lowerText == "now":
		component = helpers.Now(context.Reference())

	case lowerText == "today":
		component = helpers.Today(context.Reference())

	case lowerText == "yesterday":
		component = helpers.Yesterday(context.Reference())

	case lowerText == "tomorrow" || lowerText == "tmr" || lowerText == "tmrw":
		component = helpers.Tomorrow(context.Reference())

	case lowerText == "tonight":
		component = helpers.TonightWithHour(context.Reference(), 22)

	case lowerText == "overmorrow":
		component = helpers.TheDayAfter(context.Reference(), 2)

	case strings.Contains(lowerText, "last") && strings.Contains(lowerText, "night"):
		// Handle "last night"
		targetDate := context.RefDate()
		if targetDate.Hour() > 6 {
			// If it's after 6 AM, "last night" refers to the previous day
			targetDate = targetDate.Add(-24 * time.Hour)
		}

		component = context.CreateParsingComponents(nil)
		helpers.AssignSimilarDate(component, targetDate)
		component.Imply(kronos.ComponentHour, 0)
		component.SetPeriod(kronos.PeriodDay)

	default:
		return nil
	}

	if component != nil {
		component.AddTag("parser/ENCasualDateParser")
	}

	return component
}

// ENCasualTimeParser parses casual time expressions
// Examples: this morning, this afternoon, this evening, noon, midnight
type ENCasualTimeParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENCasualTimeParser creates a new ENCasualTimeParser
func NewENCasualTimeParser() *ENCasualTimeParser {
	parser := &ENCasualTimeParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		parser.innerPattern,
		parser.innerExtract,
		nil, // use default word boundary
	)

	return parser
}

func (p *ENCasualTimeParser) innerPattern(context *kronos.ParsingContext) *regexp.Regexp {
	pattern := `(?i)` + approximationPattern + `(?:this)?\s{0,3}(morning|afternoon|evening|night|midnight|midday|noon)`
	return regexp.MustCompile(pattern)
}

func (p *ENCasualTimeParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 2 {
		return nil
	}

	// Check if approximation words were used by examining the full match
	fullMatch := match[0]
	_, isApproximate := helpers.StripApproximationWords(fullMatch)

	timeWord := strings.ToLower(match[1])
	var component *kronos.ParsingComponents

	switch timeWord {
	case "afternoon":
		component = helpers.AfternoonWithHour(context.Reference(), 15)

	case "evening", "night":
		component = helpers.EveningWithHour(context.Reference(), 20)

	case "midnight":
		component = helpers.Midnight(context.Reference())

	case "morning":
		component = helpers.MorningWithHour(context.Reference(), 6)

	case "noon", "midday":
		component = helpers.Noon(context.Reference())

	default:
		return nil
	}

	if component != nil {
		component.AddTag("parser/ENCasualTimeParser")
		// Add approximation tag if approximation words were detected
		if isApproximate {
			component.AddTag("result/approximate")
		}
	}

	return component
}
