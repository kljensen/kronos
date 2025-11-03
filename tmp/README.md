# Kronos - Natural Language Date/Time Parser for Go

Kronos is a comprehensive natural language date/time parser for Go, ported from the JavaScript library [Chrono](https://github.com/wanasit/chrono). It can extract dates and times from casual, human-readable text in multiple formats.

## Features

- **Natural Language Parsing**: Parse casual expressions like "tomorrow", "next week", "in 3 days"
- **Multiple Formats**: Supports ISO 8601, slash formats, month names, and more
- **Date Ranges**: Extract date ranges like "Monday to Friday"
- **Time Expressions**: Parse times like "3pm", "half past 2", "quarter to 4"
- **Flexible Modes**: Casual mode for everyday language, Strict mode for formal dates
- **Locale Support**: Full English language support (UK and US formats)
- **Relative Dates**: Handles relative expressions based on reference dates

## Installation

```bash
go get github.com/markusmobius/go-chrono
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/markusmobius/go-chrono"
)

func main() {
    now := time.Now()

    // Parse a single date
    date := kronos.ParseDate("tomorrow at 3pm", now, kronos.ParsingOption{})
    fmt.Println(date) // Tomorrow's date at 15:00

    // Parse multiple dates from text
    text := "Meet me Monday at 2pm or Wednesday at 10am"
    results := kronos.Parse(text, now, kronos.ParsingOption{})
    for _, result := range results {
        fmt.Printf("%s -> %s\n", result.Text, result.Start.Date())
    }
}
```

## Usage Examples

### Basic Date Parsing

```go
now := time.Now()

// Casual expressions
kronos.ParseDate("today", now, kronos.ParsingOption{})
kronos.ParseDate("tomorrow", now, kronos.ParsingOption{})
kronos.ParseDate("yesterday", now, kronos.ParsingOption{})
kronos.ParseDate("next week", now, kronos.ParsingOption{})
kronos.ParseDate("last month", now, kronos.ParsingOption{})

// Relative dates
kronos.ParseDate("in 3 days", now, kronos.ParsingOption{})
kronos.ParseDate("2 weeks ago", now, kronos.ParsingOption{})
kronos.ParseDate("5 hours from now", now, kronos.ParsingOption{})

// Formal dates
kronos.ParseDate("2024-12-25", now, kronos.ParsingOption{})
kronos.ParseDate("Dec 25, 2024", now, kronos.ParsingOption{})
kronos.ParseDate("25/12/2024", now, kronos.ParsingOption{})
kronos.ParseDate("December 25th, 2024", now, kronos.ParsingOption{})
```

### Date Ranges

```go
text := "Monday to Friday"
results := kronos.Parse(text, now, kronos.ParsingOption{})
if len(results) > 0 && results[0].End != nil {
    fmt.Printf("Start: %s\n", results[0].Start.Date())
    fmt.Printf("End: %s\n", results[0].End.Date())
}

// Other range formats
kronos.Parse("Jan 1 - Jan 5", now, kronos.ParsingOption{})
kronos.Parse("from Monday to Wednesday", now, kronos.ParsingOption{})
```

### Time Expressions

```go
// 12-hour format
kronos.ParseDate("3pm", now, kronos.ParsingOption{})
kronos.ParseDate("3:30pm", now, kronos.ParsingOption{})

// 24-hour format
kronos.ParseDate("15:30", now, kronos.ParsingOption{})

// Casual time expressions
kronos.ParseDate("half past 3", now, kronos.ParsingOption{})
kronos.ParseDate("quarter to 4", now, kronos.ParsingOption{})
kronos.ParseDate("noon", now, kronos.ParsingOption{})
kronos.ParseDate("midnight", now, kronos.ParsingOption{})
```

### Combined Date and Time

```go
kronos.ParseDate("tomorrow at 3pm", now, kronos.ParsingOption{})
kronos.ParseDate("Dec 25 at 9:00am", now, kronos.ParsingOption{})
kronos.ParseDate("next Friday at 2pm", now, kronos.ParsingOption{})
```

### Weekday Parsing

```go
kronos.ParseDate("Monday", now, kronos.ParsingOption{})
kronos.ParseDate("next Tuesday", now, kronos.ParsingOption{})
kronos.ParseDate("last Wednesday", now, kronos.ParsingOption{})
kronos.ParseDate("this Friday", now, kronos.ParsingOption{})
```

### Casual vs Strict Mode

```go
// Casual mode (default) - accepts informal expressions
casualDate := kronos.Casual.ParseDate("tmr at 3pm", now, kronos.ParsingOption{})

// Strict mode - only formal date/time patterns
strictDate := kronos.Strict.ParseDate("2024-12-25 15:00", now, kronos.ParsingOption{})
```

### UK Date Format (Little-Endian)

```go
// US format: month/day/year (default)
usDate := kronos.Parse("12/25/2024", now, kronos.ParsingOption{})

// UK format: day/month/year
ukDate := kronos.GB.Parse("25/12/2024", now, kronos.ParsingOption{})
```

### Extracting Multiple Dates

```go
text := "The conference runs from Dec 1 to Dec 5, with a break on Dec 3"
results := kronos.Parse(text, now, kronos.ParsingOption{})

for _, result := range results {
    fmt.Printf("Found: '%s' at position %d\n", result.Text, result.Index)
    fmt.Printf("  Start: %s\n", result.Start.Date())
    if result.End != nil {
        fmt.Printf("  End: %s\n", result.End.Date())
    }
}
```

### Working with Results

```go
results := kronos.Parse("tomorrow at 3pm", now, kronos.ParsingOption{})
if len(results) > 0 {
    result := results[0]

    // Original matched text
    fmt.Println(result.Text) // "tomorrow at 3pm"

    // Position in the input string
    fmt.Println(result.Index) // 0

    // Start date/time
    fmt.Println(result.Start.Date())

    // Check individual components
    if result.Start.IsCertain(kronos.ComponentHour) {
        fmt.Println("Hour is certain:", result.Start.Get(kronos.ComponentHour))
    }

    // Check tags
    if result.Start.HasTag("ENTimeUnitCasualRelativeFormatParser") {
        fmt.Println("Parsed as casual relative time")
    }
}
```

## API Reference

### Top-Level Functions

#### `Parse(text string, ref time.Time, option ParsingOption) []*ParsingResult`
Parses text and returns all found date/time matches using casual English configuration.

#### `ParseDate(text string, ref time.Time, option ParsingOption) time.Time`
Parses text and returns the first found date using casual English configuration. Returns zero time if no date found.

#### `New(config *Configuration) *Chrono`
Creates a custom Chrono instance with specified parsers and refiners.

### Pre-configured Instances

- **`kronos.Casual`**: Casual English parsing (accepts "tomorrow", "next week", etc.)
- **`kronos.Strict`**: Strict English parsing (formal dates only)
- **`en.Casual`**: Same as `kronos.Casual`
- **`en.Strict`**: Same as `kronos.Strict`
- **`en.GB`**: UK English with little-endian date format (day/month/year)

### Types

#### `ParsingResult`
```go
type ParsingResult struct {
    Reference ReferenceWithTimezone
    Index     int
    Text      string
    Start     *ParsingComponents
    End       *ParsingComponents // nil if not a range
}
```

#### `ParsingComponents`
```go
// Methods:
Date() time.Time                         // Get the parsed date/time
Get(component Component) int             // Get a component value
IsCertain(component Component) bool      // Check if component is certain
Assign(component Component, value int)   // Assign a component value
Imply(component Component, value int)    // Imply a component value
```

#### `Component`
Constants for date/time components:
- `ComponentYear`
- `ComponentMonth`
- `ComponentDay`
- `ComponentHour`
- `ComponentMinute`
- `ComponentSecond`
- `ComponentMillisecond`
- `ComponentWeekday`
- `ComponentTimezoneOffset`
- `ComponentMeridiem`

#### `ParsingOption`
```go
type ParsingOption struct {
    ForwardDate bool    // Prefer future dates for ambiguous cases
    // Add more options as needed
}
```

## Supported Formats

### Date Formats
- ISO 8601: `2024-12-25`, `2024-12-25T15:30:00`
- Month names: `Dec 25`, `December 25th, 2024`, `25 December 2024`
- Slash format: `12/25/2024`, `25/12/2024` (US/UK)
- Casual: `today`, `tomorrow`, `yesterday`
- Relative: `next week`, `last month`, `in 3 days`, `2 weeks ago`
- Weekdays: `Monday`, `next Tuesday`, `last Friday`

### Time Formats
- 12-hour: `3pm`, `3:30pm`, `3:30:45pm`
- 24-hour: `15:30`, `15:30:45`
- Casual: `noon`, `midnight`, `morning`, `afternoon`, `evening`
- Relative: `in 2 hours`, `30 minutes ago`
- Natural: `half past 3`, `quarter to 4`

### Range Formats
- `Monday to Friday`
- `Jan 1 - Jan 5`
- `Dec 1 through Dec 5`
- `from Monday until Wednesday`

## How It Works

Kronos uses a pipeline architecture:

1. **Parsing**: Multiple parsers scan the text for date/time patterns
2. **Extraction**: Each parser extracts date/time components from matches
3. **Refinement**: Refiners post-process results (merge ranges, resolve ambiguities, etc.)
4. **Results**: Final parsed results with start/end dates

This architecture allows for:
- Modular date format support
- Easy addition of new formats
- Flexible post-processing
- High accuracy across different input styles

## Extending Kronos

You can create custom parsers and refiners:

```go
// Custom parser
type MyParser struct{}

func (p *MyParser) Pattern(context *kronos.ParsingContext) *regexp.Regexp {
    return regexp.MustCompile(`my pattern`)
}

func (p *MyParser) Extract(context *kronos.ParsingContext, match []string) interface{} {
    // Extract date components
    return map[kronos.Component]int{
        kronos.ComponentYear: 2024,
        kronos.ComponentMonth: 12,
        kronos.ComponentDay: 25,
    }
}

// Custom configuration
config := &kronos.Configuration{
    Parsers: []kronos.Parser{&MyParser{}},
    Refiners: []kronos.Refiner{},
}
chrono := kronos.New(config)
```

## Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Performance

Kronos is designed for:
- Fast pattern matching using compiled regex
- Minimal allocations
- Efficient result processing

Typical performance:
- Simple date parsing: ~1-5 µs
- Complex text with multiple dates: ~10-50 µs

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

MIT License - See LICENSE file for details

## Credits

- Original JavaScript library: [Chrono](https://github.com/wanasit/chrono) by Wanasit Tanakitrungruang
- Go port: Markus Mobius

## Related Projects

- [go-dateparser](https://github.com/araddon/dateparse): Alternative Go date parser
- [when](https://github.com/olebedev/when): Another natural language date parser

## Changelog

### v1.0.0 (Current)
- Initial release
- Full English language support
- Casual and strict parsing modes
- US and UK date format support
- Comprehensive date/time format coverage
- Date range parsing
- Relative date expressions

## Support

For issues, questions, or contributions:
- GitHub Issues: [Report bugs or request features]
- Documentation: [See examples and API reference above]

---

**Happy parsing!** 🎉
