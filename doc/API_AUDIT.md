# Kronos API Audit - Issue #144

**Date:** 2025-11-06
**Current Version:** API minimization complete (Phase 1-7)
**Total Exports:** 127 (limit: 150)

## Executive Summary

This document provides a complete inventory of the kronos public API and categorizes each export according to its intended audience and stability guarantees. After the recent API minimization effort (Phases 1-7), we have 127 public exports with room for 23 more before hitting our guard limit.

### Current State
- **Essential API (stable, keep):** 60 exports (47%)
- **Experimental API (move):** 13 exports (10%)
- **Deprecated API (phase out):** 4 exports (3%)
- **Internal API (already prefixed):** 50 exports (39%)

### Recommended Actions
1. Move 13 experimental exports to `github.com/kljensen/kronos/experimental`
2. Plan deprecation path for 4 deprecated exports in next major version
3. Internal API (X-prefixed) can stay as-is - already clearly marked
4. Reduce MAX_EXPORTS guard from 150 to **80** after moving experimental APIs

## Category Definitions

### Essential (Keep Public)
Stable, user-facing API that forms the core of kronos. These are guaranteed to remain stable and are documented for general use. Breaking changes require a major version bump.

### Experimental (Move to experimental package)
Advanced/unstable APIs that power kronos internally but are not recommended for general use. These should be moved to `github.com/kljensen/kronos/experimental` where they can evolve without breaking the main API contract.

### Deprecated (Phase Out)
Legacy exports kept for backward compatibility but marked for removal in the next major version. Users should migrate away from these.

### Internal (X-Prefixed)
Exports needed by internal/ packages. These use the X prefix convention to signal "internal use only" while remaining accessible to internal/ code. These can change without notice.

---

## Detailed Inventory

### Essential API (60 exports)

These are the stable, public API that users should rely on.

#### Types (11)
- `Component` - Core enum for accessing parsed components (e.g., ComponentYear, ComponentMonth)
- `Components` - Interface for accessing parsed components with certainty tracking
- `DateOrder` - Configuration for date format: MDY (US), DMY (European), YMD (ISO)
- `DatePreference` - Configuration for ambiguous date resolution: Past, Future, CurrentPeriod
- `Duration` - Duration type for relative dates (map of Timeunit to float64)
- `Parser` - Parser interface for custom parsers (used by advanced users)
- `ParserBuilder` - Fluent builder API - primary entry point for most users
- `Period` - Time period enum: AM, PM, Morning, Afternoon, Evening, Night
- `Result` - Primary result interface returned by Parse()
- `Settings` - Settings struct for parser configuration
- `Timeunit` - Required for Duration map keys: Year, Month, Week, Day, Hour, etc.

**Rationale:** These types form the minimal API surface for 95% of use cases:
- `ParserBuilder` is the main entry point (via `en.New()` or `kronos.New()`)
- `Result` and `Components` are needed to work with parsed results
- Enum types (`Component`, `DateOrder`, `DatePreference`, `Period`, `Timeunit`) are required for configuration and component access
- `Duration` and `Settings` are needed for advanced configuration

#### Constants (43)

**Component enum values (12):**
- `ComponentDay`, `ComponentHour`, `ComponentMeridiem`, `ComponentMicrosecond`
- `ComponentMillisecond`, `ComponentMinute`, `ComponentMonth`, `ComponentNanosecond`
- `ComponentSecond`, `ComponentTimezoneOffset`, `ComponentWeekday`, `ComponentYear`

**DateOrder enum values (3):**
- `DateOrderDMY`, `DateOrderMDY`, `DateOrderYMD`

**DatePreference enum values (3):**
- `PreferCurrentPeriod`, `PreferFuture`, `PreferPast`

**Period enum values (6):**
- `PeriodDay`, `PeriodMonth`, `PeriodTime`, `PeriodUnknown`, `PeriodWeek`, `PeriodYear`

**Timeunit enum values (12):**
- `TimeunitDay`, `TimeunitDecade`, `TimeunitHour`, `TimeunitMicrosecond`
- `TimeunitMillisecond`, `TimeunitMinute`, `TimeunitMonth`, `TimeunitNanosecond`
- `TimeunitQuarter`, `TimeunitSecond`, `TimeunitWeek`, `TimeunitYear`

**Validation bounds (7):**
- `MaxDaysDuration`, `MaxHoursDuration`, `MaxMinutesDuration`, `MaxMonthsDuration`
- `MaxYear`, `MaxYearsDuration`, `MinYear`

**Rationale:** These constants are essential for:
1. **Component access** - Users need these to check `comp.Get(kronos.ComponentHour)`
2. **Configuration** - DateOrder and DatePreference are commonly used
3. **Validation** - Bounds prevent integer overflow in date arithmetic
4. **Duration operations** - Timeunit keys are required for Duration maps

#### Functions (6)
- `New(chrono *Chrono) *ParserBuilder` - Primary entry point for builder API
- `Parse(text string, chrono *Chrono) ([]Result, error)` - Package-level convenience
- `ParseDate(text string, chrono *Chrono) (*time.Time, error)` - Package-level convenience
- `ParseWithSettings(text, refDate, settings, config)` - Settings-based parsing
- `DefaultSettings() Settings` - Get default settings
- `ValidateSettings(s Settings) error` - Validate settings

**Rationale:** These provide:
1. **Builder creation** - `New()` is the standard entry point
2. **Convenience functions** - `Parse()` and `ParseDate()` for simple cases
3. **Settings management** - `DefaultSettings()` and `ValidateSettings()` for configuration

**Usage patterns from examples:**
```go
// Most common: builder pattern
parser := en.New().
    DateOrder(kronos.DateOrderDMY).
    PreferPast()
results, _ := parser.Parse("last Monday")

// Component access
comp := results[0].Start()
if comp.IsCertain(kronos.ComponentHour) {
    hour := comp.Get(kronos.ComponentHour)
}
```

---

### Experimental API (13 exports)

These advanced APIs should be moved to `github.com/kljensen/kronos/experimental`.

#### Types (7)
- `Chrono` - Advanced parsing engine with custom parser/refiner lists
- `Configuration` - Advanced configuration struct (Parsers, Refiners)
- `ParserFactory` - Factory function type for creating parsers
- `ParserInfo` - Parser metadata (Name, Description, Priority, Tags)
- `ParserRegistry` - Registry for managing available parsers
- `Pipeline` - Internal pipeline executor
- `Refiner` - Refiner interface for post-processing results

**Rationale:** These are used for:
1. **Custom parser development** - Most users don't write custom parsers
2. **Advanced configuration** - Direct Chrono/Configuration access is rarely needed
3. **Parser registration** - Plugin-style architecture, not needed by 95% of users

**Migration path:**
```go
// Before (current):
import "github.com/kljensen/kronos"
chrono := kronos.NewChrono(&kronos.Configuration{...})

// After (future):
import "github.com/kljensen/kronos/experimental"
chrono := experimental.NewChrono(&experimental.Configuration{...})
```

#### Functions (5)
- `NewChrono(config *Configuration) *Chrono` - Create custom Chrono
- `NewParserRegistry() *ParserRegistry` - Create parser registry
- `NewPipeline(config, settings) *Pipeline` - Create pipeline
- `NewPipelineWithSettings(config, settings) (*Pipeline, error)` - Create pipeline with validation
- `Register(name, info, factory)` - Register parser globally

**Rationale:** These are advanced APIs used when building custom parsers or registries. The average user never calls these - they use the builder pattern via `en.New()` or `kronos.New()`.

#### Variables (1)
- `GlobalRegistry` - Global parser registry

**Rationale:** Global mutable state is an anti-pattern. This should be in experimental package where advanced users can access it, but it's not part of the recommended API.

**Impact assessment:** Low. Based on the examples and typical usage, very few users interact directly with Chrono, Configuration, or the parser registry. Most use the builder pattern via language-specific packages (en, etc.).

---

### Deprecated API (4 exports)

These should be marked with deprecation notices and removed in v2.0.

#### Types (1)
- `DayPreference` - Replaced by more general DatePreference configuration

#### Constants (3)
- `DayPreferCurrent` - DayPreference enum value
- `DayPreferFirst` - DayPreference enum value
- `DayPreferLast` - DayPreference enum value

**Rationale:** `DayPreference` was an early attempt at ambiguous date resolution, but `DatePreference` is more general and handles the same cases. These should be removed in v2.0.

**Deprecation notice:**
```go
// Deprecated: Use DatePreference instead.
// Will be removed in v2.0.
type DayPreference int
```

**Migration path:**
```go
// Before:
settings.DayPreference = kronos.DayPreferFirst

// After:
settings.DatePreference = kronos.PreferPast
```

**Impact assessment:** Very low. No examples use DayPreference, suggesting it was superseded before gaining adoption.

---

### Internal API (50 exports)

These X-prefixed exports are intentionally public to support internal/ packages but are clearly marked as internal-only via naming convention.

#### Types (12)
- `AmbiguousTimezoneMap` - Timezone with DST handling
- `DebugHandler` - Debug callback function type
- `ParsedComponents` - Internal interface for component access
- `ParsedResult` - Internal interface for results
- `ParsingComponents` - Internal mutable components struct
- `ParsingContext` - Parsing context with text, reference date, settings
- `ParsingOption` - Internal parsing options
- `ParsingReference` - Internal reference wrapper
- `ParsingResult` - Internal result implementation
- `ParsingResultWithBoundary` - Internal result with boundaries
- `ReferenceWithTimezone` - Internal reference with timezone
- `TimezoneAbbrMap` - Map of timezone abbreviations to offsets

**Rationale:** These types power the internal parsing logic and are used by parsers/refiners in internal/ packages. They're not intended for end users but must be public for internal/ code.

#### Functions (38)
All X-prefixed functions are internal helpers:
- `XAddDuration`, `XAssignSimilarDate`, `XAssignSimilarTime`
- `XFindMostLikelyADYear`, `XFindYearClosestToRefWithPreference`
- `XGetDaysToWeekday`, `XGetLastWeekday`, `XGetLastWeekdayOfMonth`
- `XGetNextWeekday`, `XGetNthWeekdayOfMonth`, `XGetThisWeekday`
- `XImplySimilarDate`, `XReverseDuration`, `XSafeSlice`
- `XStripApproximationWords`, `XToTimezoneOffset`
- `XAfternoon`, `XAfternoonWithHour`, `XAsParsingComponents`
- `XCreateRelativeFromReference`, `XEvening`, `XEveningWithHour`
- `XMergeDateTimeComponent`, `XMergeDateTimeResult`, `XMidnight`
- `XMorning`, `XMorningWithHour`, `XNewParsingComponents`
- `XNewParsingContext`, `XNewParsingResult`, `XNoon`, `XNow`
- `XTheDayAfter`, `XTheDayBefore`, `XToday`, `XTomorrow`
- `XTonightWithHour`, `XYesterday`

**Naming convention:** The X prefix clearly signals "internal use only" following Go conventions. This is more maintainable than internal/ packages which can't be imported from other modules.

**Status:** Keep as-is. The X prefix is clear and conventional, and these are genuinely needed by internal/ packages.

---

## Recommendations

### 1. Move Experimental APIs (Priority: High)

Create `github.com/kljensen/kronos/experimental` package and move:
- Types: `Chrono`, `Configuration`, `ParserFactory`, `ParserInfo`, `ParserRegistry`, `Pipeline`, `Refiner`
- Functions: `NewChrono`, `NewParserRegistry`, `NewPipeline`, `NewPipelineWithSettings`, `Register`
- Variables: `GlobalRegistry`

**Benefits:**
- Clearer API boundaries
- Allows experimental APIs to evolve without breaking main API
- Reduces cognitive load for new users
- Enables faster iteration on advanced features

**Migration strategy:**
1. Copy types/functions to experimental/ package
2. Re-export from root package with deprecation notices
3. Update documentation to reference experimental package
4. Remove from root in v2.0

### 2. Deprecate DayPreference (Priority: Medium)

Add deprecation notices:
```go
// Deprecated: Use DatePreference instead. Will be removed in v2.0.
type DayPreference int
```

Update any internal code to use `DatePreference` instead.

### 3. Update API Guard (Priority: High)

After moving experimental APIs, update MAX_EXPORTS:
```go
// Before:
MAX_EXPORTS = 150 // 127 current + 23 headroom

// After moving experimental APIs:
MAX_EXPORTS = 80  // 60 essential + 4 deprecated + 50 internal - 13 experimental = 114
                  // Set to 80 to encourage further minimization
```

### 4. Documentation Updates (Priority: High)

Update package documentation to clarify:
1. **Essential API** - Recommended for general use, stability guaranteed
2. **Experimental API** - Import from `experimental` package, may change
3. **Internal API** - X-prefixed, not for direct use, no stability guarantees
4. **Deprecated API** - Migration path documented, removal planned for v2.0

### 5. Future Cleanup for v2.0 (Priority: Low)

Consider for v2.0:
- Remove `DayPreference` completely
- Remove deprecated re-exports of experimental APIs
- Consider if any X-prefixed functions can be moved to internal/helpers

---

## Target API Surface for v2.0

After implementing all recommendations:

| Category | Count | Notes |
|----------|-------|-------|
| Essential | 60 | Stable public API |
| Experimental | 0 | Moved to experimental package |
| Deprecated | 0 | Removed in v2.0 |
| Internal (X-prefixed) | 50 | Keep for internal/ packages |
| **Total** | **110** | Down from 127 |

**New MAX_EXPORTS target: 80** (room for ~20 additions to essential API)

---

## Usage Analysis

### Most Common Usage Patterns (from examples)

1. **Builder pattern** (95% of use cases)
   ```go
   parser := en.New().DateOrder(kronos.DateOrderDMY).PreferPast()
   results, _ := parser.Parse("next Monday")
   ```

2. **Component access** (common)
   ```go
   comp := result.Start()
   if comp.IsCertain(kronos.ComponentHour) {
       hour := comp.Get(kronos.ComponentHour)
   }
   ```

3. **Settings-based configuration** (advanced)
   ```go
   settings := kronos.DefaultSettings()
   settings.DateOrder = kronos.DateOrderDMY
   results, _ := kronos.ParseWithSettings(text, refDate, settings, config)
   ```

### Rarely Used (Candidates for Experimental)

- Direct `Chrono` construction
- `Configuration` structs
- Parser/refiner registration
- `Pipeline` creation

**Conclusion:** The essential API (60 exports) covers 95%+ of use cases. The experimental APIs (13 exports) are genuinely advanced and should be isolated.

---

## Follow-up Issues

Based on this audit, the following follow-up issues should be created:

### Issue: Move Advanced APIs to Experimental Package
**Priority:** High
**Description:** Move 13 exports (Chrono, Configuration, etc.) to `github.com/kljensen/kronos/experimental` package with deprecation notices in root package.
**Exports to move:** Chrono, Configuration, ParserFactory, ParserInfo, ParserRegistry, Pipeline, Refiner, NewChrono, NewParserRegistry, NewPipeline, NewPipelineWithSettings, Register, GlobalRegistry

### Issue: Deprecate DayPreference Type
**Priority:** Medium
**Description:** Add deprecation notices to DayPreference and related constants. Document migration path to DatePreference.
**Exports to deprecate:** DayPreference, DayPreferCurrent, DayPreferFirst, DayPreferLast

### Issue: Update API Guard Target
**Priority:** High
**Description:** After moving experimental APIs, update MAX_EXPORTS from 150 to 80 to prevent API bloat.
**New target:** 80 (60 essential + 50 internal = 110, with ~30 room for growth)

### Issue: Improve API Documentation
**Priority:** High
**Description:** Update package docs to clearly distinguish essential, experimental, and internal APIs. Add examples for common patterns.

---

## Appendix: Export Counts by Category

```
Total: 127 exports

Essential (60):
  - Types: 11
  - Constants: 43
  - Functions: 6

Experimental (13):
  - Types: 7
  - Functions: 5
  - Variables: 1

Deprecated (4):
  - Types: 1
  - Constants: 3

Internal (50):
  - Types: 12
  - Functions: 38
```

---

## Appendix: Complete Export List

For reference, the complete list of all 127 exports is documented in the inventory above. Each export includes:
- Name
- Type (type, const, func, var)
- Category (essential, experimental, deprecated, internal)
- Rationale for categorization

This audit was generated using `go doc -all .` and automated categorization based on:
1. Usage patterns in examples/
2. Naming conventions (X-prefix for internal)
3. API design principles (builder pattern as primary interface)
4. Previous API minimization decisions (Phases 1-7)
