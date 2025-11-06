// Package parsing contains internal parsing implementation types for kronos.
// These types were previously exposed in the main kronos package but have been moved
// to internal to reduce the public API surface.
//
// The parsing package provides:
//   - Context - holds parsing context (text, options, reference)
//   - Components - represents parsed date/time components with certainty levels
//   - Result - represents a parsed result with date/time information
//   - ReferenceWithTimezone - represents a reference date/time with timezone
//
// NOTE: This package cannot import the main kronos package to avoid import cycles.
// Therefore, it uses mirrored types where necessary. The main package provides
// type aliases and wrapper functions to maintain API compatibility.
//
// External code should not use types from this package directly. Instead, use:
//   - The ParsedComponents and ParsedResult interfaces from kronos
//   - The ParsingComponents and ParsingResult type aliases from kronos (deprecated)
package parsing
