# Follow-up Issues for API Audit

This document contains the follow-up issues identified in the API audit (Issue #144).

---

## Issue #145: Move Advanced APIs to Experimental Package

**Priority:** High

### Context
Following the API audit (Issue #144), we identified 13 exports that are advanced/unstable APIs not intended for general use. These should be moved to `github.com/kljensen/kronos/experimental` where they can evolve without breaking the main API contract.

### Exports to Move

#### Types (7)
- `Chrono` - Advanced parsing engine with custom parser/refiner lists
- `Configuration` - Advanced configuration struct (Parsers, Refiners)
- `ParserFactory` - Factory function type for creating parsers
- `ParserInfo` - Parser metadata (Name, Description, Priority, Tags)
- `ParserRegistry` - Registry for managing available parsers
- `Pipeline` - Internal pipeline executor
- `Refiner` - Refiner interface for post-processing results

#### Functions (5)
- `NewChrono(config *Configuration) *Chrono`
- `NewParserRegistry() *ParserRegistry`
- `NewPipeline(config, settings) *Pipeline`
- `NewPipelineWithSettings(config, settings) (*Pipeline, error)`
- `Register(name, info, factory)`

#### Variables (1)
- `GlobalRegistry` - Global parser registry

### Tasks
1. Create `experimental/` package directory
2. Move 13 exports to experimental package
3. Add re-exports in root package with deprecation notices
4. Update any internal references to use experimental package
5. Update documentation to reference experimental package
6. Run tests to ensure nothing breaks
7. Update examples if needed

### Migration Strategy
Phase 1 (v1.x): Re-export from root with deprecation notices
```go
// Deprecated: Use github.com/kljensen/kronos/experimental instead.
// This re-export will be removed in v2.0.
type Chrono = experimental.Chrono
```

Phase 2 (v2.0): Remove re-exports from root

### Benefits
- Clearer API boundaries
- Allows experimental APIs to evolve without breaking main API
- Reduces cognitive load for new users (60 essential vs 127 total)
- Enables faster iteration on advanced features

### Definition of Done
- All 13 exports moved to experimental package
- Re-exports in root with deprecation notices
- All tests pass
- Documentation updated
- No breaking changes for existing code (v1.x compatibility)

---

## Issue #146: Deprecate DayPreference Type

**Priority:** Medium

### Context
Following the API audit (Issue #144), we identified `DayPreference` as a legacy type that has been superseded by the more general `DatePreference` configuration. It should be marked for deprecation and removed in v2.0.

### Exports to Deprecate (4)

#### Type
- `DayPreference` - Replaced by DatePreference

#### Constants
- `DayPreferCurrent` - Use PreferCurrentPeriod instead
- `DayPreferFirst` - Use PreferPast instead
- `DayPreferLast` - Use PreferFuture instead

### Tasks
1. Add deprecation notices to all 4 exports
2. Update any internal code still using DayPreference
3. Document migration path in package docs
4. Plan removal for v2.0

### Migration Path

#### Before
```go
settings.DayPreference = kronos.DayPreferFirst
```

#### After
```go
settings.DatePreference = kronos.PreferPast
```

### Deprecation Notice Template
```go
// Deprecated: Use DatePreference instead.
// DayPreferFirst is equivalent to PreferPast.
// This type will be removed in v2.0.
type DayPreference int

const (
    // Deprecated: Use PreferCurrentPeriod instead. Will be removed in v2.0.
    DayPreferCurrent DayPreference = iota
    // Deprecated: Use PreferPast instead. Will be removed in v2.0.
    DayPreferFirst
    // Deprecated: Use PreferFuture instead. Will be removed in v2.0.
    DayPreferLast
)
```

### Impact Assessment
**Very low.** No examples use DayPreference, suggesting it was superseded before gaining adoption.

### Definition of Done
- All 4 exports have deprecation notices
- Internal code migrated to DatePreference
- Migration path documented
- All tests pass

---

## Issue #147: Update API Guard Target After Experimental Move

**Priority:** High (but depends on Issue #145)

### Context
After moving 13 experimental exports to the experimental package (Issue #145), we should update the MAX_EXPORTS guard to reflect the new, tighter API surface and prevent future bloat.

### Current State
- MAX_EXPORTS: 150
- Current exports: 127
- Headroom: 23

### After Moving Experimental APIs
- Essential: 60
- Deprecated: 4
- Internal (X-prefixed): 50
- **Total: 114**
- Recommended MAX_EXPORTS: **80**

### Rationale for 80
The target of 80 provides:
- Comfortable headroom for the 60 essential exports
- Accounts for the 50 internal (X-prefixed) exports that can't be moved
- Ignores the 4 deprecated exports (will be removed in v2.0)
- Leaves room for ~20 new essential exports
- Prevents API bloat by being strict

### Tasks
1. Wait for Issue #145 to complete (move experimental APIs)
2. Update MAX_EXPORTS in api_surface_test.go from 150 to 80
3. Run tests to verify we're under the limit
4. Update test comments to reflect new target

### Code Change
```go
// Before:
const MAX_EXPORTS = 150 // Reduced from 200, allows headroom while preventing API bloat

// After:
const MAX_EXPORTS = 80 // Tight limit for essential API (60) + internal helpers (50) - experimental (13) - deprecated (4)
```

### Dependencies
- Depends on Issue #145 (Move Advanced APIs to Experimental Package)

### Definition of Done
- MAX_EXPORTS updated to 80
- api_surface_test.go passes
- Test comments updated with new rationale

---

## Summary

The three follow-up issues address different aspects of API minimization:

1. **Issue #145** (High priority) - Moves advanced APIs to experimental package, reducing the public API surface by 13 exports
2. **Issue #146** (Medium priority) - Deprecates legacy DayPreference type, removing 4 exports in v2.0
3. **Issue #147** (High priority) - Updates API guard to enforce the new, tighter limits (80 vs 150)

After completing these issues, the kronos public API will be:
- **60 essential exports** - Stable, recommended for general use
- **50 internal exports** - X-prefixed helpers for internal/ packages
- **13 experimental exports** - Available via experimental package for advanced users
- **4 deprecated exports** - To be removed in v2.0

Total public API surface: **110 exports** (down from 127)
Target MAX_EXPORTS: **80** (encourages minimalism)
