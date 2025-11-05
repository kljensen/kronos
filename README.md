<div align="center">
  <img src="doc/gopher.png" alt="Kronos Gopher" width="100"/>

  # Kronos

[![Go Reference](https://img.shields.io/badge/go.dev-reference-007d9c?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/kljensen/kronos)
[![Go Report Card](https://goreportcard.com/badge/github.com/kljensen/kronos?style=for-the-badge)](https://goreportcard.com/report/github.com/kljensen/kronos)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kljensen/kronos?style=for-the-badge&logo=go)](https://github.com/kljensen/kronos)

</div>

---

Kronos is a Go package for extracting structured datetimes from natural language.
It is heavily inspired by the excellent
[chrono](https://github.com/wanasit/chrono) JavaScript library.

## Features

Using Kronos, you can turn text like "see you Thursday at 5:30pm", "4/5/2023 around 6",
or "Thu, November 6th, 2025, 8:00 AM EST" in to Golang structures. Kronos will keep
track of ambiguity (such as am/pm in the "around 6" example above).

- **Natural Language Parsing**: Understands expressions like "tomorrow at 3pm", "next Friday", "in 2 weeks"
- **Flexible Date Formats**: Supports ISO dates, slash dates, month names, and more
- **Relative Dates**: Handles "yesterday", "last week", "3 days ago", etc.
- **Date Ranges**: Parses expressions like "from Monday to Friday"
- **Time Expressions**: Recognizes times like "3pm", "14:30", "noon", "midnight"
- **Timezone Support**: Parses timezone abbreviations and offsets
- **Configurable**: Control date order (MDY vs DMY), prefer past/future dates, strict/casual parsing
- **Multiple Locales**: Built-in support for US English and British English
- **Component Inspection**: Distinguish between explicitly mentioned vs. implied date parts
- **Builder Pattern API**: Fluent, type-safe configuration

## Installation

```bash
go get github.com/kljensen/kronos
```

## Quick Start

### Basic Usage

The simplest way to parse a date is using the `en` package convenience functions:

```go
package main

import (
    "fmt"
    "github.com/kljensen/kronos/en"
)

func main() {
    // Parse a single date
    date, err := en.ParseDateSimple("tomorrow at 3pm")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Date: %v\n", date)

    // Parse all dates in text
    results, err := en.ParseSimple("Meet me tomorrow at 3pm or next Friday at noon")
    if err != nil {
        panic(err)
    }
    for _, r := range results {
        fmt.Printf("Found '%s' at position %d: %v\n", r.Text(), r.Index(), r.Date())
    }
}
```

### Builder Pattern API

For more control, use the builder pattern:

```go
package main

import (
    "fmt"
    "time"
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

func main() {
    // Create a parser with custom configuration
    parser := en.New().
        WithReferenceDate(time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)).
        PreferPast().
        DateOrder(kronos.DateOrderDMY)

    results, err := parser.Parse("last Monday")
    if err != nil {
        panic(err)
    }

    for _, r := range results {
        fmt.Printf("Text: %s\n", r.Text())
        fmt.Printf("Date: %v\n", r.Date())

        // Inspect components
        comp := r.Start()
        if comp.IsCertain(kronos.ComponentDay) {
            day := comp.Get(kronos.ComponentDay)
            fmt.Printf("Day was explicitly mentioned: %d\n", *day)
        }
    }
}
```

### Working with Components

Kronos lets you inspect which date/time parts were explicitly mentioned vs. implied:

```go
results, _ := en.ParseSimple("tomorrow at 3pm")
comp := results[0].Start()

// Check what was explicitly mentioned
if comp.IsCertain(kronos.ComponentHour) {
    hour := comp.Get(kronos.ComponentHour)
    fmt.Printf("Hour was mentioned: %d\n", *hour)  // 15
}

// Year is implied from reference date, not mentioned in "tomorrow"
if !comp.IsCertain(kronos.ComponentYear) {
    year := comp.Get(kronos.ComponentYear)
    fmt.Printf("Year was implied: %d\n", *year)
}
```

### Date Ranges

Kronos can parse date ranges and provides access to both start and end dates:

```go
results, _ := en.ParseSimple("from Monday to Friday")
if results[0].End() != nil {
    start := results[0].Start().Date()
    end := results[0].End().Date()
    fmt.Printf("Range: %v to %v\n", start, end)
}
```

## API Overview

### Core Types

- **`ParserBuilder`**: The main entry point for parsing. Create with `en.New()` or `kronos.New(chrono)`.
- **`Result`**: Represents a parsed date/time expression. Access via `Text()`, `Index()`, `Date()`, `Start()`, `End()`.
- **`Components`**: Represents individual date/time parts. Methods: `Get()`, `IsCertain()`, `Date()`.
- **`Component`**: Constants for date/time components (`ComponentYear`, `ComponentMonth`, `ComponentDay`, etc.).

### Builder Methods

Configure parsing behavior by chaining builder methods:

```go
parser := en.New().
    WithReferenceDate(refDate).  // Set reference date for relative dates
    Strict().                     // Enable strict parsing mode
    Casual().                     // Enable casual parsing mode (default)
    DateOrder(order).            // Set date component order (MDY, DMY, YMD)
    PreferPast().                // Prefer past dates when ambiguous
    PreferFuture().              // Prefer future dates when ambiguous
    PreferCurrentPeriod().       // Prefer current period (default)
    Timezone(tz).                // Set default timezone
    ForwardDate()                // Parse only forward dates (legacy)
```

### Parsing Methods

Execute the parser:

- **`Parse(text string) ([]Result, error)`**: Parse all date/time expressions in text
- **`ParseDate(text string) (*time.Time, error)`**: Parse and return the first date only

### Package-Level Convenience Functions

The `en` package provides quick access for English parsing:

- **`en.ParseSimple(text string) ([]Result, error)`**: Parse with default settings
- **`en.ParseDateSimple(text string) (*time.Time, error)`**: Parse single date with defaults
- **`en.New()`**: Create a new parser builder with casual English
- **`en.StrictParser()`**: Create a parser builder with strict English
- **`en.GBParser()`**: Create a parser builder with British English (DMY order)

### Pre-built Configurations

The `en` package provides ready-to-use configurations:

- **`en.Casual`**: Casual English parsing (recognizes informal expressions)
- **`en.Strict`**: Strict English parsing (formal patterns only)
- **`en.GB`**: British English parsing (DMY date order)

## Examples

### Relative Dates

```go
results, _ := en.ParseSimple("yesterday")
results, _ := en.ParseSimple("tomorrow")
results, _ := en.ParseSimple("last Monday")
results, _ := en.ParseSimple("next week")
results, _ := en.ParseSimple("3 days ago")
results, _ := en.ParseSimple("in 2 weeks")
results, _ := en.ParseSimple("2 hours from now")
```

### Absolute Dates

```go
results, _ := en.ParseSimple("March 15, 2024")
results, _ := en.ParseSimple("3/15/2024")
results, _ := en.ParseSimple("2024-03-15")
results, _ := en.ParseSimple("15th of March")
results, _ := en.ParseSimple("March 2024")
```

### Time Expressions

```go
results, _ := en.ParseSimple("3pm")
results, _ := en.ParseSimple("15:30")
results, _ := en.ParseSimple("noon")
results, _ := en.ParseSimple("midnight")
results, _ := en.ParseSimple("3:30:45 PM")
results, _ := en.ParseSimple("14:30 EST")
```

### Combined Date and Time

```go
results, _ := en.ParseSimple("tomorrow at 3pm")
results, _ := en.ParseSimple("March 15 at 14:30")
results, _ := en.ParseSimple("next Friday at noon")
results, _ := en.ParseSimple("2024-03-15 15:30:00")
```

### Date Order Configuration

```go
// US format (Month/Day/Year)
parser := en.New().DateOrder(kronos.DateOrderMDY)
results, _ := parser.Parse("3/15/2024")  // March 15, 2024

// European format (Day/Month/Year)
parser = en.New().DateOrder(kronos.DateOrderDMY)
results, _ = parser.Parse("15/3/2024")   // March 15, 2024

// Or use the GB parser
results, _ = en.GBParser().Parse("15/3/2024")
```

### Prefer Past/Future

```go
// When parsing "March" in November, prefer last March
parser := en.New().
    WithReferenceDate(time.Date(2024, 11, 15, 12, 0, 0, 0, time.UTC)).
    PreferPast()
results, _ := parser.Parse("March")  // March 2024 (past)

// Prefer next March
parser = en.New().
    WithReferenceDate(time.Date(2024, 11, 15, 12, 0, 0, 0, time.UTC)).
    PreferFuture()
results, _ = parser.Parse("March")   // March 2025 (future)
```

### Multiple Dates in Text

```go
text := "The meeting is on March 15 at 2pm, with a follow-up next Friday."
results, _ := en.ParseSimple(text)

for _, r := range results {
    fmt.Printf("Position %d: '%s' = %v\n", r.Index(), r.Text(), r.Date())
}
// Output:
// Position 18: 'March 15 at 2pm' = 2024-03-15 14:00:00
// Position 52: 'next Friday' = 2024-03-22 12:00:00
```

## Advanced Usage

### Component Inspection

Determine which parts of a date were explicitly mentioned:

```go
results, _ := en.ParseSimple("March 15 at 3pm")
comp := results[0].Start()

// Check each component
components := []kronos.Component{
    kronos.ComponentYear,
    kronos.ComponentMonth,
    kronos.ComponentDay,
    kronos.ComponentHour,
    kronos.ComponentMinute,
}

for _, c := range components {
    if val := comp.Get(c); val != nil {
        certain := "implied"
        if comp.IsCertain(c) {
            certain = "certain"
        }
        fmt.Printf("%s: %d (%s)\n", c, *val, certain)
    }
}
// Output:
// month: 3 (certain)
// day: 15 (certain)
// hour: 15 (certain)
// year: 2024 (implied)
// minute: 0 (implied)
// second: 0 (implied)
```

### Custom Chrono Configuration

For advanced use cases, you can create your own `Chrono` instance with custom parsers and refiners:

```go
import (
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

// Create a custom configuration
config := &kronos.Configuration{
    Parsers: []kronos.Parser{
        // Add your custom parsers
    },
    Refiners: []kronos.Refiner{
        // Add your custom refiners
    },
}

chrono := kronos.NewChrono(config)
parser := kronos.New(chrono)
```

## Architecture

Kronos uses a pipeline architecture:

1. **Parsers**: Pattern-based parsers scan the input text for date/time expressions
2. **Components**: Extracted date/time parts are stored as components (year, month, day, etc.)
3. **Refiners**: Post-processors refine and merge results (e.g., merge "tomorrow at 3pm")
4. **Results**: Final parsed results with full date/time information

Each result includes:
- The matched text and its position
- Parsed components (with certainty tracking)
- Implied components (filled in from reference date)
- A constructed `time.Time` object

## Migrating from v1

See [MIGRATION.md](MIGRATION.md) for a detailed migration guide.

Quick summary of changes:

**Old API:**
```go
// v1 - direct Chrono usage
import "github.com/kljensen/kronos/en"

results := en.Parse("tomorrow", time.Now(), nil)
date := results[0].Date()
```

**New API:**
```go
// v2 - builder pattern
import "github.com/kljensen/kronos/en"

results, err := en.ParseSimple("tomorrow")
date := results[0].Date()
```

Key changes:
- Builder pattern for configuration
- Error handling on Parse methods
- Public `Result` and `Components` interfaces instead of internal types
- `en.ParseSimple()` convenience functions

## Testing

Run the test suite:

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## Performance

Kronos is designed for correctness and usability over raw speed, but performance is still good for most use cases. The parser pipeline processes text in a single pass, and refiners operate on already-matched results.

For high-throughput scenarios:
- Reuse `ParserBuilder` instances (they're mutable but not thread-safe)
- Create one builder per goroutine
- Consider using `ParseDate()` instead of `Parse()` if you only need the first result

## Contributing

Contributions are welcome! Please feel free to submit issues, fork the repository, and create pull requests.

## License

This project is licensed under the MIT License.

## Acknowledgments

This library is a Go port of the excellent [chrono](https://github.com/wanasit/chrono) JavaScript library by Wanasit Tanakitrungruang. The design, patterns, and test cases are heavily inspired by the original implementation.

## Related Projects

- [chrono](https://github.com/wanasit/chrono) - The original JavaScript library
- [dateparser](https://github.com/scrapinghub/dateparser) - Python date parsing library with similar goals
