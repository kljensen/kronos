# Kronos Settings System

The Kronos settings system provides comprehensive configuration for date parsing, similar to Python's dateparser library. This document explains how to use the settings system effectively.

## Overview

Settings allow you to control:
- Date interpretation (MDY vs DMY vs YMD)
- Date preferences (past, future, or current period)
- Timezone handling
- Parsing strictness
- Which parsers to use
- Parser execution order
- Performance limits

## Basic Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

func main() {
    // Create default settings
    settings := kronos.DefaultSettings()

    // Customize as needed
    settings.DateOrder = kronos.DateOrderMDY
    settings.PreferDatesFrom = kronos.PreferFuture

    // Parse with settings
    refDate := time.Now()
    results, err := en.Casual.ParseWithSettings("March 15", refDate, settings)
    if err != nil {
        panic(err)
    }

    for _, result := range results {
        fmt.Println(result.Date())
    }
}
```

## Settings Reference

### Date Interpretation

#### DateOrder
Controls how ambiguous numeric dates are interpreted:
- `DateOrderMDY`: Month-Day-Year (US format) - default
- `DateOrderDMY`: Day-Month-Year (European format)
- `DateOrderYMD`: Year-Month-Day (ISO format)

```go
settings := kronos.DefaultSettings()
settings.DateOrder = kronos.DateOrderDMY
// "3/15/2020" → March 15, 2020 (MDY)
// "15/3/2020" → March 15, 2020 (DMY)
```

#### PreferDatesFrom
Controls how dates with missing components are interpreted:
- `PreferCurrentPeriod`: Use current year/day - default
- `PreferPast`: Prefer dates in the past
- `PreferFuture`: Prefer dates in the future

```go
settings := kronos.DefaultSettings()
settings.PreferDatesFrom = kronos.PreferFuture
// If today is March 10, 2020:
// "March 15" → March 15, 2020 (future)
// "March 5" → March 5, 2021 (next occurrence)
```

#### PreferDayOfMonth
Controls day selection when ambiguous:
- `DayPreferCurrent`: Current day - default
- `DayPreferFirst`: First day of period
- `DayPreferLast`: Last day of period

### Timezone Handling

#### Timezone
Default timezone for parsing:
```go
settings := kronos.DefaultSettings()
settings.Timezone = "America/New_York"
```

#### ToTimezone
Convert parsed results to a different timezone:
```go
settings := kronos.DefaultSettings()
settings.ToTimezone = "Europe/London"
```

#### ReturnTimezoneAware
Include timezone information in results:
```go
settings := kronos.DefaultSettings()
settings.ReturnTimezoneAware = true
```

### Parsing Behavior

#### StrictParsing
Reject ambiguous or incomplete dates:
```go
settings := kronos.DefaultSettings()
settings.StrictParsing = true
// "March 15" → rejected (no year)
// "March 15, 2020" → accepted
```

#### Normalize
Enable Unicode normalization (default: true):
```go
settings := kronos.DefaultSettings()
settings.Normalize = true
```

#### SkipTokens
Words to ignore during parsing:
```go
settings := kronos.DefaultSettings()
settings.SkipTokens = []string{"at", "on", "the"}
// "on the 15th at 3pm" → "15th 3pm"
```

#### RequireParts
Require specific date/time components:
```go
settings := kronos.DefaultSettings()
settings.RequireParts = []string{"year", "month", "day"}
// "March 15" → rejected
// "March 15, 2020" → accepted
```

Valid parts: `"year"`, `"month"`, `"day"`, `"hour"`, `"minute"`, `"second"`

### Relative Dates

#### RelativeBase
Base time for relative date calculations:
```go
baseTime := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)
settings := kronos.DefaultSettings()
settings.RelativeBase = &baseTime
// "tomorrow" → March 16, 2020
```

### Parser Control

#### EnabledParsers
Specify which parsers to use:
```go
settings := kronos.DefaultSettings()
settings.EnabledParsers = []string{
    "iso8601",
    "en_casual_date",
    "en_time_expression",
}
```

Available parsers:
- `iso8601`: ISO 8601 formats
- `en_casual_date`: "today", "tomorrow", "yesterday"
- `en_time_expression`: "3:30 PM", "14:30"
- `en_time_unit_ago`: "2 days ago"
- `en_time_unit_later`: "in 3 hours"
- `en_month_name_middle_endian`: "March 15, 2020"
- `en_month_name_little_endian`: "15 March 2020"
- And many more...

List all parsers:
```go
infos := kronos.GlobalRegistry.ListParsers()
for _, info := range infos {
    fmt.Printf("%s: %s\n", info.Name, info.Description)
}
```

#### ParserOrder
Customize parser execution order:
```go
settings := kronos.DefaultSettings()
settings.ParserOrder = []string{
    "iso8601",
    "en_casual_date",
    // ... more parsers
}
```

### Performance

#### MaxParsers
Limit the number of parsers:
```go
settings := kronos.DefaultSettings()
settings.MaxParsers = 10
```

#### Timeout
Set parsing timeout:
```go
settings := kronos.DefaultSettings()
settings.Timeout = 5 * time.Second
```

## Advanced Examples

### US Date Format with Future Preference
```go
settings := kronos.DefaultSettings()
settings.DateOrder = kronos.DateOrderMDY
settings.PreferDatesFrom = kronos.PreferFuture
settings.Timezone = "America/New_York"

result, _ := en.Casual.ParseWithSettings("3/15 at 2pm", time.Now(), settings)
```

### European Format with Strict Parsing
```go
settings := kronos.DefaultSettings()
settings.DateOrder = kronos.DateOrderDMY
settings.StrictParsing = true
settings.RequireParts = []string{"year", "month", "day"}
settings.Timezone = "Europe/London"

result, _ := en.Casual.ParseWithSettings("15/3/2020", time.Now(), settings)
```

### Custom Parser Selection
```go
settings := kronos.DefaultSettings()
settings.EnabledParsers = []string{
    "relative",
    "casual_date",
    "iso8601",
}

result, _ := en.Casual.ParseWithSettings("2 days ago", time.Now(), settings)
```

### Skip Common Words
```go
settings := kronos.DefaultSettings()
settings.SkipTokens = []string{"at", "on", "the"}

result, _ := en.Casual.ParseWithSettings("on the 15th at 3pm", time.Now(), settings)
```

## Pipeline API

For more control, use the pipeline API directly:

```go
// Create a configuration
config := en.CreateConfiguration(false, false)

// Create settings
settings := kronos.DefaultSettings()
settings.StrictParsing = true

// Create pipeline
pipeline, err := kronos.NewPipelineWithSettings(config, settings)
if err != nil {
    panic(err)
}

// Execute pipeline
results, err := pipeline.Execute("March 15, 2020", time.Now())
if err != nil {
    panic(err)
}
```

## Parser Registration

Register custom parsers:

```go
kronos.Register("my_parser", kronos.ParserInfo{
    Description: "My custom parser",
    Priority:    50,
    Tags:        []string{"custom"},
}, func() kronos.Parser {
    return NewMyCustomParser()
})
```

## Backward Compatibility

The settings system is fully backward compatible. Existing code continues to work:

```go
// Old API still works
results := en.Parse("March 15, 2020", time.Now(), nil)
date := en.ParseDate("tomorrow", time.Now(), nil)
```

Default settings maintain the same behavior as before the settings system was introduced.

## Migration from ParsingOption

If you're using `ParsingOption`, you can migrate to Settings:

```go
// Old way
option := &kronos.ParsingOption{
    ForwardDate: true,
    Preference:  kronos.PreferFuture,
}
results := en.Parse("March 15", time.Now(), option)

// New way
settings := kronos.DefaultSettings()
settings.ForwardDate = true
settings.PreferDatesFrom = kronos.PreferFuture
results, _ := en.Casual.ParseWithSettings("March 15", time.Now(), settings)
```

## Best Practices

1. **Create settings once**: Create a settings object and reuse it for multiple parses
2. **Validate settings**: Always check errors from `ParseWithSettings`
3. **Use specific parsers**: If you know the format, specify `EnabledParsers` for better performance
4. **Set timeouts**: For production use, always set a timeout
5. **Test with your data**: Different settings work better for different data

## Performance Considerations

- Fewer enabled parsers = faster parsing
- StrictParsing = faster (fewer results to process)
- MaxParsers limits parser execution
- Timeout prevents runaway parsing

## Error Handling

Always handle errors:

```go
results, err := en.Casual.ParseWithSettings("March 15", time.Now(), settings)
if err != nil {
    log.Printf("Parse error: %v", err)
    return
}
```

Common errors:
- Invalid timezone
- Invalid required parts
- Timeout exceeded
- Invalid settings values

## See Also

- [Issue #89](https://github.com/kljensen/kronos/issues/89) - Original feature request
- Python's dateparser - Inspiration for this system
