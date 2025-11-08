<div align="center">
  <img src="doc/gopher.png" alt="Kronos Gopher" width="100"/>

  # Kronos

[![CI](https://img.shields.io/github/actions/workflow/status/kljensen/kronos/ci.yml?branch=main&style=for-the-badge&logo=github-actions&logoColor=white&label=CI)](https://github.com/kljensen/kronos/actions/workflows/ci.yml)
[![Go Reference](https://img.shields.io/badge/go.dev-reference-007d9c?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/kljensen/kronos)
[![Go Report Card](https://goreportcard.com/badge/github.com/kljensen/kronos?style=for-the-badge)](https://goreportcard.com/report/github.com/kljensen/kronos)
[![License: Unlicense](https://img.shields.io/badge/License-Unlicense-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/Unlicense)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kljensen/kronos?style=for-the-badge&logo=go)](https://github.com/kljensen/kronos)

</div>

## Overview

Kronos is a Go package for extracting structured datetimes from natural language.
It is heavily inspired by the excellent
[chrono](https://github.com/wanasit/chrono) JavaScript library.
Using Kronos, you can turn text like "see you Thursday at 5:30pm", "4/5/2023 around 6",
or "Thu, November 6th, 2025, 8:00 AM EST" in to Golang structures. Kronos will keep
track of ambiguity (such as am/pm in the "around 6" example above). Features include
the following.

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
- **Extensive test coverage**: test casese stolen from other date parsing libraries

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

### Time Ranges

Kronos can also parse time ranges within date expressions:

```go
// Time range on a specific date
results, _ := en.ParseSimple("Wed, Nov 5, 2025 8:30 PM - 10:45 PM")
if results[0].End() != nil {
    start := results[0].Start().Date()  // 2025-11-05 20:30:00
    end := results[0].End().Date()      // 2025-11-05 22:45:00
    fmt.Printf("Event from %s to %s\n",
        start.Format("3:04 PM"),
        end.Format("3:04 PM"))
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
    Timezone(tz)                 // Set default timezone
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

### Migration Guide

If you're upgrading from an older version of Kronos, see the [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) for detailed migration instructions. **Note**: 95%+ of users are unaffected by recent API changes and can upgrade without any code modifications.

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

## Contributing

Contributions are welcome! Please feel free to submit issues, fork the repository, and create pull requests. You must be comfy with the license (below).

## License (the Unlicense)

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org/>


## Related Projects

Kronos is inspired by excellent natural language date/time parsing libraries across many programming languages. Here are equivalent projects in major languages:

### Go
- **[when](https://github.com/olebedev/when)** - Natural language date/time parser with pluggable rules and merge strategies. Supports multiple languages (EN, RU, PT_BR, ZH, NL). Parses expressions like "next wednesday at 2:25 p.m" and "tomorrow at noon".
- **[go-dateparser](https://github.com/markusmobius/go-dateparser)** - Port of Python's dateparser library supporting 200+ language locales. Handles relative dates like "1 min ago", "2 weeks ago", "in 2 days", "tomorrow".
- **[go-naturaldate](https://github.com/tj/go-naturaldate)** - Natural date/time parsing for human-friendly relative date/time ranges. Defaults to past direction for ambiguous expressions.
- **[go-anytime](https://github.com/ijt/go-anytime)** - Parse natural and standardized dates/times and ranges without knowing the format in advance.
- **[dateparse](https://github.com/araddon/dateparse)** - Parse many date strings without knowing format in advance. While not strictly "natural language," it's very flexible for parsing various date formats automatically.

### JavaScript/TypeScript
- **[chrono](https://github.com/wanasit/chrono)** - The original JavaScript library that inspired Kronos. Supports multiple languages (en, ja, fr, nl, ru, uk, and partial support for de, es, pt, zh.hant). Available as `chrono-node` on npm.

### Python
- **[dateparser](https://github.com/scrapinghub/dateparser)** - Comprehensive Python library supporting 200+ language locales. Parses specific dates ('5:47pm 29th of December, 2015') and relative times ('10 minutes ago'). Built on top of dateutil.parser with enhanced natural language support.
- **[parsedatetime](https://github.com/bear/parsedatetime)** - Parses natural language dates with strong support for relative dates like 'tomorrow'. Supports multiple locales but requires manual specification.
- **[timefhuman](https://github.com/alvinwan/timefhuman)** - Extracts datetimes and durations from natural language text. Supports ranges, lists, and more complex expressions.

### Ruby
- **[chronic](https://github.com/mojombo/chronic)** - Pure Ruby natural language date parser. Handles a huge variety of date/time formats with case-insensitive parsing and common abbreviation/misspelling support.

### Java
- **[Natty](https://github.com/joestelmach/natty)** - Natural language date parser applying standard language recognition and translation techniques. Handles complex patterns like "2 wednesdays from now".
- **[PrettyTime::NLP](https://www.ocpsoft.org/prettytime/nlp/)** - Wraps Natty for parsing natural language date/time expressions with integration into the PrettyTime formatting library.

### C# / .NET
- **[Microsoft.Recognizers.Text](https://github.com/microsoft/Recognizers-Text)** - Microsoft's comprehensive recognizer for dates, times, numbers, and more from natural language text.
- **[Chronic .NET](https://github.com/mojombo/chronic/wiki/.NET-Ports)** - .NET port of the Ruby Chronic library for natural language date parsing.

### PHP
- **[Carbon](https://github.com/briannesbitt/Carbon)** - Popular PHP API extension for DateTime with natural language parsing support ('tomorrow', 'next wednesday', '1 year ago'). Supports 200+ languages and 500+ regional variants.

### Rust
- **[dateparser](https://crates.io/crates/dateparser)** - Parses date strings in commonly used formats, returning chrono DateTime objects. Supports unix timestamps, RFC formats, and various common date formats.
- **[temps-chrono](https://lib.rs/crates/temps-chrono)** - Parses human-readable time expressions in multiple languages (English, German). Handles "in 2 hours", "tomorrow at 2:00 pm", "next tuesday".

### Swift / iOS
- **[NSDataDetector](https://developer.apple.com/documentation/foundation/nsdatadetector)** - Apple's Foundation framework class for detecting dates, times, addresses, and other data types in natural language text. Limited natural language support (works for specific dates but not relative expressions like "last week").
- **[SoulverCore](https://github.com/soulverteam/SoulverDateFromString)** - Advanced natural language date input library for Mac and iOS apps.

### Kotlin / Android
- **[Natty](http://natty.joestelmach.com/)** - Same Java library usable in Kotlin/Android projects for natural language date parsing.
- **[Hawking](https://github.com/zoho/hawking)** - Natural Language Date Time Parser that extracts date/time from text with context.

### Perl
- **[DateTime::Format::Natural](https://metacpan.org/pod/DateTime::Format::Natural)** - Parses informal natural language date/time strings into DateTime objects. Handles "tomorrow", "next Tuesday", "1 hour ago".
- **[Date::Manip](https://metacpan.org/pod/Date::Manip)** - Powerful date manipulation and parsing module for human-formatted dates. Most comprehensive but also largest and slowest Perl date module.

### Elixir
- **[Timex](https://github.com/bitwalker/timex)** - Complete date/time library for Elixir with timezone support. Primarily handles structured formats (ISO 8601, RFC 1123, custom format strings) with limited natural language support.

### Scala
- **[nscala-time](https://github.com/nscala-time/nscala-time)** - Scala wrapper around Joda-Time with more idiomatic Scala syntax. Focused on improving expressiveness rather than advanced natural language parsing.

### C++
- **[Howard Hinnant's date library](https://github.com/HowardHinnant/date)** - Comprehensive date/time library based on C++11/14/17 `<chrono>`. Includes IANA timezone database parser. Focuses on format-based parsing rather than natural language.
