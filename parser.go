package kronos

import (
	"regexp"

	"github.com/kljensen/kronos/parser"
)

// ============================================================================
// Public Parser/Refiner Interfaces (in parser package)
// ============================================================================
//
// For new implementations supporting additional languages, use the parser package:
//   - parser.Parser - interface for custom parsers
//   - parser.Refiner - interface for custom refiners
//   - parser.Context - parsing context interface
//   - parser.Result - parsing result interface
//
// See github.com/kljensen/kronos/parser for details and examples.
//
// The types below are internal implementation details.
// ============================================================================

// internalParser is the internal interface used by built-in parsers.
// External code should implement parser.Parser instead.
type internalParser interface {
	Pattern(context *parsingContext) *regexp.Regexp
	Extract(context *parsingContext, match []string) any
}

// internalRefiner is the internal interface used by built-in refiners.
// External code should implement parser.Refiner instead.
type internalRefiner interface {
	Refine(context *parsingContext, results []*parsingResult) []*parsingResult
}

// Parser is the public interface for implementing custom parsers.
// Use this when adding support for new languages or date formats.
//
// Example:
//
//	type MyParser struct{}
//	func (p *MyParser) Pattern(ctx parser.Context) *regexp.Regexp { ... }
//	func (p *MyParser) Extract(ctx parser.Context, match []string) any { ... }
type Parser = parser.Parser

// Refiner is the public interface for implementing custom refiners.
// Use this when adding custom post-processing logic.
//
// Example:
//
//	type MyRefiner struct{}
//	func (r *MyRefiner) Refine(ctx parser.Context, results []parser.Result) []parser.Result { ... }
type Refiner = parser.Refiner
