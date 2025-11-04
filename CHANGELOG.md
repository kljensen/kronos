# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-11-03

Major API refactoring introducing the builder pattern while maintaining backward compatibility.

### Added

#### Builder Pattern API
- **New `ParserBuilder` type** - Main entry point for fluent, type-safe configuration
  - `New(chrono *Chrono) *ParserBuilder` - Create a new parser builder
  - `WithReferenceDate(time.Time)` - Set reference date for relative dates
  - `Strict()` - Enable strict parsing mode
  - `Casual()` - Enable casual parsing mode (default)
  - `DateOrder(DateOrder)` - Configure date component order (MDY, DMY, YMD)
  - `PreferPast()` - Prefer past dates when ambiguous
  - `PreferFuture()` - Prefer future dates when ambiguous
  - `PreferCurrentPeriod()` - Prefer current period (default)
  - `Timezone(string)` - Set default timezone
  - `ToTimezone(string)` - Set target timezone for conversion (experimental)
  - `ForwardDate()` - Enable forward date parsing (legacy compatibility)
  - `Parse(text string) ([]Result, error)` - Parse all dates with error handling
  - `ParseDate(text string) (*time.Time, error)` - Parse first date with error handling

#### Public Interfaces
- **`Result` interface** - Clean public API for parsed results
  - `Text() string` - Get matched text
  - `Index() int` - Get position in input
  - `Date() time.Time` - Get parsed date
  - `Start() Components` - Get start components
  - `End() Components` - Get end components (nil for non-ranges)

- **`Components` interface** - Access to date/time parts
  - `Get(component Component) *int` - Get component value
  - `IsCertain(component Component) bool` - Check if explicitly mentioned
  - `Date() time.Time` - Convert to time.Time

#### Convenience Functions
- `en.ParseSimple(text string) ([]Result, error)` - Quick parse with defaults
- `en.ParseDateSimple(text string) (*time.Time, error)` - Quick single date parse
- `en.New()` - Create parser builder with casual English
- `en.StrictParser()` - Create parser builder with strict English
- `en.GBParser()` - Create parser builder with British English

#### Package-Level Functions
- `kronos.Parse(text, chrono) ([]Result, error)` - Package-level parse
- `kronos.ParseDate(text, chrono) (*time.Time, error)` - Package-level single date parse

#### Settings System
- **`Settings` struct** - Comprehensive configuration for advanced use cases
  - `DateOrder` - Date component order (MDY, DMY, YMD)
  - `PreferDatesFrom` - Date preference (Past, Future, CurrentPeriod)
  - `PreferDayOfMonth` - Day preference (Current, First, Last)
  - `Timezone` - Default timezone
  - `ToTimezone` - Target timezone for conversion
  - `StrictParsing` - Enable strict mode
  - `Normalize` - Unicode normalization
  - `SkipTokens` - Words to ignore
  - `RequireParts` - Required components
  - `RelativeBase` - Base time for relative dates
  - And more...

- `DefaultSettings()` - Get default settings
- `ValidateSettings(Settings) error` - Validate settings
- `ParseWithSettings(text, refDate, settings, config) ([]*ParsingResult, error)` - Parse with settings

#### Pipeline Architecture
- **`Pipeline` system** - Configurable parsing pipeline
  - Input sanitization stage
  - Parser execution stage
  - Refiner execution stage
  - Settings application stage
  - Timezone conversion stage (experimental)
  - Result validation stage

#### Date Order Support
- **`DateOrder` type** with constants:
  - `DateOrderMDY` - US format (Month/Day/Year)
  - `DateOrderDMY` - European format (Day/Month/Year)
  - `DateOrderYMD` - ISO format (Year/Month/Day)

#### Documentation
- Comprehensive `README.md` with examples and API documentation
- Package-level godoc in `doc.go`
- `MIGRATION.md` guide for upgrading from v1
- `examples/` directory with runnable examples:
  - `basic/` - Basic usage
  - `configured/` - Builder configuration
  - `components/` - Component inspection
  - `ranges/` - Date ranges
  - `british/` - British English

#### Internal Improvements
- Result adapter pattern for clean public API
- Component adapter for interface implementation
- Registry system for parser/refiner management
- Improved leap year handling for ambiguous dates
- Enhanced safety and input sanitization

### Changed

#### API Changes
- Parse methods now return errors: `(results, error)` instead of just `results`
- Result type changed from `*ParsingResult` to `Result` interface
- Components type changed from `ParsedComponents` to `Components` interface
- Builder pattern is now the recommended API (Chrono direct usage still works)

#### Behavior Changes
- Leap year selection now uses smart algorithm for ambiguous dates (e.g., "Feb 29")
- Input sanitization is more robust (handles zero-width characters, etc.)
- Date order is now configurable (defaults to MDY for backward compatibility)

### Deprecated

- Direct use of `en.Parse()` without error handling (still works, but legacy)
- Direct use of `en.ParseDate()` without error handling (still works, but legacy)
- Direct access to `ParsingResult` internal fields (use Result interface methods)
- Direct construction of `ParsingOption` (use builder methods instead)

### Fixed

- Leap year handling for ambiguous dates like "Feb 29"
- Zero-width character handling in input sanitization
- Parser boundary detection for word-boundary-aware parsers
- Timezone offset parsing edge cases
- Component certainty tracking in merged results

### Backward Compatibility

**All v1 APIs remain functional:**
- `en.Parse(text, ref, option)` - Still works
- `en.ParseDate(text, ref, option)` - Still works
- `en.Casual.Parse()` - Still works
- Direct Chrono usage - Still works

**Migration is optional** - Existing code continues to work. The new builder pattern API is recommended for new code.

### Performance

- No significant performance impact from new API (builder pattern is zero-cost)
- Interface wrapping is optimized by compiler
- Pipeline architecture allows for optimization opportunities

### Internal Changes

- Refactored parser execution into reusable pipeline
- Separated public API (Result, Components) from internal implementation
- Improved code organization and modularity
- Added comprehensive test coverage for new features
- Enhanced documentation and examples

## [1.0.0] - 2024-11-02

Initial release based on chrono JavaScript library.

### Features

- Natural language date parsing for English
- Casual and strict parsing modes
- Support for relative dates (yesterday, tomorrow, last week, etc.)
- Support for absolute dates (March 15, 2024, 3/15/2024, etc.)
- Time expression parsing (3pm, 14:30, noon, midnight)
- Date range parsing (from Monday to Friday)
- Timezone support
- Component certainty tracking (certain vs implied)
- British English support (DMY date order)
- Comprehensive test suite

### Parsers

- Casual date parser (today, tomorrow, yesterday, etc.)
- Casual time parser (now, etc.)
- Weekday parser (Monday, Tuesday, etc.)
- Month name parser (March, Jan, etc.)
- Month name little endian parser (15 March 2024)
- Month name middle endian parser (March 15, 2024)
- Year-month-day parser (2024-03-15)
- Slash date parser (3/15/2024, 15/3/2024)
- Slash month parser (3/2024)
- Compact format parser (20240315)
- Time expression parser (3pm, 14:30, etc.)
- Time unit parsers (3 days ago, in 2 weeks, etc.)
- Relative date parser (last/next week, etc.)

### Refiners

- Merge weekday components
- Merge date and time
- Imply recent dates
- Imply timezone offsets
- Extract timezone offsets
- Forward dates only
- Unspecified day of month

### Supported Expressions

- Relative: yesterday, today, tomorrow, now, last/next Monday, 3 days ago, in 2 weeks
- Absolute: March 15 2024, 3/15/2024, 2024-03-15, 15th of March
- Time: 3pm, 14:30, 3:30:45 PM, noon, midnight
- Combined: tomorrow at 3pm, March 15 at 14:30
- Ranges: from Monday to Friday, between March 1 and March 15

---

## Version History

- **2.0.0** (2024-11-03) - Builder pattern API, public interfaces, comprehensive documentation
- **1.0.0** (2024-11-02) - Initial release

## Upgrade Guide

See [MIGRATION.md](MIGRATION.md) for detailed upgrade instructions from v1 to v2.

## Links

- [Repository](https://github.com/kljensen/kronos)
- [Documentation](README.md)
- [Examples](examples/)
- [Migration Guide](MIGRATION.md)
