# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

#### API Minimization and Reorganization

Kronos has undergone a major API cleanup to provide a cleaner, more focused public interface. **95%+ of users are completely unaffected** by these changes. The core functionality remains unchanged.

**Impact Summary:**
- Reduced public API surface from ~200 exports to 126 exports
- Core parsing functionality unchanged
- Builder pattern API unchanged
- Standard use cases require no code changes

**What Changed:**

1. **X-Prefixed Helper Functions Removed**
   - Functions like `XToday()`, `XTomorrow()`, `XFindMostLikelyADYear()` have been removed from the public API
   - These were internal implementation details exposed temporarily
   - **Migration:** Use `experimental` package instead (see below)
   - **Affected users:** Only those directly calling X-prefixed functions (<1%)

2. **Concrete Parsing Types Deprecated**
   - `ParsingComponents`, `ParsingResult`, and `ParsingContext` concrete types are now deprecated
   - New code should use interfaces: `ParsedComponents`, `ParsedResult`
   - **Migration:** Use interfaces from main package, or import from `experimental` package
   - **Affected users:** Those using concrete types directly (<2%)
   - **Timeline:** Deprecated types will be removed in v2.0.0

3. **Internal Data Variables Moved**
   - `ApproximationWords`, `DefaultTimezoneAbbrMap`, and `EmptyDuration` moved to internal packages
   - **Migration:** Import from `experimental` package if needed
   - **Affected users:** Those directly accessing these constants (<1%)

4. **Settings Direct Usage Deprecated**
   - Direct manipulation of `Settings` struct is deprecated
   - **Migration:** Use builder pattern methods instead
   - **Example:** Use `en.New().DateOrder(kronos.DateOrderDMY).PreferFuture()` instead of creating `Settings` directly
   - **Affected users:** Those directly creating/modifying Settings (<2%)

5. **Parser/Refiner Infrastructure Moved**
   - `Parser`, `Refiner`, `Configuration`, `Chrono` types are now deprecated in main package
   - **Migration:** Import from `experimental` package for custom parser development
   - **Affected users:** Custom parser/refiner authors (<1%)

**New: Experimental Package**

The new `experimental` package provides advanced features for power users:

```go
import "github.com/kljensen/kronos/experimental"
```

**What's in experimental:**
- Helper functions: `Today()`, `Tomorrow()`, `FindMostLikelyADYear()`, etc.
- Parsing types: `ParsingComponents`, `ParsingResult`, `ParsingContext`
- Parser infrastructure: `Parser`, `Refiner`, `Configuration`, `Chrono`
- Date math: `AddDuration()`, `ReverseDuration()`, `GetLastWeekday()`, etc.
- Internal constants: `ApproximationWords`, `DefaultTimezoneAbbrMap`, `EmptyDuration`
- Configuration options: Advanced `Settings` options via `Option` functions

**Stability Warning:** The experimental package may change between minor versions. Use it only if you need advanced features not available in the main package.

**Migration Examples:**

Before (X-prefixed functions):
```go
import "github.com/kljensen/kronos"

today := kronos.XToday(refTime, false)
year := kronos.XFindMostLikelyADYear(98)
```

After (experimental package):
```go
import "github.com/kljensen/kronos/experimental"

today := experimental.Today(refTime, false)
year := experimental.FindMostLikelyADYear(98)
```

Before (concrete types):
```go
import "github.com/kljensen/kronos"

var comp *kronos.ParsingComponents  // Deprecated
```

After (interfaces - recommended):
```go
import "github.com/kljensen/kronos"

var comp kronos.ParsedComponents  // Use interface
```

After (experimental package - alternative):
```go
import "github.com/kljensen/kronos/experimental"

var comp *experimental.ParsingComponents
```

Before (custom parser):
```go
import "github.com/kljensen/kronos"

type MyParser struct {
    kronos.Parser  // Deprecated
}
```

After (experimental package):
```go
import (
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/experimental"
)

type MyParser struct {
    experimental.Parser
}

config := &experimental.Configuration{
    Parsers: []experimental.Parser{&MyParser{}},
}
chrono := experimental.NewChrono(config)
parser := kronos.New(chrono)
```

**Unaffected Code Patterns:**

These common patterns require **NO changes**:

```go
// Basic parsing - unchanged
results, err := en.ParseSimple("tomorrow at 3pm")

// Builder pattern - unchanged
parser := en.New().
    PreferPast().
    DateOrder(kronos.DateOrderDMY)

// Pre-built configurations - unchanged
parser := en.Casual  // or en.Strict, en.GB

// Interface usage - unchanged
results, err := parser.Parse(text)
for _, r := range results {
    comp := r.Start()  // Returns Components interface
    date := r.Date()   // Returns *time.Time
}
```

**For More Information:**

See [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) for comprehensive migration instructions, examples, and troubleshooting.

**Future Plans:**

- The deprecated APIs will remain available with warnings until v2.0.0
- The main package API is now stable and focused on common use cases
- The experimental package provides an escape hatch for advanced features
- Future API surface reduction will further simplify the main package in v2.0.0 (targeting ~60-70 exports)

## [0.4.0] - 2024-11-05

### Added
- Builder pattern API for fluent configuration
- `ParserBuilder` type with chainable methods
- `en.New()`, `en.StrictParser()`, `en.GBParser()` builder constructors
- Pre-built configurations: `en.Casual`, `en.Strict`, `en.GB`
- New interface types: `ParsedComponents`, `ParsedResult`

### Changed
- Reorganized internal packages for better maintainability
- Moved parser and refiner implementations to internal packages
- Simplified English package exports to builder constructors only

### Deprecated
- Direct `Settings` manipulation (use builder pattern instead)
- Concrete types: `ParsingComponents`, `ParsingResult`, `ParsingContext`

## [0.3.0] - 2024-10-15

### Added
- Support for British English date order (DMY)
- Enhanced timezone parsing
- Improved date range detection

### Changed
- Performance improvements in parser matching
- Better handling of ambiguous dates

### Fixed
- Edge cases in relative date parsing
- Timezone offset calculation bugs

## [0.2.0] - 2024-09-01

### Added
- Comprehensive date range support
- Component certainty tracking
- Multiple locale support framework

### Changed
- Refactored parser pipeline
- Improved result merging logic

## [0.1.0] - 2024-08-01

### Added
- Initial release
- Basic natural language date parsing
- English language support
- Common date/time formats
- Relative date expressions
- Time expression parsing

[Unreleased]: https://github.com/kljensen/kronos/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/kljensen/kronos/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/kljensen/kronos/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/kljensen/kronos/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/kljensen/kronos/releases/tag/v0.1.0
