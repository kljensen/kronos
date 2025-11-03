package en

import (
	"regexp"
	"strings"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENCasualDateParser parses casual date expressions
// Examples: now, today, tonight, tomorrow, tmr, tmrw, yesterday, last night, overmorrow
type ENCasualDateParser struct {
	*common.AbstractParserWithWordBoundary
}

// NewENCasualDateParser creates a new ENCasualDateParser
func NewENCasualDateParser() *ENCasualDateParser {
	parser := &ENCasualDateParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
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
		component = kronos.Now(context.Reference())

	case lowerText == "today":
		component = kronos.Today(context.Reference())

	case lowerText == "yesterday":
		component = kronos.Yesterday(context.Reference())

	case lowerText == "tomorrow" || lowerText == "tmr" || lowerText == "tmrw":
		component = kronos.Tomorrow(context.Reference())

	case lowerText == "tonight":
		component = kronos.TonightWithHour(context.Reference(), 22)

	case lowerText == "overmorrow":
		component = kronos.TheDayAfter(context.Reference(), 2)

	case strings.Contains(lowerText, "last") && strings.Contains(lowerText, "night"):
		// Handle "last night"
		targetDate := context.RefDate()
		if targetDate.Hour() > 6 {
			// If it's after 6 AM, "last night" refers to the previous day
			targetDate = targetDate.Add(-24 * time.Hour)
		}

		component = context.CreateParsingComponents(nil)
		kronos.AssignSimilarDate(component, targetDate)
		component.Imply(kronos.ComponentHour, 0)

	default:
		return nil
	}

	if component != nil {
		component.AddTag("parser/ENCasualDateParser")
	}

	return component
}
