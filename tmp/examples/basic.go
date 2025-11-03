// Package main demonstrates basic usage of the kronos date parser.
package main

import (
	"fmt"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/en"
)

func main() {
	// Use time.Now() as the reference date
	now := time.Now()
	fmt.Printf("Reference time: %s\n\n", now.Format("2006-01-02 15:04:05"))

	// Example 1: Parse a simple date
	fmt.Println("=== Example 1: Simple parsing ===")
	date := en.ParseDate("tomorrow at 3pm", now, &kronos.ParsingOption{})
	if date != nil {
		fmt.Printf("'tomorrow at 3pm' -> %s\n\n", date.Format("2006-01-02 15:04:05"))
	}

	// Example 2: Parse multiple dates from text
	fmt.Println("=== Example 2: Multiple dates ===")
	text := "The event is next Monday at 2pm or Wednesday at 10am"
	results := en.Parse(text, now, &kronos.ParsingOption{})
	fmt.Printf("Text: %s\n", text)
	for i, result := range results {
		fmt.Printf("  [%d] '%s' -> %s\n", i+1, result.Text(), result.Start().Date().Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// Example 3: Date ranges
	fmt.Println("=== Example 3: Date ranges ===")
	text = "Monday to Friday"
	results = en.Parse(text, now, &kronos.ParsingOption{})
	if len(results) > 0 && results[0].End() != nil {
		fmt.Printf("Text: %s\n", text)
		fmt.Printf("  Start: %s\n", results[0].Start().Date().Format("2006-01-02"))
		fmt.Printf("  End:   %s\n", results[0].End().Date().Format("2006-01-02"))
	}
	fmt.Println()

	// Example 4: Casual expressions
	fmt.Println("=== Example 4: Casual expressions ===")
	expressions := []string{
		"today",
		"tomorrow",
		"yesterday",
		"next week",
		"last month",
		"in 3 days",
		"2 hours ago",
	}
	for _, expr := range expressions {
		date := en.ParseDate(expr, now, &kronos.ParsingOption{})
		if date != nil {
			fmt.Printf("%-15s -> %s\n", expr, date.Format("2006-01-02 15:04:05"))
		}
	}
	fmt.Println()

	// Example 5: Formal dates
	fmt.Println("=== Example 5: Formal dates ===")
	formalDates := []string{
		"2024-12-25",
		"Dec 25, 2024",
		"25/12/2024",
		"December 25th, 2024",
		"Wed, 25 Dec 2024",
	}
	for _, dateStr := range formalDates {
		date := en.ParseDate(dateStr, now, &kronos.ParsingOption{})
		if date != nil {
			fmt.Printf("%-25s -> %s\n", dateStr, date.Format("2006-01-02"))
		}
	}
	fmt.Println()

	// Example 6: Time expressions
	fmt.Println("=== Example 6: Time expressions ===")
	times := []string{
		"3pm",
		"3:30pm",
		"15:30",
		"half past 3",
		"quarter to 4",
	}
	for _, timeStr := range times {
		date := en.ParseDate(timeStr, now, &kronos.ParsingOption{})
		if date != nil {
			fmt.Printf("%-15s -> %s\n", timeStr, date.Format("15:04:05"))
		}
	}
	fmt.Println()

	// Example 7: Using strict mode
	fmt.Println("=== Example 7: Strict vs Casual ===")
	casualText := "I'll see you tmr"
	strictText := "Meeting on 2024-12-25"

	casualDate := en.Casual.ParseDate(casualText, now, &kronos.ParsingOption{})
	strictDate := en.Strict.ParseDate(strictText, now, &kronos.ParsingOption{})

	if casualDate != nil {
		fmt.Printf("Casual: '%s' -> %s\n", casualText, casualDate.Format("2006-01-02"))
	}
	if strictDate != nil {
		fmt.Printf("Strict: '%s' -> %s\n", strictText, strictDate.Format("2006-01-02"))
	}
	fmt.Println()

	// Example 8: Complex expressions
	fmt.Println("=== Example 8: Complex expressions ===")
	complex := []string{
		"next Friday at 2pm",
		"2 weeks from now",
		"the day after tomorrow",
		"this Saturday morning",
	}
	for _, expr := range complex {
		date := en.ParseDate(expr, now, &kronos.ParsingOption{})
		if date != nil {
			fmt.Printf("%-30s -> %s\n", expr, date.Format("2006-01-02 15:04:05"))
		}
	}
}
