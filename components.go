package kronos

import "time"

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
}
