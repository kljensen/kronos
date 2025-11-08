package kronos

import (
	"fmt"
	"maps"
	"time"
)

// ParsingComponents represents a collection of parsed date/time components.
// Components are stored as either "known" (directly parsed) or "implied" (inferred).
//
// Deprecated: This concrete type exposes internal implementation details. New code should
// use the Components interface instead, which provides a cleaner API that hides implementation.
// This type will be moved to an internal package in a future version.
type parsingComponents struct {
	knownValues   map[Component]int
	impliedValues map[Component]int
	reference     *referenceWithTimezone
	tags          map[string]bool
	period        Period
}

// NewParsingComponents creates a new ParsingComponents with the given reference.
// It initializes implied values based on the reference date.
func newParsingComponents(reference *referenceWithTimezone, knownComponents map[Component]int) *parsingComponents {
	pc := &parsingComponents{
		knownValues:   make(map[Component]int),
		impliedValues: make(map[Component]int),
		reference:     reference,
		tags:          make(map[string]bool),
	}

	maps.Copy(pc.knownValues, knownComponents)

	// Set default implied values from reference
	const defaultImpliedHour = 12 // Noon as default

	date := reference.GetDateWithAdjustedTimezone()
	pc.Imply(ComponentDay, date.Day())
	pc.Imply(ComponentMonth, int(date.Month()))
	pc.Imply(ComponentYear, date.Year())
	pc.Imply(ComponentHour, defaultImpliedHour)
	pc.Imply(ComponentMinute, 0)
	pc.Imply(ComponentSecond, 0)
	pc.Imply(ComponentMillisecond, 0)

	return pc
}

// IsCertain returns true if the component is certain (directly mentioned).
func (pc *parsingComponents) IsCertain(component Component) bool {
	_, exists := pc.knownValues[component]
	return exists
}

// Get returns the component value for either Certain or Implied values.
// Returns nil if the component is Unknown.
func (pc *parsingComponents) Get(component Component) *int {
	if val, exists := pc.knownValues[component]; exists {
		return &val
	}
	if val, exists := pc.impliedValues[component]; exists {
		return &val
	}
	return nil
}

// Assign sets a component value as certain (known).
// If the component was previously implied, it is removed from implied values.
func (pc *parsingComponents) Assign(component Component, value int) *parsingComponents {
	pc.knownValues[component] = value
	delete(pc.impliedValues, component)
	return pc
}

// Imply sets a component value as implied.
// If the component is already known, this does nothing.
func (pc *parsingComponents) Imply(component Component, value int) *parsingComponents {
	if _, exists := pc.knownValues[component]; exists {
		return pc
	}
	pc.impliedValues[component] = value
	return pc
}

// Delete removes components from both known and implied values.
func (pc *parsingComponents) Delete(components ...Component) *parsingComponents {
	for _, component := range components {
		delete(pc.knownValues, component)
		delete(pc.impliedValues, component)
	}
	return pc
}

// Clone creates a deep copy of the ParsingComponents.
func (pc *parsingComponents) Clone() *parsingComponents {
	clone := &parsingComponents{
		knownValues:   make(map[Component]int),
		impliedValues: make(map[Component]int),
		reference:     pc.reference,
		tags:          make(map[string]bool),
		period:        pc.period,
	}

	maps.Copy(clone.knownValues, pc.knownValues)
	maps.Copy(clone.impliedValues, pc.impliedValues)
	maps.Copy(clone.tags, pc.tags)

	return clone
}

// IsOnlyDate returns true if only date components are certain (no time components).
func (pc *parsingComponents) IsOnlyDate() bool {
	return !pc.IsCertain(ComponentHour) && !pc.IsCertain(ComponentMinute) && !pc.IsCertain(ComponentSecond)
}

// IsOnlyTime returns true if only time components are certain (no date components).
func (pc *parsingComponents) IsOnlyTime() bool {
	return !pc.IsCertain(ComponentWeekday) && !pc.IsCertain(ComponentDay) &&
		!pc.IsCertain(ComponentMonth) && !pc.IsCertain(ComponentYear)
}

// IsOnlyWeekdayComponent returns true if only weekday is certain without day or month.
func (pc *parsingComponents) IsOnlyWeekdayComponent() bool {
	return pc.IsCertain(ComponentWeekday) && !pc.IsCertain(ComponentDay) && !pc.IsCertain(ComponentMonth)
}

// IsDateWithUnknownYear returns true if month is certain but year is not.
func (pc *parsingComponents) IsDateWithUnknownYear() bool {
	return pc.IsCertain(ComponentMonth) && !pc.IsCertain(ComponentYear)
}

// IsValidDate validates that the components form a valid date.
func (pc *parsingComponents) IsValidDate() bool {
	date := pc.DateWithoutTimezoneAdjustment()

	yearVal := pc.Get(ComponentYear)
	if yearVal != nil && date.Year() != *yearVal {
		return false
	}

	monthVal := pc.Get(ComponentMonth)
	if monthVal != nil && int(date.Month()) != *monthVal {
		return false
	}

	dayVal := pc.Get(ComponentDay)
	if dayVal != nil && date.Day() != *dayVal {
		return false
	}

	hourVal := pc.Get(ComponentHour)
	if hourVal != nil && date.Hour() != *hourVal {
		return false
	}

	minuteVal := pc.Get(ComponentMinute)
	if minuteVal != nil && date.Minute() != *minuteVal {
		return false
	}

	return true
}

// Date returns a time.Time object constructed from the components.
// It applies timezone adjustments as needed.
func (pc *parsingComponents) Date() time.Time {
	date := pc.DateWithoutTimezoneAdjustment()

	timezoneOffsetVal := pc.Get(ComponentTimezoneOffset)
	timezoneAdjustment := pc.reference.GetSystemTimezoneAdjustmentMinute(date, timezoneOffsetVal)

	return date.Add(time.Duration(timezoneAdjustment) * time.Minute)
}

// DateWithoutTimezoneAdjustment creates a time.Time from components without timezone adjustment.
// This is useful for DST calculations where you need the "wall clock" time.
func (pc *parsingComponents) DateWithoutTimezoneAdjustment() time.Time {
	// Use the reference date's location to avoid timezone conversion issues
	location := time.Local
	if pc.reference != nil && !pc.reference.instant.IsZero() {
		location = pc.reference.instant.Location()
	}
	return pc.dateInLocation(location)
}

// DateUTC creates a UTC time.Time from components for DST calculations.
// This returns a time in UTC with the "wall clock" values from the components.
func (pc *parsingComponents) DateUTC() time.Time {
	return pc.dateInLocation(time.UTC)
}

// AddTag adds a debugging tag to the components.
func (pc *parsingComponents) AddTag(tag string) *parsingComponents {
	pc.tags[tag] = true
	return pc
}

// Tags returns all debugging tags.
func (pc *parsingComponents) Tags() map[string]bool {
	tags := make(map[string]bool)
	maps.Copy(tags, pc.tags)
	return tags
}

// String returns a string representation for debugging.
func (pc *parsingComponents) String() string {
	tagList := make([]string, 0, len(pc.tags))
	for tag := range pc.tags {
		tagList = append(tagList, tag)
	}

	return fmt.Sprintf("[ParsingComponents {tags: %v, knownValues: %v, impliedValues: %v}]",
		tagList, pc.knownValues, pc.impliedValues)
}

// Reference returns the reference.
func (pc *parsingComponents) Reference() *referenceWithTimezone {
	return pc.reference
}

// Period returns the granularity/period of the parsed date.
func (pc *parsingComponents) Period() Period {
	return pc.period
}

// SetPeriod sets the granularity/period of the parsed date.
func (pc *parsingComponents) SetPeriod(period Period) *parsingComponents {
	pc.period = period
	return pc
}

// AssignSimilarDate assigns (force updates) the parsing components to the same day as the target.
// This sets year, month, and day as certain (known) values.
func (pc *parsingComponents) AssignSimilarDate(date time.Time) {
	pc.Assign(ComponentDay, date.Day())
	pc.Assign(ComponentMonth, int(date.Month()))
	pc.Assign(ComponentYear, date.Year())
}

// breakdownNanoseconds splits nanoseconds into milliseconds, microseconds, and nanoseconds.
func breakdownNanoseconds(totalNanos int) (millisecond, microsecond, nanosecond int) {
	millisecond = totalNanos / 1000000
	remainingNanos := totalNanos % 1000000
	microsecond = remainingNanos / 1000
	nanosecond = remainingNanos % 1000
	return
}

// hourToMeridiem converts an hour (0-23) to meridiem value (0=AM, 1=PM).
func hourToMeridiem(hour int) int {
	if hour < 12 {
		return 0 // AM
	}
	return 1 // PM
}

// setSimilarTimeComponents sets time components using the provided setter function.
// This is used by both AssignSimilarTime and ImplySimilarTime to avoid duplication.
func setSimilarTimeComponents(date time.Time, setter func(Component, int) *parsingComponents) {
	setter(ComponentHour, date.Hour())
	setter(ComponentMinute, date.Minute())
	setter(ComponentSecond, date.Second())

	millisecond, microsecond, nanosecond := breakdownNanoseconds(date.Nanosecond())

	setter(ComponentMillisecond, millisecond)
	if microsecond > 0 {
		setter(ComponentMicrosecond, microsecond)
	}
	if nanosecond > 0 {
		setter(ComponentNanosecond, nanosecond)
	}

	setter(ComponentMeridiem, hourToMeridiem(date.Hour()))
}

// AssignSimilarTime assigns (force updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as certain (known) values.
func (pc *parsingComponents) AssignSimilarTime(date time.Time) {
	setSimilarTimeComponents(date, pc.Assign)
}

// ImplySimilarDate implies (weakly updates) the parsing components to the same day as the target.
// This sets year, month, and day as implied values (only if not already certain).
func (pc *parsingComponents) ImplySimilarDate(date time.Time) {
	pc.Imply(ComponentDay, date.Day())
	pc.Imply(ComponentMonth, int(date.Month()))
	pc.Imply(ComponentYear, date.Year())
}

// ImplySimilarTime implies (weakly updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as implied values (only if not already certain).
func (pc *parsingComponents) ImplySimilarTime(date time.Time) {
	setSimilarTimeComponents(date, pc.Imply)
}

// dateInLocation creates a time.Time from components in the specified location.
func (pc *parsingComponents) dateInLocation(location *time.Location) time.Time {
	const (
		defaultYear  = 2000
		defaultMonth = 1
		defaultDay   = 1
	)

	year := getValueOrDefault(pc.Get(ComponentYear), defaultYear)
	month := getValueOrDefault(pc.Get(ComponentMonth), defaultMonth)
	day := getValueOrDefault(pc.Get(ComponentDay), defaultDay)
	hour := getValueOrDefault(pc.Get(ComponentHour), 0)
	minute := getValueOrDefault(pc.Get(ComponentMinute), 0)
	second := getValueOrDefault(pc.Get(ComponentSecond), 0)
	millisecond := getValueOrDefault(pc.Get(ComponentMillisecond), 0)
	microsecond := getValueOrDefault(pc.Get(ComponentMicrosecond), 0)
	nanosecond := getValueOrDefault(pc.Get(ComponentNanosecond), 0)

	// Calculate total nanoseconds from milliseconds, microseconds, and nanoseconds
	totalNanos := (millisecond * 1000000) + (microsecond * 1000) + nanosecond

	return time.Date(year, time.Month(month), day, hour, minute, second, totalNanos, location)
}

// componentsAdapter adapts an internal ParsingComponents to implement the public Components interface.
type componentsAdapter struct {
	components *parsingComponents
}

// newComponentsAdapter creates a new componentsAdapter wrapping ParsingComponents.
func newComponentsAdapter(components *parsingComponents) *componentsAdapter {
	return &componentsAdapter{components: components}
}

// Get returns the component value.
func (c *componentsAdapter) Get(component Component) *int {
	return c.components.Get(component)
}

// IsCertain returns true if the component was explicitly mentioned in the input.
func (c *componentsAdapter) IsCertain(component Component) bool {
	return c.components.IsCertain(component)
}

// Date returns a time.Time object constructed from the components.
func (c *componentsAdapter) Date() time.Time {
	return c.components.Date()
}

// Tags returns metadata tags for these components.
// This is exposed for testing purposes to verify parser behavior.
func (c *componentsAdapter) Tags() map[string]bool {
	return c.components.Tags()
}

// Helper function to get value or default
func getValueOrDefault(ptr *int, defaultVal int) int {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}
