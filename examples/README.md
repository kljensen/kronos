# Kronos Examples

This directory contains comprehensive, runnable examples demonstrating various features of the Kronos date parsing library.

## Running the Examples

Each example is a standalone Go program. To run an example:

```bash
cd examples/basic
go run main.go
```

Or from the repository root:

```bash
go run ./examples/basic
```

## Available Examples

### basic/
**Basic usage of Kronos**

Demonstrates:
- Parsing a single date with `ParseDateSimple`
- Parsing multiple dates with `ParseSimple`
- Accessing result properties (text, index, date)
- Various date formats (relative dates, absolute dates, times)

Start here if you're new to Kronos.

### configured/
**Using ParserBuilder for custom configuration**

Demonstrates:
- Creating a parser with custom reference date
- Using `PreferPast`/`PreferFuture` settings
- Configuring date order (MDY vs DMY)
- Strict vs casual parsing modes
- Chaining multiple configuration methods

This shows how to customize parsing behavior for your specific needs.

### components/
**Working with date/time components**

Demonstrates:
- Accessing individual components (year, month, day, hour, etc.)
- Distinguishing between certain (explicit) and implied components
- Using `IsCertain()` to check component certainty
- Using `Get()` to retrieve component values
- Conditional logic based on component certainty

Essential for understanding what was actually mentioned in the input vs. what was inferred.

### ranges/
**Parsing and working with date ranges**

Demonstrates:
- Parsing date range expressions
- Accessing start and end components
- Detecting whether a result is a range
- Iterating over date ranges
- Different range formats ("from...to", "between...and")

Useful if you need to handle expressions like "from Monday to Friday".

### british/
**British English date parsing**

Demonstrates:
- Using the GB (British) parser with DMY date order
- Parsing day/month/year format
- Comparing US vs British date interpretation
- Custom DMY parsers with additional settings
- British casual expressions

Important if you're working with European date formats or international users.

## Example Output

Here's what you can expect from the basic example:

```
=== Basic Kronos Usage ===

Example 1: Parse a single date
Parsed date: 2024-03-16 15:00:00 +0000 UTC

Example 2: Parse multiple dates in text
Found 2 date(s) in: "Meet me tomorrow at 3pm or next Friday at noon"
  [1] Position 8: 'tomorrow at 3pm' => 2024-03-16 15:00:00 +0000 UTC
  [2] Position 31: 'next Friday at noon' => 2024-03-22 12:00:00 +0000 UTC

Example 3: Various date formats
  yesterday                 => 2024-03-14 12:00:00 +0000 UTC
  last Monday               => 2024-03-11 12:00:00 +0000 UTC
  3 days ago                => 2024-03-12 12:00:00 +0000 UTC
  ...
```

## Learning Path

We recommend exploring the examples in this order:

1. **basic/** - Get familiar with the API
2. **configured/** - Learn to customize parsing behavior
3. **components/** - Understand component certainty and inspection
4. **ranges/** - Work with date ranges
5. **british/** - Handle different date formats and locales

## Common Patterns

### Quick Single Date Parsing

```go
import "github.com/kljensen/kronos/en"

date, err := en.ParseDateSimple("tomorrow at 3pm")
```

### Configured Parser

```go
parser := en.New().
    WithReferenceDate(refDate).
    PreferPast().
    DateOrder(kronos.DateOrderDMY)

results, err := parser.Parse(text)
```

### Component Inspection

```go
results, _ := en.ParseSimple("tomorrow at 3pm")
comp := results[0].Start()

if comp.IsCertain(kronos.ComponentHour) {
    hour := comp.Get(kronos.ComponentHour)
    // Hour was explicitly mentioned
}
```

### Date Range Detection

```go
results, _ := en.ParseSimple("from Monday to Friday")
if results[0].End() != nil {
    start := results[0].Start().Date()
    end := results[0].End().Date()
    // It's a range!
}
```

## Next Steps

After exploring these examples, check out:
- [Main README](../README.md) for full API documentation
- [MIGRATION.md](../MIGRATION.md) if upgrading from v1
- The test files in the repository for more edge cases and scenarios

## Contributing Examples

If you have a use case that's not covered by these examples, contributions are welcome! Please submit a pull request with:
- A new example directory with a descriptive name
- A standalone `main.go` file
- Clear comments explaining what the example demonstrates
- Update this README to include your example
