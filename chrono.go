package kronos

import (
	"sort"
	"time"
)

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

// Clone creates a shallow copy of the Chrono object with the same configuration.
// The parsers and refiners slices are copied, but the parsers/refiners themselves are not.
func (c *Chrono) Clone() *Chrono {
	return &Chrono{
		parsers:  append([]Parser{}, c.parsers...),
		refiners: append([]Refiner{}, c.refiners...),
	}
}

// ParseDate is a shortcut for calling Parse and returning the first result's date.
// Returns nil if no results are found.
func (c *Chrono) ParseDate(text string, referenceDate interface{}, option *parsingOption) *time.Time {
	results := c.Parse(text, referenceDate, option)
	if len(results) > 0 {
		date := results[0].Date()
		return &date
	}
	return nil
}

// parseWithSettings parses the input text using the provided settings.
// Settings provide more comprehensive configuration than ParsingOption.
// This method creates a pipeline from the current configuration and settings.
func (c *Chrono) parseWithSettings(text string, referenceDate time.Time, settings Settings) ([]*parsingResult, error) {
	// Create a configuration from this Chrono instance
	config := &Configuration{
		Parsers:  append([]Parser{}, c.parsers...),
		Refiners: append([]Refiner{}, c.refiners...),
	}

	// Create and execute pipeline
	return parseWithSettings(text, referenceDate, settings, config)
}

// ParseDateWithSettings is a shortcut for calling parseWithSettings and returning the first result's date.
// Returns nil if no results are found or an error occurs.
func (c *Chrono) ParseDateWithSettings(text string, referenceDate time.Time, settings Settings) (*time.Time, error) {
	results, err := c.parseWithSettings(text, referenceDate, settings)
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		date := results[0].Date()
		return &date, nil
	}
	return nil, nil
}

// Parse parses the input text and returns all found date/time results.
// The parsing process:
// 1. Create a parsing context
// 2. Execute all parsers to find matches
// 3. Sort results by position in text
// 4. Apply all refiners sequentially
// 5. Return final results
func (c *Chrono) Parse(text string, referenceDate interface{}, option *parsingOption) []*parsingResult {
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
