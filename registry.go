package kronos

import (
	"sort"
	"sync"
)

// ParserInfo contains metadata about a registered parser.
// It allows parsers to be dynamically discovered and configured.
//
// Deprecated: This type is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type ParserInfo struct {
	// Name is the unique identifier for this parser.
	Name string

	// Description explains what date/time formats this parser handles.
	Description string

	// Priority determines the order in which parsers are tried.
	// Higher priority parsers are tried first.
	// Default priority is 0.
	Priority int

	// Tags are optional labels for grouping parsers.
	// Examples: "casual", "formal", "iso", "relative"
	Tags []string
}

// ParserFactory creates a parser instance.
// This allows parsers to be created with specific settings.
//
// Deprecated: This type is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type ParserFactory func() Parser

// ParserRegistry manages available parsers and their metadata.
// It provides a central place to register and discover parsers.
//
// Deprecated: This type is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type ParserRegistry struct {
	mu       sync.RWMutex
	parsers  map[string]*registeredParser
	defaults []string // Default parser order
}

// registeredParser holds the parser factory and metadata.
type registeredParser struct {
	info    ParserInfo
	factory ParserFactory
}

// newParserRegistry creates a new parser registry.
// This is unexported and only used internally. External code should use
// experimental.NewParserRegistry() instead.
func newParserRegistry() *ParserRegistry {
	return &ParserRegistry{
		parsers:  make(map[string]*registeredParser),
		defaults: []string{},
	}
}

// globalRegistry is the internal global parser registry.
// External code should use experimental.GlobalRegistry instead.
var globalRegistry = newParserRegistry()

// internalRegister registers a parser with the global registry.
// This is unexported and only used by experimental package.
// External code should use experimental.Register() instead.
func internalRegister(name string, info ParserInfo, factory ParserFactory) {
	globalRegistry.RegisterParser(name, info, factory)
}

// RegisterParser registers a parser with this registry.
// If a parser with the same name already exists, it is replaced.
func (r *ParserRegistry) RegisterParser(name string, info ParserInfo, factory ParserFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	info.Name = name
	r.parsers[name] = &registeredParser{
		info:    info,
		factory: factory,
	}
}

// GetParser retrieves a parser by name and creates an instance.
// Returns nil if the parser is not found.
func (r *ParserRegistry) GetParser(name string) Parser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if p, ok := r.parsers[name]; ok {
		return p.factory()
	}
	return nil
}

// GetParsers retrieves multiple parsers by name.
// Unknown parser names are silently skipped.
func (r *ParserRegistry) GetParsers(names []string) []Parser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parsers := make([]Parser, 0, len(names))
	for _, name := range names {
		if p, ok := r.parsers[name]; ok {
			parsers = append(parsers, p.factory())
		}
	}
	return parsers
}

// GetAllParsers returns all registered parsers in priority order.
// Higher priority parsers come first.
func (r *ParserRegistry) GetAllParsers() []Parser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a list of all parsers
	list := make([]*registeredParser, 0, len(r.parsers))
	for _, p := range r.parsers {
		list = append(list, p)
	}

	// Sort by priority (descending)
	sort.Slice(list, func(i, j int) bool {
		return list[i].info.Priority > list[j].info.Priority
	})

	// Create parser instances
	parsers := make([]Parser, 0, len(list))
	for _, p := range list {
		parsers = append(parsers, p.factory())
	}

	return parsers
}

// GetParsersByTag returns all parsers with the specified tag.
func (r *ParserRegistry) GetParsersByTag(tag string) []Parser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parsers := make([]Parser, 0)
	for _, p := range r.parsers {
		for _, t := range p.info.Tags {
			if t == tag {
				parsers = append(parsers, p.factory())
				break
			}
		}
	}
	return parsers
}

// ListParsers returns information about all registered parsers.
func (r *ParserRegistry) ListParsers() []ParserInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]ParserInfo, 0, len(r.parsers))
	for _, p := range r.parsers {
		infos = append(infos, p.info)
	}

	// Sort by name for consistent output
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return infos
}

// HasParser checks if a parser with the given name is registered.
func (r *ParserRegistry) HasParser(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.parsers[name]
	return ok
}

// SetDefaultOrder sets the default parser order for this registry.
// This order is used when no custom order is specified.
func (r *ParserRegistry) SetDefaultOrder(names []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaults = append([]string{}, names...)
}

// GetDefaultOrder returns the default parser order.
func (r *ParserRegistry) GetDefaultOrder() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]string{}, r.defaults...)
}

// Clear removes all registered parsers.
// This is mainly useful for testing.
func (r *ParserRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers = make(map[string]*registeredParser)
	r.defaults = []string{}
}
