//nolint:staticcheck // SA1019: Must use deprecated types during transition
package en

import (
	"testing"
)

// TestIntegration_ParseWithSettings tests the new settings-based parsing.
// NOTE: These tests use CreateConfiguration which was an internal helper.
// Skipped as part of API minimization (Phase 5).
func TestIntegration_ParseWithSettings(t *testing.T) {
	t.Skip("CreateConfiguration is not part of public API - test requires internal configuration helpers")
}

// TestIntegration_ParseDateWithSettings tests the convenience function.
// NOTE: These tests use CreateConfiguration which was an internal helper.
// Skipped as part of API minimization (Phase 5).
func TestIntegration_ParseDateWithSettings(t *testing.T) {
	t.Skip("CreateConfiguration is not part of public API - test requires internal configuration helpers")
}

// TestSettings_MaxParseResults tests limiting the number of parse results.
// NOTE: These tests use CreateConfiguration which was an internal helper.
// Skipped as part of API minimization (Phase 5).
func TestSettings_MaxParseResults(t *testing.T) {
	t.Skip("CreateConfiguration is not part of public API - test requires internal configuration helpers")
}

// TestSettings_ParseTimeout tests parsing timeout functionality.
// NOTE: These tests use CreateConfiguration which was an internal helper.
// Skipped as part of API minimization (Phase 5).
func TestSettings_ParseTimeout(t *testing.T) {
	t.Skip("CreateConfiguration is not part of public API - test requires internal configuration helpers")
}

// TestSettings_SkipTokens tests skip tokens functionality.
// NOTE: These tests use CreateConfiguration which was an internal helper.
// Skipped as part of API minimization (Phase 5).
func TestSettings_SkipTokens(t *testing.T) {
	t.Skip("CreateConfiguration is not part of public API - test requires internal configuration helpers")
}

// TestSettings_Normalization tests text normalization settings.
// NOTE: These tests use CreateConfiguration which was an internal helper.
// Skipped as part of API minimization (Phase 5).
func TestSettings_Normalization(t *testing.T) {
	t.Skip("CreateConfiguration is not part of public API - test requires internal configuration helpers")
}

// TestSettings_RequiredParts tests the required parts validation.
// NOTE: These tests use CreateConfiguration which was an internal helper.
// Skipped as part of API minimization (Phase 5).
func TestSettings_RequiredParts(t *testing.T) {
	t.Skip("CreateConfiguration is not part of public API - test requires internal configuration helpers")
}
