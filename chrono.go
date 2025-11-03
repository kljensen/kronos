package kronos

import (
	"sort"
	"time"
)

// Chrono is the main parsing engine that coordinates multiple parsers and refiners.
// It maintains a list of parsers (each handling a specific date format) and refiners
// (each post-processing the results).
type Chrono struct {
	parsers  []Parser
	refiners []Refiner
}

// NewChrono creates a new Chrono instance with the given configuration.
// If config is nil, an empty Chrono is created.
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
func (c *Chrono) ParseDate(text string, referenceDate interface{}, option *ParsingOption) *time.Time {
	results := c.Parse(text, referenceDate, option)
	if len(results) > 0 {
		date := results[0].Date()
		return &date
	}
	return nil
}

// Parse parses the input text and returns all found date/time results.
// The parsing process:
// 1. Create a parsing context
// 2. Execute all parsers to find matches
// 3. Sort results by position in text
// 4. Apply all refiners sequentially
// 5. Return final results
func (c *Chrono) Parse(text string, referenceDate interface{}, option *ParsingOption) []*ParsingResult {
	context := NewParsingContext(text, referenceDate, option)

	results := make([]*ParsingResult, 0)
	for _, parser := range c.parsers {
		parsedResults := c.executeParser(context, parser)
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

// executeParser executes a single parser on the text.
// It handles:
// - Finding all matches in the text
// - Proper index tracking as text is consumed
// - Three return types from Extract: map, ParsingComponents, ParsingResult
// - Overlapping matches by advancing by 1 on extract failure
func (c *Chrono) executeParser(context *ParsingContext, parser Parser) []*ParsingResult {
	results := make([]*ParsingResult, 0)
	pattern := parser.Pattern(context)

	originalText := context.Text()
	remainingText := originalText

	// Find all matches
	for {
		match := pattern.FindStringSubmatchIndex(remainingText)
		if match == nil {
			break
		}

		// Calculate match index on the full text
		index := match[0] + len(originalText) - len(remainingText)

		// Extract the matched text
		matchedText := remainingText[match[0]:match[1]]
		matchedTextLen := len(matchedText)

		// Build the match array (equivalent to RegExpMatchArray in TypeScript)
		matchArray := make([]string, len(match)/2)
		for i := 0; i < len(match); i += 2 {
			if match[i] >= 0 {
				matchArray[i/2] = remainingText[match[i]:match[i+1]]
			} else {
				matchArray[i/2] = ""
			}
		}

		// Call the parser's Extract method
		result := parser.Extract(context, matchArray)
		if result == nil {
			// If extraction fails, move on by 1
			remainingText = originalText[index+1:]
			continue
		}

		// Convert result to ParsingResult
		// Track the actual match end position for advancing remainingText
		var parsedResult *ParsingResult
		var matchEndPos int
		switch v := result.(type) {
		case *ParsingResult:
			parsedResult = v
			// Check for nil ParsingResult (typed nil)
			if parsedResult == nil {
				remainingText = originalText[index+1:]
				continue
			}
			// Update the index to reflect actual position in original text
			// Parsers using AbstractParserWithWordBoundary may have stored a header offset
			// in the index field. We add this to the match position to get the correct index.
			headerOffset := parsedResult.Index()
			parsedResult.SetIndex(index + headerOffset)
			// Match ends at the original match end position
			matchEndPos = index + matchedTextLen
		case *ParsingResultWithBoundary:
			// Parser returned ParsingComponents with boundary information
			// The AdjustedText doesn't include the boundary character, but the index
			// should point to the boundary start (where the match begins in the original text)
			var resultIndex int
			if v.IncludeBoundaryIdx {
				// Index should point past the boundary
				resultIndex = index + v.BoundaryLen
			} else {
				// Index should point at the boundary start (include the boundary in the index)
				resultIndex = index
			}
			parsedResult = context.CreateParsingResult(resultIndex, v.AdjustedText)
			parsedResult.start = v.Components
			// Match ends at the original match end position (including boundary)
			matchEndPos = index + matchedTextLen
		case *ParsingComponents:
			parsedResult = context.CreateParsingResult(index, matchedText)
			parsedResult.start = v
			// Match ends at the original match end position
			matchEndPos = index + matchedTextLen
		case map[Component]int:
			parsedResult = context.CreateParsingResult(index, matchedText, v)
			// Match ends at the original match end position
			matchEndPos = index + matchedTextLen
		default:
			// Unknown return type - skip
			remainingText = originalText[index+1:]
			continue
		}

		// Debug logging
		context.Debug(func() {
			// In a real implementation, you'd log this somewhere
			// For now, we'll skip it
		})

		results = append(results, parsedResult)

		// Move to the position after this match using the original match end position
		remainingText = originalText[matchEndPos:]
	}

	return results
}
