package experimental

import (
	"github.com/kljensen/kronos"
)

// ParserInfo contains metadata about a registered parser.
// Re-exported from the main package for the experimental API.
type ParserInfo = kronos.ParserInfo

// ParserFactory creates a parser instance.
// Re-exported from the main package for the experimental API.
type ParserFactory = kronos.ParserFactory

// ParserRegistry manages available parsers and their metadata.
// Re-exported from the main package for the experimental API.
type ParserRegistry = kronos.ParserRegistry

// GlobalRegistry is the global parser registry for all parsers.
// Re-exported from the main package for the experimental API.
var GlobalRegistry = kronos.GlobalRegistry

// NewParserRegistry creates a new parser registry.
// Re-exported from the main package for the experimental API.
var NewParserRegistry = kronos.NewParserRegistry

// Register registers a parser with the global registry.
// Re-exported from the main package for the experimental API.
var Register = kronos.Register
