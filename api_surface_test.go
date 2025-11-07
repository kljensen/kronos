package kronos

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// TestAPIExports ensures the public API surface doesn't grow unexpectedly.
// This test acts as a guard against API bloat by tracking exported symbols.
//
// The goal is to keep the API minimal and focused on what users actually need:
// - Core types: Component, DateOrder, DatePreference, Period
// - Duration support: Timeunit (needed for Duration map keys)
// - Deprecated: Weekday/Month (aliased to time package), Meridiem (internal), numeric constants
//
// If this test fails, it means new exports were added. Consider:
// 1. Is the new export truly necessary for users?
// 2. Can it be moved to internal/ package?
// 3. If it must be public, update MAX_EXPORTS with justification
func TestAPIExports(t *testing.T) {
	const (
		// Maximum allowed public exports after API minimization.
		// Current breakdown (71 exports):
		// - Essential enum types: Component (13), DateOrder (4), DatePreference (4) = 21 consts
		// - Duration/Time: Timeunit (13), Period (7) = 20 consts + Duration type
		// - Core API: Builder, Parser, Chrono, Configuration, Result/Components interfaces
		// - Helper functions: parsing, configuration, component access
		// - Internal helpers: 13 Internal* exports for internal package access
		//
		// Recent changes (Iteration 5):
		// - Removed code duplication (sanitizeInput, secondsPerMinute constant)
		// - Optimized Settings.toparsingOption to reduce allocations
		// - No API surface changes (internal refactoring only)
		//
		// Recent changes (Iteration 4):
		// - Moved Duration bounds constants to internal (8 constants)
		// - Removed Chrono.Clone, Chrono.ParseDate, Chrono.ParseDateWithSettings (3 methods)
		// - Added Tags() to Result and Components interfaces for testing (2 methods, net +2)
		// - Net reduction: 9 exports (from 83 to 74, then 71 after cleanup)
		//
		// Past achievements:
		// - Removed experimental package entirely
		// - Removed unused DayPreference feature
		// - Removed redundant builders & convenience functions
		//
		// Current: 71 exports (down from 90 in iteration 1)
		// Target: Maintain ~70-80 exports for essential API
		MAX_EXPORTS = 100 // Enforces lean essential API
	)

	// Use go doc to list all exports
	cmd := exec.CommandContext(context.Background(), "go", "doc", "-all", ".")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Logf("stderr: %s", stderr.String())
		t.Fatalf("Failed to run go doc: %v", err)
	}

	output := stdout.String()
	lines := strings.Split(output, "\n")

	// Count exported symbols (functions, types, constants)
	// Lines starting with "func ", "type ", "const (" indicate exports
	var exports []string
	var typeExports []string
	var constExports []string
	var funcExports []string

	inConstBlock := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Check for export declarations
		if strings.HasPrefix(line, "type ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[1]
				// Filter out generic type parameters
				if idx := strings.IndexAny(name, "[ "); idx > 0 {
					name = name[:idx]
				}
				if isExported(name) {
					exports = append(exports, name)
					typeExports = append(typeExports, name)
				}
			}
		} else if strings.HasPrefix(line, "func ") {
			// Extract function name
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[1]
				// Handle method receivers: func (Type) Method
				if name == "(" && len(parts) >= 4 {
					name = parts[3]
				}
				// Remove parameters
				if idx := strings.Index(name, "("); idx > 0 {
					name = name[:idx]
				}
				if isExported(name) {
					exports = append(exports, name)
					funcExports = append(funcExports, name)
				}
			}
		} else if strings.HasPrefix(line, "const (") {
			inConstBlock = true
		} else if inConstBlock {
			if line == ")" {
				inConstBlock = false
			} else if len(line) > 0 && !strings.HasPrefix(line, "//") {
				// Extract constant name
				parts := strings.Fields(line)
				if len(parts) > 0 {
					name := parts[0]
					if isExported(name) {
						exports = append(exports, name)
						constExports = append(constExports, name)
					}
				}
			}
		} else if strings.HasPrefix(line, "const ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[1]
				if isExported(name) {
					exports = append(exports, name)
					constExports = append(constExports, name)
				}
			}
		}
	}

	// Remove duplicates
	exports = uniqueStrings(exports)
	typeExports = uniqueStrings(typeExports)
	constExports = uniqueStrings(constExports)
	funcExports = uniqueStrings(funcExports)

	// Report current state
	t.Logf("Total exported symbols: %d (max: %d)", len(exports), MAX_EXPORTS)
	t.Logf("Exported types: %d", len(typeExports))
	t.Logf("Exported constants: %d", len(constExports))
	t.Logf("Exported functions: %d", len(funcExports))

	// Check if we exceed the limit
	if len(exports) > MAX_EXPORTS {
		t.Errorf("Too many exports: %d (max: %d)", len(exports), MAX_EXPORTS)
		t.Logf("New exports added. Please review if they're necessary.")
		t.Logf("\nSample types: %v", typeExports[:min(10, len(typeExports))])
		t.Logf("\nSample constants: %v", constExports[:min(20, len(constExports))])
		t.Logf("\nSample functions: %v", funcExports[:min(10, len(funcExports))])
	} else {
		reduction := MAX_EXPORTS - len(exports)
		t.Logf("✓ API surface is within limits (room for %d more exports)", reduction)
	}

	// Verify essential types are still exported
	essentialTypes := []string{
		"Component", "DateOrder", "DatePreference", "Period", "Timeunit", "Duration",
	}
	for _, typeName := range essentialTypes {
		if !contains(typeExports, typeName) {
			t.Errorf("Essential type %q is missing from exports", typeName)
		}
	}

	// Report deprecated types (for documentation)
	deprecatedTypes := []string{
		"Weekday", "Month", "Meridiem",
	}
	var stillPresent []string
	for _, typeName := range deprecatedTypes {
		if contains(typeExports, typeName) {
			stillPresent = append(stillPresent, typeName)
		}
	}
	if len(stillPresent) > 0 {
		t.Logf("Deprecated types still present (backward compatibility): %v", stillPresent)
	}
}

// TestDeprecatedExports documents which exports are deprecated and why.
// This serves as documentation for future API cleanup.
func TestDeprecatedExports(t *testing.T) {
	deprecations := map[string]string{
		// Weekday - use time.Weekday
		"Weekday":          "Use time.Weekday from standard library",
		"WeekdaySunday":    "Use time.Sunday from standard library",
		"WeekdayMonday":    "Use time.Monday from standard library",
		"WeekdayTuesday":   "Use time.Tuesday from standard library",
		"WeekdayWednesday": "Use time.Wednesday from standard library",
		"WeekdayThursday":  "Use time.Thursday from standard library",
		"WeekdayFriday":    "Use time.Friday from standard library",
		"WeekdaySaturday":  "Use time.Saturday from standard library",

		// Month - use time.Month
		"Month":          "Use time.Month from standard library",
		"MonthJanuary":   "Use time.January from standard library",
		"MonthFebruary":  "Use time.February from standard library",
		"MonthMarch":     "Use time.March from standard library",
		"MonthApril":     "Use time.April from standard library",
		"MonthMay":       "Use time.May from standard library",
		"MonthJune":      "Use time.June from standard library",
		"MonthJuly":      "Use time.July from standard library",
		"MonthAugust":    "Use time.August from standard library",
		"MonthSeptember": "Use time.September from standard library",
		"MonthOctober":   "Use time.October from standard library",
		"MonthNovember":  "Use time.November from standard library",
		"MonthDecember":  "Use time.December from standard library",

		// Meridiem - internal
		"Meridiem":   "Moved to internal/types",
		"MeridiemAM": "Moved to internal/types",
		"MeridiemPM": "Moved to internal/types",

		// Numeric constants - internal
		"HoursPerDay":            "Moved to internal/types",
		"MinutesPerHour":         "Moved to internal/types",
		"SecondsPerMinute":       "Moved to internal/types",
		"MillisecondsPerSecond":  "Moved to internal/types",
		"MicrosecondsPerMS":      "Moved to internal/types",
		"MicrosecondsPerSecond":  "Moved to internal/types",
		"NanosecondsPerMicro":    "Moved to internal/types",
		"NanosecondsPerMS":       "Moved to internal/types",
		"SecondsPerHour":         "Moved to internal/types",
		"MinutesPerDay":          "Moved to internal/types",
		"DaysPerWeek":            "Moved to internal/types",
		"MonthsPerYear":          "Moved to internal/types",
		"MonthsPerQuarter":       "Moved to internal/types",
		"WeeksPerMonthApprox":    "Moved to internal/types",
		"YearLookAheadThreshold": "Moved to internal/types",
	}

	t.Logf("Documented %d deprecated exports (kept for backward compatibility)", len(deprecations))

	// Group by category
	categories := map[string]int{
		"Weekday (use time.Weekday)":   0,
		"Month (use time.Month)":       0,
		"Meridiem (internal)":          0,
		"Numeric constants (internal)": 0,
	}

	for name, reason := range deprecations {
		if strings.Contains(reason, "Weekday") {
			categories["Weekday (use time.Weekday)"]++
		} else if strings.Contains(reason, "Month") {
			categories["Month (use time.Month)"]++
		} else if strings.Contains(name, "Meridiem") {
			categories["Meridiem (internal)"]++
		} else {
			categories["Numeric constants (internal)"]++
		}
	}

	for category, count := range categories {
		t.Logf("  %s: %d items", category, count)
	}
}

// TestPublicAPIStability verifies that essential types remain stable.
// Breaking changes to these types require a major version bump.
func TestPublicAPIStability(t *testing.T) {
	// Core types that should never be removed without a major version bump
	essentialTypes := []string{
		"Component",      // Core enum for accessing parsed components
		"DateOrder",      // Configuration for date format (MDY/DMY/YMD)
		"DatePreference", // Configuration for ambiguous date resolution
		"Timeunit",       // Required for Duration map keys
		"Duration",       // Public duration type
		"Result",         // Primary result interface
		"Components",     // Component access interface
		"Parser",         // Parser interface
		"ParserBuilder",  // Fluent builder API
	}

	cmd := exec.CommandContext(context.Background(), "go", "doc", "-short", ".")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run go doc: %v", err)
	}

	output := stdout.String()

	var missing []string
	for _, typeName := range essentialTypes {
		// Look for "type TypeName" in output
		if !strings.Contains(output, "type "+typeName) {
			missing = append(missing, typeName)
		}
	}

	if len(missing) > 0 {
		t.Errorf("Core API types are missing: %v", missing)
		t.Error("This is a BREAKING CHANGE requiring a major version bump!")
	} else {
		t.Logf("✓ All %d essential API types are present and stable", len(essentialTypes))
	}
}

// Helper functions

func isExported(name string) bool {
	if len(name) == 0 {
		return false
	}
	first := rune(name[0])
	return first >= 'A' && first <= 'Z'
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}


// Example output format for the API surface report
func ExampleTestAPIExports() {
	// This example shows the expected output format
	fmt.Println("Total exported symbols: 150 (max: 200)")
	fmt.Println("Exported types: 25")
	fmt.Println("Exported constants: 75")
	fmt.Println("Exported functions: 50")
	fmt.Println("✓ API surface is within limits (room for 50 more exports)")

	// Output:
	// Total exported symbols: 150 (max: 200)
	// Exported types: 25
	// Exported constants: 75
	// Exported functions: 50
	// ✓ API surface is within limits (room for 50 more exports)
}
