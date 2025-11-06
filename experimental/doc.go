// Package experimental contains advanced and potentially unstable features
// for the kronos date parsing library.
//
// # Purpose
//
// This package provides the "escape hatch" for advanced users who need access
// to internal parsing APIs. It re-exports types and functions from internal
// packages that are useful for:
//
//   - Writing custom date/time parsers
//   - Extending kronos with new parsing patterns
//   - Building domain-specific date parsing solutions
//   - Migrating from deprecated public APIs
//
// # Stability Warning
//
// IMPORTANT: The experimental package APIs may change between minor versions.
// While we will make efforts to maintain compatibility, breaking changes are
// possible as the library evolves. Only use this package if you need advanced
// customization that cannot be achieved through the main kronos package.
//
// # When to Use This Package
//
// Use this package if you need to:
//
//   - Create custom parsers that integrate with the kronos parsing system
//   - Access parsing internals like ParsingContext, ParsingComponents, or ParsingResult
//   - Use helper functions for date construction (Today, Tomorrow, etc.)
//   - Perform complex date math with Duration and Timeunit
//   - Configure advanced parsing behavior with Settings
//
// # For Most Users
//
// Most users should use the simpler functions in the main kronos package:
//
//   - kronos.Parse() - Parse natural language dates
//   - kronos.ParseWithOption() - Parse with custom options
//   - kronos.ParseMultiple() - Parse multiple dates from text
//
// Only use this experimental package if you need to write custom parsers or
// extend kronos functionality beyond what the main package provides.
//
// # Migration Path
//
// This package provides a migration path for users who were using deprecated
// public APIs (like the old ParsingContext, ParsingComponents, etc.). These
// types have been moved to internal packages, and this experimental package
// re-exports them with the same names for backward compatibility during the
// transition period.
//
// # Available APIs
//
// The experimental package re-exports:
//
//   - Type aliases: ParsingContext, ParsingComponents, ParsingResult, ReferenceWithTimezone
//   - Factory functions: NewParsingContext, NewParsingComponents, NewParsingResult
//   - Date constructors: Today, Tomorrow, Yesterday, Now, Midnight, Noon, Morning, Afternoon, Evening
//   - Component helpers: AssignSimilarDate, AssignSimilarTime, MergeDateTimeComponent
//   - Date math: GetLastWeekday, GetNextWeekday, GetThisWeekday, AddDuration, ReverseDuration
//   - Year helpers: FindMostLikelyADYear, FindYearClosestToRefWithPreference
//   - Data constants: DefaultTimezoneAbbrMap, ApproximationWords, EmptyDuration
//   - Configuration types: Settings, DateOrder, DayPreference, DatePreference
//   - Duration types: Duration, Timeunit, Period
//   - Advanced parsing: Chrono, Configuration, Parser, Refiner, Pipeline
//   - Registry: ParserRegistry, ParserInfo, ParserFactory, GlobalRegistry, Register
//   - Pipeline functions: NewPipeline, NewPipelineWithSettings, ParseWithSettings
//
// See the package documentation and type definitions for detailed information
// about each exported item.
package experimental
