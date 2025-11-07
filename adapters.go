package kronos

import "time"

// resultAdapter adapts an internal ParsingResult to implement the public Result interface.
// This allows us to expose a clean public API while keeping implementation details internal.
type resultAdapter struct {
	result *parsingResult
}

// newResultAdapter creates a new resultAdapter wrapping a ParsingResult.
func newResultAdapter(result *parsingResult) *resultAdapter {
	return &resultAdapter{result: result}
}

// Text returns the matched text from the input.
func (r *resultAdapter) Text() string {
	return r.result.Text()
}

// Index returns the position in the input text where this result was found.
func (r *resultAdapter) Index() int {
	return r.result.Index()
}

// Date returns a time.Time object created from the start components.
func (r *resultAdapter) Date() time.Time {
	return r.result.Date()
}

// Start returns the starting date/time components.
func (r *resultAdapter) Start() Components {
	parsed := r.result.Start()
	// The ParsedComponents interface is implemented by ParsingComponents,
	// so we can wrap it in a componentsAdapter
	if pc, ok := parsed.(*parsingComponents); ok {
		return newComponentsAdapter(pc)
	}
	return nil
}

// End returns the ending date/time components for a range, or nil for a single date/time.
func (r *resultAdapter) End() Components {
	parsed := r.result.End()
	if parsed == nil {
		return nil
	}
	// The ParsedComponents interface is implemented by ParsingComponents,
	// so we can wrap it in a componentsAdapter
	if pc, ok := parsed.(*parsingComponents); ok {
		return newComponentsAdapter(pc)
	}
	return nil
}

// Tags returns metadata tags for this result.
// This is exposed for testing purposes to verify parser behavior.
func (r *resultAdapter) Tags() map[string]bool {
	return r.result.Tags()
}

// componentsAdapter adapts an internal ParsingComponents to implement the public Components interface.
type componentsAdapter struct {
	components *parsingComponents
}

// newComponentsAdapter creates a new componentsAdapter wrapping ParsingComponents.
func newComponentsAdapter(components *parsingComponents) *componentsAdapter {
	return &componentsAdapter{components: components}
}

// Get returns the component value.
func (c *componentsAdapter) Get(component Component) *int {
	return c.components.Get(component)
}

// IsCertain returns true if the component was explicitly mentioned in the input.
func (c *componentsAdapter) IsCertain(component Component) bool {
	return c.components.IsCertain(component)
}

// Date returns a time.Time object constructed from the components.
func (c *componentsAdapter) Date() time.Time {
	return c.components.Date()
}

// Tags returns metadata tags for these components.
// This is exposed for testing purposes to verify parser behavior.
func (c *componentsAdapter) Tags() map[string]bool {
	return c.components.Tags()
}
