package kronos

import (
	"fmt"
	"sort"
	"time"
)

// Pipeline manages the parsing process with configurable parsers and settings.
// It provides a flexible way to control which parsers run and in what order.
type Pipeline struct {
	parsers  []Parser
	refiners []Refiner
	settings Settings
}

// NewPipeline creates a new parsing pipeline with the given configuration and settings.
func NewPipeline(config *Configuration, settings Settings) *Pipeline {
	if config == nil {
		config = &Configuration{
			Parsers:  []Parser{},
			Refiners: []Refiner{},
		}
	}

	return &Pipeline{
		parsers:  append([]Parser{}, config.Parsers...),
		refiners: append([]Refiner{}, config.Refiners...),
		settings: settings,
	}
}

// NewPipelineWithSettings creates a pipeline using settings to determine parsers.
// If EnabledParsers is empty, all parsers from the configuration are used.
// If EnabledParsers is specified, only those parsers are used.
// ParserOrder determines the execution order.
func NewPipelineWithSettings(config *Configuration, settings Settings) (*Pipeline, error) {
	// Validate settings
	if err := ValidateSettings(settings); err != nil {
		return nil, fmt.Errorf("invalid settings: %w", err)
	}

	// Start with empty pipeline
	pipeline := &Pipeline{
		parsers:  []Parser{},
		refiners: []Refiner{},
		settings: settings,
	}

	// Determine which parsers to use
	if len(settings.EnabledParsers) > 0 {
		// Use specified parsers from registry
		parsers := GlobalRegistry.GetParsers(settings.EnabledParsers)
		pipeline.parsers = parsers
	} else if config != nil {
		// Use all parsers from configuration
		pipeline.parsers = append([]Parser{}, config.Parsers...)
	} else {
		// Use all registered parsers
		pipeline.parsers = GlobalRegistry.GetAllParsers()
	}

	// Apply parser order if specified
	if len(settings.ParserOrder) > 0 {
		pipeline.parsers = reorderParsers(pipeline.parsers, settings.ParserOrder)
	}

	// Apply max parsers limit
	if settings.MaxParsers > 0 && len(pipeline.parsers) > settings.MaxParsers {
		pipeline.parsers = pipeline.parsers[:settings.MaxParsers]
	}

	// Copy refiners from configuration
	if config != nil {
		pipeline.refiners = append([]Refiner{}, config.Refiners...)
	}

	return pipeline, nil
}

// reorderParsers reorders parsers according to the specified order.
// Parsers not in the order list are appended at the end.
func reorderParsers(parsers []Parser, order []string) []Parser {
	// This is a simple implementation.
	// A more sophisticated version would use parser names from a registry.
	// For now, we keep the existing order since parsers don't have names yet.
	return parsers
}

// Execute runs the pipeline on the given text with a reference date.
// It returns all parsed results after applying refiners.
func (p *Pipeline) Execute(text string, refDate time.Time) ([]*ParsingResult, error) {
	// Create parsing context with settings
	ctx, err := ApplySettings(text, refDate, p.settings)
	if err != nil {
		return nil, fmt.Errorf("failed to apply settings: %w", err)
	}

	// Set up timeout if specified
	var timeoutChan <-chan time.Time
	if p.settings.Timeout > 0 {
		timer := time.NewTimer(p.settings.Timeout)
		defer timer.Stop()
		timeoutChan = timer.C
	}

	// Execute parsers
	results := make([]*ParsingResult, 0)
	for _, parser := range p.parsers {
		// Check timeout
		if timeoutChan != nil {
			select {
			case <-timeoutChan:
				return nil, fmt.Errorf("parsing timeout after %v", p.settings.Timeout)
			default:
			}
		}

		// Execute parser
		parsedResults := p.executeParser(ctx, parser)
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

	// Apply required parts validation
	if len(p.settings.RequireParts) > 0 {
		results = p.applyRequiredParts(results)
	}

	// Apply timezone conversion if needed
	if p.settings.ToTimezone != "" {
		results, err = p.applyTimezoneConversion(results)
		if err != nil {
			return nil, fmt.Errorf("failed to apply timezone conversion: %w", err)
		}
	}

	return results, nil
}

// executeParser executes a single parser on the context.
// This is similar to Chrono.executeParser but uses the pipeline's context.
func (p *Pipeline) executeParser(context *ParsingContext, parser Parser) []*ParsingResult {
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

		// Build the match array
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
			remainingText = originalText[index+1:]
			continue
		}

		// Convert result to ParsingResult
		var parsedResult *ParsingResult
		var matchEndPos int
		switch v := result.(type) {
		case *ParsingResult:
			parsedResult = v
			if parsedResult == nil {
				remainingText = originalText[index+1:]
				continue
			}
			headerOffset := parsedResult.Index()
			parsedResult.SetIndex(index + headerOffset)
			matchEndPos = index + matchedTextLen
		case *ParsingResultWithBoundary:
			var resultIndex int
			if v.IncludeBoundaryIdx {
				resultIndex = index + v.BoundaryLen
			} else {
				resultIndex = index
			}
			parsedResult = context.CreateParsingResult(resultIndex, v.AdjustedText)
			parsedResult.start = v.Components
			matchEndPos = index + matchedTextLen
		case *ParsingComponents:
			parsedResult = context.CreateParsingResult(index, matchedText)
			parsedResult.start = v
			matchEndPos = index + matchedTextLen
		case map[Component]int:
			parsedResult = context.CreateParsingResult(index, matchedText, v)
			matchEndPos = index + matchedTextLen
		default:
			remainingText = originalText[index+1:]
			continue
		}

		results = append(results, parsedResult)
		remainingText = originalText[matchEndPos:]
	}

	return results
}

// applyStrictValidation filters out results that don't meet strict parsing criteria.
// In strict mode, we reject results that are too ambiguous.
func (p *Pipeline) applyStrictValidation(results []*ParsingResult) []*ParsingResult {
	filtered := make([]*ParsingResult, 0, len(results))
	for _, result := range results {
		// In strict mode, require at least year and month
		if result.Start() != nil {
			start := result.Start()
			hasYear := start.IsCertain(ComponentYear)
			hasMonth := start.IsCertain(ComponentMonth)

			// Accept if it has year and month
			if hasYear && hasMonth {
				filtered = append(filtered, result)
			}
		}
	}
	return filtered
}

// applyRequiredParts filters out results that don't have required components.
func (p *Pipeline) applyRequiredParts(results []*ParsingResult) []*ParsingResult {
	if len(p.settings.RequireParts) == 0 {
		return results
	}

	filtered := make([]*ParsingResult, 0, len(results))
	for _, result := range results {
		if result.Start() == nil {
			continue
		}

		start := result.Start()
		hasAllParts := true

		for _, part := range p.settings.RequireParts {
			component := Component(part)
			if !start.IsCertain(component) {
				hasAllParts = false
				break
			}
		}

		if hasAllParts {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// applyTimezoneConversion converts results to the target timezone.
//
// TODO: EXPERIMENTAL/UNIMPLEMENTED
// This function currently validates the timezone but does not modify the results.
// A complete implementation needs to:
// 1. Convert the time.Time value to the target timezone
// 2. Update the ParsingResult's underlying components to reflect the conversion
// 3. Handle cases where timezone conversion affects date boundaries (DST, etc.)
func (p *Pipeline) applyTimezoneConversion(results []*ParsingResult) ([]*ParsingResult, error) {
	targetLoc, err := time.LoadLocation(p.settings.ToTimezone)
	if err != nil {
		return nil, err
	}

	// Convert each result's date to the target timezone
	for _, result := range results {
		if result.Start() != nil {
			// Get the date and convert to target timezone
			date := result.Start().Date()
			convertedDate := date.In(targetLoc)

			// TODO: Update the result's start components to reflect timezone conversion
			// This requires modifying the underlying ParsingComponents struct
			_ = convertedDate
		}
	}

	return results, nil
}

// ParseWithSettings is a convenience function that creates a pipeline
// and executes it with the given settings.
func ParseWithSettings(text string, refDate time.Time, settings Settings, config *Configuration) ([]*ParsingResult, error) {
	pipeline, err := NewPipelineWithSettings(config, settings)
	if err != nil {
		return nil, err
	}

	return pipeline.Execute(text, refDate)
}

// ParserCount returns the number of parsers in the pipeline.
// This is useful for testing and validation.
func (p *Pipeline) ParserCount() int {
	return len(p.parsers)
}
