// Package main demonstrates the builder API for kronos.
// This shows how to use the fluent API for common parsing configurations.
package main

import (
	"fmt"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/en"
)

func main() {
	demonstrateSimpleAPI()
	demonstrateCommonOptions()
	demonstrateCombinedOptions()
}

// demonstrateSimpleAPI shows the simplest way to parse dates.
func demonstrateSimpleAPI() {
	fmt.Println("=== Simple API ===")

	// Simplest usage - just parse
	results, err := en.New().Parse("tomorrow at 3pm")
	if err != nil {
		panic(err)
	}
	if len(results) > 0 {
		fmt.Printf("Parsed: %v\n", results[0].Date())
	}

	// Parse a single date
	date, err := en.New().ParseDate("next Monday")
	if err != nil {
		panic(err)
	}
	if date != nil {
		fmt.Printf("Next Monday: %v\n", *date)
	}

	fmt.Println()
}

// demonstrateCommonOptions shows the builder API for common use cases.
func demonstrateCommonOptions() {
	fmt.Println("=== Common Options ===")

	refDate := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	// US format (Month/Day/Year)
	parser := en.New().
		WithReferenceDate(refDate).
		DateOrder(kronos.DateOrderMDY).
		Casual()

	results, _ := parser.Parse("12/25/2024")
	if len(results) > 0 {
		fmt.Printf("US format 12/25/2024: %v\n", results[0].Date())
	}

	// European format (Day/Month/Year)
	parser = en.New().
		DateOrder(kronos.DateOrderDMY).
		PreferFuture() // Prefer future dates

	results, _ = parser.Parse("25/12/2024")
	if len(results) > 0 {
		fmt.Printf("EU format 25/12/2024: %v\n", results[0].Date())
	}

	// Strict mode - only formal patterns
	parser = en.New().
		Strict().
		PreferPast()

	results, _ = parser.Parse("2024-03-15")
	if len(results) > 0 {
		fmt.Printf("Strict mode: %v\n", results[0].Date())
	}

	// British English (DMY order)
	parser = en.NewGB().WithReferenceDate(refDate)

	results, _ = parser.Parse("15/3/2024")
	if len(results) > 0 {
		fmt.Printf("British format 15/3/2024: %v\n", results[0].Date())
	}

	fmt.Println()
}

// demonstrateCombinedOptions shows combining multiple configurations.
func demonstrateCombinedOptions() {
	fmt.Println("=== Combined Options ===")

	refDate := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	// Chain multiple configuration methods
	parser := en.New().
		WithReferenceDate(refDate).
		DateOrder(kronos.DateOrderDMY).
		Strict().
		PreferFuture()

	results, _ := parser.Parse("15/03/2024")
	if len(results) > 0 {
		fmt.Printf("Combined options: %v\n", results[0].Date())
	}

	fmt.Println()
}
