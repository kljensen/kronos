# Internal Packages

This directory contains implementation details for the Kronos date parsing library.
These packages are internal and not part of the public API.

## Package Structure

### `internal/core`
Core types and data structures:
- `ParsingComponents` - Internal representation of parsed date/time components
- `ParsingResult` - Internal result type with full parsing details
- `ParsingContext` - Internal parsing context with reference date and options
- `ReferenceWithTimezone` - Internal reference date handling

The public API will expose simpler `Result` and `Components` interfaces that hide
these implementation details.

### `internal/parsers`
All parser implementations:
- ISO 8601 parsers
- Slash date format parsers
- Month name parsers
- Relative date parsers
- Time expression parsers
- Casual reference parsers

Users interact with these through the builder-based `Parser` API, not directly.

### `internal/refiners`
All refiner implementations:
- Result merging (datetime, date ranges, weekday)
- Date preference handling
- Timezone resolution
- Forward date filtering
- Overlap removal
- Format validation

Users configure refinement through builder methods, not by accessing refiners directly.

## Migration Plan

The migration to internal packages will happen incrementally:

### Phase 1: Create structure (CURRENT - Issue #91)
**Status:** In Progress
**Goal:** Establish foundation and get architectural feedback

Checklist:
- [x] Create internal/ directory structure
- [x] Create internal/README.md with documentation
- [x] Document fundamental public types that stay in root
- [x] Document import hierarchy to prevent circular dependencies
- [x] Document re-export strategy for backward compatibility
- [x] Create stub `Result` and `Components` interfaces
- [x] Document all types to be migrated (including ParsingResultWithBoundary)
- [x] Add detailed checklists for each phase
- [ ] Address code review feedback
- [ ] Get approval to proceed to Phase 2

### Phase 2: Move core types (Future Issue)
**Goal:** Move implementation types to internal/core while maintaining compatibility

Checklist:
- [ ] Create internal/core package
- [ ] Move ParsingComponents to internal/core/components.go
  - [ ] Update package declaration
  - [ ] Add imports for fundamental types from root
  - [ ] Update all methods
  - [ ] Verify no circular imports
- [ ] Move ParsingResult to internal/core/result.go
  - [ ] Update to use core.ParsingComponents
  - [ ] Ensure implements ParsedResult interface
- [ ] Move ParsingContext to internal/core/context.go
  - [ ] Update references to core types
- [ ] Move ReferenceWithTimezone to internal/core/reference.go
  - [ ] Update all timezone handling
- [ ] Move ParsingResultWithBoundary to internal/core/boundary.go
  - [ ] Document this is internal-only
- [ ] Add re-exports in root results.go
  - [ ] Add type aliases: `type ParsingComponents = core.ParsingComponents`
  - [ ] Add type aliases: `type ParsingResult = core.ParsingResult`
  - [ ] Add type aliases: `type ParsingContext = core.ParsingContext`
  - [ ] Add type aliases: `type ReferenceWithTimezone = core.ReferenceWithTimezone`
  - [ ] Re-export all constructor functions
  - [ ] Add deprecation comments
- [ ] Update all imports throughout codebase
  - [ ] Update parser files
  - [ ] Update refiner files
  - [ ] Update test files
- [ ] Run full test suite
- [ ] Verify no breaking changes to public API
- [ ] Update documentation

### Phase 3: Move parser implementations (Future Issue)
**Goal:** Move all parser types to internal/parsers package

Checklist:
- [ ] Create internal/parsers package
- [ ] Audit all parser types and categorize them
  - [ ] ISO parsers
  - [ ] Slash format parsers
  - [ ] Month name parsers
  - [ ] Relative date parsers
  - [ ] Time expression parsers
  - [ ] Casual reference parsers
  - [ ] Base/abstract parser types
- [ ] Move parser types in batches (5-10 per PR)
  - [ ] Update package declarations
  - [ ] Update imports to use internal/core
  - [ ] Verify Pattern() and Extract() implementations
- [ ] Keep Parser interface in root (already public)
- [ ] Update registry.go to work with internal parsers
- [ ] Verify all parser registration works
- [ ] Run parser-specific tests
- [ ] Run full integration tests

### Phase 4: Move refiner implementations (Future Issue)
**Goal:** Move all refiner types to internal/refiners package

Checklist:
- [ ] Create internal/refiners package
- [ ] Audit all refiner types
  - [ ] Merging refiners (datetime, date range, weekday)
  - [ ] Date preference refiner
  - [ ] Timezone resolution refiner
  - [ ] Forward date filtering refiner
  - [ ] Overlap removal refiner
  - [ ] Format validation refiners
- [ ] Move refiner types in batches
  - [ ] Update package declarations
  - [ ] Update imports to use internal/core
  - [ ] Verify Refine() implementations
- [ ] Keep Refiner interface in root (already public)
- [ ] Update configuration.go to work with internal refiners
- [ ] Run refiner-specific tests
- [ ] Run full integration tests

### Phase 5: Build new public API (Future Issue)
**Goal:** Implement the new builder-based API

Checklist:
- [ ] Design builder API
  - [ ] Parser builder with fluent methods
  - [ ] Configuration options
  - [ ] Parser/Refiner registration
- [ ] Implement Result interface adapters
  - [ ] Create adapter wrapping ParsingResult
  - [ ] Implement all Result methods
  - [ ] Map Components interface to ParsingComponents
- [ ] Implement Components interface adapters
  - [ ] Create adapter wrapping ParsingComponents
  - [ ] Implement all Components methods
- [ ] Add convenience functions
  - [ ] Parse(text string, options ...Option) ([]Result, error)
  - [ ] ParseDate(text string, options ...Option) (time.Time, error)
  - [ ] New() *Parser builder
- [ ] Write comprehensive examples
- [ ] Update main documentation
- [ ] Write migration guide from old API to new API

### Phase 6: Deprecation (Future Issue - v2.0)
**Goal:** Phase out old APIs and remove deprecated code

Checklist:
- [ ] Add deprecation notices to all old APIs
  - [ ] ParseWithOptions() -> Parse()
  - [ ] Direct access to ParsingComponents -> Components interface
  - [ ] Direct access to ParsingResult -> Result interface
- [ ] Update all examples to use new API
- [ ] Create detailed migration guide
- [ ] Announce deprecation timeline
- [ ] Monitor community feedback
- [ ] In v2.0: Remove all deprecated APIs
- [ ] In v2.0: Remove type aliases
- [ ] In v2.0: Clean up any remaining backward compatibility code

## Design Principles

- **Incremental migration**: Small, reviewable changes
- **Maintain compatibility**: Keep existing APIs working during transition
- **Test coverage**: Every change verified by existing test suite
- **Clear dependencies**: Avoid circular imports between internal packages
- **Simple first**: Start with easiest migrations, tackle complex cases later

## Fundamental Public Types

The following types remain in the root `kronos` package as they are fundamental building blocks used throughout the library and will remain part of the public API per RFC #91:

### Enumerations & Constants
- `Component` - Date/time component identifiers (year, month, day, hour, etc.)
- `Timeunit` - Time unit identifiers for durations (year, month, week, day, etc.)
- `Meridiem` - AM/PM indicator
- `DatePreference` - Controls past/future/current resolution of ambiguous dates
- `Period` - Granularity of parsed expressions (year, month, week, day, time)
- `Weekday` - Day of the week
- `Month` - Calendar month

### Simple Data Types
- `Duration` - Map of timeunit to values (used by parsers and refiners)
- `TimezoneAbbrMap` - Timezone abbreviation mappings
- `AmbiguousTimezoneMap` - DST-aware timezone definitions
- `ParsingOption` - Configuration options (will be wrapped by Settings)
- `ParsingReference` - Reference date/time input

### Public Interfaces
- `ParsedResult` - Public interface for results (currently implemented by ParsingResult)
- `ParsedComponents` - Public interface for components (currently implemented by ParsingComponents)
- `Result` - Future simplified interface for results (Phase 5)
- `Components` - Future simplified interface for components (Phase 5)

### Internal Implementation Types (Will Move)
These types will move to `internal/core` but be re-exported for backward compatibility:
- `ParsingComponents` - Full implementation with tags, implies, etc.
- `ParsingResult` - Full result with reference and indexing
- `ParsingContext` - Parsing context with settings and reference
- `ReferenceWithTimezone` - Internal reference handling
- `ParsingResultWithBoundary` - Internal wrapper for word boundary handling

## Import Guidelines & Circular Dependency Prevention

### Import Hierarchy
```
root (kronos)
  ├─ Fundamental types (Component, Timeunit, Duration, etc.)
  ├─ Public interfaces (Parser, Refiner, ParsedResult, ParsedComponents)
  └─ Re-exports of internal types for backward compatibility
       ↑
       │
internal/core
  ├─ Implementation types (ParsingComponents, ParsingResult, etc.)
  └─ May import fundamental types from root
       ↑
       │
internal/parsers, internal/refiners
  ├─ Parser and refiner implementations
  └─ May import from internal/core and root fundamental types
```

### Rules to Prevent Circular Dependencies
1. **Root package** contains only:
   - Fundamental types (enums, constants, simple types)
   - Public interfaces
   - Re-exports (type aliases) to internal implementations
   - Registry and configuration code

2. **internal/core** may import:
   - Fundamental types from root (Component, Timeunit, Duration, etc.)
   - NEVER import parser or refiner implementations

3. **internal/parsers** and **internal/refiners** may import:
   - Fundamental types from root
   - Implementation types from internal/core
   - NEVER import from each other directly

4. **If a type is needed by multiple packages**, ask:
   - Is it a fundamental concept? → Keep in root
   - Is it core implementation? → Keep in internal/core
   - Is it specific functionality? → Keep in its own package

## Re-export Strategy for Backward Compatibility

During the migration, we maintain backward compatibility using type aliases. This allows existing code to continue working while we move implementations to internal packages.

### Phase 2 Example (Moving Core Types)
```go
// In internal/core/components.go
package core

import "github.com/example/kronos"

// ParsingComponents is the full implementation
type ParsingComponents struct {
    knownValues   map[kronos.Component]int
    impliedValues map[kronos.Component]int
    // ...
}

// In results.go (root package)
package kronos

import "github.com/example/kronos/internal/core"

// ParsingComponents is re-exported for backward compatibility
// Deprecated: Use internal/core.ParsingComponents directly or the new Components interface
type ParsingComponents = core.ParsingComponents

// NewParsingComponents is re-exported for backward compatibility
// Deprecated: Use internal/core.NewParsingComponents directly
func NewParsingComponents(reference *ReferenceWithTimezone, knownComponents map[Component]int) *ParsingComponents {
    return core.NewParsingComponents(reference, knownComponents)
}
```

### Type Alias Best Practices
- Use `type Alias = Original` for re-exporting types
- Re-export constructor functions that return the aliased type
- Add deprecation comments pointing to the new location or new API
- Keep aliases until v2.0 when we can break compatibility
- Test that aliases work exactly like the original types
