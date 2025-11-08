// Package main demonstrates working with date ranges.
//
// This example shows:
//   - Parsing date range expressions
//   - Accessing start and end components
//   - Detecting whether a result is a range
//   - Working with range dates
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kljensen/kronos/en"
)

func main() {
	fmt.Println("=== Date Range Examples ===")

	// Example 1: Basic date range
	fmt.Println("Example 1: Basic date range")
	expr := "from Monday to Friday"
	results, err := en.ParseSimple(expr)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	if len(results) > 0 {
		r := results[0]
		fmt.Printf("Expression: '%s'\n", expr)
		fmt.Printf("Matched text: '%s'\n", r.Text())

		if r.End() != nil {
			start := r.Start().Date()
			end := r.End().Date()
			fmt.Printf("Range: %v to %v\n", start, end)
			fmt.Printf("Duration: %v\n", end.Sub(start))
		} else {
			fmt.Println("Not a range (single date)")
		}
	}
	fmt.Println()

	// Example 2: Multiple range formats
	fmt.Println("Example 2: Various range formats")
	rangeExpressions := []string{
		"from Monday to Friday",
		"between March 1 and March 15",
		"from 9am to 5pm",
		"March 1 - March 15",
		"Wed, Nov 5, 2025 8:30 PM - 10:45 PM",
		"Thursday, Nov 6 2025 at 9:30 - 3:00 pm EST",
	}

	for _, expr := range rangeExpressions {
		results, _ := en.ParseSimple(expr)
		if len(results) > 0 {
			r := results[0]
			fmt.Printf("\nExpression: '%s'\n", expr)

			if r.End() != nil {
				start := r.Start().Date()
				end := r.End().Date()

				// Format based on the type of range
				if start.Format("2006-01-02") == end.Format("2006-01-02") {
					// Same day, different times
					fmt.Printf("  Time range: %s to %s\n",
						start.Format("3:04pm"),
						end.Format("3:04pm"))
				} else {
					// Different days
					fmt.Printf("  Start: %s\n", start.Format("2006-01-02 15:04"))
					fmt.Printf("  End:   %s\n", end.Format("2006-01-02 15:04"))
					fmt.Printf("  Duration: %v\n", end.Sub(start))
				}
			} else {
				fmt.Printf("  Single date: %v\n", r.Date())
			}
		}
	}
	fmt.Println()

	// Example 3: Detecting and handling ranges
	fmt.Println("Example 3: Detecting and handling ranges")

	testExpressions := []string{
		"tomorrow",              // Single date
		"from Monday to Friday", // Range
		"March 15",              // Single date
		"between 9am and 5pm",   // Range
		"next week",             // Single date (could be range-like)
	}

	for _, expr := range testExpressions {
		results, _ := en.ParseSimple(expr)
		if len(results) > 0 {
			r := results[0]
			if r.End() != nil {
				fmt.Printf("  %-30s => RANGE\n", expr)
			} else {
				fmt.Printf("  %-30s => SINGLE DATE\n", expr)
			}
		}
	}
	fmt.Println()

	// Example 4: Working with range components
	fmt.Println("Example 4: Inspecting range components")

	expr = "from March 1 to March 15"
	results, _ = en.ParseSimple(expr)

	if len(results) > 0 && results[0].End() != nil {
		r := results[0]
		fmt.Printf("Expression: '%s'\n\n", expr)

		// Start components
		startComp := r.Start()
		fmt.Println("Start components:")
		fmt.Printf("  Year:  %d\n", *startComp.Get("year"))
		fmt.Printf("  Month: %d\n", *startComp.Get("month"))
		fmt.Printf("  Day:   %d\n", *startComp.Get("day"))

		// End components
		endComp := r.End()
		fmt.Println("\nEnd components:")
		fmt.Printf("  Year:  %d\n", *endComp.Get("year"))
		fmt.Printf("  Month: %d\n", *endComp.Get("month"))
		fmt.Printf("  Day:   %d\n", *endComp.Get("day"))
	}
	fmt.Println()

	// Example 5: Iterating over date range
	fmt.Println("Example 5: Iterating over date range")

	expr = "from Monday to Friday"
	results, _ = en.ParseSimple(expr)

	if len(results) > 0 && results[0].End() != nil {
		r := results[0]
		start := r.Start().Date()
		end := r.End().Date()

		fmt.Printf("Expression: '%s'\n", expr)
		fmt.Printf("Iterating from %s to %s:\n\n",
			start.Format("Monday, Jan 2"),
			end.Format("Monday, Jan 2"))

		current := start
		for !current.After(end) {
			fmt.Printf("  %s\n", current.Format("Monday, January 2, 2006"))
			current = current.Add(24 * time.Hour)
		}
	}
}
