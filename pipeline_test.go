package kronos

import (
	"testing"
	"time"
)

func TestNewPipeline(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	pipeline := NewPipeline(config, settings)

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}
}

func TestNewPipeline_NilConfig(t *testing.T) {
	settings := DefaultSettings()
	pipeline := NewPipeline(nil, settings)

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}
}

func TestNewPipelineWithSettings(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	pipeline, err := NewPipelineWithSettings(config, settings)
	if err != nil {
		t.Fatalf("NewPipelineWithSettings failed: %v", err)
	}

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}
}

func TestNewPipelineWithSettings_InvalidSettings(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	_, err := NewPipelineWithSettings(config, settings)
	if err == nil {
		t.Error("Expected error for invalid settings")
	}
}

func TestNewPipelineWithSettings_EnabledParsers(t *testing.T) {
	// Register a test parser
	Register("test_parser", ParserInfo{
		Description: "Test parser",
		Priority:    50,
	}, func() Parser {
		return nil // Simplified for this test
	})

	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.EnabledParsers = []string{"test_parser"}

	pipeline, err := NewPipelineWithSettings(config, settings)
	if err != nil {
		t.Fatalf("NewPipelineWithSettings failed: %v", err)
	}

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}
}

func TestNewPipelineWithSettings_MaxParsers(t *testing.T) {
	// Create a config with multiple parsers
	config := &Configuration{
		Parsers:  []Parser{nil, nil, nil, nil, nil}, // 5 parsers
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.MaxParsers = 3

	pipeline, err := NewPipelineWithSettings(config, settings)
	if err != nil {
		t.Fatalf("NewPipelineWithSettings failed: %v", err)
	}

	if len(pipeline.parsers) != 3 {
		t.Errorf("Expected 3 parsers after applying MaxParsers, got %d", len(pipeline.parsers))
	}
}

func TestPipeline_Execute_BasicParsing(t *testing.T) {
	// This is a simplified test - full integration tests would use real parsers
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	pipeline := NewPipeline(config, settings)

	text := "March 15, 2020"
	refDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	results, err := pipeline.Execute(text, refDate)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// With no parsers, we expect no results
	if len(results) != 0 {
		t.Errorf("Expected 0 results with no parsers, got %d", len(results))
	}
}

func TestPipeline_Execute_StrictParsing(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.StrictParsing = true

	pipeline := NewPipeline(config, settings)

	text := "March 15"
	refDate := time.Now()

	results, err := pipeline.Execute(text, refDate)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Strict parsing should filter results
	_ = results
}

func TestPipeline_Execute_RequiredParts(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.RequireParts = []string{"year", "month", "day"}

	pipeline := NewPipeline(config, settings)

	text := "March 15, 2020"
	refDate := time.Now()

	results, err := pipeline.Execute(text, refDate)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// With no parsers, we expect no results
	if len(results) != 0 {
		t.Errorf("Expected 0 results with no parsers, got %d", len(results))
	}
}

func TestPipeline_Execute_TimezoneConversion(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"

	pipeline := NewPipeline(config, settings)

	text := "March 15, 2020"
	refDate := time.Now()

	_, err := pipeline.Execute(text, refDate)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestPipeline_Execute_Timeout(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.Timeout = 100 * time.Millisecond

	pipeline := NewPipeline(config, settings)

	text := "March 15, 2020"
	refDate := time.Now()

	_, err := pipeline.Execute(text, refDate)
	// Timeout shouldn't trigger with no parsers
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestParseWithSettings_ConvenienceFunction(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	text := "March 15, 2020"
	refDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	results, err := ParseWithSettings(text, refDate, settings, config)
	if err != nil {
		t.Fatalf("ParseWithSettings failed: %v", err)
	}

	// With no parsers, we expect no results
	if len(results) != 0 {
		t.Errorf("Expected 0 results with no parsers, got %d", len(results))
	}
}

func TestParseWithSettings_InvalidSettings(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	text := "March 15, 2020"
	refDate := time.Now()

	_, err := ParseWithSettings(text, refDate, settings, config)
	if err == nil {
		t.Error("Expected error for invalid settings")
	}
}

func TestPipeline_ApplyStrictValidation(t *testing.T) {
	settings := DefaultSettings()
	settings.StrictParsing = true
	pipeline := NewPipeline(nil, settings)

	// Create test results
	ref := &ReferenceWithTimezone{}

	// Result with year and month - should pass strict validation
	components1 := newParsingComponents(ref, map[Component]int{
		ComponentYear:  2020,
		ComponentMonth: 3,
	})
	components1.Assign(ComponentYear, 2020)
	components1.Assign(ComponentMonth, 3)
	result1 := &ParsingResult{start: components1}

	// Result with only month - should fail strict validation
	components2 := newParsingComponents(ref, map[Component]int{
		ComponentMonth: 3,
	})
	components2.Assign(ComponentMonth, 3)
	result2 := &ParsingResult{start: components2}

	results := []*ParsingResult{result1, result2}
	filtered := pipeline.applyStrictValidation(results)

	// Only result1 should pass
	if len(filtered) != 1 {
		t.Errorf("Expected 1 result after strict validation, got %d", len(filtered))
	}
}

func TestPipeline_ApplyRequiredParts(t *testing.T) {
	settings := DefaultSettings()
	settings.RequireParts = []string{"year", "month"}
	pipeline := NewPipeline(nil, settings)

	ref := &ReferenceWithTimezone{}

	// Result with year and month - should pass
	components1 := newParsingComponents(ref, map[Component]int{
		ComponentYear:  2020,
		ComponentMonth: 3,
	})
	components1.Assign(ComponentYear, 2020)
	components1.Assign(ComponentMonth, 3)
	result1 := &ParsingResult{start: components1}

	// Result with only month - should fail
	components2 := newParsingComponents(ref, map[Component]int{
		ComponentMonth: 3,
	})
	components2.Assign(ComponentMonth, 3)
	result2 := &ParsingResult{start: components2}

	results := []*ParsingResult{result1, result2}
	filtered := pipeline.applyRequiredParts(results)

	// Only result1 should pass
	if len(filtered) != 1 {
		t.Errorf("Expected 1 result after required parts filtering, got %d", len(filtered))
	}
}

func TestPipeline_ApplyRequiredParts_NoRequirements(t *testing.T) {
	settings := DefaultSettings()
	pipeline := NewPipeline(nil, settings)

	ref := &ReferenceWithTimezone{}
	components := newParsingComponents(ref, nil)
	result := &ParsingResult{start: components}
	results := []*ParsingResult{result}

	filtered := pipeline.applyRequiredParts(results)

	// With no requirements, all results should pass
	if len(filtered) != 1 {
		t.Errorf("Expected 1 result with no requirements, got %d", len(filtered))
	}
}

func TestPipeline_ApplyTimezoneConversion(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "America/New_York"
	pipeline := NewPipeline(nil, settings)

	ref := &ReferenceWithTimezone{}
	components := newParsingComponents(ref, nil)
	result := &ParsingResult{start: components}
	results := []*ParsingResult{result}

	converted, err := pipeline.applyTimezoneConversion(results)
	if err != nil {
		t.Fatalf("applyTimezoneConversion failed: %v", err)
	}

	if len(converted) != 1 {
		t.Errorf("Expected 1 result after timezone conversion, got %d", len(converted))
	}
}

func TestPipeline_ApplyTimezoneConversion_InvalidTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "Invalid/Timezone"
	pipeline := NewPipeline(nil, settings)

	ref := &ReferenceWithTimezone{}
	components := newParsingComponents(ref, nil)
	result := &ParsingResult{start: components}
	results := []*ParsingResult{result}

	_, err := pipeline.applyTimezoneConversion(results)
	if err == nil {
		t.Error("Expected error for invalid timezone")
	}
}
