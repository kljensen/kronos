package kronos

import (
	"testing"
	"time"
)

func TestNewParsingContext(t *testing.T) {
	t.Run("NewParsingContext with time.Time", func(t *testing.T) {
		text := "tomorrow at 3pm"
		refDate := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)

		ctx := newParsingContext(text, refDate, nil)

		if ctx.Text() != text {
			t.Errorf("Expected text '%s', got '%s'", text, ctx.Text())
		}

		if !ctx.RefDate().Equal(refDate) {
			t.Errorf("Expected refDate %v, got %v", refDate, ctx.RefDate())
		}

		if ctx.Reference() == nil {
			t.Errorf("Expected reference to be non-nil")
		}
	})

	t.Run("NewParsingContext with parsingReference", func(t *testing.T) {
		text := "tomorrow"
		instant := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
		offset := -300

		parsingRef := parsingReference{
			Instant:  &instant,
			Timezone: offset,
		}

		ctx := newParsingContext(text, parsingRef, nil)

		if !ctx.RefDate().Equal(instant) {
			t.Errorf("Expected refDate %v, got %v", instant, ctx.RefDate())
		}

		if ctx.Reference().GetTimezoneOffset() != offset {
			t.Errorf("Expected timezone offset %d, got %d", offset, ctx.Reference().GetTimezoneOffset())
		}
	})

	t.Run("NewParsingContext with nil refDate uses current time", func(t *testing.T) {
		text := "tomorrow"
		before := time.Now()
		ctx := newParsingContext(text, nil, nil)
		after := time.Now()

		if ctx.RefDate().Before(before) || ctx.RefDate().After(after) {
			t.Errorf("Expected refDate to be current time, got %v", ctx.RefDate())
		}
	})

	t.Run("NewParsingContext with nil option uses defaults", func(t *testing.T) {
		text := "tomorrow"
		ctx := newParsingContext(text, nil, nil)

		opt := ctx.Option()

		if opt.ForwardDate {
			t.Errorf("Expected ForwardDate to be false by default")
		}

		if opt.Timezones != nil {
			t.Errorf("Expected Timezones to be nil by default")
		}

		if opt.Debug != nil {
			t.Errorf("Expected Debug to be nil by default")
		}
	})

	t.Run("NewParsingContext with custom option", func(t *testing.T) {
		text := "tomorrow"
		debugCalled := false
		debugHandler := func(message string) {
			debugCalled = true
		}

		option := &parsingOption{
			ForwardDate: true,
			Debug:       debugHandler,
		}

		ctx := newParsingContext(text, nil, option)

		opt := ctx.Option()

		if !opt.ForwardDate {
			t.Errorf("Expected ForwardDate to be true")
		}

		if opt.Debug == nil {
			t.Errorf("Expected Debug handler to be set")
		}

		// Verify debug handler works
		ctx.Debug(func() {
			debugHandler("test")
		})

		if !debugCalled {
			t.Errorf("Expected debug handler to be called")
		}
	})
}

func TestCreateParsingComponents(t *testing.T) {
	refDate := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
	ctx := newParsingContext("test", refDate, nil)

	t.Run("CreateParsingComponents with nil creates default", func(t *testing.T) {
		pc := ctx.CreateParsingComponents(nil)

		if pc == nil {
			t.Fatalf("Expected non-nil ParsingComponents")
		}

		// Should have default implied values
		if pc.Get(ComponentYear) == nil {
			t.Errorf("Expected year to be implied")
		}
	})

	t.Run("CreateParsingComponents with component map", func(t *testing.T) {
		components := map[Component]int{
			ComponentYear:  2025,
			ComponentMonth: 12,
			ComponentDay:   25,
		}

		pc := ctx.CreateParsingComponents(components)

		if pc == nil {
			t.Fatalf("Expected non-nil ParsingComponents")
		}

		if !pc.IsCertain(ComponentYear) {
			t.Errorf("Expected year to be certain")
		}

		if val := pc.Get(ComponentYear); val == nil || *val != 2025 {
			t.Errorf("Expected year 2025, got %v", val)
		}

		if val := pc.Get(ComponentMonth); val == nil || *val != 12 {
			t.Errorf("Expected month 12, got %v", val)
		}
	})

	t.Run("CreateParsingComponents with existing ParsingComponents returns as-is", func(t *testing.T) {
		original := newParsingComponents(ctx.Reference(), nil)
		original.Assign(ComponentYear, 2026)

		result := ctx.CreateParsingComponents(original)

		if result != original {
			t.Errorf("Expected same ParsingComponents instance to be returned")
		}

		if val := result.Get(ComponentYear); val == nil || *val != 2026 {
			t.Errorf("Expected year 2026, got %v", val)
		}
	})
}

func TestCreateParsingResult(t *testing.T) {
	refDate := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
	ctx := newParsingContext("tomorrow at 3pm", refDate, nil)

	t.Run("CreateParsingResult with text", func(t *testing.T) {
		result := ctx.CreateParsingResult(0, "tomorrow")

		if result.Index() != 0 {
			t.Errorf("Expected index 0, got %d", result.Index())
		}

		if result.Text() != "tomorrow" {
			t.Errorf("Expected text 'tomorrow', got '%s'", result.Text())
		}

		if result.Start() == nil {
			t.Errorf("Expected start to be created")
		}
	})

	t.Run("CreateParsingResult with endIndex", func(t *testing.T) {
		result := ctx.CreateParsingResult(0, 8) // "tomorrow"

		if result.Text() != "tomorrow" {
			t.Errorf("Expected text 'tomorrow', got '%s'", result.Text())
		}
	})

	t.Run("CreateParsingResult with start components as map", func(t *testing.T) {
		startComponents := map[Component]int{
			ComponentYear:  2024,
			ComponentMonth: 11,
			ComponentDay:   3,
		}

		result := ctx.CreateParsingResult(0, "tomorrow", startComponents)

		if result.Start() == nil {
			t.Fatalf("Expected start components to be created")
		}

		start := result.Start().(*parsingComponents)
		if !start.IsCertain(ComponentYear) {
			t.Errorf("Expected year to be certain in start components")
		}

		if val := start.Get(ComponentYear); val == nil || *val != 2024 {
			t.Errorf("Expected year 2024, got %v", val)
		}
	})

	t.Run("CreateParsingResult with start and end components", func(t *testing.T) {
		startComponents := map[Component]int{
			ComponentYear:  2024,
			ComponentMonth: 11,
			ComponentDay:   3,
		}

		endComponents := map[Component]int{
			ComponentYear:  2024,
			ComponentMonth: 11,
			ComponentDay:   5,
		}

		result := ctx.CreateParsingResult(0, "Nov 3-5", startComponents, endComponents)

		if result.Start() == nil {
			t.Errorf("Expected start components to be created")
		}

		if result.End() == nil {
			t.Errorf("Expected end components to be created")
		}

		start := result.Start().(*parsingComponents)
		if val := start.Get(ComponentDay); val == nil || *val != 3 {
			t.Errorf("Expected start day 3, got %v", val)
		}

		end := result.End().(*parsingComponents)
		if val := end.Get(ComponentDay); val == nil || *val != 5 {
			t.Errorf("Expected end day 5, got %v", val)
		}
	})

	t.Run("CreateParsingResult with ParsingComponents instances", func(t *testing.T) {
		start := newParsingComponents(ctx.Reference(), nil)
		start.Assign(ComponentYear, 2024)

		end := newParsingComponents(ctx.Reference(), nil)
		end.Assign(ComponentYear, 2025)

		result := ctx.CreateParsingResult(0, "test", start, end)

		if result.Start() != start {
			t.Errorf("Expected start to be the same instance")
		}

		if result.End() != end {
			t.Errorf("Expected end to be the same instance")
		}
	})

	t.Run("CreateParsingResult with endIndex bounds checking", func(t *testing.T) {
		// endIndex beyond text length
		result := ctx.CreateParsingResult(0, 1000)

		if result.Text() != ctx.Text() {
			t.Errorf("Expected full text when endIndex exceeds length")
		}

		// negative index
		result = ctx.CreateParsingResult(-5, 5)

		if result.Index() != 0 {
			t.Errorf("Expected index to be clamped to 0")
		}
	})
}

func TestDebug(t *testing.T) {
	t.Run("Debug calls handler when set", func(t *testing.T) {
		called := false
		debugHandler := func(message string) {
			called = true
		}

		option := &parsingOption{
			Debug: debugHandler,
		}

		ctx := newParsingContext("test", nil, option)

		ctx.Debug(func() {
			debugHandler("test message")
		})

		if !called {
			t.Errorf("Expected debug handler to be called")
		}
	})

	t.Run("Debug does nothing when handler not set", func(t *testing.T) {
		ctx := newParsingContext("test", nil, nil)

		// Should not panic
		ctx.Debug(func() {
			// This should not be executed
			panic("Should not be called")
		})
	})
}

func TestContextAccessors(t *testing.T) {
	text := "tomorrow at 3pm"
	refDate := time.Date(2024, 11, 2, 14, 30, 0, 0, time.UTC)
	option := &parsingOption{
		ForwardDate: true,
	}

	ctx := newParsingContext(text, refDate, option)

	t.Run("Text returns input text", func(t *testing.T) {
		if ctx.Text() != text {
			t.Errorf("Expected text '%s', got '%s'", text, ctx.Text())
		}
	})

	t.Run("Option returns parsing option", func(t *testing.T) {
		opt := ctx.Option()

		if !opt.ForwardDate {
			t.Errorf("Expected ForwardDate to be true")
		}
	})

	t.Run("Reference returns reference with timezone", func(t *testing.T) {
		if ctx.Reference() == nil {
			t.Errorf("Expected reference to be non-nil")
		}

		if !ctx.Reference().Instant().Equal(refDate) {
			t.Errorf("Expected reference instant %v, got %v", refDate, ctx.Reference().Instant())
		}
	})

	t.Run("RefDate returns reference date", func(t *testing.T) {
		if !ctx.RefDate().Equal(refDate) {
			t.Errorf("Expected refDate %v, got %v", refDate, ctx.RefDate())
		}
	})
}
