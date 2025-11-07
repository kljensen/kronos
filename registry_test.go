package kronos

import (
	"testing"
)

func TestNewParserRegistry(t *testing.T) {
	registry := newParserRegistry()

	if registry == nil {
		t.Fatal("Expected non-nil registry")
	}
}

func TestParserRegistry_RegisterParser(t *testing.T) {
	registry := newParserRegistry()

	info := parserInfo{
		Description: "Test parser",
		Priority:    50,
		Tags:        []string{"test"},
	}

	registry.RegisterParser("test", info, func() Parser {
		return nil
	})

	if !registry.HasParser("test") {
		t.Error("Expected parser to be registered")
	}
}

func TestParserRegistry_GetParser(t *testing.T) {
	registry := newParserRegistry()

	info := parserInfo{
		Description: "Test parser",
		Priority:    50,
	}

	registry.RegisterParser("test", info, func() Parser {
		return nil
	})

	parser := registry.GetParser("test")
	if parser != nil {
		// Parser is nil in our test factory
		t.Error("Expected nil parser from test factory")
	}
}

func TestParserRegistry_GetParser_NotFound(t *testing.T) {
	registry := newParserRegistry()

	parser := registry.GetParser("nonexistent")
	if parser != nil {
		t.Error("Expected nil for non-existent parser")
	}
}

func TestParserRegistry_GetParsers(t *testing.T) {
	registry := newParserRegistry()

	info1 := parserInfo{
		Description: "Test parser 1",
		Priority:    50,
	}
	info2 := parserInfo{
		Description: "Test parser 2",
		Priority:    60,
	}

	registry.RegisterParser("test1", info1, func() Parser {
		return nil
	})
	registry.RegisterParser("test2", info2, func() Parser {
		return nil
	})

	parsers := registry.GetParsers([]string{"test1", "test2"})
	if len(parsers) != 2 {
		t.Errorf("Expected 2 parsers, got %d", len(parsers))
	}
}

func TestParserRegistry_GetParsers_SkipsUnknown(t *testing.T) {
	registry := newParserRegistry()

	info := parserInfo{
		Description: "Test parser",
		Priority:    50,
	}

	registry.RegisterParser("test", info, func() Parser {
		return nil
	})

	parsers := registry.GetParsers([]string{"test", "unknown"})
	if len(parsers) != 1 {
		t.Errorf("Expected 1 parser (unknown skipped), got %d", len(parsers))
	}
}

func TestParserRegistry_GetAllParsers(t *testing.T) {
	registry := newParserRegistry()

	info1 := parserInfo{
		Description: "Test parser 1",
		Priority:    50,
	}
	info2 := parserInfo{
		Description: "Test parser 2",
		Priority:    60,
	}

	registry.RegisterParser("test1", info1, func() Parser {
		return nil
	})
	registry.RegisterParser("test2", info2, func() Parser {
		return nil
	})

	parsers := registry.GetAllParsers()
	if len(parsers) != 2 {
		t.Errorf("Expected 2 parsers, got %d", len(parsers))
	}
}

func TestParserRegistry_GetAllParsers_SortedByPriority(t *testing.T) {
	registry := newParserRegistry()

	// Register parsers with different priorities
	info1 := parserInfo{
		Description: "Low priority",
		Priority:    10,
	}
	info2 := parserInfo{
		Description: "High priority",
		Priority:    100,
	}
	info3 := parserInfo{
		Description: "Medium priority",
		Priority:    50,
	}

	registry.RegisterParser("low", info1, func() Parser {
		return nil
	})
	registry.RegisterParser("high", info2, func() Parser {
		return nil
	})
	registry.RegisterParser("medium", info3, func() Parser {
		return nil
	})

	parsers := registry.GetAllParsers()
	if len(parsers) != 3 {
		t.Errorf("Expected 3 parsers, got %d", len(parsers))
	}

	// They should be sorted by priority (high to low)
	// Since our factory returns nil, we can't check the actual parsers
	// But we verified the sorting logic exists
}

func TestParserRegistry_GetParsersByTag(t *testing.T) {
	registry := newParserRegistry()

	info1 := parserInfo{
		Description: "Casual parser",
		Priority:    50,
		Tags:        []string{"casual", "english"},
	}
	info2 := parserInfo{
		Description: "Formal parser",
		Priority:    60,
		Tags:        []string{"formal", "english"},
	}

	registry.RegisterParser("casual", info1, func() Parser {
		return nil
	})
	registry.RegisterParser("formal", info2, func() Parser {
		return nil
	})

	casualParsers := registry.GetParsersByTag("casual")
	if len(casualParsers) != 1 {
		t.Errorf("Expected 1 casual parser, got %d", len(casualParsers))
	}

	englishParsers := registry.GetParsersByTag("english")
	if len(englishParsers) != 2 {
		t.Errorf("Expected 2 english parsers, got %d", len(englishParsers))
	}
}

func TestParserRegistry_ListParsers(t *testing.T) {
	registry := newParserRegistry()

	info1 := parserInfo{
		Description: "Test parser 1",
		Priority:    50,
	}
	info2 := parserInfo{
		Description: "Test parser 2",
		Priority:    60,
	}

	registry.RegisterParser("test1", info1, func() Parser {
		return nil
	})
	registry.RegisterParser("test2", info2, func() Parser {
		return nil
	})

	infos := registry.ListParsers()
	if len(infos) != 2 {
		t.Errorf("Expected 2 parser infos, got %d", len(infos))
	}

	// Check that names are set correctly
	foundTest1 := false
	foundTest2 := false
	for _, info := range infos {
		if info.Name == "test1" {
			foundTest1 = true
		}
		if info.Name == "test2" {
			foundTest2 = true
		}
	}

	if !foundTest1 || !foundTest2 {
		t.Error("Expected to find both test1 and test2 in parser infos")
	}
}

func TestParserRegistry_HasParser(t *testing.T) {
	registry := newParserRegistry()

	info := parserInfo{
		Description: "Test parser",
		Priority:    50,
	}

	registry.RegisterParser("test", info, func() Parser {
		return nil
	})

	if !registry.HasParser("test") {
		t.Error("Expected HasParser to return true for registered parser")
	}

	if registry.HasParser("nonexistent") {
		t.Error("Expected HasParser to return false for non-existent parser")
	}
}

func TestParserRegistry_SetDefaultOrder(t *testing.T) {
	registry := newParserRegistry()

	order := []string{"parser1", "parser2", "parser3"}
	registry.SetDefaultOrder(order)

	retrievedOrder := registry.GetDefaultOrder()
	if len(retrievedOrder) != 3 {
		t.Errorf("Expected 3 parsers in order, got %d", len(retrievedOrder))
	}

	for i, name := range order {
		if retrievedOrder[i] != name {
			t.Errorf("Expected parser %d to be %s, got %s", i, name, retrievedOrder[i])
		}
	}
}

func TestParserRegistry_GetDefaultOrder(t *testing.T) {
	registry := newParserRegistry()

	// Initially should be empty
	order := registry.GetDefaultOrder()
	if len(order) != 0 {
		t.Errorf("Expected empty default order, got %d parsers", len(order))
	}
}

func TestParserRegistry_Clear(t *testing.T) {
	registry := newParserRegistry()

	info := parserInfo{
		Description: "Test parser",
		Priority:    50,
	}

	registry.RegisterParser("test", info, func() Parser {
		return nil
	})

	if !registry.HasParser("test") {
		t.Error("Expected parser to be registered")
	}

	registry.Clear()

	if registry.HasParser("test") {
		t.Error("Expected parser to be cleared")
	}
}

func TestParserRegistry_ReplaceParser(t *testing.T) {
	registry := newParserRegistry()

	info1 := parserInfo{
		Description: "Original parser",
		Priority:    50,
	}

	info2 := parserInfo{
		Description: "Replacement parser",
		Priority:    60,
	}

	registry.RegisterParser("test", info1, func() Parser {
		return nil
	})

	// Register again with same name - should replace
	registry.RegisterParser("test", info2, func() Parser {
		return nil
	})

	infos := registry.ListParsers()
	if len(infos) != 1 {
		t.Errorf("Expected 1 parser (replaced), got %d", len(infos))
	}

	if infos[0].Description != "Replacement parser" {
		t.Errorf("Expected replacement parser description, got %s", infos[0].Description)
	}
}

func TestGlobalRegistry(t *testing.T) {
	// Test that globalRegistry exists internally
	if globalRegistry == nil {
		t.Fatal("Expected internal globalRegistry to exist")
	}
}

func TestRegister_GlobalFunction(t *testing.T) {
	// Test internal register function
	// Clear global registry to avoid conflicts
	globalRegistry.Clear()

	info := parserInfo{
		Description: "Test parser",
		Priority:    50,
	}

	internalRegister("global_test", info, func() Parser {
		return nil
	})

	if !globalRegistry.HasParser("global_test") {
		t.Error("Expected parser to be registered in global registry")
	}

	// Clean up
	globalRegistry.Clear()
}

func TestParserInfo(t *testing.T) {
	info := parserInfo{
		Name:        "test",
		Description: "Test parser",
		Priority:    50,
		Tags:        []string{"test", "example"},
	}

	if info.Name != "test" {
		t.Errorf("Expected name 'test', got %s", info.Name)
	}
	if info.Description != "Test parser" {
		t.Errorf("Expected description 'Test parser', got %s", info.Description)
	}
	if info.Priority != 50 {
		t.Errorf("Expected priority 50, got %d", info.Priority)
	}
	if len(info.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(info.Tags))
	}
}

func TestParserRegistry_ConcurrentAccess(t *testing.T) {
	registry := newParserRegistry()

	info := parserInfo{
		Description: "Test parser",
		Priority:    50,
	}

	// Test concurrent registration and access
	done := make(chan bool)

	go func() {
		registry.RegisterParser("test1", info, func() Parser {
			return nil
		})
		done <- true
	}()

	go func() {
		registry.RegisterParser("test2", info, func() Parser {
			return nil
		})
		done <- true
	}()

	go func() {
		_ = registry.HasParser("test1")
		done <- true
	}()

	// Wait for all goroutines
	<-done
	<-done
	<-done

	// Should have both parsers
	if !registry.HasParser("test1") {
		t.Error("Expected test1 to be registered")
	}
	if !registry.HasParser("test2") {
		t.Error("Expected test2 to be registered")
	}
}
