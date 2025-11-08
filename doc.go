// Package kronos provides natural language date parsing for Go.
//
// Kronos makes it easy to parse human-friendly date expressions like
// "tomorrow at 3pm", "next Friday", "in 2 weeks" into structured time.Time values.
// It's inspired by the excellent chrono JavaScript library and provides a similar
// feature set with a Go-idiomatic API.
//
// # Quick Start
//
// The simplest way to parse dates is using the en package convenience functions:
//
//	import "github.com/kljensen/kronos/en"
//
//	// Parse a single date
//	date, err := en.ParseDateSimple("tomorrow at 3pm")
//
//	// Parse all dates in text
//	results, err := en.ParseSimple("Meet me tomorrow or next Friday")
//
// # Builder Pattern API
//
// For more control over parsing behavior, use the builder pattern:
//
//	parser := en.New().
//	    WithReferenceDate(refDate).
//	    PreferPast().
//	    DateOrder(kronos.DateOrderDMY)
//
//	results, err := parser.Parse("last Monday")
//
// The builder provides a fluent API for configuration:
//
//   - WithReferenceDate(time.Time) - Set the reference date for relative dates
//   - Strict() - Enable strict parsing mode (formal patterns only)
//   - Casual() - Enable casual parsing mode (informal expressions, default)
//   - DateOrder(DateOrder) - Set date component order (MDY, DMY, YMD)
//   - PreferPast() - Prefer past dates when ambiguous
//   - PreferFuture() - Prefer future dates when ambiguous
//   - PreferCurrentPeriod() - Prefer current period (default)
//   - Timezone(string) - Set default timezone
//
// # Working with Results
//
// Each Result represents a parsed date/time expression:
//
//	for _, r := range results {
//	    fmt.Printf("Found '%s' at position %d\n", r.Text(), r.Index())
//	    fmt.Printf("Date: %v\n", r.Date())
//
//	    // Access individual components
//	    comp := r.Start()
//	    if comp.IsCertain(kronos.ComponentHour) {
//	        hour := comp.Get(kronos.ComponentHour)
//	        fmt.Printf("Hour was mentioned: %d\n", *hour)
//	    }
//	}
//
// # Component Certainty
//
// Kronos distinguishes between components that were explicitly mentioned
// and those that were implied from context:
//
//	results, _ := en.ParseSimple("tomorrow at 3pm")
//	comp := results[0].Start()
//
//	comp.IsCertain(ComponentHour)  // true - "3pm" explicitly mentions hour
//	comp.IsCertain(ComponentYear)  // false - year is implied from reference date
//
// # Date Ranges
//
// Kronos can parse date ranges with both start and end dates:
//
//	results, _ := en.ParseSimple("from Monday to Friday")
//	if results[0].End() != nil {
//	    start := results[0].Start().Date()
//	    end := results[0].End().Date()
//	}
//
// # Supported Expressions
//
// Relative dates:
//   - "yesterday", "today", "tomorrow"
//   - "last Monday", "next Friday"
//   - "3 days ago", "in 2 weeks"
//   - "2 hours from now"
//
// Absolute dates:
//   - "March 15, 2024"
//   - "3/15/2024", "15/3/2024" (configurable order)
//   - "2024-03-15" (ISO format)
//   - "15th of March"
//
// Time expressions:
//   - "3pm", "15:30", "3:30:45 PM"
//   - "noon", "midnight"
//   - "14:30 EST"
//
// Combined:
//   - "tomorrow at 3pm"
//   - "March 15 at 14:30"
//   - "next Friday at noon"
//
// # Localization
//
// The en package provides pre-configured parsers for English:
//
//   - en.Casual - Casual US English (recognizes informal expressions)
//   - en.Strict - Strict US English (formal patterns only)
//   - en.GB - British English (DMY date order)
//
// You can create parsers from these configurations:
//
//	parser := kronos.New(en.Casual)
//	parser := kronos.New(en.GB).PreferPast()
//
// Or use the convenience builders:
//
//	parser := en.New()           // Casual US English
//	parser := en.StrictParser()  // Strict US English
//	parser := en.GBParser()      // British English
//
// # Architecture
//
// Kronos uses a multi-stage pipeline:
//
// 1. Parsers scan the input text for date/time patterns
// 2. Each match is converted to Components (year, month, day, etc.)
// 3. Refiners post-process results (e.g., merge "tomorrow at 3pm")
// 4. Final Results include matched text, components, and time.Time values
//
// # Thread Safety
//
// ParserBuilder is not thread-safe. Create a separate builder for each goroutine:
//
//	// Good: one builder per goroutine
//	go func() {
//	    parser := en.New()
//	    results, _ := parser.Parse(text)
//	}()
//
//	// Bad: shared builder across goroutines
//	parser := en.New()
//	go func() {
//	    results, _ := parser.Parse(text1)  // Race condition!
//	}()
//	go func() {
//	    results, _ := parser.Parse(text2)  // Race condition!
//	}()
//
// # Performance
//
// For best performance:
//   - Reuse ParserBuilder instances within a single goroutine
//   - Use ParseDate() instead of Parse() if you only need the first result
//   - Consider using strict mode to reduce parser overhead
//
// # Extensibility
//
// The parser and refiner interfaces are used internally but not part of the public API.
// For most use cases, the pre-configured parsers (Casual, Strict, GB) provide excellent coverage.
// If you need additional parsing capabilities, please open an issue on GitHub.
package kronos
