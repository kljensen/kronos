package kronos

import (
	"time"

	"github.com/kljensen/kronos/internal/sanitization"
)

// ============================================================================
// Internal parsing types
// ============================================================================

// parsingOption contains configuration options for parsing.
type parsingOption struct {
	// ForwardDate indicates whether to parse only forward dates
	// (results should be after the reference date).
	ForwardDate bool

	// Preference specifies how ambiguous dates should be resolved.
	Preference DatePreference

	// DateOrder specifies the order of date components in ambiguous formats.
	DateOrder DateOrder

	// Timezones provides additional timezone keywords for parsers to recognize.
	Timezones TimezoneAbbrMap

	// Debug is an internal debug event handler.
	Debug DebugHandler
}

// parsingReference contains reference information for parsing dates/times.
type parsingReference struct {
	// Instant is the reference date/time when the input is written or mentioned.
	Instant *time.Time

	// Timezone is the reference timezone where the input is written or mentioned.
	Timezone any
}

// ============================================================================
// Parsing context
// ============================================================================

// parsingContext holds the context for parsing operations.
// It contains the text to parse, options, and reference information.
// This is an internal implementation type.
type parsingContext struct {
	text      string
	option    parsingOption
	reference *referenceWithTimezone
	refDate   time.Time
	settings  *Settings
}

// newParsingContext creates a new ParsingContext.
// If refDate is nil, the current time is used.
// If option is nil, default options are used.
// The input text is sanitized to normalize Unicode characters before parsing.
// The option parameter can be *parsingOption or Settings (for backward compatibility).
func newParsingContext(text string, refDate any, option any) *parsingContext {
	var opt parsingOption
	var settings *Settings

	// Handle different types of option parameter
	switch v := option.(type) {
	case *parsingOption:
		if v != nil {
			opt = *v
		}
	case parsingOption:
		opt = v
	case Settings:
		opt = v.toparsingOption(nil)
		settings = &v
	case *Settings:
		if v != nil {
			opt = v.toparsingOption(nil)
			settings = v
		}
	case nil:
		// Use default options
	}

	var timezones TimezoneAbbrMap
	if opt.Timezones != nil {
		timezones = opt.Timezones
	}

	reference := fromInput(refDate, timezones)

	// Sanitize input text to handle Unicode normalization issues
	text = sanitization.SanitizeInput(text)

	return &parsingContext{
		text:      text,
		option:    opt,
		reference: reference,
		refDate:   reference.Instant(),
		settings:  settings,
	}
}

// CreateParsingComponents creates ParsingComponents from a component map or existing components.
// If components is already a ParsingComponents, it returns it as-is.
// Otherwise, it creates new ParsingComponents with the provided values.
// Returns any to satisfy parser.Context interface, but the actual type is *parsingComponents.
func (ctx *parsingContext) CreateParsingComponents(components any) any {
	if components == nil {
		return newParsingComponents(ctx.reference, nil)
	}

	switch v := components.(type) {
	case *parsingComponents:
		return v
	case map[Component]int:
		return newParsingComponents(ctx.reference, v)
	default:
		return newParsingComponents(ctx.reference, nil)
	}
}

// CreateParsingResult creates a ParsingResult with various signatures.
// It can accept:
// - (index, text) - creates result with text
// - (index, endIndex) - creates result with substring from text
// - (index, text, startComponents) - creates result with start components
// - (index, text, startComponents, endComponents) - creates result with both start and end
func (ctx *parsingContext) CreateParsingResult(index int, textOrEndIndex any, args ...any) *parsingResult {
	var text string
	var start *parsingComponents
	var end *parsingComponents

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
		case *parsingComponents:
			if start == nil {
				start = v
			} else {
				end = v
			}
		case map[Component]int:
			if start == nil {
				if components := ctx.CreateParsingComponents(v); components != nil {
					start = components.(*parsingComponents)
				}
			} else {
				if components := ctx.CreateParsingComponents(v); components != nil {
					end = components.(*parsingComponents)
				}
			}
		}
	}

	return newParsingResult(ctx.reference, index, text, start, end)
}

// Debug executes the provided function if debugging is enabled.
func (ctx *parsingContext) Debug(fn func()) {
	if ctx.option.Debug != nil {
		fn()
	}
}

// Text returns the input text.
func (ctx *parsingContext) Text() string {
	return ctx.text
}

// Option returns the parsing options.
func (ctx *parsingContext) Option() parsingOption {
	return ctx.option
}

// Reference returns the reference with timezone.
func (ctx *parsingContext) Reference() *referenceWithTimezone {
	return ctx.reference
}

// RefDate returns the reference date.
func (ctx *parsingContext) RefDate() time.Time {
	return ctx.refDate
}

// Settings returns the parsing settings, if any.
func (ctx *parsingContext) Settings() *Settings {
	return ctx.settings
}
