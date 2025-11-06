package kronos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TeststripApproximationWords(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedClean string
		isApproximate bool
	}{
		{
			name:          "about prefix",
			input:         "about 2 hours ago",
			expectedClean: "2 hours ago",
			isApproximate: true,
		},
		{
			name:          "around prefix",
			input:         "around 3 days",
			expectedClean: "3 days",
			isApproximate: true,
		},
		{
			name:          "roughly prefix",
			input:         "roughly 5 minutes",
			expectedClean: "5 minutes",
			isApproximate: true,
		},
		{
			name:          "approximately prefix",
			input:         "approximately 1 week",
			expectedClean: "1 week",
			isApproximate: true,
		},
		{
			name:          "approx prefix",
			input:         "approx 2 months",
			expectedClean: "2 months",
			isApproximate: true,
		},
		{
			name:          "circa prefix",
			input:         "circa 1 year",
			expectedClean: "1 year",
			isApproximate: true,
		},
		{
			name:          "tilde prefix",
			input:         "~2 hours",
			expectedClean: "2 hours",
			isApproximate: true,
		},
		{
			name:          "tilde with space",
			input:         "~ 2 hours",
			expectedClean: "2 hours",
			isApproximate: true,
		},
		{
			name:          "no approximation",
			input:         "2 hours ago",
			expectedClean: "2 hours ago",
			isApproximate: false,
		},
		{
			name:          "case insensitive ABOUT",
			input:         "ABOUT 3 days",
			expectedClean: "3 days",
			isApproximate: true,
		},
		{
			name:          "case insensitive Around",
			input:         "Around 3 days",
			expectedClean: "3 days",
			isApproximate: true,
		},
		{
			name:          "multiple spaces after approximation",
			input:         "about   2 hours",
			expectedClean: "2 hours",
			isApproximate: true,
		},
		{
			name:          "word boundary - roundabout should not match",
			input:         "roundabout 2 hours",
			expectedClean: "roundabout 2 hours",
			isApproximate: false,
		},
		{
			name:          "multiple approximation words",
			input:         "about approximately 2 hours",
			expectedClean: "2 hours",
			isApproximate: true,
		},
		{
			name:          "approximation with noon",
			input:         "around noon",
			expectedClean: "noon",
			isApproximate: true,
		},
		{
			name:          "approximation with midnight",
			input:         "about midnight",
			expectedClean: "midnight",
			isApproximate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned, isApprox := stripApproximationWords(tt.input)
			assert.Equal(t, tt.expectedClean, cleaned, "Cleaned text mismatch")
			assert.Equal(t, tt.isApproximate, isApprox, "isApproximate flag mismatch")
		})
	}
}
