package kronos

import (
	"regexp"
	"sort"
)

// ============================================================================
// Advanced API - For creating custom parsers and configurations
// ============================================================================

// Parser is an abstraction for Chrono parsers.
// Each parser should recognize and handle a certain date format.
// Chrono uses multiple parsers (and refiners) together for parsing the input.
//
// Advanced API: This interface is part of the advanced API for creating custom parsers.
// Most users should use the builder pattern (en.New()) instead. This interface is primarily
// used by internal packages and users implementing custom date parsing logic.
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
// Advanced API: This interface is part of the advanced API for creating custom refiners.
// Most users should use the builder pattern (en.New()) instead. This interface is primarily
// used by internal packages and users implementing custom date parsing logic.
type Refiner interface {
	// Refine processes a list of parsing results and returns a refined list.
	Refine(context *parsingContext, results []*parsingResult) []*parsingResult
}

// Configuration holds the parsers and refiners for Chrono.
// It is simply an ordered list of parsers and refiners.
//
// Advanced API: This type is part of the advanced API for configuring custom parser combinations.
// Most users should use the builder pattern (en.New()) instead. This type is primarily used
// by language packages (like en) to define parsing configurations.
type Configuration struct {
	Parsers  []Parser
	Refiners []Refiner
}

// Chrono is the main parsing engine that coordinates multiple parsers and refiners.
// It maintains a list of parsers (each handling a specific date format) and refiners
// (each post-processing the results).
//
// Advanced API: This type is part of the advanced API for custom parsing engines.
// Most users should use the builder pattern (en.New()) instead. This type is primarily used
// by language packages (like en) to create pre-configured parsing engines.
type Chrono struct {
	parsers  []Parser
	refiners []Refiner
}

// NewChrono creates a new Chrono instance with the given configuration.
// If config is nil, an empty Chrono is created.
//
// Advanced API: This function is part of the advanced API for creating custom parsing engines.
// Most users should use the builder pattern (en.New()) instead.
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
