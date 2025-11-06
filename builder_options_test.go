package kronos

import (
	"testing"
	"time"
)

// TestBuilderWithCommonOptions tests the builder pattern with common options.
func TestBuilderWithCommonOptions(t *testing.T) {
	refDate := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)

	builder := New(nil).
		WithReferenceDate(refDate).
		WithDateOrder(DateOrderDMY).
		PreferFuture().
		Strict()

	// Verify settings are applied
	if builder.settings.DateOrder != DateOrderDMY {
		t.Errorf("Expected DateOrderDMY, got %v", builder.settings.DateOrder)
	}
	if builder.settings.PreferDatesFrom != PreferFuture {
		t.Errorf("Expected PreferFuture, got %v", builder.settings.PreferDatesFrom)
	}
	if !builder.settings.StrictParsing {
		t.Error("Expected StrictParsing to be true")
	}
}

// TestBuilderCasualMode tests the Casual() method.
func TestBuilderCasualMode(t *testing.T) {
	builder := New(nil).Casual()

	if builder.settings.StrictParsing {
		t.Error("Expected StrictParsing to be false in casual mode")
	}
}

// TestBuilderPreferences tests the preference methods.
func TestBuilderPreferences(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*ParserBuilder) *ParserBuilder
		expected DatePreference
	}{
		{
			name:     "PreferPast",
			setup:    func(b *ParserBuilder) *ParserBuilder { return b.PreferPast() },
			expected: PreferPast,
		},
		{
			name:     "PreferFuture",
			setup:    func(b *ParserBuilder) *ParserBuilder { return b.PreferFuture() },
			expected: PreferFuture,
		},
		{
			name:     "PreferCurrentPeriod",
			setup:    func(b *ParserBuilder) *ParserBuilder { return b.PreferCurrentPeriod() },
			expected: PreferCurrentPeriod,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			builder := New(nil)
			builder = tc.setup(builder)

			if builder.settings.PreferDatesFrom != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, builder.settings.PreferDatesFrom)
			}
		})
	}
}

// TestBuilderTimezone tests timezone settings.
func TestBuilderTimezone(t *testing.T) {
	builder := New(nil).
		Timezone("America/New_York").
		ToTimezone("Europe/London")

	if builder.settings.Timezone != "America/New_York" {
		t.Errorf("Expected America/New_York, got %s", builder.settings.Timezone)
	}
	if builder.settings.ToTimezone != "Europe/London" {
		t.Errorf("Expected Europe/London, got %s", builder.settings.ToTimezone)
	}
}

// TestBuilderChaining tests that methods can be chained.
func TestBuilderChaining(t *testing.T) {
	refDate := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)

	// This should not panic
	_ = New(nil).
		WithReferenceDate(refDate).
		WithDateOrder(DateOrderYMD).
		PreferPast().
		Strict().
		Timezone("UTC")
}

// TestBuilderAdvancedOptions tests advanced configuration options.
func TestBuilderAdvancedOptions(t *testing.T) {
	customTimezones := TimezoneAbbrMap{
		"CUSTOM": 123,
	}

	debugCalled := false
	debugHandler := func(msg string) {
		debugCalled = true
	}

	parser := New(nil).
		WithOption(func(s *Settings) {
			s.TimezoneOverrides = customTimezones
		}).
		WithOption(func(s *Settings) {
			s.DebugHandler = debugHandler
		})

	// Verify settings
	if parser.settings.TimezoneOverrides == nil {
		t.Error("Expected timezone overrides to be set")
	}
	if parser.settings.DebugHandler == nil {
		t.Error("Expected debug handler to be set")
	}

	// Call handler
	parser.settings.DebugHandler("test")
	if !debugCalled {
		t.Error("Expected debug handler to be called")
	}
}
