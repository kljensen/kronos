//nolint:staticcheck // SA1019: Must use deprecated types during transition
package kronos_test

import (
	"fmt"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/en"
)

// Example_parserBuilder demonstrates the new builder pattern API for parsing dates.
func Example_parserBuilder() {
	// Basic usage with the builder pattern
	refDate := time.Date(2020, 11, 15, 12, 0, 0, 0, time.UTC)

	// Parse with default settings
	results, err := en.New().
		WithReferenceDate(refDate).
		Parse("tomorrow at 3pm")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(results) > 0 {
		fmt.Printf("Parsed date: %s\n", results[0].Date().Format("2006-01-02 15:04"))
	}
	// Output: Parsed date: 2020-11-16 15:00
}

// Example_parserBuilder_preferPast demonstrates using PreferPast with the builder.
func Example_parserBuilder_preferPast() {
	refDate := time.Date(2020, 11, 15, 12, 0, 0, 0, time.UTC)

	// Parse "March" with PreferPast - should be March 2020
	date, err := en.New().
		WithReferenceDate(refDate).
		PreferPast().
		ParseDate("March")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if date != nil {
		fmt.Printf("Year: %d, Month: %s\n", date.Year(), date.Month())
	}
	// Output: Year: 2020, Month: March
}

// Example_parserBuilder_preferFuture demonstrates using PreferFuture with the builder.
func Example_parserBuilder_preferFuture() {
	refDate := time.Date(2020, 11, 15, 12, 0, 0, 0, time.UTC)

	// Parse "March" with PreferFuture - should be March 2021
	date, err := en.New().
		WithReferenceDate(refDate).
		PreferFuture().
		ParseDate("March")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if date != nil {
		fmt.Printf("Year: %d, Month: %s\n", date.Year(), date.Month())
	}
	// Output: Year: 2021, Month: March
}

// Example_parserBuilder_dateOrder demonstrates setting date order for ambiguous formats.
func Example_parserBuilder_dateOrder() {
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

	// Parse with US format (MDY)
	date, err := en.New().
		WithReferenceDate(refDate).
		DateOrder(kronos.DateOrderMDY).
		ParseDate("03/15/2020")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if date != nil {
		fmt.Printf("Month: %d, Day: %d\n", date.Month(), date.Day())
	}
	// Output: Month: 3, Day: 15
}

// Example_parserBuilder_strictMode demonstrates strict parsing mode.
func Example_parserBuilder_strictMode() {
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

	// Strict mode validates more strictly
	date, err := en.NewStrict().
		WithReferenceDate(refDate).
		ParseDate("2020-03-15")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if date != nil {
		fmt.Printf("Parsed: %s\n", date.Format("2006-01-02"))
	}
	// Output: Parsed: 2020-03-15
}

// Example_parserBuilder_gbFormat demonstrates UK date format parsing.
func Example_parserBuilder_gbFormat() {
	refDate := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

	// GB format uses day/month/year
	date, err := en.NewGB().
		WithReferenceDate(refDate).
		ParseDate("15/03/2020")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if date != nil {
		fmt.Printf("Day: %d, Month: %d\n", date.Day(), date.Month())
	}
	// Output: Day: 15, Month: 3
}

// Example_parserSimple demonstrates the simple convenience functions.
func Example_parserSimple() {
	// Simple parsing with default settings and current time as reference
	results, err := en.ParseSimple("tomorrow")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(results) > 0 {
		// Tomorrow will be in the future
		isFuture := results[0].Date().After(time.Now())
		fmt.Printf("Is future: %t\n", isFuture)
	}
	// Output: Is future: true
}

// Example_packageLevelParse demonstrates the package-level Parse function.
func Example_packageLevelParse() {
	// Using the new builder API
	results, err := en.New().Parse("tomorrow")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(results) > 0 {
		// Tomorrow will be in the future
		isFuture := results[0].Date().After(time.Now())
		fmt.Printf("Parsed successfully: %t\n", isFuture)
	}
	// Output: Parsed successfully: true
}
