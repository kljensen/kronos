package kronos

import (
	"fmt"
	"sort"
	"time"
)

// Pipeline manages the parsing process with configurable parsers and settings.
// It provides a flexible way to control which parsers run and in what order.
//
// Deprecated: This type is part of the advanced API. For most use cases, use the
// builder pattern instead (kronos.New(chrono).Parse(text)). Direct use of Pipeline
// exposes internal implementation details and will be moved to the experimental
// package in a future version.
type pipeline struct {
	parsers  []Parser
	refiners []Refiner
	settings Settings
}

// newPipeline creates a new parsing pipeline with the given configuration and settings.
//
// Deprecated: This function is part of the advanced API. For most use cases, use the
// builder pattern instead (kronos.New(chrono)). This function will be moved to the
// experimental package in a future version.
func newPipeline(config *Configuration, settings Settings) *pipeline {
	if config == nil {
		config = &Configuration{
			Parsers:  []Parser{},
			Refiners: []Refiner{},
		}
	}

	return &pipeline{
		parsers:  append([]Parser{}, config.Parsers...),
		refiners: append([]Refiner{}, config.Refiners...),
		settings: settings,
	}
}

// newPipelineWithSettings creates a pipeline using settings to determine parsers.
//
// Deprecated: This function is part of the advanced API. For most use cases, use the
// builder pattern instead (kronos.New(chrono).WithOption(...)). This function will be
// moved to the experimental package in a future version.
func newPipelineWithSettings(config *Configuration, settings Settings) (*pipeline, error) {
	// Validate settings
	if err := validateSettings(settings); err != nil {
		return nil, fmt.Errorf("invalid settings: %w", err)
	}

	// Start with empty pipeline
	pipeline := &pipeline{
		parsers:  []Parser{},
		refiners: []Refiner{},
		settings: settings,
	}

	// Use all parsers from configuration
	if config != nil {
		pipeline.parsers = append([]Parser{}, config.Parsers...)
		pipeline.refiners = append([]Refiner{}, config.Refiners...)
	}

	return pipeline, nil
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

	// Apply timezone conversion if needed
	if p.settings.ToTimezone != "" {
		results, err = p.applyTimezoneConversion(results)
		if err != nil {
			return nil, fmt.Errorf("failed to apply timezone conversion: %w", err)
		}
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
func convertToParsingResult(context *parsingContext, result interface{}, index int, matchedText string, matchedTextLen int) *parsingResult {
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

// applyStrictValidation filters out results that don't meet strict parsing criteria.
// In strict mode, we reject results that are too ambiguous.
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

// applyTimezoneConversion converts results to the target timezone.
// This function converts the parsed date/time to the specified timezone and updates
// all components (year, month, day, hour, minute, second, timezone offset) to reflect
// the new timezone. This handles DST transitions and date boundary changes correctly.
func (p *pipeline) applyTimezoneConversion(results []*parsingResult) ([]*parsingResult, error) {
	targetLoc, err := time.LoadLocation(p.settings.ToTimezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone location: %w", err)
	}

	// Convert each result's date to the target timezone
	for _, result := range results {
		if result.start != nil {
			// Convert the start date to the target timezone and update components
			startDate := result.start.Date()
			convertedStart := startDate.In(targetLoc)
			updateComponentsFromDate(result.start, convertedStart)

			_, startOffsetSeconds := convertedStart.Zone()
			startOffsetMinutes := startOffsetSeconds / 60
			startReference := newReferenceWithTimezone(convertedStart, &startOffsetMinutes)
			result.start.reference = startReference
			result.reference = startReference
			result.refDate = convertedStart
		}

		// Handle end components for range results
		if result.end != nil {
			endComponents := result.end
			endDate := endComponents.Date()
			convertedEnd := endDate.In(targetLoc)

			updateComponentsFromDate(endComponents, convertedEnd)

			_, endOffsetSeconds := convertedEnd.Zone()
			endOffsetMinutes := endOffsetSeconds / 60
			endComponents.reference = newReferenceWithTimezone(convertedEnd, &endOffsetMinutes)
		}
	}

	return results, nil
}

// updateComponentsFromDate updates all date/time components in ParsingComponents
// to match the values from the given time.Time in the target timezone.
// This preserves the "certain" vs "implied" status of each component while updating values.
func updateComponentsFromDate(components *parsingComponents, date time.Time) {
	// Helper to update a component preserving its certain/implied status
	updateComponent := func(comp Component, value int) {
		if components.IsCertain(comp) {
			components.Assign(comp, value)
		} else if components.Get(comp) != nil {
			components.Imply(comp, value)
		}
	}

	// Update date components
	updateComponent(ComponentYear, date.Year())
	updateComponent(ComponentMonth, int(date.Month()))
	updateComponent(ComponentDay, date.Day())

	// Update time components
	updateComponent(ComponentHour, date.Hour())
	updateComponent(ComponentMinute, date.Minute())
	updateComponent(ComponentSecond, date.Second())

	// Update subsecond components
	totalNanos := date.Nanosecond()
	updateComponent(ComponentMillisecond, totalNanos/1000000)
	updateComponent(ComponentMicrosecond, (totalNanos%1000000)/1000)
	updateComponent(ComponentNanosecond, totalNanos%1000)

	// Update meridiem
	meridiem := 0 // AM
	if date.Hour() >= 12 {
		meridiem = 1 // PM
	}
	updateComponent(ComponentMeridiem, meridiem)

	// Update timezone offset
	_, offset := date.Zone()
	updateComponent(ComponentTimezoneOffset, offset/60)
}

// parseWithSettings is a convenience function that creates a pipeline
// and executes it with the given settings.
//
// Deprecated: This function exposes internal implementation details (ParsingResult).
// Use the builder pattern API instead:
//
//	parser := kronos.New(chrono).WithReferenceDate(refDate)
//	results, err := parser.Parse(text)
//
// This function will be removed in a future version.
func parseWithSettings(text string, refDate time.Time, settings Settings, config *Configuration) ([]*parsingResult, error) {
	pipeline, err := newPipelineWithSettings(config, settings)
	if err != nil {
		return nil, err
	}

	return pipeline.Execute(text, refDate)
}

// ParserCount returns the number of parsers in the pipeline.
// This is useful for testing and validation.
func (p *pipeline) ParserCount() int {
	return len(p.parsers)
}
