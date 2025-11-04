// Package main demonstrates basic usage of the kronos date parser.
//
// This example shows:
//   - Parsing a single date with ParseDateSimple
//   - Parsing multiple dates with ParseSimple
//   - Accessing result properties (text, index, date)
package main

import (
	"fmt"
	"log"

	"github.com/kljensen/kronos/en"
)

func main() {
	fmt.Println("=== Basic Kronos Usage ===")

	// Example 1: Parse a single date
	fmt.Println("Example 1: Parse a single date")
	date, err := en.ParseDateSimple("tomorrow at 3pm")
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	if date != nil {
		fmt.Printf("Parsed date: %v\n\n", date)
	}

	// Example 2: Parse multiple dates in text
	fmt.Println("Example 2: Parse multiple dates in text")
	text := "Meet me tomorrow at 3pm or next Friday at noon"
	results, err := en.ParseSimple(text)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	fmt.Printf("Found %d date(s) in: \"%s\"\n", len(results), text)
	for i, r := range results {
		fmt.Printf("  [%d] Position %d: '%s' => %v\n",
			i+1, r.Index(), r.Text(), r.Date())
	}
	fmt.Println()

	// Example 3: Various date formats
	fmt.Println("Example 3: Various date formats")
	examples := []string{
		"yesterday",
		"last Monday",
		"3 days ago",
		"in 2 weeks",
		"March 15, 2024",
		"3/15/2024",
		"2024-03-15",
		"next Friday at 2:30pm",
		"noon",
		"midnight",
	}

	for _, expr := range examples {
		date, err := en.ParseDateSimple(expr)
		if err != nil {
			fmt.Printf("  %-25s => ERROR: %v\n", expr, err)
			continue
		}
		if date != nil {
			fmt.Printf("  %-25s => %v\n", expr, date)
		} else {
			fmt.Printf("  %-25s => (no date found)\n", expr)
		}
	}
}
