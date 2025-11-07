package kronos

import "time"

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
