package kronos

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockParser is a simple parser for testing
type MockParser struct {
	pattern     string
	extractFunc func(context *parsingContext, match []string) any
}

func (p *MockParser) Pattern(context *parsingContext) *regexp.Regexp {
	return regexp.MustCompile(p.pattern)
}

func (p *MockParser) Extract(context *parsingContext, match []string) any {
	if p.extractFunc != nil {
		return p.extractFunc(context, match)
	}
	return nil
}

// MockRefiner is a simple refiner for testing
type MockRefiner struct {
	refineFunc func(context *parsingContext, results []*parsingResult) []*parsingResult
}

func (r *MockRefiner) Refine(context *parsingContext, results []*parsingResult) []*parsingResult {
	if r.refineFunc != nil {
		return r.refineFunc(context, results)
	}
	return results
}

func TestNewChrono(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		c := NewChrono(nil)
		assert.NotNil(t, c)
		assert.Empty(t, c.parsers)
		assert.Empty(t, c.refiners)
	})

	t.Run("with config", func(t *testing.T) {
		parser := &MockParser{pattern: `test`}
		refiner := &MockRefiner{}
		config := &Configuration{
			Parsers:  []Parser{parser},
			Refiners: []Refiner{refiner},
		}

		c := NewChrono(config)
		assert.NotNil(t, c)
		assert.Len(t, c.parsers, 1)
		assert.Len(t, c.refiners, 1)
	})
}

// TestChronoClone was removed as Clone() method has been removed from the API.
// Chrono objects are deprecated in favor of the ParserBuilder API.

func TestChronoParse(t *testing.T) {
	t.Run("basic parsing", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\b(\d{4})-(\d{2})-(\d{2})\b`,
			extractFunc: func(context *parsingContext, match []string) any {
				return map[Component]int{
					ComponentYear:  2023,
					ComponentMonth: 10,
					ComponentDay:   15,
				}
			},
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("Today is 2023-10-15", nil, nil)
		assert.Len(t, results, 1)
		assert.Equal(t, 2023, *results[0].Start().Get(ComponentYear))
		assert.Equal(t, 10, *results[0].Start().Get(ComponentMonth))
		assert.Equal(t, 15, *results[0].Start().Get(ComponentDay))
	})

	t.Run("multiple matches", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\b(\d{4})\b`,
			extractFunc: func(context *parsingContext, match []string) any {
				year := 0
				if len(match) > 1 {
					_, _ = fmt.Sscanf(match[1], "%d", &year)
				}
				return map[Component]int{
					ComponentYear: year,
				}
			},
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("Years 2020 and 2021 were interesting", nil, nil)
		assert.Len(t, results, 2)
	})

	t.Run("sorted by index", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\b(\d{4})\b`,
			extractFunc: func(context *parsingContext, match []string) any {
				return map[Component]int{
					ComponentYear: 2020,
				}
			},
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("2021 was after 2020", nil, nil)
		assert.Len(t, results, 2)
		// Results should be sorted by position
		assert.True(t, results[0].Index() < results[1].Index())
	})
}

// TestChronoParseDate was removed as ParseDate() method has been removed from the API.
// Use ParserBuilder.ParseDate() instead, which is the recommended builder-based API.

func TestChronoRefiner(t *testing.T) {
	t.Run("refiner processes results", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\b(\d{4})\b`,
			extractFunc: func(context *parsingContext, match []string) any {
				return map[Component]int{ComponentYear: 2020}
			},
		}

		refinerCalled := false
		refiner := &MockRefiner{
			refineFunc: func(context *parsingContext, results []*parsingResult) []*parsingResult {
				refinerCalled = true
				// Filter out some results
				return results[:1]
			},
		}

		config := &Configuration{
			Parsers:  []Parser{parser},
			Refiners: []Refiner{refiner},
		}
		c := NewChrono(config)

		results := c.Parse("2020 2021", nil, nil)
		assert.True(t, refinerCalled)
		assert.Len(t, results, 1)
	})
}

func TestExecuteParser(t *testing.T) {
	t.Run("handles ParsingResult return", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\btest\b`,
			extractFunc: func(context *parsingContext, match []string) any {
				result := context.CreateParsingResult(0, "test")
				return result
			},
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("test", nil, nil)
		assert.Len(t, results, 1)
	})

	t.Run("handles ParsingComponents return", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\btest\b`,
			extractFunc: func(context *parsingContext, match []string) any {
				return context.CreateParsingComponents(map[Component]int{
					ComponentYear: 2023,
				})
			},
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("test", nil, nil)
		assert.Len(t, results, 1)
		assert.NotNil(t, results[0].Start())
	})

	t.Run("handles map return", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\btest\b`,
			extractFunc: func(context *parsingContext, match []string) any {
				return map[Component]int{
					ComponentYear: 2023,
				}
			},
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("test", nil, nil)
		assert.Len(t, results, 1)
		assert.NotNil(t, results[0].Start())
	})

	t.Run("handles nil return", func(t *testing.T) {
		parser := &MockParser{
			pattern:     `\btest\b`,
			extractFunc: func(context *parsingContext, match []string) any { return nil },
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("test", nil, nil)
		assert.Len(t, results, 0)
	})
}
