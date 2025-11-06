# Kronos API Minimization Migration Guide

This guide helps you migrate your code if you're affected by the API minimization changes in Kronos. **Good news: 95%+ of users are completely unaffected** and can upgrade without any code changes.

## Overview

Kronos has undergone a major API cleanup to provide a cleaner, more focused public interface. The core functionality remains unchanged, but some advanced and internal APIs have been reorganized.

## Are You Affected?

You are likely **NOT affected** if you:

- Use `en.ParseSimple()` or `en.ParseDateSimple()` for basic parsing
- Use `en.New()`, `en.StrictParser()`, or `en.GBParser()` with the builder pattern
- Work with `Result`, `Components`, and `Component` types from the main package
- Use pre-built configurations like `en.Casual`, `en.Strict`, or `en.GB`

You **ARE affected** if you:

- Used X-prefixed helper functions (removed)
- Accessed concrete types like `ParsingComponents`, `ParsingResult`, or `ParsingContext` (deprecated)
- Used variables `ApproximationWords`, `DefaultTimezoneAbbrMap`, or `EmptyDuration` (moved to internal)
- Created custom parsers or refiners (moved to experimental)
- Directly used the `Settings` struct (deprecated in favor of builder pattern)

## Impact Assessment

- **95%+ of users**: No changes required
- **<5% of users**: Need to migrate to experimental package or builder pattern
- **Breaking changes**: Limited to advanced/internal APIs only
- **Core functionality**: Completely unchanged

## Migration Paths

### 1. X-Prefixed Helper Functions (Removed)

These internal helper functions were removed in Phase 5. They were never intended for public use and have been moved to internal packages.

**Before (No longer works):**
```go
import "github.com/kljensen/kronos"

// These functions no longer exist in the public API
comp := kronos.XToday(refTime, false)
tomorrow := kronos.XTomorrow(refTime, false)
year := kronos.XFindMostLikelyADYear(98)
```

**After (Use experimental package):**
```go
import "github.com/kljensen/kronos/experimental"

// Helper functions are now in experimental package
comp := experimental.Today(refTime, false)
tomorrow := experimental.Tomorrow(refTime, false)
year := experimental.FindMostLikelyADYear(98)
```

**Why the change?**
These were internal helper functions exposed temporarily with X-prefixes. They're implementation details not needed by most users. Advanced users can access them via the experimental package.

### 2. Concrete Parsing Types (Deprecated)

The concrete types `ParsingComponents`, `ParsingResult`, and `ParsingContext` are now deprecated. Use interfaces or experimental package instead.

**Before (Deprecated):**
```go
import "github.com/kljensen/kronos"

// Direct use of concrete types
var comp *kronos.ParsingComponents
var result *kronos.ParsingResult
var ctx *kronos.ParsingContext
```

**After Option 1 (Use interfaces - Recommended):**
```go
import "github.com/kljensen/kronos"

// Use interfaces from the main package
var comp kronos.ParsedComponents  // Interface
var result kronos.ParsedResult    // Interface

// Most common: Just use Result from Parse()
results, err := parser.Parse(text)
for _, r := range results {
    // r is a Result, which provides access to components
    comp := r.Start()  // Returns Components interface
    date := r.Date()   // Returns *time.Time
}
```

**After Option 2 (Use experimental package):**
```go
import "github.com/kljensen/kronos/experimental"

// Concrete types available in experimental
var comp *experimental.ParsingComponents
var result *experimental.ParsingResult
var ctx *experimental.ParsingContext
```

**Why the change?**
Exposing concrete parsing types couples users to internal implementation details. The interface-based approach provides a cleaner API surface and more flexibility for future improvements.

### 3. Internal Data Variables (Moved to Internal)

Variables like `ApproximationWords`, `DefaultTimezoneAbbrMap`, and `EmptyDuration` have been moved to internal packages.

**Before (No longer works):**
```go
import "github.com/kljensen/kronos"

// These are no longer exported
words := kronos.ApproximationWords
tzMap := kronos.DefaultTimezoneAbbrMap
empty := kronos.EmptyDuration
```

**After (Use experimental package):**
```go
import "github.com/kljensen/kronos/experimental"

// Access via experimental package
words := experimental.ApproximationWords
tzMap := experimental.DefaultTimezoneAbbrMap
empty := experimental.EmptyDuration
```

**Why the change?**
These are internal constants used by parsers. Most users don't need direct access. Advanced users writing custom parsers can access them via experimental.

### 4. Custom Parsers and Refiners (Use Experimental)

If you were creating custom parsers or refiners, these types are now in the experimental package.

**Before (Deprecated):**
```go
import "github.com/kljensen/kronos"

// Direct use of parser/refiner infrastructure
type MyParser struct {
    kronos.Parser
}

config := &kronos.Configuration{
    Parsers: []kronos.Parser{myParser},
    Refiners: []kronos.Refiner{myRefiner},
}
```

**After (Use experimental package):**
```go
import (
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/experimental"
)

// Parser/Refiner types from experimental
type MyParser struct {
    experimental.Parser
}

config := &experimental.Configuration{
    Parsers: []experimental.Parser{myParser},
    Refiners: []experimental.Refiner{myRefiner},
}

chrono := experimental.NewChrono(config)
parser := kronos.New(chrono)
```

**Why the change?**
Custom parser/refiner creation is an advanced use case. Moving these to experimental keeps the main API focused on common use cases.

### 5. Direct Settings Usage (Use Builder Pattern)

Direct manipulation of the `Settings` struct is deprecated. Use the builder pattern instead.

**Before (Deprecated):**
```go
import "github.com/kljensen/kronos"

settings := &kronos.Settings{
    DateOrder: kronos.DateOrderDMY,
    ForwardDate: true,
}
```

**After (Use builder pattern - Recommended):**
```go
import (
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

parser := en.New().
    DateOrder(kronos.DateOrderDMY).
    PreferFuture()

results, err := parser.Parse(text)
```

**Why the change?**
The builder pattern provides a cleaner, more type-safe API with better documentation and discoverability.

## Experimental Package Usage Guide

The `experimental` package is your "escape hatch" for advanced features. It re-exports internal types and functions that are useful for:

- Writing custom parsers and refiners
- Performing complex date calculations
- Migrating from deprecated public APIs
- Building domain-specific parsing solutions

### When to Use Experimental

Use experimental if you need to:

1. **Create custom parsers**: Access `Parser`, `Refiner`, `Configuration`, `Chrono`
2. **Use helper functions**: `Today()`, `Tomorrow()`, `FindMostLikelyADYear()`, etc.
3. **Access parsing internals**: `ParsingComponents`, `ParsingResult`, `ParsingContext`
4. **Perform date math**: `AddDuration()`, `ReverseDuration()`, `GetLastWeekday()`, etc.
5. **Use internal constants**: `ApproximationWords`, `DefaultTimezoneAbbrMap`

### Example: Custom Parser

Here's how to create a custom parser using the experimental package:

```go
package main

import (
    "regexp"
    "time"

    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/experimental"
)

// Custom parser that recognizes "payday" as the 15th of the current month
type PaydayParser struct{}

func (p *PaydayParser) Pattern() string {
    return `\bpayday\b`
}

func (p *PaydayParser) Parse(ctx *experimental.ParsingContext) []*experimental.ParsingResult {
    text := ctx.Text
    pattern := regexp.MustCompile(p.Pattern())
    matches := pattern.FindAllStringIndex(text, -1)

    var results []*experimental.ParsingResult
    for _, match := range matches {
        comp := experimental.NewParsingComponents(ctx.Reference, nil)
        comp.Assign(experimental.ComponentDay, 15)
        comp.ImplyComponent(experimental.ComponentHour, 12)

        result := experimental.NewParsingResult(ctx.Reference, match[0], match[1]-match[0], text[match[0]:match[1]], comp)
        results = append(results, result)
    }

    return results
}

func main() {
    config := &experimental.Configuration{
        Parsers: []experimental.Parser{&PaydayParser{}},
    }

    chrono := experimental.NewChrono(config)
    parser := kronos.New(chrono)

    results, err := parser.Parse("Reminder: payday is coming soon!")
    if err != nil {
        panic(err)
    }

    for _, r := range results {
        println(r.Text(), ":", r.Date().Format("2006-01-02"))
    }
}
```

### Stability Warning

**Important**: The experimental package may change between minor versions. While we'll make efforts to maintain compatibility, breaking changes are possible as the library evolves.

Only use experimental if:
- You need advanced customization
- You're willing to potentially update your code in future versions
- The main package doesn't meet your needs

## Common Migration Scenarios

### Scenario 1: Basic User (No Changes Needed)

If your code looks like this, **no changes are required**:

```go
package main

import (
    "fmt"
    "github.com/kljensen/kronos/en"
)

func main() {
    // This code works unchanged
    results, err := en.ParseSimple("tomorrow at 3pm")
    if err != nil {
        panic(err)
    }

    for _, r := range results {
        fmt.Printf("Date: %v\n", r.Date())
    }
}
```

### Scenario 2: Builder Pattern User (No Changes Needed)

If your code looks like this, **no changes are required**:

```go
package main

import (
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

func main() {
    // This code works unchanged
    parser := en.New().
        PreferPast().
        DateOrder(kronos.DateOrderDMY)

    results, err := parser.Parse("15/3/2024")
    if err != nil {
        panic(err)
    }

    // Use results...
}
```

### Scenario 3: Advanced User with X-Functions

**Before:**
```go
package main

import "github.com/kljensen/kronos"

func myFunction(refTime time.Time) {
    today := kronos.XToday(refTime, false)
    tomorrow := kronos.XTomorrow(refTime, false)
    year := kronos.XFindMostLikelyADYear(98)
}
```

**After:**
```go
package main

import "github.com/kljensen/kronos/experimental"

func myFunction(refTime time.Time) {
    today := experimental.Today(refTime, false)
    tomorrow := experimental.Tomorrow(refTime, false)
    year := experimental.FindMostLikelyADYear(98)
}
```

**Migration steps:**
1. Import `github.com/kljensen/kronos/experimental`
2. Replace `kronos.X*` with `experimental.*` (drop the X prefix)
3. Test your code

### Scenario 4: Custom Parser Author

**Before:**
```go
package main

import "github.com/kljensen/kronos"

type MyParser struct {
    kronos.Parser
}

func (p *MyParser) Parse(ctx *kronos.ParsingContext) []*kronos.ParsingResult {
    // Implementation
}
```

**After:**
```go
package main

import (
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/experimental"
)

type MyParser struct {
    experimental.Parser
}

func (p *MyParser) Parse(ctx *experimental.ParsingContext) []*experimental.ParsingResult {
    // Implementation (unchanged)
}

func main() {
    config := &experimental.Configuration{
        Parsers: []experimental.Parser{&MyParser{}},
    }

    chrono := experimental.NewChrono(config)
    parser := kronos.New(chrono)

    // Use parser as normal
}
```

**Migration steps:**
1. Import `github.com/kljensen/kronos/experimental`
2. Change type imports from `kronos.*` to `experimental.*`
3. Use `experimental.NewChrono()` to create the engine
4. Use `kronos.New(chrono)` to create the parser builder
5. Test thoroughly

## FAQ

### Q: Why did you make these changes?

**A:** To provide a cleaner, more focused API that's easier to learn and maintain. By moving advanced features to experimental and internal packages, we can evolve the core API more confidently while still providing power users with the tools they need.

### Q: Will my existing code break?

**A:** Only if you're using advanced features (X-prefixed functions, concrete parsing types, custom parsers). 95%+ of users won't need any changes.

### Q: Is the experimental package stable?

**A:** No. The experimental package may change between minor versions. It's for advanced users who need features not available in the main package and are willing to update their code if the experimental APIs change.

### Q: What if I don't want to use experimental?

**A:** Most users don't need experimental. Stick to the main package APIs (`en.ParseSimple()`, builder pattern, interfaces). These are stable and cover 95%+ of use cases.

### Q: Can I still create custom parsers?

**A:** Yes! Use the experimental package. See the "Custom Parser" example above.

### Q: Are there performance impacts?

**A:** No. The changes are purely organizational. Performance is unchanged.

### Q: Will there be more breaking changes?

**A:** The main package API is now stable. The experimental package may evolve. We follow semantic versioning, so any future breaking changes to the main package will be in a new major version.

### Q: What's the timeline for removing deprecated APIs?

**A:** Deprecated APIs will remain available (with deprecation warnings) until the next major version. You'll have plenty of time to migrate.

### Q: Where can I get help?

**A:** Open an issue on GitHub if you have trouble migrating. We're happy to help!

## Troubleshooting

### Error: "undefined: kronos.XToday"

**Problem:** You're using X-prefixed helper functions that were removed.

**Solution:** Import `github.com/kljensen/kronos/experimental` and use `experimental.Today()` instead.

### Error: "cannot use concrete type *kronos.ParsingComponents"

**Problem:** You're using deprecated concrete types.

**Solution:** Either:
1. Use interfaces (`kronos.ParsedComponents`) instead, OR
2. Import `github.com/kljensen/kronos/experimental` and use `experimental.ParsingComponents`

### Error: "undefined: kronos.ApproximationWords"

**Problem:** Internal data variables were moved.

**Solution:** Import `github.com/kljensen/kronos/experimental` and use `experimental.ApproximationWords`.

### Error: "undefined: kronos.Parser"

**Problem:** Custom parser types were moved to experimental.

**Solution:** Import `github.com/kljensen/kronos/experimental` and use `experimental.Parser`, `experimental.Configuration`, etc.

### Deprecation Warnings

If you see deprecation warnings in your code:
1. Read the warning message - it tells you what to use instead
2. For concrete types: Switch to interfaces or experimental package
3. For Settings: Switch to builder pattern
4. You can continue using deprecated APIs until the next major version

## Summary

The API minimization makes Kronos cleaner and more focused:

- **Main package**: Stable, simple API for 95%+ of users
- **Experimental package**: Advanced features for power users
- **Internal packages**: Implementation details, not for public use

Most users won't need to change anything. If you're affected, the experimental package provides everything you need to migrate smoothly.

For questions or help, please open an issue on GitHub!
