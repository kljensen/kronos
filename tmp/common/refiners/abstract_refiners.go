package refiners

import . "github.com/markusmobius/go-chrono"

// Filter is a special type of Refiner that filters results based on validity.
type Filter interface {
	Refiner
	IsValid(context *ParsingContext, result *ParsingResult) bool
}

// BaseFilter provides a default Refine implementation for filters.
type BaseFilter struct{}

func (f *BaseFilter) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	filtered := make([]*ParsingResult, 0, len(results))
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
	Refiner
	ShouldMergeResults(textBetween string, current, next *ParsingResult, context *ParsingContext) bool
	MergeResults(textBetween string, current, next *ParsingResult, context *ParsingContext) *ParsingResult
}

// BaseMergingRefiner provides a default Refine implementation for merging refiners.
type BaseMergingRefiner struct{}

func (m *BaseMergingRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*ParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]
		textBetween := context.Text[current.Index+len(current.Text) : next.Index]

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
