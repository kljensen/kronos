package kronos

import (
	"fmt"
	"sort"
	"time"

	"github.com/kljensen/kronos/internal/chrono"
)

// ============================================================================
// Pipeline and parsing execution
// ============================================================================

// pipeline manages the parsing process with configurable parsers and settings.
// It provides a flexible way to control which parsers run and in what order.
// This is an internal implementation type used by the builder pattern.
type pipeline struct {
	parsers  []Parser
	refiners []Refiner
	settings Settings
}

// newPipeline creates a new parsing pipeline with the given configuration and settings.
// This is an internal function used by the builder pattern.
func newPipeline(config *chrono.Configuration, settings Settings) *pipeline {
	if config == nil {
		config = &chrono.Configuration{
			Parsers:  []any{},
			Refiners: []any{},
		}
	}

	// Convert []any to []Parser and []Refiner
	parsers := make([]Parser, len(config.Parsers))
	for i, p := range config.Parsers {
		if parser, ok := p.(Parser); ok {
			parsers[i] = parser
		}
	}
	refiners := make([]Refiner, len(config.Refiners))
	for i, r := range config.Refiners {
		if refiner, ok := r.(Refiner); ok {
			refiners[i] = refiner
		}
	}

	return &pipeline{
		parsers:  parsers,
		refiners: refiners,
		settings: settings,
	}
}

// newPipelineWithSettings creates a pipeline using settings to determine parsers.
// This is an internal function used by the builder pattern.
func newPipelineWithSettings(config *chrono.Configuration, settings Settings) (*pipeline, error) {
	// Validate settings
	if err := validateSettings(settings); err != nil {
		return nil, fmt.Errorf("invalid settings: %w", err)
	}

	// Convert []any to []Parser and []Refiner
	var parsers []Parser
	var refiners []Refiner
	if config != nil {
		parsers = make([]Parser, len(config.Parsers))
		for i, p := range config.Parsers {
			if parser, ok := p.(Parser); ok {
				parsers[i] = parser
			}
		}
		refiners = make([]Refiner, len(config.Refiners))
		for i, r := range config.Refiners {
			if refiner, ok := r.(Refiner); ok {
				refiners[i] = refiner
			}
		}
	}

	return &pipeline{
		parsers:  parsers,
		refiners: refiners,
		settings: settings,
	}, nil
}

// Execute runs the pipeline on the given text with a reference date.
// It returns all parsed results after applying refiners.
func (p *pipeline) Execute(text string, refDate time.Time) ([]*parsingResult, error) {
	// Create parsing context with settings
	ctx, err := applySettings(text, refDate, p.settings)
	if err != nil {
		return nil, fmt.Errorf("failed to apply settings: %w", err)
	}

	// Execute parsers
	results := make([]*parsingResult, 0)
	for _, parser := range p.parsers {
		parsedResults := executeParser(ctx, parser)
		results = append(results, parsedResults...)
	}

	// Sort results by position in text
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	// Apply refiners
	for _, refiner := range p.refiners {
		results = refiner.Refine(ctx, results)
	}

	// Apply strict parsing validation
	if p.settings.StrictParsing {
		results = p.applyStrictValidation(results)
	}

	return results, nil
}

// executeParser is the shared implementation for executing a single parser.
// It handles:
// - Finding all matches in the text
// - Proper index tracking as text is consumed
// - Multiple return types from Extract: map, ParsingComponents, ParsingResult
// - Overlapping matches by advancing by 1 on extract failure
func executeParser(context *parsingContext, parser Parser) []*parsingResult {
	results := make([]*parsingResult, 0)
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

		// Build the match array
		numGroups := len(match) / 2
		matchArray := make([]string, numGroups)
		for i := range numGroups {
			start, end := match[2*i], match[2*i+1]
			if start >= 0 {
				matchArray[i] = remainingText[start:end]
			}
		}

		// Call the parser's Extract method
		result := parser.Extract(context, matchArray)
		if result == nil {
			remainingText = originalText[index+1:]
			continue
		}

		// Convert result to ParsingResult
		parsedResult := convertToParsingResult(context, result, index, matchedText, matchedTextLen)
		if parsedResult == nil {
			remainingText = originalText[index+1:]
			continue
		}

		results = append(results, parsedResult)
		remainingText = originalText[index+matchedTextLen:]
	}

	return results
}

// convertToParsingResult converts various result types to ParsingResult.
// Returns nil if the result cannot be converted.
func convertToParsingResult(context *parsingContext, result any, index int, matchedText string, matchedTextLen int) *parsingResult {
	switch v := result.(type) {
	case *parsingResult:
		if v == nil {
			return nil
		}
		headerOffset := v.Index()
		v.SetIndex(index + headerOffset)
		return v
	case *parsingResultWithBoundary:
		resultIndex := index
		if v.IncludeBoundaryIdx {
			resultIndex = index + v.BoundaryLen
		}
		parsedResult := context.CreateParsingResult(resultIndex, v.AdjustedText)
		parsedResult.start = v.Components
		return parsedResult
	case *parsingComponents:
		parsedResult := context.CreateParsingResult(index, matchedText)
		parsedResult.start = v
		return parsedResult
	case map[Component]int:
		return context.CreateParsingResult(index, matchedText, v)
	default:
		return nil
	}
}

// ParserCount returns the number of parsers in the pipeline.
// This is useful for testing and validation.
func (p *pipeline) ParserCount() int {
	return len(p.parsers)
}

func (p *pipeline) applyStrictValidation(results []*parsingResult) []*parsingResult {
	filtered := make([]*parsingResult, 0, len(results))
	for _, result := range results {
		// In strict mode, require at least year and month
		start := result.Start()
		hasYear := start.IsCertain(ComponentYear)
		hasMonth := start.IsCertain(ComponentMonth)

		// Accept if it has year and month
		if hasYear && hasMonth {
			filtered = append(filtered, result)
		}
	}
	return filtered
}
