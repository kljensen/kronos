// Package main demonstrates working with date/time components.
//
// This example shows:
//   - Accessing individual components (year, month, day, hour, etc.)
//   - Distinguishing between certain (explicit) and implied components
//   - Checking component certainty with IsCertain
//   - Using Get to retrieve component values
package main

import (
	"fmt"
	"log"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/en"
)

func main() {
	fmt.Println("=== Component Inspection Examples ===\n")

	// Example 1: Certain vs Implied components
	fmt.Println("Example 1: Certain vs Implied components")
	expr := "tomorrow at 3pm"
	results, err := en.ParseSimple(expr)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	if len(results) == 0 {
		log.Fatal("No results found")
	}

	comp := results[0].Start()
	fmt.Printf("Expression: '%s'\n", expr)
	fmt.Printf("Parsed date: %v\n\n", comp.Date())

	// Check all main components
	components := []kronos.Component{
		kronos.ComponentYear,
		kronos.ComponentMonth,
		kronos.ComponentDay,
		kronos.ComponentHour,
		kronos.ComponentMinute,
		kronos.ComponentSecond,
	}

	fmt.Println("Component breakdown:")
	for _, c := range components {
		if val := comp.Get(c); val != nil {
			certainty := "implied"
			if comp.IsCertain(c) {
				certainty = "certain"
			}
			fmt.Printf("  %-10s: %2d (%s)\n", c, *val, certainty)
		}
	}
	fmt.Println()

	// Example 2: Different expressions, different certainty
	fmt.Println("Example 2: Comparing different expressions")

	testExpressions := []string{
		"March 15, 2024 at 3:30pm",
		"tomorrow",
		"3pm",
		"March 2024",
	}

	for _, expr := range testExpressions {
		results, _ := en.ParseSimple(expr)
		if len(results) == 0 {
			continue
		}

		comp := results[0].Start()
		fmt.Printf("\nExpression: '%s'\n", expr)

		// Show which components are certain
		certainComponents := []string{}
		for _, c := range components {
			if comp.IsCertain(c) {
				certainComponents = append(certainComponents, string(c))
			}
		}
		fmt.Printf("Certain components: %v\n", certainComponents)

		// Show full date
		fmt.Printf("Full date: %v\n", comp.Date())
	}
	fmt.Println()

	// Example 3: Conditional logic based on component certainty
	fmt.Println("Example 3: Using certainty for conditional logic")

	expr = "next Friday"
	results, _ = en.ParseSimple(expr)
	if len(results) > 0 {
		comp := results[0].Start()
		fmt.Printf("Expression: '%s'\n\n", expr)

		// Check if time was specified
		if comp.IsCertain(kronos.ComponentHour) {
			hour := comp.Get(kronos.ComponentHour)
			minute := comp.Get(kronos.ComponentMinute)
			fmt.Printf("Time was specified: %02d:%02d\n", *hour, *minute)
		} else {
			fmt.Println("Time was NOT specified (implied as noon)")
		}

		// Check if year was specified
		if comp.IsCertain(kronos.ComponentYear) {
			year := comp.Get(kronos.ComponentYear)
			fmt.Printf("Year was specified: %d\n", *year)
		} else {
			year := comp.Get(kronos.ComponentYear)
			fmt.Printf("Year was NOT specified (implied as %d)\n", *year)
		}
	}
	fmt.Println()

	// Example 4: Inspecting all available components
	fmt.Println("Example 4: Complete component inspection")

	expr = "March 15, 2024 at 3:30:45pm EST"
	results, _ = en.ParseSimple(expr)
	if len(results) > 0 {
		comp := results[0].Start()
		fmt.Printf("Expression: '%s'\n", expr)
		fmt.Printf("Parsed date: %v\n\n", comp.Date())

		// All possible components
		allComponents := []kronos.Component{
			kronos.ComponentYear,
			kronos.ComponentMonth,
			kronos.ComponentDay,
			kronos.ComponentWeekday,
			kronos.ComponentHour,
			kronos.ComponentMinute,
			kronos.ComponentSecond,
			kronos.ComponentMillisecond,
			kronos.ComponentMeridiem,
			kronos.ComponentTimezoneOffset,
		}

		fmt.Println("All components:")
		for _, c := range allComponents {
			if val := comp.Get(c); val != nil {
				certainty := "implied"
				if comp.IsCertain(c) {
					certainty = "CERTAIN"
				}
				fmt.Printf("  %-18s: %5d  [%s]\n", c, *val, certainty)
			} else {
				fmt.Printf("  %-18s: (not set)\n", c)
			}
		}
	}
}
