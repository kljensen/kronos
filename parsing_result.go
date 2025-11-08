package kronos

import (
	"fmt"
	"time"
)

// parsingResult represents a parsed result containing date/time information.
// This is an internal implementation type. External code should use the Result interface.
type parsingResult struct {
	reference *referenceWithTimezone
	refDate   time.Time
	index     int
	text      string
	start     *parsingComponents
	end       *parsingComponents
}

// NewParsingResult creates a new ParsingResult.
func newParsingResult(reference *referenceWithTimezone, index int, text string, start, end *parsingComponents) *parsingResult {
	if start == nil {
		start = newParsingComponents(reference, nil)
	}

	return &parsingResult{
		reference: reference,
		refDate:   reference.Instant(),
		index:     index,
		text:      text,
		start:     start,
		end:       end,
	}
}

// Clone creates a deep copy of the ParsingResult.
func (pr *parsingResult) Clone() *parsingResult {
	var startClone *parsingComponents
	if pr.start != nil {
		startClone = pr.start.Clone()
	}

	var endClone *parsingComponents
	if pr.end != nil {
		endClone = pr.end.Clone()
	}

	return newParsingResult(pr.reference, pr.index, pr.text, startClone, endClone)
}

// Date returns a time.Time object created from the start components.
func (pr *parsingResult) Date() time.Time {
	return pr.start.Date()
}

// AddTag adds a debugging tag to both start and end components.
func (pr *parsingResult) AddTag(tag string) *parsingResult {
	pr.start.AddTag(tag)
	if pr.end != nil {
		pr.end.AddTag(tag)
	}
	return pr
}

// Tags returns combined debugging tags from start and end components.
func (pr *parsingResult) Tags() map[string]bool {
	combinedTags := make(map[string]bool)

	for tag := range pr.start.Tags() {
		combinedTags[tag] = true
	}

	if pr.end != nil {
		for tag := range pr.end.Tags() {
			combinedTags[tag] = true
		}
	}

	return combinedTags
}

// String returns a string representation for debugging.
func (pr *parsingResult) String() string {
	tagList := make([]string, 0, len(pr.Tags()))
	for tag := range pr.Tags() {
		tagList = append(tagList, tag)
	}

	return fmt.Sprintf("[ParsingResult {index: %d, text: '%s', tags: %v}]",
		pr.index, pr.text, tagList)
}

// RefDate returns the reference date used for parsing.
func (pr *parsingResult) RefDate() time.Time {
	return pr.refDate
}

// Index returns the position in the input text.
func (pr *parsingResult) Index() int {
	return pr.index
}

// SetIndex sets the position in the input text.
// This is used internally by parsers and the chrono executor.
func (pr *parsingResult) SetIndex(index int) {
	pr.index = index
}

// SetStart sets the start component of the parsing result.
// This is used by parsers and refiners to update the parsed components.
func (pr *parsingResult) SetStart(start *parsingComponents) {
	pr.start = start
}

// parsingResultWithBoundary wraps parsingComponents with boundary information.
// This is used internally to communicate the adjusted text (without boundary) to chrono.go
// when parsers using AbstractParserWithWordBoundary return components.
type parsingResultWithBoundary struct {
	Components         *parsingComponents
	AdjustedText       string
	BoundaryLen        int
	IncludeBoundaryIdx bool // If true, index points past boundary; if false, index points at boundary start
}

// Text returns the matched text from the input.
func (pr *parsingResult) Text() string {
	return pr.text
}

// Start returns the starting date/time components.
func (pr *parsingResult) Start() Components {
	return pr.start
}

// End returns the ending date/time components.
func (pr *parsingResult) End() Components {
	if pr.end == nil {
		return nil
	}
	return pr.end
}

// ============================================================================
// Adapters for public interfaces
// ============================================================================

// resultAdapter adapts an internal ParsingResult to implement the public Result interface.
// This allows us to expose a clean public API while keeping implementation details internal.
type resultAdapter struct {
	result *parsingResult
}

// newResultAdapter creates a new resultAdapter wrapping a ParsingResult.
func newResultAdapter(result *parsingResult) *resultAdapter {
	return &resultAdapter{result: result}
}

// Text returns the matched text from the input.
func (r *resultAdapter) Text() string {
	return r.result.Text()
}

// Index returns the position in the input text where this result was found.
func (r *resultAdapter) Index() int {
	return r.result.Index()
}

// Date returns a time.Time object created from the start components.
func (r *resultAdapter) Date() time.Time {
	return r.result.Date()
}

// Start returns the starting date/time components.
func (r *resultAdapter) Start() Components {
	parsed := r.result.Start()
	// The ParsedComponents interface is implemented by ParsingComponents,
	// so we can wrap it in a componentsAdapter
	if pc, ok := parsed.(*parsingComponents); ok {
		return newComponentsAdapter(pc)
	}
	return nil
}

// End returns the ending date/time components for a range, or nil for a single date/time.
func (r *resultAdapter) End() Components {
	parsed := r.result.End()
	if parsed == nil {
		return nil
	}
	// The ParsedComponents interface is implemented by ParsingComponents,
	// so we can wrap it in a componentsAdapter
	if pc, ok := parsed.(*parsingComponents); ok {
		return newComponentsAdapter(pc)
	}
	return nil
}

// Tags returns metadata tags for this result.
// This is exposed for testing purposes to verify parser behavior.
func (r *resultAdapter) Tags() map[string]bool {
	return r.result.Tags()
}
