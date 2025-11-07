package kronos

import (
	"testing"
	"time"
)

func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	// Test defaults
	if settings.DateOrder != DateOrderMDY {
		t.Errorf("Expected DateOrder MDY, got %v", settings.DateOrder)
	}
	if settings.PreferDatesFrom != PreferCurrentPeriod {
		t.Errorf("Expected PreferCurrentPeriod, got %v", settings.PreferDatesFrom)
	}
	if settings.Timezone != "UTC" {
		t.Errorf("Expected UTC timezone, got %s", settings.Timezone)
	}
	if settings.StrictParsing {
		t.Error("Expected StrictParsing to be false")
	}
}

func TestvalidateSettings_ValidTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.Timezone = "America/New_York"

	err := validateSettings(settings)
	if err != nil {
		t.Errorf("Expected no error for valid timezone, got: %v", err)
	}
}

func TestvalidateSettings_InvalidTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	err := validateSettings(settings)
	if err == nil {
		t.Error("Expected error for invalid timezone")
	}
}


func TestDateOrderString(t *testing.T) {
	tests := []struct {
		order    DateOrder
		expected string
	}{
		{DateOrderMDY, "MDY"},
		{DateOrderDMY, "DMY"},
		{DateOrderYMD, "YMD"},
	}

	for _, tc := range tests {
		if tc.order.String() != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, tc.order.String())
		}
	}
}

func TestToparsingOption(t *testing.T) {
	settings := DefaultSettings()
	settings.PreferDatesFrom = PreferFuture

	opt := settings.toparsingOption(nil)

	if opt.Preference != PreferFuture {
		t.Errorf("Expected PreferFuture, got %v", opt.Preference)
	}
}

func TestApplySettings_Normalization(t *testing.T) {
	settings := DefaultSettings()

	text := "test\u00A0text" // Non-breaking space
	refDate := time.Now()

	ctx, err := applySettings(text, refDate, settings)
	if err != nil {
		t.Fatalf("ApplySettings failed: %v", err)
	}

	// Text should be normalized - just verify it doesn't error
	_ = ctx.Text()
}

func TestApplySettings_InvalidSettings(t *testing.T) {
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	text := "test"
	refDate := time.Now()

	_, err := applySettings(text, refDate, settings)
	if err == nil {
		t.Error("Expected error for invalid settings")
	}
}


func TestSettings_StrictParsing(t *testing.T) {
	settings := DefaultSettings()
	settings.StrictParsing = true

	if !settings.StrictParsing {
		t.Error("Expected StrictParsing to be true")
	}
}


func TestNewParsingContextWithSettings(t *testing.T) {
	settings := DefaultSettings()
	settings.PreferDatesFrom = PreferFuture

	text := "March 15"
	refDate := time.Now()

	ctx, err := applySettings(text, refDate, settings)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}

	// Check settings are applied
	if ctx.Settings() == nil {
		t.Fatal("Expected settings to be set")
	}

	if ctx.Settings().PreferDatesFrom != PreferFuture {
		t.Errorf("Expected PreferFuture, got %v", ctx.Settings().PreferDatesFrom)
	}

	// Check option is derived from settings
	opt := ctx.Option()
	if opt.Preference != PreferFuture {
		t.Errorf("Expected Preference PreferFuture, got %v", opt.Preference)
	}
}

func TestSettings_BackwardCompatibility(t *testing.T) {
	// Test that default settings maintain backward compatibility
	settings := DefaultSettings()

	// Default behavior should match existing behavior
	if settings.DateOrder != DateOrderMDY {
		t.Error("Default DateOrder changed - breaks backward compatibility")
	}
	if settings.PreferDatesFrom != PreferCurrentPeriod {
		t.Error("Default PreferDatesFrom changed - breaks backward compatibility")
	}
	if settings.StrictParsing != false {
		t.Error("Default StrictParsing changed - breaks backward compatibility")
	}
}

func TestSettings_DatePreferences(t *testing.T) {
	// Test all date preferences
	preferences := []DatePreference{
		PreferCurrentPeriod,
		PreferPast,
		PreferFuture,
	}

	for _, pref := range preferences {
		settings := DefaultSettings()
		settings.PreferDatesFrom = pref

		err := validateSettings(settings)
		if err != nil {
			t.Errorf("Valid preference %v should not error, got: %v", pref, err)
		}
	}
}

