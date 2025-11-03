package kronos

import (
	"fmt"
	"time"
)

// ReferenceWithTimezone represents a reference date/time with an optional timezone offset.
// It is used as the reference point for parsing relative dates and times.
type ReferenceWithTimezone struct {
	instant        time.Time
	timezoneOffset *int
}

// NewReferenceWithTimezone creates a new ReferenceWithTimezone with the given instant and timezone offset.
// If instant is zero, the current time is used.
// If timezoneOffset is nil, the system timezone is used.
func NewReferenceWithTimezone(instant time.Time, timezoneOffset *int) *ReferenceWithTimezone {
	if instant.IsZero() {
		instant = time.Now()
	}
	return &ReferenceWithTimezone{
		instant:        instant,
		timezoneOffset: timezoneOffset,
	}
}

// FromDate creates a ReferenceWithTimezone from a Date.
func (r *ReferenceWithTimezone) FromDate(date time.Time) *ReferenceWithTimezone {
	return NewReferenceWithTimezone(date, nil)
}

// FromInput creates a ReferenceWithTimezone from either a ParsingReference or a time.Time.
// It also handles timezone conversion using the provided timezoneOverrides.
func FromInput(input interface{}, timezoneOverrides TimezoneAbbrMap) *ReferenceWithTimezone {
	if input == nil {
		return NewReferenceWithTimezone(time.Time{}, nil)
	}

	switch v := input.(type) {
	case time.Time:
		return NewReferenceWithTimezone(v, nil)
	case ParsingReference:
		instant := time.Now()
		if v.Instant != nil {
			instant = *v.Instant
		}

		var timezoneOffset *int
		if v.Timezone != nil {
			offset := toTimezoneOffset(v.Timezone, instant, timezoneOverrides)
			timezoneOffset = &offset
		}

		return NewReferenceWithTimezone(instant, timezoneOffset)
	default:
		return NewReferenceWithTimezone(time.Time{}, nil)
	}
}

// GetDateWithAdjustedTimezone returns a time.Time with the year, month, day, hour, minute, second
// equal to the reference. The output's instant is NOT the reference's instant when the reference's
// and system's timezone are different.
func (r *ReferenceWithTimezone) GetDateWithAdjustedTimezone() time.Time {
	date := r.instant
	if r.timezoneOffset != nil {
		adjustment := r.GetSystemTimezoneAdjustmentMinute(r.instant, nil)
		date = date.Add(time.Duration(-adjustment) * time.Minute)
	}
	return date
}

// GetSystemTimezoneAdjustmentMinute returns the number of minutes difference between
// the system's timezone and the reference timezone.
func (r *ReferenceWithTimezone) GetSystemTimezoneAdjustmentMinute(date time.Time, overrideTimezoneOffset *int) int {
	if date.IsZero() || date.Unix() < 0 {
		// Javascript date timezone calculation got effect when the time epoch < 0
		date = time.Now()
	}

	_, currentOffset := date.Zone()
	currentTimezoneOffset := currentOffset / 60

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
func (r *ReferenceWithTimezone) GetTimezoneOffset() int {
	if r.timezoneOffset != nil {
		return *r.timezoneOffset
	}
	_, offset := r.instant.Zone()
	return offset / 60
}

// Instant returns the reference instant.
func (r *ReferenceWithTimezone) Instant() time.Time {
	return r.instant
}

// ParsingComponents represents a collection of parsed date/time components.
// Components are stored as either "known" (directly parsed) or "implied" (inferred).
type ParsingComponents struct {
	knownValues   map[Component]int
	impliedValues map[Component]int
	reference     *ReferenceWithTimezone
	tags          map[string]bool
}

// NewParsingComponents creates a new ParsingComponents with the given reference.
// It initializes implied values based on the reference date.
func NewParsingComponents(reference *ReferenceWithTimezone, knownComponents map[Component]int) *ParsingComponents {
	pc := &ParsingComponents{
		knownValues:   make(map[Component]int),
		impliedValues: make(map[Component]int),
		reference:     reference,
		tags:          make(map[string]bool),
	}

	if knownComponents != nil {
		for k, v := range knownComponents {
			pc.knownValues[k] = v
		}
	}

	// Set default implied values from reference
	date := reference.GetDateWithAdjustedTimezone()
	pc.Imply(ComponentDay, date.Day())
	pc.Imply(ComponentMonth, int(date.Month()))
	pc.Imply(ComponentYear, date.Year())
	pc.Imply(ComponentHour, 12)
	pc.Imply(ComponentMinute, 0)
	pc.Imply(ComponentSecond, 0)
	pc.Imply(ComponentMillisecond, 0)

	return pc
}

// IsCertain returns true if the component is certain (directly mentioned).
func (pc *ParsingComponents) IsCertain(component Component) bool {
	_, exists := pc.knownValues[component]
	return exists
}

// Get returns the component value for either Certain or Implied values.
// Returns nil if the component is Unknown.
func (pc *ParsingComponents) Get(component Component) *int {
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
func (pc *ParsingComponents) Assign(component Component, value int) *ParsingComponents {
	pc.knownValues[component] = value
	delete(pc.impliedValues, component)
	return pc
}

// Imply sets a component value as implied.
// If the component is already known, this does nothing.
func (pc *ParsingComponents) Imply(component Component, value int) *ParsingComponents {
	if _, exists := pc.knownValues[component]; exists {
		return pc
	}
	pc.impliedValues[component] = value
	return pc
}

// Delete removes components from both known and implied values.
func (pc *ParsingComponents) Delete(components ...Component) *ParsingComponents {
	for _, component := range components {
		delete(pc.knownValues, component)
		delete(pc.impliedValues, component)
	}
	return pc
}

// Clone creates a deep copy of the ParsingComponents.
func (pc *ParsingComponents) Clone() *ParsingComponents {
	clone := &ParsingComponents{
		knownValues:   make(map[Component]int),
		impliedValues: make(map[Component]int),
		reference:     pc.reference,
		tags:          make(map[string]bool),
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
func (pc *ParsingComponents) IsOnlyDate() bool {
	return !pc.IsCertain(ComponentHour) && !pc.IsCertain(ComponentMinute) && !pc.IsCertain(ComponentSecond)
}

// IsOnlyTime returns true if only time components are certain (no date components).
func (pc *ParsingComponents) IsOnlyTime() bool {
	return !pc.IsCertain(ComponentWeekday) && !pc.IsCertain(ComponentDay) &&
		!pc.IsCertain(ComponentMonth) && !pc.IsCertain(ComponentYear)
}

// IsOnlyWeekdayComponent returns true if only weekday is certain without day or month.
func (pc *ParsingComponents) IsOnlyWeekdayComponent() bool {
	return pc.IsCertain(ComponentWeekday) && !pc.IsCertain(ComponentDay) && !pc.IsCertain(ComponentMonth)
}

// IsDateWithUnknownYear returns true if month is certain but year is not.
func (pc *ParsingComponents) IsDateWithUnknownYear() bool {
	return pc.IsCertain(ComponentMonth) && !pc.IsCertain(ComponentYear)
}

// IsValidDate validates that the components form a valid date.
func (pc *ParsingComponents) IsValidDate() bool {
	date := pc.dateWithoutTimezoneAdjustment()

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
func (pc *ParsingComponents) Date() time.Time {
	date := pc.dateWithoutTimezoneAdjustment()

	timezoneOffsetVal := pc.Get(ComponentTimezoneOffset)
	timezoneAdjustment := pc.reference.GetSystemTimezoneAdjustmentMinute(date, timezoneOffsetVal)

	return date.Add(time.Duration(timezoneAdjustment) * time.Minute)
}

// dateWithoutTimezoneAdjustment creates a time.Time from components without timezone adjustment.
func (pc *ParsingComponents) dateWithoutTimezoneAdjustment() time.Time {
	year := getValueOrDefault(pc.Get(ComponentYear), 2000)
	month := getValueOrDefault(pc.Get(ComponentMonth), 1)
	day := getValueOrDefault(pc.Get(ComponentDay), 1)
	hour := getValueOrDefault(pc.Get(ComponentHour), 0)
	minute := getValueOrDefault(pc.Get(ComponentMinute), 0)
	second := getValueOrDefault(pc.Get(ComponentSecond), 0)
	millisecond := getValueOrDefault(pc.Get(ComponentMillisecond), 0)

	// Use the reference date's location to avoid timezone conversion issues
	location := time.Local
	if pc.reference != nil && !pc.reference.instant.IsZero() {
		location = pc.reference.instant.Location()
	}

	date := time.Date(year, time.Month(month), day, hour, minute, second, millisecond*1000000, location)

	// Handle years < 100 properly (time.Date can interpret them differently)
	if year < 100 {
		date = time.Date(year, time.Month(month), day, hour, minute, second, millisecond*1000000, location)
	}

	return date
}

// AddTag adds a debugging tag to the components.
func (pc *ParsingComponents) AddTag(tag string) *ParsingComponents {
	pc.tags[tag] = true
	return pc
}

// Tags returns all debugging tags.
func (pc *ParsingComponents) Tags() map[string]bool {
	tags := make(map[string]bool)
	for k, v := range pc.tags {
		tags[k] = v
	}
	return tags
}

// String returns a string representation for debugging.
func (pc *ParsingComponents) String() string {
	tagList := make([]string, 0, len(pc.tags))
	for tag := range pc.tags {
		tagList = append(tagList, tag)
	}

	return fmt.Sprintf("[ParsingComponents {tags: %v, knownValues: %v, impliedValues: %v}]",
		tagList, pc.knownValues, pc.impliedValues)
}

// Reference returns the reference.
func (pc *ParsingComponents) Reference() *ReferenceWithTimezone {
	return pc.reference
}

// CreateRelativeFromReference creates a ParsingComponents from a duration relative to the reference.
// It handles date-only durations (implies time) and time durations (assigns both date and time).
// This is used for parsing relative expressions like "in 3 days", "2 hours ago", etc.
func CreateRelativeFromReference(reference *ReferenceWithTimezone, duration Duration) *ParsingComponents {
	if duration == nil {
		duration = EmptyDuration
	}

	date := AddDuration(reference.GetDateWithAdjustedTimezone(), duration)

	components := NewParsingComponents(reference, nil)
	components.AddTag("result/relativeDate")

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
		AssignSimilarTime(components, date)
		AssignSimilarDate(components, date)
		components.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	} else {
		// Duration is date-only - imply time components
		ImplySimilarTime(components, date)
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
func (pc *ParsingComponents) AddDurationAsImplied(duration Duration) *ParsingComponents {
	// Get the current date from this component
	currentDate := pc.Date()

	// Add the duration
	newDate := AddDuration(currentDate, duration)

	// Imply the new date components
	ImplySimilarDate(pc, newDate)
	ImplySimilarTime(pc, newDate)

	return pc
}

// ParsingResult represents a parsed result containing date/time information.
type ParsingResult struct {
	reference *ReferenceWithTimezone
	refDate   time.Time
	index     int
	text      string
	start     *ParsingComponents
	end       *ParsingComponents
}

// NewParsingResult creates a new ParsingResult.
func NewParsingResult(reference *ReferenceWithTimezone, index int, text string, start, end *ParsingComponents) *ParsingResult {
	if start == nil {
		start = NewParsingComponents(reference, nil)
	}

	return &ParsingResult{
		reference: reference,
		refDate:   reference.Instant(),
		index:     index,
		text:      text,
		start:     start,
		end:       end,
	}
}

// Clone creates a deep copy of the ParsingResult.
func (pr *ParsingResult) Clone() *ParsingResult {
	var startClone *ParsingComponents
	if pr.start != nil {
		startClone = pr.start.Clone()
	}

	var endClone *ParsingComponents
	if pr.end != nil {
		endClone = pr.end.Clone()
	}

	return NewParsingResult(pr.reference, pr.index, pr.text, startClone, endClone)
}

// Date returns a time.Time object created from the start components.
func (pr *ParsingResult) Date() time.Time {
	return pr.start.Date()
}

// AddTag adds a debugging tag to both start and end components.
func (pr *ParsingResult) AddTag(tag string) *ParsingResult {
	pr.start.AddTag(tag)
	if pr.end != nil {
		pr.end.AddTag(tag)
	}
	return pr
}

// Tags returns combined debugging tags from start and end components.
func (pr *ParsingResult) Tags() map[string]bool {
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
func (pr *ParsingResult) String() string {
	tagList := make([]string, 0, len(pr.Tags()))
	for tag := range pr.Tags() {
		tagList = append(tagList, tag)
	}

	return fmt.Sprintf("[ParsingResult {index: %d, text: '%s', tags: %v}]",
		pr.index, pr.text, tagList)
}

// RefDate returns the reference date used for parsing.
func (pr *ParsingResult) RefDate() time.Time {
	return pr.refDate
}

// Index returns the position in the input text.
func (pr *ParsingResult) Index() int {
	return pr.index
}

// SetIndex sets the position in the input text.
// This is used internally by parsers and the chrono executor.
func (pr *ParsingResult) SetIndex(index int) {
	pr.index = index
}

// ParsingResultWithBoundary wraps ParsingComponents with boundary information
// This is used internally to communicate the adjusted text (without boundary) to chrono.go
// when parsers using AbstractParserWithWordBoundary return ParsingComponents.
type ParsingResultWithBoundary struct {
	Components   *ParsingComponents
	AdjustedText string
	BoundaryLen  int
}

// Text returns the matched text from the input.
func (pr *ParsingResult) Text() string {
	return pr.text
}

// Start returns the starting date/time components.
func (pr *ParsingResult) Start() ParsedComponents {
	return pr.start
}

// End returns the ending date/time components.
func (pr *ParsingResult) End() ParsedComponents {
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

// toTimezoneOffset converts various timezone representations to an offset in minutes.
// This is a simplified version - a full implementation would handle timezone names
// and ambiguous timezones.
func toTimezoneOffset(timezone interface{}, instant time.Time, timezoneOverrides TimezoneAbbrMap) int {
	if timezone == nil {
		return 0
	}

	switch v := timezone.(type) {
	case int:
		return v
	case string:
		// Try to look up in overrides first
		if timezoneOverrides != nil {
			if offset, exists := timezoneOverrides[v]; exists {
				switch o := offset.(type) {
				case int:
					return o
				case AmbiguousTimezoneMap:
					// For ambiguous timezones, we need to determine if DST is in effect
					year := instant.Year()
					dstStart := o.DstStart(year)
					dstEnd := o.DstEnd(year)

					if instant.After(dstStart) && instant.Before(dstEnd) {
						return o.TimezoneOffsetDuringDst
					}
					return o.TimezoneOffsetNonDst
				}
			}
		}

		// Try to parse as location name
		loc, err := time.LoadLocation(v)
		if err == nil {
			_, offset := instant.In(loc).Zone()
			return offset / 60
		}
	}

	return 0
}
