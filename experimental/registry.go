//nolint:staticcheck // SA1019: Must use deprecated types during transition
package experimental

import (
	"github.com/kljensen/kronos"
)

// ParserInfo contains metadata about a registered parser.
// Re-exported from the main package for the experimental API.
//
// This type provides metadata about parsers including their name, description,
// priority, and tags. It's used by the ParserRegistry to manage parser discovery
// and configuration.
type ParserInfo = kronos.ParserInfo

// ParserFactory creates a parser instance.
// Re-exported from the main package for the experimental API.
//
// This function type allows parsers to be created with specific settings.
// It's used by the ParserRegistry to instantiate parsers on demand.
type ParserFactory = kronos.ParserFactory

// ParserRegistry manages available parsers and their metadata.
// Re-exported from the main package for the experimental API.
//
// The ParserRegistry provides a central place to register and discover parsers.
// It supports querying parsers by name, tag, or priority, and maintains a
// default ordering for parser execution.
type ParserRegistry = kronos.ParserRegistry

// GlobalRegistry is the global parser registry for all parsers.
// Re-exported from the main package for the experimental API.
//
// This is the default registry used by the library. Custom parsers can be
// registered here using the Register function to make them available to
// all parsing operations.
var GlobalRegistry = kronos.GlobalRegistry

// NewParserRegistry creates a new parser registry.
// Re-exported from the main package for the experimental API.
//
// This function creates an independent parser registry that can be used
// to manage a separate set of parsers from the global registry.
var NewParserRegistry = kronos.NewParserRegistry

// Register registers a parser with the global registry.
// Re-exported from the main package for the experimental API.
//
// Use this function to add custom parsers to the global registry, making
// them available to all parsing operations in the application.
//
// Example:
//
//	Register("custom-date", ParserInfo{
//	    Name: "custom-date",
//	    Description: "Parses custom date formats",
//	    Priority: 100,
//	    Tags: []string{"custom"},
//	}, func() Parser { return &MyCustomParser{} })
var Register = kronos.Register
