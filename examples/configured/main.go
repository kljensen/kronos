// Package main demonstrates using the ParserBuilder to configure parsing behavior.
//
// This example shows:
//   - Creating a parser with custom reference date
//   - Using PreferPast/PreferFuture settings
//   - Configuring date order (MDY vs DMY)
//   - Strict vs casual parsing modes
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/en"
)

func main() {
	fmt.Println("=== Configured Parsing Examples ===\n")

	// Example 1: Custom reference date
	fmt.Println("Example 1: Custom reference date")
	refDate := time.Date(2024, time.November, 15, 12, 0, 0, 0, time.UTC)
	fmt.Printf("Reference date: %v\n", refDate)

	parser := en.New().WithReferenceDate(refDate)
	examples := []string{"yesterday", "tomorrow", "next week"}

	for _, expr := range examples {
		date, err := parser.ParseDate(expr)
		if err != nil {
			fmt.Printf("  %-15s => ERROR: %v\n", expr, err)
			continue
		}
		if date != nil {
			fmt.Printf("  %-15s => %v\n", expr, date)
		}
	}
	fmt.Println()

	// Example 2: PreferPast vs PreferFuture
	fmt.Println("Example 2: PreferPast vs PreferFuture")
	// When parsing "March" in November, should we prefer past or future?
	expr := "March"

	pastParser := en.New().
		WithReferenceDate(refDate).
		PreferPast()
	pastDate, _ := pastParser.ParseDate(expr)

	futureParser := en.New().
		WithReferenceDate(refDate).
		PreferFuture()
	futureDate, _ := futureParser.ParseDate(expr)

	currentParser := en.New().
		WithReferenceDate(refDate).
		PreferCurrentPeriod()
	currentDate, _ := currentParser.ParseDate(expr)

	fmt.Printf("  Expression: '%s' (reference: November 2024)\n", expr)
	fmt.Printf("  PreferPast:          %v\n", pastDate)
	fmt.Printf("  PreferFuture:        %v\n", futureDate)
	fmt.Printf("  PreferCurrentPeriod: %v\n", currentDate)
	fmt.Println()

	// Example 3: Date order (MDY vs DMY)
	fmt.Println("Example 3: Date order configuration")
	ambiguousDate := "3/5/2024" // March 5 or May 3?

	mdyParser := en.New().DateOrder(kronos.DateOrderMDY)
	mdyDate, _ := mdyParser.ParseDate(ambiguousDate)

	dmyParser := en.New().DateOrder(kronos.DateOrderDMY)
	dmyDate, _ := dmyParser.ParseDate(ambiguousDate)

	fmt.Printf("  Expression: '%s'\n", ambiguousDate)
	fmt.Printf("  MDY (US):       %v (March 5)\n", mdyDate)
	fmt.Printf("  DMY (European): %v (May 3)\n", dmyDate)
	fmt.Println()

	// Example 4: Strict vs Casual mode
	fmt.Println("Example 4: Strict vs Casual mode")
	// Casual mode recognizes informal expressions
	casualExamples := []string{
		"tomorrow",
		"last week",
		"in 3 days",
		"2024-03-15", // Formal date works in both modes
	}

	casualParser := en.New().Casual()
	strictParser := en.StrictParser()

	fmt.Println("  Casual mode (informal expressions allowed):")
	for _, expr := range casualExamples {
		date, _ := casualParser.ParseDate(expr)
		if date != nil {
			fmt.Printf("    %-15s => %v\n", expr, date)
		} else {
			fmt.Printf("    %-15s => (no match)\n", expr)
		}
	}

	fmt.Println("\n  Strict mode (formal patterns only):")
	for _, expr := range casualExamples {
		date, _ := strictParser.ParseDate(expr)
		if date != nil {
			fmt.Printf("    %-15s => %v\n", expr, date)
		} else {
			fmt.Printf("    %-15s => (no match)\n", expr)
		}
	}
	fmt.Println()

	// Example 5: Chaining multiple configurations
	fmt.Println("Example 5: Chaining configurations")
	complexParser := en.New().
		WithReferenceDate(refDate).
		PreferPast().
		DateOrder(kronos.DateOrderDMY)

	testExpr := "15/3"
	result, err := complexParser.ParseDate(testExpr)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	fmt.Printf("  Expression: '%s'\n", testExpr)
	fmt.Printf("  Config: DMY order, prefer past, ref = Nov 2024\n")
	fmt.Printf("  Result: %v (March 15, 2024)\n", result)
}
