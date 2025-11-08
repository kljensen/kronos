# Parser Package - Clean Interface for Language Support

## Overview

The `parser` package provides clean, public interfaces for implementing custom parsers and refiners in Kronos. This design follows Go best practices by separating interface definitions from implementation details.

## Design Pattern

This follows the `database/sql` and `database/sql/driver` pattern from the Go standard library:
- **`github.com/kljensen/kronos/parser`** - Defines interfaces (like `database/sql/driver`)
- **`github.com/kljensen/kronos`** - Main package that uses these interfaces (like `database/sql`)
- **Language packages** (e.g., `github.com/kljensen/kronos/en`) - Implement parsers for specific languages

## Why This Design?

Before this refactoring, the Parser and Refiner interfaces were in the main `kronos` package but marked as "INTERNAL USE". This was inelegant because:

1. **Circular Import Issues**: Supporting new languages would require importing kronos, but kronos needed to import the language packages
2. **Unclear API Boundary**: Exported interfaces marked "internal" confused the API surface
3. **Not Extensible**: External users couldn't easily add new language support

## Solution

Create a dedicated `parser` package with:
- `Parser` interface - for recognizing and parsing date patterns
- `Refiner` interface - for post-processing results
- `Context` interface - provides parsing context to parsers
- `Result` interface - represents parsing results
- `Reference` interface - reference date/time information
- `Option` interface - parsing options and preferences

## Usage

### For Users (Parsing Dates)

Most users never interact with the parser package directly. They use the language packages:

```go
import "github.com/kljensen/kronos/en"

results, err := en.ParseSimple("tomorrow at 3pm")
```

### For Language Implementers (Adding New Languages)

To add support for a new language (e.g., French):

```go
package fr

import (
	"regexp"
	"github.com/kljensen/kronos/parser"
)

type FrenchDateParser struct{}

func (p *FrenchDateParser) Pattern(ctx parser.Context) *regexp.Regexp {
	return regexp.MustCompile(`demain`)  // "tomorrow" in French
}

func (p *FrenchDateParser) Extract(ctx parser.Context, match []string) any {
	// Create and return components
	components := ctx.CreateParsingComponents(nil)
	// ... populate components ...
	return components
}
```

## Current Status

The parser package interfaces are defined and ready for use. The existing English language parsers continue to use internal interfaces for now, but new language implementations should use the parser package interfaces.

## Future Work

- Gradually migrate internal parsers to use the public interfaces
- Add example implementations for other languages
- Document parser and refiner patterns and best practices
