//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/parsing"
)

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
	// Add optional approximation words at the beginning
	// Tilde is handled separately because it's a symbol, not a word
	approximationPattern := `(?:~\s*|(?:about|around|roughly|approximately|approx|circa)\s+)?`
	pattern := `(?i)` + approximationPattern + `(?:this)?\s{0,3}(morning|afternoon|evening|night|midnight|midday|noon)`
	return regexp.MustCompile(pattern)
}

func (p *ENCasualTimeParser) innerExtract(context *kronos.ParsingContext, match []string) interface{} {
	if len(match) < 2 {
		return nil
	}

	// Check if approximation words were used by examining the full match
	fullMatch := match[0]
	_, isApproximate := kronos.XStripApproximationWords(fullMatch)

	timeWord := strings.ToLower(match[1])
	var component *kronos.ParsingComponents

	switch timeWord {
	case "afternoon":
		component = kronos.XAfternoonWithHour(context.Reference(), 15)

	case "evening", "night":
		component = kronos.XEveningWithHour(context.Reference(), 20)

	case "midnight":
		component = kronos.XMidnight(context.Reference())

	case "morning":
		component = kronos.XMorningWithHour(context.Reference(), 6)

	case "noon", "midday":
		component = kronos.XNoon(context.Reference())

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
