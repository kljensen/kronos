// Package main demonstrates British English date parsing.
//
// This example shows:
//   - Using the GB (British) parser with DMY date order
//   - Parsing day/month/year format
//   - Comparing US vs British date interpretation
package main

import (
	"fmt"
	"log"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/en"
)

func main() {
	fmt.Println("=== British English Date Parsing ===")

	// Example 1: Using the GB parser
	fmt.Println("Example 1: GB parser with DMY format")
	gbParser := en.GBParser()

	britishDates := []string{
		"15/3/2024",      // 15th March 2024
		"1/12/2024",      // 1st December 2024
		"25/12/2024",     // 25th December 2024
		"15th March 2024",
		"1st December",
	}

	for _, expr := range britishDates {
		date, err := gbParser.ParseDate(expr)
		if err != nil {
			fmt.Printf("  %-20s => ERROR: %v\n", expr, err)
			continue
		}
		if date != nil {
			fmt.Printf("  %-20s => %s\n", expr, date.Format("2 January 2006"))
		}
	}
	fmt.Println()

	// Example 2: US vs British interpretation
	fmt.Println("Example 2: US vs British interpretation of ambiguous dates")

	ambiguousDates := []string{
		"3/5/2024",
		"1/2/2024",
		"12/11/2024",
	}

	usParser := en.New().DateOrder(kronos.DateOrderMDY)
	gbParser = en.GBParser() // DMY order

	fmt.Println("  Date       | US (M/D/Y)         | British (D/M/Y)")
	fmt.Println("  -----------|--------------------|-----------------")

	for _, expr := range ambiguousDates {
		usDate, _ := usParser.ParseDate(expr)
		gbDate, _ := gbParser.ParseDate(expr)

		usStr := "(parse error)"
		if usDate != nil {
			usStr = usDate.Format("2 Jan 2006")
		}

		gbStr := "(parse error)"
		if gbDate != nil {
			gbStr = gbDate.Format("2 Jan 2006")
		}

		fmt.Printf("  %-10s | %-18s | %s\n", expr, usStr, gbStr)
	}
	fmt.Println()

	// Example 3: Creating a custom DMY parser
	fmt.Println("Example 3: Custom DMY parser with additional settings")

	customParser := en.New().
		DateOrder(kronos.DateOrderDMY).
		PreferPast()

	expr := "15/3" // 15th March, but which year?
	date, err := customParser.ParseDate(expr)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	fmt.Printf("Expression: '%s'\n", expr)
	fmt.Printf("Config: DMY order, prefer past\n")
	if date != nil {
		fmt.Printf("Result: %s\n", date.Format("Monday, 2 January 2006"))
	}
	fmt.Println()

	// Example 4: British casual expressions
	fmt.Println("Example 4: British casual expressions")
	// These work the same in both US and British English

	casualExpressions := []string{
		"yesterday",
		"tomorrow",
		"last Monday",
		"next Friday",
		"fortnight", // British term for "two weeks"
	}

	fmt.Println("British casual expressions:")
	for _, expr := range casualExpressions {
		date, _ := gbParser.ParseDate(expr)
		if date != nil {
			fmt.Printf("  %-15s => %s\n", expr, date.Format("Mon, 2 Jan 2006"))
		} else {
			fmt.Printf("  %-15s => (no match)\n", expr)
		}
	}
	fmt.Println()

	// Example 5: Combining British format with time
	fmt.Println("Example 5: British dates with time")

	britishDateTimes := []string{
		"15/3/2024 at 3pm",
		"1/12/2024 14:30",
		"25th December at noon",
	}

	for _, expr := range britishDateTimes {
		date, _ := gbParser.ParseDate(expr)
		if date != nil {
			fmt.Printf("  %-25s => %s\n", expr, date.Format("Mon, 2 Jan 2006 15:04"))
		}
	}
}
