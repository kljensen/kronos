// Package helpers provides internal utility functions for kronos date parsing.
// This package contains helper functions that were previously exported with X-prefixes
// from the main kronos package. These functions are implementation details and should
// not be used by external code.
//
// The helpers package includes:
//   - Casual date reference functions (Today, Tomorrow, Yesterday, Now, etc.)
//   - Date math utilities (year finding, weekday calculations)
//   - Component manipulation (AssignSimilarDate, MergeDateTimeComponent, etc.)
//   - Duration helpers (AddDuration, ReverseDuration)
//   - String utilities (SafeSlice, StripApproximationWords)
//   - Timezone helpers (ToTimezoneOffset, FromInput)
//   - Result helpers (MergeDateTimeResult, CreateRelativeFromReference)
package helpers
