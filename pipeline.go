package kronos

import (
	"fmt"
	"sort"
	"time"
)

// ============================================================================
// Internal parsing types
// ============================================================================

// parsingOption contains configuration options for parsing.
type parsingOption struct {
	// ForwardDate indicates whether to parse only forward dates
	// (results should be after the reference date).
	ForwardDate bool

	// Preference specifies how ambiguous dates should be resolved.
	Preference DatePreference

	// DateOrder specifies the order of date components in ambiguous formats.
	DateOrder DateOrder

	// Timezones provides additional timezone keywords for parsers to recognize.
	Timezones TimezoneAbbrMap

	// Debug is an internal debug event handler.
	Debug DebugHandler
}

// parsingReference contains reference information for parsing dates/times.
type parsingReference struct {
	// Instant is the reference date/time when the input is written or mentioned.
	Instant *time.Time

	// Timezone is the reference timezone where the input is written or mentioned.
	Timezone any
}

// ============================================================================
// Parsing context
// ============================================================================

// ParsingContext holds the context for parsing operations.
// It contains the text to parse, options, and reference information.
//
// Deprecated: This concrete type exposes internal implementation details. External code should
// not depend on the internal structure of this type. It will be moved to an internal package
// in a future version.
type parsingContext struct {
	text      string
	option    parsingOption
	reference *referenceWithTimezone
	refDate   time.Time
	settings  *Settings
}

// newParsingContext creates a new ParsingContext.
// If refDate is nil, the current time is used.
// If option is nil, default options are used.
// The input text is sanitized to normalize Unicode characters before parsing.
func newParsingContext(text string, refDate any, option *parsingOption) *parsingContext {
	var opt parsingOption
	if option != nil {
		opt = *option
	}

	var timezones TimezoneAbbrMap
	if opt.Timezones != nil {
		timezones = opt.Timezones
	}

	reference := fromInput(refDate, timezones)

	// Sanitize input text to handle Unicode normalization issues
	text = sanitizeInput(text)

	return &parsingContext{
		text:      text,
		option:    opt,
		reference: reference,
		refDate:   reference.Instant(),
		settings:  nil,
	}
}

// CreateParsingComponents creates ParsingComponents from a component map or existing components.
// If components is already a ParsingComponents, it returns it as-is.
// Otherwise, it creates new ParsingComponents with the provided values.
func (ctx *parsingContext) CreateParsingComponents(components any) *parsingComponents {
	if components == nil {
		return newParsingComponents(ctx.reference, nil)
	}

	switch v := components.(type) {
	case *parsingComponents:
		return v
	case map[Component]int:
		return newParsingComponents(ctx.reference, v)
	default:
		return newParsingComponents(ctx.reference, nil)
	}
}

// CreateParsingResult creates a ParsingResult with various signatures.
// It can accept:
// - (index, text) - creates result with text
// - (index, endIndex) - creates result with substring from text
// - (index, text, startComponents) - creates result with start components
// - (index, text, startComponents, endComponents) - creates result with both start and end
func (ctx *parsingContext) CreateParsingResult(index int, textOrEndIndex any, args ...any) *parsingResult {
	var text string
	var start *parsingComponents
	var end *parsingComponents

	// Determine if second arg is text or endIndex
	switch v := textOrEndIndex.(type) {
	case string:
		text = v
	case int:
		// It's an endIndex, extract substring
		endIndex := v
		if index < 0 {
			index = 0
		}
		if endIndex > len(ctx.text) {
			endIndex = len(ctx.text)
		}
		if index > len(ctx.text) {
			index = len(ctx.text)
		}
		text = ctx.text[index:endIndex]
	default:
		text = ""
	}

	// Process remaining arguments
	for _, arg := range args {
		if arg == nil {
			continue
		}

		switch v := arg.(type) {
		case *parsingComponents:
			if start == nil {
				start = v
			} else {
				end = v
			}
		case map[Component]int:
			if start == nil {
				start = ctx.CreateParsingComponents(v)
			} else {
				end = ctx.CreateParsingComponents(v)
			}
		}
	}

	return newParsingResult(ctx.reference, index, text, start, end)
}

// Debug executes the provided function if debugging is enabled.
func (ctx *parsingContext) Debug(fn func()) {
	if ctx.option.Debug != nil {
		fn()
	}
}

// Text returns the input text.
func (ctx *parsingContext) Text() string {
	return ctx.text
}

// Option returns the parsing options.
func (ctx *parsingContext) Option() parsingOption {
	return ctx.option
}

// Reference returns the reference with timezone.
func (ctx *parsingContext) Reference() *referenceWithTimezone {
	return ctx.reference
}

// RefDate returns the reference date.
func (ctx *parsingContext) RefDate() time.Time {
	return ctx.refDate
}

// Settings returns the parsing settings, if any.
func (ctx *parsingContext) Settings() *Settings {
	return ctx.settings
}

// ============================================================================
// Pipeline and parsing execution
// ============================================================================

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
