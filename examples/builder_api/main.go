// Package main demonstrates the new builder API for kronos.
// This shows how to use the simplified API for common options
// and the experimental package for advanced options.
package main

import (
	"fmt"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/en"
	"github.com/kljensen/kronos/experimental"
)

func main() {
	demonstrateSimpleAPI()
	demonstrateCommonOptions()
	demonstrateExperimentalOptions()
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
	fmt.Println("=== Common Options (Main Package) ===")

	refDate := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	// US format (Month/Day/Year)
	parser := en.New().
		WithReferenceDate(refDate).
		WithDateOrder(kronos.DateOrderMDY).
		Casual()

	results, _ := parser.Parse("12/25/2024")
	if len(results) > 0 {
		fmt.Printf("US format 12/25/2024: %v\n", results[0].Date())
	}

	// European format (Day/Month/Year)
	parser = en.New().
		WithDateOrder(kronos.DateOrderDMY).
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

	fmt.Println()
}

// demonstrateExperimentalOptions shows advanced options via the experimental package.
func demonstrateExperimentalOptions() {
	fmt.Println("=== Advanced Options (Experimental Package) ===")

	// Custom parser configuration
	parser := en.New().
		WithOption(experimental.WithMaxParsers(3)).
		WithOption(experimental.WithTimeout(5 * time.Second))

	results, _ := parser.Parse("tomorrow")
	if len(results) > 0 {
		fmt.Printf("With max parsers: %v\n", results[0].Date())
	}

	// Skip common words
	parser = en.New().
		WithOption(experimental.WithSkipTokens("at", "on"))

	results, _ = parser.Parse("on March 15 at 3pm")
	if len(results) > 0 {
		fmt.Printf("With skip tokens: %v\n", results[0].Date())
	}

	// Require specific parts
	parser = en.New().
		WithOption(experimental.WithRequiredParts("year", "month", "day"))

	// This will parse "2024-03-15" but might not parse "tomorrow"
	results, _ = parser.Parse("2024-03-15")
	if len(results) > 0 {
		fmt.Printf("With required parts: %v\n", results[0].Date())
	}

	// Custom relative base
	baseTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	parser = en.New().
		WithOption(experimental.WithRelativeBase(baseTime))

	results, _ = parser.Parse("in 3 days")
	if len(results) > 0 {
		fmt.Printf("With custom base: %v\n", results[0].Date())
	}

	fmt.Println()
}

// demonstrateCombinedOptions shows combining common and advanced options.
func demonstrateCombinedOptions() {
	fmt.Println("=== Combined Options ===")

	refDate := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	// Mix common and experimental options
	parser := en.New().
		WithReferenceDate(refDate).
		WithDateOrder(kronos.DateOrderDMY).
		Strict().
		PreferFuture().
		WithOption(experimental.WithMaxParsers(5)).
		WithOption(experimental.WithNormalization(true)).
		WithOption(experimental.WithDayPreference(kronos.DayPreferFirst))

	results, _ := parser.Parse("15/03/2024")
	if len(results) > 0 {
		fmt.Printf("Combined options: %v\n", results[0].Date())
	}

	fmt.Println()
}
