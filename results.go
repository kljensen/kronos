package kronos

import (
	"fmt"
	"time"
)

// emptyDuration is an internal copy of empty duration for use within the root package.
var emptyDuration = Duration{
	TimeunitDay:         0,
	TimeunitSecond:      0,
	TimeunitMillisecond: 0,
}

// referenceWithTimezone represents a reference date/time with an optional timezone offset.
// It is used as the reference point for parsing relative dates and times.
type referenceWithTimezone struct {
	instant        time.Time
	timezoneOffset *int
}

// newReferenceWithTimezone creates a new referenceWithTimezone with the given instant and timezone offset.
// If instant is zero, the current time is used.
// If timezoneOffset is nil, the system timezone is used.
func newReferenceWithTimezone(instant time.Time, timezoneOffset *int) *referenceWithTimezone {
	if instant.IsZero() {
		instant = time.Now()
	}
	return &referenceWithTimezone{
		instant:        instant,
		timezoneOffset: timezoneOffset,
	}
}

// FromDate creates a ReferenceWithTimezone from a Date.
func (r *referenceWithTimezone) FromDate(date time.Time) *referenceWithTimezone {
	return newReferenceWithTimezone(date, nil)
}

// FromInput creates a ReferenceWithTimezone from either a ParsingReference or a time.Time.
// It also handles timezone conversion using the provided timezoneOverrides.
func fromInput(input interface{}, timezoneOverrides TimezoneAbbrMap) *referenceWithTimezone {
	if input == nil {
		return newReferenceWithTimezone(time.Time{}, nil)
	}

	switch v := input.(type) {
	case time.Time:
		return newReferenceWithTimezone(v, nil)
	case parsingReference:
		instant := time.Now()
		if v.Instant != nil {
			instant = *v.Instant
		}

		var timezoneOffset *int
		if v.Timezone != nil {
			timezoneOffset = toTimezoneOffset(v.Timezone, instant, timezoneOverrides)
		}

		return newReferenceWithTimezone(instant, timezoneOffset)
	default:
		return newReferenceWithTimezone(time.Time{}, nil)
	}
}

// GetDateWithAdjustedTimezone returns a time.Time with the year, month, day, hour, minute, second
// equal to the reference. The output's instant is NOT the reference's instant when the reference's
// and system's timezone are different.
func (r *referenceWithTimezone) GetDateWithAdjustedTimezone() time.Time {
	date := r.instant
	if r.timezoneOffset != nil {
		adjustment := r.GetSystemTimezoneAdjustmentMinute(r.instant, nil)
		date = date.Add(time.Duration(-adjustment) * time.Minute)
	}
	return date
}

// GetSystemTimezoneAdjustmentMinute returns the number of minutes difference between
// the system's timezone and the reference timezone.
func (r *referenceWithTimezone) GetSystemTimezoneAdjustmentMinute(date time.Time, overrideTimezoneOffset *int) int {
	const secondsPerMinute = 60

	if date.IsZero() || date.Unix() < 0 {
		// Javascript date timezone calculation got effect when the time epoch < 0
		date = time.Now()
	}

	_, currentOffset := date.Zone()
	currentTimezoneOffset := currentOffset / secondsPerMinute

	targetTimezoneOffset := currentTimezoneOffset
	if overrideTimezoneOffset != nil {
		targetTimezoneOffset = *overrideTimezoneOffset
	} else if r.timezoneOffset != nil {
		targetTimezoneOffset = *r.timezoneOffset
	}

	return currentTimezoneOffset - targetTimezoneOffset
}

// GetTimezoneOffset returns the timezone offset in minutes.
// If no timezone offset is set, it returns the system timezone offset.
func (r *referenceWithTimezone) GetTimezoneOffset() int {
	const secondsPerMinute = 60

	if r.timezoneOffset != nil {
		return *r.timezoneOffset
	}
	_, offset := r.instant.Zone()
	return offset / secondsPerMinute
}

// Instant returns the reference instant.
func (r *referenceWithTimezone) Instant() time.Time {
	return r.instant
}

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

	for k, v := range knownComponents {
		pc.knownValues[k] = v
	}

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

	for k, v := range pc.knownValues {
		clone.knownValues[k] = v
	}
	for k, v := range pc.impliedValues {
		clone.impliedValues[k] = v
	}
	for k, v := range pc.tags {
		clone.tags[k] = v
	}

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

// AddTag adds a debugging tag to the components.
func (pc *parsingComponents) AddTag(tag string) *parsingComponents {
	pc.tags[tag] = true
	return pc
}

// Tags returns all debugging tags.
func (pc *parsingComponents) Tags() map[string]bool {
	tags := make(map[string]bool)
	for k, v := range pc.tags {
		tags[k] = v
	}
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

// DeterminePeriodFromDuration determines the granularity/period based on a duration.
// The period represents the finest time unit present in the duration.
// This follows the pattern from Python's dateparser.
func determinePeriodFromDuration(duration Duration) Period {
	if duration == nil {
		return PeriodDay // Default
	}

	// Check from finest to coarsest granularity
	// Time components (hour, minute, second) indicate time-level precision
	for _, timeunit := range []Timeunit{TimeunitSecond, TimeunitMinute, TimeunitHour} {
		if _, exists := duration[timeunit]; exists {
			return PeriodTime
		}
	}

	// Day indicates day-level precision
	if _, exists := duration[TimeunitDay]; exists {
		return PeriodDay
	}

	// Week indicates week-level precision
	if _, exists := duration[TimeunitWeek]; exists {
		return PeriodWeek
	}

	// Month indicates month-level precision
	if _, exists := duration[TimeunitMonth]; exists {
		return PeriodMonth
	}

	// Year, decade, or quarter indicate year-level precision
	for _, timeunit := range []Timeunit{TimeunitYear, TimeunitDecade, TimeunitQuarter} {
		if _, exists := duration[timeunit]; exists {
			return PeriodYear
		}
	}

	// Default to day if no specific duration is found
	return PeriodDay
}

// DeterminePeriodFromComponents determines the granularity/period based on which
// components are certain (explicitly mentioned). The period represents the finest
// granularity of date/time information that was directly parsed.
func determinePeriodFromComponents(pc *parsingComponents) Period {
	if pc == nil {
		return PeriodUnknown
	}

	// If any time components (hour, minute, second) are certain, it's time-level
	if pc.IsCertain(ComponentHour) || pc.IsCertain(ComponentMinute) || pc.IsCertain(ComponentSecond) {
		return PeriodTime
	}

	// If day is certain, it's day-level
	if pc.IsCertain(ComponentDay) {
		return PeriodDay
	}

	// If weekday is certain (without day/month), it's week-level
	if pc.IsCertain(ComponentWeekday) && !pc.IsCertain(ComponentDay) {
		return PeriodWeek
	}

	// If month is certain (without day), it's month-level
	if pc.IsCertain(ComponentMonth) {
		return PeriodMonth
	}

	// If only year is certain, it's year-level
	if pc.IsCertain(ComponentYear) {
		return PeriodYear
	}

	// Default to unknown if nothing is certain
	return PeriodUnknown
}

// CreateRelativeFromReference creates a ParsingComponents from a duration relative to the reference.
// It handles date-only durations (implies time) and time durations (assigns both date and time).
// This is used for parsing relative expressions like "in 3 days", "2 hours ago", etc.
// Returns nil if the duration calculation fails (e.g., overflow).
func createRelativeFromReference(reference *referenceWithTimezone, duration Duration) *parsingComponents {
	if duration == nil {
		duration = emptyDuration
	}

	date, err := addDuration(reference.GetDateWithAdjustedTimezone(), duration)
	if err != nil {
		// Duration calculation failed - return nil to indicate invalid result
		return nil
	}

	components := newParsingComponents(reference, nil)
	components.AddTag("result/relativeDate")

	// Determine and set the period based on the duration
	period := determinePeriodFromDuration(duration)
	components.SetPeriod(period)

	// Check if duration contains time components
	hasTimeComponents := false
	for _, timeunit := range []Timeunit{TimeunitHour, TimeunitMinute, TimeunitSecond, TimeunitMillisecond} {
		if _, exists := duration[timeunit]; exists {
			hasTimeComponents = true
			break
		}
	}

	if hasTimeComponents {
		// Duration includes time - assign both date and time as certain
		components.AddTag("result/relativeDateAndTime")
		assignSimilarTime(components, date)
		assignSimilarDate(components, date)
		components.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	} else {
		// Duration is date-only - imply time components
		implySimilarTime(components, date)
		components.Imply(ComponentTimezoneOffset, reference.GetTimezoneOffset())

		// Handle different date granularities
		if _, hasDayDuration := duration[TimeunitDay]; hasDayDuration {
			// Day duration - assign day, month, year and weekday
			components.Assign(ComponentDay, date.Day())
			components.Assign(ComponentMonth, int(date.Month()))
			components.Assign(ComponentYear, date.Year())
			components.Assign(ComponentWeekday, int(date.Weekday()))
		} else if _, hasWeekDuration := duration[TimeunitWeek]; hasWeekDuration {
			// Week duration - assign day, month, year and imply weekday
			components.Assign(ComponentDay, date.Day())
			components.Assign(ComponentMonth, int(date.Month()))
			components.Assign(ComponentYear, date.Year())
			components.Imply(ComponentWeekday, int(date.Weekday()))
		} else {
			// Month/year duration - imply day
			components.Imply(ComponentDay, date.Day())

			if _, hasMonthDuration := duration[TimeunitMonth]; hasMonthDuration {
				// Month duration - assign month and year
				components.Assign(ComponentMonth, int(date.Month()))
				components.Assign(ComponentYear, date.Year())
			} else {
				// Imply month
				components.Imply(ComponentMonth, int(date.Month()))

				if _, hasYearDuration := duration[TimeunitYear]; hasYearDuration {
					// Year duration - assign year
					components.Assign(ComponentYear, date.Year())
				} else if _, hasQuarterDuration := duration[TimeunitQuarter]; hasQuarterDuration {
					// Quarter duration - assign year
					components.Assign(ComponentYear, date.Year())
				} else {
					// Imply year
					components.Imply(ComponentYear, date.Year())
				}
			}
		}
	}

	return components
}

// AddDurationAsImplied adds the duration to the current components and implies the result.
// This is useful for modifying existing parsing components with a relative offset.
// Returns nil if the duration calculation fails (e.g., overflow).
func (pc *parsingComponents) AddDurationAsImplied(duration Duration) *parsingComponents {
	// Get the current date from this component
	currentDate := pc.Date()

	// Add the duration
	newDate, err := addDuration(currentDate, duration)
	if err != nil {
		// Duration calculation failed - return nil to indicate invalid result
		return nil
	}

	// Imply the new date components
	implySimilarDate(pc, newDate)
	implySimilarTime(pc, newDate)

	return pc
}

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

// Helper function to get value or default
func getValueOrDefault(ptr *int, defaultVal int) int {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}
