// Package refiners provides English-specific refiners for post-processing parsing results.
//
// This package contains refiners that:
// - Merge date ranges (merge_daterange.go)
// - Merge date and time components (merge_datetime.go)
// - Extract year suffixes (year_suffix.go)
// - Merge relative dates with absolute dates (merge_relative_after.go, merge_relative_followby.go)
// - Filter unlikely formats (unlikely_format.go)
//
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners
