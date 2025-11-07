package experimental_test

import (
	"testing"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/experimental"
)

// TestWithDayPreference tests the WithDayPreference option.
func TestWithDayPreference(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithDayPreference(kronos.DayPreferFirst)
	opt(&settings)

	if settings.PreferDayOfMonth != kronos.DayPreferFirst {
		t.Errorf("Expected DayPreferFirst, got %v", settings.PreferDayOfMonth)
	}
}

// TestWithTimezoneAware tests the WithTimezoneAware option.
func TestWithTimezoneAware(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithTimezoneAware(true)
	opt(&settings)

	if !settings.ReturnTimezoneAware {
		t.Error("Expected ReturnTimezoneAware to be true")
	}
}

// TestWithTimeAsPeriod tests the WithTimeAsPeriod option.
func TestWithTimeAsPeriod(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithTimeAsPeriod(true)
	opt(&settings)

	if !settings.ReturnTimeAsPeriod {
		t.Error("Expected ReturnTimeAsPeriod to be true")
	}
}

// TestWithTimezoneOverrides tests the WithTimezoneOverrides option.
func TestWithTimezoneOverrides(t *testing.T) {
	customTimezones := kronos.TimezoneAbbrMap{
		"CUSTOM": 123,
	}
	settings := kronos.DefaultSettings()
	opt := experimental.WithTimezoneOverrides(customTimezones)
	opt(&settings)

	if settings.TimezoneOverrides == nil {
		t.Fatal("Expected TimezoneOverrides to be set")
	}
	if val, ok := settings.TimezoneOverrides["CUSTOM"]; !ok || val != 123 {
		t.Errorf("Expected CUSTOM timezone to be 123, got %v", val)
	}
}

// TestWithDebugHandler tests the WithDebugHandler option.
func TestWithDebugHandler(t *testing.T) {
	settings := kronos.DefaultSettings()
	called := false
	handler := func(msg string) {
		called = true
	}
	opt := experimental.WithDebugHandler(handler)
	opt(&settings)

	if settings.DebugHandler == nil {
		t.Fatal("Expected DebugHandler to be set")
	}
	// Test that the handler works
	settings.DebugHandler("test")
	if !called {
		t.Error("Expected debug handler to be called")
	}
}

// TestChainedOptions tests that multiple options can be chained.
func TestChainedOptions(t *testing.T) {
	settings := kronos.DefaultSettings()

	opt1 := experimental.WithDayPreference(kronos.DayPreferFirst)
	opt2 := experimental.WithTimezoneAware(true)
	opt3 := experimental.WithTimeAsPeriod(true)

	opt1(&settings)
	opt2(&settings)
	opt3(&settings)

	if settings.PreferDayOfMonth != kronos.DayPreferFirst {
		t.Errorf("Expected DayPreferFirst, got %v", settings.PreferDayOfMonth)
	}
	if !settings.ReturnTimezoneAware {
		t.Error("Expected ReturnTimezoneAware to be true")
	}
	if !settings.ReturnTimeAsPeriod {
		t.Error("Expected ReturnTimeAsPeriod to be true")
	}
}
