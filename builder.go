package kronos

import (
	"time"
)

// ============================================================================
// ParserBuilder - Modern builder-based API
// ============================================================================

// ParserBuilder is the main entry point for configuring and executing date parsing.
// It provides a fluent API for setting up parsing options and executing the parse.
// Use New() to create a new parser builder, configure it with builder methods,
// then call Parse() or ParseDate() to execute the parsing.
//
// Thread Safety: ParserBuilder is not thread-safe and should not be shared across
// goroutines. Each goroutine should create its own ParserBuilder instance.
//
// Mutability: ParserBuilder is mutable. Calling configuration methods (WithReferenceDate,
// Strict, PreferPast, etc.) modifies the builder's state and returns the same builder
// instance to enable method chaining. If you need different configurations, create
// separate ParserBuilder instances.
//
// Example:
//
//	parser := kronos.New(en.Casual).
//	    WithReferenceDate(time.Now()).
//	    Strict().
//	    PreferPast()
//	result := parser.Parse("last Monday")
type ParserBuilder struct {
	config   *Configuration
	settings Settings
	refDate  time.Time
}

// New creates a new ParserBuilder with the given Chrono configuration.
// This is the main entry point for the builder pattern API.
//
// Example:
//
//	parser := kronos.New(en.Casual)
func New(chrono *Chrono) *ParserBuilder {
	if chrono == nil {
		chrono = NewChrono(nil)
	}

	// Extract configuration from Chrono
	config := &Configuration{
		Parsers:  append([]Parser{}, chrono.parsers...),
		Refiners: append([]Refiner{}, chrono.refiners...),
	}

	return &ParserBuilder{
		config:   config,
		settings: DefaultSettings(),
		refDate:  time.Now(),
	}
}

// WithReferenceDate sets the reference date for relative date calculations.
// This date is used as the basis for expressions like "tomorrow", "next week", etc.
func (p *ParserBuilder) WithReferenceDate(refDate time.Time) *ParserBuilder {
	p.refDate = refDate
	return p
}

// Strict enables strict parsing mode.
// In strict mode, the parser validates more strictly and rejects ambiguous dates.
func (p *ParserBuilder) Strict() *ParserBuilder {
	p.settings.StrictParsing = true
	return p
}

// Casual enables casual parsing mode (this is the default).
// In casual mode, the parser accepts informal expressions and is more lenient.
func (p *ParserBuilder) Casual() *ParserBuilder {
	p.settings.StrictParsing = false
	return p
}

// DateOrder sets the order of date components in ambiguous formats.
// Use DateOrderMDY for US format (12/31/2020), DateOrderDMY for European format (31/12/2020),
// or DateOrderYMD for ISO format (2020/12/31).
func (p *ParserBuilder) DateOrder(order DateOrder) *ParserBuilder {
	p.settings.DateOrder = order
	return p
}

// PreferPast configures the parser to prefer past dates when ambiguous.
// For example, "March" in November would be interpreted as last March.
func (p *ParserBuilder) PreferPast() *ParserBuilder {
	p.settings.PreferDatesFrom = PreferPast
	return p
}

// PreferFuture configures the parser to prefer future dates when ambiguous.
// For example, "March" in November would be interpreted as next March.
func (p *ParserBuilder) PreferFuture() *ParserBuilder {
	p.settings.PreferDatesFrom = PreferFuture
	return p
}

// PreferCurrentPeriod configures the parser to prefer dates in the current period.
// This is the default behavior.
func (p *ParserBuilder) PreferCurrentPeriod() *ParserBuilder {
	p.settings.PreferDatesFrom = PreferCurrentPeriod
	return p
}

// Timezone sets the default timezone for parsing.
// This timezone is used when no timezone is specified in the input.
func (p *ParserBuilder) Timezone(tz string) *ParserBuilder {
	p.settings.Timezone = tz
	return p
}

// WithOption applies an advanced configuration option to the parser.
// This allows fine-grained control over parsing behavior by directly modifying
// the Settings struct.
//
// Example:
//
//	parser := kronos.New(en.Casual).
//	    WithOption(func(s *kronos.Settings) {
//	        s.TimezoneOverrides = customTimezones
//	        s.DebugHandler = debugFunc
//	    })
func (p *ParserBuilder) WithOption(option func(*Settings)) *ParserBuilder {
	option(&p.settings)
	return p
}

// Parse executes the parser on the given text and returns all found date/time results.
// Returns an error if the settings are invalid or parsing fails.
func (p *ParserBuilder) Parse(text string) ([]Result, error) {
	// Create pipeline with settings
	pipeline, err := newPipelineWithSettings(p.config, p.settings)
	if err != nil {
		return nil, err
	}

	// Execute parsing
	results, err := pipeline.Execute(text, p.refDate)
	if err != nil {
		return nil, err
	}

	// Convert internal ParsingResult to public Result interface
	publicResults := make([]Result, len(results))
	for i, r := range results {
		publicResults[i] = newResultAdapter(r)
	}
	return publicResults, nil
}

// ParseDate is a shortcut for calling Parse and returning the first result's date.
// Returns nil if no results are found or an error occurs.
func (p *ParserBuilder) ParseDate(text string) (*time.Time, error) {
	results, err := p.Parse(text)
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		date := results[0].Date()
		return &date, nil
	}
	return nil, nil
}

// Settings returns a copy of the current settings for inspection.
// This is useful for advanced use cases that need to examine or validate settings.
// Note: Modifying the returned Settings will not affect the builder.
//
// Example:
//
//	parser := kronos.New(en.Casual).PreferFuture()
//	settings := parser.Settings()
//	fmt.Printf("Preference: %v\n", settings.PreferDatesFrom)
func (p *ParserBuilder) Settings() Settings {
	return p.settings
}
