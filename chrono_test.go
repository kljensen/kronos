package kronos

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockParser is a simple parser for testing
type MockParser struct {
	pattern     string
	extractFunc func(context *ParsingContext, match []string) interface{}
}

func (p *MockParser) Pattern(context *ParsingContext) *regexp.Regexp {
	return regexp.MustCompile(p.pattern)
}

func (p *MockParser) Extract(context *ParsingContext, match []string) interface{} {
	if p.extractFunc != nil {
		return p.extractFunc(context, match)
	}
	return nil
}

// MockRefiner is a simple refiner for testing
type MockRefiner struct {
	refineFunc func(context *ParsingContext, results []*ParsingResult) []*ParsingResult
}

func (r *MockRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
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

func TestChronoClone(t *testing.T) {
	parser := &MockParser{pattern: `test`}
	refiner := &MockRefiner{}
	config := &Configuration{
		Parsers:  []Parser{parser},
		Refiners: []Refiner{refiner},
	}

	c := NewChrono(config)
	clone := c.Clone()

	assert.NotNil(t, clone)
	assert.Len(t, clone.parsers, 1)
	assert.Len(t, clone.refiners, 1)

	// Ensure it's a shallow copy - modifying the clone shouldn't affect original
	clone.parsers = append(clone.parsers, &MockParser{pattern: `another`})
	assert.Len(t, c.parsers, 1)
	assert.Len(t, clone.parsers, 2)
}

func TestChronoParse(t *testing.T) {
	t.Run("basic parsing", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\b(\d{4})-(\d{2})-(\d{2})\b`,
			extractFunc: func(context *ParsingContext, match []string) interface{} {
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
			extractFunc: func(context *ParsingContext, match []string) interface{} {
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
			extractFunc: func(context *ParsingContext, match []string) interface{} {
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

func TestChronoParseDate(t *testing.T) {
	t.Run("returns first date", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\b(\d{4})-(\d{2})-(\d{2})\b`,
			extractFunc: func(context *ParsingContext, match []string) interface{} {
				return map[Component]int{
					ComponentYear:  2023,
					ComponentMonth: 10,
					ComponentDay:   15,
				}
			},
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		date := c.ParseDate("Today is 2023-10-15", nil, nil)
		assert.NotNil(t, date)
		assert.Equal(t, 2023, date.Year())
		assert.Equal(t, time.Month(10), date.Month())
		assert.Equal(t, 15, date.Day())
	})

	t.Run("returns nil when no matches", func(t *testing.T) {
		parser := &MockParser{
			pattern:     `\b(\d{4})-(\d{2})-(\d{2})\b`,
			extractFunc: func(context *ParsingContext, match []string) interface{} { return nil },
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		date := c.ParseDate("No date here", nil, nil)
		assert.Nil(t, date)
	})
}

func TestChronoRefiner(t *testing.T) {
	t.Run("refiner processes results", func(t *testing.T) {
		parser := &MockParser{
			pattern: `\b(\d{4})\b`,
			extractFunc: func(context *ParsingContext, match []string) interface{} {
				return map[Component]int{ComponentYear: 2020}
			},
		}

		refinerCalled := false
		refiner := &MockRefiner{
			refineFunc: func(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
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
			extractFunc: func(context *ParsingContext, match []string) interface{} {
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
			extractFunc: func(context *ParsingContext, match []string) interface{} {
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
			extractFunc: func(context *ParsingContext, match []string) interface{} {
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
			extractFunc: func(context *ParsingContext, match []string) interface{} { return nil },
		}

		config := &Configuration{Parsers: []Parser{parser}}
		c := NewChrono(config)

		results := c.Parse("test", nil, nil)
		assert.Len(t, results, 0)
	})
}
