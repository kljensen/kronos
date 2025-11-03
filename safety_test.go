package kronos

import (
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Apostrophe normalization tests
		{
			name:     "right single quotation mark (U+2019)",
			input:    "Aujourd'hui",
			expected: "Aujourd'hui",
		},
		{
			name:     "modifier letter apostrophe (U+02BC)",
			input:    "Aujourd\u02bchui",
			expected: "Aujourd'hui",
		},
		{
			name:     "modifier letter turned comma (U+02BB)",
			input:    "Aujourd\u02bbhui",
			expected: "Aujourd'hui",
		},
		{
			name:     "Armenian apostrophe (U+055A)",
			input:    "Aujourd\u055ahui",
			expected: "Aujourd'hui",
		},
		{
			name:     "Latin small letter saltillo (U+A78C)",
			input:    "Aujourd\ua78chui",
			expected: "Aujourd'hui",
		},
		{
			name:     "prime (U+2032)",
			input:    "Aujourd\u2032hui",
			expected: "Aujourd'hui",
		},
		{
			name:     "reversed prime (U+2035)",
			input:    "Aujourd\u2035hui",
			expected: "Aujourd'hui",
		},
		{
			name:     "modifier letter prime (U+02B9)",
			input:    "Aujourd\u02b9hui",
			expected: "Aujourd'hui",
		},
		{
			name:     "fullwidth apostrophe (U+FF07)",
			input:    "Aujourd\uff07hui",
			expected: "Aujourd'hui",
		},

		// Non-breaking space tests
		{
			name:     "non-breaking space (U+00A0)",
			input:    "March\u00a015",
			expected: "March 15",
		},
		{
			name:     "multiple non-breaking spaces",
			input:    "March\u00a0\u00a015",
			expected: "March 15",
		},
		{
			name:     "mixed regular and non-breaking spaces",
			input:    "March \u00a015",
			expected: "March 15",
		},

		// Ideographic space tests (CJK)
		{
			name:     "ideographic space (U+3000)",
			input:    "March\u300015",
			expected: "March 15",
		},

		// Zero-width character tests
		{
			name:     "zero width space (U+200B)",
			input:    "March\u200b15",
			expected: "March15",
		},
		{
			name:     "zero width non-joiner (U+200C)",
			input:    "March\u200c15",
			expected: "March15",
		},
		{
			name:     "zero width joiner (U+200D)",
			input:    "March\u200d15",
			expected: "March15",
		},
		{
			name:     "left-to-right mark (U+200E)",
			input:    "March\u200e15",
			expected: "March 15",
		},
		{
			name:     "right-to-left mark (U+200F)",
			input:    "March\u200f15",
			expected: "March 15",
		},
		{
			name:     "zero width no-break space/BOM (U+FEFF)",
			input:    "March\ufeff15",
			expected: "March15",
		},

		// Other special characters
		{
			name:     "right-pointing double angle quotation mark (U+00BB)",
			input:    "March\u00bb15",
			expected: "March15",
		},
		{
			name:     "middle dot (U+00B7)",
			input:    "March\u00b715",
			expected: "March15",
		},
		{
			name:     "Arabic fatha (U+064E)",
			input:    "March\u064e15",
			expected: "March15",
		},
		{
			name:     "Arabic damma (U+064F)",
			input:    "March\u064f15",
			expected: "March15",
		},

		// Multiple whitespace tests
		{
			name:     "multiple spaces",
			input:    "March  15",
			expected: "March 15",
		},
		{
			name:     "tabs and spaces",
			input:    "March\t\t15",
			expected: "March 15",
		},
		{
			name:     "mixed whitespace (tabs, newlines, spaces)",
			input:    "March \t\n 15",
			expected: "March 15",
		},
		{
			name:     "newlines",
			input:    "March\n15",
			expected: "March 15",
		},
		{
			name:     "carriage returns",
			input:    "March\r15",
			expected: "March 15",
		},
		{
			name:     "carriage return + newline",
			input:    "March\r\n15",
			expected: "March 15",
		},

		// Leading/trailing whitespace tests
		{
			name:     "leading space",
			input:    " March 15",
			expected: "March 15",
		},
		{
			name:     "trailing space",
			input:    "March 15 ",
			expected: "March 15",
		},
		{
			name:     "leading and trailing spaces",
			input:    " March 15 ",
			expected: "March 15",
		},
		{
			name:     "leading tabs",
			input:    "\t\tMarch 15",
			expected: "March 15",
		},
		{
			name:     "trailing tabs",
			input:    "March 15\t\t",
			expected: "March 15",
		},
		{
			name:     "leading newlines",
			input:    "\n\nMarch 15",
			expected: "March 15",
		},
		{
			name:     "trailing newlines",
			input:    "March 15\n\n",
			expected: "March 15",
		},

		// Combined issues tests
		{
			name:     "apostrophe + spaces",
			input:    " Aujourd'hui  ",
			expected: "Aujourd'hui",
		},
		{
			name:     "multiple apostrophe types in one string",
			input:    "It's don't can't",
			expected: "It's don't can't",
		},
		{
			name:     "apostrophe + non-breaking space + zero-width",
			input:    "Aujourd'hui\u00a0\u200bat\u00a010\u00a0PM",
			expected: "Aujourd'hui at 10 PM",
		},
		{
			name:     "complex real-world example",
			input:    " \tMarch\u200b\u00a015,\u00a02024\u200eat\u00a010:30\u00a0PM\t\n ",
			expected: "March 15, 2024 at 10:30 PM",
		},

		// Edge cases
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   \t\n  ",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "\u200b\u200c\u200d",
			expected: "",
		},
		{
			name:     "already clean string",
			input:    "March 15, 2024",
			expected: "March 15, 2024",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "single apostrophe variant",
			input:    "'",
			expected: "'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeApostrophes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no apostrophes",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "ASCII apostrophe",
			input:    "it's",
			expected: "it's",
		},
		{
			name:     "right single quotation mark",
			input:    "it's",
			expected: "it's",
		},
		{
			name:     "multiple different apostrophes",
			input:    "it's don't can't",
			expected: "it's don't can't",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeApostrophes(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeApostrophes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveZeroWidthChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "zero width space",
			input:    "hello\u200bworld",
			expected: "helloworld",
		},
		{
			name:     "multiple zero width characters",
			input:    "hello\u200b\u200c\u200dworld",
			expected: "helloworld",
		},
		{
			name:     "BOM",
			input:    "hello\ufeffworld",
			expected: "helloworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeZeroWidthChars(tt.input)
			if result != tt.expected {
				t.Errorf("removeZeroWidthChars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkSanitizeInput(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"clean", "March 15, 2024"},
		{"apostrophes", "Aujourd'hui"},
		{"spaces", "March  \u00a0  15"},
		{"complex", " \tMarch\u200b\u00a015,\u00a02024\u200eat\u00a010:30\u00a0PM\t\n "},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				SanitizeInput(tc.input)
			}
		})
	}
}
