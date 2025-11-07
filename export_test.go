package kronos

// Export internal functions for testing.
// This file is only compiled during tests and allows test code to access
// internal functionality without exposing it in the public API.

// Test helpers for weekday calculations
var (
	GetNthWeekdayOfMonth  = InternalGetNthWeekdayOfMonth
	GetLastWeekdayOfMonth = InternalGetLastWeekdayOfMonth
)
