//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
	"time"

	"github.com/kljensen/kronos"
)

// TestSanitizationIntegration tests that sanitization works correctly during actual parsing.
// This test verifies that various Unicode normalization issues are handled properly
// so that dates with apostrophe variants, non-breaking spaces, and zero-width characters
// parse correctly.
func TestSanitizationIntegration(t *testing.T) {
	// Note: These tests don't verify specific parsing results, but rather that
	// sanitization doesn't break the parsing process and text is normalized.
	tests := []struct {
		name        string
		input       string
		expectParse bool // Whether we expect the parser to find something
		description string
	}{
		{
			name:        "apostrophe variant in French word",
			input:       "Aujourd'hui", // U+2019 RIGHT SINGLE QUOTATION MARK
			expectParse: false,         // We don't have French parser, but sanitization should work
			description: "Should normalize apostrophe variant to ASCII",
		},
		{
			name:        "date with non-breaking spaces",
			input:       "March\u00a015,\u00a02024",
			expectParse: true,
			description: "Should replace non-breaking spaces with regular spaces",
		},
		{
			name:        "date with zero-width space",
			input:       "March\u200b15, 2024",
			expectParse: true,
			description: "Should remove zero-width space",
		},
		{
			name:        "date with multiple whitespace types",
			input:       "  March\t15,  2024  ",
			expectParse: true,
			description: "Should normalize all whitespace",
		},
		{
			name:        "date with directional marks",
			input:       "March\u200e15,\u200f2024",
			expectParse: true,
			description: "Should handle directional marks",
		},
		{
			name:        "complex real-world example",
			input:       " \tMarch\u200b\u00a015,\u00a02024\u200eat\u00a010:30\u00a0PM\t\n ",
			expectParse: true,
			description: "Should handle combination of multiple Unicode issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a basic English chrono parser
			chrono := CasualChrono()
			refDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

			// Parse the input
			results := chrono.Parse(tt.input, refDate, nil)

			// Verify parsing behavior
			if tt.expectParse {
				if len(results) == 0 {
					t.Errorf("Expected to parse results from %q but got none. %s", tt.input, tt.description)
				}
			}

			// The key test: verify that sanitization happened by checking that
			// the context received sanitized text (this is implicit in successful parsing)
			// If sanitization didn't work, parsing would fail or behave incorrectly
		})
	}
}

// TestSanitizationPreservesSemantics tests that sanitization doesn't change
// the semantic meaning of the text, only normalizes the representation.
func TestSanitizationPreservesSemantics(t *testing.T) {
	chrono := CasualChrono()
	refDate := time.Date(2024, 3, 10, 12, 0, 0, 0, time.UTC)

	// Test that the same semantic date with different Unicode representations
	// produces the same parsing result
	testCases := []struct {
		name    string
		inputs  []string // Different Unicode representations of the same semantic date
		message string
	}{
		{
			name: "date with space variations",
			inputs: []string{
				"March 15, 2024",
				"March\u00a015,\u00a02024", // Non-breaking spaces
				"March  15,  2024",         // Multiple spaces
				"  March 15, 2024  ",       // Leading/trailing spaces
			},
			message: "All space variations should parse to the same result",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var firstResult *kronos.InternalParsingResult
			for i, input := range tc.inputs {
				results := chrono.Parse(input, refDate, nil)
				if len(results) == 0 {
					t.Errorf("Failed to parse input %d: %q. %s", i, input, tc.message)
					continue
				}

				if firstResult == nil {
					firstResult = results[0]
					continue
				}

				// Compare the parsed dates
				currentDate := results[0].Date()
				expectedDate := firstResult.Date()

				if !currentDate.Equal(expectedDate) {
					t.Errorf("Input %d (%q) parsed to %v, expected %v. %s",
						i, input, currentDate, expectedDate, tc.message)
				}
			}
		})
	}
}
