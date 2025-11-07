package kronos

import (
	"regexp"
	"sort"
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

// ToTimezone sets the target timezone for converting results.
// All parsed dates will be converted to this timezone.
//
// EXPERIMENTAL: This feature is currently experimental and not fully implemented.
// The timezone conversion logic is stubbed out in the pipeline and does not yet
// modify the parsed results. Use with caution.
//
// TODO: Complete the timezone conversion implementation in pipeline.go:applyTimezoneConversion
func (p *ParserBuilder) ToTimezone(tz string) *ParserBuilder {
	p.settings.ToTimezone = tz
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
	// Use pipeline for parsing with settings
	results, err := parseWithSettings(text, p.refDate, p.settings, p.config)
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

// ============================================================================
// Deprecated Chrono API (kept for backward compatibility)
// ============================================================================

// Parser is an abstraction for Chrono parsers.
// Each parser should recognize and handle a certain date format.
// Chrono uses multiple parsers (and refiners) together for parsing the input.
//
// Deprecated: This interface is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Parser interface {
	// Pattern returns the regular expression pattern for this parser.
	// The pattern is used to find potential matches in the input text.
	Pattern(context *parsingContext) *regexp.Regexp

	// Extract is called with the pattern's match.
	// It should return parsed components, a result, a component map, or nil if extraction fails.
	Extract(context *parsingContext, match []string) interface{}
}

// Refiner is an abstraction for Chrono refiners.
// Each refiner takes the list of results (from parsers or other refiners)
// and returns another list of results.
// Chrono applies each refiner in order and returns the output from the last refiner.
//
// Deprecated: This interface is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Refiner interface {
	// Refine processes a list of parsing results and returns a refined list.
	Refine(context *parsingContext, results []*parsingResult) []*parsingResult
}

// Configuration holds the parsers and refiners for Chrono.
// It is simply an ordered list of parsers and refiners.
//
// Deprecated: This type is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Configuration struct {
	Parsers  []Parser
	Refiners []Refiner
}

// Chrono is the main parsing engine that coordinates multiple parsers and refiners.
// It maintains a list of parsers (each handling a specific date format) and refiners
// (each post-processing the results).
//
// Deprecated: This type is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
type Chrono struct {
	parsers  []Parser
	refiners []Refiner
}

// NewChrono creates a new Chrono instance with the given configuration.
// If config is nil, an empty Chrono is created.
//
// Deprecated: This function is part of the advanced API and will be moved to the
// experimental package in a future version. For new code, import and use
// github.com/kljensen/kronos/experimental instead.
func NewChrono(config *Configuration) *Chrono {
	if config == nil {
		return &Chrono{
			parsers:  []Parser{},
			refiners: []Refiner{},
		}
	}

	return &Chrono{
		parsers:  append([]Parser{}, config.Parsers...),
		refiners: append([]Refiner{}, config.Refiners...),
	}
}

// Parse parses the input text and returns all found date/time results.
// The parsing process:
// 1. Create a parsing context
// 2. Execute all parsers to find matches
// 3. Sort results by position in text
// 4. Apply all refiners sequentially
// 5. Return final results
func (c *Chrono) Parse(text string, referenceDate interface{}, option *parsingOption) []*parsingResult {
	context := newParsingContext(text, referenceDate, option)

	results := make([]*parsingResult, 0)
	for _, parser := range c.parsers {
		parsedResults := executeParser(context, parser)
		results = append(results, parsedResults...)
	}

	// Sort results by position in text
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	// Apply refiners
	for _, refiner := range c.refiners {
		results = refiner.Refine(context, results)
	}

	return results
}
