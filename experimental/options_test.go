package experimental_test

import (
	"testing"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/experimental"
)

// TestWithParserOrder tests the WithParserOrder option.
func TestWithParserOrder(t *testing.T) {
	// Test that the option can be applied without panicking
	// We test the actual settings application separately
	_ = kronos.New(nil).
		WithOption(experimental.WithParserOrder("iso8601", "en_casual_date"))
}

// TestWithEnabledParsers tests the WithEnabledParsers option.
func TestWithEnabledParsers(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithEnabledParsers("iso8601", "en_casual_date")
	opt(&settings)

	if len(settings.EnabledParsers) != 2 {
		t.Errorf("Expected 2 enabled parsers, got %d", len(settings.EnabledParsers))
	}
	if settings.EnabledParsers[0] != "iso8601" {
		t.Errorf("Expected first parser to be iso8601, got %s", settings.EnabledParsers[0])
	}
}

// TestDisableDefaultParsers tests the DisableDefaultParsers option.
func TestDisableDefaultParsers(t *testing.T) {
	settings := kronos.DefaultSettings()
	settings.EnabledParsers = []string{"default1", "default2"}

	opt := experimental.DisableDefaultParsers()
	opt(&settings)

	if len(settings.EnabledParsers) != 0 {
		t.Errorf("Expected no enabled parsers, got %d", len(settings.EnabledParsers))
	}
}

// TestWithMaxParsers tests the WithMaxParsers option.
func TestWithMaxParsers(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithMaxParsers(5)
	opt(&settings)

	if settings.MaxParsers != 5 {
		t.Errorf("Expected MaxParsers=5, got %d", settings.MaxParsers)
	}
}

// TestWithTimeout tests the WithTimeout option.
func TestWithTimeout(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithTimeout(10 * time.Second)
	opt(&settings)

	if settings.Timeout != 10*time.Second {
		t.Errorf("Expected Timeout=10s, got %v", settings.Timeout)
	}
}

// TestWithSkipTokens tests the WithSkipTokens option.
func TestWithSkipTokens(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithSkipTokens("at", "on", "the")
	opt(&settings)

	if len(settings.SkipTokens) != 3 {
		t.Errorf("Expected 3 skip tokens, got %d", len(settings.SkipTokens))
	}
	if settings.SkipTokens[0] != "at" {
		t.Errorf("Expected first token to be 'at', got %s", settings.SkipTokens[0])
	}
}

// TestWithRequiredParts tests the WithRequiredParts option.
func TestWithRequiredParts(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithRequiredParts("year", "month", "day")
	opt(&settings)

	if len(settings.RequireParts) != 3 {
		t.Errorf("Expected 3 required parts, got %d", len(settings.RequireParts))
	}
	if settings.RequireParts[0] != "year" {
		t.Errorf("Expected first part to be 'year', got %s", settings.RequireParts[0])
	}
}

// TestWithNormalization tests the WithNormalization option.
func TestWithNormalization(t *testing.T) {
	settings := kronos.DefaultSettings()
	opt := experimental.WithNormalization(false)
	opt(&settings)

	if settings.Normalize {
		t.Error("Expected Normalize to be false")
	}
}

// TestWithRelativeBase tests the WithRelativeBase option.
func TestWithRelativeBase(t *testing.T) {
	baseTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	settings := kronos.DefaultSettings()
	opt := experimental.WithRelativeBase(baseTime)
	opt(&settings)

	if settings.RelativeBase == nil {
		t.Fatal("Expected RelativeBase to be set")
	}
	if !settings.RelativeBase.Equal(baseTime) {
		t.Errorf("Expected RelativeBase to be %v, got %v", baseTime, *settings.RelativeBase)
	}
}

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

// TestChainedOptions tests that multiple options can be chained.
func TestChainedOptions(t *testing.T) {
	settings := kronos.DefaultSettings()

	opt1 := experimental.WithMaxParsers(3)
	opt2 := experimental.WithTimeout(5 * time.Second)
	opt3 := experimental.WithSkipTokens("at", "on")

	opt1(&settings)
	opt2(&settings)
	opt3(&settings)

	if settings.MaxParsers != 3 {
		t.Errorf("Expected MaxParsers=3, got %d", settings.MaxParsers)
	}
	if settings.Timeout != 5*time.Second {
		t.Errorf("Expected Timeout=5s, got %v", settings.Timeout)
	}
	if len(settings.SkipTokens) != 2 {
		t.Errorf("Expected 2 skip tokens, got %d", len(settings.SkipTokens))
	}
}
