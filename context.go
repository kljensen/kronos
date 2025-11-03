package kronos

import "time"

// ParsingContext holds the context for parsing operations.
// It contains the text to parse, options, and reference information.
type ParsingContext struct {
	text      string
	option    ParsingOption
	reference *ReferenceWithTimezone
	refDate   time.Time
}

// NewParsingContext creates a new ParsingContext.
// If refDate is nil, the current time is used.
// If option is nil, default options are used.
func NewParsingContext(text string, refDate interface{}, option *ParsingOption) *ParsingContext {
	var opt ParsingOption
	if option != nil {
		opt = *option
	}

	var timezones TimezoneAbbrMap
	if opt.Timezones != nil {
		timezones = opt.Timezones
	}

	reference := FromInput(refDate, timezones)

	return &ParsingContext{
		text:      text,
		option:    opt,
		reference: reference,
		refDate:   reference.Instant(),
	}
}

// CreateParsingComponents creates ParsingComponents from a component map or existing components.
// If components is already a ParsingComponents, it returns it as-is.
// Otherwise, it creates new ParsingComponents with the provided values.
func (ctx *ParsingContext) CreateParsingComponents(components interface{}) *ParsingComponents {
	if components == nil {
		return NewParsingComponents(ctx.reference, nil)
	}

	switch v := components.(type) {
	case *ParsingComponents:
		return v
	case map[Component]int:
		return NewParsingComponents(ctx.reference, v)
	default:
		return NewParsingComponents(ctx.reference, nil)
	}
}

// CreateParsingResult creates a ParsingResult with various signatures.
// It can accept:
// - (index, text) - creates result with text
// - (index, endIndex) - creates result with substring from text
// - (index, text, startComponents) - creates result with start components
// - (index, text, startComponents, endComponents) - creates result with both start and end
func (ctx *ParsingContext) CreateParsingResult(index int, textOrEndIndex interface{}, args ...interface{}) *ParsingResult {
	var text string
	var start *ParsingComponents
	var end *ParsingComponents

	// Determine if second arg is text or endIndex
	switch v := textOrEndIndex.(type) {
	case string:
		text = v
	case int:
		// It's an endIndex, extract substring
		endIndex := v
		if index < 0 {
			index = 0
		}
		if endIndex > len(ctx.text) {
			endIndex = len(ctx.text)
		}
		if index > len(ctx.text) {
			index = len(ctx.text)
		}
		text = ctx.text[index:endIndex]
	default:
		text = ""
	}

	// Process remaining arguments
	for _, arg := range args {
		if arg == nil {
			continue
		}

		switch v := arg.(type) {
		case *ParsingComponents:
			if start == nil {
				start = v
			} else {
				end = v
			}
		case map[Component]int:
			if start == nil {
				start = ctx.CreateParsingComponents(v)
			} else {
				end = ctx.CreateParsingComponents(v)
			}
		}
	}

	return NewParsingResult(ctx.reference, index, text, start, end)
}

// Debug executes the provided function if debugging is enabled.
func (ctx *ParsingContext) Debug(fn func()) {
	if ctx.option.Debug != nil {
		fn()
	}
}

// Text returns the input text.
func (ctx *ParsingContext) Text() string {
	return ctx.text
}

// Option returns the parsing options.
func (ctx *ParsingContext) Option() ParsingOption {
	return ctx.option
}

// Reference returns the reference with timezone.
func (ctx *ParsingContext) Reference() *ReferenceWithTimezone {
	return ctx.reference
}

// RefDate returns the reference date.
func (ctx *ParsingContext) RefDate() time.Time {
	return ctx.refDate
}
