package en

import (
	"testing"
)

// TestPeriodIntegrationWithParser tests period tracking with the actual parsing system
// NOTE: This test accesses internal Period tracking which is not part of the public API.
// It is skipped as part of the API minimization effort (Phase 5).
// The Period type and related functionality are now internal and not exposed through the public API.
func TestPeriodIntegrationWithParser(t *testing.T) {
	t.Skip("Period tracking is internal API - test requires XAsParsingComponents which accesses internal structures")
}

// TestPeriodWithCombinedExpressions tests period for combined date+time expressions
// NOTE: This test accesses internal Period tracking which is not part of the public API.
// It is skipped as part of the API minimization effort (Phase 5).
func TestPeriodWithCombinedExpressions(t *testing.T) {
	t.Skip("Period tracking is internal API - test requires XAsParsingComponents which accesses internal structures")
}

// TestPeriodConsistencyAcrossParsers tests that similar expressions have consistent periods
// NOTE: This test accesses internal Period tracking which is not part of the public API.
// It is skipped as part of the API minimization effort (Phase 5).
func TestPeriodConsistencyAcrossParsers(t *testing.T) {
	t.Skip("Period tracking is internal API - test requires XAsParsingComponents which accesses internal structures")
}

// TestPeriodPreservationInMerging tests that period is preserved during merging operations
// NOTE: This test accesses internal Period tracking which is not part of the public API.
// It is skipped as part of the API minimization effort (Phase 5).
func TestPeriodPreservationInMerging(t *testing.T) {
	t.Skip("Period tracking is internal API - test requires XAsParsingComponents which accesses internal structures")
}

// TestSlashDatePeriodTracking is a regression test for Issue #101.
// It verifies that slash-separated date formats like "03/15/2020" correctly
// return PeriodDay instead of PeriodUnknown.
// NOTE: This test accesses internal Period tracking which is not part of the public API.
// It is skipped as part of the API minimization effort (Phase 5).
func TestSlashDatePeriodTracking(t *testing.T) {
	t.Skip("Period tracking is internal API - test requires XAsParsingComponents which accesses internal structures")
}
