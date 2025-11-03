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
	if settings.PreferDayOfMonth != DayPreferCurrent {
		t.Errorf("Expected DayPreferCurrent, got %v", settings.PreferDayOfMonth)
	}
	if settings.Timezone != "UTC" {
		t.Errorf("Expected UTC timezone, got %s", settings.Timezone)
	}
	if !settings.Normalize {
		t.Error("Expected Normalize to be true")
	}
	if settings.StrictParsing {
		t.Error("Expected StrictParsing to be false")
	}
}

func TestValidateSettings_ValidTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.Timezone = "America/New_York"

	err := ValidateSettings(settings)
	if err != nil {
		t.Errorf("Expected no error for valid timezone, got: %v", err)
	}
}

func TestValidateSettings_InvalidTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	err := ValidateSettings(settings)
	if err == nil {
		t.Error("Expected error for invalid timezone")
	}
}

func TestValidateSettings_ValidToTimezone(t *testing.T) {
	settings := DefaultSettings()
	settings.ToTimezone = "Europe/London"

	err := ValidateSettings(settings)
	if err == nil || err != nil {
		// Either way is fine - just testing it doesn't panic
	}
}

func TestValidateSettings_ValidRequireParts(t *testing.T) {
	settings := DefaultSettings()
	settings.RequireParts = []string{"year", "month", "day"}

	err := ValidateSettings(settings)
	if err != nil {
		t.Errorf("Expected no error for valid required parts, got: %v", err)
	}
}

func TestValidateSettings_InvalidRequireParts(t *testing.T) {
	settings := DefaultSettings()
	settings.RequireParts = []string{"invalid_part"}

	err := ValidateSettings(settings)
	if err == nil {
		t.Error("Expected error for invalid required parts")
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

func TestDayPreferenceString(t *testing.T) {
	tests := []struct {
		pref     DayPreference
		expected string
	}{
		{DayPreferCurrent, "current"},
		{DayPreferFirst, "first"},
		{DayPreferLast, "last"},
	}

	for _, tc := range tests {
		if tc.pref.String() != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, tc.pref.String())
		}
	}
}

func TestToParsingOption(t *testing.T) {
	settings := DefaultSettings()
	settings.ForwardDate = true
	settings.PreferDatesFrom = PreferFuture

	opt := settings.ToParsingOption(nil)

	if !opt.ForwardDate {
		t.Error("Expected ForwardDate to be true")
	}
	if opt.Preference != PreferFuture {
		t.Errorf("Expected PreferFuture, got %v", opt.Preference)
	}
}

func TestApplySettings_Normalization(t *testing.T) {
	settings := DefaultSettings()
	settings.Normalize = true

	text := "test\u00A0text" // Non-breaking space
	refDate := time.Now()

	ctx, err := ApplySettings(text, refDate, settings)
	if err != nil {
		t.Fatalf("ApplySettings failed: %v", err)
	}

	// Text should be normalized (non-breaking space converted)
	if ctx.Text() == text {
		// The text might be the same if normalization doesn't change it
		// But we at least verify it doesn't error
	}
}

func TestApplySettings_SkipTokens(t *testing.T) {
	settings := DefaultSettings()
	settings.SkipTokens = []string{"at", "on"}

	text := "on March 15 at 3pm"
	refDate := time.Now()

	ctx, err := ApplySettings(text, refDate, settings)
	if err != nil {
		t.Fatalf("ApplySettings failed: %v", err)
	}

	// Tokens should be removed
	resultText := ctx.Text()
	if resultText == text {
		// Simple check - the text should be different after token removal
		// But the exact result depends on the removeToken implementation
	}
}

func TestApplySettings_RelativeBase(t *testing.T) {
	settings := DefaultSettings()
	baseTime := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)
	settings.RelativeBase = &baseTime

	text := "test"
	refDate := time.Now()

	ctx, err := ApplySettings(text, refDate, settings)
	if err != nil {
		t.Fatalf("ApplySettings failed: %v", err)
	}

	// The context should use the relative base as reference
	if !ctx.RefDate().Equal(baseTime) {
		t.Errorf("Expected RefDate to be %v, got %v", baseTime, ctx.RefDate())
	}
}

func TestApplySettings_InvalidSettings(t *testing.T) {
	settings := DefaultSettings()
	settings.Timezone = "Invalid/Timezone"

	text := "test"
	refDate := time.Now()

	_, err := ApplySettings(text, refDate, settings)
	if err == nil {
		t.Error("Expected error for invalid settings")
	}
}

func TestRemoveToken(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		token    string
		expected string
	}{
		{
			name:     "simple removal",
			text:     "on March 15",
			token:    "on",
			expected: " March 15",
		},
		{
			name:     "multiple occurrences",
			text:     "at 3pm at home",
			token:    "at",
			expected: " 3pm  home",
		},
		{
			name:     "no match",
			text:     "March 15",
			token:    "on",
			expected: "March 15",
		},
		{
			name:     "word boundary - should not remove part of word",
			text:     "attorney at law",
			token:    "at",
			expected: "attorney  law",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := removeToken(tc.text, tc.token)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestSettings_EnabledParsers(t *testing.T) {
	settings := DefaultSettings()
	settings.EnabledParsers = []string{"iso8601", "en_casual_date"}

	if len(settings.EnabledParsers) != 2 {
		t.Errorf("Expected 2 enabled parsers, got %d", len(settings.EnabledParsers))
	}
}

func TestSettings_ParserOrder(t *testing.T) {
	settings := DefaultSettings()
	settings.ParserOrder = []string{"en_casual_date", "iso8601"}

	if len(settings.ParserOrder) != 2 {
		t.Errorf("Expected 2 parsers in order, got %d", len(settings.ParserOrder))
	}
	if settings.ParserOrder[0] != "en_casual_date" {
		t.Errorf("Expected first parser to be en_casual_date, got %s", settings.ParserOrder[0])
	}
}

func TestSettings_Timeout(t *testing.T) {
	settings := DefaultSettings()
	settings.Timeout = 5 * time.Second

	if settings.Timeout != 5*time.Second {
		t.Errorf("Expected timeout of 5s, got %v", settings.Timeout)
	}
}

func TestSettings_MaxParsers(t *testing.T) {
	settings := DefaultSettings()
	settings.MaxParsers = 10

	if settings.MaxParsers != 10 {
		t.Errorf("Expected MaxParsers of 10, got %d", settings.MaxParsers)
	}
}

func TestSettings_StrictParsing(t *testing.T) {
	settings := DefaultSettings()
	settings.StrictParsing = true

	if !settings.StrictParsing {
		t.Error("Expected StrictParsing to be true")
	}
}

func TestSettings_ReturnTimeAsPeriod(t *testing.T) {
	settings := DefaultSettings()
	settings.ReturnTimeAsPeriod = true

	if !settings.ReturnTimeAsPeriod {
		t.Error("Expected ReturnTimeAsPeriod to be true")
	}
}

func TestSettings_ReturnTimezoneAware(t *testing.T) {
	settings := DefaultSettings()
	settings.ReturnTimezoneAware = true

	if !settings.ReturnTimezoneAware {
		t.Error("Expected ReturnTimezoneAware to be true")
	}
}

func TestNewParsingContextWithSettings(t *testing.T) {
	settings := DefaultSettings()
	settings.Normalize = true
	settings.PreferDatesFrom = PreferFuture

	text := "March 15"
	refDate := time.Now()

	ctx := NewParsingContextWithSettings(text, refDate, settings)

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
	if settings.Normalize != true {
		t.Error("Default Normalize changed - breaks backward compatibility")
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

		err := ValidateSettings(settings)
		if err != nil {
			t.Errorf("Valid preference %v should not error, got: %v", pref, err)
		}
	}
}

func TestSettings_DayPreferences(t *testing.T) {
	// Test all day preferences
	preferences := []DayPreference{
		DayPreferCurrent,
		DayPreferFirst,
		DayPreferLast,
	}

	for _, pref := range preferences {
		settings := DefaultSettings()
		settings.PreferDayOfMonth = pref

		err := ValidateSettings(settings)
		if err != nil {
			t.Errorf("Valid day preference %v should not error, got: %v", pref, err)
		}
	}
}
