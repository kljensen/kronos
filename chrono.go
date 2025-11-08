package kronos

import (
	"regexp"
	"sort"
)

// ============================================================================
// Deprecated Chrono API (kept for backward compatibility)
// ============================================================================

// Parser is an abstraction for Chrono parsers.
// Each parser should recognize and handle a certain date format.
// Chrono uses multiple parsers (and refiners) together for parsing the input.
//
// Deprecated: This interface is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Parser interface {
	// Pattern returns the regular expression pattern for this parser.
	// The pattern is used to find potential matches in the input text.
	Pattern(context *parsingContext) *regexp.Regexp

	// Extract is called with the pattern's match.
	// It should return parsed components, a result, a component map, or nil if extraction fails.
	Extract(context *parsingContext, match []string) any
}

// Refiner is an abstraction for Chrono refiners.
// Each refiner takes the list of results (from parsers or other refiners)
// and returns another list of results.
// Chrono applies each refiner in order and returns the output from the last refiner.
//
// Deprecated: This interface is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Refiner interface {
	// Refine processes a list of parsing results and returns a refined list.
	Refine(context *parsingContext, results []*parsingResult) []*parsingResult
}

// Configuration holds the parsers and refiners for Chrono.
// It is simply an ordered list of parsers and refiners.
//
// Deprecated: This type is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Configuration struct {
	Parsers  []Parser
	Refiners []Refiner
}

// Chrono is the main parsing engine that coordinates multiple parsers and refiners.
// It maintains a list of parsers (each handling a specific date format) and refiners
// (each post-processing the results).
//
// Deprecated: This type is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Chrono struct {
	parsers  []Parser
	refiners []Refiner
}

// NewChrono creates a new Chrono instance with the given configuration.
// If config is nil, an empty Chrono is created.
//
// Deprecated: This function is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
func NewChrono(config *Configuration) *Chrono {
	if config == nil {
		return &Chrono{
			parsers:  []Parser{},
			refiners: []Refiner{},
		}
	}

	return &Chrono{
		parsers:  append([]Parser{}, config.Parsers...),
		refiners: append([]Refiner{}, config.Refiners...),
	}
}

// Parse parses the input text and returns all found date/time results.
// The parsing process:
// 1. Create a parsing context
// 2. Execute all parsers to find matches
// 3. Sort results by position in text
// 4. Apply all refiners sequentially
// 5. Return final results
func (c *Chrono) Parse(text string, referenceDate any, option *parsingOption) []*parsingResult {
	context := newParsingContext(text, referenceDate, option)

	results := make([]*parsingResult, 0)
	for _, parser := range c.parsers {
		parsedResults := executeParser(context, parser)
		results = append(results, parsedResults...)
	}

	// Sort results by position in text
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	// Apply refiners
	for _, refiner := range c.refiners {
		results = refiner.Refine(context, results)
	}

	return results
}
