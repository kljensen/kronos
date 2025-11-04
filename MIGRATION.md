# Migration Guide: v1 to v2

This guide helps you migrate from Kronos v1 (direct Chrono API) to v2 (builder pattern API).

## Quick Summary

**v1 (Old API):**
```go
import "github.com/kljensen/kronos/en"

results := en.Parse("tomorrow", time.Now(), nil)
date := results[0].Date()
```

**v2 (New API):**
```go
import "github.com/kljensen/kronos/en"

results, err := en.ParseSimple("tomorrow")
date := results[0].Date()
```

## Key Changes

### 1. Error Handling

**v1:** Parse methods did not return errors
**v2:** Parse methods return `(results, error)`

```go
// Old
results := en.Parse(text, ref, nil)

// New
results, err := en.ParseSimple(text)
if err != nil {
    // Handle error
}
```

### 2. Builder Pattern

**v1:** Direct Chrono usage with options
**v2:** Fluent builder pattern

```go
// Old
option := &kronos.ParsingOption{
    ForwardDate: true,
    Preference: kronos.PreferPast,
}
results := en.Parse(text, ref, option)

// New
parser := en.New().
    WithReferenceDate(ref).
    PreferPast().
    ForwardDate()
results, err := parser.Parse(text)
```

### 3. Result Types

**v1:** `*ParsingResult` (internal type)
**v2:** `Result` (public interface)

```go
// Old - direct access to internal fields
result := results[0]
text := result.text
index := result.index
components := result.start  // *ParsingComponents

// New - interface methods
result := results[0]
text := result.Text()
index := result.Index()
components := result.Start()  // Components interface
```

### 4. Component Access

**v1:** `ParsedComponents` interface
**v2:** `Components` interface (same methods, cleaner naming)

```go
// Both versions work the same way
comp := result.Start()
year := comp.Get(kronos.ComponentYear)
isCertain := comp.IsCertain(kronos.ComponentYear)
date := comp.Date()
```

### 5. Convenience Functions

**v2 adds new convenience functions:**

```go
// Parse with defaults (casual English, current time as reference)
results, err := en.ParseSimple("tomorrow")

// Parse single date with defaults
date, err := en.ParseDateSimple("tomorrow")
```

## Common Migration Patterns

### Basic Parsing

```go
// Old
results := en.Parse("tomorrow", time.Now(), nil)

// New (Option 1: Simple)
results, err := en.ParseSimple("tomorrow")

// New (Option 2: With reference date)
results, err := en.New().
    WithReferenceDate(time.Now()).
    Parse("tomorrow")
```

### Parsing Single Date

```go
// Old
date := en.ParseDate("tomorrow", time.Now(), nil)

// New (Option 1: Simple)
date, err := en.ParseDateSimple("tomorrow")

// New (Option 2: With configuration)
date, err := en.New().
    WithReferenceDate(time.Now()).
    ParseDate("tomorrow")
```

### With Options

```go
// Old
option := &kronos.ParsingOption{
    ForwardDate: true,
    Preference: kronos.PreferPast,
}
results := en.Parse(text, ref, option)

// New
results, err := en.New().
    WithReferenceDate(ref).
    PreferPast().
    ForwardDate().
    Parse(text)
```

### Custom Timezone

```go
// Old
option := &kronos.ParsingOption{
    Timezones: map[string]interface{}{
        "PST": -480, // minutes offset
    },
}
results := en.Parse(text, ref, option)

// New - Use Settings for advanced timezone configuration
// (or use the legacy Chrono API if needed)
chrono := en.Casual
results := chrono.Parse(text, ref, option)
```

### Custom Configuration

```go
// Old
chrono := kronos.NewChrono(config)
results := chrono.Parse(text, ref, nil)

// New
chrono := kronos.NewChrono(config)
parser := kronos.New(chrono)
results, err := parser.Parse(text)
```

### Date Order (MDY vs DMY)

```go
// Old - used separate configurations
results := en.GB.Parse(text, ref, nil)  // DMY

// New - use DateOrder or GB parser
parser := en.GBParser()
results, err := parser.Parse(text)

// Or configure explicitly
parser := en.New().DateOrder(kronos.DateOrderDMY)
results, err := parser.Parse(text)
```

### Strict vs Casual

```go
// Old
results := en.Strict.Parse(text, ref, nil)

// New
parser := en.StrictParser()
results, err := parser.Parse(text)

// Or
parser := en.New().Strict()
results, err := parser.Parse(text)
```

## API Mapping Table

| v1 | v2 | Notes |
|----|-----|-------|
| `en.Parse(text, ref, opt)` | `en.ParseSimple(text)` | New version uses current time as ref |
| `en.ParseDate(text, ref, opt)` | `en.ParseDateSimple(text)` | New version uses current time as ref |
| `en.Casual` | `en.Casual` or `en.New()` | Can still use Chrono directly or via builder |
| `en.Strict` | `en.Strict` or `en.StrictParser()` | Can still use Chrono directly or via builder |
| `en.GB` | `en.GB` or `en.GBParser()` | Can still use Chrono directly or via builder |
| `chrono.Parse(text, ref, opt)` | `kronos.New(chrono).Parse(text)` | Builder wraps Chrono |
| `chrono.ParseDate(text, ref, opt)` | `kronos.New(chrono).ParseDate(text)` | Builder wraps Chrono |
| `option.ForwardDate` | `builder.ForwardDate()` | Builder method |
| `option.Preference` | `builder.PreferPast()` / `PreferFuture()` | Builder methods |
| `*ParsingResult` | `Result` | Interface instead of struct |
| `ParsedComponents` | `Components` | Same interface, cleaner name |

## Backward Compatibility

### Legacy Chrono API Still Available

The old Chrono API is still available and functional:

```go
// This still works in v2
import "github.com/kljensen/kronos/en"

results := en.Parse("tomorrow", time.Now(), nil)
results := en.ParseDate("tomorrow", time.Now(), nil)
results := en.Casual.Parse(text, ref, nil)
```

However, these methods:
- Don't provide error handling
- Return internal `*ParsingResult` type instead of `Result` interface
- Are marked as legacy/deprecated in documentation

### Using Both APIs

You can use both APIs in the same codebase during migration:

```go
// Legacy code
oldResults := en.Parse("tomorrow", time.Now(), nil)

// New code
newResults, err := en.ParseSimple("tomorrow")

// Both work, but prefer the new API for new code
```

## Migration Checklist

- [ ] Replace `en.Parse()` with `en.ParseSimple()` or builder pattern
- [ ] Replace `en.ParseDate()` with `en.ParseDateSimple()` or builder pattern
- [ ] Add error handling for Parse methods
- [ ] Update `*ParsingResult` usage to `Result` interface
- [ ] Replace `ParsingOption` structs with builder methods
- [ ] Update tests to check for errors
- [ ] Consider using builder pattern for complex configurations
- [ ] Update imports if needed (package structure unchanged)

## Step-by-Step Migration Example

### Before (v1)

```go
package main

import (
    "fmt"
    "time"
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

func parseDate(text string) time.Time {
    ref := time.Now()
    option := &kronos.ParsingOption{
        Preference: kronos.PreferPast,
    }

    results := en.Parse(text, ref, option)
    if len(results) > 0 {
        return results[0].Date()
    }

    return time.Time{}
}

func main() {
    date := parseDate("last Monday")
    fmt.Println(date)
}
```

### After (v2)

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/kljensen/kronos/en"
)

func parseDate(text string) (time.Time, error) {
    parser := en.New().
        WithReferenceDate(time.Now()).
        PreferPast()

    results, err := parser.Parse(text)
    if err != nil {
        return time.Time{}, err
    }

    if len(results) > 0 {
        return results[0].Date(), nil
    }

    return time.Time{}, nil
}

func main() {
    date, err := parseDate("last Monday")
    if err != nil {
        log.Fatalf("Parse error: %v", err)
    }
    fmt.Println(date)
}
```

### Key Improvements in v2

1. **Error handling** - Catches parsing failures
2. **Explicit error return** - Function signature includes error
3. **Builder pattern** - More readable configuration
4. **Type safety** - Public `Result` interface instead of internal type

## Getting Help

If you encounter issues during migration:

1. Check the [examples/](examples/) directory for working code
2. Review the [README.md](README.md) for full API documentation
3. Look at the test files for edge cases
4. Open an issue on GitHub

## Breaking Changes

### Removed/Changed

- None - v1 API is still available for backward compatibility

### Deprecated

- Direct use of `en.Parse()` and `en.ParseDate()` without error handling
- Direct access to `ParsingResult` internal fields (use interface methods)
- Direct construction of `ParsingOption` (use builder methods)

### New in v2

- Builder pattern API (`ParserBuilder`)
- Error handling on Parse methods
- Public `Result` and `Components` interfaces
- Convenience functions (`ParseSimple`, `ParseDateSimple`)
- Settings-based configuration system (advanced)

## Performance Notes

The new API has no significant performance impact:
- Builder pattern is zero-cost at runtime (same underlying implementation)
- Error handling adds minimal overhead
- Interface wrapping is optimized by the compiler

You can safely migrate without performance concerns.
