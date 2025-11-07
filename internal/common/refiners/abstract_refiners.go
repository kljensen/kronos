// Package refiners provides common refiner utilities for post-processing parsing results.
//
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package refiners

import (
	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// Filter is a special type of Refiner that filters results based on validity.
type Filter interface {
	kronos.Refiner
	IsValid(context *kronos.InternalParsingContext, result *kronos.InternalParsingResult) bool
}

// BaseFilter provides a default Refine implementation for filters.
type BaseFilter struct{}

// Refine filters results based on the IsValid method of the concrete filter implementation.
func (f *BaseFilter) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	filtered := make([]*kronos.InternalParsingResult, 0, len(results))
	for _, result := range results {
		// Each concrete filter must implement IsValid
		if filter, ok := interface{}(f).(Filter); ok {
			if filter.IsValid(context, result) {
				filtered = append(filtered, result)
			}
		}
	}
	return filtered
}

// MergingRefiner is a special type of Refiner that merges consecutive results.
type MergingRefiner interface {
	kronos.Refiner
	ShouldMergeResults(textBetween string, current, next *kronos.InternalParsingResult, context *kronos.InternalParsingContext) bool
	MergeResults(textBetween string, current, next *kronos.InternalParsingResult, context *kronos.InternalParsingContext) *kronos.InternalParsingResult
}

// BaseMergingRefiner provides a default Refine implementation for merging refiners.
type BaseMergingRefiner struct{}

// Refine merges consecutive results based on the ShouldMergeResults and MergeResults methods.
func (m *BaseMergingRefiner) Refine(context *kronos.InternalParsingContext, results []*kronos.InternalParsingResult) []*kronos.InternalParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.InternalParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]
		start := current.Index() + len(current.Text())
		end := next.Index()
		textBetween, okRange := helpers.SafeSlice(context.Text(), start, end)
		if !okRange {
			merged = append(merged, current)
			current = next
			continue
		}

		// Get the concrete implementation
		merger, ok := interface{}(m).(MergingRefiner)
		if !ok || !merger.ShouldMergeResults(textBetween, current, next, context) {
			merged = append(merged, current)
			current = next
		} else {
			current = merger.MergeResults(textBetween, current, next, context)
		}
	}

	if current != nil {
		merged = append(merged, current)
	}

	return merged
}
