package kronos

import "time"

// Component represents a date/time component that can be parsed.
// Components are used as keys in maps, so they are string constants.
type Component string

// Date/time component constants
const (
	// ComponentYear represents the year component
	ComponentYear    Component = "year"
	ComponentMonth   Component = "month"
	ComponentDay     Component = "day"
	ComponentWeekday Component = "weekday"

	// ComponentHour represents the hour component
	ComponentHour        Component = "hour"
	ComponentMinute      Component = "minute"
	ComponentSecond      Component = "second"
	ComponentMillisecond Component = "millisecond"
	ComponentMicrosecond Component = "microsecond"
	ComponentNanosecond  Component = "nanosecond"

	// ComponentMeridiem represents AM/PM
	ComponentMeridiem       Component = "meridiem"
	ComponentTimezoneOffset Component = "timezoneOffset"
)

// Components represents the individual date/time parts of a parsed expression.
// It provides access to components like year, month, day, hour, minute, etc.,
// and distinguishes between components that were explicitly mentioned (certain)
// and those that were implied or filled in by the parser.
//
// This interface is the primary way to inspect the details of parsed date/time values.
//
// Component Certainty:
//
// Each component has a certainty level:
//   - Certain (or Known): The component was explicitly mentioned in the input text.
//     For example, in "March 15, 2024 at 3pm", the month (3), day (15), year (2024),
//     and hour (15) are certain.
//   - Implied: The component was not mentioned but was inferred from context or defaults.
//     For example, in "March 15", the year might be implied from the reference date,
//     and the hour/minute/second are implied defaults (typically 12:00:00).
//   - Unknown: The component is not available. Get() returns nil for these.
//
// Example:
//
//	parser := kronos.New(en.Casual)
//	results, _ := parser.Parse("tomorrow at 3pm")
//	comp := results[0].Start()
//
//	// Check if hour was explicitly mentioned
//	if comp.IsCertain(kronos.ComponentHour) {
//	    hour := comp.Get(kronos.ComponentHour)
//	    fmt.Printf("Hour was mentioned: %d\n", *hour)
//	}
//
//	// Get year (likely implied, not explicitly mentioned in "tomorrow")
//	if year := comp.Get(kronos.ComponentYear); year != nil {
//	    fmt.Printf("Year: %d (certain=%v)\n", *year, comp.IsCertain(kronos.ComponentYear))
//	}
//
//	// Convert to time.Time
//	date := comp.Date()
type Components interface {
	// Get returns the value of the specified component.
	// Returns nil if the component is not present (neither certain nor implied).
	//
	// For components that are either certain or implied, this returns a pointer
	// to the component's integer value. Use IsCertain to distinguish between
	// certain and implied components.
	//
	// Example:
	//
	//	if hour := comp.Get(kronos.ComponentHour); hour != nil {
	//	    fmt.Printf("Hour: %d\n", *hour)
	//	}
	Get(component Component) *int

	// IsCertain returns true if the component was explicitly mentioned in the input.
	// Returns false if the component was implied or is not present.
	//
	// This allows you to distinguish between components that were actually parsed
	// from the input vs. those that were filled in by the parser based on defaults
	// or the reference date.
	//
	// Example:
	//
	//	// For input "tomorrow at 3pm"
	//	comp.IsCertain(ComponentHour)  // true - "3pm" explicitly mentions hour
	//	comp.IsCertain(ComponentDay)   // true - "tomorrow" explicitly mentions a day
	//	comp.IsCertain(ComponentYear)  // false - year is implied from reference date
	//	comp.IsCertain(ComponentMinute) // false - minute is implied as 0
	IsCertain(component Component) bool

	// Date returns a time.Time object constructed from the components.
	// This combines all certain and implied components into a concrete date/time.
	//
	// The returned time.Time includes:
	//   - All certain components (explicitly parsed from input)
	//   - All implied components (filled in based on reference date or defaults)
	//   - Timezone adjustments if a timezone was specified
	//
	// Example:
	//
	//	date := comp.Date()
	//	fmt.Printf("Parsed date: %v\n", date)
	Date() time.Time

	// Tags returns metadata tags for these components.
	// Tags are used internally for debugging and testing to track how components were parsed.
	// For example, tags might include "result/relativeDate", "result/casualTime", "forward".
	//
	// This method is primarily useful for testing and debugging parser behavior.
	Tags() map[string]bool
}

// Result represents a parsed date/time expression found in text.
// It provides access to the matched text, its position, and the parsed date/time value.
// A Result can represent either a single date/time or a range (with Start and End components).
//
// This is the primary interface for working with parsed results in the new builder-based API.
// It exposes only what users need while hiding internal implementation details.
//
// Example:
//
//	parser := kronos.New(en.Casual)
//	results, _ := parser.Parse("Meet me tomorrow at 3pm")
//	for _, r := range results {
//	    fmt.Printf("Found '%s' at index %d\n", r.Text(), r.Index())
//	    fmt.Printf("Date: %v\n", r.Date())
//	    if r.Start().IsCertain(kronos.ComponentHour) {
//	        hour := r.Start().Get(kronos.ComponentHour)
//	        fmt.Printf("Hour explicitly mentioned: %d\n", *hour)
//	    }
//	}
type Result interface {
	// Text returns the matched text from the input.
	// This is the substring that was recognized as a date/time expression.
	//
	// Example: For input "tomorrow at 3pm", this might return "tomorrow at 3pm"
	Text() string

	// Index returns the position in the input text where this result was found.
	// The index is 0-based and represents the starting position of the matched text.
	//
	// Example: For input "Meet me tomorrow", this would return 8
	Index() int

	// Date returns a time.Time object created from the start components.
	// This is a convenience method equivalent to calling Start().Date().
	// For date ranges, this returns the start date of the range.
	//
	// The returned time.Time includes all parsed and implied components.
	// Components that were not mentioned are filled in with reasonable defaults
	// based on the reference date.
	Date() time.Time

	// Start returns the starting date/time components.
	// For a single date/time (not a range), this contains all the parsed information.
	// For a date range, this represents the beginning of the range.
	//
	// Use the Components interface to inspect individual date/time parts
	// and determine which were explicitly mentioned vs. implied.
	Start() Components

	// End returns the ending date/time components for a range.
	// Returns nil for a single date/time (not a range).
	//
	// Example: For "from Monday to Friday", Start() returns Monday's components
	// and End() returns Friday's components.
	End() Components

	// Tags returns metadata tags for this result.
	// Tags are used internally for debugging and testing to track how a result was parsed.
	// For example, tags might include "parser/iso", "refiner/merged", "result/relativeDate".
	//
	// This method is primarily useful for testing and debugging parser behavior.
	Tags() map[string]bool
}
