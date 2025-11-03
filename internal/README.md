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

1. **Phase 1: Create structure** (Current)
   - Establish package structure
   - Document intended organization
   - Get feedback on approach

2. **Phase 2: Move core types**
   - Move `ParsingComponents`, `ParsingResult`, `ParsingContext`, `ReferenceWithTimezone`
   - Update root package to re-export for backward compatibility
   - Ensure all tests pass

3. **Phase 3: Move parser implementations**
   - Gradually move parser types to `internal/parsers`
   - Keep registration mechanisms in place
   - Maintain test coverage

4. **Phase 4: Move refiner implementations**
   - Move refiner types to `internal/refiners`
   - Update configuration to reference internal types
   - Verify end-to-end functionality

5. **Phase 5: Build new public API** (Future issue)
   - Implement builder-based `Parser` type
   - Create simplified `Result` and `Components` interfaces
   - Add convenience functions `Parse()`, `ParseDate()`, `New()`

6. **Phase 6: Deprecation** (Future issue)
   - Mark old APIs as deprecated
   - Provide migration guide
   - Eventually remove deprecated APIs in v2.0

## Design Principles

- **Incremental migration**: Small, reviewable changes
- **Maintain compatibility**: Keep existing APIs working during transition
- **Test coverage**: Every change verified by existing test suite
- **Clear dependencies**: Avoid circular imports between internal packages
- **Simple first**: Start with easiest migrations, tackle complex cases later

## Import Guidelines

- Internal packages MAY import from the root `kronos` package for public types
- Internal packages SHOULD NOT create circular dependencies
- Root package WILL re-export internal types during transition for backward compatibility
- New code SHOULD use the builder API once available
