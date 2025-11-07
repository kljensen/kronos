package kronos

import "regexp"

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
	Extract(context *parsingContext, match []string) interface{}
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
