package kronos

import (
	"testing"
	"time"
)

func TestnewPipeline(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	pipeline := newPipeline(config, settings)

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}
}

func TestnewPipeline_NilConfig(t *testing.T) {
	settings := DefaultSettings()
	pipeline := newPipeline(nil, settings)

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}
}

func TestnewPipelineWithSettings(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	pipeline, err := newPipelineWithSettings(config, settings)
	if err != nil {
		t.Fatalf("newPipelineWithSettings failed: %v", err)
	}

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}
}

func TestnewPipelineWithSettings_InvalidSettings(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	_, err := newPipelineWithSettings(config, settings)
	if err == nil {
		t.Error("Expected error for invalid settings")
	}
}


func TestPipeline_Execute_BasicParsing(t *testing.T) {
	// This is a simplified test - full integration tests would use real parsers
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	pipeline := newPipeline(config, settings)

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

	pipeline := newPipeline(config, settings)

	text := "March 15"
	refDate := time.Now()

	results, err := pipeline.Execute(text, refDate)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Strict parsing should filter results
	_ = results
}




func TestparseWithSettings_ConvenienceFunction(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()

	text := "March 15, 2020"
	refDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	results, err := parseWithSettings(text, refDate, settings, config)
	if err != nil {
		t.Fatalf("parseWithSettings failed: %v", err)
	}

	// With no parsers, we expect no results
	if len(results) != 0 {
		t.Errorf("Expected 0 results with no parsers, got %d", len(results))
	}
}

func TestparseWithSettings_InvalidSettings(t *testing.T) {
	config := &Configuration{
		Parsers:  []Parser{},
		Refiners: []Refiner{},
	}
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	text := "March 15, 2020"
	refDate := time.Now()

	_, err := parseWithSettings(text, refDate, settings, config)
	if err == nil {
		t.Error("Expected error for invalid settings")
	}
}

func TestPipeline_ApplyStrictValidation(t *testing.T) {
	settings := DefaultSettings()
	settings.StrictParsing = true
	pipeline := newPipeline(nil, settings)

	// Create test results
	ref := &referenceWithTimezone{}

	// Result with year and month - should pass strict validation
	components1 := newParsingComponents(ref, map[Component]int{
		ComponentYear:  2020,
		ComponentMonth: 3,
	})
	components1.Assign(ComponentYear, 2020)
	components1.Assign(ComponentMonth, 3)
	result1 := &parsingResult{start: components1}

	// Result with only month - should fail strict validation
	components2 := newParsingComponents(ref, map[Component]int{
		ComponentMonth: 3,
	})
	components2.Assign(ComponentMonth, 3)
	result2 := &parsingResult{start: components2}

	results := []*parsingResult{result1, result2}
	filtered := pipeline.applyStrictValidation(results)

	// Only result1 should pass
	if len(filtered) != 1 {
		t.Errorf("Expected 1 result after strict validation, got %d", len(filtered))
	}
	if len(filtered) > 0 && filtered[0] != result1 {
		t.Error("Expected result1 to pass strict validation")
	}
}
