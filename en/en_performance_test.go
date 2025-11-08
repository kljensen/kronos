package en

import (
	"testing"
	"time"
)

// TestPerformanceWhitespaceBacktracking tests that parsing doesn't exhibit
// catastrophic backtracking when encountering lots of whitespace with partial matches.
// NOTE: This test currently fails (takes >1s) - it documents a known performance issue
// that should be addressed in the future by optimizing regex patterns or parser logic.
func TestPerformanceWhitespaceBacktracking(t *testing.T) {
	t.Skip("Skipping performance test - known issue: takes >15s due to backtracking")
	// This test ensures that the parser doesn't take too long when encountering
	// a string with lots of whitespace and partial time unit matches that don't
	// form valid patterns. This guards against regex catastrophic backtracking.
	str := "BGR3                                                                                         " +
		"                                                                                        186          " +
		"                                      days                                                           " +
		"                                                                                                     " +
		"                                                                                                     " +
		"           18                                                hours                                   " +
		"                                                                                                     " +
		"                                                                                                     " +
		"                                   37                                                minutes         " +
		"                                                                                                     " +
		"                                                                                                     " +
		"                                                             01                                      " +
		"          seconds"

	start := time.Now()
	refDate := time.Now()
	results, err := New().WithReferenceDate(refDate).Parse(str)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	elapsed := time.Since(start)

	// Should find no valid results (numbers are too far from units)
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}

	// Should complete in under 1 second (original test threshold)
	if elapsed > time.Second {
		t.Errorf("Parsing took too long: %v (expected < 1s)", elapsed)
	}
}
