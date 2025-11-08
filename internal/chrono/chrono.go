// Package chrono provides the internal Advanced API structures for Kronos.
// This package contains Configuration and Chrono types for managing parsing pipelines.
//
// INTERNAL USE ONLY: This package is not part of the public API and may change without notice.
// Most users should use the builder pattern (en.New()) instead. This package is used by
// internal implementation and language packages (like en).
package chrono

// Configuration holds the parsers and refiners for Chrono.
// It is simply an ordered list of parsers and refiners.
//
// The Parser and Refiner interfaces are defined in the main kronos package.
type Configuration struct {
	Parsers  []any // []kronos.Parser
	Refiners []any // []kronos.Refiner
}

// Chrono is the main parsing engine configuration that holds parsers and refiners.
// It maintains a list of parsers (each handling a specific date format) and refiners
// (each post-processing the results).
//
// The actual parsing logic is provided by the main kronos package.
type Chrono struct {
	parsers  []any // []kronos.Parser
	refiners []any // []kronos.Refiner
}

// NewChrono creates a new Chrono instance with the given configuration.
// If config is nil, an empty Chrono is created.
func NewChrono(config *Configuration) *Chrono {
	if config == nil {
		return &Chrono{
			parsers:  []any{},
			refiners: []any{},
		}
	}

	return &Chrono{
		parsers:  append([]any{}, config.Parsers...),
		refiners: append([]any{}, config.Refiners...),
	}
}

// Parsers returns a copy of the parsers list for inspection.
func (c *Chrono) Parsers() []any {
	return append([]any{}, c.parsers...)
}

// Refiners returns a copy of the refiners list for inspection.
func (c *Chrono) Refiners() []any {
	return append([]any{}, c.refiners...)
}
