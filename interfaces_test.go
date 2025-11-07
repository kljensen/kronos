package kronos

import (
	"regexp"
	"testing"
	"time"
)

// Mock parser for testing
type mockParser struct {
	patternFunc func(context *parsingContext) *regexp.Regexp
	extractFunc func(context *parsingContext, match []string) interface{}
}

func (m *mockParser) Pattern(context *parsingContext) *regexp.Regexp {
	if m.patternFunc != nil {
		return m.patternFunc(context)
	}
	return regexp.MustCompile(`test`)
}

func (m *mockParser) Extract(context *parsingContext, match []string) interface{} {
	if m.extractFunc != nil {
		return m.extractFunc(context, match)
	}
	return nil
}

// Mock refiner for testing
type mockRefiner struct {
	refineFunc func(context *parsingContext, results []*parsingResult) []*parsingResult
}

func (m *mockRefiner) Refine(context *parsingContext, results []*parsingResult) []*parsingResult {
	if m.refineFunc != nil {
		return m.refineFunc(context, results)
	}
	return results
}

func TestParser(t *testing.T) {
	t.Run("Parser interface with mock implementation", func(t *testing.T) {
		refDate := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
		ctx := newParsingContext("test tomorrow", refDate, nil)

		patternCalled := false
		extractCalled := false

		parser := &mockParser{
			patternFunc: func(context *parsingContext) *regexp.Regexp {
				patternCalled = true
				return regexp.MustCompile(`tomorrow`)
			},
			extractFunc: func(context *parsingContext, match []string) interface{} {
				extractCalled = true
				components := map[Component]int{
					ComponentDay: 3,
				}
				return components
			},
		}

		// Test Pattern method
		pattern := parser.Pattern(ctx)
		if !patternCalled {
			t.Errorf("Expected Pattern to be called")
		}

		if pattern == nil {
			t.Fatalf("Expected pattern to be non-nil")
		}

		// Test Extract method
		match := pattern.FindStringSubmatch(ctx.Text())
		if match == nil {
			t.Fatalf("Expected pattern to match text")
		}

		result := parser.Extract(ctx, match)
		if !extractCalled {
			t.Errorf("Expected Extract to be called")
		}

		if result == nil {
			t.Errorf("Expected extract result to be non-nil")
		}

		if components, ok := result.(map[Component]int); ok {
			if components[ComponentDay] != 3 {
				t.Errorf("Expected day 3, got %d", components[ComponentDay])
			}
		} else {
			t.Errorf("Expected result to be component map")
		}
	})

	t.Run("Parser can return different result types", func(t *testing.T) {
		ctx := newParsingContext("test", time.Now(), nil)

		// Test returning component map
		p1 := &mockParser{
			extractFunc: func(context *parsingContext, match []string) interface{} {
				return map[Component]int{ComponentYear: 2024}
			},
		}

		result := p1.Extract(ctx, []string{"test"})
		if _, ok := result.(map[Component]int); !ok {
			t.Errorf("Expected component map")
		}

		// Test returning ParsingComponents
		p2 := &mockParser{
			extractFunc: func(context *parsingContext, match []string) interface{} {
				return newParsingComponents(ctx.Reference(), nil)
			},
		}

		result = p2.Extract(ctx, []string{"test"})
		if _, ok := result.(*parsingComponents); !ok {
			t.Errorf("Expected ParsingComponents")
		}

		// Test returning ParsingResult
		p3 := &mockParser{
			extractFunc: func(context *parsingContext, match []string) interface{} {
				return newParsingResult(ctx.Reference(), 0, "test", nil, nil)
			},
		}

		result = p3.Extract(ctx, []string{"test"})
		if _, ok := result.(*parsingResult); !ok {
			t.Errorf("Expected ParsingResult")
		}

		// Test returning nil
		p4 := &mockParser{
			extractFunc: func(context *parsingContext, match []string) interface{} {
				return nil
			},
		}

		result = p4.Extract(ctx, []string{"test"})
		if result != nil {
			t.Errorf("Expected nil result")
		}
	})
}

func TestRefiner(t *testing.T) {
	t.Run("Refiner interface with mock implementation", func(t *testing.T) {
		refDate := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
		ctx := newParsingContext("test", refDate, nil)

		refineCalled := false

		refiner := &mockRefiner{
			refineFunc: func(context *parsingContext, results []*parsingResult) []*parsingResult {
				refineCalled = true
				// Example: filter out results
				filtered := make([]*parsingResult, 0)
				for _, result := range results {
					if result.Index() > 0 {
						filtered = append(filtered, result)
					}
				}
				return filtered
			},
		}

		// Create test results
		results := []*parsingResult{
			newParsingResult(ctx.Reference(), 0, "first", nil, nil),
			newParsingResult(ctx.Reference(), 5, "second", nil, nil),
		}

		// Test Refine method
		refined := refiner.Refine(ctx, results)

		if !refineCalled {
			t.Errorf("Expected Refine to be called")
		}

		if len(refined) != 1 {
			t.Errorf("Expected 1 refined result, got %d", len(refined))
		}

		if refined[0].Index() != 5 {
			t.Errorf("Expected refined result at index 5, got %d", refined[0].Index())
		}
	})

	t.Run("Refiner can modify results", func(t *testing.T) {
		ctx := newParsingContext("test", time.Now(), nil)

		// Refiner that adds tags
		refiner := &mockRefiner{
			refineFunc: func(context *parsingContext, results []*parsingResult) []*parsingResult {
				for _, result := range results {
					result.AddTag("refined")
				}
				return results
			},
		}

		result := newParsingResult(ctx.Reference(), 0, "test", nil, nil)
		results := []*parsingResult{result}

		refined := refiner.Refine(ctx, results)

		if !refined[0].Tags()["refined"] {
			t.Errorf("Expected refined tag to be added")
		}
	})

	t.Run("Refiner can merge or split results", func(t *testing.T) {
		ctx := newParsingContext("test", time.Now(), nil)

		// Refiner that merges adjacent results
		refiner := &mockRefiner{
			refineFunc: func(context *parsingContext, results []*parsingResult) []*parsingResult {
				if len(results) < 2 {
					return results
				}

				// Simple merge: create single result from first and last
				merged := newParsingResult(
					ctx.Reference(),
					results[0].Index(),
					results[len(results)-1].Text(),
					nil,
					nil,
				)

				return []*parsingResult{merged}
			},
		}

		results := []*parsingResult{
			newParsingResult(ctx.Reference(), 0, "first", nil, nil),
			newParsingResult(ctx.Reference(), 6, "second", nil, nil),
		}

		refined := refiner.Refine(ctx, results)

		if len(refined) != 1 {
			t.Errorf("Expected 1 merged result, got %d", len(refined))
		}

		if refined[0].Index() != 0 {
			t.Errorf("Expected merged result at index 0, got %d", refined[0].Index())
		}
	})
}

func TestConfiguration(t *testing.T) {
	t.Run("Configuration with parsers and refiners", func(t *testing.T) {
		parser1 := &mockParser{}
		parser2 := &mockParser{}
		refiner1 := &mockRefiner{}
		refiner2 := &mockRefiner{}

		config := Configuration{
			Parsers:  []Parser{parser1, parser2},
			Refiners: []Refiner{refiner1, refiner2},
		}

		if len(config.Parsers) != 2 {
			t.Errorf("Expected 2 parsers, got %d", len(config.Parsers))
		}

		if len(config.Refiners) != 2 {
			t.Errorf("Expected 2 refiners, got %d", len(config.Refiners))
		}
	})

	t.Run("Configuration can be empty", func(t *testing.T) {
		config := Configuration{
			Parsers:  []Parser{},
			Refiners: []Refiner{},
		}

		if len(config.Parsers) != 0 {
			t.Errorf("Expected 0 parsers, got %d", len(config.Parsers))
		}

		if len(config.Refiners) != 0 {
			t.Errorf("Expected 0 refiners, got %d", len(config.Refiners))
		}
	})

	t.Run("Configuration with only parsers", func(t *testing.T) {
		parser := &mockParser{}

		config := Configuration{
			Parsers:  []Parser{parser},
			Refiners: []Refiner{},
		}

		if len(config.Parsers) != 1 {
			t.Errorf("Expected 1 parser, got %d", len(config.Parsers))
		}

		if len(config.Refiners) != 0 {
			t.Errorf("Expected 0 refiners, got %d", len(config.Refiners))
		}
	})

	t.Run("Configuration with only refiners", func(t *testing.T) {
		refiner := &mockRefiner{}

		config := Configuration{
			Parsers:  []Parser{},
			Refiners: []Refiner{refiner},
		}

		if len(config.Parsers) != 0 {
			t.Errorf("Expected 0 parsers, got %d", len(config.Parsers))
		}

		if len(config.Refiners) != 1 {
			t.Errorf("Expected 1 refiner, got %d", len(config.Refiners))
		}
	})
}
