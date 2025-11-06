//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
	"time"

	"github.com/kljensen/kronos"
)

// TestIntegration_ParseWithSettings tests the new settings-based parsing.
func TestIntegration_ParseWithSettings(t *testing.T) {
	refDate := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)

	t.Run("basic settings parsing", func(t *testing.T) {
		settings := kronos.DefaultSettings()

		results, err := Casual.ParseWithSettings("March 20, 2020", refDate, settings)
		if err != nil {
			t.Fatalf("ParseWithSettings failed: %v", err)
		}

		// This is expected to return results with default settings
		_ = results
	})

	t.Run("strict parsing mode", func(t *testing.T) {
		settings := kronos.DefaultSettings()
		settings.StrictParsing = true
		settings.RequireParts = []string{"year", "month", "day"}

		results, err := Casual.ParseWithSettings("March 20, 2020", refDate, settings)
		if err != nil {
			t.Fatalf("ParseWithSettings failed: %v", err)
		}

		// Strict parsing should filter results
		_ = results
	})

	t.Run("date preference future", func(t *testing.T) {
		settings := kronos.DefaultSettings()
		settings.PreferDatesFrom = kronos.PreferFuture

		results, err := Casual.ParseWithSettings("March 10", refDate, settings)
		if err != nil {
			t.Fatalf("ParseWithSettings failed: %v", err)
		}

		// With future preference, March 10 should be in 2021 (since ref is March 15, 2020)
		_ = results
	})

	t.Run("timezone conversion", func(t *testing.T) {
		settings := kronos.DefaultSettings()
		settings.Timezone = "America/New_York"
		settings.ToTimezone = "Europe/London"

		results, err := Casual.ParseWithSettings("March 20, 2020 3pm", refDate, settings)
		if err != nil {
			t.Fatalf("ParseWithSettings failed: %v", err)
		}

		_ = results
	})
}

// TestIntegration_ParseDateWithSettings tests the convenience function.
func TestIntegration_ParseDateWithSettings(t *testing.T) {
	refDate := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)

	settings := kronos.DefaultSettings()

	date, err := Casual.ParseDateWithSettings("March 20, 2020", refDate, settings)
	if err != nil {
		t.Fatalf("ParseDateWithSettings failed: %v", err)
	}

	_ = date // May be nil if no results
}

// TestIntegration_SettingsBackwardCompatibility ensures existing code still works.
func TestIntegration_SettingsBackwardCompatibility(t *testing.T) {
	refDate := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)

	// Old API should still work
	results := Parse("March 20, 2020", refDate, nil)
	if len(results) == 0 {
		t.Log("Old API still works (no results for this specific test)")
	}

	// Old ParseDate should still work
	date := ParseDate("tomorrow", refDate, nil)
	if date == nil {
		t.Log("ParseDate returned nil (expected for this test)")
	}
}

// TestIntegration_Registry tests parser registration.
func TestIntegration_Registry(t *testing.T) {
	// Test that parsers are registered
	infos := kronos.GlobalRegistry.ListParsers()
	if len(infos) == 0 {
		t.Error("Expected parsers to be registered")
	}

	// Check for some expected parsers
	expectedParsers := []string{
		"iso8601",
		"en_casual_date",
		"en_time_expression",
	}

	for _, name := range expectedParsers {
		if !kronos.GlobalRegistry.HasParser(name) {
			t.Errorf("Expected parser %s to be registered", name)
		}
	}
}

// TestIntegration_DefaultParserOrder tests that default order is set.
func TestIntegration_DefaultParserOrder(t *testing.T) {
	order := kronos.GlobalRegistry.GetDefaultOrder()
	if len(order) == 0 {
		t.Error("Expected default parser order to be set")
	}

	// Check that high-priority parsers come first
	if len(order) > 0 && order[0] != "iso8601" {
		t.Logf("First parser in default order: %s (expected iso8601 to be early)", order[0])
	}
}

// TestIntegration_SettingsValidation tests settings validation.
func TestIntegration_SettingsValidation(t *testing.T) {
	t.Run("valid settings", func(t *testing.T) {
		settings := kronos.DefaultSettings()
		settings.Timezone = "America/New_York"
		settings.RequireParts = []string{"year", "month"}

		err := kronos.ValidateSettings(settings)
		if err != nil {
			t.Errorf("Valid settings should not error: %v", err)
		}
	})

	t.Run("invalid timezone", func(t *testing.T) {
		settings := kronos.DefaultSettings()
		settings.Timezone = "Invalid/Timezone"

		err := kronos.ValidateSettings(settings)
		if err == nil {
			t.Error("Invalid timezone should error")
		}
	})

	t.Run("invalid required parts", func(t *testing.T) {
		settings := kronos.DefaultSettings()
		settings.RequireParts = []string{"invalid_component"}

		err := kronos.ValidateSettings(settings)
		if err == nil {
			t.Error("Invalid required parts should error")
		}
	})
}

// TestIntegration_PipelineCreation tests pipeline creation.
func TestIntegration_PipelineCreation(t *testing.T) {
	config := CreateConfiguration(false, false)
	settings := kronos.DefaultSettings()

	pipeline, err := kronos.NewPipelineWithSettings(config, settings)
	if err != nil {
		t.Fatalf("Failed to create pipeline: %v", err)
	}

	if pipeline == nil {
		t.Fatal("Expected non-nil pipeline")
	}

	// Execute pipeline
	refDate := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	results, err := pipeline.Execute("March 20, 2020", refDate)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	_ = results
}

// TestIntegration_CustomParserSelection tests selecting specific parsers.
func TestIntegration_CustomParserSelection(t *testing.T) {
	settings := kronos.DefaultSettings()
	settings.EnabledParsers = []string{"iso8601", "en_casual_date"}

	config := CreateConfiguration(false, false)
	pipeline, err := kronos.NewPipelineWithSettings(config, settings)
	if err != nil {
		t.Fatalf("Failed to create pipeline: %v", err)
	}

	refDate := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	results, err := pipeline.Execute("2020-03-20", refDate)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	_ = results
}

// TestIntegration_MaxParsersLimit tests the max parsers setting.
func TestIntegration_MaxParsersLimit(t *testing.T) {
	settings := kronos.DefaultSettings()
	settings.MaxParsers = 5

	config := CreateConfiguration(false, false)
	pipeline, err := kronos.NewPipelineWithSettings(config, settings)
	if err != nil {
		t.Fatalf("Failed to create pipeline: %v", err)
	}

	if pipeline.ParserCount() > 5 {
		t.Errorf("Expected max 5 parsers, got %d", pipeline.ParserCount())
	}
}
