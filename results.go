package kronos

import (
	"fmt"
	"time"
)

// ParsingResult represents a parsed result containing date/time information.
//
// Deprecated: This concrete type exposes internal implementation details. New code should
// use the Result interface instead, which provides a cleaner API that hides implementation.
// This type will be moved to an internal package in a future version.
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

// ParsingResultWithBoundary wraps ParsingComponents with boundary information.
// This is used internally to communicate the adjusted text (without boundary) to chrono.go
// when parsers using AbstractParserWithWordBoundary return ParsingComponents.
//
// Deprecated: This is an internal implementation detail that should not be used by external code.
// It remains exported only for use by internal parser implementations. This type will be moved
// to an internal package in a future version.
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

// ============================================================================
// Merging date and time results
// ============================================================================

// mergeDateTimeResult merges a date-only result with a time-only result.
func mergeDateTimeResult(dateResult, timeResult *parsingResult) *parsingResult {
	result := dateResult.Clone()
	beginDate := dateResult.start
	beginTime := timeResult.start

	result.start = mergeDateTimeComponent(beginDate, beginTime)

	if dateResult.end != nil || timeResult.end != nil {
		var endDate, endTime *parsingComponents
		if dateResult.end == nil {
			endDate = dateResult.start
		} else {
			endDate = dateResult.end
		}
		if timeResult.end == nil {
			endTime = timeResult.start
		} else {
			endTime = timeResult.end
		}

		endDateTime := mergeDateTimeComponent(endDate, endTime)

		// If date has no end and the merged end time is before start time,
		// the end should be on the next day
		if dateResult.End() == nil && endDateTime.Date().Before(result.Start().Date()) {
			nextDay := endDateTime.Date().Add(24 * time.Hour)
			if endDateTime.IsCertain(ComponentDay) {
				endDateTime.AssignSimilarDate(nextDay)
			} else {
				endDateTime.ImplySimilarDate(nextDay)
			}
		}

		result.end = endDateTime
	}

	return result
}

// assignOrImplyComponent assigns or implies a component value based on source certainty.
func assignOrImplyComponent(result, source *parsingComponents, component Component) {
	val := source.Get(component)
	if val == nil {
		return
	}
	if source.IsCertain(component) {
		result.Assign(component, *val)
	} else {
		result.Imply(component, *val)
	}
}

// mergeComponents merges multiple components from source to result.
// Uses assignOrImplyComponent to preserve certainty information.
func mergeComponents(result, source *parsingComponents, components ...Component) {
	for _, component := range components {
		assignOrImplyComponent(result, source, component)
	}
}

// mergeDateTimeComponent merges date and time components.
func mergeDateTimeComponent(dateComp, timeComp *parsingComponents) *parsingComponents {
	result := dateComp.Clone()

	// Merge all time components, preserving certainty
	mergeComponents(result, timeComp,
		ComponentHour,
		ComponentMinute,
		ComponentSecond,
		ComponentMillisecond,
		ComponentMicrosecond,
		ComponentNanosecond,
	)

	// Merge timezone
	if timeComp.IsCertain(ComponentTimezoneOffset) {
		if tzVal := timeComp.Get(ComponentTimezoneOffset); tzVal != nil {
			result.Assign(ComponentTimezoneOffset, *tzVal)
		}
	}

	// Merge meridiem
	timeMeridiem := timeComp.Get(ComponentMeridiem)
	resultMeridiem := result.Get(ComponentMeridiem)
	if timeComp.IsCertain(ComponentMeridiem) {
		if timeMeridiem != nil {
			result.Assign(ComponentMeridiem, *timeMeridiem)
		}
	} else if timeMeridiem != nil && *timeMeridiem != 0 && resultMeridiem != nil && *resultMeridiem == 0 {
		result.Imply(ComponentMeridiem, *timeMeridiem)
	}

	// Apply PM meridiem adjustment
	// Note: hour 0 (midnight) should not be converted even with PM meridiem
	// Only hours 1-11 should be adjusted for PM (becoming 13-23)
	resultMeridiem = result.Get(ComponentMeridiem)
	resultHour := result.Get(ComponentHour)
	if resultMeridiem != nil && *resultMeridiem == 1 && resultHour != nil && *resultHour > 0 && *resultHour < 12 { // PM
		if timeComp.IsCertain(ComponentHour) {
			result.Assign(ComponentHour, *resultHour+12)
		} else {
			result.Imply(ComponentHour, *resultHour+12)
		}
	}

	// Merge tags
	for tag := range dateComp.Tags() {
		result.AddTag(tag)
	}
	for tag := range timeComp.Tags() {
		result.AddTag(tag)
	}

	// Merge period - use the finest granularity
	// If time components are present, period should be PeriodTime
	datePeriod := dateComp.Period()
	timePeriod := timeComp.Period()

	// If either is time-level, the result is time-level
	if timePeriod == PeriodTime || datePeriod == PeriodTime {
		result.SetPeriod(PeriodTime)
	} else {
		// Otherwise, use the finer of the two periods
		if datePeriod > timePeriod {
			result.SetPeriod(datePeriod)
		} else {
			result.SetPeriod(timePeriod)
		}
	}

	return result
}
