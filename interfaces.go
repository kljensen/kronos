package kronos

import "regexp"

// Parser is an abstraction for Chrono parsers.
// Each parser should recognize and handle a certain date format.
// Chrono uses multiple parsers (and refiners) together for parsing the input.
type Parser interface {
	// Pattern returns the regular expression pattern for this parser.
	// The pattern is used to find potential matches in the input text.
	Pattern(context *ParsingContext) *regexp.Regexp

	// Extract is called with the pattern's match.
	// It should return parsed components, a result, a component map, or nil if extraction fails.
	Extract(context *ParsingContext, match []string) interface{}
}

// Refiner is an abstraction for Chrono refiners.
// Each refiner takes the list of results (from parsers or other refiners)
// and returns another list of results.
// Chrono applies each refiner in order and returns the output from the last refiner.
type Refiner interface {
	// Refine processes a list of parsing results and returns a refined list.
	Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult
}

// Configuration holds the parsers and refiners for Chrono.
// It is simply an ordered list of parsers and refiners.
type Configuration struct {
	Parsers  []Parser
	Refiners []Refiner
}
