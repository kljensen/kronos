package en

// Test helper functions for en package tests

// intPtr returns a pointer to the given int value.
// This is a test utility for creating pointer values in test data.
func intPtr(v int) *int {
	return &v
}
