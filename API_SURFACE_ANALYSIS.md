# Kronos Public API Surface Analysis

**Date**: 2025-11-05
**Analyst**: Claude Code
**Status**: Recommendation for API Reduction

## Executive Summary

The kronos library currently exports **200+ public identifiers** across multiple packages. Analysis shows that **most users need fewer than 20 exported symbols** to accomplish their primary use cases:
1. Parse natural language dates into `time.Time`
2. Configure date format preferences (US vs EU)
3. Configure temporal preference (past vs future dates)
4. Inspect which date components were explicitly mentioned

This document provides a comprehensive analysis of the current API surface, identifies what should remain public, and proposes a refactoring strategy.

---

## Current Public API Inventory

### Main Package (`kronos`)

#### Core Types & Entry Points (13 items)
- ✅ **KEEP** `ParserBuilder` struct - fluent API builder
- ✅ **KEEP** `New(chrono *Chrono) *ParserBuilder` - creates builder
- ✅ **KEEP** `Parse(text, chrono) ([]Result, error)` - convenience function
- ✅ **KEEP** `ParseDate(text, chrono) (*time.Time, error)` - convenience function
- ⚠️ **MOVE** `Chrono` struct - should be in experimental package
- ⚠️ **MOVE** `NewChrono(*Configuration)` - should be in experimental
- ❌ **HIDE** `ParsingResult` struct - internal implementation
- ❌ **HIDE** `ParsingComponents` struct - internal implementation
- ❌ **HIDE** `ParsingContext` struct - internal implementation
- ❌ **HIDE** `ReferenceWithTimezone` struct - internal helper
- ❌ **HIDE** `ParsingResultWithBoundary` struct - internal parser protocol
- ❌ **HIDE** `NewParsingComponents()` - internal factory
- ❌ **HIDE** `NewParsingResult()` - internal factory

#### Public Interfaces (6 items)
- ✅ **KEEP** `Result` interface - clean API for parsed results
- ✅ **KEEP** `Components` interface - clean API for date components
- ⚠️ **MOVE** `Parser` interface - extensibility, move to experimental
- ⚠️ **MOVE** `Refiner` interface - extensibility, move to experimental
- ⚠️ **MOVE** `Configuration` struct - extensibility, move to experimental
- ❌ **REMOVE** `ParsedResult` interface - duplicate of `Result`
- ❌ **REMOVE** `ParsedComponents` interface - duplicate of `Components`

#### Settings & Options (9 items)
- ✅ **KEEP** `DateOrder` type + 3 constants (MDY, DMY, YMD)
- ✅ **KEEP** `DatePreference` type + 3 constants (Past, Future, CurrentPeriod)
- ⚠️ **SIMPLIFY** `Settings` struct (20+ fields) - create minimal `Options` instead
- ❌ **HIDE** `ParsingOption` struct - legacy API
- ❌ **HIDE** `ParsingReference` struct - legacy API
- ❌ **HIDE** `DayPreference` type - advanced/internal
- ❌ **HIDE** `TimezoneAbbrMap` type - internal
- ❌ **HIDE** `AmbiguousTimezoneMap` struct - internal
- ❌ **HIDE** `DebugHandler` type - internal

#### Component Type System (48+ constants)
- ✅ **KEEP** `Component` type + 12 constants (ComponentYear, ComponentMonth, ComponentDay, etc.)
  - Users need these for `comp.Get(ComponentHour)` and `comp.IsCertain(ComponentDay)`
- ❌ **HIDE** `Timeunit` type + 12 constants (TimeunitYear, etc.)
  - Only used in internal `Duration` type
- ❌ **HIDE** `Weekday` type + 7 constants
  - Values returned as `int`, users can compare with `time.Weekday`
- ❌ **HIDE** `Month` type + 12 constants
  - Values returned as `int`, users can compare with `time.Month`
- ❌ **HIDE** `Meridiem` type + 2 constants
  - Internal AM/PM representation
- ⚠️ **CONSIDER** `Period` type + 5 constants
  - Used for granularity tracking, possibly useful for debugging

#### Utility Types & Constants (30+ items)
- ❌ **HIDE** `Duration` type (map[Timeunit]float64)
- ❌ **HIDE** `EmptyDuration` var
- ❌ **HIDE** `DefaultTimezoneAbbrMap` var
- ❌ **HIDE** `StripApproximationWords(string) (string, bool)` - internal utility
- ❌ **HIDE** All numeric constants: `HoursPerDay`, `MinutesPerHour`, `SecondsPerMinute`, etc. (20+ constants)
  - These are implementation details

#### Registry System (15+ items)
- ⚠️ **MOVE** `ParserRegistry` struct + 10 methods - move to experimental
- ⚠️ **MOVE** `GlobalRegistry` var - global state, move to experimental
- ⚠️ **MOVE** `ParserFactory` type - move to experimental
- ⚠️ **MOVE** `ParserInfo` struct - move to experimental
- ⚠️ **MOVE** `Register(name, info, factory)` - move to experimental

### English Package (`kronos/en`)

#### Pre-configured Instances (3 items)
- ⚠️ **DEPRECATE** `Casual` var - replace with `en.New()`
- ⚠️ **DEPRECATE** `Strict` var - replace with `en.StrictParser()`
- ⚠️ **DEPRECATE** `GB` var - replace with `en.GBParser()`

#### Builder Constructors (3 items)
- ✅ **KEEP** `New() *ParserBuilder` - recommended API
- ✅ **KEEP** `StrictParser() *ParserBuilder` - recommended API
- ✅ **KEEP** `GBParser() *ParserBuilder` - recommended API

#### Convenience Functions (4 items)
- ✅ **KEEP** `ParseSimple(text) ([]Result, error)` - quick parsing
- ✅ **KEEP** `ParseDateSimple(text) (*time.Time, error)` - quick date parsing
- ⚠️ **DEPRECATE** `Parse(text, ref, option)` - legacy API
- ⚠️ **DEPRECATE** `ParseDate(text, ref, option)` - legacy API

#### Data Dictionaries (8+ items)
- ❌ **HIDE** `FullMonthNameDictionary` var
- ❌ **HIDE** `MonthDictionary` var
- ❌ **HIDE** `WeekdayDictionary` var
- ❌ **HIDE** `IntegerWordDictionary` var
- ❌ **HIDE** `NumberWordDictionary` var
- ❌ **HIDE** `OrdinalWordDictionary` var
- ❌ **HIDE** `TimeUnitDictionary` var
- ❌ **HIDE** `TimeUnitRelativeDictionary` var
- ❌ **HIDE** `YearPattern` const

#### Parser Types (25+ items)
- ❌ **HIDE** All parser structs: `ENTimeUnitAgoFormatParser`, `ENWeekdayParser`, `ENMonthNameParser`, etc.
  - These should be created by factories, not directly instantiated

#### Refiner Types (6+ items)
- ❌ **HIDE** `ENMergeDateRangeRefiner`
- ❌ **HIDE** `ENMergeDateTimeRefiner`
- ❌ **HIDE** `ENMergeRelativeAfterDateRefiner`
- ❌ **HIDE** `ENMergeRelativeFollowByDateRefiner`
- ❌ **HIDE** `ENExtractYearSuffixRefiner`
- ❌ **HIDE** `ENUnlikelyFormatFilter`

### Common Package (`kronos/common`)

#### Base Parsers (15+ items)
- ❌ **MOVE TO INTERNAL** All base parser types:
  - `AbstractParserWithWordBoundary`
  - `ISOFormatParser`
  - `SlashDateFormatParser`
  - `AbstractTimeExpressionParser`
  - etc.

#### Refiners Package (`kronos/common/refiners`)
- ❌ **MOVE TO INTERNAL** All refiner infrastructure:
  - `Filter` interface
  - `MergingRefiner` interface
  - `BaseFilter`, `BaseMergingRefiner`
  - `DatePreferenceRefiner`, `ForwardDateRefiner`
  - `OverlapRemovalRefiner`
  - `AbstractMergeDateRangeRefiner`
  - `AbstractMergeDateTimeRefiner`
  - `MergeWeekdayComponentRefiner`
  - `ExtractTimezoneAbbrRefiner`
  - `ExtractTimezoneOffsetRefiner`
  - `UnlikelyFormatFilter`

---

## Recommended Minimal Public API

### Package `kronos` - Core API (15-20 exports)

```go
// Entry Points
func New(chrono *Chrono) *ParserBuilder
func Parse(text string, chrono *Chrono) ([]Result, error)
func ParseDate(text string, chrono *Chrono) (*time.Time, error)

// Builder API
type ParserBuilder struct { /* opaque */ }
func (p *ParserBuilder) WithReferenceDate(time.Time) *ParserBuilder
func (p *ParserBuilder) Strict() *ParserBuilder
func (p *ParserBuilder) Casual() *ParserBuilder
func (p *ParserBuilder) DateOrder(DateOrder) *ParserBuilder
func (p *ParserBuilder) PreferPast() *ParserBuilder
func (p *ParserBuilder) PreferFuture() *ParserBuilder
func (p *ParserBuilder) PreferCurrentPeriod() *ParserBuilder
func (p *ParserBuilder) Timezone(string) *ParserBuilder
func (p *ParserBuilder) Parse(string) ([]Result, error)
func (p *ParserBuilder) ParseDate(string) (*time.Time, error)

// Result Inspection
type Result interface {
    Text() string
    Index() int
    Date() time.Time
    Start() Components
    End() Components  // nil for single dates
}

type Components interface {
    Get(Component) *int
    IsCertain(Component) bool
    Date() time.Time
}

// Configuration Types
type Component string
const (
    ComponentYear Component = "year"
    ComponentMonth Component = "month"
    ComponentDay Component = "day"
    ComponentWeekday Component = "weekday"
    ComponentHour Component = "hour"
    ComponentMinute Component = "minute"
    ComponentSecond Component = "second"
    ComponentMillisecond Component = "millisecond"
    ComponentMicrosecond Component = "microsecond"
    ComponentNanosecond Component = "nanosecond"
    ComponentMeridiem Component = "meridiem"
    ComponentTimezoneOffset Component = "timezoneOffset"
)

type DateOrder int
const (
    DateOrderMDY DateOrder = iota  // US: 12/31/2020
    DateOrderDMY                    // EU: 31/12/2020
    DateOrderYMD                    // ISO: 2020/12/31
)

type DatePreference int
const (
    PreferCurrentPeriod DatePreference = iota
    PreferPast
    PreferFuture
)
```

### Package `kronos/en` - English Language Support (6 exports)

```go
// Builder Constructors
func New() *kronos.ParserBuilder              // Casual US English
func StrictParser() *kronos.ParserBuilder     // Strict US English
func GBParser() *kronos.ParserBuilder         // British English (DMY)

// Convenience Functions
func ParseSimple(text string) ([]kronos.Result, error)
func ParseDateSimple(text string) (*time.Time, error)
```

### Package `kronos/experimental` - Advanced API

For users who need extensibility:

```go
// Advanced Configuration
type Chrono struct { /* ... */ }
func NewChrono(*Configuration) *Chrono

type Configuration struct {
    Parsers  []Parser
    Refiners []Refiner
}

type Parser interface {
    Pattern(*ParsingContext) *regexp.Regexp
    Extract(*ParsingContext, []string) interface{}
}

type Refiner interface {
    Refine(*ParsingContext, []*ParsingResult) []*ParsingResult
}

// Parser Registry
type ParserRegistry struct { /* ... */ }
var GlobalRegistry *ParserRegistry
func Register(name string, info ParserInfo, factory ParserFactory)
```

---

## Specific Recommendations by Topic

### 1. Chrono Type - Decision: Move to Experimental

**Rationale:**
- The builder pattern (`ParserBuilder`) is the recommended API
- `Chrono` is essentially a configuration container
- Current exports (`en.Casual`, `en.Strict`, `en.GB`) leak internal types
- Advanced users who need custom configurations should opt into complexity

**Action:**
- Move `Chrono`, `NewChrono`, `Configuration` to `kronos/experimental`
- Replace `en.Casual` var with `en.New()` function
- Replace `en.Strict` var with `en.StrictParser()` function
- Replace `en.GB` var with `en.GBParser()` function

**Migration Path:**
```go
// Old way
chrono := en.Casual
results := chrono.Parse(text, refDate, nil)

// New way
parser := en.New()
results, err := parser.Parse(text)
```

### 2. Component Type System - Decision: Keep Only Component

**Rationale:**
- `Component` is the ONLY enum users pass to methods (`Get()`, `IsCertain()`)
- `Timeunit` is only used internally in `Duration` calculations
- `Weekday` and `Month` are only used for VALUES returned from `Get()`, not parameters
- Users can compare values with stdlib types: `time.Weekday`, `time.Month`

**Action:**
- Keep `Component` type + 12 constants (ComponentYear, ComponentMonth, etc.)
- Move `Timeunit`, `Weekday`, `Month`, `Meridiem` to internal packages
- Move `Period` to experimental (useful for debugging, but niche)

**Before:**
```go
type Component string       // KEEP
type Timeunit string        // HIDE
type Weekday int           // HIDE
type Month int             // HIDE
type Meridiem int          // HIDE
type Period int            // MOVE to experimental
```

**After:**
```go
// kronos package
type Component string
const (
    ComponentYear Component = "year"
    // ... 11 more
)

// Users get values as int and can compare:
weekday := comp.Get(ComponentWeekday)
if weekday != nil && time.Weekday(*weekday) == time.Monday {
    // ...
}
```

### 3. Settings Complexity - Decision: Create Minimal Options

**Rationale:**
- Current `Settings` struct has 20+ fields, many experimental or rarely used
- Builder pattern only uses 5-6 settings in practice
- Experimental features like `ToTimezone` are marked TODO in code
- Advanced settings like `SkipTokens`, `EnabledParsers` are power-user features

**Action:**
- Keep `Settings` as internal implementation detail
- Create minimal builder methods for common options
- Move experimental/advanced settings to `experimental` package

**Current Settings (20+ fields):**
```go
type Settings struct {
    DateOrder        DateOrder      // CORE
    PreferDatesFrom  DatePreference // CORE
    Timezone         string         // CORE
    StrictParsing    bool           // CORE

    // HIDE/MOVE:
    PreferDayOfMonth DayPreference
    ToTimezone       string         // Experimental, not implemented
    ReturnTimezoneAware bool
    Normalize        bool
    SkipTokens       []string       // Advanced
    RequireParts     []string       // Advanced
    EnabledParsers   []string       // Advanced
    ParserOrder      []string       // Advanced
    MaxParsers       int            // Advanced
    Timeout          time.Duration  // Advanced
    ForwardDate      bool           // Legacy
    // ...
}
```

**Proposed Builder API:**
```go
type ParserBuilder struct { /* internal */ }

// Core configuration methods (keep)
func (p *ParserBuilder) WithReferenceDate(time.Time) *ParserBuilder
func (p *ParserBuilder) DateOrder(DateOrder) *ParserBuilder
func (p *ParserBuilder) PreferPast() *ParserBuilder
func (p *ParserBuilder) PreferFuture() *ParserBuilder
func (p *ParserBuilder) PreferCurrentPeriod() *ParserBuilder
func (p *ParserBuilder) Strict() *ParserBuilder
func (p *ParserBuilder) Casual() *ParserBuilder
func (p *ParserBuilder) Timezone(string) *ParserBuilder

// Advanced options move to experimental:
// experimental.WithSkipTokens()
// experimental.WithEnabledParsers()
// experimental.WithParserOrder()
```

### 4. Package Structure - Decision: Hide Common

**Rationale:**
- `common/` package exports every shared parser/refiner implementation
- External users should never instantiate parsers directly
- Only language package builders (like `en.New()`) should use these
- Extensibility story is unclear - better to make it opt-in

**Action:**
- Move `common/` → `internal/common/`
- Move `common/refiners/` → `internal/refiners/`
- Keep `en/` public for API surface
- Hide parser/refiner implementation files in `en/`
- Hide data dictionaries in `en/`

**Current Structure:**
```
kronos/
  ├── *.go                    (main package - mixed public/private)
  ├── en/
  │   ├── *.go                (ALL exported - parsers, refiners, data)
  ├── common/
  │   ├── *.go                (ALL exported - base parsers)
  │   └── refiners/
  │       └── *.go            (ALL exported - base refiners)
  └── internal/
      └── core/               (truly internal)
```

**Proposed Structure:**
```
kronos/
  ├── builder.go              (ParserBuilder - public)
  ├── result.go               (Result interface - public)
  ├── components.go           (Components interface - public)
  ├── types.go                (Component, DateOrder, DatePreference - public)
  ├── internal/
  │   ├── parsing/
  │   │   ├── chrono.go       (Chrono - internal)
  │   │   ├── context.go      (ParsingContext - internal)
  │   │   ├── result.go       (ParsingResult - internal)
  │   │   └── components.go   (ParsingComponents - internal)
  │   ├── common/
  │   │   ├── parsers/        (base parser implementations)
  │   │   └── refiners/       (base refiner implementations)
  │   └── en/
  │       ├── parsers/        (English parser implementations)
  │       ├── refiners/       (English refiner implementations)
  │       ├── config.go       (configuration builders)
  │       └── data.go         (dictionaries, constants)
  ├── en/
  │   └── en.go               (ONLY public API: New(), ParseSimple())
  └── experimental/
      ├── chrono.go           (Chrono, Configuration)
      ├── interfaces.go       (Parser, Refiner interfaces)
      └── registry.go         (ParserRegistry)
```

### 5. Implementation Leakage - Decision: Hide Concrete Types

**Rationale:**
- Public interfaces exist (`Result`, `Components`) but concrete types also exported
- This creates two APIs: legacy (concrete) and new (interfaces)
- Concrete types (`ParsingResult`, `ParsingComponents`) are implementation details
- Legacy code depends on concrete types - need migration path

**Action:**
- Move concrete types to `internal/parsing/`
- Keep interfaces in main package
- Provide adapters for legacy code
- Deprecate direct use of concrete types

**Current Problem:**
```go
// Two ways to get the same thing - confusing!
type ParsedResult interface { ... }     // Old interface
type Result interface { ... }           // New interface (same methods)

type ParsingResult struct { ... }       // Concrete type (exported!)

// Users can do:
var r1 Result = ...           // New way (clean)
var r2 *ParsingResult = ...   // Old way (leaks implementation)
```

**Proposed Solution:**
```go
// kronos/result.go - PUBLIC
type Result interface {
    Text() string
    Index() int
    Date() time.Time
    Start() Components
    End() Components
}

type Components interface {
    Get(Component) *int
    IsCertain(Component) bool
    Date() time.Time
}

// kronos/internal/parsing/result.go - PRIVATE
type parsingResult struct { ... }
type parsingComponents struct { ... }

// kronos/adapters.go - PRIVATE (converts internal to public interface)
func newResultAdapter(r *parsingResult) Result { ... }
```

---

## Migration Strategy

### Phase 1: Create Experimental Package (Non-Breaking)

1. Create `kronos/experimental` package
2. Copy `Chrono`, `Configuration`, `Parser`, `Refiner`, `Registry` to experimental
3. Update main package to use experimental types internally
4. Add deprecation notices to old locations
5. Update documentation to recommend new locations

**Example:**
```go
// kronos/chrono.go - DEPRECATED
// Deprecated: Use experimental.Chrono instead. This will be removed in v2.0.
type Chrono = experimental.Chrono

// Deprecated: Use experimental.NewChrono instead. This will be removed in v2.0.
func NewChrono(config *Configuration) *Chrono {
    return experimental.NewChrono((*experimental.Configuration)(config))
}
```

### Phase 2: Move Common to Internal (Non-Breaking)

1. Create `internal/common/` structure
2. Move parser and refiner implementations
3. Update `en/` package imports
4. Verify all tests pass
5. Add deprecation notices on old exports

### Phase 3: Hide Concrete Types (Breaking)

1. Move `ParsingResult`, `ParsingComponents`, etc. to `internal/parsing/`
2. Keep adapters that return interface types
3. Update all internal code to use new locations
4. Provide conversion functions for legacy code
5. Update examples and documentation

**Migration helpers:**
```go
// kronos/legacy/convert.go
func AsResult(pr *ParsingResult) kronos.Result {
    return newResultAdapter(pr)
}
```

### Phase 4: Clean Up English Package (Breaking)

1. Hide parser and refiner types
2. Hide data dictionaries
3. Replace exported `Chrono` vars with functions
4. Update documentation

**Before:**
```go
package en

var Casual = kronos.NewChrono(config)  // Exposed global

type ENWeekdayParser struct { ... }   // Exposed implementation
var WeekdayDictionary = map[...]{...} // Exposed data
```

**After:**
```go
package en

// Only exported items:
func New() *kronos.ParserBuilder { ... }
func StrictParser() *kronos.ParserBuilder { ... }
func GBParser() *kronos.ParserBuilder { ... }
func ParseSimple(text string) ([]kronos.Result, error) { ... }
func ParseDateSimple(text string) (*time.Time, error) { ... }

// Everything else is internal
```

### Phase 5: Simplify Settings (Breaking)

1. Keep `Settings` as internal
2. Remove experimental/incomplete features
3. Document which settings are available via builder
4. Move advanced settings to experimental if needed

---

## Impact Analysis

### Who Will Be Affected?

#### 1. **Basic Users (90% of users) - No Impact**
These users only use:
- `en.ParseSimple(text)`
- `en.New().PreferPast().Parse(text)`
- `result.Date()`, `result.Start().Get(ComponentHour)`

**Migration:** None required. Recommended API stays the same.

#### 2. **Advanced Users Creating Custom Parsers (5% of users) - Moderate Impact**
These users:
- Import `common` package base parsers
- Implement `Parser` or `Refiner` interfaces
- Use `Registry` to register custom parsers

**Migration:**
- Import `kronos/experimental` instead of main package
- Import `kronos/internal/common` (or provide public extensibility API)
- Update import paths

#### 3. **Legacy Code Using Old API (5% of users) - High Impact**
These users:
- Use `en.Casual.Parse(text, refDate, option)`
- Use concrete types `*ParsingResult`, `*ParsingComponents`
- Access fields directly on result structs

**Migration:**
- Replace `en.Casual.Parse()` with `en.New().Parse()`
- Replace concrete types with interfaces
- Use interface methods instead of field access
- Use conversion functions during transition

### Breaking Changes Summary

| Change | Impact | Migration Path |
|--------|---------|----------------|
| Move `Chrono` to experimental | Low | Import `experimental` if needed |
| Hide concrete result types | Medium | Use `Result` interface instead |
| Move `common/` to `internal/` | Medium | Use `experimental` for extensibility |
| Hide `en` parser types | Low | Created via factories, not directly |
| Simplify `Settings` | Low | Use builder methods instead |
| Remove duplicate interfaces | Low | Use `Result`/`Components` instead |

---

## Tradeoffs & Considerations

### Extensibility vs Simplicity

**Tradeoff:** Hiding parser/refiner infrastructure makes the API cleaner but limits extensibility.

**Decision:** Provide extensibility through `experimental` package. Users who need it can opt into complexity.

**Rationale:**
- Most users (95%+) will never create custom parsers
- Those who do are sophisticated enough to use an experimental API
- Keeps main API focused on common use cases
- Can later stabilize experimental features if there's demand

### Backward Compatibility vs API Cleanliness

**Tradeoff:** Breaking changes will affect existing users, but the current API is confusing with duplicate types and leaked internals.

**Decision:** Make breaking changes but provide migration path and clear deprecation timeline.

**Rationale:**
- Current API has technical debt from evolution
- Better to fix now before more users depend on it
- Clear migration guide and helpers reduce pain
- Deprecation warnings give time to adapt
- Long-term maintainability is worth short-term disruption

### Type Safety vs Simplicity

**Tradeoff:** Exposing types like `Weekday`, `Month` provides type safety, but adds API surface.

**Decision:** Hide custom types, use stdlib types (`time.Weekday`, `time.Month`) where possible.

**Rationale:**
- Users already import `time` package
- Reduces learning curve (fewer new types)
- Values are interoperable with stdlib
- Internal code can still use custom types

### Global State vs Local Configuration

**Tradeoff:** `GlobalRegistry` provides convenient parser registration but creates global mutable state.

**Decision:** Move to experimental and recommend builder pattern for most users.

**Rationale:**
- Global state complicates testing
- Not needed for 95% of use cases
- Power users can still access it
- Builder pattern is more idiomatic Go

---

## Recommended Action Plan

### Immediate Actions (Do First)

1. **Create API reduction proposal document** ✅ (This document)
2. **Share with maintainers and key users** for feedback
3. **Create `experimental` package** with copied types
4. **Add deprecation notices** to types that will move
5. **Update documentation** to recommend new patterns

### Short-term Actions (Next Release - v1.x)

1. **Move `common/` to `internal/common/`**
   - Update imports in `en/` package
   - Verify tests pass
   - Add migration guide

2. **Clean up `en/` package**
   - Replace `Casual`, `Strict`, `GB` vars with functions
   - Hide parser/refiner implementations
   - Hide data dictionaries

3. **Hide internal types**
   - Move `ParsingResult`, `ParsingComponents` to `internal/`
   - Keep interface adapters
   - Provide conversion helpers

4. **Simplify exports**
   - Remove duplicate interfaces (`ParsedResult`, `ParsedComponents`)
   - Hide utility types (`Duration`, `Timeunit`, etc.)
   - Keep only essential constants

### Long-term Actions (v2.0)

1. **Remove deprecated exports entirely**
2. **Clean up experimental package** - promote stable features, remove incomplete ones
3. **Consider stabilizing extensibility API** if there's demand
4. **Review and optimize internal architecture** now that internals are truly private

---

## Success Metrics

### API Reduction Target

- **Before:** 200+ exported identifiers
- **After:** 20-30 exported identifiers in main package
- **Reduction:** ~85-90%

### Specific Targets

| Package | Current Exports | Target Exports | Reduction |
|---------|----------------|----------------|-----------|
| `kronos` | ~120 | ~20 | 83% |
| `kronos/en` | ~50 | ~6 | 88% |
| `kronos/common` | ~30 | 0 (moved to internal) | 100% |
| `kronos/experimental` | 0 | ~30 (moved from main) | N/A |

### Quality Metrics

- ✅ All tests pass after each phase
- ✅ Documentation updated for new API
- ✅ Migration guide with examples
- ✅ Deprecation warnings in place
- ✅ No reduction in functionality for common use cases

---

## Appendix A: Detailed Export List

### Current Main Package Exports (120+)

**Types (30):**
- Chrono, ParserBuilder, Configuration
- Parser, Refiner
- Result, Components, ParsedResult, ParsedComponents
- ParsingResult, ParsingComponents, ParsingContext, ReferenceWithTimezone
- ParsingOption, ParsingReference, Settings
- Component, Timeunit, Weekday, Month, Meridiem, Period
- DateOrder, DatePreference, DayPreference
- Duration, TimezoneAbbrMap, AmbiguousTimezoneMap
- DebugHandler
- ParserRegistry, ParserFactory, ParserInfo

**Functions (15):**
- New, NewChrono, Parse, ParseDate
- ParserBuilder methods (10)
- Registry functions

**Constants (70+):**
- Component constants (12)
- Timeunit constants (12)
- Weekday constants (7)
- Month constants (12)
- Meridiem constants (2)
- Period constants (5)
- DateOrder constants (3)
- DatePreference constants (3)
- DayPreference constants (3)
- Numeric constants (20+)

**Variables (5):**
- GlobalRegistry
- EmptyDuration
- DefaultTimezoneAbbrMap
- ApproximationWords

### Proposed Main Package Exports (20)

**Types (6):**
- ParserBuilder
- Result, Components (interfaces)
- Component (string)
- DateOrder (int)
- DatePreference (int)

**Functions (3):**
- New(*Chrono) *ParserBuilder
- Parse(string, *Chrono) ([]Result, error)
- ParseDate(string, *Chrono) (*time.Time, error)

**Methods (10):**
- ParserBuilder.WithReferenceDate
- ParserBuilder.Strict, Casual
- ParserBuilder.DateOrder
- ParserBuilder.PreferPast, PreferFuture, PreferCurrentPeriod
- ParserBuilder.Timezone
- ParserBuilder.Parse
- ParserBuilder.ParseDate

**Constants (12):**
- Component constants only (ComponentYear, ComponentMonth, etc.)
- DateOrder constants (3)
- DatePreference constants (3)

---

## Appendix B: Code Examples

### Before: Complex API

```go
package main

import (
    "fmt"
    "time"
    "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

func main() {
    // Multiple ways to do the same thing - confusing!

    // Way 1: Using exported Chrono instance
    chrono := en.Casual
    results := chrono.Parse("tomorrow at 3pm", time.Now(), nil)
    if len(results) > 0 {
        fmt.Println(results[0].Date())
    }

    // Way 2: Building custom Chrono
    config := &kronos.Configuration{
        Parsers:  []kronos.Parser{/* ... */},
        Refiners: []kronos.Refiner{/* ... */},
    }
    customChrono := kronos.NewChrono(config)

    // Way 3: Using Settings
    settings := kronos.Settings{
        DateOrder: kronos.DateOrderMDY,
        PreferDatesFrom: kronos.PreferPast,
        // ... 20 other fields
    }
    results2, _ := customChrono.ParseWithSettings("tomorrow", time.Now(), settings)

    // Way 4: Using builder (new)
    parser := en.New().PreferPast()
    results3, _ := parser.Parse("tomorrow")

    // Accessing concrete types directly
    pr := results[0]  // *kronos.ParsingResult
    pc := pr.start    // *kronos.ParsingComponents (if exported)

    // Using types that shouldn't be public
    duration := kronos.Duration{
        kronos.TimeunitDay: 1,
    }
}
```

### After: Clean API

```go
package main

import (
    "fmt"
    "github.com/kljensen/kronos/en"
    "github.com/kljensen/kronos"
)

func main() {
    // Simple case: just parse
    date, err := en.ParseDateSimple("tomorrow at 3pm")
    if err != nil {
        panic(err)
    }
    fmt.Println(date)

    // Configure parsing behavior
    parser := en.New().
        PreferPast().
        DateOrder(kronos.DateOrderDMY)

    results, err := parser.Parse("15/03/2024 at 3pm")
    if err != nil {
        panic(err)
    }

    // Inspect results using clean interfaces
    for _, result := range results {
        fmt.Printf("Found '%s' at index %d\n", result.Text(), result.Index())

        comp := result.Start()
        if comp.IsCertain(kronos.ComponentHour) {
            hour := comp.Get(kronos.ComponentHour)
            fmt.Printf("Hour was mentioned: %d\n", *hour)
        }

        fmt.Printf("Full date: %v\n", result.Date())
    }
}

// Advanced users who need custom parsers
func advancedUsage() {
    // Import experimental package for extensibility
    import "github.com/kljensen/kronos/experimental"

    config := &experimental.Configuration{
        Parsers:  []experimental.Parser{/* custom */},
        Refiners: []experimental.Refiner{/* custom */},
    }
    chrono := experimental.NewChrono(config)
    parser := kronos.New(chrono)
    // ...
}
```

---

## Appendix C: Implementation Checklist

### Phase 1: Experimental Package ✅

- [ ] Create `kronos/experimental/` directory
- [ ] Move `Chrono` type and methods
- [ ] Move `Configuration` type
- [ ] Move `Parser` and `Refiner` interfaces
- [ ] Move `ParserRegistry` and related types
- [ ] Add deprecation comments in main package
- [ ] Update tests to use new package
- [ ] Update documentation

### Phase 2: Internal Package Reorganization ✅

- [ ] Create `kronos/internal/parsing/` directory
- [ ] Move `ParsingResult`, `ParsingComponents`, `ParsingContext`
- [ ] Move `ReferenceWithTimezone`, `ParsingResultWithBoundary`
- [ ] Create `kronos/internal/common/` directory
- [ ] Move base parser types from `common/`
- [ ] Create `kronos/internal/refiners/` directory
- [ ] Move refiner implementations
- [ ] Update all internal imports
- [ ] Verify tests pass

### Phase 3: English Package Cleanup ✅

- [ ] Replace `Casual` var with `New()` function
- [ ] Replace `Strict` var with `StrictParser()` function
- [ ] Replace `GB` var with `GBParser()` function
- [ ] Hide parser implementation types
- [ ] Hide refiner implementation types
- [ ] Hide data dictionaries (move to internal)
- [ ] Hide constants (move to internal)
- [ ] Update documentation
- [ ] Verify tests pass

### Phase 4: Type System Simplification ✅

- [ ] Remove `ParsedResult` interface (use `Result`)
- [ ] Remove `ParsedComponents` interface (use `Components`)
- [ ] Move `Timeunit` to internal
- [ ] Move `Weekday` to internal (use `time.Weekday`)
- [ ] Move `Month` to internal (use `time.Month`)
- [ ] Move `Meridiem` to internal
- [ ] Move `Duration` to internal
- [ ] Remove or move `Period` to experimental
- [ ] Hide numeric constants
- [ ] Update all usages

### Phase 5: Settings Simplification ✅

- [ ] Keep `Settings` as internal implementation
- [ ] Document which settings map to builder methods
- [ ] Remove or deprecate experimental fields (`ToTimezone`)
- [ ] Move advanced settings to experimental
- [ ] Update documentation
- [ ] Add migration examples

### Phase 6: Documentation & Migration ✅

- [ ] Create migration guide
- [ ] Update README with new examples
- [ ] Update GoDoc comments
- [ ] Create deprecation timeline document
- [ ] Add examples for common use cases
- [ ] Add examples for experimental package
- [ ] Document breaking changes
- [ ] Create release notes

---

## Conclusion

The kronos library has evolved organically and accumulated a large public API surface. By reducing exports from 200+ to ~20-30 in the main package, we can:

1. **Simplify** the learning curve for new users
2. **Focus** the API on common use cases (parse text → get time.Time)
3. **Hide** implementation details and give us flexibility to refactor
4. **Provide** clear extensibility path through experimental package
5. **Improve** long-term maintainability

The recommended approach balances backward compatibility (through deprecation and migration helpers) with long-term API cleanliness. Users who just want to parse dates will have a simple API, while advanced users can opt into complexity through the experimental package.

**Next Steps:**
1. Review this document with maintainers
2. Gather feedback from existing users
3. Refine the plan based on feedback
4. Begin implementation in phases
5. Communicate changes clearly through deprecation warnings and documentation

---

**Document Version:** 1.0
**Last Updated:** 2025-11-05
**Author:** Claude Code
**Status:** Draft for Review
